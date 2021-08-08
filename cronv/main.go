package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/takumakanari/cronv"
)

const (
	version = "0.4.4"
	name    = "Cronv"
)

func main() {
	opts := cronv.NewCronvCommand()

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = fmt.Sprintf("%s v%s", name, version)
	if _, err := parser.Parse(); err != nil {
		os.Exit(0)
	}

	ctx, err := cronv.NewVisualizer(opts)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := ctx.AppendNewLine(scanner.Text()); err != nil {
			panic(err)
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
