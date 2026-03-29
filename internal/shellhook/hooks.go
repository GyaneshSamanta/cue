package shellhook

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

const hookMarkerStart = "# >>> gyanesh-help shell hooks >>>"
const hookMarkerEnd = "# <<< gyanesh-help shell hooks <<<"

// Install adds shell hooks to the user's shell config.
func Install(shell string) error {
	configPath := shellConfigPath(shell)
	if configPath == "" {
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	hookCode := generateHooks(shell)
	if hookCode == "" {
		return fmt.Errorf("no hooks available for shell: %s", shell)
	}

	// Read existing config
	data, _ := os.ReadFile(configPath)
	content := string(data)

	// Check if already installed
	if strings.Contains(content, hookMarkerStart) {
		ui.PrintInfo("Shell hooks already installed. Use 'gyanesh-help shell-hook uninstall' to remove first.")
		return nil
	}

	// Append hooks
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + hookMarkerStart + "\n" + hookCode + "\n" + hookMarkerEnd + "\n")
	if err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("Shell hooks installed in %s", configPath))
	ui.PrintInfo("Restart your shell or run: source " + configPath)
	return nil
}

// Uninstall removes shell hooks from the user's shell config.
func Uninstall(shell string) error {
	configPath := shellConfigPath(shell)
	if configPath == "" {
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	content := string(data)
	startIdx := strings.Index(content, hookMarkerStart)
	endIdx := strings.Index(content, hookMarkerEnd)

	if startIdx == -1 || endIdx == -1 {
		ui.PrintInfo("No shell hooks found to uninstall.")
		return nil
	}

	// Remove the hook block
	newContent := content[:startIdx] + content[endIdx+len(hookMarkerEnd):]
	newContent = strings.TrimRight(newContent, "\n") + "\n"

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return err
	}

	ui.PrintSuccess("Shell hooks removed. Restart your shell to apply.")
	return nil
}

func shellConfigPath(shell string) string {
	home, _ := os.UserHomeDir()
	switch shell {
	case "zsh":
		return filepath.Join(home, ".zshrc")
	case "bash":
		p := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(p); err != nil {
			p = filepath.Join(home, ".bash_profile")
		}
		return p
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish")
	case "powershell":
		return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	}
	return ""
}

func generateHooks(shell string) string {
	switch shell {
	case "zsh", "bash":
		return `# Auto-activate virtual environments on cd
_gyanesh_help_chpwd() {
  # Auto-activate .venv/venv
  if [ -d ".venv" ] && [ -f ".venv/bin/activate" ]; then
    source .venv/bin/activate 2>/dev/null
  elif [ -d "venv" ] && [ -f "venv/bin/activate" ]; then
    source venv/bin/activate 2>/dev/null
  fi

  # Auto-switch project tag from .gyanesh-help config
  if [ -f ".gyanesh-help" ]; then
    local tag=$(grep '^tag' .gyanesh-help 2>/dev/null | head -1 | sed 's/.*= *"\(.*\)"/\1/')
    if [ -n "$tag" ]; then
      gyanesh-help tag set "$tag" 2>/dev/null
    fi
  fi

  # Stack detection hint
  if [ -f "package.json" ] && ! command -v node &>/dev/null; then
    echo "  ℹ  Node.js project detected but node is not installed. Try: gyanesh-help toolkit install node"
  fi
  if [ -f "requirements.txt" ] && ! command -v python3 &>/dev/null; then
    echo "  ℹ  Python project detected but python3 is not installed. Try: gyanesh-help toolkit install python"
  fi
  if [ -f "Cargo.toml" ] && ! command -v cargo &>/dev/null; then
    echo "  ℹ  Rust project detected but cargo is not installed. Try: gyanesh-help toolkit install rust"
  fi
  if [ -f "go.mod" ] && ! command -v go &>/dev/null; then
    echo "  ℹ  Go project detected but go is not installed. Try: gyanesh-help toolkit install go"
  fi
}

if [ -n "$ZSH_VERSION" ]; then
  autoload -U add-zsh-hook
  add-zsh-hook chpwd _gyanesh_help_chpwd
else
  cd() { builtin cd "$@" && _gyanesh_help_chpwd; }
fi

# Run on shell start for current directory
_gyanesh_help_chpwd`

	case "fish":
		return `# Auto-activate virtual environments on cd
function _gyanesh_help_on_cd --on-variable PWD
  if test -d .venv; and test -f .venv/bin/activate.fish
    source .venv/bin/activate.fish 2>/dev/null
  else if test -d venv; and test -f venv/bin/activate.fish
    source venv/bin/activate.fish 2>/dev/null
  end

  if test -f .gyanesh-help
    set -l tag (grep '^tag' .gyanesh-help 2>/dev/null | head -1 | sed 's/.*= *"\(.*\)"/\1/')
    if test -n "$tag"
      gyanesh-help tag set "$tag" 2>/dev/null
    end
  end
end`

	case "powershell":
		return `# Auto-switch project tag from .gyanesh-help config
function Invoke-GyaneshHelpHook {
  if (Test-Path ".gyanesh-help") {
    $tag = (Get-Content ".gyanesh-help" | Select-String '^tag' | ForEach-Object { $_ -replace '.*= *"(.*)"', '$1' })
    if ($tag) { gyanesh-help tag set $tag 2>$null }
  }
  if (Test-Path ".venv\Scripts\Activate.ps1") {
    & .venv\Scripts\Activate.ps1
  }
}`
	}

	_ = runtime.GOOS // suppress unused
	return ""
}
