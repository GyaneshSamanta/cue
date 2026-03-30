package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/job"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume paused jobs",
	Run: func(cmd *cobra.Command, args []string) {
		jc := job.NewController(osAdapter)
		paused := jc.ListPaused()

		if len(paused) == 0 {
			ui.PrintInfo("No paused jobs.")
			return
		}

		ui.PrintHeader("Paused Jobs")
		for i, j := range paused {
			reason := j.PauseReason
			if reason == "" {
				reason = "unknown"
			}
			fmt.Printf("  [%d] %s (PID: %d, reason: %s)\n", i+1, j.Command, j.PID, reason)
		}

		if len(paused) == 1 {
			if err := jc.ResumeJob(paused[0].ID); err != nil {
				ui.PrintError(err.Error())
			}
			return
		}

		idx := ui.SelectOne("Which job to resume?", jobNames(paused))
		if idx >= 0 && idx < len(paused) {
			if err := jc.ResumeJob(paused[idx].ID); err != nil {
				ui.PrintError(err.Error())
			}
		}
	},
}

func jobNames(jobs []job.Job) []string {
	names := make([]string, len(jobs))
	for i, j := range jobs {
		names[i] = j.Command
	}
	return names
}
