// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

// A Method can optimize an objective function.
type Method interface {
	// Init initializes the method and stores the first location to evaluate
	// in xNext.
	Init(loc *Location, xNext []float64) (EvaluationType, IterationType, error)

	// Iterate performs one iteration of the method and stores the next
	// location to evaluate in xNext.
	Iterate(loc *Location, xNext []float64) (EvaluationType, IterationType, error)

	// Needs specifies information about the objective function needed by the
	// optimizer beyond just the function value. The information is used
	// internally for initialization and must match evaluation types returned
	// by Init() and Iterate() during the optimization process.
	Needs() struct {
		Gradient bool
		Hessian  bool
	}
}

// Linesearcher is a type that can perform a line search. It tries to find an
// (approximate) minimum of the objective function along the search direction
// dir_k starting at the most recent location x_k, i.e., it tries to minimize
// the function
//  φ(step) := f(x_k + step * dir_k) where step > 0.
// Typically, a Linesearcher will be used in conjuction with LinesearchMethod
// for performing gradient-based optimization through sequential line searches.
type Linesearcher interface {
	// Init initializes the linesearch method. Value and derivative contain
	// φ(0) and φ'(0), respectively, and step contains the first trial step
	// length. It returns the type of evaluation to be performed at
	// x_0 + step * dir_0.
	Init(value, derivative float64, step float64) EvaluationType

	// Finished takes in the values of φ and φ' evaluated at the previous step,
	// and returns whether a sufficiently accurate minimum of φ has been found.
	Finished(value, derivative float64) bool

	// Iterate takes in the values of φ and φ' evaluated at the previous step
	// and returns the next step size and the type of evaluation to be
	// performed at x_k + step * dir_k.
	Iterate(value, derivative float64) (step float64, e EvaluationType, err error)
}

// NextDirectioner implements a strategy for computing a new line search
// direction at each major iteration. Typically, a NextDirectioner will be
// used in conjuction with LinesearchMethod for performing gradient-based
// optimization through sequential line searches.
type NextDirectioner interface {
	// InitDirection initializes the NextDirectioner at the given starting location,
	// putting the initial direction in place into dir, and returning the initial
	// step size. InitDirection must not modify Location.
	InitDirection(loc *Location, dir []float64) (step float64)

	// NextDirection updates the search direction and step size. Location is
	// the location seen at the conclusion of the most recent linesearch. The
	// next search direction is put in place into dir, and the next step size
	// is returned. NextDirection must not modify Location.
	NextDirection(loc *Location, dir []float64) (step float64)
}

// StepSizer can set the next step size of the optimization given the last Location.
// Returned step size must be positive.
type StepSizer interface {
	Init(loc *Location, dir []float64) float64
	StepSize(loc *Location, dir []float64) float64
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
	Init() error
	Record(*Location, EvaluationType, IterationType, *Stats) error
}
