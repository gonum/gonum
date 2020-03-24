// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import "testing"

func TestAdd(t *testing.T) {

	var (
		v1   = Vec{1, 2, 3}
		v2   = Vec{-1, -2, -3}
		got  = v1.Add(v2)
		want Vec
	)

	if got != want {
		t.Fatalf("invalid v1+v2: got=%v, want=%v", got, want)
	}
}

func TestSub(t *testing.T) {
	var (
		v    = Vec{1, 2, 3}
		got  = v.Sub(v)
		want Vec
	)

	if got != want {
		t.Fatalf("invalid v-v: got=%v, want=%v", got, want)
	}
}

func TestScale(t *testing.T) {
	var (
		v    = Vec{1, 2, 3}
		got  = v.Scale(10)
		want = Vec{10, 20, 30}
	)

	if got != want {
		t.Fatalf("invalid f.v: got=%v, want=%v", got, want)
	}
}

func TestDot(t *testing.T) {
	for _, test := range []struct {
		u, v Vec
		want float64
	}{
		{
			u:    Vec{1, 2, 3},
			v:    Vec{1, 2, 3},
			want: 14,
		},
	} {
		t.Run("", func(t *testing.T) {
			got := Dot(test.u, test.v)
			if got != test.want {
				t.Fatalf("invalid dot product: got=%v, want=%v", got, test.want)
			}
		})
	}
}

func TestCross(t *testing.T) {
	for _, test := range []struct {
		u, v Vec
		want Vec
	}{
		{
			u:    Vec{1, 2, 3},
			v:    Vec{1, 2, 3},
			want: Vec{},
		},
		{
			u:    Vec{1, 2, 3},
			v:    Vec{2, 3, 4},
			want: Vec{-1, 2, -1},
		},
	} {
		t.Run("", func(t *testing.T) {
			got := Cross(test.u, test.v)
			if got != test.want {
				t.Fatalf("invalid cross product: got=%v, want=%v", got, test.want)
			}
		})
	}
}
