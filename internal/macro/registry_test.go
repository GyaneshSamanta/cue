package macro_test

import (
	"strings"
	"testing"

	"github.com/GyaneshSamanta/cue/internal/macro"
	_ "github.com/GyaneshSamanta/cue/internal/macro/builtins"
)

// TestRegistryIntegrity verifies that every registered macro is structurally sound.
func TestRegistryIntegrity(t *testing.T) {
	if len(macro.Registry) == 0 {
		t.Fatal("Macro registry is empty, expected built-in macros to be registered")
	}

	for name, m := range macro.Registry {
		if m.Name == "" {
			t.Errorf("Macro mapped as '%s' has an empty 'Name' field", name)
		}
		if m.Name != name {
			t.Errorf("Macro '%s' is registered under key '%s'", m.Name, name)
		}
		if m.Description == "" {
			t.Errorf("Macro '%s' is missing a Description", name)
		}
		
		if len(m.Commands) == 0 && m.Command == "" {
			t.Errorf("Macro '%s' has no executable Commands or Command", name)
		}

		if m.Command != "" && strings.Contains(m.Command, "&&") {
			// This is just a warning, multi-step should ideally use Commands array but strings are fine
		}
		
		for _, step := range m.Commands {
			if step.Command == "" {
				t.Errorf("Macro '%s' has an empty step command for OS: %s", name, step.OS)
			}
		}
	}
}

func TestExplanationOutputs(t *testing.T) {
	for name, m := range macro.Registry {
		if m.Explanation == "" {
			t.Logf("Warning: Macro '%s' is missing an explanation block", name)
		}
	}
}
