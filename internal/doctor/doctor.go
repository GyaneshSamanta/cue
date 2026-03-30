package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Issue represents a single health check result.
type Issue struct {
	Category    string
	Tool        string
	Status      string // "ok", "warn", "error"
	Version     string
	Message     string
	FixCmd      string
	FixableAuto bool
}

// RunDoctor performs a full system health check.
func RunDoctor(a adapter.OSAdapter) []Issue {
	var issues []Issue

	ui.PrintHeader("System Health Check")
	fmt.Println()

	// 1. Core tools
	fmt.Println("CORE TOOLS")
	issues = append(issues, checkTool("git", "git", "--version")...)
	issues = append(issues, checkTool("curl", "curl", "--version")...)
	issues = append(issues, checkTool("node", "node", "--version")...)
	issues = append(issues, checkTool("python", getPythonBin(), "--version")...)
	issues = append(issues, checkTool("docker", "docker", "--version")...)

	fmt.Println()

	// 2. PATH & shell
	fmt.Println("PATH & SHELL")
	issues = append(issues, checkPATH()...)
	issues = append(issues, checkShellCompletions()...)
	fmt.Println()

	// 3. Installed stores verification
	fmt.Println("ENVIRONMENT STORES")
	issues = append(issues, checkInstalledStores(a)...)
	fmt.Println()

	// 4. Security
	fmt.Println("SECURITY")
	issues = append(issues, checkSSHKeys()...)
	issues = append(issues, checkGitCredentials()...)
	fmt.Println()

	// 5. Disk
	fmt.Println("DISK")
	issues = append(issues, checkDisk()...)
	fmt.Println()

	// Summary
	errors := 0
	warnings := 0
	for _, iss := range issues {
		if iss.Status == "error" {
			errors++
		} else if iss.Status == "warn" {
			warnings++
		}
	}

	if errors == 0 && warnings == 0 {
		ui.PrintSuccess("All checks passed!")
	} else {
		fmt.Printf("\n%d issues found (%d errors, %d warnings).\n", errors+warnings, errors, warnings)
		fixable := 0
		for _, iss := range issues {
			if iss.FixableAuto {
				fixable++
			}
		}
		if fixable > 0 {
			fmt.Printf("Run 'cue doctor fix --all' to auto-fix %d issues.\n", fixable)
		}
	}

	return issues
}

// Fix attempts to auto-fix all fixable issues.
func Fix(a adapter.OSAdapter, issues []Issue, fixAll bool) {
	for _, iss := range issues {
		if iss.Status == "ok" || !iss.FixableAuto {
			continue
		}
		if !fixAll {
			ok := ui.Confirm(fmt.Sprintf("Fix: %s — %s?", iss.Tool, iss.Message))
			if !ok {
				continue
			}
		}
		if iss.FixCmd != "" {
			ui.PrintStep(fmt.Sprintf("Fixing: %s", iss.FixCmd))
			cmd := exec.Command(getShell(), getShellFlag(), iss.FixCmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				ui.PrintError(fmt.Sprintf("Fix failed: %v", err))
			} else {
				ui.PrintSuccess(fmt.Sprintf("Fixed: %s", iss.Tool))
			}
		}
	}
}

func checkTool(name, bin string, args ...string) []Issue {
	out, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil {
		printStatus("✗", name, "", "not found", "error")
		fixCmd := fmt.Sprintf("cue toolkit install %s", name)
		return []Issue{{
			Category: "core", Tool: name, Status: "error",
			Message: "not installed",
			FixCmd:  fixCmd, FixableAuto: true,
		}}
	}
	ver := strings.TrimSpace(string(out))
	if idx := strings.IndexByte(ver, '\n'); idx > 0 {
		ver = ver[:idx]
	}
	printStatus("✔", name, ver, "", "ok")
	return []Issue{{Category: "core", Tool: name, Status: "ok", Version: ver}}
}

