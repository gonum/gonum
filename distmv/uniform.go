// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand"
)

type Bound struct {
	Min float64
	Max float64
}

// Uniform represents a multivariate uniform distribution.
type Uniform struct {
	bounds []Bound
	dim    int
	src    *rand.Source
}

// NewUniform creates a new uniform distribution with the given bounds.
func NewUniform(bnds []Bound, src *rand.Source) *Uniform {
	dim := len(bnds)
	if dim == 0 {
		panic(badZeroDimension)
	}
	for _, b := range bnds {
		if b.Max < b.Min {
			panic("uniform: maximum less than minimum")
		}
	}
	u := &Uniform{
		bounds: make([]Bound, dim),
		dim:    dim,
		src:    src,
	}
	for i, b := range bnds {
		u.bounds[i].Min = b.Min
		u.bounds[i].Max = b.Max
	}
	return u
}

// CDF returns the multidimensional cumulative distribution function of the
// probability distribution at the point x. If p is non-nil, the CDF is stored
// in-place into the first argument, otherwise a new slice is allocated and
// returned.
//
// CDF will panic if len(x) is not equal to the dimension of the distribution,
// or if p is non-nil and len(p) is not equal to the dimension of the distribution.
func (u *Uniform) CDF(p, x []float64) []float64 {
	if len(x) != u.dim {
		panic(badSizeMismatch)
	}
	if p == nil {
		p = make([]float64, u.dim)
	}
	if len(p) != u.dim {
		panic(badSizeMismatch)
	}
	for i, v := range x {
		if v < u.bounds[i].Min {
			p[i] = 0
		} else if v > u.bounds[i].Max {
			p[i] = 1
		} else {
			p[i] = (v - u.bounds[i].Min) / (u.bounds[i].Max - u.bounds[i].Min)
		}
	}
	return p
}

// Dim returns the dimension of the distribution.
func (u *Uniform) Dim() int {
	return u.dim
}

// LogProb computes the log of the pdf of the point x.
func (u *Uniform) LogProb(x []float64) float64 {
	dim := u.dim
	if len(x) != dim {
		panic(badSizeMismatch)
	}
	var logprob float64
	for i, b := range u.bounds {
		if x[i] < b.Min || x[i] > b.Max {
			return math.Inf(-1)
		}
		logprob -= math.Log(b.Max - b.Min)
	}
	return logprob
}

// Mean returns the mean of the probability distribution at x. If the
// input argument is nil, a new slice will be allocated, otherwise the result
// will be put in-place into the receiver.
func (u *Uniform) Mean(x []float64) []float64 {
	x = reuseAs(x, u.dim)
	for i, b := range u.bounds {
		x[i] = (b.Max + b.Min) / 2
	}
	return x
}

// Prob computes the value of the probability density function at x.
func (u *Uniform) Prob(x []float64) float64 {
	return math.Exp(u.LogProb(x))
}

// Rand generates a random number according to the distributon.
// If the input slice is nil, new memory is allocated, otherwise the result is stored
// in place.
func (u *Uniform) Rand(x []float64) []float64 {
	x = reuseAs(x, u.dim)
	if u.src == nil {
		for i, b := range u.bounds {
			x[i] = rand.Float64()*(b.Max-b.Min) + b.Min
		}
		return x
	}
	for i, b := range u.bounds {
		x[i] = rand.Float64()*(b.Max-b.Min) + b.Min
	}
	return x
}

// Quantile returns the multi-dimensional inverse cumulative distribution function.
// len(x) must equal len(p), and if x is non-nil, len(x) must also equal len(p).
// If x is nil, a new slice will be allocated and returned, otherwise the quantile
// will be stored in-place into x. All of the values of p must be between 0 and 1,
// or Quantile will panic.
func (u *Uniform) Quantile(x, p []float64) []float64 {
	if len(p) != u.dim {
		panic(badSizeMismatch)
	}
	if x == nil {
		x = make([]float64, u.dim)
	}
	if len(x) != u.dim {
		panic(badSizeMismatch)
	}
	for i, v := range p {
		if v < 0 || v > 1 {
			panic(badQuantile)
		}
		x[i] = v*(u.bounds[i].Max-u.bounds[i].Min) + u.bounds[i].Min
	}
	return x
}
