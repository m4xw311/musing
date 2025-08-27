// Package cmd implements the command-line interface for the musings application.
//
// It provides commands for publishing blog posts to static websites and syncing
// posts to external platforms.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "musings",
	Short: "A tool to publish markdown-based static blogs",
	Long: `Musings is a simple tool to publish markdown based static blog 
and publish the same to platforms like substack and medium.`,
}

// Execute runs the root command and starts the CLI application.
// It returns an error if the command execution fails.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(syncCmd)
}
