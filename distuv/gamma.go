// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"
)

// Gamma implements the Gamma distribution, a two-parameter continuous distribution
// with support over the positive real numbers.
//
// The gamma distribution has density function
//  β^α / Γ(α) x^(α-1)e^(-βx)
//
// For more information, see https://en.wikipedia.org/wiki/Gamma_distribution
type Gamma struct {
	// Alpha is the shape parameter of the distribution. Alpha must be greater
	// than 0. If Alpha == 1, this is equivalent to an exponential distribution.
	Alpha float64
	// Beta is the rate parameter of the distribution. Beta must be greater than 0.
	// If Beta == 2, this is equivalent to a Chi-Squared distribution.
	Beta float64

	Source *rand.Rand
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (g Gamma) ExKurtosis() float64 {
	return 6 / g.Alpha
}

// LogProb computes the natural logarithm of the value of the probability
// density function at x.
func (g Gamma) LogProb(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	a := g.Alpha
	b := g.Beta
	lg, _ := math.Lgamma(a)
	return a*math.Log(b) - lg + (a-1)*math.Log(x) - b*x
}

// Mean returns the mean of the probability distribution.
func (g Gamma) Mean() float64 {
	return g.Alpha / g.Beta
}

// Mode returns the mode of the normal distribution.
//
// The mode is NaN in the special case where the Alpha (shape) parameter
// is less than 1.
func (g Gamma) Mode() float64 {
	if g.Alpha < 1 {
		return math.NaN()
	}
	return (g.Alpha - 1) / g.Beta
}

// NumParameters returns the number of parameters in the distribution.
func (Gamma) NumParameters() int {
	return 2
}

// Prob computes the value of the probability density function at x.
func (g Gamma) Prob(x float64) float64 {
	return math.Exp(g.LogProb(x))
}

// Rand returns a random sample drawn from the distribution.
//
// Rand panics if either alpha or beta is <= 0.
func (g Gamma) Rand() float64 {
	if g.Beta <= 0 {
		panic("gamma: beta <= 0")
	}

	unifrnd := rand.Float64
	if g.Source != nil {
		unifrnd = g.Source.Float64
	}

	exprnd := rand.ExpFloat64
	if g.Source != nil {
		exprnd = g.Source.ExpFloat64
	}
	a := g.Alpha
	b := g.Beta
	switch {
	case a <= 0:
		panic("gamma: alpha < 0")
	case a == 1:
		// Generate from exponential
		return exprnd() / b
	case a < 1:
		// Generate using
		//  Xi, Bowei, Kean Ming Tan, and Chuanhai Liu. "Logarithmic
		//  Transformation-Based Gamma Random Number Generators." Journal of
		//  Statistical Software 55.1 (2013): 1-17.
		// Algorithm 2.
		umax := math.Pow(a/math.E, a/2)
		vmin := -2 / math.E
		vmax := 2 * a / math.E / (math.E - a)
		var t, t1 float64
		for {
			u := umax * unifrnd()
			t = (unifrnd()*(vmax-vmin) + vmin) / u
			t1 = math.Exp(t / a)
			if 2*math.Log(u) <= t-t1 {
				break
			}
		}
		if a >= 0.01 {
			return t1 / b
		}
		return t / a / b
	case a > 1:
		// Generate using
		//  Martino, Luca, and David Luengo. "Extremely efficient generation of
		//  Gamma random variables for α >= 1." arXiv preprint arXiv:1304.3800 (2013).
		ap := math.Floor(a)
		var bp, lkp float64
		// The paper says ap < 2, but ap must be an integer at least 1.
		if ap == 1 {
			bp = b / a
			lkp = (1 - a) + (a-1)*(a/b)
		} else {
			bp = b * (ap - 1) / (a - 1)
			lkp = (ap - a) + (a-ap)*(a-1)/b
		}
		for {
			// Draw a sample
			x := exprnd()
			for i := 1; i < int(ap); i++ {
				x += exprnd()
			}
			x /= bp

			// Compute accept/reject
			lx := math.Log(x)
			lpx := (a-1)*lx + -b*x
			lpix := lkp + (ap-1)*lx + -bp*x
			if unifrnd() < math.Exp(lpx-lpix) {
				return x
			}
		}
	}
	panic("unreachable")
}

// StdDev returns the standard deviation of the probability distribution.
func (g Gamma) StdDev() float64 {
	return math.Sqrt(g.Variance())
}

// Variance returns the variance of the probability distribution.
func (g Gamma) Variance() float64 {
	return g.Alpha / g.Beta / g.Beta
}
