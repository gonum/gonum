// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivp

import "gonum.org/v1/gonum/mat"

// Model implements IVP interface. X vector and equations can not be nil or zero length.
type Model struct {
	x0, u0   mat.Vector
	xeq, ins func(y []float64, dom float64, x, u []float64)
}

// IV returns initial values of the IVP. First returned parameter is the
// starting x vector and second parameter are inputs when solving non-autonomous
// ODEs.
func (m *Model) IV() (mat.Vector, mat.Vector) { return m.x0, m.u0 }

// Equations returns differential equations relating to state vector x and input functions
// for non-autonomous ODEs.
//
// Input functions may be nil (ueq).
func (m *Model) Equations() (xeq, ueq func(y []float64, dom float64, x, u []float64)) {
	return m.xeq, m.ins
}

// Dims returns dimension of state and input (x length and u length, respectively).
// Dims has some dimension checking involved and can serve as a preliminary system verifier.
//
// Dims panics if x vector is nil.
func (m *Model) Dims() (nx, nu int) {
	if m.x0 == nil {
		panic("x vector can not be nil")
	}
	if m.u0 == nil {
		nu = 0
	}
	return m.x0.Len(), nu
}

// NewModel returns a IVP given initial conditions (x0,u0), differential equations (xeq) and
// input functions for non-autonomous ODEs (ueq).
func NewModel(x0, u0 mat.Vector, xeq, ueq func(y []float64, dom float64, x, u []float64)) (*Model, error) {
	return &Model{xeq: xeq, ins: ueq, x0: x0, u0: u0}, nil
}
