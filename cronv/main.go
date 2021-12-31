package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/takumakanari/cronv"
)

const (
	version = "0.4.5"
	name    = "Cronv"
)

func main() {
	opts := cronv.NewCronvOption(time.Now())

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = fmt.Sprintf("%s v%s", name, version)
	if _, err := parser.Parse(); err != nil {
		os.Exit(0)
	}

	ctx := context.Background()
	viz, err := cronv.NewVisualizer(opts)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := viz.Add(ctx, scanner.Text()); err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	path, err := viz.Dump(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] %d tasks.\n", opts.Title, len(viz.CronEntries))
	fmt.Printf("[%s] '%s' generated.\n", opts.Title, path)
}
