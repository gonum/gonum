// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"io"
	"reflect"
	"strings"
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

func TestQueryRepeat(t *testing.T) {
	const filterTestGraph = `
<ex:a> <p:1> <ex:b> .
<ex:b> <p:1> <ex:c> .
<ex:c> <p:1> <ex:d> .
<ex:d> <p:1> <ex:e> .
<ex:e> <p:1> <ex:f> .
<ex:a> <p:2> <ex:_b> .
<ex:b> <p:2> <ex:_c> .
<ex:c> <p:2> <ex:_d> .
<ex:d> <p:2> <ex:_e> .
<ex:e> <p:2> <ex:_f> .
`

	want := []string{"<ex:a>", "<ex:b>", "<ex:c>", "<ex:d>", "<ex:e>", "<ex:f>"}

	g, err := graphFromStatements(filterTestGraph)
	if err != nil {
		t.Fatalf("unexpected error constructing graph: %v", err)
	}
	start, ok := g.TermFor("<ex:a>")
	if !ok {
		t.Fatal("could not get start term")
	}
	for _, limit := range []int{0, 1, 2, 5, 100} {
		got := []string{}
		var i int
		result := g.Query(start).Repeat(func(q Query) (Query, bool) {
			if i >= limit {
				return q, false
			}
			i++
			q = q.Out(func(s *Statement) bool {
				ok := s.Predicate.Value == "<p:1>"
				if ok {
					got = append(got, s.Object.Value)
				}
				return ok
			})
			return q, true
		}).Unique().Result()

		n := limit
		if n >= len(want) {
			n = len(want) - 1
		}
		if !reflect.DeepEqual(got, want[1:n+1]) {
			t.Errorf("unexpected capture for limit=%d: got:%v want:%v",
				limit, got, want[:n])
		}

		switch {
		case limit < len(want):
			if len(result) == 0 || result[0].Value != want[i] {
				t.Errorf("unexpected result for limit=%d: got:%v want:%v",
					limit, result[0], want[i])
			}
		default:
			if len(result) != 0 {
				t.Errorf("unexpected result for limit=%d: got: %v want:none",
					limit, result[0])
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

// filterTestGraph is used to test Has*Out and Has*In. It has a symmetry
// that means that the in an out tests have the same form, just with opposite
// directions.
const filterTestGraph = `
<ex:a> <p:1> <ex:d> .
<ex:a> <p:2> <ex:f> .
<ex:b> <p:2> <ex:d> .
<ex:c> <p:2> <ex:d> .
<ex:a> <o:n> <ex:d> .
# symmetry line.
<ex:e> <p:1> <ex:a> .
<ex:g> <p:2> <ex:a> .
<ex:e> <p:2> <ex:b> .
<ex:e> <p:2> <ex:c> .
<ex:e> <o:n> <ex:a> .
`

var hasOutTests = []struct {
	name    string
	in      []Term
	fn      func(*Statement) bool
	cons    func(q Query) Query
	wantAll []Term
	wantAny []Term
}{
	{
		name: "all",
		in:   []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
		fn:   func(s *Statement) bool { return true },
		cons: func(q Query) Query {
			cond := func(s *Statement) bool { return true }
			return q.Out(cond).In(cond).Unique()
		},
		wantAll: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
		wantAny: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
	},
	{
		name: "none",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return false },
		cons: func(q Query) Query {
			cond := func(s *Statement) bool { return false }
			return q.Out(cond).In(cond).Unique()
		},
		wantAll: nil,
		wantAny: nil,
	},
	{
		name: ". <p:1> .",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return s.Predicate.Value == "<p:1>" },
		cons: func(q Query) Query {
			cond1 := func(s *Statement) bool { return s.Predicate.Value == "<p:1>" }
			cond2 := func(s *Statement) bool { return s.Predicate.Value != "<p:1>" }
			return q.Out(cond1).In(cond1).Not(q.Out(cond2).In(cond2)).Unique()
		},
		wantAll: nil,
		wantAny: []Term{{Value: "<ex:a>"}},
	},
	{
		name: "!(. <p:1> .)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return s.Predicate.Value != "<p:1>" },
		cons: func(q Query) Query {
			cond1 := func(s *Statement) bool { return s.Predicate.Value != "<p:1>" }
			cond2 := func(s *Statement) bool { return s.Predicate.Value == "<p:1>" }
			return q.Out(cond1).In(cond1).Not(q.Out(cond2).In(cond2)).Unique()
		},
		wantAll: []Term{{Value: "<ex:b>"}, {Value: "<ex:c>"}},
		wantAny: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
	},
	{
		name: "!(. <p:2> .)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return s.Predicate.Value != "<p:2>" },
		cons: func(q Query) Query {
			cond1 := func(s *Statement) bool { return s.Predicate.Value != "<p:2>" }
			cond2 := func(s *Statement) bool { return s.Predicate.Value == "<p:2>" }
			return q.Out(cond1).In(cond1).Not(q.Out(cond2).In(cond2)).Unique()
		},
		wantAll: nil,
		wantAny: []Term{{Value: "<ex:a>"}},
	},
	{
		name: "!(. <p:2>  <ex:f>)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn: func(s *Statement) bool {
			return s.Predicate.Value != "<p:2>" || (s.Predicate.Value == "<p:2>" && s.Object.Value != "<ex:f>")
		},
		cons: func(q Query) Query {
			cond := func(s *Statement) bool {
				return s.Predicate.Value == "<p:2>" && s.Object.Value != "<ex:f>"
			}
			return q.Out(cond).In(cond).Unique()
		},
		wantAll: []Term{{Value: "<ex:b>"}, {Value: "<ex:c>"}},
		wantAny: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
	},
	{
		name: "!(. <p:2>  !<ex:f>)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn: func(s *Statement) bool {
			return s.Predicate.Value != "<p:2>" || (s.Predicate.Value == "<p:2>" && s.Object.Value == "<ex:f>")
		},
		cons: func(q Query) Query {
			cond := func(s *Statement) bool {
				return s.Predicate.Value == "<p:2>" && s.Object.Value == "<ex:f>"
			}
			return q.Out(cond).In(cond).Unique()
		},
		wantAll: []Term{{Value: "<ex:a>"}},
		wantAny: []Term{{Value: "<ex:a>"}},
	},
}

