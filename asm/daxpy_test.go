// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import "testing"

func TestDaxpyUnitary(t *testing.T) {
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 1,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-1},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 1,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 1, 4, -2, 2, -10},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
	} {
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		z, zFront, zBack := newGuardedVector(test.xData, 1)

		DaxpyUnitary(test.alpha, x, y, z)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf("test %v: out-of-bounds write to x argument\nfront guard: %v\nback guard: %v\n",
				i, xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf("test %v: out-of-bounds write to y argument\nfront guard: %v\nback guard: %v\n",
				i, yFront, yBack)
		}
		if !allNaN(zFront) || !allNaN(zBack) {
			t.Errorf("test %v: out-of-bounds write to z argument\nfront guard: %v\nback guard: %v\n",
				i, zFront, zBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("test %v: modified x argument", i)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("test %v: modified y argument", i)
		}

		if !equalStrided(test.want, z, 1) {
			t.Errorf("test %v: unexpected result:\nwant: %v\ngot: %v", i, test.want, z)
		}
	}
}
