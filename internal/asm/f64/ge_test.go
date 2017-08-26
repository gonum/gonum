// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"fmt"
	"testing"

	"github.com/gonum/blas/testblas"
)

func TestGer(t *testing.T) {
	tests := []struct {
		m, n            uintptr
		alpha           float64
		x, y, a         []float64
		incX, incY, lda uintptr
		want            []float64
	}{
		{
			m: 1, n: 1, alpha: 1,
			x: []float64{2}, incX: 1,
			y: []float64{4.4}, incY: 1,
			a: []float64{10}, lda: 1,
			want: []float64{18.8},
		},
	}

	for _, test := range tests {
		Ger(test.m, test.n, test.alpha, test.x, test.incX, test.y, test.incY, test.a, test.lda)
	}
}

type dgerWrap struct{}

func (d dgerWrap) Dger(m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
	Ger(uintptr(m), uintptr(n), alpha, x, uintptr(incX), y, uintptr(incY), a, uintptr(lda))
}

func TestBlasGer(t *testing.T) {
	testblas.DgerTest(t, dgerWrap{})
}

func BenchmarkBlasGer(t *testing.B) {
	for _, dims := range newIncSet(3, 10, 30, 100, 300, 1000, 1e4, 1e5) {
		m, n := dims.x, dims.y
		if m/n >= 100 || n/m >= 100 || (m == 1e5 && n == 1e5) {
			continue
		}
		for _, inc := range newIncSet(1, 2, 3, 4, 10) {
			incX, incY := inc.x, inc.y
			t.Run(fmt.Sprintf("Dger %dx%d (%d %d)", m, n, incX, incY), func(b *testing.B) {
				for i := 0; i < t.N; i++ {
					testblas.DgerBenchmark(b, dgerWrap{}, m, n, incX, incY)
				}
			})

		}
	}
}
