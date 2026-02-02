package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the global configuration
type Config struct {
	Cluster   string `mapstructure:"cluster"`
	Provider  string `mapstructure:"provider"`
	Context   string `mapstructure:"context"`
	Namespace string `mapstructure:"namespace"`
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in ~/.fluency directory with name "config.yaml"
	configPath := filepath.Join(home, ".fluency")
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("fluency")
	viper.AutomaticEnv() // read in environment variables that match (e.g. FLUENCY_CLUSTER, FLUENCY_NAMESPACE)

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	// --- ADD THIS BLOCK ---
	// If a config file is loaded, hydrate the root variables from the active cluster
	if viper.ConfigFileUsed() != "" {
		current := viper.GetString("current-cluster")
		if current != "" {
			// Read values from the active profile
			prefix := "clusters." + current + "."

			// If the specific cluster has a value, override the root "default"
			if v := viper.GetString(prefix + "namespace"); v != "" {
				viper.Set("namespace", v)
			}
			if v := viper.GetString(prefix + "provider"); v != "" {
				viper.Set("provider", v)
			}
			if v := viper.GetString(prefix + "context"); v != "" {
				viper.Set("context", v)
			}
			// Ensure the 'cluster' key matches the current selection
			viper.Set("cluster", current)
		}
	}
}

// SaveConfig writes the current viper configuration to disk
func SaveConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".fluency")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		_ = os.Mkdir(configDir, 0755)
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// WriteConfigAs ensures we save to the specific path
	return viper.WriteConfigAs(filepath.Join(configDir, "config.yaml"))
}
