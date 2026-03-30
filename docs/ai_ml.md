# 🚀 AI & ML Stack

The AI & Machine Learning environment stack is the crown jewel of `cue`. Designed for modern ML engineers and AI application developers, this stack ensures your local environment is configured with the absolute best-in-class tools for running, fine-tuning, and interacting with Large Language Models.

## The Environment Store: `ai-dev`

Running `cue store install ai-dev` instantly provisions your machine with:
- **Ollama**: The standard for running inference locally.
- **liteLLM**: Universal API translation proxy.
- **Claude Code**: Direct CLI integration.
- **CUDA/RocM Hooks**: Automatic environment parsing for hardware acceleration.
- **HuggingFace CLI**: For pulling models directly from the Hub.

## Claude Code Orchestration

Cue seamlessly acts as an orchestrator for **Claude Code**.
Run:
```bash
cue claude-code install
```
You can choose your setup mode:
- **API Mode:** Securely connects directly to Anthropic's APIs. Sets up `promptfoo` and an internal MCP.
- **Local Mode:** Spins up Ollama, routes it through liteLLm proxies, and ties it into Claude CLI. **Zero cloud dependency!**
- **Hybrid Mode:** Combines both capabilities.

## Dedicated AI Macros

To streamline local model development, these macros are pre-configured:

### 1. `ollama-list`
- **Command:** `ollama list`
- **What it does:** View all locally pulled models and their exact memory footprints.

### 2. `ollama-chat`
- **Command:** `ollama run $1` (e.g., `cue ollama-chat --model llama3`)
- **What it does:** Instantly drop into a highly-optimized REPL for terminal chatting with an AI model.

---
*Ready to build the future? Start by running `cue store install ai-dev`.*
