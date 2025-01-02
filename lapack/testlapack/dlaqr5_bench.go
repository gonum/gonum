// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func Dlaqr5Benchmark(b *testing.B, impl Dlaqr5er) {
	const (
		random = iota
		iplusj
		laplacian
	)
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, typ := range []int{random, iplusj, laplacian} {
		for _, n := range []int{100, 200, 500, 1000} {
			h := zeros(n, n, n)
			var name string
			switch typ {
			case random:
				name = fmt.Sprintf("HessenbergRandom%d", n)
				h = randomHessenberg(n, n, rnd)
			case iplusj:
				name = fmt.Sprintf("HessenbergIPlusJ%d", n)
				for i := 0; i < n; i++ {
					for j := max(0, i-1); j < n; j++ {
						h.Data[i*h.Stride+j] = float64(i + j + 2)
					}
				}
			case laplacian:
				name = fmt.Sprintf("Laplacian%d", n)
				for i := 0; i < n; i++ {
					if i > 0 {
						h.Data[i*h.Stride+i-1] = -1
					}
					h.Data[i*h.Stride+i] = 2
					if i < n-1 {
						h.Data[i*h.Stride+i+1] = -1
					}
				}
			}
			hCopy := cloneGeneral(h)
			nshifts := 2 * n
			sr := make([]float64, nshifts)
			si := make([]float64, nshifts)
			for i := 0; i < nshifts; {
				if i == nshifts-1 || rnd.Float64() < 0.5 {
					re := rnd.NormFloat64()
					sr[i], si[i] = re, 0
					i++
					continue
				}
				re := rnd.NormFloat64()
				im := rnd.NormFloat64()
				sr[i], sr[i+1] = re, re
				si[i], si[i+1] = im, -im
				i += 2
			}
			v := zeros(nshifts/2, 3, 3)
			u := zeros(2*nshifts, 2*nshifts, 2*nshifts)
			nh := n
			wh := zeros(2*nshifts, n, n)
			nv := n
			wv := zeros(n, 2*nshifts, 2*nshifts)
			z := eye(n, n)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					copyGeneral(h, hCopy)
					impl.Dlaqr5(true, true, 1, n, 0, n-1, nshifts, sr, si, h.Data, h.Stride, 0, n-1, z.Data, z.Stride,
						v.Data, v.Stride, u.Data, u.Stride, nh, wv.Data, wv.Stride, nv, wh.Data, wh.Stride)
				}
			})
		}
	}
}
