// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// CovarianceMatrix returns the covariance matrix (also known as the
// variance-covariance matrix) calculated from a matrix of data, x, using
// a two-pass algorithm.
//
// If weights is not nil the weighted covariance of x is calculated. weights
// must have length equal to the number of rows in input data matrix and
// must not contain negative elements.
// If cov is not nil it must either be zero-sized or have the same number of
// columns as the input data matrix. cov will be used as the destination for
// the covariance data. If cov is nil, a new mat.SymDense is allocated for
// the destination.
func CovarianceMatrix(cov *mat.SymDense, x mat.Matrix, weights []float64) *mat.SymDense {
	// This is the matrix version of the two-pass algorithm. It doesn't use the
	// additional floating point error correction that the Covariance function uses
	// to reduce the impact of rounding during centering.

	r, c := x.Dims()

	if cov == nil {
		cov = mat.NewSymDense(c, nil)
	} else if n := cov.Symmetric(); n != c && n != 0 {
		panic(mat.ErrShape)
	}

	var xt mat.Dense
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

// CorrelationMatrix returns the correlation matrix calculated from a matrix
// of data, x, using a two-pass algorithm.
//
// If weights is not nil the weighted correlation of x is calculated. weights
// must have length equal to the number of rows in input data matrix and
// must not contain negative elements.
// If corr is not nil it must either be zero-sized or have the same number of
// columns as the input data matrix. corr will be used as the destination for
// the correlation data. If corr is nil, a new mat.SymDense is allocated for
// the destination.
func CorrelationMatrix(corr *mat.SymDense, x mat.Matrix, weights []float64) *mat.SymDense {
	// This will panic if the sizes don't match, or if weights is the wrong size.
	corr = CovarianceMatrix(corr, x, weights)
	covToCorr(corr)
	return corr
}

// covToCorr converts a covariance matrix to a correlation matrix.
func covToCorr(c *mat.SymDense) {
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
func corrToCov(c *mat.SymDense, sigma []float64) {
	r, _ := c.Dims()

	if r != len(sigma) {
		panic(mat.ErrShape)
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

// Mahalanobis computes the Mahalanobis distance
//  D = sqrt((x-y)^T * Σ^-1 * (x-y))
// between the column vectors x and y given the cholesky decomposition of Σ.
// Mahalanobis returns NaN if the linear solve fails.
//
// See https://en.wikipedia.org/wiki/Mahalanobis_distance for more information.
func Mahalanobis(x, y mat.Vector, chol *mat.Cholesky) float64 {
	var diff mat.VecDense
	diff.SubVec(x, y)
	var tmp mat.VecDense
	err := chol.SolveVec(&tmp, &diff)
	if err != nil {
		return math.NaN()
	}
	return math.Sqrt(mat.Dot(&tmp, &diff))
}
