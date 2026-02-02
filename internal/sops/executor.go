package sops

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executor runs SOPS commands.
type Executor struct {
	sopsPath string
}

// NewExecutor creates a new SOPS executor.
func NewExecutor(sopsPath string) *Executor {
	if sopsPath == "" {
		sopsPath = "sops"
	}
	return &Executor{sopsPath: sopsPath}
}

// Execute runs SOPS with the given arguments.
func (e *Executor) Execute(args []string) error {
	cmd := exec.Command(e.sopsPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ExecuteWithKeyFile runs SOPS with SOPS_AGE_KEY_FILE set.
func (e *Executor) ExecuteWithKeyFile(args []string, keyFile string) error {
	cmd := exec.Command(e.sopsPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if keyFile != "" {
		cmd.Env = append(os.Environ(), "SOPS_AGE_KEY_FILE="+keyFile)
	}

	return cmd.Run()
}

// DryRun prints the command that would be executed without running it.
func (e *Executor) DryRun(args []string) {
	fmt.Println("Would execute:")
	fmt.Printf("  %s %s\n", e.sopsPath, formatArgs(args))
}

// DryRunWithKeyFile prints the command with env variable that would be executed.
func (e *Executor) DryRunWithKeyFile(args []string, keyFile string) {
	fmt.Println("Would execute:")
	if keyFile != "" {
		fmt.Printf("  SOPS_AGE_KEY_FILE=%s \\\n", keyFile)
	}
	fmt.Printf("  %s %s\n", e.sopsPath, formatArgs(args))
}

// formatArgs formats arguments for display, quoting those with spaces.
func formatArgs(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		if strings.Contains(arg, " ") {
			quoted[i] = fmt.Sprintf("%q", arg)
		} else {
			quoted[i] = arg
		}
	}
	return strings.Join(quoted, " ")
}
