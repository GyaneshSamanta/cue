package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/GyaneshSamanta/cue/internal/schedule"
	"github.com/GyaneshSamanta/cue/internal/ui"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage scheduled maintenance tasks",
	Long: `Schedule recurring maintenance tasks using OS-native scheduling:
  • Linux: systemd user timers
  • macOS: launchd plists
  • Windows: Task Scheduler`,
	Run: func(cmd *cobra.Command, args []string) {
		schedule.List()
	},
}

var scheduleBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Schedule recurring workspace backups",
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetString("interval")
		if interval == "" {
			interval = "daily"
		}
		execPath, _ := os.Executable()
		task := schedule.Task{
			Name:     "backup",
			Command:  execPath + " workspace backup",
			Interval: interval,
		}
		if err := schedule.Schedule(task); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var scheduleDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Schedule recurring health checks",
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetString("interval")
		if interval == "" {
			interval = "weekly"
		}
		execPath, _ := os.Executable()
		task := schedule.Task{
			Name:     "doctor",
			Command:  execPath + " doctor",
			Interval: interval,
		}
		if err := schedule.Schedule(task); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var scheduleUpdateCheckCmd = &cobra.Command{
	Use:   "update-check",
	Short: "Schedule recurring update checks",
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetString("interval")
		if interval == "" {
			interval = "daily"
		}
		execPath, _ := os.Executable()
		task := schedule.Task{
			Name:     "update-check",
			Command:  execPath + " update --check",
			Interval: interval,
		}
		if err := schedule.Schedule(task); err != nil {
			ui.PrintError(err.Error())
		}
	},
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	Run: func(cmd *cobra.Command, args []string) {
		schedule.List()
	},
}

var scheduleRemoveCmd = &cobra.Command{
	Use:   "remove [task-name]",
	Short: "Remove a scheduled task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		schedule.Remove(args[0])
	},
}

func init() {
	scheduleBackupCmd.Flags().String("interval", "daily", "Schedule interval: daily, weekly, hourly")
	scheduleDoctorCmd.Flags().String("interval", "weekly", "Schedule interval: daily, weekly, hourly")
	scheduleUpdateCheckCmd.Flags().String("interval", "daily", "Schedule interval: daily, weekly, hourly")

	scheduleCmd.AddCommand(scheduleBackupCmd, scheduleDoctorCmd, scheduleUpdateCheckCmd, scheduleListCmd, scheduleRemoveCmd)
}
