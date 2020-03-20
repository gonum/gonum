// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve_test

import (
	"fmt"
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/linsolve"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// AllenCahnFD implements a semi-implicit finite difference scheme for the
// solution of the one-dimensional Allen-Cahn equation
//   u_t = u_xx - 1/ξ²·f'(u)  in (0,L)×(0,T)
//   u_x = 0                  on (0,T)
//  u(0) = u0                 on (0,L)
// where f is a double-well-shaped function with two minima at ±1
//  f(s) = 1/4·(s²-1)²
//
// The equation arises in materials science in the description of phase
// transitions, e.g. solidification in crystal growth, but also in other areas
// like image processing due to its connection to mean-curvature motion.
// Starting the evolution from an initial distribution u0, the solution u
// develops a thin steep layer, an interface between regions of the domain
// (0,L) where u is constant and close to one of the minima of f.
//
// AllenCahnFD approximates derivatives by finite differences and the solution
// is advanced in time by a semi-implicit Euler scheme where the nonlinear term
// is taken from the previous time step. Therefore, at each time step a linear
// system must be solved.
type AllenCahnFD struct {
	// Xi is the ξ parameter that determines the interface width.
	Xi float64

	// InitCond is the initial condition u0.
	InitCond func(x float64) float64

	h   float64 // Spatial step size
	tau float64 // Time step size

	a *mat.SymBandDense
	b *mat.VecDense
	u *mat.VecDense

	ls         linsolve.Method
	lssettings linsolve.Settings
}

// FPrime returns the value of the derivative of the double-well potential f at s.
//  f'(s) = s·(s²-1)
func FPrime(s float64) float64 {
	return s * (s*s - 1)
}

// Setup initializes the receiver for solving the Allen-Cahn equation on a
// uniform grid with n+1 nodes on the spatial interval (0,L) and with the time
// step size tau.
func (ac *AllenCahnFD) Setup(n int, L float64, tau float64) {
	ac.h = L / float64(n)
	ac.tau = tau

	// We solve this PDE numerically by replacing the derivatives with finite
	// differences. For the spatial derivative, we use a central difference
	// scheme, and for the time derivative derivative we use semi-implicit Euler
	// integration where the non-linear term with f' is taken from the previous
	// time step.
	//
	// After replacing the derivatives we get
	//  1/tau*(u^{k+1}_i - u^k_i) = 1/h²*(u^{k+1}_{i-1} - 2*u^{k+1}_i + u^{k+1}_{i+1}) - 1/ξ²*f'(u^k_i)
	// where tau is the time step size, h is the spatial step size, and u^k_i
	// denotes the numerical solution that approximates u at time level k and
	// grid node i, that is,
	//  u^k_i ≅ u(k*tau,i*h)
	// Multiplying the equation by tau and collecting the terms from the same
	// time level on each side gives
	//  -tau/h²*u^{k+1}_{i-1} + (1+2*tau/h²)*u^{k+1}_i - tau/h²*u^{k+1}_{i+1}) = u^k_i - tau/ξ²*f'(u^k_i)
	// If we denote C:=tau/h² we can simplify this to
	//  -C*u^{k+1}_{i-1} + (1+2*C)*u^{k+1}_i - C*u^{k+1}_{i+1} = u^k_i - tau/ξ²*f'(u^k_i)   (1)
	// which must hold for all i=0,...,n.
	//
	// When i=0 or i=n, the expression (1) refers to values at nodes -1 and n+1
	// which lie outside of the domain. We can eliminate them by using the fact
	// that the first derivative is zero at the boundary. We approximate it by
	// central difference:
	//      -1/h*(u^{k+1}_{-1} - u^{k+1}_1) = 0
	//  1/h*(u^{k+1}_{n+1} - u^{k+1}_{n-1}) = 0
	// which after simplifying gives
	//   u^{k+1}_{-1} = u^{k+1}_1
	//  u^{k+1}_{n+1} = u^{k+1}_{n-1}
	// By substituting these two expressions into (1) at i=0 and i=n values at
	// outside nodes do not appear in the expressions. If we further divide them
	// by 2 (to obtain a symmetric matrix), we finally get
	//  (1/2+C)*u^{k+1}_0 - C*u^{k+1}_1 = 1/2*u^k_0 - 1/2*tau/ξ²*f'(u^k_0)                  (2)
	//  -C*u^{k+1}_{n-1} + (1/2+C)*u^{k+1}_n = 1/2*u^k_n - 1/2*tau/ξ²*f'(u^k_n)             (3)
	// Note that simply means that we treat values at the boundary nodes the
	// same as the inner nodes.
	//
	// For a given k the equations (1),(2),(3) form a linear system for the
	// unknown vector [u^{k+1}_i], i=0,...,n that must be solved at each step in
	// order to advance the solution in time. The matrix of this system is
	// tridiagonal and symmetric positive-definite:
	//  ⎛1/2+C   -C    0    .    .    .    .     0⎞
	//  ⎜   -C 1+2C   -C                         .⎟
	//  ⎜    0   -C 1+2C   -C                    .⎟
	//  ⎜    .        -C    .    .               .⎟
	//  ⎜    .              .    .    .          .⎟
	//  ⎜    .                   .    .   -C     0⎟
	//  ⎜    .                       -C 1+2C    -C⎟
	//  ⎝    0    .    .    .    .    0   -C 1/2+C⎠
	// The right-hand side is:
	//  ⎛1/2*u^k_0     - 1/2*tau/ξ²*f'(u^k_0)    ⎞
	//  ⎜    u^k_1     -     tau/ξ²*f'(u^k_1)    ⎟
	//  ⎜    u^k_2     -     tau/ξ²*f'(u^k_2)    ⎟
	//  ⎜              .                         ⎟
	//  ⎜              .                         ⎟
	//  ⎜              .                         ⎟
	//  ⎜    u^k_{n-1} -     tau/ξ²*f'(u^k_{n-1})⎟
	//  ⎝1/2*u^k_n     - 1/2*tau/ξ²*f'(u^k_n)    ⎠

	// Assemble the system matrix A based on the discretization scheme above.
	// Since the matrix is symmetric, we only need to set elements in the upper
	// triangle.
	A := mat.NewSymBandDense(n+1, 1, nil)
	c := ac.tau / ac.h / ac.h
	// Boundary condition at the left node
	A.SetSymBand(0, 0, 0.5+c)
	A.SetSymBand(0, 1, -c)
	// Interior nodes
	for i := 1; i < n; i++ {
		A.SetSymBand(i, i, 1+2*c)
		A.SetSymBand(i, i+1, -c)
	}
	// Boundary condition at the right node
	A.SetSymBand(n, n, 0.5+c)
	ac.a = A

	// Allocate the right-hand side b.
	ac.b = mat.NewVecDense(n+1, nil)

	// Allocate and set up the initial condition.
	ac.u = mat.NewVecDense(n+1, nil)
	for i := 0; i < ac.u.Len(); i++ {
		ac.u.SetVec(i, ac.InitCond(float64(i)*ac.h))
	}

	// Allocate the linear solver and the settings.
	ac.ls = &linsolve.CG{}
	ac.lssettings = linsolve.Settings{
		// Solution from the previous time step will be a good initial estimate.
		InitX: ac.u,
		// Store the solution into the existing vector.
		Dst: ac.u,
		// Provide context to reduce memory allocation and GC pressure.
		Work: linsolve.NewContext(n + 1),
	}
}

// Step advances the solution one step in time.
func (ac *AllenCahnFD) Step() error {
	// Assemble the right-hand side vector b.
	tauXi2 := ac.tau / ac.Xi / ac.Xi
	n := ac.u.Len()
	for i := 0; i < ac.u.Len(); i++ {
		ui := ac.u.AtVec(i)
		bi := ui - tauXi2*FPrime(ui)
		if i == 0 || i == n-1 {
			bi *= 0.5
		}
		ac.b.SetVec(i, bi)
	}
	// Solve the system.
	_, err := linsolve.Iterative(ac, ac.b, ac.ls, &ac.lssettings)
	return err
}

// MulVecTo implements the MulVecToer interface.
func (ac *AllenCahnFD) MulVecTo(dst *mat.VecDense, _ bool, x mat.Vector) {
	dst.MulVec(ac.a, x)
}

func (ac *AllenCahnFD) Solution() *mat.VecDense {
	return ac.u
}

func output(u mat.Vector, L float64, step int) error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = fmt.Sprintf("Step %d", step)
	p.X.Label.Text = "x"
	p.X.Min = 0
	p.X.Max = L
	p.Y.Min = -1.1
	p.Y.Max = 1.1

	n := u.Len()
	h := L / float64(n-1)
	pts := make(plotter.XYs, n)
	for i := 0; i < u.Len(); i++ {
		pts[i].X = float64(i) * h
		pts[i].Y = u.AtVec(i)
	}
	err = plotutil.AddLines(p, "u", pts)
	if err != nil {
		return err
	}
	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, fmt.Sprintf("u%04d.png", step))
	if err != nil {
		return err
	}
	return nil
}

func ExampleIterative_evolutionPDE() {
	const (
		L   = 10.0
		nx  = 1000
		nt  = 200
		tau = 0.1 * L / nx
		xi  = 6.0 * L / nx
	)
	rnd := rand.New(rand.NewSource(1))
	ac := AllenCahnFD{
		Xi: xi,
		InitCond: func(x float64) float64 {
			// Initial condition is a perturbation of the (unstable) zero state
			// (the peak in the double-well function f).
			return 0.01 * rnd.NormFloat64()
		},
	}
	ac.Setup(nx, L, tau)
	for i := 1; i <= nt; i++ {
		err := ac.Step()
		if err != nil {
			log.Fatal(err)
		}
		// Generate plots of u as PNG images that can be converted to video
		// using for example
		//  ffmpeg -i u%04d.png output.webm
		err = output(ac.Solution(), L, i)
		if err != nil {
			log.Fatal(err)
		}
	}
}
