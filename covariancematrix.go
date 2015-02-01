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
// algorithm. The matrix returned will be symmetric and square.
//
// The weights wts should have the length equal to the number of rows in
// input data matrix x. If c is nil, then a new matrix with appropriate size will
// be constructed.  If c is not nil, it should be a square matrix with the same
// number of columns as the input data matrix x, and it will be used as the receiver
// for the covariance data.  Weights cannot be negative.
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

// CorrelationMatrix calculates a correlation matrix from a matrix of data,
// using a two-pass algorithm. The matrix returned will be symmetric and square.
//
// The weights wts should have the length equal to the number of rows in
// input data matrix x. If c is nil, then a new matrix with appropriate size will
// be constructed.  If c is not nil, it should be a square matrix with the same
// number of columns as the input data matrix x, and it will be used as the receiver
// for the correlation data.  Weights cannot be negative.
func CorrelationMatrix(c *mat64.Dense, x mat64.Matrix, wts []float64) *mat64.Dense {

	// TODO(jonlawlor): indicate that the resulting matrix is symmetric, and change
	// the returned type from a *mat.Dense to a *mat.Symmetric.

	// This will panic if the sizes don't match, or if wts is the wrong size.
	c = CovarianceMatrix(c, x, wts)
	covToCorr(c)
	return c
}

// covToCorr converts a covariance matrix to a correlation matrix.
func covToCorr(c *mat64.Dense) {

	// TODO(jonlawlor): use a *mat64.Symmetric as input.

	r, _ := c.Dims()

	s := make([]float64, r)
	for i := 0; i < r; i++ {
		s[i] = 1 / math.Sqrt(c.At(i, i))
	}
	for i, sx := range s {
		row := c.RawRowView(i)
		for j, sy := range s {
			if i == j {
				// Ensure that the diagonal has exactly ones.
				row[j] = 1
				continue
			}
			row[j] *= sx
			row[j] *= sy
		}
	}
}

// corrToCov converts a correlation matrix to a covariance matrix.
// The input sigma should be vector of standard deviations corresponding
// to the covariance.  It will panic if len(sigma) is not equal to the
// number of rows in the correlation matrix.
func corrToCov(c *mat64.Dense, sigma []float64) {

	// TODO(jonlawlor): use a *mat64.Symmetric as input.

	r, _ := c.Dims()

	if r != len(sigma) {
		panic(mat64.ErrShape)
	}

	for i, sx := range sigma {
		row := c.RawRowView(i)
		for j, sy := range sigma {
			if i == j {
				// Ensure that the diagonal has exactly sigma squared.
				row[j] = sx * sx
				continue
			}
			row[j] *= sx
			row[j] *= sy
		}
	}
}
