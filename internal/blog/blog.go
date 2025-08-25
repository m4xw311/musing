package blog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Post represents a blog post
type Post struct {
	Title       string    // TODO: Extract from the first # Heading
	Content     string    // TODO: Extract from the markdown file
	CreatedDate time.Time // TODO: Use should be free to omit it but then it should be updated in the md file for suture
	UpdatedDate time.Time // TODO: Use should be free to omit it but then it should be updated in the md file for suture
	Slug        string    // TODO: Derive from title
	Tags        []string
	Published   bool
}

// Blog represents a collection of blog posts
type Blog struct {
	Posts []Post
	Path  string
}

// NewBlog creates a new blog instance
func NewBlog(path string) *Blog {
	return &Blog{
		Posts: make([]Post, 0),
		Path:  path,
	}
}

// LoadPosts loads all blog posts from the blog directory
func (b *Blog) LoadPosts() error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(b.Path, 0755); err != nil {
		return err
	}

	// Read all markdown files in the directory
	err := filepath.Walk(b.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			// For now, we'll just register that the file exists
			// Actual parsing will be implemented later
			fmt.Printf("Found post: %s\n", path)
		}

		return nil
	})

	return err
}
