<div align="center">
  <img src="C:\Users\mailg\.gemini\antigravity\brain\9260aef7-b93c-4e47-873e-c27e43460e31\cue_hero_banner_1774850266139.png" alt="Cue Hero Banner" width="1200" style="border-radius: 12px; margin-bottom: 20px;">

  <h1>✨ Cue ✨</h1>
  <h3><i>Makes the terminal feel like it already knows what you need.</i></h3>
  
  <p>Queue management, semantic macros, environment stores, smart history, and Claude Code integrations — all offline, all local.</p>

  <p>
    <a href="https://github.com/GyaneshSamanta/cue/releases/latest">
      <img src="https://img.shields.io/github/v/release/GyaneshSamanta/cue?style=for-the-badge&color=00e5ff&labelColor=1e1e2e" alt="Latest Release" />
    </a>
    <a href="https://github.com/GyaneshSamanta/cue/releases">
      <img src="https://img.shields.io/github/downloads/GyaneshSamanta/cue/total?style=for-the-badge&color=bc00ff&labelColor=1e1e2e" alt="Total Downloads" />
    </a>
    <a href="https://github.com/GyaneshSamanta/cue">
      <img src="https://img.shields.io/github/repo-size/GyaneshSamanta/cue?style=for-the-badge&color=ff00ea&labelColor=1e1e2e" alt="Repo Size" />
    </a>
  </p>
</div>

<br>

<div align="center">
  <h2>📺 Welcome to v2.0</h2>
  <img src="docs/assets/onboarding_demo.gif" alt="Onboarding Demo" width="900" style="border-radius: 12px; box-shadow: 0 10px 30px rgba(0,0,0,0.5);">
</div>

---

## 🚀 The Core Concept

Modern application development requires immense cognitive overhead. Should you use `npm audit fix` or `yarn upgrade`? Did you remember to append `--force-with-lease` on your git push? Which version of python is globally overriding your deep-learning virtual environment? 

**Cue** drastically reduces this overhead. It intercepts complex needs and runs highly-optimized semantic macros, and it provisions complete declarative "Stores" to standardize how you build software. 

### 💡 Why Use Cue?

- 🎨 **Interactive TUI fallbacks:** Never check `--help` pages again. Forget arguments? The system launches a graphical TUI menu directly in the terminal to guide you.
- 📦 **Environment Stores:** Declarative dependencies. Install a full `data-science` stack in a single command.
- 🤖 **Claude Code Orchestration:** Run LLMs cleanly and directly against your file system with a single install flag.
- 🔒 **Safety & Security Checkups:** In-built `audit` commands prevent credential leaks entirely offline.

---

## 🛠 Features & User Guide

### 💿 Installation

The CLI distributes as a self-contained, statically linked fat binary. No complex Node/Python/Ruby dependencies are strictly required to run it!

<details>
<summary><b>🐧 Linux / macOS</b></summary>
<br>

```bash
curl -fsSL https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.sh | bash
```

</details>

<details>
<summary><b>🪟 Windows (PowerShell)</b></summary>
<br>

```powershell
iwr https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.ps1 -useb | iex
```

</details>

<br>

> [!TIP]
> **First Run & Onboarding:**
> Upon your first successful execution of any `cue` command, the CLI will automatically launch a dynamic **Onboarding Wizard**. This setup will gently introduce you to Environment Stores, the Macro Engine, and the LLM toolset. You can replay this tutorial at any time by running:
> ```bash
> cue onboarding
> ```

---

## 🧩 Environment Stores

Stop manually installing toolchains. Use our specialized environments, loaded dynamically depending on what directory you are currently sitting in!

**Install an Environment:**
```bash
cue store install mern
```

#### Explore the Specialized Tech Stacks:
| Stack Type | Identifier | Primary Components |
| :--- | :--- | :--- |
| **[🤖 AI & ML](./docs/ai_ml.md)** | `ai-dev` | Claude Code, Ollama, liteLLm |
| **[📊 Data Science](./docs/data_science.md)** | `data-science` | Python 3.10+, JupyterLab, Poetry |
| **[⚙️ DevOps](./docs/devops.md)** | `devops` | Terraform, K8s, Cloud CLIs, Docker |

---

## 🤖 Generative AI: Claude Code Engine

Cue elegantly orchestrates Anthropic's Claude Code for you. 
```bash
cue claude-code install
```
During installation, it offers multiple execution engines:
1. **API Mode:** Sends code direct to the cloud. Best for reasoning.
2. **Local Mode (Ollama):** Purely local. Pulls models down through `ollama` and proxies them securely. **This implementation is 100% free and extremely private.**

---

## ⚡ The Ultimate Macro Glossary

Macros encapsulate best-practices and safety constraints into readable verbs. Forget raw terminal commands—use our semantic shortcuts. 

| Macro | Category | Purpose |
| :--- | :--- | :--- |
| `cue git-oops-push` | **Git** | Overwrites a remote branch safely. |
| `cue git-undo` | **Git** | Un-does your last commit but keeps files staged. |
| `cue docker-nuke` | **Docker** | Obliterates all containers, volumes, and dangling images. |
| `cue go-mod-tidy-check` | **Go** | Formats, lints, and executes the suite. |
| `cue npm-audit-fix` | **Node.js**| Auto-patches internal Node vulnerabilities. |
| `cue python-venv-here`| **Python**| Scaffolds a virtual environment and prints activation string. |
| `cue ollama-chat` | **AI** | Drops into an optimized REPL terminal chat. |

> _To view these dynamically on your terminal, type `cue macro list`._

---

## 🫂 Developer & Maintainer Documentation

### 🔒 Security Information & Trust

When executing `cue`, **zero telemetry or user data leaves your computer**. Everything operates deterministically over standard POSIX interfaces.
- The `cue audit` command locally analyzes your SSH structures (flagging deprecated RSA signatures over Ed25519 standardizations).
- TUI menus securely scrub inputs locally before running.
- LLM API keys are strictly vaulted in local JSON structures and are never logged internally.

### 🔌 Extending the System

We welcome community contributions! Please view our **[Contributing Guidelines](CONTRIBUTING.md)** to understand the process.

**Standardization Rule:** 
All crashes must use our `ui.StructuredError` package rather than standard `panic` invocations. This system guides users on resolving their own errors dynamically. Example structure:

```go
se := ui.NewStructuredError(
    "Installation Failed",
    "Dependency 'rustc' dropped connection.",
    []string{
        "Check your internet connection",
        "Run 'cue store install rust' again",
    },
    err,
)
ui.HandleError(se)
```

---

<div align="center">
  <p><b>Built with ❤️ by Gyanesh</b></p>
  <a href="https://buymeachai.ezee.li/GyaneshOnProduct">
    <img src="https://img.shields.io/badge/Support_Project-Buy_Me_A_Chai-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black" alt="Buy Me A Chai" />
  </a>
</div>
