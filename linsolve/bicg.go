// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// BiCG implements the Bi-Conjugate Gradient method with preconditioning for
// solving systems of linear equations
//  A * x = b,
// where A is a nonsymmetric, nonsingular matrix. It uses limited memory storage
// but the convergence may be irregular and the method may break down. BiCG
// requires a multiplication with A and Aᵀ at each iteration. BiCGStab is a
// related method that does not use Aᵀ.
//
// References:
//  - Barrett, R. et al. (1994). Section 2.3.5 BiConjugate Gradient (BiCG).
//    In Templates for the Solution of Linear Systems: Building Blocks
//    for Iterative Methods (2nd ed.) (pp. 19-20). Philadelphia, PA: SIAM.
//    Retrieved from http://www.netlib.org/templates/templates.pdf
type BiCG struct {
	x     mat.VecDense
	r, rt mat.VecDense
	p, pt mat.VecDense
	z, zt mat.VecDense

	rho, rhoPrev float64

	resume int
}

// Init initializes the data for a linear solve. See the Method interface for more details.
func (b *BiCG) Init(x, residual *mat.VecDense) {
	dim := x.Len()
	if residual.Len() != dim {
		panic("bicg: vector length mismatch")
	}

	b.x.CloneVec(x)
	b.r.CloneVec(residual)
	b.rt.CloneVec(&b.r)

	b.p.Reset()
	b.p.ReuseAsVec(dim)
	b.pt.Reset()
	b.pt.ReuseAsVec(dim)
	b.z.Reset()
	b.z.ReuseAsVec(dim)
	b.zt.Reset()
	b.zt.ReuseAsVec(dim)

	b.rhoPrev = 1

	b.resume = 1
}

// Iterate performs an iteration of the linear solve. See the Method interface for more details.
//
// BiCG will command the following operations:
//  MulVec
//  MulVec|Trans
//  PreconSolve
//  PreconSolve|Trans
//  CheckResidualNorm
//  MajorIteration
//  NoOperation
func (b *BiCG) Iterate(ctx *Context) (Operation, error) {
	switch b.resume {
	case 1:
		// Solve M^{-1} * r_{i-1}.
		ctx.Src.CopyVec(&b.r)
		b.resume = 2
		return PreconSolve, nil
	case 2:
		b.z.CopyVec(ctx.Dst)
		// Solve M^{-T} * rt_{i-1}.
		ctx.Src.CopyVec(&b.rt)
		b.resume = 3
		return PreconSolve | Trans, nil
	case 3:
		b.zt.CopyVec(ctx.Dst)
		b.rho = mat.Dot(&b.z, &b.rt)
		if math.Abs(b.rho) < breakdownTol {
			b.resume = 0
			return NoOperation, &BreakdownError{math.Abs(b.rho), breakdownTol}
		}
		beta := b.rho / b.rhoPrev
		b.p.AddScaledVec(&b.z, beta, &b.p)
		b.pt.AddScaledVec(&b.zt, beta, &b.pt)
		// Compute A * p.
		ctx.Src.CopyVec(&b.p)
		b.resume = 4
		return MulVec, nil
	case 4:
		// z is overwritten and reused.
		b.z.CopyVec(ctx.Dst)
		// Compute Aᵀ * pt.
		ctx.Src.CopyVec(&b.pt)
		b.resume = 5
		return MulVec | Trans, nil
	case 5:
		// zt is overwritten and reused.
		b.zt.CopyVec(ctx.Dst)
		alpha := b.rho / mat.Dot(&b.pt, &b.z)
		ctx.X.AddScaledVec(ctx.X, alpha, &b.p)
		b.rt.AddScaledVec(&b.rt, -alpha, &b.zt)
		b.r.AddScaledVec(&b.r, -alpha, &b.z)
		ctx.ResidualNorm = mat.Norm(&b.r, 2)
		b.resume = 6
		return CheckResidualNorm, nil
	case 6:
		if ctx.Converged {
			b.resume = 0
			return MajorIteration, nil
		}
		b.rhoPrev = b.rho
		b.resume = 1
		return MajorIteration, nil

	default:
		panic("bicg: Init not called")
	}
}
