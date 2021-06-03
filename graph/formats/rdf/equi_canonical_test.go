// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestRemoveRedundantNodes(t *testing.T) {
	tests := []struct {
		name       string
		statements string
		want       string
	}{
		{
			name: "Example 5.1",
			statements: `
<ex:Chile> <ex:cabinet> _:b1 .
<ex:Chile> <ex:cabinet> _:b2 .
<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:cabinet> _:b4 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b2 <ex:members> "23" .
_:b3 <ex:members> "23" .
`,
			want: `<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
		},
		{
			name: "Example 5.2",
			statements: `
_:a <ex:p> _:b .
_:a <ex:p> _:d .
_:c <ex:p> _:b .
_:c <ex:p> _:f .
_:e <ex:p> _:b .
_:e <ex:p> _:d .
_:e <ex:p> _:f .
_:e <ex:p> _:h .
_:g <ex:p> _:d .
_:g <ex:p> _:h .
_:i <ex:p> _:f .
_:i <ex:p> _:h .
`,
			want: `_:e <ex:p> _:b .
`,
		},
	}

	for _, test := range tests {
		g := parseStatements(strings.NewReader(test.statements))
		gWant := parseStatements(strings.NewReader(test.want))

		g = removeRedundantBnodes(g)

		got := canonicalStatements(g)
		want := canonicalStatements(gWant)
		if got != want {
			got = formatStatements(g)
			t.Errorf("unexpected result for %s:\ngot: \n%s\nwant:\n%s",
				test.name, got, test.want)
		}

	}
}
func TestFindCandidates(t *testing.T) {
	tests := []struct {
		name         string
		statements   string
		want         string
		wantFixed    []map[string]bool
		wantAllFixed bool
		wantCands    []map[string]map[string]bool
	}{
		{
			name: "Example 5.1",
			statements: `
<ex:Chile> <ex:cabinet> _:b1 .
<ex:Chile> <ex:cabinet> _:b2 .
<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:cabinet> _:b4 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b2 <ex:members> "23" .
_:b3 <ex:members> "23" .
`,
			want: `<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
			// Hence, in this particular case, we have managed to fix all blank
			// nodes, and the graph is thus lean, and we need to go no further.
			// In other cases we will look at later, however, some blank nodes
			// may maintain multiple candidates.
			//
			// Note that there are two valid labellings of the graph since _:b2
			// and _:b3 are not distinguishable.
			wantFixed: []map[string]bool{
				{
					"_:a1": true, "_:a2": true, "_:a3": true, "_:a4": true,
					"_:b2": true, "_:c": true,
				},
				{
					"_:a1": true, "_:a2": true, "_:a3": true, "_:a4": true,
					"_:b3": true, "_:c": true,
				},
			},
			wantAllFixed: true,
			wantCands: []map[string]map[string]bool{
				{
					"_:a1": {"_:a1": true},
					"_:a2": {"_:a2": true},
					"_:a3": {"_:a3": true},
					"_:a4": {"_:a4": true},
					"_:b2": {"_:b2": true},
					"_:c":  {"_:c": true},
				},
				{
					"_:a1": {"_:a1": true},
					"_:a2": {"_:a2": true},
					"_:a3": {"_:a3": true},
					"_:a4": {"_:a4": true},
					"_:b3": {"_:b3": true},
					"_:c":  {"_:c": true},
				},
			},
		},
		{
			name: "Example 5.6", // This is 5.1, but simplified.
			statements: `
<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
			want: `<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
			// Hence, in this particular case, we have managed to fix all blank
			// nodes, and the graph is thus lean, and we need to go no further.
			// In other cases we will look at later, however, some blank nodes
			// may maintain multiple candidates.
			wantFixed: []map[string]bool{{
				"_:a1": true, "_:a2": true, "_:a3": true, "_:a4": true,
				"_:b3": true, "_:c": true,
			}},
			wantAllFixed: true,
			wantCands: []map[string]map[string]bool{{
				"_:a1": {"_:a1": true},
				"_:a2": {"_:a2": true},
				"_:a3": {"_:a3": true},
				"_:a4": {"_:a4": true},
				"_:b3": {"_:b3": true},
				"_:c":  {"_:c": true},
			}},
		},
		{
			name: "Example 5.9",
			statements: `
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
		`,
			want: `<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
`,
			wantFixed:    []map[string]bool{nil},
			wantAllFixed: true,
			wantCands: []map[string]map[string]bool{{
				"_:a1": {"_:a1": true, "_:a2": true, "_:a3": true},
				"_:a2": {"_:a2": true, "_:a3": true},
				"_:a3": {"_:a2": true, "_:a3": true},
				"_:a4": {"_:a2": true, "_:a3": true, "_:a4": true},
			}},
		},
		{
			name: "Example 5.10",
			statements: `
_:a <ex:p> _:b .
_:a <ex:p> _:d .
_:b <ex:q> _:e .
_:c <ex:p> _:b .
_:c <ex:p> _:f .
_:d <ex:q> _:e .
_:f <ex:q> _:e .
_:g <ex:p> _:d .
_:g <ex:p> _:h .
_:h <ex:q> _:e .
_:i <ex:p> _:f .
_:i <ex:p> _:h .
`,
			want: `_:a <ex:p> _:b .
_:a <ex:p> _:d .
_:b <ex:q> _:e .
_:c <ex:p> _:b .
_:c <ex:p> _:f .
_:d <ex:q> _:e .
_:f <ex:q> _:e .
_:g <ex:p> _:d .
_:g <ex:p> _:h .
_:h <ex:q> _:e .
_:i <ex:p> _:f .
_:i <ex:p> _:h .
`,
			wantFixed:    []map[string]bool{{"_:e": true}},
			wantAllFixed: false,
			wantCands: []map[string]map[string]bool{{
				"_:a": {"_:a": true, "_:c": true, "_:g": true, "_:i": true},
				"_:b": {"_:b": true, "_:d": true, "_:f": true, "_:h": true},
				"_:c": {"_:a": true, "_:c": true, "_:g": true, "_:i": true},
				"_:d": {"_:b": true, "_:d": true, "_:f": true, "_:h": true},
				"_:e": {"_:e": true},
				"_:f": {"_:b": true, "_:d": true, "_:f": true, "_:h": true},
				"_:g": {"_:a": true, "_:c": true, "_:g": true, "_:i": true},
				"_:h": {"_:b": true, "_:d": true, "_:f": true, "_:h": true},
				"_:i": {"_:a": true, "_:c": true, "_:g": true, "_:i": true},
			}},
		},
	}

	for _, test := range tests[:1] {
		g := parseStatements(strings.NewReader(test.statements))
		gWant := parseStatements(strings.NewReader(test.want))

		g, fixed, cands, allFixed := findCandidates(g)

		got := canonicalStatements(g)
		want := canonicalStatements(gWant)
		if got != want {
			got = formatStatements(g)
			t.Errorf("unexpected result for %s:\ngot: \n%s\nwant:\n%s",
				test.name, got, test.want)
		}

		matchedFixed := false
		for _, wantFixed := range test.wantFixed {
			if reflect.DeepEqual(fixed, wantFixed) {
				matchedFixed = true
				break
			}
		}
		if !matchedFixed {
			t.Errorf("unexpected fixed result for %s:\ngot: \n%v\nwant:\n%v",
				test.name, fixed, test.wantFixed)
		}

		if allFixed != test.wantAllFixed {
			t.Errorf("unexpected all-fixed result for %s:\ngot:%t\nwant:%t",
				test.name, allFixed, test.wantAllFixed)
		}

		matchedCands := false
		for _, wantCands := range test.wantCands {
			if reflect.DeepEqual(cands, wantCands) {
				matchedCands = true
				break
			}
		}
		if !matchedCands {
			t.Errorf("unexpected candidates result for %s:\ngot: \n%v\nwant:\n%v",
				test.name, cands, test.wantCands)
		}
	}
}

