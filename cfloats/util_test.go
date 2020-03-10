// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cfloats_test

import (
	"math"
	"math/cmplx"
	"testing"

	. "gonum.org/v1/gonum/cfloats"
	"gonum.org/v1/gonum/floats"
)

const (
	msgVal   = "%v: unexpected value at %v Got: %v Expected: %v"
	msgGuard = "%v: Guard violated in %s vector %v %v"
)

var (
	inf       = math.Inf(1)
	cinf      = cmplx.Inf()
	nan       = math.NaN()
	cnan      = cmplx.NaN()
	benchSink complex128
)

func same(x, y complex128) bool {
	// return (x == y ||
	// 	math.IsNaN(real(x)) && math.IsNaN(real(y)) && imag(x) == imag(y) ||
	// 	math.IsNaN(imag(y)) && math.IsNaN(imag(x)) && real(y) == real(x) ||
	// 	math.IsNaN(real(x)) && math.IsNaN(real(y)) && math.IsNaN(imag(y)) && math.IsNaN(imag(x)))
	return x == y || (cmplx.IsNaN(x) && cmplx.IsNaN(y)) || (cmplx.IsInf(x) && cmplx.IsInf(y))
}

// same tests for nan-aware equality.
func fsame(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

// sameApprox tests for nan-aware equality within tolerance.
func sameApprox(a, b complex128, tol float64) bool {
	return same(a, b) || EqualWithinAbsOrRel(a, b, tol, tol)
}

// sameApprox tests for nan-aware equality within tolerance.
func sameFloatApprox(a, b float64, tol float64) bool {
	return fsame(a, b) || floats.EqualWithinAbsOrRel(a, b, tol, tol)
}

func guardVector(vec []complex128, guard_val complex128, guard_len int) (guarded []complex128) {
	guarded = make([]complex128, len(vec)+guard_len*2)
	copy(guarded[guard_len:], vec)
	for i := 0; i < guard_len; i++ {
		guarded[i] = guard_val
		guarded[len(guarded)-1-i] = guard_val
	}
	return guarded
}

func isValidGuard(vec []complex128, guard_val complex128, guard_len int) bool {
	for i := 0; i < guard_len; i++ {
		if vec[i] != guard_val || vec[len(vec)-1-i] != guard_val {
			return false
		}
	}
	return true
}

func guardIncVector(vec []complex128, guard_val complex128, inc, guard_len int) (guarded []complex128) {
	s_ln := len(vec) * inc
	if inc < 0 {
		s_ln = len(vec) * -inc
	}
	guarded = make([]complex128, s_ln+guard_len*2)
	for i, cas := 0, 0; i < len(guarded); i++ {
		switch {
		case i < guard_len, i > guard_len+s_ln:
			guarded[i] = guard_val
		case (i-guard_len)%(inc) == 0 && cas < len(vec):
			guarded[i] = vec[cas]
			cas++
		default:
			guarded[i] = guard_val
		}
	}
	return guarded
}

func checkValidIncGuard(t *testing.T, vec []complex128, guard_val complex128, inc, guard_len int) {
	s_ln := len(vec) - 2*guard_len
	if inc < 0 {
		s_ln = len(vec) * -inc
	}

	for i := range vec {
		switch {
		case vec[i] == guard_val:
			// Correct value
		case i < guard_len:
			t.Errorf("Front guard violated at %d %v", i, vec[:guard_len])
		case i > guard_len+s_ln:
			t.Errorf("Back guard violated at %d %v", i-guard_len-s_ln, vec[guard_len+s_ln:])
		case (i-guard_len)%inc == 0 && (i-guard_len)/inc < len(vec):
			// Ignore input values
		default:
			t.Errorf("Internal guard violated at %d %v", i-guard_len, vec[guard_len:guard_len+s_ln])
		}
	}
}

var ( // Offset sets for testing alignment handling in Unitary assembly functions.
	align1 = []int{0, 1}
	align2 = newIncSet(0, 1)
	align3 = newIncToSet(0, 1)
)

type incSet struct {
	x, y int
}

// genInc will generate all (x,y) combinations of the input increment set.
func newIncSet(inc ...int) []incSet {
	n := len(inc)
	is := make([]incSet, n*n)
	for x := range inc {
		for y := range inc {
			is[x*n+y] = incSet{inc[x], inc[y]}
		}
	}
	return is
}

type incToSet struct {
	dst, x, y int
}

// genIncTo will generate all (dst,x,y) combinations of the input increment set.
func newIncToSet(inc ...int) []incToSet {
	n := len(inc)
	is := make([]incToSet, n*n*n)
	for i, dst := range inc {
		for x := range inc {
			for y := range inc {
				is[i*n*n+x*n+y] = incToSet{dst, inc[x], inc[y]}
			}
		}
	}
	return is
}
