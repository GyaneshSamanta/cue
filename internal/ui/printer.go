package ui

import (
	"fmt"
	"os"
	"strings"
)

// ANSI escape codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

var colorEnabled = true

func init() {
	if os.Getenv("NO_COLOR") != "" {
		colorEnabled = false
	}
}

// SetColor enables or disables color output.
func SetColor(enabled bool) { colorEnabled = enabled }

func c(code, text string) string {
	if !colorEnabled {
		return text
	}
	return code + text + reset
}

func PrintSuccess(msg string)  { fmt.Fprintln(os.Stdout, c(green, "✔ "+msg)) }
func PrintError(msg string)    { fmt.Fprintln(os.Stderr, c(red, "✖ "+msg)) }
func PrintWarning(msg string)  { fmt.Fprintln(os.Stderr, c(yellow, "⚠ "+msg)) }
func PrintInfo(msg string)     { fmt.Fprintln(os.Stdout, c(cyan, "ℹ "+msg)) }
func PrintDim(msg string)      { fmt.Fprintln(os.Stdout, c(dim, msg)) }
func PrintStep(msg string)     { fmt.Fprintln(os.Stdout, c(blue, "▸ "+msg)) }

// PrintStatus overwrites the current line (for progress updates).
func PrintStatus(msg string) {
	fmt.Fprintf(os.Stdout, "\r%s", c(cyan, msg))
}

// PrintExplanation renders a macro explanation block.
func PrintExplanation(explanation string) {
	fmt.Fprintln(os.Stdout, c(dim, explanation))
}

// PrintHeader renders a bold section header.
func PrintHeader(title string) {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, c(bold+magenta, title))
	fmt.Fprintln(os.Stdout, c(dim, strings.Repeat("─", len(title)+4)))
}

// ProgressBar tracks multi-step progress.
type ProgressBar struct {
	Total   int
	Current int
}

func NewProgressBar(total int) *ProgressBar { return &ProgressBar{Total: total} }

func (p *ProgressBar) Update(msg string) {
	p.Current++
	pct := float64(p.Current) / float64(p.Total) * 100
	filled := int(pct / 100 * 30)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", 30-filled)
	fmt.Fprintf(os.Stdout, "\r  %s %3.0f%% %s", c(green, bar), pct, msg)
	if p.Current == p.Total {
		fmt.Fprintln(os.Stdout)
	}
}

func (p *ProgressBar) Finish() {
	fmt.Fprintln(os.Stdout)
}

// PrintTable renders a formatted table.
func PrintTable(headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Header
	sep := "+"
	hdr := "|"
	for i, h := range headers {
		sep += strings.Repeat("-", widths[i]+2) + "+"
		hdr += fmt.Sprintf(" %-*s |", widths[i], h)
	}
	fmt.Fprintln(os.Stdout, sep)
	fmt.Fprintln(os.Stdout, c(bold, hdr))
	fmt.Fprintln(os.Stdout, sep)

	// Rows
	for _, row := range rows {
		line := "|"
		for i, cell := range row {
			if i < len(widths) {
				line += fmt.Sprintf(" %-*s |", widths[i], cell)
			}
		}
		fmt.Fprintln(os.Stdout, line)
	}
	fmt.Fprintln(os.Stdout, sep)
}
