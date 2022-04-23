package dur

import (
	"fmt"
	"time"
)

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

func (i *Calculator) Calculate() time.Duration {
	return i.calculate(i.eof)
}

func (i *Calculator) calculate(isEnd func() bool) time.Duration {
	var (
		v1 time.Duration
		v2 time.Duration
		op TokenType
	)

	if i.outOfRange() {
		panic("out of range")
	}

	if isEnd() {
		i.pos++
		return time.Duration(0)
	}

	v1 = i.operand()

	if isEnd() {
		i.pos++
		return v1
	}

	for {
		if i.outOfRange() {
			panic("out of range")
		}

		if isEnd() {
			i.pos++
			break
		}

		if i.isOperator() && op == TypeEmpty {
			op = i.operator()
			continue
		}

		if op == TypeEmpty {
			op = TypePlus
		}

		v2 = i.operand()

		switch op {
		case TypePlus:
			vr := v1 + v2
			i.p.print(v1, v2, vr, string(plus))
			v1 = vr
			op = TypeEmpty
		case TypeMinus:
			vr := v1 - v2
			i.p.print(v1, v2, vr, string(minus))
			v1 = vr
			op = TypeEmpty
		default:
			panic(fmt.Sprintf("unknown operator '%v'", op))
		}
	}

	return v1
}

func (i *Calculator) operand() time.Duration {
	var mod = time.Duration(1)

	if i.tokenTypeEquals(TypeMinus) {
		mod = time.Duration(-1)
		i.pos++
	} else if i.tokenTypeEquals(TypePlus) {
		i.pos++
	}

	switch i.tokenType() {
	case TypeParenClose:
		panic("unexpected closing parenthesis")
	case TypeParenOpen:
		i.pos++
		return i.calculate(i.closingParen) * mod
	case TypeDuration:
		return i.duration() * mod
	default:
		panic(fmt.Sprintf("unexpected token '%v'", i.tokenType()))
	}
}

func (i *Calculator) duration() time.Duration {
	var tok = i.tokens[i.pos]

	dur := parseDuration(tok.Literal)

	i.pos++

	return dur
}

func (i *Calculator) operator() TokenType {
	var tok = i.tokens[i.pos]

	i.pos++

	return tok.Type
}

func (i *Calculator) outOfRange() bool {
	return len(i.tokens) <= i.pos
}

func (i *Calculator) eof() bool {
	return i.tokenTypeEquals(TypeEOF)
}

func (i *Calculator) closingParen() bool {
	return i.tokenTypeEquals(TypeParenClose)
}

func (i *Calculator) tokenTypeEquals(tokenType TokenType) bool {
	return i.tokenType() == tokenType
}

func (i *Calculator) tokenType() TokenType {
	return i.tokens[i.pos].Type
}

func (i *Calculator) isOperator() bool {
	return i.tokenTypeEquals(TypeMinus) || i.tokenTypeEquals(TypePlus)
}
