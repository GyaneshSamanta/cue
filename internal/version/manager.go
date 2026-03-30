package version

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/ui"
)

// Backend defines a version manager backend.
type Backend struct {
	Name       string
	ListCmd    string
	UseCmd     string
	InstallCmd string
	RemoveCmd  string
	PinFile    string
}

var backends = map[string]Backend{
	"python": {
		Name: "pyenv", ListCmd: "pyenv versions", UseCmd: "pyenv global %s",
		InstallCmd: "pyenv install %s", RemoveCmd: "pyenv uninstall -f %s",
		PinFile: ".python-version",
	},
	"node": {
		Name: "nvm", ListCmd: "nvm list", UseCmd: "nvm use %s",
		InstallCmd: "nvm install %s", RemoveCmd: "nvm uninstall %s",
		PinFile: ".nvmrc",
	},
	"ruby": {
		Name: "rbenv", ListCmd: "rbenv versions", UseCmd: "rbenv global %s",
		InstallCmd: "rbenv install %s", RemoveCmd: "rbenv uninstall -f %s",
		PinFile: ".ruby-version",
	},
	"java": {
		Name: "sdkman", ListCmd: `bash -c 'source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk list java'`,
		UseCmd: `bash -c 'source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk use java %s'`,
		InstallCmd: `bash -c 'source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install java %s'`,
		RemoveCmd: `bash -c 'source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk uninstall java %s'`,
	},
	"rust": {
		Name: "rustup", ListCmd: "rustup toolchain list",
		UseCmd: "rustup default %s", InstallCmd: "rustup toolchain install %s",
		RemoveCmd: "rustup toolchain uninstall %s",
	},
	"go": {
		Name: "manual", ListCmd: "go version",
	},
	"terraform": {
		Name: "tfenv", ListCmd: "tfenv list", UseCmd: "tfenv use %s",
		InstallCmd: "tfenv install %s", RemoveCmd: "tfenv uninstall %s",
		PinFile: ".terraform-version",
	},
}

// ListVersions shows installed versions for a runtime.
func ListVersions(runtime string) error {
	b, ok := backends[runtime]
	if !ok {
		return fmt.Errorf("unsupported runtime: %s. Supported: %s", runtime, supportedRuntimes())
	}

	ui.PrintHeader(fmt.Sprintf("Installed versions: %s (via %s)", runtime, b.Name))
	return runShellCmd(b.ListCmd)
}

// UseVersion switches the active version for a runtime.
func UseVersion(rt, ver string) error {
	b, ok := backends[rt]
	if !ok {
		return fmt.Errorf("unsupported runtime: %s", rt)
	}
	if b.UseCmd == "" {
		return fmt.Errorf("%s does not support version switching via cue", rt)
	}

	cmd := fmt.Sprintf(b.UseCmd, ver)
	ui.PrintStep(fmt.Sprintf("Switching %s to %s...", rt, ver))
	if err := runShellCmd(cmd); err != nil {
		return err
	}
	ui.PrintSuccess(fmt.Sprintf("%s now using version %s", rt, ver))
	return nil
}

// InstallVersion installs a specific version of a runtime.
func InstallVersion(rt, ver string) error {
	b, ok := backends[rt]
	if !ok {
		return fmt.Errorf("unsupported runtime: %s", rt)
	}
	if b.InstallCmd == "" {
		return fmt.Errorf("%s does not support version installation via cue", rt)
	}

	cmd := fmt.Sprintf(b.InstallCmd, ver)
	ui.PrintStep(fmt.Sprintf("Installing %s %s...", rt, ver))
	return runShellCmd(cmd)
}

// RemoveVersion removes a specific version.
func RemoveVersion(rt, ver string) error {
	b, ok := backends[rt]
	if !ok {
		return fmt.Errorf("unsupported runtime: %s", rt)
	}
	if b.RemoveCmd == "" {
		return fmt.Errorf("%s does not support version removal via cue", rt)
	}

	cmd := fmt.Sprintf(b.RemoveCmd, ver)
	return runShellCmd(cmd)
}

// Pin writes a version pin file to the current directory.
func Pin(rt string) error {
	b, ok := backends[rt]
	if !ok {
		return fmt.Errorf("unsupported runtime: %s", rt)
	}
	if b.PinFile == "" {
		return fmt.Errorf("%s does not support version pinning", rt)
	}

	// Get current version
	ver := getCurrentVersion(rt)
	if ver == "" {
		return fmt.Errorf("cannot detect current %s version", rt)
	}

	os.WriteFile(b.PinFile, []byte(ver+"\n"), 0644)
	ui.PrintSuccess(fmt.Sprintf("Wrote %s with version %s", b.PinFile, ver))
	return nil
}

// ShowCurrent displays the active versions of all managed runtimes.
func ShowCurrent() {
	ui.PrintHeader("Active Runtime Versions")
	runtimes := []string{"python", "node", "go", "rust", "java", "ruby", "terraform"}
	for _, rt := range runtimes {
		ver := getCurrentVersion(rt)
		if ver != "" {
			fmt.Printf("  %-12s %s\n", rt, ver)
		}
	}
}

func getCurrentVersion(rt string) string {
	cmds := map[string][]string{
		"python":    {"python3", "--version"},
		"node":      {"node", "--version"},
		"go":        {"go", "version"},
		"rust":      {"rustc", "--version"},
		"java":      {"java", "-version"},
		"ruby":      {"ruby", "--version"},
		"terraform": {"terraform", "--version"},
	}
	args, ok := cmds[rt]
	if !ok {
		return ""
	}
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return ""
	}
	s := strings.TrimSpace(string(out))
	if idx := strings.IndexByte(s, '\n'); idx > 0 {
		s = s[:idx]
	}
	return s
}

func supportedRuntimes() string {
	rts := make([]string, 0, len(backends))
	for k := range backends {
		rts = append(rts, k)
	}
	return strings.Join(rts, ", ")
}

func runShellCmd(cmd string) error {
	shell := "bash"
	flag := "-c"
	if runtime.GOOS == "windows" {
		shell = "powershell"
		flag = "-Command"
	}
	c := exec.Command(shell, flag, cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// Suppress unused
var _ = filepath.Join
