// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sampleuv

import "gonum.org/v1/gonum/stat/distuv"

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
	// one sample from every few is accepted from the chain.
	rate := 50

	mh := MetropolisHastings{
		Initial:  initial,
		Target:   target,
		Proposal: proposal,
		BurnIn:   burnin,
		Rate:     rate,
	}

	samples := make([]float64, n)
	mh.Sample(samples)
}
