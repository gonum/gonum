// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func Dlaic1Benchmark(b *testing.B, impl Dlaic1er) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, j := range []int{10, 100, 1000, 10000} {
		x := make([]float64, j)
		w := make([]float64, j)
		for i := range x {
			x[i] = rnd.NormFloat64()
			w[i] = rnd.NormFloat64()
		}
		sest := rnd.Float64() + 1
		gamma := rnd.NormFloat64()

		for _, job := range []int{1, 2} {
			b.Run(fmt.Sprintf("job=%d/j=%d", job, j), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					impl.Dlaic1(job, j, x, sest, w, gamma)
				}
			})
		}
	}
}
