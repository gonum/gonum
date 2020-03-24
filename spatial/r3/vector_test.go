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
