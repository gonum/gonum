// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand"

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
