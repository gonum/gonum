// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
)

// Time represents a time in seconds
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

// Unit converts the Time to a *Unit
func (time Time) Unit() *Unit {
	return New(float64(time), Dimensions{
		TimeDim: 1,
	})
}

// Time allows Time to implement a Timer interface
func (time Time) Time() Time {
	return time
}

// From converts a Uniter to a Time. Returns an error if
// there is a mismatch in dimension
func (time *Time) From(u Uniter) error {
	if !DimensionsMatch(u, Second) {
		(*time) = Time(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*time) = Time(u.Unit().Value())
	return nil
}

func (time Time) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", time, float64(time))
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
		fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(time))
		fmt.Fprint(fs, " s")
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g s)", c, time, float64(time))
		return
	}
}
