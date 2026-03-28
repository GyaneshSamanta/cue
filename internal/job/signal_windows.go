//go:build windows

package job

import (
	"os"
)

// SignalTSTP returns os.Interrupt on Windows (no SIGTSTP equivalent).
func SignalTSTP() os.Signal {
	return os.Interrupt
}
