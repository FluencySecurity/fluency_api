package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/SecurityDo/fluency_api/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	confProvider string
	confContext  string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the fluency tool",
	Long:  `Manage fluency configuration stored in ~/.fluency/config.yaml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Subcommand: SET
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values for a cluster profile",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Identify which cluster profile we are editing
		// 'cluster' is the global flag defined in root.go
		targetCluster := cluster

		// If user didn't provide --cluster, try to edit the currently active one
		if targetCluster == "" {
			targetCluster = viper.GetString("current-cluster")
		}

		// If still empty, we can't proceed
		if targetCluster == "" {
			cmd.PrintErrln("Error: No cluster name specified. Use --cluster <name> to configure a profile.")
			return
		}

		// 2. Set "Current Cluster" to this one (Switch context)
		viper.Set("current-cluster", targetCluster)

		// 3. Save values using Dot Notation (clusters.<name>.<field>)
		// This creates a nested structure in the YAML file.
		prefix := fmt.Sprintf("clusters.%s.", targetCluster)

		// We always save the provider (defaults to 'eks' via flag if not typed)
		viper.Set(prefix+"provider", confProvider)

		if namespace != "" {
			viper.Set(prefix+"namespace", namespace)
		}
		if confContext != "" {
			viper.Set(prefix+"context", confContext)
		}

		// 4. Write to disk
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Configuration saved for cluster '%s'.\n", targetCluster)
	},
}

// Subcommand: LIST
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured clusters",
	Run: func(cmd *cobra.Command, args []string) {
		current := viper.GetString("current-cluster")
		// GetStringMap returns map[string]interface{}
		clusters := viper.GetStringMap("clusters")

		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CURRENT\tCLUSTER\tPROVIDER\tNAMESPACE")
		fmt.Fprintln(w, "-------\t-------\t--------\t---------")

		// Sort keys for consistent output
		var keys []string
		for k := range clusters {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, name := range keys {
			// Extract details from the nested map
			details, ok := clusters[name].(map[string]interface{})
			if !ok {
				continue
			}

			isCurrent := ""
			if name == current {
				isCurrent = "*"
			}

			// Safe getters for interface{} map
			prov := ""
			if v, ok := details["provider"]; ok {
				prov = fmt.Sprintf("%v", v)
			}
			ns := ""
			if v, ok := details["namespace"]; ok {
				ns = fmt.Sprintf("%v", v)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", isCurrent, name, prov, ns)
		}
		w.Flush()
	},
}

// Subcommand: DELETE
var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cluster configuration",
	Run: func(cmd *cobra.Command, args []string) {
		clusterToDelete := cluster

		if clusterToDelete == "" {
			cmd.PrintErrln("Error: --cluster is required to delete a profile.")
			return
		}

		// 1. Get the raw map
		allClusters := viper.GetStringMap("clusters")

		if _, exists := allClusters[clusterToDelete]; !exists {
			fmt.Printf("Cluster '%s' not found.\n", clusterToDelete)
			return
		}

		// 2. Delete the key
		delete(allClusters, clusterToDelete)

		// 3. Set the map back to Viper
		viper.Set("clusters", allClusters)

		// 4. Handle edge case: If we deleted the "current" cluster, pick another if available
		current := viper.GetString("current-cluster")
		if current == clusterToDelete {
			if len(allClusters) == 0 {
				viper.Set("current-cluster", "")
				fmt.Println("Warning: You deleted the currently active cluster context. No clusters remain.")
			} else {
				newCurrent := pickFirstCluster(allClusters)
				viper.Set("current-cluster", newCurrent)
				fmt.Printf("Switched current cluster to '%s'.\n", newCurrent)
			}
		} else if current == "" && len(allClusters) > 0 {
			newCurrent := pickFirstCluster(allClusters)
			viper.Set("current-cluster", newCurrent)
			fmt.Printf("Current cluster was unset. Switched current cluster to '%s'.\n", newCurrent)
		}

		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Cluster '%s' deleted.\n", clusterToDelete)
	},
}

// Subcommand: VIEW (Updated to show current-cluster logic and site config)
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "SETTING\tVALUE")
		fmt.Fprintln(w, "-------\t-----")

		// Kubernetes / cluster section
		current := viper.GetString("current-cluster")
		prefix := fmt.Sprintf("clusters.%s.", current)
		fmt.Fprintf(w, "Current Cluster\t%s\n", current)
		fmt.Fprintf(w, "Provider\t%s\n", viper.GetString(prefix+"provider"))
		fmt.Fprintf(w, "Namespace\t%s\n", viper.GetString(prefix+"namespace"))
		fmt.Fprintf(w, "Context\t%s\n", viper.GetString(prefix+"context"))

		// Site config section (default site, site-config path, available sites)
		defaultSite := viper.GetString("default-site")
		siteConfigPath := viper.GetString("site-config")
		if siteConfigPath == "" {
			if cwd, err := os.Getwd(); err == nil {
				siteConfigPath = filepath.Join(cwd, "site_credentials.json")
			}
		}
		fmt.Fprintln(w, "-------\t-----")
		fmt.Fprintf(w, "Default Site\t%s\n", defaultSite)
		fmt.Fprintf(w, "Site Config File\t%s\n", siteConfigPath)
		if siteConfigPath != "" {
			if creds, err := config.LoadSiteCredentials(siteConfigPath); err == nil {
				keys := make([]string, 0, len(creds.TokenMap))
				for k := range creds.TokenMap {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for i, k := range keys {
					mark := ""
					if k == defaultSite {
						mark = " (default)"
					}
					if i == 0 {
						fmt.Fprintf(w, "Available Sites\t%s%s\n", k, mark)
					} else {
						fmt.Fprintf(w, "\t%s%s\n", k, mark)
					}
				}
			}
		}

		fmt.Fprintln(w, "-------\t-----")
		fmt.Fprintf(w, "Config File\t%s\n", viper.ConfigFileUsed())

		w.Flush()
	},
}

// Subcommand: LIST-SITES (list sites from site_credentials.json)
var configListSitesCmd = &cobra.Command{
	Use:   "list-sites",
	Short: "List available sites from site_credentials.json",
	Long:  `Lists site hostnames (keys) from the site config file. Use --site-config to specify the file; otherwise ./site_credentials.json is used. Marks the default site with * if default-site is set.`,
	Run: func(cmd *cobra.Command, args []string) {
		siteConfigPath := viper.GetString("site-config")
		if siteConfigPath == "" {
			if cwd, err := os.Getwd(); err == nil {
				siteConfigPath = filepath.Join(cwd, "site_credentials.json")
			}
		}
		if siteConfigPath == "" {
			cmd.PrintErrln("Error: No site config path. Use --site-config <path> or run from a directory containing site_credentials.json.")
			return
		}
		creds, err := config.LoadSiteCredentials(siteConfigPath)
		if err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
		defaultSite := viper.GetString("default-site")
		keys := make([]string, 0, len(creds.TokenMap))
		for k := range creds.TokenMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SITE\tDEFAULT")
		fmt.Fprintln(w, "----\t-------")
		for _, k := range keys {
			mark := ""
			if k == defaultSite {
				mark = "*"
			}
			fmt.Fprintf(w, "%s\t%s\n", k, mark)
		}
		w.Flush()
	},
}

// Subcommand: SET-DEFAULT-SITE (persist default site to ~/.fluency/config.yaml)
var configSetDefaultSiteCmd = &cobra.Command{
	Use:   "set-default-site [site]",
	Short: "Set the default site for site config",
	Long:  `Saves the default site hostname to ~/.fluency/config.yaml. When --site is not set, this site is used. Example: fluency config set-default-site demo.cloud.fluencysecurity.com`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostname := args[0]
		viper.Set("default-site", hostname)
		if err := config.SaveConfig(); err != nil {
			cmd.PrintErrln("Error saving config:", err)
			return
		}
		fmt.Printf("Default site set to %q.\n", hostname)
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configViewCmd)

	// Add new subcommands
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configListSitesCmd)
	configCmd.AddCommand(configSetDefaultSiteCmd)
	configCmd.AddCommand(configDeleteCmd)

	// Configuration for 'config' command
	// Default value "eks" is set here for the FLAG
	configSetCmd.Flags().StringVar(&confProvider, "provider", "eks", "Provider (eks|aks|gke)")
	configSetCmd.Flags().StringVar(&confContext, "context", "", "Kubeconfig context name")

	_ = configSetCmd.MarkFlagRequired("context")
	_ = configSetCmd.MarkFlagRequired("namespace")
	_ = configSetCmd.MarkFlagRequired("cluster")

	// Set the global default for Viper as well (in case user views config without setting it)
	viper.SetDefault("provider", "eks")
}

// pickFirstCluster returns the first cluster name in sorted order.
func pickFirstCluster(clusters map[string]interface{}) string {
	var names []string
	for name := range clusters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names[0]
}

/*
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the fluency tool",
	Long:  `Sets configuration values in ~/.fluency/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set values in Viper
		// Note: cluster and namespace are already bound in root.go,
		// but we need to explicitly set them if the user provided flags here to save them to file.

		if cluster != "" {
			viper.Set("cluster", cluster)
		}
		if namespace != "" {
			viper.Set("namespace", namespace)
		}
		if confProvider != "" {
			viper.Set("provider", confProvider)
		}
		if confContext != "" {
			viper.Set("context", confContext)
		}

		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("Configuration saved.")
	},
}*/
/*
// 1. Define the 'view' subcommand
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration settings",
	Long:  "Displays the current configuration loaded from ~/.fluency/config.yaml and environment variables.",
	Run: func(cmd *cobra.Command, args []string) {
		// Use tabwriter to create a clean, aligned table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "SETTING\tVALUE")
		fmt.Fprintln(w, "-------\t-----")

		// Retrieve values from Viper (checks flags, env vars, and config file)
		fmt.Fprintf(w, "Cluster\t%s\n", viper.GetString("cluster"))
		fmt.Fprintf(w, "Namespace\t%s\n", viper.GetString("namespace"))
		fmt.Fprintf(w, "Provider\t%s\n", viper.GetString("provider"))
		fmt.Fprintf(w, "Context\t%s\n", viper.GetString("context"))

		// Access config file location
		fmt.Fprintln(w, "-------\t-----")
		fmt.Fprintf(w, "Config File\t%s\n", viper.ConfigFileUsed())

		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Local flags for the config command
	configCmd.Flags().StringVar(&confProvider, "provider", "eks", "Provider (eks|aks|gke)")
	configCmd.Flags().StringVar(&confContext, "context", "", "Kubeconfig context name")
	// Note: --cluster and --namespace are inherited from Root, but usually 'config' commands
	// might want to enforce them or treat them differently. For now, we rely on the Root persistent flags.

	// 2. Register 'view' as a child of 'config'
	configCmd.AddCommand(configViewCmd)
}*/
