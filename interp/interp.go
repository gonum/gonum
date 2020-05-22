package interp

import (
	"fmt"
	"math"
	"sort"
)

// Interpolator1D represents a 1D interpolator interface.
//
// Interpolates a sequence of (x, y) pairs on a range [begin(), end()].
// Both begin() and end() can be infinite, but begin() must be lower than end().
type Interpolator1D interface {
	// Returns the lowest allowed argument for evaluate().
	begin() float64

	// Returns the highest allowed argument for evaluate().
	end() float64

	// Evaluates interpolated sequence at x in [begin(), end()].
	// Panics if the argument is outside this range.
	eval(float64) float64
}

// ConstInterpolator1D is a constant 1D interpolator,
// It is defined over the range [-infinity, infinity].
type ConstInterpolator1D struct {
	value float64
}

func (i1d ConstInterpolator1D) begin() float64 {
	return math.Inf(-1)
}

func (i1d ConstInterpolator1D) end() float64 {
	return math.Inf(1)
}

func (i1d ConstInterpolator1D) eval(x float64) float64 {
	return i1d.value
}

// NewConstInterpolator1D creates a new constant 1D interpolator with given value.
func NewConstInterpolator1D(value float64) *ConstInterpolator1D {
	return &ConstInterpolator1D{value}
}

// LinearInterpolator1D is a linear 1D interpolator, defined over the range [xs[0], xs[len(xs) - 1]].
type LinearInterpolator1D struct {
	xs     []float64
	ys     []float64
	slopes []float64
}

func (i1d LinearInterpolator1D) begin() float64 {
	return i1d.xs[0]
}

func (i1d LinearInterpolator1D) end() float64 {
	return i1d.xs[len(i1d.xs)-1]
}

// findSegment returns a tuple of: (i such that xs[i] <= x < xs[i + 1], xs[i]),
// or panics if such i is not found. Assumes that len(xs) >= 2 without checking.
func findSegment(xs []float64, x float64) (int, float64) {
	// Find minimum i s.t. xs[i] >= x, or len(xs) if not found.
	n := len(xs)
	i := sort.Search(n, func(i int) bool { return xs[i] > x })
	if i < n {
		if i == 0 {
			// x < begin()
			panic(fmt.Sprintf("interp: x value %g below lower bound %g", x, xs[0]))
		} else {
			return i - 1, xs[i-1]
		}
	} else {
		if xs[n-1] == x {
			return n - 1, x
		}
		panic(fmt.Sprintf("interp: x value %g above upper bound %g", x, xs[n-1]))
	}
}

func (i1d LinearInterpolator1D) eval(x float64) float64 {
	i, xI := findSegment(i1d.xs, x)
	if x == xI {
		return i1d.ys[i]
	}
	// i < len(i1d.xs) - 1
	return i1d.ys[i] + i1d.slopes[i]*(x-xI)
}

// NewLinearInterpolator1D creates a new linear 1D interpolator.
// xs and ys should contain the X and Y values of interpolated nodes, respectively.
// Panics if len(xs) < 2, elements of xs are not strictly increasing or
// len(xs) != len(ys).
func NewLinearInterpolator1D(xs []float64, ys []float64) *LinearInterpolator1D {
	n := len(xs)
	if len(ys) != n {
		panic("interp: xs and ys have different lengths")
	}
	if n < 2 {
		panic("interp: too few points for interpolation")
	}
	slopes := make([]float64, n-1)
	for i := 1; i < n; i++ {
		dx := xs[i] - xs[i-1]
		if dx <= 0 {
			panic("interp: x values not strictly increasing")
		}
		slopes[i-1] = (ys[i] - ys[i-1]) / dx
	}
	return &LinearInterpolator1D{xs, ys, slopes}
}
