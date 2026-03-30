package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/ui"
	"github.com/GyaneshSamanta/cue/internal/workspace"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Backup and restore workspace configuration",
}

var workspaceBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup workspace to GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		token := ui.ReadInput("GitHub PAT: ")
		if token == "" {
			ui.PrintError("Token required. Use: cue workspace auth --token <PAT>")
			return
		}
		if err := workspace.Backup(osAdapter, token, "dev-workspace-backup"); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var workspaceRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore workspace from GitHub backup",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			ui.PrintError("Specify repo URL: --repo <url>")
			return
		}
		if err := workspace.Restore(repo, osAdapter); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

func init() {
	workspaceRestoreCmd.Flags().String("repo", "", "Backup repo URL")
	workspaceCmd.AddCommand(workspaceBackupCmd, workspaceRestoreCmd)
}
