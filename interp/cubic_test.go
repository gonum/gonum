// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
)

func TestPiecewiseCubic(t *testing.T) {
	t.Parallel()
	const (
		h        = 1e-8
		valueTol = 1e-13
		derivTol = 1e-6
		nPts     = 100
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
			xs: []float64{-1.2, -1.001, 0, 0.2, 2.01, 2.1},
			f:  func(x float64) float64 { return 4*math.Pow(x, 3) - 2*x*x + 10*x - 7 },
			df: func(x float64) float64 { return 12*x*x - 4*x + 10 },
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
	} {
		ys := applyFunc(test.xs, test.f)
		dydxs := applyFunc(test.xs, test.df)
		var pc PiecewiseCubic
		pc.FitWithDerivatives(test.xs, ys, dydxs)
		n := len(test.xs)
		m := n - 1
		x0 := test.xs[0]
		x1 := test.xs[m]
		x := x0 - 0.1
		got := pc.Predict(x)
		want := ys[0]
		if got != want {
			t.Errorf("Mismatch in value extrapolated to the left for test case %d: got %v, want %g", i, got, want)
		}
		got = pc.PredictDerivative(x)
		want = dydxs[0]
		if got != want {
			t.Errorf("Mismatch in derivative extrapolated to the left for test case %d: got %v, want %g", i, got, want)
		}
		x = x1 + 0.1
		got = pc.Predict(x)
		want = ys[m]
		if got != want {
			t.Errorf("Mismatch in value extrapolated to the right for test case %d: got %v, want %g", i, got, want)
		}
		got = pc.PredictDerivative(x)
		want = dydxs[m]
		if got != want {
			t.Errorf("Mismatch in derivative extrapolated to the right for test case %d: got %v, want %g", i, got, want)
		}
		for j := 0; j < n; j++ {
			x := test.xs[j]
			got := pc.Predict(x)
			want := test.f(x)
			if math.Abs(got-want) > valueTol {
				t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
			if j < m {
				got = pc.coeffs.At(j, 0)
				if math.Abs(got-want) > valueTol {
					t.Errorf("Mismatch in 0-th order interpolation coefficient in %d-th node for test case %d: got %v, want %g", j, i, got, want)
				}
				dx := (test.xs[j+1] - x) / nPts
				for k := 1; k < nPts; k++ {
					xk := x + float64(k)*dx
					got := pc.Predict(xk)
					want := test.f(xk)
					if math.Abs(got-want) > valueTol {
						t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
					}
					got = pc.PredictDerivative(xk)
					want = discrDerivPredict(&pc, x0, x1, xk, h)
					if math.Abs(got-want) > derivTol {
						t.Errorf("Mismatch in interpolated derivative at x == %g for test case %d: got %v, want %g", x, i, got, want)
					}
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
			got = discrDerivPredict(&pc, x0, x1, x, h)
			want = test.df(x)
			if math.Abs(got-want) > derivTol {
				t.Errorf("Mismatch in numerical derivative of interpolated function at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
			got = pc.PredictDerivative(x)
			if math.Abs(got-want) > valueTol {
				t.Errorf("Mismatch in interpolated derivative value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
		}
	}
}

func TestPiecewiseCubicFitWithDerivatives(t *testing.T) {
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
	pc.FitWithDerivatives(xs, ys, dydxs)
	lastY := rightPoly(xs[2])
	if pc.lastY != lastY {
		t.Errorf("Mismatch in lastY: got %v, want %g", pc.lastY, lastY)
	}
	lastDyDx := rightPolyDerivative(xs[2])
	if pc.lastDyDx != lastDyDx {
		t.Errorf("Mismatch in lastDxDy: got %v, want %g", pc.lastDyDx, lastDyDx)
	}
	if !floats.Equal(pc.xs, xs) {
		t.Errorf("Mismatch in xs: got %v, want %v", pc.xs, xs)
	}
	coeffs := mat.NewDense(2, 4, []float64{3, -3, 1, 0, 1, -1, 0, 1})
	if !mat.EqualApprox(&pc.coeffs, coeffs, 1e-14) {
		t.Errorf("Mismatch in coeffs: got %v, want %v", pc.coeffs, coeffs)
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
		if !panics(func() { pc.FitWithDerivatives(test.xs, test.ys, test.dydxs) }) {
			t.Errorf("expected panic for xs: %v, ys: %v and dydxs: %v", test.xs, test.ys, test.dydxs)
		}
	}
}

func TestAkimaSpline(t *testing.T) {
	t.Parallel()
	const (
		derivAbsTol = 1e-8
		derivRelTol = 1e-7
		h           = 1e-8
		nPts        = 100
		tol         = 1e-14
	)
	for i, test := range []struct {
		xs []float64
		f  func(float64) float64
	}{
		{
			xs: []float64{-5, -3, -2, -1.5, -1, 0.5, 1.5, 2.5, 3},
			f:  func(x float64) float64 { return x * x },
		},
		{
			xs: []float64{-5, -3, -2, -1.5, -1, 0.5, 1.5, 2.5, 3},
			f:  func(x float64) float64 { return math.Pow(x, 3.) - x*x + 2 },
		},
		{
			xs: []float64{-5, -3, -2, -1.5, -1, 0.5, 1.5, 2.5, 3},
			f:  func(x float64) float64 { return -10 * x },
		},
		{
			xs: []float64{-5, -3, -2, -1.5, -1, 0.5, 1.5, 2.5, 3},
			f:  math.Sin,
		},
		{
			xs: []float64{0, 1},
			f:  math.Exp,
		},
		{
			xs: []float64{-1, 0.5},
			f:  math.Cos,
		},
	} {
		var as AkimaSpline
		n := len(test.xs)
		m := n - 1
		x0 := test.xs[0]
		x1 := test.xs[m]
		ys := applyFunc(test.xs, test.f)
		err := as.Fit(test.xs, ys)
		if err != nil {
			t.Errorf("Error when fitting AkimaSpline in test case %d: %v", i, err)
		}
		for j := 0; j < n; j++ {
			x := test.xs[j]
			got := as.Predict(x)
			want := test.f(x)
			if math.Abs(got-want) > tol {
				t.Errorf("Mismatch in interpolated value at x == %g for test case %d: got %v, want %g", x, i, got, want)
			}
			if j < m {
				dx := (test.xs[j+1] - x) / nPts
				for k := 1; k < nPts; k++ {
					xk := x + float64(k)*dx
					got = as.PredictDerivative(xk)
					want = discrDerivPredict(&as, x0, x1, xk, h)
					if math.Abs(got-want) > derivRelTol*math.Abs(want)+derivAbsTol {
						t.Errorf("Mismatch in interpolated derivative at x == %g for test case %d: got %v, want %g", x, i, got, want)
					}
				}
			}
		}
		if n == 2 {
			got := as.cubic.coeffs.At(0, 1)
			want := (ys[1] - ys[0]) / (test.xs[1] - test.xs[0])
			if math.Abs(got-want) > tol {
				t.Errorf("Mismatch in approximated slope for length-2 test case %d: got %v, want %g", i, got, want)
			}
			for j := 2; i < 4; j++ {
				got := as.cubic.coeffs.At(0, j)
				if got != 0 {
					t.Errorf("Non-zero order-%d coefficient for length-2 test case %d: got %v", j, i, got)
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
			xs: []float64{0, 1},
			ys: []float64{10, 20, 30},
		},
		{
			xs: []float64{0},
			ys: []float64{0},
		},
		{
			xs: []float64{0, 1, 1},
			ys: []float64{10, 20, 10},
		},
		{
			xs: []float64{0, 2, 1},
			ys: []float64{10, 20, 10},
		},
		{
			xs: []float64{0, 0},
			ys: []float64{-1, 2},
		},
		{
			xs: []float64{0, -1},
			ys: []float64{-1, 2},
		},
	} {
		var as AkimaSpline
		if !panics(func() { _ = as.Fit(test.xs, test.ys) }) {
			t.Errorf("expected panic for xs: %v and ys: %v", test.xs, test.ys)
		}
	}
}

func TestAkimaWeightedAverage(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		v1, v2, w1, w2, want float64
		// "want" values calculated by hand.
	}{
		{
			v1:   -1,
			v2:   1,
			w1:   0,
			w2:   0,
			want: 0,
		},
		{
			v1:   -1,
			v2:   1,
			w1:   1e6,
			w2:   1e6,
			want: 0,
		},
		{
			v1:   -1,
			v2:   1,
			w1:   1e-10,
			w2:   0,
			want: -1,
		},
		{
			v1:   -1,
			v2:   1,
			w1:   0,
			w2:   1e-10,
			want: 1,
		},
		{
			v1:   0,
			v2:   1000,
			w1:   1e-13,
			w2:   3e-13,
			want: 750,
		},
		{
			v1:   0,
			v2:   1000,
			w1:   3e-13,
			w2:   1e-13,
			want: 250,
		},
	} {
		got := akimaWeightedAverage(test.v1, test.v2, test.w1, test.w2)
		if !scalar.EqualWithinAbsOrRel(got, test.want, 1e-14, 1e-14) {
			t.Errorf("Mismatch in test case %d: got %v, want %g", i, got, test.want)
		}
	}
}

func TestAkimaSlopes(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		xs, ys, want []float64
		// "want" values calculated by hand.
	}{
		{
			xs:   []float64{-2, 0, 1},
			ys:   []float64{2, 0, 1.5},
			want: []float64{-6, -3.5, -1, 1.5, 4, 6.5},
		},
		{
			xs:   []float64{-2, -0.5, 1},
			ys:   []float64{-2, -0.5, 1},
			want: []float64{1, 1, 1, 1, 1, 1},
		},
		{
			xs:   []float64{-2, -0.5, 1},
			ys:   []float64{1, 1, 1},
			want: []float64{0, 0, 0, 0, 0, 0},
		},
		{
			xs:   []float64{0, 1.5, 2, 4, 4.5, 5, 6, 7.5, 8},
			ys:   []float64{-5, -4, -3.5, -3.25, -3.25, -2.5, -1.5, -1, 2},
			want: []float64{0, 1. / 3, 2. / 3, 1, 0.125, 0, 1.5, 1, 1. / 3, 6, 12 - 1./3, 18 - 2./3},
		},
	} {
		got := akimaSlopes(test.xs, test.ys)
		if !floats.EqualApprox(got, test.want, 1e-14) {
			t.Errorf("Mismatch in test case %d: got %v, want %v", i, got, test.want)
		}
	}
}

func TestAkimaSlopesErrors(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		xs, ys []float64
	}{
		{
			xs: []float64{0, 1, 2},
			ys: []float64{10, 20},
		},
		{
			xs: []float64{0, 1},
			ys: []float64{10, 20, 30},
		},
		{
			xs: []float64{0, 2},
			ys: []float64{0, 1},
		},
		{
			xs: []float64{0, 1, 1},
			ys: []float64{10, 20, 10},
		},
		{
			xs: []float64{0, 2, 1},
			ys: []float64{10, 20, 10},
		},
		{
			xs: []float64{0, 0},
			ys: []float64{-1, 2},
		},
		{
			xs: []float64{0, -1},
			ys: []float64{-1, 2},
		},
	} {
		if !panics(func() { akimaSlopes(test.xs, test.ys) }) {
			t.Errorf("expected panic for xs: %v and ys: %v", test.xs, test.ys)
		}
	}
}

func TestAkimaWeights(t *testing.T) {
	t.Parallel()
	const tol = 1e-14
	slopes := []float64{-2, -1, -0.1, 0.2, 1.2, 2.5}
	// "want" values calculated by hand.
	want := [][]float64{
		{0.3, 1},
		{1, 0.9},
		{1.3, 0.3},
	}
	for i := 0; i < len(want); i++ {
		gotLeft, gotRight := akimaWeights(slopes, i)
		if math.Abs(gotLeft-want[i][0]) > tol {
			t.Errorf("Mismatch in left weight for node %d: got %v, want %g", i, gotLeft, want[i][0])
		}
		if math.Abs(gotRight-want[i][1]) > tol {
			t.Errorf("Mismatch in left weight for node %d: got %v, want %g", i, gotRight, want[i][1])
		}
	}
}

func TestFritschButland(t *testing.T) {
	t.Parallel()
	const (
		tol  = 1e-14
		nPts = 100
	)
	for k, test := range []struct {
		xs, ys []float64
	}{
		{
			xs: []float64{0, 2},
			ys: []float64{0, 0.5},
		},
		{
			xs: []float64{0, 2},
			ys: []float64{0, -0.5},
		},
		{
			xs: []float64{0, 2},
			ys: []float64{0, 0},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{0, 1, 2, 2.5},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{0, 1.5, 1.5, 2.5},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{0, 1.5, 1.5, 1},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{0, 2.5, 1.5, 1},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{0, 2.5, 1.5, 2},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{4, 3, 2, 1},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{4, 3, 2, 2},
		},
		{
			xs: []float64{0, 2, 3, 4},
			ys: []float64{4, 3, 2, 5},
		},
		{
			xs: []float64{0, 2, 3, 4, 5, 6},
			ys: []float64{0, 1, 0.5, 0.5, 1.5, 1.5},
		},
		{
			xs: []float64{0, 2, 3, 4, 5, 6},
			ys: []float64{0, 1, 1.5, 2.5, 1.5, 1},
		},
		{
			xs: []float64{0, 2, 3, 4, 5, 6},
			ys: []float64{0, -1, -1.5, -2.5, -1.5, -1},
		},
		{
			xs: []float64{0, 2, 3, 4, 5, 6},
			ys: []float64{0, 1, 0.5, 1.5, 1, 2},
		},
		{
			xs: []float64{0, 2, 3, 4, 5, 6},
			ys: []float64{0, 1, 1.5, 2.5, 3, 4},
		},
		{
			xs: []float64{0, 2, 3, 4, 5, 6},
			ys: []float64{0, 0.0001, -1.5, -2.5, -0.0001, 0},
		},
	} {
		var fb FritschButland
		err := fb.Fit(test.xs, test.ys)
		if err != nil {
			t.Errorf("Error when fitting FritschButland in test case %d: %v", k, err)
		}
		n := len(test.xs)
		for i := 0; i < n; i++ {
			got := fb.Predict(test.xs[i])
			want := test.ys[i]
			if got != want {
				t.Errorf("Mismatch in interpolated value for node %d in test case %d: got %v, want %g", i, k, got, want)
			}
		}
		if n == 2 {
			h := test.xs[1] - test.xs[0]
			want := (test.ys[1] - test.ys[0]) / h
			for i := 0; i < 2; i++ {
				got := fb.PredictDerivative(test.xs[i])
				if !scalar.EqualWithinAbs(got, want, tol) {
					t.Errorf("Mismatch in approximated derivative for node %d in 2-node test case %d: got %v, want %g", i, k, got, want)
				}
			}
			dx := h / (nPts + 1)
			for i := 1; i < nPts; i++ {
				x := test.xs[0] + float64(i)*dx
				got := fb.PredictDerivative(x)
				if !scalar.EqualWithinAbs(got, want, tol) {
					t.Errorf("Mismatch in interpolated derivative for x == %g in 2-node test case %d: got %v, want %g", x, k, got, want)
				}
			}
		} else {
			m := n - 1
			for i := 1; i < m; i++ {
				got := fb.PredictDerivative(test.xs[i])
				slope := (test.ys[i+1] - test.ys[i]) / (test.xs[i+1] - test.xs[i])
				prevSlope := (test.ys[i] - test.ys[i-1]) / (test.xs[i] - test.xs[i-1])
				if slope*prevSlope > 0 {
					if got == 0 {
						t.Errorf("Approximated derivative is zero for node %d in test case %d: %g", i, k, got)
					} else if math.Signbit(slope) != math.Signbit(got) {
						t.Errorf("Approximated derivative has wrong sign for node %d in test case %d: got %g, want %g", i, k, math.Copysign(1, got), math.Copysign(1, slope))
					}
				} else {
					if got != 0 {
						t.Errorf("Approximated derivative is not zero for node %d in test case %d: %g", i, k, got)
					}
				}
			}
			for i := 0; i < m; i++ {
				yL := test.ys[i]
				yR := test.ys[i+1]
				xL := test.xs[i]
				dx := (test.xs[i+1] - xL) / (nPts + 1)
				if yL == yR {
					for j := 1; j < nPts; j++ {
						x := xL + float64(j)*dx
						got := fb.Predict(x)
						if got != yL {
							t.Errorf("Mismatch in interpolated value for x == %g in test case %d: got %v, want %g", x, k, got, yL)
						}
						got = fb.PredictDerivative(x)
						if got != 0 {
							t.Errorf("Interpolated derivative not zero for x == %g in test case %d: got %v", x, k, got)
						}
					}
				} else {
					minY := math.Min(yL, yR)
					maxY := math.Max(yL, yR)
					for j := 1; j < nPts; j++ {
						x := xL + float64(j)*dx
						got := fb.Predict(x)
						if got < minY || got > maxY {
							t.Errorf("Interpolated value out of [%g, %g] bounds for x == %g in test case %d: got %v", minY, maxY, x, k, got)
						}
						got = fb.PredictDerivative(x)
						dy := yR - yL
						if got*dy < 0 {
							t.Errorf("Interpolated derivative has wrong sign for x == %g in test case %d: want %g, got %g", x, k, math.Copysign(1, dy), math.Copysign(1, got))
						}
					}
				}
			}
		}
	}
}

func TestFritschButlandErrors(t *testing.T) {
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
		var fb FritschButland
		if !panics(func() { _ = fb.Fit(test.xs, test.ys) }) {
			t.Errorf("expected panic for xs: %v and ys: %v", test.xs, test.ys)
		}
	}
}

func TestPiecewiseCubicFitWithSecondDerivatives(t *testing.T) {
	t.Parallel()
	const tol = 1e-14
	xs := []float64{-2, 0, 3}
	ys := []float64{2.5, 1, 2.5}
	d2ydx2s := []float64{1, 2, 3}
	var pc PiecewiseCubic
	pc.fitWithSecondDerivatives(xs, ys, d2ydx2s)
	m := len(xs) - 1
	if pc.lastY != ys[m] {
		t.Errorf("Mismatch in lastY: got %v, want %g", pc.lastY, ys[m])
	}
	if !floats.Equal(pc.xs, xs) {
		t.Errorf("Mismatch in xs: got %v, want %v", pc.xs, xs)
	}
	for i := 0; i < len(xs); i++ {
		yHat := pc.Predict(xs[i])
		if math.Abs(yHat-ys[i]) > tol {
			t.Errorf("Mismatch in predicted Y[%d]: got %v, want %g", i, yHat, ys[i])
		}
		var d2ydx2Hat float64
		if i < m {
			d2ydx2Hat = 2 * pc.coeffs.At(i, 2)
		} else {
			d2ydx2Hat = 2*pc.coeffs.At(m-1, 2) + 6*pc.coeffs.At(m-1, 3)*(xs[m]-xs[m-1])
		}
		if math.Abs(d2ydx2Hat-d2ydx2s[i]) > tol {
			t.Errorf("Mismatch in predicted d2Y/dX2[%d]: got %v, want %g", i, d2ydx2Hat, d2ydx2s[i])
		}
	}
	// Test pc.lastDyDx without copying verbatim the calculation from the tested method:
	lastDyDx := pc.PredictDerivative(xs[m] - tol/1000)
	if math.Abs(lastDyDx-pc.lastDyDx) > tol {
		t.Errorf("Mismatch in lastDxDy: got %v, want %g", pc.lastDyDx, lastDyDx)
	}
}

func TestPiecewiseCubicFitWithSecondDerivativesErrors(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		xs, ys, d2ydx2s []float64
	}{
		{
			xs:      []float64{0, 1, 2},
			ys:      []float64{10, 20},
			d2ydx2s: []float64{0, 0, 0},
		},
		{
			xs:      []float64{0, 1, 1},
			ys:      []float64{10, 20, 30},
			d2ydx2s: []float64{0, 0, 0, 0},
		},
		{
			xs:      []float64{0},
			ys:      []float64{0},
			d2ydx2s: []float64{0},
		},
		{
			xs:      []float64{0, 1, 1},
			ys:      []float64{10, 20, 10},
			d2ydx2s: []float64{0, 0, 0},
		},
	} {
		var pc PiecewiseCubic
		if !panics(func() { pc.fitWithSecondDerivatives(test.xs, test.ys, test.d2ydx2s) }) {
			t.Errorf("expected panic for xs: %v, ys: %v and d2ydx2s: %v", test.xs, test.ys, test.d2ydx2s)
		}
	}
}

func TestMakeCubicSplineSecondDerivativeEquations(t *testing.T) {
	t.Parallel()
	const tol = 1e-15
	xs := []float64{-1, 0, 2}
	ys := []float64{2, 0, 2}
	n := len(xs)
	a := mat.NewTridiag(n, nil, nil, nil)
	b := mat.NewVecDense(n, nil)
	makeCubicSplineSecondDerivativeEquations(a, b, xs, ys)
	if b.Len() != n {
		t.Errorf("Mismatch in b size: got %v, want %d", b.Len(), n)
	}
	r, c := a.Dims()
	if r != n || c != n {
		t.Errorf("Mismatch in A size: got %d x %d, want %d x %d", r, c, n, n)
	}
	expectedB := mat.NewVecDense(3, []float64{0, 3, 0})
	var diffB mat.VecDense
	diffB.SubVec(b, expectedB)
	if diffB.Norm(math.Inf(1)) > tol {
		t.Errorf("Mismatch in b values: got %v, want %v", b, expectedB)
	}
	expectedA := mat.NewDense(3, 3, []float64{0, 0, 0, 1 / 6., 1, 2 / 6., 0, 0, 0})
	var diffA mat.Dense
	diffA.Sub(a, expectedA)
	if diffA.Norm(math.Inf(1)) > tol {
		t.Errorf("Mismatch in A values: got %v, want %v", a, expectedA)
	}
}

func TestNaturalCubicFit(t *testing.T) {
	t.Parallel()
	xs := []float64{-1, 0, 2, 3.5}
	ys := []float64{2, 0, 2, 1.5}
	var nc NaturalCubic
	err := nc.Fit(xs, ys)
	if err != nil {
		t.Errorf("Error when fitting NaturalCubic: %v", err)
	}
	testXs := []float64{-1, -0.99, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3, 3.49, 3.5}
	// From scipy.interpolate.CubicSpline:
	want := []float64{
		2.0,
		1.9737725526315788,
		0.7664473684210527,
		0.0,
		0.027960526315789477,
		0.6184210526315789,
		1.3996710526315788,
		2.0,
		2.1403508771929824,
		1.9122807017543857,
		1.508859403508772,
		1.5}
	for i := 0; i < len(testXs); i++ {
		got := nc.Predict(testXs[i])
		if math.Abs(got-want[i]) > 1e-14 {
			t.Errorf("Mismatch in predicted Y value for x = %g: got %v, want %g", testXs[i], got, want[i])
		}
		got = nc.PredictDerivative(testXs[i])
		want := discrDerivPredict(&nc, xs[0], xs[len(xs)-1], testXs[i], 1e-9)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("Mismatch in predicted dY/dX value for x = %g: got %v, want %g", testXs[i], got, want)
		}
	}
}

func TestClampedCubicFit(t *testing.T) {
	t.Parallel()
	xs := []float64{-1, 0, 2, 3.5}
	ys := []float64{2, 0, 2, 1.5}
	var cc ClampedCubic
	err := cc.Fit(xs, ys)
	if err != nil {
		t.Errorf("Error when fitting ClampedCubic: %v", err)
	}
	testXs := []float64{-1, -0.99, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3, 3.49, 3.5}
	// From scipy.interpolate.CubicSpline:
	want := []float64{
		2.0,
		1.999564111111111,
		1.2021604938271606,
		0.0,
		-0.20833333333333343,
		0.41975308641975284,
		1.337962962962962,
		2.0,
		2.026748971193416,
		1.7078189300411522,
		1.5001129711934151,
		1.4999999999999996,
	}
	for i := 0; i < len(testXs); i++ {
		got := cc.Predict(testXs[i])
		if math.Abs(got-want[i]) > 1e-14 {
			t.Errorf("Mismatch in predicted Y value for x = %g: got %v, want %g", testXs[i], got, want[i])
		}
		got = cc.PredictDerivative(testXs[i])
		want := discrDerivPredict(&cc, xs[0], xs[len(xs)-1], testXs[i], 1e-9)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("Mismatch in predicted dY/dX value for x = %g: got %v, want %g", testXs[i], got, want)
		}
	}
}

func TestNotAKnotCubicFit(t *testing.T) {
	t.Parallel()
	xs := []float64{-1, 0, 2, 3.5}
	ys := []float64{2, 0, 2, 1.5}
	var nak NotAKnotCubic
	err := nak.Fit(xs, ys)
	if err != nil {
		t.Errorf("Error when fitting NotAKnotCubic: %v", err)
	}
	testXs := []float64{-1, -0.99, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3, 3.49, 3.5}
	// From scipy.interpolate.CubicSpline:
	want := []float64{
		2.0,
		1.961016095238095,
		0.5582010582010581,
		0.0,
		0.09523809523809526,
		0.6137566137566137,
		1.3253968253968258,
		2.0,
		2.407407407407407,
		2.3174603174603177,
		1.5249675026455023,
		1.5000000000000004,
	}
	for i := 0; i < len(testXs); i++ {
		got := nak.Predict(testXs[i])
		if math.Abs(got-want[i]) > 1e-14 {
			t.Errorf("Mismatch in predicted Y value for x = %g: got %v, want %g", testXs[i], got, want[i])
		}
		got = nak.PredictDerivative(testXs[i])
		want := discrDerivPredict(&nak, xs[0], xs[len(xs)-1], testXs[i], 1e-9)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("Mismatch in predicted dY/dX value for x = %g: got %v, want %g", testXs[i], got, want)
		}
	}
}
