// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128_test

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	. "gonum.org/v1/gonum/internal/asm/c128"
)

// nanwith copied from floats package
func nanwith(payload uint64) complex128 {
	const (
		nanBits = 0x7ff8000000000000
		nanMask = 0xfff8000000000000
	)
	nan := math.Float64frombits(nanBits | (payload &^ nanMask))
	return complex(nan, nan)
}

func TestL2NormUnitary(t *testing.T) {
	const tol = 1e-15

	var src_gd complex128 = 1
	for j, v := range []struct {
		want float64
		x    []complex128
	}{
		{want: 0, x: []complex128{}},
		{want: 2, x: []complex128{2}},
		{want: 2, x: []complex128{2i}},
		{want: math.Sqrt(8), x: []complex128{2 + 2i}},
		{want: 3.7416573867739413, x: []complex128{1, 2, 3}},
		{want: 3.7416573867739413, x: []complex128{-1, -2, -3}},
		{want: 3.7416573867739413, x: []complex128{1i, 2i, 3i}},
		{want: 3.7416573867739413, x: []complex128{-1i, -2i, -3i}},
		{want: math.Sqrt(28), x: []complex128{1 + 1i, 2 + 2i, 3 + 3i}},
		{want: math.Sqrt(28), x: []complex128{-1 - 1i, -2 - 2i, -3 - 3i}},
		{want: nan, x: []complex128{cnan}},
		{want: nan, x: []complex128{1, cinf, 3, nanwith(25), 5}},
		{want: 17.88854381999832, x: []complex128{8, -8, 8, -8, 8}},
		{want: 2.23606797749979, x: []complex128{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}},
		{want: 17.88854381999832, x: []complex128{8i, -8i, 8i, -8i, 8i}},
		{want: 2.23606797749979, x: []complex128{0, 1i, 0, -1i, 0, 1i, 0, -1i, 0, 1i}},
		{want: math.Sqrt(640), x: []complex128{8 + 8i, -8 - 8i, 8 + 8i, -8 - 8i, 8 + 8i}},
		{want: math.Sqrt(10), x: []complex128{0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i}},
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

func TestL2DistanceUnitary(t *testing.T) {
	const tol = 1e-15

	var src_gd complex128 = 1
	for j, v := range []struct {
		want float64
		x, y []complex128
	}{
		{want: 0, x: []complex128{}, y: []complex128{}},
		{want: 2, x: []complex128{3}, y: []complex128{1}},
		{want: 2, x: []complex128{3i}, y: []complex128{1i}},
		{want: 3.7416573867739413, x: []complex128{2, 4, 6}, y: []complex128{1, 2, 3}},
		{want: 3.7416573867739413, x: []complex128{1, 2, 3}, y: []complex128{2, 4, 6}},
		{want: 3.7416573867739413, x: []complex128{2i, 4i, 6i}, y: []complex128{1i, 2i, 3i}},
		{want: 3.7416573867739413, x: []complex128{1i, 2i, 3i}, y: []complex128{2i, 4i, 6i}},
		{want: math.Sqrt(28), x: []complex128{2 + 2i, 4 + 4i, 6 + 6i}, y: []complex128{1 + 1i, 2 + 2i, 3 + 3i}},
		{want: math.Sqrt(28), x: []complex128{1 + 1i, 2 + 2i, 3 + 3i}, y: []complex128{2 + 2i, 4 + 4i, 6 + 6i}},
		{want: nan, x: []complex128{cnan}, y: []complex128{0}},
		{want: 17.88854381999832, x: []complex128{9, -9, 9, -9, 9}, y: []complex128{1, -1, 1, -1, 1}},
		{want: 2.23606797749979, x: []complex128{0, 1, 0, -1, 0, 1, 0, -1, 0, 1}, y: []complex128{0, 2, 0, -2, 0, 2, 0, -2, 0, 2}},
		{want: 17.88854381999832, x: []complex128{9i, -9i, 9i, -9i, 9i}, y: []complex128{1i, -1i, 1i, -1i, 1i}},
		{want: 2.23606797749979, x: []complex128{0, 1i, 0, -1i, 0, 1i, 0, -1i, 0, 1i}, y: []complex128{0, 2i, 0, -2i, 0, 2i, 0, -2i, 0, 2i}},
		{want: math.Sqrt(640), x: []complex128{9 + 9i, -9 - 9i, 9 + 9i, -9 - 9i, 9 + 9i}, y: []complex128{1 + 1i, -1 - 1i, 1 + 1i, -1 - 1i, 1 + 1i}},
		{want: math.Sqrt(10), x: []complex128{0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i}, y: []complex128{0, 2 + 2i, 0, -2 - 2i, 0, 2 + 2i, 0, -2 - 2i, 0, 2 + 2i}},
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
	netlib := func(x []complex128) (sum float64) {
		var scale float64
		sumSquares := 1.0
		for _, v := range x {
			if v == 0 {
				continue
			}
			absxi := cmplx.Abs(v)
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
		f    func(x []complex128) float64
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
