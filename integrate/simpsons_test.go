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
	const (
		n   = 1e6
		tol = 1e-12
	)
	x := floats.Span(make([]float64, n), 0, 1)
	for i, test := range []struct {
		x         []float64
		f         func(x float64) float64
		want      float64
		tolerance float64
	}{
		{
			x:         floats.Span(make([]float64, n), 0, 1),
			f:         func(x float64) float64 { return math.Pi },
			want:      math.Pi,
			tolerance: 1e-10,
		},
		{
			x:         floats.Span(make([]float64, n), 0, 1),
			f:         func(x float64) float64 { return 1.0 },
			want:      1,
			tolerance: 1e-10,
		},
		{
			x: []float64{0.0, 0.1, 1.0},
			f: func(x float64) float64 {
				switch {
				case 1e-31 <= x || x <= 1:
					return 1.0
				default:
					return 0
				}
			},
			want:      1,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, 3), 0, 2),
			f:         func(x float64) float64 { return 2*x + 0.5 },
			want:      5,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, 10), 0, 2),
			f:         func(x float64) float64 { return 2*x + 0.5 },
			want:      5,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n), 0, 2),
			f:         func(x float64) float64 { return 2*x + 0.5 },
			want:      5,
			tolerance: tol,
		},
		{
			x:         x,
			f:         func(x float64) float64 { return x },
			want:      0.5,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n), -1, 1),
			f:         func(x float64) float64 { return x },
			want:      0,
			tolerance: tol,
		},
		{
			x:         x,
			f:         func(x float64) float64 { return x + 10 },
			want:      10.5,
			tolerance: tol,
		},
		{
			x:         x,
			f:         func(x float64) float64 { return 3*x*x + 10 },
			want:      11,
			tolerance: tol,
		},
		{
			x:         x,
			f:         func(x float64) float64 { return math.Exp(x) },
			want:      1.7182818284591876,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n), 0, math.Pi),
			f:         func(x float64) float64 { return math.Cos(x) },
			want:      0,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n), 0, 2*math.Pi),
			f:         func(x float64) float64 { return math.Cos(x) },
			want:      0,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n*10), 0, math.Pi),
			f:         func(x float64) float64 { return math.Sin(x) },
			want:      2,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n*10), 0, 0.5*math.Pi),
			f:         func(x float64) float64 { return math.Sin(x) },
			want:      1,
			tolerance: tol,
		},
		{
			x:         floats.Span(make([]float64, n), 0, 2*math.Pi),
			f:         func(x float64) float64 { return math.Sin(x) },
			want:      0,
			tolerance: tol,
		},
		{
			x: join(floats.Span(make([]float64, 3), 0, math.Pi/3),
				[]float64{4 * math.Pi / 10},
				floats.Span(make([]float64, 4), 3*math.Pi/4, math.Pi)),
			f:         func(x float64) float64 { return math.Sin(x) },
			want:      2,
			tolerance: 1e-2,
		},
		{
			x: join(floats.Span(make([]float64, 30), 0, 5*math.Pi/16),
				floats.Span(make([]float64, 100), 3*math.Pi/8, math.Pi)),
			f:         func(x float64) float64 { return math.Sin(x) },
			want:      2,
			tolerance: 1e-4,
		},
		{
			x: join(floats.Span(make([]float64, 1e5), 0, 15),
				floats.Span(make([]float64, 1e5), 23, 40),
				floats.Span(make([]float64, 1e5), 50, 80),
				floats.Span(make([]float64, 1e5), 90, 100)),
			f:         func(x float64) float64 { return 2 * x },
			want:      10000,
			tolerance: 1e-9,
		},
	} {
		y := make([]float64, len(test.x))
		for i, v := range test.x {
			y[i] = test.f(v)
		}
		v := Simpsons(test.x, y)
		if !floats.EqualWithinAbs(v, test.want, test.tolerance) {
			t.Errorf("test #%d: got=%v want=%f\n", i, v, float64(test.want))
		}
	}
}

func TestSimpsonsHandlesNaN(t *testing.T) {
	x := floats.Span(make([]float64, 1e3), 0, 1)
	y := floats.Span(make([]float64, 1e3), 0, math.NaN())

	v := Simpsons(x, y)
	if !math.IsNaN(v) {
		t.Error("integrate: test expects a NaN from the calling function")
	}
}

func TestSimpsonsPanics(t *testing.T) {
	for i, test := range []struct {
		x    []float64
		y    []float64
		want string
	}{
		{
			x:    floats.Span(make([]float64, 100), 0, 2*math.E),
			y:    []float64{0, 1, 2, 3},
			want: "integrate: slice length mismatch",
		},
		{
			x:    []float64{0},
			y:    []float64{1e2},
			want: "integrate: input data too small",
		},
		{
			x:    []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.2, 0.3, 0.1},
			y:    []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: "integrate: must be sorted",
		},
		{
			x:    []float64{0, 0, 1},
			y:    []float64{0, 0, 1},
			want: "integrate: at least three unique points are required",
		},
	} {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("test #%d did not panic: expectedPanic=\"%s\"", i, test.want)
			}
		}()
		Simpsons(test.x, test.y)
	}
}

func join(slices ...[]float64) []float64 {
	var c []float64
	for _, s := range slices {
		c = append(c, s...)
	}
	return c
}
