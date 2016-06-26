// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128

import "testing"

var tests = []struct {
	incX, incY, incDst int
	ix, iy, idst       uintptr
	a                  complex128
	dst, x, y          []complex128
	ex                 []complex128
}{
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   1 + 1i,
		dst: []complex128{5},
		x:   []complex128{1},
		y:   []complex128{1i},
		ex:  []complex128{1 + 2i}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   1 + 2i,
		dst: []complex128{0, 0, 0},
		x:   []complex128{0, 0, 0},
		y:   []complex128{1, 1, 1},
		ex:  []complex128{1, 1, 1}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   1 + 2i,
		dst: []complex128{0, 0, 0},
		x:   []complex128{0, 0},
		y:   []complex128{1, 1, 1},
		ex:  []complex128{1, 1}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   1 + 2i,
		dst: []complex128{1i, 1i, 1i},
		x:   []complex128{1i, 1i, 1i},
		y:   []complex128{1, 2, 1},
		ex:  []complex128{-1 + 1i, 1i, -1 + 1i}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   -1i,
		dst: []complex128{1i, 1i, 1i},
		x:   []complex128{1i, 1i, 1i},
		y:   []complex128{1, 2, 1},
		ex:  []complex128{2, 3, 2}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   -1i,
		dst: []complex128{1i, 1i, 1i},
		x:   []complex128{1i, 1i, 1i, 1i, 1i}[1:4],
		y:   []complex128{1, 1, 2, 1, 1}[1:4],
		ex:  []complex128{2, 3, 2}},
	{incX: 2, incY: 4, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   -2,
		dst: []complex128{1i, 1i, 1i, 1i, 1i},
		x:   []complex128{2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i},
		y:   []complex128{1, 1, 2, 1, 1},
		ex:  []complex128{-3 - 2i, -3 - 2i, -2 - 2i, -3 - 2i, -3 - 2i}},
	// Run big test twice, once aligned once unaligned.
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   1 - 1i,
		dst: make([]complex128, 10),
		x:   []complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		y:   []complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		ex:  []complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	{incX: 2, incY: 2, incDst: 3, ix: 0, iy: 0, idst: 0,
		a:   1 - 1i,
		dst: make([]complex128, 10),
		x:   []complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		y:   []complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		ex:  []complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	{incX: -2, incY: -2, incDst: -3, ix: 18, iy: 18, idst: 27,
		a:   1 - 1i,
		dst: make([]complex128, 10),
		x:   []complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		y:   []complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		ex:  []complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
	{incX: -2, incY: 2, incDst: -3, ix: 18, iy: 0, idst: 27,
		a:   1 - 1i,
		dst: make([]complex128, 10),
		x:   []complex128{1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i, 1i},
		y:   []complex128{1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		ex:  []complex128{2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 2 + 1i, 3 + 1i, 2 + 1i, 2 + 1i}},
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

func isValidGuard(v []complex128, g complex128, g_ln int) bool {
	for i := 0; i < g_ln; i++ {
		if v[i] != g || v[len(v)-1-i] != g {
			return false
		}
	}
	return true
}

func TestAxpyUnitary(t *testing.T) {
	var x_gd, y_gd complex128 = 1, 1
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardVector(v.x, x_gd, xg_ln), guardVector(v.y, y_gd, yg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		AxpyUnitary(v.a, x, y)
		for i := range v.ex {
			if y[i] != v.ex[i] {
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
	var x_gd, y_gd, dst_gd complex128 = 1, 1, 0
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardVector(v.x, x_gd, xg_ln), guardVector(v.y, y_gd, yg_ln)
		v.dst = guardVector(v.dst, dst_gd, xg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		dst := v.dst[xg_ln : len(v.dst)-xg_ln]
		AxpyUnitaryTo(dst, v.a, x, y)
		for i := range v.ex {
			if dst[i] != v.ex[i] {
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

func checkValidIncGuard(t *testing.T, v []complex128, g complex128, incV uintptr, g_ln int) {
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
	var x_gd, y_gd complex128 = 1, 1
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardIncVector(v.x, x_gd, uintptr(v.incX), xg_ln), guardIncVector(v.y, y_gd, uintptr(v.incY), yg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		AxpyInc(v.a, x, y, uintptr(len(v.ex)), uintptr(v.incX), uintptr(v.incY), v.ix, v.iy)
		for i := range v.ex {
			if y[int(v.iy)+i*int(v.incY)] != v.ex[i] {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", j, i, y[i*int(v.incY)], v.ex[i])
			}
		}
		checkValidIncGuard(t, v.x, x_gd, uintptr(v.incX), xg_ln)
		checkValidIncGuard(t, v.y, y_gd, uintptr(v.incY), yg_ln)
	}
}

func TestAxpyIncTo(t *testing.T) {
	var x_gd, y_gd, dst_gd complex128 = 1, 1, 0
	for j, v := range tests {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.x, v.y = guardIncVector(v.x, x_gd, uintptr(v.incX), xg_ln), guardIncVector(v.y, y_gd, uintptr(v.incY), yg_ln)
		v.dst = guardIncVector(v.dst, dst_gd, uintptr(v.incDst), xg_ln)
		x, y := v.x[xg_ln:len(v.x)-xg_ln], v.y[yg_ln:len(v.y)-yg_ln]
		dst := v.dst[xg_ln : len(v.dst)-xg_ln]
		AxpyIncTo(dst, uintptr(v.incDst), v.idst, v.a, x, y, uintptr(len(v.ex)), uintptr(v.incX), uintptr(v.incY), v.ix, v.iy)
		for i := range v.ex {
			if dst[int(v.idst)+i*int(v.incDst)] != v.ex[i] {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", j, i, dst[i*int(v.incDst)], v.ex[i])
			}
		}
		checkValidIncGuard(t, v.x, x_gd, uintptr(v.incX), xg_ln)
		checkValidIncGuard(t, v.y, y_gd, uintptr(v.incY), yg_ln)
		checkValidIncGuard(t, v.dst, dst_gd, uintptr(v.incDst), xg_ln)
	}
}
