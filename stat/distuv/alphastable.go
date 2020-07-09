// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"

	"golang.org/x/exp/rand"
)

// AlphaStable represents an α-stable distribution with four parameters:
//  stability 0 < α ≤ 2
//  skewness -1 ≤ β ≤ 1
//  scale c > 0
//  location μ
// More information at https://en.wikipedia.org/wiki/Stable_distribution
type AlphaStable struct {
	Alpha float64
	Beta  float64
	Scale float64
	Mu    float64
	Src   rand.Source
}

// ExKurtosis returns the excess kurtosis of the distribution.
// ExKurtosis returns NaN for Alpha < 2.
func (a AlphaStable) ExKurtosis() float64 {
	if a.Alpha == 2 {
		return 0
	}
	return math.NaN()
}

// Mean returns the mean of the probability distribution.
// Mean returns NaN if Alpha <= 1.
func (a AlphaStable) Mean() float64 {
	if a.Alpha > 1 {
		return a.Mu
	}
	return math.NaN()
}

// Median returns the median of the distribution.
// Median panics when Beta != 0, because then the mode is not analytically
// expressible.
func (a AlphaStable) Median() float64 {
	if a.Beta == 0 {
		return a.Mu
	}
	panic("Cannot compute Median for Beta != 0")
}

// Mode returns the mode of the distribution.
// Mode panics when Beta != 0, because then the mode is not analytically
// expressible.
func (a AlphaStable) Mode() float64 {
	if a.Beta == 0 {
		return a.Mu
	}
	panic("Cannot compute Mode for Beta != 0")
}

// NumParameters returns the number of parameters in the distribution.
func (a AlphaStable) NumParameters() int {
	return 4
}

// Rand returns a random sample drawn from the distribution.
func (a AlphaStable) Rand() float64 {
	// From https://en.wikipedia.org/wiki/Stable_distribution#Simulation_of_stable_variables
	const halfPi = math.Pi / 2
	u := Uniform{-halfPi, halfPi, a.Src}.Rand()
	w := Exponential{1, a.Src}.Rand()
	var y float64
	if a.Alpha == 1 {
		f := halfPi + a.Beta*u
		x := (f*math.Tan(u) - a.Beta*math.Log(halfPi*w*math.Cos(u)/f)) / halfPi
		y = a.Scale*(x+a.Beta*math.Log(a.Scale)/halfPi) + a.Mu
	} else {
		zeta := -a.Beta * math.Tan(halfPi*a.Alpha)
		xi := math.Atan(-zeta) / a.Alpha
		f := a.Alpha * (u + xi)
		g := math.Sqrt(1+zeta*zeta) * math.Pow(math.Cos(u-f)/w, 1-a.Alpha) / math.Cos(u)
		x := math.Pow(g, 1/a.Alpha) * math.Sin(f)
		y = a.Scale*x + a.Mu
	}
	return y
}

// Skewness returns the skewness of the distribution.
// Skewness returns NaN if Alpha < 2.
func (a AlphaStable) Skewness() float64 {
	if a.Alpha == 2 {
		return 0
	}
	return math.NaN()
}

// StdDev returns the standard deviation of the probability distribution.
func (a AlphaStable) StdDev() float64 {
	return math.Sqrt(a.Variance())
}

// Variance returns the variance of the probability distribution.
// Variance returns +Inf if Alpha < 2.
func (a AlphaStable) Variance() float64 {
	if a.Alpha == 2 {
		return 2 * a.Scale * a.Scale
	}
	return math.Inf(1)
}
