// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/internal/ordered"
	"github.com/gonum/graph/simple"
)

func TestIsPath(t *testing.T) {
	dg := simple.NewDirectedGraph(0, math.Inf(1))
	if !IsPathIn(dg, nil) {
		t.Error("IsPath returns false on nil path")
	}
	p := []graph.Node{simple.Node(0)}
	if IsPathIn(dg, p) {
		t.Error("IsPath returns true on nonexistant node")
	}
	dg.AddNode(p[0])
	if !IsPathIn(dg, p) {
		t.Error("IsPath returns false on single-length path with existing node")
	}
	p = append(p, simple.Node(1))
	dg.AddNode(p[1])
	if IsPathIn(dg, p) {
		t.Error("IsPath returns true on bad path of length 2")
	}
	dg.SetEdge(simple.Edge{F: p[0], T: p[1], W: 1})
	if !IsPathIn(dg, p) {
		t.Error("IsPath returns false on correct path of length 2")
	}
	p[0], p[1] = p[1], p[0]
	if IsPathIn(dg, p) {
		t.Error("IsPath erroneously returns true for a reverse path")
	}
	p = []graph.Node{p[1], p[0], simple.Node(2)}
	dg.SetEdge(simple.Edge{F: p[1], T: p[2], W: 1})
	if !IsPathIn(dg, p) {
		t.Error("IsPath does not find a correct path for path > 2 nodes")
	}
	ug := simple.NewUndirectedGraph(0, math.Inf(1))
	ug.SetEdge(simple.Edge{F: p[1], T: p[0], W: 1})
	ug.SetEdge(simple.Edge{F: p[1], T: p[2], W: 1})
	if !IsPathIn(dg, p) {
		t.Error("IsPath does not correctly account for undirected behavior")
	}
}

var pathExistsInUndirectedTests = []struct {
	g        []intset
	from, to int
	want     bool
}{
	{g: batageljZaversnikGraph, from: 0, to: 0, want: true},
	{g: batageljZaversnikGraph, from: 0, to: 1, want: false},
	{g: batageljZaversnikGraph, from: 1, to: 2, want: true},
	{g: batageljZaversnikGraph, from: 2, to: 1, want: true},
	{g: batageljZaversnikGraph, from: 2, to: 12, want: false},
	{g: batageljZaversnikGraph, from: 20, to: 6, want: true},
}

func TestPathExistsInUndirected(t *testing.T) {
	for i, test := range pathExistsInUndirectedTests {
		g := simple.NewUndirectedGraph(0, math.Inf(1))

		for u, e := range test.g {
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				if !g.Has(simple.Node(v)) {
					g.AddNode(simple.Node(v))
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		got := PathExistsIn(g, simple.Node(test.from), simple.Node(test.to))
		if got != test.want {
			t.Errorf("unexpected result for path existance in test %d: got:%t want %t", i, got, test.want)
		}
	}
}

var pathExistsInDirectedTests = []struct {
	g        []intset
	from, to int
	want     bool
}{
	// The graph definition is such that from node IDs are
	// less than to node IDs.
	{g: batageljZaversnikGraph, from: 0, to: 0, want: true},
	{g: batageljZaversnikGraph, from: 0, to: 1, want: false},
	{g: batageljZaversnikGraph, from: 1, to: 2, want: true},
	{g: batageljZaversnikGraph, from: 2, to: 1, want: false},
	{g: batageljZaversnikGraph, from: 2, to: 12, want: false},
	{g: batageljZaversnikGraph, from: 20, to: 6, want: false},
	{g: batageljZaversnikGraph, from: 6, to: 20, want: true},
}

func TestPathExistsInDirected(t *testing.T) {
	for i, test := range pathExistsInDirectedTests {
		g := simple.NewDirectedGraph(0, math.Inf(1))

		for u, e := range test.g {
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				if !g.Has(simple.Node(v)) {
					g.AddNode(simple.Node(v))
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		got := PathExistsIn(g, simple.Node(test.from), simple.Node(test.to))
		if got != test.want {
			t.Errorf("unexpected result for path existance in test %d: got:%t want %t", i, got, test.want)
		}
	}
}

var connectedComponentTests = []struct {
	g    []intset
	want [][]int
}{
	{
		g: batageljZaversnikGraph,
		want: [][]int{
			{0},
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
	},
}

func TestConnectedComponents(t *testing.T) {
	for i, test := range connectedComponentTests {
		g := simple.NewUndirectedGraph(0, math.Inf(1))

		for u, e := range test.g {
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				if !g.Has(simple.Node(v)) {
					g.AddNode(simple.Node(v))
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		cc := ConnectedComponents(g)
		got := make([][]int, len(cc))
		for j, c := range cc {
			ids := make([]int, len(c))
			for k, n := range c {
				ids[k] = n.ID()
			}
			sort.Ints(ids)
			got[j] = ids
		}
		sort.Sort(ordered.BySliceValues(got))
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected connected components for test %d %T:\ngot: %v\nwant:%v", i, g, got, test.want)
		}
	}
}
