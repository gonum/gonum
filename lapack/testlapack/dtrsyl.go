// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dtrsyler interface {
	Dtrsyl(trana, tranb blas.Transpose, isgn, m, n int, a []float64, lda int,
		b []float64, ldb int, c []float64, ldc int) (scale float64, ok bool)
}

func DtrsylTest(t *testing.T, impl Dtrsyler) {
	rnd := rand.New(rand.NewPCG(1, 1))

	for _, m := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
			for _, trana := range []blas.Transpose{blas.NoTrans, blas.Trans} {
				for _, tranb := range []blas.Transpose{blas.NoTrans, blas.Trans} {
					for _, isgn := range []int{1, -1} {
						for _, extra := range []int{0, 5} {
							testDtrsyl(t, impl, m, n, trana, tranb, isgn, extra, rnd)
						}
					}
				}
			}
		}
	}

	// Test with special matrices.
	testDtrsylSpecial(t, impl)
}

func testDtrsyl(t *testing.T, impl Dtrsyler, m, n int, trana, tranb blas.Transpose, isgn, extra int, rnd *rand.Rand) {
	const tol = 1e-12

	prefix := fmt.Sprintf("m=%v, n=%v, trana=%c, tranb=%c, isgn=%d, extra=%v",
		m, n, trana, tranb, isgn, extra)

	// Generate random quasi-triangular matrices A (m×m) and B (n×n).
	a := randomSchur(m, m+extra, rnd)
	b := randomSchur(n, n+extra, rnd)

	// Generate random RHS matrix C (m×n).
	c := randomGeneral(m, n, n+extra, rnd)
	cOrig := cloneGeneral(c)

	// Also keep originals of A and B for residual check.
	aOrig := cloneGeneral(a)
	bOrig := cloneGeneral(b)

	// Call Dtrsyl.
	scale, ok := impl.Dtrsyl(trana, tranb, isgn, m, n, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride)

	// Check that A and B are not modified.
	if !generalEqual(a, aOrig) {
		t.Errorf("%v: A was modified", prefix)
	}
	if !generalEqual(b, bOrig) {
		t.Errorf("%v: B was modified", prefix)
	}

	// Quick return for empty matrices.
	if m == 0 || n == 0 {
		if scale != 1 {
			t.Errorf("%v: scale=%v, want 1 for empty matrix", prefix, scale)
		}
		if !ok {
			t.Errorf("%v: ok=false, want true for empty matrix", prefix)
		}
		return
	}

	// Verify the solution by computing the residual:
	// R = op(A)*X + isgn*X*op(B) - scale*C
	// where X is the solution stored in c.
	resid := dtrsylResidual(aOrig, bOrig, cOrig, c, trana, tranb, isgn, scale)
	anorm := dlange(lapack.MaxColumnSum, m, m, aOrig.Data, aOrig.Stride)
	bnorm := dlange(lapack.MaxColumnSum, n, n, bOrig.Data, bOrig.Stride)
	xnorm := dlange(lapack.MaxColumnSum, m, n, c.Data, c.Stride)
	cnorm := dlange(lapack.MaxColumnSum, m, n, cOrig.Data, cOrig.Stride)

	denom := anorm*xnorm + xnorm*bnorm + cnorm
	if denom == 0 {
		denom = 1
	}
	normRes := resid / (denom * float64(max(m, n)) * dlamchE)

	if normRes > tol/dlamchE {
		t.Errorf("%v: residual too large: %v, want <= %v", prefix, normRes, tol/dlamchE)
	}

	_ = ok // ok can be false if eigenvalues are close, but solution should still be acceptable.
}

