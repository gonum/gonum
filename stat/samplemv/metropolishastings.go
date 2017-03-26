// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package samplemv

import (
	"math"
	"math/rand"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat/distmv"
)

var _ Sampler = MetropolisHastingser{}

// MHProposal defines a proposal distribution for Metropolis Hastings.
type MHProposal interface {
	// ConditionalLogProb returns the probability of the first argument
	// conditioned on being at the second argument.
	//  p(x|y)
	// ConditionalLogProb panics if the input slices are not the same length.
	ConditionalLogProb(x, y []float64) (prob float64)

	// ConditionalRand generates a new random location conditioned being at the
	// location y. If the first arguement is nil, a new slice is allocated and
	// returned. Otherwise, the random location is stored in-place into the first
	// argument, and ConditionalRand will panic if the input slice lengths differ.
	ConditionalRand(x, y []float64) []float64
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
	Initial  []float64
	Target   distmv.LogProber
	Proposal MHProposal
	Src      *rand.Rand

	BurnIn int
	Rate   int
}

// Sample generates rows(batch) samples using the Metropolis Hastings sample
// generation method. The initial location is NOT updated during the call to Sample.
//
// The number of columns in batch must equal len(m.Initial), otherwise Sample
// will panic.
func (m MetropolisHastingser) Sample(batch *mat64.Dense) {
	burnIn := m.BurnIn
	rate := m.Rate
	if rate == 0 {
		rate = 1
	}
	r, c := batch.Dims()
	if len(m.Initial) != c {
		panic("metropolishastings: length mismatch")
	}

	// Use the optimal size for the temporary memory to allow the fewest calls
	// to MetropolisHastings. The case where tmp shadows samples must be
	// aligned with the logic after burn-in so that tmp does not shadow samples
	// during the rate portion.
	tmp := batch
	if rate > r {
		tmp = mat64.NewDense(rate, c, nil)
	}
	rTmp, _ := tmp.Dims()

	// Perform burn-in.
	remaining := burnIn
	initial := make([]float64, c)
	copy(initial, m.Initial)
	for remaining != 0 {
		newSamp := min(rTmp, remaining)
		MetropolisHastings(tmp.View(0, 0, newSamp, c).(*mat64.Dense), initial, m.Target, m.Proposal, m.Src)
		copy(initial, tmp.RawRowView(newSamp-1))
		remaining -= newSamp
	}

	if rate == 1 {
		MetropolisHastings(batch, initial, m.Target, m.Proposal, m.Src)
		return
	}

	if rTmp <= r {
		tmp = mat64.NewDense(rate, c, nil)
	}

	// Take a single sample from the chain.
	MetropolisHastings(batch.View(0, 0, 1, c).(*mat64.Dense), initial, m.Target, m.Proposal, m.Src)

	copy(initial, batch.RawRowView(0))
	// For all of the other samples, first generate Rate samples and then actually
	// accept the last one.
	for i := 1; i < r; i++ {
		MetropolisHastings(tmp, initial, m.Target, m.Proposal, m.Src)
		v := tmp.RawRowView(rate - 1)
		batch.SetRow(i, v)
		copy(initial, v)
	}
}

// MetropolisHastings generates rows(batch) samples using the Metropolis Hastings
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
func MetropolisHastings(batch *mat64.Dense, initial []float64, target distmv.LogProber, proposal MHProposal, src *rand.Rand) {
	f64 := rand.Float64
	if src != nil {
		f64 = src.Float64
	}
	if len(initial) == 0 {
		panic("metropolishastings: zero length initial")
	}
	r, _ := batch.Dims()
	current := make([]float64, len(initial))
	copy(current, initial)
	proposed := make([]float64, len(initial))
	currentLogProb := target.LogProb(initial)
	for i := 0; i < r; i++ {
		proposal.ConditionalRand(proposed, current)
		proposedLogProb := target.LogProb(proposed)
		probTo := proposal.ConditionalLogProb(proposed, current)
		probBack := proposal.ConditionalLogProb(current, proposed)

		accept := math.Exp(proposedLogProb + probBack - probTo - currentLogProb)
		if accept > f64() {
			copy(current, proposed)
			currentLogProb = proposedLogProb
		}
		batch.SetRow(i, current)
	}
}

// ProposalNormal is a sampling distribution for Metropolis-Hastings. It has a
// fixed covariance matrix and changes the mean based on the current sampling
// location.
type ProposalNormal struct {
	normal *distmv.Normal
}

// NewProposalNormal constructs a new ProposalNormal for use as a proposal
// distribution for Metropolis-Hastings. ProposalNormal is a multivariate normal
// distribution (implemented by distmv.Normal) where the covariance matrix is fixed
// and the mean of the distribution changes.
//
// NewProposalNormal returns {nil, false} if the covariance matrix is not positive-definite.
func NewProposalNormal(sigma *mat64.SymDense, src *rand.Rand) (*ProposalNormal, bool) {
	mu := make([]float64, sigma.Symmetric())
	normal, ok := distmv.NewNormal(mu, sigma, src)
	if !ok {
		return nil, false
	}
	p := &ProposalNormal{
		normal: normal,
	}
	return p, true
}

// ConditionalLogProb returns the probability of the first argument conditioned on
// being at the second argument.
//  p(x|y)
// ConditionalLogProb panics if the input slices are not the same length or
// are not equal to the dimension of the covariance matrix.
func (p *ProposalNormal) ConditionalLogProb(x, y []float64) (prob float64) {
	// Either SetMean or LogProb will panic if the slice lengths are innaccurate.
	p.normal.SetMean(y)
	return p.normal.LogProb(x)
}

// ConditionalRand generates a new random location conditioned being at the
// location y. If the first arguement is nil, a new slice is allocated and
// returned. Otherwise, the random location is stored in-place into the first
// argument, and ConditionalRand will panic if the input slice lengths differ or
// if they are not equal to the dimension of the covariance matrix.
func (p *ProposalNormal) ConditionalRand(x, y []float64) []float64 {
	if x == nil {
		x = make([]float64, p.normal.Dim())
	}
	if len(x) != len(y) {
		panic(badLengthMismatch)
	}
	p.normal.SetMean(y)
	p.normal.Rand(x)
	return x
}
