package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/steipete/eightctl/internal/client"
	"github.com/steipete/eightctl/internal/output"
)

var travelCmd = &cobra.Command{Use: "travel", Short: "Travel / jetlag endpoints"}

var travelTripsCmd = &cobra.Command{Use: "trips", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Travel().Trips(context.Background())
	if err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"trips"}, []map[string]any{{"trips": res}})
}}

var travelPlansCmd = &cobra.Command{Use: "plans", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	trip := viper.GetString("trip")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Travel().Plans(context.Background(), trip)
	if err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"plans"}, []map[string]any{{"plans": res}})
}}

var travelTasksCmd = &cobra.Command{Use: "tasks", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	plan := viper.GetString("plan")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Travel().PlanTasks(context.Background(), plan)
	if err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"tasks"}, []map[string]any{{"tasks": res}})
}}

var travelAirportCmd = &cobra.Command{Use: "airport-search", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	query := viper.GetString("query")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Travel().AirportSearch(context.Background(), query)
	if err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"airports"}, []map[string]any{{"airports": res}})
}}

var travelFlightCmd = &cobra.Command{Use: "flight-status", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	flight := viper.GetString("flight")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	res, err := cl.Travel().FlightStatus(context.Background(), flight)
	if err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"flight"}, []map[string]any{{"flight": res}})
}}

func init() {
	travelPlansCmd.Flags().String("trip", "", "trip id")
	viper.BindPFlag("trip", travelPlansCmd.Flags().Lookup("trip"))
	travelTasksCmd.Flags().String("plan", "", "plan id")
	viper.BindPFlag("plan", travelTasksCmd.Flags().Lookup("plan"))
	travelAirportCmd.Flags().String("query", "", "airport query")
	viper.BindPFlag("query", travelAirportCmd.Flags().Lookup("query"))
	travelFlightCmd.Flags().String("flight", "", "flight number")
	viper.BindPFlag("flight", travelFlightCmd.Flags().Lookup("flight"))

	travelCmd.AddCommand(travelTripsCmd, travelPlansCmd, travelTasksCmd, travelAirportCmd, travelFlightCmd)
}
