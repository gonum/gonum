// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mathext"
	"gonum.org/v1/gonum/stat/combin"
)

// Binomial represents a random variable that characterizes characterizes the
// number of successes in a sequence of independent trials. It has two
// parameters: n, the number of trials, and p, the probability of success in an
// individual trial. The value of P must be between 0 and 1 and N must be
// non-negative More information at https://en.wikipedia.org/wiki/Binomial_distribution
type Binomial struct {
	P      float64
	N      int
	Source *rand.Rand
}

// Mean returns the mean of the probability distribution.
func (b Binomial) Mean() float64 {
	return b.P * float64(b.N)
}

// Prob computes the value of the probability distribution at x.
func (b Binomial) Prob(x float64) float64 {
	xi := int(x)
	return float64(combin.Binomial(b.N, xi)) *
		math.Pow(b.P, x) *
		math.Pow(1-b.P, float64(b.N-xi))
}

// CDF computes the value of the cumulative density function at x.
func (b Binomial) CDF(x float64) float64 {
	return mathext.RegIncBeta(float64(b.N)-x, x+1, 1-b.P)
}

// Entropy returns the entropy of the distribution.
func (b Binomial) Entropy() float64 {
	if b.P == 0 || b.P == 1 || b.N == 0 {
		return 0
	}
	p0 := 1 - b.P
	lg := math.Log(b.P / p0)
	lp := float64(b.N) * math.Log(p0)
	s := math.Exp(lp) * lp
	for k := 0; k < b.N; k++ {
		lp += math.Log(float64(b.N-k)/float64(k+1)) + lg
		s += math.Exp(lp) * lp
	}
	return -s
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (b Binomial) ExKurtosis() float64 {
	u := b.P * (1 - b.P)
	return (1 - 6*u) / (float64(b.N) * u)
}

// LogProb computes the natural logarithm of the value of the probability density function at x.
func (b Binomial) LogProb(x float64) float64 {
	return math.Log(b.Prob(x))
}

// Median returns the median of the probability distribution.
func (b Binomial) Median() float64 {
	return math.Ceil(float64(b.N) * b.P)
}

// NumParameters returns the number of parameters in the distribution.
func (Binomial) NumParameters() int {
	return 2
}

// Skewness returns the skewness of the distribution.
func (b Binomial) Skewness() float64 {
	p0 := 1 - b.P
	return (p0 - b.P) / math.Sqrt(float64(b.N)*p0*b.P)
}

// StdDev returns the standard deviation of the probability distribution.
func (b Binomial) StdDev() float64 {
	return math.Sqrt(b.Variance())
}

// Survival returns the survival function (complementary CDF) at x.
func (b Binomial) Survival(x float64) float64 {
	return 1 - b.CDF(x)
}

// Variance returns the variance of the probability distribution.
func (b Binomial) Variance() float64 {
	return float64(b.N) * b.P * (1 - b.P)
}

// Rand returns a random sample drawn from the distribution.
func (b Binomial) Rand() float64 {
	// TODO(sglyon): right now we just apply n independent Bernoulli...
	out := 0.0
	bern := Bernoulli{P: b.P, Source: b.Source}
	for i := 0; i < b.N; i++ {
		out += bern.Rand()
	}
	return out
}

// Quantile returns the inverse of the cumulative probability distribution.
func (b Binomial) Quantile(p float64) float64 {
	// TODO(sglyon): need to implement this
	panic("Not implemented!!!")
}
