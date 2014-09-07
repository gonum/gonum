// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dist

import (
	"math"
	"math/rand"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// Exponential represents the exponential distribution (https://en.wikipedia.org/wiki/Exponential_distribution).
type Exponential struct {
	Rate   float64
	Source *rand.Rand
}

// CDF computes the value of the cumulative density function at x.
func (e Exponential) CDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	return 1 - math.Exp(-e.Rate*x)
}

// ConjugateUpdate updates the parameters of the distribution from the sufficient
// statistics of a set of samples. The sufficient statistics, suffStat, have been
// observed with nSamples observations. The prior values of the distribution are those
// currently in the distribution, and have been observed with priorStrength samples.
//
// For the exponential distribution, the sufficient statistic is the inverse of
// the mean of the samples.
// The prior is having seen priorStrength[0] samples with inverse mean Exponential.Rate
// As a result of this function, Exponential.Rate is updated based on the weighted
// samples, and priorStrength is modified to include the new number of samples observed.
//
// This function panics if len(suffStat) != 1 or len(priorStrength) != 1.
func (e *Exponential) ConjugateUpdate(suffStat []float64, nSamples float64, priorStrength []float64) {
	if len(suffStat) != 1 {
		panic("exponential: incorrect suffStat length")
	}
	if len(priorStrength) != 1 {
		panic("exponential: incorrect priorStrength length")
	}

	totalSamples := nSamples + priorStrength[0]

	totalSum := nSamples / suffStat[0]
	if !(priorStrength[0] == 0) {
		totalSum += priorStrength[0] / e.Rate
	}
	e.Rate = totalSamples / totalSum
	priorStrength[0] = totalSamples
}

// DLogProbDX returns the derivative of the log of the probability with
// respect to the input x.
//
// Special cases are:
//  DLogProbDX(0) = NaN
func (e Exponential) DLogProbDX(x float64) float64 {
	if x > 0 {
		return -e.Rate
	}
	if x < 0 {
		return 0
	}
	return math.NaN()
}

// DLogProbDParam returns the derivative of the log of the probability with
// respect to the parameters of the distribution. The deriv slice must have length
// equal to the number of parameters of the distribution.
//
// The order is ∂LogProb / ∂Rate
//
// Special cases are:
//  The derivative at 0 is NaN.
func (e Exponential) DLogProbDParam(x float64, deriv []float64) {
	if len(deriv) != e.NumParameters() {
		panic("dist: slice length mismatch")
	}
	if x > 0 {
		deriv[0] = 1/e.Rate - x
		return
	}
	if x < 0 {
		deriv[0] = 0
		return
	}
	deriv[0] = math.NaN()
	return
}

// Entropy returns the entropy of the distribution.
func (e Exponential) Entropy() float64 {
	return 1 - math.Log(e.Rate)
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (Exponential) ExKurtosis() float64 {
	return 6
}

// Fit sets the parameters of the probability distribution from the
// data samples x with relative weights w.
// If weights is nil, then all the weights are 1.
// If weights is not nil, then the len(weights) must equal len(samples).
func (e *Exponential) Fit(samples, weights []float64) {
	suffStat := make([]float64, 1)
	nSamples := e.SuffStat(samples, weights, suffStat)
	e.ConjugateUpdate(suffStat, nSamples, []float64{0})
}

// LogProb computes the natural logarithm of the value of the probability density function at x.
func (e Exponential) LogProb(x float64) float64 {
	if x < 0 {
		return math.Inf(-1)
	}
	return math.Log(e.Rate) - e.Rate*x
}

// MarshalSlice gets the parameters of the distribution.
// Sets Rate to the first element of the slice. Panics if the length of
// the input slice is not equal to the number of parameters.
func (e Exponential) MarshalSlice(s []float64) {
	nParam := e.NumParameters()
	if len(s) != nParam {
		panic("exponential: improper parameter length")
	}
	s[0] = e.Rate
	return
}

// Mean returns the mean of the probability distribution.
func (e Exponential) Mean() float64 {
	return 1 / e.Rate
}

// Median returns the median of the probability distribution.
func (e Exponential) Median() float64 {
	return math.Ln2 / e.Rate
}

// Mode returns the mode of the probability distribution.
func (Exponential) Mode() float64 {
	return 0
}

// NumParameters returns the number of parameters in the distribution.
func (Exponential) NumParameters() int {
	return 1
}

func (Exponential) NumSuffStat() int {
	return 1
}

// ExponentialMap is the parameter mapping for the Uniform distribution.
var ExponentialMap = map[string]int{"Rate": 0}

// ParameterMap returns a mapping from fields of the distribution to elements
// of the marshaled slice. Do not edit this variable.
func (e Exponential) ParameterMap() map[string]int {
	return ExponentialMap
}

// Prob computes the value of the probability density function at x.
func (e Exponential) Prob(x float64) float64 {
	return math.Exp(e.LogProb(x))
}

// Quantile returns the inverse of the cumulative probability distribution.
func (e Exponential) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic("dist: percentile out of bounds")
	}
	return -math.Log(1-p) / e.Rate
}

// Rand returns a random sample drawn from the distribution.
func (e Exponential) Rand() float64 {
	var rnd float64
	if e.Source == nil {
		rnd = rand.ExpFloat64()
	} else {
		rnd = e.Source.ExpFloat64()
	}
	return rnd / e.Rate
}

// Skewness returns the skewness of the distribution.
func (Exponential) Skewness() float64 {
	return 2
}

// StdDev returns the standard deviation of the probability distribution.
func (e Exponential) StdDev() float64 {
	return 1 / e.Rate
}

// SuffStat computes the sufficient statistics of set of samples to update
// the distribution. The sufficient statistics are stored in place, and the
// effective number of samples are returned.
//
// The exponential distribution has one sufficient statistic, the average rate
// of the samples.
//
// If weights is nil, the weights are assumed to be 1, otherwise panics if
// len(samples) != len(weights). Panics if len(suffStat) != 1.
func (Exponential) SuffStat(samples, weights, suffStat []float64) (nSamples float64) {
	if len(weights) != 0 && len(samples) != len(weights) {
		panic("dist: slice size mismatch")
	}

	if len(suffStat) != 1 {
		panic("exponential: wrong suffStat length")
	}

	if len(weights) == 0 {
		nSamples = float64(len(samples))
	} else {
		nSamples = floats.Sum(weights)
	}

	mean := stat.Mean(samples, weights)
	suffStat[0] = 1 / mean
	return nSamples
}

// Survival returns the survival function (complementary CDF) at x.
func (e Exponential) Survival(x float64) float64 {
	if x < 0 {
		return 1
	}
	return math.Exp(-e.Rate * x)
}

// UnmarshalSlice sets the parameters of the distribution.
// This sets the rate of the distribution to the first element of the slice.
// Panics if the length of the input slice is not equal to the number of parameters.
func (e *Exponential) UnmarshalSlice(s []float64) {
	if len(s) != e.NumParameters() {
		panic("exponential: incorrect number of parameters to set")
	}
	e.Rate = s[0]
}

// Variance returns the variance of the probability distribution.
func (e Exponential) Variance() float64 {
	return 1 / (e.Rate * e.Rate)
}
