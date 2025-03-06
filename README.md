# TrendView

TrendView is a tool that leverages AI to analyze and visualize trends.

## Project Status

This project is currently in the prototype stage. It was built rapidly over a few days but is designed to be modular and extendable, allowing for future enhancements and integrations.

## Prerequisites

Before using TrendView, ensure you have the following installed:

- **Golang**: Download and install Golang from [the official website](https://golang.org/dl/).
- **Ollama**: Download and install Ollama from [Ollama's website](https://ollama.com/).

To run Ollama in server mode, use the following command:

```sh
ollama serve
```

## Example Results

Below is an example of the calculated confidence in Bitcoin based on recent news:

![Example Result](images/example_result_1.png)

## Example Commands

### Generate Trends

To generate trends, use the following command:

```sh
go run .\cmd\trendview generate -config .\config\global-rating-configuration.json -datafile .\data\global-rating-data.json -loop
```

### Convert to CSV

To convert the data to a CSV file, use the following command:

```sh
go run .\cmd\trendview convert-to-csv -datafile .\data\global-rating-data.json > .\data\global-rating-data.csv
```

### Run Unit Tests

To run unit tests, use the following command:

```sh
go clean -testcache; go test ./...
```

## Example configuration
```json
{
    "RssFeedReaders": [
        {
            "Url": "https://www.theguardian.com/uk/technology/rss",
            "ShouldCleanHtml": true
        },
        {
            "Url": "https://feeds.bloomberg.com/markets/news.rss"
        },
        {
            "Url": "https://feeds.bloomberg.com/technology/news.rss"
        }
    ],
    "RatingPrompts": [
      {
        "Identifier": "MicrosoftConfidence",
        "BasePrompt": "Based solely on the news below, rate your confidence in investing in Microsoft stocks from 0 (no confidence, unwise) to 50 (neutral) to 100 (very confident, good opportunity), considering market trends, regulations, or economic factors. News: "
      },
      {
        "Identifier": "MicrosoftRelated",
        "BasePrompt": "Based exclusively on the news provided below, evaluate the potential connection to Microsoft's stock price. Assign a rating on a scale from 0 to 100, where:  - 0 = completely unrelated - 50 = somewhat related - 100 = very much related If there is any uncertainty or insufficient information to determine relevance, default to a rating of 0. News: "
      }
    ]
}
```
