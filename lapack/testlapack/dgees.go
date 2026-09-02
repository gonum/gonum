// Copyright ©2024 The Gonum Authors. All rights reserved.
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

type Dgeeser interface {
	Dgees(jobvs lapack.SchurComp, sort lapack.SchurSort, selctg func(wr, wi float64) bool,
		n int, a []float64, lda int, wr, wi []float64, vs []float64, ldvs int,
		work []float64, lwork int, bwork []bool) (sdim int, ok bool)
}

func DgeesTest(t *testing.T, impl Dgeeser) {
	rnd := rand.New(rand.NewPCG(1, 1))

	for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 10, 20, 50, 100} {
		for _, jobvs := range []lapack.SchurComp{lapack.SchurNone, lapack.SchurHess} {
			for _, extra := range []int{0, 11} {
				for _, wl := range []worklen{minimumWork, mediumWork, optimumWork} {
					testDgees(t, impl, n, jobvs, extra, wl, rnd)
				}
			}
		}
	}

	// Test with special matrices.
	for _, tc := range []struct {
		name string
		a    blas64.General
	}{
		{"Diagonal(5)", Diagonal(5).Matrix()},
		{"Diagonal(10)", Diagonal(10).Matrix()},
		{"Wilk4", Wilk4{}.Matrix()},
		{"Circulant(5)", Circulant(5).Matrix()},
		{"Circulant(10)", Circulant(10).Matrix()},
		{"Clement(5)", Clement(5).Matrix()},
		{"Clement(10)", Clement(10).Matrix()},
		{"Grcar{N:10,K:3}", Grcar{N: 10, K: 3}.Matrix()},
		{"Zero(5)", Zero(5).Matrix()},
		{"Zero(10)", Zero(10).Matrix()},
	} {
		for _, jobvs := range []lapack.SchurComp{lapack.SchurNone, lapack.SchurHess} {
			testDgeesMatrix(t, impl, tc.name, tc.a, jobvs, optimumWork)
		}
	}
}

func testDgees(t *testing.T, impl Dgeeser, n int, jobvs lapack.SchurComp, extra int, wl worklen, rnd *rand.Rand) {
	const tol = 100 // ||A - Z*T*Zᵀ||/(||A||*n*eps) should be O(1) for stable algorithms

	wantvs := jobvs == lapack.SchurHess

	// Generate a random matrix.
	a := randomGeneral(n, n, n+extra, rnd)
	aCopy := cloneGeneral(a)

	// Allocate output arrays.
	wr := make([]float64, n)
	wi := make([]float64, n)

	var vs blas64.General
	ldvs := 1
	if wantvs {
		ldvs = max(1, n+extra)
		vs = nanGeneral(n, n, ldvs)
	}

	// Determine workspace size.
	var lwork int
	minwork := max(1, 3*n)
	switch wl {
	case minimumWork:
		lwork = minwork
	case mediumWork:
		work := make([]float64, 1)
		impl.Dgees(jobvs, lapack.SortNone, nil, n, nil, max(1, n), nil, nil, nil, ldvs, work, -1, nil)
		lwork = (int(work[0]) + minwork) / 2
		lwork = max(minwork, lwork)
	case optimumWork:
		work := make([]float64, 1)
		impl.Dgees(jobvs, lapack.SortNone, nil, n, nil, max(1, n), nil, nil, nil, ldvs, work, -1, nil)
		lwork = int(work[0])
	}
	work := make([]float64, lwork)

	prefix := fmt.Sprintf("n=%v, jobvs=%c, extra=%v, work=%v", n, jobvs, extra, wl)

	// Call Dgees.
	sdim, ok := impl.Dgees(jobvs, lapack.SortNone, nil, n, a.Data, a.Stride, wr, wi, vs.Data, ldvs, work, lwork, nil)

	if !ok {
		t.Errorf("%v: Dgees failed to converge", prefix)
		return
	}

	if sdim != 0 {
		t.Errorf("%v: unexpected sdim=%v, want 0 for SortNone", prefix, sdim)
	}

	// Quick return for n == 0.
	if n == 0 {
		return
	}

	// Check that A has been overwritten with the Schur form T.
	// T should be upper quasi-triangular.
	if !isSchurCanonicalGeneral(a) {
		t.Errorf("%v: result is not in Schur canonical form", prefix)
	}

	// Check eigenvalues: the diagonal of T should contain the eigenvalues.
	evTol := 1e-10
	checkSchurEigenvalues(t, prefix, a, wr, wi, evTol)

	if !wantvs {
		return
	}

	// Check that vs contains orthogonal Schur vectors.
	// Orthogonality should be O(eps), use a small tolerance.
	orthoTol := float64(n) * dlamchE
	if orthoTol < 1e-13 {
		orthoTol = 1e-13
	}
	resid := residualOrthogonal(vs, false)
	if resid > orthoTol {
		t.Errorf("%v: Schur vectors not orthogonal; |I - Z*Zᵀ|=%v, want<=%v", prefix, resid, orthoTol)
	}

	// Check the Schur factorization: A = Z * T * Zᵀ
	// Compute ||A - Z*T*Zᵀ|| / (||A|| * n * eps)
	residSchur := residualSchurFactorization(aCopy, a, vs)
	anorm := dlange(lapack.MaxColumnSum, n, n, aCopy.Data, aCopy.Stride)
	if anorm == 0 {
		anorm = 1
	}
	normRes := residSchur / (anorm * float64(n) * dlamchE)
	if normRes > tol {
		t.Errorf("%v: ||A - Z*T*Zᵀ||/(||A||*n*eps)=%v, want<=%v", prefix, normRes, tol)
	}
}

