// Package sops provides SOPS argument building and execution.
package sops

import (
	"fmt"

	"github.com/enbiyagoral/sopsctl/internal/config"
)

// ArgsBuilder converts a profile configuration to SOPS CLI arguments.
type ArgsBuilder struct{}

// NewArgsBuilder creates a new ArgsBuilder.
func NewArgsBuilder() *ArgsBuilder {
	return &ArgsBuilder{}
}

// Build generates SOPS CLI arguments from a profile.
func (b *ArgsBuilder) Build(profile *config.Profile, command string, file string) ([]string, error) {
	args := make([]string, 0, 16)

	// Age backend
	if profile.Age != nil {
		keys, err := profile.Age.GetAllPublicKeys()
		if err != nil {
			return nil, fmt.Errorf("failed to get age keys: %w", err)
		}
		for _, key := range keys {
			args = append(args, "--age", key)
		}
	}

	// SOPS options
	if profile.SOPS.EncryptedRegex != "" {
		args = append(args, "--encrypted-regex", profile.SOPS.EncryptedRegex)
	}
	if profile.SOPS.EncryptedSuffix != "" {
		args = append(args, "--encrypted-suffix", profile.SOPS.EncryptedSuffix)
	}
	if profile.SOPS.UnencryptedRegex != "" {
		args = append(args, "--unencrypted-regex", profile.SOPS.UnencryptedRegex)
	}
	if profile.SOPS.UnencryptedSuffix != "" {
		args = append(args, "--unencrypted-suffix", profile.SOPS.UnencryptedSuffix)
	}

	// Command and file
	args = append(args, command, file)

	return args, nil
}

// BuildDecrypt generates arguments for decrypt.
func (b *ArgsBuilder) BuildDecrypt(file string) []string {
	return []string{"decrypt", file}
}

// BuildEdit generates arguments for edit.
func (b *ArgsBuilder) BuildEdit(profile *config.Profile, file string) ([]string, error) {
	if profile == nil {
		return []string{"edit", file}, nil
	}
	return b.Build(profile, "edit", file)
}

// GetKeyFilePath returns the key file path from profile for setting SOPS_AGE_KEY_FILE.
func (b *ArgsBuilder) GetKeyFilePath(profile *config.Profile) string {
	if profile != nil && profile.Age != nil {
		return profile.Age.GetKeyFilePath()
	}
	return ""
}
