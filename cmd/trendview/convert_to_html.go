package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mouuff/TrendView/pkg/itemstore"
)

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type ConvertToHtml struct {
	flagSet *flag.FlagSet

	datafile   string
	identifier string
}

// Name gets the name of the command
func (cmd *ConvertToHtml) Name() string {
	return "convert-to-html"
}

// Init initializes the command
func (cmd *ConvertToHtml) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.datafile, "datafile", "", "file used to load data (required)")
	cmd.flagSet.StringVar(&cmd.identifier, "id", "", "the identifier for the Y axis, as configured in the RatingPrompt (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *ConvertToHtml) Run() error {
	if cmd.datafile == "" {
		log.Println("Please specify a data file using -datafile (e.g. -datafile data.json)")
		return errors.New("-datafile parameter required")
	}
	if cmd.identifier == "" {
		log.Println("Please specify a identifier using -id (e.g. -id MicrosoftConfidence)")
		return errors.New("-id parameter required")
	}

	storage := &itemstore.JsonItemStore{
		Filename: cmd.datafile,
	}

	data, err := storage.Load()

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	fmt.Println(getHtml(cmd.identifier, string(bytes)))

	return nil
}

func getHtml(identifier, jsonContent string) string {
	baseHtml := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Interactive Articles Timeline</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background: #f2f2f2;
      margin: 0;
      padding: 0;
    }
    .container {
      max-width: 1000px;
      margin: 30px auto;
      background: #fff;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    h1 {
      text-align: center;
      color: #333;
    }
    canvas {
      display: block;
      width: 100% !important;
      height: auto !important;
    }
  </style>
  <!-- Load Chart.js and date adapter -->
  <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-date-fns"></script>
</head>
<body>
  <div class="container">
    <h1>Interactive Articles Timeline</h1>
    <canvas id="myChart"></canvas>
  </div>
  <script>
    // Original JSON data (input schema unchanged)
    const dataSet = {{jsonContent}};
    const identifier = "{{identifier}}";

    // Convert the JSON object into an array and sort by DateTime
    const sortedEntries = Object.values(dataSet).sort((a, b) => new Date(a.DateTime) - new Date(b.DateTime));
    // const entries = Object.values(sortedEntries).filter((a) => a.Results[identifier].rating != 50);

    const entries = sortedEntries;

    // Get unique sources
    const sources = [...new Set(entries.map(item => item.Source))];

    // Create a dropdown for filtering by source
    const filterContainer = document.createElement('div');
    filterContainer.style.textAlign = 'center';
    filterContainer.style.marginBottom = '20px';
    const filterLabel = document.createElement('label');
    filterLabel.textContent = 'Filter by Source: ';
    const filterSelect = document.createElement('select');
    filterSelect.innerHTML = '<option value="all">All</option>' + sources.map(function(source) { return '<option value="' + source + '">' + source + '</option>'; }).join('');
    filterContainer.appendChild(filterLabel);
    filterContainer.appendChild(filterSelect);
    document.querySelector('.container').insertBefore(filterContainer, document.querySelector('canvas'));

    // Prepare chart data points (using Rating as the y-value)
    const prepareChartData = (source) => {
      const oneMonthAgo = new Date();
      oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1);
      const filteredEntries = entries.filter(item => new Date(item.DateTime) >= oneMonthAgo);
      const sourceFilteredEntries = source === 'all' ? filteredEntries : filteredEntries.filter(item => item.Source === source);
      return sourceFilteredEntries.map(item => ({
      x: new Date(item.DateTime),
      y: item.Results[identifier].rating,
      title: item.Title,
      content: item.Content,
      link: item.Link
      }));
    };

    let chartData = prepareChartData('all');

    // Create the Chart.js scatter plot with a time x-axis
    const ctx = document.getElementById('myChart').getContext('2d');
    let myChart = new Chart(ctx, {
      type: 'scatter',
      data: {
      datasets: [{
        label: 'Articles Rating',
        data: chartData,
        backgroundColor: 'rgba(75, 192, 192, 0.7)',
        borderColor: 'rgba(75, 192, 192, 1)',
        pointRadius: 6,
        pointHoverRadius: 8
      }]
      },
      options: {
      responsive: true,
      plugins: {
        tooltip: {
        callbacks: {
          title: (tooltipItems) => {
          return tooltipItems[0].raw.title;
          },
          label: (tooltipItem) => {
          const dateStr = new Date(tooltipItem.raw.x).toLocaleString();
          return 'Date: ' + dateStr + '\nRating: ' + tooltipItem.raw.y;
          },
          afterBody: (tooltipItems) => {
          return 'Content: ' + tooltipItems[0].raw.content;
          }
        }
        }
      },
      scales: {
        x: {
        type: 'time',
        time: {
          unit: 'minute'
        },
        title: {
          display: true,
          text: 'Date'
        },
        ticks: {
          color: '#333'
        },
        grid: {
          color: '#ccc'
        }
        },
        y: {
        title: {
          display: true,
          text: 'Rating'
        },
        ticks: {
          color: '#333'
        },
        grid: {
          color: '#ccc'
        }
        }
      },
      // Open the article link when a point is clicked
      onClick: (e, activeElements) => {
        if (activeElements.length) {
        const idx = activeElements[0].index;
        const point = chartData[idx];
        window.open(point.link, '_blank');
        }
      }
      }
    });

    // Update chart data when the filter changes
    filterSelect.addEventListener('change', (e) => {
      const selectedSource = e.target.value;
      chartData = prepareChartData(selectedSource);
      myChart.data.datasets[0].data = chartData;
      myChart.update();
    });
  </script>
</body>
</html>`
	result := baseHtml
	result = strings.ReplaceAll(result, "{{jsonContent}}", jsonContent)
	result = strings.ReplaceAll(result, "{{identifier}}", identifier)
	return result
}
