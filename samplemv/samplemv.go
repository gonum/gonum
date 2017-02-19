// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package samplemv implements advanced sampling routines from explicit and implicit
// probability distributions.
//
// Each sampling routine is implemented as a stateless function with a
// complementary wrapper type. The wrapper types allow the sampling routines
// to implement interfaces.
package samplemv

import (
	"errors"
	"math"
	"math/rand"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat/distmv"
)

var (
	badLengthMismatch = "samplemv: slice length mismatch"
)

var (
	_ Sampler = LatinHypercuber{}
	_ Sampler = (*Rejectioner)(nil)
	_ Sampler = IIDer{}

	_ WeightedSampler = SampleUniformWeighted{}
	_ WeightedSampler = Importancer{}
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Sampler generates a batch of samples according to the rule specified by the
// implementing type. The number of samples generated is equal to rows(batch),
// and the samples are stored in-place into the input.
type Sampler interface {
	Sample(batch *mat64.Dense)
}

// WeightedSampler generates a batch of samples and their relative weights
// according to the rule specified by the implementing type. The number of samples
// generated is equal to rows(batch), and the samples and weights
// are stored in-place into the inputs. The length of weights must equal
// rows(batch), otherwise SampleWeighted will panic.
type WeightedSampler interface {
	SampleWeighted(batch *mat64.Dense, weights []float64)
}

// SampleUniformWeighted wraps a Sampler type to create a WeightedSampler where all
// weights are equal.
type SampleUniformWeighted struct {
	Sampler
}

// SampleWeighted generates rows(batch) samples from the embedded Sampler type
// and sets all of the weights equal to 1. If rows(batch) and len(weights)
// of weights are not equal, SampleWeighted will panic.
func (w SampleUniformWeighted) SampleWeighted(batch *mat64.Dense, weights []float64) {
	r, _ := batch.Dims()
	if r != len(weights) {
		panic(badLengthMismatch)
	}
	w.Sample(batch)
	for i := range weights {
		weights[i] = 1
	}
}

// LatinHypercuber is a wrapper around the LatinHypercube sampling generation
// method.
type LatinHypercuber struct {
	Q   distmv.Quantiler
	Src *rand.Rand
}

// Sample generates rows(batch) samples using the LatinHypercube generation
// procedure.
func (l LatinHypercuber) Sample(batch *mat64.Dense) {
	LatinHypercube(batch, l.Q, l.Src)
}

// LatinHypercube generates rows(batch) samples using Latin hypercube sampling
// from the given distribution. If src is not nil, it will be used to generate
// random numbers, otherwise rand.Float64 will be used.
//
// Latin hypercube sampling divides the cumulative distribution function into equally
// spaced bins and guarantees that one sample is generated per bin. Within each bin,
// the location is randomly sampled. The distmv.NewUnitUniform function can be used
// for easy sampling from the unit hypercube.
func LatinHypercube(batch *mat64.Dense, q distmv.Quantiler, src *rand.Rand) {
	r, c := batch.Dims()
	var f64 func() float64
	var perm func(int) []int
	if src != nil {
		f64 = src.Float64
		perm = src.Perm
	} else {
		f64 = rand.Float64
		perm = rand.Perm
	}
	r64 := float64(r)
	for i := 0; i < c; i++ {
		p := perm(r)
		for j := 0; j < r; j++ {
			var v float64
			v = f64()
			v = v/r64 + float64(j)/r64
			batch.Set(p[j], i, v)
		}
	}
	p := make([]float64, c)
	for i := 0; i < r; i++ {
		copy(p, batch.RawRowView(i))
		q.Quantile(batch.RawRowView(i), p)
	}
}

// Importancer is a wrapper around the Importance sampling generation method.
type Importancer struct {
	Target   distmv.LogProber
	Proposal distmv.RandLogProber
}

// SampleWeighted generates rows(batch) samples using the Importance sampling
// generation procedure.
func (l Importancer) SampleWeighted(batch *mat64.Dense, weights []float64) {
	Importance(batch, weights, l.Target, l.Proposal)
}

// Importance sampling generates rows(batch) samples from the proposal distribution,
// and stores the locations and importance sampling weights in place.
//
// Importance sampling is a variance reduction technique where samples are
// generated from a proposal distribution, q(x), instead of the target distribution
// p(x). This allows relatively unlikely samples in p(x) to be generated more frequently.
//
// The importance sampling weight at x is given by p(x)/q(x). To reduce variance,
// a good proposal distribution will bound this sampling weight. This implies the
// support of q(x) should be at least as broad as p(x), and q(x) should be "fatter tailed"
// than p(x).
//
// If weights is nil, the weights are not stored. The length of weights must equal
// the length of batch, otherwise Importance will panic.
func Importance(batch *mat64.Dense, weights []float64, target distmv.LogProber, proposal distmv.RandLogProber) {
	r, _ := batch.Dims()
	if r != len(weights) {
		panic(badLengthMismatch)
	}
	for i := 0; i < r; i++ {
		v := batch.RawRowView(i)
		proposal.Rand(v)
		weights[i] = math.Exp(target.LogProb(v) - proposal.LogProb(v))
	}
}

// ErrRejection is returned when the constant in Rejection is not sufficiently high.
var ErrRejection = errors.New("rejection: acceptance ratio above 1")

// Rejectioner is a wrapper around the Rejection sampling generation procedure.
// If the rejection sampling fails during the call to Sample, all samples will
// be set to math.NaN() and a call to Err will return a non-nil value.
type Rejectioner struct {
	C        float64
	Target   distmv.LogProber
	Proposal distmv.RandLogProber
	Src      *rand.Rand

	err      error
	proposed int
}

// Err returns nil if the most recent call to sample was successful, and returns
// ErrRejection if it was not.
func (r *Rejectioner) Err() error {
	return r.err
}

// Proposed returns the number of samples proposed during the most recent call to
// Sample.
func (r *Rejectioner) Proposed() int {
	return r.proposed
}

// Sample generates rows(batch) using the Rejection sampling generation procedure.
// Rejection sampling may fail if the constant is insufficiently high, as described
// in the function comment for Rejection. If the generation fails, the samples
// are set to math.NaN(), and a call to Err will return a non-nil value.
func (r *Rejectioner) Sample(batch *mat64.Dense) {
	r.err = nil
	r.proposed = 0
	proposed, ok := Rejection(batch, r.Target, r.Proposal, r.C, r.Src)
	if !ok {
		r.err = ErrRejection
	}
	r.proposed = proposed
}

// Rejection generates rows(batch) samples using the rejection sampling algorithm and
// stores them in place into samples.
// Sampling continues until batch is filled. Rejection returns the total number of proposed
// locations and a boolean indicating if the rejection sampling assumption is
// violated (see details below). If the returned boolean is false, all elements
// of samples are set to NaN. If src != nil, it will be used to generate random
// numbers, otherwise rand.Float64 will be used.
//
// Rejection sampling generates points from the target distribution by using
// the proposal distribution. At each step of the algorithm, the proposaed point
// is accepted with probability
//  p = target(x) / (proposal(x) * c)
// where target(x) is the probability of the point according to the target distribution
// and proposal(x) is the probability according to the proposal distribution.
// The constant c must be chosen such that target(x) < proposal(x) * c for all x.
// The expected number of proposed samples is len(samples) * c.
//
// Target may return the true (log of) the probablity of the location, or it may return
// a value that is proportional to the probability (logprob + constant). This is
// useful for cases where the probability distribution is only known up to a normalization
// constant.
func Rejection(batch *mat64.Dense, target distmv.LogProber, proposal distmv.RandLogProber, c float64, src *rand.Rand) (nProposed int, ok bool) {
	if c < 1 {
		panic("rejection: acceptance constant must be greater than 1")
	}
	f64 := rand.Float64
	if src != nil {
		f64 = src.Float64
	}
	r, dim := batch.Dims()
	v := make([]float64, dim)
	var idx int
	for {
		nProposed++
		proposal.Rand(v)
		qx := proposal.LogProb(v)
		px := target.LogProb(v)
		accept := math.Exp(px-qx) / c
		if accept > 1 {
			// Invalidate the whole result and return a failure.
			for i := 0; i < r; i++ {
				for j := 0; j < dim; j++ {
					batch.Set(i, j, math.NaN())
				}
			}
			return nProposed, false
		}
		if accept > f64() {
			batch.SetRow(idx, v)
			idx++
			if idx == r {
				break
			}
		}
	}
	return nProposed, true
}

// IIDer is a wrapper around the IID sample generation method.
type IIDer struct {
	Dist distmv.Rander
}

// Sample generates a set of identically and independently distributed samples.
func (iid IIDer) Sample(batch *mat64.Dense) {
	IID(batch, iid.Dist)
}

// IID generates a set of independently and identically distributed samples from
// the input distribution.
func IID(batch *mat64.Dense, d distmv.Rander) {
	r, _ := batch.Dims()
	for i := 0; i < r; i++ {
		d.Rand(batch.RawRowView(i))
	}
}
