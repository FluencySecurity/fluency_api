package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	appName             string
	instanceName        string
	instanceDisplayName string
	inputParams         map[string]string
	templateContent     string
)

var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Manage Fluency Application",
}

var appTemplateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed application templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface
		templates, err := AppAPI.ListAppTemplates()
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			cmd.PrintErrf("No application template found.")
			return nil
		}

		for _, entry := range templates {
			cmd.PrintErrf("Name: %s, Description: %s\n", entry.Name, entry.Description)
		}
		return nil
	},
}

func collectInputParameters(cmd *cobra.Command, configParams map[string]string, sensitive bool) (paras []*model.InputParameter, err error) {
	for key, value := range configParams {
		if len(value) > 1 && value[0] == '@' {
			filePath := value[1:] // Remove '@'

			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file for key '%s': %w", key, err)
			}
			paras = append(paras, &model.InputParameter{
				Name:      key,
				Sensitive: sensitive,
				Value:     string(content),
			})

			// Replace the file path with the actual file content
			//config[key] = string(content)
		} else {
			//config[key] = value
			paras = append(paras, &model.InputParameter{
				Name:      key,
				Value:     string(value),
				Sensitive: sensitive,
			})
		}
	}
	return paras, nil
}

var appTemplateInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface

		//config := make(map[string]string)
		var paras []*model.InputParameter

		configs, err := collectInputParameters(cmd, configParams, false)
		if err != nil {
			return fmt.Errorf("failed to collect configuration parameters: %w", err)
		}
		secrets, err := collectInputParameters(cmd, secretParams, true)
		if err != nil {
			return fmt.Errorf("failed to collect secret parameters: %w", err)
		}
		paras = append(paras, configs...)
		paras = append(paras, secrets...)

		err = AppAPI.InstallAppInstance(appName, instanceName, instanceDisplayName, paras)
		if err != nil {
			return err
		}

		cmd.PrintErrf("Application %s:%s added successfully: ", appName, instanceName)
		//cmd.Println(id)
		return nil
	},
}

var appInstanceGetCmd = &cobra.Command{
	Use:   "get-instance",
	Short: "get application instance output",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface

		resp, err := AppAPI.GetAppInstance(appName, instanceName)
		if err != nil {
			return err
		}
		for _, output := range resp.Outputs {
			if output.Sensitive {
				cmd.PrintErrf("Output %s is sensitive and cannot be displayed\n", output.Name)
				continue
			}
			cmd.Printf("Output %s: %s\n", output.Name, output.Value)
		}

		//cmd.PrintErrf("Application %s:%s added successfully: ", appName, instanceName)
		//cmd.Println(id)
		return nil
	},
}

var appTemplateUnInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall application instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Just call the global interface
		err := AppAPI.UnInstallAppInstance(appName, instanceName)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Application instance removed successfully")
		return nil
	},
}

var appTemplateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a template",
	// Example usage:
	// 1. fluency processor add --name my-proc --content "@./my-script.js"
	// 2. fluency processor add --name my-proc --content "function process() { ... }"
	// 3. cat my-script.js | fluency processor add --name my-proc --content -
	RunE: func(cmd *cobra.Command, args []string) error {
		var content string
		//var err error

		// CHECK: Is the user asking to read from Stdin?
		if templateContent == "-" {
			// Read from the pipe
			b, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(b)
		} else if len(templateContent) > 1 && templateContent[0] == '@' {
			filePath := templateContent[1:]
			// Read from the file path provided
			b, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file '%s': %w", filePath, err)
			}
			content = string(b)
		} else {
			content = templateContent
		}

		if len(content) == 0 {
			return fmt.Errorf("template content is empty")
		}

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		//cmd.PrintErrf("Deploying template '%s' (%d bytes)...\n", appName, len(content))
		id, err := AppAPI.AddTemplate(content)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Template added successfully: ", id)
		return nil
	},
}

var appTemplateDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a template",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Deleting application template '%s'...\n", appName)
		err := AppAPI.DeleteTemplate(appName)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Application template deleted successfully")
		return nil
	},
}

var appTemplateUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a template",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Now you have the content in 'content' variable
		//  cmd.PrintErrln()
		//  cmd.Printf( )  for output data/result
		cmd.PrintErrf("Updating application template '%s'...\n", appName)
		err := AppAPI.UpdateTemplate(appName, templateContent)
		if err != nil {
			return err
		}
		cmd.PrintErrln("Application template updated successfully")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(applicationCmd)
	applicationCmd.AddCommand(appTemplateListCmd, appTemplateInstallCmd, appTemplateUnInstallCmd, appTemplateAddCmd, appTemplateDelCmd, appTemplateUpdateCmd, appInstanceGetCmd) // Add del/update similarly

	appTemplateInstallCmd.Flags().StringVar(&instanceName, "instance", "", "Instance name")
	appTemplateInstallCmd.Flags().StringVar(&appName, "app", "", "Application name")
	appTemplateInstallCmd.Flags().StringVar(&instanceDisplayName, "displayName", "", "Display name")

	appTemplateInstallCmd.Flags().StringToStringVarP(&configParams, "config", "", nil, "Configuration string type parameters")
	appTemplateInstallCmd.Flags().StringToStringVarP(&secretParams, "secret", "", nil, "Secret parameters")

	// Mark required
	_ = appTemplateInstallCmd.MarkFlagRequired("instance")
	_ = appTemplateInstallCmd.MarkFlagRequired("app")

	appTemplateUnInstallCmd.Flags().StringVar(&instanceName, "instance", "", "Instance name")
	appTemplateUnInstallCmd.Flags().StringVar(&appName, "app", "", "Application name")

	// Mark required
	_ = appTemplateUnInstallCmd.MarkFlagRequired("instance")
	_ = appTemplateUnInstallCmd.MarkFlagRequired("app")

	//streamAddCmd.AddCommand(streamAddSourceCmd)
	// Add other leaf commands: sink, router, connection

	appTemplateAddCmd.Flags().StringVar(&templateContent, "content", "", "Application template content or file path (use '-' for stdin)")
	_ = appTemplateAddCmd.MarkFlagRequired("content")

	appTemplateDelCmd.Flags().StringVar(&appName, "app", "", "Application name")
	_ = appTemplateDelCmd.MarkFlagRequired("app")

	appTemplateUpdateCmd.Flags().StringVar(&appName, "app", "", "Application name")
	appTemplateUpdateCmd.Flags().StringVar(&templateContent, "content", "", "Application template content or file path (use '-' for stdin)")

	_ = appTemplateUpdateCmd.MarkFlagRequired("app")

	_ = appTemplateUpdateCmd.MarkFlagRequired("content")

	appInstanceGetCmd.Flags().StringVar(&instanceName, "instance", "", "Instance name")
	appInstanceGetCmd.Flags().StringVar(&appName, "app", "", "Application name")

	_ = appInstanceGetCmd.MarkFlagRequired("instance")
	_ = appInstanceGetCmd.MarkFlagRequired("app")

}
