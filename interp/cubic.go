// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// PiecewiseCubic is a piecewise cubic 1-dimensional interpolator with
// continuous value and first derivative.
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

	// Last interpolated dY/dX value, corresponding to xs[len(xs) - 1].
	lastDyDx float64
}

// Predict returns the interpolation value at x.
func (pc *PiecewiseCubic) Predict(x float64) float64 {
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

// PredictDerivative returns the predicted derivative at x.
func (pc *PiecewiseCubic) PredictDerivative(x float64) float64 {
	i := findSegment(pc.xs, x)
	if i < 0 {
		return pc.coeffs.At(0, 1)
	}
	m := len(pc.xs) - 1
	if x == pc.xs[i] {
		if i < m {
			return pc.coeffs.At(i, 1)
		}
		return pc.lastDyDx
	}
	if i == m {
		return pc.lastDyDx
	}
	dx := x - pc.xs[i]
	a := pc.coeffs.RawRowView(i)
	return (3*a[3]*dx+2*a[2])*dx + a[1]
}

// FitWithDerivatives fits a piecewise cubic predictor to (X, Y, dY/dX) value
// triples provided as three slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing,
// len(xs) != len(ys) or len(xs) != len(dydxs).
func (pc *PiecewiseCubic) FitWithDerivatives(xs, ys, dydxs []float64) {
	n := len(xs)
	if len(ys) != n {
		panic(differentLengths)
	}
	if len(dydxs) != n {
		panic(differentLengths)
	}
	if n < 2 {
		panic(tooFewPoints)
	}
	m := n - 1
	pc.coeffs.Reset()
	pc.coeffs.ReuseAs(m, 4)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			panic(xsNotStrictlyIncreasing)
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
	pc.xs = append(pc.xs[:0], xs...)
	pc.lastY = ys[m]
	pc.lastDyDx = dydxs[m]
}

// AkimaSpline is a piecewise cubic 1-dimensional interpolator with
// continuous value and first derivative, which can be fitted to (X, Y)
// value pairs without providing derivatives.
// See https://www.iue.tuwien.ac.at/phd/rottinger/node60.html for more details.
type AkimaSpline struct {
	cubic PiecewiseCubic
}

// Predict returns the interpolation value at x.
func (as *AkimaSpline) Predict(x float64) float64 {
	return as.cubic.Predict(x)
}

// PredictDerivative returns the predicted derivative at x.
func (as *AkimaSpline) PredictDerivative(x float64) float64 {
	return as.cubic.PredictDerivative(x)
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing
// or len(xs) != len(ys). Always returns nil.
func (as *AkimaSpline) Fit(xs, ys []float64) error {
	n := len(xs)
	if len(ys) != n {
		panic(differentLengths)
	}
	dydxs := make([]float64, n)

	if n == 2 {
		dx := xs[1] - xs[0]
		slope := (ys[1] - ys[0]) / dx
		dydxs[0] = slope
		dydxs[1] = slope
		as.cubic.FitWithDerivatives(xs, ys, dydxs)
		return nil
	}
	slopes := akimaSlopes(xs, ys)
	for i := 0; i < n; i++ {
		wLeft, wRight := akimaWeights(slopes, i)
		dydxs[i] = akimaWeightedAverage(slopes[i+1], slopes[i+2], wLeft, wRight)
	}
	as.cubic.FitWithDerivatives(xs, ys, dydxs)
	return nil
}

// akimaSlopes returns slopes for Akima spline method, including the approximations
// of slopes outside the data range (two on each side).
// It panics if len(xs) <= 2, elements of xs are not strictly increasing
// or len(xs) != len(ys).
func akimaSlopes(xs, ys []float64) []float64 {
	n := len(xs)
	if n <= 2 {
		panic(tooFewPoints)
	}
	if len(ys) != n {
		panic(differentLengths)
	}
	m := n + 3
	slopes := make([]float64, m)
	for i := 2; i < m-2; i++ {
		dx := xs[i-1] - xs[i-2]
		if dx <= 0 {
			panic(xsNotStrictlyIncreasing)
		}
		slopes[i] = (ys[i-1] - ys[i-2]) / dx
	}
	slopes[0] = 3*slopes[2] - 2*slopes[3]
	slopes[1] = 2*slopes[2] - slopes[3]
	slopes[m-2] = 2*slopes[m-3] - slopes[m-4]
	slopes[m-1] = 3*slopes[m-3] - 2*slopes[m-4]
	return slopes
}

// akimaWeightedAverage returns (v1 * w1 + v2 * w2) / (w1 + w2) for w1, w2 >= 0 (not checked).
// If w1 == w2 == 0, it returns a simple average of v1 and v2.
func akimaWeightedAverage(v1, v2, w1, w2 float64) float64 {
	w := w1 + w2
	if w > 0 {
		return (v1*w1 + v2*w2) / w
	}
	return 0.5*v1 + 0.5*v2
}

// akimaWeights returns the left and right weight for approximating
// the i-th derivative with neighbouring slopes.
func akimaWeights(slopes []float64, i int) (float64, float64) {
	wLeft := math.Abs(slopes[i+2] - slopes[i+3])
	wRight := math.Abs(slopes[i+1] - slopes[i])
	return wLeft, wRight
}

// FritschButland is a piecewise cubic 1-dimensional interpolator with
// continuous value and first derivative, which can be fitted to (X, Y)
// value pairs without providing derivatives.
// It is monotone, local and produces a C^1 curve. Its downside is that
// exhibits high tension, flattening out unnaturally the interpolated
// curve between the nodes.
// See Fritsch, F. N. and Butland, J., "A method for constructing local
// monotone piecewise cubic interpolants" (1984), SIAM J. Sci. Statist.
// Comput., 5(2), pp. 300-304.
type FritschButland struct {
	cubic PiecewiseCubic
}

// Predict returns the interpolation value at x.
func (fb *FritschButland) Predict(x float64) float64 {
	return fb.cubic.Predict(x)
}

// PredictDerivative returns the predicted derivative at x.
func (fb *FritschButland) PredictDerivative(x float64) float64 {
	return fb.cubic.PredictDerivative(x)
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing
// or len(xs) != len(ys). Always returns nil.
func (fb *FritschButland) Fit(xs, ys []float64) error {
	n := len(xs)
	if n < 2 {
		panic(tooFewPoints)
	}
	if len(ys) != n {
		panic(differentLengths)
	}
	dydxs := make([]float64, n)

	if n == 2 {
		dx := xs[1] - xs[0]
		slope := (ys[1] - ys[0]) / dx
		dydxs[0] = slope
		dydxs[1] = slope
		fb.cubic.FitWithDerivatives(xs, ys, dydxs)
		return nil
	}
	slopes := calculateSlopes(xs, ys)
	m := len(slopes)
	prevSlope := slopes[0]
	for i := 1; i < m; i++ {
		slope := slopes[i]
		if slope*prevSlope > 0 {
			dydxs[i] = 3 * (xs[i+1] - xs[i-1]) / ((2*xs[i+1]-xs[i-1]-xs[i])/slopes[i-1] +
				(xs[i+1]+xs[i]-2*xs[i-1])/slopes[i])
		} else {
			dydxs[i] = 0
		}
		prevSlope = slope
	}
	dydxs[0] = fritschButlandEdgeDerivative(xs, ys, slopes, true)
	dydxs[m] = fritschButlandEdgeDerivative(xs, ys, slopes, false)
	fb.cubic.FitWithDerivatives(xs, ys, dydxs)
	return nil
}

// fritschButlandEdgeDerivative calculates dy/dx approximation for the
// Fritsch-Butland method for the left or right edge node.
func fritschButlandEdgeDerivative(xs, ys, slopes []float64, leftEdge bool) float64 {
	n := len(xs)
	var dE, dI, h, hE, f float64
	if leftEdge {
		dE = slopes[0]
		dI = slopes[1]
		xE := xs[0]
		xM := xs[1]
		xI := xs[2]
		hE = xM - xE
		h = xI - xE
		f = xM + xI - 2*xE
	} else {
		dE = slopes[n-2]
		dI = slopes[n-3]
		xE := xs[n-1]
		xM := xs[n-2]
		xI := xs[n-3]
		hE = xE - xM
		h = xE - xI
		f = 2*xE - xI - xM
	}
	g := (f*dE - hE*dI) / h
	if g*dE <= 0 {
		return 0
	}
	if dE*dI <= 0 && math.Abs(g) > 3*math.Abs(dE) {
		return 3 * dE
	}
	return g
}

// fitWithSecondDerivatives fits a piecewise cubic predictor to (X, Y, d^2Y/dX^2) value
// triples provided as three slices.
// It panics if any of these is true:
// - len(xs) < 2,
// - elements of xs are not strictly increasing,
// - len(xs) != len(ys),
// - len(xs) != len(d2ydx2s).
// Note that this method does not guarantee on its own the continuity of first derivatives.
func (pc *PiecewiseCubic) fitWithSecondDerivatives(xs, ys, d2ydx2s []float64) {
	n := len(xs)
	switch {
	case len(ys) != n, len(d2ydx2s) != n:
		panic(differentLengths)
	case n < 2:
		panic(tooFewPoints)
	}
	m := n - 1
	pc.coeffs.Reset()
	pc.coeffs.ReuseAs(m, 4)
	for i := 0; i < m; i++ {
		dx := xs[i+1] - xs[i]
		if dx <= 0 {
			panic(xsNotStrictlyIncreasing)
		}
		dy := ys[i+1] - ys[i]
		dm := d2ydx2s[i+1] - d2ydx2s[i]
		pc.coeffs.Set(i, 0, ys[i])                             // a_0
		pc.coeffs.Set(i, 1, (dy-(d2ydx2s[i]+dm/3)*dx*dx/2)/dx) // a_1
		pc.coeffs.Set(i, 2, d2ydx2s[i]/2)                      // a_2
		pc.coeffs.Set(i, 3, dm/6/dx)                           // a_3
	}
	pc.xs = append(pc.xs[:0], xs...)
	pc.lastY = ys[m]
	lastDx := xs[m] - xs[m-1]
	pc.lastDyDx = pc.coeffs.At(m-1, 1) + 2*pc.coeffs.At(m-1, 2)*lastDx + 3*pc.coeffs.At(m-1, 3)*lastDx*lastDx
}

// makeCubicSplineSecondDerivativeEquations generates the basic system of linear equations
// which have to be satisfied by the second derivatives to make the first derivatives of a
// cubic spline continuous. It panics if elements of xs are not strictly increasing, or
// len(xs) != len(ys).
// makeCubicSplineSecondDerivativeEquations fills a banded matrix a and a vector b
// defining a system of linear equations a*m = b for second derivatives vector m.
// Parameters a and b are assumed to have correct dimensions and initialised to zero.
func makeCubicSplineSecondDerivativeEquations(a mat.MutableBanded, b mat.MutableVector, xs, ys []float64) {
	n := len(xs)
	if len(ys) != n {
		panic(differentLengths)
	}
	m := n - 1
	if n > 2 {
		for i := 0; i < m; i++ {
			dx := xs[i+1] - xs[i]
			if dx <= 0 {
				panic(xsNotStrictlyIncreasing)
			}
			slope := (ys[i+1] - ys[i]) / dx
			if i > 0 {
				b.SetVec(i, b.AtVec(i)+slope)
				a.SetBand(i, i, a.At(i, i)+dx/3)
				a.SetBand(i, i+1, dx/6)
			}
			if i < m-1 {
				b.SetVec(i+1, b.AtVec(i+1)-slope)
				a.SetBand(i+1, i+1, a.At(i+1, i+1)+dx/3)
				a.SetBand(i+1, i, dx/6)
			}
		}
	}
}

// NaturalCubic is a piecewise cubic 1-dimensional interpolator with
// continuous value, first and second derivatives, which can be fitted to (X, Y)
// value pairs without providing derivatives. It uses the boundary conditions
// Y′′(left end ) = Y′′(right end) = 0.
// See e.g. https://www.math.drexel.edu/~tolya/cubicspline.pdf for details.
type NaturalCubic struct {
	cubic PiecewiseCubic
}

// Predict returns the interpolation value at x.
func (nc *NaturalCubic) Predict(x float64) float64 {
	return nc.cubic.Predict(x)
}

// PredictDerivative returns the predicted derivative at x.
func (nc *NaturalCubic) PredictDerivative(x float64) float64 {
	return nc.cubic.PredictDerivative(x)
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing
// or len(xs) != len(ys). It returns an error if solving the required system
// of linear equations fails.
func (nc *NaturalCubic) Fit(xs, ys []float64) error {
	n := len(xs)
	a := mat.NewTridiag(n, nil, nil, nil)
	b := mat.NewVecDense(n, nil)
	makeCubicSplineSecondDerivativeEquations(a, b, xs, ys)
	// Add boundary conditions y′′(left) = y′′(right) = 0:
	b.SetVec(0, 0)
	b.SetVec(n-1, 0)
	a.SetBand(0, 0, 1)
	a.SetBand(n-1, n-1, 1)
	x := mat.NewVecDense(n, nil)
	err := a.SolveVecTo(x, false, b)
	if err == nil {
		nc.cubic.fitWithSecondDerivatives(xs, ys, x.RawVector().Data)
	}
	return err
}

// ClampedCubic is a piecewise cubic 1-dimensional interpolator with
// continuous value, first and second derivatives, which can be fitted to (X, Y)
// value pairs without providing derivatives. It uses the boundary conditions
// Y′(left end ) = Y′(right end) = 0.
type ClampedCubic struct {
	cubic PiecewiseCubic
}

// Predict returns the interpolation value at x.
func (cc *ClampedCubic) Predict(x float64) float64 {
	return cc.cubic.Predict(x)
}

// PredictDerivative returns the predicted derivative at x.
func (cc *ClampedCubic) PredictDerivative(x float64) float64 {
	return cc.cubic.PredictDerivative(x)
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 2, elements of xs are not strictly increasing
// or len(xs) != len(ys). It returns an error if solving the required system
// of linear equations fails.
func (cc *ClampedCubic) Fit(xs, ys []float64) error {
	n := len(xs)
	a := mat.NewTridiag(n, nil, nil, nil)
	b := mat.NewVecDense(n, nil)
	makeCubicSplineSecondDerivativeEquations(a, b, xs, ys)
	// Add boundary conditions y′′(left) = y′′(right) = 0:
	// Condition Y′(left end) = 0:
	dxL := xs[1] - xs[0]
	b.SetVec(0, (ys[1]-ys[0])/dxL)
	a.SetBand(0, 0, dxL/3)
	a.SetBand(0, 1, dxL/6)
	// Condition Y′(right end) = 0:
	m := n - 1
	dxR := xs[m] - xs[m-1]
	b.SetVec(m, (ys[m]-ys[m-1])/dxR)
	a.SetBand(m, m, -dxR/3)
	a.SetBand(m, m-1, -dxR/6)
	x := mat.NewVecDense(n, nil)
	err := a.SolveVecTo(x, false, b)
	if err == nil {
		cc.cubic.fitWithSecondDerivatives(xs, ys, x.RawVector().Data)
	}
	return err
}

// NotAKnotCubic is a piecewise cubic 1-dimensional interpolator with
// continuous value, first and second derivatives, which can be fitted to (X, Y)
// value pairs without providing derivatives. It imposes the condition that
// the third derivative of the interpolant is continuous in the first and
// last interior node.
// See http://www.cs.tau.ac.il/~turkel/notes/numeng/spline_note.pdf for details.
type NotAKnotCubic struct {
	cubic PiecewiseCubic
}

// Predict returns the interpolation value at x.
func (nak *NotAKnotCubic) Predict(x float64) float64 {
	return nak.cubic.Predict(x)
}

// PredictDerivative returns the predicted derivative at x.
func (nak *NotAKnotCubic) PredictDerivative(x float64) float64 {
	return nak.cubic.PredictDerivative(x)
}

// Fit fits a predictor to (X, Y) value pairs provided as two slices.
// It panics if len(xs) < 3 (because at least one interior node is required),
// elements of xs are not strictly increasing or len(xs) != len(ys).
// It returns an error if solving the required system of linear equations fails.
func (nak *NotAKnotCubic) Fit(xs, ys []float64) error {
	n := len(xs)
	if n < 3 {
		panic(tooFewPoints)
	}
	a := mat.NewBandDense(n, n, 2, 2, nil)
	b := mat.NewVecDense(n, nil)
	makeCubicSplineSecondDerivativeEquations(a, b, xs, ys)
	// Add boundary conditions.
	// First interior node:
	dxOuter := xs[1] - xs[0]
	dxInner := xs[2] - xs[1]
	a.SetBand(0, 0, 1/dxOuter)
	a.SetBand(0, 1, -1/dxOuter-1/dxInner)
	a.SetBand(0, 2, 1/dxInner)
	if n > 3 {
		// Last interior node:
		m := n - 1
		dxOuter = xs[m] - xs[m-1]
		dxInner = xs[m-1] - xs[m-2]
		a.SetBand(m, m, 1/dxOuter)
		a.SetBand(m, m-1, -1/dxOuter-1/dxInner)
		a.SetBand(m, m-2, 1/dxInner)
	}
	x := mat.NewVecDense(n, nil)
	err := x.SolveVec(a, b)
	if err == nil {
		nak.cubic.fitWithSecondDerivatives(xs, ys, x.RawVector().Data)
	}
	return err
}
