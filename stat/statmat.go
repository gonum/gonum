// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// CovarianceMatrix calculates the covariance matrix (also known as the
// variance-covariance matrix) calculated from a matrix of data, x, using
// a two-pass algorithm. The result is stored in dst.
//
// If weights is not nil the weighted covariance of x is calculated. weights
// must have length equal to the number of rows in input data matrix and
// must not contain negative elements.
// The dst matrix must either be empty or have the same number of
// columns as the input data matrix.
func CovarianceMatrix(dst *mat.SymDense, x mat.Matrix, weights []float64) {
	// This is the matrix version of the two-pass algorithm. It doesn't use the
	// additional floating point error correction that the Covariance function uses
	// to reduce the impact of rounding during centering.

	r, c := x.Dims()

	if dst.IsEmpty() {
		*dst = *(dst.GrowSym(c).(*mat.SymDense))
	} else if n := dst.SymmetricDim(); n != c {
		panic(mat.ErrShape)
	}

	var xt mat.Dense
	xt.CloneFrom(x.T())
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
		dst.SymOuterK(1/(float64(r)-1), &xt)
		return
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
	dst.SymOuterK(1/(floats.Sum(weights)-1), &xt)
}

// CorrelationMatrix returns the correlation matrix calculated from a matrix
// of data, x, using a two-pass algorithm. The result is stored in dst.
//
// If weights is not nil the weighted correlation of x is calculated. weights
// must have length equal to the number of rows in input data matrix and
// must not contain negative elements.
// The dst matrix must either be empty or have the same number of
// columns as the input data matrix.
func CorrelationMatrix(dst *mat.SymDense, x mat.Matrix, weights []float64) {
	// This will panic if the sizes don't match, or if weights is the wrong size.
	CovarianceMatrix(dst, x, weights)
	covToCorr(dst)
}

// covToCorr converts a covariance matrix to a correlation matrix.
func covToCorr(c *mat.SymDense) {
	r := c.SymmetricDim()

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
//
//	D = sqrt((x-y)ᵀ * Σ^-1 * (x-y))
//
// between the column vectors x and y given the cholesky decomposition of Σ.
// Mahalanobis returns NaN if the linear solve fails.
//
// See https://en.wikipedia.org/wiki/Mahalanobis_distance for more information.
func Mahalanobis(x, y mat.Vector, chol *mat.Cholesky) float64 {
	var diff mat.VecDense
	diff.SubVec(x, y)
	var tmp mat.VecDense
	err := chol.SolveVecTo(&tmp, &diff)
	if err != nil {
		return math.NaN()
	}
	return math.Sqrt(mat.Dot(&tmp, &diff))
}

var (
	ErrNegativeLambda     = errors.New("negative lambda")
	ErrNegativeIterations = errors.New("negative iterations")
	ErrNegativeTolerance  = errors.New("negative tolerance")
	ErrWarmStartBetaSize  = errors.New("warm start beta does not have the same number of coefficients as training features")
)

// OLSModel represents the model options and coefficients for Ordinary Leaast Squares Regression
type OLSModel struct {
	FitIntercept bool

	Intercept float64
	Coef      []float64
}

// OLSRegression fits the linear model according to the given training data. OLS Options must be set and the training, target matrices
// must also be non nil. Number of rows in x must match the number of rows in y or this.
func (o *OLSModel) Fit(x, y mat.Matrix) (float64, []float64) {
	m, n := x.Dims()

	ym, _ := y.Dims()
	if ym != m {
		panic(mat.ErrShape)
	}

	if o.FitIntercept {
		ones := make([]float64, m)
		floats.AddConst(1.0, ones)
		onesMx := mat.NewDense(1, m, ones)
		xT := x.T()

		var xWithOnes mat.Dense
		xWithOnes.Stack(onesMx, xT)
		x = xWithOnes.T()
		_, n = x.Dims()
	}

	yT := y.T()

	qr := new(mat.QR)
	qr.Factorize(x)

	q := new(mat.Dense)
	r := new(mat.Dense)

	qr.QTo(q)
	qr.RTo(r)
	yq := new(mat.Dense)
	yq.Mul(yT, q)

	c := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		c[i] = yq.At(0, i)
		for j := i + 1; j < n; j++ {
			c[i] -= c[j] * r.At(i, j)
		}
		c[i] /= r.At(i, i)
	}

	o.Coef = c
	if o.FitIntercept {
		o.Intercept = c[0]
		o.Coef = c[1:]
	}
	return o.Intercept, o.Coef
}

