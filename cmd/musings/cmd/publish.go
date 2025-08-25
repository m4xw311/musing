package cmd

import (
	"fmt"

	"github.com/m4xw311/musing/internal/blog"
	"github.com/m4xw311/musing/internal/site"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish blog posts to static website",
	Long:  `Publish blog posts written in markdown to a static website.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Publishing blog posts...")

		// Create blog instance
		b := blog.NewBlog("posts")

		// Load existing posts
		if err := b.LoadPosts(); err != nil {
			fmt.Printf("Error loading posts: %v\n", err)
			return
		}

		// Create static site generator
		s := site.NewStaticSiteGenerator("posts", "public")

		// Generate site
		if err := s.Generate(); err != nil {
			fmt.Printf("Error generating site: %v\n", err)
			return
		}

		fmt.Println("Blog posts published successfully to public/ directory!")
	},
}
