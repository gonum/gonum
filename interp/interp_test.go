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

func TestPiecewiseCubic(t *testing.T) {
	t.Parallel()
	const (
		h        = 1e-8
		valueTol = 1e-14
		derivTol = 1e-6
	)
	for i, test := range []struct {
		xs []float64
		f  func(float64) float64
		df func(float64) float64
	}{
		{
			xs: []float64{-1.001, 0.2, 2},
			f:  func(x float64) float64 { return x * x },
			df: func(x float64) float64 { return 2 * x },
		},
		{
			xs: []float64{-1.001, 0.2, 2},
			f:  func(x float64) float64 { return x * x },
			df: func(x float64) float64 { return 2 * x },
		},
		{
			xs: []float64{-1.001, 0.2, 10},
			f:  func(x float64) float64 { return 1.5*x - 1 },
			df: func(x float64) float64 { return 1.5 },
		},
		{
			xs: []float64{-1.001, 0.2, 10},
			f:  func(x float64) float64 { return -1 },
			df: func(x float64) float64 { return 0 },
		},
		{
			xs: []float64{-1.1, 0.2, 0.99, 2.5, 2.99},
			f:  math.Cos,
			df: math.Sin,
		},
		{
			xs: []float64{-1.1, 0.2, 0.99, 2.5, 2.99},
			f:  math.Exp,
			df: math.Exp,
		},
		{
			xs: []float64{-1.1, 0.2, 0.99, 2.5, 2.99},
			f:  func(x float64) float64 { return math.Sin(x * x) },
			df: func(x float64) float64 { return -2 * x * math.Cos(x*x) },
		},
	} {
		ys := applyFunc(test.xs, test.f)
		dydxs := applyFunc(test.xs, test.df)
		var pc PiecewiseCubic
		pc.fitWithDerivatives(test.xs, ys, dydxs)
		n := len(test.xs)
		for j := 0; j < n; j++ {
			x := test.xs[j]
			got := pc.Predict(x)
			want := test.f(x)
			if math.Abs(got-want) > valueTol {
				t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
			if j < n-1 {
				got = pc.coeffs.At(j, 0)
				if math.Abs(got-want) > valueTol {
					t.Errorf("Mismatch in 0-th order interpolation coefficient in %d-th node for test case %d: got %v, want %g", j, i, got, want)
				}
			} else {
				got = pc.lastY
				if math.Abs(got-want) > valueTol {
					t.Errorf("Mismatch in lastY for test case %d: got %v, want %g", i, got, want)
				}
			}

			if j > 0 {
				dx := test.xs[j] - test.xs[j-1]
				got = ((pc.coeffs.At(j-1, 3)*dx+pc.coeffs.At(j-1, 2))*dx+pc.coeffs.At(j-1, 1))*dx + pc.coeffs.At(j-1, 0)
				if math.Abs(got-want) > valueTol {
					t.Errorf("Interpolation coefficients in %d-th node produce mismatch in interpolated value at %g for test case %d: got %v, want %g", j-1, x, i, got, want)
				}
			}
			if j == 0 {
				got = (pc.Predict(x+h) - pc.Predict(x)) / h
			} else if j == n-1 {
				got = (pc.Predict(x) - pc.Predict(x-h)) / h
			} else {
				got = (pc.Predict(x+h) - pc.Predict(x-h)) / (2 * h)
			}
			want = test.df(x)
			if math.Abs(got-want) > derivTol {
				t.Errorf("Mismatch in interpolated derivative value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
			if j < n-1 {
				got = pc.coeffs.At(j, 1)
				if math.Abs(got-want) > valueTol {
					t.Errorf("Mismatch in 1-st order interpolation coefficient in %d-th node for test case %d: got %v, want %g", j, i, got, want)
				}
			}
			if j > 0 {
				dx := test.xs[j] - test.xs[j-1]
				got = (3*pc.coeffs.At(j-1, 3)*dx+2*pc.coeffs.At(j-1, 2))*dx + pc.coeffs.At(j-1, 1)
				if math.Abs(got-want) > valueTol {
					t.Errorf("Interpolation coefficients in %d-th node produce mismatch in interpolated derivative value at %g for test case %d: got %v, want %g", j-1, x, i, got, want)
				}
			}
		}
	}
}

func TestPiecewiseCubicExactFit(t *testing.T) {
	t.Parallel()
	xs := []float64{-1, 0, 1}
	ys := make([]float64, 3)
	dydxs := make([]float64, 3)
	leftPoly := func(x float64) float64 {
		return x*x - x + 1
	}
	leftPolyDerivative := func(x float64) float64 {
		return 2*x - 1
	}
	rightPoly := func(x float64) float64 {
		return x*x*x - x + 1
	}
	rightPolyDerivative := func(x float64) float64 {
		return 3*x*x - 1
	}
	ys[0] = leftPoly(xs[0])
	ys[1] = leftPoly(xs[1])
	ys[2] = rightPoly(xs[2])
	dydxs[0] = leftPolyDerivative(xs[0])
	dydxs[1] = leftPolyDerivative(xs[1])
	dydxs[2] = rightPolyDerivative(xs[2])
	var pc PiecewiseCubic
	pc.fitWithDerivatives(xs, ys, dydxs)
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

func TestPiecewiseCubicFitWithDerivativesErrors(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		xs, ys, dydxs []float64
	}{
		{
			xs:    []float64{0, 1, 2},
			ys:    []float64{10, 20},
			dydxs: []float64{0, 0, 0},
		},
		{
			xs:    []float64{0, 1, 1},
			ys:    []float64{10, 20, 30},
			dydxs: []float64{0, 0, 0, 0},
		},
		{
			xs:    []float64{0},
			ys:    []float64{0},
			dydxs: []float64{0},
		},
		{
			xs:    []float64{0, 1, 1},
			ys:    []float64{10, 20, 10},
			dydxs: []float64{0, 0, 0},
		},
	} {
		var pc PiecewiseCubic
		if !panics(func() { pc.fitWithDerivatives(test.xs, test.ys, test.dydxs) }) {
			t.Errorf("expected panick for xs: %v, ys: %v and dydxs: %v", test.xs, test.ys, test.dydxs)
		}
	}
}

func TestAkimaSpline(t *testing.T) {
	t.Parallel()
	const (
		nPts      = 40
		wiggleTol = 1e-1
	)
	for i, test := range []struct {
		xs  []float64
		f   func(float64) float64
		tol float64
	}{
		{
			xs:  []float64{-1, 0, 1},
			f:   func(x float64) float64 { return x * x },
			tol: 1e-14,
		},
		{
			xs:  []float64{-1, 1},
			f:   func(x float64) float64 { return -10 * x },
			tol: 1e-14,
		},
		{
			xs:  []float64{-1, 0, 1},
			f:   func(x float64) float64 { return -10 * x },
			tol: 1e-14,
		},
		{
			xs:  []float64{-0.2, -0.1, 0, 0.1, 0.2},
			f:   math.Cos,
			tol: 1e-4,
		},
		{
			xs:  []float64{-0.2, -0.1, 0, 0.1, 0.2},
			f:   math.Sin,
			tol: 1e-4,
		},
		{
			xs:  []float64{-0.2, -0.1, 0, 0.1, 0.2},
			f:   math.Exp,
			tol: 1e-4,
		},
		{
			xs:  []float64{-0.2, -0.1, 0, 0.1, 0.2},
			f:   func(x float64) float64 { return 1 / (1 + math.Exp(-100*x)) },
			tol: 0.5,
		},
	} {
		var as AkimaSpline
		ys := applyFunc(test.xs, test.f)
		err := as.Fit(test.xs, ys)
		if err != nil {
			t.Errorf("Error when fitting AkimaSpline in test case %d: %v", i, err)
		}
		x0 := test.xs[0]
		x1 := test.xs[len(test.xs)-1]
		dx := (x1 - x0) / nPts
		for j := -1; j <= nPts+1; j++ {
			x := x0 + float64(j)*dx
			got := as.Predict(x)
			want := test.f(math.Min(x1, math.Max(x0, x)))
			if math.Abs(got-want) > test.tol {
				t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
		}
		n := len(test.xs)
		for j := 0; j < n; j++ {
			x := test.xs[j]
			got := as.Predict(x)
			want := test.f(x)
			if math.Abs(got-want) > 1e-14 {
				t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
		}
		m := n - 1
		for j := 0; j < m; j++ {
			x0 := test.xs[j]
			x1 := test.xs[j+1]
			yMin := math.Min(ys[j], ys[j+1])
			yMax := math.Max(ys[j], ys[j+1])
			dx := (x1 - x0) / nPts
			for k := 0; k <= nPts; k++ {
				x := x0 + float64(k)*dx
				y := as.Predict(x)
				if y < yMin-wiggleTol || y > yMax+wiggleTol {
					t.Errorf("Interpolated values show large wiggles at x == %g for test case %d: y == %g outside (%g, %g) more than %g", x, i, y, yMin, yMax, wiggleTol)
				}
			}
		}
	}
}

func TestAkimaSplineFitErrors(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		xs, ys []float64
	}{
		{
			xs: []float64{0, 1, 2},
			ys: []float64{10, 20},
		},
		{
			xs: []float64{0},
			ys: []float64{0},
		},
		{
			xs: []float64{0, 1, 1},
			ys: []float64{10, 20, 10},
		},
	} {
		var as AkimaSpline
		if !panics(func() { as.Fit(test.xs, test.ys) }) {
			t.Errorf("expected panick for xs: %v and ys: %v", test.xs, test.ys)
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
