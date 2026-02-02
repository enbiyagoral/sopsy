package cli

import (
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt <file>",
	Short: "Decrypt a file",
	Long: `Decrypt a file using SOPS.

If a profile is specified with -p, the key_file from that profile will be used.
Otherwise, SOPS will use the default key file location.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		// Build SOPS arguments
		sopsArgs := builder.BuildDecrypt(file)

		// Get key file from profile if specified
		var keyFile string
		if profileName != "" {
			profile, err := cfg.GetProfile(profileName)
			if err != nil {
				return err
			}
			keyFile = builder.GetKeyFilePath(profile)
		} else if cfg.DefaultProfile != "" {
			// Use default profile's key file
			profile, err := cfg.GetProfile(cfg.DefaultProfile)
			if err == nil {
				keyFile = builder.GetKeyFilePath(profile)
			}
		}

		// Execute or dry-run
		if dryRun {
			executor.DryRunWithKeyFile(sopsArgs, keyFile)
			return nil
		}

		return executor.ExecuteWithKeyFile(sopsArgs, keyFile)
	},
}
