// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"math"
	"testing"
)

var (
	nan = math.NaN()
	inf = math.Inf(1)
)

func guardVector(v []float64, g float64, g_ln int) (guarded []float64) {
	guarded = make([]float64, len(v)+g_ln*2)
	copy(guarded[g_ln:], v)
	for i := 0; i < g_ln; i++ {
		guarded[i] = g
		guarded[len(guarded)-1-i] = g
	}
	return guarded
}

func isValidGuard(v []float64, g float64, g_ln int) bool {
	for i := 0; i < g_ln; i++ {
		if !same(v[i], g) || v[len(v)-1-i] != g {
			return false
		}
	}
	return true
}

func guardIncVector(vec []float64, guard_val float64, inc, guard_len int) (guarded []float64) {
	s_ln := len(vec) * inc
	guarded = make([]float64, s_ln+guard_len*2)
	for i, j := 0, 0; i < len(guarded); i++ {
		switch {
		case i < guard_len, guard_len+s_ln < i:
			guarded[i] = guard_val
		case (i-guard_len)%(inc) == 0 && j < len(vec):
			guarded[i] = vec[j]
			j++
		default:
			guarded[i] = guard_val
		}
	}
	return guarded
}

func checkValidIncGuard(t *testing.T, vec []float64, guard_val float64, inc, guard_len int) {
	s_ln := len(vec) - 2*guard_len
	for i := range vec {
		switch {
		case same(vec[i], guard_val):
			// Correct value
		case i < guard_len:
			t.Errorf("Front guard violated at %d %v", i, vec[:guard_len])
		case i > guard_len+s_ln:
			t.Errorf("Back guard violated at %d %v", i-guard_len-s_ln, vec[guard_len+s_ln:])
		case (i-guard_len)%inc == 0 && (i-guard_len)/inc < len(vec):
			// Ignore input values
		default:
			t.Errorf("Internal guard violated at %d %v", i-guard_len, vec[guard_len:guard_len+s_ln])
		}
	}
}

func same(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

func TestAbsSum(t *testing.T) {
	var src_gd float64 = 1
	for j, v := range []struct {
		ex  float64
		src []float64
	}{
		{
			ex:  0,
			src: []float64{},
		},
		{
			ex:  2,
			src: []float64{2},
		},
		{
			ex:  6,
			src: []float64{1, 2, 3},
		},
		{
			ex:  6,
			src: []float64{-1, -2, -3},
		},
		{
			ex:  nan,
			src: []float64{nan},
		},
		{
			ex:  40,
			src: []float64{8, -8, 8, -8, 8},
		},
		{
			ex:  5,
			src: []float64{0, 1, 0, -1, 0, 1, 0, -1, 0, 1},
		},
	} {
		g_ln := 4 + j%2
		v.src = guardVector(v.src, src_gd, g_ln)
		src := v.src[g_ln : len(v.src)-g_ln]
		ret := AbsSum(src)
		if !same(ret, v.ex) {
			t.Errorf("Test %d AbsSum error Got: %f Expected: %f", j, ret, v.ex)
		}
		if !isValidGuard(v.src, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:g_ln], v.src[len(v.src)-g_ln:])
		}
	}
}

func TestAbsSumInc(t *testing.T) {
	var src_gd float64 = 1
	for j, v := range []struct {
		inc int
		ex  float64
		src []float64
	}{
		{
			inc: 2,
			ex:  0,
			src: []float64{},
		},
		{
			inc: 3,
			ex:  2,
			src: []float64{2},
		},
		{
			inc: 10,
			ex:  6,
			src: []float64{1, 2, 3},
		},
		{
			inc: 5,
			ex:  6,
			src: []float64{-1, -2, -3},
		},
		{
			inc: 3,
			ex:  nan,
			src: []float64{nan},
		},
		{
			inc: 15,
			ex:  40,
			src: []float64{8, -8, 8, -8, 8},
		},
		{
			inc: 1,
			ex:  5,
			src: []float64{0, 1, 0, -1, 0, 1, 0, -1, 0, 1},
		},
	} {
		g_ln, ln := 4+j%2, len(v.src)
		v.src = guardIncVector(v.src, src_gd, v.inc, g_ln)
		src := v.src[g_ln : len(v.src)-g_ln]
		ret := AbsSumInc(src, ln, v.inc)
		if !same(ret, v.ex) {
			t.Errorf("Test %d AbsSumInc error Got: %f Expected: %f", j, ret, v.ex)
		}
		checkValidIncGuard(t, v.src, src_gd, v.inc, g_ln)
	}
}