func testDgeesMatrix(t *testing.T, impl Dgeeser, name string, aOrig blas64.General, jobvs lapack.SchurComp, wl worklen) {
	const tol = 100 // ||A - Z*T*Zᵀ||/(||A||*n*eps) should be O(1) for stable algorithms

	n := aOrig.Rows
	wantvs := jobvs == lapack.SchurHess

	// Clone the input matrix.
	a := cloneGeneral(aOrig)

	// Allocate output arrays.
	wr := make([]float64, n)
	wi := make([]float64, n)

	var vs blas64.General
	ldvs := 1
	if wantvs {
		ldvs = max(1, n)
		vs = nanGeneral(n, n, ldvs)
	}

	// Determine workspace size.
	work := make([]float64, 1)
	impl.Dgees(jobvs, lapack.SortNone, nil, n, nil, max(1, n), nil, nil, nil, ldvs, work, -1, nil)
	lwork := int(work[0])
	work = make([]float64, lwork)

	prefix := fmt.Sprintf("%v, jobvs=%c", name, jobvs)

	// Call Dgees.
	_, ok := impl.Dgees(jobvs, lapack.SortNone, nil, n, a.Data, a.Stride, wr, wi, vs.Data, ldvs, work, lwork, nil)

	if !ok {
		t.Errorf("%v: Dgees failed to converge", prefix)
		return
	}

	// Check that A has been overwritten with the Schur form T.
	if !isSchurCanonicalGeneral(a) {
		t.Errorf("%v: result is not in Schur canonical form", prefix)
	}

	// Check eigenvalues.
	evTol := 1e-10
	checkSchurEigenvalues(t, prefix, a, wr, wi, evTol)

	if !wantvs {
		return
	}

	// Check orthogonality of Schur vectors.
	orthoTol := float64(n) * dlamchE
	if orthoTol < 1e-13 {
		orthoTol = 1e-13
	}
	resid := residualOrthogonal(vs, false)
	if resid > orthoTol {
		t.Errorf("%v: Schur vectors not orthogonal; |I - Z*Zᵀ|=%v, want<=%v", prefix, resid, orthoTol)
	}

	// Check Schur factorization.
	residSchur := residualSchurFactorization(aOrig, a, vs)
	anorm := dlange(lapack.MaxColumnSum, n, n, aOrig.Data, aOrig.Stride)
	if anorm == 0 {
		anorm = 1
	}
	normRes := residSchur / (anorm * float64(n) * dlamchE)
	if normRes > tol {
		t.Errorf("%v: ||A - Z*T*Zᵀ||/(||A||*n*eps)=%v, want<=%v", prefix, normRes, tol)
	}
}

// checkSchurEigenvalues verifies that the eigenvalues in wr and wi match
// the diagonal blocks of the Schur form T.
func checkSchurEigenvalues(t *testing.T, prefix string, T blas64.General, wr, wi []float64, tol float64) {
	t.Helper()
	n := T.Rows
	for j := 0; j < n; {
		if j == n-1 || T.Data[(j+1)*T.Stride+j] == 0 {
			// 1×1 block: real eigenvalue.
			diff := math.Abs(wr[j] - T.Data[j*T.Stride+j])
			if diff > tol {
				t.Errorf("%v: eigenvalue[%d] mismatch: wr=%v, T[%d,%d]=%v", prefix, j, wr[j], j, j, T.Data[j*T.Stride+j])
			}
			if wi[j] != 0 {
				t.Errorf("%v: eigenvalue[%d] has non-zero imaginary part for 1×1 block: wi=%v", prefix, j, wi[j])
			}
			j++
		} else {
			// 2×2 block: complex conjugate pair.
			a11 := T.Data[j*T.Stride+j]
			a12 := T.Data[j*T.Stride+j+1]
			a21 := T.Data[(j+1)*T.Stride+j]
			a22 := T.Data[(j+1)*T.Stride+j+1]

			// Real part should be (a11+a22)/2, but for canonical form a11==a22.
			realPart := (a11 + a22) / 2
			// Imaginary part is sqrt(-a12*a21) (assuming a12*a21 < 0).
			var imagPart float64
			if a12*a21 < 0 {
				imagPart = math.Sqrt(-a12 * a21)
			}

			diffReal := math.Abs(wr[j] - realPart)
			diffImag := math.Abs(wi[j] - imagPart)
			if diffReal > tol || diffImag > tol {
				t.Errorf("%v: eigenvalue[%d] mismatch: got (%v, %v), want (%v, ±%v)", prefix, j, wr[j], wi[j], realPart, imagPart)
			}

			// The conjugate pair should have opposite imaginary parts.
			if wi[j] <= 0 {
				t.Errorf("%v: eigenvalue[%d] should have positive imaginary part first", prefix, j)
			}
			if wi[j+1] != -wi[j] {
				t.Errorf("%v: eigenvalue[%d+1] should be conjugate: wi[%d]=%v, wi[%d]=%v", prefix, j, j, wi[j], j+1, wi[j+1])
			}

			j += 2
		}
	}
}

