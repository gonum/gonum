// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/spatial/r1"
)

// Interpolator1D represents a 1D interpolator.
// It interpolates a y(x) function based on some data over a range [begin(), end()].
// Both begin() and end() can be infinite, but begin() must be lower than end().
type Interpolator1D interface {

	// Returns the closed interval over which one can interpolate.
	Interval() r1.Interval

	// Evaluates interpolated sequence at x in [begin(), end()].
	// Panics if the argument is outside this range.
	Eval(float64) float64
}

// ConstInterpolator1D is a constant 1D interpolator,
// It is defined over the range [-infinity, infinity].
type ConstInterpolator1D struct {
	// Constant Y value.
	Value float64
}

// Interval implements Interpolator1D.Interval.
func (ci ConstInterpolator1D) Interval() r1.Interval {
	return r1.Interval{Min: math.Inf(-1), Max: math.Inf(1)}
}

// Eval implements Interpolator1D.Eval.
func (ci ConstInterpolator1D) Eval(x float64) float64 {
	return ci.Value
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

// LinearInterpolator1D is a piecewise linear 1D interpolator, defined
// over the range [xs[0], xs[len(xs) - 1]].
type LinearInterpolator1D struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64

	// Slopes of Y between neighbouring X values. len(slopes) + 1 == len(xs) == len(ys).
	slopes []float64
}

// NewLinearInterpolator1D creates a new linear 1D interpolator.
// xs and ys should contain the X and Y values of interpolated nodes, respectively.
// Panics if len(xs) < 2, elements of xs are not strictly increasing or
// len(xs) != len(ys).
func NewLinearInterpolator1D(xs []float64, ys []float64) *LinearInterpolator1D {
	validateXsAndYs(xs, ys)
	m := len(xs) - 1
	slopes := make([]float64, m)
	for i := 0; i < m; i++ {
		slopes[i] = (ys[i+1] - ys[i]) / (xs[i+1] - xs[i])
	}
	return &LinearInterpolator1D{xs, ys, slopes}
}

// Interval implements Interpolator1D.Interval.
func (li LinearInterpolator1D) Interval() r1.Interval {
	return r1.Interval{Min: li.xs[0], Max: li.xs[len(li.xs)-1]}
}

// Eval implements Interpolator1D.Eval.
func (li LinearInterpolator1D) Eval(x float64) float64 {
	i, xI := findSegment(li.xs, x)
	if x == xI {
		return li.ys[i]
	}
	// i < len(i1d.xs) - 1
	return li.ys[i] + li.slopes[i]*(x-xI)
}

// validateXsAndYs panics if xs and ys do not satify common requirements
// for piecewise 1D interpolators.
func validateXsAndYs(xs []float64, ys []float64) {
	n := len(xs)
	if len(ys) != n {
		panic("interp: xs and ys have different lengths")
	}
	if n < 2 {
		panic("interp: too few points for interpolation")
	}
	for i := 1; i < n; i++ {
		if xs[i] <= xs[i-1] {
			panic("interp: xs values not strictly increasing")
		}
	}
}

// PiecewiseConstInterpolator1D is a piecewise constant 1D interpolator.
// It is defined over [xs[0], xs[len(xs)-1]].
// If leftContinuous == true, then y(xs[i]) == y(xs[i] - eps) for small
// eps > 0. Otherwise, y(xs[i]) == y(xs[i] + eps).
type PiecewiseConstInterpolator1D struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64

	// Whether the interpolated function is left- or right-continuous.
	leftContinuous bool
}

// NewPiecewiseConstInterpolator1D creates a new piecewise constant 1D interpolator.
// xs and ys should contain the X and Y values of interpolated nodes, respectively.
// Panics if len(xs) < 2, elements of xs are not strictly increasing or
// len(xs) != len(ys).
func NewPiecewiseConstInterpolator1D(xs []float64, ys []float64, leftContinuous bool) *PiecewiseConstInterpolator1D {
	validateXsAndYs(xs, ys)
	return &PiecewiseConstInterpolator1D{xs, ys, leftContinuous}
}

// Interval implements Interpolator1D.Interval.
func (pci PiecewiseConstInterpolator1D) Interval() r1.Interval {
	return r1.Interval{Min: pci.xs[0], Max: pci.xs[len(pci.xs)-1]}
}

// Eval implements Interpolator1D.Eval.
func (pci PiecewiseConstInterpolator1D) Eval(x float64) float64 {
	i, xI := findSegment(pci.xs, x)
	if x == xI {
		return pci.ys[i]
	}
	// i < len(i1d.xs) - 1
	if pci.leftContinuous {
		return pci.ys[i+1]
	}
	return pci.ys[i]
}
