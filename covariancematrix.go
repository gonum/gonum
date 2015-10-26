// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"

	"github.com/gonum/floats"
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

// CovarianceMatrix calculates a covariance matrix (also known as a
// variance-covariance matrix) from a matrix of data, using a two-pass
// algorithm.
//
// The weights must have length equal to the number of rows in
// input data matrix x. If cov is nil, then a new matrix with appropriate size will
// be constructed. If cov is not nil, it should have the same number of columns as the
// input data matrix x, and it will be used as the destination for the covariance
// data. Weights must not be negative.
func CovarianceMatrix(cov *mat64.SymDense, x mat64.Matrix, weights []float64) *mat64.SymDense {
	// This is the matrix version of the two-pass algorithm. It doesn't use the
	// additional floating point error correction that the Covariance function uses
	// to reduce the impact of rounding during centering.

	r, c := x.Dims()

	if cov == nil {
		cov = mat64.NewSymDense(c, nil)
	} else if n := cov.Symmetric(); n != c {
		panic(matrix.ErrShape)
	}

	var xt mat64.Dense
	xt.Clone(x.T())
	// Subtract the mean of each of the columns.
	for i := 0; i < c; i++ {
		v := xt.RawRowView(i)
		// This will panic with ErrShape if len(weights) != len(v), so
		// we don't have to check the size later.
		mean := Mean(v, weights)
		floats.AddConst(-mean, v)
	}

	if weights == nil {
		// Calculate the normalization factor
		// scaled by the sample size.
		cov.SymOuterK(1/(float64(r)-1), &xt)
		return cov
	}

	// Multiply by the sqrt of the weights, so that multiplication is symmetric.
	sqrtwts := make([]float64, r)
	for i, w := range weights {
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

	// Calculate the normalization factor
	// scaled by the weighted sample size.
	cov.SymOuterK(1/(floats.Sum(weights)-1), &xt)
	return cov
}

// CorrelationMatrix calculates a correlation matrix from a matrix of data
// using a two-pass algorithm.
//
// The weights must have length equal to the number of rows in
// input data matrix x. If corr is nil, then a new matrix with appropriate size will
// be constructed. If corr is not nil, it should have the same number of columns
// as the input data matrix x, and it will be used as the destination for the
// correlation data. Weights must not be negative.
func CorrelationMatrix(corr *mat64.SymDense, x mat64.Matrix, weights []float64) *mat64.SymDense {
	// This will panic if the sizes don't match, or if weights is the wrong size.
	corr = CovarianceMatrix(corr, x, weights)
	covToCorr(corr)
	return corr
}

// covToCorr converts a covariance matrix to a correlation matrix.
func covToCorr(c *mat64.SymDense) {
	r := c.Symmetric()

	s := make([]float64, r)
	for i := 0; i < r; i++ {
		s[i] = 1 / math.Sqrt(c.At(i, i))
	}
	for i, sx := range s {
		// Ensure that the diagonal has exactly ones.
		c.SetSym(i, i, 1)
		for j := i + 1; j < r; j++ {
			v := c.At(i, j)
			c.SetSym(i, j, v*sx*s[j])
		}
	}
}

// corrToCov converts a correlation matrix to a covariance matrix.
// The input sigma should be vector of standard deviations corresponding
// to the covariance.  It will panic if len(sigma) is not equal to the
// number of rows in the correlation matrix.
func corrToCov(c *mat64.SymDense, sigma []float64) {
	r, _ := c.Dims()

	if r != len(sigma) {
		panic(matrix.ErrShape)
	}
	for i, sx := range sigma {
		// Ensure that the diagonal has exactly sigma squared.
		c.SetSym(i, i, sx*sx)
		for j := i + 1; j < r; j++ {
			v := c.At(i, j)
			c.SetSym(i, j, v*sx*sigma[j])
		}
	}
}
