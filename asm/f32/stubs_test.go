// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package f32

import (
	"math"
	"runtime"
	"testing"
)

var (
	nan = float32(math.NaN())
	inf = float32(math.Inf(1))
)

func same(x, y float32) bool {
	a, b := float64(x), float64(y)
	return !(x != y && !math.IsNaN(a) && !math.IsNaN(b) || (math.IsNaN(a) != math.IsNaN(b)))
}

func TestAxpyUnitary(t *testing.T) {
	for i, v := range []struct {
		a    float32
		x, y []float32
		ex   []float32
	}{
		{0, []float32{}, []float32{}, []float32{}},
		{nan, []float32{1, 2, 3}, []float32{1, 2, 3, 4}, []float32{nan, nan, nan}},
		{5, []float32{0, 1, 2, 3, 4, 5, 6, 7},
			[]float32{2, 3, 4, 5, 6, 7, 8, 9},
			[]float32{2, 8, 14, 20, 26, 32, 38, 44}},
		{-2, []float32{5, 4, 3}, []float32{1, 3, 5}, []float32{-9, -5, -1}},
	} {
		AxpyUnitary(v.a, v.x, v.y)
		for j := range v.ex {
			if !same(v.ex[j], v.y[j]) {
				t.Error("Test", i, "Unexpected value at", j, "Got:", v.y[j], "Expected:", v.ex[j])
			}
		}
	}
	runtime.GC()
}

func TestAxpyUnitaryTo(t *testing.T) {
	for i, v := range []struct {
		a         float32
		x, y, dst []float32
		ex        []float32
	}{
		{0, []float32{}, []float32{}, []float32{}, []float32{}},
		{nan, []float32{1, 2, 3},
			[]float32{1, 2, 3, 4},
			[]float32{0, 0, 0},
			[]float32{nan, nan, nan}},
		{5, []float32{0, 1, 2, 3, 4, 5, 6, 7},
			[]float32{2, 3, 4, 5, 6, 7, 8, 9},
			make([]float32, 8),
			[]float32{2, 8, 14, 20, 26, 32, 38, 44}},
		{-2, []float32{5, 4, 3},
			[]float32{1, 3, 5},
			[]float32{0, 0, 0},
			[]float32{-9, -5, -1}},
	} {
		AxpyUnitaryTo(v.dst, v.a, v.x, v.y)
		for j := range v.ex {
			if !same(v.ex[j], v.dst[j]) {
				t.Error("Test", i, "Unexpected value at", j, "Got:", v.dst[j], "Expected:", v.ex[j])
			}
		}
	}
	runtime.GC()
}

// func AxpyInc(alpha float32, x, y []float32, n, incX, incY, ix, iy uintptr)
func TestAxpyInc(t *testing.T) {
	for i, v := range []struct {
		a                     float32
		x, y                  []float32
		ex                    []float32
		n, incX, incY, ix, iy uintptr
	}{
		{0, []float32{}, []float32{}, []float32{}, 0, 10, 10, 5, 5},
		{nan, []float32{1, 2, 3},
			[]float32{1, 2, 3, 4},
			[]float32{nan, nan, nan}, 3, 1, 1, 0, 0},
		{0, []float32{1, 2, 3},
			[]float32{1, 2, 3, 4},
			[]float32{nan, nan, 3}, 1, 1, 1, 2, 2},
		/*{5, []float32{0, 1, 2, 3, 4, 5, 6, 7},
			[]float32{2, 3, 4, 5, 6, 7, 8, 9},
			make([]float32, 8),
			[]float32{2, 8, 14, 20, 26, 32, 38, 44}},
		{-2, []float32{5, 4, 3},
			[]float32{1, 3, 5},
			[]float32{0, 0, 0},
			[]float32{-9, -5, -1}},*/
	} {
		AxpyInc(v.a, v.x, v.y, v.n, v.incX, v.incY, v.ix, v.iy)
		for j, k := v.iy, 0; k < int(v.n); j, k = j+v.incY, k+1 {
			if !same(v.ex[j], v.y[j]) {
				t.Error("Test", i, "Unexpected value at", j, "Got:", v.y[j], "Expected:", v.ex[j])
			}
		}
	}
	runtime.GC()
}
