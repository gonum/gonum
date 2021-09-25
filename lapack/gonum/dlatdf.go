// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dlatdf uses the LU factorization of the n×n matrix Z computed by Dgetc2 and
// computes a contribution to the reciprocal Dif-estimate by solving
//  Z * x = b for x
// and choosing the rhs b such that the norm of x is as large as possible.
// On entry rhs = b holds the contribution from earlier solved sub-systems, and on return rhs = x.
//
// The factorization of Z returned by Dgetc2 has the form
//  Z = P*L*U*Q,
// where P and Q are permutation matrices. L is lower triangular with
// unit diagonal elements and U is upper triangular.
//
// It must hold that n <= len(ipiv), n <= len(jpiv). On exit the pivot indices where row i has
// been interchanged with ipiv[i] and column j interchanged with column jpiv[j].
//
// rhs must have n accesible elements.
//
// if ijob==2
//  First compute an approximative null-vector e of Z using DGECON,
//  e is normalized and solve for Zx = +-e - f with the sign giving the greater value
//  of 2-norm(x). About 5 times as expensive as Default.
// if ijob!=2
//  Local look ahead strategy where all entries of the rhs.
//  b is chosen as either +1 or -1 (Default).
//
// rdsum is the sum of squares of computed contributions to the Dif-estimate
// under computation by Dtgsyl, where the scaling factor rdscal (see below)
// has been factored out.
//
// rdscal is the scaling factor used to prevent overflow in rdsum.
//
// NOTE: rdscal and rdsum only makes sense when Dtgsy2 is called by Dtgsyl.
//
// scal is rdscal updated w.r.t. the current contributions in rdsum.
// if trans = 'T'
//  scal == rdscal
//
// sum is the sum of squares updated with the contributions from the current sub-system.
// If trans = 'T'
//  sum = rdsum
//
// Dlatdf implicitly expects z matrix to be at most an 8×8 matrix. This is fulfilled by it's only caller, Dtgsy2.
//
// Dlatdf is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlatdf(ijob, n int, z []float64, ldz int, rhs []float64, rdsum, rdscal float64, ipiv, jpiv []int) (sum, scale float64) {
	switch {
	case n < 0:
		panic(nLT0)
	case ldz < max(1, n):
		panic(badLdZ)
	case len(rhs) < n:
		panic(shortRHS)
	case len(ipiv) < n:
		panic(badLenIpiv)
	case len(jpiv) < n:
		panic(badLenJpiv)
	case n > 8:
		// Dlatdf expects z to be less than 8 in dimension.
		panic("lapack: Dlatdf expects z to be at most 8×8")
	}

	// Quick return if possible.
	if n == 0 {
		return
	}

	if len(z) < (n-1)*ldz+n {
		panic(shortZ)
	}

	bi := blas64.Implementation()
	var temp float64
	xp := make([]float64, n)
	if ijob == 2 {
		// Compute approximate nullvector xm of z when ijob == 2.
		xm := make([]float64, n)
		work := make([]float64, 4*n)
		impl.Dgecon(lapack.MaxRowSum, n, z, ldz, 1, work, make([]int, n)) // iwork is unused.
		bi.Dcopy(n, work[n:], 1, xm, 1)

		// Compute rhs.
		impl.Dlaswp(1, xm, 1, 0, n-2, ipiv[:n-1], -1)
		temp = 1 / math.Sqrt(bi.Ddot(n, xm, 1, xm, 1))
		bi.Dscal(n, temp, xm, 1)
		bi.Dcopy(n, xm, 1, xp, 1)
		bi.Daxpy(n, 1.0, rhs, 1, xp, 1)
		bi.Daxpy(n, -1.0, xm, 1, rhs, 1)
		impl.Dgesc2(n, z, ldz, rhs, ipiv, jpiv)
		impl.Dgesc2(n, z, ldz, xp, ipiv, jpiv)
		if bi.Dasum(n, xp, 1) > bi.Dasum(n, rhs, 1) {
			bi.Dcopy(n, xp, 1, rhs, 1)
		}

		// Compute sum of squares.
		rdscal, rdsum = impl.Dlassq(n, rhs, 1, rdscal, rdsum)
		return rdsum, rdscal
	}

	// Apply permutations ipiv to rhs
	// If ijob != 2 uses look-ahead strategy.
	impl.Dlaswp(1, rhs, 1, 0, n-2, ipiv[:n-1], 1)

	// Solve for L-part choosing rhs either to +1 or -1.
	pmone := -1.0
	for j := 0; j < n-2; j++ {
		bp := rhs[j] + 1
		bm := rhs[j] - 1

		// Look-ahead for L-part rhs[0:n-2] = + or -1, splus and
		// smin computed more efficiently than in Bsolve [1].
		splus := 1 + bi.Ddot(n-j-1, z[(j+1)*ldz+j:], ldz, z[(j+1)*ldz+j:], ldz)
		sminu := bi.Ddot(n-j-1, z[(j+1)*ldz+j:], ldz, rhs[j+1:], 1)
		splus *= rhs[j]
		switch {
		case splus > sminu:
			rhs[j] = bp
		case sminu > splus:
			rhs[j] = bm
		default:
			// In this case the updating sums are equal and we can
			// choose rsh[j] +1 or -1. The first time this happens
			// we choose -1, thereafter +1. This is a simple way to
			// get good estimates of matrices like Byers well-known
			// example (see [1]). (Not done in Bsolve.)
			rhs[j] += pmone
			pmone = 1.0
		}

		// Compute remaining rhs.
		temp = -rhs[j]
		bi.Daxpy(n-j-1, temp, z[(j+1)*ldz+j:], ldz, rhs[j+1:], 1)
	}

	// Solve for U-part, look-ahead for rhs[n-1] = +-1. This is not done
	// in Bsolve and will hopefully give us a better estimate because
	// any ill-conditioning of the original matrix is transferred to U
	// and not to L. U[n-1,n-1] is an approximation to sigma_min(LU).
	bi.Dcopy(n-1, rhs, 1, xp, 1)
	xp[n-1] = rhs[n-1] + 1
	rhs[n-1] -= 1
	splus := 0.0
	sminu := 0.0
	for i := n - 1; i >= 0; i-- {
		temp = 1 / z[i*ldz+i]
		xp[i] *= temp
		rhs[i] *= temp
		for k := i + 1; k < n; k++ {
			xp[i] -= float64(xp[k] * (z[i*ldz+k] * temp))
			rhs[i] -= float64(rhs[k] * (z[i*ldz+k] * temp))
		}
		splus += math.Abs(xp[i])
		sminu += math.Abs(rhs[i])
	}
	if splus > sminu {
		bi.Dcopy(n, xp, 1, rhs, 1)
	}

	// Apply the permutations jpiv to the computed solution (rhs).
	impl.Dlaswp(1, rhs, 1, 0, n-2, jpiv[:n-1], -1)
	// Compute the sum of squares.
	rdscal, rdsum = impl.Dlassq(n, rhs, 1, rdscal, rdsum)
	return rdsum, rdscal
}
