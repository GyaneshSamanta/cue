# Gyanesh-help CLI v2.0

<div align="center">
  <h3>Makes the terminal feel like it already knows what you need.</h3>
  <p>Queue management, semantic macros, environment stores, smart history, and Claude Code integrations — all offline, all local.</p>
</div>

---

## The Core Concept

Modern application development requires immense cognitive overhead. Should you use `npm audit fix` or `yarn upgrade`? Did you remember to append `--force-with-lease` on your git push? Which version of python is globally overriding your deep-learning virtual environment? 

**Gyanesh-help** drastically reduces cognitive overhead. It intercepts complex needs and runs highly-optimized semantic macros, and it provisions complete declarative "Stores" to standardize how you build software. 

### Why Use Gyanesh-help?
1. **Interactive TUI fallbacks:** Never check `--help` pages again. If you forget arguments, the system launches a graphical TUI menu directly in the terminal to guide you.
2. **Environment Stores:** Declarative dependencies. Install a full `data-science` stack in a single command.
3. **Claude Code Orchestration:** Run LLMs cleanly and directly against your file system with a single install flag.
4. **Safety & Security Checkups:** In-built `audit` commands prevent credential leaks entirely offline.

---

## Part 1: First-Time Setup & User Guide

### 💿 Installation

The CLI distributes as a self-contained, statically linked fat binary. No Node/Python/Ruby dependencies are strictly required to run it!

