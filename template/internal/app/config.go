package app

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		DB       int    `mapstructure:"db"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`

	Log struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	} `mapstructure:"log"`
}

// ConfigData holds the application configuration
var ConfigData Config

// LoadConfig loads the configuration from the specified file
func LoadConfig() error {
	// Set default configuration values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "development")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file", "logs/app.log")

	// Look for config file in the current directory
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create a default one
			fmt.Println("Config file not found, creating a default one")
			if err := createDefaultConfig(); err != nil {
				return fmt.Errorf("failed to create default config: %w", err)
			}
			// Try to read the config again
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Load environment variables
	viper.AutomaticEnv()

	// Unmarshal the configuration
	if err := viper.Unmarshal(&ConfigData); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Println("Configuration loaded successfully")
	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig() error {
	defaultConfig := `server:
  port: 8080
  mode: development

database:
  host: localhost
  port: 5432
  name: myapp
  user: postgres
  password: secret

redis:
  host: localhost
  port: 6379
  db: 0

log:
  level: info
  file: logs/app.log
`
	if err := os.WriteFile("config.yaml", []byte(defaultConfig), 0644); err != nil {
		return err
	}
	return nil
}

// init loads the configuration when the package is imported
func init() {
	if err := LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
}
