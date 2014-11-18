// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas/goblas"
	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

func init() {
	mat64.Register(goblas.Blas{})
}

func TestCovarianceMatrix(t *testing.T) {
	blasEngine := mat64.Registered()
	for i, test := range []struct {
		mat  mat64.Matrix
		r, c int
		x    []float64
	}{
		{
			mat: mat64.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			r: 2,
			c: 2,
			x: []float64{
				2.5, 3,
				3, 10,
			},
		},
	} {
		// tests with a blas engine
		mat64.Register(goblas.Blas{})
		c := CovarianceMatrix(test.mat).RawMatrix()
		if c.Rows != test.r {
			t.Errorf("BLAS %d: expected rows %d, found %d", i, test.r, c.Rows)
		}
		if c.Cols != test.c {
			t.Errorf("BLAS %d: expected cols %d, found %d", i, test.c, c.Cols)
		}
		if !floats.Equal(test.x, c.Data) {
			t.Errorf("BLAS %d: expected data %#q, found %#q", i, test.x, c.Data)
		}
		// tests without a blas engine
		mat64.Register(nil)
		c = CovarianceMatrix(test.mat).RawMatrix()
		if c.Rows != test.r {
			t.Errorf("No BLAS %d: expected rows %d, found %d", i, test.r, c.Rows)
		}
		if c.Cols != test.c {
			t.Errorf("No BLAS %d: expected cols %d, found %d", i, test.c, c.Cols)
		}
		if !floats.Equal(test.x, c.Data) {
			t.Errorf("No BLAS %d: expected data %#q, found %#q", i, test.x, c.Data)
		}
	}
	mat64.Register(blasEngine)
}

// benchmarks

func randMat(r, c int) mat64.Matrix {
	x := make([]float64, r*c)
	for i := range x {
		x[i] = rand.Float64()
	}
	return mat64.NewDense(r, c, x)
}

func benchmarkCovarianceMatrix(b *testing.B, m mat64.Matrix) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CovarianceMatrix(m)
	}
}

func BenchmarkCovarianceMatrixSmallxSmall(b *testing.B) {
	// 10 * 10 elements
	x := randMat(SMALL, SMALL)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixSmallxMedium(b *testing.B) {
	// 10 * 1000 elements
	x := randMat(SMALL, MEDIUM)
	benchmarkCovarianceMatrix(b, x)
}

/*func BenchmarkCovarianceMatrixSmallxLarge(b *testing.B) {
	x := randMat(SMALL, LARGE)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixSmallxHuge(b *testing.B) {
	x := randMat(SMALL, HUGE)
	benchmarkCovarianceMatrix(b, x)
}*/

func BenchmarkCovarianceMatrixMediumxSmall(b *testing.B) {
	// 1000 * 10 elements
	x := randMat(MEDIUM, SMALL)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixMediumxMedium(b *testing.B) {
	// 1000 * 1000 elements
	x := randMat(MEDIUM, MEDIUM)
	benchmarkCovarianceMatrix(b, x)
}

/*func BenchmarkCovarianceMatrixMediumxLarge(b *testing.B) {
	x := randMat(MEDIUM, LARGE)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixMediumxHuge(b *testing.B) {
	x := randMat(MEDIUM, HUGE)
	benchmarkCovarianceMatrix(b, x)
}*/

func BenchmarkCovarianceMatrixLargexSmall(b *testing.B) {
	// 1e5 * 10 elements
	x := randMat(LARGE, SMALL)
	benchmarkCovarianceMatrix(b, x)
}

/*func BenchmarkCovarianceMatrixLargexMedium(b *testing.B) {
	// 1e5 * 1000 elements
    x := randMat(LARGE, MEDIUM)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixLargexLarge(b *testing.B) {
	x := randMat(LARGE, LARGE)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixLargexHuge(b *testing.B) {
	x := randMat(LARGE, HUGE)
	benchmarkCovarianceMatrix(b, x)
}*/

func BenchmarkCovarianceMatrixHugexSmall(b *testing.B) {
	// 1e7 * 10 elements
	x := randMat(HUGE, SMALL)
	benchmarkCovarianceMatrix(b, x)
}

/*func BenchmarkCovarianceMatrixHugexMedium(b *testing.B) {
	// 1e7 * 1000 elements
    x := randMat(HUGE, MEDIUM)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixHugexLarge(b *testing.B) {
	x := randMat(HUGE, LARGE)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixHugexHuge(b *testing.B) {
	x := randMat(HUGE, HUGE)
	benchmarkCovarianceMatrix(b, x)
}*/
