// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dlagv2er interface {
	Dlagv2(a []float64, lda int, b []float64, ldb int, alphar, alphai, beta []float64) (csl, snl, csr, snr float64)
}

func Dlagv2Test(t *testing.T, impl Dlagv2er) {
	rnd := rand.New(rand.NewSource(1))
	const ldExtra = 0
	for _, lda := range []int{2, 2 + ldExtra} {
		for _, ldb := range []int{2, 2 + ldExtra} {
			for aKind := 9; aKind <= 20; aKind++ {
				for bKind := 9; bKind <= 20; bKind++ {
					dlagv2Test(t, impl, rnd, lda, ldb, aKind, bKind)
				}
			}
		}
	}
}

func dlagv2Test(t *testing.T, impl Dlagv2er, rnd *rand.Rand, lda, ldb int, aKind, bKind int) {
	const tol = 1e-14

	a := makeDlag2TestMatrix(rnd, lda, aKind)
	b := makeDlag2TestMatrix(rnd, ldb, bKind)
	b.Data[b.Stride] = math.NaN() // b is lower triangular.
	a.Data[a.Stride] = 0
	// Enforce positive, non-zero diagonal.
	b.Data[0] = math.Max(safmin, math.Abs(b.Data[0]))
	b.Data[b.Stride+1] = math.Max(safmin, math.Abs(b.Data[b.Stride+1]))

	if b.Data[0] < b.Data[b.Stride+1] {
		// Enforce descending order.
		b.Data[0], b.Data[b.Stride+1] = b.Data[b.Stride+1], b.Data[0]
	}

	aCopy := cloneGeneral(a)
	bCopy := cloneGeneral(b)
	bCopy.Data[bCopy.Stride] = 0 // zero out NaN for later comparisons.
	var alphar, alphai, beta [2]float64
	csl, snl, csr, snr := impl.Dlagv2(a.Data, a.Stride, b.Data, b.Stride, alphar[:], alphai[:], beta[:])
	beta1, beta2 := beta[0], beta[1]
	wi1, wi2 := alphai[0], alphai[1]
	wr1, wr2 := alphar[0], alphar[1]
	if beta1 == 0 {
		beta1 = 1
		wi1, wr1 = wi1/beta1, wr1/beta1
	}
	if beta2 == 0 {
		beta2 = 1
		wi2, wr2 = wi2/beta2, wr2/beta2
	}
	// Generate r1, r2 rotation matrices.
	// r1 = [ CSL  SNL; -SNL CSL]
	r1 := blas64.General{
		Data: []float64{csl, snl, -snl, csl},
		Rows: 2, Cols: 2, Stride: 2,
	}
	// r2 = [ CSR -SNR; SNR CSR]
	r2 := blas64.General{
		Data: []float64{csr, -snr, snr, csr},
		Rows: 2, Cols: 2, Stride: 2,
	}
	isComplexEig := wi1 != 0
	name := fmt.Sprintf("cmplx=%v,lda=%d,ldb=%d,aKind=%d,bKind=%d", isComplexEig, lda, ldb, aKind, bKind)

	aStr := sprint2x2Block("As", a.Data, a.Stride, 4)
	bStr := sprint2x2Block("Bs", b.Data, b.Stride, 4)

	if wi1 < 0 {
		t.Errorf("%s: negative wi1; wi1=%g,wi2=%g,\n%s\n%s", name, wi1, wi2, aStr, bStr)
	}
	if isComplexEig && wi2 != -wi1 {
		t.Errorf("%s: wi1 != -wi2; wi1=%g,wi2=%g,\n%s\n%s", name, wi1, wi2, aStr, bStr)
	}
	if math.Abs(b.Data[b.Stride]) > tol {
		t.Errorf("%s: expected b to remain upper triangular:\n%s", name, bStr)
	}
	if isComplexEig && b.Data[0] < b.Data[b.Stride+1] {
		t.Errorf("%s: expected b diagonal elements to be in descending order:\n%s", name, bStr)
	}

	if !isComplexEig && !isSchurCanonicalGeneral(a) {
		t.Errorf("%s: a is not Schur canonical:\n%s", name, aStr)
	}

	switch {
	case isComplexEig:
		// Complex eigenvalue pair.
		if wr1 != wr2 {
			t.Errorf("%s: complex eigenvalue but wr1 != wr2; wr1=%g, wr2=%g,\n%s\n%s", name, wr1, wr2, aStr, bStr)
		}
		if beta1 != beta2 {
			t.Errorf("%s: complex eigenvalue but scale1 != scale2; scale1=%g, scale2=%g,\n%s\n%s", name, beta1, beta2, aStr, bStr)
		}
		if b.Data[1] != 0 {
			t.Errorf("%s: expected b to be diagonal on complex pair:\n%s", name, bStr)
		}
	default:
		// Real eigenvalue pair.
		if wi1 != 0 || wi2 != 0 {
			t.Errorf("%s: real eigenvalue but wi1 != 0 or wi2 != 0; wi1=%g, wi2=%g,\n%s\n%s", name, wi1, wi2, aStr, bStr)
		}
		if a.Data[a.Stride] != 0 {
			t.Errorf("%s: expected a to be upper triangular on real pair:\n%s", name, aStr)
		}
	}

	// Test Rotations.

	aux := nanGeneral(2, 2, 2)
	result := nanGeneral(2, 2, 2)
	// Aschur = r1 * A * r2
	// Bschur = r1 * B * r2
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, r1, aCopy, 0, aux)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aux, r2, 0, result)
	if !equalApproxGeneral(result, a, tol) {
		t.Errorf("%s: bad rotation:\n%s\n%s",
			name, aStr, sprint2x2Block("A", result.Data, result.Stride, 4))
	}

	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, r1, bCopy, 0, aux)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aux, r2, 0, result)
	if !equalApproxGeneral(result, b, tol) {
		t.Errorf("%s: bad rotation:\n%s\n%s",
			name, bStr, sprint2x2Block("B", result.Data, result.Stride, 4))
	}

	// Test eigenvalues.

	res, err := residualDlag2(aCopy, bCopy, beta1, complex(wr1, wi1))
	if err != nil {
		t.Fatalf("%s: while computing residual: %v", name, err)
	}
	if res > tol || math.IsNaN(res) {
		t.Errorf("%s: unexpected first eigenvalue %g with s=%g; resid=%g, want<=%g\n%s\n%s", name, complex(wr1, wi1), beta1, res, tol, aStr, bStr)
	}

	res, err = residualDlag2(a, b, beta2, complex(wr2, wi2))
	if err != nil {
		t.Fatalf("%s: while computing residual: %v", name, err)
	}
	if res > tol || math.IsNaN(res) {
		t.Errorf("%s: unexpected second eigenvalue %g with s=%g; resid=%g, want<=%g\n%s\n%s", name, complex(wr2, wi2), beta2, res, tol, aStr, bStr)
	}
}

func sprint2x2Block(name string, a []float64, stride, prec int) string {
	const spaces = "                            "
	return fmt.Sprintf("%s = [%.*g,%.*g]\n%s[%.*g,%.*g]", name, prec, a[0], prec, a[1], spaces[:len(name)+3], prec, a[stride], prec, a[stride+1])
}
