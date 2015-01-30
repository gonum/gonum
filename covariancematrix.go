// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
	"math"
)

// CovarianceMatrix calculates a covariance matrix (also known as a
// variance-covariance matrix) from a matrix of data, using a two-pass
// algorithm. The matrix returned will be symmetric, square, and
// positive-semidefinite.
//
// The weights wts should have the length equal to the number of rows in
// input data matrix x. cov should either be a square matrix with the same
// number of columns as the input data matrix x, or nil in which case a new
// Dense matrix will be constructed.  Weights cannot be negative.
func CovarianceMatrix(cov *mat64.Dense, x mat64.Matrix, wts []float64) *mat64.Dense {
	// This is the matrix version of the two-pass algorithm. It doesn't use the
	// additional floating point error correction that the Covariance function uses
	// to reduce the impact of rounding during centering.

	// TODO(jonlawlor): indicate that the resulting matrix is symmetric, and change
	// the returned type from a *mat.Dense to a *mat.Symmetric.

	r, c := x.Dims()

	if cov == nil {
		cov = mat64.NewDense(c, c, nil)
	} else if covr, covc := cov.Dims(); covr != covc || covc != c {
		panic(mat64.ErrShape)
	}

	var xt mat64.Dense
	xt.TCopy(x)
	// Subtract the mean of each of the columns.
	for i := 0; i < c; i++ {
		v := xt.RawRowView(i)
		// This will panic with ErrShape if len(wts) != len(v), so
		// we don't have to check the size later.
		mean := Mean(v, wts)
		floats.AddConst(-mean, v)
	}

	var n float64
	if wts == nil {

		n = float64(r)

		cov.MulTrans(&xt, false, &xt, true)

		// Scale by the sample size.
		cov.Scale(1/(n-1), cov)
		return cov
	}

	// Multiply by the sqrt of the weights, so that multiplication is symmetric.
	sqrtwts := make([]float64, r)
	for i, w := range wts {
		if w < 0 {
			panic("stat: negative covariance matrix weights")
		}
		sqrtwts[i] = math.Sqrt(w)
	}
	// Weight the rows.
	for i := 0; i < c; i++ {
		v := xt.RawRowView(i)
		floats.Mul(v, sqrtwts)
	}

	// Calculate the normalization factor.
	n = floats.Sum(wts)
	cov.MulTrans(&xt, false, &xt, true)

	// Scale by the sample size.
	cov.Scale(1/(n-1), cov)
	return cov
}
