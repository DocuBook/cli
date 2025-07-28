// version 0.3.2
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	AppName        = "DocuBook CLI"
	DefaultVersion = "0.3.2" // initial version
)

var (
	Version = DefaultVersion
)

var rootCmd = &cobra.Command{
	Use:           "docubook",
	Short:         "DocuBook CLI written on GO! Initialize, Update, Push and Deploy your Docs direct into Terminal.",
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
