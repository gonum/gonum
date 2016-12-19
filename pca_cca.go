// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"errors"
	"math"

	"github.com/gonum/floats"
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

// PrincipalComponents returns the principal component vectors and the variances
// of the principal component scores of the input data which is represented as
// an n×d matrix a where each row is an observation and each column is a
// variable.
//
// PrincipalComponents centers the variables but does not scale the variance.
//
// The slice weights is used to weight the observations. If weights is nil, each
// weight is considered to have a value of one, otherwise the length of weights
// must match the number of observations or PrincipalComponents will panic.
//
// On successful completion, the principal component vectors are returned in the
// columns of vecs as a d×min(n,d) matrix, and the column variances of the
// principal component scores, b * vecs where b is a with centered columns, are
// returned in vars as a min(n,d)-long slice in descending sort order.
//
// On failure, vecs and vars are returned nil and ok is returned false.
func PrincipalComponents(a mat64.Matrix, weights []float64) (vecs *mat64.Dense, vars []float64, ok bool) {
	n, d := a.Dims()
	if weights != nil && len(weights) != n {
		panic("stat: len(weights) != observations")
	}

	svd, ok := svdFactorizeCentered(a, weights)
	if !ok {
		return nil, nil, false
	}

	vecs = &mat64.Dense{}
	vecs.VFromSVD(svd)
	if n < d {
		// Don't retain columns that are not valid direction vectors.
		vecs.Clone(vecs.Slice(0, d, 0, n))
	}
	vars = svd.Values(nil)
	var f float64
	if weights == nil {
		f = 1 / float64(n-1)
	} else {
		f = 1 / (floats.Sum(weights) - 1)
	}
	for i, v := range vars {
		vars[i] = f * v * v
	}
	return vecs, vars, true
}

// CC is the result of a canonical correlation analysis. A CC should not be used
// unless it has been returned by a successful call to CanonicalCorrelations.
type CC struct {
	// n is the number of observations used to
	// construct the canonical correlations.
	n int

	x, y, c *mat64.SVD
}

// Corrs returns the canonical correlations in dest.
func (c CC) Corrs(dest []float64) []float64 {
	if c.c == nil {
		panic("stat: canonical correlations missing or invalid")
	}

	return c.c.Values(dest)
}

// Left returns the left eigenvectors of the canonical correlation matrix if
// spheredSpace is true. If spheredSpace is false it returns these eigenvectors
// back-transformed to the original data space.
func (c CC) Left(spheredSpace bool) (vecs *mat64.Dense) {
	if c.c == nil || c.x == nil || c.n < 2 {
		panic("stat: canonical correlations missing or invalid")
	}

	var pv mat64.Dense
	pv.UFromSVD(c.c)
	if spheredSpace {
		return &pv
	}

	var xv mat64.Dense
	xs := c.x.Values(nil)
	xv.VFromSVD(c.x)

	scaleColsReciSqrt(&xv, xs)

	pv.Product(&xv, xv.T(), &pv)
	pv.Scale(math.Sqrt(float64(c.n-1)), &pv)
	return &pv
}

// Right returns the right eigenvectors of the canonical correlation matrix if
// spheredSpace is true. If spheredSpace is false it returns these eigenvectors
// back-transformed to the original data space.
func (c CC) Right(spheredSpace bool) (vecs *mat64.Dense) {
	if c.c == nil || c.y == nil || c.n < 2 {
		panic("stat: canonical correlations missing or invalid")
	}

	var qv mat64.Dense
	qv.VFromSVD(c.c)
	if spheredSpace {
		return &qv
	}

	var yv mat64.Dense
	ys := c.y.Values(nil)
	yv.VFromSVD(c.y)

	scaleColsReciSqrt(&yv, ys)

	qv.Product(&yv, yv.T(), &qv)
	qv.Scale(math.Sqrt(float64(c.n-1)), &qv)
	return &qv
}

