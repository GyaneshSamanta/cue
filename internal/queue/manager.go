package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Entry represents one queued command.
type Entry struct {
	ID           string     `json:"id"`
	Command      string     `json:"command"`
	Args         []string   `json:"args"`
	PkgManager   string     `json:"package_manager"`
	Status       string     `json:"status"` // waiting | running | done | failed
	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ExitCode     *int       `json:"exit_code,omitempty"`
	WaitSeconds  int        `json:"wait_seconds"`
}

// QueueFile holds the on-disk JSON state.
type QueueFile struct {
	Version int     `json:"version"`
	Entries []Entry `json:"entries"`
}

// Manager orchestrates lock-aware command execution.
type Manager struct {
	Adapter  adapter.OSAdapter
	Detector *LockDetector
	Poller   *Poller
	filePath string
}

// NewManager creates a queue manager.
func NewManager(a adapter.OSAdapter) *Manager {
	cfg := config.DefaultConfig()
	if config.Current != nil {
		cfg = config.Current
	}
	det := &LockDetector{adapter: a}
	return &Manager{
		Adapter:  a,
		Detector: det,
		Poller: &Poller{
			detector:     det,
			baseInterval: time.Duration(cfg.Core.LockPollIntervalSecs) * time.Second,
			maxInterval:  15 * time.Second,
			backoffAfter: 2 * time.Minute,
			timeout:      time.Duration(cfg.Core.LockTimeoutMins) * time.Minute,
		},
		filePath: filepath.Join(a.ConfigDir(), "queue.json"),
	}
}

// Enqueue checks for locks and either executes immediately or queues.
func (m *Manager) Enqueue(command string, args []string, fn func() error) error {
	locked, desc, err := m.Detector.IsLocked()
	if err != nil {
		return fmt.Errorf("lock check error: %w", err)
	}

	if !locked {
		return fn()
	}

	entry := Entry{
		ID:         fmt.Sprintf("q-%s", time.Now().Format("20060102-150405")),
		Command:    command,
		Args:       args,
		PkgManager: m.Adapter.PackageManagerName(),
		Status:     "waiting",
		CreatedAt:  time.Now(),
	}

	ui.PrintInfo(fmt.Sprintf("[QUEUED] Package manager is busy (%s). Your command will run automatically when it's free. Press Ctrl+C to cancel.", desc))

	m.save(entry)
	return m.Poller.WaitAndExecute(fn)
}

func (m *Manager) save(e Entry) {
	qf := QueueFile{Version: 1}
	data, err := os.ReadFile(m.filePath)
	if err == nil {
		json.Unmarshal(data, &qf)
	}
	qf.Entries = append(qf.Entries, e)
	out, _ := json.MarshalIndent(qf, "", "  ")
	os.MkdirAll(filepath.Dir(m.filePath), 0755)
	os.WriteFile(m.filePath, out, 0644)
}
