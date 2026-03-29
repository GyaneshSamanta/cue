package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
	vermgr "github.com/GyaneshSamanta/gyanesh-help/internal/version"
)

var versionMgrCmd = &cobra.Command{
	Use:   "version",
	Short: "Unified version manager for all runtimes",
	Long: `Switch between installed versions of Python, Node.js, Rust, Java, Go, Ruby,
and Terraform — using a single consistent interface regardless of the underlying version manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		vermgr.ShowCurrent()
	},
}

var versionListCmd = &cobra.Command{
	Use:   "list [runtime]",
	Short: "List installed versions of a runtime",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vermgr.ListVersions(args[0]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var versionUseCmd = &cobra.Command{
	Use:   "use [runtime] [version]",
	Short: "Switch to a specific version",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vermgr.UseVersion(args[0], args[1]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var versionInstallCmd = &cobra.Command{
	Use:   "install [runtime] [version]",
	Short: "Install a specific version of a runtime",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vermgr.InstallVersion(args[0], args[1]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var versionRemoveCmd = &cobra.Command{
	Use:   "remove [runtime] [version]",
	Short: "Remove a specific version of a runtime",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vermgr.RemoveVersion(args[0], args[1]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var versionPinCmd = &cobra.Command{
	Use:   "pin [runtime]",
	Short: "Write a version pin file for the current directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vermgr.Pin(args[0]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var versionCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show active versions of all managed runtimes",
	Run: func(cmd *cobra.Command, args []string) {
		vermgr.ShowCurrent()
	},
}

func init() {
	versionMgrCmd.AddCommand(versionListCmd, versionUseCmd, versionInstallCmd,
		versionRemoveCmd, versionPinCmd, versionCurrentCmd)
}
