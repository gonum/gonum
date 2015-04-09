// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"errors"
	"math"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

// Newton implements pure Newton's method without linesearch for unconstrained
// minimization of functions that satisfy Gradient and Hessian interfaces.
//
// Newton's method iteratively forms a quadratic model to the objective
// function f and tries to minimize this approximate model. It generates a
// sequence of locations x_k by means of
//  solve H_k d_k = -∇f_k for d_k,
//  x_{k+1} = x_k + d_k,
// where H_k is the Hessian matrix of f at x_k.
//
// Note that Newton's method is only locally convergent. This means that for
// starting points close to a minimizer the sequence x_k converges rapidly with
// a quadratic rate. However, away from a minimizer the quadratic approximation
// may not be valid, the Hessian may not be positive definite and d_k may not
// be a descent direction. In such case Newton fails with an error.
//
// For a globally convergent Hessian-based method, see ModifiedNewton. If the
// Hessian matrix cannot be formed explicitly or if the computational cost of
// its factorization is prohibitive, BFGS or L-BFGS methods can be used
// instead.
type Newton struct {
	chol *mat64.TriDense
	step *mat64.Dense
}

func (n *Newton) Init(loc *Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	dim := len(loc.X)
	n.chol = resizeTriDense(n.chol, dim)
	r, _ := n.step.Dims()
	if dim < r {
		n.step = n.step.View(0, 0, dim, 1).(*mat64.Dense)
	} else if dim > r {
		n.step.Grow(dim, 1)
	}
	return n.computeNextX(loc, xNext)
}

func (n *Newton) Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	return n.computeNextX(loc, xNext)
}

func (n *Newton) computeNextX(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	pd := n.chol.Cholesky(loc.Hessian, false)
	if !pd {
		return NoEvaluation, NoIteration, errors.New("optimize: Hessian is not positive definite")
	}
	n.step.SolveCholesky(n.chol, mat64.NewDense(len(loc.Gradient), 1, loc.Gradient))
	floats.AddScaledTo(xNext, loc.X, -1, n.step.RawMatrix().Data)
	return FuncEvaluation | GradEvaluation | HessEvaluation, MajorIteration, nil
}

func (n *Newton) Needs() struct {
	Gradient bool
	Hessian  bool
} {
	return struct {
		Gradient bool
		Hessian  bool
	}{true, true}
}

// ModifiedNewton implements a modified Newton's method for unconstrained
// minimization of functions that satisfy Gradient and Hessian interfaces. It
// is a globally convergent line-search method that can be applied from any
// starting point.
//
// ModifiedNewton method is based on the Newton step which computes the
// line-search direction d_k as the solution of the linear system
//  H_k d_k = -∇f_k,
// where H_k is the Hessian matrix of the objective function. Away from a
// minimizer H_k may not be positive definite and d_k may not be a descent
// direction. ModifiedNewton implements a Hessian modification strategy that
// tries to add successively larger multiples of identity to H_k until it
// becomes positive definite. Note that the repeated trial factorization of the
// modified Hessian involved in this process can be computationally expensive.
// For points in the vicinity of a minimizer ModifiedNewton reduces to a pure
// Newton's method with its quadratic convergence rate.
//
// If the Hessian matrix cannot be formed explicitly or if the computational
// cost of its factorization is prohibitive, BFGS or L-BFGS quasi-Newton method
// can be used instead.
type ModifiedNewton struct {
	// LinesearchMethod is a method used for selecting suitable steps along the
	// descent direction d. Steps should satisfy the Wolfe, Goldstein or Armijo
	// conditions. If LinesearchMethod == nil, an appropriate default is
	// chosen.
	LinesearchMethod LinesearchMethod
	// Increase is a factor by which a scalar tau is successively increased so
	// that (H + tau*I) is positive definite. Larger values reduce the number
	// of trial Hessian factorizations, but also reduce the second-order
	// information in H.
	// Increase must be greater than 1. If Increase is 0, it is defaulted to 5.
	Increase float64

	linesearch *Linesearch

	hess *mat64.SymDense
	chol *mat64.TriDense
	tau  float64
}

func (n *ModifiedNewton) Init(loc *Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	if n.Increase == 0 {
		n.Increase = 5
	}
	if n.Increase <= 1 {
		panic("optimize: ModifiedNewton.Increase must be greater than 1")
	}
	if n.LinesearchMethod == nil {
		n.LinesearchMethod = &Bisection{}
	}
	if n.linesearch == nil {
		n.linesearch = &Linesearch{}
	}
	n.linesearch.Method = n.LinesearchMethod
	n.linesearch.NextDirectioner = n

	return n.linesearch.Init(loc, f, xNext)
}

func (n *ModifiedNewton) Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	return n.linesearch.Iterate(loc, xNext)
}

func (n *ModifiedNewton) InitDirection(loc *Location, dir []float64) (stepSize float64) {
	dim := len(loc.X)
	n.chol = resizeTriDense(n.chol, dim)
	n.hess = resizeSymDense(n.hess, dim)
	n.tau = 0
	n.computeNextDir(loc, dir)
	return 1
}

func (n *ModifiedNewton) NextDirection(loc *Location, dir []float64) (stepSize float64) {
	n.computeNextDir(loc, dir)
	return 1
}

func (n *ModifiedNewton) computeNextDir(loc *Location, dir []float64) {
	dim := len(loc.X)
	n.hess.CopySym(loc.Hessian)

	// Find the smallest diagonal entry of the Hesssian.
	minA := n.hess.At(0, 0)
	for i := 1; i < dim; i++ {
		a := n.hess.At(i, i)
		if a < minA {
			minA = a
		}
	}
	// If the smallest diagonal entry is positive, the Hessian may be positive
	// definite, and so first attempt to apply the Cholesky factorization to
	// the un-modified Hessian. If the smallest entry is negative, use the
	// final tau from the last iteration if regularization was needed,
	// otherwise guess an appropriate value for tau.
	if minA > 0 {
		n.tau = 0
	} else if n.tau == 0 {
		n.tau = -minA + 0.001
	}

	for {
		if n.tau != 0 {
			// Add a multiple of identity to the Hessian.
			for i := 0; i < dim; i++ {
				n.hess.SetSym(i, i, loc.Hessian.At(i, i)+n.tau)
			}
		}
		// Try to apply the Cholesky factorization.
		pd := n.chol.Cholesky(n.hess, false)
		if pd {
			break
		}
		// Modified Hessian is not PD, so increase tau.
		n.tau = math.Max(n.Increase*n.tau, 0.001)
	}
	d := mat64.NewDense(dim, 1, dir)
	// Store the solution in d's backing array, dir.
	d.SolveCholesky(n.chol, mat64.NewDense(dim, 1, loc.Gradient))
	floats.Scale(-1, dir)
}

func (n *ModifiedNewton) Needs() struct {
	Gradient bool
	Hessian  bool
} {
	return struct {
		Gradient bool
		Hessian  bool
	}{true, true}
}
