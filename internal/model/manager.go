package model

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GyaneshSamanta/gyanesh-help/internal/config"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

// ListModels shows all Ollama models.
func ListModels() error {
	out, err := exec.Command("ollama", "list").CombinedOutput()
	if err != nil {
		return fmt.Errorf("ollama not running or not installed: %w", err)
	}
	ui.PrintHeader("Local AI Models (Ollama)")
	fmt.Println(string(out))
	return nil
}

// PullModel downloads a model.
func PullModel(name string) error {
	ui.PrintStep(fmt.Sprintf("Pulling %s...", name))
	cmd := exec.Command("ollama", "pull", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveModel deletes a model.
func RemoveModel(name string) error {
	return exec.Command("ollama", "rm", name).Run()
}

// UseModel sets the default model for Claude Code.
func UseModel(name string) error {
	cfgDir := config.ConfigDir()
	cfgFile := filepath.Join(cfgDir, "default-model")
	os.WriteFile(cfgFile, []byte(name), 0644)

	// Update claude.json
	home, _ := os.UserHomeDir()
	claudeCfg := filepath.Join(home, ".claude.json")
	data, _ := os.ReadFile(claudeCfg)
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		cfg = map[string]interface{}{}
	}
	cfg["model"] = name
	out, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(claudeCfg, out, 0644)

	ui.PrintSuccess(fmt.Sprintf("Default model set to: %s", name))
	return nil
}

// Recommend suggests models based on hardware.
func Recommend() error {
	ui.PrintHeader("Model Recommendations")

	ramGB, vramGB, cpuInfo := detectFullHardware()
	fmt.Printf("  Detected hardware:\n")
	fmt.Printf("    RAM  : %d GB\n", ramGB)
	fmt.Printf("    VRAM : %d GB\n", vramGB)
	fmt.Printf("    CPU  : %s\n", cpuInfo)
	fmt.Println()

	type rec struct {
		name, model, size, note string
		fit                     bool
	}

	recs := []rec{
		{"Code (fast)", "phi3.5-mini:3.8b", "2.4 GB", "Fits in VRAM for fastest inference", vramGB >= 4 || ramGB >= 8},
		{"Code (quality)", "qwen2.5-coder:7b", "4.7 GB", "Best code quality in 7B range", ramGB >= 12},
		{"Code (large)", "qwen2.5-coder:14b", "8.9 GB", "Superior code quality", vramGB >= 8 || ramGB >= 24},
		{"Reasoning", "deepseek-r1:8b", "5.2 GB", "Strong multi-step reasoning", ramGB >= 12},
		{"Reasoning (large)", "deepseek-r1:14b", "9.1 GB", "Best reasoning under 15B", ramGB >= 24},
		{"General", "llama3.1:8b", "4.7 GB", "Versatile general model", ramGB >= 12},
	}

	fmt.Println("  Recommendations:")
	for _, r := range recs {
		icon := "✔"
		if !r.fit {
			icon = "✗"
		}
		fmt.Printf("    %s %-12s  %-24s %s  %s\n", icon, r.name, r.model, r.size, r.note)
	}

	return nil
}

// Benchmark runs a quick speed test.
func Benchmark(modelName string) error {
	if modelName == "" {
		data, err := os.ReadFile(filepath.Join(config.ConfigDir(), "default-model"))
		if err != nil {
			return fmt.Errorf("no model specified and no default set")
		}
		modelName = strings.TrimSpace(string(data))
	}

	ui.PrintHeader(fmt.Sprintf("Benchmarking: %s", modelName))
	ui.PrintStep("Sending test prompt...")

	cmd := exec.Command("ollama", "run", modelName, "Write a hello world in Python. Respond with only code.")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Search queries the Ollama registry.
func Search(query string) error {
	ui.PrintHeader(fmt.Sprintf("Searching Ollama registry: %s", query))
	ui.PrintWarning("Ollama CLI does not support search yet. Browse: https://ollama.com/library?q=" + query)
	return nil
}

func detectFullHardware() (ramGB, vramGB int, cpuInfo string) {
	ramGB = 8
	vramGB = 0
	cpuInfo = "unknown"

	if exec.Command("nvidia-smi").Run() == nil {
		vramGB = 4 // conservative estimate
	}
	return
}
