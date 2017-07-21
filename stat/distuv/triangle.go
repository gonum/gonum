// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"math/rand"
)

// UnitTriangle is an instantiation of the triangle distribution with A=0, C=0.5, and B=1
var UnitTriangle = Triangle{A: 0, C: 0.5, B: 1}

// Triangle represents a triangle distribution (https://en.wikipedia.org/wiki/Triangular_distribution).
type Triangle struct {
	A      float64
	C      float64
	B      float64
	Source *rand.Rand
}

// CDF computes the value of the cumulative density function at x.
func (t Triangle) CDF(x float64) float64 {
	if x <= t.A {
		return 0
	}
	if x <= t.C {
		return math.Pow(x-t.A, 2) / ((t.B - t.A) * (t.C - t.A))
	}
	if x < t.B {
		return 1 - math.Pow(t.B-x, 2)/((t.B-t.A)*(t.B-t.C))
	}
	return 1
}

// Entropy returns the entropy of the distribution.
func (t Triangle) Entropy() float64 {
	return 1.0/2.0 + math.Log((t.B-t.A)/2)
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

// MarshalParameters implements the ParameterMarshaler interface
func (t Triangle) MarshalParameters(p []Parameter) {
	if len(p) != t.NumParameters() {
		panic("triangle: improper parameter length")
	}
	p[0].Name = "A"
	p[0].Value = t.A
	p[1].Name = "C"
	p[1].Value = t.C
	p[2].Name = "B"
	p[2].Value = t.B
	return
}

// Mean returns the mean of the probability distribution.
func (t Triangle) Mean() float64 {
	return (t.A + t.B + t.C) / 3
}

// Median returns the median of the probability distribution.
func (t Triangle) Median() float64 {
	if t.C >= (t.A+t.B)/2 {
		return t.A + math.Sqrt((t.B-t.A)*(t.C-t.A)/2)
	} else {
		return t.B - math.Sqrt((t.B-t.A)*(t.B-t.C)/2)
	}
}

// Mode returns the mode of the probability distribution.
func (t Triangle) Mode() float64 {
	return t.C
}

// NumParameters returns the number of parameters in the distribution.
func (Triangle) NumParameters() int {
	return 3
}

// Prob computes the value of the probability density function at x.
func (t Triangle) Prob(x float64) float64 {
	if x < t.A {
		return 0
	}
	if x < t.C {
		return 2 * (x - t.A) / ((t.B - t.A) * (t.C - t.A))
	}
	if x == t.C {
		return 2 / (t.B - t.A)
	}
	if x <= t.B {
		return 2 * (t.B - x) / ((t.B - t.A) * (t.B - t.C))
	}
	return 0
}

// Quantile returns the inverse of the cumulative probability distribution.
func (t Triangle) Quantile(p float64) float64 {
	if p < 0 || p > 1 {
		panic(badPercentile)
	}

	f := (t.C - t.A) / (t.B - t.A)

	if p < f {
		return t.A + math.Sqrt(p*(t.B-t.A)*(t.C-t.A))
	} else {
		return t.B - math.Sqrt((1-p)*(t.B-t.A)*(t.B-t.C))
	}
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
	n := math.Sqrt2 * (t.A + t.B - 2*t.C) * (2*t.A - t.B - t.C) * (t.A - 2*t.B + t.C)
	d := 5 * math.Pow(math.Pow(t.A, 2)+math.Pow(t.B, 2)+math.Pow(t.C, 2)-t.A*t.B-t.A*t.C-t.B*t.C, 3/2)

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

// UnmarshalParameters implements the ParameterMarshaler interface
func (t *Triangle) UnmarshalParameters(p []Parameter) {
	if len(p) != t.NumParameters() {
		panic("triangle: incorrect number of parameters to set")
	}
	if p[0].Name != "A" {
		panic("triangle: " + panicNameMismatch)
	}
	if p[1].Name != "C" {
		panic("triangle: " + panicNameMismatch)
	}
	if p[2].Name != "B" {
		panic("triangle: " + panicNameMismatch)
	}

	t.A = p[0].Value
	t.C = p[1].Value
	t.B = p[2].Value
}

// Variance returns the variance of the probability distribution.
func (t Triangle) Variance() float64 {
	return (math.Pow(t.A, 2) + math.Pow(t.B, 2) + math.Pow(t.C, 2) - t.A*t.B - t.A*t.C - t.B*t.C) / 18
}
