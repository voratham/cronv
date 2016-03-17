package cronv

import "strings"

type Schedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
}

func (self *Schedule) ToCrontab() string {
	dest := strings.Join([]string{self.Minute, self.Hour, self.DayOfMonth, self.Month, self.DayOfWeek, self.Year}, " ")
	return strings.Trim(dest, " ")
}

type Crontab struct {
	Line     string
	Schedule *Schedule
	Job      string
}

func ParseCrontab(line string) (*Crontab, error) {
	// https://en.wikipedia.org/wiki/Cron#Predefined_scheduling_definitions
	// TODO use regrex to parse (// https://gist.github.com/istvanp/310203)
	schedule := &Schedule{}
	job := []string{}
	c := 0
	for _, v := range strings.Split(line, " ") {
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
	return crontab, nil
}
