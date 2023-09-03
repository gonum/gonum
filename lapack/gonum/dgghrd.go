// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dgghrd reduces a pair of real matrices (A,B) to generalized upper
// Hessenberg form using orthogonal transformations, where A is a
// general matrix and B is upper triangular.  The form of the
// generalized eigenvalue problem is
//
//	A*x = lambda*B*x,
//
// and B is typically made upper triangular by computing its QR
// factorization and moving the orthogonal matrix Q to the left side
// of the equation.
// This subroutine simultaneously reduces A to a Hessenberg matrix H:
//
//	Qᵀ*A*Z = H
//
// and transforms B to another upper triangular matrix T:
//
//	Qᵀ*B*Z = T
//
// in order to reduce the problem to its standard form
//
//	H*y = lambda*T*y
//
// where y = Zᵀ*x.
//
// The orthogonal matrices Q and Z are determined as products of Givens
// rotations.  They may either be formed explicitly, or they may be
// postmultiplied into input matrices Q1 and Z1, so that
//
//	Q1 * A * Z1ᵀ = (Q1*Q) * H * (Z1*Z)ᵀ
//	Q1 * B * Z1ᵀ = (Q1*Q) * T * (Z1*Z)ᵀ
//
// If Q1 is the orthogonal matrix from the QR factorization of B in the
// original equation A*x = lambda*B*x, then Dgghrd reduces the original
// problem to generalized Hessenberg form.
//
// Dgghrd is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dgghrd(compq, compz lapack.OrthoComp, n, ilo, ihi int, a []float64, lda int, b []float64, ldb int, q []float64, ldq int, z []float64, ldz int) {
	var ilq bool
	icompq := 0
	switch compq {
	case lapack.OrthoNone:
		ilq = false
		icompq = 1
	case lapack.OrthoEntry:
		ilq = true
		icompq = 2
	case lapack.OrthoUnit:
		ilq = true
		icompq = 3
	default:
		panic(badOrthoComp)
	}

	var ilz bool
	icompz := 0
	switch compz {
	case lapack.OrthoNone:
		ilz = false
		icompz = 1
	case lapack.OrthoEntry:
		ilz = true
		icompz = 2
	case lapack.OrthoUnit:
		ilz = true
		icompz = 3
	default:
		panic(badOrthoComp)
	}

	switch {
	case n < 0:
		panic(nLT0)
	case ilo < 0:
		panic(badIlo)
	case ihi < ilo-1 || ihi >= n:
		panic(badIhi)
	case lda < max(1, n):
		panic(badLdA)
	case ldb < max(1, n):
		panic(badLdB)
	case (ilq && ldq < n) || ldq < 1:
		panic(badLdQ)
	case (ilz && ldz < n) || ldz < 1:
		panic(badLdZ)
	}

	if icompq == 3 {
		impl.Dlaset(blas.All, n, n, 0, 1, q, ldq)
	}
	if icompz == 3 {
		impl.Dlaset(blas.All, n, n, 0, 1, z, ldz)
	}
	if n < 1 {
		return // Quick return if possible.
	}

	// Zero out lower triangle of B.
	for jcol := 0; jcol < n-1; jcol++ {
		for jrow := jcol + 1; jrow < n; jrow++ {
			b[jrow*ldb+jcol] = 0
		}
	}
	bi := blas64.Implementation()
	// Reduce A and B.
	for jcol := ilo; jcol <= ihi-2; jcol++ {
		for jrow := ihi; jrow >= jcol+2; jrow-- {
			// Step 1: rotate rows JROW-1, JROW to kill A(JROW,JCOL).
			temp := a[(jrow-1)*lda+jcol]
			var c, s float64
			c, s, a[(jrow-1)*lda+jcol] = impl.Dlartg(temp, a[jrow*lda+jcol])
			a[jrow*lda+jcol] = 0
			bi.Drot(n-jcol-1, a[(jrow-1)*lda+jcol+1:], 1,
				a[jrow*lda+jcol+1:], 1, c, s)

			bi.Drot(n+2-jrow-1, b[(jrow-1)*ldb+jrow-1:], 1,
				b[jrow*ldb+jrow-1:], 1, c, s)

			if ilq {
				bi.Drot(n, q[jrow-1:], ldq, q[jrow:], ldq, c, s)
			}

			// Step 2: rotate columns JROW, JROW-1 to kill B(JROW,JROW-1).
			temp = b[jrow*ldb+jrow]
			c, s, b[jrow*ldb+jrow] = impl.Dlartg(temp, b[jrow*ldb+jrow-1])
			b[jrow*ldb+jrow-1] = 0

			bi.Drot(ihi+1, a[jrow:], lda, a[jrow-1:], lda, c, s)
			bi.Drot(jrow, b[jrow:], ldb, b[jrow-1:], ldb, c, s)

			if ilz {
				bi.Drot(n, z[jrow:], ldz, z[jrow-1:], ldz, c, s)
			}
		}
	}
}
