// Package blog provides functionality for managing blog posts.
//
// The blog package handles parsing markdown files into Post structs,
// managing collections of posts, and converting markdown to HTML with
// extended features like tables, footnotes, and syntax highlighting.
package blog

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Post represents a blog post with metadata and content.
type Post struct {
	Title              string        // Extracted from the first # Heading
	Content            string        // Extracted from the markdown file
	ContentHTML        template.HTML // HTML version of the content
	ContentSnippet     string        // Snippet of the content for cards
	ContentSnippetHTML template.HTML // HTML version of the content snippet
	CreatedDate        time.Time     // Parsed from frontmatter
	UpdatedDate        time.Time     // Parsed from frontmatter
	Slug               string        // Derived from title
	Tags               []string
	Published          bool
}

// Blog represents a collection of blog posts.
type Blog struct {
	Posts []Post
	Path  string
}

// NewBlog creates a new blog instance with the specified path.
func NewBlog(path string) *Blog {
	return &Blog{
		Posts: make([]Post, 0),
		Path:  path,
	}
}

// LoadPosts loads all blog posts from the blog directory.
// It reads all markdown files, parses them into Post structs,
// and sorts them by creation date (newest first).
func (b *Blog) LoadPosts() error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(b.Path, 0755); err != nil {
		return err
	}

	// Clear existing posts
	b.Posts = make([]Post, 0)

	// Read all markdown files in the directory
	err := filepath.Walk(b.Path, func(path string, info os.FileInfo, err error) error {
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

	if err != nil {
		return err
	}

	// Sort posts by CreatedDate in descending order (newest first)
	sort.Slice(b.Posts, func(i, j int) bool {
		return b.Posts[i].CreatedDate.After(b.Posts[j].CreatedDate)
	})

	return nil
}

// CustomMarkdownParser creates a markdown parser with extended features.
// It enables various markdown extensions including tables, footnotes,
// strikethrough, and MathJax support.
func CustomMarkdownParser() *parser.Parser {
	// Create markdown parser with extended features
	// Using all available extensions from gomarkdown for maximum compatibility
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.Footnotes |
		parser.Strikethrough |
		parser.SpaceHeadings |
		parser.HeadingIDs |
		parser.BackslashLineBreak |
		parser.DefinitionLists |
		parser.MathJax | parser.SuperSubscript

	// Create the parser with extensions
	p := parser.NewWithExtensions(extensions)
	return p
}

// CustomMarkdownRenderer creates a markdown renderer with extended features.
func CustomMarkdownRenderer() *html.Renderer {
	// Create markdown renderer with extended features
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	return html.NewRenderer(opts)
}

// ConvertMarkdownToHTML converts markdown text to HTML with extended features.
// It uses the custom parser and renderer with extended markdown support.
func ConvertMarkdownToHTML(content string) []byte {
	// Create parser and renderer with extensions
	parser := CustomMarkdownParser()
	renderer := CustomMarkdownRenderer()

	// Parse and render the markdown
	doc := parser.Parse([]byte(content))
	htmlContent := markdown.Render(doc, renderer)

	return htmlContent
}

// parsePost parses a markdown file into a Post struct.
// It extracts frontmatter metadata and converts the content to HTML.
func parsePost(filePath string) (Post, error) {
	post := Post{}

	_, err := os.Stat(filePath)
	if err != nil {
		return post, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return post, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	frontmatter := make(map[string]string)
	contentLines := make([]string, 0)
	var rawFrontmatterLines []string

	// Parse frontmatter - files start with frontmatter, then "---", then content
	inFrontmatter := false
	frontmatterProcessed := false

	// Read first line to check if it starts with ---
	if scanner.Scan() {
		firstLine := scanner.Text()
		if firstLine == "---" {
			inFrontmatter = true
			rawFrontmatterLines = append(rawFrontmatterLines, firstLine)
		} else {
			// No frontmatter, just content
			contentLines = append(contentLines, firstLine)
		}
	}

	// Continue processing the rest of the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check for frontmatter separator
		if line == "---" && inFrontmatter && !frontmatterProcessed {
			// End of frontmatter, start of content
			inFrontmatter = false
			frontmatterProcessed = true
			rawFrontmatterLines = append(rawFrontmatterLines, line)
			continue
		}

		// Process frontmatter or content
		if inFrontmatter {
			// Collect raw frontmatter lines
			rawFrontmatterLines = append(rawFrontmatterLines, line)
			// Parse frontmatter key: value
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				frontmatter[key] = value
			}
		} else {
			// Process content
			contentLines = append(contentLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return post, err
	}

	// Always use current time for missing dates to avoid any file timestamp issues
	defaultTime := time.Now()

	updatedFile := false

	// Not doing this to simplify the work of the user
	// if title, ok := frontmatter["Title"]; ok {
	// 	post.Title = title
	// }

	// Extract title from first # heading in content
	post.Title = extractTitleFromContent(contentLines)

	if createdStr, ok := frontmatter["CreatedDate"]; ok {
		if created, err := time.Parse("2006-01-02 15:04:05", createdStr); err == nil {
			post.CreatedDate = created
		} else {
			fmt.Fprintf(os.Stderr, "Invalid CreatedDate in %s. Must be in format '2006-01-02 15:04:05'\n", filePath)
		}
	} else {
		// If no CreatedDate, use current time
		post.CreatedDate = defaultTime
		updatedFile = true
		fmt.Printf("Missing CreatedDate in %s, setting to current time\n", filePath)
	}

	if updatedStr, ok := frontmatter["UpdatedDate"]; ok {
		if updated, err := time.Parse("2006-01-02 15:04:05", updatedStr); err == nil {
			post.UpdatedDate = updated
		} else {
			fmt.Fprintf(os.Stderr, "Invalid UpdatedDate in %s. Must be in format '2006-01-02 15:04:05'\n", filePath)
		}
	} else {
		// If no UpdatedDate, use current time
		post.UpdatedDate = post.CreatedDate
		updatedFile = true
		fmt.Printf("Missing UpdatedDate in %s, setting same as CreatedDate\n", filePath)
	}

	// Update the file ONLY if we added missing date fields
	// Do NOT update if both dates already existed and were valid
	if updatedFile {
		if err := updatePostFile(filePath, rawFrontmatterLines, frontmatter, contentLines, post); err != nil {
			fmt.Printf("Warning: Could not update %s with missing date fields: %v\n", filePath, err)
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
	// Convert markdown content to HTML with extended features
	htmlContent := ConvertMarkdownToHTML(post.Content)
	post.ContentHTML = template.HTML(htmlContent)

	// Generate content snippet
	post.ContentSnippet = generateContentSnippet(post.Content)
	// Convert markdown snippet to HTML with extended features
	htmlContentSnippet := ConvertMarkdownToHTML(post.ContentSnippet)
	post.ContentSnippetHTML = template.HTML(htmlContentSnippet)
	return post, nil
}

// updatePostFile updates the markdown file with missing date fields.
// It adds CreatedDate and UpdatedDate fields to the frontmatter if they're missing.
func updatePostFile(filePath string, rawFrontmatterLines []string, frontmatter map[string]string, contentLines []string, post Post) error {
	// Create the updated content
	var updatedContent strings.Builder

	// Check if we have frontmatter
	hasFrontmatter := len(rawFrontmatterLines) > 0
	hasValidFrontmatterStructure := hasFrontmatter && len(rawFrontmatterLines) >= 2 &&
		rawFrontmatterLines[0] == "---" && rawFrontmatterLines[len(rawFrontmatterLines)-1] == "---"

	if hasFrontmatter && hasValidFrontmatterStructure {
		// Check which date fields are missing
		hasCreatedDate := false
		hasUpdatedDate := false

		// Check existing frontmatter for date fields
		for key := range frontmatter {
			if key == "CreatedDate" {
				hasCreatedDate = true
			}
			if key == "UpdatedDate" {
				hasUpdatedDate = true
			}
		}

		// Write first frontmatter separator
		updatedContent.WriteString("---\n")

		// Write existing frontmatter lines (excluding separators)
		for i := 1; i < len(rawFrontmatterLines)-1; i++ {
			updatedContent.WriteString(rawFrontmatterLines[i] + "\n")
		}

		// Add missing date fields if needed
		if !hasCreatedDate {
			updatedContent.WriteString(fmt.Sprintf("CreatedDate: %s\n", post.CreatedDate.Format("2006-01-02 15:04:05")))
		}

		if !hasUpdatedDate {
			updatedContent.WriteString(fmt.Sprintf("UpdatedDate: %s\n", post.UpdatedDate.Format("2006-01-02 15:04:05")))
		}

		// Write closing frontmatter separator
		updatedContent.WriteString("---\n")
	} else if hasFrontmatter {
		// Frontmatter exists but missing proper separators
		// Check which date fields are missing
		hasCreatedDate := false
		hasUpdatedDate := false

		// Check existing frontmatter for date fields
		for key := range frontmatter {
			if key == "CreatedDate" {
				hasCreatedDate = true
			}
			if key == "UpdatedDate" {
				hasUpdatedDate = true
			}
		}

		// Write first frontmatter separator
		updatedContent.WriteString("---\n")

		// Write existing frontmatter lines
		for _, line := range rawFrontmatterLines {
			// Skip any accidental separators in the content
			if line != "---" {
				updatedContent.WriteString(line + "\n")
			}
		}

		// Add missing date fields if needed
		if !hasCreatedDate {
			updatedContent.WriteString(fmt.Sprintf("CreatedDate: %s\n", post.CreatedDate.Format("2006-01-02 15:04:05")))
		}

		if !hasUpdatedDate {
			updatedContent.WriteString(fmt.Sprintf("UpdatedDate: %s\n", post.UpdatedDate.Format("2006-01-02 15:04:05")))
		}

		// Write closing frontmatter separator
		updatedContent.WriteString("---\n")
	} else {
		// No existing frontmatter, create new one
		updatedContent.WriteString("---\n")
		updatedContent.WriteString(fmt.Sprintf("CreatedDate: %s\n", post.CreatedDate.Format("2006-01-02 15:04:05")))
		updatedContent.WriteString(fmt.Sprintf("UpdatedDate: %s\n", post.UpdatedDate.Format("2006-01-02 15:04:05")))
		updatedContent.WriteString("---\n")
	}

	// Write content
	for _, line := range contentLines {
		updatedContent.WriteString(line + "\n")
	}

	// Write the updated content back to the file
	return os.WriteFile(filePath, []byte(updatedContent.String()), 0644)
}

// extractTitleFromContent extracts the first # heading from content lines.
func extractTitleFromContent(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return "Untitled"
}

// createSlug creates a URL-friendly slug from a title.
// It converts the title to lowercase, replaces spaces with hyphens,
// and removes special characters.
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

// generateContentSnippet generates a short snippet from the content.
// Currently limits the content to first 150 characters.
func generateContentSnippet(content string) string {
	// Remove markdown headers and code blocks
	//	re := regexp.MustCompile("(?m)^(#.*|```.*|```)

	// Limit to first 150 characters
	if len(content) > 150 {
		content = content[:150] + "..."
	}

	return content
}