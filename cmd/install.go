package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/queue"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Install a package with queue management",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		qm := queue.NewManager(osAdapter)
		err := qm.Enqueue("install", args, func() error {
			return osAdapter.InstallPackage(args[0], args[1:])
		})
		if err != nil {
			ui.PrintError(err.Error())
		} else {
			ui.PrintSuccess(args[0] + " installed successfully.")
		}
	},
}
