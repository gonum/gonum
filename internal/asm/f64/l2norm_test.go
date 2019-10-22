// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import "testing"

func TestL2NormUnitary(t *testing.T) {
	var src_gd float64 = 1
	for j, v := range []struct {
		want float64
		x    []float64
	}{
		{want: 0, x: []float64{}},
		{want: 2, x: []float64{2}},
		{want: 3.7416573867739413, x: []float64{1, 2, 3}},
		{want: 3.7416573867739413, x: []float64{-1, -2, -3}},
		{want: nan, x: []float64{nan}},
		{want: 17.88854381999832, x: []float64{8, -8, 8, -8, 8}},
		{want: 2.23606797749979, x: []float64{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}},
	} {
		g_ln := 4 + j%2
		v.x = guardVector(v.x, src_gd, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L2NormUnitary(src)
		if !within(ret, v.want) {
			t.Errorf("Test %d L2Norm error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}

func TestL2NormInc(t *testing.T) {
	var src_gd float64 = 1
	for j, v := range []struct {
		inc  int
		want float64
		x    []float64
	}{
		{inc: 2, want: 0, x: []float64{}},
		{inc: 3, want: 2, x: []float64{2}},
		{inc: 10, want: 3.7416573867739413, x: []float64{1, 2, 3}},
		{inc: 5, want: 3.7416573867739413, x: []float64{-1, -2, -3}},
		{inc: 3, want: nan, x: []float64{nan}},
		{inc: 15, want: 17.88854381999832, x: []float64{8, -8, 8, -8, 8}},
		{inc: 1, want: 2.23606797749979, x: []float64{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}},
	} {
		g_ln, ln := 4+j%2, len(v.x)
		v.x = guardIncVector(v.x, src_gd, v.inc, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L2NormInc(src, uintptr(ln), uintptr(v.inc))
		if !within(ret, v.want) {
			t.Errorf("Test %d L2NormInc error Got: %f Expected: %f", j, ret, v.want)
		}
		checkValidIncGuard(t, v.x, src_gd, v.inc, g_ln)
	}
}

func TestL2DistanceUnitary(t *testing.T) {
	var src_gd float64 = 1
	for j, v := range []struct {
		want float64
		x, y []float64
	}{
		{want: 0, x: []float64{}, y: []float64{}},
		{want: 2, x: []float64{3}, y: []float64{1}},
		{want: 3.7416573867739413, x: []float64{2, 4, 6}, y: []float64{1, 2, 3}},
		{want: 3.7416573867739413, x: []float64{1, 2, 3}, y: []float64{2, 4, 6}},
		{want: nan, x: []float64{nan}, y: []float64{0}},
		{want: 17.88854381999832, x: []float64{9, -9, 9, -9, 9}, y: []float64{1, -1, 1, -1, 1}},
		{want: 2.23606797749979, x: []float64{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}, y: []float64{0, 2, 0, -2, 0, 2, 0, -2, 0, 2}},
	} {
		g_ln := 4 + j%2
		v.x = guardVector(v.x, src_gd, g_ln)
		v.y = guardVector(v.y, src_gd, g_ln)
		srcX := v.x[g_ln : len(v.x)-g_ln]
		srcY := v.y[g_ln : len(v.y)-g_ln]
		ret := L2DistanceUnitary(srcX, srcY)
		if !within(ret, v.want) {
			t.Errorf("Test %d L2Distance error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}
