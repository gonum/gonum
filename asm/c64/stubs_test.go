// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c64

import (
	"runtime"
	"testing"
)

func guardVector(v []complex64, g complex64, g_ln int) (guarded []complex64) {
	guarded = make([]complex64, len(v)+g_ln*2)
	copy(guarded[g_ln:], v)
	for i := 0; i < g_ln; i++ {
		guarded[i] = g
		guarded[len(guarded)-1-i] = g
	}
	return guarded
}

func validGuard(v []complex64, g complex64, g_ln int) bool {
	for i := 0; i < g_ln; i++ {
		if v[i] != g || v[len(v)-1-i] != g {
			return false
		}
	}
	return true
}

func TestAxpyUnitary(t *testing.T) {
	for j, v := range []struct {
		a    complex64
		x, y []complex64
		ex   []complex64
	}{
		{1 + 1i, []complex64{1}, []complex64{1i}, []complex64{1 + 2i}},
		{1 + 2i, []complex64{0, 0, 0}, []complex64{1, 1, 1}, []complex64{1, 1, 1}},
		{1 + 2i, []complex64{0, 0}, []complex64{1, 1, 1}, []complex64{1, 1}},
		{1 + 2i, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{-1 + 1i, 1i, -1 + 1i}},
		{-1i, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{2, 3, 2}},
		{-1i, []complex64{1i, 1i, 1i, 1i, 1i}, []complex64{1, 1, 2, 1, 1}, []complex64{2, 2, 3, 2, 2}},
		// Run big test twice, once aligned once unaligned.
		{1 - 1i,
			[]complex64{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
			[]complex64{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
			[]complex64{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
		{1 - 1i,
			[]complex64{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
			[]complex64{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
			[]complex64{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	} {
		g_ln := 4 + j%2
		v.x, v.y = guardVector(v.x, 1, g_ln), guardVector(v.y, 1, g_ln)
		x, y := v.x[g_ln:len(v.x)-g_ln], v.y[g_ln:len(v.y)-g_ln]
		AxpyUnitary(v.a, x, y)
		for i := range v.ex {
			if y[i] != v.ex[i] {
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
	for j, v := range []struct {
		a         complex64
		dst, x, y []complex64
		ex        []complex64
	}{
		{1 + 1i, []complex64{5}, []complex64{1}, []complex64{1i}, []complex64{1 + 2i}},
		{1 + 2i, []complex64{0, 0, 0}, []complex64{0, 0, 0}, []complex64{1, 1, 1}, []complex64{1, 1, 1}},
		{1 + 2i, []complex64{0, 0, 0}, []complex64{0, 0}, []complex64{1, 1, 1}, []complex64{1, 1}},
		{1 + 2i, []complex64{1i, 1i, 1i}, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{-1 + 1i, 1i, -1 + 1i}},
		{-1i, []complex64{1i, 1i, 1i}, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{2, 3, 2}},
		{-1i, []complex64{1i, 1i, 1i}, []complex64{1i, 1i, 1i, 1i, 1i}[1:4], []complex64{1, 1, 2, 1, 1}[1:4], []complex64{2, 3, 2}},
		// Run big test twice, once aligned once unaligned.
		{1 - 1i,
			make([]complex64, 10),
			[]complex64{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
			[]complex64{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
			[]complex64{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
		{1 - 1i,
			make([]complex64, 10),
			[]complex64{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
			[]complex64{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
			[]complex64{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	} {
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
