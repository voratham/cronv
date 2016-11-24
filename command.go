package cronv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	opt_date_format         = "2006/01/02"
	opt_time_format         = "15:04"
	opt_default_duration    = "6h"
	opt_default_output_path = "./crontab.html"
	opt_default_title       = "cron tasks"
)

type Command struct {
	OutputFilePath string `short:"o" long:"output" description:"path to .html file to output"`
	Duration       string `short:"d" long:"duration" description:"duration to visualize in N{suffix} style. e.g.) 1d(day)/1h(hour)/1m(minute)"`
	FromDate       string `long:"from-date" description:"start date in the format '2006/01/02' to visualize"`
	FromTime       string `long:"from-time" description:"start time in the format '15:04' to visualize"`
	Title          string `short:"t" long:"title" description:"title/label of output"`
}

func (self *Command) ToFromTime() (time.Time, error) {
	return time.Parse(fmt.Sprintf("%s %s", opt_date_format, opt_time_format),
		fmt.Sprintf("%s %s", self.FromDate, self.FromTime))
}

func (self *Command) ToDurationMinutes() (float64, error) {
	length := len(self.Duration)
	if length < 2 {
		return 0, errors.New(fmt.Sprintf("Invalid duration format: '%s'", self.Duration))
	}

	duration, err := strconv.Atoi(string(self.Duration[:length-1]))
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid duration format: '%s', %s", self.Duration, err))
	}

	unit := string(self.Duration[length-1])
	switch strings.ToLower(unit) {
	case "d":
		return float64(duration * 24 * 60), nil
	case "h":
		return float64(duration * 60), nil
	case "m":
		return float64(duration), nil
	}

	return 0, errors.New(fmt.Sprintf("Invalid duration format: '%s', '%s' is not in d/h/m", self.Duration, unit))
}

func NewCronvCommand() *Command {
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	return &Command{
		OutputFilePath: opt_default_output_path,
		Duration:       opt_default_duration,
		FromDate:       now.Format(opt_date_format),
		FromTime:       now.Format(opt_time_format),
		Title:          opt_default_title,
	}
}
