// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"errors"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

const (
	differentLengths        = "interp: input slices have different lengths"
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
	pl.xs = make([]float64, n)
	pl.ys = make([]float64, n)
	copy(pl.xs, xs)
	copy(pl.ys, ys)
	return nil
}

// Predict returns the interpolation value at x.
func (pl PiecewiseLinear) Predict(x float64) float64 {
	i := findSegment(pl.xs, x)
	if i < 0 {
		return pl.ys[0]
	}
	xI := pl.xs[i]
	if x == xI {
		return pl.ys[i]
	}
	n := len(pl.xs)
	if i == n-1 {
		return pl.ys[n-1]
	}
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
	pc.xs = make([]float64, n)
	pc.ys = make([]float64, n)
	copy(pc.xs, xs)
	copy(pc.ys, ys)
	return nil
}

// Predict returns the interpolation value at x.
func (pc PiecewiseConstant) Predict(x float64) float64 {
	i := findSegment(pc.xs, x)
	if i < 0 {
		return pc.ys[0]
	}
	if x == pc.xs[i] {
		return pc.ys[i]
	}
	n := len(pc.xs)
	if i == n-1 {
		return pc.ys[n-1]
	}
	return pc.ys[i+1]
}

// PiecewiseCubic is a left-continuous, piecewise cubic 1-dimensional interpolator.
type PiecewiseCubic struct {
	// Interpolated X values.
	xs []float64

	// Coefficients of interpolating cubic polynomials, with
	// len(xs) - 1 rows and 4 columns. The interpolated value
	// for xs[i] <= x < xs[i + 1] is defined as
	//   sum_{k = 0}^3 coeffs.At(i, k) * (x - xs[i])^k
	// To guarantee left-continuity, coeffs.At(i, 0) == ys[i].
	coeffs *mat.Dense

	// Last interpolated Y value, corresponding to xs[len(xs) - 1].
	lastY float64
}

// Predict returns the interpolation value at x.
func (pc PiecewiseCubic) Predict(x float64) float64 {
	i := findSegment(pc.xs, x)
	if i < 0 {
		return pc.coeffs.At(0, 0)
	}
	m := len(pc.xs) - 1
	if x == pc.xs[i] {
		if i < m {
			return pc.coeffs.At(i, 0)
		}
		return pc.lastY
	}
	if i == m {
		return pc.lastY
	}
	dx := x - pc.xs[i]
	a := pc.coeffs.RawRowView(i)
	return ((a[3]*dx+a[2])*dx+a[1])*dx + a[0]
}

// FitWithDerivatives fits a piecewise cubic predictor to (X, Y, dY/dX) value
// triples provided as three slices.
// It returns an error if len(xs) < 2, elements of xs are not strictly
// increasing, len(xs) != len(ys) or len(xs) != len(dydxs).
func (pc *PiecewiseCubic) FitWithDerivatives(xs, ys, dydxs []float64) error {
	n := len(xs)
	if len(ys) != n {
		return errors.New(differentLengths)
	}
	if len(dydxs) != n {
		return errors.New(differentLengths)
	}
	if n < 2 {
		return errors.New(tooFewPoints)
	}
	m := n - 1
	pc.coeffs = mat.NewDense(m, 4, nil)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			return errors.New(xsNotStrictlyIncreasing)
		}
		dy := ys[i+1] - ys[i]
		// a_0
		pc.coeffs.Set(i, 0, ys[i])
		// a_1
		pc.coeffs.Set(i, 1, dydxs[i])
		// Solve a linear equation system for a_2 and a_3.
		pc.coeffs.Set(i, 2, (3*dy-(2*dydxs[i]+dydxs[i+1])*dx)/dx/dx)
		pc.coeffs.Set(i, 3, (-2*dy+(dydxs[i]+dydxs[i+1])*dx)/dx/dx/dx)
	}
	pc.xs = make([]float64, n)
	copy(pc.xs, xs)
	pc.lastY = ys[m]
	return nil
}

// AkimaSplines is a left-continuous, piecewise cubic 1-dimensional interpolator
// which can be fitted to (X, Y) value pairs without providing derivatives.
// See https://www.iue.tuwien.ac.at/phd/rottinger/node60.html for more details.
type AkimaSplines struct {
	PiecewiseCubic
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It returns an error if len(xs) < 2, elements of xs are not strictly
// increasing or len(xs) != len(ys).
// If len(xs) == 2, we set both derivatives dY/dX to the slope
// (ys[1] - ys[0]) / (xs[1] - xs[0]).
func (as *AkimaSplines) Fit(xs, ys []float64) error {
	n := len(xs)
	if len(ys) != n {
		return errors.New(differentLengths)
	}
	dydxs := make([]float64, n)

	if n == 2 {
		slope := (ys[1] - ys[0]) / (xs[1] - xs[0])
		dydxs[0] = slope
		dydxs[1] = slope
		return as.FitWithDerivatives(xs, ys, dydxs)
	}

	m := n - 1
	slopes := make([]float64, m+4)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			return errors.New(xsNotStrictlyIncreasing)
		}
		slopes[i+2] = (ys[i+1] - ys[i]) / dx
	}
	slopes[1] = 2*slopes[2] - slopes[3]
	slopes[0] = 3*slopes[2] - 2*slopes[3]
	slopes[m+2] = 2*slopes[n] - slopes[m]
	slopes[m+3] = 3*slopes[n] - 2*slopes[m]
	for i := 0; i < n; i++ {
		wLeft := math.Abs(slopes[i+2] - slopes[i+3])
		wRight := math.Abs(slopes[i+1] - slopes[i])
		w := wLeft + wRight
		if w > 0 {
			dydxs[i] = (wLeft*slopes[i+1] + wRight*slopes[i+2]) / w
		} else {
			dydxs[i] = (slopes[i+1] + slopes[i+2]) / 2
		}
	}
	return as.FitWithDerivatives(xs, ys, dydxs)
}

// findSegment returns 0 <= i < len(xs) such that xs[i] <= x < xs[i + 1], where xs[len(xs)]
// is assumed to be +Inf. If no such i is found, it returns -1. It assumes that len(xs) >= 2
// without checking.
func findSegment(xs []float64, x float64) int {
	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
}
