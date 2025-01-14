// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand/v2"

	"gonum.org/v1/gonum/spatial/r1"
)

// Uniform represents a multivariate uniform distribution.
type Uniform struct {
	bounds []r1.Interval
	dim    int
	rnd    *rand.Rand
}

// NewUniform creates a new uniform distribution with the given bounds.
func NewUniform(bnds []r1.Interval, src rand.Source) *Uniform {
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
		bounds: make([]r1.Interval, dim),
		dim:    dim,
	}
	if src != nil {
		u.rnd = rand.New(src)
	}
	for i, b := range bnds {
		u.bounds[i].Min = b.Min
		u.bounds[i].Max = b.Max
	}
	return u
}

// NewUnitUniform creates a new Uniform distribution over the dim-dimensional
// unit hypercube. That is, a uniform distribution where each dimension has
// Min = 0 and Max = 1.
func NewUnitUniform(dim int, src rand.Source) *Uniform {
	if dim <= 0 {
		panic(nonPosDimension)
	}
	bounds := make([]r1.Interval, dim)
	for i := range bounds {
		bounds[i].Min = 0
		bounds[i].Max = 1
	}
	u := Uniform{
		bounds: bounds,
		dim:    dim,
	}
	if src != nil {
		u.rnd = rand.New(src)
	}
	return &u
}

// Bounds returns the bounds on the variables of the distribution.
//
// If dst is not nil, the bounds will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
func (u *Uniform) Bounds(bounds []r1.Interval) []r1.Interval {
	if bounds == nil {
		bounds = make([]r1.Interval, u.Dim())
	}
	if len(bounds) != u.Dim() {
		panic(badInputLength)
	}
	copy(bounds, u.bounds)
	return bounds
}

// CDF returns the value of the multidimensional cumulative distribution
// function of the probability distribution at the point x.
//
// If dst is not nil, the value will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution. CDF will also panic
// if the length of x is not equal to the dimension of the distribution.
func (u *Uniform) CDF(dst, x []float64) []float64 {
	if len(x) != u.dim {
		panic(badSizeMismatch)
	}
	dst = reuseAs(dst, u.dim)

	for i, v := range x {
		if v < u.bounds[i].Min {
			dst[i] = 0
		} else if v > u.bounds[i].Max {
			dst[i] = 1
		} else {
			dst[i] = (v - u.bounds[i].Min) / (u.bounds[i].Max - u.bounds[i].Min)
		}
	}
	return dst
}

// Dim returns the dimension of the distribution.
func (u *Uniform) Dim() int {
	return u.dim
}

// Entropy returns the differential entropy of the distribution.
func (u *Uniform) Entropy() float64 {
	// Entropy is log of the volume.
	var logVol float64
	for _, b := range u.bounds {
		logVol += math.Log(b.Max - b.Min)
	}
	return logVol
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

// Mean returns the mean of the probability distribution.
//
// If dst is not nil, the mean will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
func (u *Uniform) Mean(dst []float64) []float64 {
	dst = reuseAs(dst, u.dim)
	for i, b := range u.bounds {
		dst[i] = (b.Max + b.Min) / 2
	}
	return dst
}

// Prob computes the value of the probability density function at x.
func (u *Uniform) Prob(x []float64) float64 {
	return math.Exp(u.LogProb(x))
}

// Rand generates a random sample according to the distribution.
//
// If dst is not nil, the sample will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
func (u *Uniform) Rand(dst []float64) []float64 {
	dst = reuseAs(dst, u.dim)
	if u.rnd == nil {
		for i, b := range u.bounds {
			dst[i] = rand.Float64()*(b.Max-b.Min) + b.Min
		}
		return dst
	}
	for i, b := range u.bounds {
		dst[i] = u.rnd.Float64()*(b.Max-b.Min) + b.Min
	}
	return dst
}

// Quantile returns the value of the multi-dimensional inverse cumulative
// distribution function at p.
//
// If dst is not nil, the quantile will be stored in-place into dst and
// returned, otherwise a new slice will be allocated first. If dst is not nil,
// it must have length equal to the dimension of the distribution. Quantile will
// also panic if the length of p is not equal to the dimension of the
// distribution.
//
// All of the values of p must be between 0 and 1, inclusive, or Quantile will
// panic.
func (u *Uniform) Quantile(dst, p []float64) []float64 {
	if len(p) != u.dim {
		panic(badSizeMismatch)
	}
	dst = reuseAs(dst, u.dim)
	for i, v := range p {
		if v < 0 || v > 1 {
			panic(badQuantile)
		}
		dst[i] = v*(u.bounds[i].Max-u.bounds[i].Min) + u.bounds[i].Min
	}
	return dst
}
