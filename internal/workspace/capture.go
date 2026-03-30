package workspace

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/store"
)

// Manifest describes a workspace backup.
type Manifest struct {
	CapturedAt      string              `json:"captured_at"`
	OS              string              `json:"os"`
	Hostname        string              `json:"hostname"`
	InstalledStores []store.StoreEntry  `json:"installed_stores"`
	ShellFiles      []string            `json:"shell_files"`
	VSCodeFiles     []string            `json:"vscode_files,omitempty"`
	CustomMacros    bool                `json:"custom_macros"`
	HistoryIncluded bool                `json:"history_included"`
}

// Capture captures the current workspace state to a temp directory.
func Capture(a adapter.OSAdapter) (*Manifest, string, error) {
	tmpDir, err := os.MkdirTemp("", "gyanesh-workspace-*")
	if err != nil {
		return nil, "", err
	}

	manifest := &Manifest{
		CapturedAt: time.Now().UTC().Format(time.RFC3339),
		OS:         a.OSName(),
	}
	manifest.Hostname, _ = os.Hostname()

	// 1. Shell configs
	shellDir := filepath.Join(tmpDir, "shell")
	os.MkdirAll(shellDir, 0755)
	for _, sf := range a.ShellConfigPaths() {
		base := filepath.Base(sf)
		if copyIfExists(sf, filepath.Join(shellDir, base)) == nil {
			manifest.ShellFiles = append(manifest.ShellFiles, base)
		}
	}

	// Copy .gitconfig (sanitised)
	gitconfig := filepath.Join(a.HomeDir(), ".gitconfig")
	gitconfigDst := filepath.Join(shellDir, ".gitconfig")
	if copyIfExists(gitconfig, gitconfigDst) == nil {
		sanitiseGitConfig(gitconfigDst)
		manifest.ShellFiles = append(manifest.ShellFiles, ".gitconfig")
	}

	// 2. Installed stores
	manifest.InstalledStores = store.QueryInstalled(a)

	// 3. Custom macros
	macroSrc := filepath.Join(a.ConfigDir(), "macros.toml")
	if copyIfExists(macroSrc, filepath.Join(tmpDir, "macros.toml")) == nil {
		manifest.CustomMacros = true
	}

	// 4. Write .gitignore
	writeGitignore(filepath.Join(tmpDir, ".gitignore"))

	// 5. Write manifest
	data, _ := json.MarshalIndent(manifest, "", "  ")
	os.WriteFile(filepath.Join(tmpDir, "manifest.json"), data, 0644)

	return manifest, tmpDir, nil
}

func copyIfExists(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	os.MkdirAll(filepath.Dir(dst), 0755)
	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

// sanitiseGitConfig strips known secret patterns.
func sanitiseGitConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	content := string(data)
	patterns := []string{
		`(?i)(token|password|secret|key)\s*=\s*.+`,
		`(?i)(GITHUB_TOKEN|AWS_SECRET|API_KEY)\s*=\s*.+`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		content = re.ReplaceAllString(content, "# REDACTED by cue")
	}
	os.WriteFile(path, []byte(content), 0644)
}

func writeGitignore(path string) {
	ignore := strings.Join([]string{
		"*.pem", "*.key", "id_rsa", "id_ed25519",
		".env", ".env.*", "*.secret", "*.token", ".netrc",
	}, "\n") + "\n"
	os.WriteFile(path, []byte(ignore), 0644)
}
