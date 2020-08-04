// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128_test

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/cmplxs/cscalar"
	. "gonum.org/v1/gonum/internal/asm/c128"
)

func TestAdd(t *testing.T) {
	var src_gd, dst_gd complex128 = 1, 0
	for j, v := range []struct {
		dst, src, expect []complex128
	}{
		{
			dst:    []complex128{1 + 1i},
			src:    []complex128{0},
			expect: []complex128{1 + 1i},
		},
		{
			dst:    []complex128{1, 2, 3},
			src:    []complex128{1 + 1i},
			expect: []complex128{2 + 1i, 2, 3},
		},
		{
			dst:    []complex128{},
			src:    []complex128{},
			expect: []complex128{},
		},
		{
			dst:    []complex128{1},
			src:    []complex128{cnan},
			expect: []complex128{cnan},
		},
		{
			dst:    []complex128{8, 8, 8, 8, 8},
			src:    []complex128{2 + 1i, 4 - 1i, cnan, 8 + 1i, 9 - 1i},
			expect: []complex128{10 + 1i, 12 - 1i, cnan, 16 + 1i, 17 - 1i},
		},
		{
			dst:    []complex128{0, 1 + 1i, 2, 3 - 1i, 4},
			src:    []complex128{cinf, 4, cnan, 8 + 1i, 9 - 1i},
			expect: []complex128{cinf, 5 + 1i, cnan, 11, 13 - 1i},
		},
		{
			dst:    make([]complex128, 50)[1:49],
			src:    make([]complex128, 50)[1:49],
			expect: make([]complex128, 50)[1:49],
		},
	} {
		sg_ln, dg_ln := 4+j%2, 4+j%3
		v.src, v.dst = guardVector(v.src, src_gd, sg_ln), guardVector(v.dst, dst_gd, dg_ln)
		src, dst := v.src[sg_ln:len(v.src)-sg_ln], v.dst[dg_ln:len(v.dst)-dg_ln]
		Add(dst, src)
		for i := range v.expect {
			if !cscalar.Same(dst[i], v.expect[i]) {
				t.Errorf("Test %d Add error at %d Got: %v Expected: %v", j, i, dst[i], v.expect[i])
			}
		}
		if !isValidGuard(v.src, src_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:sg_ln], v.src[len(v.src)-sg_ln:])
		}
		if !isValidGuard(v.dst, dst_gd, dg_ln) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", j, v.dst[:dg_ln], v.dst[len(v.dst)-dg_ln:])
		}
	}
}

func TestAddConst(t *testing.T) {
	var src_gd complex128 = 0
	for j, v := range []struct {
		alpha       complex128
		src, expect []complex128
	}{
		{
			alpha:  1 + 1i,
			src:    []complex128{0},
			expect: []complex128{1 + 1i},
		},
		{
			alpha:  5,
			src:    []complex128{},
			expect: []complex128{},
		},
		{
			alpha:  1,
			src:    []complex128{cnan},
			expect: []complex128{cnan},
		},
		{
			alpha:  8 + 1i,
			src:    []complex128{2, 4, cnan, 8, 9},
			expect: []complex128{10 + 1i, 12 + 1i, cnan, 16 + 1i, 17 + 1i},
		},
		{
			alpha:  cinf,
			src:    []complex128{cinf, 4, cnan, 8, 9},
			expect: []complex128{cinf, cinf, cnan, cinf, cinf},
		},
	} {
		g_ln := 4 + j%2
		v.src = guardVector(v.src, src_gd, g_ln)
		src := v.src[g_ln : len(v.src)-g_ln]
		AddConst(v.alpha, src)
		for i := range v.expect {
			if !cscalar.Same(src[i], v.expect[i]) {
				t.Errorf("Test %d AddConst error at %d Got: %v Expected: %v", j, i, src[i], v.expect[i])
			}
		}
		if !isValidGuard(v.src, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:g_ln], v.src[len(v.src)-g_ln:])
		}
	}
}

