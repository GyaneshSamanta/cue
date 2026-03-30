package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/store"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var diffCmd = &cobra.Command{
	Use:   "diff [stack]",
	Short: "Show gaps between your environment and a store's requirements",
	Long: `Compare your current installed tools against what a store expects.
Shows present, missing, and outdated components with actionable fix commands.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		stack, err := store.GetStack(stackName)
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		ui.PrintHeader(fmt.Sprintf("Environment Gap Analysis: %s", stack.Name()))
		fmt.Printf("  Comparing current environment vs %s store requirements...\n\n", stack.Name())

		checks := stack.VerificationChecks()
		missing := 0
		outdated := 0
		present := 0

		for _, check := range checks {
			parts := strings.Fields(check.Command)
			if len(parts) == 0 {
				continue
			}
			out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
			output := strings.TrimSpace(string(out))
			if len(output) > 40 {
				output = output[:40]
			}

			if err != nil {
				fmt.Printf("  MISSING   ✗  %-16s —         not installed\n", check.Name)
				missing++
			} else {
				fmt.Printf("  PRESENT   ✔  %-16s %s\n", check.Name, output)
				present++
			}
		}

		fmt.Println()
		if missing == 0 {
			ui.PrintSuccess(fmt.Sprintf("All %d components present! Environment fully matches %s store.", present, stackName))
		} else {
			fmt.Printf("  %d present, %d missing, %d outdated\n\n", present, missing, outdated)
			ui.PrintInfo(fmt.Sprintf("Run 'cue store install %s --missing-only' to fill the gaps.", stackName))
		}
	},
}
