// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
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

func TestFindSegmentEdgeCases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		xs   []float64
		x    float64
		want int
	}{
		{xs: nil, x: 0, want: -1},
		{xs: []float64{0}, x: -1, want: -1},
		{xs: []float64{0}, x: 0, want: 0},
		{xs: []float64{0}, x: 1, want: 0},
	}

	for _, test := range cases {
		if got := findSegment(test.xs, test.x); got != test.want {
			t.Errorf("unexpected value of findSegment(%v, %f): got %d want: %d",
				test.xs, test.x, got, test.want)
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
		xs []float64
		ys []float64
	}
	errorParamSets := []errorParams{
		{[]float64{0, 1, 2}, []float64{-0.5, 1.5}},
		{[]float64{0.3}, []float64{0}},
		{[]float64{0.3, 0.3}, []float64{0, 0}},
		{[]float64{0.3, -0.3}, []float64{0, 0}},
	}
	for _, params := range errorParamSets {
		if !panics(func() { _ = fp.Fit(params.xs, params.ys) }) {
			t.Errorf("expected panic for xs: %v and ys: %v", params.xs, params.ys)
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

func TestCalculateSlopesErrors(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		xs, ys []float64
	}{
		{
			xs: []float64{0},
			ys: []float64{0},
		},
		{
			xs: []float64{0, 1, 2},
			ys: []float64{0, 1}},
		{
			xs: []float64{0, 0, 1},
			ys: []float64{0, 0, 0},
		},
		{
			xs: []float64{0, 1, 0},
			ys: []float64{0, 0, 0},
		},
	} {
		if !panics(func() { calculateSlopes(test.xs, test.ys) }) {
			t.Errorf("expected panic for xs: %v and ys: %v", test.xs, test.ys)
		}
	}
}

func TestCalculateSlopes(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		xs, ys, want []float64
	}{
		{
			xs:   []float64{0, 2, 3, 5},
			ys:   []float64{0, 1, 1, -1},
			want: []float64{0.5, 0, -1},
		},
		{
			xs:   []float64{10, 20},
			ys:   []float64{50, 100},
			want: []float64{5},
		},
	} {
		got := calculateSlopes(test.xs, test.ys)
		if !floats.EqualApprox(got, test.want, 1e-14) {
			t.Errorf("Mismatch in calculated slopes in case %d: got %v, want %v", i, got, test.want)
		}
	}
}

func applyFunc(xs []float64, f func(x float64) float64) []float64 {
	ys := make([]float64, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func panics(fun func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	fun()
	return
}

func discrDerivPredict(p Predictor, x0, x1, x, h float64) float64 {
	if x <= x0+h {
		return (p.Predict(x+h) - p.Predict(x)) / h
	} else if x >= x1-h {
		return (p.Predict(x) - p.Predict(x-h)) / h
	} else {
		return (p.Predict(x+h) - p.Predict(x-h)) / (2 * h)
	}
}
