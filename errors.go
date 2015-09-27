// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import "errors"

var (
	// ErrInf signifies the initial function value is infinite.
	ErrInf = errors.New("optimize: initial function value is infinite")

	// ErrNaN signifies the initial function value is NaN.
	ErrNaN = errors.New("optimize: initial function value is NaN")

	// ErrGradInf signifies that a component of the initial gradient is infinite.
	ErrGradInf = errors.New("optimize: initial gradient is infinite")

	// ErrGradNaN signifies that a component of the initial gradient is NaN.
	ErrGradNaN = errors.New("optimize: initial gradient is NaN")

	// ErrZeroDimensional signifies an optimization was called with an input of length 0.
	ErrZeroDimensional = errors.New("optimize: zero dimensional input")

	// ErrLinesearcherFailure signifies that a Linesearcher has iterated too
	// many times. This may occur if the gradient tolerance is set too low.
	ErrLinesearcherFailure = errors.New("linesearch: failed to converge")

	// ErrNonDescentDirection signifies that LinesearchMethod has received a
	// search direction from a NextDirectioner in which the function is not
	// decreasing.
	ErrNonDescentDirection = errors.New("linesearch: non-descent search direction")

	// ErrNoProgress signifies that LinesearchMethod cannot make further
	// progress because there is no change in location after Linesearcher step
	// due to floating-point arithmetic.
	ErrNoProgress = errors.New("linesearch: no change in location after Linesearcher step")
)
