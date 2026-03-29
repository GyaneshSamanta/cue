package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/gyanesh-help/internal/claudecode"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

var claudeCodeCmd = &cobra.Command{
	Use:   "claude-code",
	Short: "Manage Claude Code installation and configuration",
	Long: `Day-2 management commands for your Claude Code installation.
Switch between API/local modes, manage MCP servers, check status, and update.`,
	Run: func(cmd *cobra.Command, args []string) {
		claudecode.Status()
	},
}

var claudeCodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current Claude Code configuration and status",
	Run: func(cmd *cobra.Command, args []string) {
		claudecode.Status()
	},
}

var claudeCodeSwitchCmd = &cobra.Command{
	Use:   "switch [api|local|hybrid]",
	Short: "Switch between Claude Code modes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mode := args[0]
		switch mode {
		case "api":
			ui.PrintStep("Switching to API mode...")
			claudecode.InstallAPIMode(osAdapter)
		case "local":
			ui.PrintStep("Switching to local mode...")
			claudecode.InstallLocalMode(osAdapter)
		case "hybrid":
			ui.PrintStep("Switching to hybrid mode...")
			claudecode.InstallHybridMode(osAdapter)
		default:
			ui.PrintError(fmt.Sprintf("Unknown mode: %s. Use 'api', 'local', or 'hybrid'.", mode))
		}
	},
}

var claudeCodeMCPCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage MCP (Model Context Protocol) servers",
}

var claudeCodeMCPListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured MCP servers",
	Run: func(cmd *cobra.Command, args []string) {
		claudecode.ListMCPServers()
	},
}

var claudeCodeMCPAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new MCP server interactively",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		claudecode.AddMCPServer(args[0])
	},
}

var claudeCodeMCPRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an MCP server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		claudecode.RemoveMCPServer(args[0])
	},
}

var claudeCodeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Claude Code CLI to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		claudecode.UpdateCLI()
	},
}

func init() {
	claudeCodeMCPCmd.AddCommand(claudeCodeMCPListCmd, claudeCodeMCPAddCmd, claudeCodeMCPRemoveCmd)
	claudeCodeCmd.AddCommand(claudeCodeStatusCmd, claudeCodeSwitchCmd, claudeCodeMCPCmd, claudeCodeUpdateCmd)
}
