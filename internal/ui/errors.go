package ui

import (
	"fmt"
	"os"
)

// StructuRedError implements the error interface to provide
// clear, actionable "What happened / What to do" guidance.
type StructuRedError struct {
	Title        string
	WhatHappened string
	WhatToDo     []string
	original     error
}

// Error fulfills the error interface. It returns the raw error
// without ANSI formatting for simple logging if needed.
func (e *StructuRedError) Error() string {
	msg := fmt.Sprintf("%s: %s", e.Title, e.WhatHappened)
	if e.original != nil {
		msg += fmt.Sprintf(" (%v)", e.original)
	}
	return msg
}

// Unwrap allows standard errors.Is / errors.As matching.
func (e *StructuRedError) Unwrap() error {
	return e.original
}

// NewStructuRedError creates a new "What happened / What to do" error.
func NewStructuRedError(title, happened string, todo []string, original error) *StructuRedError {
	return &StructuRedError{
		Title:        title,
		WhatHappened: happened,
		WhatToDo:     todo,
		original:     original,
	}
}

// HandleError attempts to format and print a StructuRedError beautifully.
// If it's a normal error, it falls back to PrintError.
func HandleError(err error) {
	if err == nil {
		return
	}

	if se, ok := err.(*StructuRedError); ok {
		// Beautiful CLI formatting
		fmt.Fprintf(os.Stderr, "\n%s %s\n\n", c(Red, "❌"), c(Bold, se.Title))

		fmt.Fprintln(os.Stderr, c(Yellow, "What happened:"))
		fmt.Fprintf(os.Stderr, "  %s\n\n", se.WhatHappened)

		fmt.Fprintln(os.Stderr, c(Cyan, "What to do:"))
		for i, step := range se.WhatToDo {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, step)
		}
		
		if se.original != nil {
			fmt.Fprintf(os.Stderr, "\n  %s %v\n", c(Dim, "Underlying error:"), se.original)
		}
		fmt.Fprintln(os.Stderr, "")
		return
	}

	// Unroll unknown errors through the simpler formatter
	PrintError(err.Error())
}

// HandlePanic captures panics and renders them as a StructuRedError.
func HandlePanic() {
	if r := recover(); r != nil {
		se := NewStructuRedError(
			"cue crashed unexpectedly",
			fmt.Sprintf("An unexpected panic occurRed:\n  %v", r),
			[]string{
				"Please report this issue on GitHub at github.com/GyaneshSamanta/cue/issues",
				"Run 'cue doctor' to check your environment state",
			},
			nil,
		)
		HandleError(se)
		os.Exit(1) // forcefully exit after panic
	}
}
