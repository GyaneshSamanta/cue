package schedule

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/GyaneshSamanta/cue/internal/config"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Task represents a scheduled task.
type Task struct {
	Name     string
	Command  string
	Interval string // daily, weekly, hourly
}

// Schedule creates an OS-native scheduled task.
func Schedule(task Task) error {
	switch runtime.GOOS {
	case "linux":
		return scheduleLinux(task)
	case "darwin":
		return scheduleDarwin(task)
	case "windows":
		return scheduleWindows(task)
	default:
		return fmt.Errorf("unsupported OS for scheduling: %s", runtime.GOOS)
	}
}

// List shows all cue scheduled tasks.
func List() error {
	ui.PrintHeader("Scheduled Tasks")

	switch runtime.GOOS {
	case "linux":
		// Check systemd user timers
		cmd := exec.Command("systemctl", "--user", "list-timers", "--all")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	case "darwin":
		// Check launchd
		entries, _ := os.ReadDir(filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents"))
		for _, e := range entries {
			if len(e.Name()) > 7 && e.Name()[:7] == "gyanesh" {
				fmt.Printf("  %s\n", e.Name())
			}
		}
	case "windows":
		cmd := exec.Command("schtasks", "/query", "/tn", "cue-*", "/fo", "TABLE")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
	return nil
}

// Remove removes a scheduled task.
func Remove(taskName string) error {
	switch runtime.GOOS {
	case "linux":
		exec.Command("systemctl", "--user", "stop", "gyanesh-"+taskName+".timer").Run()
		exec.Command("systemctl", "--user", "disable", "gyanesh-"+taskName+".timer").Run()
	case "darwin":
		plistPath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "com.gyanesh."+taskName+".plist")
		exec.Command("launchctl", "unload", plistPath).Run()
		os.Remove(plistPath)
	case "windows":
		exec.Command("schtasks", "/delete", "/tn", "cue-"+taskName, "/f").Run()
	}
	ui.PrintSuccess(fmt.Sprintf("Removed scheduled task: %s", taskName))
	return nil
}

func scheduleLinux(task Task) error {
	home, _ := os.UserHomeDir()
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	os.MkdirAll(unitDir, 0755)

	// Service file
	servicePath := filepath.Join(unitDir, "gyanesh-"+task.Name+".service")
	serviceContent := fmt.Sprintf(`[Unit]
Description=cue %s

[Service]
Type=oneshot
ExecStart=%s
`, task.Name, task.Command)
	os.WriteFile(servicePath, []byte(serviceContent), 0644)

	// Timer file
	schedule := "daily"
	if task.Interval == "weekly" {
		schedule = "weekly"
	} else if task.Interval == "hourly" {
		schedule = "*:0/60"
	}

	timerPath := filepath.Join(unitDir, "gyanesh-"+task.Name+".timer")
	timerContent := fmt.Sprintf(`[Unit]
Description=cue %s timer

[Timer]
OnCalendar=%s
Persistent=true

[Install]
WantedBy=timers.target
`, task.Name, schedule)
	os.WriteFile(timerPath, []byte(timerContent), 0644)

	exec.Command("systemctl", "--user", "daemon-reload").Run()
	exec.Command("systemctl", "--user", "enable", "--now", "gyanesh-"+task.Name+".timer").Run()

	ui.PrintSuccess(fmt.Sprintf("Scheduled '%s' (%s) via systemd user timer", task.Name, task.Interval))
	return nil
}

func scheduleDarwin(task Task) error {
	home, _ := os.UserHomeDir()
	plistDir := filepath.Join(home, "Library", "LaunchAgents")
	os.MkdirAll(plistDir, 0755)

	interval := 86400 // daily
	if task.Interval == "weekly" {
		interval = 604800
	} else if task.Interval == "hourly" {
		interval = 3600
	}

	plistPath := filepath.Join(plistDir, "com.gyanesh."+task.Name+".plist")
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.gyanesh.%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>StartInterval</key>
    <integer>%d</integer>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`, task.Name, task.Command, interval)
	os.WriteFile(plistPath, []byte(plistContent), 0644)

	exec.Command("launchctl", "load", plistPath).Run()
	ui.PrintSuccess(fmt.Sprintf("Scheduled '%s' (%s) via launchd", task.Name, task.Interval))
	return nil
}

func scheduleWindows(task Task) error {
	schedMap := map[string]string{
		"daily":  "DAILY",
		"weekly": "WEEKLY",
		"hourly": "HOURLY",
	}
	sched := schedMap[task.Interval]
	if sched == "" {
		sched = "DAILY"
	}

	_ = config.ConfigDir() // ensure import used

	cmd := exec.Command("schtasks", "/create",
		"/tn", "cue-"+task.Name,
		"/tr", task.Command,
		"/sc", sched,
		"/f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Scheduled '%s' (%s) via Task Scheduler", task.Name, task.Interval))
	return nil
}
