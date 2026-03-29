package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/gyanesh-help/internal/store"
	"github.com/GyaneshSamanta/gyanesh-help/internal/ui"
)

var onboardingCmd = &cobra.Command{
	Use:   "onboarding",
	Short: "Run the interactive first-time setup wizard",
	Run: func(cmd *cobra.Command, args []string) {
		RunOnboarding()
	},
}

func init() {
	rootCmd.AddCommand(onboardingCmd)
}

// RunOnboarding executes the rich welcome experience for new users.
func RunOnboarding() {
	fmt.Println()
	ui.PrintHeader(" Welcome to Gyanesh-help v2.0 ")
	time.Sleep(500 * time.Millisecond)

	ui.PrintInfo("Your terminal is about to get significantly more capable.")
	fmt.Println("This wizard will give you a quick 3-step tour of the tools you now have installed.")
	fmt.Println()

	if !ui.Confirm("Press [Enter] or type 'y' to continue... ") {
		return
	}
	fmt.Println()

	// 1. Environment Stores
	ui.PrintHeader("Step 1/3: Environment Stores")
	fmt.Println("Are you tired of configuring paths, installing package managers, and debugging binaries?")
	fmt.Println("Gyanesh-help includes 'Stores'—verified, pre-configured software stacks (like 'devops' or 'mern').")
	fmt.Println()
	if ui.Confirm("Would you like to install a stack right now? (You can say no) [y/N] ") {
		stack := pickStack()
		if stack != "" {
			ui.PrintStep(fmt.Sprintf("Installing '%s'...", stack))
			if err := store.Install(stack, osAdapter, store.InstallOpts{Verify: true}); err != nil {
				ui.PrintError(err.Error())
			} else {
				ui.PrintSuccess("Installed successfully!")
			}
		}
	} else {
		ui.PrintDim("➔ You can explore stores later by typing: gyanesh-help store")
	}
	fmt.Println()
	time.Sleep(500 * time.Millisecond)

	// 2. Macro System
	ui.PrintHeader("Step 2/3: The Macro System")
	fmt.Println("Instead of memorizing long flags or looking up docker-compose syntax,")
	fmt.Println("gyanesh-help replaces messy workflows with highly readable 'Macros'.")
	fmt.Println()
	ui.PrintInfo("For example, run 'gyanesh-help go-mod-tidy-check' to automatically format and test Go code.")
	ui.PrintInfo("Or 'gyanesh-help nuke-docker-volume' to instantly reclaim space safely.")
	fmt.Println()
	ui.PrintDim("➔ View all 25+ built-in macros at any time by typing: gyanesh-help macro list")
	fmt.Println()
	time.Sleep(1000 * time.Millisecond)

	// 3. AI & Claude Code Integration
	ui.PrintHeader("Step 3/3: Claude Code & Local LLMs")
	fmt.Println("Gyanesh-help orchestrates direct AI functionality directly in your CLI.")
	fmt.Println("You can run Anthropic's 'Claude Code' in several modes:")
	fmt.Println()
	fmt.Println("  1. API Mode   (Directly connects to Claude over the cloud)")
	fmt.Println("  2. Local Mode (Uses Ollama + liteLLm to run entirely on your machine)")
	fmt.Println()
	ui.PrintSuccess("Highlights: The Local Mode implementation is 100% FREE and private!")
	fmt.Println()
	ui.PrintDim("➔ Install AI capabilities anytime by typing: gyanesh-help claude-code install")
	fmt.Println()
	time.Sleep(1000 * time.Millisecond)

	// Conclusion
	ui.PrintHeader(" You are all set! ")
	fmt.Println("You can replay this tutorial via: gyanesh-help onboarding")
	fmt.Println()
	fmt.Println(ui.C(ui.Bold, "Built with <3 by Gyanesh"))
	fmt.Println("Support the project here: " + ui.C(ui.Blue, "https://buymeachai.ezee.li/GyaneshOnProduct"))
	fmt.Println()
}
