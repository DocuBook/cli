package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	// AppName is the name of the application
	AppName = "DocuBook CLI"
	// DefaultVersion is the default version when not set via ldflags
	DefaultVersion = "0.2.2"
)

// These variables are set during build time via ldflags
var (
	// Version holds the application version
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
  create      Create a new DocuBook project

Flags:
  -h, --help      help for docubook
  --version   Print the version information`,
	// By setting the Version field, Cobra automatically adds the --version flag
	// and handles its execution. No Run function or manual flag is needed.
	Version:       Version,
	SilenceUsage:  true,
	SilenceErrors: true,
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
