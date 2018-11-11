// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testgraph provides a set of testing helper functions
// that test gonum graph interface implementations.
package testgraph // import "gonum.org/v1/gonum/graph/testgraph"

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/mat"
)

// BUG(kortschak): Edge equality is tested in part with reflect.DeepEqual and
// direct equality of weight values. This means that edges returned by graphs
// must not contain NaN values. Weights returned by the Weight method are
// compared with NaN-awareness, so they may be NaN when there is no edge
// associated with the Weight call.

// BUG(kortschak): The approach of using a nil return for empty sets of nodes
// and edges used prior to the introduction of the graph.Iterator types does
// not interact well with interfaces. For example, it is not possible to simply
// determine that an iterator is empty by calling it.Len without guarding that
// with a nil check. The validity of nil iterators may change depending on the
// outcome of https://github.com/gonum/gonum/issues/614.
func isValidIterator(graph.Iterator) bool {
	// TODO(kortschak): Remove nil guards in iterator
	// loops and slicer tests if this changes.
	return true
}

// A Builder function returns a graph constructed from the nodes, edges and
// default weights passed in, potentially altering the nodes and edges to
// conform to the requirements of the graph. The graph is returned along with
// the nodes, edges and default weights used to construct the graph.
// The returned edges may be any of graph.Edge, graph.WeightedEdge, graph.Line
// or graph.WeightedLine depending on what the graph requires.
// The client may skip a test case by returning ok=false when the input is not
// a valid graph construction.
type Builder func(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool)

// edgeLister is a graph that can return all its edges.
type edgeLister interface {
	// Edges returns all the edges of a graph.
	Edges() graph.Edges
}

// weightedEdgeLister is a graph that can return all its weighted edges.
type weightedEdgeLister interface {
	// WeightedEdges returns all the weighted edges of a graph.
	WeightedEdges() graph.WeightedEdges
}

// matrixer is a graph that can return an adjacency matrix.
type matrixer interface {
	// Matrix returns the graph's adjacency matrix.
	Matrix() mat.Matrix
}

// ReturnAllNodes tests the constructed graph for the ability to return all
// the nodes it claims it has used in its construction. This is a check of
// the Nodes method of graph.Graph and the iterator that is returned.
func ReturnAllNodes(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, want, _, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		it := g.Nodes()
		if !isValidIterator(it) {
			t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
			continue
		}
		var got []graph.Node
		for it != nil && it.Next() {
			got = append(got, it.Node())
		}

		sort.Sort(ordered.ByID(got))
		sort.Sort(ordered.ByID(want))

		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected nodes result for test %q:\ngot: %v\nwant:%v", test.name, got, want)
		}
	}
}

// ReturnNodeSlice tests the constructed graph for the ability to return all
// the nodes it claims it has used in its construction using the NodeSlicer
// interface. This is a check of the Nodes method of graph.Graph and the
// iterator that is returned.
func ReturnNodeSlice(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, want, _, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		it := g.Nodes()
		if !isValidIterator(it) {
			t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
			continue
		}
		if it == nil {
			continue
		}
		s, ok := it.(graph.NodeSlicer)
		if !ok {
			t.Errorf("invalid type for test %q: %T cannot return node slicer", test.name, g)
			continue
		}
		got := s.NodeSlice()

		sort.Sort(ordered.ByID(got))
		sort.Sort(ordered.ByID(want))

		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected nodes result for test %q:\ngot: %v\nwant:%v", test.name, got, want)
		}
	}
}

// NodeExistence tests the constructed graph for the ability to correctly
// return the existence of nodes within the graph. This is a check of the
// Node method of graph.Graph.
func NodeExistence(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, want, _, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		seen := make(set.Nodes)
		for _, exist := range want {
			seen.Add(exist)
			if g.Node(exist.ID()) == nil {
				t.Errorf("missing node for test %q: %v", test.name, exist)
			}
		}
		for _, ghost := range test.nonexist {
			if g.Node(ghost.ID()) != nil {
				if seen.Has(ghost) {
					// Do not fail nodes that the graph builder says can exist
					// even if the test case input thinks they should not.
					t.Logf("builder has modified non-exist node set: %v is now allowed and present", ghost)
					continue
				}
				t.Errorf("unexpected node for test %q: %v", test.name, ghost)
			}
		}
	}
}

