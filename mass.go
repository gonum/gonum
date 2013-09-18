package unit

import (
	"errors"
	"math"
)

// Represents a mass in kg
type Mass float64

const (
	Kilogram  Mass = 1.0
	Gram      Mass = 1e-3
	Centigram Mass = 1e-5
	Milligram Mass = 1e-6
	Microgram Mass = 1e-9
)

func (m Mass) Unit() *Unit {
	return CreateUnit(float64(m), &Dimensions{Mass: 1})
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
