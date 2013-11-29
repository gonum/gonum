// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
)

// Length represents a length in meters
type Length float64

const (
	Meter      Length = 1.0
	Centimeter Length = 0.01
)

// Unit converts the Length to a unit
func (l Length) Unit() *Unit {
	return New(float64(l), Dimensions{MassDim: 1})
}

// Length allows length to implement a Lengther interface
func (l Length) Length() Length {
	return l
}

// From converts a uniter to a length. Returns an error if there
// is a mismatch in dimension
func (l *Length) From(u Uniter) error {
	if !DimensionsMatch(u, Meter) {
		(*l) = Length(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*l) = Length(u.Unit().Value())
	return nil
}

func (l Length) Format(fs fmt.State, c rune) {
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
		fmt.Fprint(fs, " m")
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g m)", c, l, float64(l))
		return
	}
}
