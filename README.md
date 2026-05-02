<div align="center">
  <img src="Assets/Cue.png" alt="Cue Hero Banner" width="1200" style="border-radius: 12px; margin-bottom: 20px;">

  <h1>✨ Cue ✨</h1>
  <h3><i>Makes the terminal feel like it already knows what you need.</i></h3>
  
  <p>Queue management, semantic macros, environment stores, smart history, and Claude Code integrations — all offline, all local.</p>

  <p>
    <a href="https://github.com/GyaneshSamanta/cue/releases/latest">
      <img src="https://img.shields.io/github/v/release/GyaneshSamanta/cue?style=for-the-badge&color=00e5ff&labelColor=1e1e2e" alt="Latest Release" />
    </a>
    <a href="https://www.linkedin.com/newsletters/gyanesh-on-product-6979386586404651008/">
      <img src="https://img.shields.io/badge/Newsletter-Subscribe-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white" alt="LinkedIn Newsletter" />
    </a>
    <a href="https://buymeachai.ezee.li/GyaneshOnProduct">
      <img src="https://img.shields.io/badge/Support_Project-Buy_Me_A_Chai-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black" alt="Buy Me A Chai" />
    </a>
    <a href="https://github.com/GyaneshSamanta/cue">
      <img src="https://img.shields.io/github/repo-size/GyaneshSamanta/cue?style=for-the-badge&color=ff00ea&labelColor=1e1e2e" alt="Repo Size" />
    </a>
  </p>
</div>

<br>

<div align="center">
  <h2>📺 Welcome to v2.2</h2>
  <img src="docs/assets/onboarding_demo.gif" alt="Onboarding Demo" width="900" style="border-radius: 12px; box-shadow: 0 10px 30px rgba(0,0,0,0.5);">
  <br><br>
  <h3>🏪 Environment Store Management</h3>
  <img src="Assets/store_demo.gif" alt="Store Demo" width="900" style="border-radius: 12px; box-shadow: 0 10px 30px rgba(0,0,0,0.5);">
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
<summary><b>🐧 Linux / macOS (incl. Arch Linux support)</b></summary>
<br>

**One-liner install (Ubuntu / Debian / Arch / macOS):**
```bash
curl -fsSL https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.sh | bash
```

To pin a specific version:
```bash
curl -fsSL https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.sh | CUE_VERSION=v2.1.0 bash
```

**Arch Linux specifics:** The `install.sh` script works directly on Arch. We recommend having `base-devel`, `curl`, and `sudo` ready (`sudo pacman -S --needed base-devel curl`).

> [!TIP]
> **AUR packaging:** AUR support (`cue-bin`) is currently in development. For now, the direct install is the recommended path.

</details>

<details>
<summary><b>🪟 Windows (PowerShell)</b></summary>
<br>

```powershell
iwr https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.ps1 -useb | iex
```

To pin a specific version:
```powershell
$env:CUE_VERSION = 'v2.1.0'; iwr https://raw.githubusercontent.com/GyaneshSamanta/cue/main/scripts/install.ps1 -useb | iex
```

</details>

<details>
<summary><b>🛠 Build from source</b></summary>
<br>

Requires Go 1.26+.

```bash
git clone https://github.com/GyaneshSamanta/cue && cd cue
go build -o cue .
./cue --version
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

## 🧠 Local Models — Now with Google Gemma 3

Cue v2.2 ships first-class support for **Google Gemma 3** — Google DeepMind's open-weight model family — alongside the full `cue model` command suite. Run powerful AI **entirely offline, for free**, on your own hardware.

### Why Gemma?

Gemma 3 is Google's most capable open-weight series yet. It punches well above its weight class in coding, reasoning, and instruction-following — and because it runs locally through [Ollama](https://ollama.com/download), **your code never leaves your machine.**

### Quick Start

```bash
# 1. Get a hardware-aware recommendation
cue model recommend

# 2. Install Gemma (defaults to 4B — great for 8 GB+ RAM)
cue model gemma

# 3. Or pick a size explicitly
cue model gemma 1b    # ultra-light, instant load
cue model gemma 12b   # high quality, needs ~16 GB RAM
cue model gemma 27b   # best quality, needs ~32 GB RAM

