package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/team"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage shared team configuration",
	Long: `Sync macros, stores, and config defaults across your team
using a shared GitHub repository.`,
}

var teamInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize team config directory",
	Run: func(cmd *cobra.Command, args []string) {
		team.Init()
	},
}

var teamConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a team repo",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			ui.PrintError("Specify repo URL: --repo <url>")
			return
		}
		if err := team.Connect(repo); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var teamSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync team config (pull or push)",
	Run: func(cmd *cobra.Command, args []string) {
		push, _ := cmd.Flags().GetBool("push")
		if err := team.Sync(push); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var teamDisconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect from team repo",
	Run: func(cmd *cobra.Command, args []string) {
		team.Disconnect()
	},
}

func init() {
	teamConnectCmd.Flags().String("repo", "", "GitHub repo URL for team config")
	teamSyncCmd.Flags().Bool("push", false, "Push local config to team repo")
	teamCmd.AddCommand(teamInitCmd, teamConnectCmd, teamSyncCmd, teamDisconnectCmd)
}
