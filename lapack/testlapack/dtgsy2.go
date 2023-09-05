// Copyright ©2021 The Gonum Authors. All rights reserved.
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

type Dtgsy2er interface {
	Dtgsy2(trans blas.Transpose, ijob lapack.MaximizeNormXJob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, rdsum, rdscal float64, iwork []int) (scale, sumout, scalout float64, pq, info int)
}

func Dtgsy2Test(t *testing.T, impl Dtgsy2er) {
	const ldAdd = 5
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{2, 4, 9, 20} {
		for _, m := range []int{2, 4, 9, 20} {
			for _, lda := range []int{m, m + ldAdd} {
				for _, ldb := range []int{n, n + ldAdd} {
					for _, ldc := range []int{n, n + ldAdd} {
						for _, ldd := range []int{m, m + ldAdd} {
							for _, lde := range []int{n, n + ldAdd} {
								for _, ldf := range []int{n, n + ldAdd} {
									noTransIjob := lapack.MaximizeNormXJob(0)
									testSolveDtgsy2(t, impl, rnd, blas.NoTrans, noTransIjob, m, n, lda, ldb, ldc, ldd, lde, ldf)
									for _, ijob := range []lapack.MaximizeNormXJob{0, lapack.LocalLookAhead, lapack.NormalizedNullVector} {
										testSolveDtgsy2(t, impl, rnd, blas.Trans, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf)
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

func testSolveDtgsy2(t *testing.T, impl Dtgsy2er, rnd *rand.Rand, trans blas.Transpose, ijob lapack.MaximizeNormXJob, m, n, lda, ldb, ldc, ldd, lde, ldf int) {
	const tol = 1e-12
	name := fmt.Sprintf("trans=%v,ijob=%v,n=%v,m=%v,lda=%v,ldb=%v,ldc=%v,ldd=%v,lde=%v,ldf=%v", string(trans), ijob, n, m, lda, ldb, ldc, ldd, lde, ldf)
	lda = max(lda, max(1, m))
	ldb = max(ldb, max(1, n))
	ldc = max(ldc, max(1, n))
	ldd = max(ldd, max(1, m))
	lde = max(lde, max(1, n))
	ldf = max(ldf, max(1, n))

	// Generate random matrices (A, D) and (B, E) which must be
	// in generalized Schur canonical form, i.e. A, B are upper
	// quasi triangular and D, E are upper triangular.
	var a, b, c, d, e, f blas64.General
	// Real Schur canonical form. IF A is real, there exists a real orthogonal matrix V such that V^T A V = T is quasi-upper triangular.
	// This means that T is block upper triangular with 1-by1 and 2-by-2 blocks on the diagonal.
	// Its eigenvalues are the eigenvalues of the diagonal blocks. The 1-by-1 blocks correspond to real eigenvalues,
	// and the 2-by-2 blocks to complex conjugate pairs. From Wikipedia https://en.wikipedia.org/wiki/Talk%3ATriangular_matrix#Quasi-triangular_matrices
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
	// rdsum and rdscal only makes sense when Dtgsy2 is called by Dtgsyl.
	rdsum, rdscal := 0., 0.
	iwork := make([]int, m+n+2)
	scale, sum, scalout, pq, info := impl.Dtgsy2(trans, ijob, m, n, a.Data, lda, b.Data, ldb,
		c.Data, ldc, d.Data, ldd, e.Data, lde, f.Data, ldf, rdsum, rdscal, iwork)
	_, _, _ = sum, scalout, pq // Untested.
	if info >= 0 {
		t.Errorf("%v: got non-negative exit number. info=%d", name, info)
	}
	if m == 0 || n == 0 {
		return
	}
	// Compare block structure calculation with legacy algorithm.
	expectAIwork := calcBlockStructure(a)
	expectBIwork := calcBlockStructure(b)
	iworka := iwork[:len(expectAIwork)]
	iworkb := iwork[len(expectAIwork) : len(expectAIwork)+len(expectBIwork)]
	if !intsEqual(expectAIwork, iworka) {
		t.Errorf("%v: iwork calculation does not match expected for A. expect %d\ngot:%d", name, expectAIwork, iworka)
	}
	if !intsEqual(expectBIwork, iworkb) {
		t.Errorf("%v: iwork calculation does not match expected for B. expect %d\ngot:%d", name, expectBIwork, iworkb)
	}
	if scale == 0 {
		t.Errorf("%v: unexpected homogenous system solution", name)
	}

	// Solutions are written (R,L)->(C,F).
	r := c
	l := f
	rnorm := dlange(lapack.MaxColumnSum, r.Rows, r.Cols, r.Data, r.Stride)
	lnorm := dlange(lapack.MaxColumnSum, l.Rows, l.Cols, l.Data, l.Stride)
	rlnormmax := math.Max(rnorm, lnorm)
	if trans == blas.NoTrans {
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

// calcBlockStructure returns an array of indices which indicate the row
// at which a block begins of a Schur Canonical a matrix. The last entry is
// always the size of the matrix. len(iwork) <= m+1
//
// Consider the following 4×4 matrix:
//
//	[ -1   3   2  8]
//	[ -4  -12  1  1]
//	[ 0   0    2  8]
//	[ 0   0    0  1]
//
// The above matrix would return iwork of
//
//	[0  2  3  4]
//
// The routine was copied from the LAPACK Dtgsy2 implementation.
func calcBlockStructure(a blas64.General) (iwork []int) {
	if a.Cols != a.Rows {
		panic("block structure must be calculated for a square, quasitriangular matrix")
	}
	m := a.Cols
	p := -1
	iwork = make([]int, m+1)
	// Determine block structure of A.
	for i := 0; i < m; {
		p++
		iwork[p] = i
		if i == m-1 {
			break
		}
		if a.Data[(i+1)*a.Stride+i] != 0 {
			i += 2
			if i+2 < m && a.Data[(i+2)*a.Stride+i] != 0 {
				panic("matrix is not schur canonical")
			}
		} else {
			i++
		}
	}

	iwork[p+1] = m
	return iwork[:p+2]
}
