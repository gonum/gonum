// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"errors"
	"fmt"
)

// ErrInf signifies the initial function value is NaN.
var ErrInf = errors.New("optimize: initial function value is Inf")

// ErrLinesearchIterations signifies the linesearch has iterated too many times.
// This may occur if the gradient tolerance is set too low.
var ErrLinesearchIterations = errors.New("linesearch: too many iterations")

// ErrMismatch signifies that the optimization function did not implement the
// interfaces necessary for the supplied optimization method.
type ErrMismatch struct {
	Type EvaluationType
}

func (e ErrMismatch) Error() string {
	return fmt.Sprintf("optimizer wanted to use evaluation type %v, but the user supplied function does not implement it", e.Type)
}

// ErrNaN signifies the initial function value is NaN.
var ErrNaN = errors.New("optimize: initial function value is NaN")

// ErrNonNegativestepDirection signifies that the linesearch has received a step
// direction in which the gradient is not negative.
var ErrNonNegativeStepDirection = errors.New("linesearch: projected gradient not negative")

// ErrZeroDimensional signifies an optimization was called with an input of length 0.
var ErrZeroDimensional = errors.New("optimize: zero dimensional input")
