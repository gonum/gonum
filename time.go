package unit

import (
	"errors"
	"fmt"
	"math"
)

// Time represents a length in seconds
type Time float64

const (
	Hour        Time = 3600
	Minute      Time = 60
	Yottasecond Time = 1e24
	Zettasecond Time = 1e21
	Exasecond   Time = 1e18
	Petasecond  Time = 1e15
	Terasecond  Time = 1e12
	Gigasecond  Time = 1e9
	Megasecond  Time = 1e6
	Kilosecond  Time = 1e3
	Hectosecond Time = 1e2
	Decasecond  Time = 1e1
	Second      Time = 1.0
	Decisecond  Time = 1e-1
	Centisecond Time = 1e-2
	Millisecond Time = 1e-3
	Microsecond Time = 1e-6
	Nanosecond  Time = 1e-9
	Picosecond  Time = 1e-12
	Femtosecond Time = 1e-15
	Attosecond  Time = 1e-18
	Zeptosecond Time = 1e-21
	Yoctosecond Time = 1e-24
)

// Unit converts the Time to a unit
func (l Time) Unit() *Unit {
	return New(float64(l), Dimensions{TimeDim: 1})
}

// Time allows length to implement a Timeer interface
func (l Time) Time() Time {
	return l
}

// From converts a uniter to a length. Returns an error if there
// is a mismatch in dimension
func (l *Time) From(u Uniter) error {
	if !DimensionsMatch(u, Second) {
		(*l) = Time(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*l) = Time(u.Unit().Value())
	return nil
}

func (l Time) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", l, float64(l))
			return
		}
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p, pOk := fs.Precision()
		if !pOk {
			p = -1
		}
		w, wOk := fs.Width()
		if !wOk {
			w = -1
		}
		fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(l))
		fmt.Fprint(fs, " s")
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g s)", c, l, float64(l))
		return
	}
}
