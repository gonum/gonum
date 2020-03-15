// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmplxs

import (
	"math"
	"math/cmplx"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var (
	inf  = math.Inf(1)
	cinf = cmplx.Inf()
	nan  = math.NaN()
	cnan = cmplx.NaN()
)

// same tests for nan-aware equality.
func fsame(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

// sameApprox tests for nan-aware equality within tolerance.
func sameFloatApprox(a, b float64, tol float64) bool {
	return fsame(a, b) || floats.EqualWithinAbsOrRel(a, b, tol, tol)
}

func guardVector(vec []complex128, guardVal complex128, guardLen int) (guarded []complex128) {
	guarded = make([]complex128, len(vec)+guardLen*2)
	copy(guarded[guardLen:], vec)
	for i := 0; i < guardLen; i++ {
		guarded[i] = guardVal
		guarded[len(guarded)-1-i] = guardVal
	}
	return guarded
}

func isValidGuard(vec []complex128, guardVal complex128, guardLen int) bool {
	for i := 0; i < guardLen; i++ {
		if vec[i] != guardVal || vec[len(vec)-1-i] != guardVal {
			return false
		}
	}
	return true
}

func guardIncVector(vec []complex128, guardVal complex128, inc, guardLen int) (guarded []complex128) {
	sLn := len(vec) * inc
	if inc < 0 {
		sLn = len(vec) * -inc
	}
	guarded = make([]complex128, sLn+guardLen*2)
	for i, cas := 0, 0; i < len(guarded); i++ {
		switch {
		case i < guardLen, i > guardLen+sLn:
			guarded[i] = guardVal
		case (i-guardLen)%(inc) == 0 && cas < len(vec):
			guarded[i] = vec[cas]
			cas++
		default:
			guarded[i] = guardVal
		}
	}
	return guarded
}

func checkValidIncGuard(t *testing.T, vec []complex128, guardVal complex128, inc, guardLen int) {
	sLn := len(vec) - 2*guardLen
	if inc < 0 {
		sLn = len(vec) * -inc
	}

	for i := range vec {
		switch {
		case vec[i] == guardVal:
			// Correct value
		case i < guardLen:
			t.Errorf("Front guard violated at %d %v", i, vec[:guardLen])
		case i > guardLen+sLn:
			t.Errorf("Back guard violated at %d %v", i-guardLen-sLn, vec[guardLen+sLn:])
		case (i-guardLen)%inc == 0 && (i-guardLen)/inc < len(vec):
			// Ignore input values
		default:
			t.Errorf("Internal guard violated at %d %v", i-guardLen, vec[guardLen:guardLen+sLn])
		}
	}
}