// Predict using the OLS model. OLS options must be initialized before calling and input x must not be nil.
// The number of columns in the matrix must match the number of model coefficients.
func (o *OLSModel) Predict(x mat.Matrix) []float64 {
	coef := o.Coef
	if o.FitIntercept {
		coef = append([]float64{o.Intercept}, o.Coef...)

		m, _ := x.Dims()
		ones := make([]float64, m)
		floats.AddConst(1.0, ones)
		onesMx := mat.NewDense(1, m, ones)
		xT := x.T()

		var xWithOnes mat.Dense
		xWithOnes.Stack(onesMx, xT)
		x = xWithOnes.T()
	}
	n := len(coef)

	xT := x.T()
	xn, _ := xT.Dims()
	if xn != n {
		panic(mat.ErrShape)
	}
	coefMx := mat.NewDense(1, n, coef)

	var res mat.Dense
	res.Mul(coefMx, xT)
	return res.RawRowView(0)
}

// Score computes the coefficient of determination of the prediction. This will panic if no
// OLS options have been initialized or if x, y are nil. x and y must have the same number of
// rows or will panic.
func (o *OLSModel) Score(x, y mat.Matrix) float64 {
	m, _ := x.Dims()

	ym, _ := y.Dims()
	if m != ym {
		panic(mat.ErrShape)
	}

	res := o.Predict(x)

	ySlice := mat.Col(nil, 0, y)

	return RSquaredFrom(res, ySlice, nil)
}

// LassoModel represents model options and coefficients for Lasso Regression
type LassoModel struct {
	// WarmStartBeta is used to prime the coordinate descent to reduce the training time if a previous
	// fit has been performed.
	WarmStartBeta []float64

	// Lambda represents the L1 multiplier, controlling the regularization. Must be a non-negative. 0.0 results in converging
	// to Ordinary Least Squares (OLS).
	Lambda float64

	// Iterations is the maximum number of times the fit loops through training all coefficients.
	Iterations int

	// Tolerance is the smallest coefficient channge on each iteration to determine when to stop iterating.
	Tolerance float64

	// FitIntercept adds a constant 1.0 feature as the first column if set to true
	FitIntercept bool

	Intercept float64
	Coef      []float64
}

