// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func DtzrzfBenchmark(b *testing.B, impl Dtzrzfer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, m := range []int{10, 50, 100, 200} {
		n := 2 * m
		lda := n
		aOrig := make([]float64, m*lda)
		for i := range aOrig {
			aOrig[i] = rnd.NormFloat64()
		}
		a := make([]float64, len(aOrig))
		tau := make([]float64, m)
		work := make([]float64, m)

		b.Run(fmt.Sprintf("m=%d/n=%d", m, n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				copy(a, aOrig)
				b.StartTimer()
				impl.Dtzrzf(m, n, a, lda, tau, work, len(work))
			}
		})
	}
}
