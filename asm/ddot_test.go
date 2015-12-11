// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
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
			xData: []float64{2, 3},
			yData: []float64{-3, 4},
			want:  6,
		},
		{
			xData: []float64{2, 3, -4},
			yData: []float64{-3, 4, 5},
			want:  -14,
		},
		{
			xData: []float64{2, 3, -4, -5},
			yData: []float64{-3, 4, 5, -6},
			want:  16,
		},
		{
			xData: []float64{0, 2, 3, -4, -5},
			yData: []float64{0, -3, 4, 5, -6},
			want:  16,
		},
		{
			xData: []float64{0, 0, 2, 3, -4, -5},
			yData: []float64{0, 1, -3, 4, 5, -6},
			want:  16,
		},
		{
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  4,
		},
		{
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  39,
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		got := DdotUnitary(x, y)

		msg := "test %v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"
		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msg, i, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msg, i, "y", yFront, yBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("test %v: modified read-only x argument", i)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("test %v: modified read-only y argument", i)
		}
		if math.IsNaN(got) {
			t.Errorf("test %v: invalid memory read", i)
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
			xData:   []float64{2, 3},
			yData:   []float64{-3, 4},
			want:    6,
			wantRev: -1,
		},
		{
			xData:   []float64{2, 3, -4},
			yData:   []float64{-3, 4, 5},
			want:    -14,
			wantRev: 34,
		},
		{
			xData:   []float64{2, 3, -4, -5},
			yData:   []float64{-3, 4, 5, -6},
			want:    16,
			wantRev: 2,
		},
		{
			xData:   []float64{0, 2, 3, -4, -5},
			yData:   []float64{0, -3, 4, 5, -6},
			want:    16,
			wantRev: 34,
		},
		{
			xData:   []float64{0, 0, 2, 3, -4, -5},
			yData:   []float64{0, 1, -3, 4, 5, -6},
			want:    16,
			wantRev: -5,
		},
		{
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    4,
			wantRev: -4,
		},
		{
			xData:   []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:    39,
			wantRev: 3,
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

				prefix := fmt.Sprintf("test %v, incX = %v, incY = %v", i, incX, incY)
				msg := "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"
				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msg, prefix, "x", xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf(msg, prefix, "y", yFront, yBack)
				}
				if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}
				if nonStridedWrite(y, incY) || !equalStrided(test.yData, y, incY) {
					t.Errorf("%v: modified read-only y argument", prefix)
				}
				if math.IsNaN(got) {
					t.Errorf("%v: invalid memory read", prefix)
					continue
				}

				want := test.want
				if incX*incY < 0 {
					want = test.wantRev
				}
				if got != want {
					t.Errorf("%v: unexpected result. want %v, got %v", prefix, want, got)
				}
			}
		}
	}
}
