# cronv

===================


## Description
Visualize your cron schedules in crontab

![cronv output 1d](https://raw.github.com/wiki/takumakanari/cronv/images/outputs/cronv-1d.png)
![cronv output 30m](https://raw.github.com/wiki/takumakanari/cronv/images/outputs/cronv-30m.png)


## Installation

```shell
$ go install github.com/takumakanari/cronv/cronv
```

## Basic usage

Cronv can parse your crontab from stdin like as follows:

```shell
$ crontab -l | cronv -o ./my_cron_schedule.html
```

You can also specify the duration to analysis job schedules.

In a case like the follows, the job schedules will be analyzed from now to 24 hours later:


```shell
$ crontab -l | cronv -o ./my_cron_schedule.html -d 24h
```



> Cronv can parse cron entry written in basic cron format.
You can see the basically crontab specofication in [https://en.wikipedia.org/wiki/Cron#CRON_expression](https://en.wikipedia.org/wiki/Cron#CRON_expression).


## Options

- -d, -duration=1h : Duration to visualize, in %d{d,h,m} style, '6h' is used by default.
- -o, -output=./my_cron_schedule.html : Path to html file for output, './crontab.html' is used by default.
- -from-date=2016-11-08 : Start date in the format '2006/01/02' to visualize, current date is used by default.
- -from-time=01:30 : Start time in the format '15:04' to visualize, current time is used by default.
- -t, -title=mycrontab : Title/label in output html file.
- -h, -help : Show help message.

## Examples
Analyze crontab for 6 hours (by default) from now, , output html file to default path:
```shell
$ crontab -l | cronv
```

For 1 day from now, output html file to default path:

```shell
$ crontab -l | cronv -d 1d
```

For 12 hours from 21:00, today:

```shell
$ crontab -l | cronv -from-time 21:00 -d 12h
```

For 30 minuts from now, output html file to path/to/output.html:

```shell
$ crontab -l | cronv -d 30m -o path/to/output.html
```

For 2 hours from 2016/12/24 17:30, output html file to path/to/output2.html:

```shell
$ crontab -l | cronv -from-date '2016/12/24' -from-time 17:30 -d 2h -o path/to/output2.html
```

With original title/label:

```shell
$ crontab -l | cronv -d 1d -t "crontab@`hostname`"  # title/label of html file will be 'crontab@myhost'
```


## TODO

- Add output format/style other than HTML.
- Filter entries in output HTML file.


## Patch

Welcome!

