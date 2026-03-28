package queue

import (
	"fmt"
	"time"

	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

// Poller implements adaptive lock polling with backoff.
type Poller struct {
	detector     *LockDetector
	baseInterval time.Duration
	maxInterval  time.Duration
	backoffAfter time.Duration
	timeout      time.Duration
}

// WaitAndExecute polls until the lock is released, then runs fn.
func (p *Poller) WaitAndExecute(fn func() error) error {
	start := time.Now()
	interval := p.baseInterval

	for {
		locked, desc, err := p.detector.IsLocked()
		if err != nil {
			return fmt.Errorf("lock check error: %w", err)
		}

		if !locked {
			ui.PrintSuccess("Lock released! Executing now...")
			return fn()
		}

		if time.Since(start) > p.timeout {
			return fmt.Errorf("lock wait timeout after %v: %s", p.timeout, desc)
		}

		// Adaptive backoff after 2 minutes
		if time.Since(start) > p.backoffAfter && interval < p.maxInterval {
			interval += 5 * time.Second
			if interval > p.maxInterval {
				interval = p.maxInterval
			}
		}

		elapsed := time.Since(start).Round(time.Second)
		ui.PrintStatus(fmt.Sprintf("[QUEUED] Waiting %ds... (elapsed: %s)",
			int(interval.Seconds()), elapsed))

		time.Sleep(interval)
	}
}
