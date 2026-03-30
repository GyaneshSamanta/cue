package workspace

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Backup captures workspace and pushes to GitHub.
func Backup(a adapter.OSAdapter, token, repoName string) error {
	ui.PrintStep("Capturing workspace state...")
	manifest, tmpDir, err := Capture(a)
	if err != nil {
		return fmt.Errorf("capture failed: %w", err)
	}

	ui.PrintInfo(fmt.Sprintf("Captured: %d shell files, %d stores, macros: %v",
		len(manifest.ShellFiles), len(manifest.InstalledStores), manifest.CustomMacros))

	ui.PrintStep("Pushing to GitHub...")
	url, err := pushToGitHub(tmpDir, token, repoName)
	if err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Backup complete: %s", url))
	return nil
}

func pushToGitHub(srcDir, token, repoName string) (string, error) {
	// Initialize git repo in temp dir
	if err := runGit(srcDir, "init"); err != nil {
		return "", err
	}
	if err := runGit(srcDir, "add", "."); err != nil {
		return "", err
	}

	commitMsg := fmt.Sprintf("cue backup: %s", time.Now().UTC().Format(time.RFC3339))
	if err := runGit(srcDir, "commit", "-m", commitMsg); err != nil {
		return "", err
	}

	remoteURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", token, "user", repoName)
	runGit(srcDir, "remote", "add", "origin", remoteURL)

	if err := runGit(srcDir, "push", "--force", "origin", "HEAD:main"); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://github.com/user/%s", repoName), nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}
