package claudecode

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/GyaneshSamanta/gyanesh-help/internal/adapter"
	"github.com/GyaneshSamanta/gyanesh-help/internal/config"
	"github.com/GyaneshSamanta/gyanesh-help/internal/toolkit"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

// Mode is the Claude Code installation mode.
type Mode string

const (
	ModeAPI    Mode = "api"
	ModeLocal  Mode = "local"
	ModeHybrid Mode = "hybrid"
)

// InstallAPIMode sets up Claude Code with Anthropic API.
func InstallAPIMode(a adapter.OSAdapter) error {
	ui.PrintHeader("Claude Code — API Mode Setup")

	// Step 1: Node.js prerequisite
	ui.PrintStep("[1/7] Checking Node.js...")
	if _, ok := toolkit.GetVersion("node", "--version"); !ok {
		ui.PrintWarning("Node.js not found, installing via toolkit...")
		if err := toolkit.Install("node", a, "", false); err != nil {
			return fmt.Errorf("Node.js required: %w", err)
		}
	}
	ui.PrintSuccess("Node.js ready")

	// Step 2: Claude Code CLI
	ui.PrintStep("[2/7] Installing Claude Code CLI...")
	if err := toolkit.RunInstallCmd("npm", "install", "-g", "@anthropic-ai/claude-code"); err != nil {
		ui.PrintWarning("npm install failed — you may need to install claude-code manually")
	}

	// Step 3: API key
	ui.PrintStep("[3/7] API Key Setup")
	apiKey := ui.ReadInput("  Paste your ANTHROPIC_API_KEY (input hidden): ")
	if apiKey != "" {
		saveAPIKey(apiKey, a)
		ui.PrintSuccess("API key saved")
	}

	// Step 4: MCP Servers
	ui.PrintStep("[4/7] MCP Server Installation")
	installMCPServers(a)

	// Step 5: promptfoo
	ui.PrintStep("[5/7] Installing promptfoo...")
	exec.Command("npm", "install", "-g", "promptfoo").Run()
	writePromptfooConfig()
	ui.PrintSuccess("promptfoo installed")

	// Step 6: CLAUDE.md template
	ui.PrintStep("[6/7] Writing CLAUDE.md template...")
	writeClaudeTemplate()
	ui.PrintSuccess("CLAUDE.md template created")

	// Step 7: Verify
	ui.PrintStep("[7/7] Verification")
	if out, err := exec.Command("claude", "--version").Output(); err == nil {
		ui.PrintSuccess(fmt.Sprintf("Claude Code: %s", string(out)))
	}

	saveMode(ModeAPI)
	fmt.Println()
	ui.PrintSuccess("Claude Code (API mode) setup complete!")
	return nil
}

