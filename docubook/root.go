// version 0.3.5
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	AppName        = "DocuBook CLI"
	DefaultVersion = "0.3.5" // initial version
)

var (
	Version = DefaultVersion
)

var rootCmd = &cobra.Command{
	Use:           "docubook",
	Short:         "DocuBook CLI written on GO! Initialize, Update, Push and Deploy your Docs direct into Terminal.",
	Version:       Version, // Cobra handles the --version flag automatically
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	// FIX: Change the template to one line with a separator " - "
	// The first newline character is removed and replaced with " - ".
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s - Version: {{.Version}}\n", AppName))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
