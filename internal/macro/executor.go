package macro

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Execute runs a named macro after handling safety checks.
func Execute(name string, flags map[string]string) error {
	m, ok := Registry[name]
	if !ok {
		// Fuzzy match suggestion
		closest := findClosest(name)
		msg := fmt.Sprintf("Unknown macro: '%s'.", name)
		if closest != "" {
			msg += fmt.Sprintf(" Did you mean '%s'?", closest)
		}
		msg += " Run 'cue explain --list' to see all."
		return fmt.Errorf("%s", msg)
	}

	// Dangerous action gate
	isDangerous := m.Dangerous || (name == "git-undo" && flags["hard"] == "true")
	if isDangerous {
		ui.PrintWarning(fmt.Sprintf("'%s' is a destructive operation.", name))
		if !ui.Confirm("Are you sure you want to continue? [y/N] ") {
			ui.PrintInfo("Aborted. No changes were made.")
			return nil
		}
	}

	// Select and execute OS-appropriate steps
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "darwin"
	}
	steps := m.StepsForOS(osName, flags)

	// Handle --hard flag for git-undo
	if name == "git-undo" && flags["hard"] == "true" {
		steps = []Step{{OS: "all", Command: "git reset --hard HEAD~1"}}
	}

	for _, step := range steps {
		ui.PrintDim(fmt.Sprintf("$ %s", step.Command))
		if err := runShell(step.Command); err != nil {
			return fmt.Errorf("macro step failed: %w", err)
		}
	}

	// Print explanation
	cfg := config.DefaultConfig()
	if config.Current != nil {
		cfg = config.Current
	}
	if cfg.UI.ExplainAfterMacro {
		ui.PrintExplanation(m.Explanation)
	}

	return nil
}

// Explain prints a macro's explanation without executing.
func Explain(name string) error {
	m, ok := Registry[name]
	if !ok {
		return fmt.Errorf("unknown macro: %s", name)
	}
	ui.PrintHeader(m.Name)
	fmt.Printf("Category:    %s\n", m.Category)
	fmt.Printf("Description: %s\n", m.Description)
	fmt.Printf("Dangerous:   %v\n", m.Dangerous)
	if m.BuiltIn {
		fmt.Println("Source:      [built-in]")
	} else {
		fmt.Println("Source:      [custom]")
	}
	fmt.Println()
	fmt.Println("Commands:")
	for _, s := range m.Commands {
		fmt.Printf("  [%s] %s\n", s.OS, s.Command)
	}
	fmt.Println()
	fmt.Println("Explanation:")
	ui.PrintExplanation(m.Explanation)
	return nil
}

// ListAll prints all macros with their descriptions.
func ListAll() {
	ui.PrintHeader("Available Macros")
	headers := []string{"Macro", "Category", "Description", "Source"}
	var rows [][]string
	for _, m := range Registry {
		source := "[built-in]"
		if !m.BuiltIn {
			source = "[custom]"
		}
		rows = append(rows, []string{m.Name, m.Category, m.Description, source})
	}
	ui.PrintTable(headers, rows)
}

func runShell(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// findClosest returns the closest macro name via simple prefix matching.
func findClosest(input string) string {
	bestMatch := ""
	bestScore := 0
	for name := range Registry {
		score := commonPrefix(input, name)
		if score > bestScore {
			bestScore = score
			bestMatch = name
		}
	}
	if bestScore >= 2 {
		return bestMatch
	}
	return ""
}

func commonPrefix(a, b string) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}
