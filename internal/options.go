package internal

type options struct {
	p printer
}

type Option func(o *options)

func DiscardPrinter(o *options) {
	o.p = discardPrinter{}
}
func HumanReadablePrinter(o *options) {
	o.p = humanReadablePrinter{}
}
func NanoPrinter(o *options) {
	o.p = nanoPrinter{}
}
