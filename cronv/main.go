package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/takumakanari/cronv"
	"os"
)

const (
	VERSION = "0.2.0"
	NAME    = "Cronv"
)

func main() {
	opts := cronv.NewCronvCommand()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = fmt.Sprintf("%s v%s", NAME, VERSION)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	timeFrom, err := opts.ToFromTime()
	if err != nil {
		panic(err)
	}

	durationMinutes, err := opts.ToDurationMinutes()
	if err != nil {
		panic(err)
	}

	output, err := os.Create(opts.OutputFilePath)
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
		"Duration":    opts.Duration,
	})

	fmt.Printf("'%s' generated successfully.\n", opts.OutputFilePath)
}
