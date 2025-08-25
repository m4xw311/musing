package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/m4xw311/musing/internal/blog"
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

	// Load blog posts
	b := blog.NewBlog(s.PostsDir)
	if err := b.LoadPosts(); err != nil {
		return fmt.Errorf("error loading posts: %w", err)
	}

	// Generate index page
	if err := s.generateIndex(b); err != nil {
		return fmt.Errorf("error generating index: %w", err)
	}

	// Generate individual post pages
	if err := s.generatePosts(b); err != nil {
		return fmt.Errorf("error generating posts: %w", err)
	}

	return nil
}

// generateIndex creates the index page with a list of all posts
func (s *StaticSiteGenerator) generateIndex(b *blog.Blog) error {
	indexTemplate := `<!DOCTYPE html>
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
        <h2>Posts</h2>
        <ul>
        {{range .Posts}}
            <li><a href="{{.Slug}}.html">{{.Title}}</a> - {{.CreatedDate.Format "2006-01-02"}}</li>
        {{end}}
        </ul>
    </main>
</body>
</html>`

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(s.OutputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b)
}

// generatePosts creates individual HTML pages for each post
func (s *StaticSiteGenerator) generatePosts(b *blog.Blog) error {
	postTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - My Blog</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <header>
        <h1><a href="index.html">My Blog</a></h1>
    </header>
    <main>
        <article>
            <h1>{{.Title}}</h1>
            <p><em>Published: {{.CreatedDate.Format "2006-01-02"}}</em></p>
            <div>
                {{.ContentHTML}}
            </div>
        </article>
    </main>
</body>
</html>`

	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return err
	}

	for _, post := range b.Posts {
		// Convert markdown content to HTML
		htmlContent := markdown.ToHTML([]byte(post.Content), nil, nil)
		post.ContentHTML = template.HTML(htmlContent)

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