var axpyTests = []struct {
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

func TestAxpyUnitary(t *testing.T) {
	const xGdVal, yGdVal = 1, 1
	for cas, test := range axpyTests {
		xgLn, ygLn := 4+cas%2, 4+cas%3
		test.x, test.y = guardVector(test.x, xGdVal, xgLn), guardVector(test.y, yGdVal, ygLn)
		x, y := test.x[xgLn:len(test.x)-xgLn], test.y[ygLn:len(test.y)-ygLn]
		AxpyUnitary(test.a, x, y)
		for i := range test.ex {
			if y[i] != test.ex[i] {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", cas, i, y[i], test.ex[i])
			}
		}
		if !isValidGuard(test.x, xGdVal, xgLn) {
			t.Errorf("Test %d Guard violated in x vector %v %v", cas, test.x[:xgLn], test.x[len(test.x)-xgLn:])
		}
		if !isValidGuard(test.y, yGdVal, ygLn) {
			t.Errorf("Test %d Guard violated in y vector %v %v", cas, test.y[:ygLn], test.y[len(test.y)-ygLn:])
		}
	}
}

func TestAxpyUnitaryTo(t *testing.T) {
	const xGdVal, yGdVal, dstGdVal = 1, 1, 0
	for cas, test := range axpyTests {
		xgLn, ygLn := 4+cas%2, 4+cas%3
		test.x, test.y = guardVector(test.x, xGdVal, xgLn), guardVector(test.y, yGdVal, ygLn)
		test.dst = guardVector(test.dst, dstGdVal, xgLn)
		x, y := test.x[xgLn:len(test.x)-xgLn], test.y[ygLn:len(test.y)-ygLn]
		dst := test.dst[xgLn : len(test.dst)-xgLn]
		AxpyUnitaryTo(dst, test.a, x, y)
		for i := range test.ex {
			if dst[i] != test.ex[i] {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", cas, i, dst[i], test.ex[i])
			}
		}
		if !isValidGuard(test.x, xGdVal, xgLn) {
			t.Errorf("Test %d Guard violated in x vector %v %v", cas, test.x[:xgLn], test.x[len(test.x)-xgLn:])
		}
		if !isValidGuard(test.y, yGdVal, ygLn) {
			t.Errorf("Test %d Guard violated in y vector %v %v", cas, test.y[:ygLn], test.y[len(test.y)-ygLn:])
		}
		if !isValidGuard(test.dst, dstGdVal, xgLn) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", cas, test.dst[:xgLn], test.dst[len(test.dst)-xgLn:])
		}

	}
}

func TestAxpyInc(t *testing.T) {
	const xGdVal, yGdVal = 1, 1
	for cas, test := range axpyTests {
		xgLn, ygLn := 4+cas%2, 4+cas%3
		test.x, test.y = guardIncVector(test.x, xGdVal, test.incX, xgLn), guardIncVector(test.y, yGdVal, test.incY, ygLn)
		x, y := test.x[xgLn:len(test.x)-xgLn], test.y[ygLn:len(test.y)-ygLn]
		AxpyInc(test.a, x, y, uintptr(len(test.ex)), uintptr(test.incX), uintptr(test.incY), test.ix, test.iy)
		for i := range test.ex {
			if y[int(test.iy)+i*int(test.incY)] != test.ex[i] {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", cas, i, y[i*int(test.incY)], test.ex[i])
			}
		}
		checkValidIncGuard(t, test.x, xGdVal, test.incX, xgLn)
		checkValidIncGuard(t, test.y, yGdVal, test.incY, ygLn)
	}
}

func TestAxpyIncTo(t *testing.T) {
	const xGdVal, yGdVal, dstGdVal = 1, 1, 0
	for cas, test := range axpyTests {
		xgLn, ygLn := 4+cas%2, 4+cas%3
		test.x, test.y = guardIncVector(test.x, xGdVal, test.incX, xgLn), guardIncVector(test.y, yGdVal, test.incY, ygLn)
		test.dst = guardIncVector(test.dst, dstGdVal, test.incDst, xgLn)
		x, y := test.x[xgLn:len(test.x)-xgLn], test.y[ygLn:len(test.y)-ygLn]
		dst := test.dst[xgLn : len(test.dst)-xgLn]
		AxpyIncTo(dst, uintptr(test.incDst), test.idst, test.a, x, y, uintptr(len(test.ex)), uintptr(test.incX), uintptr(test.incY), test.ix, test.iy)
		for i := range test.ex {
			if dst[int(test.idst)+i*int(test.incDst)] != test.ex[i] {
				t.Errorf("Test %d Unexpected result at %d Got: %v Expected: %v", cas, i, dst[i*int(test.incDst)], test.ex[i])
			}
		}
		checkValidIncGuard(t, test.x, xGdVal, test.incX, xgLn)
		checkValidIncGuard(t, test.y, yGdVal, test.incY, ygLn)
		checkValidIncGuard(t, test.dst, dstGdVal, test.incDst, xgLn)
	}
}

