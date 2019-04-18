// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bound

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var isValidTests = []struct {
	b    Bound
	want bool
}{
	{b: Bound{Min: math.Inf(-1), Max: 5}, want: true},
	{b: Bound{Min: 0, Max: 5}, want: true},
	{b: Bound{Min: 0, Max: math.Inf(1)}, want: true},
	{b: Bound{Min: 5, Max: 0}, want: false},
	{b: Bound{Min: math.NaN(), Max: 5}, want: false},
	{b: Bound{Min: 5, Max: math.NaN()}, want: false},
	{b: Bound{Min: math.NaN(), Max: math.NaN()}, want: false},
}

func TestIsValid(t *testing.T) {
	for _, test := range isValidTests {
		got := test.b.IsValid()
		if got != test.want {
			t.Errorf("unexpected validity of %+v: got:%t want:%t", test.b, got, test.want)
		}
	}
}

var intersectionTests = []struct {
	bounds []Bound
	want   Bound
}{
	{
		bounds: []Bound{{Min: 0, Max: 5}},
		want:   Bound{Min: 0, Max: 5},
	},
	{
		bounds: []Bound{{Min: 0, Max: 5}, {Min: 3, Max: 4}},
		want:   Bound{Min: 3, Max: 4},
	},
	{
		bounds: []Bound{{Min: 0, Max: 5}, {Min: -1, Max: 1}},
		want:   Bound{Min: 0, Max: 1},
	},
	{
		bounds: []Bound{{Min: 0, Max: 5}, {Min: 6, Max: 8}},
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
	{
		bounds: []Bound{{Min: 0, Max: 5}, {Min: 3, Max: 4}, {Min: -1, Max: 1}},
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
	{
		bounds: []Bound{},
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
	{
		bounds: nil,
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
}

func TestIntersection(t *testing.T) {
	for _, test := range intersectionTests {
		got := Intersection(test.bounds...)
		if !same(got, test.want) {
			t.Errorf("unexpected result from Intersection(%#v...): got:%+v want:%+v", test.bounds, got, test.want)
		}
	}
}

var unionTests = []struct {
	bounds []Bound
	want   Bound
}{
	{
		bounds: []Bound{{Min: 0, Max: 5}},
		want:   Bound{Min: 0, Max: 5},
	},
	{
		bounds: []Bound{{Min: 0, Max: 5}, {Min: 3, Max: 4}},
		want:   Bound{Min: 0, Max: 5},
	},
	{
		bounds: []Bound{{Min: 0, Max: 5}, {Min: -1, Max: 1}},
		want:   Bound{Min: -1, Max: 5},
	},
	{
		bounds: []Bound{{Min: -1, Max: 1}, {Min: 1, Max: 2}},
		want:   Bound{Min: -1, Max: 2},
	},
	{
		bounds: []Bound{{Min: -1, Max: 1}, {Min: 1.1, Max: 2}},
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
	{
		bounds: []Bound{{Min: 0, Max: 1}, {Min: 2, Max: 3}, {Min: 0.5, Max: 2.5}},
		want:   Bound{Min: 0, Max: 3},
	},
	{
		bounds: []Bound{},
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
	{
		bounds: nil,
		want:   Bound{Min: math.NaN(), Max: math.NaN()},
	},
}

func TestUnion(t *testing.T) {
	for _, test := range unionTests {
		got := Union(test.bounds...)
		if !same(got, test.want) {
			t.Errorf("unexpected result from Union(%#v...): got:%+v want:%+v", test.bounds, got, test.want)
		}
	}
}

func same(a, b Bound) bool {
	return floats.Same([]float64{a.Min, a.Max}, []float64{b.Min, b.Max})
}
