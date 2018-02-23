package cronv

import (
	"fmt"
	"strings"
)

type Schedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
}

func (self *Schedule) toCrontab() string {
	dest := strings.Join([]string{self.Minute, self.Hour, self.DayOfMonth,
		self.Month, self.DayOfWeek, self.Year}, " ")
	return strings.Trim(dest, " ")
}

type Extra struct {
	Line  string
	Label string
	Job   string
}

type Crontab struct {
	Line     string
	Schedule *Schedule
	Job      string
}

func (c *Crontab) isRunningEveryMinutes() bool {
	for i, v := range strings.Split(c.Schedule.toCrontab(), " ") {
		if v != "*" && (i > 0 || v != "*/1") {
			return false
		}
	}
	return true
}

type InvalidTaskError struct {
	Line string
}

func (e *InvalidTaskError) Error() string {
	return fmt.Sprintf("Invalid task: '%s'", e.Line)
}

func parseCrontab(line string) (*Crontab, *Extra, error) {
	// TODO use regrex to parse: https://gist.github.com/istvanp/310203
	parts := strings.Split(line, " ")

	// @reboot /something/to/do
	if parts[0] == "@reboot" {
		extra := &Extra{
			Line:  line,
			Label: parts[0],
			Job:   strings.Join(parts[1:], " "),
		}
		return nil, extra, nil
	}

	if len(parts) < 5 {
		return nil, nil, &InvalidTaskError{line}
	}

	// https://en.wikipedia.org/wiki/Cron#Predefined_scheduling_definitions
	schedule := &Schedule{}
	job := []string{}
	c := 0
	for _, v := range parts {
		if len(v) == 0 {
			continue
		}
		switch c {
		case 0:
			schedule.Minute = v
		case 1:
			schedule.Hour = v
		case 2:
			schedule.DayOfMonth = v
		case 3:
			schedule.Month = v
		case 4:
			schedule.DayOfWeek = v
		default:
			job = append(job, v)
		}
		c++
	}
	crontab := &Crontab{
		Line:     line,
		Schedule: schedule,
		Job:      strings.Join(job, " "),
	}
	return crontab, nil, nil
}