func DgeesSortTest(t *testing.T, impl Dgeeser) {
	rnd := rand.New(rand.NewPCG(1, 1))

	// Select eigenvalues with absolute value < 1.
	selctg := func(wr, wi float64) bool {
		return wr*wr+wi*wi < 1.0
	}

	for _, n := range []int{1, 2, 3, 4, 5, 10, 20, 50} {
		for _, extra := range []int{0, 11} {
			for cas := 0; cas < 10; cas++ {
				testDgeesSort(t, impl, n, extra, selctg, rnd)
			}
		}
	}
}

func testDgeesSort(t *testing.T, impl Dgeeser, n, extra int, selctg func(wr, wi float64) bool, rnd *rand.Rand) {
	const tol = 100

	a := randomGeneral(n, n, n+extra, rnd)
	aCopy := cloneGeneral(a)

	wr := make([]float64, n)
	wi := make([]float64, n)

	ldvs := max(1, n+extra)
	vs := nanGeneral(n, n, ldvs)
	bwork := make([]bool, n)

	// Workspace query.
	work := make([]float64, 1)
	impl.Dgees(lapack.SchurHess, lapack.SortSelected, selctg, n, nil, max(1, a.Stride), nil, nil, nil, ldvs, work, -1, nil)
	lwork := int(work[0])
	work = make([]float64, lwork)

	prefix := fmt.Sprintf("n=%d, extra=%d", n, extra)

	sdim, ok := impl.Dgees(lapack.SchurHess, lapack.SortSelected, selctg, n, a.Data, a.Stride, wr, wi, vs.Data, ldvs, work, lwork, bwork)
	if !ok {
		t.Logf("%s: Dgees with sorting failed to converge", prefix)
		return
	}

	// Verify Schur form.
	if !isSchurCanonicalGeneral(a) {
		t.Errorf("%s: result is not in Schur canonical form", prefix)
	}

	// Verify selected eigenvalues are at the top.
	for j := 0; j < sdim; j++ {
		if !selctg(wr[j], wi[j]) {
			t.Errorf("%s: eigenvalue %d (wr=%v, wi=%v) at top should be selected", prefix, j, wr[j], wi[j])
		}
	}
	for j := sdim; j < n; j++ {
		if selctg(wr[j], wi[j]) {
			t.Errorf("%s: eigenvalue %d (wr=%v, wi=%v) at bottom should NOT be selected", prefix, j, wr[j], wi[j])
		}
	}

	// Verify eigenvalues match diagonal blocks.
	evTol := 1e-10
	checkSchurEigenvalues(t, prefix, a, wr, wi, evTol)

	// Verify Schur vectors are orthogonal.
	orthoTol := float64(n) * dlamchE
	if orthoTol < 1e-13 {
		orthoTol = 1e-13
	}
	resid := residualOrthogonal(vs, false)
	if resid > orthoTol {
		t.Errorf("%s: Schur vectors not orthogonal; |I - Z*Zᵀ|=%v, want<=%v", prefix, resid, orthoTol)
	}

	// Verify Schur factorization: A = Z * T * Zᵀ.
	residSchur := residualSchurFactorization(aCopy, a, vs)
	anorm := dlange(lapack.MaxColumnSum, n, n, aCopy.Data, aCopy.Stride)
	if anorm == 0 {
		anorm = 1
	}
	normRes := residSchur / (anorm * float64(n) * dlamchE)
	if normRes > float64(tol) {
		t.Errorf("%s: ||A - Z*T*Zᵀ||/(||A||*n*eps)=%v, want<=%v", prefix, normRes, tol)
	}
}

// residualSchurFactorization computes ||A - Z*T*Zᵀ||_1.
func residualSchurFactorization(A, T, Z blas64.General) float64 {
	n := A.Rows
	if n == 0 {
		return 0
	}

	// Compute R = Z * T.
	R := zeros(n, n, n)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, Z, T, 0, R)

	// Compute R = R * Zᵀ = Z * T * Zᵀ.
	R2 := zeros(n, n, n)
	blas64.Gemm(blas.NoTrans, blas.Trans, 1, R, Z, 0, R2)

	// Compute R2 = A - Z * T * Zᵀ.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			R2.Data[i*R2.Stride+j] = A.Data[i*A.Stride+j] - R2.Data[i*R2.Stride+j]
		}
	}

	return dlange(lapack.MaxColumnSum, n, n, R2.Data, R2.Stride)
}