// ReturnAllEdges tests the constructed graph for the ability to return all
// the edges it claims it has used in its construction. This is a check of
// the Edges method of graph.Graph and the iterator that is returned.
// ReturnAllEdges  also checks that the edge end nodes exist within the graph,
// checking the Node method of graph.Graph.
func ReturnAllEdges(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, _, want, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		var got []graph.Edge
		switch eg := g.(type) {
		case edgeLister:
			it := eg.Edges()
			if !isValidIterator(it) {
				t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
				continue
			}
			for it != nil && it.Next() {
				e := it.Edge()
				got = append(got, e)
				if g.Edge(e.From().ID(), e.To().ID()) == nil {
					t.Errorf("missing edge for test %q: %v", test.name, e)
				}
				if g.Node(e.From().ID()) == nil {
					t.Errorf("missing from node for test %q: %v", test.name, e.From().ID())
				}
				if g.Node(e.To().ID()) == nil {
					t.Errorf("missing to node for test %q: %v", test.name, e.To().ID())
				}
			}

		default:
			t.Errorf("invalid type for test %q: %T cannot return edge iterator", test.name, g)
			continue
		}

		checkEdges(t, test.name, g, got, want)
	}
}

// ReturnEdgeSlice tests the constructed graph for the ability to return all
// the edges it claims it has used in its construction using the EdgeSlicer
// interface. This is a check of the Edges method of graph.Graph and the
// iterator that is returned. ReturnEdgeSlice also checks that the edge end
// nodes exist within the graph, checking the Node method of graph.Graph.
func ReturnEdgeSlice(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, _, want, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		var got []graph.Edge
		switch eg := g.(type) {
		case edgeLister:
			it := eg.Edges()
			if !isValidIterator(it) {
				t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
				continue
			}
			if it == nil {
				continue
			}
			s, ok := it.(graph.EdgeSlicer)
			if !ok {
				t.Errorf("invalid type for test %q: %T cannot return edge slicer", test.name, g)
				continue
			}
			got = s.EdgeSlice()
			for _, e := range got {
				if g.Edge(e.From().ID(), e.To().ID()) == nil {
					t.Errorf("missing edge for test %q: %v", test.name, e)
				}
				if g.Node(e.From().ID()) == nil {
					t.Errorf("missing from node for test %q: %v", test.name, e.From().ID())
				}
				if g.Node(e.To().ID()) == nil {
					t.Errorf("missing to node for test %q: %v", test.name, e.To().ID())
				}
			}

		default:
			t.Errorf("invalid type for test %T: cannot return edge iterator", g)
			continue
		}

		checkEdges(t, test.name, g, got, want)
	}
}

// ReturnAllLines tests the constructed graph for the ability to return all
// the edges it claims it has used in its construction and then recover all
// the lines that contribute to those edges. This is a check of the Edges
// method of graph.Graph and the iterator that is returned and the graph.Lines
// implementation of those edges. ReturnAllLines also checks that the edge
// end nodes exist within the graph, checking the Node method of graph.Graph.
//
// The edges used within and returned by the Builder function should be
// graph.Line. The edge parameter passed to b will contain only graph.Line.
func ReturnAllLines(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, _, want, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		var got []graph.Edge
		switch eg := g.(type) {
		case edgeLister:
			it := eg.Edges()
			if !isValidIterator(it) {
				t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
				continue
			}
			for _, e := range graph.EdgesOf(it) {
				if g.Edge(e.From().ID(), e.To().ID()) == nil {
					t.Errorf("missing edge for test %q: %v", test.name, e)
				}

				// FIXME(kortschak): This would not be necessary
				// if graph.WeightedLines (and by symmetry)
				// graph.WeightedEdges also were graph.Lines
				// and graph.Edges.
				switch lit := e.(type) {
				case graph.Lines:
					for lit != nil && lit.Next() {
						got = append(got, lit.Line())
					}
				case graph.WeightedLines:
					for lit != nil && lit.Next() {
						got = append(got, lit.WeightedLine())
					}
				default:
					continue
				}

				if g.Node(e.From().ID()) == nil {
					t.Errorf("missing from node for test %q: %v", test.name, e.From().ID())
				}
				if g.Node(e.To().ID()) == nil {
					t.Errorf("missing to node for test %q: %v", test.name, e.To().ID())
				}
			}

		default:
			t.Errorf("invalid type for test: %T cannot return edge iterator", g)
			continue
		}

		checkEdges(t, test.name, g, got, want)
	}
}

