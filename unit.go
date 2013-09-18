package unit

// Uniter is an interface representing a type that can be converted
// to a unit.
type Uniter interface {
	Unit() *Unit
}

// Dimensions is a struct containing the SI base dimensions. Dimensions
// can be used in conjuntion with CreateUnit to create a
type Dimensions struct {
	Current     int
	Length      int
	Luminosity  int
	Mass        int
	Temperature int
	Time        int
	Chemamt     int
}

// Unit is a type a value with generic SI units. Most useful for
// translating between dimensions, for example, by multiplying
// an acceleration with a mass to get a force
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

// CreateUnit creates a new variable of type Unit which has the value
// specified by value and the dimensions specified by the
// base units struct. The value is always in SI Units.
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

// DimensionsMatch checks if the dimensions of two Uniters are the same
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

// Add adds the function argument to the reciever. Panics if the units of
// the receiver and the argument don't match.
func (u *Unit) Add(aU Uniter) *Unit {
	a := aU.Unit()
	if !DimensionsMatch(u, a) {
		panic("Attempted to add the values of two units whose dimensions do not match.")
	}
	u.value += a.value
	return u
}

// Unit allows unit to satisfy the uniter interface
func (u *Unit) Unit() *Unit {
	return u
}

// Mul multiply the receiver by the unit changing the dimensions
// of the receiver as appropriate
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

// Div divides the receiver by the argument changing the
// dimensions of the receiver as appropriate
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

// Value return the raw value of the unit as a float64. Use of this
// method is not recommended, instead it is recommended to use a
// FromUnit type of a specific dimension
func (u *Unit) Value() float64 {
	return u.value
}
