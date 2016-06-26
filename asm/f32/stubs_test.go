// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32

import (
	"math"
	"testing"
)

var (
	nan = float32(math.NaN())
	inf = float32(math.Inf(1))
)

var tests = []struct {
	incX, incY, incDst uintptr
	ix, iy, idst       uintptr
	a                  float32
	dst, x, y          []float32
	ex                 []float32
}{
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   3,
		dst: []float32{5},
		x:   []float32{2},
		y:   []float32{1},
		ex:  []float32{7}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   5,
		dst: []float32{0, 0, 0},
		x:   []float32{0, 0, 0},
		y:   []float32{1, 1, 1},
		ex:  []float32{1, 1, 1}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   5,
		dst: []float32{0, 0, 0},
		x:   []float32{0, 0},
		y:   []float32{1, 1, 1},
		ex:  []float32{1, 1}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   -1,
		dst: []float32{-1, -1, -1},
		x:   []float32{1, 1, 1},
		y:   []float32{1, 2, 1},
		ex:  []float32{0, 1, 0}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   -1,
		dst: []float32{1, 1, 1},
		x:   []float32{1, 2, 1},
		y:   []float32{-1, -2, -1},
		ex:  []float32{-2, -4, -2}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   2.5,
		dst: []float32{1, 1, 1, 1, 1},
		x:   []float32{1, 2, 3, 2, 1},
		y:   []float32{0, 0, 0, 0, 0},
		ex:  []float32{2.5, 5, 7.5, 5, 2.5}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0, // Run big test twice, once aligned once unaligned.
		a:   16.5,
		dst: make([]float32, 20),
		x:   []float32{.5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5, .5},
		y:   []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		ex:  []float32{9.25, 10.25, 11.25, 12.25, 13.25, 14.25, 15.25, 16.25, 17.25, 18.25, 9.25, 10.25, 11.25, 12.25, 13.25, 14.25, 15.25, 16.25, 17.25, 18.25}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   16.5,
		dst: make([]float32, 10),
		x:   []float32{.5, .5, .5, .5, .5, .5, .5, .5, .5, .5},
		y:   []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		ex:  []float32{9.25, 10.25, 11.25, 12.25, 13.25, 14.25, 15.25, 16.25, 17.25, 18.25}},
}

func guardVector(v []float32, g float32, g_ln int) (guarded []float32) {
	guarded = make([]float32, len(v)+g_ln*2)
	copy(guarded[g_ln:], v)
	for i := 0; i < g_ln; i++ {
		guarded[i] = g
		guarded[len(guarded)-1-i] = g
	}
	return guarded
}

func isValidGuard(v []float32, g float32, g_ln int) bool {
	for i := 0; i < g_ln; i++ {
		if v[i] != g || v[len(v)-1-i] != g {
			return false
		}
	}
	return true
}

func same(x, y float32) bool {
	a, b := float64(x), float64(y)
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

func TestAxpyUnitary(t *testing.T) {
	var x_gd, y_gd float32 = 1, 1
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardVector(v.x, x_gd, xg_ln), guardVector(v.y, y_gd, yg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		AxpyUnitary(v.a, x, y)
		for i := range v.ex {
			if !same(y[i], v.ex[i]) {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", j, i, y[i], v.ex[i])
			}
		}
		if !isValidGuard(v.x, x_gd, xg_ln) {
			t.Errorf("Test %d Guard violated in x vector %v %v", j, v.x[:xg_ln], v.x[len(v.x)-xg_ln:])
		}
		if !isValidGuard(v.y, y_gd, yg_ln) {
			t.Errorf("Test %d Guard violated in y vector %v %v", j, v.y[:yg_ln], v.y[len(v.y)-yg_ln:])
		}
	}
}

func TestAxpyUnitaryTo(t *testing.T) {
	var x_gd, y_gd, dst_gd float32 = 1, 1, 0
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardVector(v.x, x_gd, xg_ln), guardVector(v.y, y_gd, yg_ln)
		v.dst = guardVector(v.dst, dst_gd, xg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		dst := v.dst[xg_ln : len(v.dst)-xg_ln]
		AxpyUnitaryTo(dst, v.a, x, y)
		for i := range v.ex {
			if !same(v.ex[i], dst[i]) {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", j, i, dst[i], v.ex[i])
			}
		}
		if !isValidGuard(v.x, x_gd, xg_ln) {
			t.Errorf("Test %d Guard violated in x vector %v %v", j, v.x[:xg_ln], v.x[len(v.x)-xg_ln:])
		}
		if !isValidGuard(v.y, y_gd, yg_ln) {
			t.Errorf("Test %d Guard violated in y vector %v %v", j, v.y[:yg_ln], v.y[len(v.y)-yg_ln:])
		}
		if !isValidGuard(v.dst, dst_gd, xg_ln) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", j, v.dst[:xg_ln], v.dst[len(v.dst)-xg_ln:])
		}
	}
}

