package job

import (
	"encoding/json"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Job represents a managed child process.
type Job struct {
	ID          string     `json:"id"`
	PID         int        `json:"pid"`
	Command     string     `json:"command"`
	Status      string     `json:"status"` // running | paused | completed | failed | orphaned
	PauseReason string     `json:"pause_reason,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	PausedAt    *time.Time `json:"paused_at,omitempty"`
	ResumedAt   *time.Time `json:"resumed_at,omitempty"`
	ProjectTag  string     `json:"project_tag,omitempty"`
	Progress    string     `json:"progress_hint,omitempty"`
}

// JobsFile is the on-disk state.
type JobsFile struct {
	Version int   `json:"version"`
	Jobs    []Job `json:"jobs"`
}

// Controller manages child processes with network-aware pause/resume.
type Controller struct {
	Adapter    adapter.OSAdapter
	NetMonitor *NetworkMonitor
	filePath   string
}

// NewController creates a job controller.
func NewController(a adapter.OSAdapter) *Controller {
	return &Controller{
		Adapter:    a,
		NetMonitor: NewNetworkMonitor(),
		filePath:   filepath.Join(a.ConfigDir(), "jobs.json"),
	}
}

// ManagedRun spawns cmd and wraps it with network-aware auto-pause/resume.
func (jc *Controller) ManagedRun(cmd *exec.Cmd, jobID string) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	pid := cmd.Process.Pid
	job := &Job{
		ID:        jobID,
		PID:       pid,
		Command:   strings.Join(cmd.Args, " "),
		Status:    "running",
		CreatedAt: time.Now(),
	}
	jc.persist(job)

	netEvents := jc.NetMonitor.Watch()

	// Signal listener for manual pause
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, SignalTSTP())

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case err := <-done:
			job.Status = "completed"
			if err != nil {
				job.Status = "failed"
			}
			jc.persist(job)
			jc.NetMonitor.Stop()
			return err

		case event := <-netEvents:
			switch event {
			case NetworkLost:
				ui.PrintWarning("[PAUSED] Network lost. Installation paused. Waiting for connectivity...")
				jc.Adapter.SuspendProcess(pid)
				job.Status = "paused"
				job.PauseReason = "network_loss"
				now := time.Now()
				job.PausedAt = &now
				jc.persist(job)

			case NetworkRestored:
				ui.PrintSuccess("[RESUMED] Network restored. Continuing...")
				jc.Adapter.ResumeProcess(pid)
				job.Status = "running"
				now := time.Now()
				job.ResumedAt = &now
				jc.persist(job)
			}

		case <-sigCh:
			ui.PrintInfo("[PAUSED] Manual pause. Run 'cue resume' to continue.")
			jc.Adapter.SuspendProcess(pid)
			job.Status = "paused"
			job.PauseReason = "manual"
			now := time.Now()
			job.PausedAt = &now
			jc.persist(job)
		}
	}
}

// ListPaused returns all paused jobs.
func (jc *Controller) ListPaused() []Job {
	jf := jc.load()
	var paused []Job
	for _, j := range jf.Jobs {
		if j.Status == "paused" {
			paused = append(paused, j)
		}
	}
	return paused
}

// ResumeJob resumes a specific paused job by ID.
func (jc *Controller) ResumeJob(jobID string) error {
	jf := jc.load()
	for i, j := range jf.Jobs {
		if j.ID == jobID && j.Status == "paused" {
			if err := jc.Adapter.ResumeProcess(j.PID); err != nil {
				jf.Jobs[i].Status = "orphaned"
				jc.saveFile(jf)
				return err
			}
			jf.Jobs[i].Status = "running"
			now := time.Now()
			jf.Jobs[i].ResumedAt = &now
			jc.saveFile(jf)
			ui.PrintSuccess("Job resumed: " + j.Command)
			return nil
		}
	}
	return nil
}

func (jc *Controller) persist(job *Job) {
	jf := jc.load()
	found := false
	for i, j := range jf.Jobs {
		if j.ID == job.ID {
			jf.Jobs[i] = *job
			found = true
			break
		}
	}
	if !found {
		jf.Jobs = append(jf.Jobs, *job)
	}
	jc.saveFile(jf)
}

func (jc *Controller) load() JobsFile {
	jf := JobsFile{Version: 1}
	data, err := os.ReadFile(jc.filePath)
	if err == nil {
		json.Unmarshal(data, &jf)
	}
	return jf
}

func (jc *Controller) saveFile(jf JobsFile) {
	out, _ := json.MarshalIndent(jf, "", "  ")
	os.MkdirAll(filepath.Dir(jc.filePath), 0755)
	os.WriteFile(jc.filePath, out, 0644)
}
