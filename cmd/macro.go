package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/macro"
	"github.com/GyaneshSamanta/cue/internal/tui"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

func pickMacro() string {
	items := make([]tui.PickerItem, 0, len(macro.Registry))
	for _, m := range macro.Registry {
		tag := "built-in"
		if !m.BuiltIn {
			tag = "custom"
		}
		items = append(items, tui.PickerItem{
			Name:        m.Name,
			Description: m.Description,
			Tag:         tag,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	idx, err := tui.RunPicker("Select Macro", items)
	if err != nil || idx == -1 {
		return ""
	}
	return items[idx].Name
}

var macroCmd = &cobra.Command{
	Use:   "macro",
	Short: "Manage macros (add, list, remove)",
	Run: func(cmd *cobra.Command, args []string) {
		name := pickMacro()
		if name == "" {
			return
		}
		if err := macro.Execute(name, map[string]string{}); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var macroListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available macros",
	Run: func(cmd *cobra.Command, args []string) {
		macro.ListAll()
	},
}

var macroAddCmd = &cobra.Command{
	Use:   "add [name] [command] [explanation]",
	Short: "Add a custom macro",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		command := args[1]
		explanation := args[2]

		tm := macro.TomlMacro{
			Name:        name,
			Command:     command,
			Explanation: explanation,
			Category:    "custom",
			Description: "User-defined macro",
		}

		path := filepath.Join(config.ConfigDir(), "macros.toml")
		if err := macro.SaveUserMacro(path, tm); err != nil {
			ui.PrintError(fmt.Sprintf("Failed to save macro: %v", err))
			return
		}

		// Also register it in the current session
		macro.Register(&macro.Macro{
			Name:        name,
			Category:    "custom",
			Description: "User-defined macro",
			Commands:    []macro.Step{{OS: "all", Command: command}},
			Explanation: explanation,
			BuiltIn:     false,
		})

		ui.PrintSuccess(fmt.Sprintf("Macro '%s' added and saved successfully.", name))
	},
}

var explainCmd = &cobra.Command{
	Use:   "explain [macro-name]",
	Short: "Explain a macro without executing it",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			macro.ListAll()
			return
		}

		name := ""
		if len(args) == 1 {
			name = args[0]
		} else {
			name = pickMacro()
			if name == "" {
				return
			}
		}

		if err := macro.Explain(name); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

func init() {
	explainCmd.Flags().Bool("list", false, "List all macros with descriptions")
	macroCmd.AddCommand(macroListCmd, macroAddCmd)

	// Register each built-in macro as a root-level command for direct invocation
	for name, m := range macro.Registry {
		macroName := name
		macroMeta := m
		macroExecCmd := &cobra.Command{
			Use:   macroName,
			Short: macroMeta.Description,
			RunE: func(cmd *cobra.Command, args []string) error {
				flags := map[string]string{}
				for _, f := range macroMeta.Flags {
					val, _ := cmd.Flags().GetString(f.Name)
					if val != "" {
						flags[f.Name] = val
					}
				}
				// Check --hard flag specifically
				hard, _ := cmd.Flags().GetBool("hard")
				if hard {
					flags["hard"] = "true"
				}
				return macro.Execute(macroName, flags)
			},
		}
		for _, f := range macroMeta.Flags {
			if f.Name == "hard" {
				macroExecCmd.Flags().Bool(f.Name, false, f.Description)
			} else {
				macroExecCmd.Flags().String(f.Name, f.Default, f.Description)
			}
		}
		rootCmd.AddCommand(macroExecCmd)
	}
}
