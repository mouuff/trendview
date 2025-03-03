package provider

import (
	"encoding/xml"
	"fmt"
	"net/http"
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
}

type GoogleNewsProvider struct {
	url string
}

func NewGoogleNewsProvider() *GoogleNewsProvider {
	return &GoogleNewsProvider{
		url: "https://news.google.com/rss/search?q=bitcoin&hl=en-US&gl=US&ceid=US:en",
	}
}
func (p *GoogleNewsProvider) GetReports() ([]Report, error) {
	// Make HTTP request
	resp, err := http.Get(p.url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse XML
	var rss Rss
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %v", err)
	}

	// Convert to Reports
	reports := make([]Report, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		// Parse the publication date from the RSS feed
		parsedDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// Log the error and use current time as fallback, but this shouldn't happen with Google News RSS
			fmt.Printf("Warning: Failed to parse date '%s': %v. Using original string.\n", item.PubDate, err)
			reports = append(reports, Report{
				Title:    item.Title,
				Content:  item.Description,
				DateTime: item.PubDate, // Use original string if parsing fails
				Link:     item.Link,
			})
			continue
		}

		report := Report{
			Title:    item.Title,
			Content:  item.Description,
			DateTime: parsedDate.Format(time.RFC3339),
			Link:     item.Link,
		}
		reports = append(reports, report)
	}

	return reports, nil
}
