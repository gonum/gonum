// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"fmt"
	"math"
	"testing"
)

func TestConstant(t *testing.T) {
	t.Parallel()
	const value = 42.0
	c := Constant(value)
	xs := []float64{math.Inf(-1), -11, 0.4, 1e9, math.Inf(1)}
	for _, x := range xs {
		y := c.Predict(x)
		if y != value {
			t.Errorf("unexpected Predict(%g) value: got: %g want: %g", x, y, value)
		}
	}
}

func TestFunction(t *testing.T) {
	fn := func(x float64) float64 { return math.Exp(x) }
	predictor := Function(fn)
	xs := []float64{-100, -1, 0, 0.5, 15}
	for _, x := range xs {
		want := fn(x)
		got := predictor.Predict(x)
		if got != want {
			t.Errorf("unexpected Predict(%g) value: got: %g want: %g", x, got, want)
		}
	}
}

func TestFindSegment(t *testing.T) {
	t.Parallel()
	xs := []float64{0, 1, 2}
	testXs := []float64{-0.6, 0, 0.3, 1, 1.5, 2, 2.8}
	expectedIs := []int{-1, 0, 0, 1, 1, 2, 2}
	for k, x := range testXs {
		i := findSegment(xs, x)
		if i != expectedIs[k] {
			t.Errorf("unexpected value of findSegment(xs, %g): got %d want: %d", x, i, expectedIs[k])
		}
	}
}

func BenchmarkFindSegment(b *testing.B) {
	xs := []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 12, 13.5, 16.5}
	for i := 0; i < b.N; i++ {
		findSegment(xs, 0)
		findSegment(xs, 16.5)
		findSegment(xs, -1)
		findSegment(xs, 8.25)
		findSegment(xs, 4.125)
		findSegment(xs, 13.6)
		findSegment(xs, 23.6)
		findSegment(xs, 13.5)
		findSegment(xs, 6)
		findSegment(xs, 4.5)
	}
}

// testPiecewiseInterpolatorCreation tests common functionality in creating piecewise  interpolators.
func testPiecewiseInterpolatorCreation(t *testing.T, fp FittablePredictor) {
	type errorParams struct {
		xs              []float64
		ys              []float64
		expectedMessage string
	}
	errorParamSets := []errorParams{
		{[]float64{0, 1, 2}, []float64{-0.5, 1.5}, "xs and ys have different lengths"},
		{[]float64{0.3}, []float64{0}, "too few points for interpolation"},
		{[]float64{0.3, 0.3}, []float64{0, 0}, "xs values not strictly increasing"},
		{[]float64{0.3, -0.3}, []float64{0, 0}, "xs values not strictly increasing"},
	}
	for _, params := range errorParamSets {
		err := fp.Fit(params.xs, params.ys)
		expectedMessage := fmt.Sprintf("interp: %s", params.expectedMessage)
		if err == nil || err.Error() != expectedMessage {
			t.Errorf("expected error for xs: %v and ys: %v with message: %s", params.xs, params.ys, expectedMessage)
		}
	}
}

func TestPiecewiseLinearFit(t *testing.T) {
	t.Parallel()
	testPiecewiseInterpolatorCreation(t, &PiecewiseLinear{})
}

// testInterpolatorPredict tests evaluation of a  interpolator.
func testInterpolatorPredict(t *testing.T, p Predictor, xs []float64, expectedYs []float64, tol float64) {
	for i, x := range xs {
		y := p.Predict(x)
		yErr := math.Abs(y - expectedYs[i])
		if yErr > tol {
			if tol == 0 {
				t.Errorf("unexpected Predict(%g) value: got: %g want: %g", x, y, expectedYs[i])
			} else {
				t.Errorf("unexpected Predict(%g) value: got: %g want: %g with tolerance: %g", x, y, expectedYs[i], tol)
			}
		}
	}
}

func TestPiecewiseLinearPredict(t *testing.T) {
	t.Parallel()
	xs := []float64{0, 1, 2}
	ys := []float64{-0.5, 1.5, 1}
	var pl PiecewiseLinear
	err := pl.Fit(xs, ys)
	if err != nil {
		t.Errorf("Fit error: %s", err.Error())
	}
	testInterpolatorPredict(t, pl, xs, ys, 0)
	testInterpolatorPredict(t, pl, []float64{-0.4, 2.6}, []float64{-0.5, 1}, 0)
	testInterpolatorPredict(t, pl, []float64{0.1, 0.5, 0.8, 1.2}, []float64{-0.3, 0.5, 1.1, 1.4}, 1e-15)
}

func BenchmarkNewPiecewiseLinear(b *testing.B) {
	xs := []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 12, 13.5, 16.5}
	ys := []float64{0, 1, 2, 2.5, 2, 1.5, 4, 10, -2, 2}
	var pl PiecewiseLinear
	for i := 0; i < b.N; i++ {
		_ = pl.Fit(xs, ys)
	}
}

func BenchmarkPiecewiseLinearPredict(b *testing.B) {
	xs := []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 12, 13.5, 16.5}
	ys := []float64{0, 1, 2, 2.5, 2, 1.5, 4, 10, -2, 2}
	var pl PiecewiseLinear
	_ = pl.Fit(xs, ys)
	for i := 0; i < b.N; i++ {
		pl.Predict(0)
		pl.Predict(16.5)
		pl.Predict(-2)
		pl.Predict(4)
		pl.Predict(7.32)
		pl.Predict(9.0001)
		pl.Predict(1.4)
		pl.Predict(1.6)
		pl.Predict(30)
		pl.Predict(13.5)
		pl.Predict(4.5)
	}
}

func TestNewPiecewiseConstant(t *testing.T) {
	var pc PiecewiseConstant
	testPiecewiseInterpolatorCreation(t, &pc)
}

func benchmarkPiecewiseConstantPredict(b *testing.B) {
	xs := []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 12, 13.5, 16.5}
	ys := []float64{0, 1, 2, 2.5, 2, 1.5, 4, 10, -2, 2}
	var pc PiecewiseConstant
	_ = pc.Fit(xs, ys)
	for i := 0; i < b.N; i++ {
		pc.Predict(0)
		pc.Predict(16.5)
		pc.Predict(4)
		pc.Predict(7.32)
		pc.Predict(9.0001)
		pc.Predict(1.4)
		pc.Predict(1.6)
		pc.Predict(13.5)
		pc.Predict(4.5)
	}
}

func BenchmarkPiecewiseConstantPredict(b *testing.B) {
	benchmarkPiecewiseConstantPredict(b)
}

func TestPiecewiseConstantPredict(t *testing.T) {
	t.Parallel()
	xs := []float64{0, 1, 2}
	ys := []float64{-0.5, 1.5, 1}
	var pc PiecewiseConstant
	err := pc.Fit(xs, ys)
	if err != nil {
		t.Errorf("Fit error: %s", err.Error())
	}
	testInterpolatorPredict(t, pc, xs, ys, 0)
	testXs := []float64{-0.9, 0.1, 0.5, 0.8, 1.2, 3.1}
	leftYs := []float64{-0.5, 1.5, 1.5, 1.5, 1, 1}
	testInterpolatorPredict(t, pc, testXs, leftYs, 0)
}