// ReturnAllWeightedEdges tests the constructed graph for the ability to return
// all the edges it claims it has used in its construction. This is a check of
// the Edges method of graph.Graph and the iterator that is returned.
// ReturnAllWeightedEdges also checks that the edge end nodes exist within the
// graph, checking the Node method of graph.Graph.
//
// The edges used within and returned by the Builder function should be
// graph.WeightedEdge. The edge parameter passed to b will contain only
// graph.WeightedEdge.
func ReturnAllWeightedEdges(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, _, want, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		var got []graph.Edge
		switch eg := g.(type) {
		case weightedEdgeLister:
			it := eg.WeightedEdges()
			if !isValidIterator(it) {
				t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
				continue
			}
			for it != nil && it.Next() {
				e := it.WeightedEdge()
				got = append(got, e)
				switch g := g.(type) {
				case graph.Weighted:
					if g.WeightedEdge(e.From().ID(), e.To().ID()) == nil {
						t.Errorf("missing edge for test %q: %v", test.name, e)
					}
				default:
					t.Logf("weighted edge lister is not a weighted graph - are you sure?: %T", g)
					if g.Edge(e.From().ID(), e.To().ID()) == nil {
						t.Errorf("missing edge for test %q: %v", test.name, e)
					}
				}
				if g.Node(e.From().ID()) == nil {
					t.Errorf("missing from node for test %q: %v", test.name, e.From().ID())
				}
				if g.Node(e.To().ID()) == nil {
					t.Errorf("missing to node for test %q: %v", test.name, e.To().ID())
				}
			}

		default:
			t.Errorf("invalid type for test: %T cannot return weighted edge iterator", g)
			continue
		}

		checkEdges(t, test.name, g, got, want)
	}
}

// ReturnWeightedEdgeSlice tests the constructed graph for the ability to
// return all the edges it claims it has used in its construction using the
// WeightedEdgeSlicer interface. This is a check of the Edges method of
// graph.Graph and the iterator that is returned. ReturnWeightedEdgeSlice
// also checks that the edge end nodes exist within the graph, checking
// the Node method of graph.Graph.
//
// The edges used within and returned by the Builder function should be
// graph.WeightedEdge. The edge parameter passed to b will contain only
// graph.WeightedEdge.
func ReturnWeightedEdgeSlice(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, _, want, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		var got []graph.Edge
		switch eg := g.(type) {
		case weightedEdgeLister:
			it := eg.WeightedEdges()
			if !isValidIterator(it) {
				t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
				continue
			}
			s, ok := it.(graph.WeightedEdgeSlicer)
			if !ok {
				t.Errorf("invalid type for test %T: cannot return weighted edge slice", g)
				continue
			}
			for _, e := range s.WeightedEdgeSlice() {
				got = append(got, e)
				if g.Edge(e.From().ID(), e.To().ID()) == nil {
					t.Errorf("missing edge for test %q: %v", test.name, e)
				}
				if g.Node(e.From().ID()) == nil {
					t.Errorf("missing from node for test %q: %v", test.name, e.From().ID())
				}
				if g.Node(e.To().ID()) == nil {
					t.Errorf("missing to node for test %q: %v", test.name, e.To().ID())
				}
			}

		default:
			t.Errorf("invalid type for test: %T cannot return weighted edge iterator", g)
			continue
		}

		checkEdges(t, test.name, g, got, want)
	}
}

