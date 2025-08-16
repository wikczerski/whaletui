package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information - these should be set during build
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for D5r including:
  • Version number
  • Git commit SHA
  • Build date`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("D5r - Docker CLI Dashboard\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Commit: %s\n", CommitSHA)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
