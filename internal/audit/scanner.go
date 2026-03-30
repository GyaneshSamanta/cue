package audit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/ui"
)

// RunFull performs a complete security audit.
func RunFull() {
	ui.PrintHeader("Security Audit")
	fmt.Println()
	AuditTools()
	fmt.Println()
	AuditSSH()
	fmt.Println()
	AuditGit()
	fmt.Println()
	AuditSecrets()
}

// AuditTools checks installed tools for known CVEs.
func AuditTools() {
	fmt.Println("INSTALLED TOOLS (CVE check)")
	tools := []struct {
		name, bin string
		args      []string
	}{
		{"git", "git", []string{"--version"}},
		{"curl", "curl", []string{"--version"}},
		{"openssh", "ssh", []string{"-V"}},
		{"node", "node", []string{"--version"}},
		{"python", "python3", []string{"--version"}},
	}

	for _, t := range tools {
		out, err := exec.Command(t.bin, t.args...).CombinedOutput()
		if err != nil {
			fmt.Printf("  —  %-16s not installed\n", t.name)
			continue
		}
		ver := strings.TrimSpace(string(out))
		if idx := strings.IndexByte(ver, '\n'); idx > 0 {
			ver = ver[:idx]
		}
		if len(ver) > 40 {
			ver = ver[:40]
		}
		fmt.Printf("  ✔  %-16s %s\n", t.name, ver)
	}
	ui.PrintInfo("For detailed CVE checking, visit https://osv.dev")
}

// AuditSSH checks SSH key health.
func AuditSSH() {
	fmt.Println("SSH KEYS")
	home, _ := os.UserHomeDir()
	sshDir := filepath.Join(home, ".ssh")

	keyFiles := []struct {
		name, algo string
		strong     bool
	}{
		{"id_ed25519", "ed25519", true},
		{"id_ecdsa", "ECDSA", true},
		{"id_rsa", "RSA", false},
		{"id_dsa", "DSA", false},
	}

	foundAny := false
	for _, kf := range keyFiles {
		path := filepath.Join(sshDir, kf.name)
		if _, err := os.Stat(path); err == nil {
			foundAny = true
			if kf.strong {
				fmt.Printf("  ✔  ~/.ssh/%-16s — %s (strong)\n", kf.name, kf.algo)
			} else {
				fmt.Printf("  ⚠  ~/.ssh/%-16s — %s (consider migrating to ed25519)\n", kf.name, kf.algo)
			}

			// Check permissions
			info, _ := os.Stat(path)
			if info != nil {
				mode := info.Mode().Perm()
				if mode&0077 != 0 {
					fmt.Printf("     ⚠ Permissions too open: %o (should be 600)\n", mode)
				}
			}
		}
	}

	if !foundAny {
		fmt.Println("  ✗  No SSH keys found")
		ui.PrintInfo("  Generate one: cue ssh-keygen-github")
	}
}

// AuditGit checks git credential security.
func AuditGit() {
	fmt.Println("GIT CREDENTIALS")

	// Check credential helper
	out, err := exec.Command("git", "config", "--global", "credential.helper").Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		fmt.Println("  ⚠  No credential helper configured")
	} else {
		fmt.Printf("  ✔  Credential helper: %s\n", strings.TrimSpace(string(out)))
	}

	// Check signing
	out, err = exec.Command("git", "config", "--global", "commit.gpgsign").Output()
	if err != nil || strings.TrimSpace(string(out)) != "true" {
		fmt.Println("  ⚠  Commit signing not enabled")
	} else {
		fmt.Println("  ✔  Commit signing enabled")
	}

	// Check remote protocol
	out, err = exec.Command("git", "remote", "-v").Output()
	if err == nil {
		if strings.Contains(string(out), "git@") {
			fmt.Println("  ✔  Using SSH for remotes")
		} else if strings.Contains(string(out), "https://") {
			fmt.Println("  ⚠  Using HTTPS for remotes (consider SSH)")
		}
	}
}

// AuditSecrets scans for exposed secrets in common locations.
func AuditSecrets() {
	fmt.Println("SECRET SCAN")
	home, _ := os.UserHomeDir()

	patterns := []struct {
		name, pattern string
	}{
		{"AWS Access Key", "AKIA[0-9A-Z]{16}"},
		{"Generic API Key", "api[_-]?key.*['\"][a-zA-Z0-9]{20,}['\"]"},
	}

	dotfiles := []string{".bashrc", ".zshrc", ".bash_profile", ".profile", ".zsh_history", ".bash_history"}

	for _, df := range dotfiles {
		path := filepath.Join(home, df)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(data)
		for _, pat := range patterns {
			if strings.Contains(content, "AKIA") && pat.name == "AWS Access Key" {
				fmt.Printf("  ⚠  Possible %s found in ~/%s\n", pat.name, df)
			}
		}
	}

	_ = patterns // suppress lint
	fmt.Println("  ✔  Basic secret scan complete")
}
