// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"errors"
	"sort"

	"gonum.org/v1/gonum/mat"
)

const (
	differentLengths        = "interp: xs and ys have different lengths"
	tooFewPoints            = "interp: too few points for interpolation"
	xsNotStrictlyIncreasing = "interp: xs values not strictly increasing"
)

// Predictor predicts the value of a function. It handles both
// interpolation and extrapolation.
type Predictor interface {
	// Predict returns the predicted value at x.
	Predict(x float64) float64
}

// Fitter fits a predictor to data.
type Fitter interface {
	// Fit fits a predictor to (X, Y) value pairs provided as two slices.
	// It returns an error if len(xs) < 2, elements of xs are not strictly
	// increasing or len(xs) != len(ys).
	Fit(xs, ys []float64) error
}

// FittablePredictor is a Predictor which can fit itself to data.
type FittablePredictor interface {
	Fitter
	Predictor
}

// Constant predicts a constant value.
type Constant float64

// Predict returns the predicted value at x.
func (c Constant) Predict(x float64) float64 {
	return float64(c)
}

// Function predicts by evaluating itself.
type Function func(float64) float64

// Predict returns the predicted value at x by evaluating fn(x).
func (fn Function) Predict(x float64) float64 {
	return fn(x)
}

// PiecewiseLinear is a piecewise linear 1-dimensional interpolator.
type PiecewiseLinear struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64

	// Slopes of Y between neighbouring X values. len(slopes) + 1 == len(xs) == len(ys).
	slopes []float64
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It returns an error if len(xs) < 2, elements of xs are not strictly
// increasing or len(xs) != len(ys).
func (pl *PiecewiseLinear) Fit(xs, ys []float64) error {
	n := len(xs)
	if len(ys) != n {
		return errors.New(differentLengths)
	}
	if n < 2 {
		return errors.New(tooFewPoints)
	}
	m := n - 1
	pl.slopes = make([]float64, m)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			return errors.New(xsNotStrictlyIncreasing)
		}
		pl.slopes[i] = (ys[i+1] - ys[i]) / dx
	}
	pl.xs = xs
	pl.ys = ys
	return nil
}

// Predict returns the interpolation value at x.
func (pl PiecewiseLinear) Predict(x float64) float64 {
	i := findSegment(pl.xs, x)
	if i < 0 {
		return pl.ys[0]
	}
	// i < len(pl.xs)
	xI := pl.xs[i]
	if x == xI {
		return pl.ys[i]
	}
	n := len(pl.xs)
	if i == n-1 {
		// x > pl.xs[i]
		return pl.ys[n-1]
	}
	// i < len(pl.xs) - 1
	return pl.ys[i] + pl.slopes[i]*(x-xI)
}

// PiecewiseConstant is a left-continous, piecewise constant
// 1-dimensional interpolator.
type PiecewiseConstant struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It returns an error if len(xs) < 2, elements of xs are not strictly
// increasing or len(xs) != len(ys).
func (pc *PiecewiseConstant) Fit(xs, ys []float64) error {
	n := len(xs)
	if len(ys) != n {
		return errors.New(differentLengths)
	}
	if n < 2 {
		return errors.New(tooFewPoints)
	}
	for i := 1; i < n; i++ {
		if xs[i] <= xs[i-1] {
			return errors.New(xsNotStrictlyIncreasing)
		}
	}
	pc.xs = xs
	pc.ys = ys
	return nil
}

// Predict returns the interpolation value at x.
func (pc PiecewiseConstant) Predict(x float64) float64 {
	i := findSegment(pc.xs, x)
	if i < 0 {
		return pc.ys[0]
	}
	// i < len(pc.xs)
	if x == pc.xs[i] {
		return pc.ys[i]
	}
	n := len(pc.xs)
	if i == n-1 {
		// x > pc.xs[i]
		return pc.ys[n-1]
	}
	return pc.ys[i+1]
}

// PiecewiseCubic is a left-continuous, piecewise cubic
// 1-dimensional interpolator.
type PiecewiseCubic struct {
	// Interpolated X values.
	xs []float64

	// Coefficients of interpolating cubic polynomials, with
	// len(xs) - 1 rows and 4 columns. The interpolated value
	// for xs[i] <= x < xs[i + 1] is defined as
	//   sum_{k = 0}^3 coeffs.At(i, k) * (x - xs[i])^k
	// To guarantee left-continuity, coeffs.At(i, 0) == ys[i].
	coeffs mat.Dense

	// Last interpolated Y value, corresponding to xs[len(xs) - 1].
	lastY float64
}

// Predict returns the interpolation value at x.
func (pc *PiecewiseCubic) Predict(x float64) float64 {
	i := findSegment(pc.xs, x)
	if i < 0 {
		return pc.coeffs.At(0, 0)
	}
	// i < len(pc.xs)
	if x == pc.xs[i] {
		return pc.coeffs.At(0, 0)
	}
	n := len(pc.xs)
	if i == n-1 {
		// x > pc.xs[i]
		return pc.lastY
	}
	dx := x - pc.xs[i]
	a := pc.coeffs.RawRowView(i)
	return ((a[3]*dx+a[2])*dx+a[1])*dx + a[0]
}

// findSegment returns 0 <= i < len(xs) such that xs[i] <= x < xs[i + 1], where xs[len(xs)]
// is assumed to be +Inf. If no such i is found, it returns -1. It assumes that len(xs) >= 2
// without checking.
func findSegment(xs []float64, x float64) int {
	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
}
