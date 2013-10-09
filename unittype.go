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
// domain-specific problem, NewDimension should be used. Integers
// should never be cast as type dimension
//		// Good: Create a package constant with an init function
//		var MyDimension unit.Dimension
//		init(){
//				MyDimension = NewDimension("my")
//		}
//		main(){
//			var := MyDimension(28.2)
//      }
type Dimension int

const (
	reserved Dimension = iota
	// SI Base Units
	CurrentDim Dimension
	LengthDim
	LuminousIntensityDim
	MassDim
	TemperatureDim
	TimeDim
	// Start of other SI Units
	AngleDim             // e.g. radians
	lastPackageDimension // Used in create dimension
)

var lastCreatedDimension Dimension = lastPackageDimension
var dimensionToSymbol map[Dimension]string = make(map[Dimension]string) // for printing
var symbolToDimension map[string]Dimension = make(map[string]Dimension) // for guaranteeing there aren't two identical symbols

// TODO: Should we actually reserve "common" SI unit symbols ("N", "J", etc.) so there isn't confusion
// TODO: If we have a fancier ParseUnit, maybe the 'reserved' symbols should be a different map
// 		map[string]string which says how they go?
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

	// Reserve common SI symbols
	// base units
	symbolToDimension["mol"] = reserved
	// prefixes
	symbolToDimension["Y"] = reserved
	symbolToDimension["Z"] = reserved
	symbolToDimension["E"] = reserved
	symbolToDimension["P"] = reserved
	symbolToDimension["T"] = reserved
	symbolToDimension["G"] = reserved
	symbolToDimension["M"] = reserved
	symbolToDimension["k"] = reserved
	symbolToDimension["h"] = reserved
	symbolToDimension["da"] = reserved
	symbolToDimension["d"] = reserved
	symbolToDimension["c"] = reserved
	symbolToDimension["m"] = reserved
	symbolToDimension["μ"] = reserved
	symbolToDimension["n"] = reserved
	symbolToDimension["p"] = reserved
	symbolToDimension["f"] = reserved
	symbolToDimension["a"] = reserved
	symbolToDimension["z"] = reserved
	symbolToDimension["y"] = reserved
	// SI Derived units with special symbols
	symbolToDimension["sr"] = reserved
	symbolToDimension["F"] = reserved
	symbolToDimension["C"] = reserved
	symbolToDimension["S"] = reserved
	symbolToDimension["H"] = reserved
	symbolToDimension["V"] = reserved
	symbolToDimension["Ω"] = reserved
	symbolToDimension["J"] = reserved
	symbolToDimension["N"] = reserved
	symbolToDimension["Hz"] = reserved
	symbolToDimension["lx"] = reserved
	symbolToDimension["lm"] = reserved
	symbolToDimension["Wb"] = reserved
	symbolToDimension["T"] = reserved
	symbolToDimension["W"] = reserved
	symbolToDimension["Pa"] = reserved
	symbolToDimension["Bq"] = reserved
	symbolToDimension["Gy"] = reserved
	symbolToDimension["Sv"] = reserved
	symbolToDimension["kat"] = reserved
	// Units in use with SI
	symbolToDimension["ha"] = reserved
	symbolToDimension["L"] = reserved
	symbolToDimension["l"] = reserved
	// Units in Use Temporarily with SI
	symbolToDimension["bar"] = reserved
	symbolToDimension["b"] = reserved
	symbolToDimension["Ci"] = reserved
	symbolToDimension["R"] = reserved
	symbolToDimension["rd"] = reserved
	symbolToDimension["rem"] = reserved
}

// Dimensions represent the dimensionality of the unit in powers
// of that dimension. If a key is not present, the power of that
// dimension is zero. Dimensions is used in conjuction with NewUnit
type Dimensions map[Dimension]int

var newUnitMutex *sync.Mutex = &sync.Mutex{} // so there is no race condition for dimension

// NewDimension returns a new dimension variable which will have a
// unique representation across packages to prevent accidental overlap.
// The input string represents a symbol name which will be used for printing
// Unit types. This symbol may not overlap with any of the SI base units
// or other symbols of common use in SI ("kg", "J", "μ", etc.). A list of
// such symbols can be found at http://lamar.colostate.edu/~hillger/basic.htm or
// by consulting the package source. Furthermore, the provided symbol is also
// forbidden from overlapping with other packages calling NewDimension. NewDimension
// is expecting to be used only during initialization, and as such it will panic
// if the symbol matching an existing symbol
// NewDimension should only be called for unit types that are actually orthogonal
// to the base dimensions defined in this package. Please see the package-level
// documentation for further explanation
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
// an acceleration with a mass to get a force. Please see the
// package documentation for further explanation.
type Unit struct {
	dimensions map[Dimension]int // Map for custom dimensions
	value      float64
}

// NewUnit creates a new variable of type Unit which has the value
// specified by value and the dimensions specified by the
// base units struct. The value is always in SI Units.
//
// Example: To create an acceleration of 3 m/s^2, one could do
// myvar := CreateUnit(3.0, &Dimensions{unit.LengthDim: 1, unit.TimeDim: -2})
func NewUnit(value float64, d Dimensions) *Unit {
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

// Unit implements the Uniter interface
func (u *Unit) Unit() *Unit {
	return u
}

// Mul multiply the receiver by the input changing the dimensions
// of the receiver as appropriate. The input is not changed
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
// method is, in general, not recommended, though it can be useful
// for printing. Instead, the FromUnit type of a specific
// dimension should be used to guarantee dimension consistency
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

// String makes Unit satisfy the stringer interface. The unit is printed
// using strconv.FormatFloat(unit.value, 'e', -1, 64), with dimensions
// appended. If the power if the dimension is not zero or one,
// symbol^power is appended, if the power is one, just the symbol is appended
// and if the power is zero, nothing is appended. Dimensions are appended
// in order by symbol name.
func (u Unit) String() string {
	str := strconv.FormatFloat(u.value, 'e', -1, 64)
	// Map iterates randomly, but print should be in a fixed order. Can't use
	// dimension number, because for user-defined dimension that number may
	// not be fixed from run to run.
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
