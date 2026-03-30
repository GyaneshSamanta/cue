package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func init() {
	registerGoMacros()
}

func registerGoMacros() {
	macro.Register(&macro.Macro{
		Name:        "go-mod-tidy-check",
		Command:     "go mod tidy && go vet ./... && go test ./...",
		Description: "Go module cleanup + vet + test",
		Explanation: `Runs the standard Go project health check:
1. go mod tidy — removes unused dependencies, adds missing ones
2. go vet ./... — reports suspicious constructs (potential bugs)
3. go test ./... — runs all tests in all packages
This is the recommended pre-commit/pre-push check for Go projects.`,
		Dangerous: false,
	})
}
