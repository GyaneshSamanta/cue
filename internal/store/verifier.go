package store

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Verify checks that all stack components are reachable and working.
func Verify(stackName string, a adapter.OSAdapter) error {
	stack, err := GetStack(stackName)
	if err != nil {
		return err
	}

	ui.PrintHeader(fmt.Sprintf("Verifying: %s", stack.Name()))
	checks := stack.VerificationChecks()
	headers := []string{"Check", "Status", "Output"}
	var rows [][]string
	passed := 0

	for _, check := range checks {
		parts := strings.Fields(check.Command)
		if len(parts) == 0 {
			continue
		}
		out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
		output := strings.TrimSpace(string(out))
		if len(output) > 60 {
			output = output[:60] + "..."
		}

		status := "✔ PASS"
		if err != nil {
			status = "✖ FAIL"
		} else if check.Pattern != "" {
			matched, _ := regexp.MatchString(check.Pattern, string(out))
			if !matched {
				status = "⚠ WARN"
			} else {
				passed++
			}
		} else {
			passed++
		}
		rows = append(rows, []string{check.Name, status, output})
	}

	ui.PrintTable(headers, rows)
	fmt.Printf("\n  %d/%d checks passed.\n", passed, len(checks))
	return nil
}

// QueryInstalled returns which stores have been installed (based on verification).
func QueryInstalled(a adapter.OSAdapter) []StoreEntry {
	var installed []StoreEntry
	for _, stack := range ListStacks() {
		checks := stack.VerificationChecks()
		passing := 0
		for _, check := range checks {
			parts := strings.Fields(check.Command)
			if len(parts) > 0 && exec.Command(parts[0], parts[1:]...).Run() == nil {
				passing++
			}
		}
		if passing > len(checks)/2 {
			installed = append(installed, StoreEntry{Name: stack.Name()})
		}
	}
	return installed
}

// StoreEntry represents an installed store in the manifest.
type StoreEntry struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}
