package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerWorkspaceMacros() {
	reg(&macro.Macro{
		Name: "backup-now", Category: "workspace",
		Description: "Quick alias for workspace backup",
		Commands:    []macro.Step{{OS: "all", Command: "cue workspace backup"}},
		Explanation: `
✔ Workspace backup triggered.
─────────────────────────────────────────────────────
This is a shortcut for 'cue workspace backup'.
Your shell configs, macros, and store manifests will
be pushed to your private GitHub backup repo.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
