// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"math"
	"runtime"
	"testing"
)

var (
	nan, inf, ninf float64
)

func init() {
	nan, inf, ninf = math.NaN(), math.Inf(1), math.Inf(-1)
}

func diff(a, b float64) bool {
	return a != b && !math.IsNaN(a) && !math.IsNaN(b) || (math.IsNaN(a) != math.IsNaN(b))
}

func TestAdd(t *testing.T) {
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{[]float64{1}, []float64{0}, []float64{1}},
		{[]float64{1, 2, 3}, []float64{1}, []float64{2, 2, 3}},
		{[]float64{}, []float64{}, []float64{}},
		{[]float64{1}, []float64{nan}, []float64{nan}},
		{[]float64{8, 8, 8, 8, 8},
			[]float64{2, 4, nan, 8, 9},
			[]float64{10, 12, nan, 16, 17}},
		{[]float64{0, 1, 2, 3, 4},
			[]float64{ninf, 4, nan, 8, 9},
			[]float64{ninf, 5, nan, 11, 13}},
	} {
		Add(v.dst, v.src)
		for i := range v.expect {
			if diff(v.dst[i], v.expect[i]) {

				t.Log("Test", j, "Add error at", i, "Got:", v.dst[i], "Expected:", v.expect[i])
				t.Fail()
			}
		}
	}
	runtime.GC()
}

func TestAddConst(t *testing.T) {
	for j, v := range []struct {
		alpha       float64
		src, expect []float64
	}{
		{1, []float64{0}, []float64{1}},
		{5, []float64{}, []float64{}},
		{1, []float64{nan}, []float64{nan}},
		{8, []float64{2, 4, nan, 8, 9}, []float64{10, 12, nan, 16, 17}},
		{inf, []float64{ninf, 4, nan, 8, 9}, []float64{nan, inf, nan, inf, inf}},
	} {
		AddConst(v.alpha, v.src)
		for i := range v.expect {
			if diff(v.src[i], v.expect[i]) {
				t.Log("Test", j, "AddConst error at", i, "Got:", v.src[i], "Expected:", v.expect[i])
				t.Fail()
			}
		}
	}
	runtime.GC()
}

func TestCumSum(t *testing.T) {
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{[]float64{0}, []float64{1}, []float64{1}},
		{[]float64{nan}, []float64{nan}, []float64{nan}},
		{[]float64{0, 0, 0}, []float64{1, 2, 3, 4}, []float64{1, 3, 6}},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3}, []float64{1, 3, 6}},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3, 4}, []float64{1, 3, 6, 10}},
		{[]float64{1, nan, nan, 1, 1},
			[]float64{1, 1, nan, 1, 1},
			[]float64{1, 2, nan, nan, nan}},
		{[]float64{nan, 4, inf, ninf, 9},
			[]float64{inf, 4, nan, ninf, 9},
			[]float64{inf, inf, nan, nan, nan}},
	} {
		ret := CumSum(v.dst, v.src)
		for i := range v.expect {
			if diff(ret[i], v.expect[i]) {
				t.Log("Test", j, "CumSum error at", i, "Got:", ret[i], "Expected:", v.expect[i])
				t.Fail()
			}
			if diff(ret[i], v.dst[i]) {
				t.Log("Test", j, "CumSum ret/dst mismatch", i, "Ret:", ret[i], "Dst:", v.dst[i])
				t.Fail()
			}
		}
	}
	runtime.GC()
}

func TestCumProd(t *testing.T) {
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{[]float64{1}, []float64{1}, []float64{1}},
		{[]float64{nan}, []float64{nan}, []float64{nan}},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3, 4}, []float64{1, 2, 6, 24}},
		{[]float64{0, 0, 0}, []float64{1, 2, 3, 4}, []float64{1, 2, 6}},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3}, []float64{1, 2, 6}},
		{[]float64{nan, 1, nan, 1, 0},
			[]float64{1, 1, nan, 1, 1},
			[]float64{1, 1, nan, nan, nan}},
		{[]float64{nan, 4, nan, ninf, 9},
			[]float64{inf, 4, nan, ninf, 9},
			[]float64{inf, inf, nan, nan, nan}},
	} {
		ret := CumProd(v.dst, v.src)
		for i := range v.expect {
			if diff(ret[i], v.expect[i]) {
				t.Log("Test", j, "CumProd error at", i, "Got:", ret[i], "Expected:", v.expect[i])
				t.Fail()
			}
			if diff(ret[i], v.dst[i]) {
				t.Log("Test", j, "CumProd ret/dst mismatch", i, "Ret:", ret[i], "Dst:", v.dst[i])
				t.Fail()
			}
		}
	}
	runtime.GC()
}

