// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fd

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Rosenbrock struct {
	nDim int
}

func (r Rosenbrock) F(x []float64) (sum float64) {
	deriv := make([]float64, len(x))
	return r.FDf(x, deriv)
}

func (r Rosenbrock) FDf(x []float64, deriv []float64) (sum float64) {
	sum = 0

	for i := range deriv {
		deriv[i] = 0
	}

	for i := 0; i < len(x)-1; i++ {
		sum += math.Pow(1-x[i], 2) + 100*math.Pow(x[i+1]-math.Pow(x[i], 2), 2)
	}
	for i := 0; i < len(x)-1; i++ {
		deriv[i] += -1 * 2 * (1 - x[i])
		deriv[i] += 2 * 100 * (x[i+1] - math.Pow(x[i], 2)) * (-2 * x[i])
	}
	for i := 1; i < len(x); i++ {
		deriv[i] += 2 * 100 * (x[i] - math.Pow(x[i-1], 2))
	}

	return sum
}

func TestGradient(t *testing.T) {
	for i, test := range []struct {
		nDim   int
		tol    float64
		method Method
	}{
		{
			nDim:   2,
			tol:    1e-4,
			method: Forward,
		},
		{
			nDim:   2,
			tol:    1e-6,
			method: Central,
		},
		{
			nDim:   40,
			tol:    1e-4,
			method: Forward,
		},
		{
			nDim:   40,
			tol:    1e-6,
			method: Central,
		},
	} {
		x := make([]float64, test.nDim)
		for i := range x {
			x[i] = rand.Float64()
		}
		xcopy := make([]float64, len(x))
		copy(xcopy, x)

		r := Rosenbrock{len(x)}
		trueGradient := make([]float64, len(x))
		r.FDf(x, trueGradient)

		settings := DefaultSettings()
		settings.Method = test.method
		gradient := make([]float64, len(x))
		for i := range gradient {
			gradient[i] = rand.Float64()
		}

		Gradient(r.F, x, settings, gradient)
		if !floats.EqualApprox(gradient, trueGradient, test.tol) {
			t.Errorf("Case %v: gradient mismatch in serial. Want: %v, Got: %v.", i, trueGradient, gradient)
		}
		if !floats.Equal(x, xcopy) {
			t.Errorf("Case %v: x modified during call to gradient in serial")
		}

		// Try with known value
		for i := range gradient {
			gradient[i] = rand.Float64()
		}
		settings.OriginKnown = true
		settings.OriginValue = r.F(x)
		Gradient(r.F, x, settings, gradient)
		if !floats.EqualApprox(gradient, trueGradient, test.tol) {
			t.Errorf("Case %v: gradient mismatch with known origin in serial. Want: %v, Got: %v.", i, trueGradient, gradient)
		}

		// Concurrently
		for i := range gradient {
			gradient[i] = rand.Float64()
		}
		settings.Concurrent = true
		settings.OriginKnown = false
		Gradient(r.F, x, settings, gradient)
		if !floats.EqualApprox(gradient, trueGradient, test.tol) {
			t.Errorf("Case %v: gradient mismatch with unknown origin in parallel. Want: %v, Got: %v.", i, trueGradient, gradient)
		}
		if !floats.Equal(x, xcopy) {
			t.Errorf("Case %v: x modified during call to gradient in parallel")
		}

		// Concurrently with origin known
		for i := range gradient {
			gradient[i] = rand.Float64()
		}
		settings.OriginKnown = true
		Gradient(r.F, x, settings, gradient)
		if !floats.EqualApprox(gradient, trueGradient, test.tol) {
			t.Errorf("Case %v: gradient mismatch with known origin in parallel. Want: %v, Got: %v.", i, trueGradient, gradient)
		}

	}
}
