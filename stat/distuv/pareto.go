// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"

	"golang.org/x/exp/rand"
)

// Pareto implements the Pareto (Type I) distribution, a one parameter distribution
// with support above the scale parameter.
//
// The density function is given by
//  (α x_m^{α})/(x^{α+1}) for x >= x_m.
//
// For more information, see https://en.wikipedia.org/wiki/Pareto_distribution.
type Pareto struct {
	// Xm is the scale parameter.
	// Xm must be greater than 0.
	Xm float64

	// Alpha is the shape parameter.
	// Alpha must be greater than 0.
	Alpha float64

	Src *rand.Rand
}

// CDF computes the value of the cumulative density function at x.
func (p Pareto) CDF(x float64) float64 {
	if x < p.Xm {
		return 0
	}
	return 1 - p.Survival(x)
}

// Entropy returns the differential entropy of the distribution.
func (p Pareto) Entropy() float64 {
	return math.Log(p.Xm) - math.Log(p.Alpha) + (1 + 1/p.Alpha)
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (p Pareto) ExKurtosis() float64 {
	if p.Alpha <= 4 {
		return 0
	}
	return 6 * (p.Alpha*p.Alpha*p.Alpha + p.Alpha*p.Alpha - 6*p.Alpha - 2) / (p.Alpha * (p.Alpha - 3) * (p.Alpha - 4))

}

// LogProb computes the natural logarithm of the value of the probability
// density function at x.
func (p Pareto) LogProb(x float64) float64 {
	if x < p.Xm {
		return math.Inf(-1)
	}
	return math.Log(p.Alpha) + p.Alpha*math.Log(p.Xm) - (p.Alpha+1)*math.Log(x)
}

// Mean returns the mean of the probability distribution.
func (p Pareto) Mean() float64 {
	if p.Alpha <= 1 {
		return math.Inf(1)
	}
	return p.Alpha * p.Xm / (p.Alpha - 1)
}

// Median returns the median of the pareto distribution.
func (p Pareto) Median() float64 {
	return p.Xm * math.Pow(2, 1/p.Alpha)
}

// Mode returns the mode of the distribution.
func (p Pareto) Mode() float64 {
	return p.Xm
}

// NumParameters returns the number of parameters in the distribution.
func (p Pareto) NumParameters() int {
	return 2
}

// Prob computes the value of the probability density function at x.
func (p Pareto) Prob(x float64) float64 {
	return math.Exp(p.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
func (p Pareto) Rand() float64 {
	var rnd float64
	if p.Src == nil {
		rnd = rand.ExpFloat64()
	} else {
		rnd = p.Src.ExpFloat64()
	}
	return math.Exp(math.Log(p.Xm) + 1/p.Alpha*rnd)
}

// StdDev returns the standard deviation of the probability distribution.
func (p Pareto) StdDev() float64 {
	return math.Sqrt(p.Variance())
}

// Survival returns the survival function (complementary CDF) at x.
func (p Pareto) Survival(x float64) float64 {
	if x < p.Xm {
		return 1
	}
	return math.Exp(p.Alpha * (math.Log(p.Xm) - math.Log(x)))
}

// Variance returns the variance of the probability distribution.
func (p Pareto) Variance() float64 {
	if p.Alpha <= 2 {
		return math.Inf(1)
	}
	am1 := p.Alpha - 1
	return p.Xm * p.Xm * p.Alpha / (am1 * am1 * (p.Alpha - 2))
}
