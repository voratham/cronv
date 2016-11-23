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
<link href="http://visjs.org/dist/vis.min.css" rel="stylesheet" type="text/css" />
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
<style>
body,html {
  background: linear-gradient(180deg, #172429,#43565E);
	color: #d3d3d3;
}
.header {
  color: white;
}
.vis-labelset .vis-label, .vis-time-axis .vis-text {
	color: #d3d3d3;
}
.control {
	margin-top: 25px;
}
#cronv-timeline {
	height:100%;
	width:100%;
	margin-top:5px;
}
</style>
</head>
<body>
  <div class="container-fluid">
    <h1 class="header">
			{{ .Opts.Title }}
    	<small>From {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</small>
		</h1>
		<div class="row control">
			<div class="col-sm-12">
				<button
					type="button"
					class="btn btn-primary btn-sm fit-tl"
					onclick="(function(){window.tl.fit()})();">Fit</button>
			</div>
		</div>
		<div class="row">
    	<div class="col-sm-12" id="cronv-timeline" style=""></div>
		</div>
  </div>
  <script type="text/javascript">
		var groupCache = {},
				groups = new vis.DataSet(),
				items = new vis.DataSet(),
				itemId = 1;
		{{ $shorten := AvgJobNameLen .CronEntries }}
		{{ range $index, $cronv := .CronEntries }}
			{{ $job := JSEscapeString $cronv.Crontab.Job }}
			var itemSize = 0,
					jobId = '{{ Md5Sum $job }}';
			{{ range CronvIter $cronv }}
				{{ $startFormatted := DateFormat .Start "2006-01-02 15:04" }}
				items.add({
					id: itemId++,
					group: jobId,
					start: '{{ $startFormatted }}'
				});
				itemSize++;
			{{ end }}
			if (itemSize > 0) { // TODO add 'show all tasks' option
				var groupRef = groupCache[jobId];
				if (!groupRef) {
					{{ $jobNameShort := Shorten $cronv.Crontab.Job $shorten "..." }}
					var g = {id: jobId, content: '{{ JSEscapeString $jobNameShort }}', itemSize: itemSize};
					groups.add(g);
					groupCache[jobId] = g;
				} else {
					groupRef.itemSize = groupRef.itemSize + itemSize;
				}
			}
		{{ end }}
		var options = {
	    showCurrentTime: true,
	    start: '{{DateFormat .TimeFrom "2006/1/2 15:04"}}',
	    end: '{{DateFormat .TimeTo "2006/1/2 15:04"}}',
	    zoomMax: {{ .DurationMinutes }} * 60 * 1000,
			stack: false,
			margin: {
				item: 30
			},
			groupOrder: function(a, b) {
				return a.itemSize - b.itemSize;
			}
	  };
		window.tl = new vis.Timeline(document.getElementById('cronv-timeline'), items, options);
		window.tl.setGroups(groups);
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
		"AvgJobNameLen": func(cronv []*Cronv) int {
			s := len(cronv)
			if s == 0 {
				return 0
			}
			i := 0
			for _, c := range cronv {
				i += len([]rune(c.Crontab.Job))
			}
			v := 0
			if i%s > 0 {
				v = 1
			}
			// FIXME devide by active tasks
			return i/s + v
		},
		"Shorten": func(v string, size int, suffix string) string {
			return Shorten(v, size, suffix)
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(TEMPLATE))
}
