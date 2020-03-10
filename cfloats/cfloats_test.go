// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cfloats_test

import (
	"math"
	"math/cmplx"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"

	. "gonum.org/v1/gonum/cfloats"
	"gonum.org/v1/gonum/floats"
)

const (
	EqTolerance = 1e-14
	Small       = 10
	Medium      = 1000
	Large       = 100000
	Huge        = 10000000
)

func areSlicesEqual(t *testing.T, truth, comp []complex128, str string) {
	if !EqualApprox(comp, truth, EqTolerance) {
		t.Errorf(str+". Expected %v, returned %v", truth, comp)
	}
}

func areFloatSlicesEqual(t *testing.T, truth, comp []float64, str string) {
	if !floats.EqualApprox(comp, truth, EqTolerance) {
		t.Errorf(str+". Expected %v, returned %v", truth, comp)
	}
}

func Panics(fun func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	fun()
	return
}

func TestAbs(t *testing.T) {
	test := []struct {
		a     []complex128
		truth []float64
	}{
		{
			[]complex128{1, -2, 3},
			[]float64{1, 2, 3},
		},
		{
			[]complex128{4 + 4i, 5 - 5i, -6 + 6i},
			[]float64{5.656854249492381, 7.0710678118654755, 8.485281374238571},
		},
		{
			[]complex128{7i, -8i, 9i},
			[]float64{7, 8, 9},
		},
	}

	for _, v := range test {
		n := Abs(v.a)
		areFloatSlicesEqual(t, v.truth, n, "cfloats: Abs: incorrect values")
	}
}