// Fit the model according to the given training data. LassoOptions must be set and the training, target matrices
// must also be non nil. Number of rows in x must match the number of rows in y or this. If WarmStartBeta has been
// set through the initialized options, the length must match that of the number of columns in x.
func (l *LassoModel) Fit(x, y mat.Matrix) (float64, []float64) {
	if l.Lambda < 0 {
		panic(ErrNegativeLambda)
	}
	if l.Iterations < 0 {
		panic(ErrNegativeIterations)
	}
	if l.Tolerance < 0 {
		panic(ErrNegativeTolerance)
	}

	m, n := x.Dims()

	ym, _ := y.Dims()
	if ym != m {
		panic(mat.ErrShape)
	}

	if l.FitIntercept {
		ones := make([]float64, m)
		floats.AddConst(1.0, ones)
		onesMx := mat.NewDense(1, m, ones)
		xT := x.T()

		var xWithOnes mat.Dense
		xWithOnes.Stack(onesMx, xT)
		x = xWithOnes.T()
		_, n = x.Dims()
	}

	if l.WarmStartBeta != nil && len(l.WarmStartBeta) != n {
		panic(ErrWarmStartBetaSize)
	}

	// tracks current betas
	beta := make([]float64, n)
	if l.WarmStartBeta != nil {
		copy(beta, l.WarmStartBeta)
	}

	xcols := make([][]float64, n)

	// precompute the per feature dot product
	xdot := make([]float64, n)
	for i := 0; i < n; i++ {
		xi := mat.Col(nil, i, x)
		xdot[i] = floats.Dot(xi, xi)
		xcols[i] = xi
	}

	// tracks the per coordinate residual
	residual := make([]float64, m)

	// tracks the current beta * x by adding the deltas on each beta iteration
	betaX := make([]float64, m)

	// tracks the delta of the beta * x of each iteration by computing the next beta
	// multiplied by the feature observations of that beta. will be added to betaX on
	// the next beta iteration
	betaXDelta := make([]float64, m)

	yArr := mat.Col(nil, 0, y)
	for i := 0; i < l.Iterations; i++ {
		maxCoef := 0.0
		maxUpdate := 0.0
		betaDiff := 0.0

		// loop through all features and minimize loss function
		for j := 0; j < n; j++ {
			betaCurr := beta[j]
			if i != 0 {
				if betaCurr == 0 {
					continue
				}
			}

			floats.Add(betaX, betaXDelta)
			floats.SubTo(residual, yArr, betaX)

			obsCol := xcols[j]
			num := floats.Dot(obsCol, residual)
			betaNext := num/xdot[j] + betaCurr

			gamma := l.Lambda / xdot[j]
			betaNext = softThreshold(betaNext, gamma)

			maxCoef = math.Max(maxCoef, betaNext)
			maxUpdate = math.Max(maxUpdate, math.Abs(betaNext-betaCurr))
			betaDiff = betaNext - betaCurr
			floats.ScaleTo(betaXDelta, betaDiff, obsCol)
			beta[j] = betaNext
		}

		// break early if we've achieved the desired tolerance
		if maxUpdate < l.Tolerance*maxCoef {
			break
		}
	}

	l.Coef = beta
	if l.FitIntercept {
		l.Intercept = beta[0]
		l.Coef = beta[1:]
	}
	return l.Intercept, l.Coef
}

// Predict using the Lasso model. Lasso options must be initialized and input matrix must be non-nil.
// The number of columns in the matrix must match the number of model coefficients.
func (l *LassoModel) Predict(x mat.Matrix) []float64 {
	coef := l.Coef
	if l.FitIntercept {
		coef = append([]float64{l.Intercept}, l.Coef...)

		m, _ := x.Dims()
		ones := make([]float64, m)
		floats.AddConst(1.0, ones)
		onesMx := mat.NewDense(1, m, ones)
		xT := x.T()

		var xWithOnes mat.Dense
		xWithOnes.Stack(onesMx, xT)
		x = xWithOnes.T()
	}
	n := len(coef)

	_, xn := x.Dims()
	if xn != n {
		panic(mat.ErrShape)
	}

	xT := x.T()
	coefMx := mat.NewDense(1, n, coef)

	var res mat.Dense
	res.Mul(coefMx, xT)
	return res.RawRowView(0)
}

// Score computes the coefficient of determination of the prediction. Lasso options must be initialized and x, y
// must be non-nil. x and y must have the same number of rows.
func (l *LassoModel) Score(x, y mat.Matrix) float64 {
	m, _ := x.Dims()

	ym, _ := y.Dims()
	if m != ym {
		panic(mat.ErrShape)
	}

	res := l.Predict(x)

	ySlice := mat.Col(nil, 0, y)

	return RSquaredFrom(res, ySlice, nil)
}

// softThreshold returns 0.0 if the value is less than or equal to the gamma input
func softThreshold(x, gamma float64) float64 {
	switch {
	case x < -gamma:
		return x + gamma
	case gamma < x:
		return x - gamma
	default:
		return 0
	}
}
