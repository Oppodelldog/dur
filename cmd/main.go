package main

import (
	"dur"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		options []dur.Option
		fs      = flag.NewFlagSet("dur", flag.ExitOnError)
		printer = fs.String("p", "", "-p=h -p=n (h - human readable, n - nanoseconds)")
		err     = fs.Parse(os.Args[1:])
	)

	if err != nil {
		panic(err)
	}

	input := strings.Join(os.Args[len(os.Args)-fs.NArg():], "")

	switch *printer {
	case "h":
		options = append(options, dur.HumanReadablePrinter)
	case "n":
		options = append(options, dur.NanoPrinter)
	}

	fmt.Println(dur.NewCalculator(input, options...).Calculate())
}