func TestLean(t *testing.T) {
	var tests = []struct {
		name       string
		statements string
		want       string
		wantErr    error
	}{
		{
			name: "Example 5.1",
			statements: `
<ex:Chile> <ex:cabinet> _:b1 .
<ex:Chile> <ex:cabinet> _:b2 .
<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:cabinet> _:b4 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b2 <ex:members> "23" .
_:b3 <ex:members> "23" .
`,
			want: `<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
		},
		{
			name: "Example 5.6",
			statements: `
<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
			want: `<ex:Chile> <ex:cabinet> _:b3 .
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
<ex:MBachelet> <ex:spouse> _:c .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a2 <ex:president> <ex:MBachelet> .
_:a3 <ex:next> _:a4 .
_:a4 <ex:president> <ex:MBachelet> .
_:b3 <ex:members> "23" .
`,
		},
		{
			name: "Example 5.9",
			statements: `
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
		`,
			want: `<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
`,
		},
		{
			name: "Example 5.10",
			statements: `
_:a <ex:p> _:b .
_:a <ex:p> _:d .
_:b <ex:q> _:e .
_:c <ex:p> _:b .
_:c <ex:p> _:f .
_:d <ex:q> _:e .
_:f <ex:q> _:e .
_:g <ex:p> _:d .
_:g <ex:p> _:h .
_:h <ex:q> _:e .
_:i <ex:p> _:f .
_:i <ex:p> _:h .
`,
			want: `_:a <ex:p> _:b .
_:b <ex:q> _:e .
`,
		},
		{
			name: "Example 5.10 halved",
			statements: `
_:a <ex:p> _:b .
_:a <ex:p> _:d .
_:b <ex:q> _:e .
_:c <ex:p> _:b .
_:c <ex:p> _:f .
_:d <ex:q> _:e .
_:f <ex:q> _:e .
`,
			want: `_:a <ex:p> _:b .
_:b <ex:q> _:e .
`,
		},
		{
			name: "Example 5.10 quartered",
			statements: `
_:a <ex:p> _:b .
_:a <ex:p> _:d .
_:b <ex:q> _:e .
_:d <ex:q> _:e .
`,
			want: `_:a <ex:p> _:b .
_:b <ex:q> _:e .
`,
		},
	}

	for _, test := range tests {
		g := parseStatements(strings.NewReader(test.statements))
		gWant := parseStatements(strings.NewReader(test.want))

		lean, err := Lean(g)
		if err != test.wantErr {
			t.Errorf("unexpected error for %v: got:%v want:%v",
				test.name, err, test.wantErr)
		}

		got := canonicalStatements(lean)
		want := canonicalStatements(gWant)

		if got != want {
			got = formatStatements(g)
			t.Errorf("unexpected result for %s:\ngot: \n%s\nwant:\n%s",
				test.name, got, test.want)
		}
	}
}

