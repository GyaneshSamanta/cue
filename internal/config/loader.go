package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Current holds the active merged configuration.
var Current *Config

// Config holds all user-configurable settings.
type Config struct {
	Core      CoreConfig      `toml:"core"`
	Network   NetworkConfig   `toml:"network"`
	History   HistoryConfig   `toml:"history"`
	Workspace WorkspaceConfig `toml:"workspace"`
	UI        UIConfig        `toml:"ui"`
}

type CoreConfig struct {
	LockPollIntervalSecs int  `toml:"lock_poll_interval_secs"`
	LockTimeoutMins      int  `toml:"lock_timeout_mins"`
	AdaptiveBackoff      bool `toml:"adaptive_backoff"`
	NotifyOnCompletion   bool `toml:"notify_on_completion"`
}

type NetworkConfig struct {
	ProbeHost         string `toml:"probe_host"`
	ProbeFallbackHost string `toml:"probe_fallback_host"`
	ProbeFallbackPort int    `toml:"probe_fallback_port"`
	FailThreshold     int    `toml:"fail_threshold"`
	RecoveryThreshold int    `toml:"recovery_threshold"`
	ProbeIntervalSecs int    `toml:"probe_interval_secs"`
}

type HistoryConfig struct {
	MaxEntries          int    `toml:"max_entries"`
	DefaultDisplayCount int    `toml:"default_display_count"`
	ExportDir           string `toml:"export_dir"`
}

type WorkspaceConfig struct {
	GithubRepoName          string `toml:"github_repo_name"`
	BackupShellConfigs      bool   `toml:"backup_shell_configs"`
	BackupVSCode            bool   `toml:"backup_vscode"`
	BackupHistory           bool   `toml:"backup_history"`
	AutoBackupIntervalHours int    `toml:"auto_backup_interval_hours"`
}

type UIConfig struct {
	Color             bool   `toml:"color"`
	ProgressStyle     string `toml:"progress_style"`
	ExplainAfterMacro bool   `toml:"explain_after_macro"`
}

// DefaultConfig returns the compiled-in defaults.
func DefaultConfig() *Config {
	return &Config{
		Core: CoreConfig{
			LockPollIntervalSecs: 5,
			LockTimeoutMins:      30,
			AdaptiveBackoff:      true,
			NotifyOnCompletion:   true,
		},
		Network: NetworkConfig{
			ProbeHost:         "1.1.1.1",
			ProbeFallbackHost: "8.8.8.8",
			ProbeFallbackPort: 53,
			FailThreshold:     3,
			RecoveryThreshold: 1,
			ProbeIntervalSecs: 10,
		},
		History: HistoryConfig{
			MaxEntries:          50000,
			DefaultDisplayCount: 20,
			ExportDir:           filepath.Join(ConfigDir(), "exports"),
		},
		Workspace: WorkspaceConfig{
			GithubRepoName:     "dev-workspace-backup",
			BackupShellConfigs: true,
		},
		UI: UIConfig{
			Color:             true,
			ProgressStyle:     "bar",
			ExplainAfterMacro: true,
		},
	}
}

// ConfigDir returns the gyanesh-help config directory path.
func ConfigDir() string {
	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "gyanesh-help")
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".gyanesh-help")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gyanesh-help")
}

// Load applies the 3-tier config priority: defaults → user config → project-local.
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Ensure config dir exists
	dir := ConfigDir()
	os.MkdirAll(dir, 0755)

	// Layer: user config
	userCfg := filepath.Join(dir, "config.toml")
	if _, err := os.Stat(userCfg); err == nil {
		if _, err := toml.DecodeFile(userCfg, cfg); err != nil {
			return nil, err
		}
	}

	// Layer: project-local config
	localCfg := ".gyanesh-help.toml"
	if _, err := os.Stat(localCfg); err == nil {
		toml.DecodeFile(localCfg, cfg)
	}

	Current = cfg
	return cfg, nil
}
