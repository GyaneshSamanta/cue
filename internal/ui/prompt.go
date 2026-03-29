package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks a y/N question and returns true if user confirms.
func Confirm(msg string) bool {
	fmt.Fprint(os.Stdout, c(Yellow, msg))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// SelectOne presents a numbered list and returns the chosen index.
func SelectOne(prompt string, options []string) int {
	fmt.Fprintln(os.Stdout, c(Cyan, prompt))
	for i, opt := range options {
		fmt.Fprintf(os.Stdout, "  %s %s\n", c(Bold, fmt.Sprintf("[%d]", i+1)), opt)
	}
	fmt.Fprint(os.Stdout, c(Cyan, "  Choice: "))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	for i := range options {
		if input == fmt.Sprintf("%d", i+1) {
			return i
		}
	}
	return -1
}

// ReadInput reads a single line of input from the user.
func ReadInput(prompt string) string {
	fmt.Fprint(os.Stdout, c(Cyan, prompt))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
