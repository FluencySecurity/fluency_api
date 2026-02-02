package commands

import (
	//"fmt"

	"github.com/spf13/cobra"
	//"github.com/SecurityDo/fluency_api/model"
)

var ()

// Parent command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "System status commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding user: %s (Role: %s)\n", authName, authRole)
		err := AppAPI.CheckStatus()
		if err != nil {
			cmd.PrintErrf("Error adding user: %s %v\n", authName, err)
			return err
		}
		return nil
	},
}

// ... Repeat similar logic for del/update user/token ...

func init() {
	RootCmd.AddCommand(statusCmd)

}
