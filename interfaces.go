// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

// Function evaluates the objective function at the given location. F
// must not modify x.
type Function interface {
	F(x []float64) (obj float64)
}

// Gradient evaluates the gradient at x and stores the result in place. Df
// must not modify x.
type Gradient interface {
	Df(x []float64, grad []float64)
}

// FunctionGradient evaluates both the derivative and the function at x, storing
// the gradient in place. FDf must not modify x.
type FunctionGradient interface {
	FDf(x []float64, grad []float64) (obj float64)
}

// LinesearchMethod is a type that can perform a line search. Typically, these
// methods will not be called by the user directly, as they will be called by
// a Linesearch struct.
type LinesearchMethod interface {
	// Init initializes the linesearch method. LinesearchLocation contains the
	// function information at step == 0, and step contains the first step length
	// as specified by the NextDirectioner.
	Init(init LinesearchLocation, step float64, f *FunctionInfo) EvaluationType

	// Finished takes in the function result at the most recent linesearch location,
	// and returns true if the line search has been concluded.
	Finished(l LinesearchLocation) bool

	// Iterate takes in the function results
	// from evaluating the function at the previous step, and returns the
	// next step size and EvaluationType to evaluate.
	Iterate(l LinesearchLocation) (nextStep float64, e EvaluationType, err error)
}

// NextDirectioner implements a strategy for computing a new line search direction
// at each major iteration. Typically, these methods will not be called by the user directly,
// as they will be called by a Linesearch struct.
type NextDirectioner interface {
	// InitDirection initializes the NextDirectioner at the given starting location,
	// putting the initial direction in place into dir, and returning the initial
	// step size. InitDirection must not modify Location.
	InitDirection(l Location, dir []float64) (stepSize float64)

	// NextDirection updates the search direction and step size. Location is
	// the location seen at the conclusion of the most recent linesearch. The
	// next search direction is put in place into dir, and the next stepsize
	// is returned. NextDirection must not modify Location.
	NextDirection(l Location, dir []float64) (stepSize float64)
}

// A Method can optimize an objective function.
type Method interface {
	// Initializes the method and returns the first location to evaluate
	Init(l Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error)

	// Stores the next location to evaluate in xNext
	Iterate(l Location, xNext []float64) (EvaluationType, IterationType, error)
}

// StepSizer can set the next step size of the optimization given the last Location.
// Returned step size must be positive.
type StepSizer interface {
	Init(l Location, dir []float64) float64
	StepSize(l Location, dir []float64) float64
}

// Statuser returns the status of the Function being optimized. This can be used
// by the Function to terminate early, for example with an error. The user can
// use one of the pre-provided Status constants, or may call NewStatus to create
// a custom Status value.
type Statuser interface {
	Status() (Status, error)
}

// A Recorder can record the progress of the optimization, for example to print
// the progress to StdOut or to a log file. A Recorder must not modify any data.
type Recorder interface {
	Init(*FunctionInfo) error
	Record(Location, EvaluationType, IterationType, *Stats) error
}
