package store

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Stack defines a complete environment setup.
type Stack interface {
	Name() string
	Description() string
	EstimatedSizeMB() int
	Components() []Component
	VerificationChecks() []Check
}

// Component is a single installable tool.
type Component struct {
	Name           string
	Version        string
	OS             []string
	InstallMethod  InstallMethod
	Optional       bool
	OptionalPrompt string
	DependsOn      []string
}

// InstallMethod defines OS-specific install commands.
type InstallMethod struct {
	Linux   []string
	Darwin  []string
	Windows []string
	Script  string // fallback install command
}

// Check validates a tool is installed and working.
type Check struct {
	Name    string
	Command string
	Pattern string
}

// InstallOpts controls installation behavior.
type InstallOpts struct {
	Verify  bool
	DryRun  bool
	Force   bool
}

// ComponentResult tracks per-component install outcome.
type ComponentResult struct {
	Name    string
	Status  string // success | failed | skipped
	Version string
	Error   error
}

var stacks = map[string]Stack{}

// RegisterStack adds a stack to the registry.
func RegisterStack(s Stack) { stacks[s.Name()] = s }

// GetStack retrieves a stack by name.
func GetStack(name string) (Stack, error) {
	s, ok := stacks[name]
	if !ok {
		available := make([]string, 0, len(stacks))
		for k := range stacks {
			available = append(available, k)
		}
		return nil, fmt.Errorf("unknown store: '%s'. Available: %s", name, strings.Join(available, ", "))
	}
	return s, nil
}

// ListStacks returns all registered stack names.
func ListStacks() []Stack {
	result := make([]Stack, 0, len(stacks))
	for _, s := range stacks {
		result = append(result, s)
	}
	return result
}

// Install executes a full stack installation.
func Install(stackName string, a adapter.OSAdapter, opts InstallOpts) error {
	stack, err := GetStack(stackName)
	if err != nil {
		return err
	}

	components := stack.Components()
	osName := a.OSName()

	// Filter to OS-applicable components
	var applicable []Component
	for _, c := range components {
		if len(c.OS) == 0 || containsOS(c.OS, osName) {
			applicable = append(applicable, c)
		}
	}

	// Prompt for optional components
	var selected []Component
	for _, c := range applicable {
		if c.Optional {
			if ui.Confirm(fmt.Sprintf("  Install %s? %s [y/N] ", c.Name, c.OptionalPrompt)) {
				selected = append(selected, c)
			}
		} else {
			selected = append(selected, c)
		}
	}

	if opts.DryRun {
		ui.PrintHeader("Dry Run — No installations will occur")
		for _, c := range selected {
			fmt.Printf("  Would install: %s (%s)\n", c.Name, c.Version)
		}
		return nil
	}

	progress := ui.NewProgressBar(len(selected))
	results := make([]ComponentResult, 0, len(selected))

	for _, comp := range selected {
		progress.Update(fmt.Sprintf("Installing %s...", comp.Name))
		err := installComponent(comp, a)
		status := "success"
		if err != nil {
			status = "failed"
		}
		ver := probeVersion(comp)
		results = append(results, ComponentResult{Name: comp.Name, Status: status, Version: ver, Error: err})
	}

	// Print results table
	printResults(results)

	if opts.Verify {
		return Verify(stackName, a)
	}
	return nil
}

// Preview shows what a store would install.
func Preview(stackName string, a adapter.OSAdapter) error {
	stack, err := GetStack(stackName)
	if err != nil {
		return err
	}

	ui.PrintHeader(fmt.Sprintf("Store Preview: %s", stack.Name()))
	fmt.Printf("  %s\n", stack.Description())
	fmt.Printf("  Estimated download: ~%d MB\n\n", stack.EstimatedSizeMB())

	headers := []string{"Component", "Version", "Optional", "Install Method"}
	var rows [][]string
	osName := a.OSName()
	for _, c := range stack.Components() {
		if len(c.OS) == 0 || containsOS(c.OS, osName) {
			opt := "No"
			if c.Optional {
				opt = "Yes"
			}
			method := getMethodDesc(c.InstallMethod, osName)
			rows = append(rows, []string{c.Name, c.Version, opt, method})
		}
	}
	ui.PrintTable(headers, rows)
	return nil
}

// Remove uninstalls a full stack.
func Remove(stackName string, a adapter.OSAdapter, force bool) error {
	stack, err := GetStack(stackName)
	if err != nil {
		return err
	}
	ui.PrintHeader(fmt.Sprintf("Removing store: %s", stack.Name()))
	for _, comp := range stack.Components() {
		ui.PrintStep(fmt.Sprintf("Removing %s...", comp.Name))
		// attempt uninstall, continue on error
		if err := a.UninstallPackage(comp.Name); err != nil {
			ui.PrintWarning(fmt.Sprintf("Could not remove %s: %v", comp.Name, err))
		}
	}
	ui.PrintSuccess("Store removal complete.")
	return nil
}

func installComponent(comp Component, a adapter.OSAdapter) error {
	osName := a.OSName()
	// Try OS-specific packages first
	var pkgs []string
	switch osName {
	case "linux":
		pkgs = comp.InstallMethod.Linux
	case "darwin":
		pkgs = comp.InstallMethod.Darwin
	case "windows":
		pkgs = comp.InstallMethod.Windows
	}

	if len(pkgs) > 0 {
		for _, pkg := range pkgs {
			if err := a.InstallPackage(pkg, nil); err != nil {
				return err
			}
		}
		return nil
	}

	// Fallback: script/command
	if comp.InstallMethod.Script != "" {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", comp.InstallMethod.Script)
		} else {
			cmd = exec.Command("sh", "-c", comp.InstallMethod.Script)
		}
		return cmd.Run()
	}

	return fmt.Errorf("no install method for %s on %s", comp.Name, osName)
}

func probeVersion(comp Component) string {
	// Try common version flags
	for _, vFlag := range []string{"--version", "-v", "version"} {
		name := strings.Fields(comp.Name)[0]
		out, err := exec.Command(strings.ToLower(name), vFlag).Output()
		if err == nil {
			return strings.TrimSpace(strings.Split(string(out), "\n")[0])
		}
	}
	return "installed"
}

func printResults(results []ComponentResult) {
	headers := []string{"Component", "Status", "Version"}
	var rows [][]string
	for _, r := range results {
		rows = append(rows, []string{r.Name, r.Status, r.Version})
	}
	fmt.Println()
	ui.PrintTable(headers, rows)
}

func containsOS(list []string, os string) bool {
	for _, s := range list {
		if s == os {
			return true
		}
	}
	return false
}

func getMethodDesc(m InstallMethod, os string) string {
	switch os {
	case "linux":
		if len(m.Linux) > 0 {
			return strings.Join(m.Linux, ", ")
		}
	case "darwin":
		if len(m.Darwin) > 0 {
			return strings.Join(m.Darwin, ", ")
		}
	case "windows":
		if len(m.Windows) > 0 {
			return strings.Join(m.Windows, ", ")
		}
	}
	if m.Script != "" {
		return m.Script
	}
	return "N/A"
}
