# TrendView
View trends using AI

`go run .\cmd\trendview generate -config .\config\global-rating-configuration.json -datafile .\data\global-rating-data.json -loop`

`go run .\cmd\trendview convert-to-html -datafile .\data\global-rating-data.json -id MicrosoftConfidence > .\data\view.html; .\data\view.html`

`go run .\cmd\trendview convert-to-csv -datafile .\data\global-rating-data.json > .\data\global-rating-data.csv`

`go clean -testcache; go test ./...`