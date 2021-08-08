package cronv

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tkmgo/cronexpr"
)

type Record struct {
	Crontab         *Crontab
	expr            *cronexpr.Expression
	startTime       time.Time
	durationMinutes float64
}

func newRecord(ctx context.Context, line string, startTime time.Time, durationMinutes float64) (*Record, *Extra, error) {
	crontab, extra, err := parse(ctx, line)
	if err != nil {
		return nil, nil, err
	}

	// Maybe the line was extra (@reboot, ENV etc ...)
	if crontab == nil {
		return nil, extra, nil
	}

	expr, err := cronexpr.Parse(crontab.Schedule.toCrontab())
	if err != nil {
		return nil, nil, err
	}

	return &Record{
		Crontab:         crontab,
		expr:            expr,
		startTime:       startTime,
		durationMinutes: durationMinutes,
	}, extra, nil
}

type Exec struct {
	Start time.Time
	End   time.Time
}

func (r *Record) iter() <-chan *Exec {
	ch := make(chan *Exec)
	eneTime := r.startTime.Add(time.Duration(r.durationMinutes) * time.Minute)
	next := r.expr.Next(r.startTime)
	go func() {
		for next.Equal(eneTime) || eneTime.After(next) {
			ch <- &Exec{
				Start: next,
				End:   next.Add(time.Duration(1) * time.Minute),
			}
			next = r.expr.Next(next)
		}
		close(ch)
	}()
	return ch
}

type Visualizer struct {
	Opts            *Command
	TimeFrom        time.Time
	TimeTo          time.Time
	CronEntries     []*Record
	Extras          []*Extra
	durationMinutes float64
}

func NewVisualizer(opts *Command) (*Visualizer, error) {
	timeFrom, err := opts.toFromTime()
	if err != nil {
		return nil, err
	}

	durationMinutes, err := opts.toDurationMinutes()
	if err != nil {
		return nil, err
	}

	return &Visualizer{
		Opts:            opts,
		TimeFrom:        timeFrom,
		TimeTo:          timeFrom.Add(time.Duration(durationMinutes) * time.Minute),
		durationMinutes: durationMinutes,
	}, nil
}

func (v *Visualizer) Add(ctx context.Context, line string) (bool, error) {
	trimed := strings.TrimSpace(line)
	if len(trimed) == 0 || string(trimed[0]) == "#" {
		return false, nil
	}

	record, extra, err := newRecord(ctx, trimed, v.TimeFrom, v.durationMinutes)
	if err != nil {
		switch err.(type) {
		case *InvalidTaskError:
			return false, nil // pass
		default:
			return false, fmt.Errorf("failed to analyze cron '%s': %s", line, err)
		}
	}

	if record != nil {
		v.CronEntries = append(v.CronEntries, record)
	}

	if extra != nil {
		v.Extras = append(v.Extras, extra)
	}

	return true, nil
}

func (v *Visualizer) Dump(ctx context.Context) (string, error) {
	output, err := os.Create(v.Opts.OutputFilePath)
	if err != nil {
		return "", err
	}
	makeTemplate().Execute(output, v)
	return v.Opts.OutputFilePath, nil
}
