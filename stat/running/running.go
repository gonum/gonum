// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package running

import "math"

// Mean helps compute running mean
type Mean struct {
	// Decay sets the decay rate. A value of 0.99 means that the effective count
	// is multiplied by 0.99 for every added value. If Decay is zero, a default
	// value of 1 is used (the count is not decayed).
	Decay     float64
	InitCount float64
	InitMean  float64

	count float64
	mean  float64
}

// Reset resets the counter. The next accumulate will use the
// initial mean and count.
func (m *Mean) Reset() {
	m.count = 0
	m.mean = 0
}

func (m *Mean) isZero() bool {
	return m.count == 0
}

// Mean returns the current estimate of the running mean.
func (m *Mean) Mean() float64 {
	if m.isZero() {
		return m.InitMean
	}
	return m.mean
}

// Count returns the current count of values.
func (m *Mean) Count() float64 {
	if m.isZero() {
		return m.InitCount
	}
	return m.count
}

// Accum adds the value to the running total.
func (m *Mean) Accum(v float64) {
	m.AccumWeighted(v, 1)
}

// AccumWeighted adds the weighted value to the running total.
func (m *Mean) AccumWeighted(v, weight float64) {
	if weight < 0 {
		panic("running: negative weight")
	}
	if m.isZero() && weight == 0 {
		return
	}
	decay := m.decay()
	decay = math.Pow(decay, weight)
	m.count *= decay
	m.mean = (m.mean*m.count + v*weight) / (m.count + weight)
	m.count += weight
}

func (m *Mean) decay() float64 {
	if m.Decay == 0 {
		return 1
	}
	return m.Decay
}
