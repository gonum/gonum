package r3

import "testing"

func TestVector_Cross(t *testing.T) {
	a := Vec{X: 2, Y: 3, Z: 4}
	b := Vec{X: 5, Y: 6, Z: 7}
	expected := Vec{X: -3, Y: 6, Z: -3}
	actual := a.Cross(b)
	if expected != actual {
		t.Fatalf("unexpected result from cross product of %+v and %+v: got:%+v want:%+v", a, b, actual, expected)
	}
}
