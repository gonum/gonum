// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64_test

import (
	"fmt"
	"math"
	"testing"

	. "gonum.org/v1/gonum/internal/asm/f64"
)

// nanwith copied from floats package
func nanwith(payload uint64) float64 {
	const (
		nanBits = 0x7ff8000000000000
		nanMask = 0xfff8000000000000
	)
	return math.Float64frombits(nanBits | (payload &^ nanMask))
}

func TestL2NormUnitary(t *testing.T) {
	const tol = 1e-15

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
		{want: nan, x: []float64{1, inf, 3, nanwith(25), 5}},
		{want: 17.88854381999832, x: []float64{8, -8, 8, -8, 8}},
		{want: 2.23606797749979, x: []float64{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}},
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
	const tol = 1e-15

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
		if !sameApprox(ret, v.want, tol) {
			t.Errorf("Test %d L2NormInc error Got: %f Expected: %f", j, ret, v.want)
		}
		checkValidIncGuard(t, v.x, src_gd, v.inc, g_ln)
	}
}

func TestL2DistanceUnitary(t *testing.T) {
	const tol = 1e-15

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
		if !sameApprox(ret, v.want, tol) {
			t.Errorf("Test %d L2Distance error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}

func BenchmarkL2NormNetlib(b *testing.B) {
	netlib := func(x []float64) (sum float64) {
		var scale float64
		sumSquares := 1.0
		for _, v := range x {
			if v == 0 {
				continue
			}
			absxi := math.Abs(v)
			if math.IsNaN(absxi) {
				return math.NaN()
			}
			if scale < absxi {
				s := scale / absxi
				sumSquares = 1 + sumSquares*s*s
				scale = absxi
			} else {
				s := absxi / scale
				sumSquares += s * s
			}
		}
		if math.IsInf(scale, 1) {
			return math.Inf(1)
		}
		return scale * math.Sqrt(sumSquares)
	}

	tests := []struct {
		name string
		f    func(x []float64) float64
	}{
		{"L2NormUnitaryNetlib", netlib},
		{"L2NormUnitary", L2NormUnitary},
	}
	x[0] = randomSlice(1, 1)[0] // replace the leading zero (edge case)
	for _, test := range tests {
		for _, ln := range []uintptr{1, 3, 10, 30, 1e2, 3e2, 1e3, 3e3, 1e4, 3e4, 1e5} {
			b.Run(fmt.Sprintf("%s-%d", test.name, ln), func(b *testing.B) {
				b.SetBytes(int64(64 * ln))
				x := x[:ln]
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					test.f(x)
				}
			})
		}
	}
}
