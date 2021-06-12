// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "math"

// DLAG2 computes the eigenvalues of a 2 x 2 generalized eigenvalue
// problem
//  A - w B,
// with scaling as necessary to avoid over-/underflow.
// The scaling factor "s" results in a modified eigenvalue equation
//  s A - w B
// where  s  is a non-negative scaling factor chosen so that  w,  w B,
// and  s A  do not overflow and, if possible, do not underflow, either.
//
// Dlag2 is an internal routine. It is exported for testing purposes.
func (Implementation) Dlag2(a []float64, lda int, b []float64, ldb int, safmin float64) (scale1, scale2, wr1, wr2, wi float64) {
	rtmin := math.Sqrt(safmin)
	rtmax := 1. / rtmin
	safmax := 1. / safmin
	mmax := math.Max
	fabs := math.Abs
	sign := math.Copysign
	// Scale a.
	anorm := mmax(fabs(a[0*lda+0])+fabs(a[1*lda+0]),
		fabs(a[0*lda+1])+fabs(a[1*lda+1]))
	ascale := 1. / anorm
	a11 := ascale * a[0*lda+0]
	a21 := ascale * a[1*lda+0]
	a12 := ascale * a[0*lda+1]
	a22 := ascale * a[1*lda+1]

	// Perturb B if necessary to insure non-singularity.
	b11 := b[0*ldb+0]
	// b21 := b[1*ldb+0]
	b12 := b[0*ldb+1]
	b22 := b[1*ldb+1]

	bmin := rtmin * mmax(mmax(fabs(b11), fabs(b12)), mmax(fabs(b22), fabs(rtmin)))
	if fabs(b11) < bmin {
		b11 = sign(bmin, b11)
	}
	if fabs(b22) < bmin {
		b22 = sign(bmin, b22)
	}

	// Scale B.
	bnorm := mmax(mmax(fabs(b11), fabs(b12)+fabs(b22)), safmin)
	bsize := mmax(fabs(b11), fabs(b22))
	bscale := 1. / bsize
	b11 = bscale * b11
	b12 = bscale * b12
	b22 = bscale * b22

	// Compute larger eigenvalue by method described by C. van Loan.
}
