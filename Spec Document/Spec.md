# Technical Specification Document
## `cue` — Cross-Platform CLI Developer Utility
**Version:** 1.0-draft  
**Status:** For Engineering Review  
**Author:** Gyanesh Samanta 
**Last Updated:** 2026-03-28

---

## Table of Contents

1. [High-Level System Architecture](#1-high-level-system-architecture)
2. [Module Breakdown](#2-module-breakdown)
3. [Local Data Schema](#3-local-data-schema)
4. [Queue / Lock-Polling Implementation](#4-queue--lock-polling-implementation)
5. [Pause / Resume Implementation](#5-pause--resume-implementation)
6. [Semantic Macro Engine](#6-semantic-macro-engine)
7. [Environment Store Engine](#7-environment-store-engine)
8. [History Engine](#8-history-engine)
9. [Workspace Backup & Sync Engine](#9-workspace-backup--sync-engine)
10. [Cross-Platform Abstractions](#10-cross-platform-abstractions)
11. [Configuration System](#11-configuration-system)
12. [Edge Cases & Error Handling](#12-edge-cases--error-handling)
13. [Build, Distribution & Packaging](#13-build-distribution--packaging)
14. [Testing Strategy](#14-testing-strategy)
15. [Security Considerations](#15-security-considerations)

---

## 1. High-Level System Architecture

### 1.1 Conceptual Data Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         USER TERMINAL                                    │
│   $ cue store install mern                                      │
└────────────────────────────────┬─────────────────────────────────────────┘
                                 │ argv
                                 ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  ENTRY POINT  (main.go / main.rs / antigravity entrypoint)               │
│  • Parse raw argv                                                         │
│  • Load ~/.cue/config.toml                                      │
│  • Hydrate session state from ~/.cue/session.json               │
└──────────────────────────────────┬───────────────────────────────────────┘
                                   │
                                   ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  COMMAND ROUTER / PARSER                                                  │
│  • Matches sub-commands: store | macro | history | workspace | tag | ... │
│  • Validates flags and positional arguments                               │
│  • Returns a typed Command struct                                         │
└───────┬──────────┬───────────┬───────────┬──────────┬────────────────────┘
        │          │           │           │          │
        ▼          ▼           ▼           ▼          ▼
   [Store]    [Macro]    [History]   [Workspace]  [Queue/
   Engine     Engine      Engine      Engine      Job Ctrl]
        │          │           │           │          │
        └──────────┴───────────┴─────┬─────┘          │
                                     ▼                 │
                          ┌──────────────────┐         │
                          │  OS ADAPTER LAYER │◄────────┘
                          │  • detect_os()    │
                          │  • pkg_manager()  │
                          │  • lock_path()    │
                          │  • elevate()      │
                          └──────────┬───────┘
                                     │
                    ┌────────────────┼─────────────────┐
                    ▼                ▼                  ▼
             [Linux Adapter]  [macOS Adapter]  [Windows Adapter]
               apt/dnf/pacman    brew/port       winget/choco
                    │                │                  │
                    └────────────────┼──────────────────┘
                                     ▼
                          ┌──────────────────────┐
                          │   CHILD PROCESS MGR   │
                          │  • spawn()            │
                          │  • suspend()          │
                          │  • resume()           │
                          │  • kill()             │
                          └──────────┬────────────┘
                                     │
                          ┌──────────▼────────────┐
                          │   LOCAL DATA LAYER    │
                          │  SQLite history.db    │
                          │  TOML config files    │
                          │  JSON job/queue state │
                          └───────────────────────┘
```

### 1.2 Technology Stack

| Concern | Choice | Rationale |
|---------|--------|-----------|
| **Primary Language** | Go 1.22+ | Compiles to a single static binary; excellent cross-platform stdlib; goroutines for concurrency; fast startup (<10ms overhead) |
| **Build Framework** | Antigravity (as specified) | Project requirement; handles task orchestration and cross-compilation |
| **CLI Framework** | `cobra` (Go) | Industry-standard; handles sub-commands, flags, help generation |
| **Config Format** | TOML | Human-readable; `BurntSushi/toml` parser in Go |
| **Local Database** | SQLite via `mattn/go-sqlite3` | Zero-server; single file; full SQL query support |
| **Job State** | JSON files | Simple; human-readable; no DB overhead for low-cardinality data |
| **Crypto / Keychain** | OS keychain via `zalando/go-keyring` | Cross-platform secret storage without rolling custom crypto |
| **GitHub API** | `google/go-github` client | PAT-authenticated REST calls for backup/restore |
| **Network Probing** | `go-ping` (ICMP) + `net.DialTimeout` fallback | Both ICMP and TCP probe strategies |

### 1.3 Directory Layout (Source)

```
cue/
├── main.go                     # Entrypoint
├── cmd/                        # Cobra command definitions
│   ├── root.go
│   ├── store.go
│   ├── macro.go
│   ├── history.go
│   ├── workspace.go
│   ├── tag.go
│   ├── queue.go
│   └── resume.go
├── internal/
│   ├── adapter/                # OS Adapter Layer
│   │   ├── adapter.go          # Interface definition
│   │   ├── linux.go
│   │   ├── macos.go
│   │   └── windows.go
│   ├── queue/                  # Queue & Lock Polling
│   │   ├── manager.go
│   │   ├── lock_detector.go
│   │   └── poller.go
│   ├── job/                    # Pause/Resume
│   │   ├── controller.go
│   │   ├── network_monitor.go
│   │   └── signal_handler.go
│   ├── macro/                  # Macro Engine
│   │   ├── registry.go
│   │   ├── executor.go
│   │   └── builtins/           # One .go file per macro category
│   │       ├── git.go
│   │       ├── docker.go
│   │       ├── filesystem.go
│   │       ├── network.go
│   │       └── system.go
│   ├── store/                  # Environment Stores
│   │   ├── engine.go
│   │   ├── verifier.go
│   │   └── stacks/
│   │       ├── data_science.go
│   │       ├── frontend.go
│   │       ├── lamp.go
│   │       ├── mern.go
│   │       ├── backend.go
│   │       ├── ai_dev.go
│   │       └── claude_setup.go
│   ├── history/                # History Engine
│   │   ├── db.go
│   │   ├── writer.go
│   │   └── query.go
│   ├── workspace/              # Backup & Sync
│   │   ├── capture.go
│   │   ├── backup.go
│   │   └── restore.go
│   ├── config/                 # Config Loading
│   │   └── loader.go
│   └── ui/                     # Terminal Output Helpers
│       ├── printer.go          # Colour, tables, progress bars
│       └── prompt.go           # Y/N, selection prompts
├── data/
│   └── macros.toml             # Built-in macro definitions
├── Antigravityfile             # Build task definitions
├── go.mod
└── go.sum
```

---

## 2. Module Breakdown

### 2.1 OS Adapter Interface

All platform-specific behaviour is hidden behind a single `Adapter` interface. The correct implementation is selected at runtime (not compile time), keeping the binary universal.

```go
// internal/adapter/adapter.go

type OSAdapter interface {
    // Package management
    PackageManagerName() string                 // "apt", "brew", "winget", etc.
    InstallPackage(pkg string, args []string) error
    UninstallPackage(pkg string) error
    IsPackageInstalled(pkg string) bool
    PackageVersion(pkg string) (string, error)

    // Lock detection
    LockPaths() []string                        // Paths to check for lock files
    IsLocked() (bool, string, error)            // locked?, lock holder description, error

    // Process management
    SuspendProcess(pid int) error
    ResumeProcess(pid int) error
    KillProcess(pid int) error
    RunElevated(cmd string, args []string) error // sudo on Unix; runas on Windows

    // System info
    HomeDir() string
    ConfigDir() string                          // ~/.cue on Unix, %APPDATA%\cue on Windows
    OSName() string                             // "linux", "darwin", "windows"
    OSDistro() string                           // "ubuntu", "fedora", "arch", "" on non-Linux
    HasGPU() bool                               // Checks for nvidia-smi or equivalent
}
```

**Runtime Detection Logic:**

```go
// internal/adapter/adapter.go

func Detect() OSAdapter {
    switch runtime.GOOS {
    case "linux":
        return newLinuxAdapter()    // further detects distro via /etc/os-release
    case "darwin":
        return newMacOSAdapter()
    case "windows":
        return newWindowsAdapter()
    default:
        log.Fatalf("Unsupported OS: %s", runtime.GOOS)
        return nil
    }
}
```

**Linux Adapter — Distro Detection:**

```go
// internal/adapter/linux.go

func newLinuxAdapter() *LinuxAdapter {
    distro := detectLinuxDistro()  // reads /etc/os-release
    switch distro {
    case "ubuntu", "debian", "linuxmint", "pop":
        return &LinuxAdapter{pkgMgr: "apt", lockPaths: aptLockPaths}
    case "fedora", "rhel", "centos", "rocky", "alma":
        return &LinuxAdapter{pkgMgr: "dnf", lockPaths: dnfLockPaths}
    case "arch", "manjaro", "endeavouros":
        return &LinuxAdapter{pkgMgr: "pacman", lockPaths: pacmanLockPaths}
    default:
        // Fallback: try to detect by binary presence
        return detectByBinaryPresence()
    }
}

var aptLockPaths = []string{
    "/var/lib/dpkg/lock-frontend",
    "/var/lib/dpkg/lock",
    "/var/lib/apt/lists/lock",
    "/var/cache/apt/archives/lock",
}

var dnfLockPaths = []string{
    "/var/lib/rpm/.rpm.lock",
    "/var/run/dnf.pid",
}

var pacmanLockPaths = []string{
    "/var/lib/pacman/db.lck",
}
```

---

## 3. Local Data Schema

### 3.1 Directory Structure (Runtime)

```
~/.cue/
├── config.toml           # User configuration
├── history.db            # SQLite command history
├── macros.toml           # User-defined custom macros
├── session.json          # Active session state (current tag, etc.)
├── queue.json            # Pending queued commands
├── jobs.json             # Active/paused managed jobs
└── workspace/
    ├── manifest.json     # Last-captured store manifest
    └── backup.log        # Backup/restore operation log
```

### 3.2 config.toml

```toml
[core]
lock_poll_interval_secs = 5          # How often to poll for lock release
lock_timeout_mins = 30               # Maximum wait time before aborting queue
adaptive_backoff = true              # Enable adaptive polling backoff
notify_on_completion = true          # Desktop notification when queued cmd completes

[network]
probe_host = "1.1.1.1"              # Primary ICMP probe target
probe_fallback_host = "8.8.8.8"    # TCP fallback probe target
probe_fallback_port = 53
fail_threshold = 3                   # Consecutive failures before pausing
recovery_threshold = 1               # Successes before resuming
probe_interval_secs = 10

[history]
max_entries = 50000                  # Rotate after this many rows
default_display_count = 20
export_dir = "~/.cue/exports"

[workspace]
github_repo_name = "dev-workspace-backup"
backup_shell_configs = true
backup_vscode = false                # Opt-in; can be large
backup_history = false               # Opt-in; can be large
auto_backup_interval_hours = 0       # 0 = manual only; >0 = scheduled

[ui]
color = true
progress_style = "bar"               # "bar" | "spinner" | "none"
explain_after_macro = true           # Print explanation block after every macro
```

### 3.3 history.db — SQLite Schema

```sql
-- Main command history table
CREATE TABLE IF NOT EXISTS history (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp    TEXT    NOT NULL,           -- ISO 8601: 2026-03-28T14:32:00Z
    command      TEXT    NOT NULL,           -- Full command string as typed
    subcommand   TEXT,                       -- Parsed sub-command (store, macro, etc.)
    args         TEXT,                       -- JSON array of args
    exit_code    INTEGER NOT NULL DEFAULT 0,
    duration_ms  INTEGER NOT NULL DEFAULT 0,
    project_tag  TEXT,                       -- Value from session.json at time of execution
    stack        TEXT,                       -- Active store stack at time of execution
    os           TEXT    NOT NULL,           -- "linux/ubuntu", "darwin", "windows"
    notes        TEXT                        -- Optional user annotation
);

-- Full-text search index for fast substring queries
CREATE VIRTUAL TABLE IF NOT EXISTS history_fts
    USING fts5(command, args, notes, content='history', content_rowid='id');

-- Triggers to keep FTS index in sync
CREATE TRIGGER IF NOT EXISTS history_ai AFTER INSERT ON history BEGIN
    INSERT INTO history_fts(rowid, command, args, notes)
    VALUES (new.id, new.command, new.args, new.notes);
END;
CREATE TRIGGER IF NOT EXISTS history_ad AFTER DELETE ON history BEGIN
    INSERT INTO history_fts(history_fts, rowid, command, args, notes)
    VALUES ('delete', old.id, old.command, old.args, old.notes);
END;

-- Macro execution analytics
CREATE TABLE IF NOT EXISTS macro_stats (
    macro_name   TEXT    NOT NULL,
    run_count    INTEGER NOT NULL DEFAULT 0,
    last_run     TEXT,
    PRIMARY KEY (macro_name)
);

-- Store installation log
CREATE TABLE IF NOT EXISTS store_installs (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp    TEXT    NOT NULL,
    stack        TEXT    NOT NULL,           -- "data-science", "mern", etc.
    action       TEXT    NOT NULL,           -- "install" | "remove" | "verify"
    status       TEXT    NOT NULL,           -- "success" | "partial" | "failed"
    components   TEXT,                       -- JSON array of {name, version, status}
    duration_ms  INTEGER
);

-- Standard indexes
CREATE INDEX IF NOT EXISTS idx_history_timestamp   ON history(timestamp);
CREATE INDEX IF NOT EXISTS idx_history_project_tag ON history(project_tag);
CREATE INDEX IF NOT EXISTS idx_history_exit_code   ON history(exit_code);
CREATE INDEX IF NOT EXISTS idx_history_subcommand  ON history(subcommand);
```

### 3.4 queue.json

```json
{
  "version": 1,
  "entries": [
    {
      "id": "q-20260328-001",
      "command": "install",
      "args": ["vim"],
      "package_manager": "apt",
      "status": "waiting",
      "created_at": "2026-03-28T14:30:00Z",
      "started_at": null,
      "completed_at": null,
      "exit_code": null,
      "wait_seconds": 47
    }
  ]
}
```

### 3.5 jobs.json

```json
{
  "version": 1,
  "jobs": [
    {
      "id": "job-20260328-001",
      "pid": 18342,
      "command": "store install ai-dev",
      "status": "paused",
      "pause_reason": "network_loss",
      "created_at": "2026-03-28T13:00:00Z",
      "paused_at": "2026-03-28T13:22:17Z",
      "resumed_at": null,
      "project_tag": "ml-project",
      "progress_hint": "Installing PyTorch..."
    }
  ]
}
```

### 3.6 macros.toml (Custom Macros)

```toml
[[macro]]
name = "nuke-venv"
command = "rm -rf .venv && python -m venv .venv && source .venv/bin/activate"
category = "python"
description = "Nuke and recreate the virtual environment from scratch"
explanation = """
Your old virtual environment was deleted entirely.
A clean new one was created in .venv/ and activated.
All your old packages are gone — reinstall from requirements.txt if needed.
This is useful when dependencies get into a confused, broken state.
"""
tags = ["python", "cleanup"]
added_at = "2026-03-28T10:00:00Z"
```

---

## 4. Queue / Lock-Polling Implementation

### 4.1 Flow Diagram

```
cue install vim
         │
         ▼
   LockDetector.Check()
         │
    ┌────┴────┐
    │ Locked? │
    └────┬────┘
    No   │  Yes
    │    ▼
    │  Print: "[QUEUED] Package manager busy..."
    │  Write entry to queue.json (status: "waiting")
    │  Start Poller goroutine (non-blocking)
    │         │
    │         ▼
    │   ┌─────────────┐
    │   │ Poll loop   │◄──────────────────────┐
    │   │  sleep(N)   │                       │
    │   │  Check lock │                       │
    │   └──────┬──────┘                       │
    │          │ Still locked?          ┌──────┴──────┐
    │          │ Yes → Update timer UI  │ Adaptive    │
    │          │ and continue           │ Backoff:    │
    │          │                        │ 5s→10s→15s │
    │          │ No ──────────────►──► └─────────────┘
    │          ▼
    └──► Execute command (update queue.json status: "running")
               │
               ▼
         Record result (status: "done" | "failed")
         Write to history.db
         Print: "[DONE] vim installed successfully."
         Optional: Desktop notification
```

### 4.2 Lock Detector Implementation

```go
// internal/queue/lock_detector.go

type LockDetector struct {
    adapter adapter.OSAdapter
}

// IsLocked checks all known lock paths AND active processes.
// Returns (locked bool, description string, error).
func (d *LockDetector) IsLocked() (bool, string, error) {
    // Method 1: File-based lock check
    for _, path := range d.adapter.LockPaths() {
        if fileIsLocked(path) {
            return true, fmt.Sprintf("lock file held: %s", path), nil
        }
    }

    // Method 2: Process-based lock check (belt-and-suspenders)
    if d.adapter.OSName() == "windows" {
        return d.checkWindowsProcessLock()
    }

    return false, "", nil
}

// fileIsLocked checks existence AND attempts an exclusive flock.
// Existence alone is insufficient — stale lockfiles can remain after crashes.
func fileIsLocked(path string) bool {
    f, err := os.Open(path)
    if err != nil {
        return false // File doesn't exist — no lock
    }
    defer f.Close()

    // Try non-blocking exclusive lock
    err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
    if err == syscall.EWOULDBLOCK {
        return true // Actively locked by another process
    }
    // If we got the lock, release it immediately — it was stale
    syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
    return false
}
```

**Windows Process Lock Check:**

```go
// internal/queue/lock_detector.go (Windows branch)

func (d *LockDetector) checkWindowsProcessLock() (bool, string, error) {
    // WMI query for active msiexec or winget processes
    processes := []string{"msiexec.exe", "winget.exe", "choco.exe"}
    for _, proc := range processes {
        if isProcessRunning(proc) {
            return true, fmt.Sprintf("active installer process: %s", proc), nil
        }
    }
    return false, "", nil
}

func isProcessRunning(name string) bool {
    out, err := exec.Command("tasklist", "/FI",
        fmt.Sprintf("IMAGENAME eq %s", name), "/NH").Output()
    if err != nil {
        return false
    }
    return strings.Contains(strings.ToLower(string(out)), strings.ToLower(name))
}
```

### 4.3 Adaptive Polling Backoff

```go
// internal/queue/poller.go

type Poller struct {
    detector      *LockDetector
    baseInterval  time.Duration  // From config: default 5s
    maxInterval   time.Duration  // Caps at 15s
    backoffAfter  time.Duration  // Start backing off after 2 minutes
    timeout       time.Duration  // From config: default 30 minutes
}

func (p *Poller) WaitAndExecute(ctx context.Context, fn func() error) error {
    start := time.Now()
    interval := p.baseInterval

    for {
        select {
        case <-ctx.Done():
            return ctx.Err() // User cancelled via Ctrl+C

        default:
            locked, desc, err := p.detector.IsLocked()
            if err != nil {
                return fmt.Errorf("lock check error: %w", err)
            }

            if !locked {
                return fn() // Execute the queued command
            }

            // Timeout check
            if time.Since(start) > p.timeout {
                return fmt.Errorf("lock wait timeout after %v: %s", p.timeout, desc)
            }

            // Adaptive backoff after 2 minutes
            if time.Since(start) > p.backoffAfter && interval < p.maxInterval {
                interval = min(interval+5*time.Second, p.maxInterval)
            }

            // Update terminal status line in place
            ui.PrintStatus(fmt.Sprintf("[QUEUED] Waiting %ds... (held for %s)",
                int(interval.Seconds()), formatDuration(time.Since(start))))

            time.Sleep(interval)
        }
    }
}
```

---

## 5. Pause / Resume Implementation

### 5.1 Job Controller

```go
// internal/job/controller.go

type JobController struct {
    adapter    adapter.OSAdapter
    netMonitor *NetworkMonitor
    jobFile    string              // Path to jobs.json
}

// ManagedRun spawns cmd as a child process and wraps it with:
//   - Network-aware auto-pause/resume
//   - SIGTERM / Ctrl+C signal handling
//   - State persistence to jobs.json
func (jc *JobController) ManagedRun(cmd *exec.Cmd, jobID string) error {
    if err := cmd.Start(); err != nil {
        return err
    }

    pid := cmd.Process.Pid
    job := &Job{
        ID:        jobID,
        PID:       pid,
        Command:   strings.Join(cmd.Args, " "),
        Status:    "running",
        CreatedAt: time.Now(),
    }
    jc.persist(job)

    // Start network monitor
    netEvents := jc.netMonitor.Watch(context.Background())

    // Signal listener for manual pause (Ctrl+Z)
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGTSTP)

    // Wait for process in a goroutine
    done := make(chan error, 1)
    go func() { done <- cmd.Wait() }()

    for {
        select {
        case err := <-done:
            job.Status = "completed"
            if err != nil {
                job.Status = "failed"
            }
            jc.persist(job)
            return err

        case event := <-netEvents:
            switch event {
            case NetworkLost:
                ui.PrintWarning("[PAUSED] Network lost. Installation paused...")
                jc.adapter.SuspendProcess(pid)
                job.Status = "paused"
                job.PauseReason = "network_loss"
                job.PausedAt = timePtr(time.Now())
                jc.persist(job)

            case NetworkRestored:
                ui.PrintInfo("[RESUMED] Network restored. Continuing...")
                jc.adapter.ResumeProcess(pid)
                job.Status = "running"
                job.ResumedAt = timePtr(time.Now())
                jc.persist(job)
            }

        case <-sigCh:
            ui.PrintInfo("[PAUSED] Manual pause. Run 'cue resume' to continue.")
            jc.adapter.SuspendProcess(pid)
            job.Status = "paused"
            job.PauseReason = "manual"
            job.PausedAt = timePtr(time.Now())
            jc.persist(job)
        }
    }
}
```

### 5.2 OS-Specific Process Control

| Platform | Suspend | Resume |
|----------|---------|--------|
| Linux | `syscall.Kill(pid, syscall.SIGSTOP)` | `syscall.Kill(pid, syscall.SIGCONT)` |
| macOS | `syscall.Kill(pid, syscall.SIGSTOP)` | `syscall.Kill(pid, syscall.SIGCONT)` |
| Windows | `NtSuspendProcess` via `SuspendThread` on all threads | `ResumeThread` on all suspended threads |

**Windows Suspend/Resume (CGo or syscall wrapper):**

```go
// internal/adapter/windows.go

import "golang.org/x/sys/windows"

func (a *WindowsAdapter) SuspendProcess(pid int) error {
    handle, err := windows.OpenProcess(
        windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
    if err != nil {
        return err
    }
    defer windows.CloseHandle(handle)
    // NtSuspendProcess is an undocumented but stable NT API
    ntdll := windows.NewLazySystemDLL("ntdll.dll")
    proc := ntdll.NewProc("NtSuspendProcess")
    r, _, err := proc.Call(uintptr(handle))
    if r != 0 {
        return fmt.Errorf("NtSuspendProcess failed: %w", err)
    }
    return nil
}
```

### 5.3 Network Monitor

```go
// internal/job/network_monitor.go

type NetworkEvent int
const (
    NetworkLost     NetworkEvent = iota
    NetworkRestored NetworkEvent = iota
)

type NetworkMonitor struct {
    probeHost     string         // "1.1.1.1"
    fallbackHost  string         // "8.8.8.8"
    fallbackPort  int            // 53
    failThreshold int            // 3
    probeInterval time.Duration  // 10s
}

func (nm *NetworkMonitor) Watch(ctx context.Context) <-chan NetworkEvent {
    events := make(chan NetworkEvent, 1)
    go func() {
        failures := 0
        wasDown := false
        for {
            select {
            case <-ctx.Done():
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

func (nm *NetworkMonitor) probe() bool {
    // Method 1: ICMP ping (may require root/admin on some systems)
    pinger, err := ping.NewPinger(nm.probeHost)
    if err == nil {
        pinger.Count = 1
        pinger.Timeout = 3 * time.Second
        pinger.SetPrivileged(true)
        err = pinger.Run()
        if err == nil && pinger.Statistics().PacketsRecv > 0 {
            return true
        }
    }

    // Method 2: TCP dial fallback
    conn, err := net.DialTimeout("tcp",
        fmt.Sprintf("%s:%d", nm.fallbackHost, nm.fallbackPort),
        3*time.Second)
    if err == nil {
        conn.Close()
        return true
    }
    return false
}
```

---

## 6. Semantic Macro Engine

### 6.1 Macro Registry

Macros are defined in Go as structs, compiled into the binary. User-defined macros in `~/.cue/macros.toml` are loaded at startup and merged into the registry.

```go
// internal/macro/registry.go

type Macro struct {
    Name        string
    Category    string
    Description string     // One-line summary for --list
    Commands    []Step     // Ordered steps to execute
    Explanation string     // Multi-line plain-English explanation
    Flags       []Flag     // Supported flags (e.g., --hard)
    Dangerous   bool       // Triggers Y/N confirmation before execution
    BuiltIn     bool       // true for compiled-in; false for user-defined
}

type Step struct {
    OS       string   // "all" | "linux" | "darwin" | "windows"
    Command  string   // Shell command template; supports {flags} placeholder
    Args     []string
}

type Flag struct {
    Name        string
    Description string
    Default     string
}

var Registry = map[string]*Macro{}

func init() {
    // Register all built-in macros
    for _, m := range builtins.All() {
        Registry[m.Name] = m
    }
}

// LoadUserMacros merges macros.toml into the registry at startup
func LoadUserMacros(path string) error {
    var file struct {
        Macro []tomlMacro `toml:"macro"`
    }
    if _, err := toml.DecodeFile(path, &file); err != nil {
        if os.IsNotExist(err) { return nil }
        return err
    }
    for _, tm := range file.Macro {
        Registry[tm.Name] = tm.toMacro()
    }
    return nil
}
```

### 6.2 Built-in Macro Definition Example

```go
// internal/macro/builtins/git.go

func init() {
    register(&macro.Macro{
        Name:        "git-undo",
        Category:    "git",
        Description: "Safely undo the last commit, keeping changes staged",
        Flags: []macro.Flag{
            {Name: "hard", Description: "Discard staged changes too (DESTRUCTIVE)", Default: "false"},
        },
        Dangerous: false, // --hard variant overrides this at execution time
        Commands: []macro.Step{
            {OS: "all", Command: "git reset --soft HEAD~1"},
        },
        Explanation: `
✔ Done. Here's what happened:
─────────────────────────────────────────────────────────────────
Your last commit was "undone," but your file changes are SAFE
and still staged. The commit message is gone, but your work is not.

You can re-commit when ready with:  git commit -m "your message"

This rewrites local history only. If you had already pushed this
commit, you will need to force-push (use 'cue git-oops-push').
─────────────────────────────────────────────────────────────────`,
        BuiltIn: true,
    })

    register(&macro.Macro{
        Name:        "git-branch-clean",
        Category:    "git",
        Description: "Delete all local branches that have been merged into main/master",
        Dangerous:   true,
        Commands: []macro.Step{
            {OS: "linux",   Command: `git branch --merged | grep -v '\\*\\|main\\|master\\|develop' | xargs -r git branch -d`},
            {OS: "darwin",  Command: `git branch --merged | grep -v '\\*\\|main\\|master\\|develop' | xargs git branch -d`},
            {OS: "windows", Command: `FOR /F "tokens=*" %b IN ('git branch --merged') DO git branch -d %b`},
        },
        Explanation: `
✔ Done. Here's what happened:
─────────────────────────────────────────────────────────────────
All local branches that had already been merged into your main/
master/develop branch were deleted.

These branches still exist on your remote (GitHub/GitLab) if
you had pushed them. Only your local copies were removed.

To clean remote branches too, use: git remote prune origin
─────────────────────────────────────────────────────────────────`,
        BuiltIn: true,
    })
}
```

### 6.3 Macro Executor

```go
// internal/macro/executor.go

func Execute(name string, flags map[string]string, adapter adapter.OSAdapter) error {
    m, ok := Registry[name]
    if !ok {
        return fmt.Errorf("unknown macro: %s. Run 'cue explain --list' to see all.", name)
    }

    // Dangerous action gate
    isDangerous := m.Dangerous || (name == "git-undo" && flags["hard"] == "true")
    if isDangerous {
        ui.PrintWarning(fmt.Sprintf("⚠ WARNING: '%s' is a destructive operation.", name))
        if !ui.Confirm("Are you sure you want to continue? [y/N] ") {
            ui.PrintInfo("Aborted. No changes were made.")
            return nil
        }
    }

    // Select OS-appropriate steps
    os := adapter.OSName()
    for _, step := range m.Steps(os, flags) {
        ui.PrintDim(fmt.Sprintf("$ %s", step.Command))
        if err := runShellCommand(step.Command, adapter); err != nil {
            return fmt.Errorf("macro step failed: %w", err)
        }
    }

    // Print explanation if enabled
    if config.Current.UI.ExplainAfterMacro {
        ui.PrintExplanation(m.Explanation)
    }

    return nil
}
```

---

## 7. Environment Store Engine

### 7.1 Store Definition Interface

```go
// internal/store/engine.go

type Stack interface {
    Name() string
    Description() string
    EstimatedSizeMB() int
    EstimatedMinutes(os string) int
    Components() []Component
    PostInstallSteps(adapter adapter.OSAdapter) []Step
    VerificationChecks() []Check
}

type Component struct {
    Name           string
    Version        string      // "latest", "lts", "8.x", specific semver
    OS             []string    // ["linux", "darwin", "windows"] or subset
    InstallMethod  InstallMethod
    Optional       bool
    OptionalPrompt string      // Message shown when asking user about optional component
    DependsOn      []string    // Other component names that must be installed first
}

type InstallMethod struct {
    Linux   []string    // apt/dnf/pacman package names or script URL
    Darwin  []string    // brew formulae or cask names
    Windows []string    // winget package IDs or choco packages
    Script  string      // Fallback: URL to install script (curl | bash pattern)
}

type Check struct {
    Name    string
    Command string      // e.g., "python3 --version"
    Pattern string      // Regex to validate output; empty = just check exit code
}
```

### 7.2 Installation Engine

```go
// internal/store/engine.go

func Install(stackName string, adapter adapter.OSAdapter, opts InstallOpts) error {
    stack, err := GetStack(stackName)
    if err != nil {
        return err
    }

    components := stack.Components()

    // Prompt for optional components
    for i, c := range components {
        if c.Optional {
            if ui.Confirm(fmt.Sprintf("Install %s? %s [y/N] ", c.Name, c.OptionalPrompt)) {
                components[i].Optional = false // Mark as selected
            }
        }
    }

    // Topological sort by DependsOn
    ordered := topoSort(components)

    progress := ui.NewProgressBar(len(ordered))
    results := make([]ComponentResult, 0, len(ordered))

    for _, comp := range ordered {
        progress.Update(fmt.Sprintf("Installing %s...", comp.Name))

        // Route to OS-specific install method
        method := comp.InstallMethod.ForOS(adapter.OSName())
        var err error

        if len(method.Packages) > 0 {
            err = adapter.InstallPackages(method.Packages)
        } else if method.Script != "" {
            err = runInstallScript(method.Script, adapter)
        }

        results = append(results, ComponentResult{
            Name:    comp.Name,
            Status:  resultStatus(err),
            Version: probeVersion(comp, adapter),
            Error:   err,
        })

        // Write to store_installs table regardless of individual component result
        history.RecordStoreComponent(stackName, comp.Name, resultStatus(err))
    }

    // Post-install steps (e.g., start services, configure PATH)
    for _, step := range stack.PostInstallSteps(adapter) {
        runStep(step, adapter)
    }

    // Print results table
    ui.PrintStoreResults(results)

    // Run verification
    if opts.Verify {
        return Verify(stackName, adapter)
    }
    return nil
}
```

### 7.3 Cross-Platform Install Routing — MERN Example

```go
// internal/store/stacks/mern.go

func (s *MERNStack) Components() []Component {
    return []Component{
        {
            Name:    "Node.js LTS",
            Version: "lts",
            OS:      []string{"linux", "darwin", "windows"},
            InstallMethod: InstallMethod{
                Script: "https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh",
                // Post-script: nvm install --lts
                Darwin:  []string{"nvm"},       // brew install nvm
                Windows: []string{"CoreyButler.NVMforWindows"}, // winget
            },
        },
        {
            Name:    "MongoDB Community Server",
            Version: "7.x",
            OS:      []string{"linux", "darwin", "windows"},
            InstallMethod: InstallMethod{
                Linux: []string{"mongodb-org"},  // Requires adding MongoDB APT/DNF repo first
                Darwin: []string{"mongodb/brew/mongodb-community"},
                Windows: []string{"MongoDB.Server"},  // winget
            },
        },
        {
            Name:           "MongoDB Compass",
            Version:        "latest",
            Optional:       true,
            OptionalPrompt: "(graphical MongoDB GUI, ~200MB)",
            InstallMethod: InstallMethod{
                Linux:   []string{"mongodb-compass"},
                Darwin:  []string{"mongodb-compass"},  // brew --cask
                Windows: []string{"MongoDB.Compass"},
            },
        },
        {
            Name:    "PM2",
            Version: "latest",
            InstallMethod: InstallMethod{
                Script: "npm install -g pm2",
            },
        },
        // ... additional components
    }
}
```

---

## 8. History Engine

### 8.1 Write Path

Every command issued through `cue` passes through a **middleware wrapper** that records a history entry after execution.

```go
// internal/history/writer.go

type Writer struct {
    db      *sql.DB
    session *SessionState  // Loaded from session.json
}

func (w *Writer) Record(cmd string, args []string, exitCode int,
    durationMs int64, subCmd string) error {

    _, err := w.db.Exec(`
        INSERT INTO history
            (timestamp, command, subcommand, args, exit_code, duration_ms,
             project_tag, stack, os)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        time.Now().UTC().Format(time.RFC3339),
        cmd,
        subCmd,
        mustMarshal(args),
        exitCode,
        durationMs,
        w.session.ActiveTag,
        w.session.ActiveStack,
        detectOSLabel(),
    )
    return err
}
```

### 8.2 Query Interface

```go
// internal/history/query.go

type QueryOptions struct {
    Tag      string
    Search   string
    Since    time.Time
    FailOnly bool
    Limit    int    // Default 20
    Offset   int
}

func Query(db *sql.DB, opts QueryOptions) ([]HistoryEntry, error) {
    query := `SELECT id, timestamp, command, exit_code, duration_ms,
                     project_tag, stack FROM history`
    conditions := []string{}
    args := []interface{}{}

    if opts.Tag != "" {
        conditions = append(conditions, "project_tag = ?")
        args = append(args, opts.Tag)
    }
    if opts.FailOnly {
        conditions = append(conditions, "exit_code != 0")
    }
    if !opts.Since.IsZero() {
        conditions = append(conditions, "timestamp >= ?")
        args = append(args, opts.Since.UTC().Format(time.RFC3339))
    }
    if opts.Search != "" {
        // Use FTS for search
        return ftsQuery(db, opts)
    }

    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }
    query += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
    args = append(args, opts.Limit, opts.Offset)

    // ... execute and scan rows
}

func ftsQuery(db *sql.DB, opts QueryOptions) ([]HistoryEntry, error) {
    return execQuery(db, `
        SELECT h.id, h.timestamp, h.command, h.exit_code, h.duration_ms,
               h.project_tag, h.stack
        FROM history h
        JOIN history_fts f ON h.id = f.rowid
        WHERE history_fts MATCH ?
        ORDER BY rank
        LIMIT ? OFFSET ?`,
        opts.Search, opts.Limit, opts.Offset)
}
```

---

## 9. Workspace Backup & Sync Engine

### 9.1 Capture Phase

```go
// internal/workspace/capture.go

type Manifest struct {
    CapturedAt    string              `json:"captured_at"`
    OS            string              `json:"os"`
    Hostname      string              `json:"hostname"`
    InstalledStores []StoreEntry      `json:"installed_stores"`
    ShellFiles    []string            `json:"shell_files"`
    VSCodeFiles   []string            `json:"vscode_files,omitempty"`
    CustomMacros  bool                `json:"custom_macros"`
    HistoryIncluded bool              `json:"history_included"`
}

func Capture(adapter adapter.OSAdapter, opts CaptureOptions) (*Manifest, string, error) {
    tmpDir, _ := os.MkdirTemp("", "gyanesh-workspace-*")

    manifest := &Manifest{
        CapturedAt: time.Now().UTC().Format(time.RFC3339),
        OS:         adapter.OSName(),
    }
    manifest.Hostname, _ = os.Hostname()

    // 1. Capture shell config files
    shellFiles := shellConfigPaths(adapter)
    for _, sf := range shellFiles {
        if copyIfExists(sf, filepath.Join(tmpDir, "shell", filepath.Base(sf))) == nil {
            manifest.ShellFiles = append(manifest.ShellFiles, filepath.Base(sf))
        }
    }

    // 2. Capture installed store manifest
    manifest.InstalledStores = store.QueryInstalled(adapter)

    // 3. Capture custom macros
    macroDst := filepath.Join(tmpDir, "macros.toml")
    if copyIfExists(macroPath(adapter), macroDst) == nil {
        manifest.CustomMacros = true
    }

    // 4. Sanitise: strip secrets from .gitconfig
    if err := sanitiseGitConfig(filepath.Join(tmpDir, "shell", ".gitconfig")); err != nil {
        // Non-fatal; log warning
    }

    // 5. Write manifest.json
    writeJSON(filepath.Join(tmpDir, "manifest.json"), manifest)

    return manifest, tmpDir, nil
}
```

### 9.2 GitHub Push

```go
// internal/workspace/backup.go

func PushToGitHub(srcDir string, token string, repoName string) (string, error) {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    gh := github.NewClient(oauth2.NewClient(ctx, ts))

    // Create repo if it doesn't exist
    repo, err := ensureRepo(ctx, gh, repoName)
    if err != nil {
        return "", err
    }

    // Use local git to commit and push
    // This avoids the GitHub Content API 25MB file limit for history.db
    runGit(srcDir, "init")
    copyFile(builtinGitignorePath(), filepath.Join(srcDir, ".gitignore"))
    runGit(srcDir, "add", ".")
    runGit(srcDir, "commit", "-m",
        fmt.Sprintf("cue backup: %s", time.Now().UTC().Format(time.RFC3339)))
    runGit(srcDir, "remote", "add", "origin", *repo.CloneURL)
    runGit(srcDir, "push", "--force", "origin", "main")

    return *repo.HTMLURL, nil
}
```

### 9.3 Restore Flow

```go
// internal/workspace/restore.go

func Restore(repoURL string, adapter adapter.OSAdapter) error {
    tmpDir, _ := os.MkdirTemp("", "gyanesh-restore-*")

    // Clone the backup repo
    if err := runCmd("git", "clone", repoURL, tmpDir); err != nil {
        return fmt.Errorf("failed to clone backup: %w", err)
    }

    // Read manifest
    var manifest Manifest
    if err := readJSON(filepath.Join(tmpDir, "manifest.json"), &manifest); err != nil {
        return fmt.Errorf("cannot read backup manifest: %w", err)
    }

    ui.PrintInfo(fmt.Sprintf("Restoring backup from %s (OS: %s)",
        manifest.CapturedAt, manifest.OS))

    // 1. Run store installs
    for _, s := range manifest.InstalledStores {
        ui.PrintStep(fmt.Sprintf("Installing store: %s", s.Name))
        if err := store.Install(s.Name, adapter, store.InstallOpts{Verify: false}); err != nil {
            ui.PrintWarning(fmt.Sprintf("Store '%s' failed: %v. Continuing...", s.Name, err))
        }
    }

    // 2. Copy shell config files (prompt on conflict)
    for _, sf := range manifest.ShellFiles {
        dst := shellFileDest(sf, adapter)
        if fileExists(dst) {
            if !ui.Confirm(fmt.Sprintf("Overwrite existing %s? [y/N] ", dst)) {
                continue
            }
        }
        copyFile(filepath.Join(tmpDir, "shell", sf), dst)
    }

    // 3. Restore custom macros
    if manifest.CustomMacros {
        copyFile(filepath.Join(tmpDir, "macros.toml"), macroPath(adapter))
    }

    ui.PrintSuccess("Restore complete. Open a new terminal for changes to take effect.")
    return nil
}
```

---

## 10. Cross-Platform Abstractions

### 10.1 Privilege Elevation

Privilege elevation (sudo on Unix, UAC on Windows) is handled by the OS Adapter.

```go
// internal/adapter/linux.go
func (a *LinuxAdapter) RunElevated(cmd string, args []string) error {
    // Check if already root
    if os.Getuid() == 0 {
        return exec.Command(cmd, args...).Run()
    }
    // Use sudo; password prompt appears in terminal naturally
    return exec.Command("sudo", append([]string{cmd}, args...)...).Run()
}

// internal/adapter/windows.go
func (a *WindowsAdapter) RunElevated(cmd string, args []string) error {
    // ShellExecute with "runas" verb triggers UAC prompt
    verb := "runas"
    exePath, _ := exec.LookPath(cmd)
    return windows.ShellExecute(0, &verb,
        windows.StringToUTF16Ptr(exePath),
        windows.StringToUTF16Ptr(strings.Join(args, " ")),
        nil, windows.SW_SHOW)
}
```

**Elevation Detection — Pre-flight Check:**

Before attempting any operation requiring elevation, the tool probes whether it has sufficient privileges and provides a clear, actionable message if not:

```
✖ This operation requires administrator/sudo access.
  On Linux/macOS: run with 'sudo cue store install lamp'
  On Windows:     re-open this terminal as Administrator
  
  If you do not have admin rights on this machine, ask your
  IT administrator, or check 'cue store preview lamp'
  to see what would be installed.
```

### 10.2 Shell Profile Detection

```go
// internal/adapter/linux.go (and darwin equivalent)

func shellConfigPaths(homeDir string) []string {
    paths := []string{}
    candidates := []string{
        filepath.Join(homeDir, ".bashrc"),
        filepath.Join(homeDir, ".bash_profile"),
        filepath.Join(homeDir, ".zshrc"),
        filepath.Join(homeDir, ".profile"),
        filepath.Join(homeDir, ".config", "fish", "config.fish"),
    }
    for _, p := range candidates {
        if _, err := os.Stat(p); err == nil {
            paths = append(paths, p)
        }
    }
    return paths
}
```

### 10.3 PATH Manipulation — Cross-Platform

```go
// Macro: path-add

func PathAdd(dir string, adapter adapter.OSAdapter) error {
    switch adapter.OSName() {
    case "linux", "darwin":
        shell := detectShell()  // reads $SHELL env var
        profile := shellProfile(shell, adapter.HomeDir())
        line := fmt.Sprintf(`export PATH="%s:$PATH"`, dir)
        return appendLineIfAbsent(profile, line)

    case "windows":
        // Modifies HKCU\Environment\Path via registry
        return addToWindowsUserPath(dir)
    }
    return nil
}
```

---

## 11. Configuration System

### 11.1 Config Loading Priority (highest to lowest)

1. **CLI flags** — e.g., `--poll-interval 10` overrides everything.
2. **Project-local config** — `.cue.toml` in the current working directory.
3. **User config** — `~/.cue/config.toml`.
4. **Built-in defaults** — compiled into the binary.

```go
// internal/config/loader.go

func Load(args []string) (*Config, error) {
    cfg := DefaultConfig()

    // Layer 3: User config
    userCfg := filepath.Join(userConfigDir(), "config.toml")
    if _, err := os.Stat(userCfg); err == nil {
        if _, err := toml.DecodeFile(userCfg, cfg); err != nil {
            return nil, fmt.Errorf("invalid config at %s: %w", userCfg, err)
        }
    }

    // Layer 2: Project-local config
    localCfg := ".cue.toml"
    if _, err := os.Stat(localCfg); err == nil {
        if _, err := toml.DecodeFile(localCfg, cfg); err != nil {
            return nil, fmt.Errorf("invalid project config: %w", err)
        }
    }

    // Layer 1: CLI flags handled by cobra (merged externally)
    return cfg, nil
}
```

---

## 12. Edge Cases & Error Handling

### 12.1 Lock Management Edge Cases

| Edge Case | Handling Strategy |
|-----------|------------------|
| **Stale lock file** (package manager crashed, lock remains) | `flock()` test proves lock is not actively held; proceed without queuing. On Windows, check if holding PID is still alive. |
| **Lock held by the user themselves** (two tabs) | Detect by checking if the locking PID is owned by the same user. Warn: `[INFO] You have another cue install running in another terminal. Queuing...` |
| **Lock wait timeout** | Print clear error, remove queue entry, suggest manual check: `sudo lsof /var/lib/dpkg/lock-frontend` |
| **Lock acquired mid-poll** (TOCTOU race) | The underlying package manager handles this; `cue` does not attempt to "steal" locks. If the second install fails with a lock error, re-queue automatically (max 1 auto-requeue). |
| **Multiple commands queued simultaneously** | FIFO queue; sequential execution. Commands are never run in parallel through the queue. |

### 12.2 Pause/Resume Edge Cases

| Edge Case | Handling Strategy |
|-----------|------------------|
| **Package manager doesn't support mid-download resume** (e.g., old pip) | Detect by checking if partial download files exist on resume; restart download; warn user: `[WARN] pip does not support partial resume. Restarting download from cached index.` |
| **Process exits during network outage** (package manager self-terminates) | `cmd.Wait()` returns; record exit code; report failure; suggest `cue store install <stack> --resume` which skips already-installed components. |
| **SIGKILL received while paused** (system shutdown) | `jobs.json` persists the paused state. On next startup, `cue` prints: `[INFO] 1 paused job detected from previous session. Run 'cue resume' to continue.` |
| **Process PID recycled between sessions** | On resume, validate PID against process name before sending signals. If PID is gone, mark job as `orphaned` and surface to user. |
| **Insufficient permissions to SIGSTOP** | Log warning; fall back to killing child processes and offering a re-run. |
| **VPN-connected machine** | VPN routes traffic through a tunnel; ICMP to 1.1.1.1 may fail even with VPN up. Fallback TCP probe to 8.8.8.8:53 handles this. Users can set `probe_host` in config to a VPN-internal host for best accuracy. |

### 12.3 Environment Store Edge Cases

| Edge Case | Handling Strategy |
|-----------|------------------|
| **Component already installed at a different version** | Detect via version probe before install. Prompt: `[WARN] Node.js v18.1.0 already installed. Store requires LTS (v20.x). Upgrade? [y/N]` |
| **Disk space insufficient** | Check `df` / `GetDiskFreeSpaceEx` before starting. Abort early with: `[ERROR] Insufficient disk space. Store requires ~2.5 GB; 1.1 GB available.` |
| **No internet access during store install** | Fail fast with a clear error. Suggest offline install option (future v1.1 feature). |
| **Corporate proxy required** | Read `HTTP_PROXY` / `HTTPS_PROXY` env vars and pass through to child processes. Print proxy info in verbose mode. |
| **Non-admin user** | For components requiring admin (e.g., Docker Engine on Linux), clearly list which components need elevation and which don't. Offer to install what's possible without elevation. |
| **ARM architecture** (Apple Silicon, Raspberry Pi) | Detect via `runtime.GOARCH`. Route to ARM-specific install paths where available (e.g., `brew` on Apple Silicon is `/opt/homebrew`). Flag components without ARM builds. |
| **WSL2 on Windows** | Detected as Linux (Ubuntu). Treat as standard Linux install. Note WSL-specific quirks (e.g., systemd may not be available; Docker Desktop WSL integration required). |

### 12.4 Workspace Backup Edge Cases

| Edge Case | Handling Strategy |
|-----------|------------------|
| **GitHub PAT expired or revoked** | API call returns 401; display: `[ERROR] GitHub token is invalid or expired. Re-authenticate with 'cue workspace auth --token <new-PAT>'.` |
| **history.db too large for Git** | Warn if `> 50 MB`. Offer to export as CSV and include the CSV instead. |
| **Shell config contains plaintext secrets** | Pattern-scan for common secret patterns (`export AWS_SECRET`, `GITHUB_TOKEN=`, API key regexes) before committing. Warn and redact (replace with `# REDACTED by cue`) if found. |
| **Restore on a different OS** | Manifest records source OS. If restoring Linux backup on macOS, warn and skip Linux-only shell files; attempt equivalent macOS installs where possible. |
| **No git installed on new machine** | `cue workspace restore` requires git to clone the repo. Provide the OS-appropriate install command before failing. |

### 12.5 General CLI Edge Cases

| Edge Case | Handling Strategy |
|-----------|------------------|
| **Unknown sub-command** | Show closest match using Levenshtein distance: `Unknown command 'stor'. Did you mean 'store'?` |
| **Concurrent `cue` invocations** | SQLite handles concurrent writes via WAL mode. `jobs.json` and `queue.json` use file-level locking via `flock`. |
| **Filesystem permissions on config dir** | If `~/.cue/` cannot be created or written to, fail clearly and explain the problem. Never silently swallow write errors. |
| **Missing git in macro context** | Before running any `git-*` macro, check that `git` is on PATH. If not: `[ERROR] git not found. Install it with 'cue install git'.` |
| **Running inside a Docker container** | Many store components (Docker itself, desktop apps) will fail inside a container. Detect via `/.dockerenv` existence and warn proactively. |

---

## 13. Build, Distribution & Packaging

### 13.1 Antigravity Build Configuration

```yaml
# Antigravityfile

project: cue
version: 1.0.0

build:
  main: ./main.go
  output: dist/cue
  ldflags:
    - "-s -w"                           # Strip debug symbols; minimise binary size
    - "-X main.Version={{.Version}}"
    - "-X main.BuildDate={{.Date}}"

cross_compile:
  targets:
    - os: linux
      arch: [amd64, arm64]
    - os: darwin
      arch: [amd64, arm64]
    - os: windows
      arch: [amd64]
      ext: .exe

tasks:
  test:
    command: go test ./... -race -cover
  lint:
    command: golangci-lint run
  size-check:
    command: |
      ls -lh dist/cue-linux-amd64
      echo "Must be under 15 MB"
  release:
    depends: [test, lint, build, size-check]
    command: gh release create v{{.Version}} dist/*
```

### 13.2 Installation Methods

| Platform | Method | Command |
|----------|--------|---------|
| **Linux/macOS** | Shell script | `curl -fsSL https://get.gyanesh.help \| bash` |
| **Windows** | PowerShell | `iwr https://get.gyanesh.help/win \| iex` |
| **Homebrew** | Tap | `brew install cue/tap/cue` |
| **Manual** | GitHub Releases | Download pre-built binary from GitHub Releases page |
| **Go install** | Go toolchain | `go install github.com/cue/cue@latest` |

The install script:
1. Detects OS and architecture.
2. Downloads the correct pre-built binary from GitHub Releases.
3. Verifies the SHA256 checksum.
4. Places the binary in `/usr/local/bin` (Linux/macOS) or `%LOCALAPPDATA%\cue\` (Windows) and adds to PATH.
5. Creates `~/.cue/` with default `config.toml`.

### 13.3 Binary Size Budget

| Component | Estimated Size |
|-----------|---------------|
| Go runtime + stdlib | ~4.5 MB |
| cobra + viper | ~1.2 MB |
| go-sqlite3 (CGo) | ~2.0 MB |
| go-github, go-keyring | ~0.8 MB |
| All other dependencies | ~1.5 MB |
| **Estimated total (stripped)** | **~10 MB** |
| **Target ceiling** | **15 MB** |

*Note: CGo (required by go-sqlite3) prevents full static linking on some platforms. A `-tags pure_go` build using `modernc.org/sqlite` (pure Go SQLite) achieves a fully static binary at the cost of ~10–15% query performance — acceptable for this use case.*

---

## 14. Testing Strategy

### 14.1 Unit Tests

| Module | Key Test Cases |
|--------|---------------|
| `adapter/linux` | Correct lock paths per distro; `IsLocked()` returns true on held flock; false on stale file |
| `adapter/windows` | `isProcessRunning("msiexec.exe")` positive and negative cases |
| `queue/poller` | Timeout fires correctly; backoff interval increases; ctx.Done() exits cleanly |
| `job/network_monitor` | `probe()` returns false on unreachable host; event channel emits `NetworkLost` after 3 failures |
| `macro/executor` | All 30 built-in macros resolve without errors; dangerous macros require confirmation |
| `history/query` | FTS search returns correct rows; tag filter excludes untagged entries |
| `workspace/capture` | Known secret patterns are redacted from .gitconfig; `.env` files are excluded |

### 14.2 Integration Tests

- **Store install dry-run:** Mock the OS adapter; verify each store's component list resolves and is topologically sorted.
- **Lock detection integration:** Create a real lockfile in a temp dir; verify `IsLocked()` returns true; release lock; verify returns false.
- **History write/read round-trip:** Write 1000 entries; query by tag, search, date filter; verify row counts.
- **Macro execution (git):** Run `git-undo` in a temp git repo with a known commit; verify commit is reverted.

### 14.3 Cross-Platform Test Matrix

| OS | Arch | Package Manager | CI Provider |
|----|------|----------------|-------------|
| Ubuntu 22.04 | amd64 | apt | GitHub Actions |
| Ubuntu 24.04 | arm64 | apt | GitHub Actions (ARM runner) |
| Fedora 40 | amd64 | dnf | GitHub Actions |
| Arch Linux | amd64 | pacman | GitHub Actions (custom container) |
| macOS 14 (Sonoma) | arm64 | brew | GitHub Actions (macOS runner) |
| Windows Server 2022 | amd64 | winget | GitHub Actions |

---

## 15. Security Considerations

### 15.1 Secret Storage

- GitHub PAT is stored using the OS native keyring (`go-keyring`): macOS Keychain, Windows Credential Manager, Linux Secret Service (libsecret).
- **Never** stored in plaintext in `config.toml` or any file in `~/.cue/`.
- Token is retrieved from keyring only at the moment of use; never logged or printed.

### 15.2 Install Script Safety

- All install scripts fetched from URLs are verified against a SHA256 hash before execution.
- The built-in `.gitignore` for workspace backups excludes: `*.pem`, `*.key`, `id_rsa`, `id_ed25519`, `.env`, `.env.*`, `*.secret`, `*.token`, `.netrc`.
- Pattern scan for common secrets before any `git commit` during backup.

### 15.3 Privilege Escalation Safety

- `sudo` / UAC elevation is invoked only for specific, listed operations (package installs, `/etc/hosts` edit, service management).
- `cue` never stores or caches sudo credentials.
- Every operation requiring elevation prints the exact command that will be run with elevated permissions before requesting elevation.

### 15.4 Child Process Safety

- All user-provided arguments passed to child processes are passed as `args []string` (never interpolated into a shell string), preventing shell injection.
- The `--` separator is used when passing user arguments to underlying tools.
- Custom macros defined by users that use shell interpolation (`{arg}`) are sandboxed via shell quoting rules.

### 15.5 No Telemetry

- Zero telemetry, analytics, or crash reporting in v1.0.
- No outbound network calls during normal operation except:
  - Package manager installs (controlled by the user's invocation).
  - GitHub API calls (only when `workspace backup/restore` is explicitly invoked).
  - Network probes to 1.1.1.1 / 8.8.8.8 (only during active managed installations).
- All three categories are transparent to the user and documented in `cue --privacy`.

---

## Appendix A — Supported Package Manager Matrix

| OS | Primary | Secondary | Tertiary |
|----|---------|-----------|----------|
| Ubuntu/Debian | `apt` | `snap` | `dpkg` (direct) |
| Fedora | `dnf` | `rpm` | — |
| Arch | `pacman` | `yay` (AUR, optional) | — |
| macOS | `brew` | `port` (fallback) | — |
| Windows | `winget` | `choco` | `scoop` |

For each secondary/tertiary manager, the adapter falls back only if the primary is not available.

---

## Appendix B — Macro Quick Reference (v1.0)

| Macro | Category | Dangerous? |
|-------|----------|------------|
| `git-undo` | git | No (`--hard` flag: Yes) |
| `git-undo --hard` | git | Yes |
| `git-clean` | git | Yes |
| `git-save` | git | No |
| `git-unsave` | git | No |
| `git-whoops` | git | No |
| `git-oops-push` | git | Yes |
| `git-log-pretty` | git | No |
| `git-branch-clean` | git | Yes |
| `git-diff-staged` | git | No |
| `find-big-files` | filesystem | No |
| `find-old-logs` | filesystem | No |
| `nuke-node` | nodejs | Yes |
| `port-kill <port>` | network | Yes |
| `port-check <port>` | network | No |
| `env-check` | system | No |
| `disk-check` | system | No |
| `process-find <n>` | system | No |
| `docker-nuke` | docker | Yes |
| `docker-shell <container>` | docker | No |
| `pip-freeze-clean` | python | No |
| `venv-create` | python | No |
| `npm-audit-fix` | nodejs | No |
| `ssh-keygen-github` | ssh | No |
| `hosts-edit` | system | No (requires elevation) |
| `path-add <dir>` | system | No |
| `kill-port` | network | Yes |
| `ip-info` | network | No |
| `cert-check <domain>` | security | No |
| `backup-now` | workspace | No |

---

## Appendix C — Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1 | 2026-03-28 | Architecture | Initial draft |
| 1.0 | 2026-03-28 | Architecture | Full spec with all modules, schemas, edge cases |