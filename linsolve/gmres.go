// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

// GMRES implements the Generalized Minimum Residual method with the modified
// Gram-Schmidt orthogonalization for solving systems of linear equations
//  A * x = b,
// where A is a nonsymmetric, nonsingular matrix. It may find a solution in
// fewer iterations and with fewer matrix-vector products compared to BiCG or
// BiCGStab at the price of much large memory storage. This implementation uses
// restarts to limit the memory requirements. GMRES does not need the
// multiplication with Aᵀ.
//
// References:
//  - Barrett, R. et al. (1994). Section 2.3.4 Generalized Minimal Residual
//    (GMRES). In Templates for the Solution of Linear Systems: Building Blocks
//    for Iterative Methods (2nd ed.) (pp. 17-19). Philadelphia, PA: SIAM.
//    Retrieved from http://www.netlib.org/templates/templates.pdf
//  - Saad, Y., and Schultz, M. (1986). GMRES: A generalized minimal residual
//    algorithm for solving nonsymmetric linear systems. SIAM J. Sci. Stat.
//    Comput., 7(3), 856. doi:10.6028/jres.049.044
//    Retrieved from https://web.stanford.edu/class/cme324/saad-schultz.pdf
type GMRES struct {
	// Restart is the restart parameter which limits the computation and
	// storage costs. It must hold that
	//  1 <= Restart <= n
	// where n is the dimension of the problem. If Restart is 0, n will be
	// used instead. This guarantees convergence of GMRES and increases
	// robustness. Many specific problems however, particularly for large
	// n, will benefit in efficiency by setting Restart to
	// a problem-dependent value less than n.
	Restart int

	// m is the used value of Restart.
	m int
	// v is an n×(m+1) matrix V whose columns form an orthonormal basis of the
	// Krylov subspace.
	v mat.Dense
	// h is an (m+1)×m upper Hessenberg matrix H.
	h mat.Dense
	// givs holds Givens rotations that are used to reduce H to upper triangular
	// form.
	givs []givens

	x mat.VecDense
	y mat.VecDense
	s mat.VecDense

	k      int // Loop variable for inner iterations.
	resume int
}

// Init initializes the data for a linear solve. See the Method interface for more details.
func (g *GMRES) Init(x, residual *mat.VecDense) {
	dim := x.Len()
	if residual.Len() != dim {
		panic("gmres: vector length mismatch")
	}

	g.m = g.Restart
	if g.m == 0 {
		g.m = dim
	}
	if g.m <= 0 || dim < g.m {
		panic("gmres: invalid value of Restart")
	}

	g.v.Reset()
	g.v.ReuseAs(dim, g.m+1)
	// Store the residual in the first column of V.
	g.vcol(0).CopyVec(residual)

	g.h.Reset()
	g.h.ReuseAs(g.m+1, g.m)

	if cap(g.givs) < g.m {
		g.givs = make([]givens, g.m)
	} else {
		g.givs = g.givs[:g.m]
		for i := range g.givs {
			g.givs[i].c = 0
			g.givs[i].s = 0
		}
	}

	g.x.CloneVec(x)
	g.y.Reset()
	g.y.ReuseAsVec(g.m + 1)
	g.s.Reset()
	g.s.ReuseAsVec(g.m + 1)

	g.resume = 1
}

// Iterate performs an iteration of the linear solve. See the Method interface for more details.
//
// GMRES will command the following operations:
//  MulVec
//  PreconSolve
//  CheckResidualNorm
//  MajorIteration
//  NoOperation
func (g *GMRES) Iterate(ctx *Context) (Operation, error) {
	switch g.resume {
	case 1:
		// The initial residual is in the first column of V.
		ctx.Src.CopyVec(g.vcol(0))
		g.resume = 2
		// Solve M^{-1} * r_0.
		return PreconSolve, nil
	case 2:
		// v_0 = M^{-1} * r_0
		v0 := g.vcol(0)
		v0.CopyVec(ctx.Dst)
		// Normalize v_0.
		norm := mat.Norm(v0, 2)
		v0.ScaleVec(1/norm, v0)
		// Initialize s to the elementary vector e_1 scaled by norm.
		g.s.Zero()
		g.s.SetVec(0, norm)

		// Begin the inner for-loop for k going from 0 to m-1.
		g.k = 0
		fallthrough
	case 3:
		ctx.Src.CopyVec(g.vcol(g.k))
		g.resume = 4
		// Compute A * v_k.
		return MulVec, nil
	case 4:
		ctx.Src.CopyVec(ctx.Dst)
		g.resume = 5
		// Solve M^{-1} * (A * v_k).
		return PreconSolve, nil
	case 5:
		// v_{k+1} = M^{-1} * (A * v_k)
		vk1 := g.vcol(g.k + 1)
		vk1.CopyVec(ctx.Dst)
		// Construct the k-th column of the upper Hessenberg matrix H
		// using the modified Gram-Schmidt process to make v_{k+1}
		// orthonormal to the first k+1 columns of V.
		g.modifiedGS(g.k, &g.h, &g.v, vk1)
		// Reduce H back to upper triangular form and update the vector s.
		g.qr(g.k, g.givs, &g.h, &g.s)
		// Check the approximate residual norm.
		ctx.ResidualNorm = math.Abs(g.s.AtVec(g.k + 1))
		g.resume = 6
		return CheckResidualNorm, nil
	case 6:
		g.k++
		if g.k < g.m && !ctx.Converged {
			// Continue the inner for-loop.
			g.resume = 3
			return NoOperation, nil
		}
		// Either restarting or converged, we have to update the solution.
		// Solve the upper triangular system H*y=s.
		g.solveLeastSquares(g.k, &g.y, &g.h, &g.s)
		// Compute x as a linear combination of columns of V.
		g.updateSolution(g.k, &g.x, &g.v, &g.y)
		ctx.X.CopyVec(&g.x)
		if ctx.Converged {
			g.resume = 0
			return MajorIteration, nil
		}
		// We are restarting, so we have to also compute the residual.
		g.resume = 7
		return ComputeResidual, nil
	case 7:
		// Store the residual again in the first column of V.
		g.vcol(0).CopyVec(ctx.Dst)
		g.resume = 1
		return MajorIteration, nil

	default:
		panic("gmres: Init not called")
	}
}

