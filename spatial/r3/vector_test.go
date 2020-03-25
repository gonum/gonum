// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import "testing"

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
		{Vec{0, 0, 0}, Vec{0, 0, 0}, Vec{0, 0, 0}},
		{Vec{1, 0, 0}, Vec{0, 0, 0}, Vec{1, 0, 0}},
		{Vec{1, 2, 3}, Vec{4, 5, 7}, Vec{-3, -3, -4}},
		{Vec{1, -3, 5}, Vec{1, -6, -6}, Vec{0, 3, 11}},
		{Vec{1, 2, 3}, Vec{1, 2, 3}, Vec{}},
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
		{3, Vec{0, 0, 0}, Vec{0, 0, 0}},
		{1, Vec{1, 0, 0}, Vec{1, 0, 0}},
		{0, Vec{1, 0, 0}, Vec{0, 0, 0}},
		{3, Vec{1, 0, 0}, Vec{3, 0, 0}},
		{-1, Vec{1, -3, 5}, Vec{-1, 3, -5}},
		{2, Vec{1, -3, 5}, Vec{2, -6, 10}},
		{10, Vec{1, 2, 3}, Vec{10, 20, 30}},
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
