package blog

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Post represents a blog post
type Post struct {
	Title       string         // Extracted from the first # Heading
	Content     string         // Extracted from the markdown file
	ContentHTML template.HTML  // HTML version of the content
	CreatedDate time.Time      // Parsed from frontmatter
	UpdatedDate time.Time      // Parsed from frontmatter
	Slug        string         // Derived from title
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
	return filepath.Walk(b.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			post, err := parsePost(path)
			if err != nil {
				return fmt.Errorf("error parsing post %s: %w", path, err)
			}
			b.Posts = append(b.Posts, post)
			fmt.Printf("Loaded post: %s\n", post.Title)
		}

		return nil
	})
}

// parsePost parses a markdown file into a Post struct
func parsePost(filePath string) (Post, error) {
	post := Post{}

	file, err := os.Open(filePath)
	if err != nil {
		return post, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	frontmatter := make(map[string]string)
	contentLines := make([]string, 0)

	// Parse frontmatter - files start with frontmatter, then "---", then content
	inFrontmatter := true

	for scanner.Scan() {
		line := scanner.Text()

		// Check for frontmatter separator
		if line == "---" {
			if inFrontmatter {
				// End of frontmatter, start of content
				inFrontmatter = false
				continue
			}
		}

		// Process frontmatter or content
		if inFrontmatter {
			// Parse frontmatter key: value
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				frontmatter[key] = value
			}
		} else {
			// Process content - now it will get here
			contentLines = append(contentLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return post, err
	}

	// Process frontmatter values
	if title, ok := frontmatter["Title"]; ok {
		post.Title = title
	} else {
		// Extract title from first # heading in content
		post.Title = extractTitleFromContent(contentLines)
	}

	if createdStr, ok := frontmatter["CreatedDate"]; ok {
		if created, err := time.Parse("2006-01-02 15:04:05", createdStr); err == nil {
			post.CreatedDate = created
		}
	}

	if updatedStr, ok := frontmatter["UpdatedDate"]; ok {
		if updated, err := time.Parse("2006-01-02 15:04:05", updatedStr); err == nil {
			post.UpdatedDate = updated
		}
	}

	if tagsStr, ok := frontmatter["Tags"]; ok {
		tags := strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		post.Tags = tags
	}

	if publishedStr, ok := frontmatter["Published"]; ok {
		post.Published = strings.ToLower(publishedStr) == "true"
	}

	// Derive slug from title
	post.Slug = createSlug(post.Title)

	// Set content
	post.Content = strings.Join(contentLines, "\n")

	return post, nil
}

// extractTitleFromContent extracts the first # heading from content lines
func extractTitleFromContent(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return "Untitled"
}

// createSlug creates a URL-friendly slug from a title
func createSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	// Replace spaces and special characters with hyphens
	slug = strings.Join(strings.Fields(slug), "-")
	// Remove any remaining non-alphanumeric characters (except hyphens)
	slug = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			return r
		}
		return -1
	}, slug)
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}