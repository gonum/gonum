package mass

import "github.com/gonum/unit"

// Represents a mass in kg
type Mass float64

const (
	Kilogram Mass = 1.0
	Gram     Mass = 0.001
	Pound    Mass = 0.45359237
)

func (m Mass) Unit() *unit.Unit {
	return unit.CreateUnit(float64(m), &unit.Dimensions{Mass: 1})
}

func (m Mass) In(m2 Mass) float64 {
	return float64(m) / float64(m2)
}

// So it can implement a masser interface
func (m Mass) Mass() Mass {
	return m
}
