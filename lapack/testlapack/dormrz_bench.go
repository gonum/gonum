// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
)

func DormrzBenchmark(b *testing.B, impl Dormrzer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, sz := range []int{10, 50, 100, 200} {
		for _, side := range []blas.Side{blas.Left, blas.Right} {
			sideName := "Left"
			if side == blas.Right {
				sideName = "Right"
			}

			// Z is of order nq where nq=m (Left) or nq=n (Right).
			// Dtzrzf(k, nq, ...) produces k reflectors with l = nq - k.
			mc := sz
			nc := sz
			k := sz / 2
			nq := mc
			if side == blas.Right {
				nq = nc
			}
			// Expand nq so factorization is non-trivial.
			nq = 2 * nq
			if side == blas.Left {
				mc = nq
			} else {
				nc = nq
			}
			l := nq - k
			lda := nq
			ldc := nc

			aOrig := make([]float64, k*lda)
			for i := range aOrig {
				aOrig[i] = rnd.NormFloat64()
			}
			a := make([]float64, len(aOrig))
			tau := make([]float64, k)

			copy(a, aOrig)
			fwork := make([]float64, k)
			impl.Dtzrzf(k, nq, a, lda, tau, fwork, len(fwork))
			aFact := make([]float64, len(a))
			copy(aFact, a)
			tauFact := make([]float64, len(tau))
			copy(tauFact, tau)

			cOrig := make([]float64, mc*ldc)
			for i := range cOrig {
				cOrig[i] = rnd.NormFloat64()
			}
			c := make([]float64, len(cOrig))

			nw := nc
			if side == blas.Right {
				nw = mc
			}

			b.Run(fmt.Sprintf("%s/m=%d/n=%d", sideName, mc, nc), func(b *testing.B) {
				workDormrz := make([]float64, nw)
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					copy(c, cOrig)
					copy(a, aFact)
					copy(tau, tauFact)
					b.StartTimer()
					impl.Dormrz(side, blas.NoTrans, mc, nc, k, l, a, lda, tau, c, ldc, workDormrz, len(workDormrz))
				}
			})
		}
	}
}
