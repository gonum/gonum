// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// UnitNormal is an instantiation of the standard normal distribution
var UnitNormal = Normal{Mu: 0, Sigma: 1}

// Normal respresents a normal (Gaussian) distribution (https://en.wikipedia.org/wiki/Normal_distribution).
type Normal struct {
	Mu     float64 // Mean of the normal distribution
	Sigma  float64 // Standard deviation of the normal distribution
	Source *rand.Rand

	// Needs to be Mu and Sigma and not Mean and StdDev because Normal has functions
	// Mean and StdDev
}

// CDF computes the value of the cumulative density function at x.
func (n Normal) CDF(x float64) float64 {
	return 0.5 * (1 + math.Erf((x-n.Mu)/(n.Sigma*math.Sqrt2)))
}

// ConjugateUpdate updates the parameters of the distribution from the sufficient
// statistics of a set of samples. The sufficient statistics, suffStat, have been
// observed with nSamples observations. The prior values of the distribution are those
// currently in the distribution, and have been observed with priorStrength samples.
//
// For the normal distribution, the sufficient statistics are the mean and
// uncorrected standard deviation of the samples.
// The prior is having seen strength[0] samples with mean Normal.Mu
// and strength[1] samples with standard deviation Normal.Sigma. As a result of
// this function, Normal.Mu and Normal.Sigma are updated based on the weighted
// samples, and strength is modified to include the new number of samples observed.
//
// This function panics if len(suffStat) != 2 or len(priorStrength) != 2.
func (n *Normal) ConjugateUpdate(suffStat []float64, nSamples float64, priorStrength []float64) {

	// TODO: Support prior strength with math.Inf(1) to allow updating with
	// a known mean/standard deviation

	totalMeanSamples := nSamples + priorStrength[0]
	totalSum := suffStat[0]*nSamples + n.Mu*priorStrength[0]

	totalVarianceSamples := nSamples + priorStrength[1]
	// sample variance
	totalVariance := nSamples * suffStat[1] * suffStat[1]
	// add prior variance
	totalVariance += priorStrength[1] * n.Sigma * n.Sigma
	// add cross variance from the difference of the means
	meanDiff := (suffStat[0] - n.Mu)
	totalVariance += priorStrength[0] * nSamples * meanDiff * meanDiff / totalMeanSamples

	n.Mu = totalSum / totalMeanSamples
	n.Sigma = math.Sqrt(totalVariance / totalVarianceSamples)
	floats.AddConst(nSamples, priorStrength)
}

// Entropy returns the differential entropy of the distribution.
func (n Normal) Entropy() float64 {
	return 0.5 * (log2Pi + 1 + 2*math.Log(n.Sigma))
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (Normal) ExKurtosis() float64 {
	return 0
}

// Fit sets the parameters of the probability distribution from the
// data samples x with relative weights w. If weights is nil, then all the weights
// are 1. If weights is not nil, then the len(weights) must equal len(samples).
func (n *Normal) Fit(samples, weights []float64) {
	suffStat := make([]float64, 1)
	nSamples := n.SuffStat(samples, weights, suffStat)
	n.ConjugateUpdate(suffStat, nSamples, []float64{0, 0})
}

// LogProb computes the natural logarithm of the value of the probability density function at x.
func (n Normal) LogProb(x float64) float64 {
	return negLogRoot2Pi - math.Log(n.Sigma) - (x-n.Mu)*(x-n.Mu)/(2*n.Sigma*n.Sigma)
}

// Mean returns the mean of the probability distribution.
func (n Normal) Mean() float64 {
	return n.Mu
}

// Median returns the median of the normal distribution.
func (n Normal) Median() float64 {
	return n.Mu
}

// Mode returns the mode of the normal distribution.
func (n Normal) Mode() float64 {
	return n.Mu
}

// NumParameters returns the number of parameters in the distribution.
func (Normal) NumParameters() int {
	return 2
}

func (Normal) NumSuffStat() int {
	return 2
}

// Prob computes the value of the probability density function at x.
func (n Normal) Prob(x float64) float64 {
	return math.Exp(n.LogProb(x))
}

// Quantile returns the inverse of the cumulative probability distribution.
func (n Normal) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic(badPercentile)
	}
	return n.Mu + n.Sigma*zQuantile(p)
}

