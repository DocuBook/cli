package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docubook",
	Short: "DocuBook CLI - Create beautiful documentation",
	Long: `DocuBook CLI helps you create and manage beautiful documentation sites.
With DocuBook, you can quickly set up a documentation site with a modern look and feel.`,
}

// Execute executes the root command and returns any error encountered
func Execute() error {
	return rootCmd.Execute()
}
