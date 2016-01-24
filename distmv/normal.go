// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand"
	"sync"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

// Normal is a multivariate normal distribution (also known as the multivariate
// Gaussian distribution). Its pdf in k dimensions is given by
//  (2 π)^(-k/2) |Σ|^(-1/2) exp(-1/2 (x-μ)'Σ^-1(x-μ))
// where μ is the mean vector and Σ the covariance matrix. Σ must be symmetric
// and positive definite. Use NewNormal to construct.
type Normal struct {
	mu []float64

	once  sync.Once
	sigma *mat64.SymDense // only stored if needed

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
	copy(n.mu, mu)
	n.lower.LFromCholesky(chol)
	n.logSqrtDet = 0.5 * n.chol.LogDet()
	return n
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
		panic("normal: input slice length mismatch")
	}
	for _, v := range observed {
		if v < 0 || v >= n.Dim() {
			panic("normal: observed value out of bounds")
		}
	}

	ob := len(observed)
	unob := n.Dim() - ob
	obMap := make(map[int]struct{})
	for _, v := range observed {
		if _, ok := obMap[v]; ok {
			panic("normal: observed dimension occurs twice")
		}
		obMap[v] = struct{}{}
	}
	if len(observed) == n.Dim() {
		panic("normal: all dimensions observed")
	}
	unobserved := make([]int, 0, unob)
	for i := 0; i < n.Dim(); i++ {
		if _, ok := obMap[i]; !ok {
			unobserved = append(unobserved, i)
		}
	}
	mu1 := make([]float64, unob)
	for i, v := range unobserved {
		mu1[i] = n.mu[v]
	}
	mu2 := make([]float64, ob) // really v - mu2
	for i, v := range observed {
		mu2[i] = values[i] - n.mu[v]
	}

	n.setSigma()

	var sigma11, sigma22 mat64.SymDense
	sigma11.SubsetSym(n.sigma, unobserved)
	sigma22.SubsetSym(n.sigma, observed)

	sigma21 := mat64.NewDense(ob, unob, nil)
	for i, r := range observed {
		for j, c := range unobserved {
			v := n.sigma.At(r, c)
			sigma21.Set(i, j, v)
		}
	}

	var chol mat64.Cholesky
	ok := chol.Factorize(&sigma22)
	if !ok {
		return nil, ok
	}

	// Compute sigma_{2,1}^T * sigma_{2,2}^-1 (v - mu_2).
	v := mat64.NewVector(ob, mu2)
	var tmp, tmp2 mat64.Vector
	err := tmp.SolveCholeskyVec(&chol, v)
	if err != nil {
		return nil, false
	}
	tmp2.MulVec(sigma21.T(), &tmp)

	// Compute sigma_{2,1}^T * sigma_{2,2}^-1 * sigma_{2,1}.
	// TODO(btracey): Should this be a method of SymDense?
	var tmp3, tmp4 mat64.Dense
	err = tmp3.SolveCholesky(&chol, sigma21)
	if err != nil {
		return nil, false
	}
	tmp4.Mul(sigma21.T(), &tmp3)

	for i := range mu1 {
		mu1[i] += tmp2.At(i, 0)
	}

	// TODO(btracey): If tmp2 can constructed with a method, then this can be
	// replaced with SubSym.
	for i := 0; i < len(unobserved); i++ {
		for j := i; j < len(unobserved); j++ {
			v := sigma11.At(i, j)
			sigma11.SetSym(i, j, v-tmp4.At(i, j))
		}
	}
	return NewNormal(mu1, &sigma11, src)
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
	n.setSigma()
	s.CopySym(n.sigma)
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
	// Compute the normalization constant
	c := -0.5*float64(dim)*logTwoPi - n.logSqrtDet

	// Compute (x-mu)'Sigma^-1 (x-mu)
	xMinusMu := make([]float64, dim)
	floats.SubTo(xMinusMu, x, n.mu)
	d := mat64.NewVector(dim, xMinusMu)
	tmp := make([]float64, dim)
	tmpVec := mat64.NewVector(dim, tmp)
	tmpVec.SolveCholeskyVec(&n.chol, d)
	return c - 0.5*floats.Dot(tmp, xMinusMu)
}

// MarginalNormal returns the marginal distribution of the given input variables.
// That is, MarginalNormal returns
//  p(x_i) = \int_{x_o} p(x_i | x_o) p(x_o) dx_o
// where x_i are the dimensions in the input, and x_o are the remaining dimensions.
// The input src is passed to the call to NewNormal.
func (n *Normal) MarginalNormal(vars []int, src *rand.Rand) (*Normal, bool) {
	newMean := make([]float64, len(vars))
	for i, v := range vars {
		newMean[i] = n.mu[v]
	}
	n.setSigma()
	var s mat64.SymDense
	s.SubsetSym(n.sigma, vars)
	return NewNormal(newMean, &s, src)
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
	tmpVec := mat64.NewVector(n.dim, tmp)
	xVec := mat64.NewVector(n.dim, x)
	xVec.MulVec(&n.lower, tmpVec)
	floats.Add(x, n.mu)
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

// setSigma computes and stores the covariance matrix of the distribution.
func (n *Normal) setSigma() {
	n.once.Do(func() {
		n.sigma = mat64.NewSymDense(n.Dim(), nil)
		n.sigma.FromCholesky(&n.chol)
	})
}
