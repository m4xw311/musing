// Package site provides functionality for generating static websites from blog posts.
//
// The site package handles the generation of static HTML files from markdown blog posts,
// including creating index pages, individual post pages, and RSS/Atom feeds. It also
// manages copying assets like CSS files and images to the output directory.
package site

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/m4xw311/musing/internal/blog"
)

// StaticSiteGenerator handles generating static HTML from markdown posts.
type StaticSiteGenerator struct {
	PostsDir  string
	OutputDir string
}

// NewStaticSiteGenerator creates a new static site generator with the specified
// posts directory and output directory.
func NewStaticSiteGenerator(postsDir, outputDir string) *StaticSiteGenerator {
	return &StaticSiteGenerator{
		PostsDir:  postsDir,
		OutputDir: outputDir,
	}
}

// IndexData holds the data for the index page.
type IndexData struct {
	Posts       []blog.Post
	LatestPosts []blog.Post
}

// Generate generates the complete static site.
// It creates the output directory, loads blog posts, copies assets,
// generates the index page, individual post pages, and RSS/Atom feeds.
func (s *StaticSiteGenerator) Generate() error {
	// Create output directory
	if err := os.MkdirAll(s.OutputDir, 0755); err != nil {
		return err
	}

	// Load blog posts
	b := blog.NewBlog(s.PostsDir)
	if err := b.LoadPosts(); err != nil {
		return fmt.Errorf("error loading posts: %w", err)
	}

	// Copy style.css to the output directory
	if err := s.copyStyleCSS(); err != nil {
		return fmt.Errorf("error copying style.css: %w", err)
	}

	// Copy images directory to the output directory
	if err := s.copyImages(); err != nil {
		return fmt.Errorf("error copying images: %w", err)
	}

	// Generate index page
	if err := s.generateIndex(b); err != nil {
		return fmt.Errorf("error generating index: %w", err)
	}

	// Generate individual post pages
	if err := s.generatePosts(b); err != nil {
		return fmt.Errorf("error generating posts: %w", err)
	}

	// Generate RSS feed
	if err := s.generateRSSFeed(b); err != nil {
		return fmt.Errorf("error generating RSS feed: %w", err)
	}

	// Generate Atom feed
	if err := s.generateAtomFeed(b); err != nil {
		return fmt.Errorf("error generating Atom feed: %w", err)
	}

	return nil
}

// copyStyleCSS copies the style.css file to the output directory.
func (s *StaticSiteGenerator) copyStyleCSS() error {
	src := "internal/template/style.css"
	dst := filepath.Join(s.OutputDir, "style.css")

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// copyImages copies the images directory from posts to the output directory.
func (s *StaticSiteGenerator) copyImages() error {
	src := filepath.Join(s.PostsDir, "images")
	dst := filepath.Join(s.OutputDir, "images")

	// Check if source directory exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		// No images directory, nothing to copy
		return nil
	} else if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Copy all files in the images directory
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create destination file path
		dstPath := filepath.Join(dst, relPath)

		// Create destination directory if it doesn't exist
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return err
		}

		// Copy file
		return s.copyFile(path, dstPath)
	})
}

// copyFile copies a file from src to dst.
func (s *StaticSiteGenerator) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// generateIndex creates the index page with a list of all posts.
func (s *StaticSiteGenerator) generateIndex(b *blog.Blog) error {
	// Prepare data for the index page
	var latestPosts []blog.Post
	if len(b.Posts) > 4 {
		latestPosts = b.Posts[:4]
	} else {
		latestPosts = b.Posts
	}

	indexData := IndexData{
		Posts:       b.Posts,
		LatestPosts: latestPosts,
	}

	tmpl, err := template.ParseFiles("internal/template/index.html")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(s.OutputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, indexData)
}

// generatePosts creates individual HTML pages for each post.
func (s *StaticSiteGenerator) generatePosts(b *blog.Blog) error {
	tmpl, err := template.ParseFiles("internal/template/post.html")
	if err != nil {
		return err
	}

	for _, post := range b.Posts {
		file, err := os.Create(filepath.Join(s.OutputDir, post.Slug+".html"))
		if err != nil {
			return err
		}

		err = tmpl.Execute(file, post)
		file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}