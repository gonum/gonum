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

func TestL1NormIEEEAccumulation(t *testing.T) {
	lengths := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 31, 32, 33, 127, 128, 129, 4095, 4096, 4097}
	patterns := []string{"zero", "negative-zero", "dyadic", "all-positive-inf", "all-negative-inf", "mixed-inf", "paired-positive-inf", "paired-negative-inf", "same-lane-inf", "nan-first", "nan-last", "inf-then-nan", "nan-then-inf"}
	for _, n := range lengths {
		for _, offset := range []int{0, 1} {
			for _, pattern := range patterns {
				t.Run(fmt.Sprintf("n%d/offset%d/%s", n, offset, pattern), func(t *testing.T) {
					backing := make([]float64, n+2)
					x := backing[offset : offset+n]
					for i := range x {
						switch pattern {
						case "negative-zero":
							x[i] = math.Copysign(0, -1)
						case "dyadic":
							x[i] = float64(i%17-8) / 4
						case "all-positive-inf":
							x[i] = math.Inf(1)
						case "all-negative-inf":
							x[i] = math.Inf(-1)
						case "mixed-inf":
							x[i] = math.Inf(2*(i%2) - 1)
						}
					}
					if n > 0 {
						nan := math.Float64frombits(0xfff8000000000042)
						switch pattern {
						case "paired-positive-inf":
							x[0], x[n-1] = math.Inf(1), math.Inf(1)
						case "paired-negative-inf":
							x[0], x[n-1] = math.Inf(-1), math.Inf(-1)
						case "same-lane-inf":
							for i := 0; i < n; i += 8 {
								x[i] = math.Inf(1)
							}
						case "nan-first":
							x[0] = nan
						case "nan-last":
							x[n-1] = nan
						case "inf-then-nan":
							x[0] = math.Inf(1)
							x[n-1] = nan
						case "nan-then-inf":
							x[n-1] = math.Inf(1)
							x[0] = nan
						}
					}
					before := make([]uint64, len(backing))
					for i, v := range backing {
						before[i] = math.Float64bits(v)
					}
					// Finite values in these fixtures are exact dyadic sums.
					var want float64
					for _, v := range x {
						want += math.Abs(v)
					}
					got := f64.L1Norm(x)
					if math.IsNaN(want) {
						if !math.IsNaN(got) {
							t.Fatalf("got %g, want NaN", got)
						}
					} else if math.Float64bits(got) != math.Float64bits(want) {
						t.Fatalf("got %g (%016x), want %g (%016x)", got, math.Float64bits(got), want, math.Float64bits(want))
					}
					for i, v := range backing {
						if math.Float64bits(v) != before[i] {
							t.Fatalf("input or guard modified at %d", i)
						}
					}
				})
			}
		}
	}
	if got := f64.L1Norm(nil); math.Float64bits(got) != 0 {
		t.Fatalf("nil returned %g", got)
	}
}
