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

var tests = []struct {
	incX, incY, incDst uintptr
	ix, iy, idst       uintptr
	a                  float32
	dst, x, y          []float32
	ex                 []float32
}{
	{2, 2, 3, 0, 0, 0,
		3,
		[]float32{5},
		[]float32{2},
		[]float32{1},
		[]float32{7}},
	{2, 2, 3, 0, 0, 0,
		5,
		[]float32{0, 0, 0},
		[]float32{0, 0, 0},
		[]float32{1, 1, 1},
		[]float32{1, 1, 1}},
	{2, 2, 3, 0, 0, 0,
		5,
		[]float32{0, 0, 0},
		[]float32{0, 0},
		[]float32{1, 1, 1},
		[]float32{1, 1}},
	{2, 2, 3, 0, 0, 0,
		-1,
		[]float32{-1, -1, -1},
		[]float32{1, 1, 1},
		[]float32{1, 2, 1},
		[]float32{0, 1, 0}},
	{2, 2, 3, 0, 0, 0,
		-1,
		[]float32{1, 1, 1},
		[]float32{1, 2, 1},
		[]float32{-1, -2, -1},
		[]float32{-2, -4, -2}},
	{2, 2, 3, 0, 0, 0,
		2.5,
		[]float32{1, 1, 1, 1, 1},
		[]float32{1, 2, 3, 2, 1},
		[]float32{0, 0, 0, 0, 0},
		[]float32{2.5, 5, 7.5, 5, 2.5}},
	// Run big test twice, once aligned once unaligned.
	{2, 2, 3, 0, 0, 0,
		16.5,
		make([]float32, 10),
		[]float32{.5, .5, .5, .5, .5, .5, .5, .5, .5, .5},
		[]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]float32{9.25, 10.25, 11.25, 12.25, 13.25, 14.25, 15.25, 16.25, 17.25, 18.25}},
	{2, 2, 3, 0, 0, 0,
		16.5,
		make([]float32, 10),
		[]float32{.5, .5, .5, .5, .5, .5, .5, .5, .5, .5},
		[]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]float32{9.25, 10.25, 11.25, 12.25, 13.25, 14.25, 15.25, 16.25, 17.25, 18.25}},
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

func validGuard(v []float32, g float32, g_ln int) bool {
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
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardVector(v.x, 1, g_ln), guardVector(v.y, 1, g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		AxpyUnitary(v.a, x, y)
		for i := range v.ex {
			if !same(y[i], v.ex[i]) {
				t.Error("Test", j, "Unexpected result at", i, "Got:", y[i], "Expected:", v.ex[i])
			}
		}
		if !validGuard(v.x, 1, g_ln) {
			t.Error("Test", j, "Guard violated in x vector", v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
		if !validGuard(v.y, 1, g_ln) {
			t.Error("Test", j, "Guard violated in y vector", v.y[:g_ln], v.y[len(v.x)-g_ln:])
		}
		runtime.GC()
	}
}

func TestAxpyUnitaryTo(t *testing.T) {
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardVector(v.x, 1, g_ln), guardVector(v.y, 1, g_ln)
		v.dst = guardVector(v.dst, 0, g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		dst := v.dst[g_ln : len(v.dst)-g_ln]
		AxpyUnitaryTo(dst, v.a, x, y)
		for i := range v.ex {
			if !same(v.ex[i], dst[i]) {
				t.Error("Test", j, "Unexpected result at", i, "Got:", dst[i], "Expected:", v.ex[i])
			}
		}
		if !validGuard(v.x, 1, g_ln) {
			t.Error("Test", j, "Guard violated in x vector", v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
		if !validGuard(v.y, 1, g_ln) {
			t.Error("Test", j, "Guard violated in y vector", v.y[:g_ln], v.y[len(v.x)-g_ln:])
		}
		if !validGuard(v.dst, 0, g_ln) {
			t.Error("Test", j, "Guard violated in x vector", v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
		runtime.GC()
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

func validIncGuard(t *testing.T, v []float32, g float32, incV uintptr, g_ln int) {
	inc := int(incV)
	s_ln := len(v) - 2*g_ln
	for i := range v {
		switch {
		case same(v[i], g):
			// Correct value
		case i < g_ln:
			t.Error("Front guard violated at", i, v[:g_ln])
		case i > g_ln+s_ln:
			t.Error("Back guard violated at", i-g_ln-s_ln, v[g_ln+s_ln:])
		case (i-g_ln)%inc == 0 && (i-g_ln)/inc < len(v):
		default:
			t.Error("Internal guard violated at", i-g_ln, v[g_ln:g_ln+s_ln])
		}
	}
}

func TestAxpyInc(t *testing.T) {
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardIncVector(v.x, 1, v.incX, g_ln), guardIncVector(v.y, 1, v.incY, g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		AxpyInc(v.a, x, y, uintptr(len(v.ex)), v.incX, v.incY, v.ix, v.iy)
		for i := range v.ex {
			if !same(y[i*int(v.incY)], v.ex[i]) {
				t.Error("Test", j, "Unexpected result at", i, "Got:", y[i*int(v.incY)], "Expected:", v.ex[i])
				t.Error("Result:", y)
				t.Error("Expect:", v.ex)
			}
		}
		validIncGuard(t, v.x, 1, v.incX, g_ln)
		validIncGuard(t, v.y, 1, v.incY, g_ln)
		runtime.GC()
	}
}

func TestAxpyIncTo(t *testing.T) {
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardIncVector(v.x, 1, v.incX, g_ln), guardIncVector(v.y, 1, v.incY, g_ln)
		v.dst = guardIncVector(v.dst, 0, v.incDst, g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		dst := v.dst[g_ln : len(v.dst)-g_ln]
		AxpyIncTo(dst, v.incDst, v.idst, v.a, x, y, uintptr(len(v.ex)), v.incX, v.incY, v.ix, v.iy)
		for i := range v.ex {
			if !same(dst[i*int(v.incDst)], v.ex[i]) {
				t.Error("Test", j, "Unexpected result at", i, "Got:", dst[i*int(v.incDst)], "Expected:", v.ex[i])
				t.Error(v.dst)
				t.Error(v.ex)
			}
		}
		validIncGuard(t, v.x, 1, v.incX, g_ln)
		validIncGuard(t, v.y, 1, v.incY, g_ln)
		validIncGuard(t, v.dst, 0, v.incDst, g_ln)
		runtime.GC()
	}
}
