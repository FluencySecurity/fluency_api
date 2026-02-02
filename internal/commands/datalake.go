package commands

import (
	"github.com/spf13/cobra"
)

var (
	datalake    string
	index       string
	managed     bool
	integration string
	schema      string
)

var lakeCmd = &cobra.Command{
	Use:   "datalake",
	Short: "Manage datalake",
}

var lakeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		//	cmd.PrintErrf("Adding datalake %s\n", datalake)
		lakes, err := AppAPI.ListDatalakes()
		if err != nil {
			return err
		}

		if len(lakes) == 0 {
			cmd.PrintErrf("No datalake found.")
			return nil
		}

		for _, entry := range lakes {
			cmd.PrintErrf("Name: %s, Description: %s\n", entry.Name, entry.StorageDescription)
		}

		//cmd.PrintErrln("Datalake added successfully")
		return nil
	},
}

var lakeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrf("Adding datalake %s  %t\n", datalake, managed)
		err := AppAPI.AddDatalake(datalake, managed, integration)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Datalake added successfully")
		return nil
	},
}

var lakeAddIndexCmd = &cobra.Command{
	Use:   "add-index",
	Short: "Add an index to the datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrf("Adding index %s to datalake %s\n", index, datalake)
		err := AppAPI.AddDatalakeIndex(datalake, index, schema)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Index added successfully")
		return nil
	},
}

var lakeListIndexCmd = &cobra.Command{
	Use:   "list-index",
	Short: "list all index ine one datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding index % to datalake %s\n", index, datalake)
		entries, err := AppAPI.ListDatalakeIndex(datalake)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.Println("No datalake index found.")
			return nil
		}

		for _, entry := range entries {
			cmd.PrintErrf("Datalake: %s, Index: %s, Description: %s\n", entry.Datalake, entry.DatalakeIndex, entry.StorageDescription)
		}
		return nil

	},
}

var lakeDeleteIndexCmd = &cobra.Command{
	Use:   "del-index",
	Short: "delete index from one datalake",
	RunE: func(cmd *cobra.Command, args []string) error {
		//cmd.PrintErrf("Adding index % to datalake %s\n", index, datalake)
		err := AppAPI.DeleteDatalakeIndex(datalake, index)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Index deleted successfully")
		return nil

	},
}

func init() {
	RootCmd.AddCommand(lakeCmd)
	lakeCmd.AddCommand(lakeAddCmd, lakeListCmd, lakeAddIndexCmd, lakeListIndexCmd, lakeDeleteIndexCmd)
	//lakeAddCmd.AddCommand(lakeAddIndexCmd)

	lakeAddCmd.Flags().StringVar(&datalake, "datalake", "", "datalake name")
	lakeAddCmd.Flags().BoolVar(&managed, "managed", true, "Managed datalake (default: true)")
	lakeAddCmd.Flags().StringVar(&integration, "integration", "", "integration id")

	//_ = lakeAddCmd.MarkFlagRequired("datalake")
	// _ = lakeAddCmd.MarkFlagRequired("name")

	// Flags for 'lake add index'
	lakeAddIndexCmd.Flags().StringVar(&datalake, "datalake", "managed", "datalake name")
	lakeAddIndexCmd.Flags().StringVar(&index, "index", "", "datalake index name")
	lakeAddIndexCmd.Flags().StringVar(&schema, "schema", "fluency default", "schema name")

	//_ = lakeAddIndexCmd.MarkFlagRequired("datalake")
	_ = lakeAddIndexCmd.MarkFlagRequired("index")

	lakeListIndexCmd.Flags().StringVar(&datalake, "datalake", "", "datalake name")
	//_ = lakeListIndexCmd.MarkFlagRequired("datalake")

	lakeDeleteIndexCmd.Flags().StringVar(&datalake, "datalake", "", "datalake name")
	lakeDeleteIndexCmd.Flags().StringVar(&index, "index", "", "datalake index name")

	_ = lakeDeleteIndexCmd.MarkFlagRequired("datalake")
	_ = lakeDeleteIndexCmd.MarkFlagRequired("index")

}
