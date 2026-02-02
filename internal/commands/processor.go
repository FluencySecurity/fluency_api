package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	procName    string
	procContent string
	procType    string // Default to "fpl_processor" if not specified
	procDesc    string // Optional description for the processor
)

var processorCmd = &cobra.Command{
	Use:   "processor",
	Short: "Manage processors",
}

var listProcessorCmd = &cobra.Command{
	Use:   "list",
	Short: "List all processors",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing processors...")

		entries, err := AppAPI.ListProcessor()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			cmd.PrintErrln("No processor found.")
			return nil
		}
		for _, entry := range entries {
			cmd.PrintErrf("Name: %s, Type: %s\n", entry.Name, entry.Type)
		}
		return nil
	},
}

var processorAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a processor",
	// Example usage:
	// 1. fluency processor add --name my-proc --content "@./my-script.js"
	// 2. fluency processor add --name my-proc --content "function process() { ... }"
	// 3. cat my-script.js | fluency processor add --name my-proc --content -
	RunE: func(cmd *cobra.Command, args []string) error {
		var content string
		//var err error

		// CHECK: Is the user asking to read from Stdin?
		if procContent == "-" {
			// Read from the pipe
			b, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(b)
		} else if len(procContent) > 1 && procContent[0] == '@' {
			filePath := procContent[1:]
			// Read from the file path provided
			b, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file '%s': %w", filePath, err)
			}
			content = string(b)
		} else {
			content = procContent
		}

		if len(content) == 0 {
			return fmt.Errorf("processor content is empty")
		}

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Deploying processor '%s' (%d bytes)...\n", procName, len(content))
		err := AppAPI.AddProcessor(procName, content, procType, procDesc)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Processor added successfully")
		return nil
	},
}

var processorDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a processor",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Deleting processor '%s'...\n", procName)
		err := AppAPI.DeleteProcessor(procName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Processor deleted successfully")
		return nil
	},
}

// fluency processor add --name filter --content "@./scripts/filter.js"
// echo "function process() { ... }" | fluency processor add --name filter --content -
func init() {
	RootCmd.AddCommand(processorCmd)
	processorCmd.AddCommand(processorAddCmd, listProcessorCmd, processorDelCmd) // Add del similarly

	//processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	//processorAddCmd.Flags().StringVar(&procFile, "file", "", "Processor file path")

	processorAddCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	processorAddCmd.Flags().StringVar(&procContent, "content", "", "Processor content or file path (use '-' for stdin)")
	processorAddCmd.Flags().StringVar(&procType, "type", "fpl_processor", "Processor type (fpl_processor|fpl_receiver|fpl_packer|fpl_report)")
	processorAddCmd.Flags().StringVar(&procDesc, "desc", "", "Processor description (optional)")

	_ = processorAddCmd.MarkFlagRequired("name")
	_ = processorAddCmd.MarkFlagRequired("content")

	processorDelCmd.Flags().StringVar(&procName, "name", "", "Processor name")
	_ = processorDelCmd.MarkFlagRequired("name")

}