func TestAdd(t *testing.T) {
	var src_gd, dst_gd float64 = 1, 0
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{
			dst:    []float64{1},
			src:    []float64{0},
			expect: []float64{1},
		},
		{
			dst:    []float64{1, 2, 3},
			src:    []float64{1},
			expect: []float64{2, 2, 3},
		},
		{
			dst:    []float64{},
			src:    []float64{},
			expect: []float64{},
		},
		{
			dst:    []float64{1},
			src:    []float64{nan},
			expect: []float64{nan},
		},
		{
			dst:    []float64{8, 8, 8, 8, 8},
			src:    []float64{2, 4, nan, 8, 9},
			expect: []float64{10, 12, nan, 16, 17},
		},
		{
			dst:    []float64{0, 1, 2, 3, 4},
			src:    []float64{-inf, 4, nan, 8, 9},
			expect: []float64{-inf, 5, nan, 11, 13},
		},
		{
			dst:    make([]float64, 50)[1:49],
			src:    make([]float64, 50)[1:49],
			expect: make([]float64, 50)[1:49],
		},
	} {
		sg_ln, dg_ln := 4+j%2, 4+j%3
		v.src, v.dst = guardVector(v.src, src_gd, sg_ln), guardVector(v.dst, dst_gd, dg_ln)
		src, dst := v.src[sg_ln:len(v.src)-sg_ln], v.dst[dg_ln:len(v.dst)-dg_ln]
		Add(dst, src)
		for i := range v.expect {
			if !same(dst[i], v.expect[i]) {
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
	var src_gd float64 = 0
	for j, v := range []struct {
		alpha       float64
		src, expect []float64
	}{
		{
			alpha:  1,
			src:    []float64{0},
			expect: []float64{1},
		},
		{
			alpha:  5,
			src:    []float64{},
			expect: []float64{},
		},
		{
			alpha:  1,
			src:    []float64{nan},
			expect: []float64{nan},
		},
		{
			alpha:  8,
			src:    []float64{2, 4, nan, 8, 9},
			expect: []float64{10, 12, nan, 16, 17},
		},
		{
			alpha:  inf,
			src:    []float64{-inf, 4, nan, 8, 9},
			expect: []float64{nan, inf, nan, inf, inf},
		},
	} {
		g_ln := 4 + j%2
		v.src = guardVector(v.src, src_gd, g_ln)
		src := v.src[g_ln : len(v.src)-g_ln]
		AddConst(v.alpha, src)
		for i := range v.expect {
			if !same(src[i], v.expect[i]) {
				t.Errorf("Test %d AddConst error at %d Got: %v Expected: %v", j, i, src[i], v.expect[i])
			}
		}
		if !isValidGuard(v.src, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.src[:g_ln], v.src[len(v.src)-g_ln:])
		}
	}
}

func TestCumSum(t *testing.T) {
	var src_gd, dst_gd float64 = -1, 0
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{
			dst:    []float64{},
			src:    []float64{},
			expect: []float64{},
		},
		{
			dst:    []float64{0},
			src:    []float64{1},
			expect: []float64{1},
		},
		{
			dst:    []float64{nan},
			src:    []float64{nan},
			expect: []float64{nan},
		},
		{
			dst:    []float64{0, 0, 0},
			src:    []float64{1, 2, 3},
			expect: []float64{1, 3, 6},
		},
		{
			dst:    []float64{0, 0, 0, 0},
			src:    []float64{1, 2, 3},
			expect: []float64{1, 3, 6},
		},
		{
			dst:    []float64{0, 0, 0, 0},
			src:    []float64{1, 2, 3, 4},
			expect: []float64{1, 3, 6, 10},
		},
		{
			dst:    []float64{1, nan, nan, 1, 1},
			src:    []float64{1, 1, nan, 1, 1},
			expect: []float64{1, 2, nan, nan, nan},
		},
		{
			dst:    []float64{nan, 4, inf, -inf, 9},
			src:    []float64{inf, 4, nan, -inf, 9},
			expect: []float64{inf, inf, nan, nan, nan},
		},
		{
			dst:    make([]float64, 16),
			src:    []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expect: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	} {
		g_ln := 4 + j%2
		v.src, v.dst = guardVector(v.src, src_gd, g_ln), guardVector(v.dst, dst_gd, g_ln)
		src, dst := v.src[g_ln:len(v.src)-g_ln], v.dst[g_ln:len(v.dst)-g_ln]
		ret := CumSum(dst, src)
		for i := range v.expect {
			if !same(ret[i], v.expect[i]) {
				t.Errorf("Test %d CumSum error at %d Got: %v Expected: %v", j, i, ret[i], v.expect[i])
			}
			if !same(ret[i], dst[i]) {
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
	var src_gd, dst_gd float64 = -1, 1
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{
			dst:    []float64{},
			src:    []float64{},
			expect: []float64{},
		},
		{
			dst:    []float64{1},
			src:    []float64{1},
			expect: []float64{1},
		},
		{
			dst:    []float64{nan},
			src:    []float64{nan},
			expect: []float64{nan},
		},
		{
			dst:    []float64{0, 0, 0, 0},
			src:    []float64{1, 2, 3, 4},
			expect: []float64{1, 2, 6, 24},
		},
		{
			dst:    []float64{0, 0, 0},
			src:    []float64{1, 2, 3},
			expect: []float64{1, 2, 6},
		},
		{
			dst:    []float64{0, 0, 0, 0},
			src:    []float64{1, 2, 3},
			expect: []float64{1, 2, 6},
		},
		{
			dst:    []float64{nan, 1, nan, 1, 0},
			src:    []float64{1, 1, nan, 1, 1},
			expect: []float64{1, 1, nan, nan, nan},
		},
		{
			dst:    []float64{nan, 4, nan, -inf, 9},
			src:    []float64{inf, 4, nan, -inf, 9},
			expect: []float64{inf, inf, nan, nan, nan},
		},
		{
			dst:    make([]float64, 18),
			src:    []float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
			expect: []float64{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536},
		},
	} {
		sg_ln, dg_ln := 4+j%2, 4+j%3
		v.src, v.dst = guardVector(v.src, src_gd, sg_ln), guardVector(v.dst, dst_gd, dg_ln)
		src, dst := v.src[sg_ln:len(v.src)-sg_ln], v.dst[dg_ln:len(v.dst)-dg_ln]
		ret := CumProd(dst, src)
		for i := range v.expect {
			if !same(ret[i], v.expect[i]) {
				t.Errorf("Test %d CumProd error at %d Got: %v Expected: %v", j, i, ret[i], v.expect[i])
			}
			if !same(ret[i], dst[i]) {
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
	var src_gd, dst_gd float64 = -1, 0.5
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{
			dst:    []float64{1},
			src:    []float64{1},
			expect: []float64{1},
		},
		{
			dst:    []float64{nan},
			src:    []float64{nan},
			expect: []float64{nan},
		},
		{
			dst:    []float64{1, 2, 3, 4},
			src:    []float64{1, 2, 3, 4},
			expect: []float64{1, 1, 1, 1},
		},
		{
			dst:    []float64{1, 2, 3, 4, 2, 4, 6, 8},
			src:    []float64{1, 2, 3, 4, 1, 2, 3, 4},
			expect: []float64{1, 1, 1, 1, 2, 2, 2, 2},
		},
		{
			dst:    []float64{2, 4, 6},
			src:    []float64{1, 2, 3},
			expect: []float64{2, 2, 2},
		},
		{
			dst:    []float64{0, 0, 0, 0},
			src:    []float64{1, 2, 3},
			expect: []float64{0, 0, 0},
		},
		{
			dst:    []float64{nan, 1, nan, 1, 0},
			src:    []float64{1, 1, nan, 1, 1},
			expect: []float64{nan, 1, nan, 1, 0},
		},
		{
			dst:    []float64{inf, 4, nan, -inf, 9},
			src:    []float64{inf, 4, nan, -inf, 3},
			expect: []float64{nan, 1, nan, nan, 3},
		},
	} {
		sg_ln, dg_ln := 4+j%2, 4+j%3
		v.src, v.dst = guardVector(v.src, src_gd, sg_ln), guardVector(v.dst, dst_gd, dg_ln)
		src, dst := v.src[sg_ln:len(v.src)-sg_ln], v.dst[dg_ln:len(v.dst)-dg_ln]
		Div(dst, src)
		for i := range v.expect {
			if !same(dst[i], v.expect[i]) {
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
	var dst_gd, x_gd, y_gd float64 = -1, 0.5, 0.25
	for j, v := range []struct {
		dst, x, y, expect []float64
	}{
		{
			dst:    []float64{1},
			x:      []float64{1},
			y:      []float64{1},
			expect: []float64{1},
		},
		{
			dst:    []float64{1},
			x:      []float64{nan},
			y:      []float64{nan},
			expect: []float64{nan},
		},
		{
			dst:    []float64{-2, -2, -2},
			x:      []float64{1, 2, 3},
			y:      []float64{1, 2, 3},
			expect: []float64{1, 1, 1},
		},
		{
			dst:    []float64{0, 0, 0},
			x:      []float64{2, 4, 6},
			y:      []float64{1, 2, 3, 4},
			expect: []float64{2, 2, 2},
		},
		{
			dst:    []float64{-1, -1, -1},
			x:      []float64{0, 0, 0},
			y:      []float64{1, 2, 3},
			expect: []float64{0, 0, 0},
		},
		{
			dst:    []float64{inf, inf, inf, inf, inf},
			x:      []float64{nan, 1, nan, 1, 0},
			y:      []float64{1, 1, nan, 1, 1},
			expect: []float64{nan, 1, nan, 1, 0},
		},
		{
			dst:    []float64{0, 0, 0, 0, 0},
			x:      []float64{inf, 4, nan, -inf, 9},
			y:      []float64{inf, 4, nan, -inf, 3},
			expect: []float64{nan, 1, nan, nan, 3},
		},
	} {
		xg_ln, yg_ln := 4+j%2, 4+j%3
		v.y, v.x = guardVector(v.y, y_gd, yg_ln), guardVector(v.x, x_gd, xg_ln)
		y, x := v.y[yg_ln:len(v.y)-yg_ln], v.x[xg_ln:len(v.x)-xg_ln]
		v.dst = guardVector(v.dst, dst_gd, xg_ln)
		dst := v.dst[xg_ln : len(v.dst)-xg_ln]
		ret := DivTo(dst, x, y)
		for i := range v.expect {
			if !same(ret[i], v.expect[i]) {
				t.Errorf("Test %d DivTo error at %d Got: %v Expected: %v", j, i, ret[i], v.expect[i])
			}
			if !same(ret[i], dst[i]) {
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

func TestL1Norm(t *testing.T) {
	var t_gd, s_gd float64 = -inf, inf
	for j, v := range []struct {
		s, t   []float64
		expect float64
	}{
		{
			s:      []float64{1},
			t:      []float64{1},
			expect: 0,
		},
		{
			s:      []float64{nan},
			t:      []float64{nan},
			expect: nan,
		},
		{
			s:      []float64{1, 2, 3, 4},
			t:      []float64{1, 2, 3, 4},
			expect: 0,
		},
		{
			s:      []float64{2, 4, 6},
			t:      []float64{1, 2, 3, 4},
			expect: 6,
		},
		{
			s:      []float64{0, 0, 0},
			t:      []float64{1, 2, 3},
			expect: 6,
		},
		{
			s:      []float64{0, -4, -10},
			t:      []float64{1, 2, 3},
			expect: 20,
		},
		{
			s:      []float64{0, 1, 0, 1, 0},
			t:      []float64{1, 1, inf, 1, 1},
			expect: inf,
		},
		{
			s:      []float64{inf, 4, nan, -inf, 9},
			t:      []float64{inf, 4, nan, -inf, 3},
			expect: nan,
		},
	} {
		sg_ln, tg_ln := 4+j%2, 4+j%3
		v.s, v.t = guardVector(v.s, s_gd, sg_ln), guardVector(v.t, t_gd, tg_ln)
		s_lc, t_lc := v.s[sg_ln:len(v.s)-sg_ln], v.t[tg_ln:len(v.t)-tg_ln]
		ret := L1Norm(s_lc, t_lc)
		if !same(ret, v.expect) {
			t.Errorf("Test %d L1Norm error Got: %f Expected: %f", j, ret, v.expect)
		}
		if !isValidGuard(v.s, s_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in s vector %v %v", j, v.s[:sg_ln], v.s[len(v.s)-sg_ln:])
		}
		if !isValidGuard(v.t, t_gd, tg_ln) {
			t.Errorf("Test %d Guard violated in t vector %v %v", j, v.t[:tg_ln], v.t[len(v.t)-tg_ln:])
		}
	}
}

func TestLinfNorm(t *testing.T) {
	var t_gd, s_gd float64 = 0, inf
	for j, v := range []struct {
		s, t   []float64
		expect float64
	}{
		{
			s:      []float64{},
			t:      []float64{},
			expect: 0,
		},
		{
			s:      []float64{1},
			t:      []float64{1},
			expect: 0,
		},
		{
			s:      []float64{nan},
			t:      []float64{nan},
			expect: nan,
		},
		{
			s:      []float64{1, 2, 3, 4},
			t:      []float64{1, 2, 3, 4},
			expect: 0,
		},
		{
			s:      []float64{2, 4, 6},
			t:      []float64{1, 2, 3, 4},
			expect: 3,
		},
		{
			s:      []float64{0, 0, 0},
			t:      []float64{1, 2, 3},
			expect: 3,
		},
		{
			s:      []float64{0, 1, 0, 1, 0},
			t:      []float64{1, 1, inf, 1, 1},
			expect: inf,
		},
		{
			s:      []float64{inf, 4, nan, -inf, 9},
			t:      []float64{inf, 4, nan, -inf, 3},
			expect: 6,
		},
	} {
		sg_ln, tg_ln := 4+j%2, 4+j%3
		v.s, v.t = guardVector(v.s, s_gd, sg_ln), guardVector(v.t, t_gd, tg_ln)
		s_lc, t_lc := v.s[sg_ln:len(v.s)-sg_ln], v.t[tg_ln:len(v.t)-tg_ln]
		ret := LinfNorm(s_lc, t_lc)
		if !same(ret, v.expect) {
			t.Errorf("Test %d LinfNorm error Got: %f Expected: %f", j, ret, v.expect)
		}
		if !isValidGuard(v.s, s_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in s vector %v %v", j, v.s[:sg_ln], v.s[len(v.s)-sg_ln:])
		}
		if !isValidGuard(v.t, t_gd, tg_ln) {
			t.Errorf("Test %d Guard violated in t vector %v %v", j, v.t[:tg_ln], v.t[len(v.t)-tg_ln:])
		}
	}
}
