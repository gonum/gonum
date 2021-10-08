// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterator_test

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

var orderedNodesTests = []struct {
	nodes []graph.Node
}{
	{nodes: nil},
	{nodes: []graph.Node{simple.Node(1)}},
	{nodes: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)}},
	{nodes: []graph.Node{simple.Node(5), simple.Node(3), simple.Node(2), simple.Node(1)}},
}

func TestOrderedNodesIterate(t *testing.T) {
	for _, test := range orderedNodesTests {
		it := iterator.NewOrderedNodes(test.nodes)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.nodes) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.nodes))
			}
			var got []graph.Node
			for it.Next() {
				got = append(got, it.Node())
				if len(got)+it.Len() != len(test.nodes) {
					t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(test.nodes)-len(got))
				}
			}
			want := test.nodes
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

func TestOrderedNodesSlice(t *testing.T) {
	for _, test := range orderedNodesTests {
		it := iterator.NewOrderedNodes(test.nodes)
		for i := 0; i < 2; i++ {
			got := it.NodeSlice()
			want := test.nodes
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
			}
			it.Reset()
		}
	}
}

var implicitNodesTests = []struct {
	beg, end int
	new      func(int) graph.Node
	want     []graph.Node
}{
	{
		beg: 1, end: 1,
		want: nil,
	},
	{
		beg: 1, end: 2,
		new:  newSimpleNode,
		want: []graph.Node{simple.Node(1)},
	},
	{
		beg: 1, end: 5,
		new:  newSimpleNode,
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(4)},
	},
}

func newSimpleNode(id int) graph.Node { return simple.Node(id) }

func TestImplicitNodesIterate(t *testing.T) {
	for _, test := range implicitNodesTests {
		it := iterator.NewImplicitNodes(test.beg, test.end, test.new)
		for i := 0; i < 2; i++ {
			if it.Len() != len(test.want) {
				t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.want))
			}
			var got []graph.Node
			for it.Next() {
				got = append(got, it.Node())
				if len(got)+it.Len() != test.end-test.beg {
					t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), (test.end-test.beg)-len(got))
				}
			}
			if it.Len() != 0 {
				t.Errorf("unexpected depleted iterator length for round %d: got:%d want:0", i, it.Len())
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, test.want)
			}
			it.Reset()
		}
	}
}

var nodesTests = []struct {
	nodes map[int64]graph.Node
}{
	{nodes: nil},
	{nodes: make(map[int64]graph.Node)},
	{nodes: map[int64]graph.Node{1: simple.Node(1)}},
	{nodes: map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 5: simple.Node(5)}},
	{nodes: map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)}},
}

func TestIterateNodes(t *testing.T) {
	for _, typ := range []struct {
		name string
		new  func(map[int64]graph.Node) graph.Nodes
	}{
		{name: "Nodes", new: func(n map[int64]graph.Node) graph.Nodes { return iterator.NewNodes(n) }},
		{name: "LazyOrderedNodes", new: func(n map[int64]graph.Node) graph.Nodes { return iterator.NewLazyOrderedNodes(n) }},
	} {
		t.Run(typ.name, func(t *testing.T) {
			for _, test := range nodesTests {
				it := typ.new(test.nodes)
				for i := 0; i < 2; i++ {
					if it.Len() != len(test.nodes) {
						t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(test.nodes))
					}
					var got map[int64]graph.Node
					if test.nodes != nil {
						got = make(map[int64]graph.Node)
					}
					for it.Next() {
						n := it.Node()
						got[n.ID()] = n
						if len(got)+it.Len() != len(test.nodes) {
							t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(test.nodes))
						}
					}
					want := test.nodes
					if !reflect.DeepEqual(got, want) {
						t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, want)
					}
					func() {
						defer func() {
							r := recover()
							if r != nil {
								t.Errorf("unexpected panic: %v", r)
							}
						}()
						it.Next()
					}()
					it.Reset()
				}
			}
		})
	}
}

