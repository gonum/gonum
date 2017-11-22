// Copyright 2017 The Gonum Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fit provides functions to fit data.
package fit // import "gonum.org/v1/gonum/fit"

import (
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
)

// Func1D describes a 1D function to fit some data.
type Func1D struct {
	// F is the function to minimize.
	// ps is the slice of parameters to optimize during the fit.
	F func(x float64, ps []float64) float64

	// N is the number of parameters to optimize during the fit.
	// If N is 0, Ps must not be nil.
	N int

	// Ps is the initial values for the parameters.
	// If Ps is nil, the set of initial parameters values is a slice of
	// length N filled with zeros.
	Ps []float64

	X   []float64
	Y   []float64
	Err []float64

	sig2 []float64 // inverse of squares of measurement errors along Y.

	fct  func(ps []float64) float64 // cost function (objective function)
	grad func(grad, ps []float64)
	hess func(hess mat.MutableSymmetric, x []float64)
}

func (f *Func1D) init() {

	f.sig2 = make([]float64, len(f.Y))
	switch {
	default:
		for i := range f.Y {
			f.sig2[i] = 1
		}
	case f.Err != nil:
		for i, v := range f.Err {
			f.sig2[i] = 1 / (v * v)
		}
	}

	if f.Ps == nil {
		f.Ps = make([]float64, f.N)
	}

	if len(f.Ps) == 0 {
		panic("fit: invalid number of initial parameters")
	}

	if len(f.X) != len(f.Y) {
		panic("fit: mismatch length")
	}

	if len(f.sig2) != len(f.Y) {
		panic("fit: mismatch length")
	}

	f.fct = func(ps []float64) float64 {
		var chi2 float64
		for i := range f.X {
			res := f.F(f.X[i], ps) - f.Y[i]
			chi2 += res * res * f.sig2[i]
		}
		return 0.5 * chi2
	}

	f.grad = func(grad, ps []float64) {
		fd.Gradient(grad, f.fct, ps, nil)
	}

	f.hess = func(hess mat.MutableSymmetric, x []float64) {
		fd.Hessian(hess.(*mat.SymDense), f.fct, x, nil)
	}
}
