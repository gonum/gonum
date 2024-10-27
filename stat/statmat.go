// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"errors"
	"fmt"
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
	ErrNoOptions          = errors.New("no initialized model options")
	ErrTargetLenMismatch  = errors.New("target length does not match target rows")
	ErrNoTrainingMatrix   = errors.New("no training matrix")
	ErrNoTargetMatrix     = errors.New("no target matrix")
	ErrNoDesignMatrix     = errors.New("no design matrix for inference")
	ErrFeatureLenMismatch = errors.New("number of features does not match number of model coefficients")

	// Lasso Errors
	ErrNegativeLambda     = errors.New("negative lambda")
	ErrNegativeIterations = errors.New("negative iterations")
	ErrNegativeTolerance  = errors.New("negative tolerance")
	ErrWarmStartBetaSize  = errors.New("warm start beta does not have the same number of coefficients as training features")
)

// OLSOptions represents input options to run the OLS Regression
type OLSOptions struct {
	// FitIntercept adds a constant 1.0 feature as the first column if set to true
	FitIntercept bool
}

// Validate runs basic validation on OLS options
func (o *OLSOptions) Validate() (*OLSOptions, error) {
	if o == nil {
		o = NewDefaultOLSOptions()
	}

	return o, nil
}

// NewDefaultOLSOptions returns a default set of OLS Regression options
func NewDefaultOLSOptions() *OLSOptions {
	return &OLSOptions{
		FitIntercept: true,
	}
}

// OLSRegression computes ordinary least squares using QR factorization
type OLSRegression struct {
	opt       *OLSOptions
	coef      []float64
	intercept float64
}

// NewOLSRegression initializes an ordinary least squares model ready for fitting
func NewOLSRegression(opt *OLSOptions) (*OLSRegression, error) {
	opt, err := opt.Validate()
	if err != nil {
		return nil, err
	}
	return &OLSRegression{
		opt: opt,
	}, nil
}