**Linux / macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/GyaneshSamanta/gyanesh-help/main/install.sh | bash
```

**Windows** (PowerShell)
```powershell
iwr https://raw.githubusercontent.com/GyaneshSamanta/gyanesh-help/main/install.ps1 -useb | iex
```

#### First Run & Onboarding

Upon your first successful execution of any `gyanesh-help` command, the CLI will automatically launch a dynamic **Onboarding Wizard**. This setup will gently introduce you to Environment Stores, the Macro Engine, and the LLM toolset. 

You can replay this tutorial at any time by running:
```bash
gyanesh-help onboarding
```

### 🧩 Environment Stores

Stop manually installing toolchains. Use our specialized environments, loaded dynamically depending on what directory you are currently sitting in!

**Install an Environment:**
```bash
gyanesh-help store install mern
```

**Explore the Specialized Tech Stacks:**
- **[🤖 AI & Machine Learning Stack](./docs/ai_ml.md)** (`ai-dev`: Claude Code, Ollama, liteLLm)
- **[📊 Data Science Stack](./docs/data_science.md)** (`data-science`: Python 3.10+, JupyterLab, Poetry)
- **[⚙️ DevOps Stack](./docs/devops.md)** (`devops`: Terraform, K8s, Cloud CLIs, Docker)

### 🤖 Generative AI: Claude Code Engine

Gyanesh-help elegantly orchestrates Anthropic's Claude Code for you. 
```bash
gyanesh-help claude-code install
```
During installation, it offers multiple execution engines:
1. **API Mode:** Sends code direct to the cloud. Best for reasoning.
2. **Local Mode (Ollama):** Purely local. Pulls models down through `ollama` and proxies them securely. **This implementation is 100% free and extremely private.**

---

## Part 2: The Ultimate Macro Glossary

Macros encapsulate best-practices and safety constraints into readable verbs. Forget raw terminal commands—use our semantic shortcuts. 

| Macro Identifier | Category | Underlying Bash / Execution | Definition & Result |
| :--- | :--- | :--- | :--- |
| `git-oops-push` | **Git** | `git push --force-with-lease` | Overwrites a remote branch safely by ensuring you don't delete colleagues' undocumented commits. |
| `git-undo` | **Git** | `git reset --soft HEAD~1` | Un-does your last commit but gracefully keeps all files staged. |
| `git-diff-staged` | **Git** | `git diff --cached` | Shows exactly what code changes are staged for the next commit. |
| `git-pr` | **Git** | `gh pr create --web` | Uses the GitHub CLI to instantly scaffold a Pull Request in the browser. |
| `git-sync` | **Git** | `git pull --rebase ... && git push` | Fetches main, rebases local, and pushes linearly without merge commits. |
| `git-contributors`| **Git** | `git shortlog -sn --no-merges` | Prints a leaderboard of contributors based on pure commit frequency. |
| `docker-nuke` | **Docker** | `docker system prune -af --volumes` | Absolutely obliterates all containers, volumes, and dangling images to aggressively save disk space. Prompts for confirmation. |
| `docker-shell` | **Docker** | `docker exec -it <id> bash` | Instantly drops you into a bash terminal inside a running container. |
| `docker-compose-restart`| **Docker**| `docker-compose down && ... up -d`| Gracefully hard-restarts all mapped services in the active compose file. |
| `nuke-docker-volume` | **Docker**| `docker volume rm $(... dangling=true)` | Safely removes only orphaned, dangling volumes to optimize space. |
| `go-mod-tidy-check` | **Go** | `go mod tidy && go vet && go test` | The holy-trinity checker for Go code. Formats, lints, and executes the suite. |
| `cargo-release`| **Rust** | `cargo clippy -D warnings && ...`| Bulletproof release compiler that fails instantly if lint warnings exist. |
| `npm-audit-fix`| **Node.js**| `npm audit fix --force` | Auto-patches internal Node vulnerabilities. |
| `node-version-check` | **Node.js**| `node -v && npm -v` | Instantly validates local runtime versions against package constraints. |
| `tf-plan-clean`| **Terraform**| `terraform init && ... validate && ...` | Formats and validates HCL infrastructure code before speculatively planning it. |
| `k8s-pod-shell` | **K8s** | `kubectl exec -it ...` | Executes into an alpine pod. Falls back to `sh` if `bash` isn't found. |
| `k8s-logs` | **K8s** | `kubectl logs -f ... --tail=100` | Safely tails cluster logs without overflowing the buffer length. |
| `python-venv-here`| **Python**| `python3 -m venv .venv` | Scaffolds a virtual environment and prints OS-specific activation strings. |
| `pip-freeze-clean`| **Python**| `pip freeze > requirements.txt` | Dumps a deterministic lockfile representing your current data-science dependencies. |
| `ollama-list` | **AI** | `ollama list` | View all local neural networks and their memory footprint. |
| `ollama-chat` | **AI** | `ollama run <model>` | Instantly drops into an optimized REPL terminal chat. |

_To view these dynamically on your terminal, just type `gyanesh-help macro list`._

---

## Part 3: Developer & Maintainer Documentation

### 🔒 Security Information & Trust

When executing `gyanesh-help`, **zero telemetry or user data leaves your computer**. Everything operates deterministically over standard POSIX interfaces.
- The `gyanesh-help audit` command locally analyzes your SSH structures (flagging deprecated RSA signatures over Ed25519 standardizations).
- TUI menus securely scrub inputs locally before running.
- LLM API keys (if using API Mode) are strictly vaulted in local JSON structures and are never logged internally in the `.gyanesh-help/exports` logs.

### Extending the System

To build new stores or macros, please submit a Pull Request altering `internal/macro/builtins/*.go` or `internal/store/stacks/*.go`.

**Standardization Rule:** 
All crashes must use our `ui.StructuredError` package rather than standard `panic` invocations. This system guides users on resolving their own errors dynamically. Example structure:

```go
se := ui.NewStructuredError(
    "Installation Failed",
    "Dependency 'rustc' dropped connection.",
    []string{
        "Check your internet connection",
        "Run 'gyanesh-help store install rust' again",
    },
    err, // Original underlying error trace
)
ui.HandleError(se)
```

---

<p align="center">
  <b>Built with <3 by Gyanesh</b> <br>
  Support the project and see more tooling: <a href="https://buymeachai.ezee.li/GyaneshOnProduct">https://buymeachai.ezee.li/GyaneshOnProduct</a>
</p>
