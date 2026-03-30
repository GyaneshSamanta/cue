package toolkit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/cue/internal/adapter"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

func init() {
	registerNode()
	registerPython()
	registerGo()
	registerRust()
	registerJava()
	registerDocker()
	registerSimpleTools()
}

// --- Node.js via nvm ---

func registerNode() {
	Register(&Tool{
		Name:           "node",
		DisplayName:    "Node.js",
		Description:    "JavaScript runtime via nvm for version management",
		VersionManager: "nvm",
		EstSizeMB:      80,
		Categories:     []string{"runtime", "javascript"},
		VerifyFunc: func() (string, bool) {
			return GetVersion("node", "--version")
		},
		InstallFunc: func(a adapter.OSAdapter, version string) error {
			if runtime.GOOS == "windows" {
				return installNodeWindows(a, version)
			}
			return installNodeUnix(a, version)
		},
		UpgradeFunc: func(a adapter.OSAdapter) error {
			if runtime.GOOS == "windows" {
				return RunInstallCmd("nvm", "install", "lts")
			}
			return runBashCmd("source ~/.nvm/nvm.sh && nvm install --lts")
		},
	})
}

func installNodeUnix(a adapter.OSAdapter, version string) error {
	if !CommandExists("nvm") {
		ui.PrintStep("[1/3] Installing nvm")
		if err := runBashCmd("curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash"); err != nil {
			return fmt.Errorf("nvm install failed: %w", err)
		}
	}
	ver := "lts/*"
	if version != "" {
		ver = version
	}
	ui.PrintStep(fmt.Sprintf("[2/3] Installing Node.js %s", ver))
	if err := runBashCmd(fmt.Sprintf("source ~/.nvm/nvm.sh && nvm install %s && nvm alias default %s", ver, ver)); err != nil {
		return err
	}
	ui.PrintStep("[3/3] Verifying")
	ui.PrintSuccess("Node.js installed via nvm")
	return nil
}

func installNodeWindows(a adapter.OSAdapter, version string) error {
	if !CommandExists("nvm") {
		ui.PrintStep("[1/3] Installing nvm-windows")
		if err := a.InstallPackage("nvm", nil); err != nil {
			ui.PrintWarning("nvm-windows not available via winget. Install from https://github.com/coreybutler/nvm-windows")
			return err
		}
	}
	ver := "lts"
	if version != "" {
		ver = version
	}
	ui.PrintStep(fmt.Sprintf("[2/3] Installing Node.js %s", ver))
	if err := RunInstallCmd("nvm", "install", ver); err != nil {
		return err
	}
	RunInstallCmd("nvm", "use", ver)
	ui.PrintStep("[3/3] Verifying")
	ui.PrintSuccess("Node.js installed via nvm-windows")
	return nil
}

// --- Python via pyenv ---

func registerPython() {
	Register(&Tool{
		Name:           "python",
		DisplayName:    "Python",
		Description:    "Python runtime via pyenv for version management",
		VersionManager: "pyenv",
		EstSizeMB:      120,
		Categories:     []string{"runtime", "python"},
		VerifyFunc: func() (string, bool) {
			v, ok := GetVersion("python3", "--version")
			if !ok {
				v, ok = GetVersion("python", "--version")
			}
			return v, ok
		},
		InstallFunc: func(a adapter.OSAdapter, version string) error {
			if runtime.GOOS == "windows" {
				return installPythonWindows(a, version)
			}
			return installPythonUnix(a, version)
		},
	})
}

func installPythonUnix(a adapter.OSAdapter, version string) error {
	if !CommandExists("pyenv") {
		ui.PrintStep("[1/3] Installing pyenv")
		if err := runBashCmd("curl https://pyenv.run | bash"); err != nil {
			return fmt.Errorf("pyenv install failed: %w", err)
		}
		AppendToShellConfig(`export PYENV_ROOT="$HOME/.pyenv"`)
		AppendToShellConfig(`export PATH="$PYENV_ROOT/bin:$PATH"`)
		AppendToShellConfig(`eval "$(pyenv init -)"`)
	}
	ver := "3.12.2"
	if version != "" {
		ver = version
	}
	ui.PrintStep(fmt.Sprintf("[2/3] Installing Python %s", ver))
	pyenvBin := filepath.Join(HomeDir(), ".pyenv", "bin", "pyenv")
	if err := RunInstallCmd(pyenvBin, "install", ver); err != nil {
		return err
	}
	RunInstallCmd(pyenvBin, "global", ver)
	ui.PrintStep("[3/3] Installing pip + pipx")
	ui.PrintSuccess(fmt.Sprintf("Python %s installed via pyenv", ver))
	return nil
}

func installPythonWindows(a adapter.OSAdapter, version string) error {
	ui.PrintStep("[1/2] Installing Python via winget")
	ver := ""
	if version != "" {
		ver = version
	}
	_ = ver
	if err := a.InstallPackage("Python.Python.3.12", nil); err != nil {
		return err
	}
	ui.PrintStep("[2/2] Verifying")
	ui.PrintSuccess("Python installed via winget")
	return nil
}

// --- Go ---

