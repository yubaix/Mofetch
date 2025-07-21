// Pretty much this entire file is a copy of Ashish Kumar's Mufetch config.go file
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds OMDB API credentials
type Config struct {
	OmdbApiKey string `mapstructure:"omdb_api_key"`
}

// InitConfig sets up configuration directory and default values
func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Create config directory in user's home/.config/mofetch
	configDir := filepath.Join(home, ".config", "mofetch")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set default empty values for credentials
	viper.SetDefault("omdb_api_key", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.SafeWriteConfig() // Create config file if not found
		}
		return err
	}

	return nil
}

// GetConfig unmarshals configuration into Config struct
func GetConfig() (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SetCredentials saves OMDB API credentials to config file
func SetCredentials(OmdbApiKey string) error {
	viper.Set("omdb_api_key", OmdbApiKey)
	return viper.WriteConfig()
}

// HasCredentials checks if valid OMDB credentials are configured
func HasCredentials() bool {
	config, err := GetConfig()
	if err != nil {
		return false
	}
	return config.OmdbApiKey != ""
}

// init initializes the viper configuration
func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MOFETCH")
}
