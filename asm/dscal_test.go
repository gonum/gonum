// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
	"testing"
)

var dscalTests = []struct {
	alpha float64
	x     []float64
	want  []float64
}{
	{
		alpha: 0,
		x:     []float64{1},
		want:  []float64{0},
	},
	{
		alpha: 1,
		x:     []float64{1},
		want:  []float64{1},
	},
	{
		alpha: 2,
		x:     []float64{1, -2},
		want:  []float64{2, -4},
	},
	{
		alpha: 2,
		x:     []float64{1, -2, 3},
		want:  []float64{2, -4, 6},
	},
	{
		alpha: 2,
		x:     []float64{1, -2, 3, 4},
		want:  []float64{2, -4, 6, 8},
	},
	{
		alpha: 2,
		x:     []float64{1, -2, 3, 4, -5},
		want:  []float64{2, -4, 6, 8, -10},
	},
	{
		alpha: 2,
		x:     []float64{0, 1, -2, 3, 4, -5, 6, -7},
		want:  []float64{0, 2, -4, 6, 8, -10, 12, -14},
	},
	{
		alpha: 2,
		x:     []float64{0, 1, -2, 3, 4, -5, 6, -7, 8},
		want:  []float64{0, 2, -4, 6, 8, -10, 12, -14, 16},
	},
	{
		alpha: 2,
		x:     []float64{0, 1, -2, 3, 4, -5, 6, -7, 8, 9},
		want:  []float64{0, 2, -4, 6, 8, -10, 12, -14, 16, 18},
	},
}

func TestDscalUnitary(t *testing.T) {
	for i, test := range dscalTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		prefix := fmt.Sprintf("test %v (x*=a)", i)
		x, xFront, xBack := newGuardedVector(test.x, 1)
		DscalUnitary(test.alpha, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}

func TestDscalUnitaryTo(t *testing.T) {
	for i, test := range dscalTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		// Test dst = alpha * x.
		prefix := fmt.Sprintf("test %v (dst=a*x)", i)
		x, xFront, xBack := newGuardedVector(test.x, 1)
		dst, dstFront, dstBack := newGuardedVector(test.x, 1)
		DscalUnitaryTo(dst, test.alpha, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(dstFront) || !allNaN(dstBack) {
			t.Errorf(msgGuard, prefix, "dst", dstFront, dstBack)
		}
		if !equalStrided(test.x, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}

		if !equalStrided(test.want, dst, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, dst)
		}

		// Test x = alpha * x.
		prefix = fmt.Sprintf("test %v (x=a*x)", i)
		x, xFront, xBack = newGuardedVector(test.x, 1)
		DscalUnitaryTo(x, test.alpha, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}
