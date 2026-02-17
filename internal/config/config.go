package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// ConfigFileName is the name of the config file
	ConfigFileName = ".ahrefsrc"
)

// Config represents the CLI configuration
type Config struct {
	APIKey string `json:"api_key"`
}

// Load loads the configuration from file
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ConfigFileName), nil
}

// GetAPIKey gets the API key from config, env var, or returns empty string
func GetAPIKey() string {
	// First check env var
	if key := os.Getenv("AHREFS_API_KEY"); key != "" {
		return key
	}

	// Then check config file
	cfg, err := Load()
	if err != nil {
		return ""
	}

	return cfg.APIKey
}
