// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dlag2er interface {
	Dlag2(a []float64, lda int, b []float64, ldb int) (scale1, scale2, wr1, wr2, wi float64)
}

func Dlag2Test(t *testing.T, impl Dlag2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, lda := range []int{2, 5} {
		for _, ldb := range []int{2, 5} {
			for i := 0; i <= 80; i++ { // 80 iterations for variability.
				dlag2Test(t, impl, rnd, lda, ldb)
			}
		}
	}
}

func dlag2Test(t *testing.T, impl Dlag2er, rnd *rand.Rand, lda, ldb int) {
	const tol = 1e-14
	a := randomGeneral(2, 2, lda, rnd)
	b := randomGeneral(2, 2, ldb, rnd)
	a11 := a.Data[0]
	a12 := a.Data[1]
	a21 := a.Data[lda]
	a22 := a.Data[lda+1]
	b11 := b.Data[0]
	b12 := b.Data[1]
	b22 := b.Data[ldb+1]
	name := fmt.Sprintf("lda=%d, ldb=%d", lda, ldb)
	scale1, scale2, wr1, wr2, wi := impl.Dlag2(a.Data, a.Stride, b.Data, b.Stride)
	if wi < 0 {
		t.Errorf("%v: wi can not be negative wi=%v", name, wi)
	}

	dget53Test(t, a, b, scale1, wr1, wi)
	dget53Test(t, a, b, scale2, wr2, wi)

	// Complex eigenvalue solution.
	if wi > 0 {
		if wr1 != wr2 {
			t.Errorf("%v: wr1 should equal wr2 for complex eigenvalues. got %v!=%v", name, wr1, wr2)
		}
		idet1 := cmplxdet2x2(complex(scale1*a11-wr1*b11, -wi*b11), complex(scale1*a12-wr1*b12, -wi*b12),
			complex(scale1*a21, 0), complex(scale1*a22-wr1*b22, -wi*b22))
		idet2 := cmplxdet2x2(complex(scale1*a11-wr1*b11, wi*b11), complex(scale1*a12-wr1*b12, wi*b12),
			complex(scale1*a21, 0), complex(scale1*a22-wr1*b22, wi*b22))

		if math.Abs(real(idet1))+math.Abs(imag(idet1)) > tol || math.Abs(real(idet2))+math.Abs(imag(idet2)) > tol {
			t.Errorf("%v: (%.4v±%.4vi)/%.4v did not solve eigenvalue problem", name, wr1, wi, scale1)
		}
		return // Complex solution verified
	}

	// Real eigenvalue solution.
	// |s*A - w*B| = 0   is solution for eigenvalue problem. Calculate distance from solution.
	res1 := det2x2(scale1*a11-wr1*b11, scale1*a12-wr1*b12,
		scale1*a21, scale1*a22-wr1*b22)
	if math.Abs(res1) > tol {
		t.Errorf("%v: got a residual from |%0.4v*A - w1*B| > tol: %v", name, scale1, res1)
	}
	res2 := det2x2(scale2*a11-wr2*b11, scale2*a12-wr2*b12,
		scale2*a21, scale2*a22-wr2*b22)
	if math.Abs(res2) > tol {
		t.Errorf("%v: got a residual from |s*A - w2*B| > tol: %v", name, res1)
	}
}

// dget53Test checks the generalized eigenvalues computed by dlag2.
// The basic test for an eigenvalue is:
//                           | det( s*A - w*B ) |
//  result =  ---------------------------------------------------
//            ulp max( s*norm(A), |w|*norm(B) )*norm( s*A - w*B )
//
// If s and w cant be scaled result will be 1/ulp.
func dget53Test(t *testing.T, a, b blas64.General, scale, wr, wi float64) (result float64) {
	lda := a.Stride
	ldb := b.Stride
	// Machine constants and norms.
	safmin := dlamchS
	ulp := dlamchE * dlamchB
	absw := math.Abs(wr) + math.Abs(wi)
	anorm := math.Max(math.Abs(a.Data[0])+math.Abs(a.Data[1]), math.Max(math.Abs(a.Data[lda])+math.Abs(a.Data[lda+1]), safmin))
	bnorm := math.Max(math.Abs(b.Data[0]), math.Max(math.Abs(b.Data[ldb])+math.Abs(b.Data[ldb+1]), safmin))

	// Check for possible overflow.
	temp := (safmin*bnorm)*absw + (safmin*anorm)*scale

	scales := scale
	wrs := wr
	wis := wi
	if temp >= 1 {
		t.Error("dget53: s*norm(A) + |w|*norm(B) > 1/safe_minimum")
		temp = 1 / temp
		scales *= temp
		wrs *= temp
		wis *= temp
		absw = math.Abs(wrs) + math.Abs(wis)
	}
	s1 := math.Max(ulp*math.Max(scales*anorm, absw*bnorm), safmin*math.Max(scales, absw))

	// Check for W and SCALE essentially zero.
	if s1 < safmin {
		t.Error("dget53: ulp*max( s*norm(A), |w|*norm(B) ) < safe_minimum")
		if scales < safmin && absw < safmin {
			t.Error("dget53: s and w could not be scaled so as to compute test")
			return 1 / ulp
		}
		// Scale up to avoid underflow.
		temp = 1 / math.Max(scales*anorm+absw*bnorm, safmin)
		scales *= temp
		wrs *= temp
		wis *= temp
		absw = math.Abs(wrs) + math.Abs(wis)
		s1 = math.Max(ulp*math.Max(scales*anorm, absw*bnorm),
			safmin*math.Max(scales, absw))
		if s1 < safmin {
			t.Error("dget53: s and w could not be scaled so as to compute test")
			return 1 / ulp
		}
	}

	// Compute C = s*A - w*B.
	cr11 := scales*a.Data[0] - wrs*b.Data[0]
	ci11 := -wis * b.Data[0]
	cr21 := scales * a.Data[1]
	cr12 := scales*a.Data[lda] - wrs*b.Data[ldb]
	ci12 := -wis * b.Data[ldb]
	cr22 := scales*a.Data[lda+1] - wrs*b.Data[ldb+1]
	ci22 := -wis * b.Data[ldb+1]

	// Compute the smallest singular value of s*A - w*B:
	//                 |det( s*A - w*B )|
	//     sigma_min = ------------------
	//                 norm( s*A - w*B )
	cnorm := math.Max(math.Abs(cr11)+math.Abs(ci11)+math.Abs(cr21),
		math.Max(math.Abs(cr12)+math.Abs(ci12)+math.Abs(cr22)+math.Abs(ci22), safmin))
	cscale := 1 / math.Sqrt(cnorm)
	detr := (cscale*cr11)*(cscale*cr22) -
		(cscale*ci11)*(cscale*ci22) -
		(cscale*cr12)*(cscale*cr21)
	deti := (cscale*cr11)*(cscale*ci22) +
		(cscale*ci11)*(cscale*cr22) -
		(cscale*ci12)*(cscale*cr21)
	sigmin := math.Abs(detr) + math.Abs(deti)
	result = sigmin / s1
	return result
}

// cmplxdet2x2 returns the determinant of
//  |a11 a12|
//  |a21 a22|
func cmplxdet2x2(a11, a12, a21, a22 complex128) complex128 {
	return a11*a22 - a12*a21
}
