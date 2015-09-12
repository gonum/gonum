// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gonum/matrix/mat64"
)

const defaultGradientAbsTol = 1e-6

// EvaluationType is used by a Method to specify the objective-function
// information needed at an x location.
type EvaluationType uint

// Evaluation types can be composed together using the binary or operator, for
// example 'FuncEvaluation | GradEvaluation' to evaluate both the function
// value and the gradient.
const (
	NoEvaluation   EvaluationType = 0
	FuncEvaluation EvaluationType = 1 << (iota - 1)
	GradEvaluation
	HessEvaluation
)

func (e EvaluationType) String() string {
	return fmt.Sprintf("EvaluationType(Func: %t, Grad: %t, Hess: %t, Extra: 0b%b)",
		e&FuncEvaluation != 0,
		e&GradEvaluation != 0,
		e&HessEvaluation != 0,
		e&^(FuncEvaluation|GradEvaluation|HessEvaluation))
}

// IterationType specifies the type of iteration.
type IterationType int

const (
	NoIteration IterationType = iota
	MajorIteration
	MinorIteration
	SubIteration
	InitIteration
	PostIteration // Iteration after the optimization. Sent to Recorder.
)

func (i IterationType) String() string {
	if i < 0 || int(i) >= len(iterationStrings) {
		return fmt.Sprintf("IterationType(%d)", i)
	}
	return iterationStrings[i]
}

var iterationStrings = [...]string{
	"NoIteration",
	"MajorIteration",
	"MinorIteration",
	"SubIteration",
	"InitIteration",
	"PostIteration",
}

// Location represents a location in the optimization procedure.
type Location struct {
	X        []float64
	F        float64
	Gradient []float64
	Hessian  *mat64.SymDense
}

// Result represents the answer of an optimization run. It contains the optimum
// location as well as the Status at convergence and Statistics taken during the
// run.
type Result struct {
	Location
	Stats
	Status Status
}

// Stats contains the statistics of the run.
type Stats struct {
	MajorIterations int           // Total number of major iterations
	FuncEvaluations int           // Number of evaluations of Func()
	GradEvaluations int           // Number of evaluations of Grad()
	HessEvaluations int           // Number of evaluations of Hess()
	Runtime         time.Duration // Total runtime of the optimization
}

// complementEval returns an evaluation type that evaluates fields of loc not
// evaluated by eval.
func complementEval(loc *Location, eval EvaluationType) (complEval EvaluationType) {
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
	// Grad evaluates the gradient at x and stores the result in-place in grad.
	// Grad must not modify x.
	Grad func(x []float64, grad []float64)
	// Hess evaluates the Hessian at x and stores the result in-place in hess.
	// Hess must not modify x.
	Hess func(x []float64, hess *mat64.SymDense)
	// Status reports the status of the optimization problem and reports
	// any error.
	Status func() (Status, error)
}

// TODO(btracey): Think about making this an exported function when the
// constraint interface is designed.
func (p Problem) satisfies(method Method) error {
	if method.Needs().Gradient && p.Grad == nil {
		return errors.New("optimize: problem does not provide needed Grad function")
	}
	if method.Needs().Hessian && p.Hess == nil {
		return errors.New("optimize: problem does not provide needed Hess function")
	}
	return nil
}

// Settings represents settings of the optimization run. It contains initial
// settings, convergence information, and Recorder information. In general, users
// should use DefaultSettings() rather than constructing a Settings literal.
//
// If UseInitData is true, InitialValue, InitialGradient and InitialHessian
// specify function information at the initial location.
//
// If Recorder is nil, no information will be recorded.
type Settings struct {
	UseInitialData  bool            // Use supplied information about the conditions at the initial x.
	InitialValue    float64         // Func(x) at the initial x.
	InitialGradient []float64       // Grad(x) at the initial x.
	InitialHessian  *mat64.SymDense // Hess(x) at the initial x.

	// FunctionThreshold is the threshold for acceptably small values of the
	// objective function. FunctionThreshold status is returned if
	// the objective function is less than this value.
	// The default value is -inf.
	FunctionThreshold float64

	// GradientThreshold determines the accuracy to which the minimum is found.
	// GradientThreshold status is returned if the infinity norm of
	// the gradient is less than this value.
	// Has no effect if gradient information is not used.
	// The default value is 1e-6.
	GradientThreshold float64

	// FunctionConverge tests that the function value decreases by a significant
	// amount over the specified number of iterations. If
	//  f < f_best && f_best - f > Relative * maxabs(f, f_best) + Absolute
	// then a significant decrease has occured, and f_best is updated. If there is
	// no significant decrease for Iterations major iterations, FunctionConvergence
	// is returned. If this is nil or if Iterations == 0, it has no effect.
	FunctionConverge *FunctionConverge

	// MajorIterations is the maximum number of iterations allowed.
	// IterationLimit status is returned if the number of major iterations
	// equals or exceeds this value.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	MajorIterations int

	// Runtime is the maximum runtime allowed. RuntimeLimit status is returned
	// if the duration of the run is longer than this value. Runtime is only
	// checked at iterations of the Method.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	Runtime time.Duration

	// FuncEvaluations is the maximum allowed number of function evaluations.
	// FunctionEvaluationLimit status is returned if the total number of calls
	// to Func() equals or exceeds this number.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	FuncEvaluations int

	// GradEvaluations is the maximum allowed number of gradient evaluations.
	// GradientEvaluationLimit status is returned if the total number of calls
	// to Grad() equals or exceeds this number.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	GradEvaluations int

	// HessEvaluations is the maximum allowed number of Hessian evaluations.
	// HessianEvaluationLimit status is returned if the total number of calls
	// to Hess() equals or exceeds this number.
	// If it equals zero, this setting has no effect.
	// The default value is 0.
	HessEvaluations int

	Recorder Recorder
}

// DefaultSettings returns a new Settings struct containing the default settings.
func DefaultSettings() *Settings {
	return &Settings{
		GradientThreshold: defaultGradientAbsTol,
		FunctionThreshold: math.Inf(-1),
		FunctionConverge: &FunctionConverge{
			Absolute:   1e-10,
			Iterations: 20,
		},
	}
}

// resize takes x and returns a slice of length dim. It returns a resliced x
// if cap(x) >= dim, and a new slice otherwise.
func resize(x []float64, dim int) []float64 {
	if dim > cap(x) {
		return make([]float64, dim)
	}
	return x[:dim]
}

func resizeSymDense(m *mat64.SymDense, dim int) *mat64.SymDense {
	if m == nil || cap(m.RawSymmetric().Data) < dim*dim {
		return mat64.NewSymDense(dim, nil)
	}
	return mat64.NewSymDense(dim, m.RawSymmetric().Data[:dim*dim])
}
