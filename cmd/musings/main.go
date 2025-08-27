// Package main provides the entry point for the musings command-line application.
//
// The musings application is a tool for publishing markdown-based static blogs
// and syncing them to external platforms. It provides two main commands:
// - publish: Generates a static website from markdown blog posts
// - sync: Syncs blog posts to external platforms
package main

import (
	"log"
	"os"

	"github.com/m4xw311/musing/cmd/musings/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}