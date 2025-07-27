package docubook

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// These will be set during build time via ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "docubook",
	Short: "DocuBook CLI - Create beautiful documentation",
	Long: `DocuBook CLI helps you create and manage beautiful documentation sites.
With DocuBook, you can quickly set up a documentation site with a modern look and feel.`,
	Version: Version,
}

func init() {
	rootCmd.SetVersionTemplate(`DocuBook CLI
Version: {{.Version}}
Commit: ` + Commit + `
Build Time: ` + BuildTime + `
`)
	rootCmd.Flags().BoolP("version", "v", false, "Print the version information")
	rootCmd.SetVersionTemplate(`DocuBook CLI
Version: {{.Version}}
Commit: ` + Commit + `
Build Time: ` + BuildTime + `
`)
}

// Execute executes the root command and returns any error encountered
func Execute() error {
	return rootCmd.Execute()
}
