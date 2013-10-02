package unit

import (
	"errors"
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
	return NewUnit(float64(m), Dimensions{MassDim: 1})
}

// Mass allows Mass to implement a masser interface
func (m Mass) Mass() Mass {
	return m
}

// FromUnit converts the unit to a mass. Returns an error if there
// is a mismatch in dimension
func (m *Mass) FromUnit(u Uniter) error {
	if !DimensionsMatch(u, Kilogram) {
		(*m) = Mass(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*m) = Mass(u.Unit().Value())
	return nil
}

// Length represents a length in meters
type Length float64

const (
	Meter      Length = 1.0
	Centimeter Length = 0.01
)

// Unit converts the Length to a unit
func (l Length) Unit() *Unit {
	return NewUnit(float64(l), Dimensions{Mass: 1})
}

// So it can implement a lengther interface
func (l Length) Length() Length {
	return l
}

// FromUnit converts a uniter to a length. Returns an error if there
// is a mismatch in dimension
func (l *Length) FromUnit(u Uniter) error {
	if !DimensionsMatch(u, Meter) {
		(*l) = Length(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*l) = Length(u.Unit().Value())
	return nil
}

// Dimless represents a dimensionless constant
type Dimless float64

const (
	One Dimless = 1.0
)

// Unit converts the Dimless to a unit
func (d Dimless) Unit() *Unit {
	return NewUnit(float64(d), Dimensions{})
}

// So it can implement a Dimless interface
func (d Dimless) Dimless() Dimless {
	return d
}

// FromUnit converts the unit to a dimless. Returns an error if there
// is a mismatch in dimension
func (d *Dimless) FromUnit(u *Unit) error {
	if !DimensionsMatch(u, One) {
		(*d) = Dimless(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*d) = Dimless(u.Unit().Value())
	return nil
}