// Rand returns a random sample drawn from the distribution.
func (n Normal) Rand() float64 {
	var rnd float64
	if n.Source == nil {
		rnd = rand.NormFloat64()
	} else {
		rnd = n.Source.NormFloat64()
	}
	return rnd*n.Sigma + n.Mu
}

// Score returns the score function with respect to the parameters of the
// distribution at the input location x. The score function is the derivative
// of the log-likelihood at x with respect to the parameters
//  (∂/∂θ) log(p(x;θ))
// If deriv is non-nil, len(deriv) must equal the number of parameters otherwise
// Score will panic, and the derivative is stored in-place into deriv. If deriv
// is nil a new slice will be allocated and returned.
//
// The order is [∂LogProb / ∂Mu, ∂LogProb / ∂Sigma].
//
// For more information, see https://en.wikipedia.org/wiki/Score_%28statistics%29.
func (n Normal) Score(deriv []float64, x float64) []float64 {
	if deriv == nil {
		deriv = make([]float64, n.NumParameters())
	}
	if len(deriv) != n.NumParameters() {
		panic(badLength)
	}
	deriv[0] = (x - n.Mu) / (n.Sigma * n.Sigma)
	deriv[1] = 1 / n.Sigma * (-1 + ((x-n.Mu)/n.Sigma)*((x-n.Mu)/n.Sigma))
	return deriv
}

// ScoreInput returns the score function with respect to the input of the
// distribution at the input location specified by x. The score function is the
// derivative of the log-likelihood
//  (d/dx) log(p(x)) .
func (n Normal) ScoreInput(x float64) float64 {
	return -(1 / (2 * n.Sigma * n.Sigma)) * 2 * (x - n.Mu)
}

// Skewness returns the skewness of the distribution.
func (Normal) Skewness() float64 {
	return 0
}

// StdDev returns the standard deviation of the probability distribution.
func (n Normal) StdDev() float64 {
	return n.Sigma
}

// SuffStat computes the sufficient statistics of a set of samples to update
// the distribution. The sufficient statistics are stored in place, and the
// effective number of samples are returned.
//
// The normal distribution has two sufficient statistics, the mean of the samples
// and the standard deviation of the samples.
//
// If weights is nil, the weights are assumed to be 1, otherwise panics if
// len(samples) != len(weights). Panics if len(suffStat) != 2.
func (Normal) SuffStat(samples, weights, suffStat []float64) (nSamples float64) {
	lenSamp := len(samples)
	if len(weights) != 0 && len(samples) != len(weights) {
		panic(badLength)
	}
	if len(suffStat) != 2 {
		panic(badSuffStat)
	}

	if len(weights) == 0 {
		nSamples = float64(lenSamp)
	} else {
		nSamples = floats.Sum(weights)
	}

	mean := stat.Mean(samples, weights)
	suffStat[0] = mean

	// Use Moment and not StdDev because we want it to be uncorrected
	variance := stat.MomentAbout(2, samples, mean, weights)
	suffStat[1] = math.Sqrt(variance)
	return nSamples
}

// Survival returns the survival function (complementary CDF) at x.
func (n Normal) Survival(x float64) float64 {
	return 0.5 * (1 - math.Erf((x-n.Mu)/(n.Sigma*math.Sqrt2)))
}

// setParameters modifies the parameters of the distribution.
func (n *Normal) setParameters(p []Parameter) {
	if len(p) != n.NumParameters() {
		panic("normal: incorrect number of parameters to set")
	}
	if p[0].Name != "Mu" {
		panic("normal: " + panicNameMismatch)
	}
	if p[1].Name != "Sigma" {
		panic("normal: " + panicNameMismatch)
	}
	n.Mu = p[0].Value
	n.Sigma = p[1].Value
}

