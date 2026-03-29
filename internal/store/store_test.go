package store_test

import (
	"testing"

	"github.com/GyaneshSamanta/gyanesh-help/internal/store"
	_ "github.com/GyaneshSamanta/gyanesh-help/internal/store/stacks"
)

// TestStoreStacks verifies that the engine loads stacks correctly.
func TestStoreStacks(t *testing.T) {
	stacksList := store.ListStacks()
	if len(stacksList) == 0 {
		t.Fatal("No stacks available in the store registry")
	}

	for _, s := range stacksList {
		if s.Name() == "" {
			t.Error("Found a stack with an empty name")
		}
		if s.Description() == "" {
			t.Errorf("Stack %s has an empty description", s.Name())
		}
		if len(s.Components()) == 0 {
			t.Errorf("Stack %s has no components defined", s.Name())
		}
		
		// Validate that the stack can be retrieved back by name
		retrieved, err := store.GetStack(s.Name())
		if err != nil {
			t.Errorf("GetStack failed to retrieve '%s': %v", s.Name(), err)
		}
		if retrieved.Name() != s.Name() {
			t.Errorf("GetStack returned mismatched stack: expected %s, got %s", s.Name(), retrieved.Name())
		}
		
		// Validate components
		for _, comp := range s.Components() {
			if comp.Name == "" {
				t.Errorf("Stack %s contains a component with an empty name", s.Name())
			}
			if len(comp.InstallMethod.Linux) == 0 && len(comp.InstallMethod.Darwin) == 0 && len(comp.InstallMethod.Windows) == 0 && comp.InstallMethod.Script == "" {
				t.Errorf("Component %s in stack %s has no install methods defined", comp.Name, s.Name())
			}
		}
	}
}

// TestUnknownStack ensures getting a non-existent stack errors cleanly.
func TestUnknownStack(t *testing.T) {
	_, err := store.GetStack("this-should-not-exist-1234")
	if err == nil {
		t.Error("Expected error when requesting unknown stack, got nil")
	}
}
