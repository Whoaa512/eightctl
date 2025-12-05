package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage device temperature schedules (cloud)",
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schedules",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		scheds, err := cl.ListSchedules(context.Background())
		if err != nil {
			return err
		}
		rows := make([]map[string]any, 0, len(scheds))
		for _, s := range scheds {
			rows = append(rows, map[string]any{
				"id":      s.ID,
				"start":   s.StartTime,
				"level":   s.Level,
				"days":    s.DaysOfWeek,
				"enabled": s.Enabled,
			})
		}
		rows = output.FilterFields(rows, viper.GetStringSlice("fields"))
		return output.Print(output.Format(viper.GetString("output")), []string{"id", "start", "level", "days", "enabled"}, rows)
	},
}

var scheduleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create schedule",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		start := viper.GetString("start")
		if start == "" {
			return fmt.Errorf("--start HH:MM required")
		}
		level := viper.GetInt("level")
		days := viper.GetIntSlice("days")
		if len(days) == 0 {
			return fmt.Errorf("--days required")
		}
		enabled := !viper.GetBool("disabled")
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		s := client.TemperatureSchedule{StartTime: start, Level: level, DaysOfWeek: days, Enabled: enabled}
		res, err := cl.CreateSchedule(context.Background(), s)
		if err != nil {
			return err
		}
		fmt.Printf("created schedule %s\n", res.ID)
		return nil
	},
}

var scheduleUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update schedule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		patch := map[string]any{}
		if cmd.Flags().Changed("start") {
			patch["startTime"] = viper.GetString("start")
		}
		if cmd.Flags().Changed("level") {
			patch["level"] = viper.GetInt("level")
		}
		if cmd.Flags().Changed("days") {
			patch["daysOfWeek"] = viper.GetIntSlice("days")
		}
		if cmd.Flags().Changed("enabled") {
			patch["enabled"] = viper.GetBool("enabled")
		}
		if len(patch) == 0 {
			return fmt.Errorf("no fields to update")
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if _, err := cl.UpdateSchedule(context.Background(), args[0], patch); err != nil {
			return err
		}
		fmt.Println("updated")
		return nil
	},
}

var scheduleDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete schedule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.DeleteSchedule(context.Background(), args[0]); err != nil {
			return err
		}
		fmt.Println("deleted")
		return nil
	},
}

func init() {
	scheduleCreateCmd.Flags().String("start", "", "HH:MM start time")
	scheduleCreateCmd.Flags().Int("level", 0, "Temperature level -100..100")
	scheduleCreateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	scheduleCreateCmd.Flags().Bool("disabled", false, "Create disabled")
	viper.BindPFlag("start", scheduleCreateCmd.Flags().Lookup("start"))
	viper.BindPFlag("level", scheduleCreateCmd.Flags().Lookup("level"))
	viper.BindPFlag("days", scheduleCreateCmd.Flags().Lookup("days"))
	viper.BindPFlag("disabled", scheduleCreateCmd.Flags().Lookup("disabled"))

	scheduleUpdateCmd.Flags().String("start", "", "HH:MM start time")
	scheduleUpdateCmd.Flags().Int("level", 0, "Temperature level -100..100")
	scheduleUpdateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	scheduleUpdateCmd.Flags().Bool("enabled", true, "Enable/disable schedule")
	viper.BindPFlag("start", scheduleUpdateCmd.Flags().Lookup("start"))
	viper.BindPFlag("level", scheduleUpdateCmd.Flags().Lookup("level"))
	viper.BindPFlag("days", scheduleUpdateCmd.Flags().Lookup("days"))
	viper.BindPFlag("enabled", scheduleUpdateCmd.Flags().Lookup("enabled"))

	scheduleCmd.AddCommand(scheduleListCmd, scheduleCreateCmd, scheduleUpdateCmd, scheduleDeleteCmd)
}