// ReturnAllWeightedLines tests the constructed graph for the ability to return
// all the edges it claims it has used in its construction and then recover all
// the lines that contribute to those edges. This is a check of the Edges
// method of graph.Graph and the iterator that is returned and the graph.Lines
// implementation of those edges. ReturnAllWeightedLines also checks that the
// edge end nodes exist within the graph, checking the Node method of
// graph.Graph.
//
// The edges used within and returned by the Builder function should be
// graph.WeightedLine. The edge parameter passed to b will contain only
// graph.WeightedLine.
func ReturnAllWeightedLines(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, _, want, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		var got []graph.Edge
		switch eg := g.(type) {
		case weightedEdgeLister:
			it := eg.WeightedEdges()
			if !isValidIterator(it) {
				t.Errorf("invalid iterator for test %q: got:%#v", test.name, it)
				continue
			}
			for _, e := range graph.WeightedEdgesOf(it) {
				if g.Edge(e.From().ID(), e.To().ID()) == nil {
					t.Errorf("missing edge for test %q: %v", test.name, e)
				}

				// FIXME(kortschak): This would not be necessary
				// if graph.WeightedLines (and by symmetry)
				// graph.WeightedEdges also were graph.Lines
				// and graph.Edges.
				switch lit := e.(type) {
				case graph.Lines:
					for lit != nil && lit.Next() {
						got = append(got, lit.Line())
					}
				case graph.WeightedLines:
					for lit != nil && lit.Next() {
						got = append(got, lit.WeightedLine())
					}
				default:
					continue
				}

				if g.Node(e.From().ID()) == nil {
					t.Errorf("missing from node for test %q: %v", test.name, e.From().ID())
				}
				if g.Node(e.To().ID()) == nil {
					t.Errorf("missing to node for test %q: %v", test.name, e.To().ID())
				}
			}

		default:
			t.Errorf("invalid type for test: %T cannot return edge iterator", g)
			continue
		}

		checkEdges(t, test.name, g, got, want)
	}
}

// checkEdges compares got and want for the given graph type.
func checkEdges(t *testing.T, name string, g graph.Graph, got, want []graph.Edge) {
	t.Helper()
	switch g.(type) {
	case graph.Undirected:
		sort.Sort(lexicalUndirectedEdges(got))
		sort.Sort(lexicalUndirectedEdges(want))
		if !undirectedEdgeSetEqual(got, want) {
			t.Errorf("unexpected edges result for test %q:\ngot: %#v\nwant:%#v", name, got, want)
		}
	default:
		sort.Sort(lexicalEdges(got))
		sort.Sort(lexicalEdges(want))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected edges result for test %q:\ngot: %#v\nwant:%#v", name, got, want)
		}
	}
}

