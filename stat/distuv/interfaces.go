// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

// LogProber interface
type LogProber interface {
	LogProb(float64) float64
}

// Rander interface
type Rander interface {
	Rand() float64
}

// RandLogProber interface
type RandLogProber interface {
	Rander
	LogProber
}

// Quantiler interface
type Quantiler interface {
	Quantile(p float64) float64
}
