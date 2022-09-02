// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
)

func TestAdd(t *testing.T) {
	for _, test := range []struct {
		v1, v2 Vec
		want   Vec
	}{
		{Vec{0, 0, 0}, Vec{0, 0, 0}, Vec{0, 0, 0}},
		{Vec{1, 0, 0}, Vec{0, 0, 0}, Vec{1, 0, 0}},
		{Vec{1, 2, 3}, Vec{4, 5, 7}, Vec{5, 7, 10}},
		{Vec{1, -3, 5}, Vec{1, -6, -6}, Vec{2, -9, -1}},
		{Vec{1, 2, 3}, Vec{-1, -2, -3}, Vec{}},
	} {
		got := Add(test.v1, test.v2)
		if got != test.want {
			t.Errorf(
				"error: %v + %v: got=%v, want=%v",
				test.v1, test.v2, got, test.want,
			)
		}
	}
}

func TestSub(t *testing.T) {
	for _, test := range []struct {
		v1, v2 Vec
		want   Vec
	}{
		{Vec{0, 0, 0}, Vec{0, 0, 0}, Vec{0, 0, 0}},
		{Vec{1, 0, 0}, Vec{0, 0, 0}, Vec{1, 0, 0}},
		{Vec{1, 2, 3}, Vec{4, 5, 7}, Vec{-3, -3, -4}},
		{Vec{1, -3, 5}, Vec{1, -6, -6}, Vec{0, 3, 11}},
		{Vec{1, 2, 3}, Vec{1, 2, 3}, Vec{}},
	} {
		got := Sub(test.v1, test.v2)
		if got != test.want {
			t.Errorf(
				"error: %v - %v: got=%v, want=%v",
				test.v1, test.v2, got, test.want,
			)
		}
	}
}

func TestScale(t *testing.T) {
	for _, test := range []struct {
		a    float64
		v    Vec
		want Vec
	}{
		{3, Vec{0, 0, 0}, Vec{0, 0, 0}},
		{1, Vec{1, 0, 0}, Vec{1, 0, 0}},
		{0, Vec{1, 0, 0}, Vec{0, 0, 0}},
		{3, Vec{1, 0, 0}, Vec{3, 0, 0}},
		{-1, Vec{1, -3, 5}, Vec{-1, 3, -5}},
		{2, Vec{1, -3, 5}, Vec{2, -6, 10}},
		{10, Vec{1, 2, 3}, Vec{10, 20, 30}},
	} {
		got := Scale(test.a, test.v)
		if got != test.want {
			t.Errorf(
				"error: %v * %v: got=%v, want=%v",
				test.a, test.v, got, test.want)
		}
	}
}

func TestDot(t *testing.T) {
	for _, test := range []struct {
		u, v Vec
		want float64
	}{
		{Vec{1, 2, 3}, Vec{1, 2, 3}, 14},
		{Vec{1, 0, 0}, Vec{1, 0, 0}, 1},
		{Vec{1, 0, 0}, Vec{0, 1, 0}, 0},
		{Vec{1, 0, 0}, Vec{0, 1, 1}, 0},
		{Vec{1, 1, 1}, Vec{-1, -1, -1}, -3},
		{Vec{1, 2, 2}, Vec{-0.3, 0.4, -1.2}, -1.9},
	} {
		{
			got := Dot(test.u, test.v)
			if got != test.want {
				t.Errorf(
					"error: %v · %v: got=%v, want=%v",
					test.u, test.v, got, test.want,
				)
			}
		}
		{
			got := Dot(test.v, test.u)
			if got != test.want {
				t.Errorf(
					"error: %v · %v: got=%v, want=%v",
					test.v, test.u, got, test.want,
				)
			}
		}
	}
}

func TestCross(t *testing.T) {
	for _, test := range []struct {
		v1, v2, want Vec
	}{
		{Vec{1, 0, 0}, Vec{1, 0, 0}, Vec{0, 0, 0}},
		{Vec{1, 0, 0}, Vec{0, 1, 0}, Vec{0, 0, 1}},
		{Vec{0, 1, 0}, Vec{1, 0, 0}, Vec{0, 0, -1}},
		{Vec{1, 2, 3}, Vec{-4, 5, -6}, Vec{-27, -6, 13}},
		{Vec{1, 2, 3}, Vec{1, 2, 3}, Vec{}},
		{Vec{1, 2, 3}, Vec{2, 3, 4}, Vec{-1, 2, -1}},
	} {
		got := Cross(test.v1, test.v2)
		if got != test.want {
			t.Errorf(
				"error: %v × %v = %v, want %v",
				test.v1, test.v2, got, test.want,
			)
		}
	}
}

