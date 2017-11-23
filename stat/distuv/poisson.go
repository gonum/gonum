// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mathext"
)

// Poisson implements the Poisson distribution, a discrete probability distribution
// that expresses the probability of a given number of events occurring in a fixed
// interval of time and/or space.
// The poisson distribution has density function:
//  f(k) = λ^k / k! e^(-λ)
// For more information, see https://en.wikipedia.org/wiki/Poisson_distribution.
type Poisson struct {
	// Lambda is the average number of events in an interval.
	// Lambda must be greater than 0.
	Lambda float64

	Source *rand.Rand
}

// CDF computes the value of the cumulative distribution function at x.
func (p Poisson) CDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	return mathext.GammaIncComp(math.Floor(x+1), p.Lambda)
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (p Poisson) ExKurtosis() float64 {
	return 1 / p.Lambda
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x.
func (p Poisson) LogProb(x float64) float64 {
	if x < 0 || math.Floor(x) != x {
		return math.Inf(-1)
	}
	lg, _ := math.Lgamma(math.Floor(x) + 1)
	return x*math.Log(p.Lambda) - p.Lambda - lg
}

// Mean returns the mean of the probability distribution.
func (p Poisson) Mean() float64 {
	return p.Lambda
}

// NumParameters returns the number of parameters in the distribution.
func (Poisson) NumParameters() int {
	return 1
}

// Prob computes the value of the probability density function at x.
func (p Poisson) Prob(x float64) float64 {
	return math.Exp(p.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
func (p Poisson) Rand() float64 {
	// poisson generator based upon the multiplication of
	// uniform random variates.
	// see:
	//  Non-Uniform Random Variate Generation,
	//  Luc Devroye (p504)
	//  http://www.eirene.de/Devroye.pdf
	x := 0.0
	prod := 1.0
	exp := math.Exp(-p.Lambda)
	rnd := rand.Float64
	if p.Source != nil {
		rnd = p.Source.Float64
	}

	for {
		prod *= rnd()
		if prod <= exp {
			return x
		}
		x++
	}
}

// Skewness returns the skewness of the distribution.
func (p Poisson) Skewness() float64 {
	return 1 / math.Sqrt(p.Lambda)
}

// StdDev returns the standard deviation of the probability distribution.
func (p Poisson) StdDev() float64 {
	return math.Sqrt(p.Variance())
}

// Survival returns the survival function (complementary CDF) at x.
func (p Poisson) Survival(x float64) float64 {
	return 1 - p.CDF(x)
}

// Variance returns the variance of the probability distribution.
func (p Poisson) Variance() float64 {
	return p.Lambda
}
