package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/toolkit"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var (
	toolkitDryRun bool
)

var toolkitCmd = &cobra.Command{
	Use:   "toolkit",
	Short: "Install and manage developer tools",
	Long: `Install and manage developer tools with automatic version manager bootstrapping.

Supports Node.js, Python, Go, Rust, Java, Docker, and more.
Run 'cue toolkit list' to see available tools.`,
}

var toolkitListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tools",
	Run: func(cmd *cobra.Command, args []string) {
		tools := toolkit.List()
		sort.Slice(tools, func(i, j int) bool {
			return tools[i].Name < tools[j].Name
		})

		ui.PrintHeader("Available Tools")
		fmt.Printf("Total: %d tools\n\n", len(tools))

		for _, t := range tools {
			ver, installed := t.VerifyFunc()
			status := "not installed"
			if installed {
				status = ver
				ui.PrintSuccess(fmt.Sprintf("%-12s %-30s [%s]", t.Name, t.DisplayName, status))
			} else {
				ui.PrintWarning(fmt.Sprintf("%-12s %-30s [%s]", t.Name, t.DisplayName, status))
			}
		}
		fmt.Println()
		fmt.Println("Run 'cue toolkit install <tool>' to install a tool.")
		fmt.Println("Run 'cue toolkit info <tool>' for details.")
	},
}

var toolkitInstallCmd = &cobra.Command{
	Use:   "install <tool> [version]",
	Short: "Install a developer tool",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]
		version := ""
		if len(args) > 1 {
			version = args[1]
		}
		return toolkit.Install(toolName, osAdapter, version, toolkitDryRun)
	},
}

var toolkitUpgradeCmd = &cobra.Command{
	Use:   "upgrade <tool>",
	Short: "Upgrade a tool to the latest version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolkit.Upgrade(args[0], osAdapter)
	},
}

var toolkitRemoveCmd = &cobra.Command{
	Use:   "remove <tool>",
	Short: "Remove an installed tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolkit.Remove(args[0], osAdapter)
	},
}

var toolkitInfoCmd = &cobra.Command{
	Use:   "info <tool>",
	Short: "Show detailed information about a tool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return toolkit.Info(args[0])
	},
}

func init() {
	toolkitCmd.AddCommand(toolkitListCmd)
	toolkitCmd.AddCommand(toolkitInstallCmd)
	toolkitCmd.AddCommand(toolkitUpgradeCmd)
	toolkitCmd.AddCommand(toolkitRemoveCmd)
	toolkitCmd.AddCommand(toolkitInfoCmd)

	toolkitInstallCmd.Flags().BoolVar(&toolkitDryRun, "dry-run", false, "Show what would be installed without installing")

	rootCmd.AddCommand(toolkitCmd)
}
