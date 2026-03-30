package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerNodejsMacros() {
	reg(&macro.Macro{
		Name: "npm-audit-fix", Category: "nodejs",
		Description: "Auto-fix npm audit vulnerabilities",
		Commands:    []macro.Step{{OS: "all", Command: "npm audit fix --force"}},
		Explanation: `
✔ Done. npm vulnerabilities were auto-fixed.
─────────────────────────────────────────────────────
--force may have upgraded or downgraded some packages.
Run 'npm test' to verify nothing is broken.
Review changes: git diff package.json package-lock.json
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})

	reg(&macro.Macro{
		Name: "node-version-check", Category: "nodejs",
		Description: "Verify installation by printing Node and NPM versions",
		Commands:    []macro.Step{{OS: "all", Command: "node -v && npm -v"}},
		Explanation: `
✔ Displayed current Node and NPM versions.
─────────────────────────────────────────────────────
If these versions do not match your project requirements,
consider running 'cue version' to switch runtimes.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
