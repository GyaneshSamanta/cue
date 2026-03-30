package queue

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/adapter"
)

// LockDetector checks for active package manager locks.
type LockDetector struct {
	adapter adapter.OSAdapter
}

// IsLocked checks all known lock mechanisms for the current OS.
func (d *LockDetector) IsLocked() (bool, string, error) {
	return d.adapter.IsLocked()
}

// CheckWindowsProcesses checks for active installer processes on Windows.
func CheckWindowsProcesses() (bool, string) {
	for _, proc := range []string{"msiexec.exe", "winget.exe", "choco.exe"} {
		out, err := exec.Command("tasklist", "/FI",
			fmt.Sprintf("IMAGENAME eq %s", proc), "/NH").Output()
		if err == nil && strings.Contains(strings.ToLower(string(out)), strings.ToLower(proc)) {
			return true, proc
		}
	}
	return false, ""
}
