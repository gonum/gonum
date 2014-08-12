// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dist

import (
	"math"
	"math/rand"
	"sort"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// Laplace represents the Laplace distribution (https://en.wikipedia.org/wiki/Laplace_distribution).
type Laplace struct {
	Mu     float64 // Mean of the Laplace distribution
	Scale  float64 // Scale of the Laplace distribution
	Source *rand.Rand
}

// CDF computes the value of the cumulative density function at x.
func (l Laplace) CDF(x float64) float64 {
	if x < l.Mu {
		return 0.5 * math.Exp((x-l.Mu)/l.Scale)
	}
	return 1 - 0.5*math.Exp(-(x-l.Mu)/l.Scale)
}

// DLogProbDX returns the derivative of the log of the probability with
// respect to the input x. Returns 0 if x == l.Mu.
func (l Laplace) DLogProbDX(x float64) float64 {
	diff := x - l.Mu
	if diff == 0 {
		return 0
	}
	if diff > 0 {
		return -1 / l.Scale
	}
	return 1 / l.Scale
}

// DLogProbDParam returns the derivative of the log of the probability at x with
// respect to the parameters of the distribution. If deriv is nil, new memory
// will be allocated, otherwise len(deriv) must equal NumParameters
//  The first is ∂LogProb / ∂μ
//  Second is ∂LogProb / ∂b
func (l Laplace) DLogProbDParam(x float64, deriv []float64) []float64 {
	if deriv == nil {
		deriv = make([]float64, l.NumParameters())
	}
	if len(deriv) != l.NumParameters() {
		panic("dist: slice length mismatch")
	}
	diff := x - l.Mu
	if diff > 0 {
		deriv[0] = 1 / l.Scale
	} else if diff < 0 {
		deriv[0] = -1 / l.Scale
	} else if diff == 0 {
		deriv[0] = 0
	} else {
		// must be NaN
		deriv[0] = math.NaN()
	}

	deriv[1] = math.Abs(diff)/(l.Scale*l.Scale) - 0.5/(l.Scale)
	return deriv
}

// Entropy returns the entropy of the distribution.
func (l Laplace) Entropy() float64 {
	return 1 + math.Log(2*l.Scale)
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (l Laplace) ExKurtosis() float64 {
	return 3
}

// Fit sets the parameters of the probability distribution from the
// data samples x with relative weights w.
// If weights is nil, then all the weights are 1.
// If weights is not nil, then the len(weights) must equal len(samples).
//
// Note: Laplace distribution has no FitPrior because it has no sufficient
// statistics.
func (l *Laplace) Fit(samples, weights []float64) {
	if len(samples) != len(weights) {
		panic("dist: length of samples and weights must match")
	}

	if len(samples) == 0 {
		panic("dist: must have at least one sample")
	}
	if len(samples) == 1 {
		l.Mu = samples[0]
		l.Scale = 0
		return
	}

	var (
		sortedSamples []float64
		sortedWeights []float64
	)
	if sort.Float64sAreSorted(samples) {
		sortedSamples = samples
		sortedWeights = weights
	} else {
		// Need to copy variables so the input variables aren't effected by the sorting
		sortedSamples = make([]float64, len(samples))
		copy(sortedSamples, samples)
		sortedWeights := make([]float64, len(samples))
		copy(sortedWeights, weights)

		stat.SortWeighted(sortedSamples, sortedWeights)
	}

	// The (weighted) median of the samples is the maximum likelihood estimate
	// of the mean parameter
	// TODO: Rethink quantile type when stat has more options
	l.Mu = stat.Quantile(0.5, stat.Empirical, sortedSamples, sortedWeights)

	sumWeights := floats.Sum(weights)

	// The scale parameter is the average absolute distance
	// between the sample and the mean
	absError := stat.Moment(1, samples, l.Mu, weights)

	l.Scale = absError / sumWeights
}

// LogProb computes the natural logarithm of the value of the probability density
// function at x.
func (l Laplace) LogProb(x float64) float64 {
	return -math.Ln2 - math.Log(l.Scale) - math.Abs(x-l.Mu)/l.Scale
}

// Mean returns the mean of the probability distribution.
func (l Laplace) Mean() float64 {
	return l.Mu
}

// Median returns the median of the LaPlace distribution.
func (l Laplace) Median() float64 {
	return l.Mu
}

// Mode returns the mode of the LaPlace distribution.
func (l Laplace) Mode() float64 {
	return l.Mu
}

// NumParameters returns the number of parameters in the distribution.
func (l Laplace) NumParameters() int {
	return 2
}

// Quantile returns the inverse of the cumulative probability distribution.
func (l Laplace) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic("dist: percentile out of bounds")
	}
	if p < 0.5 {
		return l.Mu + l.Scale*math.Log(1+2*(p-0.5))
	}
	return l.Mu - l.Scale*math.Log(1-2*(p-0.5))
}

// Parameters gets the parameters of the distribution. Panics if the length of
// the input slice is non-zero and not equal to the number of parameters
// The first element of Parameters is the Mean, the second is the scale variable.
func (l Laplace) Parameters(s []float64) []float64 {
	if s == nil {
		s = make([]float64, l.NumParameters())
	}
	if len(s) != l.NumParameters() {
		panic("dist: slice length mismatch")
	}
	s[0] = l.Mu
	s[1] = l.Scale
	return s
}

// Prob computes the value of the probability density function at x.
func (l Laplace) Prob(x float64) float64 {
	return math.Exp(l.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
func (l Laplace) Rand() float64 {
	var rnd float64
	if l.Source == nil {
		rnd = rand.Float64()
	} else {
		rnd = l.Source.Float64()
	}
	u := rnd - 0.5
	if u < 0 {
		return l.Mu + l.Scale*math.Log(1+2*u)
	}
	return l.Mu - l.Scale*math.Log(1-2*u)
}

// SetParameters gets the parameters of the distribution. Panics if the length of
// the input slice is not equal to the number of parameters.
// This sets Mu to be the first element of the slice and Scale to be the second
// element of the slice
func (l Laplace) SetParameters(s []float64) {
	if len(s) != l.NumParameters() {
		panic("dist: slice length mismatch")
	}
	l.Mu = s[0]
	l.Scale = s[1]
}

// Skewness returns the skewness of the distribution.
func (Laplace) Skewness() float64 {
	return 0
}

// StdDev returns the standard deviation of the distribution.
func (l Laplace) StdDev() float64 {
	return math.Sqrt2 * l.Scale
}

// Survival returns the survival function (complementary CDF) at x.
func (l Laplace) Survival(x float64) float64 {
	if x < l.Mu {
		return 1 - 0.5*math.Exp((x-l.Mu)/l.Scale)
	}
	return 0.5 * math.Exp(-(x-l.Mu)/l.Scale)
}

// Variance returns the variance of the probability distribution.
func (l Laplace) Variance() float64 {
	return 2 * l.Scale * l.Scale
}
