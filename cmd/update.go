package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
	"github.com/GyaneshSamanta/gyanesh-help/internal/updater"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install gyanesh-help updates",
	Long: `Check for the latest version of gyanesh-help on GitHub and optionally
download and install it. The current binary is backed up before replacement.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkOnly, _ := cmd.Flags().GetBool("check")
		rollback, _ := cmd.Flags().GetBool("rollback")

		if rollback {
			if err := updater.Rollback(appVersion); err != nil {
				ui.PrintError(err.Error())
			}
			return
		}

		ui.PrintStep("Checking for updates...")
		release, hasUpdate, err := updater.CheckUpdate(appVersion)
		if err != nil {
			ui.PrintError(err.Error())
			return
		}

		if !hasUpdate {
			ui.PrintSuccess(fmt.Sprintf("You're on the latest version (v%s)", appVersion))
			return
		}

		fmt.Printf("\n  Current version : v%s\n", appVersion)
		fmt.Printf("  Latest version  : %s\n\n", release.TagName)

		if release.Body != "" {
			body := release.Body
			if len(body) > 300 {
				body = body[:300] + "..."
			}
			fmt.Printf("  Release notes:\n  %s\n\n", body)
		}

		if checkOnly {
			ui.PrintInfo(fmt.Sprintf("Run 'gyanesh-help update' to install %s", release.TagName))
			return
		}

		if !ui.Confirm(fmt.Sprintf("Download and install %s? [Y/n] ", release.TagName)) {
			return
		}

		if err := updater.Update(release, appVersion); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

func init() {
	updateCmd.Flags().Bool("check", false, "Only check version, don't install")
	updateCmd.Flags().Bool("rollback", false, "Rollback to previous version")
}
