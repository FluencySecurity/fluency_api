package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var collectorCmd = &cobra.Command{
	Use:   "collector",
	Short: "Collector service",
}

var collectorStatusName string

var collectorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collectors (CollectorForWeb)",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := AppAPI.CollectorList()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.PrintErrln("No collectors found.")
			return nil
		}
		for _, e := range entries {
			cmd.Printf("Name: %s, LastPoll: %d\n", e.Name, e.LastPoll)
		}
		return nil
	},
}

// collectorStatusCargs is fixed for the status command (always el9).
var collectorStatusCargs = map[string]interface{}{"status": "el9"}

var collectorStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get collector status (kargs: collector + cargs status=el9)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if collectorStatusName == "" {
			return fmt.Errorf("collector name is required (--collector)")
		}
		out, err := AppAPI.CollectorStatus(collectorStatusName, collectorStatusCargs)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		cmd.Printf("%s\n", b)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(collectorCmd)
	collectorCmd.AddCommand(collectorListCmd, collectorStatusCmd)

	collectorStatusCmd.Flags().StringVar(&collectorStatusName, "collector", "", "Collector name")
}
