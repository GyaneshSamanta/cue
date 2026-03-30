package builtins

import "github.com/GyaneshSamanta/cue/internal/macro"

func init() {
	registerSecurityMacros()
}

func registerSecurityMacros() {
	macro.Register(&macro.Macro{
		Name:        "certs-local",
		Command:     "openssl req -x509 -newkey rsa:4096 -keyout localhost-key.pem -out localhost-cert.pem -days 365 -nodes -subj '/CN=localhost'",
		Description: "Generate self-signed SSL cert for localhost",
		Explanation: `Creates a self-signed SSL certificate for local development.
Generates two files:
  • localhost-key.pem — private key
  • localhost-cert.pem — certificate (valid 365 days)
For production, use Let's Encrypt or a proper CA.
If mkcert is installed, prefer: mkcert localhost`,
		Dangerous: false,
	})

	macro.Register(&macro.Macro{
		Name:        "ssh-copy-id-github",
		Command:     "cat ~/.ssh/id_ed25519.pub | gh ssh-key add -t 'cue-added'",
		Description: "Copy SSH public key to GitHub (via gh CLI)",
		Explanation: `Reads your ed25519 public key and adds it to your GitHub account.
Requires: gh CLI authenticated (run 'gh auth login' first).
This replaces the manual process of copying your key and pasting it in GitHub Settings.`,
		Dangerous: false,
	})

	macro.Register(&macro.Macro{
		Name:        "env-diff",
		Command:     "comm -23 <(grep -v '^#' .env.example | grep '=' | cut -d= -f1 | sort) <(grep -v '^#' .env | grep '=' | cut -d= -f1 | sort)",
		Description: "Compare .env vs .env.example, show missing keys",
		Explanation: `Shows which environment variables are defined in .env.example but missing from .env.
This prevents the common bug of deploying without required configuration.
Note: Only compares key names, not values.`,
		Dangerous: false,
	})
}
