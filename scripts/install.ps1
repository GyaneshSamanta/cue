# cue installer for Windows
# Usage: iwr https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.ps1 | iex
# Or:    .\scripts\install.ps1

$ErrorActionPreference = "Stop"

$Repo = "GyaneshSamanta/cue"
$BinaryName = "cue.exe"
$InstallDir = "$env:LOCALAPPDATA\cue"
$ConfigDir = "$env:APPDATA\cue"

function Write-Step($msg)  { Write-Host "▸ $msg" -ForegroundColor Cyan }
function Write-Ok($msg)    { Write-Host "✔ $msg" -ForegroundColor Green }
function Write-Warn($msg)  { Write-Host "⚠ $msg" -ForegroundColor Yellow }
function Write-Err($msg)   { Write-Host "✖ $msg" -ForegroundColor Red }

# Banner
Write-Host ""
Write-Host "  ╔══════════════════════════════════════╗" -ForegroundColor Magenta
Write-Host "  ║       cue installer         ║" -ForegroundColor Magenta
Write-Host "  ║   Cross-Platform CLI Dev Utility     ║" -ForegroundColor Magenta
Write-Host "  ╚══════════════════════════════════════╝" -ForegroundColor Magenta
Write-Host ""

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
Write-Ok "Detected: windows/$Arch"

# Get latest version
Write-Step "Fetching latest version..."
try {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
    $Version = $Release.tag_name -replace '^v', ''
    Write-Ok "Latest version: v$Version"
} catch {
    $Version = "1.0.0"
    Write-Warn "Could not fetch latest version, using v$Version"
}

# Download
$FileName = "cue-windows-$Arch.exe"
$Url = "https://github.com/$Repo/releases/download/v$Version/$FileName"
$TempFile = Join-Path $env:TEMP $BinaryName

Write-Step "Downloading cue v$Version..."
try {
    Invoke-WebRequest -Uri $Url -OutFile $TempFile -UseBasicParsing
    Write-Ok "Downloaded successfully"
} catch {
    Write-Err "Download failed from: $Url"
    Write-Warn "Try manual download from: https://github.com/$Repo/releases"
    exit 1
}

# Install
Write-Step "Installing to $InstallDir..."
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}
Move-Item -Path $TempFile -Destination (Join-Path $InstallDir $BinaryName) -Force
Write-Ok "Installed to $InstallDir\$BinaryName"

# Add to PATH
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Step "Adding to user PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Ok "Added $InstallDir to PATH"
} else {
    Write-Ok "$InstallDir already in PATH"
}

# Config directory
if (-not (Test-Path $ConfigDir)) {
    New-Item -ItemType Directory -Path $ConfigDir -Force | Out-Null
    Write-Ok "Created config directory: $ConfigDir"
}

# Default config
$ConfigFile = Join-Path $ConfigDir "config.toml"
if (-not (Test-Path $ConfigFile)) {
    @"
[core]
lock_poll_interval_secs = 5
lock_timeout_mins = 30
adaptive_backoff = true
notify_on_completion = true

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
"@ | Set-Content -Path $ConfigFile -Encoding UTF8
    Write-Ok "Created default config"
}

# Verify
Write-Step "Verifying installation..."
try {
    $ver = & (Join-Path $InstallDir $BinaryName) --version 2>&1
    Write-Ok "Installation verified: $ver"
} catch {
    Write-Warn "Installed but verification failed. Try opening a new terminal."
}

Write-Host ""
Write-Ok "Installation complete! Run 'cue --help' to get started."
Write-Host "  (You may need to open a new terminal for PATH changes to take effect)" -ForegroundColor Gray
Write-Host ""
