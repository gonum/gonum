// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivp

import (
	"errors"

	"gonum.org/v1/gonum/floats"
)

// RK4 implements Integrator interface for Runke-Kutta
// 4th order multivariable method (non adaptive)
type RK4 struct {
	dom        float64
	x, u       []float64
	a, b, c, d []float64
	fx, fu     func(y []float64, t float64, x, u []float64)
}

// Step implements Integrator interface
func (rk *RK4) Step(y []float64, h float64) (float64, error) {
	const overSix = 1. / 6.
	t := rk.dom
	// set a, b, c, d (in some literatures these are k1,k2,k3,k4)
	rk.fx(rk.a, t, rk.x, rk.u)
	rk.fx(rk.b, t+0.5*h, floats.AddScaledTo(rk.b, rk.x, 0.5*h, rk.a), rk.u)
	rk.fx(rk.c, t+0.5*h, floats.AddScaledTo(rk.c, rk.x, 0.5*h, rk.b), rk.u)
	rk.fx(rk.d, t+h, floats.AddScaledTo(rk.d, rk.x, h, rk.c), rk.u)

	floats.Add(rk.a, rk.d)
	floats.Add(rk.b, rk.c)
	floats.AddScaled(rk.a, 2, rk.b)
	floats.AddScaledTo(y, rk.x, h*overSix, rk.a)
	// finished integrating. Now we update RK4 structure for future Steps
	copy(rk.x, y) // store results
	rk.dom += h
	if len(rk.u) > 0 {
		rk.fu(rk.u, t, rk.x, rk.u)
	}
	return h, nil
}

// Set implements integrator interface. All RK4 data
// is reset when calling Set.
func (rk *RK4) Set(d0 float64, ivp IVP) error {
	if ivp == nil {
		return errors.New("IVP is nil")
	}
	nx, nu := ivp.Dims()
	rk.dom = d0 //set start domain
	rk.a = make([]float64, nx)
	rk.b = make([]float64, nx)
	rk.c = make([]float64, nx)
	rk.d = make([]float64, nx)
	rk.x = make([]float64, nx)
	rk.u = make([]float64, nu)
	// set initial values
	x0, u0 := ivp.IV()
	for i := 0; i < x0.Len(); i++ {
		rk.x[i] = x0.AtVec(i)
	}
	if nu > 0 {
		for i := 0; i < u0.Len(); i++ {
			rk.u[i] = u0.AtVec(i)
		}
	}
	rk.fx, rk.fu = ivp.Equations()
	return nil
}
