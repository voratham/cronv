# cronv

===================


## Description
Visualize your cron schedules in crontab


The cron specofication what cronv can handle is as follows:

[https://en.wikipedia.org/wiki/Cron#CRON_expression](https://en.wikipedia.org/wiki/Cron#CRON_expression)

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


## Options

- -d, -duration=1h : Duration to visualize, in %d{d,h,m} style, '6h' is used by default.
- -o, -output=./my_cron_schedule.html : Path to html file for output, './crontab.html' is used by default.


## TODO
