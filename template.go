package cronv

import (
	"fmt"
	"text/template"
	"time"
)

const TEMPLATE = `
<html>
<head>
<title>{{.Opts.Title}} | {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</title>
<script src="http://visjs.org/dist/vis.js"></script>
<link href="http://visjs.org/dist/vis-timeline-graph2d.min.css" rel="stylesheet" type="text/css" />
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
</head>
<body>
  <div class="container-fluid">
    <h1>{{.Opts.Title}}</h1>
    <p>From {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</p>
    <div id="cronv-timeline" style="height:100%; width:100%;"></div>
  </div>
  <script type="text/javascript">
		var groupCache = {},
				groups = new vis.DataSet(),
				items = new vis.DataSet(),
				itemId = 1;
		{{ range $index, $cronv := .CronEntries }}
			{{ $job := JSEscapeString $cronv.Crontab.Job }}
			var itemSize = 0,
					jobId = '{{ Md5Sum $job }}';
			{{ range CronvIter $cronv }}
				{{ $job := JSEscapeString $cronv.Crontab.Job }}
				{{ $startFormatted := DateFormat .Start "2006-01-02 15:04" }}
				items.add({
					id: itemId++,
					group: jobId,
					start: '{{ $startFormatted }}'
				});
				itemSize++;
			{{ end }}
			if (itemSize > 0) { // TODO add 'show all tasks' option
				if (!groupCache[jobId]) {
					groups.add({id: jobId, content: '{{ $job }}'});
					groupCache[jobId] = true;
				}
			}
		{{ end }}
		var options = {
	    showCurrentTime: true,
	    start: '{{DateFormat .TimeFrom "2006/1/2 15:04"}}',
	    end: '{{DateFormat .TimeTo "2006/1/2 15:04"}}',
	    zoomMax: {{ .DurationMinutes }} * 60 * 1000,
			stack: false
	  };
		new vis.Timeline(document.getElementById('cronv-timeline'), items, options).setGroups(groups);
  </script>
</body>
</html>
`

func MakeTemplate() *template.Template {
	funcMap := template.FuncMap{
		"CronvIter": func(cronv *Cronv) <-chan *Exec {
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
		"Md5Sum": func(data string) string {
			return Md5Sum(data)
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(TEMPLATE))
}
