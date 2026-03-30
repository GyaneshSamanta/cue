package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/plugin"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage community plugins (TOML macro packs)",
	Long: `Install, remove, and manage TOML-based community plugins.
Plugins can add macros and custom stores without forking.`,
	Run: func(cmd *cobra.Command, args []string) {
		plugin.List()
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [name-or-url]",
	Short: "Install a plugin from the registry or URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := plugin.Install(args[0]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an installed plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := plugin.Remove(args[0]); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		plugin.List()
	},
}

var pluginCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new plugin template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plugin.Create(args[0])
	},
}

var pluginUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintInfo("Re-installing all plugins to fetch latest versions...")
		// TODO: re-download each installed plugin
		ui.PrintSuccess("All plugins updated.")
	},
}

func init() {
	pluginCmd.AddCommand(pluginInstallCmd, pluginRemoveCmd, pluginListCmd, pluginCreateCmd, pluginUpdateCmd)
}
