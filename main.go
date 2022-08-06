package main

import (
	"flag"
	"fmt"
	"github.com/Oppodelldog/dur/internal"
	"os"
	"strings"
)

func main() {
	var (
		options []internal.Option
		fs      = flag.NewFlagSet("dur", flag.ExitOnError)
		printer = fs.String("p", "", "prints a line for each calculation that is performed.\nOutput options:\n  h - human readable\n  n - nanoseconds\nexample: -p=h")
	)

	fs.Usage = usage(fs)

	var err = fs.Parse(os.Args[1:])

	input := strings.Join(os.Args[len(os.Args)-fs.NArg():], "")

	if len(input) == 0 || err != nil {
		fs.Usage()
	}

	switch *printer {
	case "h":
		options = append(options, internal.HumanReadablePrinter)
	case "n":
		options = append(options, internal.NanoPrinter)
	}

	fmt.Println(internal.NewCalculator(input, options...).Calculate())
}

func usage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Println("dur - duration calculation")
		fmt.Println("usage: dur [-p] OPERAND [OPERATION OPERAND ...]")
		fmt.Println()
		fmt.Println("example: dur 40h + 2h")
		fmt.Println(" output: 42h0m0s")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}
}
