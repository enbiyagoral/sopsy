package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile represents a SOPS encryption profile.
type Profile struct {
	Name        string `yaml:"-"` // Populated from map key
	Description string `yaml:"description,omitempty"`

	// Encryption backend
	Age *AgeConfig `yaml:"age,omitempty"`

	// SOPS-specific options
	SOPS SOPSOptions `yaml:"sops,omitempty"`
}

// AgeConfig represents age encryption configuration.
type AgeConfig struct {
	// KeyFile is the path to the age key file (contains both public and private keys)
	KeyFile string `yaml:"key_file,omitempty"`
	// Recipients are explicit public keys (alternative to KeyFile)
	Recipients []string `yaml:"recipients,omitempty"`
}

// SOPSOptions represents SOPS-specific encryption options.
type SOPSOptions struct {
	EncryptedRegex    string `yaml:"encrypted_regex,omitempty"`
	EncryptedSuffix   string `yaml:"encrypted_suffix,omitempty"`
	UnencryptedRegex  string `yaml:"unencrypted_regex,omitempty"`
	UnencryptedSuffix string `yaml:"unencrypted_suffix,omitempty"`
}

// GetBackendSummary returns a human-readable summary of configured backends.
func (p *Profile) GetBackendSummary() string {
	if p.Age != nil && (p.Age.KeyFile != "" || len(p.Age.Recipients) > 0) {
		return "age"
	}
	return "none"
}

// HasBackends returns true if the profile has at least one backend configured.
func (p *Profile) HasBackends() bool {
	return p.Age != nil && (p.Age.KeyFile != "" || len(p.Age.Recipients) > 0)
}

// GetPublicKey extracts the public key from an age key file or returns recipients.
func (a *AgeConfig) GetPublicKey() (string, error) {
	if len(a.Recipients) > 0 {
		return a.Recipients[0], nil
	}

	if a.KeyFile == "" {
		return "", fmt.Errorf("no key_file or recipients configured")
	}

	keyFile := expandPath(a.KeyFile)
	file, err := os.Open(keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to open key file: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Public key is in comment: "# public key: age1..."
		if strings.HasPrefix(line, "# public key:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("no public key found in key file: %s", keyFile)
}

// GetAllPublicKeys returns all public keys (from file and recipients).
func (a *AgeConfig) GetAllPublicKeys() ([]string, error) {
	var keys []string

	// Add recipients first
	keys = append(keys, a.Recipients...)

	// Add key from file if specified
	if a.KeyFile != "" {
		key, err := a.GetPublicKey()
		if err != nil {
			return nil, err
		}
		// Avoid duplicates
		found := false
		for _, k := range keys {
			if k == key {
				found = true
				break
			}
		}
		if !found {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// GetKeyFilePath returns the expanded key file path.
func (a *AgeConfig) GetKeyFilePath() string {
	if a.KeyFile == "" {
		return ""
	}
	return expandPath(a.KeyFile)
}

// expandPath expands ~ to home directory.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	return filepath.Clean(path)
}
