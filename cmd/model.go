package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/model"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage local AI models (Ollama, Gemma, etc.)",
	Long: `Manage local AI models served by Ollama, including Google Gemma.

Cue treats Ollama as the canonical local-model runtime. Use this command
to pull models, set the active model for Claude Code, benchmark inference
speed, or get hardware-aware recommendations (including Gemma sizes).`,
}

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed local models",
	RunE: func(cmd *cobra.Command, args []string) error {
		return model.ListModels()
	},
}

var modelPullCmd = &cobra.Command{
	Use:   "pull <name>",
	Short: "Download a model (e.g. gemma3:4b, llama3.1:8b)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return model.PullModel(args[0])
	},
}

var modelRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a downloaded model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return model.RemoveModel(args[0])
	},
}

var modelUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the default model used by Claude Code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return model.UseModel(args[0])
	},
}

var modelRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Suggest models based on detected hardware",
	RunE: func(cmd *cobra.Command, args []string) error {
		return model.Recommend()
	},
}

var modelBenchmarkCmd = &cobra.Command{
	Use:   "benchmark [name]",
	Short: "Run a quick inference speed test",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) == 1 {
			name = args[0]
		}
		return model.Benchmark(name)
	},
}

var modelSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search the Ollama registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return model.Search(args[0])
	},
}

var modelGemmaCmd = &cobra.Command{
	Use:   "gemma [size]",
	Short: "Quick-install Google Gemma (sizes: 1b, 4b, 12b, 27b — default 4b)",
	Long: `Quick-install Google Gemma 3 via Ollama.

Examples:
  cue model gemma          # installs gemma3:4b (good default for 16GB RAM)
  cue model gemma 1b       # tiny — fits low-RAM laptops
  cue model gemma 12b      # quality — needs ~16GB RAM or 12GB VRAM
  cue model gemma 27b      # best — needs ~32GB RAM or 24GB VRAM

After install, run 'cue model use gemma3:<size>' to make it the default
for Claude Code.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		size := "4b"
		if len(args) == 1 {
			size = args[0]
		}
		valid := map[string]bool{"1b": true, "4b": true, "12b": true, "27b": true}
		if !valid[size] {
			ui.PrintWarning("Unknown Gemma size '" + size + "'. Supported: 1b, 4b, 12b, 27b. Proceeding anyway.")
		}
		tag := "gemma3:" + size
		if err := model.EnsureOllama(); err != nil {
			return err
		}
		return model.PullModel(tag)
	},
}

func init() {
	modelCmd.AddCommand(modelListCmd, modelPullCmd, modelRemoveCmd, modelUseCmd,
		modelRecommendCmd, modelBenchmarkCmd, modelSearchCmd, modelGemmaCmd)
}
