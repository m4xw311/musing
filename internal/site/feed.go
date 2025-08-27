package site

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/m4xw311/musing/internal/blog"
)

// RSSFeed represents an RSS feed
type RSSFeed struct {
	XMLName xml.Name    `xml:"rss"`
	Version string      `xml:"version,attr"`
	Channel *RSSChannel `xml:"channel"`
}

// RSSChannel represents an RSS channel
type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language,omitempty"`
	PubDate     string    `xml:"pubDate,omitempty"`
	Items       []RSSItem `xml:"item"`
}

// RSSItem represents an RSS item
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

// AtomFeed represents an Atom feed
type AtomFeed struct {
	XMLName  xml.Name    `xml:"feed"`
	Xmlns    string      `xml:"xmlns,attr"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle,omitempty"`
	ID       string      `xml:"id"`
	Updated  string      `xml:"updated"`
	Author   *AtomAuthor `xml:"author,omitempty"`
	Entries  []AtomEntry `xml:"entry"`
}

// AtomAuthor represents an Atom author
type AtomAuthor struct {
	Name string `xml:"name"`
}

// AtomEntry represents an Atom entry
type AtomEntry struct {
	Title   string      `xml:"title"`
	ID      string      `xml:"id"`
	Link    AtomLink    `xml:"link"`
	Updated string      `xml:"updated"`
	Summary string      `xml:"summary"`
	Content AtomContent `xml:"content"`
}

// AtomLink represents an Atom link
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

// AtomContent represents Atom content
type AtomContent struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

// generateRSSFeed creates an RSS feed from blog posts
func (s *StaticSiteGenerator) generateRSSFeed(b *blog.Blog) error {
	channel := &RSSChannel{
		Title:       "My Blog",               // This should be configurable
		Link:        "http://localhost:8080", // This should be configurable
		Description: "A blog about technology and programming",
		Language:    "en-us",
	}

	// Add items for each post
	for _, post := range b.Posts {
		if !post.Published {
			continue
		}

		pubDate := post.CreatedDate.Format(time.RFC1123Z)
		if channel.PubDate == "" {
			channel.PubDate = pubDate
		}

		item := RSSItem{
			Title:       post.Title,
			Link:        fmt.Sprintf("http://localhost:8080/%s.html", post.Slug), // This should be configurable
			Description: string(post.ContentSnippetHTML),
			PubDate:     pubDate,
			GUID:        fmt.Sprintf("http://localhost:8080/%s.html", post.Slug), // This should be configurable
		}

		channel.Items = append(channel.Items, item)
	}

	rss := &RSSFeed{
		Version: "2.0",
		Channel: channel,
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.OutputDir, "rss.xml")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add XML header
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}

	_, err = file.Write(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated RSS feed: %s\n", filePath)
	return nil
}

// generateAtomFeed creates an Atom feed from blog posts
func (s *StaticSiteGenerator) generateAtomFeed(b *blog.Blog) error {
	feed := &AtomFeed{
		Xmlns:    "http://www.w3.org/2005/Atom",
		Title:    "My Blog", // This should be configurable
		Subtitle: "A blog about technology and programming",
		ID:       "http://localhost:8080", // This should be configurable
	}

	// Add entries for each post
	for _, post := range b.Posts {
		if !post.Published {
			continue
		}

		updated := post.UpdatedDate.Format(time.RFC3339)
		if feed.Updated == "" {
			feed.Updated = updated
		}

		entry := AtomEntry{
			Title:   post.Title,
			ID:      fmt.Sprintf("http://localhost:8080/%s.html", post.Slug),                 // This should be configurable
			Link:    AtomLink{Href: fmt.Sprintf("http://localhost:8080/%s.html", post.Slug)}, // This should be configurable
			Updated: updated,
			Summary: string(post.ContentSnippetHTML),
			Content: AtomContent{
				Type: "html",
				Body: string(post.ContentHTML),
			},
		}

		feed.Entries = append(feed.Entries, entry)
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.OutputDir, "atom.xml")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add XML header
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}

	_, err = file.Write(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated Atom feed: %s\n", filePath)
	return nil
}
