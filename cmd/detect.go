package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/gyanesh-help/internal/project"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Auto-detect project type in the current directory",
	Long: `Scans the current directory for project type indicators
(package.json, go.mod, Cargo.toml, etc.) and recommends the appropriate store.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		detections := project.Detect(dir)
		project.PrintDetection(detections)
	},
}

var initProjectCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Scaffold a new project with stack-aware templates",
	Long: `Create a new project directory with the right scaffolding, git init,
.gitignore, and .gyanesh-help project configuration.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		scaffolds := project.ListScaffolds()

		fmt.Println("\nWhat kind of project?")
		options := make([]string, len(scaffolds))
		for i, s := range scaffolds {
			options[i] = fmt.Sprintf("%s — %s", s.Name, s.Description)
		}
		idx := ui.SelectOne("", options)
		if idx < 0 {
			return
		}

		if err := project.Scaffold(name, idx); err != nil {
			ui.PrintError(err.Error())
		}
	},
}
