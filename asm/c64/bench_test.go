// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c64

import (
	"runtime"
	"testing"
)

var (
	a = complex64(2 + 2i)
	x = make([]complex64, 50000)
	y = make([]complex64, 50000)
	z = make([]complex64, 50000)
)

func init() {
	tmp := complex64(1 + 1i)
	for n := range x {
		x[n] = complex(float32(n), float32(n)) * tmp
		y[n] = complex(float32(n), float32(n)) * tmp
	}
}

func benchaxpyu(t *testing.B, n int, f func(a complex64, x, y []complex64)) {
	x, y := x[:n], y[:n]

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		f(a, x, y)
	}
	t.StopTimer()
	runtime.GC()
}

func naiveaxpyu(a complex64, x, y []complex64) {
	for i, v := range x {
		y[i] += a * v
	}
}

func BenchmarkAxpyUnitary1(t *testing.B)     { benchaxpyu(t, 1, AxpyUnitary) }
func BenchmarkAxpyUnitary2(t *testing.B)     { benchaxpyu(t, 2, AxpyUnitary) }
func BenchmarkAxpyUnitary3(t *testing.B)     { benchaxpyu(t, 3, AxpyUnitary) }
func BenchmarkAxpyUnitary4(t *testing.B)     { benchaxpyu(t, 4, AxpyUnitary) }
func BenchmarkAxpyUnitary5(t *testing.B)     { benchaxpyu(t, 5, AxpyUnitary) }
func BenchmarkAxpyUnitary10(t *testing.B)    { benchaxpyu(t, 10, AxpyUnitary) }
func BenchmarkAxpyUnitary100(t *testing.B)   { benchaxpyu(t, 100, AxpyUnitary) }
func BenchmarkAxpyUnitary1000(t *testing.B)  { benchaxpyu(t, 1000, AxpyUnitary) }
func BenchmarkAxpyUnitary5000(t *testing.B)  { benchaxpyu(t, 5000, AxpyUnitary) }
func BenchmarkAxpyUnitary10000(t *testing.B) { benchaxpyu(t, 10000, AxpyUnitary) }
func BenchmarkAxpyUnitary50000(t *testing.B) { benchaxpyu(t, 50000, AxpyUnitary) }

func BenchmarkFAxpyUnitary1(t *testing.B)     { benchaxpyu(t, 1, naiveaxpyu) }
func BenchmarkFAxpyUnitary2(t *testing.B)     { benchaxpyu(t, 2, naiveaxpyu) }
func BenchmarkFAxpyUnitary3(t *testing.B)     { benchaxpyu(t, 3, naiveaxpyu) }
func BenchmarkFAxpyUnitary4(t *testing.B)     { benchaxpyu(t, 4, naiveaxpyu) }
func BenchmarkFAxpyUnitary5(t *testing.B)     { benchaxpyu(t, 5, naiveaxpyu) }
func BenchmarkFAxpyUnitary10(t *testing.B)    { benchaxpyu(t, 10, naiveaxpyu) }
func BenchmarkFAxpyUnitary100(t *testing.B)   { benchaxpyu(t, 100, naiveaxpyu) }
func BenchmarkFAxpyUnitary1000(t *testing.B)  { benchaxpyu(t, 1000, naiveaxpyu) }
func BenchmarkFAxpyUnitary5000(t *testing.B)  { benchaxpyu(t, 5000, naiveaxpyu) }
func BenchmarkFAxpyUnitary10000(t *testing.B) { benchaxpyu(t, 10000, naiveaxpyu) }
func BenchmarkFAxpyUnitary50000(t *testing.B) { benchaxpyu(t, 50000, naiveaxpyu) }

func benchaxpyut(t *testing.B, n int, f func(d []complex64, a complex64, x, y []complex64)) {
	x, y, z := x[:n], y[:n], z[:n]

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		f(z, a, x, y)
	}
	t.StopTimer()
	runtime.GC()
}

func naiveaxpyut(d []complex64, a complex64, x, y []complex64) {
	for i, v := range x {
		d[i] = y[i] + a*v
	}
}

func BenchmarkAxpyUnitaryTo1(t *testing.B)     { benchaxpyut(t, 1, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo2(t *testing.B)     { benchaxpyut(t, 2, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo3(t *testing.B)     { benchaxpyut(t, 3, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo4(t *testing.B)     { benchaxpyut(t, 4, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo5(t *testing.B)     { benchaxpyut(t, 5, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo10(t *testing.B)    { benchaxpyut(t, 10, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo100(t *testing.B)   { benchaxpyut(t, 100, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo1000(t *testing.B)  { benchaxpyut(t, 1000, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo5000(t *testing.B)  { benchaxpyut(t, 5000, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo10000(t *testing.B) { benchaxpyut(t, 10000, AxpyUnitaryTo) }
func BenchmarkAxpyUnitaryTo50000(t *testing.B) { benchaxpyut(t, 50000, AxpyUnitaryTo) }

func BenchmarkFAxpyUnitaryTo1(t *testing.B)     { benchaxpyut(t, 1, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo2(t *testing.B)     { benchaxpyut(t, 2, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo3(t *testing.B)     { benchaxpyut(t, 3, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo4(t *testing.B)     { benchaxpyut(t, 4, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo5(t *testing.B)     { benchaxpyut(t, 5, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo10(t *testing.B)    { benchaxpyut(t, 10, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo100(t *testing.B)   { benchaxpyut(t, 100, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo1000(t *testing.B)  { benchaxpyut(t, 1000, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo5000(t *testing.B)  { benchaxpyut(t, 5000, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo10000(t *testing.B) { benchaxpyut(t, 10000, naiveaxpyut) }
func BenchmarkFAxpyUnitaryTo50000(t *testing.B) { benchaxpyut(t, 50000, naiveaxpyut) }
