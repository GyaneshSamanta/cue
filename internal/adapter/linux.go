//go:build linux

package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// Detect returns the Linux adapter on Linux builds.
func Detect() OSAdapter {
	return newLinuxAdapter()
}

type LinuxAdapter struct {
	distro    string
	pkgMgr    string
	lockPaths []string
	homeDir   string
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

func newLinuxAdapter() *LinuxAdapter {
	home, _ := os.UserHomeDir()
	distro := detectLinuxDistro()
	a := &LinuxAdapter{distro: distro, homeDir: home}

	switch {
	case contains([]string{"ubuntu", "debian", "linuxmint", "pop", "elementary"}, distro):
		a.pkgMgr = "apt"
		a.lockPaths = aptLockPaths
	case contains([]string{"fedora", "rhel", "centos", "rocky", "alma"}, distro):
		a.pkgMgr = "dnf"
		a.lockPaths = dnfLockPaths
	case contains([]string{"arch", "manjaro", "endeavouros"}, distro):
		a.pkgMgr = "pacman"
		a.lockPaths = pacmanLockPaths
	default:
		a.pkgMgr = detectByBinary()
		a.lockPaths = nil
	}
	return a
}

func detectLinuxDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
	}
	return ""
}

func detectByBinary() string {
	for _, p := range []struct{ bin, name string }{
		{"apt", "apt"}, {"dnf", "dnf"}, {"pacman", "pacman"}, {"zypper", "zypper"},
	} {
		if _, err := exec.LookPath(p.bin); err == nil {
			return p.name
		}
	}
	return "unknown"
}

func (a *LinuxAdapter) PackageManagerName() string { return a.pkgMgr }

func (a *LinuxAdapter) InstallPackage(pkg string, args []string) error {
	var cmd *exec.Cmd
	switch a.pkgMgr {
	case "apt":
		cmdArgs := append([]string{"install", "-y"}, args...)
		cmdArgs = append(cmdArgs, pkg)
		cmd = exec.Command("sudo", append([]string{"apt"}, cmdArgs...)...)
	case "dnf":
		cmdArgs := append([]string{"install", "-y"}, args...)
		cmdArgs = append(cmdArgs, pkg)
		cmd = exec.Command("sudo", append([]string{"dnf"}, cmdArgs...)...)
	case "pacman":
		cmdArgs := append([]string{"-S", "--noconfirm"}, args...)
		cmdArgs = append(cmdArgs, pkg)
		cmd = exec.Command("sudo", append([]string{"pacman"}, cmdArgs...)...)
	default:
		return fmt.Errorf("unsupported package manager: %s", a.pkgMgr)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *LinuxAdapter) UninstallPackage(pkg string) error {
	var cmd *exec.Cmd
	switch a.pkgMgr {
	case "apt":
		cmd = exec.Command("sudo", "apt", "remove", "-y", pkg)
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "remove", "-y", pkg)
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-R", "--noconfirm", pkg)
	default:
		return fmt.Errorf("unsupported package manager: %s", a.pkgMgr)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *LinuxAdapter) IsPackageInstalled(pkg string) bool {
	switch a.pkgMgr {
	case "apt":
		return exec.Command("dpkg", "-s", pkg).Run() == nil
	case "dnf":
		return exec.Command("rpm", "-q", pkg).Run() == nil
	case "pacman":
		return exec.Command("pacman", "-Qi", pkg).Run() == nil
	}
	return false
}

func (a *LinuxAdapter) PackageVersion(pkg string) (string, error) {
	var cmd *exec.Cmd
	switch a.pkgMgr {
	case "apt":
		cmd = exec.Command("dpkg-query", "-W", "-f=${Version}", pkg)
	case "dnf":
		cmd = exec.Command("rpm", "-q", "--qf", "%{VERSION}", pkg)
	case "pacman":
		cmd = exec.Command("pacman", "-Qi", pkg)
	default:
		return "", fmt.Errorf("unsupported")
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (a *LinuxAdapter) LockPaths() []string { return a.lockPaths }

func (a *LinuxAdapter) IsLocked() (bool, string, error) {
	for _, path := range a.lockPaths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			f.Close()
			return true, fmt.Sprintf("lock held: %s", path), nil
		}
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}
	return false, "", nil
}

func (a *LinuxAdapter) SuspendProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGSTOP)
}

func (a *LinuxAdapter) ResumeProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGCONT)
}

func (a *LinuxAdapter) KillProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGKILL)
}

func (a *LinuxAdapter) RunElevated(cmd string, args []string) error {
	if os.Getuid() == 0 {
		return exec.Command(cmd, args...).Run()
	}
	return exec.Command("sudo", append([]string{cmd}, args...)...).Run()
}

func (a *LinuxAdapter) HomeDir() string   { return a.homeDir }
func (a *LinuxAdapter) ConfigDir() string  { return filepath.Join(a.homeDir, ".cue") }
func (a *LinuxAdapter) OSName() string     { return "linux" }
func (a *LinuxAdapter) OSDistro() string   { return a.distro }

func (a *LinuxAdapter) HasGPU() bool {
	return exec.Command("nvidia-smi").Run() == nil
}

func (a *LinuxAdapter) ShellConfigPaths() []string {
	paths := []string{}
	candidates := []string{
		filepath.Join(a.homeDir, ".bashrc"),
		filepath.Join(a.homeDir, ".bash_profile"),
		filepath.Join(a.homeDir, ".zshrc"),
		filepath.Join(a.homeDir, ".profile"),
		filepath.Join(a.homeDir, ".config", "fish", "config.fish"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			paths = append(paths, p)
		}
	}
	return paths
}