func TestCumSum(t *testing.T) {
	var src_gd, dst_gd complex128 = -1, 0
	for j, v := range []struct {
		dst, src, expect []complex128
	}{
		{
			dst:    []complex128{},
			src:    []complex128{},
			expect: []complex128{},
		},
		{
			dst:    []complex128{0},
			src:    []complex128{1 + 1i},
			expect: []complex128{1 + 1i},
		},
		{
			dst:    []complex128{cnan},
			src:    []complex128{cnan},
			expect: []complex128{cnan},
		},
		{
			dst:    []complex128{0, 0, 0},
			src:    []complex128{1, 2 + 1i, 3 + 2i},
			expect: []complex128{1, 3 + 1i, 6 + 3i},
		},
		{
			dst:    []complex128{0, 0, 0, 0},
			src:    []complex128{1, 2 + 1i, 3 + 2i},
			expect: []complex128{1, 3 + 1i, 6 + 3i},
		},
		{
			dst:    []complex128{0, 0, 0, 0},
			src:    []complex128{1, 2 + 1i, 3 + 2i, 4 + 3i},
			expect: []complex128{1, 3 + 1i, 6 + 3i, 10 + 6i},
		},
		{
			dst:    []complex128{1, cnan, cnan, 1, 1},
			src:    []complex128{1, 1, cnan, 1, 1},
			expect: []complex128{1, 2, cnan, cnan, cnan},
		},
		{
			dst:    []complex128{cnan, 4, cinf, cinf, 9},
			src:    []complex128{cinf, 4, cnan, cinf, 9},
			expect: []complex128{cinf, cinf, cnan, cnan, cnan},
		},
		{
			dst:    make([]complex128, 16),
			src:    []complex128{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expect: []complex128{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	} {
		g_ln := 4 + j%2
		v.src, v.dst = guardVector(v.src, src_gd, g_ln), guardVector(v.dst, dst_gd, g_ln)
		src, dst := v.src[g_ln:len(v.src)-g_ln], v.dst[g_ln:len(v.dst)-g_ln]
		ret := CumSum(dst, src)
		for i := range v.expect {
			if !cscalar.Same(ret[i], v.expect[i]) {
				t.Errorf("Test %d CumSum error at %d Got: %v Expected: %v", j, i, ret[i], v.expect[i])
			}
			if !cscalar.Same(ret[i], dst[i]) {
				t.Errorf("Test %d CumSum ret/dst mismatch %d Ret: %v Dst: %v", j, i, ret[i], dst[i])
			}
		}
		if !isValidGuard(v.src, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:g_ln], v.src[len(v.src)-g_ln:])
		}
		if !isValidGuard(v.dst, dst_gd, g_ln) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", j, v.dst[:g_ln], v.dst[len(v.dst)-g_ln:])
		}
	}
}

