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
	DefaultVersion = "1.0.0"
)

// These variables are set during build time via ldflags
var (
	// Version holds the application version
	Version = DefaultVersion
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docubook",
	Short: "DocuBook CLI - Create beautiful documentation",
	Long: `DocuBook CLI helps you create and manage beautiful documentation sites.
With DocuBook, you can quickly set up a documentation site with a modern look and feel.

Find more information at https://docubook.pro`,
	Version:       Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// Set up version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\nVersion: {{.Version}}\n", AppName))

	// Add global flags
	rootCmd.Flags().BoolP("version", "v", false, "Print the version information")
}

// PrintVersion prints the application version information
func PrintVersion() {
	fmt.Printf("%s\nVersion: %s\n", AppName, Version)
}

// Execute executes the root command and handles any errors
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