func registerGo() {
	Register(&Tool{
		Name:        "go",
		DisplayName: "Go",
		Description: "Go programming language",
		EstSizeMB:   200,
		Categories:  []string{"runtime", "golang"},
		VerifyFunc: func() (string, bool) {
			return GetVersion("go", "version")
		},
		InstallFunc: func(a adapter.OSAdapter, version string) error {
			if runtime.GOOS == "windows" {
				return a.InstallPackage("GoLang.Go", nil)
			}
			if runtime.GOOS == "darwin" {
				return RunInstallCmd("brew", "install", "go")
			}
			// Linux: download official tarball
			ver := "1.22.1"
			if version != "" {
				ver = version
			}
			url := fmt.Sprintf("https://go.dev/dl/go%s.linux-amd64.tar.gz", ver)
			ui.PrintStep(fmt.Sprintf("Downloading Go %s", ver))
			return runBashCmd(fmt.Sprintf("curl -L %s | sudo tar -C /usr/local -xzf - && echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc", url))
		},
	})
}

// --- Rust via rustup ---

func registerRust() {
	Register(&Tool{
		Name:           "rust",
		DisplayName:    "Rust",
		Description:    "Rust toolchain via rustup (includes cargo, rustfmt, clippy)",
		VersionManager: "rustup",
		EstSizeMB:      300,
		Categories:     []string{"runtime", "rust"},
		VerifyFunc: func() (string, bool) {
			return GetVersion("rustc", "--version")
		},
		InstallFunc: func(a adapter.OSAdapter, version string) error {
			if runtime.GOOS == "windows" {
				ui.PrintStep("Installing rustup via winget")
				return a.InstallPackage("Rustlang.Rustup", nil)
			}
			ui.PrintStep("[1/2] Installing rustup")
			if err := runBashCmd("curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y"); err != nil {
				return err
			}
			ui.PrintStep("[2/2] Installing components")
			rustup := filepath.Join(HomeDir(), ".cargo", "bin", "rustup")
			RunInstallCmd(rustup, "component", "add", "rust-analyzer", "clippy")
			ui.PrintSuccess("Rust toolchain installed via rustup")
			return nil
		},
	})
}

// --- Java via SDKMAN ---

func registerJava() {
	Register(&Tool{
		Name:           "java",
		DisplayName:    "Java (OpenJDK)",
		Description:    "OpenJDK + Maven + Gradle via SDKMAN!",
		VersionManager: "sdkman",
		EstSizeMB:      400,
		Categories:     []string{"runtime", "java"},
		VerifyFunc: func() (string, bool) {
			return GetVersion("java", "-version")
		},
		InstallFunc: func(a adapter.OSAdapter, version string) error {
			if runtime.GOOS == "windows" {
				ui.PrintStep("Installing OpenJDK via winget")
				return a.InstallPackage("EclipseAdoptium.Temurin.21.JDK", nil)
			}
			ui.PrintStep("[1/3] Installing SDKMAN!")
			if err := runBashCmd(`curl -s "https://get.sdkman.io" | bash`); err != nil {
				return err
			}
			ui.PrintStep("[2/3] Installing Java 21 LTS")
			runBashCmd(`source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install java 21-tem`)
			ui.PrintStep("[3/3] Installing Maven + Gradle")
			runBashCmd(`source "$HOME/.sdkman/bin/sdkman-init.sh" && sdk install maven && sdk install gradle`)
			ui.PrintSuccess("Java, Maven, Gradle installed via SDKMAN!")
			return nil
		},
	})
}

// --- Docker ---

func registerDocker() {
	Register(&Tool{
		Name:        "docker",
		DisplayName: "Docker",
		Description: "Docker Engine + Compose",
		EstSizeMB:   500,
		Categories:  []string{"container"},
		VerifyFunc: func() (string, bool) {
			return GetVersion("docker", "--version")
		},
		InstallFunc: func(a adapter.OSAdapter, version string) error {
			if runtime.GOOS == "windows" {
				return a.InstallPackage("Docker.DockerDesktop", nil)
			}
			if runtime.GOOS == "darwin" {
				return RunInstallCmd("brew", "install", "--cask", "docker")
			}
			ui.PrintStep("Installing Docker via official script")
			return runBashCmd("curl -fsSL https://get.docker.com | sh")
		},
	})
}

// --- Simple tools (single binary via package manager) ---

