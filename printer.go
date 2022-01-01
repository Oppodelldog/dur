package dur

import (
	"fmt"
	"time"
)

type discardPrinter struct{}

func (p discardPrinter) print(_ time.Duration, _ time.Duration, _ time.Duration, _ string) {
}

type nanoPrinter struct{}

func (p nanoPrinter) print(v1 time.Duration, v2 time.Duration, vr time.Duration, op string) {
	fmt.Printf("%18d %s %18d = %18d\n", v1, op, v2, vr)
}

type humanReadablePrinter struct{}

func (p humanReadablePrinter) print(v1 time.Duration, v2 time.Duration, vr time.Duration, op string) {
	fmt.Printf("%12v %s %12v = %12v\n", v1, op, v2, vr)
}
