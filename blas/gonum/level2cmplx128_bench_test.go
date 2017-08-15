// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/blas"
)

func benchmarkZher(b *testing.B, uplo blas.Uplo, n, inc int) {
	rnd := rand.New(rand.NewSource(1))
	alpha := rnd.NormFloat64()
	x := make([]complex128, (n-1)*inc+1)
	for i := range x {
		x[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}
	a := make([]complex128, n*n)
	for i := range a {
		a[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		impl.Zher(uplo, n, alpha, x, inc, a, n)
	}
}

func BenchmarkZherUpperN10Inc1(b *testing.B)    { benchmarkZher(b, blas.Upper, 10, 1) }
func BenchmarkZherUpperN100Inc1(b *testing.B)   { benchmarkZher(b, blas.Upper, 100, 1) }
func BenchmarkZherUpperN1000Inc1(b *testing.B)  { benchmarkZher(b, blas.Upper, 1000, 1) }
func BenchmarkZherUpperN10000Inc1(b *testing.B) { benchmarkZher(b, blas.Upper, 10000, 1) }

func BenchmarkZherUpperN10Inc10(b *testing.B)    { benchmarkZher(b, blas.Upper, 10, 10) }
func BenchmarkZherUpperN100Inc10(b *testing.B)   { benchmarkZher(b, blas.Upper, 100, 10) }
func BenchmarkZherUpperN1000Inc10(b *testing.B)  { benchmarkZher(b, blas.Upper, 1000, 10) }
func BenchmarkZherUpperN10000Inc10(b *testing.B) { benchmarkZher(b, blas.Upper, 10000, 10) }

func BenchmarkZherUpperN10Inc1000(b *testing.B)    { benchmarkZher(b, blas.Upper, 10, 1000) }
func BenchmarkZherUpperN100Inc1000(b *testing.B)   { benchmarkZher(b, blas.Upper, 100, 1000) }
func BenchmarkZherUpperN1000Inc1000(b *testing.B)  { benchmarkZher(b, blas.Upper, 1000, 1000) }
func BenchmarkZherUpperN10000Inc1000(b *testing.B) { benchmarkZher(b, blas.Upper, 10000, 1000) }

func BenchmarkZherLowerN10Inc1(b *testing.B)    { benchmarkZher(b, blas.Lower, 10, 1) }
func BenchmarkZherLowerN100Inc1(b *testing.B)   { benchmarkZher(b, blas.Lower, 100, 1) }
func BenchmarkZherLowerN1000Inc1(b *testing.B)  { benchmarkZher(b, blas.Lower, 1000, 1) }
func BenchmarkZherLowerN10000Inc1(b *testing.B) { benchmarkZher(b, blas.Lower, 10000, 1) }

func BenchmarkZherLowerN10Inc10(b *testing.B)    { benchmarkZher(b, blas.Lower, 10, 10) }
func BenchmarkZherLowerN100Inc10(b *testing.B)   { benchmarkZher(b, blas.Lower, 100, 10) }
func BenchmarkZherLowerN1000Inc10(b *testing.B)  { benchmarkZher(b, blas.Lower, 1000, 10) }
func BenchmarkZherLowerN10000Inc10(b *testing.B) { benchmarkZher(b, blas.Lower, 10000, 10) }

func BenchmarkZherLowerN10Inc1000(b *testing.B)    { benchmarkZher(b, blas.Lower, 10, 1000) }
func BenchmarkZherLowerN100Inc1000(b *testing.B)   { benchmarkZher(b, blas.Lower, 100, 1000) }
func BenchmarkZherLowerN1000Inc1000(b *testing.B)  { benchmarkZher(b, blas.Lower, 1000, 1000) }
func BenchmarkZherLowerN10000Inc1000(b *testing.B) { benchmarkZher(b, blas.Lower, 10000, 1000) }
