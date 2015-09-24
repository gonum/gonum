// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

// A Method can optimize an objective function.
//
// It uses a reverse-communication interface between the optimization method
// and the caller. Method acts as a client that asks the caller to perform
// needed operations via RequestType returned from Init and Iterate methods.
// This provides independence of the optimization algorithm on user-supplied
// data and their representation, and enables automation of common operations
// like checking for (various types of) convergence and maintaining statistics.
//
// A Method can command an Evaluation, a MajorIteration or NoOperation operations.
// An evaluation operation is one or more of the Evaluation operations
// (FuncEvaluation, GradEvaluation, etc.) which can be combined with
// the bitwise or operator. In an evaluation operation, the requested routines
// will be evaluated at the point specified in Location.X. The corresponding
// fields of Location will be filled with the results from the routine and can
// be retrieved upon the next call to Iterate. Alternatively, a Method can
// declare a MajorIteration. In a MajorIteration, all values in Location must
// be valid and consistent, and are interpreted as a new minimum. Convergence
// of the optimization (GradientThreshold, etc.) will be checked using this new
// minimum.
//
// A Method must not return InitIteration and PostIteration operations. These are
// reserved for the clients to be passed to Recorders. A Method must also not
// combine the Evaluation operations with the Iteration operations.
type Method interface {
	// Init initializes the method based on the initial data in loc, updates it
	// and returns the first operation to be carried out by the caller.
	// The initial location must be valid as specified by Needs().
	Init(loc *Location) (Operation, error)

	// Iterate retrieves data from loc, performs one iteration of the method,
	// updates loc and returns the next operation.
	// TODO(vladimir-ch): When decided, say something whether the contents of
	// Location is preserved between calls to Iterate().
	Iterate(loc *Location) (Operation, error)

	// Needs specifies information about the objective function needed by the
	// optimizer beyond just the function value. The information is used
	// internally for initialization and must match evaluation types returned
	// by Init() and Iterate() during the optimization process.
	Needs() struct {
		Gradient bool
		Hessian  bool
	}
}

// Statuser can report the status and any error. It is intended for methods as
// an additional error reporting mechanism apart from the errors returned from
// Init() and Iterate().
type Statuser interface {
	Status() (Status, error)
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
	Init(value, derivative float64, step float64) Operation

	// Finished takes in the values of φ and φ' evaluated at the previous step,
	// and returns whether a sufficiently accurate minimum of φ has been found.
	Finished(value, derivative float64) bool

	// Iterate takes in the values of φ and φ' evaluated at the previous step
	// and returns the next step size and the type of evaluation to be
	// performed at x_k + step * dir_k.
	Iterate(value, derivative float64) (step float64, op Operation, err error)
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

// A Recorder can record the progress of the optimization, for example to print
// the progress to StdOut or to a log file. A Recorder must not modify any data.
type Recorder interface {
	Init() error
	Record(*Location, Operation, *Stats) error
}
