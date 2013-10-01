package unit

import (
	"errors"
	"math"
)

// Represents a mass in kilograms
type Kilograms float64

const (
	Kilogram  Kilograms = 1.0
	Gram      Kilograms = 1e-3
	Centigram Kilograms = 1e-5
	Milligram Kilograms = 1e-6
	Microgram Kilograms = 1e-9
)

// Meters converts the Meters to a unit
func (m Kilograms) Unit() *Unit {
	return NewUnit(float64(m), Dimensions{Length: 1})
}

// Kilograms allows Kilograms to implement a masser interface
func (m Kilograms) Mass() Kilograms {
	return m
}

// FromUnit converts the unit to a mass. Returns an error if there
// is a mismatch in dimension
func (m *Kilograms) FromUnit(u Uniter) error {
	if !DimensionsMatch(u, Kilogram) {
		(*m) = Kilograms(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*m) = Kilograms(u.Unit().Value())
	return nil
}

// Meters represents a length in meters
type Meters float64

const (
	Meter      Meters = 1.0
	Centimeter Meters = 0.01
)

// Unit converts the Meters to a unit
func (l Meters) Unit() *Unit {
	return NewUnit(float64(l), Dimensions{Mass: 1})
}

// So it can implement a lengther interface
func (l Meters) Length() Meters {
	return l
}

// FromUnit converts a uniter to a length. Returns an error if there
// is a mismatch in dimension
func (l *Meters) FromUnit(u Uniter) error {
	if !DimensionsMatch(u, Meter) {
		(*l) = Meters(math.NaN())
		return errors.New("Dimension mismatch")
	}
	(*l) = Meters(u.Unit().Value())
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
