package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/store"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Restore clones a backup repo and restores the workspace.
func Restore(repoURL string, a adapter.OSAdapter) error {
	tmpDir, err := os.MkdirTemp("", "gyanesh-restore-*")
	if err != nil {
		return err
	}

	ui.PrintStep("Cloning backup repository...")
	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// Read manifest
	var manifest Manifest
	data, err := os.ReadFile(filepath.Join(tmpDir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("no manifest.json in backup: %w", err)
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	ui.PrintInfo(fmt.Sprintf("Restoring backup from %s (OS: %s, Host: %s)",
		manifest.CapturedAt, manifest.OS, manifest.Hostname))

	// 1. Install stores
	for _, s := range manifest.InstalledStores {
		ui.PrintStep(fmt.Sprintf("Installing store: %s", s.Name))
		if err := store.Install(s.Name, a, store.InstallOpts{Verify: false}); err != nil {
			ui.PrintWarning(fmt.Sprintf("Store '%s' failed: %v. Continuing...", s.Name, err))
		}
	}

	// 2. Copy shell config files
	for _, sf := range manifest.ShellFiles {
		src := filepath.Join(tmpDir, "shell", sf)
		dst := shellFileDest(sf, a)
		if _, err := os.Stat(dst); err == nil {
			if !ui.Confirm(fmt.Sprintf("  Overwrite existing %s? [y/N] ", dst)) {
				continue
			}
		}
		copyIfExists(src, dst)
		ui.PrintDim(fmt.Sprintf("  Restored: %s", dst))
	}

	// 3. Custom macros
	if manifest.CustomMacros {
		src := filepath.Join(tmpDir, "macros.toml")
		dst := filepath.Join(a.ConfigDir(), "macros.toml")
		copyIfExists(src, dst)
		ui.PrintDim("  Restored: custom macros")
	}

	ui.PrintSuccess("Restore complete. Open a new terminal for changes to take effect.")
	return nil
}

func shellFileDest(name string, a adapter.OSAdapter) string {
	return filepath.Join(a.HomeDir(), name)
}
