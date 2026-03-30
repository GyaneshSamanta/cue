package toolkit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Tool represents a single installable developer tool.
type Tool struct {
	Name           string
	DisplayName    string
	Description    string
	VersionManager string // e.g. "nvm", "pyenv", "rustup", ""
	EstSizeMB      int
	Categories     []string
	InstallFunc    func(a adapter.OSAdapter, version string) error
	VerifyFunc     func() (string, bool) // returns version, isInstalled
	UpgradeFunc    func(a adapter.OSAdapter) error
	RemoveFunc     func(a adapter.OSAdapter) error
}

// Registry holds all available toolkit tools.
var Registry = map[string]*Tool{}

// Register adds a tool to the registry.
func Register(t *Tool) {
	Registry[t.Name] = t
}

// List returns all tools sorted by name.
func List() []*Tool {
	tools := make([]*Tool, 0, len(Registry))
	for _, t := range Registry {
		tools = append(tools, t)
	}
	return tools
}

// Install installs a single tool, bootstrapping its version manager first if needed.
func Install(name string, a adapter.OSAdapter, version string, dryRun bool) error {
	t, ok := Registry[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s. Run 'cue toolkit list' to see available tools", name)
	}

	if dryRun {
		ui.PrintHeader(fmt.Sprintf("[DRY RUN] toolkit install %s", name))
		fmt.Printf("  Tool:      %s\n", t.DisplayName)
		fmt.Printf("  Method:    %s\n", installMethod(t))
		fmt.Printf("  Size:      ~%d MB\n", t.EstSizeMB)
		if t.VersionManager != "" {
			fmt.Printf("  Version Mgr: %s (will be bootstrapped if missing)\n", t.VersionManager)
		}
		if version != "" {
			fmt.Printf("  Version:   %s\n", version)
		} else {
			fmt.Printf("  Version:   latest stable\n")
		}
		return nil
	}

	// Check if already installed
	if ver, installed := t.VerifyFunc(); installed {
		ui.PrintSuccess(fmt.Sprintf("%s is already installed (version %s)", t.DisplayName, ver))
		return nil
	}

	ui.PrintHeader(fmt.Sprintf("Installing %s", t.DisplayName))
	return t.InstallFunc(a, version)
}

// Upgrade upgrades a tool to the latest version.
func Upgrade(name string, a adapter.OSAdapter) error {
	t, ok := Registry[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}
	if t.UpgradeFunc == nil {
		return fmt.Errorf("%s does not support upgrade. Reinstall with 'toolkit install %s'", name, name)
	}
	return t.UpgradeFunc(a)
}

// Remove removes a tool.
func Remove(name string, a adapter.OSAdapter) error {
	t, ok := Registry[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}
	if t.RemoveFunc == nil {
		return fmt.Errorf("%s does not support removal via cue", name)
	}
	return t.RemoveFunc(a)
}

// Info returns details about a tool.
func Info(name string) error {
	t, ok := Registry[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}
	ui.PrintHeader(t.DisplayName)
	fmt.Printf("  Name:        %s\n", t.Name)
	fmt.Printf("  Description: %s\n", t.Description)
	fmt.Printf("  Install:     %s\n", installMethod(t))
	fmt.Printf("  Est. Size:   ~%d MB\n", t.EstSizeMB)
	if t.VersionManager != "" {
		fmt.Printf("  Version Mgr: %s\n", t.VersionManager)
	}
	ver, installed := t.VerifyFunc()
	if installed {
		ui.PrintSuccess(fmt.Sprintf("  Installed:   %s", ver))
	} else {
		ui.PrintWarning("  Not installed")
	}
	return nil
}

func installMethod(t *Tool) string {
	if t.VersionManager != "" {
		return fmt.Sprintf("via %s", t.VersionManager)
	}
	return "OS package manager / official installer"
}

// --- Helper functions for tool definitions ---

// CommandExists checks if a binary is on PATH.
func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetVersion runs a command and extracts the first version-like string.
func GetVersion(bin string, args ...string) (string, bool) {
	out, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil {
		return "", false
	}
	s := strings.TrimSpace(string(out))
	// Take first line
	if idx := strings.IndexByte(s, '\n'); idx > 0 {
		s = s[:idx]
	}
	return s, true
}

// RunInstallCmd runs an install command with stdout/stderr piped.
func RunInstallCmd(name string, args ...string) error {
	ui.PrintStep(fmt.Sprintf("→ %s %s", name, strings.Join(args, " ")))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PkgInstall installs via the OS package manager.
func PkgInstall(a adapter.OSAdapter, pkg string) error {
	return a.InstallPackage(pkg, nil)
}

// HomeDir returns the user's home directory.
func HomeDir() string {
	h, _ := os.UserHomeDir()
	return h
}

// ShellConfigPath returns the most likely shell config file.
func ShellConfigPath() string {
	home := HomeDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	}
	// Check zsh first, then bash
	for _, f := range []string{".zshrc", ".bashrc", ".bash_profile", ".profile"} {
		p := filepath.Join(home, f)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join(home, ".bashrc")
}

// AppendToShellConfig appends a line to the shell config if it doesn't exist.
func AppendToShellConfig(line string) error {
	path := ShellConfigPath()
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), line) {
		return nil // Already there
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("\n" + line + "\n")
	return err
}
