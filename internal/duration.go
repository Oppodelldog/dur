package internal

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func parseDuration(lit string) time.Duration {
	if strings.Contains(lit, ",") || strings.Contains(lit, ".") {
		return durationFromFloat(lit)
	}

	dur, err := time.ParseDuration(lit)
	if err != nil {
		panic(fmt.Sprintf("invalid duration %v: %v", lit, err))
	}

	return dur
}

func durationFromFloat(lit string) time.Duration {
	var d time.Duration

	lit = strings.ReplaceAll(lit, ",", ".")
	number := lit[:len(lit)-1]
	unit := string(lit[len(lit)-1])

	if !isDigit(lit[len(lit)-2]) {
		number = lit[:len(lit)-2]
		unit = lit[len(lit)-2:]
	}

	v, err := strconv.ParseFloat(number, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid duration %v: %v", lit, err))
	}

	if unit == "ns" {
		panic("floating point values for unit ns is not supported")
	}

	switch unit {
	case "h":
		h, m := math.Modf(v)
		d += time.Hour * time.Duration(h)
		v = m * 60

		fallthrough
	case "m":
		m, s := math.Modf(v)
		d += time.Minute * time.Duration(m)
		v = s * 60

		fallthrough
	case "s":
		s, ms := math.Modf(v)
		d += time.Second * time.Duration(s)
		v = ms * 1000

		fallthrough
	case "ms":
		ms, mc := math.Modf(v)
		d += time.Millisecond * time.Duration(ms)
		v = mc * 1000

		fallthrough
	case "us":
		us, ns := math.Modf(v)
		d += time.Microsecond * time.Duration(us)
		v = ns * 1000

		fallthrough
	case "ns":
		d += time.Nanosecond * time.Duration(v)
	}

	return d
}
