package cronv

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"os"
	"time"
)

type Cronv struct {
	Crontab         *Crontab
	expr            *cronexpr.Expression
	startTime       time.Time
	durationMinutes float64
}

func NewCronv(line string, startTime time.Time, durationMinutes float64) (*Cronv, error) {
	crontab, err := ParseCrontab(line)
	if err != nil {
		return nil, err
	}

	expr, err := cronexpr.Parse(crontab.Schedule.ToCrontab())
	if err != nil {
		return nil, err
	}

	cronv := &Cronv{
		Crontab:         crontab,
		expr:            expr,
		startTime:       startTime,
		durationMinutes: durationMinutes,
	}
	return cronv, nil
}

type Exec struct {
	Start time.Time
	End   time.Time
}

func (self *Cronv) Iter() <-chan *Exec {
	ch := make(chan *Exec)
	eneTime := self.startTime.Add(time.Duration(self.durationMinutes) * time.Minute)
	next := self.expr.Next(self.startTime)
	go func() {
		for next.Equal(eneTime) || eneTime.After(next) {
			ch <- &Exec{
				Start: next,
				End:   next.Add(time.Duration(1) * time.Minute),
			}
			next = self.expr.Next(next)
		}
		close(ch)
	}()
	return ch
}

type CronvCtx struct {
	Opts            *Command
	TimeFrom        time.Time
	TimeTo          time.Time
	CronEntries     []*Cronv
	durationMinutes float64
}

func (self *CronvCtx) AppendNewLine(line string) error {
	cronv, err := NewCronv(line, self.TimeFrom, self.durationMinutes)
	if err != nil {
		switch err.(type) {
		case *InvalidTaskError:
			return nil // pass
		default:
			return fmt.Errorf("Failed to analyze cron '%s': %s", line, err)
		}
	}
	self.CronEntries = append(self.CronEntries, cronv)
	return nil
}

func (self *CronvCtx) Dump() (string, error) {
	output, err := os.Create(self.Opts.OutputFilePath)
	if err != nil {
		return "", err
	}
	MakeTemplate().Execute(output, self)
	return self.Opts.OutputFilePath, nil
}

func NewCtx(opts *Command) (*CronvCtx, error) {
	timeFrom, err := opts.ToFromTime()
	if err != nil {
		return nil, err
	}

	durationMinutes, err := opts.ToDurationMinutes()
	if err != nil {
		return nil, err
	}

	return &CronvCtx{
		Opts:            opts,
		TimeFrom:        timeFrom,
		TimeTo:          timeFrom.Add(time.Duration(durationMinutes) * time.Minute),
		CronEntries:     []*Cronv{},
		durationMinutes: durationMinutes,
	}, nil
}