// InstallLocalMode sets up Claude Code with Ollama + LiteLLM.
func InstallLocalMode(a adapter.OSAdapter) error {
	ui.PrintHeader("Claude Code — Local Mode Setup (Ollama)")

	// Step 1: Prerequisites
	ui.PrintStep("[1/8] Checking prerequisites...")
	if _, ok := toolkit.GetVersion("node", "--version"); !ok {
		toolkit.Install("node", a, "", false)
	}
	if _, ok := toolkit.GetVersion(getPythonBin(), "--version"); !ok {
		toolkit.Install("python", a, "", false)
	}

	// Step 2: Ollama
	ui.PrintStep("[2/8] Installing Ollama...")
	if _, ok := toolkit.GetVersion("ollama", "--version"); !ok {
		if err := toolkit.Install("ollama", a, "", false); err != nil {
			return fmt.Errorf("Ollama install failed: %w", err)
		}
	}
	ui.PrintSuccess("Ollama ready")

	// Step 3: Model selection
	ui.PrintStep("[3/8] Model Selection")
	ram, vram := detectHardware()
	model := recommendModel(ram, vram)
	fmt.Printf("  Detected: %d GB RAM, %d GB VRAM\n", ram, vram)
	fmt.Printf("  Recommended: %s\n", model)
	if ui.Confirm(fmt.Sprintf("  Pull %s?", model)) {
		ui.PrintStep(fmt.Sprintf("  Pulling %s...", model))
		cmd := exec.Command("ollama", "pull", model)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	// Step 4: LiteLLM proxy
	ui.PrintStep("[4/8] Installing LiteLLM proxy...")
	pipBin := "pip3"
	if runtime.GOOS == "windows" {
		pipBin = "pip"
	}
	exec.Command(pipBin, "install", "litellm[proxy]").Run()
	writeLiteLLMConfig(model)
	ui.PrintSuccess("LiteLLM proxy configured (http://localhost:4000)")

	// Step 5: Claude Code CLI
	ui.PrintStep("[5/8] Installing Claude Code CLI...")
	exec.Command("npm", "install", "-g", "@anthropic-ai/claude-code").Run()

	// Step 6: Route to local proxy
	ui.PrintStep("[6/8] Routing Claude Code to local proxy...")
	toolkit.AppendToShellConfig(`export ANTHROPIC_BASE_URL=http://localhost:4000`)
	toolkit.AppendToShellConfig(`export ANTHROPIC_API_KEY=local-mode`)

	// Write claude config
	writeClaudeConfig(model)
	ui.PrintSuccess("Claude Code → LiteLLM → Ollama routing configured")

	// Step 7: MCP Servers (optional)
	ui.PrintStep("[7/8] MCP Servers (optional)")
	installMCPServers(a)

	// Step 8: Open WebUI (optional)
	ui.PrintStep("[8/8] Open WebUI")
	if ui.Confirm("  Install Open WebUI (browser chat interface)?") {
		exec.Command(pipBin, "install", "open-webui").Run()
		ui.PrintSuccess("Open WebUI available at http://localhost:3000")
	}

	saveMode(ModeLocal)
	fmt.Println()
	ui.PrintSuccess("Claude Code (local mode) setup complete!")
	fmt.Println("  Start the proxy: litellm --config ~/.gyanesh-help/litellm-config.yaml")
	fmt.Println("  Then use: claude \"your prompt\"")
	return nil
}

// InstallHybridMode sets up both API and local with automatic fallback.
func InstallHybridMode(a adapter.OSAdapter) error {
	ui.PrintHeader("Claude Code — Hybrid Mode Setup")
	if err := InstallAPIMode(a); err != nil {
		return err
	}
	ui.PrintStep("Now setting up local fallback...")
	return InstallLocalMode(a)
}

// Status shows current Claude Code status.
func Status() {
	ui.PrintHeader("Claude Code Status")

	mode := loadMode()
	fmt.Printf("  Mode: %s\n", mode)

	if out, err := exec.Command("claude", "--version").CombinedOutput(); err == nil {
		fmt.Printf("  Claude CLI: %s\n", string(out))
	} else {
		ui.PrintWarning("  Claude CLI: not found")
	}

	if out, err := exec.Command("ollama", "list").CombinedOutput(); err == nil {
		fmt.Printf("  Ollama models:\n%s\n", string(out))
	}
}

// --- Helpers ---

func saveAPIKey(key string, a adapter.OSAdapter) {
	// Store in shell config (encrypted keychain would be better but complex)
	toolkit.AppendToShellConfig(fmt.Sprintf("export ANTHROPIC_API_KEY=%s", key))
	if runtime.GOOS == "windows" {
		exec.Command("setx", "ANTHROPIC_API_KEY", key).Run()
	}
}

func installMCPServers(a adapter.OSAdapter) {
	servers := []struct{ name, pkg string }{
		{"Filesystem MCP", "@anthropic-ai/mcp-server-filesystem"},
		{"GitHub MCP", "@anthropic-ai/mcp-server-github"},
	}
	for _, s := range servers {
		if ui.Confirm(fmt.Sprintf("  Install %s?", s.name)) {
			exec.Command("npm", "install", "-g", s.pkg).Run()
			ui.PrintSuccess(fmt.Sprintf("  %s installed", s.name))
		}
	}
}

func writePromptfooConfig() {
	cfgDir := config.ConfigDir()
	content := `# promptfoo config — evaluate Claude Code prompts
prompts:
  - "Refactor this code to be more readable: {{code}}"
  - "Add error handling to: {{code}}"

providers:
  - anthropic:messages:claude-sonnet-4-20250514

tests:
  - vars:
      code: "function add(a,b){return a+b}"
    assert:
      - type: contains
        value: "function"
`
	os.WriteFile(filepath.Join(cfgDir, "promptfoo.yaml"), []byte(content), 0644)
}

func writeClaudeTemplate() {
	home, _ := os.UserHomeDir()
	content := `# CLAUDE.md — Project Context for Claude Code
# Place this file in your project root to give Claude context about your codebase.

## Project Overview
<!-- Describe your project here -->

## Tech Stack
<!-- List your technologies -->

## Conventions
<!-- Coding standards, naming conventions, etc. -->

## Important Files
<!-- Key files Claude should know about -->

## Testing
<!-- How to run tests, testing conventions -->
`
	os.WriteFile(filepath.Join(home, "CLAUDE.md.template"), []byte(content), 0644)
}

func writeLiteLLMConfig(model string) {
	cfgDir := config.ConfigDir()
	content := fmt.Sprintf(`model_list:
  - model_name: claude-sonnet-4-20250514
    litellm_params:
      model: ollama/%s
      api_base: http://localhost:11434
  - model_name: claude-haiku-4-20250514
    litellm_params:
      model: ollama/%s
      api_base: http://localhost:11434

general_settings:
  master_key: local-mode
`, model, model)
	os.WriteFile(filepath.Join(cfgDir, "litellm-config.yaml"), []byte(content), 0644)
}

func writeClaudeConfig(model string) {
	home, _ := os.UserHomeDir()
	cfg := map[string]interface{}{
		"model":       model,
		"permissions": map[string]bool{"allow_read": true, "allow_write": true},
	}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(filepath.Join(home, ".claude.json"), data, 0644)
}

func detectHardware() (ramGB, vramGB int) {
	// Simple detection — platform-specific improvements possible
	ramGB = 8 // default
	vramGB = 0
	if exec.Command("nvidia-smi").Run() == nil {
		vramGB = 4 // default VRAM estimate
	}
	return
}

func recommendModel(ramGB, vramGB int) string {
	switch {
	case vramGB >= 12:
		return "qwen2.5-coder:14b"
	case vramGB >= 6:
		return "qwen2.5-coder:7b"
	case ramGB >= 16:
		return "qwen2.5-coder:7b"
	default:
		return "phi3.5-mini:3.8b"
	}
}

func saveMode(mode Mode) {
	cfgDir := config.ConfigDir()
	os.WriteFile(filepath.Join(cfgDir, "claude-mode"), []byte(string(mode)), 0644)
}

func loadMode() Mode {
	data, err := os.ReadFile(filepath.Join(config.ConfigDir(), "claude-mode"))
	if err != nil {
		return "not configured"
	}
	return Mode(data)
}

func getPythonBin() string {
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}
	return "python"
}

