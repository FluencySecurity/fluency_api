package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	integType        string
	integName        string
	integDesc        string
	integID          string
	configParams     map[string]string
	configBoolParams map[string]string
	configIntParams  map[string]int64
	configJsonFlags  []string
	secretParams     map[string]string
	addSource        bool
)

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Manage integrations",
}

func getSecretRaw() (raw json.RawMessage, err error) {
	// This function is a placeholder for the actual implementation
	// that retrieves the raw JSON configuration.
	// For now, it returns an empty JSON object.
	secret := make(map[string]interface{})
	for key, value := range secretParams {
		if len(value) > 1 && value[0] == '@' {
			filePath := value[1:] // Remove '@'

			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file for key '%s': %w", key, err)
			}

			// Replace the file path with the actual file content
			secret[key] = string(content)
		} else {
			secret[key] = value
		}
	}
	b, _ := json.Marshal(secret)
	return b, nil
}
func getConfigRaw() (raw json.RawMessage, err error) {
	// This function is a placeholder for the actual implementation
	// that retrieves the raw JSON configuration.
	// For now, it returns an empty JSON object.
	config := make(map[string]interface{})
	for key, value := range configParams {
		if len(value) > 1 && value[0] == '@' {
			filePath := value[1:] // Remove '@'

			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file for key '%s': %w", key, err)
			}

			// Replace the file path with the actual file content
			config[key] = string(content)
		} else {
			config[key] = value
		}
	}
	for k, v := range configIntParams {
		config[k] = v
	}
	for k, v := range configBoolParams {

		if strings.EqualFold("true", v) {
			config[k] = true
		} else if strings.EqualFold("false", v) {
			config[k] = false
		} else {
			return nil, fmt.Errorf("invalid boolean value '%s' for key '%s', expected 'true' or 'false'", v, k)
		}

		//config[k] = v
	}
	// Iterate over raw flags: ["config={\"a\":1}", "retries=3", "debug=true"]
	for _, flag := range configJsonFlags {
		// 1. Split key and value manually on the first "="
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format '%s', expected key=json_value", flag)
		}
		key := parts[0]
		rawValue := parts[1]

		// 2. Unmarshal the value part as JSON
		var typedValue interface{}
		err := json.Unmarshal([]byte(rawValue), &typedValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON for key '%s': %w", key, err)
		}

		config[key] = typedValue
	}

	b, _ := json.Marshal(config)

	return b, nil
}

//return json.RawMessage("{}")

var integrationAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrf("Adding integration %s of type %s\n", integName, integType)

		entry := &model.Integration{
			Integration: integType,
			Name:        integName,
			Description: integDesc,
			//ConfigParameters:   configParams,
			//SecretParameters:   secretParams,
		}
		configRaw, err := getConfigRaw()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}
		secretRaw, err := getSecretRaw()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}
		entry.Config = configRaw
		entry.Secret = secretRaw
		id, err := AppAPI.AddIntegration(entry)
		if err != nil {
			return err
		}

		cmd.PrintErrln("Integration added successfully: ", id)
		if addSource {
			resp, err := AppAPI.AddDataSource(&model.DataSourceConfig{
				Type: "plugin",
				Name: integName,
				Plugin: &model.PluginSourceConfig{
					ID: id,
				},
			})
			if err != nil {
				cmd.PrintErrf("Error adding plugin data source: with integration %s %v\n", id, err)
				return err
			}
			cmd.PrintErrln("Plugin data source added successfully: ", resp.ID)
			fmt.Println(resp.ID)
			return nil
		}
		fmt.Println(id)
		return nil
	},
}

var integrationDelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete an integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Deleting integration %s\n", integID)
		err := AppAPI.DeleteIntegration(integID)
		if err != nil {
			cmd.PrintErrf("Error deleting integration: %s %v\n", integID, err)
			return err
		}
		cmd.PrintErrln("Integration deleted successfully: ", integID)
		return nil
	},
}

var integrationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all integrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Printf("Deleting integration %s\n", integName)
		entries, err := AppAPI.ListIntegration()
		if err != nil {
			return err
		}
		cmd.PrintErrln("Listing Integration...")
		if len(entries) == 0 {
			cmd.PrintErrln("No integration found.")
			return nil
		}

		for _, entry := range entries {
			cmd.PrintErrf("ID: %s, Name: %s, Integration: %s, Description: %s\n", entry.ID, entry.Name, entry.Integration, entry.Description)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(integrationCmd)
	integrationCmd.AddCommand(integrationAddCmd, integrationDelCmd, integrationListCmd)

	// Flags
	integrationAddCmd.Flags().StringVar(&integType, "integration", "", "Integration type")
	integrationAddCmd.Flags().StringVar(&integName, "name", "", "Name")
	integrationAddCmd.Flags().StringVar(&integDesc, "description", "", "Description")
	integrationAddCmd.Flags().BoolVar(&addSource, "add-source", false, "add plugin-type data source")

	integrationAddCmd.Flags().StringToStringVarP(&configParams, "config", "", nil, "Configuration string type parameters")
	integrationAddCmd.Flags().StringToStringVarP(&configBoolParams, "config-bool", "", nil, "Configuration boolean type parameters")
	integrationAddCmd.Flags().StringToInt64VarP(&configIntParams, "config-int", "", nil, "Configuration int type parameters")
	integrationAddCmd.Flags().StringArrayVar(&configJsonFlags, "config-json", []string{}, "Set JSON values (e.g. key=[1,2])")
	integrationAddCmd.Flags().StringToStringVarP(&secretParams, "secret", "", nil, "Secret parameters")

	_ = integrationAddCmd.MarkFlagRequired("integration")
	_ = integrationAddCmd.MarkFlagRequired("name")

	integrationDelCmd.Flags().StringVar(&integID, "id", "", "Integration id")
	_ = integrationDelCmd.MarkFlagRequired("id")

}
