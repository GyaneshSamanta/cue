package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func registerSSHMacros() {
	reg(&macro.Macro{
		Name: "ssh-keygen-github", Category: "ssh",
		Description: "Generate an ED25519 SSH key for GitHub",
		Commands: []macro.Step{
			{OS: "linux", Command: `ssh-keygen -t ed25519 -C "github" -f ~/.ssh/id_ed25519 -N "" && cat ~/.ssh/id_ed25519.pub`},
			{OS: "darwin", Command: `ssh-keygen -t ed25519 -C "github" -f ~/.ssh/id_ed25519 -N "" && cat ~/.ssh/id_ed25519.pub`},
			{OS: "windows", Command: `ssh-keygen -t ed25519 -C "github" -f %USERPROFILE%\.ssh\id_ed25519 -N "" & type %USERPROFILE%\.ssh\id_ed25519.pub`},
		},
		Explanation: `
✔ Done. SSH key generated and public key displayed.
─────────────────────────────────────────────────────
Copy the public key above and add it to GitHub:
  Settings → SSH Keys → New SSH Key

The private key is at ~/.ssh/id_ed25519 — NEVER share it.
The public key (.pub) is safe to share anywhere.
─────────────────────────────────────────────────────`,
		BuiltIn: true,
	})
}
