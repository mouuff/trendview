# TrendView
View trends using AI

`go run .\cmd\trendview generate-trend -config .\data\bitcoin-confidence-configuration.json -datafile .\data\bitcoin-confidence-data.json -loop`
`go run .\cmd\trendview convert-to-html -datafile .\data\bitcoin-confidence-data.json > .\data\view.html; .\data\view.html`
`go run .\cmd\trendview convert-to-csv -datafile .\data\bitcoin-confidence-data.json > .\data\bitcoin-confidence-data.csv`