// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
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
		{[]float64{0, 1, 2}, []float64{-0.5, 1.5}, "input slices have different lengths"},
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

func TestPiecewiseCubicPredict(t *testing.T) {
	t.Parallel()
	xs := []float64{-1, 0, 1}
	lastY := rightPoly(xs[2])
	coeffs := mat.NewDense(2, 4, []float64{3, -3, 1, 0, 1, -1, 0, 1})
	pc := PiecewiseCubic{xs, coeffs, lastY}
	testFittedPolys(t, &pc)
}

func TestPiecewiseCubicFitWithDerivatives(t *testing.T) {
	t.Parallel()
	xs := []float64{-1, 0, 1}
	ys := make([]float64, 3)
	dydxs := make([]float64, 3)
	ys[0] = leftPoly(xs[0])
	ys[1] = leftPoly(xs[1])
	ys[2] = rightPoly(xs[2])
	dydxs[0] = leftPolyDerivative(xs[0])
	dydxs[1] = leftPolyDerivative(xs[1])
	dydxs[2] = rightPolyDerivative(xs[2])
	var pc PiecewiseCubic
	err := pc.FitWithDerivatives(xs, ys, dydxs)
	if err != nil {
		t.Errorf("Error when fitting piecewise cubic interpolator: %v", err)
	}
	testFittedPolys(t, &pc)
	lastY := rightPoly(xs[2])
	if pc.lastY != lastY {
		t.Errorf("Mismatch in lastY: got %v, want %g", pc.lastY, lastY)
	}
	if !floats.Equal(pc.xs, xs) {
		t.Errorf("Mismatch in xs: got %v, want %v", pc.xs, xs)
	}
	coeffs := mat.NewDense(2, 4, []float64{3, -3, 1, 0, 1, -1, 0, 1})
	if !mat.EqualApprox(pc.coeffs, coeffs, 1e-14) {
		t.Errorf("Mismatch in coeffs: got %v, want %v", pc.coeffs, coeffs)
	}
}

func TestPiecewiseCubicFitWithDerivativesErrors(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		xs, ys, dydxs   []float64
		expectedMessage string
	}{
		{
			xs:              []float64{0, 1, 2},
			ys:              []float64{10, 20},
			dydxs:           []float64{0, 0, 0},
			expectedMessage: differentLengths,
		},
		{
			xs:              []float64{0, 1, 1},
			ys:              []float64{10, 20, 30},
			dydxs:           []float64{0, 0, 0, 0},
			expectedMessage: differentLengths,
		},
		{
			xs:              []float64{0},
			ys:              []float64{0},
			dydxs:           []float64{0},
			expectedMessage: tooFewPoints,
		},
		{
			xs:              []float64{0, 1, 1},
			ys:              []float64{10, 20, 10},
			dydxs:           []float64{0, 0, 0},
			expectedMessage: xsNotStrictlyIncreasing,
		},
	} {
		var pc PiecewiseCubic
		err := pc.FitWithDerivatives(test.xs, test.ys, test.dydxs)
		if err == nil || err.Error() != test.expectedMessage {
			t.Errorf("expected error for xs: %v, ys: %v and dydxs: %v with message: %s", test.xs, test.ys, test.dydxs, test.expectedMessage)
		}
	}
}

func leftPoly(x float64) float64 {
	return x*x - x + 1
}

func leftPolyDerivative(x float64) float64 {
	return 2*x - 1
}

func rightPoly(x float64) float64 {
	return x*x*x - x + 1
}

func rightPolyDerivative(x float64) float64 {
	return 3*x*x - 1
}

func testFittedPolys(t *testing.T, pc *PiecewiseCubic) {
	lastY := rightPoly(1)
	for i, test := range []struct {
		x    float64
		want float64
	}{
		{-2, 3},
		{-1, 3},
		{-0.9, leftPoly(-0.9)},
		{-0.75, leftPoly(-0.75)},
		{-0.5, leftPoly(-0.5)},
		{-0.25, leftPoly(-0.25)},
		{-0.1, leftPoly(-0.1)},
		{0, 1},
		{0.1, rightPoly(0.1)},
		{0.25, rightPoly(0.25)},
		{0.5, rightPoly(0.5)},
		{0.75, rightPoly(0.75)},
		{0.9, rightPoly(0.9)},
		{1, lastY},
		{2, lastY},
	} {
		got := pc.Predict(test.x)
		if math.Abs(got-test.want) > 1e-14 {
			t.Errorf("Mismatch in test case %d for x = %g: got %v, want %g", i, test.x, got, test.want)
		}
	}
}

func TestAkimaSplinesSingleFunction(t *testing.T) {
	t.Parallel()
	const (
		nPts = 40
		tol  = 1e-14
	)
	for i, test := range []struct {
		xs, ys []float64
		f      func(float64) float64
	}{
		{
			xs: []float64{-1, 0, 1},
			ys: []float64{1, 0, 1},
			f:  func(x float64) float64 { return x * x },
		},
		{
			xs: []float64{-1, 1},
			ys: []float64{10, -10},
			f:  func(x float64) float64 { return -10 * x },
		},
		{
			xs: []float64{-1, 0, 1},
			ys: []float64{10, 0, -10},
			f:  func(x float64) float64 { return -10 * x },
		},
	} {
		var as AkimaSplines
		err := as.Fit(test.xs, test.ys)
		if err != nil {
			t.Errorf("Error when fitting AkimaSplines in test case %d: %v", i, err)
		}
		x0 := test.xs[0]
		x1 := test.xs[len(test.xs)-1]
		dx := (x1 - x0) / nPts
		for j := -1; j <= nPts+1; j++ {
			x := x0 + float64(j)*dx
			got := as.Predict(x)
			want := test.f(math.Min(x1, math.Max(x0, x)))
			if math.Abs(got-want) > tol {
				t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
		}
	}
}

func TestAkimaSplinesNoWiggles(t *testing.T) {
	const nPts = 40
	const wiggleTol = 1e-1
	xs := []float64{0, 1, 2, 3, 4, 5}
	ys := []float64{-2, 0.6, 1.4, -3.8, -4.2, -3.5}
	var as AkimaSplines
	err := as.Fit(xs, ys)
	if err != nil {
		t.Errorf("Error when fitting AkimaSplines: %v", err)
	}
	m := len(xs) - 1
	for i := 0; i < m; i++ {
		x0 := xs[i]
		x1 := xs[i+1]
		yMin := math.Min(ys[i], ys[i+1])
		yMax := math.Max(ys[i], ys[i+1])
		dx := (x1 - x0) / nPts
		for j := 0; j <= nPts; j++ {
			x := x0 + float64(j)*dx
			y := as.Predict(x)
			if y < yMin-wiggleTol || y > yMax+wiggleTol {
				t.Errorf("Interpolated values show large wiggles for x == %g: y == %g", x, y)
			}
		}
	}
}
