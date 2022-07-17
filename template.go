package cronv

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

const TEMPLATE = `
<html>
<head>
<title>{{.Opts.Title}} | {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</title>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>
<body>
  <div class="container-fluid">
    <h1>
			{{.Opts.Title}}
      <br/>
      <small class="text-muted" style="font-size:16px;">⏰ Assume start {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</small>
		</h1>
    


    {{if .Extras}}
      <h3>Extra</h3>
      <div id="cronv-extra" style="width:{{.Opts.Width}}%;">
        <dl class="row">
          {{range $index, $extra := .Extras}}
            <dt class="col-sm-1">{{$extra.Label}}</dt>
            <dd class="col-sm-11">{{$extra.Job}}</dd>
          {{end}}
        </dl>
      </div>
      <hr>
    {{end}}

    <div class="input-group mb-3">
      <input id="cronv-input-filter" type="text" class="form-control" placeholder="Enter Cronjob name" aria-label="Enter Cronjob name" aria-describedby="cron-job-name">
      <div class="input-group-append">
        <span class="input-group-text" id="cron-job-name">🕵️‍♂️</span>
      </div>
    </div>

    <h3>Timeline</h3>
    <div id="cronv-timeline" style="height:60vh; width:{{.Opts.Width}}%; padding-right:20px;">
      <b>Loading...</b>
    </div>

    <h3>🔴 Cronjob out of bounds 1 week</h3>
    <div id="cronv-cannot-render">
  </div>


  </div>
  <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
  <script type="text/javascript">
     google.charts.load("current", {packages:["timeline"]});
     google.charts.setOnLoadCallback(function() {
       const container = document.getElementById('cronv-timeline');

				var tasks = {};
        var tasksWithDayOfMonths = {};
				{{ $timeFrom := .TimeFrom }}
				{{ $timeTo := .TimeTo }}
        {{ $durationString := .Opts.Duration}}
				{{range $index, $cronv := .CronEntries}}
					{{ $job := JSEscapeString $cronv.Crontab.Job }}
					tasks['{{$job}}'] = tasks['{{$job}}'] || [];

					{{if IsRunningEveryMinutes $cronv.Crontab }}
						tasks['{{$job}}'].push(['{{$job}}', '', 'Every minutes {{$job}}', {{NewJsDate $timeFrom}}, {{NewJsDate $timeTo}}]);
					{{else if not (ShouldRenderCronTabInDayOfMonthLevel $cronv.Crontab $durationString) }}
              tasksWithDayOfMonths['{{$job}}'] = '{{ GetCronScheduleString $cronv.Crontab }}'
          {{ else }}
            {{range CronvIter $cronv}}tasks['{{$job}}'].push(['{{$job}}', '', '{{DateFormat .Start "15:04"}} {{$job}}', {{NewJsDate .Start}}, {{NewJsDate .End}}]);{{end}}
          {{ end }}
				{{end}}

        const rows = transformTaskToRows(tasks)
        if (rows.length > 0) {
          drawChart(rows)
        } else {
          container.innerHTML = '<div class="alert alert-success"><strong>Woops!</strong> There is no data!</div>';
        }

        function transformTaskToRows(_tasks){
          let taskByJobCount = [];
          for (let k in _tasks){
            taskByJobCount.push({name: k, size: _tasks[k].length});
          }
          taskByJobCount.sort(function(a, b) {
            if (a.size == b.size) return 0;
            return a.size > b.size ? -1 : 1;
          });

          let rows = [];
          for (let i = 0; i < taskByJobCount.length; i++) {
            jobs = _tasks[taskByJobCount[i].name];
            let jl = jobs.length;
            for (let j = 0; j < jl; j++) rows.push(jobs[j]);
          }
          return rows
        }

        function drawChart(_rows){
          let chart = new google.visualization.Timeline(container);
          let dataTable = new google.visualization.DataTable();
          dataTable.addColumn({ type: 'string', id: 'job' });
          dataTable.addColumn({ type: 'string', id: 'dummy bar label' });
          dataTable.addColumn({ type: 'string', role: 'tooltip' });
          dataTable.addColumn({ type: 'date', id: 'Start' });
          dataTable.addColumn({ type: 'date', id: 'End' });
          dataTable.addRows(_rows);
          chart.draw(dataTable, {
            timeline: { colorByRowLabel: true },
            avoidOverlappingGridLines: false
          });
          google.visualization.events.addListener(chart, 'onmouseover', function(e) {
            var t = document.getElementsByClassName("google-visualization-tooltip")[0];
            if (mousePosX) t.style.left = mousePosX + 'px';
            if (mousePosY) t.style.top = mousePosY - 120 + 'px';
          });
        }

        var mousePosX = undefined,
            mousePosY = undefined;

        document.addEventListener('mousemove', function(e) {
          mousePosX = e.pageX;
          mousePosY = e.pageY;
        });

        document.getElementById("cronv-cannot-render").innerHTML = Object.keys(tasksWithDayOfMonths).map( key => {
          return "<div><b>"+ key + "</b> : '" + tasksWithDayOfMonths[key]+ "'</div>"
        }).join("")


        
        const debounce = (func, wait) => {
          let timeout;
          return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };

          clearTimeout(timeout);
          timeout = setTimeout(later, wait);
        };
      };


      document.getElementById("cronv-input-filter").addEventListener("keyup", debounce(function(e){  
        const serachResult = e.target.value
        const foundKeyPartial = Object.keys(tasks).filter( key => {
           const result = key.includes(serachResult)
           return result
        }).map( val => tasks[val])


        if(foundKeyPartial.length > 0){
            const rows = transformTaskToRows(foundKeyPartial)
            if (rows.length > 0) {
            drawChart(rows)
            return
          }

        }

          const rows = transformTaskToRows(tasks)
          if (rows.length > 0) drawChart(rows)
      }, 300))

     });

  </script>
</body>
</html>
`

func makeTemplate() *template.Template {

	funcMap := template.FuncMap{
		"CronvIter": func(cronv *Record) <-chan *Exec {
			return cronv.iter()
		},
		"JSEscapeString": func(v string) string {
			return template.JSEscapeString(strings.TrimSpace(v))
		},
		"NewJsDate": func(v time.Time) string {
			return fmt.Sprintf("new Date(%d,%d,%d,%d,%d)", v.Year(), v.Month()-1, v.Day(), v.Hour(), v.Minute())
		},
		"DateFormat": func(v time.Time, format string) string {
			return v.Format(format)
		},
		"IsRunningEveryMinutes": func(c *Crontab) bool {
			return c.isRunningEveryMinutes()
		},
		"GetCronScheduleString": func(c *Crontab) string {
			return c.Schedule.toCrontab()
		},
		"ShouldRenderCronTabInDayOfMonthLevel": func(c *Crontab, dStr string) bool {
			if c.Schedule.DayOfMonth != "*" && strings.HasSuffix(dStr, "d") {
				fmt.Printf("🔴 warning: found crontab equal day of month level %s \n", c.Schedule.toCrontab())
				return false
			}
			return true
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(TEMPLATE))
}
