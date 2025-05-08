// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"errors"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

const (
	// ErrMsgSVDFailed is the error message for SVD factorization failure.
	ErrMsgSVDFailed = "transform: SVD factorization failed"
)

// Umeyama estimates the similarity transformation parameters between two matrices X and Y.
//
// This is an implementation of the algorithm presented in:
// "Least-Squares Estimation of Transformation Parameters Between Two Point Patterns"
// by Shinji Umeyama, IEEE Transactions on Pattern Analysis and Machine Intelligence,
// Vol. 13, No. 4, April 1991, which can be found here: https://doi.org/10.1109/34.88573
//
// The algorithm finds the optimal similarity transformation [c, R, t] ∈ Sim(m)
// that minimizes the mean squared error between two point patterns.
//
// The transformation relates the point sets as:
// Y ≈ c * R * X + t
//
// The dimensions of X and Y must be equal. The function will panic if they are not.
// The points require consistent indexing. This means that point i of X needs to correspond
// to point i of Y.
//
// In this implementation, rows represent points and columns represent dimensions.
//
// Umeyama returns the scale factor c, the rotation matrix r and the translation vector t.
// If a computation fails, Umeyama will return an error.
// varThreshold is used for detecting a degenerate input by comparing it with the variance of x. This is necessary
// because a variance equal or close to zero could cause numerical instability and/or division by zero.
// In case of variance <= threshold, Umeyama will return an error.
// The threshold should be >= 0. If a negative value is passed, the default threshold of 1e-10 will be used.
func Umeyama(x, y *mat.Dense, varThreshold float64) (c float64, r *mat.Dense, t *mat.VecDense, err error) {
	if varThreshold < 0 {
		varThreshold = 1e-10
	}

	rowsX, colsX := x.Dims()
	rowsY, colsY := y.Dims()

	// Check dimensions.
	if rowsX != rowsY || colsX != colsY {
		panic("transform: dimensions of x and y do not match")
	}

	n := rowsX // number of points
	m := colsX // number of dimensions

	// Calculate means.
	muX := mat.NewVecDense(m, nil)
	muY := mat.NewVecDense(m, nil)

	colX := make([]float64, n)
	colY := make([]float64, n)

	for j := 0; j < m; j++ {
		mat.Col(colX, j, x)
		mat.Col(colY, j, y)
		muX.SetVec(j, stat.Mean(colX, nil))
		muY.SetVec(j, stat.Mean(colY, nil))
	}

	// Center the matrices and calculate variance of x.
	var varX float64
	xc := mat.NewDense(n, m, nil)
	yc := mat.NewDense(n, m, nil)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			xc.Set(i, j, x.At(i, j)-muX.AtVec(j))
			yc.Set(i, j, y.At(i, j)-muY.AtVec(j))

			varX += float64(xc.At(i, j) * xc.At(i, j))
		}
	}
	varX /= float64(n)

	// Check for degenerate case. This prevents cases of division by zero and mathematical instability due to
	// very low variance.
	if varX <= varThreshold {
		return 0, nil, nil, mat.DegenerateInputError{
			Variance:  varX,
			Threshold: varThreshold,
		}
	}

	// Calculate covariance matrix.
	covXY := mat.NewDense(m, m, nil)
	covXY.Mul(yc.T(), xc)
	covXY.Scale(1/float64(n), covXY)

	// Singular Value Decomposition
	var svd mat.SVD
	if !svd.Factorize(covXY, mat.SVDFull) {
		return 0, nil, nil, errors.New(ErrMsgSVDFailed)
	}

	// Get U and V.
	var u, v mat.Dense
	svd.UTo(&u)
	svd.VTo(&v)

	// Create identity matrix.
	s := mat.NewDiagDense(m, nil)
	for i := 0; i < m; i++ {
		s.SetDiag(i, 1)
	}

	// Check determinants to ensure proper rotation matrix (not reflection).
	if mat.Det(&u)*mat.Det(&v) < 0 {
		s.SetDiag(m-1, -1)
	}

	// Calculate scale factor c.
	singularValues := svd.Values(nil)
	for i := 0; i < m; i++ {
		c += singularValues[i] * s.At(i, i)
	}
	c /= varX

	// Calculate rotation matrix R.
	r = mat.NewDense(m, m, nil)
	r.Product(&u, s, v.T())

	// Calculate translation vector t.
	t = mat.NewVecDense(m, nil)
	rMuX := mat.NewVecDense(m, nil)
	rMuX.MulVec(r, muX)

	t.CopyVec(muY)
	t.AddScaledVec(t, -c, rMuX)

	return c, r, t, nil
}
