// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"
)

// Beta implements the Beta distribution, a two-parameter continuous distribution
// with support over the positive real numbers.
//
// The beta distribution has density function
//  x^(α-1) * (1-x)^(β-1) * Γ(α+β) / (Γ(α)*Γ(β))
//
// For more information, see https://en.wikipedia.org/wiki/Beta_distribution
type Beta struct {
	// Alpha is the left shape parameter of the distribution. Alpha must be greater
	// than 0.
	Alpha float64
	// Beta is the right shape parameter of the distribution. Beta must be greater
	// than 0.
	Beta float64

	Source *rand.Rand
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (be Beta) ExKurtosis() float64 {
	a := be.Alpha
	b := be.Beta
	num := 6 * (math.Pow(a-b, 2)*(a+b+1) - a*b*(a+b+2))
	den := a * b * (a + b + 2) * (a + b + 3)
	return num / den
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x.
func (be Beta) LogProb(x float64) float64 {
	if x < 0 || x > 1 {
		return math.NaN()
	}
	a := be.Alpha
	b := be.Beta
	lgab, sign := math.Lgamma(a + b)
	lgab = lgab * float64(sign)
	lga, sign := math.Lgamma(a)
	lga = lga * float64(sign)
	lgb, sign := math.Lgamma(b)
	lgb = lgb * float64(sign)
	return lgab - lga - lgb + (a-1)*math.Log(x) + (b-1)*math.Log(1-x)
}

// Mean returns the mean of the probability distribution.
func (be Beta) Mean() float64 {
	return be.Alpha / (be.Alpha + be.Beta)
}

// Mode returns the mode of the distribution.
//
// The mode is NaN in the special case where one parameter is less than or
// equal to 1.
func (be Beta) Mode() float64 {
	if be.Alpha <= 1 || be.Beta <= 1 {
		return math.NaN()
	}
	return (be.Alpha - 1) / (be.Alpha + be.Beta - 2)
}

// NumParameters returns the number of parameters in the distribution.
func (be Beta) NumParameters() int {
	return 2
}

// Prob computes the value of the probability density function at x.
func (be Beta) Prob(x float64) float64 {
	return math.Exp(be.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
func (be Beta) Rand() float64 {
	ga := Gamma{Alpha: be.Alpha, Beta: 1, Source: be.Source}
	gb := Gamma{Alpha: be.Beta, Beta: 1, Source: be.Source}
	gaRand := ga.Rand()
	gbRand := gb.Rand()
	return gaRand / (gaRand + gbRand)
}

// StdDev returns the standard deviation of the probability distribution.
func (be Beta) StdDev() float64 {
	return math.Sqrt(be.Variance())
}

// Variance returns the variance of the probability distribution.
func (be Beta) Variance() float64 {
	return be.Alpha * be.Beta / (math.Pow(be.Alpha+be.Beta, 2) * (be.Alpha + be.Beta + 1))
}
