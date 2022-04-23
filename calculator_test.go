package dur_test

import (
	"dur"
	"reflect"
	"testing"
)

func TestCalculator_Calculate(t *testing.T) {
	type testCase struct {
		name  string
		input string
		want  string
	}

	tests := []testCase{
		{name: "single value hours", input: "0h", want: "0s"},
		{name: "single value hours", input: "1h", want: "1h0m0s"},
		{name: "single value hours", input: "61h", want: "61h0m0s"},

		{name: "single value minutes", input: "0m", want: "0s"},
		{name: "single value minutes", input: "1m", want: "1m0s"},
		{name: "single value minutes", input: "61m", want: "1h1m0s"},

		{name: "single value seconds", input: "0s", want: "0s"},
		{name: "single value seconds", input: "1s", want: "1s"},
		{name: "single value seconds", input: "61s", want: "1m1s"},

		{name: "single value float hours", input: "0,166666666666667h", want: "10m0s"},
		{name: "single value float hours", input: "0,16666666666666h", want: "9m59.999999999s"},
		{name: "single value float hours", input: "1h", want: "1h0m0s"},
		{name: "single value float hours", input: "1,25h", want: "1h15m0s"},
		{name: "single value float hours", input: "1,5h", want: "1h30m0s"},
		{name: "single value float hours", input: "1,75h", want: "1h45m0s"},
		{name: "single value float hours", input: "2.0h", want: "2h0m0s"},

		{name: "single value float minutes", input: "10,5m", want: "10m30s"},
		{name: "single value float seconds", input: "10,005s", want: "10.005s"},

		{name: "single value float milliseconds", input: "10,5ms", want: "10.5ms"},
		{name: "single value float milliseconds", input: "10,5ms+0,5ms", want: "11ms"},

		{name: "single value float microseconds", input: "10,5us", want: "10.5µs"},
		{name: "single value float microseconds", input: "10,5us+0,5us", want: "11µs"},

		{name: "all units", input: "1h1m1s1ms1us1ns", want: "1h1m1.001001001s"},

		{name: "add", input: "1h+12m", want: "1h12m0s"},
		{name: "subtract", input: "1h-12m", want: "48m0s"},
		{name: "subtract add", input: "1h-12m+3h", want: "3h48m0s"},

		{name: "subtract float hours", input: "2,0h-0,5h", want: "1h30m0s"},
		{name: "subtract float hour from value", input: "30m-0,5h", want: "0s"},
		{name: "subtract float value from float hour", input: "0,5h-30m", want: "0s"},

		{name: "add float hours", input: "2,0h+0,5h", want: "2h30m0s"},
		{name: "add float hour to value", input: "30m+0,5h", want: "1h0m0s"},
		{name: "add float value to float hour", input: "0,5h+30m", want: "1h0m0s"},

		{name: "add milliseconds", input: "1s+2ms", want: "1.002s"},
		{name: "add microseconds", input: "1s+2us", want: "1.000002s"},
		{name: "add nanoseconds", input: "1s+2ns", want: "1.000000002s"},
		{name: "add nanoseconds", input: "1s+2000ns", want: "1.000002s"},
		{name: "add nanoseconds", input: "1s+1000000000ns", want: "2s"},
		{name: "add all", input: "1h1m1s999ms999us999ns + 1ns", want: "1h1m2s"},

		{name: "add float seconds", input: "9m59.999999999s+0.000000001s", want: "10m0s"},

		{name: "be aware of floating precision", input: "3.01s + 0s", want: "3.009999999s"},

		{name: "multiple value concat", input: "1h30m", want: "1h30m0s"},
		{name: "multiple value concat in subtraction", input: "10h - 1h30m", want: "9h30m0s"},

		{name: "parentheses", input: "()", want: "0s"},
		{name: "parentheses", input: "(1h)", want: "1h0m0s"},
		{name: "parentheses", input: "(-1h)", want: "-1h0m0s"},
		{name: "parentheses", input: "-(1h)", want: "-1h0m0s"},
		{name: "parentheses", input: "-(+1h)", want: "-1h0m0s"},
		{name: "parentheses", input: "(-1h-1h)", want: "-2h0m0s"},
		{name: "parentheses", input: "(-1h)+(-1h)", want: "-2h0m0s"},
		{name: "parentheses", input: "(-1h)-(-1h)", want: "0s"},
		{name: "parentheses", input: "(-1h+-1h)", want: "-2h0m0s"},
		{name: "parentheses", input: "(-1h-+1h)", want: "-2h0m0s"},
		{name: "parentheses", input: "1h(())1h()1h(())", want: "3h0m0s"},
		{name: "parentheses", input: "1h(1h(1h))", want: "3h0m0s"},
		{name: "parentheses", input: "1h(10m)", want: "1h10m0s"},
		{name: "parentheses", input: "1h(10m)", want: "1h10m0s"},
		{name: "parentheses", input: "1h(10m+20m)", want: "1h30m0s"},
		{name: "parentheses", input: "2h-(1h30m)", want: "30m0s"},
		{name: "parentheses", input: "2h-(1h+30m)", want: "30m0s"},
		{name: "parentheses", input: "2h-(1h+30m)", want: "30m0s"},
		{name: "parentheses", input: "(2h-(1h+30m))", want: "30m0s"},
		{name: "parentheses", input: "(2h)-(1h+30m)", want: "30m0s"},
		{name: "parentheses", input: "((2h)-(1h+30m))", want: "30m0s"},

		{name: "empty", input: "", want: "0s"},
		{name: "missing operand", input: "0h+", want: "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dur.NewCalculator(tt.input).Calculate().String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculator_Calculate_Panics(t *testing.T) {
	type testCase struct {
		name  string
		input string
		want  string
	}

	tests := []testCase{
		{name: "floating nanoseconds", input: "0,1ns", want: "floating point values for unit ns is not supported"},
		{name: "unexpected character", input: "0hh", want: "unexpected character 'h'"},
		{name: "floating nanoseconds", input: "0,,1s", want: "expected digit, but got ','"},
		{name: "missing end of term", input: ")", want: "unexpected closing parenthesis"},
		{name: "unexpected end of term", input: ")1h", want: "unexpected closing parenthesis"},
		{name: "missing begin of term", input: "1h)", want: "unexpected closing parenthesis"},
		{name: "missing end of term", input: "(", want: "unexpected token 'EOF'"},
		{name: "missing end of term", input: "(1h", want: "unexpected token 'EOF'"},
		{name: "missing end of term", input: "1h(", want: "unexpected token 'EOF'"},
		{name: "invalid operator", input: "--1h", want: "unexpected token 'MINUS'"},
		{name: "invalid operator", input: "-+1h", want: "unexpected token 'PLUS'"},
		{name: "parentheses", input: "1h(1h", want: "unexpected token 'EOF'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					r := recover()
					if r == nil {
						t.Errorf("The code did not panic")
					} else if r != tt.want {
						t.Errorf("panic = %v, want %v", r, tt.want)
					}

				}()
				if got := dur.NewCalculator(tt.input).Calculate(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Calculate() = %v, want %v", got, tt.want)
				}
			}()
		})
	}
}
