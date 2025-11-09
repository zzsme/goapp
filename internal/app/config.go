package app

import (
	"encoding/json"
	"fmt"
	"os"
)

// ServerConfig contains server configuration
type ServerConfig struct {
	Port int    `json:"port"`
	Mode string `json:"mode"`
}

// LogConfig contains logging configuration
type LogConfig struct {
	Filename   string `json:"filename"`
	Level      string `json:"level"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Config is the main configuration struct
type Config struct {
	Server   ServerConfig   `json:"server"`
	Log      LogConfig      `json:"log"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
}

// ConfigData holds the application configuration
var ConfigData Config

func init() {
	// Set default configuration
	ConfigData = Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "release",
		},
		Log: LogConfig{
			Filename:   "logs/app.log",
			Level:      "info",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "",
			DBName:   "goapp",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
	}

	// Try to load configuration from file
	if _, err := os.Stat("config.json"); err == nil {
		file, err := os.Open("config.json")
		if err != nil {
			fmt.Printf("Warning: Could not open config file: %v\n", err)
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&ConfigData); err != nil {
			fmt.Printf("Warning: Could not decode config file: %v\n", err)
		}
	}
}
