package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/store"
	"github.com/GyaneshSamanta/cue/internal/tui"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

func pickStack() string {
	stacks := store.ListStacks()
	items := make([]tui.PickerItem, len(stacks))
	for i, s := range stacks {
		items[i] = tui.PickerItem{
			Name:        s.Name(),
			Description: s.Description(),
			SizeMB:      s.EstimatedSizeMB(),
		}
	}
	idx, err := tui.RunPicker("Select Environment Stack", items)
	if err != nil || idx == -1 {
		return ""
	}
	return stacks[idx].Name()
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage environment stores (install, preview, verify, remove)",
	Run: func(cmd *cobra.Command, args []string) {
		actions := []tui.PickerItem{
			{Name: "Install", Description: "Set up a new environment stack"},
			{Name: "Preview", Description: "Preview what a store will install"},
			{Name: "Verify", Description: "Check if installed components are healthy"},
			{Name: "Remove", Description: "Uninstall a stack and its components"},
		}
		idx, err := tui.RunPicker("Store Actions", actions)
		if err != nil || idx == -1 {
			// Fallback: Just list them
			ui.PrintHeader("Available Environment Stores")
			headers := []string{"Store", "Description", "Size"}
			var rows [][]string
			for _, s := range store.ListStacks() {
				rows = append(rows, []string{
					s.Name(), s.Description(),
					formatSize(s.EstimatedSizeMB()),
				})
			}
			ui.PrintTable(headers, rows)
			ui.PrintInfo("Use 'cue store install <name>' to set up a stack.")
			return
		}

		action := actions[idx].Name
		stack := pickStack()
		if stack == "" {
			return
		}

		verify, _ := cmd.Flags().GetBool("verify")
		force, _ := cmd.Flags().GetBool("force")

		var doErr error
		switch action {
		case "Install":
			doErr = store.Install(stack, osAdapter, store.InstallOpts{Verify: verify})
		case "Preview":
			doErr = store.Preview(stack, osAdapter)
		case "Verify":
			doErr = store.Verify(stack, osAdapter)
		case "Remove":
			doErr = store.Remove(stack, osAdapter, force)
		}

		if doErr != nil {
			ui.PrintError(doErr.Error())
		}
	},
}

var storeInstallCmd = &cobra.Command{
	Use:   "install [stack]",
	Short: "Install an environment store",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stack := ""
		if len(args) == 1 {
			stack = args[0]
		} else {
			stack = pickStack()
			if stack == "" {
				return
			}
		}
		verify, _ := cmd.Flags().GetBool("verify")
		err := store.Install(stack, osAdapter, store.InstallOpts{Verify: verify})
		if err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var storePreviewCmd = &cobra.Command{
	Use:   "preview [stack]",
	Short: "Preview what a store will install",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stack := ""
		if len(args) == 1 {
			stack = args[0]
		} else {
			stack = pickStack()
			if stack == "" {
				return
			}
		}
		if err := store.Preview(stack, osAdapter); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var storeVerifyCmd = &cobra.Command{
	Use:   "verify [stack]",
	Short: "Verify installed store components",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stack := ""
		if len(args) == 1 {
			stack = args[0]
		} else {
			stack = pickStack()
			if stack == "" {
				return
			}
		}
		if err := store.Verify(stack, osAdapter); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var storeRemoveCmd = &cobra.Command{
	Use:   "remove [stack]",
	Short: "Remove a store and its components",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stack := ""
		if len(args) == 1 {
			stack = args[0]
		} else {
			stack = pickStack()
			if stack == "" {
				return
			}
		}
		force, _ := cmd.Flags().GetBool("force")
		if err := store.Remove(stack, osAdapter, force); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

func init() {
	storeInstallCmd.Flags().Bool("verify", true, "Run verification after install")
	storeRemoveCmd.Flags().Bool("force", false, "Force remove shared components")

	// Apply to root store command too so flags exist when called via picker
	storeCmd.Flags().Bool("verify", true, "Run verification after install")
	storeCmd.Flags().Bool("force", false, "Force remove shared components")

	storeCmd.AddCommand(storeInstallCmd, storePreviewCmd, storeVerifyCmd, storeRemoveCmd)
}

func formatSize(mb int) string {
	if mb >= 1000 {
		return fmt.Sprintf("~%.1f GB", float64(mb)/1000)
	}
	return fmt.Sprintf("~%d MB", mb)
}

// Ignore unused error
var _ = strings.Contains
