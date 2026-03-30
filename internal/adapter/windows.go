//go:build windows

package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Detect returns the Windows adapter on Windows builds.
func Detect() OSAdapter {
	return newWindowsAdapter()
}

type WindowsAdapter struct {
	homeDir string
	pkgMgr  string
}

func newWindowsAdapter() *WindowsAdapter {
	home, _ := os.UserHomeDir()
	mgr := "winget"
	if _, err := exec.LookPath("winget"); err != nil {
		if _, err := exec.LookPath("choco"); err == nil {
			mgr = "choco"
		}
	}
	return &WindowsAdapter{homeDir: home, pkgMgr: mgr}
}

func (a *WindowsAdapter) PackageManagerName() string { return a.pkgMgr }

func (a *WindowsAdapter) InstallPackage(pkg string, args []string) error {
	var cmd *exec.Cmd
	switch a.pkgMgr {
	case "winget":
		cmdArgs := append([]string{"install", "--accept-package-agreements", "--accept-source-agreements"}, args...)
		cmdArgs = append(cmdArgs, pkg)
		cmd = exec.Command("winget", cmdArgs...)
	case "choco":
		cmdArgs := append([]string{"install", "-y"}, args...)
		cmdArgs = append(cmdArgs, pkg)
		cmd = exec.Command("choco", cmdArgs...)
	default:
		return fmt.Errorf("no package manager available")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *WindowsAdapter) UninstallPackage(pkg string) error {
	var cmd *exec.Cmd
	switch a.pkgMgr {
	case "winget":
		cmd = exec.Command("winget", "uninstall", pkg)
	case "choco":
		cmd = exec.Command("choco", "uninstall", "-y", pkg)
	default:
		return fmt.Errorf("no package manager available")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *WindowsAdapter) IsPackageInstalled(pkg string) bool {
	switch a.pkgMgr {
	case "winget":
		out, err := exec.Command("winget", "list", "--id", pkg).Output()
		return err == nil && strings.Contains(string(out), pkg)
	case "choco":
		return exec.Command("choco", "list", "--local-only", pkg).Run() == nil
	}
	return false
}

func (a *WindowsAdapter) PackageVersion(pkg string) (string, error) {
	switch a.pkgMgr {
	case "winget":
		out, err := exec.Command("winget", "list", "--id", pkg).Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	return "", fmt.Errorf("unsupported")
}

func (a *WindowsAdapter) LockPaths() []string {
	return []string{`C:\ProgramData\chocolatey\.chocolatey.lock`}
}

func (a *WindowsAdapter) IsLocked() (bool, string, error) {
	// Check for active msiexec, winget, or choco processes
	processes := []string{"msiexec.exe", "winget.exe", "choco.exe"}
	for _, proc := range processes {
		if isProcessRunning(proc) {
			return true, fmt.Sprintf("active process: %s", proc), nil
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

func (a *WindowsAdapter) SuspendProcess(pid int) error {
	handle, err := windows.OpenProcess(
		windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	proc := ntdll.NewProc("NtSuspendProcess")
	r, _, _ := proc.Call(uintptr(handle))
	if r != 0 {
		return fmt.Errorf("NtSuspendProcess failed: NTSTATUS 0x%x", r)
	}
	return nil
}

func (a *WindowsAdapter) ResumeProcess(pid int) error {
	handle, err := windows.OpenProcess(
		windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	proc := ntdll.NewProc("NtResumeProcess")
	r, _, _ := proc.Call(uintptr(handle))
	if r != 0 {
		return fmt.Errorf("NtResumeProcess failed: NTSTATUS 0x%x", r)
	}
	return nil
}

func (a *WindowsAdapter) KillProcess(pid int) error {
	handle, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)
	return windows.TerminateProcess(handle, 1)
}

func (a *WindowsAdapter) RunElevated(cmdPath string, args []string) error {
	verb, _ := windows.UTF16PtrFromString("runas")
	exe, _ := windows.UTF16PtrFromString(cmdPath)
	argStr, _ := windows.UTF16PtrFromString(strings.Join(args, " "))
	err := windows.ShellExecute(0, verb, exe, argStr, nil, windows.SW_SHOW)
	if err != nil {
		return err
	}
	return nil
}

func (a *WindowsAdapter) HomeDir() string { return a.homeDir }

func (a *WindowsAdapter) ConfigDir() string {
	if appdata := os.Getenv("APPDATA"); appdata != "" {
		return filepath.Join(appdata, "cue")
	}
	return filepath.Join(a.homeDir, ".cue")
}

func (a *WindowsAdapter) OSName() string   { return "windows" }
func (a *WindowsAdapter) OSDistro() string  { return "" }

func (a *WindowsAdapter) HasGPU() bool {
	return exec.Command("nvidia-smi").Run() == nil
}

func (a *WindowsAdapter) ShellConfigPaths() []string {
	paths := []string{}
	// PowerShell profile
	psProfile := filepath.Join(a.homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	if _, err := os.Stat(psProfile); err == nil {
		paths = append(paths, psProfile)
	}
	// Windows Terminal settings
	wtSettings := filepath.Join(os.Getenv("LOCALAPPDATA"), "Packages",
		"Microsoft.WindowsTerminal_8wekyb3d8bbwe", "LocalState", "settings.json")
	if _, err := os.Stat(wtSettings); err == nil {
		paths = append(paths, wtSettings)
	}
	return paths
}

// Ensure unsafe import is used (for ShellExecute pointer casting)
var _ = unsafe.Sizeof(0)
