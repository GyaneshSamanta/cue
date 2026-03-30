package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerPythonMacros() {
	reg(&macro.Macro{
		Name: "pip-freeze-clean", Category: "python",
		Description: "Export clean requirements.txt from pip",
		Commands:    []macro.Step{{OS: "all", Command: "pip freeze > requirements.txt"}},
		Explanation: `
✔ Done. requirements.txt created.
─────────────────────────────────────────────────────
All currently installed packages with pinned versions
are saved. Share this file for reproducible installs:
  pip install -r requirements.txt
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "venv-create", Category: "python",
		Description: "Create and activate a Python virtual environment",
		Commands: []macro.Step{
			{OS: "linux", Command: "python3 -m venv .venv && source .venv/bin/activate"},
			{OS: "darwin", Command: "python3 -m venv .venv && source .venv/bin/activate"},
			{OS: "windows", Command: `python -m venv .venv && .venv\Scripts\activate`},
		},
		Explanation: `
✔ Done. Virtual environment created and activated.
─────────────────────────────────────────────────────
A .venv directory was created. All pip installs now
go into this isolated environment. Deactivate with:
  deactivate
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	// Add an alias macro for 'python-venv-here' to match spec
	reg(&macro.Macro{
		Name: "python-venv-here", Category: "python",
		Description: "Create a Python virtual environment in current directory",
		Commands: []macro.Step{
			{OS: "linux", Command: "python3 -m venv .venv"},
			{OS: "darwin", Command: "python3 -m venv .venv"},
			{OS: "windows", Command: `python -m venv .venv`},
		},
		Explanation: `
✔ Virtual environment created at ./.venv!
─────────────────────────────────────────────────────
To activate it, run the following:
  macOS/Linux: source .venv/bin/activate
  Windows:     .venv\Scripts\activate
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
