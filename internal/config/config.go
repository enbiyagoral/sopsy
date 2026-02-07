// Package config handles sopsy configuration loading, saving, and management.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the main sopsy configuration.
type Config struct {
	Version        string              `yaml:"version"`
	DefaultProfile string              `yaml:"default_profile,omitempty"`
	Profiles       map[string]*Profile `yaml:"profiles,omitempty"`
	Settings       Settings            `yaml:"settings,omitempty"`
}

// Settings contains global sopsy settings.
type Settings struct {
	FZFOptions string `yaml:"fzf_options,omitempty"`
	SOPSPath   string `yaml:"sops_path,omitempty"`
}

// DefaultConfigPath returns the default configuration file path.
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "sopsy", "config.yaml"), nil
}

// NewConfig creates a new empty configuration with defaults.
func NewConfig() *Config {
	return &Config{
		Version:  "1",
		Profiles: make(map[string]*Profile),
		Settings: Settings{
			SOPSPath: "sops",
		},
	}
}

// Load loads configuration from the specified path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := NewConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Populate profile names from map keys
	for name, profile := range cfg.Profiles {
		profile.Name = name
	}

	return cfg, nil
}

// Save saves the configuration to the specified path.
func Save(cfg *Config, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetProfile returns a profile by name.
func (c *Config) GetProfile(name string) (*Profile, error) {
	profile, ok := c.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", name)
	}
	return profile, nil
}

// AddProfile adds a new profile.
func (c *Config) AddProfile(profile *Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name is required")
	}
	if _, exists := c.Profiles[profile.Name]; exists {
		return fmt.Errorf("profile already exists: %s", profile.Name)
	}
	c.Profiles[profile.Name] = profile
	return nil
}

// RemoveProfile removes a profile by name.
func (c *Config) RemoveProfile(name string) error {
	if _, exists := c.Profiles[name]; !exists {
		return fmt.Errorf("profile not found: %s", name)
	}
	delete(c.Profiles, name)
	return nil
}

// ListProfiles returns all profiles as a slice.
func (c *Config) ListProfiles() []*Profile {
	profiles := make([]*Profile, 0, len(c.Profiles))
	for _, profile := range c.Profiles {
		profiles = append(profiles, profile)
	}
	return profiles
}