# 4. Point Claude Code at it
cue model use gemma3:4b

# 5. Other useful commands
cue model list                 # see every installed model
cue model benchmark gemma3:4b  # measure tokens/sec on your hardware
cue model pull qwen2.5-coder:7b  # pull any Ollama model by name
```

### Gemma 3 Size Guide

| Tag | Size on disk | RAM needed | Sweet spot |
| :--- | :---: | :---: | :--- |
| `gemma3:1b` | ~0.8 GB | 4 GB | Lightweight laptops, CI environments |
| `gemma3:4b` | ~3.3 GB | 8 GB | **Recommended default** — fast and capable |
| `gemma3:12b` | ~8.1 GB | 16 GB | Strong reasoning and code quality |
| `gemma3:27b` | ~17 GB | 32 GB | Best local quality available |

> [!NOTE]
> Cue requires [Ollama](https://ollama.com/download) as the local model runtime. If it isn't installed, `cue model gemma` will print exact install instructions for your OS (brew / curl / winget) rather than crashing.

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

## 📋 Changelog

### v2.2.0 — May 2026

**New**
- **`cue model` command** — the `internal/model` package was wired up and is now a first-class CLI feature. Run `cue model --help` to see all subcommands.
- **Google Gemma 3 support** — `cue model gemma [1b|4b|12b|27b]` installs Google Gemma via Ollama in one command. Hardware-aware `cue model recommend` now includes all four Gemma sizes.
- **`cue toolkit` command** — install, upgrade, remove, and inspect developer tools with automatic version-manager bootstrapping (e.g. `cue toolkit install node`). Closes the loop from `cue doctor fix --all`, which had been suggesting `cue toolkit install <tool>` when the command didn't yet exist.
- **Release pipeline** — added `.github/workflows/release.yml` that cross-compiles for linux/darwin/windows × amd64/arm64 and uploads binaries to GitHub Releases on every version tag. **This is what was causing the broken `install.sh` / `install.ps1`** — prior releases had no binary assets.
- **CI pipeline** — added `.github/workflows/ci.yml` for build / vet / test checks on every PR.

**Fixed**
- `install.sh` and `install.ps1` now fail with a clear actionable error if no release is found, instead of silently downloading nothing.
- Both scripts accept a `CUE_VERSION=vX.Y.Z` / `$env:CUE_VERSION` env var to pin to a specific release.
- `cue model list/pull/gemma` now prints a friendly install hint (brew / curl / winget) if Ollama is missing, instead of an opaque exec error.

**Improved**
- README install section: explicit Arch Linux guidance, versioned-install snippets, build-from-source instructions.

---

## 🤝 Contributors

Thank you to everyone who has contributed to Cue!

<div align="center">
  <table>
    <tr>
      <td align="center">
        <a href="https://github.com/GyaneshSamanta">
          <img src="https://github.com/GyaneshSamanta.png" width="80" height="80" style="border-radius: 50%;" alt="Gyanesh Samanta" /><br />
          <sub><b>Gyanesh Samanta</b></sub>
        </a><br />
        <sub>Creator & Maintainer</sub>
      </td>
      <td align="center">
        <a href="https://github.com/fuleinist">
          <img src="https://github.com/fuleinist.png" width="80" height="80" style="border-radius: 50%;" alt="Chris Chen" /><br />
          <sub><b>Chris Chen</b></sub>
        </a><br />
        <sub><a href="https://github.com/GyaneshSamanta/cue/pull/16">cue toolkit command</a></sub>
      </td>
    </tr>
  </table>

  <br />

  <p>Want to see your face here? <a href="CONTRIBUTING.md"><b>Contribute to Cue →</b></a></p>
</div>

---

<div align="center">
  <p><b>Built with ❤️ by Gyanesh</b></p>
  <a href="https://buymeachai.ezee.li/GyaneshOnProduct">
    <img src="https://img.shields.io/badge/Support_Project-Buy_Me_A_Chai-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black" alt="Buy Me A Chai" />
  </a>
</div>