func TestNorm(t *testing.T) {
	for _, test := range []struct {
		v    Vec
		want float64
	}{
		{Vec{0, 0, 0}, 0},
		{Vec{0, 1, 0}, 1},
		{Vec{3, -4, 12}, 13},
		{Vec{1, 1e-16, 1e-32}, 1},
		{Vec{-0, 4.3145006366056343748277397783556100978621924913975e-196, 4.3145006366056343748277397783556100978621924913975e-196}, 6.101625315155041e-196},
	} {
		if got, want := Norm(test.v), test.want; got != want {
			t.Errorf("|%v| = %v, want %v", test.v, got, want)
		}
	}
}

func TestNorm2(t *testing.T) {
	for _, test := range []struct {
		v    Vec
		want float64
	}{
		{Vec{0, 0, 0}, 0},
		{Vec{0, 1, 0}, 1},
		{Vec{1, 1, 1}, 3},
		{Vec{1, 2, 3}, 14},
		{Vec{3, -4, 12}, 169},
		{Vec{1, 1e-16, 1e-32}, 1},
		// This will underflow and return zero.
		{Vec{-0, 4.3145006366056343748277397783556100978621924913975e-196, 4.3145006366056343748277397783556100978621924913975e-196}, 0},
	} {
		if got, want := Norm2(test.v), test.want; got != want {
			t.Errorf("|%v|^2 = %v, want %v", test.v, got, want)
		}
	}
}

func TestUnit(t *testing.T) {
	for _, test := range []struct {
		v, want Vec
	}{
		{Vec{}, Vec{math.NaN(), math.NaN(), math.NaN()}},
		{Vec{1, 0, 0}, Vec{1, 0, 0}},
		{Vec{0, 1, 0}, Vec{0, 1, 0}},
		{Vec{0, 0, 1}, Vec{0, 0, 1}},
		{Vec{1, 1, 1}, Vec{1. / math.Sqrt(3), 1. / math.Sqrt(3), 1. / math.Sqrt(3)}},
		{Vec{1, 1e-16, 1e-32}, Vec{1, 1e-16, 1e-32}},
	} {
		got := Unit(test.v)
		if !vecEqual(got, test.want) {
			t.Errorf(
				"Normalize(%v) = %v, want %v",
				test.v, got, test.want,
			)
		}
		if test.v == (Vec{}) {
			return
		}
		if n, want := Norm(got), 1.0; n != want {
			t.Errorf("|%v| = %v, want 1", got, n)
		}
	}
}

func TestCos(t *testing.T) {
	for _, test := range []struct {
		v1, v2 Vec
		want   float64
	}{
		{Vec{1, 1, 1}, Vec{1, 1, 1}, 1},
		{Vec{1, 1, 1}, Vec{-1, -1, -1}, -1},
		{Vec{1, 1, 1}, Vec{1, -1, 1}, 1.0 / 3},
		{Vec{1, 0, 0}, Vec{1, 0, 0}, 1},
		{Vec{1, 0, 0}, Vec{0, 1, 0}, 0},
		{Vec{1, 0, 0}, Vec{0, 1, 1}, 0},
		{Vec{1, 0, 0}, Vec{-1, 0, 0}, -1},
	} {
		tol := 1e-14
		got := Cos(test.v1, test.v2)
		if !scalar.EqualWithinAbs(got, test.want, tol) {
			t.Errorf("cos(%v, %v)= %v, want %v",
				test.v1, test.v2, got, test.want,
			)
		}
	}
}

func TestRotate(t *testing.T) {
	const tol = 1e-14
	for _, test := range []struct {
		v, axis Vec
		alpha   float64
		want    Vec
	}{
		{Vec{1, 0, 0}, Vec{1, 0, 0}, math.Pi / 2, Vec{1, 0, 0}},
		{Vec{1, 0, 0}, Vec{1, 0, 0}, 0, Vec{1, 0, 0}},
		{Vec{1, 0, 0}, Vec{1, 0, 0}, 2 * math.Pi, Vec{1, 0, 0}},
		{Vec{1, 0, 0}, Vec{0, 0, 0}, math.Pi / 2, Vec{math.NaN(), math.NaN(), math.NaN()}},
		{Vec{1, 0, 0}, Vec{0, 1, 0}, math.Pi / 2, Vec{0, 0, -1}},
		{Vec{1, 0, 0}, Vec{0, 1, 0}, math.Pi, Vec{-1, 0, 0}},
		{Vec{2, 0, 0}, Vec{0, 1, 0}, math.Pi, Vec{-2, 0, 0}},
		{Vec{1, 2, 3}, Vec{1, 1, 1}, 2. / 3. * math.Pi, Vec{3, 1, 2}},
	} {
		got := Rotate(test.v, test.alpha, test.axis)
		if !vecApproxEqual(got, test.want, tol) {
			t.Errorf(
				"quat rotate(%v, %v, %v)= %v, want=%v",
				test.v, test.alpha, test.axis, got, test.want,
			)
		}

		var gotv mat.VecDense
		gotv.MulVec(NewRotation(test.alpha, test.axis).Mat(), vecDense(test.v))
		got = vec(gotv)
		if !vecApproxEqual(got, test.want, tol) {
			t.Errorf(
				"matrix rotate(%v, %v, %v)= %v, want=%v",
				test.v, test.alpha, test.axis, got, test.want,
			)
		}
	}
}

