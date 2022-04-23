package dur

import (
	"fmt"
	"strings"
)

const (
	TypeEOF        TokenType = "EOF"
	TypeWhitespace TokenType = "WHITESPACE"
	TypePlus       TokenType = "PLUS"
	TypeMinus      TokenType = "MINUS"
	TypeMultiply   TokenType = "MULTIPLY"
	TypeDivide     TokenType = "DIVIDE"
	TypeParenOpen  TokenType = "PAREN_OPEN"
	TypeParenClose TokenType = "PAREN_CLOSE"

	TypeValue = "VALUE"

	TypeEmpty = ""
)

const (
	minus      = '-'
	plus       = '+'
	multiply   = '*'
	divide     = '/'
	space      = ' '
	parenOpen  = '('
	parenClose = ')'
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewScanner(input string) *Scanner {
	return &Scanner{
		input: input,
		pos:   0,
		len:   len(input),
	}
}

type Scanner struct {
	input string
	pos   int
	len   int
}

func (s *Scanner) Tokens() []Token {
	var tokens []Token

	for {
		token := s.nextToken()
		if token.Type == TypeWhitespace {
			continue
		}

		tokens = append(tokens, token)

		if token.Type == TypeEOF {
			break
		}
	}

	return tokens
}

func (s *Scanner) eof(offset int) bool {
	return s.pos+offset >= s.len
}

func (s *Scanner) peek(offset int) byte {
	if s.eof(offset) {
		panic(fmt.Sprintf("peek out of bounds at pos: %v, offset: %v", s.pos, offset))
	}

	return s.input[s.pos+offset]
}

func (s *Scanner) nextChar() {
	s.pos++
}

func (s *Scanner) read() byte {
	if s.eof(0) {
		panic(fmt.Sprintf("read out of bounds at pos: %v", s.pos))
	}

	ch := s.input[s.pos]
	s.nextChar()

	return ch
}

func (s *Scanner) nextToken() Token {
	var tok Token

	if s.eof(0) {
		return Token{Type: TypeEOF}
	}

	ch := s.peek(0)

	switch {
	case ch == space:
		tok = Token{Type: TypeWhitespace}

		s.nextChar()
	case ch == minus:
		tok = Token{Type: TypeMinus}

		s.nextChar()
	case ch == plus:
		tok = Token{Type: TypePlus}

		s.nextChar()
	case ch == multiply:
		tok = Token{Type: TypeMultiply}

		s.nextChar()
	case ch == divide:
		tok = Token{Type: TypeDivide}

		s.nextChar()
	case ch == parenOpen:
		tok = Token{Type: TypeParenOpen}

		s.nextChar()
	case ch == parenClose:
		tok = Token{Type: TypeParenClose}

		s.nextChar()
	case isDigit(ch):
		tok = s.readValue()
	default:
		panic(fmt.Sprintf("unexpected character '%v'", string(ch)))
	}

	return tok
}

func (s *Scanner) readValue() Token {
	const (
		uh   = 'h'
		um   = 'm'
		us   = 's'
		umc  = 'u'
		un   = 'n'
		dec1 = ','
		dec2 = '.'
	)

	var (
		numDec = 0
		sb     = strings.Builder{}
		ch     = s.read()
	)

	if !isDigit(ch) {
		panic(fmt.Sprintf("first character of value must be digit, got '%s'", string(ch)))
	}

	sb.WriteByte(ch)

loop:
	for {
		ch = s.read()
		switch {
		case ch == umc:
			sb.WriteByte(ch)
			ch = s.read()
			if ch != us {
				panic(fmt.Sprintf("invalid character for microseconds '%s'", string(ch)))
			}
			sb.WriteByte(ch)

			break loop
		case ch == un:
			sb.WriteByte(ch)
			ch = s.read()
			if ch != us {
				panic(fmt.Sprintf("invalid character for nanoseconds '%s'", string(ch)))
			}
			sb.WriteByte(ch)

			break loop
		case ch == uh || ch == um || ch == us:
			sb.WriteByte(ch)
			if !s.eof(0) && s.peek(0) == us && ch == um {
				sb.WriteByte(s.read())
			}

			break loop
		case (ch == dec1 || ch == dec2) && numDec == 0:
			sb.WriteByte(ch)
			numDec++
		case isDigit(ch):
			sb.WriteByte(ch)
		default:
			panic(fmt.Sprintf("expected digit, but got '%s'", string(ch)))
		}
	}

	return Token{Type: TypeValue, Literal: sb.String()}
}

func isDigit(ch byte) bool {
	return ch >= 48 && ch <= 57
}