func TestAdd(t *testing.T) {
	a := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	b := []complex128{4 + 4i, 5 + 5i, 6 + 6i}
	c := []complex128{7 + 7i, 8 + 8i, 9 + 9i}
	truth := []complex128{12 + 12i, 15 + 15i, 18 + 18i}
	n := make([]complex128, len(a))

	Add(n, a)
	Add(n, b)
	Add(n, c)
	areSlicesEqual(t, truth, n, "Wrong addition of slices new receiver")
	Add(a, b)
	Add(a, c)
	areSlicesEqual(t, truth, n, "Wrong addition of slices for no new receiver")

	// Test that it panics
	if !Panics(func() { Add(make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestAddConst(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	c := 6 - 6i
	truth := []complex128{9 - 3i, 10 - 2i, 7 - 5i, 13 + 1i, 11 - 1i}
	AddConst(c, s)
	areSlicesEqual(t, truth, s, "Wrong addition of constant")
}

func TestAddTo(t *testing.T) {
	a := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	b := []complex128{4 + 4i, 5 + 5i, 6 + 6i}
	truth := []complex128{5 + 5i, 7 + 7i, 9 + 9i}
	n1 := make([]complex128, len(a))

	n2 := AddTo(n1, a, b)
	areSlicesEqual(t, truth, n1, "Bad addition from mutator")
	areSlicesEqual(t, truth, n2, "Bad addition from returned slice")

	// Test that it panics
	if !Panics(func() { AddTo(make([]complex128, 2), make([]complex128, 3), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
	if !Panics(func() { AddTo(make([]complex128, 3), make([]complex128, 3), make([]complex128, 2)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestAddScaled(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	alpha := 6 + 6i
	dst := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i}
	ans := []complex128{1 + 37i, 2 + 50i, 3 + 15i, 4 + 88i, 5 + 65i}
	AddScaled(dst, alpha, s)
	if !EqualApprox(dst, ans, EqTolerance) {
		t.Errorf("Adding scaled did not match")
	}
	short := []complex128{1}
	if !Panics(func() { AddScaled(dst, alpha, short) }) {
		t.Errorf("Doesn't panic if s is smaller than dst")
	}
	if !Panics(func() { AddScaled(short, alpha, s) }) {
		t.Errorf("Doesn't panic if dst is smaller than s")
	}
}

func TestAddScaledTo(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	alpha := 6 + 6i
	y := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i}
	dst1 := make([]complex128, 5)
	ans := []complex128{1 + 37i, 2 + 50i, 3 + 15i, 4 + 88i, 5 + 65i}
	dst2 := AddScaledTo(dst1, y, alpha, s)
	if !EqualApprox(dst1, ans, EqTolerance) {
		t.Errorf("AddScaledTo did not match for mutator")
	}
	if !EqualApprox(dst2, ans, EqTolerance) {
		t.Errorf("AddScaledTo did not match for returned slice")
	}
	AddScaledTo(dst1, y, alpha, s)
	if !EqualApprox(dst1, ans, EqTolerance) {
		t.Errorf("Reusing dst did not match")
	}
	short := []complex128{1}
	if !Panics(func() { AddScaledTo(dst1, y, alpha, short) }) {
		t.Errorf("Doesn't panic if s is smaller than dst")
	}
	if !Panics(func() { AddScaledTo(short, y, alpha, s) }) {
		t.Errorf("Doesn't panic if dst is smaller than s")
	}
	if !Panics(func() { AddScaledTo(dst1, short, alpha, s) }) {
		t.Errorf("Doesn't panic if y is smaller than dst")
	}
}

func TestCumProd(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	receiver := make([]complex128, len(s))
	result := CumProd(receiver, s)
	truth := []complex128{3 + 3i, 24i, -24 + 24i, -336, -1680 - 1680i}
	areSlicesEqual(t, truth, receiver, "Wrong cumprod mutated with new receiver")
	areSlicesEqual(t, truth, result, "Wrong cumprod result with new receiver")
	CumProd(receiver, s)
	areSlicesEqual(t, truth, receiver, "Wrong cumprod returned with reused receiver")

	// Test that it panics
	if !Panics(func() { CumProd(make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}

	// Test empty CumProd
	emptyReceiver := make([]complex128, 0)
	truth = []complex128{}
	CumProd(emptyReceiver, emptyReceiver)
	areSlicesEqual(t, truth, emptyReceiver, "Wrong cumprod returned with empty receiver")
}

func TestCumSum(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	receiver := make([]complex128, len(s))
	result := CumSum(receiver, s)
	truth := []complex128{3 + 3i, 7 + 7i, 8 + 8i, 15 + 15i, 20 + 20i}
	areSlicesEqual(t, truth, receiver, "Wrong cumsum mutated with new receiver")
	areSlicesEqual(t, truth, result, "Wrong cumsum returned with new receiver")
	CumSum(receiver, s)
	areSlicesEqual(t, truth, receiver, "Wrong cumsum returned with reused receiver")

	// Test that it panics
	if !Panics(func() { CumSum(make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}

	// Test empty CumSum
	emptyReceiver := make([]complex128, 0)
	truth = []complex128{}
	CumSum(emptyReceiver, emptyReceiver)
	areSlicesEqual(t, truth, emptyReceiver, "Wrong cumsum returned with empty receiver")
}

func TestDistance(t *testing.T) {
	norms := []float64{1, 2, 4, math.Inf(1)}
	slices := []struct {
		s []complex128
		t []complex128
	}{
		{
			nil,
			nil,
		},
		{
			[]complex128{8 + 8i, 9 + 9i, 10 + 10i, -12 - 12i},
			[]complex128{8 + 8i, 9 + 9i, 10 + 10i, -12 - 12i},
		},
		{
			[]complex128{1 + 1i, 2 + 2i, 3 + 3i, -4 - 4i, -5 - 5i, 8 + 8i},
			[]complex128{-9.2 - 9.2i, -6.8 - 6.8i, 9 + 9i, -3 - 3i, -2 - 2i, 1 + 1i},
		},
	}

	for j, test := range slices {
		tmp := make([]complex128, len(test.s))
		for i, L := range norms {
			dist := Distance(test.s, test.t, L)
			copy(tmp, test.s)
			Sub(tmp, test.t)
			norm := Norm(tmp, L)
			if dist != norm { // Use equality because they should be identical.
				t.Errorf("Distance does not match norm for case %v, %v. Expected %v, Found %v.", i, j, norm, dist)
			}
		}
	}

	if !Panics(func() { Distance([]complex128{}, []complex128{1 + 1i}, 1) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestDiv(t *testing.T) {
	s1 := []complex128{5 + 5i, 12 + 12i, 27 + 27i}
	s2 := []complex128{1 + 1i, 2 - 2i, 3 + 3i}
	ans := []complex128{5, 6i, 9}
	Div(s1, s2)
	if !EqualApprox(s1, ans, EqTolerance) {
		t.Errorf("Mul doesn't give correct answer")
	}
	s1short := []complex128{1}
	if !Panics(func() { Div(s1short, s2) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
	s2short := []complex128{1}
	if !Panics(func() { Div(s1, s2short) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestDivTo(t *testing.T) {
	s1 := []complex128{5 + 5i, 12 + 12i, 27 + 27i}
	s1orig := []complex128{5 + 5i, 12 + 12i, 27 + 27i}
	s2 := []complex128{1 + 1i, 2 - 2i, 3 + 3i}
	s2orig := []complex128{1 + 1i, 2 - 2i, 3 + 3i}
	dst1 := make([]complex128, 3)
	ans := []complex128{5, 6i, 9}
	dst2 := DivTo(dst1, s1, s2)
	if !EqualApprox(dst1, ans, EqTolerance) {
		t.Errorf("DivTo doesn't give correct answer in mutated slice")
	}
	if !EqualApprox(dst2, ans, EqTolerance) {
		t.Errorf("DivTo doesn't give correct answer in returned slice")
	}
	if !EqualApprox(s1, s1orig, EqTolerance) {
		t.Errorf("S1 changes during multo")
	}
	if !EqualApprox(s2, s2orig, EqTolerance) {
		t.Errorf("s2 changes during multo")
	}
	DivTo(dst1, s1, s2)
	if !EqualApprox(dst1, ans, EqTolerance) {
		t.Errorf("DivTo doesn't give correct answer reusing dst")
	}
	dstShort := []complex128{1}
	if !Panics(func() { DivTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []complex128{1}
	if !Panics(func() { DivTo(dst1, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []complex128{1}
	if !Panics(func() { DivTo(dst1, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestDot(t *testing.T) {
	s1 := []complex128{1 - 1i, 2 + 2i, 3 - 3i, 4 + 4i}
	s2 := []complex128{-3 - 3i, 4 + 4i, 5 - 5i, -6 + 6i}
	truth := -54 - 14i
	ans := Dot(s1, s2, false)
	if ans != truth {
		t.Errorf("Dot product computed incorrectly")
	}
	truth = 46 + 42i
	ans = Dot(s1, s2, true)
	if ans != truth {
		t.Errorf("Dot product computed incorrectly")
	}

	// Test that it panics
	if !Panics(func() { Dot(make([]complex128, 2), make([]complex128, 3), false) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestEquals(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	if !Equal(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i + 1e-14}
	if Equal(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
	if Equal(s1, []complex128{}) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualApprox(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i + 1e-10}
	if EqualApprox(s1, s2, 1e-13) {
		t.Errorf("Unequal slices returned as equal for absolute")
	}
	if !EqualApprox(s1, s2, 1e-5) {
		t.Errorf("Equal slices returned as unequal for absolute")
	}
	s1 = []complex128{1 + 1i, 2 + 2i, 3 + 3i, 1000 + 1000i}
	s2 = []complex128{1 + 1i, 2 + 2i, 3 + 3i, (1000 + 1000i) * (1 + 1e-7)}
	if EqualApprox(s1, s2, 1e-8) {
		t.Errorf("Unequal slices returned as equal for relative")
	}
	if !EqualApprox(s1, s2, 1e-5) {
		t.Errorf("Equal slices returned as unequal for relative")
	}
	if EqualApprox(s1, []complex128{}, 1e-5) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualFunc(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	eq := func(x, y complex128) bool { return x == y }
	if !EqualFunc(s1, s2, eq) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i + 1e-14}
	if EqualFunc(s1, s2, eq) {
		t.Errorf("Unequal slices returned as equal")
	}
	if EqualFunc(s1, []complex128{}, eq) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualsRelative(t *testing.T) {
	equalityTests := []struct {
		a, b  complex128
		tol   float64
		equal bool
	}{
		{1000000 + 1000000i, 1000001 + 1000001i, 0, true},
		{1000001 + 1000001i, 1000000 + 1000000i, 0, true},
		{10000 + 10000i, 10001 + 10001i, 0, false},
		{10001 + 10001i, 10000 + 10000i, 0, false},
		{-1000000 - 1000000i, -1000001 - 1000001i, 0, true},
		{-1000001 + 1000001i, -1000000 + 1000000i, 0, true},
		{-10000 + 10000i, -10001 + 10001i, 0, false},
		{-10001 + 10001i, -10000 + 10000i, 0, false},
		{1.0000001 + 1.0000001i, 1.0000002 + 1.0000002i, 0, true},
		{1.0000002 + 1.0000002i, 1.0000001 + 1.0000001i, 0, true},
		{1.0002 + 1.0002i, 1.0001 + 1.0001i, 0, false},
		{1.0001 + 1.0001i, 1.0002 + 1.0002i, 0, false},
		{-1.000001 - 1.000001i, -1.000002 - 1.000002i, 0, true},
		{-1.000002 - 1.000002i, -1.000001 - 1.000001i, 0, true},
		{-1.0001 - 1.0001i, -1.0002 - 1.0002i, 0, false},
		{-1.0002 - 1.0002i, -1.0001 - 1.0001i, 0, false},
		{0.000000001000001 + 0.000000001000001i, 0.000000001000002 + 0.000000001000002i, 0, true},
		{0.000000001000002 + 0.000000001000002i, 0.000000001000001 + 0.000000001000001i, 0, true},
		{0.000000000001002 + 0.000000000001002i, 0.000000000001001 + 0.000000000001001i, 0, false},
		{0.000000000001001 + 0.000000000001001i, 0.000000000001002 + 0.000000000001002i, 0, false},
		{-0.000000001000001 - 0.000000001000001i, -0.000000001000002 - 0.000000001000002i, 0, true},
		{-0.000000001000002 - 0.000000001000002i, -0.000000001000001 - 0.000000001000001i, 0, true},
		{-0.000000000001002 - 0.000000000001002i, -0.000000000001001 - 0.000000000001001i, 0, false},
		{-0.000000000001001 - 0.000000000001001i, -0.000000000001002 - 0.000000000001002i, 0, false},
		{0, 0, 0, true},
		{0, -0, 0, true},
		{-0, -0, 0, true},
		{0.00000001 + 0.00000001i, 0, 0, false},
		{0, 0.00000001 + 0.00000001i, 0, false},
		{-0.00000001 - 0.00000001i, 0, 0, false},
		{0, -0.00000001 - 0.00000001i, 0, false},
		{0, 1e-310 + 1e-310i, 0.01, true},
		{1e-310 + 1e-310i, 0, 0.01, true},
		{1e-310 + 1e-310i, 0, 0.000001, false},
		{0, 1e-310 + 1e-310i, 0.000001, false},
		{0, -1e-310 - 1e-310i, 0.1, true},
		{-1e-310 - 1e-310i, 0, 0.1, true},
		{-1e-310 - 1e-310i, 0, 0.00000001, false},
		{0, -1e-310 - 1e-310i, 0.00000001, false},
		{cmplx.Inf(), cmplx.Inf(), 0, true},
		{-cmplx.Inf(), -cmplx.Inf(), 0, true},
		{-cmplx.Inf(), cmplx.Inf(), 0, false},
		{cmplx.Inf(), math.MaxFloat64, 0, false},
		{-cmplx.Inf(), -math.MaxFloat64, 0, false},
		{cmplx.NaN(), cmplx.NaN(), 0, false},
		{cmplx.NaN(), 0, 0, false},
		{-0, cmplx.NaN(), 0, false},
		{cmplx.NaN(), -0, 0, false},
		{0, cmplx.NaN(), 0, false},
		{cmplx.NaN(), cmplx.Inf(), 0, false},
		{cmplx.Inf(), cmplx.NaN(), 0, false},
		{cmplx.NaN(), -cmplx.Inf(), 0, false},
		{-cmplx.Inf(), cmplx.NaN(), 0, false},
		{complex(math.Inf(1), 0), cmplx.Inf(), 0, false},
		{cmplx.Inf(), complex(math.Inf(1), 0), 0, false},
		{complex(0, math.Inf(1)), cmplx.Inf(), 0, false},
		{cmplx.Inf(), complex(0, math.Inf(1)), 0, false},
		{complex(math.Inf(1), math.Inf(1)), cmplx.Inf(), 0, true},
		{cmplx.Inf(), complex(math.Inf(1), math.Inf(1)), 0, true},
		{cmplx.NaN(), complex(math.MaxFloat64, math.MaxFloat64), 0, false},
		{complex(math.MaxFloat64, math.MaxFloat64), cmplx.NaN(), 0, false},
		{cmplx.NaN(), -complex(math.MaxFloat64, math.MaxFloat64), 0, false},
		{-complex(math.MaxFloat64, math.MaxFloat64), cmplx.NaN(), 0, false},
		{cmplx.NaN(), complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, false},
		{complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), cmplx.NaN(), 0, false},
		{cmplx.NaN(), -complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, false},
		{-complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), cmplx.NaN(), 0, false},
		{1.000000001 + 1.000000001i, -1.0 - 1.0i, 0, false},
		{-1.0 - 1.0i, 1.000000001 + 1.000000001i, 0, false},
		{-1.000000001 - 1.000000001i, 1.0 + 1.0i, 0, false},
		{1.0 + 1.0i, -1.000000001 - 1.000000001i, 0, false},
		{10 * complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 10 * -complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, true},
		{1e11 * complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 1e11 * -complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, false},
		{complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), -complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, true},
		{-complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, true},
		{complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, 0, true},
		{0, complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, true},
		{-complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, 0, true},
		{0, -complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, true},
		{0.000000001 + 0.000000001i, -complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, false},
		{0.000000001 + 0.000000001i, complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0, false},
		{complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0.000000001 + 0.000000001i, 0, false},
		{-complex(math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64), 0.000000001 + 0.000000001i, 0, false},
	}
	for _, ts := range equalityTests {
		if ts.tol == 0 {
			ts.tol = 1e-5
		}
		if equal := EqualWithinRel(ts.a, ts.b, ts.tol); equal != ts.equal {
			t.Errorf("Relative equality of %g and %g with tolerance %g returned: %v. Expected: %v",
				ts.a, ts.b, ts.tol, equal, ts.equal)
		}
	}
}

func nextAfterN(x, y complex128, n int) complex128 {
	for i := 0; i < n; i++ {
		x = complex(math.Nextafter(real(x), real(y)), math.Nextafter(imag(x), imag(y)))
	}
	return x
}

func TestEqualsULP(t *testing.T) {
	if f := 67329.242 + 67329.242i; !EqualWithinULP(f, nextAfterN(f, cmplx.Inf(), 10), 10) {
		t.Errorf("Equal values returned as unequal")
	}
	if f := 67329.242 + 67329.242i; EqualWithinULP(f, nextAfterN(f, cmplx.Inf(), 5), 1) {
		t.Errorf("Unequal values returned as equal")
	}
	if f := 67329.242 + 67329.242i; EqualWithinULP(nextAfterN(f, cmplx.Inf(), 5), f, 1) {
		t.Errorf("Unequal values returned as equal")
	}
	if f := nextAfterN(0+0i, cmplx.Inf(), 2); !EqualWithinULP(f, nextAfterN(f, -cmplx.Inf(), 5), 10) {
		t.Errorf("Equal values returned as unequal")
	}
	if !EqualWithinULP(67329.242+67329.242i, 67329.242+67329.242i, 10) {
		t.Errorf("Equal complex128s not returned as equal")
	}
	if EqualWithinULP(1+1i, cmplx.NaN(), 10) {
		t.Errorf("NaN returned as equal")
	}
}

func TestEqualLengths(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s3 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	if !EqualLengths(s1, s2) {
		t.Errorf("Equal lengths returned as unequal")
	}
	if EqualLengths(s1, s3) {
		t.Errorf("Unequal lengths returned as equal")
	}
	if !EqualLengths(s1) {
		t.Errorf("Single slice returned as unequal")
	}
	if !EqualLengths() {
		t.Errorf("No slices returned as unequal")
	}
}

func eqIntSlice(one, two []int) string {
	if len(one) != len(two) {
		return "Length mismatch"
	}
	for i, val := range one {
		if val != two[i] {
			return "Index " + strconv.Itoa(i) + " mismatch"
		}
	}
	return ""
}

func TestFind(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	f := func(v complex128) bool { return real(v) > 3.5 && imag(v) < 7 }
	allTrueInds := []int{1, 4}

	// Test finding first two elements
	inds, err := Find(nil, f, s, 2)
	if err != nil {
		t.Errorf("Find first two: Improper error return")
	}
	trueInds := allTrueInds[:2]
	str := eqIntSlice(inds, trueInds)
	if str != "" {
		t.Errorf("Find first two: " + str)
	}

	// Test finding no elements with non nil slice
	inds = []int{1, 2, 3, 4, 5, 6}
	inds, err = Find(inds, f, s, 0)
	if err != nil {
		t.Errorf("Find no elements: Improper error return")
	}
	str = eqIntSlice(inds, []int{})
	if str != "" {
		t.Errorf("Find no non-nil: " + str)
	}

	// Test finding first two elements with non nil slice
	inds = []int{1, 2, 3, 4, 5, 6}
	inds, err = Find(inds, f, s, 2)
	if err != nil {
		t.Errorf("Find first two non-nil: Improper error return")
	}
	str = eqIntSlice(inds, trueInds)
	if str != "" {
		t.Errorf("Find first two non-nil: " + str)
	}

	// Test finding too many elements
	inds, err = Find(inds, f, s, 4)
	if err == nil {
		t.Errorf("Request too many: No error returned")
	}
	str = eqIntSlice(inds, allTrueInds)
	if str != "" {
		t.Errorf("Request too many: Does not match all of the inds: " + str)
	}

	// Test finding all elements
	inds, err = Find(nil, f, s, -1)
	if err != nil {
		t.Errorf("Find all: Improper error returned")
	}
	str = eqIntSlice(inds, allTrueInds)
	if str != "" {
		t.Errorf("Find all: Does not match all of the inds: " + str)
	}
}

func TestHasNaN(t *testing.T) {
	for i, test := range []struct {
		s   []complex128
		ans bool
	}{
		{},
		{
			s: []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
		},
		{
			s:   []complex128{1 + 1i, cmplx.NaN(), 3 + 3i, 4 + 4i},
			ans: true,
		},
		{
			s:   []complex128{1 + 1i, 2 + 2i, 3 + 3i, cmplx.NaN()},
			ans: true,
		},
	} {
		b := HasNaN(test.s)
		if b != test.ans {
			t.Errorf("HasNaN mismatch case %d. Expected %v, Found %v", i, test.ans, b)
		}
	}
}

func TestL1Dist(t *testing.T) {
	var t_gd, s_gd complex128 = -cinf, cinf
	for j, v := range []struct {
		s, t   []complex128
		expect float64
	}{
		{
			s:      []complex128{1 + 1i},
			t:      []complex128{1 + 1i},
			expect: 0,
		},
		{
			s:      []complex128{cnan},
			t:      []complex128{cnan},
			expect: nan,
		},
		{
			s:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i},
			t:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			expect: 0,
		},
		{
			s:      []complex128{1 + 2i, 4 + 2i, 3 + 6i},
			t:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			expect: 6,
		},
		{
			s:      []complex128{0, 0, 0},
			t:      []complex128{1i, 2, 3i},
			expect: 6,
		},
		{
			s:      []complex128{0, -4, -10},
			t:      []complex128{1, 2, 3},
			expect: 20,
		},
		{
			s:      []complex128{0, 1, 0, 1, 0},
			t:      []complex128{1, 1, cinf, 1, 1},
			expect: inf,
		},
		{
			s:      []complex128{cinf, 4, cnan, -cinf, 9},
			t:      []complex128{cinf, 4, cnan, -cinf, 3},
			expect: nan,
		},
	} {
		sg_ln, tg_ln := 4+j%2, 4+j%3
		v.s, v.t = guardVector(v.s, s_gd, sg_ln), guardVector(v.t, t_gd, tg_ln)
		s_lc, t_lc := v.s[sg_ln:len(v.s)-sg_ln], v.t[tg_ln:len(v.t)-tg_ln]
		ret := L1Dist(s_lc, t_lc)
		if !fsame(ret, v.expect) {
			t.Errorf("Test %d L1Dist error Got: %f Expected: %f", j, ret, v.expect)
		}
		if !isValidGuard(v.s, s_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in s vector %v %v", j, v.s[:sg_ln], v.s[len(v.s)-sg_ln:])
		}
		if !isValidGuard(v.t, t_gd, tg_ln) {
			t.Errorf("Test %d Guard violated in t vector %v %v", j, v.t[:tg_ln], v.t[len(v.t)-tg_ln:])
		}
	}
}

func TestL1Norm(t *testing.T) {
	var src_gd complex128 = 1 + 1i
	for j, v := range []struct {
		want float64
		x    []complex128
	}{
		{want: 0, x: []complex128{}},
		{want: 2.82842712474619009760, x: []complex128{2 + 2i}},
		{want: 8.48528137423857131694, x: []complex128{1 + 1i, 2 + 2i, 3 + 3i}},
		{want: 8.48528137423857131694, x: []complex128{-1 - 1i, -2 - 2i, -3 - 3i}},
		{want: nan, x: []complex128{cnan}},
		{want: 56.56854249492380195206, x: []complex128{8 + 8i, -8 - 8i, 8 + 8i, -8 - 8i, 8 + 8i}},
		{want: 7.07106781186547524400, x: []complex128{0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i}},
	} {
		g_ln := 4 + j%2
		v.x = guardVector(v.x, src_gd, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L1Norm(src)
		if !fsame(ret, v.want) {
			t.Errorf("Test %d L1Norm error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}

func TestL1NormInc(t *testing.T) {
	var src_gd complex128 = 1
	for j, v := range []struct {
		inc  int
		want float64
		x    []complex128
	}{
		{inc: 2, want: 0, x: []complex128{}},
		{inc: 3, want: 2.82842712474619009760, x: []complex128{2 + 2i}},
		{inc: 10, want: 8.48528137423857131694, x: []complex128{1 + 1i, 2 + 2i, 3 + 3i}},
		{inc: 5, want: 8.48528137423857131694, x: []complex128{-1 - 1i, -2 - 2i, -3 - 3i}},
		{inc: 3, want: math.NaN(), x: []complex128{cmplx.NaN()}},
		{inc: 15, want: 56.56854249492380195206, x: []complex128{8 + 8i, -8 - 8i, 8 + 8i, -8 - 8i, 8 + 8i}},
		{inc: 1, want: 7.07106781186547524400, x: []complex128{0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i}},
	} {
		g_ln, ln := 4+j%2, len(v.x)
		v.x = guardIncVector(v.x, src_gd, v.inc, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L1NormInc(src, ln, v.inc)
		if !fsame(ret, v.want) {
			t.Errorf("Test %d L1NormInc error Got: %f Expected: %f", j, ret, v.want)
		}
		checkValidIncGuard(t, v.x, src_gd, v.inc, g_ln)
	}
}

func TestL2NormUnitary(t *testing.T) {
	const tol = 1e-15

	var src_gd complex128 = 1
	for j, v := range []struct {
		want float64
		x    []complex128
	}{
		{want: 0, x: []complex128{}},
		{want: 2.8284271247461900976021, x: []complex128{2 + 2i}},
		{want: 5.2915026221291811810038, x: []complex128{1 + 1i, 2 + 2i, 3 + 3i}},
		{want: 5.2915026221291811810038, x: []complex128{-1 - 1i, -2 - 2i, -3 - 3i}},
		{want: nan, x: []complex128{cnan}},
		{want: nan, x: []complex128{1 + 1i, cinf, 3 + 3i, complex(floats.NaNWith(25), floats.NaNWith(25)), 5 + 5i}},
		{want: 25.298221281347034655984, x: []complex128{8 + 8i, -8 - 8i, 8 + 8i, -8 - 8i, 8 + 8i}},
		{want: 3.162277660168379331998, x: []complex128{0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i, 0, -1 - 1i, 0, 1 + 1i}},
	} {
		g_ln := 4 + j%2
		v.x = guardVector(v.x, src_gd, g_ln)
		src := v.x[g_ln : len(v.x)-g_ln]
		ret := L2NormUnitary(src)
		if !sameFloatApprox(ret, v.want, tol) {
			t.Errorf("Test %d L2Norm error Got: %f Expected: %f", j, ret, v.want)
		}
		if !isValidGuard(v.x, src_gd, g_ln) {
			t.Errorf("Test %d Guard violated in src vector %v %v", j, v.x[:g_ln], v.x[len(v.x)-g_ln:])
		}
	}
}

func TestLinfDist(t *testing.T) {
	var t_gd, s_gd complex128 = 0, cinf
	for j, v := range []struct {
		s, t   []complex128
		expect float64
	}{
		{
			s:      []complex128{1 + 1i},
			t:      []complex128{1 + 1i},
			expect: 0,
		},
		{
			s:      []complex128{cnan},
			t:      []complex128{cnan},
			expect: nan,
		},
		{
			s:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i},
			t:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			expect: 0,
		},
		{
			s:      []complex128{1 + 2i, 4 + 2i, 3 + 6i},
			t:      []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			expect: 3,
		},
		{
			s:      []complex128{0, 0, 0},
			t:      []complex128{1i, 2, 3i},
			expect: 3,
		},
		{
			s:      []complex128{0, 1, 0, 1, 0},
			t:      []complex128{1, 1, cinf, 1, 1},
			expect: inf,
		},
		{
			s:      []complex128{cinf, 4 + 4i, cnan, -cinf, 9i},
			t:      []complex128{cinf, 4 + 4i, cnan, -cinf, 3i},
			expect: 6,
		},
	} {
		sg_ln, tg_ln := 4+j%2, 4+j%3
		v.s, v.t = guardVector(v.s, s_gd, sg_ln), guardVector(v.t, t_gd, tg_ln)
		s_lc, t_lc := v.s[sg_ln:len(v.s)-sg_ln], v.t[tg_ln:len(v.t)-tg_ln]
		ret := LinfDist(s_lc, t_lc)
		if !fsame(ret, v.expect) {
			t.Errorf("Test %d LcinfDist error Got: %f Expected: %f", j, ret, v.expect)
		}
		if !isValidGuard(v.s, s_gd, sg_ln) {
			t.Errorf("Test %d Guard violated in s vector %v %v", j, v.s[:sg_ln], v.s[len(v.s)-sg_ln:])
		}
		if !isValidGuard(v.t, t_gd, tg_ln) {
			t.Errorf("Test %d Guard violated in t vector %v %v", j, v.t[:tg_ln], v.t[len(v.t)-tg_ln:])
		}
	}
}

func TestMul(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 - 3i}
	ans := []complex128{2i, 8i, 18}
	Mul(s1, s2)
	if !EqualApprox(s1, ans, EqTolerance) {
		t.Errorf("Mul doesn't give correct answer")
	}
	s1short := []complex128{1}
	if !Panics(func() { Mul(s1short, s2) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
	s2short := []complex128{1}
	if !Panics(func() { Mul(s1, s2short) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestMulTo(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	s1orig := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 - 3i}
	s2orig := []complex128{1 + 1i, 2 + 2i, 3 - 3i}
	dst1 := make([]complex128, 3)
	ans := []complex128{2i, 8i, 18}
	dst2 := MulTo(dst1, s1, s2)
	if !EqualApprox(dst1, ans, EqTolerance) {
		t.Errorf("MulTo doesn't give correct answer in mutated slice")
	}
	if !EqualApprox(dst2, ans, EqTolerance) {
		t.Errorf("MulTo doesn't give correct answer in returned slice")
	}
	if !EqualApprox(s1, s1orig, EqTolerance) {
		t.Errorf("S1 changes during multo")
	}
	if !EqualApprox(s2, s2orig, EqTolerance) {
		t.Errorf("s2 changes during multo")
	}
	MulTo(dst1, s1, s2)
	if !EqualApprox(dst1, ans, EqTolerance) {
		t.Errorf("MulTo doesn't give correct answer reusing dst")
	}
	dstShort := []complex128{1}
	if !Panics(func() { MulTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []complex128{1}
	if !Panics(func() { MulTo(dst1, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []complex128{1}
	if !Panics(func() { MulTo(dst1, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestNorm(t *testing.T) {
	s := []complex128{-1 - 1i, -3.4 + 3.4i, 5 + 5i, -6 - 6i}
	val := Norm(s, math.Inf(1))
	truth := 8.485281374238570292810132345258188471418031252261688439060078427944394870772642233102325205965849436353
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28-1%29%5E2+%2B++%28-3.4%29%5E2+%2B+5%5E2%2B++6%5E2%29%5E%281%2F2%29
	val = Norm(s, 2)
	truth = 12.12930336004504418595378095272330054901971147600040528843911630812107638517655523756683456134739177167
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28%7C-1%7C%29%5E3+%2B++%28%7C-3.4%7C%29%5E3+%2B+%7C5%7C%5E3%2B++%7C6%7C%5E3%29%5E%281%2F3%29
	val = Norm(s, 3)
	truth = 10.25507367732196054450277323048144438247898652676749379210000865661923721606947258806948751041799248425
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}

	//http://www.wolframalpha.com/input/?i=%7C-1%7C+%2B+%7C-3.4%7C+%2B+%7C5%7C%2B++%7C6%7C
	val = Norm(s, 1)
	truth = 21.77888886054566362593860965812734565850235712371992101304700246184694836827909127962633573768355303665
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
}

func TestProd(t *testing.T) {
	s := []complex128{}
	val := Prod(s)
	if val != 1 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	val = Prod(s)
	if val != (-1680 - 1680i) {
		t.Errorf("Wrong prod returned. Expected %v returned %v", 420, val)
	}
}

func TestReverse(t *testing.T) {
	for _, s := range [][]complex128{
		{0},
		{1 + 1i, 0},
		{2 + 2i, 1 + 1i, 0},
		{3 + 3i, 2 + 2i, 1 + 1i, 0},
		{9 + 9i, 8 + 8i, 7 + 7i, 6 + 6i, 5 + 5i, 4 + 4i, 3 + 3i, 2 + 2i, 1 + 1i, 0},
	} {
		Reverse(s)
		for i, v := range s {
			if v != complex(float64(i), float64(i)) {
				t.Errorf("unexpected values for element %d: got:%v want:%v", i, v, i)
			}
		}
	}
}

func TestRound(t *testing.T) {
	for _, test := range []struct {
		x    complex128
		prec int
		want complex128
	}{
		{x: 0, prec: 1, want: 0},
		{x: cmplx.Inf(), prec: 1, want: cmplx.Inf()},
		{x: cmplx.NaN(), prec: 1, want: cmplx.NaN()},
		{x: func() complex128 { var f complex128; return -f }(), prec: 1, want: 0},
		{x: complex(math.MaxFloat64, math.MaxFloat64) / 2, prec: 1, want: complex(math.MaxFloat64, math.MaxFloat64) / 2},
		{x: complex(1<<64, 1<<64), prec: 1, want: complex(1<<64, 1<<64)},
		{x: 454.4445 + 454.4445i, prec: 3, want: 454.445 + 454.445i},
		{x: 454.44445 + 454.44445i, prec: 4, want: 454.4445 + 454.4445i},
		{x: 0.42499 + 0.42499i, prec: 4, want: 0.425 + 0.425i},
		{x: 0.42599 + 0.42599i, prec: 4, want: 0.426 + 0.426i},
		{x: 0.424999999999993 + 0.424999999999993i, prec: 2, want: 0.42 + 0.42i},
		{x: 0.425 + 0.425i, prec: 2, want: 0.43 + 0.43i},
		{x: 0.425000000000001 + 0.425000000000001i, prec: 2, want: 0.43 + 0.43i},
		{x: 123.4244999999999 + 123.4244999999999i, prec: 3, want: 123.424 + 123.424i},
		{x: 123.4245 + 123.4245i, prec: 3, want: 123.425 + 123.425i},
		{x: 123.4245000000001 + 123.4245000000001i, prec: 3, want: 123.425 + 123.425i},

		{x: 454.45 + 454.45i, prec: 0, want: 454 + 454i},
		{x: 454.45 + 454.45i, prec: 1, want: 454.5 + 454.5i},
		{x: 454.45 + 454.45i, prec: 2, want: 454.45 + 454.45i},
		{x: 454.45 + 454.45i, prec: 3, want: 454.45 + 454.45i},
		{x: 454.445 + 454.445i, prec: 0, want: 454 + 454i},
		{x: 454.445 + 454.445i, prec: 1, want: 454.4 + 454.4i},
		{x: 454.445 + 454.445i, prec: 2, want: 454.45 + 454.45i},
		{x: 454.445 + 454.445i, prec: 3, want: 454.445 + 454.445i},
		{x: 454.445 + 454.445i, prec: 4, want: 454.445 + 454.445i},
		{x: 454.55 + 454.55i, prec: 0, want: 455 + 455i},
		{x: 454.55 + 454.55i, prec: 1, want: 454.6 + 454.6i},
		{x: 454.55 + 454.55i, prec: 2, want: 454.55 + 454.55i},
		{x: 454.55 + 454.55i, prec: 3, want: 454.55 + 454.55i},
		{x: 454.455 + 454.455i, prec: 0, want: 454 + 454i},
		{x: 454.455 + 454.455i, prec: 1, want: 454.5 + 454.5i},
		{x: 454.455 + 454.455i, prec: 2, want: 454.46 + 454.46i},
		{x: 454.455 + 454.455i, prec: 3, want: 454.455 + 454.455i},
		{x: 454.455 + 454.455i, prec: 4, want: 454.455 + 454.455i},

		// Negative precision.
		{x: 454.45 + 454.45i, prec: -1, want: 450 + 450i},
		{x: 454.45 + 454.45i, prec: -2, want: 500 + 500i},
		{x: 500 + 500i, prec: -3, want: 1000 + 1000i},
		{x: 500 + 500i, prec: -4, want: 0},
		{x: 1500 + 1500i, prec: -3, want: 2000 + 2000i},
		{x: 1500 + 1500i, prec: -4, want: 0},
	} {
		for _, sign := range []complex128{1, -1} {
			got := Round(sign*test.x, test.prec)
			want := sign * test.want
			if want == 0 {
				want = 0
			}
			if (got != want || math.Signbit(real(got)) != math.Signbit(real(want)) || math.Signbit(imag(got)) != math.Signbit(imag(want))) && !(cmplx.IsNaN(got) && cmplx.IsNaN(want)) {
				t.Errorf("unexpected result for Round(%g, %d): got: %g, want: %g", sign*test.x, test.prec, got, want)
			}
		}
	}
}

func TestRoundEven(t *testing.T) {
	for _, test := range []struct {
		x    complex128
		prec int
		want complex128
	}{
		{x: 0, prec: 1, want: 0},
		{x: cmplx.Inf(), prec: 1, want: cmplx.Inf()},
		{x: cmplx.NaN(), prec: 1, want: cmplx.NaN()},
		{x: func() complex128 { var f complex128; return -f }(), prec: 1, want: 0},
		{x: complex(math.MaxFloat64, 0) / 2, prec: 1, want: complex(math.MaxFloat64, 0) / 2},
		{x: complex(1<<64, 1<<64), prec: 1, want: complex(1<<64, 1<<64)},
		{x: 454.4445 + 454.4445i, prec: 3, want: 454.444 + 454.444i},
		{x: 454.44445 + 454.44445i, prec: 4, want: 454.4444 + 454.4444i},
		{x: 0.42499 + 0.42499i, prec: 4, want: 0.425 + 0.425i},
		{x: 0.42599 + 0.42599i, prec: 4, want: 0.426 + 0.426i},
		{x: 0.424999999999993 + 0.424999999999993i, prec: 2, want: 0.42 + 0.42i},
		{x: 0.425 + 0.425i, prec: 2, want: 0.42 + 0.42i},
		{x: 0.425000000000001 + 0.425000000000001i, prec: 2, want: 0.43 + 0.43i},
		{x: 123.4244999999999 + 123.4244999999999i, prec: 3, want: 123.424 + 123.424i},
		{x: 123.4245 + 123.4245i, prec: 3, want: 123.424 + 123.424i},
		{x: 123.4245000000001 + 123.4245000000001i, prec: 3, want: 123.425 + 123.425i},

		{x: 454.45 + 454.45i, prec: 0, want: 454 + 454i},
		{x: 454.45 + 454.45i, prec: 1, want: 454.4 + 454.4i},
		{x: 454.45 + 454.45i, prec: 2, want: 454.45 + 454.45i},
		{x: 454.45 + 454.45i, prec: 3, want: 454.45 + 454.45i},
		{x: 454.445 + 454.445i, prec: 0, want: 454 + 454i},
		{x: 454.445 + 454.445i, prec: 1, want: 454.4 + 454.4i},
		{x: 454.445 + 454.445i, prec: 2, want: 454.44 + 454.44i},
		{x: 454.445 + 454.445i, prec: 3, want: 454.445 + 454.445i},
		{x: 454.445 + 454.445i, prec: 4, want: 454.445 + 454.445i},
		{x: 454.55 + 454.55i, prec: 0, want: 455 + 455i},
		{x: 454.55 + 454.55i, prec: 1, want: 454.6 + 454.6i},
		{x: 454.55 + 454.55i, prec: 2, want: 454.55 + 454.55i},
		{x: 454.55 + 454.55i, prec: 3, want: 454.55 + 454.55i},
		{x: 454.455 + 454.455i, prec: 0, want: 454 + 454i},
		{x: 454.455 + 454.455i, prec: 1, want: 454.5 + 454.5i},
		{x: 454.455 + 454.455i, prec: 2, want: 454.46 + 454.46i},
		{x: 454.455 + 454.455i, prec: 3, want: 454.455 + 454.455i},
		{x: 454.455 + 454.455i, prec: 4, want: 454.455 + 454.455i},

		// Negative precision.
		{x: 454.45 + 454.45i, prec: -1, want: 450 + 450i},
		{x: 454.45 + 454.45i, prec: -2, want: 500 + 500i},
		{x: 500 + 500i, prec: -3, want: 0},
		{x: 500 + 500i, prec: -4, want: 0},
		{x: 1500 + 1500i, prec: -3, want: 2000 + 2000i},
		{x: 1500 + 1500i, prec: -4, want: 0},
	} {
		for _, sign := range []complex128{1, -1} {
			got := RoundEven(sign*test.x, test.prec)
			want := sign * test.want
			if want == 0 {
				want = 0
			}
			if (got != want || math.Signbit(real(got)) != math.Signbit(real(want)) || math.Signbit(imag(got)) != math.Signbit(imag(want))) && !(cmplx.IsNaN(got) && cmplx.IsNaN(want)) {
				t.Errorf("unexpected result for RoundEven(%g, %d): got: %g, want: %g", sign*test.x, test.prec, got, want)
			}
		}
	}
}

func TestSame(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	if !Same(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i + 1e-14}
	if Same(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
	if Same(s1, []complex128{}) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
	s1 = []complex128{1 + 1i, 2 + 2i, cmplx.NaN(), 4 + 4i}
	s2 = []complex128{1 + 1i, 2 + 2i, cmplx.NaN(), 4 + 4i}
	if !Same(s1, s2) {
		t.Errorf("Slices with matching NaN values returned as unequal")
	}
	s1 = []complex128{1 + 1i, 2 + 2i, cmplx.NaN(), 4 + 4i}
	s2 = []complex128{1 + 1i, cmplx.NaN(), 3 + 3i, 4 + 4i}
	if Same(s1, s2) {
		t.Errorf("Slices with unmatching NaN values returned as equal")
	}
}

func TestScale(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	c := 5.0 + 5.0i
	truth := []complex128{30i, 40i, 10i, 70i, 50i}
	Scale(c, s)
	areSlicesEqual(t, truth, s, "Bad scaling")
}

func TestScaleTo(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 + 7i, 5 + 5i}
	sCopy := make([]complex128, len(s))
	copy(sCopy, s)
	c := 5.0 + 5.0i
	truth := []complex128{30i, 40i, 10i, 70i, 50i}
	dst := make([]complex128, len(s))
	ScaleTo(dst, c, s)
	if !Same(dst, truth) {
		t.Errorf("Scale to does not match. Got %v, want %v", dst, truth)
	}
	if !Same(s, sCopy) {
		t.Errorf("Source modified during call. Got %v, want %v", s, sCopy)
	}
}

func TestSub(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 - 7i, 5 + 5i}
	v := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i}
	truth := []complex128{2 + 2i, 2 + 2i, -2 - 2i, 3 - 11i, 0}
	Sub(s, v)
	areSlicesEqual(t, truth, s, "Bad subtract")
	// Test that it panics
	if !Panics(func() { Sub(make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestSubTo(t *testing.T) {
	s := []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 - 7i, 5 + 5i}
	v := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i}
	truth := []complex128{2 + 2i, 2 + 2i, -2 - 2i, 3 - 11i, 0}
	dst1 := make([]complex128, len(s))
	dst2 := SubTo(dst1, s, v)
	areSlicesEqual(t, truth, dst1, "Bad subtract from mutator")
	areSlicesEqual(t, truth, dst2, "Bad subtract from returned slice")
	// Test that all mismatch combinations panic
	if !Panics(func() { SubTo(make([]complex128, 2), make([]complex128, 3), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with dst different length")
	}
	if !Panics(func() { SubTo(make([]complex128, 3), make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with subtractor different length")
	}
	if !Panics(func() { SubTo(make([]complex128, 3), make([]complex128, 3), make([]complex128, 2)) }) {
		t.Errorf("Did not panic with subtractee different length")
	}
}

func TestSum(t *testing.T) {
	s := []complex128{}
	val := Sum(s)
	if val != 0 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []complex128{3 + 3i, 4 + 4i, 1 + 1i, 7 - 7i, 5 + 5i}
	val = Sum(s)
	if val != (20 + 6i) {
		t.Errorf("Wrong sum returned")
	}
}

func randomSlice(l int) []complex128 {
	s := make([]complex128, l)
	for i := range s {
		s[i] = complex(rand.Float64(), rand.Float64())
	}
	return s
}

func benchmarkAdd(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Add(s1, s2)
	}
}
func BenchmarkAddSmall(b *testing.B) { benchmarkAdd(b, Small) }
func BenchmarkAddMed(b *testing.B)   { benchmarkAdd(b, Medium) }
func BenchmarkAddLarge(b *testing.B) { benchmarkAdd(b, Large) }
func BenchmarkAddHuge(b *testing.B)  { benchmarkAdd(b, Huge) }

func benchmarkAddTo(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddTo(dst, s1, s2)
	}
}
func BenchmarkAddToSmall(b *testing.B) { benchmarkAddTo(b, Small) }
func BenchmarkAddToMed(b *testing.B)   { benchmarkAddTo(b, Medium) }
func BenchmarkAddToLarge(b *testing.B) { benchmarkAddTo(b, Large) }
func BenchmarkAddToHuge(b *testing.B)  { benchmarkAddTo(b, Huge) }

func benchmarkCumProd(b *testing.B, size int) {
	s := randomSlice(size)
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CumProd(dst, s)
	}
}
func BenchmarkCumProdSmall(b *testing.B) { benchmarkCumProd(b, Small) }
func BenchmarkCumProdMed(b *testing.B)   { benchmarkCumProd(b, Medium) }
func BenchmarkCumProdLarge(b *testing.B) { benchmarkCumProd(b, Large) }
func BenchmarkCumProdHuge(b *testing.B)  { benchmarkCumProd(b, Huge) }

func benchmarkCumSum(b *testing.B, size int) {
	s := randomSlice(size)
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CumSum(dst, s)
	}
}
func BenchmarkCumSumSmall(b *testing.B) { benchmarkCumSum(b, Small) }
func BenchmarkCumSumMed(b *testing.B)   { benchmarkCumSum(b, Medium) }
func BenchmarkCumSumLarge(b *testing.B) { benchmarkCumSum(b, Large) }
func BenchmarkCumSumHuge(b *testing.B)  { benchmarkCumSum(b, Huge) }

func benchmarkDiv(b *testing.B, size int) {
	s := randomSlice(size)
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Div(dst, s)
	}
}
func BenchmarkDivSmall(b *testing.B) { benchmarkDiv(b, Small) }
func BenchmarkDivMed(b *testing.B)   { benchmarkDiv(b, Medium) }
func BenchmarkDivLarge(b *testing.B) { benchmarkDiv(b, Large) }
func BenchmarkDivHuge(b *testing.B)  { benchmarkDiv(b, Huge) }

func benchmarkDivTo(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DivTo(dst, s1, s2)
	}
}
func BenchmarkDivToSmall(b *testing.B) { benchmarkDivTo(b, Small) }
func BenchmarkDivToMed(b *testing.B)   { benchmarkDivTo(b, Medium) }
func BenchmarkDivToLarge(b *testing.B) { benchmarkDivTo(b, Large) }
func BenchmarkDivToHuge(b *testing.B)  { benchmarkDivTo(b, Huge) }

func benchmarkSub(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sub(s1, s2)
	}
}
func BenchmarkSubSmall(b *testing.B) { benchmarkSub(b, Small) }
func BenchmarkSubMed(b *testing.B)   { benchmarkSub(b, Medium) }
func BenchmarkSubLarge(b *testing.B) { benchmarkSub(b, Large) }
func BenchmarkSubHuge(b *testing.B)  { benchmarkSub(b, Huge) }

func benchmarkSubTo(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SubTo(dst, s1, s2)
	}
}
func BenchmarkSubToSmall(b *testing.B) { benchmarkSubTo(b, Small) }
func BenchmarkSubToMed(b *testing.B)   { benchmarkSubTo(b, Medium) }
func BenchmarkSubToLarge(b *testing.B) { benchmarkSubTo(b, Large) }
func BenchmarkSubToHuge(b *testing.B)  { benchmarkSubTo(b, Huge) }

func benchmarkDot(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dot(s1, s2, false)
	}
}
func BenchmarkDotSmall(b *testing.B) { benchmarkDot(b, Small) }
func BenchmarkDotMed(b *testing.B)   { benchmarkDot(b, Medium) }
func BenchmarkDotLarge(b *testing.B) { benchmarkDot(b, Large) }
func BenchmarkDotHuge(b *testing.B)  { benchmarkDot(b, Huge) }

func benchmarkAddScaledTo(b *testing.B, size int) {
	dst := randomSlice(size)
	y := randomSlice(size)
	s := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddScaledTo(dst, y, 2.3, s)
	}
}
func BenchmarkAddScaledToSmall(b *testing.B)  { benchmarkAddScaledTo(b, Small) }
func BenchmarkAddScaledToMedium(b *testing.B) { benchmarkAddScaledTo(b, Medium) }
func BenchmarkAddScaledToLarge(b *testing.B)  { benchmarkAddScaledTo(b, Large) }
func BenchmarkAddScaledToHuge(b *testing.B)   { benchmarkAddScaledTo(b, Huge) }

func benchmarkScale(b *testing.B, size int) {
	dst := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i += 2 {
		Scale(2.0, dst)
		Scale(0.5, dst)
	}
}
func BenchmarkScaleSmall(b *testing.B)  { benchmarkScale(b, Small) }
func BenchmarkScaleMedium(b *testing.B) { benchmarkScale(b, Medium) }
func BenchmarkScaleLarge(b *testing.B)  { benchmarkScale(b, Large) }
func BenchmarkScaleHuge(b *testing.B)   { benchmarkScale(b, Huge) }

func benchmarkNorm2(b *testing.B, size int) {
	s := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Norm(s, 2)
	}
}
func BenchmarkNorm2Small(b *testing.B)  { benchmarkNorm2(b, Small) }
func BenchmarkNorm2Medium(b *testing.B) { benchmarkNorm2(b, Medium) }
func BenchmarkNorm2Large(b *testing.B)  { benchmarkNorm2(b, Large) }
func BenchmarkNorm2Huge(b *testing.B)   { benchmarkNorm2(b, Huge) }
