// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

// LogProber calculates the log-probability for a given value of a random variable.
type LogProber interface {
	LogProb(float64) float64
}

// Rander samples a random variable.
type Rander interface {
	Rand() float64
}

// RandLogProber is a Rander and a LogProber.
type RandLogProber interface {
	Rander
	LogProber
}

// Quantiler calculates the quantile of a distribution for given probability.
type Quantiler interface {
	Quantile(p float64) float64
}
