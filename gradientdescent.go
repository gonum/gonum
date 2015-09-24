// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "github.com/gonum/floats"

// GradientDescent is a Method that performs gradient-based optimization.
// Gradient Descent performs successive steps along the direction of the
// gradient. The Linesearcher specifies the kind of linesearch to be done, and
// StepSizer determines the initial step size of each direction. If either
// Linesearcher or StepSizer are nil, a reasonable value will be chosen.
type GradientDescent struct {
	Linesearcher Linesearcher
	StepSizer    StepSizer

	ls *LinesearchMethod
}

func (g *GradientDescent) Init(loc *Location) (Operation, error) {
	if g.StepSizer == nil {
		g.StepSizer = &QuadraticStepSize{}
	}
	if g.Linesearcher == nil {
		g.Linesearcher = &Backtracking{}
	}
	if g.ls == nil {
		g.ls = &LinesearchMethod{}
	}
	g.ls.Linesearcher = g.Linesearcher
	g.ls.NextDirectioner = g

	return g.ls.Init(loc)
}

func (g *GradientDescent) Iterate(loc *Location) (Operation, error) {
	return g.ls.Iterate(loc)
}

func (g *GradientDescent) InitDirection(loc *Location, dir []float64) (stepSize float64) {
	copy(dir, loc.Gradient)
	floats.Scale(-1, dir)
	return g.StepSizer.Init(loc, dir)
}

func (g *GradientDescent) NextDirection(loc *Location, dir []float64) (stepSize float64) {
	copy(dir, loc.Gradient)
	floats.Scale(-1, dir)
	return g.StepSizer.StepSize(loc, dir)
}

func (*GradientDescent) Needs() struct {
	Gradient bool
	Hessian  bool
} {
	return struct {
		Gradient bool
		Hessian  bool
	}{true, false}
}
