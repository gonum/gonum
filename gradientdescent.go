// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "github.com/gonum/floats"

// GradientDescent is a Method that performs gradient-based optimization. Gradient
// Descent performs successive steps along the direction of the gradient. The
// LinesearchMethod specifies the kind of linesearch to be done, and StepSizer determines
// the initial step size of each direction. If either LinesearchMethod or StepSizer
// are nil, a reasonable value will be chosen.
type GradientDescent struct {
	LinesearchMethod LinesearchMethod
	StepSizer        StepSizer

	linesearch *Linesearch
}

func (g *GradientDescent) Init(loc *Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	if g.StepSizer == nil {
		g.StepSizer = &QuadraticStepSize{}
	}
	if g.LinesearchMethod == nil {
		g.LinesearchMethod = &Backtracking{}
	}
	if g.linesearch == nil {
		g.linesearch = &Linesearch{}
	}
	g.linesearch.Method = g.LinesearchMethod
	g.linesearch.NextDirectioner = g

	return g.linesearch.Init(loc, f, xNext)
}

func (g *GradientDescent) Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error) {
	return g.linesearch.Iterate(loc, xNext)
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