func testDtrsylSpecial(t *testing.T, impl Dtrsyler) {
	const tol = 1e-10

	// Test case: diagonal matrices (all 1x1 blocks).
	for _, n := range []int{1, 2, 5, 10} {
		a := eye(n, n)
		for i := 0; i < n; i++ {
			a.Data[i*a.Stride+i] = float64(i + 1)
		}
		b := eye(n, n)
		for i := 0; i < n; i++ {
			b.Data[i*b.Stride+i] = float64(i + 10)
		}
		c := ones(n, n)
		cOrig := cloneGeneral(c)

		scale, _ := impl.Dtrsyl(blas.NoTrans, blas.NoTrans, 1, n, n, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride)

		resid := dtrsylResidual(a, b, cOrig, c, blas.NoTrans, blas.NoTrans, 1, scale)
		anorm := dlange(lapack.MaxColumnSum, n, n, a.Data, a.Stride)
		bnorm := dlange(lapack.MaxColumnSum, n, n, b.Data, b.Stride)
		xnorm := dlange(lapack.MaxColumnSum, n, n, c.Data, c.Stride)
		cnorm := dlange(lapack.MaxColumnSum, n, n, cOrig.Data, cOrig.Stride)

		denom := anorm*xnorm + xnorm*bnorm + cnorm
		if denom == 0 {
			denom = 1
		}
		normRes := resid / (denom * float64(n) * dlamchE)
		if normRes > tol/dlamchE {
			t.Errorf("Diagonal n=%d: residual too large: %v", n, normRes)
		}
	}

	// Test case: 2x2 blocks (complex eigenvalues).
	// A has one 2x2 block representing complex eigenvalue 1±2i.
	a := blas64.General{
		Rows:   2,
		Cols:   2,
		Stride: 2,
		Data:   []float64{1, 2, -2, 1},
	}
	b := blas64.General{
		Rows:   2,
		Cols:   2,
		Stride: 2,
		Data:   []float64{3, 1, -1, 3},
	}
	c := ones(2, 2)
	cOrig := cloneGeneral(c)

	scale, _ := impl.Dtrsyl(blas.NoTrans, blas.NoTrans, 1, 2, 2, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride)

	resid := dtrsylResidual(a, b, cOrig, c, blas.NoTrans, blas.NoTrans, 1, scale)
	anorm := dlange(lapack.MaxColumnSum, 2, 2, a.Data, a.Stride)
	bnorm := dlange(lapack.MaxColumnSum, 2, 2, b.Data, b.Stride)
	xnorm := dlange(lapack.MaxColumnSum, 2, 2, c.Data, c.Stride)
	cnorm := dlange(lapack.MaxColumnSum, 2, 2, cOrig.Data, cOrig.Stride)

	denom := anorm*xnorm + xnorm*bnorm + cnorm
	if denom == 0 {
		denom = 1
	}
	normRes := resid / (denom * 2 * dlamchE)
	if normRes > tol/dlamchE {
		t.Errorf("2x2 block: residual too large: %v", normRes)
	}
}

// dtrsylResidual computes ||op(A)*X + isgn*X*op(B) - scale*C||_1.
func dtrsylResidual(A, B, C, X blas64.General, trana, tranb blas.Transpose, isgn int, scale float64) float64 {
	m := A.Rows
	n := B.Rows

	if m == 0 || n == 0 {
		return 0
	}

	// R = op(A)*X.
	R := zeros(m, n, n)
	blas64.Gemm(trana, blas.NoTrans, 1, A, X, 0, R)

	// R += isgn * X * op(B).
	blas64.Gemm(blas.NoTrans, tranb, float64(isgn), X, B, 1, R)

	// R -= scale * C.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			R.Data[i*R.Stride+j] -= scale * C.Data[i*C.Stride+j]
		}
	}

	return dlange(lapack.MaxColumnSum, m, n, R.Data, R.Stride)
}

// randomSchur generates a random upper quasi-triangular matrix (Schur form).
func randomSchur(n, stride int, rnd *rand.Rand) blas64.General {
	if n == 0 {
		return blas64.General{Rows: 0, Cols: 0, Stride: max(1, stride)}
	}

	a := zeros(n, n, stride)

	// Fill upper triangular part randomly.
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			a.Data[i*a.Stride+j] = rnd.NormFloat64()
		}
	}

	// Randomly create 2x2 blocks on the diagonal (complex eigenvalue pairs).
	i := 0
	for i < n-1 {
		if rnd.IntN(3) == 0 {
			// Create a 2x2 block with complex eigenvalues.
			// Use form [a b; -c a] where c > 0 for complex eigenvalues a ± i*sqrt(bc).
			a11 := a.Data[i*a.Stride+i]
			a.Data[(i+1)*a.Stride+i+1] = a11 // Make diagonal elements equal.
			a.Data[(i+1)*a.Stride+i] = -math.Abs(rnd.NormFloat64()) - 0.1
			i += 2
		} else {
			i++
		}
	}

	return a
}

// ones creates an m×n matrix filled with 1s.
func ones(m, n int) blas64.General {
	stride := max(1, n)
	data := make([]float64, m*stride)
	for i := range data {
		data[i] = 1
	}
	return blas64.General{Rows: m, Cols: n, Stride: stride, Data: data}
}

// generalEqual checks if two matrices are equal.
func generalEqual(a, b blas64.General) bool {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		return false
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			if a.Data[i*a.Stride+j] != b.Data[i*b.Stride+j] {
				return false
			}
		}
	}
	return true
}
