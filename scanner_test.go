package dur_test

import (
	"dur"
	"reflect"
	"testing"
)

func TestScanner_Tokens(t *testing.T) {
	type testCase struct {
		name  string
		input string
		want  []dur.Token
	}

	tests := []testCase{
		{name: "empty input", input: "", want: []dur.Token{{Type: dur.TypeEOF}}},
		{name: "operators", input: "+-*/", want: []dur.Token{{Type: dur.TypePlus}, {Type: dur.TypeMinus}, {Type: dur.TypeMultiply}, {Type: dur.TypeDivide}, {Type: dur.TypeEOF}}},
		{name: "parentheses", input: "()", want: []dur.Token{{Type: dur.TypeParenOpen}, {Type: dur.TypeParenClose}, {Type: dur.TypeEOF}}},
		{name: "hours", input: "12h", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12h"}, {Type: dur.TypeEOF}}},
		{name: "minutes", input: "12m", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12m"}, {Type: dur.TypeEOF}}},
		{name: "seconds", input: "12s", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12s"}, {Type: dur.TypeEOF}}},
		{name: "milliseconds", input: "12ms", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12ms"}, {Type: dur.TypeEOF}}},
		{name: "milliseconds", input: "12us", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12us"}, {Type: dur.TypeEOF}}},
		{name: "nanoseconds", input: "12ns", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12ns"}, {Type: dur.TypeEOF}}},
		{name: "hours floating number 1", input: "12,5h", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12,5h"}, {Type: dur.TypeEOF}}},
		{name: "hours floating number 2", input: "12,333333h", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12,333333h"}, {Type: dur.TypeEOF}}},
		{name: "hours floating number 3", input: "12.333333h", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12.333333h"}, {Type: dur.TypeEOF}}},
		{name: "combined values", input: "12h11m2s", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12h"}, {Type: dur.TypeDuration, Literal: "11m"}, {Type: dur.TypeDuration, Literal: "2s"}, {Type: dur.TypeEOF}}},
		{name: "combined durations with operators", input: "12h-11m+10m*4/2", want: []dur.Token{{Type: dur.TypeDuration, Literal: "12h"}, {Type: dur.TypeMinus}, {Type: dur.TypeDuration, Literal: "11m"}, {Type: dur.TypePlus}, {Type: dur.TypeDuration, Literal: "10m"}, {Type: dur.TypeMultiply}, {Type: dur.TypeInteger, Literal: "4"}, {Type: dur.TypeDivide}, {Type: dur.TypeInteger, Literal: "2"}, {Type: dur.TypeEOF}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dur.NewScanner(tt.input)
			if got := s.Tokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
