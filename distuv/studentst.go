// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"
)

const logPi = 1.1447298858494001741 // http://oeis.org/A053510

// StudentsT implements the Student's T distribution, a 1-parameter distribution
// over the real numbers.
//
// The Student's T distribution has density function
//  Γ((ν+1)/2) / (sqrt(νπ) Γ(ν/2)) (1 + x^2/ν)^(-(ν+1)/2)
//
// The Student's T distribution approaches the standard normal as ν → ∞.
//
// For more information, see https://en.wikipedia.org/wiki/Student%27s_t-distribution.
type StudentsT struct {
	// Nu is the shape prameter of the distribution, representing the number of
	// degrees of the distribution, and one less than the number of observations
	// from a Normal distribution.
	Nu float64

	Src *rand.Rand
}

// ExKurtosis returns the excess kurtosis of the distribution.
//
// The excess Kurtosis is undefined for ν <= 2, and this returns math.NaN().
func (s StudentsT) ExKurtosis() float64 {
	if s.Nu <= 2 {
		return math.NaN()
	}
	if s.Nu <= 4 {
		return math.Inf(1)
	}
	return 6 / (s.Nu - 4)
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x.
func (s StudentsT) LogProb(x float64) float64 {
	g1, _ := math.Lgamma((s.Nu + 1) / 2)
	g2, _ := math.Lgamma(s.Nu / 2)
	return g1 - g2 - 0.5*math.Log(s.Nu) - 0.5*logPi - ((s.Nu+1)/2)*math.Log(1+x*x/s.Nu)
}

// Mean returns the mean of the probability distribution.
func (StudentsT) Mean() float64 {
	return 0
}

// Mode returns the mode of the distribution.
func (StudentsT) Mode() float64 {
	return 0
}

// NumParameters returns the number of parameters in the distribution.
func (StudentsT) NumParameters() int {
	return 1
}

// Prob computes the value of the probability density function at x.
func (s StudentsT) Prob(x float64) float64 {
	return math.Exp(s.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
func (s StudentsT) Rand() float64 {
	// http://www.math.uah.edu/stat/special/Student.html
	n := Normal{0, 1, s.Src}.Rand()
	c := Gamma{s.Nu / 2, 0.5, s.Src}.Rand()
	return n / math.Sqrt(c/s.Nu)
}

// StdDev returns the standard deviation of the probability distribution.
//
// The standard deviation is undefined for ν <= 1, and this returns math.NaN().
func (s StudentsT) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// Variance returns the variance of the probability distribution.
//
// The variance is undefined for ν <= 1, and this returns math.NaN().
func (s StudentsT) Variance() float64 {
	if s.Nu < 1 {
		return math.NaN()
	}
	if s.Nu <= 2 {
		return math.Inf(1)
	}
	return s.Nu / (s.Nu - 2)
}
