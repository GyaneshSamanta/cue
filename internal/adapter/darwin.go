//go:build darwin

package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// Detect returns the macOS adapter on Darwin builds.
func Detect() OSAdapter {
	return newDarwinAdapter()
}

type DarwinAdapter struct {
	homeDir string
}

func newDarwinAdapter() *DarwinAdapter {
	home, _ := os.UserHomeDir()
	return &DarwinAdapter{homeDir: home}
}

func (a *DarwinAdapter) PackageManagerName() string { return "brew" }

func (a *DarwinAdapter) InstallPackage(pkg string, args []string) error {
	cmdArgs := append([]string{"install"}, args...)
	cmdArgs = append(cmdArgs, pkg)
	cmd := exec.Command("brew", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *DarwinAdapter) UninstallPackage(pkg string) error {
	cmd := exec.Command("brew", "uninstall", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *DarwinAdapter) IsPackageInstalled(pkg string) bool {
	return exec.Command("brew", "list", pkg).Run() == nil
}

func (a *DarwinAdapter) PackageVersion(pkg string) (string, error) {
	out, err := exec.Command("brew", "list", "--versions", pkg).Output()
	if err != nil {
		return "", err
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) >= 2 {
		return parts[1], nil
	}
	return strings.TrimSpace(string(out)), nil
}

func (a *DarwinAdapter) LockPaths() []string {
	return []string{filepath.Join(a.homeDir, ".homebrew", "locks")}
}

func (a *DarwinAdapter) IsLocked() (bool, string, error) {
	// Check for active brew processes
	out, err := exec.Command("pgrep", "-f", "brew").Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return true, "active brew process detected", nil
	}
	// Check lock directory
	for _, lp := range a.LockPaths() {
		entries, _ := os.ReadDir(lp)
		if len(entries) > 0 {
			return true, fmt.Sprintf("lock files in %s", lp), nil
		}
	}
	return false, "", nil
}

func (a *DarwinAdapter) SuspendProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGSTOP)
}

func (a *DarwinAdapter) ResumeProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGCONT)
}

func (a *DarwinAdapter) KillProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGKILL)
}

func (a *DarwinAdapter) RunElevated(cmd string, args []string) error {
	if os.Getuid() == 0 {
		return exec.Command(cmd, args...).Run()
	}
	return exec.Command("sudo", append([]string{cmd}, args...)...).Run()
}

func (a *DarwinAdapter) HomeDir() string   { return a.homeDir }
func (a *DarwinAdapter) ConfigDir() string  { return filepath.Join(a.homeDir, ".cue") }
func (a *DarwinAdapter) OSName() string     { return "darwin" }
func (a *DarwinAdapter) OSDistro() string   { return "" }

func (a *DarwinAdapter) HasGPU() bool {
	// Apple Silicon has GPU but not NVIDIA
	return exec.Command("system_profiler", "SPDisplaysDataType").Run() == nil
}

func (a *DarwinAdapter) ShellConfigPaths() []string {
	paths := []string{}
	candidates := []string{
		filepath.Join(a.homeDir, ".zshrc"),
		filepath.Join(a.homeDir, ".bash_profile"),
		filepath.Join(a.homeDir, ".bashrc"),
		filepath.Join(a.homeDir, ".profile"),
		filepath.Join(a.homeDir, ".config", "fish", "config.fish"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			paths = append(paths, p)
		}
	}
	return paths
}
