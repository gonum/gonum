// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && !noasm && !safe && !gccgo

package f64_test

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/internal/asm/f64"
)

func TestL1NormFiniteGroupingAMD64(t *testing.T) {
	for _, n := range []int{1, 2, 7, 8, 9, 15, 16, 17, 31, 32, 33, 127, 128, 129, 4095, 4096, 4097} {
		for _, kind := range []string{"wide", "subnormal", "near-overflow", "large-small"} {
			t.Run(fmt.Sprintf("n%d/%s", n, kind), func(t *testing.T) {
				x := make([]float64, n)
				for i := range x {
					switch kind {
					case "wide":
						x[i] = math.Ldexp(1+float64(i%13)/1024, (73*i)%1901-1000)
					case "subnormal":
						x[i] = math.Float64frombits(uint64(i%127 + 1))
					case "near-overflow":
						x[i] = math.MaxFloat64 / float64(n)
					case "large-small":
						x[i] = 1
						if i%5 == 0 {
							x[i] = 0x1p54
						}
					}
					if i%2 != 0 {
						x[i] = -x[i]
					}
				}
				// The old ASM computes this positive-addition graph for finite
				// inputs. Preserve its rounding, including finite-input overflow.
				var lanes [8]float64
				i := 0
				for ; i+8 <= n; i += 8 {
					for j := range lanes {
						lanes[j] += math.Abs(x[i+j])
					}
				}
				lo := (lanes[0] + lanes[2]) + (lanes[6] + lanes[4])
				hi := (lanes[1] + lanes[3]) + (lanes[7] + lanes[5])
				want := hi + lo
				for ; i < n; i++ {
					want += math.Abs(x[i])
				}
				got := f64.L1Norm(x)
				if math.Float64bits(got) != math.Float64bits(want) {
					t.Fatalf("got %g (%016x), want %g (%016x)", got, math.Float64bits(got), want, math.Float64bits(want))
				}
				stridedWant := ((lanes[4] + lanes[5]) + (lanes[7] + lanes[6])) + ((lanes[0] + lanes[1]) + (lanes[3] + lanes[2]))
				for j := n / 8 * 8; j < n; j++ {
					stridedWant += math.Abs(x[j])
				}
				for _, inc := range []int{1, 3} {
					strided := make([]float64, 1+(n-1)*inc)
					for j, v := range x {
						strided[j*inc] = v
					}
					got := f64.L1NormInc(strided, n, inc)
					if math.Float64bits(got) != math.Float64bits(stridedWant) {
						t.Fatalf("inc=%d: got %g, want %g", inc, got, stridedWant)
					}
				}
			})
		}
	}
}
