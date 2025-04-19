package jot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration settings.
// These settings can be customized in the ~/.jot/config.yaml file.
type Config struct {
	// Editor specifies the command to use for editing notes and templates.
	Editor string `yaml:"editor"`

	// DefaultContext specifies the default context to use when creating notes.
	DefaultContext string `yaml:"default_context"`

	// StoragePath specifies the base directory for storing notes and templates.
	StoragePath string `yaml:"storage_path"`
}

// LoadConfig loads the configuration from the config file.
// If the config file doesn't exist, default values are returned.
// If the config file exists but is invalid, an error is returned.
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory for config initialization: %w", err)
	}
	configDir := filepath.Join(home, ".jot")
	configPath := filepath.Join(configDir, "config.yaml")

	cfg := &Config{
		Editor:         "vi",
		DefaultContext: "",
		StoragePath:    configDir,
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file at path '%s': %w", configPath, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file at path '%s': %w", configPath, err)
	}

	expandedPath, err := expandHome(cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand storage path '%s': %w", cfg.StoragePath, err)
	}
	cfg.StoragePath = expandedPath

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in file '%s': %w", configPath, err)
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the config file.
func (c *Config) SaveConfig() error {
	configDir := filepath.Dir(filepath.Join(c.StoragePath, "config.yaml"))
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory at path '%s': %w", configDir, err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config data to YAML: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file to path '%s': %w", configPath, err)
	}

	return nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Editor == "" {
		return fmt.Errorf("editor cannot be empty")
	}
	if c.StoragePath == "" {
		return fmt.Errorf("storage path cannot be empty")
	}
	return nil
}

// NotesDir returns the path to the notes directory.
func (c *Config) NotesDir() string {
	return filepath.Join(c.StoragePath, "notes")
}

// TemplatesDir returns the path to the templates directory.
func (c *Config) TemplatesDir() string {
	return filepath.Join(c.StoragePath, "templates")
}

// EnsureDirectories creates the necessary directories for the application.
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		c.NotesDir(),
		c.TemplatesDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create required directory at path '%s': %w", dir, err)
		}
	}

	return nil
}

// expandHome expands the tilde in a path to the user's home directory.
func expandHome(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
