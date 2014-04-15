// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"testing"
)

func TestSame(t *testing.T) {
	var (
		a = make(set)
		b = make(set)
		c = a
	)

	if same(a, b) {
		t.Error("Independently created sets test as same")
	}
	if !same(a, c) {
		t.Error("Set copy and original test as not same.")
	}
	a.add(1)
	if !same(a, c) {
		t.Error("Set copy and original test as not same after addition.")
	}
	if !same(nil, nil) {
		t.Error("nil sets test as not same.")
	}
	if same(b, nil) {
		t.Error("nil and empty sets test as same.")
	}
}

func TestAdd(t *testing.T) {
	s := make(set)
	if s == nil {
		t.Fatal("Set cannot be created successfully")
	}

	if s.count() != 0 {
		t.Error("Set somehow contains new elements upon creation")
	}

	s.add(1)
	s.add(3)
	s.add(5)

	if s.count() != 3 {
		t.Error("Incorrect number of set elements after adding")
	}

	if !s.has(1) || !s.has(3) || !s.has(5) {
		t.Error("Set doesn't contain element that was added")
	}

	s.add(1)

	if s.count() > 3 {
		t.Error("Set double-adds element (element not unique)")
	} else if s.count() < 3 {
		t.Error("Set double-add lowered len")
	}

	if !s.has(1) {
		t.Error("Set doesn't contain double-added element")
	}

	if !s.has(3) || !s.has(5) {
		t.Error("Set removes element on double-add")
	}

}

func TestRemove(t *testing.T) {
	s := make(set)

	s.add(1)
	s.add(3)
	s.add(5)

	s.remove(1)

	if s.count() != 2 {
		t.Error("Incorrect number of set elements after removing an element")
	}

	if s.has(1) {
		t.Error("Element present after removal")
	}

	if !s.has(3) || !s.has(5) {
		t.Error("Set remove removed wrong element")
	}

	s.remove(1)

	if s.count() != 2 || s.has(1) {
		t.Error("Double set remove does something strange")
	}

	s.add(1)

	if s.count() != 3 || !s.has(1) {
		t.Error("Cannot add element after removal")
	}

}

func TestElements(t *testing.T) {
	s := make(set)
	el := s.elements()
	if el == nil {
		t.Errorf("elements of empty set incorrectly returns nil and not zero-length slice")
	}

	if len(el) != 0 {
		t.Errorf("elements of empty set has len greater than 0")
	}

	s.add(1)
	s.add(2)

	el = s.elements()
	if len(el) != 2 {
		t.Fatalf("elements not of same len as set that spawned it")
	}

	if e := el[0]; e != 1 && e != 2 {
		t.Errorf("Element in elements has incorrect value %d", e)
	}

	if e := el[1]; e != 1 && e != 2 {
		t.Errorf("Element in elements has incorrect value %d", e)
	}

	el[0] = 19
	el[1] = 52

	if !s.has(1) || !s.has(2) || s.count() != 2 {
		t.Error("Mutating elements slice mutates set")
	}
}

func TestClear(t *testing.T) {
	s := make(set)

	s.add(8)
	s.add(9)
	s.add(10)

	s = clear(s)

	if s.count() != 0 {
		t.Error("Clear did not properly reset set to size 0")
	}
}

func TestSelfEqual(t *testing.T) {
	s := make(set)

	if !equal(s, s) {
		t.Error("Set is not equal to itself")
	}

	s.add(1)

	if !equal(s, s) {
		t.Error("Set ceases self equality after adding element")
	}
}

func TestEqual(t *testing.T) {
	s1 := make(set)
	s2 := make(set)

	if !equal(s1, s2) {
		t.Error("Two different empty sets not equal")
	}

	s1.add(1)
	if equal(s1, s2) {
		t.Error("Two different sets with different elements not equal")
	}

	s2.add(1)
	if !equal(s1, s2) {
		t.Error("Two sets with same element not equal")
	}
}

func TestCopy(t *testing.T) {
	s1 := make(set)
	s2 := make(set)

	s1.add(1)
	s1.add(2)
	s1.add(3)

	s2.copy(s1)

	if !equal(s1, s2) {
		t.Fatalf("Two sets not equal after copy")
	}

	s2.remove(1)

	if equal(s1, s2) {
		t.Errorf("Mutating one set mutated another after copy")
	}
}

