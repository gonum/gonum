// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r2

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	for _, test := range []struct {
		v1, v2 Vec
		want   Vec
	}{
		{Vec{0, 0}, Vec{0, 0}, Vec{0, 0}},
		{Vec{1, 0}, Vec{0, 0}, Vec{1, 0}},
		{Vec{1, 2}, Vec{3, 4}, Vec{4, 6}},
		{Vec{1, -3}, Vec{1, -6}, Vec{2, -9}},
		{Vec{1, 2}, Vec{-1, -2}, Vec{}},
	} {
		t.Run("", func(t *testing.T) {
			got := test.v1.Add(test.v2)

			if got != test.want {
				t.Fatalf(
					"error: %v + %v: got=%v, want=%v",
					test.v1, test.v2, got, test.want,
				)
			}
		})
	}
}

func TestSub(t *testing.T) {
	for _, test := range []struct {
		v1, v2 Vec
		want   Vec
	}{
		{Vec{0, 0}, Vec{0, 0}, Vec{0, 0}},
		{Vec{1, 0}, Vec{0, 0}, Vec{1, 0}},
		{Vec{1, 2}, Vec{3, 4}, Vec{-2, -2}},
		{Vec{1, -3}, Vec{1, -6}, Vec{0, 3}},
		{Vec{1, 2}, Vec{1, 2}, Vec{}},
	} {
		t.Run("", func(t *testing.T) {
			got := test.v1.Sub(test.v2)
			if got != test.want {
				t.Fatalf(
					"error: %v - %v: got=%v, want=%v",
					test.v1, test.v2, got, test.want,
				)
			}
		})
	}
}

func TestScale(t *testing.T) {
	for _, test := range []struct {
		a    float64
		v    Vec
		want Vec
	}{
		{3, Vec{0, 0}, Vec{0, 0}},
		{1, Vec{1, 0}, Vec{1, 0}},
		{0, Vec{1, 0}, Vec{0, 0}},
		{3, Vec{1, 0}, Vec{3, 0}},
		{-1, Vec{1, -3}, Vec{-1, 3}},
		{2, Vec{1, -3}, Vec{2, -6}},
		{10, Vec{1, 2}, Vec{10, 20}},
	} {
		t.Run("", func(t *testing.T) {
			got := test.v.Scale(test.a)
			if got != test.want {
				t.Fatalf(
					"error: %v * %v: got=%v, want=%v",
					test.a, test.v, got, test.want)
			}
		})
	}
}

func TestDot(t *testing.T) {
	for _, test := range []struct {
		u, v Vec
		want float64
	}{
		{Vec{1, 2}, Vec{1, 2}, 5},
		{Vec{1, 0}, Vec{1, 0}, 1},
		{Vec{1, 0}, Vec{0, 1}, 0},
		{Vec{1, 0}, Vec{0, 1}, 0},
		{Vec{1, 1}, Vec{-1, -1}, -2},
		{Vec{1, 2}, Vec{-0.3, 0.4}, 0.5},
	} {
		t.Run("", func(t *testing.T) {
			{
				got := test.u.Dot(test.v)
				if got != test.want {
					t.Fatalf(
						"error: %v · %v: got=%v, want=%v",
						test.u, test.v, got, test.want,
					)
				}
			}
			{
				got := test.v.Dot(test.u)
				if got != test.want {
					t.Fatalf(
						"error: %v · %v: got=%v, want=%v",
						test.v, test.u, got, test.want,
					)
				}
			}
		})
	}
}

func TestCross(t *testing.T) {
	for _, test := range []struct {
		v1, v2 Vec
		want   float64
	}{
		{Vec{1, 0}, Vec{1, 0}, 0},
		{Vec{1, 0}, Vec{0, 1}, 1},
		{Vec{0, 1}, Vec{1, 0}, -1},
		{Vec{1, 2}, Vec{-4, 5}, 13},
		{Vec{1, 2}, Vec{2, 3}, -1},
	} {
		t.Run("", func(t *testing.T) {
			got := test.v1.Cross(test.v2)
			if got != test.want {
				t.Fatalf(
					"error: %v × %v = %v, want %v",
					test.v1, test.v2, got, test.want,
				)
			}
		})
	}
}

func TestNorm(t *testing.T) {
	for _, test := range []struct {
		v    Vec
		want float64
	}{
		{Vec{0, 0}, 0},
		{Vec{0, 1}, 1},
		{Vec{1, 1}, math.Sqrt2},
		{Vec{1, 2}, math.Sqrt(5)},
		{Vec{3, -4}, 5},
		{Vec{1, 1e-16}, 1},
		{Vec{4.3145006366056343748277397783556100978621924913975e-196, 4.3145006366056343748277397783556100978621924913975e-196}, 6.101625315155041e-196},
	} {
		t.Run("", func(t *testing.T) {
			if got, want := Norm(test.v), test.want; got != want {
				t.Fatalf("|%v| = %v, want %v", test.v, got, want)
			}
		})
	}
}

func TestNorm2(t *testing.T) {
	for _, test := range []struct {
		v    Vec
		want float64
	}{
		{Vec{0, 0}, 0},
		{Vec{0, 1}, 1},
		{Vec{1, 1}, 2},
		{Vec{1, 2}, 5},
		{Vec{3, -4}, 25},
		{Vec{1, 1e-16}, 1},
		// This will underflow and return zero.
		{Vec{4.3145006366056343748277397783556100978621924913975e-196, 4.3145006366056343748277397783556100978621924913975e-196}, 0},
	} {
		t.Run("", func(t *testing.T) {
			if got, want := Norm2(test.v), test.want; got != want {
				t.Fatalf("|%v|^2 = %v, want %v", test.v, got, want)
			}
		})
	}
}
