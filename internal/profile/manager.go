package profile

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/GyaneshSamanta/gyanesh-help/internal/config"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

// Profile represents a named environment profile.
type Profile struct {
	Identity IdentityConfig `toml:"identity"`
	Git      GitConfig      `toml:"git"`
	Env      EnvConfig      `toml:"env"`
	Shell    ShellConfig    `toml:"shell"`
}

type IdentityConfig struct {
	Name string `toml:"name"`
}

type GitConfig struct {
	Name       string `toml:"name"`
	Email      string `toml:"email"`
	SigningKey string `toml:"signing_key"`
}

type EnvConfig struct {
	PathPrepend []string          `toml:"PATH_prepend"`
	Variables   map[string]string `toml:"variables"`
}

type ShellConfig struct {
	DefaultTag   string `toml:"default_tag"`
	PromptPrefix string `toml:"prompt_prefix"`
}

// ProfilesDir returns the profiles directory.
func ProfilesDir() string {
	return filepath.Join(config.ConfigDir(), "profiles")
}

// Create creates a new profile.
func Create(name string) error {
	dir := ProfilesDir()
	os.MkdirAll(dir, 0755)

	profilePath := filepath.Join(dir, name+".toml")
	if _, err := os.Stat(profilePath); err == nil {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	content := fmt.Sprintf(`[identity]
name = "%s"

[git]
name = ""
email = ""
signing_key = ""

[env]
PATH_prepend = []

[shell]
default_tag = "%s"
prompt_prefix = "[%s]"
`, name, name, strings.ToUpper(name))

	os.WriteFile(profilePath, []byte(content), 0644)
	ui.PrintSuccess(fmt.Sprintf("Profile '%s' created at %s", name, profilePath))
	ui.PrintInfo("Edit the file to configure git, env variables, and shell settings.")
	return nil
}

// Switch activates a profile.
func Switch(name string) error {
	profilePath := filepath.Join(ProfilesDir(), name+".toml")
	if _, err := os.Stat(profilePath); err != nil {
		return fmt.Errorf("profile '%s' not found", name)
	}

	var profile Profile
	if _, err := toml.DecodeFile(profilePath, &profile); err != nil {
		return fmt.Errorf("invalid profile: %w", err)
	}

	// Apply git config
	if profile.Git.Name != "" {
		exec.Command("git", "config", "--global", "user.name", profile.Git.Name).Run()
		ui.PrintStep(fmt.Sprintf("Git name → %s", profile.Git.Name))
	}
	if profile.Git.Email != "" {
		exec.Command("git", "config", "--global", "user.email", profile.Git.Email).Run()
		ui.PrintStep(fmt.Sprintf("Git email → %s", profile.Git.Email))
	}
	if profile.Git.SigningKey != "" {
		exec.Command("git", "config", "--global", "user.signingkey", profile.Git.SigningKey).Run()
	}

	// Apply env variables
	for key, val := range profile.Env.Variables {
		os.Setenv(key, val)
		if runtime.GOOS == "windows" {
			exec.Command("setx", key, val).Run()
		}
	}

	// Save active profile
	os.WriteFile(filepath.Join(config.ConfigDir(), "active-profile"), []byte(name), 0644)

	ui.PrintSuccess(fmt.Sprintf("Switched to profile: %s", name))
	return nil
}

// ListProfiles shows all available profiles.
func ListProfiles() error {
	dir := ProfilesDir()
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		ui.PrintInfo("No profiles created. Use 'gyanesh-help profile create <name>' to create one.")
		return nil
	}

	active := ""
	data, _ := os.ReadFile(filepath.Join(config.ConfigDir(), "active-profile"))
	active = strings.TrimSpace(string(data))

	ui.PrintHeader("Profiles")
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".toml")
		marker := "  "
		if name == active {
			marker = "● "
		}
		fmt.Printf("  %s%s\n", marker, name)
	}
	return nil
}

// Current shows the active profile.
func Current() {
	data, err := os.ReadFile(filepath.Join(config.ConfigDir(), "active-profile"))
	if err != nil {
		ui.PrintInfo("No active profile. Use 'gyanesh-help profile switch <name>'.")
		return
	}
	fmt.Printf("Active profile: %s\n", strings.TrimSpace(string(data)))
}

// DiffProfiles shows differences between two profiles.
func DiffProfiles(a, b string) error {
	pathA := filepath.Join(ProfilesDir(), a+".toml")
	pathB := filepath.Join(ProfilesDir(), b+".toml")

	var profileA, profileB Profile
	if _, err := toml.DecodeFile(pathA, &profileA); err != nil {
		return fmt.Errorf("cannot read profile '%s': %w", a, err)
	}
	if _, err := toml.DecodeFile(pathB, &profileB); err != nil {
		return fmt.Errorf("cannot read profile '%s': %w", b, err)
	}

	ui.PrintHeader(fmt.Sprintf("Profile Diff: %s vs %s", a, b))

	if profileA.Git.Name != profileB.Git.Name {
		fmt.Printf("  git.name:  %s → %s\n", profileA.Git.Name, profileB.Git.Name)
	}
	if profileA.Git.Email != profileB.Git.Email {
		fmt.Printf("  git.email: %s → %s\n", profileA.Git.Email, profileB.Git.Email)
	}
	if profileA.Shell.DefaultTag != profileB.Shell.DefaultTag {
		fmt.Printf("  tag:       %s → %s\n", profileA.Shell.DefaultTag, profileB.Shell.DefaultTag)
	}

	return nil
}
