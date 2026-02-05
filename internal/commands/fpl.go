package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	fplReportFile string
	fplID         uint
)

var fplCmd = &cobra.Command{
	Use:   "fpl",
	Short: "FPL report and task operations",
}

var fplRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an FPL v2 report (run_report)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if fplReportFile == "" {
			return fmt.Errorf("--file is required")
		}
		data, err := os.ReadFile(fplReportFile)
		if err != nil {
			return err
		}
		var req model.RunFPLV2Report
		if err := json.Unmarshal(data, &req); err != nil {
			return err
		}
		id, err := AppAPI.RunReport(&req)
		if err != nil {
			return err
		}
		cmd.Printf("%d\n", id)
		return nil
	},
}

var fplGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get FPL task by ID (get_task)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.GetTaskByID(fplID)
		if err != nil {
			return err
		}
		if resp.Entry == nil {
			return fmt.Errorf("task %d: no entry in response", fplID)
		}
		cmd.Printf("name: %s state: %s\n", resp.Entry.Name, resp.Entry.State)
		return nil
	},
}

var fplResultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Get FPL results by task ID (get_results)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.GetResultsByID(fplID)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(resp.Result)
	},
}

func init() {
	RootCmd.AddCommand(fplCmd)
	fplCmd.AddCommand(fplRunCmd, fplGetCmd, fplResultsCmd)

	fplRunCmd.Flags().StringVarP(&fplReportFile, "file", "f", "", "Path to JSON file with RunFPLV2Report (reportName, arguments, etc.)")
	_ = fplRunCmd.MarkFlagRequired("file")

	fplGetCmd.Flags().UintVar(&fplID, "id", 0, "Task ID")
	_ = fplGetCmd.MarkFlagRequired("id")

	fplResultsCmd.Flags().UintVar(&fplID, "id", 0, "Task ID")
	_ = fplResultsCmd.MarkFlagRequired("id")
}
