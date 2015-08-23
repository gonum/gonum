// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "errors"

var (
	// ErrInf signifies the initial function value is Inf.
	ErrInf = errors.New("optimize: initial function value is Inf")

	// ErrGradInf signifies the initial function value is Inf.
	ErrGradInf = errors.New("optimize: initial gradient is Inf")

	// ErrLinesearchFailure signifies the linesearch has iterated too many
	// times. This may occur if the gradient tolerance is set too low.
	ErrLinesearchFailure = errors.New("linesearch: failed to converge")

	// ErrNaN signifies the initial function value is NaN.
	ErrNaN = errors.New("optimize: initial function value is NaN")

	// ErrGradNaN signifies the initial function value is NaN.
	ErrGradNaN = errors.New("optimize: initial gradient is NaN")

	// ErrNonNegativeStepDirection signifies that the linesearch has received a
	// step direction in which the gradient is not negative.
	ErrNonNegativeStepDirection = errors.New("linesearch: projected gradient not negative")

	// ErrZeroDimensional signifies an optimization was called with an input of length 0.
	ErrZeroDimensional = errors.New("optimize: zero dimensional input")

	// ErrNoProgress signifies that Linesearch cannot make further progress
	// because there is no change in location after LinesearchMethod step due
	// to floating-point arithmetic.
	ErrNoProgress = errors.New("linesearch: no change in location after linesearch step")
)
