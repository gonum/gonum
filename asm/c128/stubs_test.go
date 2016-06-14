// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128

import (
	"runtime"
	"testing"
)

var tests = []struct {
	incX, incY, incDst int
	ix, iy, idst       uintptr
	a                  complex128
	dst, x, y          []complex128
	ex                 []complex128
}{
	{2, 2, 3, 0, 0, 0,
		1 + 1i,
		[]complex128{5},
		[]complex128{1},
		[]complex128{1i},
		[]complex128{1 + 2i}},
	{2, 2, 3, 0, 0, 0,
		1 + 2i,
		[]complex128{0, 0, 0},
		[]complex128{0, 0, 0},
		[]complex128{1, 1, 1},
		[]complex128{1, 1, 1}},
	{2, 2, 3, 0, 0, 0,
		1 + 2i,
		[]complex128{0, 0, 0},
		[]complex128{0, 0},
		[]complex128{1, 1, 1},
		[]complex128{1, 1}},
	{2, 2, 3, 0, 0, 0,
		1 + 2i,
		[]complex128{1i, 1i, 1i},
		[]complex128{1i, 1i, 1i},
		[]complex128{1, 2, 1},
		[]complex128{-1 + 1i, 1i, -1 + 1i}},
	{2, 2, 3, 0, 0, 0,
		-1i,
		[]complex128{1i, 1i, 1i},
		[]complex128{1i, 1i, 1i},
		[]complex128{1, 2, 1},
		[]complex128{2, 3, 2}},
	{2, 2, 3, 0, 0, 0,
		-1i,
		[]complex128{1i, 1i, 1i},
		[]complex128{1i, 1i, 1i, 1i, 1i}[1:4],
		[]complex128{1, 1, 2, 1, 1}[1:4],
		[]complex128{2, 3, 2}},
	// Run big test twice, once aligned once unaligned.
	{2, 2, 3, 0, 0, 0,
		1 - 1i,
		make([]complex128, 10),
		[]complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		[]complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		[]complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	{2, 2, 3, 0, 0, 0,
		1 - 1i,
		make([]complex128, 10),
		[]complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		[]complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		[]complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	{-2, -2, -3, 18, 18, 27,
		1 - 1i,
		make([]complex128, 10),
		[]complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		[]complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		[]complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	{-2, 2, -3, 18, 0, 27,
		1 - 1i,
		make([]complex128, 10),
		[]complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		[]complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		[]complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
}

func guardVector(v []complex128, g complex128, g_ln int) (guarded []complex128) {
	guarded = make([]complex128, len(v)+g_ln*2)
	copy(guarded[g_ln:], v)
	for i := 0; i < g_ln; i++ {
		guarded[i] = g
		guarded[len(guarded)-1-i] = g
	}
	return guarded
}

func validGuard(v []complex128, g complex128, g_ln int) bool {
	for i := 0; i < g_ln; i++ {
		if v[i] != g || v[len(v)-1-i] != g {
			return false
		}
	}
	return true
}

func TestAxpyUnitary(t *testing.T) {
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardVector(v.x, 1, g_ln), guardVector(v.y, 1, g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		AxpyUnitary(v.a, x, y)
		for i := range v.ex {
			if y[i] != v.ex[i] {
				t.Error("Test", j, "Unexpected result at", i, "Got:", y[i], "Expected:", v.ex[i])
				t.Error(v.y)
				t.FailNow()
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
			if dst[i] != v.ex[i] {
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

func guardIncVector(v []complex128, g complex128, incV uintptr, g_ln int) (guarded []complex128) {
	inc := int(incV)
	s_ln := len(v) * inc
	if inc < 0 {
		s_ln = len(v) * -inc
	}
	guarded = make([]complex128, s_ln+g_ln*2)
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

func validIncGuard(t *testing.T, v []complex128, g complex128, incV uintptr, g_ln int) {
	inc := int(incV)
	s_ln := len(v) - 2*g_ln
	if inc < 0 {
		s_ln = len(v) * -inc
	}

	for i := range v {
		switch {
		case v[i] == g:
			// Correct value
		case i < g_ln:
			t.Error("Front guard violated at", i, v[:g_ln])
		case i > g_ln+s_ln:
			t.Error("Back guard violated at", i-g_ln-s_ln, v[g_ln+s_ln:])
		case (i-g_ln)%inc == 0 && (i-g_ln)/inc < len(v):
		case v[i] != g:
			t.Error("Internal guard violated at", i-g_ln, v[g_ln:g_ln+s_ln])
		default:
		}
	}
}

func TestAxpyInc(t *testing.T) {
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardIncVector(v.x, 1, uintptr(v.incX), g_ln), guardIncVector(v.y, 1, uintptr(v.incY), g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		AxpyInc(v.a, x, y, uintptr(len(v.ex)), uintptr(v.incX), uintptr(v.incY), v.ix, v.iy)
		for i := range v.ex {
			if y[int(v.iy)+i*int(v.incY)] != v.ex[i] {
				t.Error("Test", j, "Unexpected result at", i, "Got:", y[i*int(v.incY)], "Expected:", v.ex[i])
				t.Error("Result:", y)
				t.Error("Expect:", v.ex)
			}
		}
		validIncGuard(t, v.x, 1, uintptr(v.incX), g_ln)
		validIncGuard(t, v.y, 1, uintptr(v.incY), g_ln)
		runtime.GC()
	}
}

func TestAxpyIncTo(t *testing.T) {
	for j, v := range tests {
		g_ln := 4 + j%2
		v.x, v.y = guardIncVector(v.x, 1, uintptr(v.incX), g_ln), guardIncVector(v.y, 1, uintptr(v.incY), g_ln)
		v.dst = guardIncVector(v.dst, 0, uintptr(v.incDst), g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		dst := v.dst[g_ln : len(v.dst)-g_ln]
		AxpyIncTo(dst, uintptr(v.incDst), v.idst, v.a, x, y, uintptr(len(v.ex)), uintptr(v.incX), uintptr(v.incY), v.ix, v.iy)
		for i := range v.ex {
			if dst[int(v.idst)+i*int(v.incDst)] != v.ex[i] {
				t.Error("Test", j, "Unexpected result at", i, "Got:", dst[i*int(v.incDst)], "Expected:", v.ex[i])
			}
		}
		validIncGuard(t, v.x, 1, uintptr(v.incX), g_ln)
		validIncGuard(t, v.y, 1, uintptr(v.incY), g_ln)
		validIncGuard(t, v.dst, 0, uintptr(v.incDst), g_ln)
		runtime.GC()
	}
}