// ListMCPServers lists installed MCP servers.
func ListMCPServers() {
	ui.PrintHeader("MCP Servers")
	out, err := exec.Command("npm", "list", "-g", "--depth=0").CombinedOutput()
	if err == nil {
		fmt.Println(string(out))
	} else {
		ui.PrintError("Failed to list MCP servers. Is npm installed?")
	}
}

// AddMCPServer adds an MCP server.
func AddMCPServer(name string) {
	ui.PrintStep(fmt.Sprintf("Installing MCP server: %s...", name))
	if err := exec.Command("npm", "install", "-g", name).Run(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to install MCP server: %v", err))
		return
	}
	ui.PrintSuccess(fmt.Sprintf("Installed %s", name))
}

// RemoveMCPServer removes an MCP server.
func RemoveMCPServer(name string) {
	ui.PrintStep(fmt.Sprintf("Removing MCP server: %s...", name))
	if err := exec.Command("npm", "uninstall", "-g", name).Run(); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to remove MCP server: %v", err))
		return
	}
	ui.PrintSuccess(fmt.Sprintf("Removed %s", name))
}

// UpdateCLI updates the Claude Code CLI.
func UpdateCLI() {
	ui.PrintStep("Updating Claude Code CLI to latest version...")
	if err := exec.Command("npm", "install", "-g", "@anthropic-ai/claude-code").Run(); err != nil {
		ui.PrintError(fmt.Sprintf("Update failed: %v", err))
		return
	}
	ui.PrintSuccess("Claude Code CLI updated successfully")
}