// Fit the model according to the given training data
func (o *OLSRegression) Fit(x, y mat.Matrix) error {
	if o.opt == nil {
		return ErrNoOptions
	}
	if x == nil {
		return ErrNoTrainingMatrix
	}
	if y == nil {
		return ErrNoTargetMatrix
	}
	m, n := x.Dims()

	ym, _ := y.Dims()
	if ym != m {
		return fmt.Errorf("training data has %d rows and target has %d row, %w", m, ym, ErrTargetLenMismatch)
	}

	if o.opt.FitIntercept {
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

	if o.opt.FitIntercept {
		o.intercept = c[0]
		o.coef = c[1:]
	} else {
		o.coef = c
	}

	return nil
}

// Predict using the OLS model
func (o *OLSRegression) Predict(x mat.Matrix) ([]float64, error) {
	if o.opt == nil {
		return nil, ErrNoOptions
	}
	if x == nil {
		return nil, ErrNoDesignMatrix
	}

	coef := o.coef
	if o.opt.FitIntercept {
		coef = append([]float64{o.intercept}, o.coef...)

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
		return nil, fmt.Errorf("got %d features in design matrix, but expected %d, %w", xn, n, ErrFeatureLenMismatch)
	}
	coefMx := mat.NewDense(1, n, coef)

	var res mat.Dense
	res.Mul(coefMx, xT)
	return res.RawRowView(0), nil
}

// Score computes the coefficient of determination of the prediction
func (o *OLSRegression) Score(x, y mat.Matrix) (float64, error) {
	if o.opt == nil {
		return 0.0, ErrNoOptions
	}
	if x == nil {
		return 0.0, ErrNoDesignMatrix
	}
	if y == nil {
		return 0.0, ErrNoTargetMatrix
	}

	m, _ := x.Dims()

	ym, _ := y.Dims()
	if m != ym {
		return 0.0, fmt.Errorf("design matrix has %d rows and target has %d rows, %w", m, ym, ErrTargetLenMismatch)
	}

	res, err := o.Predict(x)
	if err != nil {
		return 0.0, err
	}

	ySlice := mat.Col(nil, 0, y)

	return RSquaredFrom(res, ySlice, nil), nil
}

// Intercept returns the computed intercept if FitIntercept is set to true. Defaults to 0.0 if not set.
func (o *OLSRegression) Intercept() float64 {
	return o.intercept
}

// Coef returns a slice of the trained coefficients in the same order of the training feature Matrix by column.
func (o *OLSRegression) Coef() []float64 {
	c := make([]float64, len(o.coef))
	copy(c, o.coef)
	return c
}

// LassoOptions represents input options to run the Lasso Regression
type LassoOptions struct {
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
}

// Validate runs basic validation on Lasso options
func (l *LassoOptions) Validate() (*LassoOptions, error) {
	if l == nil {
		l = NewDefaultLassoOptions()
	}

	if l.Lambda < 0 {
		return nil, ErrNegativeLambda
	}
	if l.Iterations < 0 {
		return nil, ErrNegativeIterations
	}
	if l.Tolerance < 0 {
		return nil, ErrNegativeTolerance
	}
	return l, nil
}

// NewDefaultLassoOptions returns a default set of Lasso Regression options
func NewDefaultLassoOptions() *LassoOptions {
	return &LassoOptions{
		Lambda:        1.0,
		Iterations:    1000,
		Tolerance:     1e-4,
		WarmStartBeta: nil,
		FitIntercept:  true,
	}
}

// LassoRegression computes the lasso regression using coordinate descent. lambda = 0 converges to OLS
type LassoRegression struct {
	opt *LassoOptions

	coef      []float64
	intercept float64
}

// NewLassoRegression initializes a Lasso model ready for fitting
func NewLassoRegression(opt *LassoOptions) (*LassoRegression, error) {
	opt, err := opt.Validate()
	if err != nil {
		return nil, err
	}
	return &LassoRegression{
		opt: opt,
	}, nil
}

// Fit the model according to the given training data
func (l *LassoRegression) Fit(x, y mat.Matrix) error {
	if l.opt == nil {
		return ErrNoOptions
	}
	if x == nil {
		return ErrNoTrainingMatrix
	}
	if y == nil {
		return ErrNoTargetMatrix
	}

	m, n := x.Dims()

	ym, _ := y.Dims()
	if ym != m {
		return fmt.Errorf("training data has %d rows and target has %d row, %w", m, ym, ErrTargetLenMismatch)
	}

	if l.opt.FitIntercept {
		ones := make([]float64, m)
		floats.AddConst(1.0, ones)
		onesMx := mat.NewDense(1, m, ones)
		xT := x.T()

		var xWithOnes mat.Dense
		xWithOnes.Stack(onesMx, xT)
		x = xWithOnes.T()
		_, n = x.Dims()
	}

	if l.opt.WarmStartBeta != nil && len(l.opt.WarmStartBeta) != n {
		return fmt.Errorf("warm start beta has %d features instead of %d, %w", len(l.opt.WarmStartBeta), n, ErrWarmStartBetaSize)
	}

	// tracks current betas
	beta := make([]float64, n)
	if l.opt.WarmStartBeta != nil {
		copy(beta, l.opt.WarmStartBeta)
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
	for i := 0; i < l.opt.Iterations; i++ {
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

			gamma := l.opt.Lambda / xdot[j]
			betaNext = SoftThreshold(betaNext, gamma)

			maxCoef = math.Max(maxCoef, betaNext)
			maxUpdate = math.Max(maxUpdate, math.Abs(betaNext-betaCurr))
			betaDiff = betaNext - betaCurr
			floats.ScaleTo(betaXDelta, betaDiff, obsCol)
			beta[j] = betaNext
		}

		// break early if we've achieved the desired tolerance
		if maxUpdate < l.opt.Tolerance*maxCoef {
			break
		}
	}

	if l.opt.FitIntercept {
		l.intercept = beta[0]
		l.coef = beta[1:]
	} else {
		l.coef = beta
	}

	return nil
}

// Predict using the Lasso model
func (l *LassoRegression) Predict(x mat.Matrix) ([]float64, error) {
	if l.opt == nil {
		return nil, ErrNoOptions
	}
	if x == nil {
		return nil, ErrNoDesignMatrix
	}

	coef := l.coef
	if l.opt.FitIntercept {
		coef = append([]float64{l.intercept}, l.coef...)

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
		return nil, fmt.Errorf("got %d features in design matrix, but expected %d, %w", xn, n, ErrFeatureLenMismatch)
	}

	xT := x.T()
	coefMx := mat.NewDense(1, n, coef)

	var res mat.Dense
	res.Mul(coefMx, xT)
	return res.RawRowView(0), nil
}

// Score computes the coefficient of determination of the prediction
func (l *LassoRegression) Score(x, y mat.Matrix) (float64, error) {
	if l.opt == nil {
		return 0.0, ErrNoOptions
	}
	if x == nil {
		return 0.0, ErrNoDesignMatrix
	}
	if y == nil {
		return 0.0, ErrNoTargetMatrix
	}

	m, _ := x.Dims()

	ym, _ := y.Dims()
	if m != ym {
		return 0.0, fmt.Errorf("design matrix has %d rows and target has %d rows, %w", m, ym, ErrTargetLenMismatch)
	}

	res, err := l.Predict(x)
	if err != nil {
		return 0.0, err
	}

	ySlice := mat.Col(nil, 0, y)

	return RSquaredFrom(res, ySlice, nil), nil
}

// Intercept returns the computed intercept if FitIntercept is set to true. Defaults to 0.0 if not set.
func (l *LassoRegression) Intercept() float64 {
	return l.intercept
}

// Coef returns a slice of the trained coefficients in the same order of the training feature Matrix by column.
func (l *LassoRegression) Coef() []float64 {
	return l.coef
}

// SoftThreshold returns 0.0 if the value is less than or equal to the gamma input
func SoftThreshold(x, gamma float64) float64 {
	res := math.Max(0, math.Abs(x)-gamma)
	if math.Signbit(x) {
		return -res
	}
	return res
}
