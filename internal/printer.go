package internal

import (
	"fmt"
)

type printer interface {
	print(v1 interface{}, v2 interface{}, vr interface{}, op string)
}

type discardPrinter struct{}

func (p discardPrinter) print(_ interface{}, _ interface{}, _ interface{}, _ string) {
}

type nanoPrinter struct{}

func (p nanoPrinter) print(v1 interface{}, v2 interface{}, vr interface{}, op string) {
	fmt.Printf("%18d %s %18d = %18d\n", v1, op, v2, vr)
}

type humanReadablePrinter struct{}

func (p humanReadablePrinter) print(v1 interface{}, v2 interface{}, vr interface{}, op string) {
	fmt.Printf("%12v %s %12v = %12v\n", v1, op, v2, vr)
}
