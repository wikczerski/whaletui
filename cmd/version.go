package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information - these should be set during build
	Version = "dev"
	// CommitSHA is the Git commit SHA hash
	CommitSHA = "unknown"
	// BuildDate is the build timestamp
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for whaletui including:
	- Version number
	- Build date
	- Git commit hash
	- Go version`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("whaletui - Docker CLI Dashboard\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Commit: %s\n", CommitSHA)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
