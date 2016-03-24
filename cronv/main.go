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
	"time"
)

const VERSION = "0.1.0"

const (
	OPT_DATE_FORMAT         = "2006/01/02"
	OPT_TIME_FORMAT         = "15:04"
	OPT_OUTPUT_PATH_DEFAULT = "./crontab.html"
)

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

func toFromTime(d string, t string) (time.Time, error) {
	return time.Parse(fmt.Sprintf("%s %s", OPT_DATE_FORMAT, OPT_TIME_FORMAT),
		fmt.Sprintf("%s %s", d, t))
}

func main() {
	var (
		outputFilePath string
		duration       string
		fromDate       string
		fromTime       string
	)
	now := optimizeTime(time.Now())

	for _, f := range []string{"o", "output"} {
		flag.StringVar(&outputFilePath, f, OPT_OUTPUT_PATH_DEFAULT, "path to .html file to output.")
	}
	for _, f := range []string{"d", "duration"} {
		flag.StringVar(&duration, f, "6h",
			"duration to visualize in N{suffix} style. e.g.) 1d(day)/1h(hour)/1m(minute).")
	}
	flag.StringVar(&fromDate, "from-date", now.Format(OPT_DATE_FORMAT),
		fmt.Sprintf("start date in the format '%s' to visualize.", OPT_DATE_FORMAT))
	flag.StringVar(&fromTime, "from-time", now.Format(OPT_TIME_FORMAT),
		fmt.Sprintf("start time in the format '%s' to visualize.", OPT_TIME_FORMAT))
	flag.Parse()

	timeFrom, err := toFromTime(fromDate, fromTime)
	if err != nil {
		panic(err)
	}

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

	cronv.MakeTemplate().Execute(output, map[string]interface{}{
		"CronEntries": cronEntries,
		"TimeFrom":    timeFrom,
		"Duration":    duration,
	})

	fmt.Printf("'%s' generated successfully.\n", outputFilePath)
}
