// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat"
	"github.com/gonum/stat/distuv"
)

var (
	badInputLength = "distmv: input slice length mismatch"
)

// Normal is a multivariate normal distribution (also known as the multivariate
// Gaussian distribution). Its pdf in k dimensions is given by
//  (2 π)^(-k/2) |Σ|^(-1/2) exp(-1/2 (x-μ)'Σ^-1(x-μ))
// where μ is the mean vector and Σ the covariance matrix. Σ must be symmetric
// and positive definite. Use NewNormal to construct.
type Normal struct {
	mu []float64

	sigma mat64.SymDense

	chol       mat64.Cholesky
	lower      mat64.TriDense
	logSqrtDet float64
	dim        int

	src *rand.Rand
}

// NewNormal creates a new Normal with the given mean and covariance matrix.
// NewNormal panics if len(mu) == 0, or if len(mu) != sigma.N. If the covariance
// matrix is not positive-definite, the returned boolean is false.
func NewNormal(mu []float64, sigma mat64.Symmetric, src *rand.Rand) (*Normal, bool) {
	if len(mu) == 0 {
		panic(badZeroDimension)
	}
	dim := sigma.Symmetric()
	if dim != len(mu) {
		panic(badSizeMismatch)
	}
	n := &Normal{
		src: src,
		dim: dim,
		mu:  make([]float64, dim),
	}
	copy(n.mu, mu)
	ok := n.chol.Factorize(sigma)
	if !ok {
		return nil, false
	}
	n.sigma = *mat64.NewSymDense(dim, nil)
	n.sigma.CopySym(sigma)
	n.lower.LFromCholesky(&n.chol)
	n.logSqrtDet = 0.5 * n.chol.LogDet()
	return n, true
}

// NewNormalChol creates a new Normal distribution with the given mean and
// covariance matrix represented by its Cholesky decomposition. NewNormalChol
// panics if len(mu) is not equal to chol.Size().
func NewNormalChol(mu []float64, chol *mat64.Cholesky, src *rand.Rand) *Normal {
	dim := len(mu)
	if dim != chol.Size() {
		panic(badSizeMismatch)
	}
	n := &Normal{
		src: src,
		dim: dim,
		mu:  make([]float64, dim),
	}
	n.chol.Clone(chol)
	copy(n.mu, mu)
	n.lower.LFromCholesky(chol)
	n.logSqrtDet = 0.5 * n.chol.LogDet()
	return n
}

// NewNormalPrecision creates a new Normal distribution with the given mean and
// precision matrix (inverse of the covariance matrix). NewNormalPrecision
// panics if len(mu) is not equal to prec.Symmetric(). If the precision matrix
// is not positive-definite, NewNormalPrecision returns nil for norm and false
// for ok.
func NewNormalPrecision(mu []float64, prec *mat64.SymDense, src *rand.Rand) (norm *Normal, ok bool) {
	if len(mu) == 0 {
		panic(badZeroDimension)
	}
	dim := prec.Symmetric()
	if dim != len(mu) {
		panic(badSizeMismatch)
	}
	// TODO(btracey): Computing a matrix inverse is generally numerically instable.
	// This only has to compute the inverse of a positive definite matrix, which
	// is much better, but this still loses precision. It is worth considering if
	// instead the precision matrix should be stored explicitly and used instead
	// of the Cholesky decomposition of the covariance matrix where appropriate.
	var chol mat64.Cholesky
	ok = chol.Factorize(prec)
	if !ok {
		return nil, false
	}
	var sigma mat64.SymDense
	sigma.InverseCholesky(&chol)
	return NewNormal(mu, &sigma, src)
}

// ConditionNormal returns the Normal distribution that is the receiver conditioned
// on the input evidence. The returned multivariate normal has dimension
// n - len(observed), where n is the dimension of the original receiver. The updated
// mean and covariance are
//  mu = mu_un + sigma_{ob,un}^T * sigma_{ob,ob}^-1 (v - mu_ob)
//  sigma = sigma_{un,un} - sigma_{ob,un}^T * sigma_{ob,ob}^-1 * sigma_{ob,un}
// where mu_un and mu_ob are the original means of the unobserved and observed
// variables respectively, sigma_{un,un} is the unobserved subset of the covariance
// matrix, sigma_{ob,ob} is the observed subset of the covariance matrix, and
// sigma_{un,ob} are the cross terms. The elements of x_2 have been observed with
// values v. The dimension order is preserved during conditioning, so if the value
// of dimension 1 is observed, the returned normal represents dimensions {0, 2, ...}
// of the original Normal distribution.
//
// ConditionNormal returns {nil, false} if there is a failure during the update.
// Mathematically this is impossible, but can occur with finite precision arithmetic.
func (n *Normal) ConditionNormal(observed []int, values []float64, src *rand.Rand) (*Normal, bool) {
	if len(observed) == 0 {
		panic("normal: no observed value")
	}
	if len(observed) != len(values) {
		panic(badInputLength)
	}
	for _, v := range observed {
		if v < 0 || v >= n.Dim() {
			panic("normal: observed value out of bounds")
		}
	}

	_, mu1, sigma11 := studentsTConditional(observed, values, math.Inf(1), n.mu, &n.sigma)
	if mu1 == nil {
		return nil, false
	}
	return NewNormal(mu1, sigma11, src)
}

// CovarianceMatrix returns the covariance matrix of the distribution. Upon
// return, the value at element {i, j} of the covariance matrix is equal to
// the covariance of the i^th and j^th variables.
//  covariance(i, j) = E[(x_i - E[x_i])(x_j - E[x_j])]
// If the input matrix is nil a new matrix is allocated, otherwise the result
// is stored in-place into the input.
func (n *Normal) CovarianceMatrix(s *mat64.SymDense) *mat64.SymDense {
	if s == nil {
		s = mat64.NewSymDense(n.Dim(), nil)
	}
	sn := s.Symmetric()
	if sn != n.Dim() {
		panic("normal: input matrix size mismatch")
	}
	s.CopySym(&n.sigma)
	return s
}