func TestCumProd(t *testing.T) {
	var src_gd, dst_gd complex128 = -1, 1
	for j, v := range []struct {
		dst, src, expect []complex128
	}{
		{
			dst:    []complex128{},
			src:    []complex128{},
			expect: []complex128{},
		},
		{
			dst:    []complex128{1},
			src:    []complex128{1 + 1i},
			expect: []complex128{1 + 1i},
		},
		{
			dst:    []complex128{cnan},
			src:    []complex128{cnan},
			expect: []complex128{cnan},
		},
		{
			dst:    []complex128{0, 0, 0, 0},
			src:    []complex128{1, 2 + 1i, 3 + 2i, 4 + 3i},
			expect: []complex128{1, 2 + 1i, 4 + 7i, -5 + 40i},
		},
		{
			dst:    []complex128{0, 0, 0},
			src:    []complex128{1, 2 + 1i, 3 + 2i},
			expect: []complex128{1, 2 + 1i, 4 + 7i},
		},
		{
			dst:    []complex128{0, 0, 0, 0},
			src:    []complex128{1, 2 + 1i, 3 + 2i},
			expect: []complex128{1, 2 + 1i, 4 + 7i},
		},
		{
			dst:    []complex128{cnan, 1, cnan, 1, 0},
			src:    []complex128{1, 1, cnan, 1, 1},
			expect: []complex128{1, 1, cnan, cnan, cnan},
		},
		{
			dst:    []complex128{cnan, 4, cnan, cinf, 9},
			src:    []complex128{cinf, 4, cnan, cinf, 9},
			expect: []complex128{cinf, cnan, cnan, cnan, cnan},
		},
		{
			dst:    make([]complex128, 18),
			src:    []complex128{2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i, 2i},
			expect: []complex128{2i, -4, -8i, 16, 32i, -64, -128i, 256, 512i, -1024, -2048i, 4096, 8192i, -16384, -32768i, 65536},
		},
	} {
		sg_ln, dg_ln := 4+j%2, 4+j%3
		v.src, v.dst = guardVector(v.src, src_gd, sg_ln), guardVector(v.dst, dst_gd, dg_ln)
		src, dst := v.src[sg_ln:len(v.src)-sg_ln], v.dst[dg_ln:len(v.dst)-dg_ln]
		ret := CumProd(dst, src)
		for i := range v.expect {
			if !cscalar.Same(ret[i], v.expect[i]) {
				t.Errorf("Test %d CumProd error at %d Got: %v Expected: %v", j, i, ret[i], v.expect[i])
			}
			if !cscalar.Same(ret[i], dst[i]) {
				t.Errorf("Test %d CumProd ret/dst mismatch %d Ret: %v Dst: %v", j, i, ret[i], dst[i])
			}
		}
		if !isValidGuard(v.src, src_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:sg_ln], v.src[len(v.src)-sg_ln:])
		}
		if !isValidGuard(v.dst, dst_gd, dg_ln) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", j, v.dst[:dg_ln], v.dst[len(v.dst)-dg_ln:])
		}
	}
}

func TestDiv(t *testing.T) {
	const tol = 1e-15

	var src_gd, dst_gd complex128 = -1, 0.5
	for j, v := range []struct {
		dst, src, expect []complex128
	}{
		{
			dst:    []complex128{1 + 1i},
			src:    []complex128{1 + 1i},
			expect: []complex128{1},
		},
		{
			dst:    []complex128{cnan},
			src:    []complex128{cnan},
			expect: []complex128{cnan},
		},
		{
			dst:    []complex128{1 + 1i, 2 + 1i, 3 + 1i, 4 + 1i},
			src:    []complex128{1 + 1i, 2 + 1i, 3 + 1i, 4 + 1i},
			expect: []complex128{1, 1, 1, 1},
		},
		{
			dst:    []complex128{1 + 1i, 2 + 1i, 3 + 1i, 4 + 1i, 2 + 2i, 4 + 2i, 6 + 2i, 8 + 2i},
			src:    []complex128{1 + 1i, 2 + 1i, 3 + 1i, 4 + 1i, 1 + 1i, 2 + 1i, 3 + 1i, 4 + 1i},
			expect: []complex128{1, 1, 1, 1, 2, 2, 2, 2},
		},
		{
			dst:    []complex128{2 + 2i, 4 + 8i, 6 - 12i},
			src:    []complex128{1 + 1i, 2 + 4i, 3 - 6i},
			expect: []complex128{2, 2, 2},
		},
		{
			dst:    []complex128{0, 0, 0, 0},
			src:    []complex128{1 + 1i, 2 + 2i, 3 + 3i},
			expect: []complex128{0, 0, 0},
		},
		{
			dst:    []complex128{cnan, 1, cnan, 1, 0, cnan, 1, cnan, 1, 0},
			src:    []complex128{1, 1, cnan, 1, 1, 1, 1, cnan, 1, 1},
			expect: []complex128{cnan, 1, cnan, 1, 0, cnan, 1, cnan, 1, 0},
		},
		{
			dst:    []complex128{cinf, 4, cnan, cinf, 9, cinf, 4, cnan, cinf, 9},
			src:    []complex128{cinf, 4, cnan, cinf, 3, cinf, 4, cnan, cinf, 3},
			expect: []complex128{cnan, 1, cnan, cnan, 3, cnan, 1, cnan, cnan, 3},
		},
	} {
		sg_ln, dg_ln := 4+j%2, 4+j%3
		v.src, v.dst = guardVector(v.src, src_gd, sg_ln), guardVector(v.dst, dst_gd, dg_ln)
		src, dst := v.src[sg_ln:len(v.src)-sg_ln], v.dst[dg_ln:len(v.dst)-dg_ln]
		Div(dst, src)
		for i := range v.expect {
			if !sameCmplxApprox(dst[i], v.expect[i], tol) {
				t.Errorf("Test %d Div error at %d Got: %v Expected: %v", j, i, dst[i], v.expect[i])
			}
		}
		if !isValidGuard(v.src, src_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:sg_ln], v.src[len(v.src)-sg_ln:])
		}
		if !isValidGuard(v.dst, dst_gd, dg_ln) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", j, v.dst[:dg_ln], v.dst[len(v.dst)-dg_ln:])
		}
	}
}

