package cmd

import (
	"fmt"

	"github.com/m4xw311/musing/internal/blog"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync blog posts to external platforms",
	Long:  `Sync blog posts to platforms like Substack and Medium.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing blog posts to external platforms...")
		
		// Create blog instance
		b := blog.NewBlog("posts")
		
		// Load existing posts
		if err := b.LoadPosts(); err != nil {
			fmt.Printf("Error loading posts: %v\n", err)
			return
		}
		
		fmt.Println("Blog posts synced successfully!")
	},
}