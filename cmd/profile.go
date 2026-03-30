package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/profile"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage environment profiles (work, personal, etc.)",
	Long: `Named profiles for developers who maintain separate identities.
Each profile stores git config, env variables, SSH keys, and shell settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile.ListProfiles()
	},
}

var profileCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.Create(args[0]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var profileSwitchCmd = &cobra.Command{
	Use:   "switch [name]",
	Short: "Switch to a profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.Switch(args[0]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profile.ListProfiles()
	},
}

var profileCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show active profile",
	Run: func(cmd *cobra.Command, args []string) {
		profile.Current()
	},
}

var profileDiffCmd = &cobra.Command{
	Use:   "diff [profile-a] [profile-b]",
	Short: "Compare two profiles",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := profile.DiffProfiles(args[0], args[1]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

func init() {
	profileCmd.AddCommand(profileCreateCmd, profileSwitchCmd, profileListCmd, profileCurrentCmd, profileDiffCmd)
}
