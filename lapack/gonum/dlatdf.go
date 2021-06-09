package gonum

import "gonum.org/v1/gonum/blas/blas64"

// Dlatdf uses the LU factorization of the n-by-n matrix Z computed by
// Dgetc2 and computes a contribution to the reciprocal Dif-estimate
// by solving
//  Z * x = b for x
// and choosing the r.h.s. b such that
// the norm of x is as large as possible. On entry RHS = b holds the
// contribution from earlier solved sub-systems, and on return RHS = x.

// The factorization of Z returned by Dgetc2 has the form
//  Z = P*L*U*Q,
// where P and Q are permutation matrices. L is lower triangular with
// unit diagonal elements and U is upper triangular.
func (Implementation) Dlatdf(ijob, n int, z []float64, ldz int, rhs []float64, rdsum, rdscal float64, ipiv, jpiv []int) (sumout, scalout float64) {
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
		panic(badLdA)
	}

	// Compute approximate nullvector XM of Z.
	bi := blas64.Implementation()
	if ijob == 2 {

		Dgecon
	}
}
