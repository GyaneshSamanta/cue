package team

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/GyaneshSamanta/gyanesh-help/internal/config"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

// TeamDir returns the team sync directory.
func TeamDir() string {
	return filepath.Join(config.ConfigDir(), "team")
}

// Init initializes a team config directory.
func Init() error {
	dir := TeamDir()
	os.MkdirAll(filepath.Join(dir, "macros"), 0755)
	os.MkdirAll(filepath.Join(dir, "stores"), 0755)
	os.MkdirAll(filepath.Join(dir, "config"), 0755)

	ui.PrintSuccess("Team directory initialized at " + dir)
	ui.PrintInfo("Use 'gyanesh-help team connect --repo <url>' to link a shared repo.")
	return nil
}

// Connect links a GitHub repo for team sync.
func Connect(repoURL string) error {
	dir := TeamDir()
	os.MkdirAll(dir, 0755)

	// Save repo URL
	os.WriteFile(filepath.Join(dir, ".repo"), []byte(repoURL), 0644)

	// Clone the repo
	ui.PrintStep("Cloning team repo...")
	cmd := exec.Command("git", "clone", repoURL, filepath.Join(dir, "repo"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone team repo: %w", err)
	}

	ui.PrintSuccess("Connected to team repo: " + repoURL)
	return nil
}

// Sync pulls latest team config.
func Sync(push bool) error {
	repoDir := filepath.Join(TeamDir(), "repo")
	if _, err := os.Stat(repoDir); err != nil {
		return fmt.Errorf("no team repo connected. Run 'gyanesh-help team connect --repo <url>'")
	}

	if push {
		ui.PrintStep("Pushing local team config...")
		cmd := exec.Command("git", "-C", repoDir, "add", "-A")
		cmd.Run()
		exec.Command("git", "-C", repoDir, "commit", "-m", "sync: update team config").Run()
		cmd = exec.Command("git", "-C", repoDir, "push")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	ui.PrintStep("Pulling latest team config...")
	cmd := exec.Command("git", "-C", repoDir, "pull", "--rebase")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	ui.PrintSuccess("Team config synced")
	return nil
}

// Disconnect removes the team repo link.
func Disconnect() error {
	repoDir := filepath.Join(TeamDir(), "repo")
	os.RemoveAll(repoDir)
	os.Remove(filepath.Join(TeamDir(), ".repo"))
	ui.PrintSuccess("Disconnected from team repo")
	return nil
}
