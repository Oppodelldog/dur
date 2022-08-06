package internal_test

import (
	"github.com/Oppodelldog/dur/internal"
	"reflect"
	"testing"
)

func TestScanner_Tokens(t *testing.T) {
	type testCase struct {
		name  string
		input string
		want  []internal.Token
	}

	tests := []testCase{
		{name: "empty input", input: "", want: []internal.Token{{Type: internal.TypeEOF}}},
		{name: "operators", input: "+-*/", want: []internal.Token{{Type: internal.TypePlus}, {Type: internal.TypeMinus}, {Type: internal.TypeMultiply}, {Type: internal.TypeDivide}, {Type: internal.TypeEOF}}},
		{name: "parentheses", input: "()", want: []internal.Token{{Type: internal.TypeParenOpen}, {Type: internal.TypeParenClose}, {Type: internal.TypeEOF}}},
		{name: "hours", input: "12h", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12h"}, {Type: internal.TypeEOF}}},
		{name: "minutes", input: "12m", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12m"}, {Type: internal.TypeEOF}}},
		{name: "seconds", input: "12s", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12s"}, {Type: internal.TypeEOF}}},
		{name: "milliseconds", input: "12ms", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12ms"}, {Type: internal.TypeEOF}}},
		{name: "milliseconds", input: "12us", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12us"}, {Type: internal.TypeEOF}}},
		{name: "nanoseconds", input: "12ns", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12ns"}, {Type: internal.TypeEOF}}},
		{name: "hours floating number 1", input: "12,5h", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12,5h"}, {Type: internal.TypeEOF}}},
		{name: "hours floating number 2", input: "12,333333h", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12,333333h"}, {Type: internal.TypeEOF}}},
		{name: "hours floating number 3", input: "12.333333h", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12.333333h"}, {Type: internal.TypeEOF}}},
		{name: "combined values", input: "12h11m2s", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12h"}, {Type: internal.TypeDuration, Literal: "11m"}, {Type: internal.TypeDuration, Literal: "2s"}, {Type: internal.TypeEOF}}},
		{name: "combined durations with operators", input: "12h-11m+10m*4/2", want: []internal.Token{{Type: internal.TypeDuration, Literal: "12h"}, {Type: internal.TypeMinus}, {Type: internal.TypeDuration, Literal: "11m"}, {Type: internal.TypePlus}, {Type: internal.TypeDuration, Literal: "10m"}, {Type: internal.TypeMultiply}, {Type: internal.TypeInteger, Literal: "4"}, {Type: internal.TypeDivide}, {Type: internal.TypeInteger, Literal: "2"}, {Type: internal.TypeEOF}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := internal.NewScanner(tt.input)
			if got := s.Tokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
