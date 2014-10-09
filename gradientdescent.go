// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "github.com/gonum/floats"

// GradientDescent is a Method that performs gradient-based optimization. Gradient
// Descent performs successive steps along the direction of the gradient. The
// LinesearchMethod specifies the kind of linesearch to be done, and StepSizer determines
// the initial step size of each direction. If NoLinesearch is true, no linesearch
// will be done. If either LinesearchMethod or StepSizer are nil, a reasonable
// value will be chosen.
type GradientDescent struct {
	LinesearchMethod LinesearchMethod
	StepSizer        StepSizer
	NoLinesearch     bool

	linesearch *Linesearch
}

func (g *GradientDescent) Init(l Location, f *FunctionStats, xNext []float64) (EvaluationType, IterationType, error) {
	if g.StepSizer == nil {
		g.StepSizer = ConstantStepSize{1}
	}
	if g.NoLinesearch {
		stepSize := g.StepSizer.Init(l)
		floats.AddScaledTo(xNext, l.X, -stepSize, l.Gradient)
		return FunctionAndGradient, Major, nil
	}
	if g.LinesearchMethod == nil {
		g.LinesearchMethod = &Backtracking{}
	}
	if g.linesearch == nil {
		g.linesearch = &Linesearch{}
	}
	g.linesearch.Method = g.LinesearchMethod
	g.linesearch.NextDirectioner = g

	return g.linesearch.Init(l, f, xNext)
}

func (g *GradientDescent) Iterate(l Location, xNext []float64) (EvaluationType, IterationType, error) {
	if g.NoLinesearch {
		stepSize := g.StepSizer.StepSize(l)
		floats.AddScaledTo(xNext, l.X, -stepSize, l.Gradient)
		return FunctionAndGradient, Major, nil
	}

	return g.linesearch.Iterate(l, xNext)
}

func (g *GradientDescent) InitDirection(l Location, direction []float64) (stepSize float64) {
	copy(direction, l.Gradient)
	floats.Scale(-1, direction)
	return g.StepSizer.Init(l)
}

func (g *GradientDescent) NextDirection(l Location, direction []float64) (stepSize float64) {
	copy(direction, l.Gradient)
	floats.Scale(-1, direction)
	return g.StepSizer.StepSize(l)
}
