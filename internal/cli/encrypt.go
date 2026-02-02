package cli

import (
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt <file>",
	Short: "Encrypt a file using the selected profile",
	Long: `Encrypt a file using SOPS with the selected profile's configuration.

If no profile is specified, you'll be prompted to select one using fzf.
The profile determines which encryption backends (age, KMS, etc.) are used.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		// Select profile
		profile, err := selectProfile()
		if err != nil {
			return err
		}

		// Build SOPS arguments
		sopsArgs, err := builder.Build(profile, "encrypt", file)
		if err != nil {
			return err
		}

		// Get key file for decrypt operations within encrypt (re-encryption scenarios)
		keyFile := builder.GetKeyFilePath(profile)

		// Execute or dry-run
		if dryRun {
			executor.DryRunWithKeyFile(sopsArgs, keyFile)
			return nil
		}

		return executor.ExecuteWithKeyFile(sopsArgs, keyFile)
	},
}
