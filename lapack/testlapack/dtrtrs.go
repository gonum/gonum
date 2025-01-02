// Copyright ©2020 The Gonum Authors. All rights reserved.
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
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dtrtrser interface {
	Dtrtrs(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, nrhs int, a []float64, lda int, b []float64, ldb int) bool
}

func DtrtrsTest(t *testing.T, impl Dtrtrser) {
	rnd := rand.New(rand.NewPCG(1, 1))

	for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans, blas.ConjTrans} {
		name := transToString(trans)
		t.Run(name, func(t *testing.T) {
			for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
				for _, diag := range []blas.Diag{blas.Unit, blas.NonUnit} {
					for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
						for _, nrhs := range []int{0, 1, 2, 3, 4, 5, 10, 15} {
							for _, lda := range []int{max(1, n), n + 3} {
								for _, ldb := range []int{max(1, nrhs), nrhs + 3} {
									if diag == blas.Unit {
										dtrtrsTest(t, impl, rnd, uplo, trans, diag, n, nrhs, lda, ldb, false)
									} else {
										dtrtrsTest(t, impl, rnd, uplo, trans, diag, n, nrhs, lda, ldb, true)
										dtrtrsTest(t, impl, rnd, uplo, trans, diag, n, nrhs, lda, ldb, false)
									}
								}
							}
						}
					}
				}
			}
		})
	}
}

func dtrtrsTest(t *testing.T, impl Dtrtrser, rnd *rand.Rand, uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, nrhs int, lda, ldb int, singular bool) {
	if singular && diag == blas.Unit {
		panic("blas.Unit triangular matrix cannot be singular")
	}

	const tol = 1e-14

	if n == 0 {
		singular = false
	}
	name := fmt.Sprintf("uplo=%v,diag=%v,n=%v,nrhs=%v,lda=%v,ldb=%v,sing=%v", string(uplo), string(diag), n, nrhs, lda, ldb, singular)

	// Generate a random triangular matrix A. One of its triangles won't be
	// referenced.
	a := make([]float64, n*lda)
	for i := range a {
		a[i] = rnd.NormFloat64()
	}
	if singular {
		i := rnd.IntN(n)
		a[i*lda+i] = 0
	}
	aCopy := make([]float64, len(a))
	copy(aCopy, a)

	// Generate a random solution matrix X.
	x := make([]float64, n*ldb)
	for i := range x {
		x[i] = rnd.NormFloat64()
	}

	// Generate the right-hand side as A * X  or  Aᵀ * X.
	b := make([]float64, len(x))
	copy(b, x)
	bi := blas64.Implementation()
	bi.Dtrmm(blas.Left, uplo, trans, diag, n, nrhs, 1, a, lda, b, ldb)

	got := make([]float64, len(b))
	copy(got, b)

	ok := impl.Dtrtrs(uplo, trans, diag, n, nrhs, a, lda, got, ldb)

	if !floats.Equal(a, aCopy) {
		t.Errorf("%v: unexpected modification of A", name)
	}

	if ok == singular {
		t.Errorf("%v: misdetected singular matrix, ok=%v", name, ok)
	}

	if !ok {
		if !floats.Equal(got, b) {
			t.Errorf("%v: unexpected modification of B when singular", name)
		}
		return
	}

	if n == 0 || nrhs == 0 {
		return
	}

	work := make([]float64, n)

	// Compute the 1-norm of A or Aᵀ.
	var aNorm float64
	if trans == blas.NoTrans {
		aNorm = dlantr(lapack.MaxColumnSum, uplo, diag, n, n, a, lda, work)
	} else {
		aNorm = dlantr(lapack.MaxRowSum, uplo, diag, n, n, a, lda, work)
	}

	// Compute the maximum over the number of right-hand sides of
	//  |op(A)*x-b| / (|op(A)| * |x|)
	var resid float64
	for j := 0; j < nrhs; j++ {
		bi.Dcopy(n, got[j:], ldb, work, 1)
		bi.Dtrmv(uplo, trans, diag, n, a, lda, work, 1)
		bi.Daxpy(n, -1, b[j:], ldb, work, 1)
		rjNorm := bi.Dasum(n, work, 1)
		xNorm := bi.Dasum(n, got[j:], ldb)
		resid = math.Max(resid, rjNorm/aNorm/xNorm)
	}
	if resid > tol {
		t.Errorf("%v: unexpected result; resid=%v,want<=%v", name, resid, tol)
	}
}
