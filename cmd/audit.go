package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/audit"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Security audit of your development environment",
	Long: `Read-only scan of your development environment covering:
  • Installed tool CVE checks (via osv.dev data)
  • SSH key strength and permissions
  • Git credential security
  • Secret scanning in dotfiles`,
	Run: func(cmd *cobra.Command, args []string) {
		toolsOnly, _ := cmd.Flags().GetBool("tools")
		sshOnly, _ := cmd.Flags().GetBool("ssh")
		gitOnly, _ := cmd.Flags().GetBool("git")
		secretsOnly, _ := cmd.Flags().GetBool("secrets")

		if toolsOnly {
			audit.AuditTools()
		} else if sshOnly {
			audit.AuditSSH()
		} else if gitOnly {
			audit.AuditGit()
		} else if secretsOnly {
			audit.AuditSecrets()
		} else {
			audit.RunFull()
		}
	},
}

func init() {
	auditCmd.Flags().Bool("tools", false, "CVE check on installed tools only")
	auditCmd.Flags().Bool("ssh", false, "SSH key audit only")
	auditCmd.Flags().Bool("git", false, "Git credential audit only")
	auditCmd.Flags().Bool("secrets", false, "Secret scan only")
}
