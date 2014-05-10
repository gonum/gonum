// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"errors"
	"fmt"
	"math"
)

// Mass represents a mass in kilograms
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
	Hectogram Mass = 1e-1
	Decagram  Mass = 1e-2
	Gram      Mass = 1e-3
	Decigram  Mass = 1e-4
	Centigram Mass = 1e-5
	Milligram Mass = 1e-6
	Microgram Mass = 1e-9
	Nanogram  Mass = 1e-12
	Picogram  Mass = 1e-15
	Femtogram Mass = 1e-18
	Attogram  Mass = 1e-21
	Zeptogram Mass = 1e-24
	Yoctogram Mass = 1e-27
)

// Unit converts the Mass to a *Unit
func (mass Mass) Unit() *Unit {
	return New(float64(mass), Dimensions{
		MassDim: 1,
	})
}

// Mass allows Mass to implement a Masser interface
func (mass Mass) Mass() Mass {
	return mass
}

// From converts a Uniter to a Mass. Returns an error if
// there is a mismatch in dimension
func (mass *Mass) From(u Uniter) error {
	if !DimensionsMatch(u, Gram) {
		(*mass) = Mass(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*mass) = Mass(u.Unit().Value())
	return nil
}

func (mass Mass) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", mass, float64(mass))
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
		fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(mass))
		fmt.Fprint(fs, " kg")
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g kg)", c, mass, float64(mass))
		return
	}
}
