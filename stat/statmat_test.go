// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestCovarianceMatrix(t *testing.T) {
	// An alternative way to test this is to call the Variance and
	// Covariance functions and ensure that the results are identical.
	for i, test := range []struct {
		data    *mat.Dense
		weights []float64
		ans     *mat.Dense
	}{
		{
			data: mat.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			weights: nil,
			ans: mat.NewDense(2, 2, []float64{
				2.5, 3,
				3, 10,
			}),
		}, {
			data: mat.NewDense(3, 2, []float64{
				1, 1,
				2, 4,
				3, 9,
			}),
			weights: []float64{
				1,
				1.5,
				1,
			},
			ans: mat.NewDense(2, 2, []float64{
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

		var cov mat.SymDense
		CovarianceMatrix(&cov, test.data, test.weights)
		if !mat.Equal(&cov, test.ans) {
			t.Errorf("%d: expected cov %v, found %v", i, test.ans, &cov)
		}
		if !floats.Equal(d, r.Data) {
			t.Errorf("%d: data was modified during execution", i)
		}
		if !floats.Equal(w, test.weights) {
			t.Errorf("%d: weights was modified during execution", i)
		}

		// compare with call to Covariance
		_, cols := cov.Dims()
		for ci := 0; ci < cols; ci++ {
			for cj := 0; cj < cols; cj++ {
				x := mat.Col(nil, ci, test.data)
				y := mat.Col(nil, cj, test.data)
				wc := Covariance(x, y, test.weights)
				if math.Abs(wc-cov.At(ci, cj)) > 1e-14 {
					t.Errorf("CovMat does not match at (%v, %v). Want %v, got %v.", ci, cj, wc, cov.At(ci, cj))
				}
			}
		}

	}
	if !panics(func() { CovarianceMatrix(nil, mat.NewDense(5, 2, nil), []float64{}) }) {
		t.Errorf("CovarianceMatrix did not panic with weight size mismatch")
	}
	if !panics(func() { CovarianceMatrix(mat.NewSymDense(1, nil), mat.NewDense(5, 2, nil), nil) }) {
		t.Errorf("CovarianceMatrix did not panic with preallocation size mismatch")
	}
	if !panics(func() { CovarianceMatrix(nil, mat.NewDense(2, 2, []float64{1, 2, 3, 4}), []float64{1, -1}) }) {
		t.Errorf("CovarianceMatrix did not panic with negative weights")
	}
}

func TestCorrelationMatrix(t *testing.T) {
	for i, test := range []struct {
		data    *mat.Dense
		weights []float64
		ans     *mat.Dense
	}{
		{
			data: mat.NewDense(3, 3, []float64{
				1, 2, 3,
				3, 4, 5,
				5, 6, 7,
			}),
			weights: nil,
			ans: mat.NewDense(3, 3, []float64{
				1, 1, 1,
				1, 1, 1,
				1, 1, 1,
			}),
		},
		{
			data: mat.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			weights: nil,
			ans: mat.NewDense(2, 2, []float64{
				1, 0.6,
				0.6, 1,
			}),
		}, {
			data: mat.NewDense(3, 2, []float64{
				1, 1,
				2, 4,
				3, 9,
			}),
			weights: []float64{
				1,
				1.5,
				1,
			},
			ans: mat.NewDense(2, 2, []float64{
				1, 0.9868703275903379,
				0.9868703275903379, 1,
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

		var corr mat.SymDense
		CorrelationMatrix(&corr, test.data, test.weights)
		if !mat.Equal(&corr, test.ans) {
			t.Errorf("%d: expected corr %v, found %v", i, test.ans, &corr)
		}
		if !floats.Equal(d, r.Data) {
			t.Errorf("%d: data was modified during execution", i)
		}
		if !floats.Equal(w, test.weights) {
			t.Errorf("%d: weights was modified during execution", i)
		}

		// compare with call to Covariance
		_, cols := corr.Dims()
		for ci := 0; ci < cols; ci++ {
			for cj := 0; cj < cols; cj++ {
				x := mat.Col(nil, ci, test.data)
				y := mat.Col(nil, cj, test.data)
				wc := Correlation(x, y, test.weights)
				if math.Abs(wc-corr.At(ci, cj)) > 1e-14 {
					t.Errorf("CorrMat does not match at (%v, %v). Want %v, got %v.", ci, cj, wc, corr.At(ci, cj))
				}
			}
		}

	}
	if !panics(func() { CorrelationMatrix(nil, mat.NewDense(5, 2, nil), []float64{}) }) {
		t.Errorf("CorrelationMatrix did not panic with weight size mismatch")
	}
	if !panics(func() { CorrelationMatrix(mat.NewSymDense(1, nil), mat.NewDense(5, 2, nil), nil) }) {
		t.Errorf("CorrelationMatrix did not panic with preallocation size mismatch")
	}
	if !panics(func() { CorrelationMatrix(nil, mat.NewDense(2, 2, []float64{1, 2, 3, 4}), []float64{1, -1}) }) {
		t.Errorf("CorrelationMatrix did not panic with negative weights")
	}
}

func TestCorrCov(t *testing.T) {
	// test both Cov2Corr and Cov2Corr
	for i, test := range []struct {
		data    *mat.Dense
		weights []float64
	}{
		{
			data: mat.NewDense(3, 3, []float64{
				1, 2, 3,
				3, 4, 5,
				5, 6, 7,
			}),
			weights: nil,
		},
		{
			data: mat.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			weights: nil,
		}, {
			data: mat.NewDense(3, 2, []float64{
				1, 1,
				2, 4,
				3, 9,
			}),
			weights: []float64{
				1,
				1.5,
				1,
			},
		},
	} {
		var corr, cov mat.SymDense
		CorrelationMatrix(&corr, test.data, test.weights)
		CovarianceMatrix(&cov, test.data, test.weights)

		r := cov.Symmetric()

		// Get the diagonal elements from cov to determine the sigmas.
		sigmas := make([]float64, r)
		for i := range sigmas {
			sigmas[i] = math.Sqrt(cov.At(i, i))
		}

		covFromCorr := mat.NewSymDense(corr.Symmetric(), nil)
		covFromCorr.CopySym(&corr)
		corrToCov(covFromCorr, sigmas)

		corrFromCov := mat.NewSymDense(cov.Symmetric(), nil)
		corrFromCov.CopySym(&cov)
		covToCorr(corrFromCov)

		if !mat.EqualApprox(&corr, corrFromCov, 1e-14) {
			t.Errorf("%d: corrToCov did not match direct Correlation calculation.  Want: %v, got: %v. ", i, corr, corrFromCov)
		}
		if !mat.EqualApprox(&cov, covFromCorr, 1e-14) {
			t.Errorf("%d: covToCorr did not match direct Covariance calculation.  Want: %v, got: %v. ", i, cov, covFromCorr)
		}

		if !panics(func() { corrToCov(mat.NewSymDense(2, nil), []float64{}) }) {
			t.Errorf("CorrelationMatrix did not panic with sigma size mismatch")
		}
	}
}

func TestMahalanobis(t *testing.T) {
	// Comparison with scipy.
	for cas, test := range []struct {
		x, y  *mat.VecDense
		Sigma *mat.SymDense
		ans   float64
	}{
		{
			x: mat.NewVecDense(3, []float64{1, 2, 3}),
			y: mat.NewVecDense(3, []float64{0.8, 1.1, -1}),
			Sigma: mat.NewSymDense(3,
				[]float64{
					0.8, 0.3, 0.1,
					0.3, 0.7, -0.1,
					0.1, -0.1, 7}),
			ans: 1.9251757377680914,
		},
	} {
		var chol mat.Cholesky
		ok := chol.Factorize(test.Sigma)
		if !ok {
			panic("bad test")
		}
		ans := Mahalanobis(test.x, test.y, &chol)
		if math.Abs(ans-test.ans) > 1e-14 {
			t.Errorf("Cas %d: got %v, want %v", cas, ans, test.ans)
		}
	}
}

// benchmarks

func randMat(r, c int) mat.Matrix {
	x := make([]float64, r*c)
	for i := range x {
		x[i] = rand.Float64()
	}
	return mat.NewDense(r, c, x)
}

func benchmarkCovarianceMatrix(b *testing.B, m mat.Matrix) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CovarianceMatrix(nil, m, nil)
	}
}
func benchmarkCovarianceMatrixWeighted(b *testing.B, m mat.Matrix) {
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
func benchmarkCovarianceMatrixInPlace(b *testing.B, m mat.Matrix) {
	_, c := m.Dims()
	res := mat.NewSymDense(c, nil)
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

func BenchmarkCovToCorr(b *testing.B) {
	// generate a 10x10 covariance matrix
	m := randMat(small, small)
	var cov mat.SymDense
	CovarianceMatrix(&cov, m, nil)
	cc := mat.NewSymDense(cov.Symmetric(), nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cc.CopySym(&cov)
		b.StartTimer()
		covToCorr(cc)
	}
}

func BenchmarkCorrToCov(b *testing.B) {
	// generate a 10x10 correlation matrix
	m := randMat(small, small)
	var corr mat.SymDense
	CorrelationMatrix(&corr, m, nil)
	cc := mat.NewSymDense(corr.Symmetric(), nil)
	sigma := make([]float64, small)
	for i := range sigma {
		sigma[i] = 2
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cc.CopySym(&corr)
		b.StartTimer()
		corrToCov(cc, sigma)
	}
}
