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

// Interpolator1D represents a 1-dimensional interpolator.
// It interpolates a function over a defined closed interval range
// (which can be infinite on both sides).
type Interpolator1D interface {
	// Interval returns the valid closed interval for interpolation.
	Interval() r1.Interval

	// ValueAt returns the interpolation value at x.
	// It panics if x is not in the interpolation interval.
	ValueAt(float64) float64
}

// Constant1D is a constant 1D interpolator,
// It is defined over the range [-infinity, infinity].
type Constant1D struct {
	// Constant Y value.
	Value float64
}

// Interval returns the valid closed interval for interpolation.
func (ci Constant1D) Interval() r1.Interval {
	return r1.Interval{Min: math.Inf(-1), Max: math.Inf(1)}
}

// ValueAt returns the interpolation value at x.
func (ci Constant1D) ValueAt(x float64) float64 {
	return ci.Value
}

// Linear1D is a piecewise linear 1-dimensional interpolator.
type Linear1D struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64

	// Slopes of Y between neighbouring X values. len(slopes) + 1 == len(xs) == len(ys).
	slopes []float64
}

// NewLinear1D creates a new linear 1-dimensional interpolator.
// xs and ys should contain the X and Y values of interpolated nodes, respectively.
// NewLinear1D panics if len(xs) < 2, elements of xs are not strictly increasing or
// len(xs) != len(ys).
func NewLinear1D(xs []float64, ys []float64) *Linear1D {
	n := len(xs)
	if len(ys) != n {
		panic("interp: xs and ys have different lengths")
	}
	if n < 2 {
		panic("interp: too few points for interpolation")
	}
	m := n - 1
	slopes := make([]float64, m)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			panic("interp: xs values not strictly increasing")
		}
		slopes[i] = (ys[i+1] - ys[i]) / dx
	}
	return &Linear1D{xs, ys, slopes}
}

// Interval returns the valid closed interval for interpolation.
func (li Linear1D) Interval() r1.Interval {
	return r1.Interval{Min: li.xs[0], Max: li.xs[len(li.xs)-1]}
}

// ValueAt returns the interpolation value at x.
func (li Linear1D) ValueAt(x float64) float64 {
	i, xI := findSegment(li.xs, x)
	if x == xI {
		return li.ys[i]
	}
	// i < len(i1d.xs) - 1
	return li.ys[i] + li.slopes[i]*(x-xI)
}

// PiecewiseConstant1D is a piecewise constant 1-dimensional interpolator.
type PiecewiseConstant1D struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64

	// Whether the interpolated function is left- or right-continuous.
	leftContinuous bool
}

// NewPiecewiseConstant1D creates a new piecewise constant 1-dimensional interpolator.
// xs and ys should contain the X and Y values of interpolated nodes, respectively.
// If leftContinuous == true, then y(xs[i]) == y(xs[i] - eps) for small eps > 0. Otherwise,
// y(xs[i]) == y(xs[i] + eps).
// NewPiecewiseConstant1D panics if len(xs) < 2, elements of xs are not strictly increasing or
// len(xs) != len(ys).
func NewPiecewiseConstant1D(xs []float64, ys []float64, leftContinuous bool) *PiecewiseConstant1D {
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
	return &PiecewiseConstant1D{xs, ys, leftContinuous}
}

// Interval returns the valid closed interval for interpolation.
func (pci PiecewiseConstant1D) Interval() r1.Interval {
	return r1.Interval{Min: pci.xs[0], Max: pci.xs[len(pci.xs)-1]}
}

// ValueAt returns the interpolation value at x.
func (pci PiecewiseConstant1D) ValueAt(x float64) float64 {
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
