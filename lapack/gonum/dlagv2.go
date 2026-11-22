// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// Dlagv2 computes the Generalized Schur factorization of a real 2×2
// matrix pencil (A,B) where B is upper triangular. This routine
// computes orthogonal (rotation) matrices given by csl, SNL and CSR,
// SNR such that
//
//  1. if the pencil (A,B) has two real eigenvalues (include 0/0 or 1/0 types), then
//     [ a11 a12 ] = [  csl  snl ] [ a11 a12 ] [ csr -snr ]
//     [  0  a22 ]   [ -snl  csl ] [ a21 a22 ] [ snr  csr ]
//
//     [ b11 b12 ] = [  csl  snl ] [ b11 b12 ] [ csr -snr ]
//     [  0  b22 ]   [ -snl  csl ] [  0  b22 ] [ snr  csr ],
//
//  2. if the pencil (A,B) has a pair of complex conjugate eigenvalues, then
//     [ a11 a12 ] = [  csl  snl ] [ a11 a12 ] [ csr -snr ]
//     [ a21 a22 ]   [ -snl  csl ] [ a21 a22 ] [ snr  csr ]
//
//     [ b11  0  ] = [  csl  snl ] [ b11 b12 ] [ csr -snr ]
//     [  0  b22 ]   [ -snl  csl ] [  0  b22 ] [ snr  csr ]
//
// where b11 >= b22 > 0.
//
// On exit, A is overwritten by the A-part of the generalized Schur form
// and B is overwritten by the B-part of the generalized Schur form.
//
// (alphar(k)+imag*alphai(k))/beta(k) are the eigenvalues of the
// pencil (A,B), k=0,1 where imag = sqrt(-1). Note that beta(k) may
// be zero. They are solely outputs and will be overwritten.
//
// Returned floats are cosines (cs) and sines (sn) of left (l) and right (r)
// rotation matrices.
//
// Dlagv2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlagv2(a []float64, lda int, b []float64, ldb int, alphar, alphai, beta []float64) (csl, snl, csr, snr float64) {
	switch {
	case lda < 2:
		panic(badLdA)
	case ldb < 2:
		panic(badLdB)
	case len(a) < 4:
		panic(shortA)
	case len(b) < 4:
		panic(shortB)
	case len(alphar) < 2 || len(alphai) < 2:
		panic(badLenAlpha)
	case len(beta) < 2:
		panic(badLenBeta)
		// Documentation in reference requires B diagonal to be non-negative and
		// non-zero and in descending order but it seems to not be strictly
		// necessary looking at code.
	}

	const (
		safmin = dlamchS
		ulp    = dlamchP
	)

	// Scale A.
	anorm := math.Max(math.Abs(a[0])+math.Abs(a[lda]),
		math.Abs(a[1])+math.Abs(a[lda+1]))
	anorm = math.Max(anorm, safmin)
	ascale := 1 / anorm
	a[0] *= ascale
	a[1] *= ascale
	a[lda] *= ascale
	a[lda+1] *= ascale

	// Scale B.
	bnorm := math.Max(math.Abs(b[0]),
		math.Abs(b[1])+math.Abs(b[ldb+1]))
	bnorm = math.Max(bnorm, safmin)
	bscale := 1 / bnorm
	b[0] *= bscale
	b[1] *= bscale
	b[ldb+1] *= bscale

	bi := blas64.Implementation()
	var wi, wr1, scale1 float64
	switch {
	// Check if A can be deflated.
	case math.Abs(a[lda]) <= ulp:
		csl, csr = 1, 1
		a[lda] = 0
		b[ldb] = 0

	// Check if B is singular.
	case math.Abs(b[0]) <= ulp:
		csl, snl, _ = impl.Dlartg(a[0], a[lda])
		csr = 1
		bi.Drot(2, a, 1, a[lda:], 1, csl, snl)
		bi.Drot(2, b, 1, b[ldb:], 1, csl, snl)
		a[lda] = 0
		b[0] = 0
		b[ldb] = 0

	case math.Abs(b[ldb+1]) <= ulp:
		csr, snr, _ = impl.Dlartg(a[lda+1], a[lda+0])
		snr *= -1
		bi.Drot(2, a, lda, a[1:], lda, csr, snr)
		bi.Drot(2, b, ldb, b[1:], ldb, csr, snr)
		csl = 1
		a[lda] = 0
		b[ldb] = 0
		b[ldb+1] = 0

	// B is nonsingular, first compute the eigenvalues of (A,B).
	default:
		scale1, _, wr1, _, wi = impl.Dlag2(a, lda, b, ldb)

		// Ensure upper triangular form before rotations
		// as Drot will access this element and may propagate NaNs.
		// Reference does not contain this zeroing.
		b[ldb] = 0

		if wi != 0 {
			// A pair of complex conjugate eigenvalues:
			// first compute the SVD of the matrix B.
			_, _, snr, csr, snl, csl = impl.Dlasv2(b[0], b[1], b[ldb+1])

			// Form (A,B) := Q(A,B)Zᵀ where Q is left rotation matrix and
			// Z is right rotation matrix computed from Dlasv2
			bi.Drot(2, a, 1, a[lda:], 1, csl, snl)
			bi.Drot(2, b, 1, b[ldb:], 1, csl, snl)

			bi.Drot(2, a, lda, a[1:], lda, csr, snr)
			bi.Drot(2, b, ldb, b[1:], ldb, csr, snr)

			b[ldb] = 0
			b[1] = 0
			break // Exit from switch.
		}

		// Got real eigenvalues from Dlag2, compute s*A-w*B.
		h1 := scale1*a[0] - wr1*b[0]
		h2 := scale1*a[1] - wr1*b[1]
		h3 := scale1*a[lda+1] - wr1*b[ldb+1]

		rr := impl.Dlapy2(h1, h2)
		qq := impl.Dlapy2(scale1*a[lda], h3)
		if rr > qq {
			// Find right rotation matrix to zero 0,0 element of (sA - wB).
			csr, snr, _ = impl.Dlartg(h2, h1)
		} else {
			// Find right rotation matrix to zero 1,0 element of (sA - wB).
			csr, snr, _ = impl.Dlartg(h3, scale1*a[lda])
		}
		snr *= -1
		bi.Drot(2, a, lda, a[1:], lda, csr, snr)
		bi.Drot(2, b, ldb, b[1:], ldb, csr, snr)

		// Compute inf. norms of A and B.
		h1 = math.Max(math.Abs(a[0])+math.Abs(a[1]),
			math.Abs(a[lda])+math.Abs(a[lda+1]))

		h2 = math.Max(math.Abs(b[0])+math.Abs(b[1]),
			math.Abs(b[ldb])+math.Abs(b[ldb+1])) // b[ldb] may be non-zero after rotations.

		if scale1*h1 >= math.Abs(wr1)*h2 {
			// Find left rotation matrix Q to zero out B[1,0].
			csl, snl, _ = impl.Dlartg(b[0], b[ldb])
		} else {
			// Find left rotation matrix Q to zero out A[1,0]
			csl, snl, _ = impl.Dlartg(a[0], a[lda])
		}

		bi.Drot(2, a, 1, a[lda:], 1, csl, snl)
		bi.Drot(2, b, 1, b[ldb:], 1, csl, snl)

		a[lda] = 0
		b[ldb] = 0
	}

	// Unscaling.

	a[0] *= anorm
	a[1] *= anorm
	a[lda] *= anorm
	a[lda+1] *= anorm

	b[0] *= bnorm
	b[1] *= bnorm
	b[ldb] *= bnorm
	b[ldb+1] *= bnorm

	if wi == 0 {
		alphar[0] = a[0]
		alphar[1] = a[lda+1]
		alphai[0] = 0
		alphai[1] = 0
		beta[0] = b[0]
		beta[1] = b[ldb+1]
	} else {
		alphar[0] = anorm * wr1 / scale1 / bnorm
		alphai[0] = anorm * wi / scale1 / bnorm
		alphar[1] = alphar[0]
		alphai[1] = -alphai[0]
		beta[0] = 1
		beta[1] = 1
	}

	return csl, snl, csr, snr
}
