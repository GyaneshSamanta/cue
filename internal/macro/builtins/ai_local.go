package builtins

import "github.com/GyaneshSamanta/gyanesh-help/internal/macro"

func init() {
	registerAILocalMacros()
}

func registerAILocalMacros() {
	macro.Register(&macro.Macro{
		Name:        "ollama-list",
		Command:     "ollama list",
		Description: "List all pulled + running Ollama models",
		Explanation: `Shows all AI models downloaded to your machine via Ollama.
Includes model name, size, and when it was last modified.
To pull a new model: ollama pull <model-name>
Browse available models: https://ollama.com/library`,
		Dangerous: false,
	})

	macro.Register(&macro.Macro{
		Name:        "ollama-chat",
		Command:     "ollama run {{.model}}",
		Description: "Start an interactive chat with a local AI model",
		Explanation: `Opens an interactive chat session with a locally running Ollama model.
Type your messages and press Enter. Use /bye to exit.
Popular models: qwen2.5-coder:7b, deepseek-r1:8b, llama3.1:8b`,
		Dangerous: false,
		Flags:     []macro.Flag{{Name: "model", Description: "Model name", Default: "qwen2.5-coder:7b"}},
	})
}
