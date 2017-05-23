// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/matrix/mat64"
	"gonum.org/v1/gonum/stat"
)

// Bhattacharyya is a type for computing the Bhattacharyya distance between
// probability distributions.
//
// The Battachara distance is defined as
//  D_B = -ln(BC(l,r))
//  BC = \int_x (p(x)q(x))^(1/2) dx
// Where BC is known as the Bhattacharyya coefficient.
// The Bhattacharyya distance is related to the Hellinger distance by
//  H = sqrt(1-BC)
// For more information, see
//  https://en.wikipedia.org/wiki/Bhattacharyya_distance
type Bhattacharyya struct{}

// DistNormal computes the Bhattacharyya distance between normal distributions l and r.
// The dimensions of the input distributions must match or DistNormal will panic.
//
// For Normal distributions, the Bhattacharyya distance is
//  Σ = (Σ_l + Σ_r)/2
//  D_B = (1/8)*(μ_l - μ_r)^T*Σ^-1*(μ_l - μ_r) + (1/2)*ln(det(Σ)/(det(Σ_l)*det(Σ_r))^(1/2))
func (Bhattacharyya) DistNormal(l, r *Normal) float64 {
	dim := l.Dim()
	if dim != r.Dim() {
		panic(badSizeMismatch)
	}

	var sigma mat64.SymDense
	sigma.AddSym(&l.sigma, &r.sigma)
	sigma.ScaleSym(0.5, &sigma)

	var chol mat64.Cholesky
	chol.Factorize(&sigma)

	mahalanobis := stat.Mahalanobis(mat64.NewVector(dim, l.mu), mat64.NewVector(dim, r.mu), &chol)
	mahalanobisSq := mahalanobis * mahalanobis

	dl := l.chol.LogDet()
	dr := r.chol.LogDet()
	ds := chol.LogDet()

	return 0.125*mahalanobisSq + 0.5*ds - 0.25*dl - 0.25*dr
}

// CrossEntropy is a type for computing the cross-entropy between probability
// distributions.
//
// The cross-entropy is defined as
//  - \int_x l(x) log(r(x)) dx = KL(l || r) + H(l)
// where KL is the Kullback-Leibler divergence and H is the entropy.
// For more information, see
//  https://en.wikipedia.org/wiki/Cross_entropy
type CrossEntropy struct{}

// DistNormal returns the cross-entropy between normal distributions l and r.
// The dimensions of the input distributions must match or DistNormal will panic.
func (CrossEntropy) DistNormal(l, r *Normal) float64 {
	if l.Dim() != r.Dim() {
		panic(badSizeMismatch)
	}
	kl := KullbackLeibler{}.DistNormal(l, r)
	return kl + l.Entropy()
}

// Hellinger is a type for computing the Hellinger distance between probability
// distributions.
//
// The Hellinger distance is defined as
//  H^2(l,r) = 1/2 * int_x (\sqrt(l(x)) - \sqrt(r(x)))^2 dx
// and is bounded between 0 and 1.
// The Hellinger distance is related to the Bhattacharyya distance by
//  H^2 = 1 - exp(-Db)
// For more information, see
//  https://en.wikipedia.org/wiki/Hellinger_distance
type Hellinger struct{}

// DistNormal returns the Hellinger distance between normal distributions l and r.
// The dimensions of the input distributions must match or DistNormal will panic.
//
// See the documentation of Bhattacharyya.DistNormal for the formula for Normal
// distributions.
func (Hellinger) DistNormal(l, r *Normal) float64 {
	if l.Dim() != r.Dim() {
		panic(badSizeMismatch)
	}
	db := Bhattacharyya{}.DistNormal(l, r)
	bc := math.Exp(-db)
	return math.Sqrt(1 - bc)
}

// KullbackLiebler is a type for computing the Kullback-Leibler divergence from l to r.
// The dimensions of the input distributions must match or the function will panic.
//
// The Kullback-Liebler divergence is defined as
//  D_KL(l || r ) = \int_x p(x) log(p(x)/q(x)) dx
// Note that the Kullback-Liebler divergence is not symmetric with respect to
// the order of the input arguments.
type KullbackLeibler struct{}

// DistNormal returns the KullbackLeibler distance between normal distributions l and r.
// The dimensions of the input distributions must match or DistNormal will panic.
//
// For two normal distributions, the KL divergence is computed as
//   D_KL(l || r) = 0.5*[ln(|Σ_r|) - ln(|Σ_l|) + (μ_l - μ_r)^T*Σ_r^-1*(μ_l - μ_r) + tr(Σ_r^-1*Σ_l)-d]
func (KullbackLeibler) DistNormal(l, r *Normal) float64 {
	dim := l.Dim()
	if dim != r.Dim() {
		panic(badSizeMismatch)
	}

	mahalanobis := stat.Mahalanobis(mat64.NewVector(dim, l.mu), mat64.NewVector(dim, r.mu), &r.chol)
	mahalanobisSq := mahalanobis * mahalanobis

	// TODO(btracey): Optimize where there is a SolveCholeskySym
	// TODO(btracey): There may be a more efficient way to just compute the trace
	// Compute tr(Σ_r^-1*Σ_l) using the fact that Σ_l = U^T * U
	var u mat64.TriDense
	u.UFromCholesky(&l.chol)
	var m mat64.Dense
	err := m.SolveCholesky(&r.chol, u.T())
	if err != nil {
		return math.NaN()
	}
	m.Mul(&m, &u)
	tr := mat64.Trace(&m)

	return r.logSqrtDet - l.logSqrtDet + 0.5*(mahalanobisSq+tr-float64(l.dim))
}

// Wasserstein is a type for computing the Wasserstein distance between two
// probability distributions.
//
// The Wasserstein distance is defined as
//  W(l,r) := inf 𝔼(||X-Y||_2^2)^1/2
// For more information, see
//  https://en.wikipedia.org/wiki/Wasserstein_metric
type Wasserstein struct{}

// DistNormal returns the Wasserstein distance between normal distributions l and r.
// The dimensions of the input distributions must match or DistNormal will panic.
//
// The Wasserstein distance for Normal distributions is
//  d^2 = ||m_l - m_r||_2^2 + Tr(Σ_l + Σ_r - 2(Σ_l^(1/2)*Σ_r*Σ_l^(1/2))^(1/2))
// For more information, see
//  http://djalil.chafai.net/blog/2010/04/30/wasserstein-distance-between-two-gaussians/
func (Wasserstein) DistNormal(l, r *Normal) float64 {
	dim := l.Dim()
	if dim != r.Dim() {
		panic(badSizeMismatch)
	}

	d := floats.Distance(l.mu, r.mu, 2)
	d = d * d

	// Compute Σ_l^(1/2)
	var ssl mat64.SymDense
	ssl.PowPSD(&l.sigma, 0.5)
	// Compute Σ_l^(1/2)*Σ_r*Σ_l^(1/2)
	var mean mat64.Dense
	mean.Mul(&ssl, &r.sigma)
	mean.Mul(&mean, &ssl)

	// Reinterpret as symdense, and take Σ^(1/2)
	meanSym := mat64.NewSymDense(dim, mean.RawMatrix().Data)
	ssl.PowPSD(meanSym, 0.5)

	tr := mat64.Trace(&r.sigma)
	tl := mat64.Trace(&l.sigma)
	tm := mat64.Trace(&ssl)

	return d + tl + tr - 2*tm
}
