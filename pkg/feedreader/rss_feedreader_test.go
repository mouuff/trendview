package feedreader_test

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/mouuff/TrendView/pkg/feedreader"
	"github.com/stretchr/testify/assert"
)

func TestGetFeedItems(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockRSS := `
	<rss>
		<channel>
			<item>
				<title>Test Title 1</title>
				<description>Test Description 1</description>
				<pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
				<link>http://example.com/1</link>
				<guid>1</guid>
			</item>
			<item>
				<title>Test Title 2</title>
				<description>Test Description 2</description>
				<pubDate>Tue, 03 Jan 2006 15:04:05 MST</pubDate>
				<link>http://example.com/2</link>
				<guid>2</guid>
			</item>
		</channel>
	</rss>`

	httpmock.RegisterResponder("GET", "https://news.google.com/rss/search?q=test&hl=en-US&gl=US&ceid=US:en",
		httpmock.NewStringResponder(200, mockRSS))

	provider := feedreader.NewGoogleRssFeedReader("test")
	items, err := provider.GetFeedItems()

	assert.NoError(t, err)
	assert.Len(t, items, 2)

	expectedDate1, _ := time.Parse(time.RFC1123, "Mon, 02 Jan 2006 15:04:05 MST")
	expectedDate2, _ := time.Parse(time.RFC1123, "Tue, 03 Jan 2006 15:04:05 MST")

	assert.Equal(t, "Test Title 1", items[0].Title)
	assert.Equal(t, "Test Description 1", items[0].Content)
	assert.Equal(t, expectedDate1, items[0].DateTime)
	assert.Equal(t, "http://example.com/1", items[0].Link)
	assert.Equal(t, "1", items[0].GUID)

	assert.Equal(t, "Test Title 2", items[1].Title)
	assert.Equal(t, "Test Description 2", items[1].Content)
	assert.Equal(t, expectedDate2, items[1].DateTime)
	assert.Equal(t, "http://example.com/2", items[1].Link)
	assert.Equal(t, "2", items[1].GUID)
}

func TestGetFeedItemsWithInvalidDate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockRSS := `
	<rss>
		<channel>
			<item>
				<title>Test Title 1</title>
				<description>Test Description 1</description>
				<pubDate>Invalid Date</pubDate>
				<link>http://example.com/1</link>
				<guid>1</guid>
			</item>
		</channel>
	</rss>`

	httpmock.RegisterResponder("GET", "https://news.google.com/rss/search?q=test&hl=en-US&gl=US&ceid=US:en",
		httpmock.NewStringResponder(200, mockRSS))

	provider := feedreader.NewGoogleRssFeedReader("test")
	items, err := provider.GetFeedItems()

	assert.NoError(t, err)
	assert.Len(t, items, 0)
}
