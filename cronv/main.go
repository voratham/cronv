package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/takumakanari/cronv"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const TEMPLATE = `
<html>
<head>
<title>cronv | {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Duration}}</title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
</head>
<body>
  <div class="container-fluid">
    <h1>cronv</h1>
    <p>From {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Duration}}</p>
    <div id="cronv-timeline" style="height:100%; width:100%;">
      <b>Loading...</b>
    </div>
  </div>
  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
  <script type="text/javascript">
     google.charts.load("current", {packages:["timeline"]});
     google.charts.setOnLoadCallback(function() {
       var container = document.getElementById('cronv-timeline');
       var chart = new google.visualization.Timeline(container);
       var dataTable = new google.visualization.DataTable();
        dataTable.addColumn({ type: 'string', id: 'job' });
        dataTable.addColumn({ type: 'string', id: 'dummy bar label' });
        dataTable.addColumn({ type: 'string', role: 'tooltip' });
        dataTable.addColumn({ type: 'date', id: 'Start' });
        dataTable.addColumn({ type: 'date', id: 'End' });
        var rows = [
          {{range $index, $cronv := .CronEntries}}
            {{range CronvIter $cronv}}
              {{ $job := JSEscapeString $cronv.Crontab.Job }}
              {{ $startFormatted := DateFormat .Start "15:04" }}
              ['{{$job}}', '', '{{$startFormatted}} {{$job}}', {{NewJsDate .Start}}, {{NewJsDate .End}}],
            {{end}}
          {{end}}
        ];
        if (rows.length > 0) {
          dataTable.addRows(rows);
          chart.draw(dataTable, {
            timeline: {
              colorByRowLabel: true,
            },
            avoidOverlappingGridLines: false  
          });
        } else {
          container.innerHTML = '<div class="alert alert-success"><strong>Woops!</strong> There is no data!</div>';
        }
     });
  </script>
</body>
</html>
`

func makeTemplate() *template.Template {
	funcMap := template.FuncMap{
		"CronvIter": func(cronv *cronv.Cronv) <-chan *cronv.Exec {
			return cronv.Iter()
		},
		"JSEscapeString": func(v string) string {
			return template.JSEscapeString(v)
		},
		"NewJsDate": func(v time.Time) string {
			return fmt.Sprintf("new Date(%d,%d,%d,%d,%d)", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute())
		},
		"DateFormat": func(v time.Time, format string) string {
			return v.Format(format)
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(TEMPLATE))
}

func optimizeTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
}

func durationToMinutes(s string) (float64, error) {
	length := len(s)
	if length < 2 {
		return 0, errors.New(fmt.Sprintf("Invalid duration format: '%s'", s))
	}

	duration, err := strconv.Atoi(string(s[:length-1]))
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid duration format: '%s', %s", s, err))
	}

	unit := string(s[length-1])
	switch strings.ToLower(unit) {
	case "d":
		return float64(duration * 24 * 60), nil
	case "h":
		return float64(duration * 60), nil
	case "m":
		return float64(duration), nil
	}

	return 0, errors.New(fmt.Sprintf("Invalid duration format: '%s', '%s' is not in d/h/m", s, unit))
}

func main() {
	var (
		outputFilePath string
		duration       string
	)
	for _, f := range []string{"o", "output"} {
		flag.StringVar(&outputFilePath, f, "./crontab.html", "path/to/htmlfile to output.")
	}
	for _, f := range []string{"d", "duration"} {
		flag.StringVar(&duration, f, "6h", "duration to visualize in N{suffix} style. e.g.) 1d(day)/1h(hour)/1m(minute)")
	}
	flag.Parse()

	durationMinutes, err := durationToMinutes(duration)
	if err != nil {
		panic(err)
	}

	output, err := os.Create(outputFilePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to handle output file: %s", err))
	}

	cronEntries := []*cronv.Cronv{}
	scanner := bufio.NewScanner(os.Stdin)
	timeFrom := optimizeTime(time.Now())

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && string(line[0]) != "#" {
			cronv, err := cronv.NewCronv(line, timeFrom, durationMinutes)
			if err != nil {
				panic(fmt.Sprintf("Failed to analyze cron '%s': %s", line, err))
			}
			cronEntries = append(cronEntries, cronv)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	makeTemplate().Execute(output, map[string]interface{}{
		"CronEntries": cronEntries,
		"TimeFrom":    timeFrom,
		"Duration":    duration,
	})

	fmt.Printf("'%s' generated successfully.\n", outputFilePath)
}
