// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dpotrfer interface {
	Dpotrf(ul blas.Uplo, n int, a []float64, lda int) (ok bool)
}

func DpotrfTest(t *testing.T, impl Dpotrfer) {
	const tol = 1e-13
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for tc, test := range []struct {
			n   int
			lda int
		}{
			{1, 0},
			{2, 0},
			{3, 0},
			{10, 0},
			{30, 0},
			{63, 0},
			{65, 0},
			{127, 0},
			{129, 0},
			{500, 0},
			{1, 10},
			{2, 10},
			{3, 10},
			{10, 20},
			{30, 50},
			{63, 100},
			{65, 100},
			{127, 200},
			{129, 200},
			{500, 600},
		} {
			n := test.n
			lda := test.lda
			if lda == 0 {
				lda = n
			}
			// Construct a diagonally-dominant symmetric matrix.
			// Such a matrix is positive definite.
			a := make([]float64, n*lda)
			for i := range a {
				a[i] = rnd.Float64()
			}
			for i := 0; i < n; i++ {
				a[i*lda+i] += float64(n)
				for j := 0; j < i; j++ {
					a[i*lda+j] = a[j*lda+i]
				}
			}

			aCopy := make([]float64, len(a))
			copy(aCopy, a)

			ok := impl.Dpotrf(uplo, n, a, lda)

			if !ok {
				t.Errorf("Case %v: unexpected failure for positive definite matrix", tc)
				continue
			}

			switch uplo {
			case blas.Upper:
				for i := 0; i < n; i++ {
					for j := 0; j < i; j++ {
						a[i*lda+j] = 0
					}
				}
			case blas.Lower:
				for i := 0; i < n; i++ {
					for j := i + 1; j < n; j++ {
						a[i*lda+j] = 0
					}
				}
			default:
				panic("bad uplo")
			}

			ans := make([]float64, len(a))
			switch uplo {
			case blas.Upper:
				// Multiply U^T * U.
				bi.Dsyrk(uplo, blas.Trans, n, n, 1, a, lda, 0, ans, lda)
			case blas.Lower:
				// Multiply L * L^T.
				bi.Dsyrk(uplo, blas.NoTrans, n, n, 1, a, lda, 0, ans, lda)
			}

			match := true
			switch uplo {
			case blas.Upper:
				for i := 0; i < n; i++ {
					for j := i; j < n; j++ {
						if !floats.EqualWithinAbsOrRel(ans[i*lda+j], aCopy[i*lda+j], tol, tol) {
							match = false
						}
					}
				}
			case blas.Lower:
				for i := 0; i < n; i++ {
					for j := 0; j <= i; j++ {
						if !floats.EqualWithinAbsOrRel(ans[i*lda+j], aCopy[i*lda+j], tol, tol) {
							match = false
						}
					}
				}
			}
			if !match {
				t.Errorf("Case %v (uplo=%v,n=%v,lda=%v): unexpected result\n%v\n%v", tc, uplo, n, lda, ans, aCopy)
			}
		}
	}
}
