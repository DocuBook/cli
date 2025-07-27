// root.go (Improved Version)

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	AppName        = "DocuBook CLI"
	DefaultVersion = "0.2.0"
)

var (
	Version = DefaultVersion
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docubook",
	Short: "DocuBook CLI helps you create and manage beautiful documentation sites.",
	Long: `DocuBook CLI helps you create and manage beautiful documentation sites.
With DocuBook, you can quickly set up a documentation site with a modern look and feel.

Find more information at https://docubook.pro

Usage:
  docubook [command]

Available Commands:
  cli         to initialize the CLI equipment
  create      Create a new DocuBook project

Flags:
  -h, --help      help for docubook
  --version   Print the version information`,
	Version:       Version, // This is all you need for version handling
	SilenceUsage:  true,
	SilenceErrors: true,
	// No 'Run' function is needed. If no subcommand is given,
	// Cobra will automatically show the help text.
}

func init() {
	// Customize the output of the automatic --version flag
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\nVersion: {{.Version}}\n", AppName))
}

// Execute executes the root command and handles any errors
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
