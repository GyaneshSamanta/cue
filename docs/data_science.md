# 📊 Data Science Stack

Data Science requires absolute stability across library versions and C++ bindings. The `cue` CLI completely eliminates the "it works on my machine" problem for Data Science teams by standardizing the toolchain.

## The Environment Store: `data-science`

Running `cue store install data-science` provisions a scientifically robust environment:
- **Python 3.10+**: Optimized for SciPy/NumPy compatability.
- **JupyterLab**: Interactive graphing and data exploration.
- **Poetry**: Deterministic dependency management.
- **Pandas Data-Stack**: Pre-compiled binaries for your architecture.

## Dedicated Data Science Macros

Data scientists shouldn't be memorizing `pip` freeze commands. Our macros handle best-practices automatically:

### 1. `pip-freeze-clean`
- **Command:** `pip freeze > requirements.txt`
- **Why it matters:** Generates a deterministic lockfile of all currently installed packages. Use this before pushing any notebook to ensure the next researcher can reproduce your environment exactly.

### 2. `python-venv-here`
- **Command:** `python3 -m venv .venv`
- **Why it matters:** Creates an isolated environment. It prevents global pip pollution, meaning your deep learning project doesn't accidentally overwrite the dependencies of your data scraping project.

### 3. `jupyter-serve` (Custom Macro Example)
You can easily add your own macros to boot up environments.
```bash
cue macro add jupyter-serve "jupyter lab --no-browser --port 8080" "Boot up local analysis suite"
```

---
*Spend less time resolving NumPy C-binding errors and more time building models. Run `cue store install data-science`.*
