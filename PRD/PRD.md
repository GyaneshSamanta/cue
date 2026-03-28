# Product Requirements Document (PRD)
## `gyanesh-help` — Cross-Platform CLI Developer Utility
**Version:** 1.0-draft  
**Status:** For Review  
**Author:** Gyanesh Samanta
**Last Updated:** 2026-03-28

---

## Table of Contents

1. [Executive Summary & Problem Statement](#1-executive-summary--problem-statement)
2. [Target Audience](#2-target-audience)
3. [Goals & Non-Goals (v1.0)](#3-goals--non-goals-v10)
4. [Detailed User Stories](#4-detailed-user-stories)
5. [Feature Specifications](#5-feature-specifications)
6. [Success Metrics](#6-success-metrics)
7. [Out of Scope — v1.0](#7-out-of-scope--v10)
8. [Appendix](#8-appendix)

---

## 1. Executive Summary & Problem Statement

### 1.1 Executive Summary

`gyanesh-help` is a **zero-dependency, ultra-lightweight command-line interface (CLI) utility** that acts as a universal developer productivity layer across Windows, macOS, and Linux. It replaces the need for a developer to memorise OS-specific package manager syntax, manually handle installation lock conflicts, babysit large downloads over flaky networks, or reverse-engineer complex git or shell operations from documentation.

The tool operates **entirely offline and locally**, using hardcoded, rules-based logic. It has no LLM runtime, no telemetry, and no cloud dependency in its core execution path.

### 1.2 Problem Statement

Developer time is chronically wasted at the seams — not during deep work, but before and after it. Five specific frictions account for the majority of lost time during environment setup and daily terminal use:

| # | Pain Point | Manifestation |
|---|-----------|---------------|
| 1 | **Lock conflicts** | `dpkg` or `apt` holds a lock; a second install attempt crashes rather than queues. On Windows, `msiexec` serialisation failures are cryptic and silent. |
| 2 | **Network fragility** | A 2 GB CUDA toolkit download silently fails at 94% due to a Wi-Fi drop. The user must restart from zero. |
| 3 | **Syntax opacity** | Commands like `git reset --hard HEAD~1` or `find . -name "*.log" -mtime +7 -delete` are correct but unmemorable. Developers spend more time searching Stack Overflow than coding. |
| 4 | **Environment setup cost** | Setting up a full MERN or Data Science stack on a new machine routinely takes 2–4 hours of manual, error-prone installation across multiple tools. |
| 5 | **History & context loss** | Terminal history is a flat, unsearchable dump with no project tagging, no annotation, and no portability across machines. |

### 1.3 Product Vision Statement

> `gyanesh-help` makes the terminal feel like it already knows what you need — queuing what it can't do yet, recovering what would otherwise break, explaining what it just did, and remembering everything it has ever done for you.

---

## 2. Target Audience

### 2.1 Primary Users

| Persona | Description | Primary Use Cases |
|---------|-------------|-------------------|
| **Student Developer** | CS/engineering students, bootcamp learners; comfortable with basic terminal but not OS internals. | Environment Stores, Semantic Macros, History |
| **Solo Developer / Freelancer** | Works across multiple client projects on a single machine, often switching stacks. | Workspace Backup/Sync, Queue Manager, Macro Explainers |
| **Data Scientist** | Python/R-first; often on constrained corporate laptops; uncomfortable with system administration. | Data Science Store, Pause/Resume for large conda installs |
| **Backend / DevOps Engineer** | Power user; values automation over hand-holding but hates repetitive setup rituals. | Environment Stores, Queue Manager, Backup/Sync |
| **AI/ML Practitioner** | Needs CUDA, HuggingFace, and ML library stacks; routinely hits driver/version conflict issues. | AI Dev Store, Pause/Resume, Queue Manager |

### 2.2 Secondary Audience

- Technical educators and workshop facilitators who need repeatable student environment setups.
- Open-source maintainers writing `CONTRIBUTING.md` who want to link to a single `gyanesh-help store install frontend` command instead of maintaining a multi-page setup guide.

### 2.3 Audience Constraints & Assumptions

- Users have a basic familiarity with opening a terminal / command prompt.
- Users may **not** have administrator/sudo access on corporate machines (the tool must gracefully degrade and inform rather than silently fail).
- Users operate on any of: Windows 10/11, macOS 12+, Ubuntu/Debian-family Linux, Fedora/RHEL-family Linux, Arch Linux.

---

## 3. Goals & Non-Goals (v1.0)

### 3.1 Goals

- **G1:** Detect and resolve package manager lock conflicts automatically without user intervention.
- **G2:** Pause and resume large installations across network interruptions without restarting from zero.
- **G3:** Provide a library of ≥30 semantic macro commands with inline human-readable explanations.
- **G4:** Deliver 7 pre-built, tested, cross-platform Environment Stores.
- **G5:** Maintain a queryable, tagged local command history database.
- **G6:** Enable workspace state capture, GitHub backup, and 1-click restore.
- **G7:** Compile to a single binary or script bundle under **15 MB** with zero mandatory cloud calls at runtime.

### 3.2 Non-Goals (v1.0)

- No GUI or web dashboard.
- No LLM or ML-based command suggestion.
- No paid tier, subscription, or telemetry in v1.0.
- No management of remote servers or SSH orchestration.
- No support for containerised dev environments (Docker-in-Docker, Kubernetes) beyond Docker client installation within the Backend Store.

---

## 4. Detailed User Stories

### 4.1 Cross-Platform Queue Management System

**US-Q1:** As a developer on Linux, I want `gyanesh-help install vim` to automatically wait if `apt` is locked by an ongoing system update, so that I don't have to manually retry or decipher lock error messages.

- **Acceptance Criteria:**
  - When a lock is detected, the tool prints: `[QUEUED] Package manager is busy. Your command will run automatically when it's free. Press Ctrl+C to cancel.`
  - A non-blocking background process polls the lock state every configurable interval (default: 5 seconds).
  - The command executes without user re-intervention once the lock releases.
  - If the user cancels, the queued command is removed and a cancellation message is shown.
  - Elapsed wait time is displayed in the terminal in real time.

**US-Q2:** As a developer on Windows, I want `gyanesh-help install git` to detect an active `msiexec` or `winget` process and queue my installation, so that I don't see cryptic "another installation is in progress" errors.

- **Acceptance Criteria:**
  - Detects active `msiexec.exe` processes via WMI/tasklist query before attempting `winget`.
  - Queued commands survive a terminal window resize without dropping.
  - On success, a desktop notification (Windows Toast, if available) is shown: `[gyanesh-help] vim installed successfully.`

**US-Q3:** As a developer on macOS, I want `gyanesh-help install wget` to wait if Homebrew's lockfile is held, so that running two install commands in separate tabs doesn't crash either.

- **Acceptance Criteria:**
  - Detects `~/.homebrew` lockfile or active `brew` processes.
  - Lock polling is adaptive: starts at 3s intervals, backs off to 15s after 2 minutes of waiting.
  - Lock wait timeout is configurable in `~/.gyanesh-help/config.toml` (default: 30 minutes).

---

### 4.2 Resilient Pause/Resume System

**US-R1:** As a data scientist, I want large conda or pip package installations to automatically pause when my Wi-Fi disconnects and resume when connectivity is restored, so that I don't lose hours of download progress.

- **Acceptance Criteria:**
  - The tool wraps the underlying package manager call with a network monitor thread.
  - On disconnect detection (ICMP ping to 1.1.1.1 or 8.8.8.8 fails 3 consecutive times), the installation process is sent a SIGSTOP (Linux/macOS) or suspended (Windows).
  - A visible message is shown: `[PAUSED] Network lost. Installation paused. Waiting for connectivity...`
  - Connectivity is polled every 10 seconds.
  - On restore, the process receives SIGCONT (Linux/macOS) or is resumed (Windows).
  - If the underlying package manager does not support mid-download resume (e.g., older pip versions), the tool falls back to restarting from a cached checkpoint, and warns the user.

**US-R2:** As a developer, I want to manually pause and resume any `gyanesh-help`-managed installation using `Ctrl+Z` / `gyanesh-help resume`, so that I can defer a large install during a video call.

- **Acceptance Criteria:**
  - `Ctrl+Z` within a `gyanesh-help`-managed session suspends the underlying child process.
  - `gyanesh-help resume` lists all paused jobs and allows selection for resumption.
  - Paused job state is persisted to `~/.gyanesh-help/jobs.json` so it survives terminal closure.

---

### 4.3 Semantic Macro Assistant & Explainers

**US-M1:** As a git user, I want to run `gyanesh-help git-undo` to safely revert my last commit, and I want the tool to explain in plain English what it just did, so that I learn and don't feel anxious about data loss.

- **Acceptance Criteria:**
  - `gyanesh-help git-undo` runs `git reset --soft HEAD~1` by default.
  - The tool prints the actual git command executed.
  - It then prints a hardcoded explanation block:
    ```
    ✔ Done. Here's what happened:
    Your last commit was "undone," but your file changes are SAFE and still staged.
    The commit message is gone, but your work is not. You can re-commit when ready.
    This is the safe version of undo. Your history was rewritten locally only.
    ```
  - A `--hard` flag variant (`gyanesh-help git-undo --hard`) executes `git reset --hard HEAD~1` with a destructive-action warning and Y/N confirmation prompt.

**US-M2:** As a developer, I want to run `gyanesh-help explain <command>` to get a plain-English breakdown of any command in the macro library, before executing it, so that I stay in control.

- **Acceptance Criteria:**
  - `gyanesh-help explain git-undo` prints the full explanation without executing anything.
  - `gyanesh-help explain --list` prints all available macros with one-line descriptions.

**US-M3:** As a developer, I want to define my own custom macros with `gyanesh-help macro add`, so that I can extend the library with my own project-specific shortcuts.

- **Acceptance Criteria:**
  - `gyanesh-help macro add <name> "<command string>" "<explanation text>"` persists the macro to `~/.gyanesh-help/macros.toml`.
  - Custom macros are callable exactly like built-in macros.
  - `gyanesh-help macro list` shows both built-in (tagged `[built-in]`) and custom (tagged `[custom]`) macros.
  - Custom macros can be exported and shared as a `.toml` file.

---

### 4.4 Environment Stores

**US-E1:** As a new data science student, I want to run a single command to set up a complete Python/R/Jupyter data science environment, so that I can go from zero to coding in under 20 minutes.

- **Acceptance Criteria:**
  - `gyanesh-help store install data-science` installs all required toolchain components (see Section 5.4).
  - Installation is OS-aware and routes to the correct package manager.
  - A progress bar shows per-component installation status.
  - On completion, a `gyanesh-help store verify data-science` command checks that all components are reachable on PATH and prints a pass/fail table.

**US-E2:** As a developer, I want to preview what a store will install before committing, so that I can make informed decisions on shared or corporate machines.

- **Acceptance Criteria:**
  - `gyanesh-help store preview <stack>` prints the full list of tools, versions, and estimated download size.
  - No installations occur during a `preview` invocation.

**US-E3:** As a developer, I want to uninstall a full environment store in one command, so that I can cleanly remove a stack I no longer need.

- **Acceptance Criteria:**
  - `gyanesh-help store remove <stack>` attempts uninstallation of all components installed by that store.
  - Components shared with other stores are flagged and skipped unless `--force` is passed.

---

### 4.5 Smart History Maintenance

**US-H1:** As a developer, I want every command I run through `gyanesh-help` to be saved with timestamp, project tag, and exit code, so that I can query my history intelligently.

- **Acceptance Criteria:**
  - Every execution writes a record to `~/.gyanesh-help/history.db` (SQLite).
  - Records include: `id`, `timestamp`, `command`, `exit_code`, `duration_ms`, `project_tag`, `stack_context`.
  - `gyanesh-help history` shows the last 20 entries in a formatted table.
  - `gyanesh-help history --tag mern` filters entries to the `mern` project tag.
  - `gyanesh-help history --search "docker"` performs a substring search on the `command` field.

**US-H2:** As a developer, I want to tag the current working session with a project name, so that subsequent commands are logically grouped.

- **Acceptance Criteria:**
  - `gyanesh-help tag set <project-name>` sets the active session tag.
  - `gyanesh-help tag clear` removes the tag.
  - The active tag is shown in the `gyanesh-help status` output.

---

### 4.6 Workspace Backup & Sync

**US-B1:** As a developer moving to a new machine, I want to back up my current toolchain state and config files to a private GitHub repo with one command, so that I can restore my setup in minutes.

- **Acceptance Criteria:**
  - `gyanesh-help workspace backup` captures the current environment state (see Section 5.7) and pushes to a new or existing private GitHub repo.
  - Requires a one-time `gyanesh-help workspace auth --token <PAT>` to store a GitHub Personal Access Token locally (encrypted at rest using OS keychain).
  - The backup commit message is timestamped: `gyanesh-help backup: 2026-03-28T14:32:00Z`.
  - Sensitive files (`.env`, private keys) are excluded by a built-in `.gitignore` template.

**US-B2:** As a developer, I want to restore my workspace from my GitHub backup on a new machine with a single command, so that I don't spend hours manually reinstalling and configuring tools.

- **Acceptance Criteria:**
  - `gyanesh-help workspace restore --repo <github-url>` clones the backup repo, reads the manifest, and runs the appropriate store installs and config file placements.
  - The restore process shows a step-by-step progress log.
  - On conflict (a tool is already installed at a different version), the user is prompted.

---

## 5. Feature Specifications

### 5.1 Cross-Platform Queue Management

The queue system wraps every install/system command issued through `gyanesh-help`. It operates via a pre-execution hook that checks for active locks **before** spawning the underlying command.

**Lock Detection Matrix:**

| OS | Package Manager | Lock Mechanism | Detection Method |
|----|----------------|----------------|-----------------|
| Linux (Debian/Ubuntu) | `apt` / `dpkg` | `/var/lib/dpkg/lock-frontend`, `/var/lib/apt/lists/lock` | File existence + `flock` test |
| Linux (Fedora/RHEL) | `dnf` / `rpm` | `/var/lib/rpm/.rpm.lock` | File existence |
| Linux (Arch) | `pacman` | `/var/lib/pacman/db.lck` | File existence |
| macOS | `brew` | `~/.homebrew/locks/` or `brew` process check | Process list + lockfile |
| Windows | `winget` | `msiexec.exe` process presence | WMI query / `tasklist` |
| Windows | `choco` | `C:\ProgramData\chocolatey\.chocolatey.lock` | File existence |

**Queue Persistence:** The queue is maintained as a FIFO list in `~/.gyanesh-help/queue.json`. Each entry is: `{id, command, args, status, created_at, started_at}`.

---

### 5.2 Resilient Pause/Resume System

The pause/resume system wraps child processes using a **Job Controller** module. It:

1. Spawns the underlying command as a managed child process (not a fire-and-forget shell).
2. Monitors network state on a background thread using ICMP probe packets.
3. Issues OS-appropriate suspend/resume signals on state change.
4. Writes job state to `~/.gyanesh-help/jobs.json` after every state transition.

**Network Probe Strategy:**
- Primary: ICMP ping to `1.1.1.1` (Cloudflare DNS) — fast and reliable.
- Fallback: TCP connect attempt to `8.8.8.8:53` — works on networks that block ICMP.
- Failure threshold: 3 consecutive failures = network lost.
- Recovery threshold: 1 success = network restored.

---

### 5.3 Semantic Macro Library (v1.0 Built-ins)

| Macro | Underlying Command(s) | Category |
|-------|-----------------------|----------|
| `git-undo` | `git reset --soft HEAD~1` | Git |
| `git-undo --hard` | `git reset --hard HEAD~1` | Git |
| `git-clean` | `git clean -fd` | Git |
| `git-save` | `git stash push -m "<msg>"` | Git |
| `git-unsave` | `git stash pop` | Git |
| `git-whoops` | `git commit --amend --no-edit` | Git |
| `git-oops-push` | `git push --force-with-lease` | Git |
| `git-log-pretty` | `git log --oneline --graph --decorate` | Git |
| `git-branch-clean` | Delete all local merged branches | Git |
| `git-diff-staged` | `git diff --cached` | Git |
| `find-big-files` | `find . -size +100M` | Filesystem |
| `find-old-logs` | `find . -name "*.log" -mtime +7` | Filesystem |
| `nuke-node` | Remove `node_modules` + `package-lock.json` | Node.js |
| `port-kill <port>` | Find and kill process on given port | Network |
| `port-check <port>` | Check if a port is in use | Network |
| `env-check` | Print all PATH entries, one per line | System |
| `disk-check` | Human-readable disk usage summary | System |
| `process-find <name>` | Find all running processes matching name | System |
| `docker-nuke` | Stop all containers, prune images/volumes | Docker |
| `docker-shell <container>` | `docker exec -it <container> bash` | Docker |
| `pip-freeze-clean` | `pip freeze > requirements.txt` | Python |
| `venv-create` | `python -m venv .venv && source .venv/bin/activate` | Python |
| `npm-audit-fix` | `npm audit fix --force` | Node.js |
| `ssh-keygen-github` | Generate and display SSH key for GitHub | SSH |
| `hosts-edit` | Open `/etc/hosts` in default editor with elevation | System |
| `path-add <dir>` | Persist a new directory to PATH in shell profile | System |
| `kill-port` | Interactive port-to-process killer | Network |
| `ip-info` | Local + public IP addresses | Network |
| `cert-check <domain>` | SSL cert expiry check | Security |
| `backup-now` | Alias for `workspace backup` | Workspace |

Each macro hardcodes both the **command executed** and a **plain-English explanation block**, displayed after execution.

---

### 5.4 Environment Stores — Detailed Toolchain Specifications

Each store targets the following component categories: **Runtime**, **Package Manager**, **Framework/Libraries**, **Dev Tools**, **Optional Extras**.

---

#### 5.4.1 Data Science Store (`data-science`)

| Component | Tool | Notes |
|-----------|------|-------|
| Runtime | Python 3.11+ | Via `pyenv` on Linux/macOS; `winget` on Windows |
| Runtime | R 4.x | From CRAN mirror |
| Package Manager | `pip` (bundled with Python) | — |
| Environment Manager | Miniconda (Anaconda-lite) | Full Anaconda optional via `--full` flag |
| Notebook | JupyterLab | `pip install jupyterlab` |
| Core Libraries | `numpy`, `pandas`, `matplotlib`, `seaborn`, `scipy`, `scikit-learn` | Via conda or pip |
| Data Access | `sqlalchemy`, `psycopg2`, `pymysql` | DB connectors |
| Stats (R) | `tidyverse`, `ggplot2`, `caret` | Via `Rscript -e "install.packages(...)"` |
| IDE Integration | VS Code + Python extension | Optional; prompted |
| Validation | Python, R, Jupyter all reachable on PATH | — |

**Install Command:** `gyanesh-help store install data-science`  
**Estimated Setup Time:** 12–25 min (network dependent)  
**Estimated Download Size:** ~2.5 GB (Miniconda + libraries)

---

#### 5.4.2 Front End Store (`frontend`)

| Component | Tool | Notes |
|-----------|------|-------|
| Runtime | Node.js LTS | Via `nvm` on Linux/macOS; `nvm-windows` on Windows |
| Package Managers | npm (bundled) + yarn (Corepack) + pnpm | — |
| Bundler | Vite | `npm install -g vite` |
| Linters | ESLint + Prettier | `npm install -g eslint prettier` |
| Type Checking | TypeScript | `npm install -g typescript` |
| Browser Testing | Playwright (optional) | Prompted; ~100 MB |
| Frameworks | React, Vue, Svelte starters | Via `create-` scaffolding tools; optional |
| Dev Proxy | `http-server`, `serve` | `npm install -g serve` |
| Validation | `node -v`, `npm -v`, `yarn -v`, `tsc -v` pass | — |

**Install Command:** `gyanesh-help store install frontend`  
**Estimated Setup Time:** 5–10 min  
**Estimated Download Size:** ~400 MB

---

#### 5.4.3 Full Stack — LAMP Store (`lamp`)

| Component | Tool | Notes |
|-----------|------|-------|
| OS (Linux) | Ubuntu/Debian assumed; Fedora/Arch branches | Conditionally routed |
| Web Server | Apache 2.4 | `apt install apache2` / `brew install httpd` (macOS) |
| Database | MySQL 8.x | `apt install mysql-server` / `winget install Oracle.MySQL` |
| Language | PHP 8.2+ | With `php-fpm`, `php-mysql`, `php-curl`, `php-mbstring` |
| Composer | PHP Composer 2.x | Auto-installed from `getcomposer.org` |
| Admin UI | phpMyAdmin (optional) | Prompted |
| Virtual Hosts | Auto-configures a `localhost` virtual host | — |
| Service Start | Enables and starts `apache2`, `mysql` on boot | `systemctl enable` on Linux |
| Validation | `apache2 -v`, `mysql --version`, `php -v`, `composer -V` | — |

**Install Command:** `gyanesh-help store install lamp`  
**Estimated Setup Time:** 8–15 min  
**Estimated Download Size:** ~350 MB

---

#### 5.4.4 Full Stack — MERN Store (`mern`)

| Component | Tool | Notes |
|-----------|------|-------|
| Runtime | Node.js LTS | Via `nvm` |
| Database | MongoDB Community 7.x | Official MongoDB repos for Linux; `brew tap mongodb/brew` for macOS; MSI for Windows |
| Database GUI | MongoDB Compass (optional) | ~200 MB, prompted |
| Backend Framework | Express.js | `npm install -g express-generator` |
| Frontend Framework | React | Via `create-react-app` or `vite` template |
| State Management | Redux Toolkit (optional) | Prompted |
| API Testing | Postman CLI (`newman`) | `npm install -g newman` |
| Process Manager | PM2 | `npm install -g pm2` |
| Validation | `mongod --version`, `node -v`, `pm2 --version` | — |

**Install Command:** `gyanesh-help store install mern`  
**Estimated Setup Time:** 10–18 min  
**Estimated Download Size:** ~800 MB

---

#### 5.4.5 Backend Store (`backend`)

| Component | Tool | Notes |
|-----------|------|-------|
| Containerisation | Docker Desktop (macOS/Windows) / Docker Engine (Linux) | Official install scripts |
| Container Orchestration | Docker Compose v2 | Bundled with Docker Desktop; standalone on Linux |
| DB Clients (CLI) | `psql` (PostgreSQL), `mysql` client, `redis-cli` | Package manager installs |
| DB Admin GUI | TablePlus or DBeaver (optional) | Prompted |
| API Testing | HTTPie, `curl` (usually pre-installed) | `pip install httpie` |
| Secrets Management | `pass` (Linux), macOS Keychain CLI | — |
| Task Runner | `make` | Usually pre-installed; ensured |
| Shell Enhancements | `zsh` + `oh-my-zsh` + `zsh-autosuggestions` (optional) | Prompted |
| Validation | `docker run hello-world`, all CLI clients respond to `--version` | — |

**Install Command:** `gyanesh-help store install backend`  
**Estimated Setup Time:** 10–20 min (Docker is the heavyweight)  
**Estimated Download Size:** ~1.2 GB

---

#### 5.4.6 AI Development Store (`ai-dev`)

| Component | Tool | Notes |
|-----------|------|-------|
| Runtime | Python 3.11+ | Via `pyenv` |
| GPU Driver | NVIDIA Driver (if GPU detected) | Prompted; driver version matched to CUDA |
| CUDA Toolkit | CUDA 12.x | Routed via `cuda-keyring` (Linux) / NVIDIA installer (Windows) |
| cuDNN | cuDNN 8.x/9.x | Matched to CUDA version |
| Deep Learning | PyTorch (CUDA-enabled) | Official install URL per platform |
| Deep Learning | TensorFlow 2.x | `pip install tensorflow` |
| Hugging Face | `transformers`, `datasets`, `accelerate`, `huggingface_hub` | `pip install` |
| HuggingFace CLI | `huggingface-cli` | For model download and repo management |
| Local Inference | `ollama` | From `ollama.ai` install script |
| Vector Store | `faiss-cpu` / `faiss-gpu` | Conditional on GPU detection |
| Experiment Tracking | MLflow | `pip install mlflow` |
| Environment | Jupyter + `ipywidgets` | For notebook-based experimentation |
| CPU Fallback | All installs gracefully fall back to CPU-only if no GPU is detected | Detected via `nvidia-smi` |
| Validation | `python -c "import torch; print(torch.cuda.is_available())"` etc. | — |

**Install Command:** `gyanesh-help store install ai-dev`  
**Estimated Setup Time:** 20–45 min  
**Estimated Download Size:** ~5–12 GB (CUDA-enabled builds)

---

#### 5.4.7 Claude Setup Store (`claude`)

| Component | Tool | Notes |
|-----------|------|-------|
| CLI Foundation | `anthropic` Python SDK | `pip install anthropic` |
| API Key Setup | Guided `ANTHROPIC_API_KEY` configuration | Written to shell profile + `.env` template |
| MCP SDK | `@anthropic-ai/mcp` (Node.js) | `npm install -g @anthropic-ai/mcp` |
| MCP Server Scaffolding | Official `create-mcp-server` scaffold | `npx @anthropic-ai/create-mcp-server` |
| Prompt Testing | `promptfoo` | `npm install -g promptfoo` |
| CLI Integration | Claude CLI (if officially released) | Conditioned on availability; fallback to Python SDK shell wrapper |
| Config Management | `.claude-config.json` template in home dir | — |
| Useful Extras | `jq` (for parsing API JSON responses), `httpie` | — |
| Model Listing | `gyanesh-help claude list-models` (calls Anthropic API) | Single network call; fails gracefully offline |
| Validation | `anthropic --version`, `promptfoo --version`, `mcp --version` | — |

**Install Command:** `gyanesh-help store install claude`  
**Estimated Setup Time:** 5–8 min  
**Estimated Download Size:** ~150 MB

---

### 5.5 Smart History Maintenance

History is stored in a local **SQLite database** at `~/.gyanesh-help/history.db`. Every command issued through `gyanesh-help` is recorded. The schema and query interface are detailed in the Technical Specification (Part B).

**Key Query Commands:**

| Command | Behaviour |
|---------|-----------|
| `gyanesh-help history` | Last 20 entries, tabular view |
| `gyanesh-help history --all` | Full history, paginated |
| `gyanesh-help history --tag <name>` | Filtered by project tag |
| `gyanesh-help history --search <term>` | Full-text search on command string |
| `gyanesh-help history --since <date>` | ISO date filter |
| `gyanesh-help history --failed` | Only commands with non-zero exit codes |
| `gyanesh-help history export --format csv` | Export to CSV |

---

### 5.6 Smart History — Session Tagging

Session tagging allows the user to annotate a block of related commands with a project name.

- Active tag is stored in `~/.gyanesh-help/session.json` and persists across terminal windows (machine-scoped, not session-scoped).
- The tag is written to every history record produced while it is active.

---

### 5.7 Workspace Backup & Sync

The `workspace backup` command captures the following artefacts:

| Category | What Is Captured | Notes |
|----------|-----------------|-------|
| Installed Tools | Output of `gyanesh-help store verify --all` | Versions of all store-managed tools |
| Shell Config | `~/.bashrc`, `~/.zshrc`, `~/.profile`, `~/.config/fish/config.fish` | All that exist |
| Git Config | `~/.gitconfig` | Sanitised (excludes tokens) |
| VS Code Settings | `settings.json`, `keybindings.json`, `extensions list` | Optional; prompted |
| SSH Keys | **Public keys only** (`~/.ssh/*.pub`) | Never backs up private keys |
| Custom Macros | `~/.gyanesh-help/macros.toml` | — |
| History DB | `~/.gyanesh-help/history.db` | Optional; prompted (may be large) |
| Store Manifest | JSON file listing all installed stores and versions | Used by `restore` command |

A `.gitignore` in the backup repo always excludes: `*.pem`, `*.key`, `id_rsa`, `id_ed25519`, `.env`, `.env.*`, `*.secret`.

**Restore Flow:**

1. `gyanesh-help workspace restore --repo <url>` — clones the backup repo.
2. Reads the **Store Manifest** to determine which environment stores are needed.
3. Runs `gyanesh-help store install <stack>` for each entry in the manifest.
4. Copies shell config files to their correct home locations (prompts before overwriting).
5. Copies custom macros to `~/.gyanesh-help/macros.toml`.
6. Prints a final validation summary.

---

## 6. Success Metrics

### 6.1 Adoption Metrics (3 months post-launch)

| Metric | Target |
|--------|--------|
| GitHub Stars | ≥ 500 |
| Weekly Active Users (CLI invocations tracked via opt-in anonymous counter) | ≥ 200 |
| Environment Stores installed | ≥ 1,000 total across all stacks |
| Community-submitted custom macros | ≥ 50 |

### 6.2 Performance Metrics

| Metric | Target |
|--------|--------|
| Binary / install size | < 15 MB |
| CLI startup time (cold) | < 150 ms |
| Lock detection latency (Linux/macOS/Windows) | < 100 ms |
| Network drop detection latency | < 35 seconds (3 probes × 10s + threshold logic) |
| History query response (10k records) | < 200 ms |

### 6.3 Reliability Metrics

| Metric | Target |
|--------|--------|
| Lock detection false-positive rate | < 2% |
| Pause/resume data loss incidents | 0 (measured via beta user reports) |
| Workspace restore success rate | ≥ 95% on clean machines |

---

## 7. Out of Scope — v1.0

The following items are explicitly deferred to v1.1 or later:

- **GUI / TUI dashboard** (ncurses or Electron-based).
- **Plugin marketplace** for community-distributed macro packs.
- **Remote sync** beyond GitHub (GitLab, Bitbucket, S3 backup destinations).
- **Team / org-level** shared macro libraries.
- **Scheduled task integration** (cron / Windows Task Scheduler wrapping).
- **LLM-powered command suggestion** (deliberately excluded; may be considered for a `--ai` opt-in layer in v2.0).
- **Container-based environment isolation** (dev containers, Nix flakes, Devbox).
- **Windows Subsystem for Linux (WSL)** specific optimisations (WSL is supported as a Linux target, not as a Windows-specific variant).
- **Mobile / tablet terminal support** (iOS, Android via Termux — deferred).

---

## 8. Appendix

### 8.1 Competitive Landscape

| Tool | What It Does | Gap `gyanesh-help` Fills |
|------|-------------|--------------------------|
| `homebrew` | macOS package manager | Not cross-platform; no queuing; no macros |
| `mise` / `asdf` | Version manager | No lock management; no history; no stores |
| `chezmoi` | Dotfile manager | No environment stores; no macros |
| `fig` / `warp` | AI-powered terminal | Heavy; cloud-dependent; no package manager |
| Custom shell scripts | Project-specific | Not distributable; not cross-platform |

### 8.2 Glossary

| Term | Definition |
|------|-----------|
| **Lock** | A file or process that signals a package manager is actively running, preventing concurrent writes. |
| **Store** | A curated, pre-tested set of tools that constitute a complete developer environment stack. |
| **Macro** | A memorable shorthand command that maps to one or more complex terminal operations, with a hardcoded explanation. |
| **Job** | A `gyanesh-help`-managed process that can be queued, paused, and resumed. |
| **Stack Context** | The active environment store tag associated with the current session. |
| **PAT** | GitHub Personal Access Token, used for authenticated backup/restore operations. |

### 8.3 Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1 | 2026-03-28 | PM | Initial draft |
| 1.0 | 2026-03-28 | PM | Complete PRD with all 7 stores, full user stories |