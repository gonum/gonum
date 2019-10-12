// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrate

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestSimpsons(t *testing.T) {
	x := floats.Span(make([]float64, 1e6), 0, 1)
	for i, test := range []struct {
		x    []float64
		f    func(x float64) float64
		want float64
		tol  float64
	}{
		{
			x:    floats.Span(make([]float64, 1e6), 0, 1),
			f:    func(x float64) float64 { return math.Pi },
			want: math.Pi,
			tol:  1e-10,
		},
		{
			x:    floats.Span(make([]float64, 1e6), 0, 1),
			f:    func(x float64) float64 { return 1.0 },
			want: 1,
			tol:  1e-10,
		},
		{
			x:    floats.Span(make([]float64, 3), 0, 2),
			f:    func(x float64) float64 { return 2*x + 0.5 },
			want: 5,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 10), 0, 2),
			f:    func(x float64) float64 { return 2*x + 0.5 },
			want: 5,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6), 0, 2),
			f:    func(x float64) float64 { return 2*x + 0.5 },
			want: 5,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return x },
			want: 0.5,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6), -1, 1),
			f:    func(x float64) float64 { return x },
			want: 0,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return x + 10 },
			want: 10.5,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return 3*x*x + 10 },
			want: 11,
			tol:  1e-12,
		},
		{
			x:    x,
			f:    func(x float64) float64 { return math.Exp(x) },
			want: 1.7182818284591876,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6), 0, math.Pi),
			f:    func(x float64) float64 { return math.Cos(x) },
			want: 0,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6), 0, 2*math.Pi),
			f:    func(x float64) float64 { return math.Cos(x) },
			want: 0,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6*10), 0, math.Pi),
			f:    func(x float64) float64 { return math.Sin(x) },
			want: 2,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6*10), 0, 0.5*math.Pi),
			f:    func(x float64) float64 { return math.Sin(x) },
			want: 1,
			tol:  1e-12,
		},
		{
			x:    floats.Span(make([]float64, 1e6), 0, 2*math.Pi),
			f:    func(x float64) float64 { return math.Sin(x) },
			want: 0,
			tol:  1e-12,
		},
		{
			x: join(floats.Span(make([]float64, 3), 0, math.Pi/3),
				[]float64{4 * math.Pi / 10},
				floats.Span(make([]float64, 4), 3*math.Pi/4, math.Pi)),
			f:    func(x float64) float64 { return math.Sin(x) },
			want: 2,
			tol:  1e-2,
		},
		{
			x: join(floats.Span(make([]float64, 30), 0, 5*math.Pi/16),
				floats.Span(make([]float64, 100), 3*math.Pi/8, math.Pi)),
			f:    func(x float64) float64 { return math.Sin(x) },
			want: 2,
			tol:  1e-4,
		},
		{
			x: join(floats.Span(make([]float64, 1e5), 0, 15),
				floats.Span(make([]float64, 1e5), 23, 40),
				floats.Span(make([]float64, 1e5), 50, 80),
				floats.Span(make([]float64, 1e5), 90, 100)),
			f:    func(x float64) float64 { return 2 * x },
			want: 10000,
			tol:  1e-9,
		},
		{
			x:    []float64{0, 1, 2},
			f:    func(x float64) float64 { return 2 * x },
			want: 4,
			tol:  1e-12,
		},
		{
			x:    []float64{0, 0.2, 2},
			f:    func(x float64) float64 { return 2 * x },
			want: 4,
			tol:  1e-12,
		},
		{
			x:    []float64{0, 1.89, 2},
			f:    func(x float64) float64 { return 2 * x },
			want: 4,
			tol:  1e-12,
		},
	} {
		y := make([]float64, len(test.x))
		for i, v := range test.x {
			y[i] = test.f(v)
		}
		v := Simpsons(test.x, y)
		if !floats.EqualWithinAbs(v, test.want, test.tol) {
			t.Errorf("test #%d: got=%v want=%f\n", i, v, float64(test.want))
		}
	}
}

func join(slices ...[]float64) []float64 {
	var c []float64
	for _, s := range slices {
		c = append(c, s...)
	}
	return c
}
