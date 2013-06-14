package unit

// A uniter is a type that can be converted
// to a unit. These functions declare the
// power of the basic dimension of the unit
type Uniter interface {
	Unit() *Unit
}

// A list of the base units as defined by the SI system.
// Used to create a new unit
type Dimensions struct {
	Current     int
	Length      int
	Luminosity  int
	Mass        int
	Temperature int
	Time        int
	Chemamt     int
}

// A generic unit type. Mostly used for doing math involving units
type Unit struct {
	current     int
	length      int
	luminosity  int
	mass        int
	temperature int
	time        int
	chemamt     int // For mol
	value       float64
}

// TODO: Some oddities with chemamt. Not sure what to do

// Create a new variable of type Unit having the value
// specified by value and the dimensions specified by the
// base units struct.
//
// Example: To create an acceleration of 3 m/s^2, one could do
// myvar := CreateUnit(3.0, &Dimensions{length: 1, time: -2})
func CreateUnit(value float64, d *Dimensions) *Unit {
	return &Unit{
		current:     d.Current,
		length:      d.Length,
		luminosity:  d.Luminosity,
		mass:        d.Mass,
		temperature: d.Temperature,
		time:        d.Time,
		chemamt:     d.Chemamt,
		value:       value,
	}
}

// Check if the dimensions of two units are the same
func DimensionsMatch(aU, bU Uniter) bool {
	a := aU.Unit()
	b := bU.Unit()
	if a.length != b.length {
		return false
	}
	if a.time != b.time {
		return false
	}
	if a.mass != b.mass {
		return false
	}
	if a.current != b.current {
		return false
	}
	if a.temperature != b.temperature {
		return false
	}
	if a.luminosity != b.luminosity {
		return false
	}
	if a.chemamt != b.chemamt {
		return false
	}
	return true
}

func (u *Unit) Add(aU Uniter) *Unit {
	a := aU.Unit()
	if !DimensionsMatch(u, a) {
		panic("Attempted to add the values of two units whose dimensions do not match.")
	}
	u.value += a.value
	return u
}

func (u *Unit) Unit() *Unit {
	return u
}

// Multiply the receiver by the unit
func (u *Unit) Mul(aU Uniter) *Unit {
	a := aU.Unit()
	u.length += a.length
	u.time += a.time
	u.mass += a.mass
	u.current += a.current
	u.temperature += a.temperature
	u.luminosity += a.luminosity
	u.value *= a.value
	return u
}

// Divide the receive by the unit
func (u *Unit) Div(aU Uniter) *Unit {
	a := aU.Unit()
	u.length -= a.length
	u.time -= a.time
	u.mass -= a.mass
	u.current -= a.current
	u.temperature -= a.temperature
	u.luminosity -= a.luminosity
	u.value /= a.value
	return u
}

// Convert to a different dimension
func (u *Unit) In(a Uniter) float64 {
	u2 := a.Unit()
	if !DimensionsMatch(u, u2) {
		panic("Attempt to assign to the wrong dimension")
	}
	return u.value / u2.value
}

// Return the value of the unit (will always be in SI units).
// If it is wanted as a specific dimension, see ToLength, etc.
func (u *Unit) Value() float64 {
	return u.value
}
