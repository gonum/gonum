// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"

	"github.com/gonum/floats"
)

// ConstantStepSize is a StepSizer that returns the same step size for
// every iteration.
type ConstantStepSize struct {
	Size float64
}

func (c ConstantStepSize) Init(l Location, dir []float64) float64 {
	return c.Size
}

func (c ConstantStepSize) StepSize(l Location, dir []float64) float64 {
	return c.Size
}

// QuadraticInterpolateStepSize estimates the initial step size for line search
// as the minimum of a quadratic that interpolates f(x_{k-1}), f(x_k) and grad
// f_{k-1} \dot p_k. This is useful for line search methods that do not produce
// well-scaled descent directions, such as gradient descent or conjugate
// gradient methods.
//
// See also Nocedal, Wright (2006), Numerical Optimization (2nd ed.), sec.
// 3.5, page 59.
type QuadraticInterpolateStepSize struct {
	fPrev float64
}

func (q *QuadraticInterpolateStepSize) Init(l Location, dir []float64) float64 {
	q.fPrev = l.F
	return math.Min(1, 1/floats.Norm(l.Gradient, 2))
}

func (q *QuadraticInterpolateStepSize) StepSize(l Location, dir []float64) (s float64) {
	s = 2 * (l.F - q.fPrev) / floats.Dot(l.Gradient, dir)
	s = math.Min(1, 1.01*s)
	q.fPrev = l.F
	return s
}
