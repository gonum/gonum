// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"time"

	"gonum.org/v1/gonum/mat"
)

const defaultGradientAbsTol = 1e-12

// Operation represents the set of operations commanded by Method at each
// iteration. It is a bitmap of various Iteration and Evaluation constants.
// Individual constants must NOT be combined together by the binary OR operator
// except for the Evaluation operations.
type Operation uint64

// Supported Operations.
const (
	// NoOperation specifies that no evaluation or convergence check should
	// take place.
	NoOperation Operation = 0
	// InitIteration is sent to Recorder to indicate the initial location.
	// All fields of the location to record must be valid.
	// Method must not return it.
	InitIteration Operation = 1 << (iota - 1)
	// PostIteration is sent to Recorder to indicate the final location
	// reached during an optimization run.
	// All fields of the location to record must be valid.
	// Method must not return it.
	PostIteration
	// MajorIteration indicates that the next candidate location for
	// an optimum has been found and convergence should be checked.
	MajorIteration
	// MethodDone declares that the method is done running. A method must
	// be a Statuser in order to use this iteration, and after returning
	// MethodDone, the Status must return other than NotTerminated.
	MethodDone
	// FuncEvaluation specifies that the objective function
	// should be evaluated.
	FuncEvaluation
	// GradEvaluation specifies that the gradient
	// of the objective function should be evaluated.
	GradEvaluation
	// HessEvaluation specifies that the Hessian
	// of the objective function should be evaluated.
	HessEvaluation
	// signalDone is used internally to signal completion.
	signalDone

	// Mask for the evaluating operations.
	evalMask = FuncEvaluation | GradEvaluation | HessEvaluation
)

func (op Operation) isEvaluation() bool {
	return op&evalMask != 0 && op&^evalMask == 0
}

func (op Operation) String() string {
	if op&evalMask != 0 {
		return fmt.Sprintf("Evaluation(Func: %t, Grad: %t, Hess: %t, Extra: 0b%b)",
			op&FuncEvaluation != 0,
			op&GradEvaluation != 0,
			op&HessEvaluation != 0,
			op&^(evalMask))
	}
	s, ok := operationNames[op]
	if ok {
		return s
	}
	return fmt.Sprintf("Operation(%d)", op)
}

var operationNames = map[Operation]string{
	NoOperation:    "NoOperation",
	InitIteration:  "InitIteration",
	MajorIteration: "MajorIteration",
	PostIteration:  "PostIteration",
	MethodDone:     "MethodDone",
	signalDone:     "signalDone",
}

// Result represents the answer of an optimization run. It contains the optimum
// function value, X location, and gradient as well as the Status at convergence
// and Statistics taken during the run.
type Result struct {
	Location
	Stats
	Status Status
}

// Stats contains the statistics of the run.
type Stats struct {
	MajorIterations int           // Total number of major iterations
	FuncEvaluations int           // Number of evaluations of Func
	GradEvaluations int           // Number of evaluations of Grad
	HessEvaluations int           // Number of evaluations of Hess
	Runtime         time.Duration // Total runtime of the optimization
}

// complementEval returns an evaluating operation that evaluates fields of loc
// not evaluated by eval.
func complementEval(loc *Location, eval Operation) (complEval Operation) {
	if eval&FuncEvaluation == 0 {
		complEval = FuncEvaluation
	}
	if loc.Gradient != nil && eval&GradEvaluation == 0 {
		complEval |= GradEvaluation
	}
	if loc.Hessian != nil && eval&HessEvaluation == 0 {
		complEval |= HessEvaluation
	}
	return complEval
}

// Problem describes the optimization problem to be solved.
type Problem struct {
	// Func evaluates the objective function at the given location. Func
	// must not modify x.
	Func func(x []float64) float64

	// Grad evaluates the gradient at x and returns the result. Grad may use
	// (and return) the provided slice if it is non-nil, or must allocate a new
	// slice otherwise. Grad must not modify x.
	Grad func(grad []float64, x []float64) []float64

	// Hess evaluates the Hessian at x and stores the result in-place in hess.
	// Hess must not modify x. Hess may use (and return) the provided Symmetric
	// if it is non-nil, or must allocate a new Symmetric otherwise. Minimize
	// will 'give back' the returned Symmetric where possible, allowing Hess
	// to use a type assertion on the provided matrix.
	Hess func(hess mat.Symmetric, x []float64) mat.Symmetric

	// Status reports the status of the objective function being optimized and any
	// error. This can be used to terminate early, for example when the function is
	// not able to evaluate itself. The user can use one of the pre-provided Status
	// constants, or may call NewStatus to create a custom Status value.
	Status func() (Status, error)
}

