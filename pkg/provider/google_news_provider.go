package provider

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

type GoogleNewsProvider struct {
	url string
}

func NewGoogleNewsProvider() *GoogleNewsProvider {
	return &GoogleNewsProvider{
		url: "https://news.google.com/rss/search?q=bitcoin&hl=en-US&gl=US&ceid=US:en",
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

func (p *GoogleNewsProvider) GetReports() ([]Report, error) {
	resp, err := http.Get(p.url)
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
		parsedDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("Warning: Failed to parse date '%s': %v. Using original string.\n", item.PubDate, err)
			reports = append(reports, Report{
				Title:    item.Title,
				Content:  cleanHTML(item.Description),
				DateTime: item.PubDate,
				Link:     item.Link,
				GUID:     item.GUID,
			})
			continue
		}

		report := Report{
			Title:    item.Title,
			Content:  cleanHTML(item.Description),
			DateTime: parsedDate.Format(time.RFC3339),
			Link:     item.Link,
			GUID:     item.GUID,
		}
		reports = append(reports, report)
	}

	return reports, nil
}
