package job

import (
	"fmt"
	"net"
	"time"

	"github.com/GyaneshSamanta/cue/internal/config"
)

// NetworkEvent represents a network state change.
type NetworkEvent int

const (
	NetworkLost     NetworkEvent = iota
	NetworkRestored NetworkEvent = iota
)

// NetworkMonitor polls network connectivity and emits events.
type NetworkMonitor struct {
	probeHost     string
	fallbackHost  string
	fallbackPort  int
	failThreshold int
	probeInterval time.Duration
	stopCh        chan struct{}
}

// NewNetworkMonitor creates a monitor from config.
func NewNetworkMonitor() *NetworkMonitor {
	cfg := config.DefaultConfig()
	if config.Current != nil {
		cfg = config.Current
	}
	return &NetworkMonitor{
		probeHost:     cfg.Network.ProbeHost,
		fallbackHost:  cfg.Network.ProbeFallbackHost,
		fallbackPort:  cfg.Network.ProbeFallbackPort,
		failThreshold: cfg.Network.FailThreshold,
		probeInterval: time.Duration(cfg.Network.ProbeIntervalSecs) * time.Second,
		stopCh:        make(chan struct{}),
	}
}

// Watch starts monitoring and returns a channel of events.
func (nm *NetworkMonitor) Watch() <-chan NetworkEvent {
	events := make(chan NetworkEvent, 1)
	go func() {
		failures := 0
		wasDown := false
		for {
			select {
			case <-nm.stopCh:
				close(events)
				return
			case <-time.After(nm.probeInterval):
				up := nm.probe()
				if !up {
					failures++
					if failures >= nm.failThreshold && !wasDown {
						wasDown = true
						events <- NetworkLost
					}
				} else {
					if wasDown {
						wasDown = false
						failures = 0
						events <- NetworkRestored
					} else {
						failures = 0
					}
				}
			}
		}
	}()
	return events
}

// Stop ceases monitoring.
func (nm *NetworkMonitor) Stop() {
	select {
	case nm.stopCh <- struct{}{}:
	default:
	}
}

// probe checks connectivity via TCP dial (works without root/admin unlike ICMP).
func (nm *NetworkMonitor) probe() bool {
	// Primary: TCP dial to probe host port 53
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:53", nm.probeHost), 3*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	// Fallback: TCP dial to fallback host
	conn, err = net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", nm.fallbackHost, nm.fallbackPort),
		3*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	return false
}
