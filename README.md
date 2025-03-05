# TrendView
View trends using AI

`go run .\cmd\trendview generate -config .\config\msft-rating-configuration.json -datafile .\data\msft-rating-data.json -loop`

`go run .\cmd\trendview convert-to-html -datafile .\data\msft-rating-data.json -id MicrosoftConfidence > .\data\view.html; .\data\view.html`

`go run .\cmd\trendview convert-to-csv -datafile .\data\msft-rating-data.json > .\data\msft-rating-data.csv`

`go clean -testcache; go test ./...`