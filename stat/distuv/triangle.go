// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"
)

// Triangle represents a triangle distribution (https://en.wikipedia.org/wiki/Triangular_distribution).
type Triangle struct {
	a, b, c float64
	Source  *rand.Rand
}

// NewTriangle constructs a new triangle distribution with lower limit a, upper limit b, and mode c.
// Constraints are a < b and a ≤ c ≤ b.
// This distribution is uncommon in nature, but may be useful for simulation.
func NewTriangle(a, b, c float64) Triangle {
	checkTriangleParameters(a, b, c)
	return Triangle{a, b, c, nil}
}

func checkTriangleParameters(a, b, c float64) {
	if a >= b {
		panic("triangle: constraint of a < b violated")
	}
	if a > c {
		panic("triangle: constraint of a <= c violated")
	}
	if c > b {
		panic("triangle: constraint of c <= b violated")
	}
}

// CDF computes the value of the cumulative density function at x.
func (t Triangle) CDF(x float64) float64 {
	switch {
	case x <= t.a:
		return 0
	case x <= t.c:
		d := x - t.a
		return (d * d) / ((t.b - t.a) * (t.c - t.a))
	case x < t.b:
		d := t.b - x
		return 1 - (d*d)/((t.b-t.a)*(t.b-t.c))
	default:
		return 1
	}
}

// Entropy returns the entropy of the distribution.
func (t Triangle) Entropy() float64 {
	return 0.5 + math.Log(t.b-t.a) - math.Ln2
}

// ExKurtosis returns the excess kurtosis of the distribution.
func (Triangle) ExKurtosis() float64 {
	return -3.0 / 5.0
}

// Fit is not appropriate for Triangle, because the distribution is generally used when there is little data.

// LogProb computes the natural logarithm of the value of the probability density function at x.
func (t Triangle) LogProb(x float64) float64 {
	return math.Log(t.Prob(x))
}

// Mean returns the mean of the probability distribution.
func (t Triangle) Mean() float64 {
	return (t.a + t.b + t.c) / 3
}

// Median returns the median of the probability distribution.
func (t Triangle) Median() float64 {
	if t.c >= (t.a+t.b)/2 {
		return t.a + math.Sqrt((t.b-t.a)*(t.c-t.a)/2)
	}
	return t.b - math.Sqrt((t.b-t.a)*(t.b-t.c)/2)
}

// Mode returns the mode of the probability distribution.
func (t Triangle) Mode() float64 {
	return t.c
}

// NumParameters returns the number of parameters in the distribution.
func (Triangle) NumParameters() int {
	return 3
}

// Prob computes the value of the probability density function at x.
func (t Triangle) Prob(x float64) float64 {
	switch {
	case x < t.a:
		return 0
	case x < t.c:
		return 2 * (x - t.a) / ((t.b - t.a) * (t.c - t.a))
	case x == t.c:
		return 2 / (t.b - t.a)
	case x <= t.b:
		return 2 * (t.b - x) / ((t.b - t.a) * (t.b - t.c))
	default:
		return 0
	}
}

// Quantile returns the inverse of the cumulative probability distribution.
func (t Triangle) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic(badPercentile)
	}

	f := (t.c - t.a) / (t.b - t.a)

	if p < f {
		return t.a + math.Sqrt(p*(t.b-t.a)*(t.c-t.a))
	}
	return t.b - math.Sqrt((1-p)*(t.b-t.a)*(t.b-t.c))
}

// Rand returns a random sample drawn from the distribution.
func (t Triangle) Rand() float64 {
	var rnd float64
	if t.Source == nil {
		rnd = rand.Float64()
	} else {
		rnd = t.Source.Float64()
	}

	return t.Quantile(rnd)
}

// Skewness returns the skewness of the distribution.
func (t Triangle) Skewness() float64 {
	n := math.Sqrt2 * (t.a + t.b - 2*t.c) * (2*t.a - t.b - t.c) * (t.a - 2*t.b + t.c)
	d := 5 * math.Pow(t.a*t.a+t.b*t.b+t.c*t.c-t.a*t.b-t.a*t.c-t.b*t.c, 3.0/2.0)

	return n / d
}

// StdDev returns the standard deviation of the probability distribution.
func (t Triangle) StdDev() float64 {
	return math.Sqrt(t.Variance())
}

// Survival returns the survival function (complementary CDF) at x.
func (t Triangle) Survival(x float64) float64 {
	return 1 - t.CDF(x)
}

// MarshalParameters implements the ParameterMarshaler interface
func (t Triangle) MarshalParameters(p []Parameter) {
	if len(p) != t.NumParameters() {
		panic("triangle: improper parameter length")
	}
	p[0].Name = "A"
	p[0].Value = t.a
	p[1].Name = "B"
	p[1].Value = t.b
	p[2].Name = "C"
	p[2].Value = t.c
}

// UnmarshalParameters implements the ParameterMarshaler interface
func (t *Triangle) UnmarshalParameters(p []Parameter) {
	if len(p) != t.NumParameters() {
		panic("triangle: incorrect number of parameters to set")
	}
	if p[0].Name != "A" {
		panic("triangle: " + panicNameMismatch)
	}
	if p[1].Name != "B" {
		panic("triangle: " + panicNameMismatch)
	}
	if p[2].Name != "C" {
		panic("triangle: " + panicNameMismatch)
	}

	checkTriangleParameters(p[0].Value, p[1].Value, p[2].Value)

	t.a = p[0].Value
	t.b = p[1].Value
	t.c = p[2].Value
}

// Variance returns the variance of the probability distribution.
func (t Triangle) Variance() float64 {
	return (t.a*t.a + t.b*t.b + t.c*t.c - t.a*t.b - t.a*t.c - t.b*t.c) / 18
}