// EdgeExistence tests the constructed graph for the ability to correctly
// return the existence of edges within the graph. This is a check of the
// Edge method of graph.Graph, the EdgeBetween method of graph.Undirected
// and the EdgeFromTo method of graph.Directed. EdgeExistence also checks
// that the nodes and traversed edges exist within the graph, checking the
// Node, Edge, EdgeBetween and HasEdgeBetween methods of graph.Graph, the
// EdgeBetween method of graph.Undirected and the HasEdgeFromTo method of
// graph.Directed.
func EdgeExistence(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, nodes, edges, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		want := make(map[edge]bool)
		for _, e := range edges {
			want[edge{f: e.From().ID(), t: e.To().ID()}] = true
		}
		for _, x := range nodes {
			for _, y := range nodes {
				between := want[edge{f: x.ID(), t: y.ID()}] || want[edge{f: y.ID(), t: x.ID()}]

				if has := g.HasEdgeBetween(x.ID(), y.ID()); has != between {
					if has {
						t.Errorf("unexpected edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					} else {
						t.Errorf("missing edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					}
				} else {
					if want[edge{f: x.ID(), t: y.ID()}] && g.Edge(x.ID(), y.ID()) == nil {
						t.Errorf("missing edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					}
					if between && !g.HasEdgeBetween(x.ID(), y.ID()) {
						t.Errorf("missing edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					}
					if g.Node(x.ID()) == nil {
						t.Errorf("missing from node for test %q: %v", test.name, x.ID())
					}
					if g.Node(y.ID()) == nil {
						t.Errorf("missing to node for test %q: %v", test.name, y.ID())
					}
				}

				switch g := g.(type) {
				case graph.Directed:
					u := x
					v := y
					if has := g.HasEdgeFromTo(u.ID(), v.ID()); has != want[edge{f: u.ID(), t: v.ID()}] {
						if has {
							t.Errorf("unexpected edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
						} else {
							t.Errorf("missing edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
						}
						continue
					}
					// Edge has already been tested above.
					if g.Node(u.ID()) == nil {
						t.Errorf("missing from node for test %q: %v", test.name, u.ID())
					}
					if g.Node(v.ID()) == nil {
						t.Errorf("missing to node for test %q: %v", test.name, v.ID())
					}

				case graph.Undirected:
					// HasEdgeBetween is already tested above.
					if between && g.Edge(x.ID(), y.ID()) == nil {
						t.Errorf("missing edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					}
					if between && g.EdgeBetween(x.ID(), y.ID()) == nil {
						t.Errorf("missing edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					}
				}
			}
		}
	}
}

// ReturnAdjacentNodes tests the constructed graph for the ability to correctly
// return the nodes reachable from each node within the graph. This is a check
// of the From method of graph.Graph and the To method of graph.Directed.
// ReturnAdjacentNodes also checks that the nodes and traversed edges exist
// within the graph, checking the Node, Edge, EdgeBetween and HasEdgeBetween
// methods of graph.Graph, the EdgeBetween method of graph.Undirected and the
// HasEdgeFromTo method of graph.Directed.
func ReturnAdjacentNodes(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, nodes, edges, _, _, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}

		want := make(map[edge]bool)
		for _, e := range edges {
			want[edge{f: e.From().ID(), t: e.To().ID()}] = true
		}
		for _, x := range nodes {
			switch g := g.(type) {
			case graph.Directed:
				// Test forward.
				u := x
				it := g.From(u.ID())
				for i := 0; it != nil && it.Next(); i++ {
					v := it.Node()
					if i == 0 && g.Node(u.ID()) == nil {
						t.Errorf("missing from node for test %q: %v", test.name, u.ID())
					}
					if g.Node(v.ID()) == nil {
						t.Errorf("missing to node for test %q: %v", test.name, v.ID())
					}
					if g.Edge(u.ID(), v.ID()) == nil {
						t.Errorf("missing from edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
					}
					if !g.HasEdgeBetween(u.ID(), v.ID()) {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
					}
					if !g.HasEdgeFromTo(u.ID(), v.ID()) {
						t.Errorf("missing from edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
					}
					if !want[edge{f: u.ID(), t: v.ID()}] {
						t.Errorf("unexpected edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
					}
				}

				// Test backward.
				v := x
				it = g.To(v.ID())
				for i := 0; it != nil && it.Next(); i++ {
					u := it.Node()
					if i == 0 && g.Node(v.ID()) == nil {
						t.Errorf("missing to node for test %q: %v", test.name, v.ID())
					}
					if g.Node(u.ID()) == nil {
						t.Errorf("missing from node for test %q: %v", test.name, u.ID())
					}
					if g.Edge(u.ID(), v.ID()) == nil {
						t.Errorf("missing from edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
						continue
					}
					if !g.HasEdgeBetween(u.ID(), v.ID()) {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
						continue
					}
					if !g.HasEdgeFromTo(u.ID(), v.ID()) {
						t.Errorf("missing from edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
						continue
					}
					if !want[edge{f: u.ID(), t: v.ID()}] {
						t.Errorf("unexpected edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
					}
				}

			case graph.Undirected:
				u := x
				it := g.From(u.ID())
				for i := 0; it != nil && it.Next(); i++ {
					v := it.Node()
					if i == 0 && g.Node(u.ID()) == nil {
						t.Errorf("missing from node for test %q: %v", test.name, u.ID())
					}
					if g.Edge(u.ID(), v.ID()) == nil {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
						continue
					}
					if g.EdgeBetween(u.ID(), v.ID()) == nil {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
						continue
					}
					if !g.HasEdgeBetween(u.ID(), v.ID()) {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
						continue
					}
					between := want[edge{f: u.ID(), t: v.ID()}] || want[edge{f: v.ID(), t: u.ID()}]
					if !between {
						t.Errorf("unexpected edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
					}
				}

			default:
				u := x
				it := g.From(u.ID())
				for i := 0; it != nil && it.Next(); i++ {
					v := it.Node()
					if i == 0 && g.Node(u.ID()) == nil {
						t.Errorf("missing from node for test %q: %v", test.name, u.ID())
					}
					if g.Edge(u.ID(), v.ID()) == nil {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
						continue
					}
					if !g.HasEdgeBetween(u.ID(), v.ID()) {
						t.Errorf("missing from edge for test %q: (%v)--(%v)", test.name, u.ID(), v.ID())
						continue
					}
					between := want[edge{f: u.ID(), t: v.ID()}] || want[edge{f: v.ID(), t: u.ID()}]
					if !between {
						t.Errorf("unexpected edge for test %q: (%v)->(%v)", test.name, u.ID(), v.ID())
					}
				}
			}
		}
	}
}

// Weight tests the constructed graph for the ability to correctly return
// the weight between to nodes, checking the Weight method of graph.Weighted.
//
// The self and absent values returned by the Builder should match the values
// used by the Weight method.
func Weight(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, nodes, _, self, absent, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}
		wg, ok := g.(graph.Weighted)
		if !ok {
			t.Errorf("invalid graph type for test %q: %T is not graph.Weighted", test.name, g)
		}
		_, multi := g.(graph.Multigraph)

		for _, x := range nodes {
			for _, y := range nodes {
				w, ok := wg.Weight(x.ID(), y.ID())
				e := wg.WeightedEdge(x.ID(), y.ID())
				switch {
				case !ok:
					if e != nil {
						t.Errorf("missing edge weight for existing edge for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())
					}
					if !same(w, absent) {
						t.Errorf("unexpected absent weight for test %q: got:%v want:%v", test.name, w, absent)
					}

				case !multi && x.ID() == y.ID():
					if !same(w, self) {
						t.Errorf("unexpected self weight for test %q: got:%v want:%v", test.name, w, self)
					}

				case e == nil:
					t.Errorf("missing edge for existing non-self weight for test %q: (%v)--(%v)", test.name, x.ID(), y.ID())

				case e.Weight() != w:
					t.Errorf("weight mismatch for test %q: edge=%v graph=%v", test.name, e.Weight(), w)
				}
			}
		}
	}
}

// AdjacencyMatrix tests the constructed graph for the ability to correctly
// return an adjacency matrix that matches the weights returned by the graphs
// Weight method.
//
// The self and absent values returned by the Builder should match the values
// used by the Weight method.
func AdjacencyMatrix(t *testing.T, b Builder) {
	for _, test := range testCases {
		g, nodes, _, self, absent, ok := b(test.nodes, test.edges, test.self, test.absent)
		if !ok {
			t.Logf("skipping test case: %q", test.name)
			continue
		}
		wg, ok := g.(graph.Weighted)
		if !ok {
			t.Errorf("invalid graph type for test %q: %T is not graph.Weighted", test.name, g)
		}
		mg, ok := g.(matrixer)
		if !ok {
			t.Errorf("invalid graph type for test %q: %T cannot return adjacency matrix", test.name, g)
		}
		m := mg.Matrix()

		r, c := m.Dims()
		if r != c || r != len(nodes) {
			t.Errorf("dimension mismatch for test %q: r=%d c=%d order=%d", test.name, r, c, len(nodes))
		}

		for _, x := range nodes {
			i := int(x.ID())
			for _, y := range nodes {
				j := int(y.ID())
				w, ok := wg.Weight(x.ID(), y.ID())
				switch {
				case !ok:
					if !same(m.At(i, j), absent) {
						t.Errorf("weight mismatch for test %q: (%v)--(%v) matrix=%v graph=%v", test.name, x.ID(), y.ID(), m.At(i, j), w)
					}
				case x.ID() == y.ID():
					if !same(m.At(i, j), self) {
						t.Errorf("weight mismatch for test %q: (%v)--(%v) matrix=%v graph=%v", test.name, x.ID(), y.ID(), m.At(i, j), w)
					}
				default:
					if !same(m.At(i, j), w) {
						t.Errorf("weight mismatch for test %q: (%v)--(%v) matrix=%v graph=%v", test.name, x.ID(), y.ID(), m.At(i, j), w)
					}
				}
			}
		}
	}
}

// lexicalEdges sorts a collection of edges lexically on the
// keys: from.ID > to.ID > [line.ID] > [weight].
type lexicalEdges []graph.Edge

func (e lexicalEdges) Len() int { return len(e) }
func (e lexicalEdges) Less(i, j int) bool {
	if e[i].From().ID() < e[j].From().ID() {
		return true
	}
	sf := e[i].From().ID() == e[j].From().ID()
	if sf && e[i].To().ID() < e[j].To().ID() {
		return true
	}
	st := e[i].To().ID() == e[j].To().ID()
	li, oki := e[i].(graph.Line)
	lj, okj := e[j].(graph.Line)
	if oki != okj {
		panic(fmt.Sprintf("testgraph: mismatched types %T != %T", e[i], e[j]))
	}
	if !oki {
		return sf && st && lessWeight(e[i], e[j])
	}
	if sf && st && li.ID() < lj.ID() {
		return true
	}
	return sf && st && li.ID() == lj.ID() && lessWeight(e[i], e[j])
}
func (e lexicalEdges) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// lexicalUndirectedEdges sorts a collection of edges lexically on the
// keys: lo.ID > hi.ID > [line.ID] > [weight].
type lexicalUndirectedEdges []graph.Edge

func (e lexicalUndirectedEdges) Len() int { return len(e) }
func (e lexicalUndirectedEdges) Less(i, j int) bool {
	lidi, hidi, _ := undirectedIDs(e[i])
	lidj, hidj, _ := undirectedIDs(e[j])

	if lidi < lidj {
		return true
	}
	sl := lidi == lidj
	if sl && hidi < hidj {
		return true
	}
	sh := hidi == hidj
	li, oki := e[i].(graph.Line)
	lj, okj := e[j].(graph.Line)
	if oki != okj {
		panic(fmt.Sprintf("testgraph: mismatched types %T != %T", e[i], e[j]))
	}
	if !oki {
		return sl && sh && lessWeight(e[i], e[j])
	}
	if sl && sh && li.ID() < lj.ID() {
		return true
	}
	return sl && sh && li.ID() == lj.ID() && lessWeight(e[i], e[j])
}
func (e lexicalUndirectedEdges) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func lessWeight(ei, ej graph.Edge) bool {
	wei, oki := ei.(graph.WeightedEdge)
	wej, okj := ej.(graph.WeightedEdge)
	if oki != okj {
		panic(fmt.Sprintf("testgraph: mismatched types %T != %T", ei, ej))
	}
	if !oki {
		return false
	}
	return wei.Weight() < wej.Weight()
}

// undirectedEdgeSetEqual returned whether a pair of undirected edge
// slices sorted by lexicalUndirectedEdges are equal.
func undirectedEdgeSetEqual(a, b []graph.Edge) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	if !undirectedEdgeEqual(a[0], b[0]) {
		return false
	}
	i, j := 0, 0
	for {
		switch {
		case i == len(a)-1 && j == len(b)-1:
			return true

		case i < len(a)-1 && undirectedEdgeEqual(a[i+1], b[j]):
			i++

		case j < len(b)-1 && undirectedEdgeEqual(a[i], b[j+1]):
			j++

		case i < len(a)-1 && j < len(b)-1 && undirectedEdgeEqual(a[i+1], b[j+1]):
			i++
			j++

		default:
			return false
		}
	}
}

// undirectedEdgeEqual returns whether a pair of undirected edges are equal
// after canonicalising from and to IDs by numerical sort order.
func undirectedEdgeEqual(a, b graph.Edge) bool {
	loa, hia, inva := undirectedIDs(a)
	lob, hib, invb := undirectedIDs(b)
	// Use reflect.DeepEqual if the edges are parallel
	// rather anti-parallel.
	if inva == invb {
		return reflect.DeepEqual(a, b)
	}
	if loa != lob || hia != hib {
		return false
	}
	la, oka := a.(graph.Line)
	lb, okb := b.(graph.Line)
	if !oka && !okb {
		return true
	}
	if la.ID() != lb.ID() {
		return false
	}
	wea, oka := a.(graph.WeightedEdge)
	web, okb := b.(graph.WeightedEdge)
	if !oka && !okb {
		return true
	}
	return wea.Weight() == web.Weight()
}

// undirectedIDs returns a numerical sort ordered canonicalisation of the
// IDs of e.
func undirectedIDs(e graph.Edge) (lo, hi int64, inverted bool) {
	lid := e.From().ID()
	hid := e.To().ID()
	if hid < lid {
		inverted = true
		hid, lid = lid, hid
	}
	return lid, hid, inverted
}

type edge struct {
	f, t int64
}

func same(a, b float64) bool {
	return (math.IsNaN(a) && math.IsNaN(b)) || a == b
}
