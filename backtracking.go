// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

const (
	defaultBacktrackingDecrease = 0.5
	defaultBacktrackingFunConst = 1e-4
	minimumBacktrackingStepSize = 1e-20
)

// Backtracking is a type that implements LinesearchMethod using a backtracking
// line search. A backtracking line search checks that the Armijo condition has
// been met with the given function constant. If the Armijo condition has not
// been met, the step size is decreased by a factor of Decrease.
//
// The Armijo conditions only require the gradient at the initial condition
// (not successive step locations), and so Backtracking may be a good linesearch
// method for functions with expensive gradients. Backtracking is not appropriate
// for optimizers that require the Wolfe conditions to be met, such as BFGS.
//
// Both FunConst and Decrease must be between zero and one, and Backtracking will
// panic otherwise. If either FunConst or Decrease are zero, it will be set to a
// reasonable default.
type Backtracking struct {
	FunConst float64 // Necessary function descrease for Armijo condition.
	Decrease float64 // Step size multiplier at each iteration (stepSize *= Decrease).

	stepSize float64
	initF    float64
	initG    float64
}

func (b *Backtracking) Init(loc LinesearchLocation, step float64, _ *FunctionInfo) EvaluationType {
	if step <= 0 {
		panic("backtracking: bad step size")
	}
	if loc.Derivative >= 0 {
		panic("Backtracking: init G non-negative")
	}

	if b.Decrease == 0 {
		b.Decrease = defaultBacktrackingDecrease
	}
	if b.FunConst == 0 {
		b.FunConst = defaultBacktrackingFunConst
	}
	if b.Decrease <= 0 || b.Decrease >= 1 {
		panic("backtracking: Decrease must be between 0 and 1")
	}
	if b.FunConst <= 0 || b.FunConst >= 1 {
		panic("backtracking: FunConst must be between 0 and 1")
	}

	b.stepSize = step
	b.initF = loc.F
	b.initG = loc.Derivative
	return FuncEvaluation
}

func (b *Backtracking) Finished(loc LinesearchLocation) bool {
	return ArmijoConditionMet(loc.F, b.initF, b.initG, b.stepSize, b.FunConst)
}

func (b *Backtracking) Iterate(_ LinesearchLocation) (float64, EvaluationType, error) {
	b.stepSize *= b.Decrease
	if b.stepSize < minimumBacktrackingStepSize {
		return 0, NoEvaluation, ErrLinesearchFailure
	}
	return b.stepSize, FuncEvaluation, nil
}