// Dim returns the dimension of the distribution.
func (n *Normal) Dim() int {
	return n.dim
}

// Entropy returns the differential entropy of the distribution.
func (n *Normal) Entropy() float64 {
	return float64(n.dim)/2*(1+logTwoPi) + n.logSqrtDet
}

// LogProb computes the log of the pdf of the point x.
func (n *Normal) LogProb(x []float64) float64 {
	dim := n.dim
	if len(x) != dim {
		panic(badSizeMismatch)
	}
	c := -0.5*float64(dim)*logTwoPi - n.logSqrtDet
	dst := stat.Mahalanobis(mat64.NewVector(dim, x), mat64.NewVector(dim, n.mu), &n.chol)
	return c - 0.5*dst*dst
}

// MarginalNormal returns the marginal distribution of the given input variables.
// That is, MarginalNormal returns
//  p(x_i) = \int_{x_o} p(x_i | x_o) p(x_o) dx_o
// where x_i are the dimensions in the input, and x_o are the remaining dimensions.
// See https://en.wikipedia.org/wiki/Marginal_distribution for more information.
//
// The input src is passed to the call to NewNormal.
func (n *Normal) MarginalNormal(vars []int, src *rand.Rand) (*Normal, bool) {
	newMean := make([]float64, len(vars))
	for i, v := range vars {
		newMean[i] = n.mu[v]
	}
	var s mat64.SymDense
	s.SubsetSym(&n.sigma, vars)
	return NewNormal(newMean, &s, src)
}

// MarginalNormalSingle returns the marginal of the given input variable.
// That is, MarginalNormal returns
//  p(x_i) = \int_{x_¬i} p(x_i | x_¬i) p(x_¬i) dx_¬i
// where i is the input index.
// See https://en.wikipedia.org/wiki/Marginal_distribution for more information.
//
// The input src is passed to the constructed distuv.Normal.
func (n *Normal) MarginalNormalSingle(i int, src *rand.Rand) distuv.Normal {
	return distuv.Normal{
		Mu:     n.mu[i],
		Sigma:  math.Sqrt(n.sigma.At(i, i)),
		Source: src,
	}
}

// Mean returns the mean of the probability distribution at x. If the
// input argument is nil, a new slice will be allocated, otherwise the result
// will be put in-place into the receiver.
func (n *Normal) Mean(x []float64) []float64 {
	x = reuseAs(x, n.dim)
	copy(x, n.mu)
	return x
}

// Prob computes the value of the probability density function at x.
func (n *Normal) Prob(x []float64) float64 {
	return math.Exp(n.LogProb(x))
}

// Quantile returns the multi-dimensional inverse cumulative distribution function.
// If x is nil, a new slice will be allocated and returned. If x is non-nil,
// len(x) must equal len(p) and the quantile will be stored in-place into x.
// All of the values of p must be between 0 and 1, inclusive, or Quantile will panic.
func (n *Normal) Quantile(x, p []float64) []float64 {
	dim := n.Dim()
	if len(p) != dim {
		panic(badInputLength)
	}
	if x == nil {
		x = make([]float64, dim)
	}
	if len(x) != len(p) {
		panic(badInputLength)
	}

	// Transform to a standard normal and then transform to a multivariate Gaussian.
	tmp := make([]float64, len(x))
	for i, v := range p {
		tmp[i] = distuv.UnitNormal.Quantile(v)
	}
	n.TransformNormal(x, tmp)
	return x
}

// Rand generates a random number according to the distributon.
// If the input slice is nil, new memory is allocated, otherwise the result is stored
// in place.
func (n *Normal) Rand(x []float64) []float64 {
	x = reuseAs(x, n.dim)
	tmp := make([]float64, n.dim)
	if n.src == nil {
		for i := range x {
			tmp[i] = rand.NormFloat64()
		}
	} else {
		for i := range x {
			tmp[i] = n.src.NormFloat64()
		}
	}
	n.transformNormal(x, tmp)
	return x
}

// SetMean changes the mean of the normal distribution. SetMean panics if len(mu)
// does not equal the dimension of the normal distribution.
func (n *Normal) SetMean(mu []float64) {
	if len(mu) != n.Dim() {
		panic(badSizeMismatch)
	}
	copy(n.mu, mu)
}

// TransformNormal transforms the vector, normal, generated from a standard
// multidimensional normal into a vector that has been generated under the
// distribution of the receiver.
//
// If dst is non-nil, the result will be stored into dst, otherwise a new slice
// will be allocated. TransformNormal will panic if the length of normal is not
// the dimension of the receiver, or if dst is non-nil and len(dist) != len(normal).
func (n *Normal) TransformNormal(dst, normal []float64) []float64 {
	if len(normal) != n.dim {
		panic(badInputLength)
	}
	dst = reuseAs(dst, n.dim)
	if len(dst) != len(normal) {
		panic(badInputLength)
	}
	n.transformNormal(dst, normal)
	return dst
}

// transformNormal performs the same operation as TransformNormal except no
// safety checks are performed and both input slices must be non-nil.
func (n *Normal) transformNormal(dst, normal []float64) []float64 {
	srcVec := mat64.NewVector(n.dim, normal)
	dstVec := mat64.NewVector(n.dim, dst)
	dstVec.MulVec(&n.lower, srcVec)
	floats.Add(dst, n.mu)
	return dst
}