func TestQueryHasAllOut(t *testing.T) {
	g, err := graphFromStatements(filterTestGraph)
	if err != nil {
		t.Fatalf("unexpected error constructing graph: %v", err)
	}
	for _, test := range hasOutTests {
		ids := make(map[string]int64)
		for i, v := range test.in {
			term, ok := g.TermFor(v.Value)
			if !ok {
				t.Fatalf("unexpected error constructing graph: could not get UID for term: %v", v.Value)
			}
			test.in[i].UID = term.UID
			ids[term.Value] = term.UID
		}
		for i, v := range test.wantAll {
			test.wantAll[i].UID = ids[v.Value]
		}

		a := Query{g: g, terms: test.in}

		got := a.HasAllOut(test.fn).Result()
		sortByID(got)
		sortByID(test.wantAll)

		if !reflect.DeepEqual(got, test.wantAll) {
			t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
				test.name, got, test.wantAll)
		}

		cons := test.cons(a).Result()
		sortByID(cons)
		if !reflect.DeepEqual(got, cons) {
			t.Errorf("unexpected construction result for test %q:\ngot: %v\nwant:%v",
				test.name, got, cons)
		}
	}
}

func TestQueryHasAnyOut(t *testing.T) {
	g, err := graphFromStatements(filterTestGraph)
	if err != nil {
		t.Fatalf("unexpected error constructing graph: %v", err)
	}
	for _, test := range hasOutTests {
		ids := make(map[string]int64)
		for i, v := range test.in {
			term, ok := g.TermFor(v.Value)
			if !ok {
				t.Fatalf("unexpected error constructing graph: could not get UID for term: %v", v.Value)
			}
			test.in[i].UID = term.UID
			ids[term.Value] = term.UID
		}
		for i, v := range test.wantAny {
			test.wantAny[i].UID = ids[v.Value]
		}

		a := Query{g: g, terms: test.in}

		got := a.HasAnyOut(test.fn).Result()
		sortByID(got)
		sortByID(test.wantAny)

		if !reflect.DeepEqual(got, test.wantAny) {
			t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
				test.name, got, test.wantAny)
		}
	}
}

