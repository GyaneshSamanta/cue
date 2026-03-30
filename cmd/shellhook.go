package cmd

import (
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/shellhook"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var shellHookCmd = &cobra.Command{
	Use:   "shell-hook",
	Short: "Manage shell hooks for auto-activation and project detection",
	Long: `Install or uninstall shell hooks that provide:
  • Auto-activate virtual environments on cd
  • Auto-switch project tags from .cue config
  • Stack detection hints when missing tools
  • Prompt enrichment with active project tag`,
}

var shellHookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install shell hooks into your shell config",
	Run: func(cmd *cobra.Command, args []string) {
		shell := detectCurrentShell()
		if err := shellhook.Install(shell); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var shellHookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove shell hooks from your shell config",
	Run: func(cmd *cobra.Command, args []string) {
		shell := detectCurrentShell()
		if err := shellhook.Uninstall(shell); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

func detectCurrentShell() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	shell := os.Getenv("SHELL")
	switch {
	case len(shell) >= 3 && shell[len(shell)-3:] == "zsh":
		return "zsh"
	case len(shell) >= 4 && shell[len(shell)-4:] == "fish":
		return "fish"
	default:
		return "bash"
	}
}

func init() {
	shellHookCmd.AddCommand(shellHookInstallCmd, shellHookUninstallCmd)
}
