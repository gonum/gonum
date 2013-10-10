package unit

import (
	"errors"
	"fmt"
	"math"
)

// Represents a mass in kilograms
type Mass float64

const (
	Kilogram  Mass = 1.0
	Gram      Mass = 1e-3
	Centigram Mass = 1e-5
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

// Dimless represents a dimensionless constant
type Dimless float64

const (
	One Dimless = 1.0
)

// Unit converts the Dimless to a unit
func (d Dimless) Unit() *Unit {
	return New(float64(d), Dimensions{})
}

// Dimless allows Dimless to implement a Dimlesser interface
func (d Dimless) Dimless() Dimless {
	return d
}

// From converts the unit to a dimless. Returns an error if there
// is a mismatch in dimension
func (d *Dimless) From(u *Unit) error {
	if !DimensionsMatch(u, One) {
		(*d) = Dimless(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*d) = Dimless(u.Unit().Value())
	return nil
}

func (d Dimless) Format(fs fmt.State, c rune) {
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T(%v)", d, float64(d))
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
		fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(d))
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%g)", c, d, float64(d))
		return
	}
}
