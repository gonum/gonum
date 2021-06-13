// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
)

type Dlag2er interface {
	// DLAG2 computes the eigenvalues of a 2 x 2 generalized eigenvalue
	// problem
	//  A - w B,
	// with scaling as necessary to avoid over-/underflow.
	// The scaling factor "s" results in a modified eigenvalue equation
	//  s A - w B
	// where  s  is a non-negative scaling factor chosen so that  w,  w B,
	// and  s A  do not overflow and, if possible, do not underflow, either.
	//
	// B is a  ldb, 2)
	// On entry, the 2 x 2 upper triangular matrix B.  It is
	// assumed that the one-norm of B is less than 1/SAFMIN.  The
	// diagonals should be at least sqrt(SAFMIN) times the largest
	// element of B (in absolute value); if a diagonal is smaller
	// than that, then  +/- sqrt(SAFMIN) will be used instead of
	// that diagonal.
	//
	// Dlag2 is an internal routine. It is exported for testing purposes.
	Dlag2(a []float64, lda int, b []float64, ldb int) (scale1, scale2, wr1, wr2, wi float64)
}

func Dlag2Test(t *testing.T, impl Dlag2er) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		// Generate randomly the elements of a 2×2 matrix A
		//  [ a11 a12 ]
		//  [ a21 a22 ]

		a11, a12, a21, a22 := rnd.Float64(), rnd.Float64(), rnd.Float64(), rnd.Float64()
		// Generate randomly the elements of a 2×2 upper-triangular matrix B.
		//  [ b11 b12 ]
		//  [ 0   b22 ]
		b11, b12, b22 := rnd.Float64(), rnd.Float64(), rnd.Float64()
		scale1, scale2, wr1, wr2, wi := impl.Dlag2([]float64{a11, a12, a21, a22}, 2, []float64{b11, b12, 0, b22}, 2)
		if wi < 0 {
			t.Errorf("wi can not be negative wi=%v", wi)
		}
		// Complex eigenvalue solution.
		if wi > 0 {
			if wr1 != wr2 {
				t.Errorf("wr1 should equal wr2 for complex eigenvalues. got %v!=%v", wr1, wr2)
			}
			idet1 := cmplxdet2x2(complex(scale1*a11-wr1*b11, -wi*b11), complex(scale1*a12-wr1*b12, -wi*b12),
				complex(scale1*a21, 0), complex(scale1*a22-wr1*b22, -wi*b22))
			idet2 := cmplxdet2x2(complex(scale1*a11-wr1*b11, wi*b11), complex(scale1*a12-wr1*b12, wi*b12),
				complex(scale1*a21, 0), complex(scale1*a22-wr1*b22, wi*b22))

			if math.Abs(real(idet1))+math.Abs(imag(idet1)) > tol || math.Abs(real(idet2))+math.Abs(imag(idet2)) > tol {
				t.Errorf(`I thought imag eigenvalue %v (scaled %v) solved 
			[%.6v  %.6v]    [%.6v  %.6v]
			[%.6v  %.6v] - w[    0   %.6v]`, complex(wr1, wi), scale1, a11, a12, b11, b12, a21, a22, b22)
			}
			continue // End complex branch.
		}

		// Real eigenvalue solution.
		cerr := 0
		// |s A - w B| = 0   is solution for eigenvalue problem. Calculate distance from solution.
		res1 := det2x2(scale1*a11-wr1*b11, scale1*a12-wr1*b12,
			scale1*a21, scale1*a22-wr1*b22)
		if math.Abs(res1) > tol {
			t.Errorf("got a residual from |%0.2g*A - w1*B| > tol: %v", scale1, res1)
			cerr++
		}
		res2 := det2x2(scale2*a11-wr2*b11, scale2*a12-wr2*b12,
			scale2*a21, scale2*a22-wr2*b22)
		if math.Abs(res2) > tol {
			t.Errorf("got a residual from |s*A - w2*B| > tol: %v", res1)
			cerr++
		}
		if cerr > 0 {
			t.Errorf(`I thought the eigenvalues %v, %v (scaled %v, %v) solved 
[%.6v  %.6v]    [%.6v  %.6v]
[%.6v  %.6v] - w[    0   %.6v]`, wr1, wr2, scale1, scale2, a11, a12, b11, b12, a21, a22, b22)
		}
	}
}

// Solves determinant of
//  |a11 a12|
//  |a21 a22|
func cmplxdet2x2(a11, a12, a21, a22 complex128) complex128 { return a11*a22 - a12*a21 }
