package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

type sessionState struct {
	ActiveTag string `json:"active_tag"`
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage session tags for history grouping",
}

var tagSetCmd = &cobra.Command{
	Use:   "set [name]",
	Short: "Set the active project tag",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := sessionState{ActiveTag: args[0]}
		saveSession(s)
		ui.PrintSuccess(fmt.Sprintf("Active tag set to '%s'", args[0]))
	},
}

var tagClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the active project tag",
	Run: func(cmd *cobra.Command, args []string) {
		saveSession(sessionState{})
		ui.PrintSuccess("Active tag cleared")
	},
}

func init() {
	tagCmd.AddCommand(tagSetCmd, tagClearCmd)
}

func saveSession(s sessionState) {
	dir := config.ConfigDir()
	os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(filepath.Join(dir, "session.json"), data, 0644)
}

// LoadSessionTag reads the current session tag.
func LoadSessionTag() string {
	data, err := os.ReadFile(filepath.Join(config.ConfigDir(), "session.json"))
	if err != nil {
		return ""
	}
	var s sessionState
	json.Unmarshal(data, &s)
	return s.ActiveTag
}
