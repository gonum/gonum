// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterator_test

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
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

func TestNodesIterate(t *testing.T) {
	for _, test := range nodesTests {
		it := iterator.NewNodes(test.nodes)
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
	for _, test := range nodesByEdgeTests {
		nodes := make(map[int64]graph.Node)
		for i := int64(0); i < test.n; i++ {
			nodes[i] = simple.Node(i)
		}

		it := iterator.NewNodesByEdge(nodes, test.edges)
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
	for _, test := range nodesByWeightedEdgeTests {
		nodes := make(map[int64]graph.Node)
		for i := int64(0); i < test.n; i++ {
			nodes[i] = simple.Node(i)
		}

		it := iterator.NewNodesByWeightedEdge(nodes, test.edges)
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
	for _, test := range nodesByLinesTests {
		nodes := make(map[int64]graph.Node)
		for i := int64(0); i < test.n; i++ {
			nodes[i] = simple.Node(i)
		}

		it := iterator.NewNodesByLines(nodes, test.lines)
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
	for _, test := range nodesByWeightedLinesTests {
		nodes := make(map[int64]graph.Node)
		for i := int64(0); i < test.n; i++ {
			nodes[i] = simple.Node(i)
		}

		it := iterator.NewNodesByWeightedLines(nodes, test.lines)
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
}
