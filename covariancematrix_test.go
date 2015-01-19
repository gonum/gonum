// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

func TestCovarianceMatrix(t *testing.T) {
	// An alternate way to test this is to call the Variance
	// and Covariance functions and ensure that the results are identical.
	for i, test := range []struct {
		data    *mat64.Dense
		weights []float64
		ans     *mat64.Dense
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
		// Make a copy of the data to check that it isn't changing.
		r := test.data.RawMatrix()
		d := make([]float64, len(r.Data))
		copy(d, r.Data)

		w := make([]float64, len(test.weights))
		if test.weights != nil {
			copy(w, test.weights)
		}
		c := CovarianceMatrix(nil, test.data, test.weights)
		if !c.Equals(test.ans) {
			t.Errorf("%d: expected cov %v, found %v", i, test.ans, c)
		}
		if !floats.Equal(d, r.Data) {
			t.Errorf("%d: data was modified during execution", i)
		}
		if !floats.Equal(w, test.weights) {
			t.Errorf("%d: weights was modified during execution", i)
		}

		// compare with call to Covariance
		_, cols := c.Dims()
		for ci := 0; ci < cols; ci++ {
			for cj := 0; cj < cols; cj++ {
				x := test.data.Col(nil, ci)
				y := test.data.Col(nil, cj)
				cov := Covariance(x, y, test.weights)
				if math.Abs(cov-c.At(ci, cj)) > 1e-14 {
					t.Errorf("CovMat does not match at (%v, %v). Want %v, got %v.", ci, cj, cov, c.At(ci, cj))
				}
			}
		}

	}
	if !Panics(func() { CovarianceMatrix(nil, mat64.NewDense(5, 2, nil), []float64{}) }) {
		t.Errorf("CovarianceMatrix did not panic with weight size mismatch")
	}
	if !Panics(func() { CovarianceMatrix(mat64.NewDense(1, 1, nil), mat64.NewDense(5, 2, nil), nil) }) {
		t.Errorf("CovarianceMatrix did not panic with preallocation size mismatch")
	}
	if !Panics(func() { CovarianceMatrix(nil, mat64.NewDense(2, 2, []float64{1, 2, 3, 4}), []float64{1, -1}) }) {
		t.Errorf("CovarianceMatrix did not panic with negative weights")
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
func benchmarkCovarianceMatrixWeighted(b *testing.B, m mat64.Matrix) {
	r, _ := m.Dims()
	wts := make([]float64, r)
	for i := range wts {
		wts[i] = 0.5
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CovarianceMatrix(nil, m, wts)
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

func BenchmarkCovarianceMatrixSmallxSmallWeighted(b *testing.B) {
	// 10 * 10 elements
	x := randMat(small, small)
	benchmarkCovarianceMatrixWeighted(b, x)
}
func BenchmarkCovarianceMatrixSmallxMediumWeighted(b *testing.B) {
	// 10 * 1000 elements
	x := randMat(small, medium)
	benchmarkCovarianceMatrixWeighted(b, x)
}

func BenchmarkCovarianceMatrixMediumxSmallWeighted(b *testing.B) {
	// 1000 * 10 elements
	x := randMat(medium, small)
	benchmarkCovarianceMatrixWeighted(b, x)
}
func BenchmarkCovarianceMatrixMediumxMediumWeighted(b *testing.B) {
	// 1000 * 1000 elements
	x := randMat(medium, medium)
	benchmarkCovarianceMatrixWeighted(b, x)
}

func BenchmarkCovarianceMatrixLargexSmallWeighted(b *testing.B) {
	// 1e5 * 10 elements
	x := randMat(large, small)
	benchmarkCovarianceMatrixWeighted(b, x)
}

func BenchmarkCovarianceMatrixHugexSmallWeighted(b *testing.B) {
	// 1e7 * 10 elements
	x := randMat(huge, small)
	benchmarkCovarianceMatrixWeighted(b, x)
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