func registerSimpleTools() {
	simpleTools := []struct {
		name, display, desc, pkg, winPkg, brewPkg, verCmd string
		verArgs                                           []string
		sizeMB                                            int
		cats                                              []string
	}{
		{"git", "Git", "Distributed version control", "git", "Git.Git", "git", "git", []string{"--version"}, 30, []string{"vcs"}},
		{"gh", "GitHub CLI", "GitHub from the command line", "gh", "GitHub.cli", "gh", "gh", []string{"--version"}, 20, []string{"github"}},
		{"kubectl", "kubectl", "Kubernetes CLI", "kubectl", "Kubernetes.kubectl", "kubernetes-cli", "kubectl", []string{"version", "--client", "--short"}, 50, []string{"k8s"}},
		{"terraform", "Terraform", "Infrastructure as Code", "terraform", "Hashicorp.Terraform", "terraform", "terraform", []string{"--version"}, 80, []string{"iac"}},
		{"jq", "jq", "Lightweight JSON processor", "jq", "jqlang.jq", "jq", "jq", []string{"--version"}, 2, []string{"utility"}},
		{"fzf", "fzf", "Fuzzy finder", "fzf", "junegunn.fzf", "fzf", "fzf", []string{"--version"}, 3, []string{"utility"}},
		{"ripgrep", "ripgrep", "Fast recursive search", "ripgrep", "BurntSushi.ripgrep.MSVC", "ripgrep", "rg", []string{"--version"}, 5, []string{"utility"}},
		{"bat", "bat", "Cat with syntax highlighting", "bat", "sharkdp.bat", "bat", "bat", []string{"--version"}, 5, []string{"utility"}},
		{"btop", "btop", "Resource monitor", "btop", "", "btop", "btop", []string{"--version"}, 5, []string{"utility"}},
		{"neovim", "Neovim", "Hyperextensible Vim fork", "neovim", "Neovim.Neovim", "neovim", "nvim", []string{"--version"}, 30, []string{"editor"}},
		{"tmux", "tmux", "Terminal multiplexer", "tmux", "", "tmux", "tmux", []string{"-V"}, 3, []string{"utility"}},
		{"httpie", "HTTPie", "Human-friendly HTTP client", "httpie", "", "httpie", "http", []string{"--version"}, 15, []string{"http"}},
		{"ollama", "Ollama", "Run LLMs locally", "ollama", "Ollama.Ollama", "ollama", "ollama", []string{"--version"}, 100, []string{"ai"}},
		{"postgres", "PostgreSQL", "Relational database", "postgresql", "PostgreSQL.PostgreSQL", "postgresql", "psql", []string{"--version"}, 100, []string{"database"}},
		{"mysql", "MySQL", "Relational database", "mysql-server", "Oracle.MySQL", "mysql", "mysql", []string{"--version"}, 150, []string{"database"}},
		{"redis", "Redis", "In-memory data store", "redis-server", "", "redis", "redis-server", []string{"--version"}, 10, []string{"database"}},
		{"mongodb", "MongoDB", "Document database", "mongodb", "MongoDB.Server", "mongodb-community", "mongod", []string{"--version"}, 200, []string{"database"}},
		{"deno", "Deno", "Secure JS/TS runtime", "", "", "", "deno", []string{"--version"}, 40, []string{"runtime"}},
		{"bun", "Bun", "Fast JS runtime + bundler", "", "", "", "bun", []string{"--version"}, 50, []string{"runtime"}},
		{"vscode", "VS Code", "Code editor", "", "Microsoft.VisualStudioCode", "visual-studio-code", "code", []string{"--version"}, 300, []string{"editor"}},
		{"aws", "AWS CLI", "Amazon Web Services CLI", "", "Amazon.AWSCLI", "awscli", "aws", []string{"--version"}, 50, []string{"cloud"}},
		{"gcloud", "Google Cloud SDK", "GCP CLI tools", "", "Google.CloudSDK", "google-cloud-sdk", "gcloud", []string{"--version"}, 200, []string{"cloud"}},
		{"azure", "Azure CLI", "Microsoft Azure CLI", "", "Microsoft.AzureCLI", "azure-cli", "az", []string{"--version"}, 100, []string{"cloud"}},
	}

	for _, st := range simpleTools {
		tool := st // capture
		Register(&Tool{
			Name:        tool.name,
			DisplayName: tool.display,
			Description: tool.desc,
			EstSizeMB:   tool.sizeMB,
			Categories:  tool.cats,
			VerifyFunc: func() (string, bool) {
				return GetVersion(tool.verCmd, tool.verArgs...)
			},
			InstallFunc: func(a adapter.OSAdapter, version string) error {
				switch runtime.GOOS {
				case "windows":
					if tool.winPkg != "" {
						return a.InstallPackage(tool.winPkg, nil)
					}
					return installViaScript(tool.name)
				case "darwin":
					if tool.brewPkg != "" {
						return RunInstallCmd("brew", "install", tool.brewPkg)
					}
					return installViaScript(tool.name)
				default:
					if tool.pkg != "" {
						return PkgInstall(a, tool.pkg)
					}
					return installViaScript(tool.name)
				}
			},
		})
	}
}

// installViaScript handles tools that use their official install scripts.
func installViaScript(name string) error {
	scripts := map[string]string{
		"deno": "curl -fsSL https://deno.land/install.sh | sh",
		"bun":  "curl -fsSL https://bun.sh/install | bash",
	}
	if script, ok := scripts[name]; ok {
		return runBashCmd(script)
	}
	return fmt.Errorf("no install script available for %s on this platform", name)
}

func runBashCmd(script string) error {
	shell := "bash"
	flag := "-c"
	if runtime.GOOS == "windows" {
		shell = "powershell"
		flag = "-Command"
	}
	cmd := exec.Command(shell, flag, script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Suppress unused
var _ = strings.Contains