var hasInTests = []struct {
	name    string
	in      []Term
	fn      func(*Statement) bool
	cons    func(q Query) Query
	wantAll []Term
	wantAny []Term
}{
	{
		name: "all",
		in:   []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
		fn:   func(s *Statement) bool { return true },
		cons: func(q Query) Query {
			cond := func(s *Statement) bool { return true }
			return q.In(cond).Out(cond).Unique()
		},
		wantAll: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
		wantAny: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
	},
	{
		name: "none",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return false },
		cons: func(q Query) Query {
			cond := func(s *Statement) bool { return false }
			return q.In(cond).Out(cond).Unique()
		},
		wantAll: nil,
		wantAny: nil,
	},
	{
		name: ". <p:1> .",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return s.Predicate.Value == "<p:1>" },
		cons: func(q Query) Query {
			cond1 := func(s *Statement) bool { return s.Predicate.Value == "<p:1>" }
			cond2 := func(s *Statement) bool { return s.Predicate.Value != "<p:1>" }
			return q.In(cond1).Out(cond1).Not(q.In(cond2).Out(cond2)).Unique()
		},
		wantAll: nil,
		wantAny: []Term{{Value: "<ex:a>"}},
	},
	{
		name: "!(. <p:1> .)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return s.Predicate.Value != "<p:1>" },
		cons: func(q Query) Query {
			cond1 := func(s *Statement) bool { return s.Predicate.Value != "<p:1>" }
			cond2 := func(s *Statement) bool { return s.Predicate.Value == "<p:1>" }
			return q.In(cond1).Out(cond1).Not(q.In(cond2).Out(cond2)).Unique()
		},
		wantAll: []Term{{Value: "<ex:b>"}, {Value: "<ex:c>"}},
		wantAny: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
	},
	{
		name: "!(. <p:2> .)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn:   func(s *Statement) bool { return s.Predicate.Value != "<p:2>" },
		cons: func(q Query) Query {
			cond1 := func(s *Statement) bool { return s.Predicate.Value != "<p:2>" }
			cond2 := func(s *Statement) bool { return s.Predicate.Value == "<p:2>" }
			return q.In(cond1).Out(cond1).Not(q.In(cond2).Out(cond2)).Unique()
		},
		wantAll: nil,
		wantAny: []Term{{Value: "<ex:a>"}},
	},
	{
		name: "!(<ex:f> <p:2>  .)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn: func(s *Statement) bool {
			return s.Predicate.Value != "<p:2>" || (s.Predicate.Value == "<p:2>" && s.Subject.Value != "<ex:g>")
		},
		cons: func(q Query) Query {
			cond := func(s *Statement) bool {
				return s.Predicate.Value == "<p:2>" && s.Subject.Value != "<ex:g>"
			}
			return q.In(cond).Out(cond).Unique()
		},
		wantAll: []Term{{Value: "<ex:b>"}, {Value: "<ex:c>"}},
		wantAny: []Term{{Value: "<ex:a>"}, {Value: "<ex:b>"}, {Value: "<ex:c>"}},
	},
	{
		name: "!(!<ex:f> <p:2>  .)",
		in:   []Term{{Value: "<ex:a>", UID: 1}, {Value: "<ex:b>", UID: 2}, {Value: "<ex:c>", UID: 3}},
		fn: func(s *Statement) bool {
			return s.Predicate.Value != "<p:2>" || (s.Predicate.Value == "<p:2>" && s.Subject.Value == "<ex:g>")
		},
		cons: func(q Query) Query {
			cond := func(s *Statement) bool {
				return s.Predicate.Value == "<p:2>" && s.Subject.Value == "<ex:g>"
			}
			return q.In(cond).Out(cond).Unique()
		},
		wantAll: []Term{{Value: "<ex:a>"}},
		wantAny: []Term{{Value: "<ex:a>"}},
	},
}

func TestQueryHasAllIn(t *testing.T) {
	g, err := graphFromStatements(filterTestGraph)
	if err != nil {
		t.Fatalf("unexpected error constructing graph: %v", err)
	}
	for _, test := range hasInTests {
		ids := make(map[string]int64)
		for i, v := range test.in {
			term, ok := g.TermFor(v.Value)
			if !ok {
				t.Fatalf("unexpected error constructing graph: could not get UID for term: %v", v.Value)
			}
			test.in[i].UID = term.UID
			ids[term.Value] = term.UID
		}
		for i, v := range test.wantAll {
			test.wantAll[i].UID = ids[v.Value]
		}

		a := Query{g: g, terms: test.in}

		got := a.HasAllIn(test.fn).Result()
		sortByID(got)
		sortByID(test.wantAll)

		if !reflect.DeepEqual(got, test.wantAll) {
			t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
				test.name, got, test.wantAll)
		}

		cons := test.cons(a).Result()
		sortByID(cons)
		if !reflect.DeepEqual(got, cons) {
			t.Errorf("unexpected construction result for test %q:\ngot: %v\nwant:%v",
				test.name, got, cons)
		}
	}
}

func TestQueryHasAnyIn(t *testing.T) {
	g, err := graphFromStatements(filterTestGraph)
	if err != nil {
		t.Fatalf("unexpected error constructing graph: %v", err)
	}
	for _, test := range hasInTests {
		ids := make(map[string]int64)
		for i, v := range test.in {
			term, ok := g.TermFor(v.Value)
			if !ok {
				t.Fatalf("unexpected error constructing graph: could not get UID for term: %v", v.Value)
			}
			test.in[i].UID = term.UID
			ids[term.Value] = term.UID
		}
		for i, v := range test.wantAny {
			test.wantAny[i].UID = ids[v.Value]
		}

		a := Query{g: g, terms: test.in}

		got := a.HasAnyIn(test.fn).Result()
		sortByID(got)
		sortByID(test.wantAny)

		if !reflect.DeepEqual(got, test.wantAny) {
			t.Errorf("unexpected result for test %q:\ngot: %v\nwant:%v",
				test.name, got, test.wantAny)
		}
	}
}

func graphFromStatements(s string) (*Graph, error) {
	g := NewGraph()
	dec := NewDecoder(strings.NewReader(s))
	for {
		s, err := dec.Unmarshal()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		g.AddStatement(s)
	}
	return g, nil
}

func permutedTerms(t []Term, src rand.Source) []Term {
	rnd := rand.New(src)
	p := make([]Term, len(t))
	for i, j := range rnd.Perm(len(t)) {
		p[i] = t[j]
	}
	return p
}
