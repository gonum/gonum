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
	"gonum.org/v1/gonum/floats"
)

type Dlag2er interface {
	Dlag2(a []float64, lda int, b []float64, ldb int) (scale1, scale2, wr1, wr2, wi float64)
}

func Dlag2Test(t *testing.T, impl Dlag2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, lda := range []int{2, 5} {
		for _, ldb := range []int{2, 5} {
			for aKind := 0; aKind <= 20; aKind++ {
				for bKind := 0; bKind <= 20; bKind++ {
					dlag2Test(t, impl, rnd, lda, ldb, aKind, bKind)
				}
			}
		}
	}
}

func dlag2Test(t *testing.T, impl Dlag2er, rnd *rand.Rand, lda, ldb int, aKind, bKind int) {
	const tol = 1e-14

	a := makeDlag2TestMatrix(rnd, lda, aKind)
	b := makeDlag2TestMatrix(rnd, ldb, bKind)

	aCopy := cloneGeneral(a)
	bCopy := cloneGeneral(b)

	scale1, scale2, wr1, wr2, wi := impl.Dlag2(a.Data, a.Stride, b.Data, b.Stride)

	name := fmt.Sprintf("lda=%d,ldb=%d,aKind=%d,bKind=%d", lda, ldb, aKind, bKind)
	aStr := fmt.Sprintf("A = [%g,%g]\n    [%g,%g]", a.Data[0], a.Data[1], a.Data[a.Stride], a.Data[a.Stride+1])
	bStr := fmt.Sprintf("B = [%g,%g]\n    [%g,%g]", b.Data[0], b.Data[1], 0.0, b.Data[b.Stride+1])

	if !floats.Same(a.Data, aCopy.Data) {
		t.Errorf("%s: unexpected modification of a", name)
	}
	if !floats.Same(b.Data, bCopy.Data) {
		t.Errorf("%s: unexpected modification of b", name)
	}

	if wi < 0 {
		t.Fatalf("%s: wi is negative; wi=%g,\n%s\n%s", name, wi, aStr, bStr)
		return
	}

	if wi > 0 {
		if wr1 != wr2 {
			t.Fatalf("%s: complex eigenvalue but wr1 != wr2; wr1=%g, wr2=%g,\n%s\n%s", name, wr1, wr2, aStr, bStr)
			return
		}
		if scale1 != scale2 {
			t.Fatalf("%s: complex eigenvalue but scale1 != scale2; scale1=%g, scale2=%g,\n%s\n%s", name, scale1, scale2, aStr, bStr)
			return
		}
	}

	resid, err := residualDlag2(a, b, scale1, complex(wr1, wi))
	if err != nil {
		t.Logf("%s: invalid input data: %v\n%s\n%s", name, err, aStr, bStr)
		return
	}
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%s: unexpected first eigenvalue %g with s=%g; resid=%g, want<=%g\n%s\n%s", name, complex(wr1, wi), scale1, resid, tol, aStr, bStr)
	}

	resid, err = residualDlag2(a, b, scale2, complex(wr2, -wi))
	if err != nil {
		t.Logf("%s: invalid input data: %s\n%s\n%s", name, err, aStr, bStr)
		return
	}
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%s: unexpected second eigenvalue %g with s=%g; resid=%g, want<=%g\n%s\n%s", name, complex(wr2, -wi), scale2, resid, tol, aStr, bStr)
	}
}

func makeDlag2TestMatrix(rnd *rand.Rand, ld, kind int) blas64.General {
	a := zeros(2, 2, ld)
	switch kind {
	case 0:
		// Zero matrix.
	case 1:
		// Identity.
		a.Data[0] = 1
		a.Data[a.Stride+1] = 1
	case 2:
		// Large diagonal.
		a.Data[0] = 2 * safmax
		a.Data[a.Stride+1] = 2 * safmax
	case 3:
		// Tiny diagonal.
		a.Data[0] = safmin
		a.Data[a.Stride+1] = safmin
	case 4:
		// Tiny and large diagonal.
		a.Data[0] = safmin
		a.Data[a.Stride+1] = safmax
	case 5:
		// Large and tiny diagonal.
		a.Data[0] = safmax
		a.Data[a.Stride+1] = safmin
	case 6:
		// Large complex eigenvalue.
		a.Data[0] = safmax
		a.Data[1] = safmax
		a.Data[a.Stride] = -safmax
		a.Data[a.Stride+1] = safmax
	case 7:
		// Tiny complex eigenvalue.
		a.Data[0] = safmin
		a.Data[1] = safmin
		a.Data[a.Stride] = -safmin
		a.Data[a.Stride+1] = safmin
	case 8:
		// Random matrix with large elements.
		a.Data[0] = safmax * (2*rnd.Float64() - 1)
		a.Data[1] = safmax * (2*rnd.Float64() - 1)
		a.Data[a.Stride] = safmax * (2*rnd.Float64() - 1)
		a.Data[a.Stride+1] = safmax * (2*rnd.Float64() - 1)
	case 9:
		// Random matrix with tiny elements.
		a.Data[0] = safmin * (2*rnd.Float64() - 1)
		a.Data[1] = safmin * (2*rnd.Float64() - 1)
		a.Data[a.Stride] = safmin * (2*rnd.Float64() - 1)
		a.Data[a.Stride+1] = safmin * (2*rnd.Float64() - 1)
	default:
		// Random matrix.
		a = randomGeneral(2, 2, ld, rnd)
	}
	return a
}

