// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"fmt"
	"sort"
)

// notIncreasingError represents an error where v1 <= v2.
type notIncreasingError struct {
	v1 float64
	v2 float64
}

// newNotIncreasingError returns a new notIncreasingError.
func newNotIncreasingError(v1, v2 float64) *notIncreasingError {
	return &notIncreasingError{v1: v1, v2: v2}
}

// Error implements error.
func (e notIncreasingError) Error() string {
	return fmt.Sprintf("interp: xs values not strictly increasing: %f <= %f", e.v1, e.v2)
}

// validate returns and error if:
// * len(xs) < min
// * len(xs) != len of any of slices
// otherwise it returns nil
func validate(min int, xs []float64, slices ...[]float64) error {
	if len(xs) < min {
		return fmt.Errorf("interp: xs has too few points %d for interpolation need %d", len(xs), min)
	}

	for _, s := range slices {
		if len(xs) != len(s) {
			return fmt.Errorf("interp: input slices have different lengths: %d != %d", len(xs), len(s))
		}
	}

	return nil
}

// Predictor predicts the value of a function. It handles both
// interpolation and extrapolation.
type Predictor interface {
	// Predict returns the predicted value at x.
	Predict(x float64) float64
}

// Fitter fits a predictor to data.
type Fitter interface {
	// Fit fits a predictor to (X, Y) value pairs provided as two slices.
	// It returns and error if len(xs) < 2, elements of xs are not strictly increasing,
	// len(xs) != len(ys) or if fitting fails.
	Fit(xs, ys []float64) error
}

// FittablePredictor is a Predictor which can fit itself to data.
type FittablePredictor interface {
	Fitter
	Predictor
}

// DerivativePredictor predicts both the value and the derivative of
// a function. It handles both interpolation and extrapolation.
type DerivativePredictor interface {
	Predictor

	// PredictDerivative returns the predicted derivative at x.
	PredictDerivative(x float64) float64
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
// It returns an error if len(xs) < 2, elements of xs are not strictly increasing,
// len(xs) != len(ys).
func (pl *PiecewiseLinear) Fit(xs, ys []float64) (err error) {
	pl.slopes, err = calculateSlopes(xs, ys)
	if err != nil {
		return err
	}

	pl.xs = make([]float64, len(xs))
	pl.ys = make([]float64, len(ys))
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

// PiecewiseConstant is a left-continuous, piecewise constant
// 1-dimensional interpolator.
type PiecewiseConstant struct {
	// Interpolated X values.
	xs []float64

	// Interpolated Y data values, same len as ys.
	ys []float64
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It returns an error if len(xs) < 2, elements of xs are not strictly increasing,
// len(xs) != len(ys).
func (pc *PiecewiseConstant) Fit(xs, ys []float64) error {
	if err := validate(2, xs, ys); err != nil {
		return err
	}

	n := len(xs)
	for i := 1; i < n; i++ {
		if xs[i] <= xs[i-1] {
			return newNotIncreasingError(xs[i], xs[i-1])
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

// findSegment returns 0 <= i < len(xs) such that xs[i] <= x < xs[i + 1], where xs[len(xs)]
// is assumed to be +Inf. If no such i is found, it returns -1. It assumes that len(xs) >= 2
// without checking.
func findSegment(xs []float64, x float64) int {
	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
}

// calculateSlopes calculates slopes (ys[i+1] - ys[i]) / (xs[i+1] - xs[i]).
// It returns an error if len(xs) < 2, elements of xs are not strictly increasing
// or len(xs) != len(ys).
func calculateSlopes(xs, ys []float64) ([]float64, error) {
	if err := validate(2, xs, ys); err != nil {
		return nil, err
	}

	n := len(xs)
	m := n - 1
	slopes := make([]float64, m)
	prevX := xs[0]
	prevY := ys[0]
	for i := 0; i < m; i++ {
		x := xs[i+1]
		y := ys[i+1]
		dx := x - prevX
		if dx <= 0 {
			return nil, newNotIncreasingError(x, prevX)
		}
		slopes[i] = (y - prevY) / dx
		prevX = x
		prevY = y
	}
	return slopes, nil
}