// modifiedGS orthonormalizes the vector w with respect to the first k+1 columns
// of V using the modified Gram-Schmidt algorithm, and stores the computed
// coefficients in the k-th column of H.
func (g *GMRES) modifiedGS(k int, h, v *mat.Dense, w *mat.VecDense) {
	hk := h.ColView(k).(*mat.VecDense)
	for j := 0; j <= k; j++ {
		vj := v.ColView(j).(*mat.VecDense)
		hkj := mat.Dot(vj, w)
		hk.SetVec(j, hkj)           // H[j,k] = v_j · w
		w.AddScaledVec(w, -hkj, vj) // w -= H[j,k] * v_j
	}
	norm := mat.Norm(w, 2)
	hk.SetVec(k+1, norm)  // H[k+1,k] = |w|
	w.ScaleVec(1/norm, w) // Normalize w.
}

// qr applies previous Givens rotations to the k-th column of H, computes the
// next Givens rotation to zero out H[k+1,k] and applies it also to the vector s.
func (g *GMRES) qr(k int, givs []givens, h *mat.Dense, s *mat.VecDense) {
	// Apply previous Givens rotations to the k-th column of H.
	hk := h.ColView(k).(*mat.VecDense)
	for i, giv := range givs[:k] {
		hki, hki1 := giv.apply(hk.AtVec(i), hk.AtVec(i+1))
		hk.SetVec(i, hki)
		hk.SetVec(i+1, hki1)
	}

	// Compute the k-th Givens rotation that zeros H[k+1,k].
	givs[k].c, givs[k].s, _, _ = blas64.Rotg(hk.AtVec(k), hk.AtVec(k+1))

	// Apply the k-th Givens rotation to (H[k,k], H[k+1,k]).
	hkk, _ := givs[k].apply(hk.AtVec(k), hk.AtVec(k+1))
	hk.SetVec(k, hkk)
	// We don't call hk.SetVec(k+1, 0) because the element won't be accessed.

	// Apply the k-th Givens rotation to (s[k], s[k+1]).
	sk, sk1 := givs[k].apply(s.AtVec(k), s.AtVec(k+1))
	s.SetVec(k, sk)
	s.SetVec(k+1, sk1)
}

// solveLeastSquares solves the upper triangular linear system
//  H * y = s
func (g *GMRES) solveLeastSquares(k int, y *mat.VecDense, h *mat.Dense, s *mat.VecDense) {
	// Copy the first k elements of s into y.
	y.CopyVec(s.SliceVec(0, k))
	// Convert H into an upper triangular matrix.
	hraw := h.RawMatrix()
	htri := blas64.Triangular{
		Uplo:   blas.Upper,
		Diag:   blas.NonUnit,
		N:      k,
		Data:   hraw.Data,
		Stride: hraw.Stride,
	}
	// Solve the upper triangular system storing the solution into y.
	blas64.Trsv(blas.NoTrans, htri, y.RawVector())
}

// updateSolution updates the current solution vector x with a linear
// combination of the first k columns of V:
//  x = x + V * y = x + \sum y_j * v_j
func (g *GMRES) updateSolution(k int, x *mat.VecDense, v *mat.Dense, y *mat.VecDense) {
	for j := 0; j < k; j++ {
		vj := v.ColView(j)
		x.AddScaledVec(x, y.AtVec(j), vj)
	}
}

// vcol returns a view of the j-th column of the matrix V.
func (g *GMRES) vcol(j int) *mat.VecDense {
	return g.v.ColView(j).(*mat.VecDense)
}

// givens is a Givens rotation.
type givens struct {
	c, s float64
}

func (giv *givens) apply(x, y float64) (float64, float64) {
	return giv.c*x + giv.s*y, giv.c*y - giv.s*x
}
