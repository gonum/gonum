// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64_test

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/internal/asm/f64"
)

func TestL1NormIncIEEE(t *testing.T) {
	for _, n := range []int{1, 2, 3, 7, 8, 9, 15, 16, 17, 31, 32, 33} {
		for _, inc := range []int{1, 2, 3, 7} {
			for _, value := range []float64{0, math.Copysign(0, -1), -0.25, math.Inf(1), math.Inf(-1), math.NaN(), math.SmallestNonzeroFloat64} {
				t.Run(fmt.Sprintf("n=%d/inc=%d/value=%g", n, inc, value), func(t *testing.T) {
					x := make([]float64, 1+(n-1)*inc)
					for i := range x {
						x[i] = math.NaN()
					}
					var want float64
					for i := 0; i < n; i++ {
						x[i*inc] = value
						want += math.Abs(value)
					}
					got := f64.L1NormInc(x, n, inc)
					if math.IsNaN(want) {
						if !math.IsNaN(got) {
							t.Fatalf("got %g, want NaN", got)
						}
					} else if math.Float64bits(got) != math.Float64bits(want) {
						t.Fatalf("got %g, want %g", got, want)
					}
				})
			}
		}
	}
}

func TestL1NormIncZeroStride(t *testing.T) {
	for _, n := range []int{0, 1, 2, 7, 8, 9, 17} {
		for _, x := range [][]float64{nil, {}, {-3}, {math.NaN()}, {math.Inf(1)}} {
			if got := f64.L1NormInc(x, n, 0); math.Float64bits(got) != 0 {
				t.Fatalf("n=%d x=%v: got %g, want +0", n, x, got)
			}
		}
	}
}
