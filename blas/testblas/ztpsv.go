// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/blas"
)

type Ztpsver interface {
	Ztpsv(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n int, ap []complex128, x []complex128, incX int)

	Ztpmver
}

func ZtpsvTest(t *testing.T, impl Ztpsver) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans, blas.ConjTrans} {
			for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
				for _, n := range []int{0, 1, 2, 3, 4, 10} {
					for _, incX := range []int{-11, -3, -2, -1, 1, 2, 3, 7} {
						ztpsvTest(t, impl, uplo, trans, diag, n, incX, rnd)
					}
				}
			}
		}
	}
}

func ztpsvTest(t *testing.T, impl Ztpsver, uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, incX int, rnd *rand.Rand) {
	const tol = 1e-10

	lda := n
	a := makeZGeneral(nil, n, n, max(1, lda))
	if uplo == blas.Upper {
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				re := rnd.NormFloat64()
				im := rnd.NormFloat64()
				a[i*lda+j] = complex(re, im)
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j <= i; j++ {
				re := rnd.NormFloat64()
				im := rnd.NormFloat64()
				a[i*lda+j] = complex(re, im)
			}
		}
	}
	if diag == blas.Unit {
		for i := 0; i < n; i++ {
			a[i*lda+i] = znan
		}
	}
	ap := zPack(uplo, n, a, n)
	apCopy := make([]complex128, len(ap))
	copy(apCopy, ap)

	xtest := make([]complex128, n)
	for i := range xtest {
		re := rnd.NormFloat64()
		im := rnd.NormFloat64()
		xtest[i] = complex(re, im)
	}
	x := makeZVector(xtest, incX)
	want := make([]complex128, len(x))
	copy(want, x)

	impl.Ztpmv(uplo, trans, diag, n, ap, x, incX)
	impl.Ztpsv(uplo, trans, diag, n, ap, x, incX)

	if !zsame(ap, apCopy) {
		t.Errorf("Case uplo=%v,trans=%v,diag=%v,n=%v,incX=%v: unexpected modification of ap", uplo, trans, diag, n, incX)
	}
	if !zSameAtNonstrided(x, want, incX) {
		t.Errorf("Case uplo=%v,trans=%v,diag=%v,n=%v,incX=%v: unexpected modification of x\nwant %v\ngot  %v", uplo, trans, diag, n, incX, want, x)
	}
	if !zEqualApproxAtStrided(x, want, incX, tol) {
		t.Errorf("Case uplo=%v,trans=%v,diag=%v,n=%v,incX=%v: unexpected result\nwant %v\ngot  %v", uplo, trans, diag, n, incX, want, x)
	}
}
