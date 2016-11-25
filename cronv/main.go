package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/takumakanari/cronv"
	"os"
)

const (
	VERSION = "0.2.2"
	NAME    = "Cronv"
)

func main() {
	opts := cronv.NewCronvCommand()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = fmt.Sprintf("%s v%s", NAME, VERSION)
	if _, err := parser.Parse(); err != nil {
		os.Exit(0)
	}

	ctx, err := cronv.NewCtx(opts)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && string(line[0]) != "#" {
			if err := ctx.AppendNewLine(line); err != nil {
				panic(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	path, err := ctx.Dump()
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] %d tasks.\n", opts.Title, len(ctx.CronEntries))
	fmt.Printf("[%s] '%s' generated.\n", opts.Title, path)
}
