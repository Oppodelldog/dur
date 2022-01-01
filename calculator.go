package dur

import (
	"fmt"
	"time"
)

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

type printer interface {
	print(v1 time.Duration, v2 time.Duration, vr time.Duration, op string)
}

func NewCalculator(input string, opts ...Option) *Calculator {
	var options options

	DiscardPrinter(&options)

	for _, opt := range opts {
		opt(&options)
	}

	return &Calculator{
		tokens: NewScanner(input).Tokens(),
		p:      options.p,
	}
}

type Calculator struct {
	tokens []Token
	pos    int
	p      printer
}

func (i *Calculator) Calculate() string {
	var (
		v1 time.Duration
		v2 time.Duration
		op TokenType
	)

	if i.eof() {
		return ""
	}

	v1 = i.duration()

	if i.eof() {
		return v1.String()
	}

	for {
		if i.eof() {
			break
		}

		if i.tokenType() == TypeValue {
			op = TypePlus
			v2 = i.duration()
		} else {
			op = i.operation()
			v2 = i.duration()
		}

		switch op {
		case TypePlus:
			vr := v1 + v2
			i.p.print(v1, v2, vr, "+")
			v1 = vr
		case TypeMinus:
			vr := v1 - v2
			i.p.print(v1, v2, vr, "-")
			v1 = vr
		}
	}

	return v1.String()
}

func (i *Calculator) duration() time.Duration {
	var tok = i.tokens[i.pos]
	if tok.Type != TypeValue {
		panic(fmt.Sprintf("expected value, but got %v (%v'')", tok.Type, tok.Literal))
	}

	dur := parseDuration(tok.Literal)

	i.pos++

	return dur
}

func (i *Calculator) operation() TokenType {
	var tok = i.tokens[i.pos]
	if tok.Type != TypePlus && tok.Type != TypeMinus {
		panic(fmt.Sprintf("expected operation, but got %v ('%v')", tok.Type, tok.Literal))
	}
	i.pos++

	return tok.Type
}

func (i *Calculator) eof() bool {
	return i.tokens[i.pos].Type == TypeEOF
}

func (i *Calculator) tokenType() TokenType {
	return i.tokens[i.pos].Type
}
