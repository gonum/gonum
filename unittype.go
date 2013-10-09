package unit

import (
	"sort"
	"strconv"
	"sync"
)

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
	CurrentDim Dimension = iota
	LengthDim
	LuminousIntensityDim
	MassDim
	TemperatureDim
	TimeDim
	// Start of other SI Units
	AngleDim             // e.g. radians
	lastPackageDimension // Used in create dimension
)

//TODO: Should there be some number reserved? We don't want users ever using integer literals
var lastCreatedDimension Dimension = lastPackageDimension

var dimensionToSymbol map[Dimension]string = make(map[Dimension]string)
var symbolToDimension map[string]Dimension = make(map[string]Dimension)

func init() {
	dimensionToSymbol[CurrentDim] = "A"
	symbolToDimension["A"] = CurrentDim
	dimensionToSymbol[LengthDim] = "m"
	symbolToDimension["m"] = LengthDim
	dimensionToSymbol[LuminousIntensityDim] = "cd"
	symbolToDimension["cd"] = LuminousIntensityDim
	dimensionToSymbol[MassDim] = "kg"
	symbolToDimension["kg"] = MassDim
	dimensionToSymbol[TemperatureDim] = "K"
	symbolToDimension["K"] = TemperatureDim
	dimensionToSymbol[TimeDim] = "s"
	symbolToDimension["s"] = TimeDim
	dimensionToSymbol[AngleDim] = "rad"
	symbolToDimension["rad"] = AngleDim
}

// Dimensions represent the dimensionality of the unit in powers
// of that dimension. If a key is not present, the power of that
// dimension is zero. Dimensions is used in conjuction with NewUnit
type Dimensions map[Dimension]int

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
//  Expecting to be used only during initialization, it panics if the mapping between types and names is not a bijection.
func NewDimension(symbol string) Dimension {
	newUnitMutex.Lock()
	defer newUnitMutex.Unlock()
	lastCreatedDimension++
	_, ok := symbolToDimension[symbol]
	if ok {
		panic("unit: dimension string " + symbol + " already used")
	}
	dimensionToSymbol[lastCreatedDimension] = symbol
	symbolToDimension[symbol] = lastCreatedDimension
	return lastCreatedDimension
}

// Unit is a type a value with generic SI units. Most useful for
// translating between dimensions, for example, by multiplying
// an acceleration with a mass to get a force
type Unit struct {
	dimensions map[Dimension]int // Map for custom dimensions
	value      float64
}

// NewUnit creates a new variable of type Unit which has the value
// specified by value and the dimensions specified by the
// base units struct. The value is always in SI Units.
//
// Example: To create an acceleration of 3 m/s^2, one could do
// myvar := CreateUnit(3.0, &Dimensions{length: 1, time: -2})
func NewUnit(value float64, d Dimensions) *Unit {

	// TODO: Find most efficient way of doing this
	// I think copy is necessary in case the input
	// dimension map is changed later
	u := &Unit{
		dimensions: make(map[Dimension]int),
	}
	for key, val := range d {
		u.dimensions[key] = val
	}
	u.value = value
	return u
}

// DimensionsMatch checks if the dimensions of two Uniters are the same
func DimensionsMatch(a, b Uniter) bool {
	aUnit := a.Unit()
	bUnit := b.Unit()
	if len(aUnit.dimensions) != len(bUnit.dimensions) {
		return false
	}
	for key, val := range aUnit.dimensions {
		if bUnit.dimensions[key] != val {
			return false
		}
	}
	return true
}

// Add adds the function argument to the reciever. Panics if the units of
// the receiver and the argument don't match.
func (u *Unit) Add(uniter Uniter) *Unit {
	a := uniter.Unit()
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
func (u *Unit) Mul(uniter Uniter) *Unit {
	a := uniter.Unit()
	for key, val := range a.dimensions {
		u.dimensions[key] += val
	}
	u.value *= a.value
	return u
}

// Div divides the receiver by the argument changing the
// dimensions of the receiver as appropriate
func (u *Unit) Div(uniter Uniter) *Unit {
	a := uniter.Unit()
	u.value /= a.value
	for key, val := range a.dimensions {
		u.dimensions[key] -= val
	}
	return u
}

// Value return the raw value of the unit as a float64. Use of this
// method is not recommended, instead it is recommended to use a
// FromUnit type of a specific dimension
func (u *Unit) Value() float64 {
	return u.value
}

type symbolString struct {
	symbol string
	pow    int
}

type unitPrinters []symbolString

func (u unitPrinters) Len() int {
	return len(u)
}

func (u unitPrinters) Less(i, j int) bool {
	return u[i].symbol < u[j].symbol
}

func (u unitPrinters) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u Unit) String() string {
	str := strconv.FormatFloat(u.value, 'e', -1, 64)
	// Map iterates randomly, but we should output the symbols in a logical order
	data := make(unitPrinters, 0, 10)
	for dimension, power := range u.dimensions {
		if power != 0 {
			data = append(data, symbolString{dimensionToSymbol[dimension], power})
		}
	}
	sort.Sort(data)
	for _, s := range data {
		str += " " + s.symbol
		if s.pow != 1 {
			str += "^" + strconv.Itoa(s.pow)
		}
	}
	return str
}
