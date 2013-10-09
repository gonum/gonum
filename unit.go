package unit

import (
	"errors"
	"math"
	"strconv"
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

func (m Mass) String() string {
	return strconv.FormatFloat(float64(m), 'e', -1, 64) + " kg"
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

func (l Length) String() string {
	return strconv.FormatFloat(float64(l), 'e', -1, 64) + " m"
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

func (d Dimless) String() string {
	return strconv.FormatFloat(float64(d), 'e', -1, 64)
}
