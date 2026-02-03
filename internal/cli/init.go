package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const banner = `
  ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
  ┃                                                                ┃
  ┃   ███████  ██████  ██████  ███████  ██████ ████████ ██         ┃
  ┃   ██      ██    ██ ██   ██ ██      ██         ██    ██         ┃
  ┃   ███████ ██    ██ ██████  ███████ ██         ██    ██         ┃
  ┃        ██ ██    ██ ██           ██ ██         ██    ██         ┃
  ┃   ███████  ██████  ██      ███████  ██████    ██    ███████    ┃
  ┃                                                                ┃
  ┃   🔐 SOPS Profile Manager                                      ┃
  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`

var initCmd = &cobra.Command{
	Use:   "init <shell>",
	Short: "Install shell integration",
	Long: `Install shell integration to your shell config file.

This adds a shell function that makes 'sopsctl profile use <name>' 
automatically set SOPS_AGE_KEY_FILE in your current shell.

Examples:
  sopsctl init zsh    # Install to ~/.zshrc
  sopsctl init bash   # Install to ~/.bashrc`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"zsh", "bash"},
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]

		var configFile string
		switch shell {
		case "zsh":
			configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
		case "bash":
			configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
		default:
			return fmt.Errorf("unsupported shell: %s (supported: zsh, bash)", shell)
		}

		// Check if already installed
		content, err := os.ReadFile(configFile)
		if err == nil && strings.Contains(string(content), "# sopsctl shell integration") {
			fmt.Printf("✓ Already installed in %s\n", configFile)
			return nil
		}

		// Show banner for first install
		fmt.Print(banner)
		fmt.Println()

		// Append shell function to config file
		f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", configFile, err)
		}
		defer func() { _ = f.Close() }()

		_, err = f.WriteString(shellFunction)
		if err != nil {
			return fmt.Errorf("failed to write to %s: %w", configFile, err)
		}

		fmt.Printf("✓ Shell integration installed to %s\n", configFile)
		fmt.Println()
		fmt.Println("Restart your terminal or run:")
		fmt.Printf("  source %s\n", configFile)
		fmt.Println()
		fmt.Println("Then switch profiles with:")
		fmt.Println("  sopsctl profile use <name>")
		return nil
	},
}

const shellFunction = `
# sopsctl shell integration
sopsctl() {
  if [[ "${1:-}" == "profile" && "${2:-}" == "use" ]]; then
    local output; output=$(command sopsctl "$@" 2>&1); local rc=$?
    if [[ $rc -eq 0 ]]; then
      while IFS= read -r line; do
        [[ "$line" == export\ * ]] && eval "$line" && echo "✓ Set ${line#export }"
      done <<< "$output"
    else echo "$output" >&2; return $rc; fi
  else command sopsctl "$@"; fi
}
`

func init() {
	rootCmd.AddCommand(initCmd)
}