func TestDivTo(t *testing.T) {
	const tol = 1e-15

	var dst_gd, x_gd, y_gd complex128 = -1, 0.5, 0.25
	for j, v := range []struct {
		dst, x, y, expect []complex128
	}{
		{
			dst:    []complex128{1 - 1i},
			x:      []complex128{1 + 1i},
			y:      []complex128{1 + 1i},
			expect: []complex128{1},
		},
		{
			dst:    []complex128{1},
			x:      []complex128{cnan},
			y:      []complex128{cnan},
			expect: []complex128{cnan},
		},
		{
			dst:    []complex128{-2, -2, -2},
			x:      []complex128{1 + 1i, 2 + 1i, 3 + 1i},
			y:      []complex128{1 + 1i, 2 + 1i, 3 + 1i},
			expect: []complex128{1, 1, 1},
		},
		{
			dst:    []complex128{0, 0, 0},
			x:      []complex128{2 + 2i, 4 + 4i, 6 + 6i},
			y:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			expect: []complex128{2, 2, 2},
		},
		{
			dst:    []complex128{-1, -1, -1},
			x:      []complex128{0, 0, 0},
			y:      []complex128{1 + 1i, 2 + 1i, 3 + 1i},
			expect: []complex128{0, 0, 0},
		},
		{
			dst:    []complex128{cinf, cinf, cinf, cinf, cinf, cinf, cinf, cinf, cinf, cinf},
			x:      []complex128{cnan, 1, cnan, 1, 0, cnan, 1, cnan, 1, 0},
			y:      []complex128{1, 1, cnan, 1, 1, 1, 1, cnan, 1, 1},
			expect: []complex128{cnan, 1, cnan, 1, 0, cnan, 1, cnan, 1, 0},
		},
		{
			dst:    []complex128{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:      []complex128{cinf, 4, cnan, cinf, 9, cinf, 4, cnan, cinf, 9},
			y:      []complex128{cinf, 4, cnan, cinf, 3, cinf, 4, cnan, cinf, 3},
			expect: []complex128{cnan, 1, cnan, cnan, 3, cnan, 1, cnan, cnan, 3},
		},
	} {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.y, v.x = guardVector(v.y, y_gd, yg_ln), guardVector(v.x, x_gd, xg_ln)
		y, x := v.y[yg_ln:len(v.y)-yg_ln], v.x[xg_ln:len(v.x)-xg_ln]
		v.dst = guardVector(v.dst, dst_gd, xg_ln)
		dst := v.dst[xg_ln : len(v.dst)-xg_ln]
		ret := DivTo(dst, x, y)
		for i := range v.expect {
			if !sameCmplxApprox(ret[i], v.expect[i], tol) {
				t.Errorf("Test %d DivTo error at %d Got: %v Expected: %v", j, i, ret[i], v.expect[i])
			}
			if !cscalar.Same(ret[i], dst[i]) {
				t.Errorf("Test %d DivTo ret/dst mismatch %d Ret: %v Dst: %v", j, i, ret[i], dst[i])
			}
		}
		if !isValidGuard(v.y, y_gd, yg_ln) {
			t.Errorf("Test %d Guard violated in y vector %v %v", j, v.y[:yg_ln], v.y[len(v.y)-yg_ln:])
		}
		if !isValidGuard(v.x, x_gd, xg_ln) {
			t.Errorf("Test %d Guard violated in x vector %v %v", j, v.x[:xg_ln], v.x[len(v.x)-xg_ln:])
		}
		if !isValidGuard(v.dst, dst_gd, xg_ln) {
			t.Errorf("Test %d Guard violated in dst vector %v %v", j, v.dst[:xg_ln], v.dst[len(v.dst)-xg_ln:])
		}
	}
}

var dscalTests = []struct {
	alpha float64
	x     []complex128
	want  []complex128
}{
	{
		alpha: 0,
		x:     []complex128{},
		want:  []complex128{},
	},
	{
		alpha: 1,
		x:     []complex128{1 + 2i},
		want:  []complex128{1 + 2i},
	},
	{
		alpha: 2,
		x:     []complex128{1 + 2i},
		want:  []complex128{2 + 4i},
	},
	{
		alpha: 2,
		x:     []complex128{1 + 2i, 3 + 5i, 6 + 11i, 12 - 23i},
		want:  []complex128{2 + 4i, 6 + 10i, 12 + 22i, 24 - 46i},
	},
	{
		alpha: 3,
		x:     []complex128{1 + 2i, 5 + 4i, 3 + 6i, 8 + 12i, -3 - 2i, -5 + 5i},
		want:  []complex128{3 + 6i, 15 + 12i, 9 + 18i, 24 + 36i, -9 - 6i, -15 + 15i},
	},
	{
		alpha: 5,
		x:     []complex128{1 + 2i, 5 + 4i, 3 + 6i, 8 + 12i, -3 - 2i, -5 + 5i, 1 + 2i, 5 + 4i, 3 + 6i, 8 + 12i, -3 - 2i, -5 + 5i},
		want:  []complex128{5 + 10i, 25 + 20i, 15 + 30i, 40 + 60i, -15 - 10i, -25 + 25i, 5 + 10i, 25 + 20i, 15 + 30i, 40 + 60i, -15 - 10i, -25 + 25i},
	},
}

func TestDscalUnitary(t *testing.T) {
	const xGdVal = -0.5
	for i, test := range dscalTests {
		for _, align := range align1 {
			prefix := fmt.Sprintf("Test %v (x:%v)", i, align)
			xgLn := 4 + align
			xg := guardVector(test.x, xGdVal, xgLn)
			x := xg[xgLn : len(xg)-xgLn]

			DscalUnitary(test.alpha, x)

			for i := range test.want {
				if !cscalar.Same(x[i], test.want[i]) {
					t.Errorf(msgVal, prefix, i, x[i], test.want[i])
				}
			}
			if !isValidGuard(xg, xGdVal, xgLn) {
				t.Errorf(msgGuard, prefix, "x", xg[:xgLn], xg[len(xg)-xgLn:])
			}
		}
	}
}

func TestDscalInc(t *testing.T) {
	const xGdVal = -0.5
	gdLn := 4
	for i, test := range dscalTests {
		n := len(test.x)
		for _, incX := range []int{1, 2, 3, 4, 7, 10} {
			prefix := fmt.Sprintf("Test %v (x:%v)", i, incX)
			xg := guardIncVector(test.x, xGdVal, incX, gdLn)
			x := xg[gdLn : len(xg)-gdLn]

			DscalInc(test.alpha, x, uintptr(n), uintptr(incX))

			for i := range test.want {
				if !cscalar.Same(x[i*incX], test.want[i]) {
					t.Errorf(msgVal, prefix, i, x[i*incX], test.want[i])
				}
			}
			checkValidIncGuard(t, xg, xGdVal, incX, gdLn)
		}
	}
}

var scalTests = []struct {
	alpha complex128
	x     []complex128
	want  []complex128
}{
	{
		alpha: 0,
		x:     []complex128{},
		want:  []complex128{},
	},
	{
		alpha: 1 + 1i,
		x:     []complex128{1 + 2i},
		want:  []complex128{-1 + 3i},
	},
	{
		alpha: 2 + 3i,
		x:     []complex128{1 + 2i},
		want:  []complex128{-4 + 7i},
	},
	{
		alpha: 2 - 4i,
		x:     []complex128{1 + 2i},
		want:  []complex128{10},
	},
	{
		alpha: 2 + 8i,
		x:     []complex128{1 + 2i, 5 + 4i, 3 + 6i, 8 + 12i, -3 - 2i, -5 + 5i},
		want:  []complex128{-14 + 12i, -22 + 48i, -42 + 36i, -80 + 88i, 10 - 28i, -50 - 30i},
	},
	{
		alpha: 5 - 10i,
		x:     []complex128{1 + 2i, 5 + 4i, 3 + 6i, 8 + 12i, -3 - 2i, -5 + 5i, 1 + 2i, 5 + 4i, 3 + 6i, 8 + 12i, -3 - 2i, -5 + 5i},
		want:  []complex128{25, 65 - 30i, 75, 160 - 20i, -35 + 20i, 25 + 75i, 25, 65 - 30i, 75, 160 - 20i, -35 + 20i, 25 + 75i},
	},
}

func TestScalUnitary(t *testing.T) {
	const xGdVal = -0.5
	for i, test := range scalTests {
		for _, align := range align1 {
			prefix := fmt.Sprintf("Test %v (x:%v)", i, align)
			xgLn := 4 + align
			xg := guardVector(test.x, xGdVal, xgLn)
			x := xg[xgLn : len(xg)-xgLn]

			ScalUnitary(test.alpha, x)

			for i := range test.want {
				if !cscalar.Same(x[i], test.want[i]) {
					t.Errorf(msgVal, prefix, i, x[i], test.want[i])
				}
			}
			if !isValidGuard(xg, xGdVal, xgLn) {
				t.Errorf(msgGuard, prefix, "x", xg[:xgLn], xg[len(xg)-xgLn:])
			}
		}
	}
}

func TestScalInc(t *testing.T) {
	const xGdVal = -0.5
	gdLn := 4
	for i, test := range scalTests {
		n := len(test.x)
		for _, inc := range []int{1, 2, 3, 4, 7, 10} {
			prefix := fmt.Sprintf("Test %v (x:%v)", i, inc)
			xg := guardIncVector(test.x, xGdVal, inc, gdLn)
			x := xg[gdLn : len(xg)-gdLn]

			ScalInc(test.alpha, x, uintptr(n), uintptr(inc))

			for i := range test.want {
				if !cscalar.Same(x[i*inc], test.want[i]) {
					t.Errorf(msgVal, prefix, i, x[i*inc], test.want[i])
				}
			}
			checkValidIncGuard(t, xg, xGdVal, inc, gdLn)
		}
	}
}

func TestSum(t *testing.T) {
	var srcGd complex128 = -1
	for j, v := range []struct {
		src    []complex128
		expect complex128
	}{
		{
			src:    []complex128{},
			expect: 0,
		},
		{
			src:    []complex128{1},
			expect: 1,
		},
		{
			src:    []complex128{cnan},
			expect: cnan,
		},
		{
			src:    []complex128{1 + 1i, 2 + 2i, 3 + 3i},
			expect: 6 + 6i,
		},
		{
			src:    []complex128{1 + 1i, -4, 3 - 1i},
			expect: 0,
		},
		{
			src:    []complex128{1 - 1i, 2 + 2i, 3 - 3i, 4 + 4i},
			expect: 10 + 2i,
		},
		{
			src:    []complex128{1, 1, cnan, 1, 1},
			expect: cnan,
		},
		{
			src:    []complex128{cinf, 4, cnan, cinf, 9},
			expect: cnan,
		},
		{
			src:    []complex128{1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 9 + 9i, 1 + 1i, 1 + 1i, 1 + 1i, 2 + 2i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 5 + 5i, 1 + 1i},
			expect: 29 + 29i,
		},
		{
			src:    []complex128{1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 9 + 9i, 1 + 1i, 1 + 1i, 1 + 1i, 2 + 2i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 5 + 5i, 11 + 11i, 1 + 1i, 1 + 1i, 1 + 1i, 9 + 9i, 1 + 1i, 1 + 1i, 1 + 1i, 2 + 2i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 1 + 1i, 5 + 5i, 1 + 1i},
			expect: 67 + 67i,
		},
	} {
		gdLn := 4 + j%2
		gsrc := guardVector(v.src, srcGd, gdLn)
		src := gsrc[gdLn : len(gsrc)-gdLn]
		ret := Sum(src)
		if !cscalar.Same(ret, v.expect) {
			t.Errorf("Test %d Sum error Got: %v Expected: %v", j, ret, v.expect)
		}
		if !isValidGuard(gsrc, srcGd, gdLn) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, gsrc[:gdLn], gsrc[len(gsrc)-gdLn:])
		}
	}
}