var nodesByEdgeTests = []struct {
	n     int64
	edges map[int64]graph.Edge
	want  map[int64]graph.Node
}{
	// The actual values of the edge stored in the edge
	// map leading to each node are not used, so they are
	// filled with nil values.
	{
		n:     6,
		edges: nil,
		want:  nil,
	},
	{
		n:     6,
		edges: make(map[int64]graph.Edge),
		want:  make(map[int64]graph.Node),
	},
	{
		n:     6,
		edges: map[int64]graph.Edge{1: nil},
		want:  map[int64]graph.Node{1: simple.Node(1)},
	},
	{
		n:     6,
		edges: map[int64]graph.Edge{1: nil, 2: nil, 3: nil, 5: nil},
		want:  map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 5: simple.Node(5)},
	},
	{
		n:     6,
		edges: map[int64]graph.Edge{5: nil, 3: nil, 2: nil, 1: nil},
		want:  map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)},
	},
}

func TestNodesByEdgeIterate(t *testing.T) {
	for _, typ := range []struct {
		name string
		new  func(map[int64]graph.Node, map[int64]graph.Edge) graph.Nodes
	}{
		{
			name: "NodesByEdge",
			new: func(n map[int64]graph.Node, e map[int64]graph.Edge) graph.Nodes {
				return iterator.NewNodesByEdge(n, e)
			}},
		{
			name: "LazyOrderedNodesByEdge",
			new: func(n map[int64]graph.Node, e map[int64]graph.Edge) graph.Nodes {
				return iterator.NewLazyOrderedNodesByEdge(n, e)
			},
		},
	} {
		t.Run(typ.name, func(t *testing.T) {
			for _, test := range nodesByEdgeTests {
				nodes := make(map[int64]graph.Node)
				for i := int64(0); i < test.n; i++ {
					nodes[i] = simple.Node(i)
				}

				it := typ.new(nodes, test.edges)
				for i := 0; i < 2; i++ {
					if it.Len() != len(test.edges) {
						t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(nodes))
					}
					var got map[int64]graph.Node
					if test.edges != nil {
						got = make(map[int64]graph.Node)
					}
					for it.Next() {
						n := it.Node()
						got[n.ID()] = n
						if len(got)+it.Len() != len(test.edges) {
							t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(nodes))
						}
					}
					if !reflect.DeepEqual(got, test.want) {
						t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, test.want)
					}
					func() {
						defer func() {
							r := recover()
							if r != nil {
								t.Errorf("unexpected panic: %v", r)
							}
						}()
						it.Next()
					}()
					it.Reset()
				}
			}
		})
	}
}

var nodesByWeightedEdgeTests = []struct {
	n     int64
	edges map[int64]graph.WeightedEdge
	want  map[int64]graph.Node
}{
	// The actual values of the edges stored in the edge
	// map leading to each node are not used, so they are
	// filled with nil values.
	{
		n:     6,
		edges: nil,
		want:  nil,
	},
	{
		n:     6,
		edges: make(map[int64]graph.WeightedEdge),
		want:  make(map[int64]graph.Node),
	},
	{
		n:     6,
		edges: map[int64]graph.WeightedEdge{1: nil},
		want:  map[int64]graph.Node{1: simple.Node(1)},
	},
	{
		n:     6,
		edges: map[int64]graph.WeightedEdge{1: nil, 2: nil, 3: nil, 5: nil},
		want:  map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 5: simple.Node(5)},
	},
	{
		n:     6,
		edges: map[int64]graph.WeightedEdge{5: nil, 3: nil, 2: nil, 1: nil},
		want:  map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)},
	},
}

