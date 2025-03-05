# TrendView
View trends using AI

`go run .\cmd\trendview generate -config .\data\bitcoin-rating-configuration.json -datafile .\data\bitcoin-rating-data.json -loop`
`go run .\cmd\trendview convert-to-html -datafile .\data\bitcoin-rating-data.json > .\data\view.html; .\data\view.html`
`go run .\cmd\trendview convert-to-csv -datafile .\data\bitcoin-rating-data.json > .\data\bitcoin-rating-data.csv`
`go clean -testcache; go test ./...`