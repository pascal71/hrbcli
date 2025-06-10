package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the Harbor CLI configuration
type Config struct {
	HarborURL      string `yaml:"harbor_url" json:"harbor_url"`
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password,omitempty" json:"password,omitempty"`
	APIVersion     string `yaml:"api_version" json:"api_version"`
	Insecure       bool   `yaml:"insecure" json:"insecure"`
	OutputFormat   string `yaml:"output_format" json:"output_format"`
	DefaultProject string `yaml:"default_project,omitempty" json:"default_project,omitempty"`
	NoColor        bool   `yaml:"no_color" json:"no_color"`
	Debug          bool   `yaml:"debug" json:"debug"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	if cfgFile := viper.ConfigFileUsed(); cfgFile != "" {
		return cfgFile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".hrbcli.yaml")
}

// Load loads the configuration from file
func Load() (*Config, error) {
	cfg := &Config{
		HarborURL:      viper.GetString("harbor_url"),
		Username:       viper.GetString("username"),
		Password:       viper.GetString("password"),
		APIVersion:     viper.GetString("api_version"),
		Insecure:       viper.GetBool("insecure"),
		OutputFormat:   viper.GetString("output_format"),
		DefaultProject: viper.GetString("default_project"),
		NoColor:        viper.GetBool("no_color"),
		Debug:          viper.GetBool("debug"),
	}

	// Set defaults
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v2.0"
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "table"
	}

	return cfg, nil
}

// Save saves the configuration to file
func Save(cfg *Config) error {
	configPath := GetConfigPath()

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Don't save password to file
	cfgToSave := *cfg
	cfgToSave.Password = ""

	// Marshal to YAML
	data, err := yaml.Marshal(&cfgToSave)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Set sets a configuration value
func Set(key string, value interface{}) error {
	viper.Set(key, value)

	// Load current config
	cfg, err := Load()
	if err != nil {
		return err
	}

	// Update the specific field
	switch key {
	case "harbor_url":
		cfg.HarborURL = value.(string)
	case "username":
		cfg.Username = value.(string)
	case "password":
		cfg.Password = value.(string)
	case "api_version":
		cfg.APIVersion = value.(string)
	case "insecure":
		cfg.Insecure = value.(bool)
	case "output_format":
		cfg.OutputFormat = value.(string)
	case "default_project":
		cfg.DefaultProject = value.(string)
	case "no_color":
		cfg.NoColor = value.(bool)
	case "debug":
		cfg.Debug = value.(bool)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	// Save to file
	return Save(cfg)
}

// Get gets a configuration value
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString gets a string configuration value
func GetString(key string) string {
	return viper.GetString(key)
}

// GetBool gets a boolean configuration value
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// ValidateConfig validates the configuration
func ValidateConfig() error {
	if viper.GetString("harbor_url") == "" {
		return fmt.Errorf("harbor_url is required")
	}

	// Validate output format
	format := viper.GetString("output_format")
	switch format {
	case "table", "json", "yaml", "":
		// Valid formats
	default:
		return fmt.Errorf("invalid output format: %s (valid: table, json, yaml)", format)
	}

	return nil
}
