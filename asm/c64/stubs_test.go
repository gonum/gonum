// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c64

import (
	"runtime"
	"testing"
)

func TestAxpyUnitary(t *testing.T) {
	for j, v := range []struct {
		a    complex64
		x, y []complex64
		ex   []complex64
	}{
		{1 + 1i, []complex64{1}, []complex64{1i}, []complex64{1 + 2i}},
		{1 + 2i, []complex64{0, 0, 0}, []complex64{1, 1, 1}, []complex64{1, 1, 1}},
		{1 + 2i, []complex64{0, 0}, []complex64{1, 1, 1}, []complex64{1, 1}},
		{1 + 2i, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{-1 + 1i, 1i, -1 + 1i}},
		{-1i, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{2, 3, 2}},
	} {
		AxpyUnitary(v.a, v.x, v.y)
		for i := range v.ex {
			if v.y[i] != v.ex[i] {
				t.Log("Test", j, "Unexpected result at", i, "Got:", v.y[i], "Expected:", v.ex[i])
				t.Fail()
			}
		}
		runtime.GC()
	}
}

func TestAxpyUnitaryTo(t *testing.T) {
	for j, v := range []struct {
		a         complex64
		dst, x, y []complex64
		ex        []complex64
	}{
		{1 + 1i, []complex64{5}, []complex64{1}, []complex64{1i}, []complex64{1 + 2i}},
		{1 + 2i, []complex64{0, 0, 0}, []complex64{0, 0, 0}, []complex64{1, 1, 1}, []complex64{1, 1, 1}},
		{1 + 2i, []complex64{0, 0, 0}, []complex64{0, 0}, []complex64{1, 1, 1}, []complex64{1, 1}},
		{1 + 2i, []complex64{1i, 1i, 1i}, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{-1 + 1i, 1i, -1 + 1i}},
		{-1i, []complex64{1i, 1i, 1i}, []complex64{1i, 1i, 1i}, []complex64{1, 2, 1}, []complex64{2, 3, 2}},
	} {
		AxpyUnitaryTo(v.dst, v.a, v.x, v.y)
		for i := range v.ex {
			if v.dst[i] != v.ex[i] {
				t.Log("Test", j, "Unexpected result at", i, "Got:", v.dst[i], "Expected:", v.ex[i])
				t.Fail()
			}
		}
	}
}
