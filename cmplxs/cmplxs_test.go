// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cmplxs

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/cmplxs/cscalar"
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

func areSlicesSame(t *testing.T, truth, comp []complex128, str string) {
	ok := len(truth) == len(comp)
	if ok {
		for i, a := range truth {
			if !cscalar.EqualWithinAbsOrRel(a, comp[i], EqTolerance, EqTolerance) && !cscalar.Same(a, comp[i]) {
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

func TestAddConst(t *testing.T) {
	s := []complex128{3 + 1i, 4 + 2i, 1 + 3i, 7 + 4i, 5 + 5i}
	c := 6 + 1i
	truth := []complex128{9 + 2i, 10 + 3i, 7 + 4i, 13 + 5i, 11 + 6i}
	AddConst(c, s)
	areSlicesEqual(t, truth, s, "Wrong addition of constant")
}

func TestAddScaled(t *testing.T) {
	s := []complex128{3, 4, 1, 7, 5}
	alpha := 6 + 1i
	dst := []complex128{1, 2, 3, 4, 5}
	ans := []complex128{19 + 3i, 26 + 4i, 9 + 1i, 46 + 7i, 35 + 5i}
	AddScaled(dst, alpha, s)
	if !EqualApprox(dst, ans, EqTolerance) {
		t.Errorf("Adding scaled did not match. Expected %v, returned %v", ans, dst)
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
	s := []complex128{3, 4, 1, 7, 5}
	alpha := 6 + 1i
	y := []complex128{1, 2, 3, 4, 5}
	dst1 := make([]complex128, 5)
	ans := []complex128{19 + 3i, 26 + 4i, 9 + 1i, 46 + 7i, 35 + 5i}
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

func TestCount(t *testing.T) {
	s := []complex128{3, 4, 1, 7, 5}
	f := func(v complex128) bool { return cmplx.Abs(v) > 3.5 }
	truth := 3
	n := Count(f, s)
	if n != truth {
		t.Errorf("Wrong number of elements counted")
	}
}

func TestCumProd(t *testing.T) {
	s := []complex128{3 + 1i, 4 + 2i, 1 + 3i, 7 + 4i, 5 + 5i}
	receiver := make([]complex128, len(s))
	result := CumProd(receiver, s)
	truth := []complex128{3 + 1i, 10 + 10i, -20 + 40i, -300 + 200i, -2500 - 500i}
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

func TestComplex(t *testing.T) {
	for i, test := range []struct {
		dst        []complex128
		real, imag []float64
		want       []complex128
		panics     bool
	}{
		{},
		{
			dst:  make([]complex128, 4),
			real: []float64{1, 2, 3, 4},
			imag: []float64{1, 2, 3, 4},
			want: []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
		},
		{
			dst:    make([]complex128, 3),
			real:   []float64{1, 2, 3, 4},
			imag:   []float64{1, 2, 3, 4},
			panics: true,
		},
		{
			dst:    make([]complex128, 4),
			real:   []float64{1, 2, 3},
			imag:   []float64{1, 2, 3, 4},
			panics: true,
		},
		{
			dst:    make([]complex128, 4),
			real:   []float64{1, 2, 3, 4},
			imag:   []float64{1, 2, 3},
			panics: true,
		},
		{
			dst:  make([]complex128, 4),
			real: []float64{1, 2, 3, 4},
			imag: []float64{1, 2, 3, math.NaN()},
			want: []complex128{1 + 1i, 2 + 2i, 3 + 3i, cmplx.NaN()},
		},
	} {
		var got []complex128
		panicked := Panics(func() {
			got = Complex(test.dst, test.real, test.imag)
		})
		if panicked != test.panics {
			if panicked {
				t.Errorf("unexpected panic for test %d", i)
			} else {
				t.Errorf("expected panic for test %d", i)
			}
		}
		if panicked || test.panics {
			continue
		}
		if !Same(got, test.dst) {
			t.Errorf("mismatch between dst and return test %d: got:%v want:%v", i, got, test.dst)
		}
		if !Same(got, test.want) {
			t.Errorf("unexpected result for test %d: got:%v want:%v", i, got, test.want)
		}
	}

}

func TestCumSum(t *testing.T) {
	s := []complex128{3 + 1i, 4 + 2i, 1 + 3i, 7 + 4i, 5 + 5i}
	receiver := make([]complex128, len(s))
	result := CumSum(receiver, s)
	truth := []complex128{3 + 1i, 7 + 3i, 8 + 6i, 15 + 10i, 20 + 15i}
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
			[]complex128{8 + 1i, 9 + 2i, 10 + 3i, -12 + 4i},
			[]complex128{8 + 1i, 9 + 2i, 10 + 3i, -12 + 4i},
		},
		{
			[]complex128{1 + 1i, 2 + 2i, 3 + 3i, -4 + 4i, -5 + 5i, 8 + 6i},
			[]complex128{-9.2 - 1i, -6.8 - 2i, 9 - 3i, -3 - 4i, -2 - 5i, 1 - 6i},
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

	if !Panics(func() { Distance([]complex128{}, []complex128{1}, 1) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestDiv(t *testing.T) {
	s1 := []complex128{5 - 5i, 12 + 2i, 27 - 3i}
	s2 := []complex128{1 - 1i, 2 + 2i, 3 - 1i}
	ans := []complex128{5 + 0i, 3.5 - 2.5i, 8.4 + 1.8i}
	Div(s1, s2)
	if !EqualApprox(s1, ans, EqTolerance) {
		t.Errorf("Div doesn't give correct answer. Expected %v, Found %v.", ans, s1)
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
	s1 := []complex128{5 - 5i, 12 + 2i, 27 - 3i}
	s1orig := []complex128{5 - 5i, 12 + 2i, 27 - 3i}
	s2 := []complex128{1 - 1i, 2 + 2i, 3 - 1i}
	s2orig := []complex128{1 - 1i, 2 + 2i, 3 - 1i}
	dst1 := make([]complex128, 3)
	ans := []complex128{5 + 0i, 3.5 - 2.5i, 8.4 + 1.8i}
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
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i}
	s2 := []complex128{-3 + 4i, 4 + 3i, 5 + 2i, -6 + 1i}
	truth := 16 + 24i
	ans := Dot(s1, s2)
	if ans != truth {
		t.Errorf("Dot product computed incorrectly. Expected %v, Found %v.", truth, ans)
	}

	// Test that it panics
	if !Panics(func() { Dot(make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestEquals(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i}
	s2 := []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i}
	if !Equal(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i + 1e-14}
	if Equal(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
	if Equal(s1, []complex128{}) {
		t.Errorf("Unequal slice lengths returned as equal")
	}
}

func TestEqualApprox(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i}
	s2 := []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i + 1e-10}
	if EqualApprox(s1, s2, 1e-13) {
		t.Errorf("Unequal slices returned as equal for absolute")
	}
	if !EqualApprox(s1, s2, 1e-5) {
		t.Errorf("Equal slices returned as unequal for absolute")
	}
	s1 = []complex128{1 + 1i, 2 + 4i, 3 + 8i, 1000 + 1000i}
	s2 = []complex128{1 + 1i, 2 + 4i, 3 + 8i, (1000 + 1000i) * (1 + 1e-7)}
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
	s1 := []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i}
	s2 := []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i}
	eq := func(x, y complex128) bool { return x == y }
	if !EqualFunc(s1, s2, eq) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []complex128{1 + 1i, 2 + 4i, 3 + 8i, 4 + 16i + 1e-14}
	if EqualFunc(s1, s2, eq) {
		t.Errorf("Unequal slices returned as equal")
	}
	if EqualFunc(s1, []complex128{}, eq) {
		t.Errorf("Unequal slice lengths returned as equal")
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
	s := []complex128{3 + 1i, 4 - 1i, 1 + 2i, 7 + 10i, 5 - 100i}
	f := func(v complex128) bool { return cmplx.Abs(v) > 3.5 }
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

func TestImag(t *testing.T) {
	for i, test := range []struct {
		dst    []float64
		src    []complex128
		want   []float64
		panics bool
	}{
		{},
		{
			dst:  make([]float64, 4),
			src:  []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			want: []float64{1, 2, 3, 4},
		},
		{
			dst:    make([]float64, 3),
			src:    []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			panics: true,
		},
		{
			dst:  make([]float64, 4),
			src:  []complex128{1 + 1i, 2 + 2i, 3 + 3i, cmplx.NaN()},
			want: []float64{1, 2, 3, math.NaN()},
		},
	} {
		var got []float64
		panicked := Panics(func() {
			got = Imag(test.dst, test.src)
		})
		if panicked != test.panics {
			if panicked {
				t.Errorf("unexpected panic for test %d", i)
			} else {
				t.Errorf("expected panic for test %d", i)
			}
		}
		if panicked || test.panics {
			continue
		}
		if !floats.Same(got, test.dst) {
			t.Errorf("mismatch between dst and return test %d: got:%v want:%v", i, got, test.dst)
		}
		if !floats.Same(got, test.want) {
			t.Errorf("unexpected result for test %d: got:%v want:%v", i, got, test.want)
		}
	}

}

func TestLogSpan(t *testing.T) {
	// FIXME(kortschak)
	receiver1 := make([]complex128, 6)
	truth := []complex128{0.001, 0.01, 0.1, 1, 10, 100}
	receiver2 := LogSpan(receiver1, 0.001, 100)
	tst := make([]complex128, 6)
	for i := range truth {
		tst[i] = receiver1[i] / truth[i]
	}
	comp := make([]complex128, 6)
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
	if !Panics(func() { LogSpan(make([]complex128, 1), 1, 5) }) {
		t.Errorf("Span accepts argument of len = 1")
	}
}

func TestMaxAbsAndIdx(t *testing.T) {
	for _, test := range []struct {
		in      []complex128
		wantIdx int
		wantVal complex128
		desc    string
	}{
		{
			in:      []complex128{3 + 1i, 4 + 1i, 1 + 1i, 7 + 1i, 5 + 1i},
			wantIdx: 3,
			wantVal: 7 + 1i,
			desc:    "with only finite entries",
		},
		{
			in:      []complex128{cmplx.NaN(), 4 + 1i, 1 + 1i, 7 + 1i, 5 + 1i},
			wantIdx: 3,
			wantVal: 7 + 1i,
			desc:    "with leading NaN",
		},
		{
			in:      []complex128{cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
			wantIdx: 0,
			wantVal: cmplx.NaN(),
			desc:    "when only NaN elements exist",
		},
		{
			in:      []complex128{cmplx.NaN(), cmplx.Inf()},
			wantIdx: 1,
			wantVal: cmplx.Inf(),
			desc:    "leading NaN followed by Inf",
		},
	} {
		ind := MaxAbsIdx(test.in)
		if ind != test.wantIdx {
			t.Errorf("Wrong index "+test.desc+": got:%d want:%d", ind, test.wantIdx)
		}
		val := MaxAbs(test.in)
		if !cscalar.Same(val, test.wantVal) {
			t.Errorf("Wrong value "+test.desc+": got:%f want:%f", val, test.wantVal)
		}
	}
}

func TestMinAbsAndIdx(t *testing.T) {
	for _, test := range []struct {
		in      []complex128
		wantIdx int
		wantVal complex128
		desc    string
	}{
		{
			in:      []complex128{3 + 1i, 4 + 1i, 1 + 1i, 7 + 1i, 5 + 1i},
			wantIdx: 2,
			wantVal: 1 + 1i,
			desc:    "with only finite entries",
		},
		{
			in:      []complex128{cmplx.NaN(), 4 + 1i, 1 + 1i, 7 + 1i, 5 + 1i},
			wantIdx: 2,
			wantVal: 1 + 1i,
			desc:    "with leading NaN",
		},
		{
			in:      []complex128{cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
			wantIdx: 0,
			wantVal: cmplx.NaN(),
			desc:    "when only NaN elements exist",
		},
		{
			in:      []complex128{cmplx.NaN(), cmplx.Inf()},
			wantIdx: 1,
			wantVal: cmplx.Inf(),
			desc:    "leading NaN followed by Inf",
		},
	} {
		ind := MinAbsIdx(test.in)
		if ind != test.wantIdx {
			t.Errorf("Wrong index "+test.desc+": got:%d want:%d", ind, test.wantIdx)
		}
		val := MinAbs(test.in)
		if !cscalar.Same(val, test.wantVal) {
			t.Errorf("Wrong value "+test.desc+": got:%f want:%f", val, test.wantVal)
		}
	}
}

func TestMul(t *testing.T) {
	s1 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	ans := []complex128{0 + 2i, 0 + 8i, 0 + 18i}
	Mul(s1, s2)
	if !EqualApprox(s1, ans, EqTolerance) {
		t.Errorf("Mul doesn't give correct answer. Expected %v, Found %v", ans, s1)
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
	s2 := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	s2orig := []complex128{1 + 1i, 2 + 2i, 3 + 3i}
	dst1 := make([]complex128, 3)
	ans := []complex128{0 + 2i, 0 + 8i, 0 + 18i}
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

// Complexify

func TestNearestIdx(t *testing.T) {
	for _, test := range []struct {
		in    []complex128
		query complex128
		want  int
		desc  string
	}{
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 2,
			want:  1,
			desc:  "Wrong index returned when value is less than all of elements",
		},
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 9,
			want:  4,
			desc:  "Wrong index returned when value is greater than all of elements",
		},
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 3.1,
			want:  1,
			desc:  "Wrong index returned when value is greater than closest element",
		},
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 2.9,
			want:  1,
			desc:  "Wrong index returned when value is less than closest element",
		},
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 3,
			want:  1,
			desc:  "Wrong index returned when value is equal to element",
		},
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 6.2,
			want:  0,
			desc:  "Wrong index returned when value is equal to several elements",
		},
		{
			in:    []complex128{6.2, 3, 5, 6.2, 8},
			query: 4,
			want:  1,
			desc:  "Wrong index returned when value is exactly between two closest elements",
		},
		{
			in:    []complex128{cmplx.NaN(), 3, 2, -1},
			query: 2,
			want:  2,
			desc:  "Wrong index returned when initial element is NaN",
		},
		{
			in:    []complex128{0, cmplx.NaN(), -1, 2},
			query: cmplx.NaN(),
			want:  0,
			desc:  "Wrong index returned when query is NaN and a NaN element exists",
		},
		{
			in:    []complex128{0, cmplx.NaN(), -1, 2},
			query: cmplx.Inf(),
			want:  3,
			desc:  "Wrong index returned when query is Inf and no Inf element exists",
		},
		{
			in:    []complex128{cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
			query: 1,
			want:  0,
			desc:  "Wrong index returned when query is a number and only NaN elements exist",
		},
		{
			in:    []complex128{cmplx.NaN(), cmplx.Inf()},
			query: 1,
			want:  1,
			desc:  "Wrong index returned when query is a number and single NaN precedes Inf",
		},
	} {
		ind := NearestIdx(test.in, test.query)
		if ind != test.want {
			t.Errorf(test.desc+": got:%d want:%d", ind, test.want)
		}
	}
}

func TestNorm(t *testing.T) {
	s := []complex128{-1, -3.4, 5, -6}
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
	s := []complex128{}
	val := Prod(s)
	if val != 1 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []complex128{3, 4, 1, 7, 5}
	val = Prod(s)
	if val != 420 {
		t.Errorf("Wrong prod returned. Expected %v returned %v", 420, val)
	}
}

func TestReverse(t *testing.T) {
	for _, s := range [][]complex128{
		{0},
		{1, 0},
		{2, 1, 0},
		{3, 2, 1, 0},
		{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	} {
		Reverse(s)
		for i, v := range s {
			if v != complex(float64(i), 0) {
				t.Errorf("unexpected values for element %d: got:%v want:%v", i, v, i)
			}
		}
	}
}

func TestReal(t *testing.T) {
	for i, test := range []struct {
		dst    []float64
		src    []complex128
		want   []float64
		panics bool
	}{
		{},
		{
			dst:  make([]float64, 4),
			src:  []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			want: []float64{1, 2, 3, 4},
		},
		{
			dst:    make([]float64, 3),
			src:    []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i},
			panics: true,
		},
		{
			dst:  make([]float64, 4),
			src:  []complex128{1 + 1i, 2 + 2i, 3 + 3i, cmplx.NaN()},
			want: []float64{1, 2, 3, math.NaN()},
		},
	} {
		var got []float64
		panicked := Panics(func() {
			got = Real(test.dst, test.src)
		})
		if panicked != test.panics {
			if panicked {
				t.Errorf("unexpected panic for test %d", i)
			} else {
				t.Errorf("expected panic for test %d", i)
			}
		}
		if panicked || test.panics {
			continue
		}
		if !floats.Same(got, test.dst) {
			t.Errorf("mismatch between dst and return test %d: got:%v want:%v", i, got, test.dst)
		}
		if !floats.Same(got, test.want) {
			t.Errorf("unexpected result for test %d: got:%v want:%v", i, got, test.want)
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
	s := []complex128{3, 4, 1, 7, 5}
	c := 5 + 5i
	truth := []complex128{15 + 15i, 20 + 20i, 5 + 5i, 35 + 35i, 25 + 25i}
	Scale(c, s)
	areSlicesEqual(t, truth, s, "Bad scaling")
}

func TestScaleTo(t *testing.T) {
	s := []complex128{3, 4, 1, 7, 5}
	sCopy := make([]complex128, len(s))
	copy(sCopy, s)
	c := 5 + 5i
	truth := []complex128{15 + 15i, 20 + 20i, 5 + 5i, 35 + 35i, 25 + 25i}
	dst := make([]complex128, len(s))
	ScaleTo(dst, c, s)
	if !Same(dst, truth) {
		t.Errorf("Scale to does not match. Got %v, want %v", dst, truth)
	}
	if !Same(s, sCopy) {
		t.Errorf("Source modified during call. Got %v, want %v", s, sCopy)
	}
}

func TestSpan(t *testing.T) {
	receiver1 := make([]complex128, 5)
	truth := []complex128{1 + 1i, 2 + 2i, 3 + 3i, 4 + 4i, 5 + 5i}
	receiver2 := Span(receiver1, 1+1i, 5+5i)
	areSlicesEqual(t, truth, receiver1, "Improper linspace from mutator")
	areSlicesEqual(t, truth, receiver2, "Improper linspace from returned slice")
	receiver1 = make([]complex128, 6)
	truth = []complex128{0, 0.2 + 0.4i, 0.4 + 0.8i, 0.6 + 1.2i, 0.8 + 1.6i, 1 + 2i}
	Span(receiver1, 0, 1+2i)
	areSlicesEqual(t, truth, receiver1, "Improper linspace")
	if !Panics(func() { Span(nil, 1, 5) }) {
		t.Errorf("Span accepts nil argument")
	}
	if !Panics(func() { Span(make([]complex128, 1), 1, 5) }) {
		t.Errorf("Span accepts argument of len = 1")
	}

	for _, test := range []struct {
		n    int
		l, u complex128
		want []complex128
	}{
		{
			n: 5, l: cmplx.Inf(), u: cmplx.Inf(),
			want: []complex128{cmplx.Inf(), cmplx.Inf(), cmplx.Inf(), cmplx.Inf(), cmplx.Inf()},
		},
		{
			n: 5, l: cmplx.Inf(), u: cmplx.NaN(),
			want: []complex128{cmplx.Inf(), cmplx.NaN(), cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
		},
		{
			n: 5, l: cmplx.NaN(), u: cmplx.Inf(),
			want: []complex128{cmplx.NaN(), cmplx.NaN(), cmplx.NaN(), cmplx.NaN(), cmplx.Inf()},
		},
		{
			n: 5, l: 42, u: cmplx.Inf(),
			want: []complex128{42, cmplx.Inf(), cmplx.Inf(), cmplx.Inf(), cmplx.Inf()},
		},
		{
			n: 5, l: 42, u: cmplx.NaN(),
			want: []complex128{42, cmplx.NaN(), cmplx.NaN(), cmplx.NaN(), cmplx.NaN()},
		},
		{
			n: 5, l: cmplx.Inf(), u: 42,
			want: []complex128{cmplx.Inf(), cmplx.Inf(), cmplx.Inf(), cmplx.Inf(), 42},
		},
		{
			n: 5, l: cmplx.NaN(), u: 42,
			want: []complex128{cmplx.NaN(), cmplx.NaN(), cmplx.NaN(), cmplx.NaN(), 42},
		},
	} {
		got := Span(make([]complex128, test.n), test.l, test.u)
		areSlicesSame(t, test.want, got,
			fmt.Sprintf("Unexpected slice of length %d for %f to %f", test.n, test.l, test.u))
	}
}

func TestSub(t *testing.T) {
	s := []complex128{3 + 2i, 4 + 3i, 1 + 7i, 7 + 1i, 5 - 1i}
	v := []complex128{1 + 1i, 2 + 4i, 3, 4, 5 - 1i}
	truth := []complex128{2 + 1i, 2 - 1i, -2 + 7i, 3 + 1i, 0}
	Sub(s, v)
	areSlicesEqual(t, truth, s, "Bad subtract")
	// Test that it panics
	if !Panics(func() { Sub(make([]complex128, 2), make([]complex128, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestSubTo(t *testing.T) {
	s := []complex128{3 + 2i, 4 + 3i, 1 + 7i, 7 + 1i, 5 - 1i}
	v := []complex128{1 + 1i, 2 + 4i, 3, 4, 5 - 1i}
	truth := []complex128{2 + 1i, 2 - 1i, -2 + 7i, 3 + 1i, 0}
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
	s = []complex128{3 + 1i, 4 + 2i, 1 + 3i, 7 + 4i, 5 + 5i}
	val = Sum(s)
	if val != 20+15i {
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
