// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sampleuv implements advanced sampling routines from explicit and implicit
// probability distributions.
//
// Each sampling routine is implemented as a stateless function with a
// complementary wrapper type. The wrapper types allow the sampling routines
// to implement interfaces.
package sampleuv

import (
	"errors"
	"math"
	"math/rand"

	"github.com/gonum/stat/distuv"
)

var (
	badLengthMismatch = "sample: slice length mismatch"
)

var (
	_ Sampler = LatinHypercuber{}
	_ Sampler = MetropolisHastingser{}
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
// implementing type. The number of samples generated is equal to len(batch),
// and the samples are stored in-place into the input.
type Sampler interface {
	Sample(batch []float64)
}

// WeightedSampler generates a batch of samples and their relative weights
// according to the rule specified by the implementing type. The number of samples
// generated is equal to len(batch), and the samples and weights
// are stored in-place into the inputs. The length of weights must equal
// len(batch), otherwise SampleWeighted will panic.
type WeightedSampler interface {
	SampleWeighted(batch, weights []float64)
}

// SampleUniformWeighted wraps a Sampler type to create a WeightedSampler where all
// weights are equal.
type SampleUniformWeighted struct {
	Sampler
}

// SampleWeighted generates len(batch) samples from the embedded Sampler type
// and sets all of the weights equal to 1. If len(batch) and len(weights)
// are not equal, SampleWeighted will panic.
func (w SampleUniformWeighted) SampleWeighted(batch, weights []float64) {
	if len(batch) != len(weights) {
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
	Q   distuv.Quantiler
	Src *rand.Rand
}

// Sample generates len(batch) samples using the LatinHypercube generation
// procedure.
func (l LatinHypercuber) Sample(batch []float64) {
	LatinHypercube(batch, l.Q, l.Src)
}

// LatinHypercube generates len(batch) samples using Latin hypercube sampling
// from the given distribution. If src != nil, it will be used to generate
// random numbers, otherwise rand.Float64 will be used.
//
// Latin hypercube sampling divides the cumulative distribution function into equally
// spaced bins and guarantees that one sample is generated per bin. Within each bin,
// the location is randomly sampled. The distuv.UnitUniform variable can be used
// for easy generation from the unit interval.
func LatinHypercube(batch []float64, q distuv.Quantiler, src *rand.Rand) {
	n := len(batch)
	var perm []int
	var f64 func() float64
	if src != nil {
		f64 = src.Float64
		perm = src.Perm(n)
	} else {
		f64 = rand.Float64
		perm = rand.Perm(n)
	}
	for i := range batch {
		v := f64()/float64(n) + float64(i)/float64(n)
		batch[perm[i]] = q.Quantile(v)
	}
}

// Importancer is a wrapper around the Importance sampling generation method.
type Importancer struct {
	Target   distuv.LogProber
	Proposal distuv.RandLogProber
}

// Sample generates len(batch) samples using the Importance sampling generation
// procedure.
func (l Importancer) SampleWeighted(batch, weights []float64) {
	Importance(batch, weights, l.Target, l.Proposal)
}

// Importance sampling generates len(batch) samples from the proposal distribution,
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
func Importance(batch, weights []float64, target distuv.LogProber, proposal distuv.RandLogProber) {
	if len(batch) != len(weights) {
		panic(badLengthMismatch)
	}
	for i := range batch {
		v := proposal.Rand()
		batch[i] = v
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
	Target   distuv.LogProber
	Proposal distuv.RandLogProber
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

// Sample generates len(batch) using the Rejection sampling generation procedure.
// Rejection sampling may fail if the constant is insufficiently high, as described
// in the function comment for Rejection. If the generation fails, the samples
// are set to math.NaN(), and a call to Err will return a non-nil value.
func (r *Rejectioner) Sample(batch []float64) {
	r.err = nil
	r.proposed = 0
	proposed, ok := Rejection(batch, r.Target, r.Proposal, r.C, r.Src)
	if !ok {
		r.err = ErrRejection
	}
	r.proposed = proposed
}

// Rejection generates len(batch) samples using the rejection sampling algorithm
// and stores them in place into samples. Sampling continues until batch is
// filled. Rejection returns the total number of proposed locations and a boolean
// indicating if the rejection sampling assumption is violated (see details
// below). If the returned boolean is false, all elements of samples are set to
// NaN. If src is not nil, it will be used to generate random numbers, otherwise
// rand.Float64 will be used.
//
// Rejection sampling generates points from the target distribution by using
// the proposal distribution. At each step of the algorithm, the proposed point
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
func Rejection(batch []float64, target distuv.LogProber, proposal distuv.RandLogProber, c float64, src *rand.Rand) (nProposed int, ok bool) {
	if c < 1 {
		panic("rejection: acceptance constant must be greater than 1")
	}
	f64 := rand.Float64
	if src != nil {
		f64 = src.Float64
	}
	var idx int
	for {
		nProposed++
		v := proposal.Rand()
		qx := proposal.LogProb(v)
		px := target.LogProb(v)
		accept := math.Exp(px-qx) / c
		if accept > 1 {
			// Invalidate the whole result and return a failure.
			for i := range batch {
				batch[i] = math.NaN()
			}
			return nProposed, false
		}
		if accept > f64() {
			batch[idx] = v
			idx++
			if idx == len(batch) {
				break
			}
		}
	}
	return nProposed, true
}

// MHProposal defines a proposal distribution for Metropolis Hastings.
type MHProposal interface {
	// ConditionalDist returns the probability of the first argument conditioned on
	// being at the second argument
	//  p(x|y)
	ConditionalLogProb(x, y float64) (prob float64)

	// ConditionalRand generates a new random location conditioned being at the
	// location y.
	ConditionalRand(y float64) (x float64)
}

// MetropolisHastingser is a wrapper around the MetropolisHastings sampling type.
//
// BurnIn sets the number of samples to discard before keeping the first sample.
// A properly set BurnIn rate will decorrelate the sampling chain from the initial
// location. The proper BurnIn value will depend on the mixing time of the
// Markov chain defined by the target and proposal distributions.
//
// Rate sets the number of samples to discard in between each kept sample. A
// higher rate will better approximate independently and identically distributed
// samples, while a lower rate will keep more information (at the cost of
// higher correlation between samples). If Rate is 0 it is defaulted to 1.
//
// The initial value is NOT changed during calls to Sample.
type MetropolisHastingser struct {
	Initial  float64
	Target   distuv.LogProber
	Proposal MHProposal
	Src      *rand.Rand

	BurnIn int
	Rate   int
}

// Sample generates len(batch) samples using the Metropolis Hastings sample
// generation method. The initial location is NOT updated during the call to Sample.
func (m MetropolisHastingser) Sample(batch []float64) {
	burnIn := m.BurnIn
	rate := m.Rate
	if rate == 0 {
		rate = 1
	}

	// Use the optimal size for the temporary memory to allow the fewest calls
	// to MetropolisHastings. The case where tmp shadows samples must be
	// aligned with the logic after burn-in so that tmp does not shadow samples
	// during the rate portion.
	tmp := batch
	if rate > len(batch) {
		tmp = make([]float64, rate)
	}

	// Perform burn-in.
	remaining := burnIn
	initial := m.Initial
	for remaining != 0 {
		newSamp := min(len(tmp), remaining)
		MetropolisHastings(tmp[newSamp:], initial, m.Target, m.Proposal, m.Src)
		initial = tmp[newSamp-1]
		remaining -= newSamp
	}

	if rate == 1 {
		MetropolisHastings(batch, initial, m.Target, m.Proposal, m.Src)
		return
	}

	if len(tmp) <= len(batch) {
		tmp = make([]float64, rate)
	}

	// Take a single sample from the chain
	MetropolisHastings(batch[0:1], initial, m.Target, m.Proposal, m.Src)
	initial = batch[0]

	// For all of the other samples, first generate Rate samples and then actually
	// accept the last one.
	for i := 1; i < len(batch); i++ {
		MetropolisHastings(tmp, initial, m.Target, m.Proposal, m.Src)
		v := tmp[rate-1]
		batch[i] = v
		initial = v
	}
}

// MetropolisHastings generates len(batch) samples using the Metropolis Hastings
// algorithm (http://en.wikipedia.org/wiki/Metropolis%E2%80%93Hastings_algorithm),
// with the given target and proposal distributions, starting at the intial location
// and storing the results in-place into samples. If src != nil, it will be used to generate random
// numbers, otherwise rand.Float64 will be used.
//
// Metropolis-Hastings is a Markov-chain Monte Carlo algorithm that generates
// samples according to the distribution specified by target by using the Markov
// chain implicitly defined by the proposal distribution. At each
// iteration, a proposal point is generated randomly from the current location.
// This proposal point is accepted with probability
//  p = min(1, (target(new) * proposal(current|new)) / (target(current) * proposal(new|current)))
// If the new location is accepted, it is stored into batch and becomes the
// new current location. If it is rejected, the current location remains and
// is stored into samples. Thus, a location is stored into batch at every iteration.
//
// The samples in Metropolis Hastings are correlated with one another through the
// Markov chain. As a result, the initial value can have a significant influence
// on the early samples, and so, typically, the first samples generated by the chain
// are ignored. This is known as "burn-in", and can be accomplished with slicing.
// The best choice for burn-in length will depend on the sampling and target
// distributions.
//
// Many choose to have a sampling "rate" where a number of samples
// are ignored in between each kept sample. This helps decorrelate
// the samples from one another, but also reduces the number of available samples.
// A sampling rate can be implemented with successive calls to MetropolisHastings.
func MetropolisHastings(batch []float64, initial float64, target distuv.LogProber, proposal MHProposal, src *rand.Rand) {
	f64 := rand.Float64
	if src != nil {
		f64 = src.Float64
	}
	current := initial
	currentLogProb := target.LogProb(initial)
	for i := range batch {
		proposed := proposal.ConditionalRand(current)
		proposedLogProb := target.LogProb(proposed)
		probTo := proposal.ConditionalLogProb(proposed, current)
		probBack := proposal.ConditionalLogProb(current, proposed)

		accept := math.Exp(proposedLogProb + probBack - probTo - currentLogProb)
		if accept > f64() {
			current = proposed
			currentLogProb = proposedLogProb
		}
		batch[i] = current
	}
}

// IIDer is a wrapper around the IID sample generation method.
type IIDer struct {
	Dist distuv.Rander
}

// Sample generates a set of identically and independently distributed samples.
func (iid IIDer) Sample(batch []float64) {
	IID(batch, iid.Dist)
}

// IID generates a set of independently and identically distributed samples from
// the input distribution.
func IID(batch []float64, d distuv.Rander) {
	for i := range batch {
		batch[i] = d.Rand()
	}
}
