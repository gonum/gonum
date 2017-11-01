// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "math"

// FunctionConverge tests for the convergence of function values. See comment
// in Settings.
type FunctionConverge struct {
	Absolute   float64
	Relative   float64
	Iterations int

	best float64
	iter int
}

func (fc *FunctionConverge) Init(f float64) {
	fc.best = f
	fc.iter = 0
}

func (fc *FunctionConverge) FunctionConverged(f float64) Status {
	if fc.Iterations == 0 {
		return NotTerminated
	}
	maxAbs := math.Max(math.Abs(f), math.Abs(fc.best))
	if f < fc.best && fc.best-f > fc.Relative*maxAbs+fc.Absolute {
		fc.best = f
		fc.iter = 0
		return NotTerminated
	}
	fc.iter++
	if fc.iter < fc.Iterations {
		return NotTerminated
	}
	return FunctionConvergence
}
