package commands

import (
	"github.com/spf13/cobra"
)

// Parent command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "System status commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := AppAPI.CheckStatus()
		if err != nil {
			cmd.PrintErrf("Error checking status: %v\n", err)
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
