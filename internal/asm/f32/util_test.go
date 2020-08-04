// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32_test

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

const (
	msgRes      = "%v: unexpected result Got: %v Expected: %v"
	msgVal      = "%v: unexpected value at %v Got: %v Expected: %v"
	msgGuard    = "%v: guard violated in %s vector %v %v"
	msgReadOnly = "%v: modified read-only %v argument"
)

var (
	nan = float32(math.NaN())
	inf = float32(math.Inf(1))
)

// sameApprox tests for nan-aware equality within tolerance.
func sameApprox(x, y, tol float32) bool {
	a, b := float64(x), float64(y)
	return same(x, y) || scalar.EqualWithinAbsOrRel(a, b, float64(tol), float64(tol))
}

func same(x, y float32) bool {
	return scalar.Same(float64(x), float64(y))
}

// sameStrided returns true if the strided vector x contains elements of the
// dense vector ref at indices i*inc, false otherwise.
func sameStrided(ref, x []float32, inc int) bool {
	if inc < 0 {
		inc = -inc
	}
	for i, v := range ref {
		if !same(x[i*inc], v) {
			return false
		}
	}
	return true
}

func guardVector(v []float32, g float32, gdLn int) (guarded []float32) {
	guarded = make([]float32, len(v)+gdLn*2)
	copy(guarded[gdLn:], v)
	for i := 0; i < gdLn; i++ {
		guarded[i] = g
		guarded[len(guarded)-1-i] = g
	}
	return guarded
}

func isValidGuard(v []float32, g float32, gdLn int) bool {
	for i := 0; i < gdLn; i++ {
		if !same(v[i], g) || !same(v[len(v)-1-i], g) {
			return false
		}
	}
	return true
}

func guardIncVector(vec []float32, gdVal float32, inc, gdLen int) (guarded []float32) {
	if inc < 0 {
		inc = -inc
	}
	inrLen := len(vec) * inc
	guarded = make([]float32, inrLen+gdLen*2)
	for i := range guarded {
		guarded[i] = gdVal
	}
	for i, v := range vec {
		guarded[gdLen+i*inc] = v
	}
	return guarded
}

func checkValidIncGuard(t *testing.T, v []float32, g float32, inc, gdLn int) {
	srcLn := len(v) - 2*gdLn
	for i := range v {
		switch {
		case same(v[i], g):
			// Correct value
		case i < gdLn:
			t.Error("Front guard violated at", i, v[:gdLn])
		case i > gdLn+srcLn:
			t.Error("Back guard violated at", i-gdLn-srcLn, v[gdLn+srcLn:])
		case (i-gdLn)%inc == 0 && (i-gdLn)/inc < len(v):
		default:
			t.Error("Internal guard violated at", i-gdLn, v[gdLn:gdLn+srcLn])
		}
	}
}

var ( // Offset sets for testing alignment handling in Unitary assembly functions.
	align2 = newIncSet(0, 1, 2, 3)
	align3 = newIncToSet(0, 1, 2, 3)
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
