package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/gyanesh-help/internal/config"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

const (
	repoOwner = "GyaneshSamanta"
	repoName  = "gyanesh-help"
	apiURL    = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases/latest"
)

// Release represents a GitHub release.
type Release struct {
	TagName string  `json:"tag_name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a release binary.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// CheckUpdate queries GitHub for the latest version.
func CheckUpdate(currentVersion string) (*Release, bool, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, fmt.Errorf("failed to parse release: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	if latest != current && latest > current {
		return &release, true, nil
	}
	return &release, false, nil
}

// Update downloads and installs the latest release binary.
func Update(release *Release, currentVersion string) error {
	// Find the right asset for this platform
	osName := runtime.GOOS
	arch := runtime.GOARCH
	ext := ".tar.gz"
	if osName == "windows" {
		ext = ".zip"
	}

	assetName := fmt.Sprintf("gyanesh-help-%s-%s%s", osName, arch, ext)
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s in release %s", osName, arch, release.TagName)
	}

	// Backup current binary
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary: %w", err)
	}

	backupDir := filepath.Join(config.ConfigDir(), "bin")
	os.MkdirAll(backupDir, 0755)
	backupPath := filepath.Join(backupDir, fmt.Sprintf("gyanesh-help-%s", currentVersion))

	data, err := os.ReadFile(execPath)
	if err == nil {
		os.WriteFile(backupPath, data, 0755)
		ui.PrintStep(fmt.Sprintf("Backed up current binary to %s", backupPath))
	}

	// Download new binary
	ui.PrintStep(fmt.Sprintf("Downloading %s...", release.TagName))

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(os.TempDir(), assetName)
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	io.Copy(out, resp.Body)
	out.Close()

	// Extract and replace
	if osName == "windows" {
		// Use PowerShell to extract zip
		extractDir := filepath.Join(os.TempDir(), "gyanesh-help-update")
		os.MkdirAll(extractDir, 0755)
		exec.Command("powershell", "-Command",
			fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force", tmpFile, extractDir)).Run()

		newBin := filepath.Join(extractDir, "gyanesh-help.exe")
		if _, err := os.Stat(newBin); err == nil {
			// Copy new binary (can't replace running exe on Windows directly)
			ui.PrintWarning("On Windows, please replace the binary manually from: " + extractDir)
		}
	} else {
		// Extract tar.gz
		extractDir := filepath.Join(os.TempDir(), "gyanesh-help-update")
		os.MkdirAll(extractDir, 0755)
		exec.Command("tar", "-xzf", tmpFile, "-C", extractDir).Run()

		newBin := filepath.Join(extractDir, "gyanesh-help")
		if _, err := os.Stat(newBin); err == nil {
			os.Rename(newBin, execPath)
			os.Chmod(execPath, 0755)
		}
	}

	os.Remove(tmpFile)
	ui.PrintSuccess(fmt.Sprintf("Updated to %s", release.TagName))
	return nil
}

// Rollback restores the previous binary version.
func Rollback(currentVersion string) error {
	backupDir := filepath.Join(config.ConfigDir(), "bin")
	entries, err := os.ReadDir(backupDir)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("no backup found to rollback to")
	}

	// Get the most recent backup
	latest := entries[len(entries)-1]
	backupPath := filepath.Join(backupDir, latest.Name())

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary: %w", err)
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(execPath, data, 0755); err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Rolled back to: %s", latest.Name()))
	return nil
}
