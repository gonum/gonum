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
	"gonum.org/v1/gonum/lapack"
)

type Dtgsyler interface {
	Dtgsyl(trans blas.Transpose, ijob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, work []float64, iwork []int, workspaceQuery bool) (difOut, scaleOut float64, infoOut int)
}

func DtgsylTest(t *testing.T, impl Dtgsyler) {
	const ldAdd = 5
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{4, 9, 20} {
		for _, m := range []int{4, 9, 20} {
			for _, lda := range []int{m, m + ldAdd} {
				for _, ldb := range []int{n, n + ldAdd} {
					for _, ldc := range []int{n, n + ldAdd} {
						for _, ldd := range []int{m, m + ldAdd} {
							for _, lde := range []int{n, n + ldAdd} {
								for _, ldf := range []int{n, n + ldAdd} {
									for _, ijob := range []int{2, 1, 0} {
										testSolveDtgsyl(t, impl, rnd, blas.NoTrans, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf)
										testSolveDtgsyl(t, impl, rnd, blas.Trans, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func testSolveDtgsyl(t *testing.T, impl Dtgsyler, rnd *rand.Rand, trans blas.Transpose, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf int) {
	const tol = 1e-12
	name := fmt.Sprintf("trans=%v,ijob=%v,n=%v,m=%v,lda=%v,ldb=%v,ldc=%v,ldd=%v,lde=%v,ldf=%v", string(trans), ijob, n, m, lda, ldb, ldc, ldd, lde, ldf)
	lda = max(lda, max(1, m))
	ldb = max(ldb, max(1, n))
	ldc = max(ldc, max(1, n))
	ldd = max(ldd, max(1, m))
	lde = max(lde, max(1, n))
	ldf = max(ldf, max(1, n))
	notrans := trans == blas.NoTrans
	// Generate random matrices (A, D) and (B, E) which must be
	// in generalized Schur canonical form, i.e. A, B are upper
	// quasi triangular and D, E are upper triangular.
	var a, b, c, d, e, f blas64.General
	a, _, _ = randomSchurCanonical(m, lda, false, rnd)
	b, _, _ = randomSchurCanonical(n, ldb, false, rnd)

	d = randomUpperTriGeneral(m, ldd, rnd)
	e = randomUpperTriGeneral(n, lde, rnd)

	// Generate random general matrix.
	c = randomGeneral(m, n, ldc, rnd)
	f = randomGeneral(m, n, ldf, rnd)
	cCopy := cloneGeneral(c)
	fCopy := cloneGeneral(f)
	// Calculate norms
	anorm := dlange(lapack.MaxColumnSum, a.Rows, a.Cols, a.Data, a.Stride)
	bnorm := dlange(lapack.MaxColumnSum, b.Rows, b.Cols, b.Data, b.Stride)
	cnorm := dlange(lapack.MaxColumnSum, c.Rows, c.Cols, c.Data, c.Stride)
	dnorm := dlange(lapack.MaxColumnSum, d.Rows, d.Cols, d.Data, d.Stride)
	enorm := dlange(lapack.MaxColumnSum, e.Rows, e.Cols, e.Data, e.Stride)
	fnorm := dlange(lapack.MaxColumnSum, f.Rows, f.Cols, f.Data, f.Stride)

	// Query for optimum workspace size.
	var query [1]float64
	impl.Dtgsyl(trans, ijob, m, n, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride, d.Data, d.Stride, e.Data, e.Stride, f.Data, f.Stride, query[:], nil, true)
	lwork := int(query[0] + dlamchE)
	if lwork < 1 {
		t.Fatalf("%v: bad workspace query lwork=%d", name, lwork)
	}
	lworkMin := 1
	if notrans && (ijob == 1 || ijob == 2) {
		lworkMin = 2 * m * n
	}
	if lwork < lworkMin {
		t.Fatalf("%v: bad workspace query lwork=%d, expected >=%d", name, lwork, lworkMin)
	}
	iwork := make([]int, m+n+6)
	work := make([]float64, lwork)
	dif, scale, info := impl.Dtgsyl(trans, ijob, m, n, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride, d.Data, d.Stride, e.Data, e.Stride, f.Data, f.Stride, work, iwork, false)
	_, _ = dif, scale // untested.
	if info >= 0 {
		t.Errorf("%v: info>=0: matrix was perturbed", name)
	}
	lwork = int(work[0])
	if lwork < 1 {
		t.Fatalf("%v: bad workspace query lwork=%d", name, lwork)
	}
	// Solutions are written (R,L)->(C,F).
	r := c
	l := f
	rnorm := dlange(lapack.MaxColumnSum, r.Rows, r.Cols, r.Data, r.Stride)
	lnorm := dlange(lapack.MaxColumnSum, l.Rows, l.Cols, l.Data, l.Stride)
	rlnormmax := math.Max(rnorm, lnorm)
	if notrans {
		// Calculate residuals
		// | A * R - L * B - scale * C |  from (1)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, a, r, -scale, cCopy)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, l, b, 1, cCopy)
		res := dlange(lapack.MaxColumnSum, m, n, cCopy.Data, cCopy.Stride) / math.Max(math.Max(anorm, rlnormmax), math.Max(bnorm, cnorm))
		if res > tol || math.IsNaN(res) {
			t.Errorf("%v: | A * R - L * B - scale * C | residual large or NaN %v", name, res)
		}

		// | D * R - L * E - scale * F |  from (1)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, d, r, -scale, fCopy)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, l, e, 1, fCopy)
		res = dlange(lapack.MaxColumnSum, m, n, fCopy.Data, fCopy.Stride) / math.Max(math.Max(dnorm, rlnormmax), math.Max(enorm, fnorm))
		if res > tol || math.IsNaN(res) {
			t.Errorf("%v: | D * R - L * E - scale * F | residual large or NaN %v", name, res)
		}
	} else {
		// Calculate residuals
		// | Aᵀ * R + Dᵀ * L - scale * C |  from (3)
		blas64.Gemm(trans, blas.NoTrans, 1, a, r, -scale, cCopy)
		blas64.Gemm(trans, blas.NoTrans, 1, d, l, 1, cCopy)
		res := dlange(lapack.MaxColumnSum, m, n, cCopy.Data, cCopy.Stride) / math.Max(math.Max(anorm, rlnormmax), math.Max(dnorm, cnorm))
		if res > tol || math.IsNaN(res) {
			t.Errorf("%v: | Aᵀ * R + Dᵀ * L - scale * C | residual large or NaN %v", name, res)
		}

		// | R * Bᵀ + L * Eᵀ - scale * -F |  from (3)
		blas64.Gemm(blas.NoTrans, trans, 1, r, b, scale, fCopy)
		blas64.Gemm(blas.NoTrans, trans, 1, l, e, 1, fCopy)
		res = dlange(lapack.MaxColumnSum, m, n, fCopy.Data, fCopy.Stride) / math.Max(math.Max(bnorm, rlnormmax), math.Max(enorm, fnorm))

		if res > tol || math.IsNaN(res) {
			t.Errorf("%v: | R * Bᵀ + L * Eᵀ - scale * -F | residual large or NaN %v", name, res)
		}
	}
}
