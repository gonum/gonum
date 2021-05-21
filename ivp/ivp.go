// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivp

import (
	"gonum.org/v1/gonum/mat"
)

// IVP defines an initial value problem.
type IVP interface {
	// Initial values vector for state variables x and inputs u.
	IV() (x0, u0 mat.Vector)
	// Equations returns the coupled, non-linear algebraic differential
	// equations for the state variables (x) and the functions for inputs (u).
	// The input functions (ueq) are not differential equations but rather
	// calculated directly from a given x vector and current input vector.
	// Results are stored in y which are the length of x and u, respectively.
	// The scalar (float64) argument is the domain over which the
	// problem is integrated, which is usually time for most physical problems.
	//
	// If problem has no input functions then u supplied and ufunc returned
	// may be nil. x equations my not be nil.
	Equations() (xeq, ufunc func(y []float64, domain float64, x, u []float64))
	// Dimensions of x state variables and u inputs
	Dims() (nx, nu int)
}

// Integrator abstracts algorithm specifics. For anyone looking to
// implement it, Set(ivp) should be called first to initialize the IVP with
// initial values. Step will calculate the next x values and store them in y
// u values should not be stored as they can easily be obtained if one has
// x values. Integrator should store 1 or more (depending on algorithm used)
// of previously calculated x values to be able to integrate.
type Integrator interface {
	// Set initializes an initial value problem. First argument
	// is the initial domain integration point, is usually zero.
	Set(float64, IVP) error
	// Step integrates IVP and stores result in y. step is a suggested step
	// for the algorithm to take. The algorithm may decide that it is not sufficiently
	// small or big enough (these are adaptive algorithms) and take a different step.
	// The resulting step is returned as the first parameter
	Step(y []float64, step float64) (float64, error)
}
