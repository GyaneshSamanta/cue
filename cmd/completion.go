package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for gyanesh-help.

To load completions:

Bash:
  $ source <(gyanesh-help completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ gyanesh-help completion bash > /etc/bash_completion.d/gyanesh-help
  # macOS:
  $ gyanesh-help completion bash > $(brew --prefix)/etc/bash_completion.d/gyanesh-help

Zsh:
  $ gyanesh-help completion zsh > "${fpath[1]}/_gyanesh-help"
  # You might need to start a new shell or run: compinit

Fish:
  $ gyanesh-help completion fish | source
  $ gyanesh-help completion fish > ~/.config/fish/completions/gyanesh-help.fish

PowerShell:
  PS> gyanesh-help completion powershell | Out-String | Invoke-Expression
  PS> gyanesh-help completion powershell > gyanesh-help.ps1
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
