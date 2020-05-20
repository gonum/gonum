package interp

import "math"

// 1D interpolator interface.
//
// Interpolates a sequence of (x, y) pairs on a range [start(), end()].
// Both start() and end() can be infinite, but start() must be lower
// than end().
type Interpolator1D interface {
	// Returns the lowest allowed argument for evaluate().
	start() float64

	// Returns the highest allowed argument for evaluate().
	end() float64

	// Evaluates interpolated sequence at x in [start(), end()].
	// Panics if the argument is outside this range.
	eval(float64) float64
}

// Constant 1D interpolator.
// Defined over the range [-infinity, infinity].
type ConstInterpolator1D struct {
	value float64
}

func (i1d ConstInterpolator1D) start() float64 {
	return math.Inf(-1)
}

func (i1d ConstInterpolator1D) end() float64 {
	return math.Inf(1)
}

func (i1d ConstInterpolator1D) eval(x float64) float64 {
	return i1d.value
}

// Creates new constant 1D interpolator with given value.
func NewConstInterpolator1D(value float64) *ConstInterpolator1D {
	return &ConstInterpolator1D{value}
}
