package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dlatdf uses the LU factorization of the n×n matrix Z computed by
// Dgetc2 and computes a contribution to the reciprocal Dif-estimate
// by solving
//  Z * x = b for x
// and choosing the r.h.s. b such that
// the norm of x is as large as possible. On entry RHS = b holds the
// contribution from earlier solved sub-systems, and on return RHS = x.
//
// The factorization of Z returned by Dgetc2 has the form
//  Z = P*L*U*Q,
// where P and Q are permutation matrices. L is lower triangular with
// unit diagonal elements and U is upper triangular.
//
// On entry ipiv and jpiv are n length slices. On exit the pivot indices
// where row i has been interchanged with ipiv[i] and column j interchanged
// with column jpiv[j].
//
// Dlatdf is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlatdf(ijob, n int, z []float64, ldz int, rhs []float64, rdsum, rdscal float64, ipiv, jpiv []int) (sum, scale float64) {
	// ijob info:
	// IJOB = 2: First compute an approximative null-vector e
	// of Z using DGECON, e is normalized and solve for
	// Zx = +-e - f with the sign giving the greater value
	// of 2-norm(x). About 5 times as expensive as Default.
	// IJOB .ne. 2: Local look ahead strategy where all entries of
	// the r.h.s. b is chosen as either +1 or -1 (Default).
	//
	// rdsum info:
	// On entry, the sum of squares of computed contributions to
	// the Dif-estimate under computation by DTGSYL, where the
	// scaling factor RDSCAL (see below) has been factored out.
	// On exit, the corresponding sum of squares updated with the
	// contributions from the current sub-system.
	// If TRANS = 'T' RDSUM is not touched.
	// NOTE: RDSUM only makes sense when DTGSY2 is called by STGSYL.
	//
	// rdscal info:
	// On entry, scaling factor used to prevent overflow in RDSUM.
	// On exit, RDSCAL is updated w.r.t. the current contributions
	// in RDSUM.
	// If TRANS = 'T', RDSCAL is not touched.
	// NOTE: RDSCAL only makes sense when DTGSY2 is called by
	// DTGSYL.
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
	}

	bi := blas64.Implementation()
	var temp float64
	xp := make([]float64, n)
	// Compute approximate nullvector xm of Z.
	if ijob == 2 {
		xm := make([]float64, n)
		work := make([]float64, 4*n)
		impl.Dgecon(lapack.MaxRowSum, n, z, ldz, 1.0, work, make([]int, n))
		bi.Dcopy(n, work[n:], 1, xm, 1)

		// Compute rhs.
		impl.Dlaswp(1, xm, n, 0, n-2, ipiv, -1)
		temp = 1.0 / math.Sqrt(bi.Ddot(n, xm, 1, xm, 1))
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

	// If ijob != 2 uses look-ahead strategy.
	impl.Dlaswp(1, rhs, 1, 0, n-1, ipiv, 1)
	// Solve for L-part choosing rhs either to +1 or -1.
	var bp, bm, splus, sminu float64
	pmone := -1.0
	for j := 0; j < n-2; j++ {
		bp = rhs[j] + 1.0
		bm = rhs[j] - 1.0
		// Look-ahead for L-part rhs[0:n-2] = + or -1, splus and
		// smin computed more efficiently than in Bsolve [1].
		splus = 2.0 + bi.Ddot(n-j-1, z[(j+1)*ldz+j:], ldz, z[(j+1)*ldz+j:], ldz)
		sminu = bi.Ddot(n-j-1, z[(j+1)*ldz+j:], ldz, rhs[j+1:], 1)
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
	xp[n-1] = rhs[n-1] + 1.0
	rhs[n-1] -= 1.0
	splus = 0
	sminu = 0
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
	impl.Dlaswp(1, rhs, 1, 0, n-1, jpiv, -1)
	// Compute the sum of squares.
	rdscal, rdsum = impl.Dlassq(n, rhs, 1, rdscal, rdsum)
	return rdsum, rdscal
}