var vectorFields = []struct {
	field      func(Vec) Vec
	divergence func(Vec) float64
	jacobian   func(Vec) *Mat
}{
	{
		field: func(v Vec) Vec {
			// (x*y*z, y*z, z*x)
			return Vec{X: v.X * v.Y * v.Z, Y: v.Y * v.Z, Z: v.Z * v.X}
		},
		divergence: func(v Vec) float64 {
			return v.X + v.Y*v.Z + v.Z
		},
		jacobian: func(v Vec) *Mat {
			return NewMat([]float64{
				v.Y * v.Z, v.X * v.Z, v.X * v.Y,
				0, v.Z, v.Y,
				v.Z, 0, v.X,
			})
		},
	},
	{
		field: func(v Vec) Vec {
			// (x*y*z*cos(y), y*z+sin(x), z*x*sin(y))
			sx := math.Sin(v.X)
			sy, cy := math.Sincos(v.Y)
			return Vec{
				X: v.X * v.Y * v.Z * cy,
				Y: v.Y*v.Z + sx,
				Z: v.Z * v.X * sy,
			}
		},
		divergence: func(v Vec) float64 {
			sy, cy := math.Sincos(v.Y)
			return v.X*sy + v.Y*v.Z*cy + v.Z
		},
		jacobian: func(v Vec) *Mat {
			cx := math.Cos(v.X)
			sy, cy := math.Sincos(v.Y)
			return NewMat([]float64{
				v.Y * v.Z * cy, v.X*v.Z*cy - v.X*v.Y*v.Z*sy, v.X * v.Y * cy,
				cx, v.Z, v.Y,
				v.Z * sy, v.X * v.Z * cy, v.X * sy,
			})
		},
	},
}

func TestDivergence(t *testing.T) {
	const (
		tol = 1e-10
		h   = 1e-2
	)
	step := Vec{X: h, Y: h, Z: h}
	rnd := rand.New(rand.NewSource(1))
	for _, test := range vectorFields {
		for i := 0; i < 30; i++ {
			p := randomVec(rnd)
			got := Divergence(p, step, test.field)
			want := test.divergence(p)
			if math.Abs(got-want) > tol {
				t.Errorf("result out of tolerance. got %v, want %v", got, want)
			}
		}
	}
}

func TestGradient(t *testing.T) {
	const (
		tol = 1e-6
		h   = 1e-5
	)
	step := Vec{X: h, Y: h, Z: h}
	rnd := rand.New(rand.NewSource(1))
	for _, test := range scalarFields {
		for i := 0; i < 30; i++ {
			p := randomVec(rnd)
			got := Gradient(p, step, test.field)
			want := test.gradient(p)
			if !vecApproxEqual(got, want, tol) {
				t.Errorf("result out of tolerance. got %v, want %v", got, want)
			}
		}
	}
}

func vecDense(v Vec) *mat.VecDense {
	return mat.NewVecDense(3, []float64{v.X, v.Y, v.Z})
}

func vec(v mat.VecDense) Vec {
	if v.Len() != 3 {
		panic(mat.ErrShape)
	}
	return Vec{v.AtVec(0), v.AtVec(1), v.AtVec(2)}
}

func vecIsNaN(v Vec) bool {
	return math.IsNaN(v.X) && math.IsNaN(v.Y) && math.IsNaN(v.Z)
}

func vecIsNaNAny(v Vec) bool {
	return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z)
}

func vecEqual(a, b Vec) bool {
	if vecIsNaNAny(a) || vecIsNaNAny(b) {
		return vecIsNaN(a) && vecIsNaN(b)
	}
	return a == b
}

func vecApproxEqual(a, b Vec, tol float64) bool {
	if vecIsNaNAny(a) || vecIsNaNAny(b) {
		return vecIsNaN(a) && vecIsNaN(b)
	}
	return scalar.EqualWithinAbs(a.X, b.X, tol) &&
		scalar.EqualWithinAbs(a.Y, b.Y, tol) &&
		scalar.EqualWithinAbs(a.Z, b.Z, tol)
}
