package cli

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <file>",
	Short: "Edit an encrypted file",
	Long: `Edit an encrypted file using SOPS.

Opens the file in your $EDITOR after decrypting, then re-encrypts on save.
If a profile is specified, uses the key_file from that profile.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		var sopsArgs []string
		var keyFile string
		var err error

		if profileName != "" {
			profile, err := cfg.GetProfile(profileName)
			if err != nil {
				return err
			}
			sopsArgs, err = builder.BuildEdit(profile, file)
			if err != nil {
				return err
			}
			keyFile = builder.GetKeyFilePath(profile)
		} else if cfg.DefaultProfile != "" {
			profile, err := cfg.GetProfile(cfg.DefaultProfile)
			if err == nil {
				sopsArgs, err = builder.BuildEdit(profile, file)
				if err != nil {
					return err
				}
				keyFile = builder.GetKeyFilePath(profile)
			} else {
				sopsArgs = []string{"edit", file}
			}
		} else {
			sopsArgs = []string{"edit", file}
		}

		// Avoid unused variable error
		_ = err

		// Execute or dry-run
		if dryRun {
			executor.DryRunWithKeyFile(sopsArgs, keyFile)
			return nil
		}

		return executor.ExecuteWithKeyFile(sopsArgs, keyFile)
	},
}
