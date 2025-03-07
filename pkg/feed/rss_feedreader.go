package feed

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mouuff/TrendView/pkg/model"
	"golang.org/x/text/encoding/charmap"
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

// RssFeedReader is a feed provider that fetches RSS feeds
type RssFeedReader struct {
	Url             string
	ShouldCleanHtml bool
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

// cleanSpecialChars removes special characters
func cleanSpecialChars(text string) string {
	text = strings.ReplaceAll(text, "’", "'")
	text = strings.ReplaceAll(text, "‘", "'")
	text = strings.ReplaceAll(text, "–", "-")
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

func (p *RssFeedReader) GetSource() string {
	return p.Url
}

func (p *RssFeedReader) GetFeedItems() ([]model.FeedItem, error) {
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
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "iso-8859-1":
			return charmap.ISO8859_1.NewDecoder().Reader(input), nil
		case "utf-8":
			return input, nil // UTF-8 is handled by default
		default:
			return nil, fmt.Errorf("unsupported charset: %s", charset)
		}
	}

	err = decoder.Decode(&rss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %v", err)
	}

	reports := make([]model.FeedItem, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		parsedDate, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Printf("Warning: Failed to parse date '%s': %v. Skipping this item.\n", item.PubDate, err)
			continue
		}

		var source string
		title := cleanSpecialChars(item.Title)
		description := cleanSpecialChars(item.Description)

		if p.ShouldCleanHtml {
			description = cleanHTML(description)
		}

		if item.Link != "" {
			url, err := url.Parse(item.Link)
			if err != nil {
				log.Printf("Warning: Failed to parse link '%s': %v.\n", item.Link, err)
			} else {
				source = url.Hostname()
			}
		}

		report := model.FeedItem{
			Title:    title,
			Content:  description,
			DateTime: *parsedDate,
			Link:     item.Link,
			GUID:     item.GUID,
			Source:   source,
		}
		reports = append(reports, report)
	}

	return reports, nil
}
