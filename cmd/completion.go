package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for cue.

To load completions:

Bash:
  $ source <(cue completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ cue completion bash > /etc/bash_completion.d/cue
  # macOS:
  $ cue completion bash > $(brew --prefix)/etc/bash_completion.d/cue

Zsh:
  $ cue completion zsh > "${fpath[1]}/_cue"
  # You might need to start a new shell or run: compinit

Fish:
  $ cue completion fish | source
  $ cue completion fish > ~/.config/fish/completions/cue.fish

PowerShell:
  PS> cue completion powershell | Out-String | Invoke-Expression
  PS> cue completion powershell > cue.ps1
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
