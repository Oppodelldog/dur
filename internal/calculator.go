package internal

import (
	"fmt"
	"strconv"
	"time"
)

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
		v1 interface{}
		v2 interface{}
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
		return mustDuration(v1)
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
			vr := add(v1, v2)
			i.p.print(v1, v2, vr, string(plus))
			v1 = vr
			op = TypeEmpty
		case TypeMinus:
			vr := sub(v1, v2)
			i.p.print(v1, v2, vr, string(minus))
			v1 = vr
			op = TypeEmpty
		case TypeMultiply:
			vr := mul(v1, v2)
			i.p.print(v1, v2, vr, string(multiply))
			v1 = vr
			op = TypeEmpty
		case TypeDivide:
			vr := div(v1, v2)
			i.p.print(v1, v2, vr, string(divide))
			v1 = vr
			op = TypeEmpty
		default:
			panic(fmt.Sprintf("unknown operator '%v'", op))
		}
	}

	return mustDuration(v1)
}

func (i *Calculator) operand() interface{} {
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
	case TypeInteger:
		return i.integer() * int(mod)
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

func (i *Calculator) integer() int {
	var tok = i.tokens[i.pos]

	dur, err := strconv.Atoi(tok.Literal)
	if err != nil {
		panic(err)
	}

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
	return i.tokenTypeEquals(TypeMinus) || i.tokenTypeEquals(TypePlus) || i.tokenTypeEquals(TypeMultiply) || i.tokenTypeEquals(TypeDivide)
}

func add(v1, v2 interface{}) time.Duration {
	mustDuration(v1)
	mustDuration(v2)

	return v1.(time.Duration) + v2.(time.Duration)
}

func sub(v1, v2 interface{}) time.Duration {
	mustDuration(v1)
	mustDuration(v2)

	return v1.(time.Duration) - v2.(time.Duration)
}

func div(v1 interface{}, v2 interface{}) interface{} {
	var vr interface{}

	if i1, ok1 := v1.(int); ok1 {
		if i2, ok2 := v2.(int); ok2 {
			vr = i1 / i2
		} else if _, ok2 := v2.(time.Duration); ok2 {
			panic("cannot divide by a duration")
		}
	} else if i1, ok1 := v1.(time.Duration); ok1 {
		if i2, ok2 := v2.(int); ok2 {
			vr = i1 / time.Duration(i2)
		} else if _, ok2 := v2.(time.Duration); ok2 {
			panic("cannot calculate 2 durations")
		}
	}

	return vr
}

func mul(v1 interface{}, v2 interface{}) interface{} {
	var vr interface{}

	if i1, ok1 := v1.(int); ok1 {
		if i2, ok2 := v2.(int); ok2 {
			vr = i1 * i2
		} else if i2, ok2 := v2.(time.Duration); ok2 {
			vr = time.Duration(i1) * i2
		}
	} else if i1, ok1 := v1.(time.Duration); ok1 {
		if i2, ok2 := v2.(int); ok2 {
			vr = i1 * time.Duration(i2)
		} else if _, ok2 := v2.(time.Duration); ok2 {
			panic("cannot calculate 2 durations")
		}
	}

	return vr
}

func mustDuration(v interface{}) time.Duration {
	if duration, ok := v.(time.Duration); ok {
		return duration
	}

	panic("result is no duration")
}
