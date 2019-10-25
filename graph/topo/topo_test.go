// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/simple"
)

func TestIsPath(t *testing.T) {
	dg := simple.NewDirectedGraph()
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
	dg.SetEdge(simple.Edge{F: p[0], T: p[1]})
	if !IsPathIn(dg, p) {
		t.Error("IsPath returns false on correct path of length 2")
	}
	p[0], p[1] = p[1], p[0]
	if IsPathIn(dg, p) {
		t.Error("IsPath erroneously returns true for a reverse path")
	}
	p = []graph.Node{p[1], p[0], simple.Node(2)}
	dg.SetEdge(simple.Edge{F: p[1], T: p[2]})
	if !IsPathIn(dg, p) {
		t.Error("IsPath does not find a correct path for path > 2 nodes")
	}
	ug := simple.NewUndirectedGraph()
	ug.SetEdge(simple.Edge{F: p[1], T: p[0]})
	ug.SetEdge(simple.Edge{F: p[1], T: p[2]})
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
		g := simple.NewUndirectedGraph()

		for u, e := range test.g {
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				if g.Node(int64(v)) == nil {
					g.AddNode(simple.Node(v))
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		got := PathExistsIn(g, simple.Node(test.from), simple.Node(test.to))
		if got != test.want {
			t.Errorf("unexpected result for path existence in test %d: got:%t want %t", i, got, test.want)
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
		g := simple.NewDirectedGraph()

		for u, e := range test.g {
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				if g.Node(int64(v)) == nil {
					g.AddNode(simple.Node(v))
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		got := PathExistsIn(g, simple.Node(test.from), simple.Node(test.to))
		if got != test.want {
			t.Errorf("unexpected result for path existence in test %d: got:%t want %t", i, got, test.want)
		}
	}
}

var connectedComponentTests = []struct {
	g    []intset
	want [][]int64
}{
	{
		g: batageljZaversnikGraph,
		want: [][]int64{
			{0},
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
	},
}

func TestConnectedComponents(t *testing.T) {
	for i, test := range connectedComponentTests {
		g := simple.NewUndirectedGraph()

		for u, e := range test.g {
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				if g.Node(int64(v)) == nil {
					g.AddNode(simple.Node(v))
				}
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		cc := ConnectedComponents(g)
		got := make([][]int64, len(cc))
		for j, c := range cc {
			ids := make([]int64, len(c))
			for k, n := range c {
				ids[k] = n.ID()
			}
			sort.Sort(ordered.Int64s(ids))
			got[j] = ids
		}
		sort.Sort(ordered.BySliceValues(got))
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected connected components for test %d %T:\ngot: %v\nwant:%v", i, g, got, test.want)
		}
	}
}

var equalTests = []struct {
	name string
	a, b graph.Graph
	want bool
}{
	{name: "empty g=g", a: simple.NewUndirectedGraph(), b: simple.NewUndirectedGraph(), want: true},
	{name: "empty dg=dg", a: simple.NewDirectedGraph(), b: simple.NewDirectedGraph(), want: true},
	{name: "empty g=dg", a: simple.NewUndirectedGraph(), b: simple.NewDirectedGraph(), want: true},

	{
		name: "1 g=g", want: true,
		a: addNodes(simple.NewUndirectedGraph(), simple.Node(1)),
		b: addNodes(simple.NewUndirectedGraph(), simple.Node(1)),
	},
	{
		name: "1 dg=dg", want: true,
		a: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
		b: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
	},
	{
		name: "1 g=dg", want: true,
		a: addNodes(simple.NewUndirectedGraph(), simple.Node(1)),
		b: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
	},

	{
		name: "0/1 g≠g", want: false,
		a: simple.NewUndirectedGraph(),
		b: addNodes(simple.NewUndirectedGraph(), simple.Node(1)),
	},
	{
		name: "0/1 dg≠dg", want: false,
		a: simple.NewDirectedGraph(),
		b: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
	},
	{
		name: "0/1 g≠dg", want: false,
		a: simple.NewUndirectedGraph(),
		b: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
	},
	{
		name: "0/1 g≠dg", want: false,
		a: addNodes(simple.NewUndirectedGraph(), simple.Node(1)),
		b: simple.NewDirectedGraph(),
	},

	{
		name: "1 g≠g", want: false,
		a: addNodes(simple.NewUndirectedGraph(), simple.Node(0)),
		b: addNodes(simple.NewUndirectedGraph(), simple.Node(1)),
	},
	{
		name: "1 dg≠dg", want: false,
		a: addNodes(simple.NewDirectedGraph(), simple.Node(0)),
		b: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
	},
	{
		name: "1 g≠dg", want: false,
		a: addNodes(simple.NewUndirectedGraph(), simple.Node(0)),
		b: addNodes(simple.NewDirectedGraph(), simple.Node(1)),
	},

	{
		name: "box g=g", want: true,
		a: setEdges(simple.NewUndirectedGraph(),
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
		),
		b: setEdges(simple.NewUndirectedGraph(),
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
		),
	},
	{
		name: "box dg=dg", want: true,
		a: setEdges(simple.NewDirectedGraph(),
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
		),
		b: setEdges(simple.NewDirectedGraph(),
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
		),
	},
	{
		name: "box reversed dg≠dg", want: false,
		a: setEdges(simple.NewDirectedGraph(),
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
		),
		b: setEdges(simple.NewDirectedGraph(),
			simple.Edge{F: simple.Node(1), T: simple.Node(0)},
			simple.Edge{F: simple.Node(2), T: simple.Node(1)},
			simple.Edge{F: simple.Node(3), T: simple.Node(2)},
			simple.Edge{F: simple.Node(0), T: simple.Node(3)},
		),
	},
	{
		name: "box g=dg", want: true,
		a: setEdges(simple.NewUndirectedGraph(),
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
		),
		b: setEdges(simple.NewDirectedGraph(),
			simple.Edge{F: simple.Node(0), T: simple.Node(1)},
			simple.Edge{F: simple.Node(1), T: simple.Node(0)},
			simple.Edge{F: simple.Node(1), T: simple.Node(2)},
			simple.Edge{F: simple.Node(2), T: simple.Node(3)},
			simple.Edge{F: simple.Node(2), T: simple.Node(1)},
			simple.Edge{F: simple.Node(3), T: simple.Node(2)},
			simple.Edge{F: simple.Node(3), T: simple.Node(0)},
			simple.Edge{F: simple.Node(0), T: simple.Node(3)},
		),
	},
}

func TestEqual(t *testing.T) {
	for _, test := range equalTests {
		if got := Equal(test.a, test.b); got != test.want {
			t.Errorf("unexpected result for %q equality test: got:%t want:%t", test.name, got, test.want)
		}
		if got := Equal(plainGraph{test.a}, plainGraph{test.b}); got != test.want {
			t.Errorf("unexpected result for %q equality test with filtered method set: got:%t want:%t", test.name, got, test.want)
		}
	}
}

type plainGraph struct {
	graph.Graph
}

type builder interface {
	graph.Graph
	graph.Builder
}

func addNodes(dst builder, nodes ...graph.Node) builder {
	for _, n := range nodes {
		dst.AddNode(n)
	}
	return dst
}

func setEdges(dst builder, edges ...graph.Edge) builder {
	for _, e := range edges {
		dst.SetEdge(e)
	}
	return dst
}
