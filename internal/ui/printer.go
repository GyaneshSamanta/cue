package ui

import (
	"fmt"
	"os"
	"strings"
)

// ANSI escape codes
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

var colorEnabled = true

func init() {
	if os.Getenv("NO_COLOR") != "" {
		colorEnabled = false
	}
}

// SetColor enables or disables color output.
func SetColor(enabled bool) { colorEnabled = enabled }

// C surrounds text with ANSI codes if colors are enabled.
func C(code, text string) string {
	if !colorEnabled {
		return text
	}
	return code + text + Reset
}

func c(code, text string) string {
	return C(code, text)
}

func PrintSuccess(msg string)  { fmt.Fprintln(os.Stdout, C(Green, "✔ "+msg)) }
func PrintError(msg string)    { fmt.Fprintln(os.Stderr, C(Red, "✖ "+msg)) }
func PrintWarning(msg string)  { fmt.Fprintln(os.Stderr, C(Yellow, "⚠ "+msg)) }
func PrintInfo(msg string)     { fmt.Fprintln(os.Stdout, C(Cyan, "ℹ "+msg)) }
func PrintDim(msg string)      { fmt.Fprintln(os.Stdout, C(Dim, msg)) }
func PrintStep(msg string)     { fmt.Fprintln(os.Stdout, C(Blue, "▸ "+msg)) }

// PrintStatus overwrites the current line (for progress updates).
func PrintStatus(msg string) {
	fmt.Fprintf(os.Stdout, "\r%s", C(Cyan, msg))
}

// PrintExplanation renders a macro explanation block.
func PrintExplanation(explanation string) {
	fmt.Fprintln(os.Stdout, C(Dim, explanation))
}

// PrintHeader renders a Bold section header.
func PrintHeader(title string) {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, C(Bold+Magenta, title))
	fmt.Fprintln(os.Stdout, C(Dim, strings.Repeat("─", len(title)+4)))
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
	fmt.Fprintf(os.Stdout, "\r  %s %3.0f%% %s", c(Green, bar), pct, msg)
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
	fmt.Fprintln(os.Stdout, c(Bold, hdr))
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