func TestJoin(t *testing.T) {
	var tests = []struct {
		name       string
		q          string
		statements string
		cands      map[string]map[string]bool
		mu         map[string]string
		want       []map[string]string
	}{
		{
			name: "Indentity",
			q:    `_:a <ex:p> _:b .`,
			statements: `
_:a <ex:p> _:b .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true},
				"_:b": {"_:b": true},
			},
			mu: nil,
			want: []map[string]string{
				{"_:a": "_:a", "_:b": "_:b"},
			},
		},
		{
			name: "Cross identity",
			q:    `_:a <ex:p> _:b .`,
			statements: `
_:a <ex:p> _:b .
_:b <ex:p> _:a .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: nil,
			want: []map[string]string{
				{"_:a": "_:a", "_:b": "_:b"},
				{"_:a": "_:b", "_:b": "_:a"},
			},
		},
		{
			name: "Cross identity with restriction",
			q:    `_:a <ex:p> _:b .`,
			statements: `
_:a <ex:p> _:b .
_:b <ex:p> _:a .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: map[string]string{"_:a": "_:a"},
			want: []map[string]string{
				{"_:a": "_:a", "_:b": "_:b"},
			},
		},
		{
			name: "Cross identity with complete restriction",
			q:    `_:a <ex:p> _:b .`,
			statements: `
_:a <ex:p> _:b .
_:b <ex:p> _:a .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu:   map[string]string{"_:a": "_:a", "_:b": "_:a"},
			want: nil,
		},
		{
			name: "Cross identity with extension",
			q:    `_:a <ex:p> _:b .`,
			statements: `
_:a <ex:p> _:b .
_:b <ex:p> _:a .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: map[string]string{"_:c": "_:a"},
			want: []map[string]string{
				{"_:a": "_:a", "_:b": "_:b", "_:c": "_:a"},
				{"_:a": "_:b", "_:b": "_:a", "_:c": "_:a"},
			},
		},

		{
			name: "Loop",
			q:    `_:a <ex:p> _:a .`,
			statements: `
_:a <ex:p> _:a .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true},
			},
			mu: nil,
			want: []map[string]string{
				{"_:a": "_:a"},
			},
		},
		{
			name: "Cross identity loop",
			q:    `_:a <ex:p> _:a .`,
			statements: `
_:a <ex:p> _:a .
_:b <ex:p> _:b .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: nil,
			want: []map[string]string{
				{"_:a": "_:a"},
				{"_:a": "_:b"},
			},
		},
		{
			name: "Cross identity loop with restriction",
			q:    `_:a <ex:p> _:a .`,
			statements: `
_:a <ex:p> _:a .
_:b <ex:p> _:b .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: map[string]string{"_:a": "_:a"},
			want: []map[string]string{
				{"_:a": "_:a"},
			},
		},
		{
			name: "Cross identity loop with complete restriction",
			q:    `_:a <ex:p> _:a .`,
			statements: `
_:a <ex:p> _:a .
_:b <ex:p> _:b .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: map[string]string{"_:a": "_:b", "_:b": "_:a"},
			want: []map[string]string{
				{"_:a": "_:b", "_:b": "_:a"},
			},
		},
		{
			name: "Cross identity loop with extension",
			q:    `_:a <ex:p> _:a .`,
			statements: `
_:a <ex:p> _:a .
_:b <ex:p> _:b .
`,
			cands: map[string]map[string]bool{
				"_:a": {"_:a": true, "_:b": true},
				"_:b": {"_:b": true, "_:a": true},
			},
			mu: map[string]string{"_:c": "_:a"},
			want: []map[string]string{
				{"_:a": "_:a", "_:c": "_:a"},
				{"_:a": "_:b", "_:c": "_:a"},
			},
		},

		{
			name: "Example 5.9 step 1",
			q:    `_:a1 <ex:next> _:a2 .`,
			statements: `
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
`,
			cands: map[string]map[string]bool{
				"_:a1": {"_:a1": true, "_:a2": true, "_:a3": true},
				"_:a2": {"_:a2": true, "_:a3": true},
				"_:a3": {"_:a2": true, "_:a3": true},
				"_:a4": {"_:a2": true, "_:a3": true, "_:a4": true},
			},
			mu: map[string]string{},
			want: []map[string]string{
				{"_:a1": "_:a1", "_:a2": "_:a2"},
				{"_:a1": "_:a2", "_:a2": "_:a3"},
			},
		},
		{
			name: "Example 5.9 step 2",
			q:    `_:a2 <ex:next> _:a3 .`,
			statements: `
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
`,
			cands: map[string]map[string]bool{
				"_:a1": {"_:a1": true, "_:a2": true, "_:a3": true},
				"_:a2": {"_:a2": true, "_:a3": true},
				"_:a3": {"_:a2": true, "_:a3": true},
				"_:a4": {"_:a2": true, "_:a3": true, "_:a4": true},
			},
			mu:   map[string]string{"_:a1": "_:a2", "_:a2": "_:a3"},
			want: nil,
		},
		{
			name: "Example 5.9 step 3",
			q:    `_:a2 <ex:next> _:a3 .`,
			statements: `
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
`,
			cands: map[string]map[string]bool{
				"_:a1": {"_:a1": true, "_:a2": true, "_:a3": true},
				"_:a2": {"_:a2": true, "_:a3": true},
				"_:a3": {"_:a2": true, "_:a3": true},
				"_:a4": {"_:a2": true, "_:a3": true, "_:a4": true},
			},
			mu: map[string]string{"_:a1": "_:a1", "_:a2": "_:a2"},
			want: []map[string]string{
				{"_:a1": "_:a1", "_:a2": "_:a2", "_:a3": "_:a3"},
			},
		},
		{
			name: "Example 5.9 step 4",
			q:    `_:a3 <ex:next> _:a4 .`,
			statements: `
<ex:Chile> <ex:presidency> _:a1 .
<ex:Chile> <ex:presidency> _:a2 .
<ex:Chile> <ex:presidency> _:a3 .
<ex:Chile> <ex:presidency> _:a4 .
_:a1 <ex:next> _:a2 .
_:a2 <ex:next> _:a3 .
_:a3 <ex:next> _:a4 .
`,
			cands: map[string]map[string]bool{
				"_:a1": {"_:a1": true, "_:a2": true, "_:a3": true},
				"_:a2": {"_:a2": true, "_:a3": true},
				"_:a3": {"_:a2": true, "_:a3": true},
				"_:a4": {"_:a2": true, "_:a3": true, "_:a4": true},
			},
			mu: map[string]string{"_:a1": "_:a1", "_:a2": "_:a2", "_:a3": "_:a3"},
			want: []map[string]string{
				{"_:a1": "_:a1", "_:a2": "_:a2", "_:a3": "_:a3", "_:a4": "_:a4"},
			},
		},
	}

	for _, test := range tests {
		q := parseStatement(strings.NewReader(test.q))
		g := parseStatements(strings.NewReader(test.statements))

		st := dfs{}
		got := st.join(q, g, test.cands, test.mu)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected result for %s:\ngot:\n%#v\nwant:\n%#v",
				test.name, got, test.want)
		}

		naive := joinNaive(q, g, test.cands, []map[string]string{test.mu})
		if !reflect.DeepEqual(naive, test.want) {
			t.Errorf("unexpected naive result for %s:\ngot:\n%#v\nwant:\n%#v",
				test.name, naive, test.want)
		}

	}
}

// joinNaive is a direct translation of lines 47-51 of algorithm 6 in doi:10.1145/3068333.
func joinNaive(q *Statement, G []*Statement, cands map[string]map[string]bool, M []map[string]string) []map[string]string {
	isLoop := q.Subject.Value == q.Object.Value
	// Line 48: M_q ← {µ | µ(q) ∈ G}
	var M_q []map[string]string
	for _, s := range G {
		// µ(q) ∈ G ↔ (µ(q_s),q_p,µ(q_o)) ∈ G
		if q.Predicate.Value != s.Predicate.Value {
			continue
		}
		// q_s = q_o ↔ µ(q_s) =_µ(q_o)
		if isLoop && s.Subject.Value != s.Object.Value {
			continue
		}

		var µ map[string]string
		if isLoop {
			µ = map[string]string{
				q.Subject.Value: s.Subject.Value,
			}
		} else {
			µ = map[string]string{
				q.Subject.Value: s.Subject.Value,
				q.Object.Value:  s.Object.Value,
			}
		}
		M_q = append(M_q, µ)
	}

	// Line 49: M_q' ← {µ ∈ M_q | for all b ∈ bnodes({q}), µ(b) ∈ cands[b]}
	var M_qPrime []map[string]string
	for _, µ := range M_q {
		if !cands[q.Subject.Value][µ[q.Subject.Value]] {
			continue
		}
		if !cands[q.Object.Value][µ[q.Object.Value]] {
			continue
		}
		M_qPrime = append(M_qPrime, µ)
	}

	// Line 50: M' ← M_q' ⋈ M
	// M₁ ⋈ M₂ = {μ₁ ∪ μ₂ | μ₁ ∈ M₁, μ₂ ∈ M₂ and μ₁, μ₂ are compatible mappings}
	var MPrime []map[string]string
	for _, µ := range M {
	join:
		for _, µ_qPrime := range M_qPrime {
			for b, x_qPrime := range µ_qPrime {
				if x, ok := µ[b]; ok && x != x_qPrime {
					continue join
				}
			}
			// Line 50: μ₁ ∪ μ₂
			for b, x := range µ {
				µ_qPrime[b] = x
			}
			MPrime = append(MPrime, µ_qPrime)
		}
	}
	return MPrime
}

func parseStatement(r io.Reader) *Statement {
	g := parseStatements(r)
	if len(g) != 1 {
		panic(fmt.Sprintf("invalid statement stream length %d != 1", len(g)))
	}
	return g[0]
}

func parseStatements(r io.Reader) []*Statement {
	var g []*Statement
	dec := NewDecoder(r)
	for {
		s, err := dec.Unmarshal()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		g = append(g, s)
	}
	return g
}

func canonicalStatements(g []*Statement) string {
	g, _ = URDNA2015(nil, g)
	return formatStatements(g)
}

func formatStatements(g []*Statement) string {
	var buf strings.Builder
	for _, s := range g {
		fmt.Fprintln(&buf, s)
	}
	return buf.String()
}
