package site

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// StaticSiteGenerator handles generating static HTML from markdown posts
type StaticSiteGenerator struct {
	PostsDir  string
	OutputDir string
}

// NewStaticSiteGenerator creates a new static site generator
func NewStaticSiteGenerator(postsDir, outputDir string) *StaticSiteGenerator {
	return &StaticSiteGenerator{
		PostsDir:  postsDir,
		OutputDir: outputDir,
	}
}

// Generate generates the static site
func (s *StaticSiteGenerator) Generate() error {
	// Create output directory
	if err := os.MkdirAll(s.OutputDir, 0755); err != nil {
		return err
	}

	// TODO: Rather than hard coding here, use a template engine like Go's text/template or html/template
	// 		 Maintain templates in separate files
	//       Maintain template for home page, posts page and error page
	// Create a simple index.html
	indexContent := `<!DOCTYPE html>
<html>
<head>
    <title>My Blog</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <header>
        <h1>My Blog</h1>
    </header>
    <main>
        <p>Welcome to my blog!</p>
    </main>
</body>
</html>`

	return ioutil.WriteFile(filepath.Join(s.OutputDir, "index.html"), []byte(indexContent), 0644)
}

// TODO: Use Blog in blog package rather than reading files directly here
// ProcessPosts processes all markdown posts
func (s *StaticSiteGenerator) ProcessPosts() error {
	return filepath.Walk(s.PostsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			fmt.Printf("Processing post: %s\n", path)
			// Implementation will be added later
		}

		return nil
	})
}
