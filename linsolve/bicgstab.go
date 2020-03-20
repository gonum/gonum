// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// BiCGStab implements the BiConjugate Gradient Stabilized method with
// preconditioning for solving systems of linear equations
//  A * x = b,
// where A is a nonsymmetric, nonsingular matrix. The method is a variant of
// BiCG but offers a smoother convergence and does not require multiplication
// with Aᵀ.
//
// References:
//  - Barrett, R. et al. (1994). Section 2.3.8 BiConjugate Gradient Stabilized (Bi-CGSTAB).
//    In Templates for the Solution of Linear Systems: Building Blocks
//    for Iterative Methods (2nd ed.) (pp. 24-25). Philadelphia, PA: SIAM.
//    Retrieved from http://www.netlib.org/templates/templates.pdf
type BiCGStab struct {
	x     mat.VecDense
	r, rt mat.VecDense
	p     mat.VecDense
	phat  mat.VecDense
	shat  mat.VecDense
	t     mat.VecDense
	v     mat.VecDense

	rho, rhoPrev float64
	alpha        float64
	omega        float64

	resume int
}

// Init initializes the data for a linear solve. See the Method interface for more details.
func (b *BiCGStab) Init(x, residual *mat.VecDense) {
	dim := x.Len()
	if residual.Len() != dim {
		panic("bicgstab: vector length mismatch")
	}

	b.x.CloneVec(x)
	b.r.CloneVec(residual)
	b.rt.CloneVec(&b.r)

	b.p.Reset()
	b.p.ReuseAsVec(dim)
	b.phat.Reset()
	b.phat.ReuseAsVec(dim)
	b.shat.Reset()
	b.shat.ReuseAsVec(dim)
	b.t.Reset()
	b.t.ReuseAsVec(dim)
	b.v.Reset()
	b.v.ReuseAsVec(dim)

	b.rhoPrev = 1
	b.alpha = 0
	b.omega = 1

	b.resume = 1
}

// Iterate performs an iteration of the linear solve. See the Method interface for more details.
//
// BiCGStab will command the following operations:
//  MulVec
//  PreconSolve
//  CheckResidualNorm
//  MajorIteration
//  NoOperation
func (b *BiCGStab) Iterate(ctx *Context) (Operation, error) {
	switch b.resume {
	case 1:
		b.rho = mat.Dot(&b.rt, &b.r)
		if math.Abs(b.rho) < breakdownTol {
			b.resume = 0
			return NoOperation, &BreakdownError{math.Abs(b.rho), breakdownTol}
		}
		// p_i = r_{i-1} + beta*(p_{i-1} - omega * v_{i-1})
		beta := (b.rho / b.rhoPrev) * (b.alpha / b.omega)
		b.p.AddScaledVec(&b.p, -b.omega, &b.v)
		b.p.AddScaledVec(&b.r, beta, &b.p)
		// Solve M^{-1} * p_i.
		ctx.Src.CopyVec(&b.p)
		b.resume = 2
		return PreconSolve, nil
	case 2:
		b.phat.CopyVec(ctx.Dst)
		// Compute A * \hat{p}_i.
		ctx.Src.CopyVec(&b.phat)
		b.resume = 3
		return MulVec, nil
	case 3:
		b.v.CopyVec(ctx.Dst)
		rtv := mat.Dot(&b.rt, &b.v)
		if rtv == 0 {
			b.resume = 0
			return NoOperation, &BreakdownError{}
		}
		b.alpha = b.rho / rtv
		// Form the residual and X so that we can check for tolerance early.
		ctx.X.AddScaledVec(ctx.X, b.alpha, &b.phat)
		b.r.AddScaledVec(&b.r, -b.alpha, &b.v)
		ctx.ResidualNorm = mat.Norm(&b.r, 2)
		b.resume = 4
		return CheckResidualNorm, nil
	case 4:
		if ctx.Converged {
			b.resume = 0
			return MajorIteration, nil
		}
		// Solve M^{-1} * r_i.
		ctx.Src.CopyVec(&b.r)
		b.resume = 5
		return PreconSolve, nil
	case 5:
		b.shat.CopyVec(ctx.Dst)
		// Compute A * \hat{s}_i.
		ctx.Src.CopyVec(&b.shat)
		b.resume = 6
		return MulVec, nil
	case 6:
		b.t.CopyVec(ctx.Dst)
		b.omega = mat.Dot(&b.t, &b.r) / mat.Dot(&b.t, &b.t)
		ctx.X.AddScaledVec(ctx.X, b.omega, &b.shat)
		b.r.AddScaledVec(&b.r, -b.omega, &b.t)
		ctx.ResidualNorm = mat.Norm(&b.r, 2)
		b.resume = 7
		return CheckResidualNorm, nil
	case 7:
		if ctx.Converged {
			b.resume = 0
			return MajorIteration, nil
		}
		if math.Abs(b.omega) < breakdownTol {
			b.resume = 0
			return NoOperation, &BreakdownError{math.Abs(b.omega), breakdownTol}
		}
		b.rhoPrev = b.rho
		b.resume = 1
		return MajorIteration, nil

	default:
		panic("bicgstab: Init not called")
	}
}
