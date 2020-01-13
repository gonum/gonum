// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32_test

import (
	"testing"

	. "gonum.org/v1/gonum/internal/asm/f32"
)

func TestL2NormUnitary(t *testing.T) {
	const tol = 1e-7

	var src_gd float32 = 1
	for j, v := range []struct {
		want float32
		x    []float32
	}{
		{want: 0, x: []float32{}},
		{want: 2, x: []float32{2}},
		{want: 3.7416573867739413, x: []float32{1, 2, 3}},
		{want: 3.7416573867739413, x: []float32{-1, -2, -3}},
		{want: nan, x: []float32{nan}},
		{want: 17.88854381999832, x: []float32{8, -8, 8, -8, 8}},
		{want: 2.23606797749979, x: []float32{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}},
	} {
		g_ln := 4 + j%2
		v.x = guardVector(v.x, src_gd, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L2NormUnitary(src)
		if !sameApprox(ret, v.want, tol) {
			t.Errorf("Test %d L2Norm error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}

func TestL2NormInc(t *testing.T) {
	const tol = 1e-7

	var src_gd float32 = 1
	for j, v := range []struct {
		inc  int
		want float32
		x    []float32
	}{
		{inc: 2, want: 0, x: []float32{}},
		{inc: 3, want: 2, x: []float32{2}},
		{inc: 10, want: 3.7416573867739413, x: []float32{1, 2, 3}},
		{inc: 5, want: 3.7416573867739413, x: []float32{-1, -2, -3}},
		{inc: 3, want: nan, x: []float32{nan}},
		{inc: 15, want: 17.88854381999832, x: []float32{8, -8, 8, -8, 8}},
		{inc: 1, want: 2.23606797749979, x: []float32{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}},
	} {
		g_ln, ln := 4+j%2, len(v.x)
		v.x = guardIncVector(v.x, src_gd, v.inc, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L2NormInc(src, uintptr(ln), uintptr(v.inc))
		if !sameApprox(ret, v.want, tol) {
			t.Errorf("Test %d L2NormInc error Got: %f Expected: %f", j, ret, v.want)
		}
		checkValidIncGuard(t, v.x, src_gd, v.inc, g_ln)
	}
}

func TestL2DistanceUnitary(t *testing.T) {
	const tol = 1e-7

	var src_gd float32 = 1
	for j, v := range []struct {
		want float32
		x, y []float32
	}{
		{want: 0, x: []float32{}, y: []float32{}},
		{want: 2, x: []float32{3}, y: []float32{1}},
		{want: 3.7416573867739413, x: []float32{2, 4, 6}, y: []float32{1, 2, 3}},
		{want: 3.7416573867739413, x: []float32{1, 2, 3}, y: []float32{2, 4, 6}},
		{want: nan, x: []float32{nan}, y: []float32{0}},
		{want: 17.88854381999832, x: []float32{9, -9, 9, -9, 9}, y: []float32{1, -1, 1, -1, 1}},
		{want: 2.23606797749979, x: []float32{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}, y: []float32{0, 2, 0, -2, 0, 2, 0, -2, 0, 2}},
	} {
		g_ln := 4 + j%2
		v.x = guardVector(v.x, src_gd, g_ln)
		v.y = guardVector(v.y, src_gd, g_ln)
		srcX := v.x[g_ln : len(v.x)-g_ln]
		srcY := v.y[g_ln : len(v.y)-g_ln]
		ret := L2DistanceUnitary(srcX, srcY)
		if !sameApprox(ret, v.want, tol) {
			t.Errorf("Test %d L2Distance error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}
