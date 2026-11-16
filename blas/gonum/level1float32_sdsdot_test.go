// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"
	"testing"
)

func TestSdsdot(t *testing.T) {
	for _, test := range []struct {
		name  string
		alpha float32
		x, y  []float32
		want  float32
	}{
		{name: "empty", alpha: 1.5, want: 1.5},
		{name: "empty-negative", alpha: -1.5, want: -1.5},
		{name: "empty-infinity", alpha: float32(math.Inf(1)), want: float32(math.Inf(1))},
		// The exact sum is 1 + 2^-24 + 2^-48, just above the halfway point.
		// Rounding the dot product before adding alpha would round down to 1.
		{name: "round-once", alpha: 1, x: []float32{1, 1}, y: []float32{0x1p-24, 0x1p-48}, want: 1 + 0x1p-23},
		{name: "round-once-negative", alpha: -1, x: []float32{-1, -1}, y: []float32{0x1p-24, 0x1p-48}, want: -1 - 0x1p-23},
		// The dot product exceeds the float32 range, but adding alpha brings
		// the final result back into range before the conversion to float32.
		{name: "finite-cancellation", alpha: -math.MaxFloat32, x: []float32{math.MaxFloat32, math.MaxFloat32}, y: []float32{1, 1}, want: math.MaxFloat32},
		{name: "ordinary", alpha: 2, x: []float32{1, 2}, y: []float32{3, 4}, want: 13},
	} {
		for _, incX := range []int{-2, -1, 1, 2} {
			for _, incY := range []int{-2, -1, 1, 2} {
				x := sdsdotStrided(test.x, incX)
				y := sdsdotStrided(test.y, incY)
				got := (Implementation{}).Sdsdot(len(test.x), test.alpha, x, incX, y, incY)
				if got != test.want {
					t.Errorf("%s incX=%d incY=%d: got %g, want %g", test.name, incX, incY, got, test.want)
				}
			}
		}
	}
}

func sdsdotStrided(values []float32, inc int) []float32 {
	if len(values) == 0 {
		return nil
	}
	step := inc
	if step < 0 {
		step = -step
	}
	result := make([]float32, 1+(len(values)-1)*step)
	i := 0
	if inc < 0 {
		i = len(result) - 1
	}
	for _, value := range values {
		result[i] = value
		i += inc
	}
	return result
}