func TestDiv(t *testing.T) {
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{[]float64{1}, []float64{1}, []float64{1}},
		{[]float64{nan}, []float64{nan}, []float64{nan}},
		{[]float64{1, 2, 3, 4}, []float64{1, 2, 3, 4}, []float64{1, 1, 1, 1}},
		{[]float64{2, 4, 6}, []float64{1, 2, 3, 4}, []float64{2, 2, 2}},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3}, []float64{0, 0, 0}},
		{[]float64{nan, 1, nan, 1, 0},
			[]float64{1, 1, nan, 1, 1},
			[]float64{nan, 1, nan, 1, 0}},
		{[]float64{inf, 4, nan, ninf, 9},
			[]float64{inf, 4, nan, ninf, 3},
			[]float64{nan, 1, nan, nan, 3}},
	} {
		Div(v.dst, v.src)
		for i := range v.expect {
			if diff(v.dst[i], v.expect[i]) {
				t.Log("Test", j, "Div error at", i, "Got:", v.dst[i], "Expected:", v.expect[i])
				t.Fail()
			}
		}
	}
	runtime.GC()
}

func TestDivTo(t *testing.T) {
	for j, v := range []struct {
		dst, src, expect []float64
	}{
		{[]float64{1}, []float64{1}, []float64{1}},
		{[]float64{nan}, []float64{nan}, []float64{nan}},
		{[]float64{1, 2, 3, 4}, []float64{1, 2, 3, 4}, []float64{1, 1, 1, 1}},
		{[]float64{2, 4, 6}, []float64{1, 2, 3, 4}, []float64{2, 2, 2}},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3}, []float64{0, 0, 0}},
		{[]float64{nan, 1, nan, 1, 0},
			[]float64{1, 1, nan, 1, 1},
			[]float64{nan, 1, nan, 1, 0}},
		{[]float64{inf, 4, nan, ninf, 9},
			[]float64{inf, 4, nan, ninf, 3},
			[]float64{nan, 1, nan, nan, 3}},
	} {
		ret := DivTo(v.dst, v.dst, v.src)
		for i := range v.expect {
			if diff(ret[i], v.expect[i]) {
				t.Log("Test", j, "DivTo error at", i, "Got:", v.dst[i], "Expected:", v.expect[i])
				t.Fail()
			}
			if diff(ret[i], v.dst[i]) {
				t.Log("Test", j, "DivTo ret/dst mismatch", i, "Ret:", ret[i], "Dst:", v.dst[i])
				t.Fail()
			}
		}
	}
	runtime.GC()
}

func TestL1norm(t *testing.T) {
	for j, v := range []struct {
		s, t   []float64
		expect float64
	}{
		{[]float64{1}, []float64{1}, 0},
		{[]float64{nan}, []float64{nan}, nan},
		{[]float64{1, 2, 3, 4}, []float64{1, 2, 3, 4}, 0},
		{[]float64{2, 4, 6}, []float64{1, 2, 3, 4}, 6},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3}, 6},
		{[]float64{0, -4, -10, 0}, []float64{1, 2, 3}, 20},
		{[]float64{0, 1, 0, 1, 0}, []float64{1, 1, inf, 1, 1}, inf},
		{[]float64{inf, 4, nan, ninf, 9}, []float64{inf, 4, nan, ninf, 3}, nan},
	} {
		ret := L1norm(v.s, v.t)
		if diff(ret, v.expect) {
			t.Log("Test", j, "L1norm error. Got:", ret, "Expected:", v.expect)
			t.Fail()
		}
	}
	runtime.GC()
}

func TestLinfNorm(t *testing.T) {
	for j, v := range []struct {
		s, t   []float64
		expect float64
	}{
		{[]float64{1}, []float64{1}, 0},
		{[]float64{nan}, []float64{nan}, nan},
		{[]float64{1, 2, 3, 4}, []float64{1, 2, 3, 4}, 0},
		{[]float64{2, 4, 6}, []float64{1, 2, 3, 4}, 3},
		{[]float64{0, 0, 0, 0}, []float64{1, 2, 3}, 3},
		{[]float64{0, 1, 0, 1, 0}, []float64{1, 1, inf, 1, 1}, inf},
		{[]float64{inf, 4, nan, ninf, 9}, []float64{inf, 4, nan, ninf, 3}, 6},
	} {
		ret := LinfNorm(v.s, v.t)
		if diff(ret, v.expect) {
			t.Log("Test", j, "LinfNorm error. Got:", ret, "Expected:", v.expect)
			t.Fail()
		}
	}
	runtime.GC()
}
