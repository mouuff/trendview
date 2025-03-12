# TrendView
TrendView is an AI-powered tool designed to analyze and visualize trends effectively.

### Features:
- **Configurable Automation**: Retrieve RSS content and obtain ratings from your preferred AI.
- **Database Integration**: Store and manage results efficiently.
- **CSV Export**: Easily export data to CSV format.
- **Real-time Visualization**: Use the built-in web server to visualize results in real-time.

## Project Status

This project is currently in the prototype stage. It was built over a few days but is designed to be modular and extendable, allowing for future enhancements and integrations.

You can easily set up and run this project entirely locally.

You can also explore a live demo of TrendView hosted temporarily at: [TrendView Demo](http://79.160.3.137/).

## Example Results

Below is an example of the calculated confidence in Bitcoin based on recent news. The confidence score ranges from 0 to 100, with higher scores indicating a stronger positive outlook for investing in Bitcoin.

![Example Result](images/example_result_1.png)

## Prerequisites

Before using TrendView, ensure you have the following installed:

- **Golang**: Download and install Golang from [the official website](https://golang.org/dl/).
- **Ollama**: Download and install Ollama from [Ollama's website](https://ollama.com/).

To run Ollama in server mode, use the following command:

```sh
ollama serve
```

## Example Commands

### Generate Trends

To generate trends, use the following command:

```sh
go run .\cmd\trendview generate -config .\config\global-rating-configuration.json -datafile .\data\global.db -loop
```

### Run the Web Server Locally

To run the web server, use the following command:

```sh
go run .\cmd\trendview serve -datafile .\data\global.db
```

Once the server is running, you can access the endpoint on your local machine at:

[http://localhost:8081/](http://localhost:8081/)

### Convert to CSV

To convert the data to a CSV file, use the following command:

```sh
go run .\cmd\trendview convert-to-csv -datafile .\data\global.db > .\data\global.csv
```

### Run Unit Tests

To run unit tests, use the following command:

```sh
go clean -testcache; go test ./...
```

## Example configuration
```json
{
  "Model": "llama3.2",
  "RssFeedReaders": [
    {
      "Url": "https://www.theguardian.com/uk/technology/rss",
      "ShouldCleanHtml": true
    },
    {
      "Url": "https://bitcoinmagazine.com/feed",
      "ShouldCleanHtml": true
    },
    {
      "Url": "https://www.cnbc.com/id/19854910/device/rss/rss.html"
    },
    {
      "Url": "http://rss.cnn.com/rss/money_markets.rss"
    },
    {
      "Url": "http://rss.cnn.com/rss/money_technology.rss"
    },
    {
      "Url": "https://feeds.bloomberg.com/markets/news.rss"
    },
    {
      "Url": "https://feeds.bloomberg.com/technology/news.rss"
    },
    {
      "Url": "https://www.lemonde.fr/en/economy/rss_full.xml"
    }
  ],
  "RatingPrompts": [
    {
      "SubjectName": "Microsoft",
      "InsightName": "Confidence",
      "BasePrompt": "Based solely on the news provided below, give a rating on how it might affect Microsoft's stock price on a scale from 0 to 100, where: 0 indicates a very negative confidence (likely price drop), 50 indicates a neutral confidence (no significant change), and 100 indicates a positive confidence (likely price increase). For the rating, consider market trends, regulations, economic factors, and any other relevant information. News: "
    },
    {
      "SubjectName": "Microsoft",
      "InsightName": "Relevance",
      "BasePrompt": "Based solely on the news provided below, give a rating on how related it is to Microsoft on a scale from 0 to 100, where: 0 indicates no relevance (completely unrelated), 50 indicates moderate relevance (somewhat related), and 100 indicates high relevance (directly related). News:"
    },
    {
      "SubjectName": "Bitcoin",
      "InsightName": "Confidence",
      "BasePrompt": "Based solely on the news provided below, give a rating on how it might affect Bitcoin's price on a scale from 0 to 100, where: 0 indicates a very negative confidence (likely price drop), 50 indicates a neutral confidence (no significant change), and 100 indicates a positive confidence (likely price increase). For the rating, consider market trends, regulations, economic factors, and any other relevant information. News: "
    },
    {
      "SubjectName": "Bitcoin",
      "InsightName": "Relevance",
      "BasePrompt": "Based solely on the news provided below, give a rating on how related it is to Bitcoin on a scale from 0 to 100, where: 0 indicates no relevance (completely unrelated), 50 indicates moderate relevance (somewhat related), and 100 indicates high relevance (directly related). For the rating, consider the content of the news, its potential impact on Bitcoin, market perception, and any other relevant factors. News: "
    }
  ]
}
```

