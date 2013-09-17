package unit

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

func (m Mass) In(m2 Mass) float64 {
	return float64(m) / float64(m2)
}

// So it can implement a masser interface
func (m Mass) Mass() Mass {
	return m
}
