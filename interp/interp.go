// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"sort"
)

// Predictor predicts the value of a function. It handles both
// interpolation and extrapolation.
type Predictor interface {
	// Predict returns the predicted value at x.
	Predict(float64) float64
}

// Fitter fits a predictor to data.
type Fitter interface {
	// Fit fits a predictor to (X, Y) value pairs provided as two slices.
	// It panics if len(xs) < 2, elements of xs are not strictly increasing or len(xs) != len(ys).
	Fit(xs, ys []float64)
}

// FittablePredictor is a Predictor which can fit itself to data.
type FittablePredictor interface {
	Fitter
	Predictor
}

// Constant predicts a constant value.
type Constant struct {
	// Constant Y value.
	Value float64
}

// Predict returns the predicted value at x.
func (c Constant) Predict(x float64) float64 {
	return c.Value
}

// PiecewiseLinear is a piecewise linear 1-dimensional interpolator.
// It extrapolates flat forwards (backwards) the last (first) known Y value.
type PiecewiseLinear struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64

	// Slopes of Y between neighbouring X values. len(slopes) + 1 == len(xs) == len(ys).
	slopes []float64
}

// Fit fits a piecewise linear predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing or len(xs) != len(ys).
func (pl *PiecewiseLinear) Fit(xs, ys []float64) {
	n := len(xs)
	if len(ys) != n {
		panic("interp: xs and ys have different lengths")
	}
	if n < 2 {
		panic("interp: too few points for interpolation")
	}
	m := n - 1
	pl.slopes = make([]float64, m)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			panic("interp: xs values not strictly increasing")
		}
		pl.slopes[i] = (ys[i+1] - ys[i]) / dx
	}
	pl.xs = xs
	pl.ys = ys
}

// Predict returns the interpolation value at x.
func (pl PiecewiseLinear) Predict(x float64) float64 {
	i := findSegment(pl.xs, x)
	if i < 0 {
		return pl.ys[0]
	}
	// i < len(pci.xs)
	xI := pl.xs[i]
	if x == xI {
		return pl.ys[i]
	}
	n := len(pl.xs)
	if i == n-1 {
		// x > li.xs[i]
		return pl.ys[n-1]
	}
	// i < len(i1d.xs) - 1
	return pl.ys[i] + pl.slopes[i]*(x-xI)
}

// PiecewiseConstant is a piecewise constant 1-dimensional interpolator.
// It extrapolates flat forwards (backwards) the last (first) known Y value.
type PiecewiseConstant struct {
	// Whether the interpolated function is left- or right-continuous.
	// If LeftContinuous == true, then y(xs[i]) == y(xs[i] - eps) for small eps > 0. Otherwise,
	// y(xs[i]) == y(xs[i] + eps).
	LeftContinuous bool

	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64
}

// Fit fits a piecewise constant predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing or len(xs) != len(ys).
func (pc *PiecewiseConstant) Fit(xs, ys []float64) {
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
	pc.xs = xs
	pc.ys = ys
}

// Predict returns the interpolation value at x.
func (pc PiecewiseConstant) Predict(x float64) float64 {
	i := findSegment(pc.xs, x)
	if i < 0 {
		return pc.ys[0]
	}
	// i < len(pci.xs)
	if x == pc.xs[i] {
		return pc.ys[i]
	}
	n := len(pc.xs)
	if i == n-1 {
		// x > pci.xs[i]
		return pc.ys[n-1]
	}
	if pc.LeftContinuous {
		return pc.ys[i+1]
	}
	return pc.ys[i]
}

// findSegment returns 0 <= i < len(xs) such that xs[i] <= x < xs[i + 1], or -1
// if no such interval containing x is found. It assumes that len(xs) >= 2
// without checking.
func findSegment(xs []float64, x float64) int {
	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
}
