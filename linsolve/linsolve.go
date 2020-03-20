// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// BreakdownError signifies that a breakdown occured and the method cannot continue.
type BreakdownError struct {
	Value     float64
	Tolerance float64
}

func (e *BreakdownError) Error() string {
	return fmt.Sprintf("linsolve: breakdown, value=%v tolerance=%v", e.Value, e.Tolerance)
}

// Method is an iterative method that produces a sequence of vectors that
// converge to the solution of the system of linear equations
//  A * x = b,
// where A is non-singular n×n matrix, and x and b are vectors of dimension n.
//
// Method uses a reverse-communication interface between the iterative algorithm
// and the caller. Method acts as a client that commands the caller to perform
// needed operations via an Operation returned from the Iterate method. This
// provides independence of Method on representation of the matrix A, and
// enables automation of common operations like checking for convergence and
// maintaining statistics.
type Method interface {
	// Init initializes the method for solving an n×n
	// linear system with an initial estimate x and
	// the corresponding residual vector.
	//
	// Method will not retain x or residual.
	Init(x, residual *mat.VecDense)

	// Iterate performs a step in converging to the
	// solution of a linear system.
	//
	// Iterate retrieves data from Context, updates it,
	// and returns the next operation. The caller must
	// perform the Operation using data in Context, and
	// depending on the state call Iterate again.
	Iterate(*Context) (Operation, error)
}

// Context mediates the communication between the Method and
// the caller. The caller must not modify Context apart from
// the commanded Operations.
type Context struct {
	// X will be set by Method to the current approximate
	// solution when it commands ComputeResidual and MajorIteration.
	X *mat.VecDense

	// ResidualNorm is (an estimate of) a norm of
	// the residual. Method will set it to the current
	// value when it commands CheckResidualNorm.
	ResidualNorm float64

	// Converged indicates to Method whether ResidualNorm
	// satisfies a stopping criterion as a result of
	// CheckResidualNorm operation.
	Converged bool

	// Src and Dst are the source and destination vectors
	// for various Operations. Src will be set by Method
	// and the caller must store the result in Dst. See
	// the Operation documentation for more information.
	Src, Dst *mat.VecDense
}

// NewContext returns a new Context for work on problems of dimension n.
// NewContext will panic if n is not positive.
func NewContext(n int) *Context {
	if n <= 0 {
		panic("linsolve: context size is not positive")
	}
	return &Context{
		X:   mat.NewVecDense(n, nil),
		Src: mat.NewVecDense(n, nil),
		Dst: mat.NewVecDense(n, nil),
	}
}

// Reset reinitializes the Context for work on problems of dimension n.
// Reset will panic if n is not positive.
func (ctx *Context) Reset(n int) {
	if n <= 0 {
		panic("linsolve: dimension not positive")
	}
	ctx.X.Reset()
	ctx.X.ReuseAsVec(n)
	ctx.Src.Reset()
	ctx.Src.ReuseAsVec(n)
	ctx.Dst.Reset()
	ctx.Dst.ReuseAsVec(n)
}

// Operation specifies the type of operation.
type Operation uint

// Operations commanded by Method.Iterate.
const (
	NoOperation Operation = 0

	// Compute A*x where x is stored in Context.Src. The
	// result must be placed in Context.Dst.
	MulVec Operation = 1 << (iota - 1)

	// Perform a preconditioner solve
	//  M z = r
	// where r is stored in Context.Src. The solution z
	// must be placed in Context.Dst.
	PreconSolve

	// Trans indicates that MulVec or PreconSolve
	// operation must be performed wih the transpose,
	// that is, compute Aᵀ*x or solve Mᵀ*z = r. Method
	// will command Trans only in bitwise OR combination
	// with MulVec and PreconSolve.
	Trans

	// Compute b-A*x where x is stored in Context.X,
	// and store the result in Context.Dst.
	ComputeResidual

	// Check convergence using (an estimate of) a
	// residual norm in Context.ResidualNorm. Context.X
	// does not need to be valid. The caller must set
	// Context.Converged to indicate whether convergence
	// has been determined, and then it must call
	// Method.Iterate again.
	CheckResidualNorm

	// MajorIteration indicates that Method has finished
	// what it considers to be one iteration. Method
	// will make sure that Context.X is updated.
	// If Context.Converged is true, the caller must
	// terminate the iterative process, otherwise it
	// should call Method.Iterate again.
	MajorIteration
)

const (
	// Machine epsilon.
	eps = 1.0 / (1 << 53)

	// Breakdown tolerance for BiCG and BiCGStab methods.
	breakdownTol = eps * eps
)
