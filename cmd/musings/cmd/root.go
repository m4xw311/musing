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

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(syncCmd)
}