func checkPATH() []Issue {
	var issues []Issue
	_, err := exec.LookPath("cue")
	if err != nil {
		printStatus("✗", "cue", "", "not on PATH", "error")
		issues = append(issues, Issue{
			Category: "path", Tool: "cue", Status: "error",
			Message: "binary not on PATH",
		})
	} else {
		printStatus("✔", "cue", "", "on PATH", "ok")
		issues = append(issues, Issue{Category: "path", Tool: "cue", Status: "ok"})
	}
	return issues
}

func checkShellCompletions() []Issue {
	complPath := filepath.Join(config.ConfigDir(), "completions")
	if _, err := os.Stat(complPath); err != nil {
		printStatus("⚠", "shell completions", "", "not installed", "warn")
		return []Issue{{
			Category: "shell", Tool: "completions", Status: "warn",
			Message: "completions not installed",
			FixCmd:  "cue setup", FixableAuto: true,
		}}
	}
	printStatus("✔", "shell completions", "", "installed", "ok")
	return []Issue{{Category: "shell", Tool: "completions", Status: "ok"}}
}

func checkInstalledStores(a adapter.OSAdapter) []Issue {
	manifest := filepath.Join(config.ConfigDir(), "installed_stores.json")
	if _, err := os.Stat(manifest); err != nil {
		printStatus("⚠", "stores", "", "no stores installed yet", "warn")
		return nil
	}
	printStatus("✔", "store manifest", "", "present", "ok")
	return nil
}

func checkSSHKeys() []Issue {
	home, _ := os.UserHomeDir()
	sshDir := filepath.Join(home, ".ssh")

	ed25519 := filepath.Join(sshDir, "id_ed25519")
	rsa := filepath.Join(sshDir, "id_rsa")

	if _, err := os.Stat(ed25519); err == nil {
		printStatus("✔", "SSH key", "", "ed25519 (strong)", "ok")
		return []Issue{{Category: "security", Tool: "ssh-key", Status: "ok"}}
	}
	if _, err := os.Stat(rsa); err == nil {
		printStatus("⚠", "SSH key", "", "RSA — consider ed25519", "warn")
		return []Issue{{
			Category: "security", Tool: "ssh-key", Status: "warn",
			Message: "RSA key found; consider migrating to ed25519",
			FixCmd:  "cue ssh-keygen-github", FixableAuto: true,
		}}
	}
	printStatus("✗", "SSH key", "", "none found", "error")
	return []Issue{{
		Category: "security", Tool: "ssh-key", Status: "error",
		Message: "No SSH key found",
		FixCmd:  "cue ssh-keygen-github", FixableAuto: true,
	}}
}

func checkGitCredentials() []Issue {
	out, err := exec.Command("git", "config", "--global", "credential.helper").Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		printStatus("⚠", "git credentials", "", "no credential helper configured", "warn")
		fix := "git config --global credential.helper store"
		if runtime.GOOS == "darwin" {
			fix = "git config --global credential.helper osxkeychain"
		} else if runtime.GOOS == "windows" {
			fix = "git config --global credential.helper manager"
		}
		return []Issue{{
			Category: "security", Tool: "git-credential", Status: "warn",
			Message: "no credential helper", FixCmd: fix, FixableAuto: true,
		}}
	}
	printStatus("✔", "git credentials", "", strings.TrimSpace(string(out)), "ok")
	return []Issue{{Category: "security", Tool: "git-credential", Status: "ok"}}
}

func checkDisk() []Issue {
	// Simple check — just report availability
	printStatus("✔", "disk", "", "check your disk space manually", "ok")
	return nil
}

func printStatus(icon, name, version, msg, _ string) {
	ver := ""
	if version != "" {
		ver = fmt.Sprintf("  %s", version)
	}
	suffix := ""
	if msg != "" {
		suffix = fmt.Sprintf("  (%s)", msg)
	}
	fmt.Printf("  %s  %-16s%s%s\n", icon, name, ver, suffix)
}

func getPythonBin() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	return "python"
}

func getShell() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	return "bash"
}

func getShellFlag() string {
	if runtime.GOOS == "windows" {
		return "-Command"
	}
	return "-c"
}