func TestSelfCopy(t *testing.T) {
	s1 := make(set)

	s1.add(1)
	s1.add(2)

	s1.copy(s1)

	if s1.count() != 2 {
		t.Error("Something strange happened when copying into self")
	}
}

func TestUnionSame(t *testing.T) {
	s1 := make(set)
	s2 := make(set)
	s3 := make(set)

	s1.add(1)
	s1.add(2)

	s2.add(1)
	s2.add(2)

	s3.union(s1, s2)

	if s3.count() != 2 {
		t.Error("Union of same sets yields set with wrong len")
	}

	if !s3.has(1) || !s3.has(2) {
		t.Error("Union of same sets yields wrong elements")
	}
}

func TestUnionDiff(t *testing.T) {
	s1 := make(set)
	s2 := make(set)
	s3 := make(set)

	s1.add(1)
	s1.add(2)

	s2.add(3)

	s3.union(s1, s2)

	if s3.count() != 3 {
		t.Error("Union of different sets yields set with wrong len")
	}

	if !s3.has(1) || !s3.has(2) || !s3.has(3) {
		t.Error("Union of different sets yields set with wrong elements")
	}

	if s1.has(3) || !s1.has(2) || !s1.has(1) || s1.count() != 2 {
		t.Error("Union of sets mutates non-destination set (argument 1)")
	}

	if !s2.has(3) || s2.has(1) || s2.has(2) || s2.count() != 1 {
		t.Error("Union of sets mutates non-destination set (argument 2)")
	}
}

func TestUnionOverlapping(t *testing.T) {
	s1 := make(set)
	s2 := make(set)
	s3 := make(set)

	s1.add(1)
	s1.add(2)

	s2.add(2)
	s2.add(3)

	s3.union(s1, s2)

	if s3.count() != 3 {
		t.Error("Union of overlapping sets yields set with wrong len")
	}

	if !s3.has(1) || !s3.has(2) || !s3.has(3) {
		t.Error("Union of overlapping sets yields set with wrong elements")
	}

	if s1.has(3) || !s1.has(2) || !s1.has(1) || s1.count() != 2 {
		t.Error("Union of sets mutates non-destination set (argument 1)")
	}

	if !s2.has(3) || s2.has(1) || !s2.has(2) || s2.count() != 2 {
		t.Error("Union of sets mutates non-destination set (argument 2)")
	}
}

func TestIntersectSame(t *testing.T) {
	s1 := make(set)
	s2 := make(set)
	s3 := make(set)

	s1.add(2)
	s1.add(3)

	s2.add(2)
	s2.add(3)

	s3.intersect(s1, s2)

	if card := s3.count(); card != 2 {
		t.Errorf("Intersection of identical sets yields set of wrong len %d", card)
	}

	if !s3.has(2) || !s3.has(3) {
		t.Error("Intersection of identical sets yields set of wrong elements")
	}
}

func TestIntersectDiff(t *testing.T) {
	s1 := make(set)
	s2 := make(set)
	s3 := make(set)

	s1.add(2)
	s1.add(3)

	s2.add(1)
	s2.add(4)

	s3.intersect(s1, s2)

	if card := s3.count(); card != 0 {
		t.Errorf("Intersection of different yields non-empty set %d", card)
	}

	if !s1.has(2) || !s1.has(3) || s1.has(1) || s1.has(4) || s1.count() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}

	if s2.has(2) || s2.has(3) || !s2.has(1) || !s2.has(4) || s2.count() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}
}

func TestIntersectOverlapping(t *testing.T) {
	s1 := make(set)
	s2 := make(set)
	s3 := make(set)

	s1.add(2)
	s1.add(3)

	s2.add(3)
	s2.add(4)

	s3.intersect(s1, s2)

	if card := s3.count(); card != 1 {
		t.Errorf("Intersection of overlapping sets yields set of incorrect len %d", card)
	}

	if !s3.has(3) {
		t.Errorf("Intersection of overlapping sets yields set with wrong element")
	}

	if !s1.has(2) || !s1.has(3) || s1.has(4) || s1.count() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}

	if s2.has(2) || !s2.has(3) || !s2.has(4) || s2.count() != 2 {
		t.Error("Intersection of sets mutates non-destination set (argument 1)")
	}
}
