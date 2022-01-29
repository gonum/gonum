// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"reflect"
	"testing"

	"golang.org/x/exp/rand"
)

var andTests = []struct {
	name string
	a, b []Term
	want []Term
}{
	{
		name: "identical",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "identical with excess a",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "identical with excess b",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "b less",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}},
	},
	{
		name: "a less",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:c>", UID: 3}},
	},
}

func TestQueryAnd(t *testing.T) {
	src := rand.NewSource(1)
	for _, test := range andTests {
		for i := 0; i < 10; i++ {
			a := Query{terms: permutedTerms(test.a, src)}
			b := Query{terms: permutedTerms(test.b, src)}

			got := a.And(b).Result()
			sortByID(got)
			sortByID(test.want)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
					test.name, got, test.want)
			}
		}
	}
}

var orTests = []struct {
	name string
	a, b []Term
	want []Term
}{
	{
		name: "identical",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "identical with excess a",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "identical with excess b",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "b less",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "a less",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
}

func TestQueryOr(t *testing.T) {
	src := rand.NewSource(1)
	for _, test := range orTests {
		for i := 0; i < 10; i++ {
			a := Query{terms: permutedTerms(test.a, src)}
			b := Query{terms: permutedTerms(test.b, src)}

			got := a.Or(b).Result()
			sortByID(got)
			sortByID(test.want)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
					test.name, got, test.want)
			}
		}
	}
}

var notTests = []struct {
	name string
	a, b []Term
	want []Term
}{
	{
		name: "identical",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: nil,
	},
	{
		name: "identical with excess a",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: nil,
	},
	{
		name: "identical with excess b",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: nil,
	},
	{
		name: "b less",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}},
		want: []Term{{Value: "<ex:c>", UID: 3}},
	},
	{
		name: "a less",
		a:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:c>", UID: 3}},
		b:    []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: nil,
	},
}

func TestQueryNot(t *testing.T) {
	src := rand.NewSource(1)
	for _, test := range notTests {
		for i := 0; i < 10; i++ {
			a := Query{terms: permutedTerms(test.a, src)}
			b := Query{terms: permutedTerms(test.b, src)}

			got := a.Not(b).Result()
			sortByID(got)
			sortByID(test.want)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
					test.name, got, test.want)
			}
		}
	}
}

var uniqueTests = []struct {
	name string
	in   []Term
	want []Term
}{
	{
		name: "excess a",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "excess b",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
	{
		name: "excess c",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}, {Value: "<ex:c>", UID: 3}},
		want: []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
	},
}

func TestQueryUnique(t *testing.T) {
	src := rand.NewSource(1)
	for _, test := range uniqueTests {
		for i := 0; i < 10; i++ {
			a := Query{terms: permutedTerms(test.in, src)}

			got := a.Unique().Result()
			sortByID(got)
			sortByID(test.want)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
					test.name, got, test.want)
			}
		}
	}
}

func permutedTerms(t []Term, src rand.Source) []Term {
	rnd := rand.New(src)
	p := make([]Term, len(t))
	for i, j := range rnd.Perm(len(t)) {
		p[i] = t[j]
	}
	return p
}
