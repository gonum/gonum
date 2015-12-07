// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"math"
	"testing"
)

func TestDdotUnitary(t *testing.T) {
	for i, test := range []struct {
		xData []float64
		yData []float64

		want float64
	}{
		{
			xData: []float64{2},
			yData: []float64{-3},
			want:  -6,
		},
		{
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  4,
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		got := DdotUnitary(x, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf("test %v: out-of-bounds write to first argument\nfront guard: %v\nback guard: %v\n",
				i, xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf("test %v: out-of-bounds write to second argument\nfront guard: %v\nback guard: %v\n",
				i, yFront, yBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("test %v: modified first argument", i)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("test %v: modified second argument", i)
		}
		if math.IsNaN(got) {
			t.Errorf("test %v: invalid memory access", i)
			continue
		}

		if got != test.want {
			t.Errorf("test %v: unexpected result. want %v, got %v", i, test.want, got)
		}
	}
}

func TestDdotInc(t *testing.T) {
	for i, test := range []struct {
		xData []float64
		yData []float64

		want    float64
		wantRev float64 // Result when one of the vectors is reversed.
	}{
		{
			xData:   []float64{2},
			yData:   []float64{-3},
			want:    -6,
			wantRev: -6,
		},
		{
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    4,
			wantRev: -4,
		},
	} {
		for _, incX := range []int{-7, -3, -2, -1, 1, 2, 3, 7} {
			for _, incY := range []int{-7, -3, -2, -1, 1, 2, 3, 7} {
				n := len(test.xData)
				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)

				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}
				got := DdotInc(x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf("test %v: out-of-bounds write to first argument. incX = %v, incY = %v\nfront guard: %v\nback guard: %v\n",
						i, incX, incY, xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf("test %v: out-of-bounds write to second argument. incX = %v, incY = %v\nfront guard: %v\nback guard: %v\n",
						i, incX, incY, xFront, yBack)
				}
				if !equalStrided(test.xData, x, incX) {
					t.Errorf("test %v: modified first argument. incX = %v, incY = %v", i, incX, incY)
				}
				if !equalStrided(test.yData, y, incY) {
					t.Errorf("test %v: modified second argument. incX = %v, incY = %v", i, incX, incY)
				}
				if math.IsNaN(got) {
					t.Errorf("test %v: invalid memory read. incX = %v, incY = %v", i, incX, incY)
					continue
				}

				want := test.want
				if incX*incY < 0 {
					want = test.wantRev
				}
				if got != want {
					t.Errorf("test %v: unexpected result. incX = %v, incY = %v, want %v, got %v",
						i, incX, incY, want, got)
				}
			}
		}
	}
}

// newGuardedVector allocates a new slice and returns it as three subslices.
// v is a strided vector that contains elements of data at indices i*inc and
// NaN elsewhere. frontGuard and backGuard are filled with NaN values, and
// their backing arrays are directly adjacent to v in memory. The three slices
// can be used to detect invalid memory reads and writes.
func newGuardedVector(data []float64, inc int) (v, frontGuard, backGuard []float64) {
	if inc < 0 {
		inc = -inc
	}
	guard := 2 * inc
	size := (len(data)-1)*inc + 1
	whole := make([]float64, size+2*guard)
	v = whole[guard : len(whole)-guard]
	for i := range whole {
		whole[i] = math.NaN()
	}
	for i, d := range data {
		v[i*inc] = d
	}
	return v, whole[:guard], whole[len(whole)-guard:]
}

// allNaN returns true if x contains only NaN values, and false otherwise.
func allNaN(x []float64) bool {
	for _, v := range x {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// equalStrided returns true if the strided vector x contains elements of the
// dense vector ref at indices i*inc and NaN values elsewhere, false otherwise.
func equalStrided(ref, x []float64, inc int) bool {
	if inc < 0 {
		inc = -inc
	}
	for i, v := range x {
		if i%inc == 0 {
			if ref[i/inc] != v {
				return false

			}
		} else if !math.IsNaN(v) {
			return false
		}
	}
	return true
}
