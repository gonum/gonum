// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

var result float64

func DlantbBenchmark(b *testing.B, impl Dlantber) {
	rnd := rand.New(rand.NewSource(1))
	for _, bm := range []struct {
		n, k int
	}{
		{n: 10000, k: 1},
		{n: 10000, k: 2},
		{n: 10000, k: 100},
	} {
		n := bm.n
		k := bm.k
		lda := k + 1
		for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
			var work []float64
			if norm == lapack.MaxColumnSum {
				work = make([]float64, n)
			}
			for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
				for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
					name := fmt.Sprintf("%v%v%vN=%vK=%v", normToString(norm), uploToString(uplo), diagToString(diag), n, k)
					b.Run(name, func(b *testing.B) {
						for i := 0; i < b.N; i++ {
							b.StopTimer()
							a := make([]float64, n*lda)
							for i := range a {
								a[i] = rnd.NormFloat64()
							}
							b.StartTimer()
							result = impl.Dlantb(norm, uplo, diag, bm.n, bm.k, a, lda, work)
						}
					})
				}
			}
		}
	}
}