// CanonicalCorrelations returns a CC which can provide the results of canonical
// correlation analysis of the input data x and y, columns of which should be
// interpretable as two sets of measurements on the same observations (rows).
// These observations are optionally weighted by weights.
//
// Canonical correlation analysis finds associations between two sets of
// variables on the same observations by finding linear combinations of the two
// sphered datasets that maximise the correlation between them.
//
// Some notation: let Xc and Yc denote the centered input data matrices x
// and y (column means subtracted from each column), let Sx and Sy denote the
// sample covariance matrices within x and y respectively, and let Sxy denote
// the covariance matrix between x and y. The sphered data can then be expressed
// as Xc * Sx^{-1/2} and Yc * Sy^{-1/2} respectively, and the correlation matrix
// between the sphered data is called the canonical correlation matrix,
// Sx^{-1/2} * Sxy * Sy^{-1/2}. In cases where S^{-1/2} is ambiguous for some
// covariance matrix S, S^{-1/2} is taken to be E * D^{-1/2} * E^T where S can
// be eigendecomposed as S = E * D * E^T.
//
// The canonical correlations are the correlations between the corresponding
// pairs of canonical variables and can be obtained with c.Corrs(). Canonical
// variables can be obtained by projecting the sphered data into the left and
// right eigenvectors of the canonical correlation matrix, and these
// eigenvectors can be obtained with c.Left(true) and c.Right(true)
// respectively. The canonical variables can also be obtained directly from the
// centered raw data by using the back-transformed eigenvectors which can be
// obtained with c.Left(false) and c.Right(false) respectively.
//
// The first pair of left and right eigenvectors of the canonical correlation
// matrix can be interpreted as directions into which the respective sphered
// data can be projected such that the correlation between the two projections
// is maximised. The second pair and onwards solve the same optimization but
// under the constraint that they are uncorrelated (orthogonal in sphered space)
// to previous projections.
//
// CanonicalCorrelations will panic if the inputs x and y do not have the same
// number of rows.
//
// The slice weights is used to weight the observations. If weights is nil, each
// weight is considered to have a value of one, otherwise the length of weights
// must match the number of observations (rows of both x and y) or CanonicalCorrelations will panic.
//
// More details can be found at
// https://en.wikipedia.org/wiki/Canonical_correlation
// or in Chapter 3 of
// Koch, Inge. Analysis of multivariate and high-dimensional data.
// Vol. 32. Cambridge University Press, 2013. ISBN: 9780521887939
func CanonicalCorrelations(x, y mat64.Matrix, weights []float64) (c CC, err error) {
	n, _ := x.Dims()
	yn, _ := y.Dims()
	if n != yn {
		panic("stat: unequal number of observations")
	}
	if weights != nil && len(weights) != n {
		panic("stat: len(weights) != observations")
	}

	// Center and factorize x and y.
	xsvd, ok := svdFactorizeCentered(x, weights)
	if !ok {
		return CC{}, errors.New("stat: failed to factorize x")
	}
	ysvd, ok := svdFactorizeCentered(y, weights)
	if !ok {
		return CC{}, errors.New("stat: failed to factorize y")
	}
	var xu, xv, yu, yv mat64.Dense
	xu.UFromSVD(xsvd)
	xv.VFromSVD(xsvd)
	yu.UFromSVD(ysvd)
	yv.VFromSVD(ysvd)

	// Calculate and factorise the canonical correlation matrix.
	var ccor mat64.Dense
	ccor.Product(&xv, xu.T(), &yu, yv.T())
	var csvd mat64.SVD
	ok = csvd.Factorize(&ccor, svdKind(ccor.Dims()))
	if !ok {
		return CC{}, errors.New("stat: failed to factorize ccor")
	}

	return CC{n: n, x: xsvd, y: ysvd, c: &csvd}, nil
}

func svdKind(n, d int) matrix.SVDKind {
	if n > d {
		return matrix.SVDThin
	}
	return matrix.SVDFull
}

func svdFactorizeCentered(m mat64.Matrix, weights []float64) (svd *mat64.SVD, ok bool) {
	n, d := m.Dims()
	centered := mat64.NewDense(n, d, nil)
	col := make([]float64, n)
	for j := 0; j < d; j++ {
		mat64.Col(col, j, m)
		floats.AddConst(-Mean(col, weights), col)
		centered.SetCol(j, col)
	}
	for i, w := range weights {
		floats.Scale(math.Sqrt(w), centered.RawRowView(i))
	}
	svd = &mat64.SVD{}
	ok = svd.Factorize(centered, svdKind(n, d))
	if !ok {
		svd = nil
	}
	return svd, ok
}

// scaleColsReciSqrt scales the columns of cols
// by the reciprocal square-root of vals.
func scaleColsReciSqrt(cols *mat64.Dense, vals []float64) {
	if cols == nil {
		panic("stat: input nil")
	}
	n, d := cols.Dims()
	if len(vals) != d {
		panic("stat: input length mismatch")
	}
	col := make([]float64, n)
	for j := 0; j < d; j++ {
		mat64.Col(col, j, cols)
		floats.Scale(math.Sqrt(1/vals[j]), col)
		cols.SetCol(j, col)
	}
}
