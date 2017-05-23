// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sampleuv

import "github.com/gonum/stat/distuv"

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func ExampleMetropolisHastings_samplingRate() {
	// See Burnin example for a description of these quantities.
	n := 1000
	burnin := 300
	var initial float64
	target := distuv.Weibull{K: 5, Lambda: 0.5}
	proposal := ProposalDist{Sigma: 0.2}

	// Successive samples are correlated with one another through the
	// Markov Chain defined by the proposal distribution. To get less
	// correlated samples, one may use a sampling rate, in which only
	// one sample from every few is accepted from the chain. This can
	// be accomplished through a for loop.
	rate := 50

	tmp := make([]float64, max(rate, burnin))

	// First deal with burnin.
	tmp = tmp[:burnin]
	MetropolisHastings(tmp, initial, target, proposal, nil)
	// The final sample in tmp in the final point in the chain.
	// Use it as the new initial location.
	initial = tmp[len(tmp)-1]

	// Now, generate samples by using one every rate samples.
	tmp = tmp[:rate]
	samples := make([]float64, n)
	samples[0] = initial
	for i := 1; i < len(samples); i++ {
		MetropolisHastings(tmp, initial, target, proposal, nil)
		initial = tmp[len(tmp)-1]
		samples[i] = initial
	}
}
