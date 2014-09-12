// Copyright 2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package floats

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
)

const (
	EQTOLERANCE = 1E-14
	SMALL       = 10
	MEDIUM      = 1000
	LARGE       = 100000
	HUGE        = 10000000
)

func AreSlicesEqual(t *testing.T, truth, comp []float64, str string) {
	if !EqualApprox(comp, truth, EQTOLERANCE) {
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
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	c := []float64{7, 8, 9}
	truth := []float64{12, 15, 18}
	n := make([]float64, len(a))
	Add(n, a, b, c)
	AreSlicesEqual(t, truth, n, "Wrong addition of slices new receiver")
	Add(a, b, c)
	AreSlicesEqual(t, truth, n, "Wrong addition of slices for no new receiver")
	// Test that it panics
	if !Panics(func() { Add(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestAddConst(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	c := 6.0
	truth := []float64{9, 10, 7, 13, 11}
	AddConst(c, s)
	AreSlicesEqual(t, truth, s, "Wrong addition of constant")
}

func TestAddScaled(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	alpha := 6.0
	dst := []float64{1, 2, 3, 4, 5}
	ans := []float64{19, 26, 9, 46, 35}
	AddScaled(dst, alpha, s)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
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
	s := []float64{3, 4, 1, 7, 5}
	alpha := 6.0
	y := []float64{1, 2, 3, 4, 5}
	dst := make([]float64, 5)
	ans := []float64{19, 26, 9, 46, 35}
	AddScaledTo(dst, y, alpha, s)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
		t.Errorf("AddScaledTo did not match")
	}
	AddScaledTo(dst, y, alpha, s)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
		t.Errorf("Reusing dst did not match")
	}
	short := []float64{1}
	if !Panics(func() { AddScaledTo(dst, y, alpha, short) }) {
		t.Errorf("Doesn't panic if s is smaller than dst")
	}
	if !Panics(func() { AddScaledTo(short, y, alpha, s) }) {
		t.Errorf("Doesn't panic if dst is smaller than s")
	}
	if !Panics(func() { AddScaledTo(dst, short, alpha, s) }) {
		t.Errorf("Doesn't panic if y is smaller than dst")
	}
}

func TestApply(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	f := math.Sin
	truth := make([]float64, len(s))
	for i, val := range s {
		truth[i] = math.Sin(val)
	}
	Apply(f, s)
	AreSlicesEqual(t, truth, s, "Wrong application of function")
}

func TestArgsort(t *testing.T) {
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
	s := []float64{3, 4, 1, 7, 5}
	f := func(v float64) bool { return v > 3.5 }
	truth := 3
	n := Count(f, s)
	if n != truth {
		t.Errorf("Wrong number of elements counted")
	}
}

func TestCumProd(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	receiver := make([]float64, len(s))
	CumProd(receiver, s)
	truth := []float64{3, 12, 12, 84, 420}
	AreSlicesEqual(t, truth, receiver, "Wrong cumprod returned with new receiver")
	CumProd(receiver, s)
	AreSlicesEqual(t, truth, receiver, "Wrong cumprod returned with reused receiver")
	// Test that it panics
	if !Panics(func() { CumProd(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestCumSum(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	receiver := make([]float64, len(s))
	CumSum(receiver, s)
	truth := []float64{3, 7, 8, 15, 20}
	AreSlicesEqual(t, truth, receiver, "Wrong cumsum returned with new receiver")
	CumSum(receiver, s)
	AreSlicesEqual(t, truth, receiver, "Wrong cumsum returned with reused receiver")

	// Test that it panics
	if !Panics(func() { CumSum(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestDistance(t *testing.T) {
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
			if dist != norm { // Use equality because they should be identical
				t.Errorf("Distance does not match norm for case %v, %v. Expected %v, Found %v.", i, j, norm, dist)
			}
		}
	}

}

func TestDiv(t *testing.T) {
	s1 := []float64{5, 12, 27}
	s2 := []float64{1, 2, 3}
	ans := []float64{5, 6, 9}
	Div(s1, s2)
	if !EqualApprox(s1, ans, EQTOLERANCE) {
		t.Errorf("Mul doesn't give correct answer")
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
	s1 := []float64{5, 12, 27}
	s1orig := []float64{5, 12, 27}
	s2 := []float64{1, 2, 3}
	s2orig := []float64{1, 2, 3}
	dst := make([]float64, 3)
	ans := []float64{5, 6, 9}
	DivTo(dst, s1, s2)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
		t.Errorf("DivTo doesn't give correct answer")
	}
	if !EqualApprox(s1, s1orig, EQTOLERANCE) {
		t.Errorf("S1 changes during multo")
	}
	if !EqualApprox(s2, s2orig, EQTOLERANCE) {
		t.Errorf("s2 changes during multo")
	}
	DivTo(dst, s1, s2)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
		t.Errorf("DivTo doesn't give correct answer reusing dst")
	}
	dstShort := []float64{1}
	if !Panics(func() { DivTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []float64{1}
	if !Panics(func() { DivTo(dst, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []float64{1}
	if !Panics(func() { DivTo(dst, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestDot(t *testing.T) {
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
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{1, 2, 3, 4}
	if !Equal(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []float64{1, 2, 3, 4 + 1e-14}
	if Equal(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
}

func TestEqualApprox(t *testing.T) {
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
}

func TestEqualsRelative(t *testing.T) {
	var equalityTests = []struct {
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
		if equal := EqualWithinRel(ts.a, ts.b, ts.tol); equal != ts.equal {
			t.Errorf("Relative equality of %g and %g with tolerance %g returned: %v. Expected: %v",
				ts.a, ts.b, ts.tol, equal, ts.equal)
		}
	}
}

func nextAfterN(x, y float64, n int) float64 {
	for i := 0; i < n; i++ {
		x = math.Nextafter(x, y)
	}
	return x
}

func TestEqualsULP(t *testing.T) {
	if f := 67329.242; !EqualWithinULP(f, nextAfterN(f, math.Inf(1), 10), 10) {
		t.Errorf("Equal values returned as unequal")
	}
	if f := 67329.242; EqualWithinULP(f, nextAfterN(f, math.Inf(1), 5), 1) {
		t.Errorf("Unequal values returned as equal")
	}
	if f := nextAfterN(0, math.Inf(1), 2); !EqualWithinULP(f, nextAfterN(f, math.Inf(-1), 5), 10) {
		t.Errorf("Equal values returned as unequal")
	}
}

func TestEqualLengths(t *testing.T) {
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
	receiver := make([]float64, 6)
	truth := []float64{0.001, 0.01, 0.1, 1, 10, 100}
	LogSpan(receiver, 0.001, 100)
	tst := make([]float64, 6)
	for i := range truth {
		tst[i] = receiver[i] / truth[i]
	}
	comp := make([]float64, 6)
	for i := range comp {
		comp[i] = 1
	}
	AreSlicesEqual(t, comp, tst, "Improper logspace")
	if !Panics(func() { LogSpan(nil, 1, 5) }) {
		t.Errorf("Span accepts nil argument")
	}
	if !Panics(func() { LogSpan(make([]float64, 1), 1, 5) }) {
		t.Errorf("Span accepts argument of len = 1")
	}
}

func TestLogSumExp(t *testing.T) {
	s := []float64{1, 2, 3, 4, 5}
	val := LogSumExp(s)
	// http://www.wolframalpha.com/input/?i=log%28exp%281%29+%2B+exp%282%29+%2B+exp%283%29+%2B+exp%284%29+%2B+exp%285%29%29
	truth := 5.4519143959375933331957225109748087179338972737576824
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Wrong logsumexp for many values")
	}
	s = []float64{1, 2}
	// http://www.wolframalpha.com/input/?i=log%28exp%281%29+%2B+exp%282%29%29
	truth = 2.3132616875182228340489954949678556419152800856703483
	val = LogSumExp(s)
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Wrong logsumexp for two values. %v expected, %v found", truth, val)
	}
	// This case would normally underflow
	s = []float64{-1001, -1002, -1003, -1004, -1005}
	// http://www.wolframalpha.com/input/?i=log%28exp%28-1001%29%2Bexp%28-1002%29%2Bexp%28-1003%29%2Bexp%28-1004%29%2Bexp%28-1005%29%29
	truth = -1000.54808560406240666680427748902519128206610272624
	val = LogSumExp(s)
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for underflow case. %v expected, %v found", truth, val)
	}
}

func TestMax(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	val, ind := Max(s)
	if val != 7 {
		t.Errorf("Wrong value returned")
	}
	if ind != 3 {
		t.Errorf("Wrong index returned")
	}
}

func TestMin(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	val, ind := Min(s)
	if val != 1 {
		t.Errorf("Wrong value returned")
	}
	if ind != 2 {
		t.Errorf("Wrong index returned")
	}
}

func TestMul(t *testing.T) {
	s1 := []float64{1, 2, 3}
	s2 := []float64{1, 2, 3}
	ans := []float64{1, 4, 9}
	Mul(s1, s2)
	if !EqualApprox(s1, ans, EQTOLERANCE) {
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
	s1 := []float64{1, 2, 3}
	s1orig := []float64{1, 2, 3}
	s2 := []float64{1, 2, 3}
	s2orig := []float64{1, 2, 3}
	dst := make([]float64, 3)
	ans := []float64{1, 4, 9}
	MulTo(dst, s1, s2)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
		t.Errorf("MulTo doesn't give correct answer")
	}
	if !EqualApprox(s1, s1orig, EQTOLERANCE) {
		t.Errorf("S1 changes during multo")
	}
	if !EqualApprox(s2, s2orig, EQTOLERANCE) {
		t.Errorf("s2 changes during multo")
	}
	MulTo(dst, s1, s2)
	if !EqualApprox(dst, ans, EQTOLERANCE) {
		t.Errorf("MulTo doesn't give correct answer reusing dst")
	}
	dstShort := []float64{1}
	if !Panics(func() { MulTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []float64{1}
	if !Panics(func() { MulTo(dst, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []float64{1}
	if !Panics(func() { MulTo(dst, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestNearest(t *testing.T) {
	s := []float64{6.2, 3, 5, 6.2, 8}
	ind := Nearest(s, 2.0)
	if ind != 1 {
		t.Errorf("Wrong index returned when value is less than all of elements")
	}
	ind = Nearest(s, 9.0)
	if ind != 4 {
		t.Errorf("Wrong index returned when value is greater than all of elements")
	}
	ind = Nearest(s, 3.1)
	if ind != 1 {
		t.Errorf("Wrong index returned when value is greater than closest element")
	}
	ind = Nearest(s, 3.1)
	if ind != 1 {
		t.Errorf("Wrong index returned when value is greater than closest element")
	}
	ind = Nearest(s, 2.9)
	if ind != 1 {
		t.Errorf("Wrong index returned when value is less than closest element")
	}
	ind = Nearest(s, 3)
	if ind != 1 {
		t.Errorf("Wrong index returned when value is equal to element")
	}
	ind = Nearest(s, 6.2)
	if ind != 0 {
		t.Errorf("Wrong index returned when value is equal to several elements")
	}
	ind = Nearest(s, 4)
	if ind != 1 {
		t.Errorf("Wrong index returned when value is exactly between two closest elements")
	}
}

func TestNearestWithinSpan(t *testing.T) {

	if !Panics(func() { NearestWithinSpan(13, 7, 8.2, 10) }) {
		t.Errorf("Did not panic below lower bound")
	}
	if !Panics(func() { NearestWithinSpan(13, 7, 8.2, 10) }) {
		t.Errorf("Did not panic above upper bound")
	}
	ind := NearestWithinSpan(13, 7, 8.2, 7.19)
	if ind != 2 {
		t.Errorf("Wrong value when just below the bucket. %i found, %i expected", ind, 2)
	}
	ind = NearestWithinSpan(13, 7, 8.2, 7.21)
	if ind != 2 {
		t.Errorf("Wrong value when just above the bucket. %i found, %i expected", ind, 2)
	}
	ind = NearestWithinSpan(13, 7, 8.2, 7.2)
	if ind != 2 {
		t.Errorf("Wrong value when equal to bucket. %i found, %i expected", ind, 2)
	}
	ind = NearestWithinSpan(13, 7, 8.2, 7.151)
	if ind != 2 {
		t.Errorf("Wrong value when just above halfway point. %i found, %i expected", ind, 2)
	}
	ind = NearestWithinSpan(13, 7, 8.2, 7.249)
	if ind != 2 {
		t.Errorf("Wrong value when just below halfway point. %i found, %i expected", ind, 2)
	}
}

func TestNorm(t *testing.T) {
	s := []float64{-1, -3.4, 5, 6}
	val := Norm(s, math.Inf(1))
	truth := 6.0
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28-1%29%5E2+%2B++%28-3.4%29%5E2+%2B+5%5E2%2B++6%5E2%29%5E%281%2F2%29
	val = Norm(s, 2)
	truth = 8.5767126569566267590651614132751986658027271236078592
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
	// http://www.wolframalpha.com/input/?i=%28%28%7C-1%7C%29%5E3+%2B++%28%7C-3.4%7C%29%5E3+%2B+%7C5%7C%5E3%2B++%7C6%7C%5E3%29%5E%281%2F3%29
	val = Norm(s, 3)
	truth = 7.2514321388020228478109121239004816430071237369356233
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}

	//http://www.wolframalpha.com/input/?i=%7C-1%7C+%2B+%7C-3.4%7C+%2B+%7C5%7C%2B++%7C6%7C
	val = Norm(s, 1)
	truth = 15.4
	if math.Abs(val-truth) > EQTOLERANCE {
		t.Errorf("Doesn't match for inf norm. %v expected, %v found", truth, val)
	}
}

func TestProd(t *testing.T) {
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

func TestScale(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	c := 5.0
	truth := []float64{15, 20, 5, 35, 25}
	Scale(c, s)
	AreSlicesEqual(t, truth, s, "Bad scaling")
}

func TestSpan(t *testing.T) {
	receiver := make([]float64, 5)
	truth := []float64{1, 2, 3, 4, 5}
	Span(receiver, 1, 5)
	AreSlicesEqual(t, truth, receiver, "Improper linspace")
	receiver = make([]float64, 6)
	truth = []float64{0, 0.2, 0.4, 0.6, 0.8, 1.0}
	Span(receiver, 0, 1)
	AreSlicesEqual(t, truth, receiver, "Improper linspace")
	if !Panics(func() { Span(nil, 1, 5) }) {
		t.Errorf("Span accepts nil argument")
	}
	if !Panics(func() { Span(make([]float64, 1), 1, 5) }) {
		t.Errorf("Span accepts argument of len = 1")
	}
}

func TestSub(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	v := []float64{1, 2, 3, 4, 5}
	truth := []float64{2, 2, -2, 3, 0}
	Sub(s, v)
	AreSlicesEqual(t, truth, s, "Bad subtract")
	// Test that it panics
	if !Panics(func() { Sub(make([]float64, 2), make([]float64, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestSubTo(t *testing.T) {
	s := []float64{3, 4, 1, 7, 5}
	v := []float64{1, 2, 3, 4, 5}
	truth := []float64{2, 2, -2, 3, 0}
	dst := make([]float64, len(s))
	SubTo(dst, s, v)
	AreSlicesEqual(t, truth, dst, "Bad subtract")
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

func RandomSlice(l int) []float64 {
	s := make([]float64, l)
	for i := range s {
		s[i] = rand.Float64()
	}
	return s
}

func benchmarkMin(b *testing.B, s []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Min(s)
	}
}

func BenchmarkMinSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMin(b, s)
}

func BenchmarkMinMed(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMin(b, s)
}

func BenchmarkMinLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMin(b, s)
}
func BenchmarkMinHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMin(b, s)
}

func benchmarkAdd(b *testing.B, s ...[]float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Add(s[0], s[1:]...)
	}
}

func BenchmarkAddTwoSmall(b *testing.B) {
	i := SMALL
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourSmall(b *testing.B) {
	i := SMALL
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func BenchmarkAddTwoMed(b *testing.B) {
	i := MEDIUM
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourMed(b *testing.B) {
	i := MEDIUM
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func BenchmarkAddTwoLarge(b *testing.B) {
	i := LARGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourLarge(b *testing.B) {
	i := LARGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func BenchmarkAddTwoHuge(b *testing.B) {
	i := HUGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourHuge(b *testing.B) {
	i := HUGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func benchmarkLogSumExp(b *testing.B, s []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LogSumExp(s)
	}
}

func BenchmarkLogSumExpSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkLogSumExp(b, s)
}

func BenchmarkLogSumExpMed(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkLogSumExp(b, s)
}

func BenchmarkLogSumExpLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkLogSumExp(b, s)
}
func BenchmarkLogSumExpHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkLogSumExp(b, s)
}

func benchmarkDot(b *testing.B, s1 []float64, s2 []float64) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dot(s1, s2)
	}
}

func BenchmarkDotSmall(b *testing.B) {
	s1 := RandomSlice(SMALL)
	s2 := RandomSlice(SMALL)
	benchmarkDot(b, s1, s2)
}

func BenchmarkDotMed(b *testing.B) {
	s1 := RandomSlice(MEDIUM)
	s2 := RandomSlice(MEDIUM)
	benchmarkDot(b, s1, s2)
}

func BenchmarkDotLarge(b *testing.B) {
	s1 := RandomSlice(LARGE)
	s2 := RandomSlice(LARGE)
	benchmarkDot(b, s1, s2)
}
func BenchmarkDotHuge(b *testing.B) {
	s1 := RandomSlice(HUGE)
	s2 := RandomSlice(HUGE)
	benchmarkDot(b, s1, s2)
}
