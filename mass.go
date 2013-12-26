// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
)

// Represents a mass in kilograms
type Mass float64

const (
	Yottagram Mass = 1e21
	Zettagram Mass = 1e18
	Exagram   Mass = 1e15
	Petagram  Mass = 1e12
	Teragram  Mass = 1e9
	Gigagram  Mass = 1e6
	Megagram  Mass = 1e3
	Kilogram  Mass = 1.0
	Gram      Mass = 1e-3
	Milligram Mass = 1e-6
	Microgram Mass = 1e-9
)

// Mass converts the Mass to a unit
func (m Mass) Unit() *Unit {
	return New(float64(m), Dimensions{MassDim: 1})
}

// Mass allows Mass to implement a Masser interface
func (m Mass) Mass() Mass {
	return m
}

// From converts the unit to a mass. Returns an error if there
// is a mismatch in dimension
func (m *Mass) From(u Uniter) error {
	if !DimensionsMatch(u, Kilogram) {
		(*m) = Mass(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*m) = Mass(u.Unit().Value())
	return nil
}

func (m Mass) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", m, float64(m))
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
		fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(m))
		fmt.Fprint(fs, " kg")
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g kg)", c, m, float64(m))
		return
	}
}
