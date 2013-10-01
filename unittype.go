package unit

import "sync"

// Uniter is an interface representing a type that can be converted
// to a unit.
type Uniter interface {
	Unit() *Unit
}

// Dimension is a type representing an SI base dimension or other
// orthogonal dimension. If a new dimension is desired for a
// domain-specific problem, NewDimension should be called.
type Dimension int

const (
	// SI Base Units
	Chemamt Dimension = iota // e.g. mol
	Current
	Length
	Luminosity
	Mass
	Temperature
	Time
	// Start of other SI Units
	Angle                // e.g. radians
	lastPackageDimension // Used in create dimension
)

// Dimensions represent the dimensionality of the unit in powers
// of that dimension. If a key is not present, the power of that
// dimension is zero. Dimensions is used in conjuction with NewUnit
type Dimensions map[Dimension]int

//TODO: Should there be some number reserved? We don't want users ever using integer literals
var lastCreatedDimension Dimension = 64      // Reserve first 63 for our use
var newUnitMutex *sync.Mutex = &sync.Mutex{} // so there is no race condition for dimension

// NewDimension returns a new dimension variable which will have a
// unique representation across packages to prevent accidental overlap.
// NewDimension should only be called for unit types that are orthogonal
// to the base dimensions defined in this package. For example, one unit
// that comes up in blood work is "White blood cells per microscope slide".
// NewDimension is appropriate for "White blood cells", as they are not
// representable in SI base units. However, NewDimension is not appropriate
// for "Slide", as slide is really a unit of area. Slide should instead be
// defined as a constant of type unit.Area
func NewDimension() Dimension {
	newUnitMutex.Lock()
	defer newUnitMutex.Unlock()
	lastCreatedDimension++
	return lastCreatedDimension
}

// Unit is a type a value with generic SI units. Most useful for
// translating between dimensions, for example, by multiplying
// an acceleration with a mass to get a force
type Unit struct {
	dimensions  map[Dimension]int // Map for custom dimensions
	value       float64
	current     int
	length      int
	luminosity  int
	mass        int
	temperature int
	time        int
	chemamt     int
}

// NewUnit creates a new variable of type Unit which has the value
// specified by value and the dimensions specified by the
// base units struct. The value is always in SI Units.
//
// Example: To create an acceleration of 3 m/s^2, one could do
// myvar := CreateUnit(3.0, &Dimensions{length: 1, time: -2})
func NewUnit(value float64, d Dimensions) *Unit {
	u := &Unit{
		current:     d[Current],
		length:      d[Length],
		luminosity:  d[Luminosity],
		mass:        d[Mass],
		temperature: d[Temperature],
		time:        d[Time],
		chemamt:     d[Chemamt],
		value:       value,
	}
	for key, val := range d {
		if key < lastPackageDimension {
			u.dimensions[key] = val
		}
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
