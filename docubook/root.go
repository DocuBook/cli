// version 0.3.1
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	AppName        = "DocuBook CLI"
	DefaultVersion = "0.3.1" // initial version
)

var (
	Version = DefaultVersion
)

var rootCmd = &cobra.Command{
	Use:           "docubook",
	Short:         "DocuBook CLI helps you create and manage beautiful documentation sites.",
	Version:       Version, // Cobra menangani flag --version secara otomatis
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\nVersion: {{.Version}}\n", AppName))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