func TestNodesByWeightedEdgeIterate(t *testing.T) {
	for _, typ := range []struct {
		name string
		new  func(map[int64]graph.Node, map[int64]graph.WeightedEdge) graph.Nodes
	}{
		{
			name: "NodesByWeightedEdge",
			new: func(n map[int64]graph.Node, e map[int64]graph.WeightedEdge) graph.Nodes {
				return iterator.NewNodesByWeightedEdge(n, e)
			}},
		{
			name: "LazyOrderedNodesByWeightedEdge",
			new: func(n map[int64]graph.Node, e map[int64]graph.WeightedEdge) graph.Nodes {
				return iterator.NewLazyOrderedNodesByWeightedEdge(n, e)
			},
		},
	} {
		t.Run(typ.name, func(t *testing.T) {
			for _, test := range nodesByWeightedEdgeTests {
				nodes := make(map[int64]graph.Node)
				for i := int64(0); i < test.n; i++ {
					nodes[i] = simple.Node(i)
				}

				it := typ.new(nodes, test.edges)
				for i := 0; i < 2; i++ {
					if it.Len() != len(test.edges) {
						t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(nodes))
					}
					var got map[int64]graph.Node
					if test.edges != nil {
						got = make(map[int64]graph.Node)
					}
					for it.Next() {
						n := it.Node()
						got[n.ID()] = n
						if len(got)+it.Len() != len(test.edges) {
							t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(nodes))
						}
					}
					if !reflect.DeepEqual(got, test.want) {
						t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, test.want)
					}
					func() {
						defer func() {
							r := recover()
							if r != nil {
								t.Errorf("unexpected panic: %v", r)
							}
						}()
						it.Next()
					}()
					it.Reset()
				}
			}
		})
	}
}

var nodesByLinesTests = []struct {
	n     int64
	lines map[int64]map[int64]graph.Line
	want  map[int64]graph.Node
}{
	// The actual values of the lines stored in the line
	// collection leading to each node are not used, so
	// they are filled with nil.
	{
		n:     6,
		lines: nil,
		want:  nil,
	},
	{
		n:     6,
		lines: make(map[int64]map[int64]graph.Line),
		want:  make(map[int64]graph.Node),
	},
	{
		n:     6,
		lines: map[int64]map[int64]graph.Line{1: nil},
		want:  map[int64]graph.Node{1: simple.Node(1)},
	},
	{
		n:     6,
		lines: map[int64]map[int64]graph.Line{1: nil, 2: nil, 3: nil, 5: nil},
		want:  map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 5: simple.Node(5)},
	},
	{
		n:     6,
		lines: map[int64]map[int64]graph.Line{5: nil, 3: nil, 2: nil, 1: nil},
		want:  map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)},
	},
}

func TestNodesByLinesIterate(t *testing.T) {
	for _, typ := range []struct {
		name string
		new  func(map[int64]graph.Node, map[int64]map[int64]graph.Line) graph.Nodes
	}{
		{
			name: "NodesByLines",
			new: func(n map[int64]graph.Node, e map[int64]map[int64]graph.Line) graph.Nodes {
				return iterator.NewNodesByLines(n, e)
			}},
		{
			name: "LazyOrderedNodesByLines",
			new: func(n map[int64]graph.Node, e map[int64]map[int64]graph.Line) graph.Nodes {
				return iterator.NewLazyOrderedNodesByLines(n, e)
			},
		},
	} {
		t.Run(typ.name, func(t *testing.T) {
			for _, test := range nodesByLinesTests {
				nodes := make(map[int64]graph.Node)
				for i := int64(0); i < test.n; i++ {
					nodes[i] = simple.Node(i)
				}

				it := typ.new(nodes, test.lines)
				for i := 0; i < 2; i++ {
					if it.Len() != len(test.lines) {
						t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(nodes))
					}
					var got map[int64]graph.Node
					if test.lines != nil {
						got = make(map[int64]graph.Node)
					}
					for it.Next() {
						n := it.Node()
						got[n.ID()] = n
						if len(got)+it.Len() != len(test.lines) {
							t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(nodes))
						}
					}
					if !reflect.DeepEqual(got, test.want) {
						t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, test.want)
					}
					func() {
						defer func() {
							r := recover()
							if r != nil {
								t.Errorf("unexpected panic: %v", r)
							}
						}()
						it.Next()
					}()
					it.Reset()
				}
			}
		})
	}
}

