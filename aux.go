// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import "math"

// ArmijioConditionMet returns true if the Armijio condition (aka sufficient decrease)
// has been met. Under normal conditions, the following should be true, though this is not enforced.:
//  initGrad < 0
//  step > 0
//  0 < funConst < 1
func ArmijioConditionMet(currObj, initObj, initGrad, step, funConst float64) bool {
	if currObj > initObj+funConst*step*initGrad {
		return false
	}
	return true
}

// StrongWolfeConditionsMet returns true if the strong Wolfe conditions have been met.
// The strong wolfe conditions ensure sufficient decrease in the function value,
// and sufficient decrease in the magnitude of the projected gradient. Under normal
// conditions, the following should be true, though this is not enforced:
//  initGrad < 0
//  step > 0
//  0 <= funConst < gradConst < 1
func StrongWolfeConditionsMet(currObj, currGrad, initObj, initGrad, step, funConst, gradConst float64) bool {
	if currObj > initObj+funConst*step*initGrad {
		return false
	}
	if math.Abs(currGrad) >= gradConst*math.Abs(initGrad) {
		return false
	}
	return true
}

// WeakWolfeConditionsMet returns true if the weak Wolfe conditions have been met.
// The weak wolfe conditions ensure sufficient decrease in the function value,
// and sufficient decrease in the value of the projected gradient. Under normal
// conditions, the following should be true, though this is not enforced:
//  initGrad < 0
//  step > 0
//  0 <= funConst < gradConst < 1
func WeakWolfeConditionsMet(currObj, currGrad, initObj, initGrad, step, funConst, gradConst float64) bool {
	if currObj > initObj+funConst*step*initGrad {
		return false
	}
	if currGrad < gradConst*initGrad {
		return false
	}
	return true
}

// resize takes x and returns a slice of length dim.
// It returns a resliced x if cap(x) >= dim, and a new
// slice otherwies
func resize(x []float64, dim int) []float64 {
	if cap(x) < dim {
		return make([]float64, dim)
	}
	return x[:dim]
}
