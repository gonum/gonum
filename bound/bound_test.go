// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: Rewrite these tests as table-driven.

package bound

import (
	"math"
	"testing"
)

func TestIntersection(t *testing.T) {
	b0 := Bound{Min: 0.0, Max: 5.0}
	b1 := Bound{Min: 3.0, Max: 4.0}
	b2 := Bound{Min: -1.0, Max: 1.0}
	b3 := Bound{Min: 6.0, Max: 8.0}

	ret := Intersection(b0)
	expected := Bound{Min: 0.0, Max: 5.0}
	if ret != expected {
		t.Errorf("Intersection(b0) should be %v but got %v", expected, ret)
	}

	ret = Intersection(b0, b1)
	expected = Bound{Min: 3.0, Max: 4.0}
	if ret != expected {
		t.Errorf("Intersection(b0, b1) should be %v but got %v", expected, ret)
	}

	ret = Intersection(b0, b2)
	expected = Bound{Min: 0.0, Max: 1.0}
	if ret != expected {
		t.Errorf("Intersection(b0, b2) should be %v but got %v", expected, ret)
	}

	ret = Intersection(b0, b3)
	if ret.IsValid() {
		t.Error("Intersection(b0, b3) should not be valid")
	}

	ret = Intersection(b0, b1, b2)
	if ret.IsValid() {
		t.Error("Intersection(b0, b1, b2) should not be valid")
	}

	ret = Intersection()
	if ret.IsValid() {
		t.Error("Intersection() with zero input length should not be valid")
	}
}

func TestIsValid(t *testing.T) {
	b0 := Bound{Min: 0.0, Max: 5.0}
	if !b0.IsValid() {
		t.Error("b0 is valid")
	}

	b1 := Bound{Min: 5.0, Max: 0.0}
	if b1.IsValid() {
		t.Error("b1 is invalid")
	}

	b2 := Bound{Min: math.NaN(), Max: 5.0}
	if b2.IsValid() {
		t.Error("b2 is invalid")
	}

	b3 := Bound{Min: 5.0, Max: math.NaN()}
	if b3.IsValid() {
		t.Error("b3 is invalid")
	}

	b4 := Bound{Min: math.NaN(), Max: math.NaN()}
	if b4.IsValid() {
		t.Error("b4 is invalid")
	}
}

func TestUnion(t *testing.T) {
	b0 := Bound{Min: 0.0, Max: 5.0}
	b1 := Bound{Min: 3.0, Max: 4.0}
	b2 := Bound{Min: -1.0, Max: 1.0}
	b3 := Bound{Min: 1.0, Max: 2.0}
	b4 := Bound{Min: 1.1, Max: 2.0}

	ret := Union(b0)
	expected := Bound{Min: 0.0, Max: 5.0}
	if ret != expected {
		t.Errorf("Intersection(b0) should be %v but got %v", expected, ret)
	}

	ret = Union(b0, b1)
	expected = Bound{Min: 0.0, Max: 5.0}
	if ret != expected {
		t.Errorf("Intersection(b0, b1) should be %v but got %v", expected, ret)
	}

	ret = Union(b0, b2)
	expected = Bound{Min: -1.0, Max: 5.0}
	if ret != expected {
		t.Errorf("Intersection(b0, b2) should be %v but got %v", expected, ret)
	}

	ret = Union(b2, b3)
	expected = Bound{Min: -1.0, Max: 2.0}
	if ret != expected {
		t.Errorf("Intersection(b0, b2) should be %v but got %v", expected, ret)
	}

	ret = Union(b2, b4)
	if ret.IsValid() {
		t.Errorf("Intersection(b2, b4) should be be valid but got %v", ret)
	}

	ret = Union()
	if ret.IsValid() {
		t.Error("Intersection() with zero input length should not be valid")
	}
}