var nodesByWeightedLinesTests = []struct {
	n     int64
	lines map[int64]map[int64]graph.WeightedLine
	want  map[int64]graph.Node
}{
	// The actual values of the lines stored in the line
	// collection leading to each node are not used, so
	// they are filled with nil.
	{
		n:     6,
		lines: nil,
		want:  nil,
	},
	{
		n:     6,
		lines: make(map[int64]map[int64]graph.WeightedLine),
		want:  make(map[int64]graph.Node),
	},
	{
		n:     6,
		lines: map[int64]map[int64]graph.WeightedLine{1: nil},
		want:  map[int64]graph.Node{1: simple.Node(1)},
	},
	{
		n:     6,
		lines: map[int64]map[int64]graph.WeightedLine{1: nil, 2: nil, 3: nil, 5: nil},
		want:  map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 5: simple.Node(5)},
	},
	{
		n:     6,
		lines: map[int64]map[int64]graph.WeightedLine{5: nil, 3: nil, 2: nil, 1: nil},
		want:  map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)},
	},
}

func TestNodesByWeightedLinesIterate(t *testing.T) {
	for _, typ := range []struct {
		name string
		new  func(map[int64]graph.Node, map[int64]map[int64]graph.WeightedLine) graph.Nodes
	}{
		{
			name: "NodesByWeightedLines",
			new: func(n map[int64]graph.Node, e map[int64]map[int64]graph.WeightedLine) graph.Nodes {
				return iterator.NewNodesByWeightedLines(n, e)
			}},
		{
			name: "LazyOrderedNodesByWeightedLines",
			new: func(n map[int64]graph.Node, e map[int64]map[int64]graph.WeightedLine) graph.Nodes {
				return iterator.NewLazyOrderedNodesByWeightedLines(n, e)
			},
		},
	} {
		t.Run(typ.name, func(t *testing.T) {
			for _, test := range nodesByWeightedLinesTests {
				nodes := make(map[int64]graph.Node)
				for i := int64(0); i < test.n; i++ {
					nodes[i] = simple.Node(i)
				}

				it := typ.new(nodes, test.lines)
				for i := 0; i < 2; i++ {
					if it.Len() != len(test.lines) {
						t.Errorf("unexpected iterator length for round %d: got:%d want:%d", i, it.Len(), len(nodes))
					}
					var got map[int64]graph.Node
					if test.lines != nil {
						got = make(map[int64]graph.Node)
					}
					for it.Next() {
						n := it.Node()
						got[n.ID()] = n
						if len(got)+it.Len() != len(test.lines) {
							t.Errorf("unexpected iterator length during iteration for round %d: got:%d want:%d", i, it.Len(), len(nodes))
						}
					}
					if !reflect.DeepEqual(got, test.want) {
						t.Errorf("unexpected iterator output for round %d: got:%#v want:%#v", i, got, test.want)
					}
					func() {
						defer func() {
							r := recover()
							if r != nil {
								t.Errorf("unexpected panic: %v", r)
							}
						}()
						it.Next()
					}()
					it.Reset()
				}
			}
		})
	}
}

type nodeSlicer interface {
	graph.Nodes
	graph.NodeSlicer
}

