// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas/goblas"
	"github.com/gonum/matrix/mat64"
)

func init() {
	mat64.Register(goblas.Blas{})
}

func TestCovarianceMatrix(t *testing.T) {

	// An alternate way to test this is to call the Variance
	// and Covariance functions and ensure that the results are identical.
	for i, test := range []struct {
		data    mat64.Matrix
		weights mat64.Vec
		ans     mat64.Matrix
	}{
		{
			data: mat64.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			weights: nil,
			ans: mat64.NewDense(2, 2, []float64{
				2.5, 3,
				3, 10,
			}),
		}, {
			data: mat64.NewDense(3, 2, []float64{
				1, 1,
				2, 4,
				3, 9,
			}),
			weights: []float64{
				1,
				1.5,
				1,
			},
			ans: mat64.NewDense(2, 2, []float64{
				.8, 3.2,
				3.2, 13.142857142857146,
			}),
		},
	} {
		c := CovarianceMatrix(nil, test.data, test.weights)
		if !c.Equals(test.ans) {
			t.Errorf("%d: expected cov %v, found %v", i, test.ans, c)
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
	x := randMat(small, small)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixSmallxMedium(b *testing.B) {
	// 10 * 1000 elements
	x := randMat(small, medium)
	benchmarkCovarianceMatrix(b, x)
}

func BenchmarkCovarianceMatrixMediumxSmall(b *testing.B) {
	// 1000 * 10 elements
	x := randMat(medium, small)
	benchmarkCovarianceMatrix(b, x)
}
func BenchmarkCovarianceMatrixMediumxMedium(b *testing.B) {
	// 1000 * 1000 elements
	x := randMat(medium, medium)
	benchmarkCovarianceMatrix(b, x)
}

func BenchmarkCovarianceMatrixLargexSmall(b *testing.B) {
	// 1e5 * 10 elements
	x := randMat(large, small)
	benchmarkCovarianceMatrix(b, x)
}

func BenchmarkCovarianceMatrixHugexSmall(b *testing.B) {
	// 1e7 * 10 elements
	x := randMat(huge, small)
	benchmarkCovarianceMatrix(b, x)
}

func BenchmarkCovarianceMatrixSmallxSmallInPlace(b *testing.B) {
	// 10 * 10 elements
	x := randMat(small, small)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixSmallxMediumInPlace(b *testing.B) {
	// 10 * 1000 elements
	x := randMat(small, medium)
	benchmarkCovarianceMatrixInPlace(b, x)
}

func BenchmarkCovarianceMatrixMediumxSmallInPlace(b *testing.B) {
	// 1000 * 10 elements
	x := randMat(medium, small)
	benchmarkCovarianceMatrixInPlace(b, x)
}
func BenchmarkCovarianceMatrixMediumxMediumInPlace(b *testing.B) {
	// 1000 * 1000 elements
	x := randMat(medium, medium)
	benchmarkCovarianceMatrixInPlace(b, x)
}

func BenchmarkCovarianceMatrixLargexSmallInPlace(b *testing.B) {
	// 1e5 * 10 elements
	x := randMat(large, small)
	benchmarkCovarianceMatrixInPlace(b, x)
}

func BenchmarkCovarianceMatrixHugexSmallInPlace(b *testing.B) {
	// 1e7 * 10 elements
	x := randMat(huge, small)
	benchmarkCovarianceMatrixInPlace(b, x)
}
