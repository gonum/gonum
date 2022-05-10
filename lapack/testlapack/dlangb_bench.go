// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/lapack"
)

func DlangbBenchmark(b *testing.B, impl Dlangber) {
	var result float64
	rnd := rand.New(rand.NewSource(1))
	for _, bm := range []struct {
		n, k int
	}{
		{n: 1000, k: 0},
		{n: 1000, k: 1},
		{n: 1000, k: 2},
		{n: 1000, k: 5},
		{n: 1000, k: 8},
		{n: 1000, k: 10},
		{n: 1000, k: 20},
		{n: 1000, k: 30},
		{n: 10000, k: 0},
		{n: 10000, k: 1},
		{n: 10000, k: 2},
		{n: 10000, k: 5},
		{n: 10000, k: 8},
		{n: 10000, k: 10},
		{n: 10000, k: 30},
		{n: 10000, k: 60},
		{n: 10000, k: 100},
	} {
		n := bm.n
		k := bm.k
		lda := 2*k + 1
		aCopy := make([]float64, n*lda)
		for i := range aCopy {
			aCopy[i] = 1 - 2*rnd.Float64()
		}
		a := make([]float64, len(aCopy))

		for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum} {
			name := fmt.Sprintf("%v_N=%v_K=%v", normToString(norm), n, k)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = impl.Dlangb(norm, n, n, k, k, a, lda)
				}
			})
		}

		// Frobenius norm is benchmarked separately because its execution time
		// depends on the element magnitude.
		norm := lapack.Frobenius
		for _, scale := range []string{"Small", "Medium", "Big"} {
			name := fmt.Sprintf("%v_N=%v_K=%v_%v", normToString(norm), n, k, scale)
			var scl float64
			switch scale {
			default:
				scl = 1
			case "Small":
				scl = smlnum
			case "Big":
				scl = bignum
			}
			// Scale some elements so that the matrix contains a mix of small
			// and medium, all medium, or big and medium values.
			copy(a, aCopy)
			for i := range a {
				if i%2 == 0 {
					a[i] *= scl
				}
			}
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					result = impl.Dlangb(norm, n, n, k, k, a, lda)
				}
			})
		}
	}
	if math.IsNaN(result) {
		b.Error("unexpected NaN result")
	}
}