// residualDlag2 returns the value of
//
//	           | det( s*A - w*B ) |
//	-------------------------------------------
//	max(s*norm(A), |w|*norm(B))*norm(s*A - w*B)
//
// that can be used to check the generalized eigenvalues computed by Dlag2 and
// an error that indicates invalid input data.
func residualDlag2(a, b blas64.General, s float64, w complex128) (float64, error) {
	const ulp = dlamchP

	a11, a12 := a.Data[0], a.Data[1]
	a21, a22 := a.Data[a.Stride], a.Data[a.Stride+1]

	b11, b12 := b.Data[0], b.Data[1]
	b22 := b.Data[b.Stride+1]

	// Compute norms.
	absw := zabs(w)
	anorm := math.Max(math.Abs(a11)+math.Abs(a21), math.Abs(a12)+math.Abs(a22))
	anorm = math.Max(anorm, safmin)
	bnorm := math.Max(math.Abs(b11), math.Abs(b12)+math.Abs(b22))
	bnorm = math.Max(bnorm, safmin)

	// Check for possible overflow.
	temp := (safmin*anorm)*s + (safmin*bnorm)*absw
	if temp >= 1 {
		// Scale down to avoid overflow.
		s /= temp
		w = scale(1/temp, w)
		absw = zabs(w)
	}

	// Check for w and s essentially zero.
	s1 := math.Max(ulp*math.Max(s*anorm, absw*bnorm), safmin*math.Max(s, absw))
	if s1 < safmin {
		if s < safmin && absw < safmin {
			return 1 / ulp, fmt.Errorf("ulp*max(s*|A|,|w|*|B|) < safmin and s and w could not be scaled; s=%g, |w|=%g", s, absw)
		}
		// Scale up to avoid underflow.
		temp = 1 / math.Max(s*anorm+absw*bnorm, safmin)
		s *= temp
		w = scale(temp, w)
		absw = zabs(w)
		s1 = math.Max(ulp*math.Max(s*anorm, absw*bnorm), safmin*math.Max(s, absw))
		if s1 < safmin {
			return 1 / ulp, fmt.Errorf("ulp*max(s*|A|,|w|*|B|) < safmin and s and w could not be scaled; s=%g, |w|=%g", s, absw)
		}
	}

	// Compute C = s*A - w*B.
	c11 := complex(s*a11, 0) - w*complex(b11, 0)
	c12 := complex(s*a12, 0) - w*complex(b12, 0)
	c21 := complex(s*a21, 0)
	c22 := complex(s*a22, 0) - w*complex(b22, 0)
	// Compute norm(s*A - w*B).
	cnorm := math.Max(zabs(c11)+zabs(c21), zabs(c12)+zabs(c22))
	// Compute det(s*A - w*B)/norm(s*A - w*B).
	cs := 1 / math.Sqrt(math.Max(cnorm, safmin))
	det := cmplxdet2x2(scale(cs, c11), scale(cs, c12), scale(cs, c21), scale(cs, c22))
	// Compute |det(s*A - w*B)|/(norm(s*A - w*B)*max(s*norm(A), |w|*norm(B))).
	return zabs(det) / s1 * ulp, nil
}

func zabs(z complex128) float64 {
	return math.Abs(real(z)) + math.Abs(imag(z))
}

// scale scales the complex number c by f.
func scale(f float64, c complex128) complex128 {
	return complex(f*real(c), f*imag(c))
}

// cmplxdet2x2 returns the determinant of
//
//	|a11 a12|
//	|a21 a22|
func cmplxdet2x2(a11, a12, a21, a22 complex128) complex128 {
	return a11*a22 - a12*a21
}
