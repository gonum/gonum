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
// respect to the parameters of the distribution. The deriv slice may either
// be nil (in which case new memory is allocated), or may have length equal
// to the number of parameters.
//
// Special cases are:
//  DLogProbDParam(0, nil) == []float64{math.NaN()}
func (e Exponential) DLogProbDParam(x float64, deriv []float64) []float64 {
	nParam := e.NumParameters()
	if deriv == nil {
		deriv = make([]float64, nParam)
	}
	if len(deriv) != nParam {
		panic("dist: slice length mismatch")
	}
	if x > 0 {
		deriv[0] = 1/e.Rate - x
		return deriv
	}
	if x < 0 {
		deriv[0] = 0
		return deriv
	}
	deriv[0] = math.NaN()
	return deriv
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
	e.FitPrior(samples, weights, nil, nil)
}

// FitPrior fits the distribution with a set of priors for the sufficient
// statistics. If priorValue and priorWeights both have length 0, no prior is used.
// For the exponential distribution, there is one prior value specifying
// the prior on the sample rate and the number of samples that observed that value.
func (e *Exponential) FitPrior(samples, weights, priorValue, priorWeight []float64) (newPriorValue, newPriorWeight []float64) {
	lenValue := len(priorValue)
	lenPriorWeight := len(priorWeight)
	if lenValue != lenPriorWeight {
		panic("exponential: mismatch in prior lengths")
	}
	if lenValue > 1 {
		panic("exponential: too many prior values")
	}
	prior := true
	if lenValue == 0 || lenPriorWeight == 0 {
		if lenValue == 0 && lenPriorWeight == 0 {
			prior = false
		} else if lenValue == 0 && lenPriorWeight != 0 {
			panic("exponential: prior weight provided but not the value")
		} else {
			panic("exponential: prior value provided but not the weight")
		}
	}

	lenSamp := len(samples)
	lenWeight := len(weights)
	if lenWeight != 0 && lenSamp != lenWeight {
		panic("exponential: length of samples and length of weights does not match")
	}

	var sumWeights float64
	if lenWeight == 0 {
		sumWeights = float64(lenSamp)
	} else {
		sumWeights = floats.Sum(weights)
	}

	sampleMean := stat.Mean(samples, weights)

	totalSum := sampleMean * sumWeights
	totalWeights := sumWeights
	if prior {
		totalWeights += priorWeight[0]
		totalSum += (1 / priorValue[0]) * priorWeight[0]
	}

	e.Rate = totalWeights / totalSum

	newPriorValue = []float64{e.Rate}
	newPriorWeight = []float64{totalWeights}
	return newPriorValue, newPriorWeight
}

// LogProb computes the natural logarithm of the value of the probability density function at x.
func (e Exponential) LogProb(x float64) float64 {
	if x < 0 {
		return math.Inf(-1)
	}
	return math.Log(e.Rate) - e.Rate*x
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

// Parameters gets the parameters of the distribution. Panics if the length of
// the input slice is non-zero and not equal to the number of parameters
// This returns a slice with length 1 containing the rate of the distribution.
func (e Exponential) Parameters(s []float64) []float64 {
	nParam := e.NumParameters()
	if s == nil {
		s = make([]float64, nParam)
	}
	if len(s) != nParam {
		panic("exponential: improper parameter length")
	}
	s[0] = e.Rate
	return s
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

// SetParameters gets the parameters of the distribution. Panics if the length of
// the input slice is not equal to the number of parameters.
// This sets the rate of the distribution to the first element of the slice.
func (e *Exponential) SetParameters(s []float64) {
	if len(s) != e.NumParameters() {
		panic("exponential: incorrect number of parameters to set")
	}
	e.Rate = s[0]
}

// Skewness returns the skewness of the distribution.
func (Exponential) Skewness() float64 {
	return 2
}

// StdDev returns the standard deviation of the probability distribution.
func (e Exponential) StdDev() float64 {
	return 1 / e.Rate
}

// Survival returns the survival function (complementary CDF) at x.
func (e Exponential) Survival(x float64) float64 {
	if x < 0 {
		return 1
	}
	return math.Exp(-e.Rate * x)
}

// Variance returns the variance of the probability distribution.
func (e Exponential) Variance() float64 {
	return 1 / (e.Rate * e.Rate)
}
