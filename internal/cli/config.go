package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/enbiyagoral/sopsctl/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage sopsctl configuration",
	Long:  `Manage sopsctl configuration (init, show, edit).`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long: `Create a new sopsctl configuration file.

The config file will be created at ~/.config/sopsctl/config.yaml`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.DefaultConfigPath()
		if err != nil {
			return err
		}
		if cfgFile != "" {
			path = cfgFile
		}

		// Check if already exists
		if _, err := os.Stat(path); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				return fmt.Errorf("config already exists at %s (use --force to overwrite)", path)
			}
		}

		// Create empty config
		newCfg := config.NewConfig()

		if err := config.Save(newCfg, path); err != nil {
			return err
		}

		fmt.Printf("Configuration created at %s\n", path)
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Add a profile:   sopsctl profile add <name> --age-key-file <path>")
		fmt.Println("  2. Set default:     sopsctl profile use <name>")
		fmt.Println("  3. Encrypt a file:  sopsctl encrypt <file>")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration in your editor",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.DefaultConfigPath()
		if err != nil {
			return err
		}
		if cfgFile != "" {
			path = cfgFile
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		c := exec.Command(editor, path)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	configInitCmd.Flags().Bool("force", false, "overwrite existing config")

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configEditCmd)
}
