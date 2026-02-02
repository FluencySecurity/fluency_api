package commands

import (
	"github.com/spf13/cobra"
)

var (
	repoName string
	//procType    string // Default to "parser" if not specified
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import resources from github repository",
}

var importProcessorCmd = &cobra.Command{
	Use:   "processor",
	Short: "Import processors",
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Println("Listing processors...")

		err := AppAPI.ImportProcessor(procType, repoName)
		if err != nil {
			return err
		}
		return nil
	},
}

var importAppTemplateCmd = &cobra.Command{
	Use:   "application",
	Short: "Import application template",
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Println("Listing processors...")

		err := AppAPI.ImportAppTemplate(repoName)
		if err != nil {
			return err
		}
		return nil
	},
}

// fluency processor add --name filter --content "@./scripts/filter.js"
// echo "function process() { ... }" | fluency processor add --name filter --content -
func init() {
	RootCmd.AddCommand(importCmd)
	importCmd.AddCommand(importProcessorCmd, importAppTemplateCmd)

	importProcessorCmd.Flags().StringVar(&procType, "type", "fpl_processor", "Processor type (fpl_processor|fpl_receiver|fpl_packer|fpl_report)")

}
