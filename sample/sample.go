// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package sample contains a set of advanced routines for sampling from
// probability distributions.
package sample

import (
	"math"
	"math/rand"

	"github.com/gonum/stat/dist"
)

var (
	badLengthMismatch = "sample: slice length mismatch"
)

// LatinHypercube generates len(samples) samples using Latin hypercube sampling
// from the given distribution. If src != nil, it will be used to generate
// random numbers, otherwise rand.Float64 will be used.
//
// Latin hypercube sampling divides the cumulative distribution function into equally
// spaced bins and guarantees that one sample is generated per bin. Within each bin,
// the location is randomly sampled. The dist.UnitNormal variable can be used
// for easy generation from the unit interval.
func LatinHypercube(samples []float64, q dist.Quantiler, src *rand.Rand) {
	n := len(samples)
	var perm []int
	var f64 func() float64
	if src != nil {
		f64 = src.Float64
		perm = src.Perm(n)
	} else {
		f64 = rand.Float64
		perm = rand.Perm(n)
	}
	for i := range samples {
		v := f64()/float64(n) + float64(i)/float64(n)
		samples[perm[i]] = q.Quantile(v)
	}
}

// Importance sampling generates len(x) samples from the proposal distribution,
// and stores the locations and importance sampling weights in place.
//
// Importance sampling is a variance reduction technique where samples are
// generated from a proposal distribution, q(x), instead of the target distribution
// p(x). This allows relatively unlikely samples in p(x) to be generated more frequently
//
// The importance sampling weight at x is given by p(x)/q(x). To reduce variance,
// a good proposal distribution will bound this sampling weight. This implies the
// support of q(x) should be at least as broad as p(x), and q(x) should be "fatter tailed"
// than p(x).
func Importance(samples, weights []float64, target dist.LogProber, proposal dist.RandLogProber) {
	if len(samples) != len(weights) {
		panic(badLengthMismatch)
	}
	for i := range samples {
		v := proposal.Rand()
		samples[i] = v
		weights[i] = math.Exp(target.LogProb(v) - proposal.LogProb(v))
	}
}

// Rejection generates len(x) samples using the rejection sampling algorithm and
// stores them in place into samples.
// Sampling continues until x is filled. Rejection the total number of proposed
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
func Rejection(samples []float64, target dist.LogProber, proposal dist.RandLogProber, c float64, src *rand.Rand) (nProposed int, ok bool) {
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
			for i := range samples {
				samples[i] = math.NaN()
			}
			return nProposed, false
		}
		if accept > f64() {
			samples[idx] = v
			idx++
			if idx == len(samples) {
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

// MetropolisHastings generates len(samples) samples using the Metropolis Hastings
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
// If the new location is accepted, it is stored into samples and becomes the
// new current location. If it is rejected, the current location remains and
// is stored into samples. Thus, a location is stored into samples at every iteration.
//
// The samples in Metropolis Hastings are correlated with one another through the
// Markov-Chain. As a result, the initial value can have a significant influence
// on the early samples, and so typically, the first sapmles generated by the chain.
// are ignored. This is known as "burn-in", and can be accomplished with slicing.
// The best choice for burn-in length will depend on the sampling and the target
// distribution.
//
// Many choose to have a sampling "rate" where a number of samples
// are ignored in between each kept sample. This helps decorrelate
// the samples from one another, but also reduces the number of available samples.
// A sampling rate can be implemented with successive calls to MetropolisHastings.
func MetropolisHastings(samples []float64, initial float64, target dist.LogProber, proposal MHProposal, src *rand.Rand) {
	f64 := rand.Float64
	if src != nil {
		f64 = src.Float64
	}
	current := initial
	currentLogProb := target.LogProb(initial)
	for i := range samples {
		proposed := proposal.ConditionalRand(current)
		proposedLogProb := target.LogProb(proposed)
		probTo := proposal.ConditionalLogProb(proposed, current)
		probBack := proposal.ConditionalLogProb(current, proposed)

		accept := math.Exp(proposedLogProb + probBack - probTo - currentLogProb)
		if accept > f64() {
			current = proposed
			currentLogProb = proposedLogProb
		}
		samples[i] = current
	}
}
