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
func CovarianceMatrix(x mat64.Matrix) *mat64.Dense {

	// matrix version of the two pass algorithm.  This doesn't use
	// the correction found in the Covariance and Variance functions.

	r, _ := x.Dims()

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

	// TODO: indicate that the resulting matrix is symmetric, which
	// should improve performance.
	var ss mat64.Dense
	ss.Mul(&xt, xc)
	ss.Scale(1/float64(r-1), &ss)
	return &ss
}
