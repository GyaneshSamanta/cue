package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/gyanesh-help/internal/history"
	"github.com/GyaneshSamanta/gyanesh-help/internal/tui"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View command history interactively",
	Run: func(cmd *cobra.Command, args []string) {
		tag, _ := cmd.Flags().GetString("tag")
		search, _ := cmd.Flags().GetString("search")
		since, _ := cmd.Flags().GetString("since")
		failOnly, _ := cmd.Flags().GetBool("failed")
		all, _ := cmd.Flags().GetBool("all")
		exportCSV, _ := cmd.Flags().GetString("export")
		listMode, _ := cmd.Flags().GetBool("list")

		opts := history.QueryOptions{
			Tag:      tag,
			Search:   search,
			FailOnly: failOnly,
			Limit:    20,
		}
		if all {
			opts.Limit = 50000
		}
		if since != "" {
			t, err := time.Parse("2006-01-02", since)
			if err == nil {
				opts.Since = t
			}
		}

		if exportCSV != "" {
			if err := history.ExportCSV(exportCSV, opts); err != nil {
				ui.PrintError(err.Error())
				return
			}
			ui.PrintSuccess(fmt.Sprintf("Exported to %s", exportCSV))
			return
		}

		entries, err := history.Query(opts)
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		if len(entries) == 0 {
			ui.PrintInfo("No history entries found.")
			return
		}

		// List mode prints the table
		if listMode || len(entries) == 0 {
			headers := []string{"#", "Timestamp", "Command", "Exit", "Duration", "Tag"}
			var rows [][]string
			for _, e := range entries {
				dur := fmt.Sprintf("%dms", e.DurationMs)
				rows = append(rows, []string{
					fmt.Sprintf("%d", e.ID),
					e.Timestamp[:19],
					truncate(e.Command, 40),
					fmt.Sprintf("%d", e.ExitCode),
					dur,
					e.ProjectTag,
				})
			}
			ui.PrintTable(headers, rows)
			return
		}

		// TUI interactive mode
		items := make([]tui.PickerItem, len(entries))
		for i, e := range entries {
			tagStr := "✓"
			if e.ExitCode != 0 {
				tagStr = "✗"
			}

			items[i] = tui.PickerItem{
				Name:        truncate(e.Command, 60),
				Description: fmt.Sprintf("[%s] %s", e.ProjectTag, e.Timestamp[:19]),
				Tag:         tagStr,
			}
		}

		idx, err := tui.RunPicker("Command History (Enter to run)", items)
		if err != nil || idx == -1 {
			return
		}

		selected := entries[idx]
		ui.PrintInfo(fmt.Sprintf("\nSelected command: %s", selected.Command))
		
		if ui.Confirm("Execute this command? [y/N] ") {
			// re-run the chosen command (handling sh/cmd.exe properly)
			var reruncmd *exec.Cmd
			if osAdapter.OSName() == "Windows" {
				reruncmd = exec.Command("cmd", "/C", selected.Command)
			} else {
				reruncmd = exec.Command("sh", "-c", selected.Command)
			}
			reruncmd.Stdout = os.Stdout
			reruncmd.Stderr = os.Stderr
			reruncmd.Stdin = os.Stdin
			err := reruncmd.Run()
			if err != nil {
				ui.PrintError(err.Error())
			}
		}
	},
}

func init() {
	historyCmd.Flags().String("tag", "", "Filter by project tag")
	historyCmd.Flags().String("search", "", "Full-text search")
	historyCmd.Flags().String("since", "", "Filter entries since date (YYYY-MM-DD)")
	historyCmd.Flags().Bool("failed", false, "Show only failed commands")
	historyCmd.Flags().Bool("all", false, "Show all entries")
	historyCmd.Flags().String("export", "", "Export to CSV file path")
	historyCmd.Flags().Bool("list", false, "Print as a table instead of TUI")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