// Available describes the functions available to call in Problem.
type Available struct {
	Grad bool
	Hess bool
}

func availFromProblem(prob Problem) Available {
	return Available{Grad: prob.Grad != nil, Hess: prob.Hess != nil}
}

// function tests if the Problem described by the receiver is suitable for an
// unconstrained Method that only calls the function, and returns the result.
func (has Available) function() (uses Available, err error) {
	// TODO(btracey): This needs to be modified when optimize supports
	// constrained optimization.
	return Available{}, nil
}

// gradient tests if the Problem described by the receiver is suitable for an
// unconstrained gradient-based Method, and returns the result.
func (has Available) gradient() (uses Available, err error) {
	// TODO(btracey): This needs to be modified when optimize supports
	// constrained optimization.
	if !has.Grad {
		return Available{}, ErrMissingGrad
	}
	return Available{Grad: true}, nil
}

// hessian tests if the Problem described by the receiver is suitable for an
// unconstrained Hessian-based Method, and returns the result.
func (has Available) hessian() (uses Available, err error) {
	// TODO(btracey): This needs to be modified when optimize supports
	// constrained optimization.
	if !has.Grad {
		return Available{}, ErrMissingGrad
	}
	if !has.Hess {
		return Available{}, ErrMissingHess
	}
	return Available{Grad: true, Hess: true}, nil
}

// Settings represents settings of the optimization run. It contains initial
// settings, convergence information, and Recorder information. Convergence
// settings are only checked at MajorIterations, while Evaluation thresholds
// are checked at every Operation. See the field comments for default values.
type Settings struct {
	// InitValues specifies properties (function value, gradient, etc.) known
	// at the initial location passed to Minimize. If InitValues is non-nil, then
	// the function value F must be provided, the location X must not be specified
	// and other fields may be specified. The values in Location may be modified
	// during the call to Minimize.
	InitValues *Location

	// GradientThreshold stops optimization with GradientThreshold status if the
	// infinity norm of the gradient is less than this value. This defaults to
	// a value of 0 (and so gradient convergence is not checked), however note
	// that many Methods (LBFGS, CG, etc.) will converge with a small value of
	// the gradient, and so to fully disable this setting the Method may need to
	// be modified.
	// This setting has no effect if the gradient is not used by the Method.
	GradientThreshold float64

	// Converger checks if the optimization has converged based on the (history
	// of) locations found during the optimization. Minimize will pass the
	// Location at every MajorIteration to the Converger.
	//
	// If the Converger is nil, a default value of
	//  FunctionConverge {
	//		Absolute: 1e-10,
	//		Iterations: 100,
	//  }
	// will be used. NeverTerminated can be used to always return a
	// NotTerminated status.
	Converger Converger

	// MajorIterations is the maximum number of iterations allowed.
	// IterationLimit status is returned if the number of major iterations
	// equals or exceeds this value.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	MajorIterations int

	// Runtime is the maximum runtime allowed. RuntimeLimit status is returned
	// if the duration of the run is longer than this value. Runtime is only
	// checked at MajorIterations of the Method.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	Runtime time.Duration

	// FuncEvaluations is the maximum allowed number of function evaluations.
	// FunctionEvaluationLimit status is returned if the total number of calls
	// to Func equals or exceeds this number.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	FuncEvaluations int

	// GradEvaluations is the maximum allowed number of gradient evaluations.
	// GradientEvaluationLimit status is returned if the total number of calls
	// to Grad equals or exceeds this number.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	GradEvaluations int

	// HessEvaluations is the maximum allowed number of Hessian evaluations.
	// HessianEvaluationLimit status is returned if the total number of calls
	// to Hess equals or exceeds this number.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	HessEvaluations int

	Recorder Recorder

	// Concurrent represents how many concurrent evaluations are possible.
	Concurrent int
}

// resize takes x and returns a slice of length dim. It returns a resliced x
// if cap(x) >= dim, and a new slice otherwise.
func resize(x []float64, dim int) []float64 {
	if dim > cap(x) {
		return make([]float64, dim)
	}
	return x[:dim]
}

func resizeSymDense(m *mat.SymDense, dim int) *mat.SymDense {
	if m == nil || cap(m.RawSymmetric().Data) < dim*dim {
		return mat.NewSymDense(dim, nil)
	}
	return mat.NewSymDense(dim, m.RawSymmetric().Data[:dim*dim])
}