// Variance returns the variance of the probability distribution.
func (n Normal) Variance() float64 {
	return n.Sigma * n.Sigma
}

// TODO: Is the right way to compute inverf?
// It seems to me like the precision is not high enough, but I don't
// know the correct version. It would be nice if this were built into the
// math package in the standard library (issue 6359)

/*
Copyright (c) 2012 The Probab Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
* Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

var (
	zQuantSmallA = []float64{3.387132872796366608, 133.14166789178437745, 1971.5909503065514427, 13731.693765509461125, 45921.953931549871457, 67265.770927008700853, 33430.575583588128105, 2509.0809287301226727}
	zQuantSmallB = []float64{1.0, 42.313330701600911252, 687.1870074920579083, 5394.1960214247511077, 21213.794301586595867, 39307.89580009271061, 28729.085735721942674, 5226.495278852854561}
	zQuantInterA = []float64{1.42343711074968357734, 4.6303378461565452959, 5.7694972214606914055, 3.64784832476320460504, 1.27045825245236838258, 0.24178072517745061177, 0.0227238449892691845833, 7.7454501427834140764e-4}
	zQuantInterB = []float64{1.0, 2.05319162663775882187, 1.6763848301838038494, 0.68976733498510000455, 0.14810397642748007459, 0.0151986665636164571966, 5.475938084995344946e-4, 1.05075007164441684324e-9}
	zQuantTailA  = []float64{6.6579046435011037772, 5.4637849111641143699, 1.7848265399172913358, 0.29656057182850489123, 0.026532189526576123093, 0.0012426609473880784386, 2.71155556874348757815e-5, 2.01033439929228813265e-7}
	zQuantTailB  = []float64{1.0, 0.59983220655588793769, 0.13692988092273580531, 0.0148753612908506148525, 7.868691311456132591e-4, 1.8463183175100546818e-5, 1.4215117583164458887e-7, 2.04426310338993978564e-15}
)

func rateval(a []float64, na int64, b []float64, nb int64, x float64) float64 {
	var (
		u, v, r float64
	)
	u = a[na-1]

	for i := na - 1; i > 0; i-- {
		u = x*u + a[i-1]
	}

	v = b[nb-1]

	for j := nb - 1; j > 0; j-- {
		v = x*v + b[j-1]
	}

	r = u / v

	return r
}

func zQuantSmall(q float64) float64 {
	r := 0.180625 - q*q
	return q * rateval(zQuantSmallA, 8, zQuantSmallB, 8, r)
}

func zQuantIntermediate(r float64) float64 {
	return rateval(zQuantInterA, 8, zQuantInterB, 8, (r - 1.6))
}

func zQuantTail(r float64) float64 {
	return rateval(zQuantTailA, 8, zQuantTailB, 8, (r - 5.0))
}

// Compute the quantile in normalized units
func zQuantile(p float64) float64 {
	switch {
	case p == 1.0:
		return math.Inf(1)
	case p == 0.0:
		return math.Inf(-1)
	}
	var r, x, pp, dp float64
	dp = p - 0.5
	if math.Abs(dp) <= 0.425 {
		return zQuantSmall(dp)
	}
	if p < 0.5 {
		pp = p
	} else {
		pp = 1.0 - p
	}
	r = math.Sqrt(-math.Log(pp))
	if r <= 5.0 {
		x = zQuantIntermediate(r)
	} else {
		x = zQuantTail(r)
	}
	if p < 0.5 {
		return -x
	}
	return x
}

// parameters returns the parameters of the distribution.
func (n Normal) parameters(p []Parameter) []Parameter {
	nParam := n.NumParameters()
	if p == nil {
		p = make([]Parameter, nParam)
	} else if len(p) != nParam {
		panic("normal: improper parameter length")
	}
	p[0].Name = "Mu"
	p[0].Value = n.Mu
	p[1].Name = "Sigma"
	p[1].Value = n.Sigma
	return p
}