var nodeSlicerTests = []struct {
	nodes nodeSlicer
	want  []graph.Node
}{
	{
		nodes: iterator.NewOrderedNodes([]graph.Node{simple.Node(1)}),
		want:  []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewOrderedNodes([]graph.Node{simple.Node(1), simple.Node(2), simple.Node(3)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3)},
	},

	{
		nodes: iterator.NewImplicitNodes(1, 2, func(id int) graph.Node { return simple.Node(id) }),
		want:  []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewImplicitNodes(1, 4, func(id int) graph.Node { return simple.Node(id) }),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3)},
	},

	{
		nodes: iterator.NewNodes(map[int64]graph.Node{1: simple.Node(1)}),
		want:  []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewNodes(map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3)},
	},
	{
		nodes: iterator.NewNodes(map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(4), simple.Node(5)},
	},
	{
		nodes: iterator.NewNodes(map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},

	{
		nodes: iterator.NewLazyOrderedNodes(map[int64]graph.Node{1: simple.Node(1)}),
		want:  []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewLazyOrderedNodes(map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3)},
	},
	{
		nodes: iterator.NewLazyOrderedNodes(map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(4), simple.Node(5)},
	},
	{
		nodes: iterator.NewLazyOrderedNodes(map[int64]graph.Node{5: simple.Node(5), 3: simple.Node(3), 2: simple.Node(2), 1: simple.Node(1)}),
		want:  []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},

	// The actual values of the edges stored in the edge
	// map leading to each node are not used, so they are
	// filled with nil values.
	//
	// The three other constructors for NodesByEdge are not
	// tested for this behaviour since they have already
	// been tested above.
	{
		nodes: iterator.NewNodesByEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.Edge{1: nil},
		),
		want: []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewNodesByEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.Edge{1: nil, 2: nil, 3: nil, 5: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},
	{
		nodes: iterator.NewNodesByEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.Edge{5: nil, 3: nil, 2: nil, 1: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},

	{
		nodes: iterator.NewLazyOrderedNodesByEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.Edge{1: nil},
		),
		want: []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.Edge{1: nil, 2: nil, 3: nil, 5: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.Edge{5: nil, 3: nil, 2: nil, 1: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},

	{
		nodes: iterator.NewLazyOrderedNodesByWeightedEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.WeightedEdge{1: nil},
		),
		want: []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByWeightedEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.WeightedEdge{1: nil, 2: nil, 3: nil, 5: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByWeightedEdge(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]graph.WeightedEdge{5: nil, 3: nil, 2: nil, 1: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},

	{
		nodes: iterator.NewLazyOrderedNodesByLines(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]map[int64]graph.Line{1: nil},
		),
		want: []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByLines(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]map[int64]graph.Line{1: nil, 2: nil, 3: nil, 5: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByLines(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]map[int64]graph.Line{5: nil, 3: nil, 2: nil, 1: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},

	{
		nodes: iterator.NewLazyOrderedNodesByWeightedLines(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]map[int64]graph.WeightedLine{1: nil},
		),
		want: []graph.Node{simple.Node(1)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByWeightedLines(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]map[int64]graph.WeightedLine{1: nil, 2: nil, 3: nil, 5: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},
	{
		nodes: iterator.NewLazyOrderedNodesByWeightedLines(
			map[int64]graph.Node{1: simple.Node(1), 2: simple.Node(2), 3: simple.Node(3), 4: simple.Node(4), 5: simple.Node(5)},
			map[int64]map[int64]graph.WeightedLine{5: nil, 3: nil, 2: nil, 1: nil},
		),
		want: []graph.Node{simple.Node(1), simple.Node(2), simple.Node(3), simple.Node(5)},
	},
}

func TestNodeSlicers(t *testing.T) {
	for k, test := range nodeSlicerTests {
		wantLen := test.nodes.Len()
		for i := 0; i < wantLen; i++ {
			var gotIter []graph.Node
			for n := 0; n < i; n++ {
				ok := test.nodes.Next()
				if !ok {
					t.Errorf("test %d: unexpected failed Next call at position %d of len %d", k, n, wantLen)
				}
				gotIter = append(gotIter, test.nodes.Node())
			}
			gotSlice := test.nodes.NodeSlice()
			if test.nodes.Next() {
				t.Errorf("test %d: expected no further iteration possible after NodeSlice with %d pre-iterations of %d", k, i, wantLen)
			}

			if gotLen := len(gotIter) + len(gotSlice); gotLen != wantLen {
				t.Errorf("test %d: unexpected total node count: got:%d want:%d", k, gotLen, wantLen)
			}
			got := append(gotIter, gotSlice...)
			ordered.ByID(got)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("test %d: unexpected node slice:\ngot: %v\nwant:%v", k, got, test.want)
			}

			test.nodes.Reset()
		}
	}
}
