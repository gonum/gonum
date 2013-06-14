package unit

// Length units

// Represents a length in meters
type Length float64

const (
	Meter      Length = 1.0
	Centimeter Length = 0.01
	Foot       Length = 0.3048
	Furlong    Length = 600 * Foot
	Yard       Length = 3 * Foot
	Mile       Length = 5280 * Foot
	Inch       Length = 1.0 / 12 * Foot
	Smoot      Length = 67.0 / 12 * Foot
	Au         Length = 149597870700
	LightYear  Length = 9460730472580800
)

func (l Length) Unit() *Unit {
	return CreateUnit(float64(l), &Dimensions{Length: 1})
}

func (l Length) In(l2 Length) float64 {
	return float64(l) / float64(l2)
}

// So it can implement a lengther interface
func (l Length) Length() Length {
	return l
}
