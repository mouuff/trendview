package feedreader

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// RSS structs for XML parsing
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
}

type RssProvider struct {
	Url             string
	ShouldCleanHtml bool
}

func NewGoogleRssProvider(query string) *RssProvider {
	url := fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=en-US&gl=US&ceid=US:en", query)
	return &RssProvider{
		Url: url,
	}
}

// cleanHTML removes HTML tags and extra whitespace from text
func cleanHTML(html string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")

	// Clean up extra whitespace
	text = strings.Join(strings.Fields(text), " ")

	// Remove any remaining special markup
	text = strings.ReplaceAll(text, "&nbsp;", " ") // Replace non-breaking spaces
	return text
}

// parsePubDate attempts to parse RSS publication date strings
func parsePubDate(pubDate string) (*time.Time, error) {
	// Try RFC1123 first (handles GMT and other timezone names)
	if parsed, err := time.Parse(time.RFC1123, pubDate); err == nil {
		return &parsed, nil
	}

	// Fallback to RFC1123Z (handles numeric offsets)
	if parsed, err := time.Parse(time.RFC1123Z, pubDate); err == nil {
		return &parsed, nil
	}

	return nil, fmt.Errorf("failed to parse date '%s'", pubDate)
}

func (p *RssProvider) GetReports() ([]Report, error) {
	resp, err := http.Get(p.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rss Rss
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %v", err)
	}

	reports := make([]Report, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		parsedDate, err := parsePubDate(item.PubDate)
		if err != nil {
			fmt.Printf("Warning: Failed to parse date '%s': %v. Skipping this item.\n", item.PubDate, err)
			continue
		}

		description := item.Description
		if p.ShouldCleanHtml {
			description = cleanHTML(description)
		}

		report := Report{
			Title:    item.Title,
			Content:  description,
			DateTime: parsedDate,
			Link:     item.Link,
			GUID:     item.GUID,
		}
		reports = append(reports, report)
	}

	return reports, nil
}
