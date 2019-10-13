// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrate

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestRomberg(t *testing.T) {
	const n = 1<<8 + 1
	x := floats.Span(make([]float64, n), 0, 1)

	for i, test := range []struct {
		x    []float64
		f    func(x float64) float64
		want float64
		tol  float64
	}{
		{
			x:    x,
			f:    func(x float64) float64 { return x },
			want: 0.5,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return x * x },
			want: 1.0 / 3.0,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return x * x * x },
			want: 1.0 / 4.0,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return math.Sqrt(x) },
			want: 2.0 / 3.0,
			tol:  1e-4,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return math.Sin(math.Pi * x) },
			want: 2.0 / math.Pi,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 3), 0, 1),
			f:    func(x float64) float64 { return x * x },
			want: 1.0 / 3.0,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 3), 0, 1),
			f:    func(x float64) float64 { return x * x * x },
			want: 1.0 / 4.0,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return x * math.Exp(-x) },
			want: (math.Exp(1) - 2) / math.Exp(1),
			tol:  1e-12,
		},
	} {
		n := len(test.x)
		y := make([]float64, n)
		for i, v := range test.x {
			y[i] = test.f(v)
		}

		dx := (test.x[n-1] - test.x[0]) / float64(n-1)
		v := Romberg(y, dx)
		diff := math.Abs(v - test.want)
		if diff > test.tol {
			t.Errorf("test #%d: got=%v want=%v diff=%v\n", i, v, test.want, diff)
		}
	}
}
