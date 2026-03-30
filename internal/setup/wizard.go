package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// NeedsSetup returns true if first-run wizard should trigger.
func NeedsSetup() bool {
	_, err := os.Stat(filepath.Join(config.ConfigDir(), "config.toml"))
	return os.IsNotExist(err)
}

// RunWizard executes the welcome wizard (5-step interactive setup).
func RunWizard(a adapter.OSAdapter) error {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║     Welcome to cue v2.0     ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Let's configure a few things. This takes about 60 seconds.")
	fmt.Println()

	// Step 1: OS detection confirm
	fmt.Printf("  [1/5] Detected: %s (%s) — is this right? ", a.OSDistro(), a.PackageManagerName())
	confirmed := ui.Confirm("")
	if !confirmed {
		fmt.Println("  No worries — everything will still work. You can edit config.toml later.")
	}
	fmt.Println()

	// Step 2: Shell detection
	shell := detectShell()
	fmt.Printf("  [2/5] Detected shell: %s\n", shell)
	fmt.Println()

	// Step 3: Shell completions
	fmt.Print("  [3/5] Install shell completions? ")
	installCompletions := ui.Confirm("")
	if installCompletions {
		if err := writeCompletions(shell); err != nil {
			ui.PrintWarning(fmt.Sprintf("Completion install: %v", err))
		} else {
			ui.PrintSuccess("Shell completions installed")
		}
	}
	fmt.Println()

	// Step 4: Notifications
	fmt.Print("  [4/5] Enable desktop notifications when queued commands finish? ")
	notifications := ui.Confirm("")
	fmt.Println()

	// Step 5: Default tag
	fmt.Print("  [5/5] Set a default project tag? (leave blank to skip): ")
	tag := ui.ReadInput("")

	// Write config
	cfgDir := config.ConfigDir()
	os.MkdirAll(cfgDir, 0755)

	notifyStr := "false"
	if notifications {
		notifyStr = "true"
	}

	cfgContent := fmt.Sprintf(`[core]
lock_poll_interval_secs = 5
lock_timeout_mins = 30
adaptive_backoff = true
notify_on_completion = %s

[network]
probe_host = "1.1.1.1"
probe_fallback_host = "8.8.8.8"
probe_fallback_port = 53
fail_threshold = 3
recovery_threshold = 1
probe_interval_secs = 10

[history]
max_entries = 50000
default_display_count = 20

[workspace]
github_repo_name = "dev-workspace-backup"
backup_shell_configs = true
backup_vscode = false
backup_history = false

[ui]
color = true
progress_style = "bar"
explain_after_macro = true

[session]
default_tag = "%s"
shell = "%s"
`, notifyStr, tag, shell)

	os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(cfgContent), 0644)

	fmt.Println()
	ui.PrintSuccess("Setup complete!")
	fmt.Println()
	fmt.Println("  Try these commands next:")
	fmt.Println("  ┌──────────────────────────────────────────────────────────────┐")
	fmt.Println("  │  cue doctor          — check your environment      │")
	fmt.Println("  │  cue store            — browse environment stacks  │")
	fmt.Println("  │  cue toolkit list     — available dev tools        │")
	fmt.Println("  │  cue explain --list   — explore macros             │")
	fmt.Println("  │  cue status           — system overview            │")
	fmt.Println("  └──────────────────────────────────────────────────────────────┘")
	fmt.Println()

	return nil
}

func detectShell() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	shell := os.Getenv("SHELL")
	switch {
	case containsStr(shell, "zsh"):
		return "zsh"
	case containsStr(shell, "fish"):
		return "fish"
	default:
		return "bash"
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > len(sub) && s[len(s)-len(sub):] == sub || findSubstring(s, sub))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func writeCompletions(shell string) error {
	dir := filepath.Join(config.ConfigDir(), "completions")
	os.MkdirAll(dir, 0755)

	var cmd *exec.Cmd
	var outFile string

	switch shell {
	case "zsh":
		outFile = filepath.Join(dir, "_cue")
		cmd = exec.Command("cue", "completion", "zsh")
	case "bash":
		outFile = filepath.Join(dir, "cue.bash")
		cmd = exec.Command("cue", "completion", "bash")
	case "fish":
		outFile = filepath.Join(dir, "cue.fish")
		cmd = exec.Command("cue", "completion", "fish")
	case "powershell":
		outFile = filepath.Join(dir, "cue.ps1")
		cmd = exec.Command("cue", "completion", "powershell")
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	return os.WriteFile(outFile, out, 0644)
}
