// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

// CovarianceMatrix calculates a covariance matrix (also known as a
// variance-covariance matrix) from a matrix of data, using a two-pass
// algorithm. The matrix returned will be symmetric, square, and
// positive-semidefinite.
//
// The weights wts should have the length equal to the number of rows in
// input data matrix x. cov should either be a square matrix with the same
// number of columns as the input data matrix x, or nil in which case a new
// Dense matrix will be constructed.
func CovarianceMatrix(cov *mat64.Dense, x mat64.Matrix, wts []float64) *mat64.Dense {
	// This is the matrix version of the two-pass algorithm. It doesn't use the
	// additional floating point error correction that the Covariance function uses
	// to reduce the impact of rounding during centering.

	r, c := x.Dims()

	// TODO(jonlawlor): indicate that the resulting matrix is symmetric, which
	// should improve performance.
	if cov == nil {
		cov = mat64.NewDense(c, c, nil)
	} else if covr, covc := cov.Dims(); covr != covc || covc != c {
		panic(mat64.ErrShape)
	}

	var xt mat64.Dense
	xt.TCopy(x)
	// Subtract the mean of each of the columns.
	for i := 0; i < c; i++ {
		v := xt.RowView(i)
		mean := Mean(v, wts)
		floats.AddConst(-mean, v)
	}

	var xc mat64.Dense
	xc.TCopy(&xt)

	var n, scale float64
	if wts != nil {
		if wr := len(wts); wr != r {
			panic(mat64.ErrShape)
		}

		// Weight the rows.
		for i := 0; i < c; i++ {
			v := xt.RowView(i)
			floats.Mul(v, wts)
		}

		// Calculate the normalization factor.
		n = floats.Sum(wts)
	} else {
		n = float64(r)
	}

	cov.Mul(&xt, &xc)

	// Scale by the sample size.
	scale = 1 / (n - 1)
	cov.Scale(scale, cov)

	return cov
}
