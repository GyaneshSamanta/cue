//go:build !windows

package job

import (
	"os"
	"syscall"
)

// SignalTSTP returns SIGTSTP for Unix systems (Ctrl+Z).
func SignalTSTP() os.Signal {
	return syscall.SIGTSTP
}
