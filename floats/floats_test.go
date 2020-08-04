// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package floats

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats/scalar"
)

const (
	EqTolerance = 1e-14
	Small       = 10
	Medium      = 1000
	Large       = 100000
	Huge        = 10000000
)

func areSlicesEqual(t *testing.T, truth, comp []float64, str string) {
	if !EqualApprox(comp, truth, EqTolerance) {
		t.Errorf(str+". Expected %v, returned %v", truth, comp)
	}
}

func areSlicesSame(t *testing.T, truth, comp []float64, str string) {
	ok := len(truth) == len(comp)
	if ok {
		for i, a := range truth {
			if !scalar.EqualWithinAbsOrRel(a, comp[i], EqTolerance, EqTolerance) && !scalar.Same(a, comp[i]) {
				ok = false
				break
			}
		}
	}
	if !ok {
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

func TestAdd(t *testing.T) {
	t.Parallel()
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	c := []float64{7, 8, 9}
	truth := []float64{12, 15, 18}
	n := make([]float64, len(a))

	Add(n, a)
	Add(n, b)
	Add(n, c)
	areSlicesEqual(t, truth, n, "Wrong addition of slices new receiver")
	Add(a, b)
	Add(a, c)
	areSlicesEqual(t, truth, n, "Wrong addition of slices for no new receiver")

	// Test that it panics
	if !Panics(func() { Add(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestAddTo(t *testing.T) {
	t.Parallel()
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	truth := []float64{5, 7, 9}
	n1 := make([]float64, len(a))

	n2 := AddTo(n1, a, b)
	areSlicesEqual(t, truth, n1, "Bad addition from mutator")
	areSlicesEqual(t, truth, n2, "Bad addition from returned slice")

	// Test that it panics
	if !Panics(func() { AddTo(make([]float64, 2), make([]float64, 3), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
	if !Panics(func() { AddTo(make([]float64, 3), make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestAddConst(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	c := 6.0
	truth := []float64{9, 10, 7, 13, 11}
	AddConst(c, s)
	areSlicesEqual(t, truth, s, "Wrong addition of constant")
}

func TestAddScaled(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	alpha := 6.0
	dst := []float64{1, 2, 3, 4, 5}
	ans := []float64{19, 26, 9, 46, 35}
	AddScaled(dst, alpha, s)
	if !EqualApprox(dst, ans, EqTolerance) {
		t.Errorf("Adding scaled did not match")
	}
	short := []float64{1}
	if !Panics(func() { AddScaled(dst, alpha, short) }) {
		t.Errorf("Doesn't panic if s is smaller than dst")
	}
	if !Panics(func() { AddScaled(short, alpha, s) }) {
		t.Errorf("Doesn't panic if dst is smaller than s")
	}
}

func TestAddScaledTo(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	alpha := 6.0
	y := []float64{1, 2, 3, 4, 5}
	dst1 := make([]float64, 5)
	ans := []float64{19, 26, 9, 46, 35}
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
	short := []float64{1}
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

func TestArgsort(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	inds := make([]int, len(s))

	Argsort(s, inds)

	sortedS := []float64{1, 3, 4, 5, 7}
	trueInds := []int{2, 0, 1, 4, 3}

	if !Equal(s, sortedS) {
		t.Error("elements not sorted correctly")
	}
	for i := range trueInds {
		if trueInds[i] != inds[i] {
			t.Error("inds not correct")
		}
	}

	inds = []int{1, 2}
	if !Panics(func() { Argsort(s, inds) }) {
		t.Error("does not panic if lengths do not match")
	}
}

func TestCount(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	f := func(v float64) bool { return v > 3.5 }
	truth := 3
	n := Count(f, s)
	if n != truth {
		t.Errorf("Wrong number of elements counted")
	}
}

func TestCumProd(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	receiver := make([]float64, len(s))
	result := CumProd(receiver, s)
	truth := []float64{3, 12, 12, 84, 420}
	areSlicesEqual(t, truth, receiver, "Wrong cumprod mutated with new receiver")
	areSlicesEqual(t, truth, result, "Wrong cumprod result with new receiver")
	CumProd(receiver, s)
	areSlicesEqual(t, truth, receiver, "Wrong cumprod returned with reused receiver")

	// Test that it panics
	if !Panics(func() { CumProd(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}

	// Test empty CumProd
	emptyReceiver := make([]float64, 0)
	truth = []float64{}
	CumProd(emptyReceiver, emptyReceiver)
	areSlicesEqual(t, truth, emptyReceiver, "Wrong cumprod returned with empty receiver")
}

func TestCumSum(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	receiver := make([]float64, len(s))
	result := CumSum(receiver, s)
	truth := []float64{3, 7, 8, 15, 20}
	areSlicesEqual(t, truth, receiver, "Wrong cumsum mutated with new receiver")
	areSlicesEqual(t, truth, result, "Wrong cumsum returned with new receiver")
	CumSum(receiver, s)
	areSlicesEqual(t, truth, receiver, "Wrong cumsum returned with reused receiver")

	// Test that it panics
	if !Panics(func() { CumSum(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}

	// Test empty CumSum
	emptyReceiver := make([]float64, 0)
	truth = []float64{}
	CumSum(emptyReceiver, emptyReceiver)
	areSlicesEqual(t, truth, emptyReceiver, "Wrong cumsum returned with empty receiver")
}

func TestDistance(t *testing.T) {
	t.Parallel()
	norms := []float64{1, 2, 4, math.Inf(1)}
	slices := []struct {
		s []float64
		t []float64
	}{
		{
			nil,
			nil,
		},
		{
			[]float64{8, 9, 10, -12},
			[]float64{8, 9, 10, -12},
		},
		{
			[]float64{1, 2, 3, -4, -5, 8},
			[]float64{-9.2, -6.8, 9, -3, -2, 1},
		},
	}

	for j, test := range slices {
		tmp := make([]float64, len(test.s))
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

	if !Panics(func() { Distance([]float64{}, norms, 1) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestDiv(t *testing.T) {
	t.Parallel()
	s1 := []float64{5, 12, 27}
	s2 := []float64{1, 2, 3}
	ans := []float64{5, 6, 9}
	Div(s1, s2)
	if !EqualApprox(s1, ans, EqTolerance) {
		t.Errorf("Div doesn't give correct answer")
	}
	s1short := []float64{1}
	if !Panics(func() { Div(s1short, s2) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
	s2short := []float64{1}
	if !Panics(func() { Div(s1, s2short) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestDivTo(t *testing.T) {
	t.Parallel()
	s1 := []float64{5, 12, 27}
	s1orig := []float64{5, 12, 27}
	s2 := []float64{1, 2, 3}
	s2orig := []float64{1, 2, 3}
	dst1 := make([]float64, 3)
	ans := []float64{5, 6, 9}
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
	dstShort := []float64{1}
	if !Panics(func() { DivTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []float64{1}
	if !Panics(func() { DivTo(dst1, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []float64{1}
	if !Panics(func() { DivTo(dst1, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestDot(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{-3, 4, 5, -6}
	truth := -4.0
	ans := Dot(s1, s2)
	if ans != truth {
		t.Errorf("Dot product computed incorrectly")
	}

	// Test that it panics
	if !Panics(func() { Dot(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestEquals(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	if !Equal(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []float64{1, 2, 3, 4 + 1e-14}
	if Equal(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
	if Equal(s1, []float64{}) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualApprox(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4 + 1e-10}
	if EqualApprox(s1, s2, 1e-13) {
		t.Errorf("Unequal slices returned as equal for absolute")
	}
	if !EqualApprox(s1, s2, 1e-5) {
		t.Errorf("Equal slices returned as unequal for absolute")
	}
	s1 = []float64{1, 2, 3, 1000}
	s2 = []float64{1, 2, 3, 1000 * (1 + 1e-7)}
	if EqualApprox(s1, s2, 1e-8) {
		t.Errorf("Unequal slices returned as equal for relative")
	}
	if !EqualApprox(s1, s2, 1e-5) {
		t.Errorf("Equal slices returned as unequal for relative")
	}
	if EqualApprox(s1, []float64{}, 1e-5) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualFunc(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	eq := func(x, y float64) bool { return x == y }
	if !EqualFunc(s1, s2, eq) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []float64{1, 2, 3, 4 + 1e-14}
	if EqualFunc(s1, s2, eq) {
		t.Errorf("Unequal slices returned as equal")
	}
	if EqualFunc(s1, []float64{}, eq) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualsRelative(t *testing.T) {
	t.Parallel()
	equalityTests := []struct {
		a, b  float64
		tol   float64
		equal bool
	}{
		{1000000, 1000001, 0, true},
		{1000001, 1000000, 0, true},
		{10000, 10001, 0, false},
		{10001, 10000, 0, false},
		{-1000000, -1000001, 0, true},
		{-1000001, -1000000, 0, true},
		{-10000, -10001, 0, false},
		{-10001, -10000, 0, false},
		{1.0000001, 1.0000002, 0, true},
		{1.0000002, 1.0000001, 0, true},
		{1.0002, 1.0001, 0, false},
		{1.0001, 1.0002, 0, false},
		{-1.000001, -1.000002, 0, true},
		{-1.000002, -1.000001, 0, true},
		{-1.0001, -1.0002, 0, false},
		{-1.0002, -1.0001, 0, false},
		{0.000000001000001, 0.000000001000002, 0, true},
		{0.000000001000002, 0.000000001000001, 0, true},
		{0.000000000001002, 0.000000000001001, 0, false},
		{0.000000000001001, 0.000000000001002, 0, false},
		{-0.000000001000001, -0.000000001000002, 0, true},
		{-0.000000001000002, -0.000000001000001, 0, true},
		{-0.000000000001002, -0.000000000001001, 0, false},
		{-0.000000000001001, -0.000000000001002, 0, false},
		{0, 0, 0, true},
		{0, -0, 0, true},
		{-0, -0, 0, true},
		{0.00000001, 0, 0, false},
		{0, 0.00000001, 0, false},
		{-0.00000001, 0, 0, false},
		{0, -0.00000001, 0, false},
		{0, 1e-310, 0.01, true},
		{1e-310, 0, 0.01, true},
		{1e-310, 0, 0.000001, false},
		{0, 1e-310, 0.000001, false},
		{0, -1e-310, 0.1, true},
		{-1e-310, 0, 0.1, true},
		{-1e-310, 0, 0.00000001, false},
		{0, -1e-310, 0.00000001, false},
		{math.Inf(1), math.Inf(1), 0, true},
		{math.Inf(-1), math.Inf(-1), 0, true},
		{math.Inf(-1), math.Inf(1), 0, false},
		{math.Inf(1), math.MaxFloat64, 0, false},
		{math.Inf(-1), -math.MaxFloat64, 0, false},
		{math.NaN(), math.NaN(), 0, false},
		{math.NaN(), 0, 0, false},
		{-0, math.NaN(), 0, false},
		{math.NaN(), -0, 0, false},
		{0, math.NaN(), 0, false},
		{math.NaN(), math.Inf(1), 0, false},
		{math.Inf(1), math.NaN(), 0, false},
		{math.NaN(), math.Inf(-1), 0, false},
		{math.Inf(-1), math.NaN(), 0, false},
		{math.NaN(), math.MaxFloat64, 0, false},
		{math.MaxFloat64, math.NaN(), 0, false},
		{math.NaN(), -math.MaxFloat64, 0, false},
		{-math.MaxFloat64, math.NaN(), 0, false},
		{math.NaN(), math.SmallestNonzeroFloat64, 0, false},
		{math.SmallestNonzeroFloat64, math.NaN(), 0, false},
		{math.NaN(), -math.SmallestNonzeroFloat64, 0, false},
		{-math.SmallestNonzeroFloat64, math.NaN(), 0, false},
		{1.000000001, -1.0, 0, false},
		{-1.0, 1.000000001, 0, false},
		{-1.000000001, 1.0, 0, false},
		{1.0, -1.000000001, 0, false},
		{10 * math.SmallestNonzeroFloat64, 10 * -math.SmallestNonzeroFloat64, 0, true},
		{1e11 * math.SmallestNonzeroFloat64, 1e11 * -math.SmallestNonzeroFloat64, 0, false},
		{math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, 0, true},
		{-math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, 0, true},
		{math.SmallestNonzeroFloat64, 0, 0, true},
		{0, math.SmallestNonzeroFloat64, 0, true},
		{-math.SmallestNonzeroFloat64, 0, 0, true},
		{0, -math.SmallestNonzeroFloat64, 0, true},
		{0.000000001, -math.SmallestNonzeroFloat64, 0, false},
		{0.000000001, math.SmallestNonzeroFloat64, 0, false},
		{math.SmallestNonzeroFloat64, 0.000000001, 0, false},
		{-math.SmallestNonzeroFloat64, 0.000000001, 0, false},
	}
	for _, ts := range equalityTests {
		if ts.tol == 0 {
			ts.tol = 1e-5
		}
		if equal := scalar.EqualWithinRel(ts.a, ts.b, ts.tol); equal != ts.equal {
			t.Errorf("Relative equality of %g and %g with tolerance %g returned: %v. Expected: %v",
				ts.a, ts.b, ts.tol, equal, ts.equal)
		}
	}
}

func TestEqualLengths(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	s3 := []float64{1, 2, 3}
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
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	f := func(v float64) bool { return v > 3.5 }
	allTrueInds := []int{1, 3, 4}

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
	t.Parallel()
	for i, test := range []struct {
		s   []float64
		ans bool
	}{
		{},
		{
			s: []float64{1, 2, 3, 4},
		},
		{
			s:   []float64{1, math.NaN(), 3, 4},
			ans: true,
		},
		{
			s:   []float64{1, 2, 3, math.NaN()},
			ans: true,
		},
	} {
		b := HasNaN(test.s)
		if b != test.ans {
			t.Errorf("HasNaN mismatch case %d. Expected %v, Found %v", i, test.ans, b)
		}
	}
}

func TestLogSpan(t *testing.T) {
	t.Parallel()
	receiver1 := make([]float64, 6)
	truth := []float64{0.001, 0.01, 0.1, 1, 10, 100}
	receiver2 := LogSpan(receiver1, 0.001, 100)
	tst := make([]float64, 6)
	for i := range truth {
		tst[i] = receiver1[i] / truth[i]
	}
	comp := make([]float64, 6)
	for i := range comp {
		comp[i] = 1
	}
	areSlicesEqual(t, comp, tst, "Improper logspace from mutator")

	for i := range truth {
		tst[i] = receiver2[i] / truth[i]
	}
	areSlicesEqual(t, comp, tst, "Improper logspace from returned slice")

	if !Panics(func() { LogSpan(nil, 1, 5) }) {
		t.Errorf("Span accepts nil argument")
	}
	if !Panics(func() { LogSpan(make([]float64, 1), 1, 5) }) {
		t.Errorf("Span accepts argument of len = 1")
	}
}

func TestLogSumExp(t *testing.T) {
	t.Parallel()
	s := []float64{1, 2, 3, 4, 5}
	val := LogSumExp(s)
	// http://www.wolframalpha.com/input/?i=log%28exp%281%29+%2B+exp%282%29+%2B+exp%283%29+%2B+exp%284%29+%2B+exp%285%29%29
	truth := 5.4519143959375933331957225109748087179338972737576824
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Wrong logsumexp for many values")
	}
	s = []float64{1, 2}
	// http://www.wolframalpha.com/input/?i=log%28exp%281%29+%2B+exp%282%29%29
	truth = 2.3132616875182228340489954949678556419152800856703483
	val = LogSumExp(s)
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Wrong logsumexp for two values. %v expected, %v found", truth, val)
	}
	// This case would normally underflow
	s = []float64{-1001, -1002, -1003, -1004, -1005}
	// http://www.wolframalpha.com/input/?i=log%28exp%28-1001%29%2Bexp%28-1002%29%2Bexp%28-1003%29%2Bexp%28-1004%29%2Bexp%28-1005%29%29
	truth = -1000.54808560406240666680427748902519128206610272624
	val = LogSumExp(s)
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for underflow case. %v expected, %v found", truth, val)
	}
	// positive infinite case
	s = []float64{1, 2, 3, 4, 5, math.Inf(1)}
	val = LogSumExp(s)
	truth = math.Inf(1)
	if val != truth {
		t.Errorf("Doesn't match for pos Infinity case. %v expected, %v found", truth, val)
	}
	// negative infinite case
	s = []float64{1, 2, 3, 4, 5, math.Inf(-1)}
	val = LogSumExp(s)
	truth = 5.4519143959375933331957225109748087179338972737576824 // same as first case
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Wrong logsumexp for values with negative infinity")
	}
}

func TestMaxAndIdx(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		in      []float64
		wantIdx int
		wantVal float64
		desc    string
	}{
		{
			in:      []float64{3, 4, 1, 7, 5},
			wantIdx: 3,
			wantVal: 7,
			desc:    "with only finite entries",
		},
		{
			in:      []float64{math.NaN(), 4, 1, 7, 5},
			wantIdx: 3,
			wantVal: 7,
			desc:    "with leading NaN",
		},
		{
			in:      []float64{math.NaN(), math.NaN(), math.NaN()},
			wantIdx: 0,
			wantVal: math.NaN(),
			desc:    "when only NaN elements exist",
		},
		{
			in:      []float64{math.NaN(), math.Inf(-1)},
			wantIdx: 1,
			wantVal: math.Inf(-1),
			desc:    "leading NaN followed by -Inf",
		},
		{
			in:      []float64{math.NaN(), math.Inf(1)},
			wantIdx: 1,
			wantVal: math.Inf(1),
			desc:    "leading NaN followed by +Inf",
		},
	} {
		ind := MaxIdx(test.in)
		if ind != test.wantIdx {
			t.Errorf("Wrong index "+test.desc+": got:%d want:%d", ind, test.wantIdx)
		}
		val := Max(test.in)
		if !scalar.Same(val, test.wantVal) {
			t.Errorf("Wrong value "+test.desc+": got:%f want:%f", val, test.wantVal)
		}
	}
	if !Panics(func() { MaxIdx([]float64{}) }) {
		t.Errorf("Expected panic with zero length")
	}
}

func TestMinAndIdx(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		in      []float64
		wantIdx int
		wantVal float64
		desc    string
	}{
		{
			in:      []float64{3, 4, 1, 7, 5},
			wantIdx: 2,
			wantVal: 1,
			desc:    "with only finite entries",
		},
		{
			in:      []float64{math.NaN(), 4, 1, 7, 5},
			wantIdx: 2,
			wantVal: 1,
			desc:    "with leading NaN",
		},
		{
			in:      []float64{math.NaN(), math.NaN(), math.NaN()},
			wantIdx: 0,
			wantVal: math.NaN(),
			desc:    "when only NaN elements exist",
		},
		{
			in:      []float64{math.NaN(), math.Inf(-1)},
			wantIdx: 1,
			wantVal: math.Inf(-1),
			desc:    "leading NaN followed by -Inf",
		},
		{
			in:      []float64{math.NaN(), math.Inf(1)},
			wantIdx: 1,
			wantVal: math.Inf(1),
			desc:    "leading NaN followed by +Inf",
		},
	} {
		ind := MinIdx(test.in)
		if ind != test.wantIdx {
			t.Errorf("Wrong index "+test.desc+": got:%d want:%d", ind, test.wantIdx)
		}
		val := Min(test.in)
		if !scalar.Same(val, test.wantVal) {
			t.Errorf("Wrong value "+test.desc+": got:%f want:%f", val, test.wantVal)
		}
	}
	if !Panics(func() { MinIdx([]float64{}) }) {
		t.Errorf("Expected panic with zero length")
	}
}

func TestMul(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3}
	s2 := []float64{1, 2, 3}
	ans := []float64{1, 4, 9}
	Mul(s1, s2)
	if !EqualApprox(s1, ans, EqTolerance) {
		t.Errorf("Mul doesn't give correct answer")
	}
	s1short := []float64{1}
	if !Panics(func() { Mul(s1short, s2) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
	s2short := []float64{1}
	if !Panics(func() { Mul(s1, s2short) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestMulTo(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3}
	s1orig := []float64{1, 2, 3}
	s2 := []float64{1, 2, 3}
	s2orig := []float64{1, 2, 3}
	dst1 := make([]float64, 3)
	ans := []float64{1, 4, 9}
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
	dstShort := []float64{1}
	if !Panics(func() { MulTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []float64{1}
	if !Panics(func() { MulTo(dst1, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []float64{1}
	if !Panics(func() { MulTo(dst1, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestNearestIdx(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		in    []float64
		query float64
		want  int
		desc  string
	}{
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 2,
			want:  1,
			desc:  "Wrong index returned when value is less than all of elements",
		},
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 9,
			want:  4,
			desc:  "Wrong index returned when value is greater than all of elements",
		},
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 3.1,
			want:  1,
			desc:  "Wrong index returned when value is greater than closest element",
		},
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 2.9,
			want:  1,
			desc:  "Wrong index returned when value is less than closest element",
		},
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 3,
			want:  1,
			desc:  "Wrong index returned when value is equal to element",
		},
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 6.2,
			want:  0,
			desc:  "Wrong index returned when value is equal to several elements",
		},
		{
			in:    []float64{6.2, 3, 5, 6.2, 8},
			query: 4,
			want:  1,
			desc:  "Wrong index returned when value is exactly between two closest elements",
		},
		{
			in:    []float64{math.NaN(), 3, 2, -1},
			query: 2,
			want:  2,
			desc:  "Wrong index returned when initial element is NaN",
		},
		{
			in:    []float64{0, math.NaN(), -1, 2},
			query: math.NaN(),
			want:  0,
			desc:  "Wrong index returned when query is NaN and a NaN element exists",
		},
		{
			in:    []float64{0, math.NaN(), -1, 2},
			query: math.Inf(1),
			want:  3,
			desc:  "Wrong index returned when query is +Inf and no +Inf element exists",
		},
		{
			in:    []float64{0, math.NaN(), -1, 2},
			query: math.Inf(-1),
			want:  2,
			desc:  "Wrong index returned when query is -Inf and no -Inf element exists",
		},
		{
			in:    []float64{math.NaN(), math.NaN(), math.NaN()},
			query: 1,
			want:  0,
			desc:  "Wrong index returned when query is a number and only NaN elements exist",
		},
		{
			in:    []float64{math.NaN(), math.Inf(-1)},
			query: 1,
			want:  1,
			desc:  "Wrong index returned when query is a number and single NaN precedes -Inf",
		},
	} {
		ind := NearestIdx(test.in, test.query)
		if ind != test.want {
			t.Errorf(test.desc+": got:%d want:%d", ind, test.want)
		}
	}
	if !Panics(func() { NearestIdx([]float64{}, 0) }) {
		t.Errorf("Expected panic with zero length")
	}
}

func TestNearestIdxForSpan(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		length int
		lower  float64
		upper  float64
		value  float64
		idx    int
	}{
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  6,
			idx:    0,
		},
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  10,
			idx:    12,
		},
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  7.19,
			idx:    2,
		},
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  7.21,
			idx:    2,
		},
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  7.2,
			idx:    2,
		},
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  7.151,
			idx:    2,
		},
		{
			length: 13,
			lower:  7,
			upper:  8.2,
			value:  7.249,
			idx:    2,
		},
		{
			length: 4,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  math.Copysign(0, -1),
			idx:    0,
		},
		{
			length: 5,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  0,
			idx:    2,
		},
		{
			length: 5,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  math.Inf(-1),
			idx:    0,
		},
		{
			length: 5,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  math.Inf(1),
			idx:    3,
		},
		{
			length: 4,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  0,
			idx:    2,
		},
		{
			length: 4,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  math.Inf(1),
			idx:    2,
		},
		{
			length: 4,
			lower:  math.Inf(-1),
			upper:  math.Inf(1),
			value:  math.Inf(-1),
			idx:    0,
		},
		{
			length: 5,
			lower:  math.Inf(1),
			upper:  math.Inf(1),
			value:  1,
			idx:    0,
		},
		{
			length: 5,
			lower:  math.NaN(),
			upper:  math.NaN(),
			value:  1,
			idx:    0,
		},
		{
			length: 5,
			lower:  0,
			upper:  1,
			value:  math.NaN(),
			idx:    0,
		},
		{
			length: 5,
			lower:  math.NaN(),
			upper:  1,
			value:  0,
			idx:    4,
		},
		{
			length: 5,
			lower:  math.Inf(-1),
			upper:  1,
			value:  math.Inf(-1),
			idx:    0,
		},
		{
			length: 5,
			lower:  math.Inf(-1),
			upper:  1,
			value:  0,
			idx:    4,
		},
		{
			length: 5,
			lower:  math.Inf(1),
			upper:  1,
			value:  math.Inf(1),
			idx:    0,
		},
		{
			length: 5,
			lower:  math.Inf(1),
			upper:  1,
			value:  0,
			idx:    4,
		},
		{
			length: 5,
			lower:  100,
			upper:  math.Inf(-1),
			value:  math.Inf(-1),
			idx:    4,
		},
		{
			length: 5,
			lower:  100,
			upper:  math.Inf(-1),
			value:  200,
			idx:    0,
		},
		{
			length: 5,
			lower:  100,
			upper:  math.Inf(1),
			value:  math.Inf(1),
			idx:    4,
		},
		{
			length: 5,
			lower:  100,
			upper:  math.Inf(1),
			value:  200,
			idx:    0,
		},
		{
			length: 5,
			lower:  -1,
			upper:  2,
			value:  math.Inf(-1),
			idx:    0,
		},
		{
			length: 5,
			lower:  -1,
			upper:  2,
			value:  math.Inf(1),
			idx:    4,
		},
		{
			length: 5,
			lower:  1,
			upper:  -2,
			value:  math.Inf(-1),
			idx:    4,
		},
		{
			length: 5,
			lower:  1,
			upper:  -2,
			value:  math.Inf(1),
			idx:    0,
		},
		{
			length: 5,
			lower:  2,
			upper:  0,
			value:  3,
			idx:    0,
		},
		{
			length: 5,
			lower:  2,
			upper:  0,
			value:  -1,
			idx:    4,
		},
	} {
		if idx := NearestIdxForSpan(test.length, test.lower, test.upper, test.value); test.idx != idx {
			t.Errorf("Case %v mismatch: Want: %v, Got: %v", i, test.idx, idx)
		}
	}
	if !Panics(func() { NearestIdxForSpan(1, 0, 1, 0.5) }) {
		t.Errorf("Expected panic for short span length")
	}
}

func TestNorm(t *testing.T) {
	t.Parallel()
	s := []float64{-1, -3.4, 5, -6}
	val := Norm(s, math.Inf(1))
	truth := 6.0
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28-1%29%5E2+%2B++%28-3.4%29%5E2+%2B+5%5E2%2B++6%5E2%29%5E%281%2F2%29
	val = Norm(s, 2)
	truth = 8.5767126569566267590651614132751986658027271236078592
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28%7C-1%7C%29%5E3+%2B++%28%7C-3.4%7C%29%5E3+%2B+%7C5%7C%5E3%2B++%7C6%7C%5E3%29%5E%281%2F3%29
	val = Norm(s, 3)
	truth = 7.2514321388020228478109121239004816430071237369356233
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}

	//http://www.wolframalpha.com/input/?i=%7C-1%7C+%2B+%7C-3.4%7C+%2B+%7C5%7C%2B++%7C6%7C
	val = Norm(s, 1)
	truth = 15.4
	if math.Abs(val-truth) > EqTolerance {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
}

func TestProd(t *testing.T) {
	t.Parallel()
	s := []float64{}
	val := Prod(s)
	if val != 1 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []float64{3, 4, 1, 7, 5}
	val = Prod(s)
	if val != 420 {
		t.Errorf("Wrong prod returned. Expected %v returned %v", 420, val)
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()
	for _, s := range [][]float64{
		{0},
		{1, 0},
		{2, 1, 0},
		{3, 2, 1, 0},
		{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	} {
		Reverse(s)
		for i, v := range s {
			if v != float64(i) {
				t.Errorf("unexpected values for element %d: got:%v want:%v", i, v, i)
			}
		}
	}
}

func TestSame(t *testing.T) {
	t.Parallel()
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	if !Same(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []float64{1, 2, 3, 4 + 1e-14}
	if Same(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
	if Same(s1, []float64{}) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
	s1 = []float64{1, 2, math.NaN(), 4}
	s2 = []float64{1, 2, math.NaN(), 4}
	if !Same(s1, s2) {
		t.Errorf("Slices with matching NaN values returned as unequal")
	}
	s1 = []float64{1, 2, math.NaN(), 4}
	s2 = []float64{1, math.NaN(), 3, 4}
	if Same(s1, s2) {
		t.Errorf("Slices with unmatching NaN values returned as equal")
	}
}

func TestScale(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	c := 5.0
	truth := []float64{15, 20, 5, 35, 25}
	Scale(c, s)
	areSlicesEqual(t, truth, s, "Bad scaling")
}

func TestScaleTo(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	sCopy := make([]float64, len(s))
	copy(sCopy, s)
	c := 5.0
	truth := []float64{15, 20, 5, 35, 25}
	dst := make([]float64, len(s))
	ScaleTo(dst, c, s)
	if !Same(dst, truth) {
		t.Errorf("Scale to does not match. Got %v, want %v", dst, truth)
	}
	if !Same(s, sCopy) {
		t.Errorf("Source modified during call. Got %v, want %v", s, sCopy)
	}
	if !Panics(func() { ScaleTo(dst, 0, []float64{1}) }) {
		t.Errorf("Expected panic with different slice lengths")
	}
}

func TestSpan(t *testing.T) {
	t.Parallel()
	receiver1 := make([]float64, 5)
	truth := []float64{1, 2, 3, 4, 5}
	receiver2 := Span(receiver1, 1, 5)
	areSlicesEqual(t, truth, receiver1, "Improper linspace from mutator")
	areSlicesEqual(t, truth, receiver2, "Improper linspace from returned slice")
	receiver1 = make([]float64, 6)
	truth = []float64{0, 0.2, 0.4, 0.6, 0.8, 1.0}
	Span(receiver1, 0, 1)
	areSlicesEqual(t, truth, receiver1, "Improper linspace")
	if !Panics(func() { Span(nil, 1, 5) }) {
		t.Errorf("Span accepts nil argument")
	}
	if !Panics(func() { Span(make([]float64, 1), 1, 5) }) {
		t.Errorf("Span accepts argument of len = 1")
	}

	for _, test := range []struct {
		n    int
		l, u float64
		want []float64
	}{
		{
			n: 4, l: math.Inf(-1), u: math.Inf(1),
			want: []float64{math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1)},
		},
		{
			n: 4, l: math.Inf(1), u: math.Inf(-1),
			want: []float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)},
		},
		{
			n: 5, l: math.Inf(-1), u: math.Inf(1),
			want: []float64{math.Inf(-1), math.Inf(-1), 0, math.Inf(1), math.Inf(1)},
		},
		{
			n: 5, l: math.Inf(1), u: math.Inf(-1),
			want: []float64{math.Inf(1), math.Inf(1), 0, math.Inf(-1), math.Inf(-1)},
		},
		{
			n: 5, l: math.Inf(1), u: math.Inf(1),
			want: []float64{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)},
		},
		{
			n: 5, l: math.Inf(-1), u: math.Inf(-1),
			want: []float64{math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1)},
		},
		{
			n: 5, l: math.Inf(-1), u: math.NaN(),
			want: []float64{math.Inf(-1), math.NaN(), math.NaN(), math.NaN(), math.NaN()},
		},
		{
			n: 5, l: math.Inf(1), u: math.NaN(),
			want: []float64{math.Inf(1), math.NaN(), math.NaN(), math.NaN(), math.NaN()},
		},
		{
			n: 5, l: math.NaN(), u: math.Inf(-1),
			want: []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.Inf(-1)},
		},
		{
			n: 5, l: math.NaN(), u: math.Inf(1),
			want: []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.Inf(1)},
		},
		{
			n: 5, l: 42, u: math.Inf(-1),
			want: []float64{42, math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1)},
		},
		{
			n: 5, l: 42, u: math.Inf(1),
			want: []float64{42, math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)},
		},
		{
			n: 5, l: 42, u: math.NaN(),
			want: []float64{42, math.NaN(), math.NaN(), math.NaN(), math.NaN()},
		},
		{
			n: 5, l: math.Inf(-1), u: 42,
			want: []float64{math.Inf(-1), math.Inf(-1), math.Inf(-1), math.Inf(-1), 42},
		},
		{
			n: 5, l: math.Inf(1), u: 42,
			want: []float64{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), 42},
		},
		{
			n: 5, l: math.NaN(), u: 42,
			want: []float64{math.NaN(), math.NaN(), math.NaN(), math.NaN(), 42},
		},
	} {
		got := Span(make([]float64, test.n), test.l, test.u)
		areSlicesSame(t, test.want, got,
			fmt.Sprintf("Unexpected slice of length %d for %f to %f", test.n, test.l, test.u))
	}
}

func TestSub(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	v := []float64{1, 2, 3, 4, 5}
	truth := []float64{2, 2, -2, 3, 0}
	Sub(s, v)
	areSlicesEqual(t, truth, s, "Bad subtract")
	// Test that it panics
	if !Panics(func() { Sub(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestSubTo(t *testing.T) {
	t.Parallel()
	s := []float64{3, 4, 1, 7, 5}
	v := []float64{1, 2, 3, 4, 5}
	truth := []float64{2, 2, -2, 3, 0}
	dst1 := make([]float64, len(s))
	dst2 := SubTo(dst1, s, v)
	areSlicesEqual(t, truth, dst1, "Bad subtract from mutator")
	areSlicesEqual(t, truth, dst2, "Bad subtract from returned slice")
	// Test that all mismatch combinations panic
	if !Panics(func() { SubTo(make([]float64, 2), make([]float64, 3), make([]float64, 3)) }) {
		t.Errorf("Did not panic with dst different length")
	}
	if !Panics(func() { SubTo(make([]float64, 3), make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with subtractor different length")
	}
	if !Panics(func() { SubTo(make([]float64, 3), make([]float64, 3), make([]float64, 2)) }) {
		t.Errorf("Did not panic with subtractee different length")
	}
}

func TestSum(t *testing.T) {
	t.Parallel()
	s := []float64{}
	val := Sum(s)
	if val != 0 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []float64{3, 4, 1, 7, 5}
	val = Sum(s)
	if val != 20 {
		t.Errorf("Wrong sum returned")
	}
}

func TestWithin(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		s      []float64
		v      float64
		idx    int
		panics bool
	}{
		{
			s:   []float64{1, 2, 5, 9},
			v:   1,
			idx: 0,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   9,
			idx: -1,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   1.5,
			idx: 0,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   2,
			idx: 1,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   2.5,
			idx: 1,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   -3,
			idx: -1,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   15,
			idx: -1,
		},
		{
			s:   []float64{1, 2, 5, 9},
			v:   math.NaN(),
			idx: -1,
		},
		{
			s:      []float64{5, 2, 6},
			panics: true,
		},
		{
			panics: true,
		},
		{
			s:      []float64{1},
			panics: true,
		},
	} {
		var idx int
		panics := Panics(func() { idx = Within(test.s, test.v) })
		if panics {
			if !test.panics {
				t.Errorf("Case %v: bad panic", i)
			}
			continue
		}
		if test.panics {
			if !panics {
				t.Errorf("Case %v: did not panic when it should", i)
			}
			continue
		}
		if idx != test.idx {
			t.Errorf("Case %v: Idx mismatch. Want: %v, got: %v", i, test.idx, idx)
		}
	}
}

func TestSumCompensated(t *testing.T) {
	t.Parallel()
	k := 100000
	s1 := make([]float64, 2*k+1)
	for i := -k; i <= k; i++ {
		s1[i+k] = 0.2 * float64(i)
	}
	s2 := make([]float64, k+1)
	for i := 0; i < k; i++ {
		s2[i] = 10. / float64(k)
	}
	s2[k] = -10

	for i, test := range []struct {
		s    []float64
		want float64
	}{
		{
			// Fails if we use simple Sum.
			s:    s1,
			want: 0,
		},
		{
			// Fails if we use simple Sum.
			s:    s2,
			want: 0,
		},
		{
			s:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want: 55,
		},
		{
			s:    []float64{1.2e20, 0.1, -2.4e20, -0.1, 1.2e20, 0.2, 0.2},
			want: 0.4,
		},
		{
			s:    []float64{1, 1e100, 1, -1e100},
			want: 2,
		},
	} {
		got := SumCompensated(test.s)
		if math.Abs(got-test.want) > EqTolerance {
			t.Errorf("Wrong sum returned in test case %d. Want: %g, got: %g", i, test.want, got)
		}
	}
}

func randomSlice(l int) []float64 {
	s := make([]float64, l)
	for i := range s {
		s[i] = rand.Float64()
	}
	return s
}

func benchmarkMin(b *testing.B, size int) {
	s := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Min(s)
	}
}
func BenchmarkMinSmall(b *testing.B) { benchmarkMin(b, Small) }
func BenchmarkMinMed(b *testing.B)   { benchmarkMin(b, Medium) }
func BenchmarkMinLarge(b *testing.B) { benchmarkMin(b, Large) }
func BenchmarkMinHuge(b *testing.B)  { benchmarkMin(b, Huge) }

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

func benchmarkLogSumExp(b *testing.B, size int) {
	s := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LogSumExp(s)
	}
}
func BenchmarkLogSumExpSmall(b *testing.B) { benchmarkLogSumExp(b, Small) }
func BenchmarkLogSumExpMed(b *testing.B)   { benchmarkLogSumExp(b, Medium) }
func BenchmarkLogSumExpLarge(b *testing.B) { benchmarkLogSumExp(b, Large) }
func BenchmarkLogSumExpHuge(b *testing.B)  { benchmarkLogSumExp(b, Huge) }

func benchmarkDot(b *testing.B, size int) {
	s1 := randomSlice(size)
	s2 := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dot(s1, s2)
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

func benchmarkSumCompensated(b *testing.B, size int) {
	s := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumCompensated(s)
	}
}

func BenchmarkSumCompensatedSmall(b *testing.B)  { benchmarkSumCompensated(b, Small) }
func BenchmarkSumCompensatedMedium(b *testing.B) { benchmarkSumCompensated(b, Medium) }
func BenchmarkSumCompensatedLarge(b *testing.B)  { benchmarkSumCompensated(b, Large) }
func BenchmarkSumCompensatedHuge(b *testing.B)   { benchmarkSumCompensated(b, Huge) }

func benchmarkSum(b *testing.B, size int) {
	s := randomSlice(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(s)
	}
}

func BenchmarkSumSmall(b *testing.B)  { benchmarkSum(b, Small) }
func BenchmarkSumMedium(b *testing.B) { benchmarkSum(b, Medium) }
func BenchmarkSumLarge(b *testing.B)  { benchmarkSum(b, Large) }
func BenchmarkSumHuge(b *testing.B)   { benchmarkSum(b, Huge) }
