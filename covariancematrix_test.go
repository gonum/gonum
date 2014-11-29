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
	for i, test := range []struct {
		mat     mat64.Matrix
		weights mat64.Vec
		r, c    int
		x       []float64
	}{
		{
			mat: mat64.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			weights: nil,
			r:       2,
			c:       2,
			x: []float64{
				2.5, 3,
				3, 10,
			},
		}, {
			mat: mat64.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			weights: mat64.Vec([]float64{
				1.5,
				.5,
				1.5,
				.5,
				1,
			}),
			r: 2,
			c: 2,
			x: []float64{
				2.75, 4.5,
				4.5, 11,
			},
		},
	} {
		c := CovarianceMatrix(nil, test.mat, test.weights).RawMatrix()
		if c.Rows != test.r {
			t.Errorf("%d: expected rows %d, found %d", i, test.r, c.Rows)
		}
		if c.Cols != test.c {
			t.Errorf("%d: expected cols %d, found %d", i, test.c, c.Cols)
		}
		if !floats.Equal(test.x, c.Data) {
			t.Errorf("%d: expected data %#q, found %#q", i, test.x, c.Data)
		}
	}
	if !Panics(func() { CovarianceMatrix(nil, mat64.NewDense(5, 2, nil), mat64.Vec([]float64{})) }) {
		t.Errorf("CovarianceMatrix did not panic with weight size mismatch")
	}
	if !Panics(func() { CovarianceMatrix(mat64.NewDense(1, 1, nil), mat64.NewDense(5, 2, nil), nil) }) {
		t.Errorf("CovarianceMatrix did not panic with preallocation size mismatch")
	}

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
		CovarianceMatrix(nil, m, nil)
	}
}
func benchmarkCovarianceMatrixInPlace(b *testing.B, m mat64.Matrix) {
	_, c := m.Dims()
	res := mat64.NewDense(c, c, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CovarianceMatrix(res, m, nil)
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

func BenchmarkCovarianceMatrixSmallxSmallInPlace(b *testing.B) {
	// 10 * 10 elements
	x := randMat(SMALL, SMALL)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixSmallxMediumInPlace(b *testing.B) {
	// 10 * 1000 elements
	x := randMat(SMALL, MEDIUM)
	benchmarkCovarianceMatrixInPlace(b, x)
}

/*func BenchmarkCovarianceMatrixSmallxLargeInPlace(b *testing.B) {
	x := randMat(SMALL, LARGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixSmallxHugeInPlace(b *testing.B) {
	x := randMat(SMALL, HUGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}*/

func BenchmarkCovarianceMatrixMediumxSmallInPlace(b *testing.B) {
	// 1000 * 10 elements
	x := randMat(MEDIUM, SMALL)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixMediumxMediumInPlace(b *testing.B) {
	// 1000 * 1000 elements
	x := randMat(MEDIUM, MEDIUM)
	benchmarkCovarianceMatrixInPlace(b, x)
}

/*func BenchmarkCovarianceMatrixMediumxLargeInPlace(b *testing.B) {
	x := randMat(MEDIUM, LARGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixMediumxHugeInPlace(b *testing.B) {
	x := randMat(MEDIUM, HUGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}*/

func BenchmarkCovarianceMatrixLargexSmallInPlace(b *testing.B) {
	// 1e5 * 10 elements
	x := randMat(LARGE, SMALL)
	benchmarkCovarianceMatrixInPlace(b, x)
}

/*func BenchmarkCovarianceMatrixLargexMediumInPlace(b *testing.B) {
	// 1e5 * 1000 elements
    x := randMat(LARGE, MEDIUM)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixLargexLargeInPlace(b *testing.B) {
	x := randMat(LARGE, LARGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixLargexHugeInPlace(b *testing.B) {
	x := randMat(LARGE, HUGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}*/

func BenchmarkCovarianceMatrixHugexSmallInPlace(b *testing.B) {
	// 1e7 * 10 elements
	x := randMat(HUGE, SMALL)
	benchmarkCovarianceMatrixInPlace(b, x)
}

/*func BenchmarkCovarianceMatrixHugexMediumInPlace(b *testing.B) {
	// 1e7 * 1000 elements
    x := randMat(HUGE, MEDIUM)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixHugexLargeInPlace(b *testing.B) {
	x := randMat(HUGE, LARGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixHugexHugeInPlace(b *testing.B) {
	x := randMat(HUGE, HUGE)
	benchmarkCovarianceMatrixInPlace(b, x)
}*/
