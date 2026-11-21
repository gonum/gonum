// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package running

import "math"

// Mean is a running mean accumulator.
type Mean struct {
	decay float64
	count float64
	mean  float64
}

// NewMean returns a new accumulator for computing the running mean,
// with the internal mean and count initialized to zero.
// decay sets the decay rate for the internal sample counter. The count
// is multiplied by the decay value for each sample added.
func NewMean(decay float64) *Mean {
	return &Mean{decay: decay}
}

// NewMeanInitialized returns a new Mean accumulator with initialized
// values. decay sets the decay rate for the internal sample counter. The count
// is multiplied by the decay value for each sample added. mean and count
// specify the initial values for internal running mean and running count
// respectively.
func NewMeanInitialized(decay, mean, count float64) *Mean {
	return &Mean{
		decay: decay,
		mean:  mean,
		count: count,
	}
}

// Mean returns the current estimate of the running mean.
// Note that if the count is zero this returns the value
// of the initial mean.
func (m *Mean) Mean() float64 {
	return m.mean
}

// Count returns the current count of values.
func (m *Mean) Count() float64 {
	return m.count
}

// Accum adds the value to the running total.
func (m *Mean) Accum(v float64) {
	m.AccumWeighted(v, 1)
}

// AccumWeighted adds the weighted value to the running total.
// Weights must be positive.
func (m *Mean) AccumWeighted(v, weight float64) {
	if weight < 0 {
		panic("running: negative weight")
	}
	if m.count == 0 && weight == 0 {
		// Avoid NaN.
		return
	}
	decay := m.decay
	decay = math.Pow(decay, weight)
	m.count *= decay
	m.mean = (m.mean*m.count + v*weight) / (m.count + weight)
	m.count += weight
}
