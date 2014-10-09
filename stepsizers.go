// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

// ConstantStepSize is a StepSizer that returns the same step size for
// every iteration.
type ConstantStepSize struct {
	Size float64
}

func (c ConstantStepSize) Init(l Location) float64 {
	return c.Size
}

func (c ConstantStepSize) StepSize(l Location) float64 {
	return c.Size
}