func guardIncVector(v []float32, g float32, incV uintptr, g_ln int) (guarded []float32) {
	inc := int(incV)
	s_ln := len(v) * (inc)
	guarded = make([]float32, s_ln+g_ln*2)
	for i, j := 0, 0; i < len(guarded); i++ {
		switch {
		case i < g_ln, i > g_ln+s_ln:
			guarded[i] = g
		case (i-g_ln)%(inc) == 0 && j < len(v):
			guarded[i] = v[j]
			j++
		default:
			guarded[i] = g
		}
	}
	return guarded
}

func checkValidIncGuard(t *testing.T, v []float32, g float32, incV uintptr, g_ln int) {
	inc := int(incV)
	s_ln := len(v) - 2*g_ln
	for i := range v {
		switch {
		case same(v[i], g):
			// Correct value
		case i < g_ln:
			t.Errorf("Front guard violated at %d %v", i, v[:g_ln])
		case i > g_ln+s_ln:
			t.Errorf("Back guard violated at %d %v", i-g_ln-s_ln, v[g_ln+s_ln:])
		case (i-g_ln)%inc == 0 && (i-g_ln)/inc < len(v):
			// Ignore input values
		default:
			t.Errorf("Internal guard violated at %d %v", i-g_ln, v[g_ln:g_ln+s_ln])
		}
	}
}

func TestAxpyInc(t *testing.T) {
	var x_gd, y_gd float32 = 1, 1
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardIncVector(v.x, x_gd, uintptr(v.incX), xg_ln), guardIncVector(v.y, y_gd, uintptr(v.incY), yg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		AxpyInc(v.a, x, y, uintptr(len(v.ex)), v.incX, v.incY, v.ix, v.iy)
		for i := range v.ex {
			if !same(y[i*int(v.incY)], v.ex[i]) {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", j, i, y[i*int(v.incY)], v.ex[i])
			}
		}
		checkValidIncGuard(t, v.x, x_gd, uintptr(v.incX), xg_ln)
		checkValidIncGuard(t, v.y, y_gd, uintptr(v.incY), yg_ln)
	}
}

func TestAxpyIncTo(t *testing.T) {
	var x_gd, y_gd, dst_gd float32 = 1, 1, 0
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardIncVector(v.x, x_gd, uintptr(v.incX), xg_ln), guardIncVector(v.y, y_gd, uintptr(v.incY), yg_ln)
		v.dst = guardIncVector(v.dst, dst_gd, uintptr(v.incDst), xg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		dst := v.dst[xg_ln : len(v.dst)-xg_ln]
		AxpyIncTo(dst, v.incDst, v.idst, v.a, x, y, uintptr(len(v.ex)), v.incX, v.incY, v.ix, v.iy)
		for i := range v.ex {
			if !same(dst[i*int(v.incDst)], v.ex[i]) {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", j, i, dst[i*int(v.incDst)], v.ex[i])
			}
		}
		checkValidIncGuard(t, v.x, x_gd, uintptr(v.incX), xg_ln)
		checkValidIncGuard(t, v.y, y_gd, uintptr(v.incY), yg_ln)
		checkValidIncGuard(t, v.dst, dst_gd, uintptr(v.incDst), xg_ln)
	}
}
