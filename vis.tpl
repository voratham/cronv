<html>
<head>
<title>{{.Opts.Title}} | {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
</head>
<body>
  <div class="container-fluid">
    <h1>{{.Opts.Title}}</h1>
    <p>From {{DateFormat .TimeFrom "2006/1/2 15:04"}}, +{{.Opts.Duration}}</p>
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
