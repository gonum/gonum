// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"github.com/gonum/matrix/mat64"
)

// CovarianceMatrix calculates a covariance matrix (also known as a
// variance-covariance matrix) from a matrix of data, using a two-pass
// algorithm.  It requires a registered BLAS engine in gonum/matrix/mat64.
//
// The matrix returned will be symmetric, square, and positive-semidefinite.
//
// The weights matrix wts should have the same number of rows as the input
// data matrix x.  cov should be a square matrix with the same number of
// columns as the input data matrix x, or if it is nil then a new Dense
// matrix will be constructed.
func CovarianceMatrix(cov *mat64.Dense, x, wts mat64.Matrix) *mat64.Dense {

	// matrix version of the two pass algorithm.  This doesn't use
	// the correction found in the Covariance and Variance functions.

	r, c := x.Dims()

	// determine the mean of each of the columns
	ones := make([]float64, r)
	for i := range ones {
		ones[i] = 1
	}
	b := mat64.NewDense(1, r, ones)
	b.Mul(b, x)
	b.Scale(1/float64(r), b)
	mu := b.RowView(0)

	// subtract the mean from the data
	xc := mat64.DenseCopyOf(x)

	for i := 0; i < r; i++ {
		rv := xc.RowView(i)
		for j, mean := range mu {
			rv[j] -= mean
		}
	}
	var xt mat64.Dense
	xt.TCopy(xc)

	// normalization factor, typical n-1
	var N float64
	if wts != nil {
		// should this be col major or row major?
		if wr, wc := wts.Dims(); wr != r || wc != 1 {
			panic("matrix length mismatch")
		}

		for i := 0; i < r; i++ {
			rv := xc.RowView(i)
			w := wts.At(i, 0)
			N += w
			for j := 0; j < c; j++ {
				rv[j] *= w
			}
		}
		N = 1 / (N - 1)
	} else {
		N = 1 / float64(r-1)
	}

	// TODO: indicate that the resulting matrix is symmetric, which
	// should improve performance.
	if cov == nil {
		cov = mat64.NewDense(c, c, nil)
	} else if covr, covc := cov.Dims(); covr != covc || covc != c {
		panic("matrix size mismatch")
	}

	cov.Mul(&xt, xc)
	cov.Scale(N, cov)
	return cov
}
