// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"math"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

func init() {
	for _, test := range spanningTreeTests {
		var w float64
		for _, e := range test.treeEdges {
			w += e.W
		}
		if w != test.want {
			panic(fmt.Sprintf("bad test: %s weight mismatch: %v != %v", test.name, w, test.want))
		}
	}
}

type spanningGraph interface {
	graph.MutableUndirected
	graph.Weighter
	Edges() []graph.Edge
}

var spanningTreeTests = []struct {
	name      string
	graph     func() spanningGraph
	edges     []concrete.Edge
	want      float64
	treeEdges []concrete.Edge
}{
	{
		name:  "Empty",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		want:  0,
	},
	{
		// https://upload.wikimedia.org/wikipedia/commons/f/f7/Prim%27s_algorithm.svg
		// Modified to make edge weights unique; A--B is increased to 2.5 otherwise
		// to prevent the alternative solution being found.
		name:  "Prim WP figure 1",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('B'), W: 2.5},
			{F: concrete.Node('A'), T: concrete.Node('D'), W: 1},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 2},
			{F: concrete.Node('C'), T: concrete.Node('D'), W: 3},
		},

		want: 6,
		treeEdges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('D'), W: 1},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 2},
			{F: concrete.Node('C'), T: concrete.Node('D'), W: 3},
		},
	},
	{
		// https://upload.wikimedia.org/wikipedia/commons/5/5c/MST_kruskal_en.gif
		name:  "Kruskal WP figure 1",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('a'), T: concrete.Node('b'), W: 3},
			{F: concrete.Node('a'), T: concrete.Node('e'), W: 1},
			{F: concrete.Node('b'), T: concrete.Node('c'), W: 5},
			{F: concrete.Node('b'), T: concrete.Node('e'), W: 4},
			{F: concrete.Node('c'), T: concrete.Node('d'), W: 2},
			{F: concrete.Node('c'), T: concrete.Node('e'), W: 6},
			{F: concrete.Node('d'), T: concrete.Node('e'), W: 7},
		},

		want: 11,
		treeEdges: []concrete.Edge{
			{F: concrete.Node('a'), T: concrete.Node('b'), W: 3},
			{F: concrete.Node('a'), T: concrete.Node('e'), W: 1},
			{F: concrete.Node('b'), T: concrete.Node('c'), W: 5},
			{F: concrete.Node('c'), T: concrete.Node('d'), W: 2},
		},
	},
	{
		// https://upload.wikimedia.org/wikipedia/commons/8/87/Kruskal_Algorithm_6.svg
		name:  "Kruskal WP example",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('B'), W: 7},
			{F: concrete.Node('A'), T: concrete.Node('D'), W: 5},
			{F: concrete.Node('B'), T: concrete.Node('C'), W: 8},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 9},
			{F: concrete.Node('B'), T: concrete.Node('E'), W: 7},
			{F: concrete.Node('C'), T: concrete.Node('E'), W: 5},
			{F: concrete.Node('D'), T: concrete.Node('E'), W: 15},
			{F: concrete.Node('D'), T: concrete.Node('F'), W: 6},
			{F: concrete.Node('E'), T: concrete.Node('F'), W: 8},
			{F: concrete.Node('E'), T: concrete.Node('G'), W: 9},
			{F: concrete.Node('F'), T: concrete.Node('G'), W: 11},
		},

		want: 39,
		treeEdges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('B'), W: 7},
			{F: concrete.Node('A'), T: concrete.Node('D'), W: 5},
			{F: concrete.Node('B'), T: concrete.Node('E'), W: 7},
			{F: concrete.Node('C'), T: concrete.Node('E'), W: 5},
			{F: concrete.Node('D'), T: concrete.Node('F'), W: 6},
			{F: concrete.Node('E'), T: concrete.Node('G'), W: 9},
		},
	},
	{
		// https://upload.wikimedia.org/wikipedia/commons/2/2e/Boruvka%27s_algorithm_%28Sollin%27s_algorithm%29_Anim.gif
		name:  "Borůvka WP example",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('B'), W: 13},
			{F: concrete.Node('A'), T: concrete.Node('C'), W: 6},
			{F: concrete.Node('B'), T: concrete.Node('C'), W: 7},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 1},
			{F: concrete.Node('C'), T: concrete.Node('D'), W: 14},
			{F: concrete.Node('C'), T: concrete.Node('E'), W: 8},
			{F: concrete.Node('C'), T: concrete.Node('H'), W: 20},
			{F: concrete.Node('D'), T: concrete.Node('E'), W: 9},
			{F: concrete.Node('D'), T: concrete.Node('F'), W: 3},
			{F: concrete.Node('E'), T: concrete.Node('F'), W: 2},
			{F: concrete.Node('E'), T: concrete.Node('J'), W: 18},
			{F: concrete.Node('G'), T: concrete.Node('H'), W: 15},
			{F: concrete.Node('G'), T: concrete.Node('I'), W: 5},
			{F: concrete.Node('G'), T: concrete.Node('J'), W: 19},
			{F: concrete.Node('G'), T: concrete.Node('K'), W: 10},
			{F: concrete.Node('H'), T: concrete.Node('J'), W: 17},
			{F: concrete.Node('I'), T: concrete.Node('K'), W: 11},
			{F: concrete.Node('J'), T: concrete.Node('K'), W: 16},
			{F: concrete.Node('J'), T: concrete.Node('L'), W: 4},
			{F: concrete.Node('K'), T: concrete.Node('L'), W: 12},
		},

		want: 83,
		treeEdges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('C'), W: 6},
			{F: concrete.Node('B'), T: concrete.Node('C'), W: 7},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 1},
			{F: concrete.Node('D'), T: concrete.Node('F'), W: 3},
			{F: concrete.Node('E'), T: concrete.Node('F'), W: 2},
			{F: concrete.Node('E'), T: concrete.Node('J'), W: 18},
			{F: concrete.Node('G'), T: concrete.Node('H'), W: 15},
			{F: concrete.Node('G'), T: concrete.Node('I'), W: 5},
			{F: concrete.Node('G'), T: concrete.Node('K'), W: 10},
			{F: concrete.Node('J'), T: concrete.Node('L'), W: 4},
			{F: concrete.Node('K'), T: concrete.Node('L'), W: 12},
		},
	},
	{
		// https://upload.wikimedia.org/wikipedia/commons/d/d2/Minimum_spanning_tree.svg
		// Nodes labelled row major.
		name:  "Minimum Spanning Tree WP figure 1",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node(1), T: concrete.Node(2), W: 4},
			{F: concrete.Node(1), T: concrete.Node(3), W: 1},
			{F: concrete.Node(1), T: concrete.Node(4), W: 4},
			{F: concrete.Node(2), T: concrete.Node(3), W: 5},
			{F: concrete.Node(2), T: concrete.Node(5), W: 9},
			{F: concrete.Node(2), T: concrete.Node(6), W: 9},
			{F: concrete.Node(2), T: concrete.Node(8), W: 7},
			{F: concrete.Node(3), T: concrete.Node(4), W: 3},
			{F: concrete.Node(3), T: concrete.Node(8), W: 9},
			{F: concrete.Node(4), T: concrete.Node(8), W: 10},
			{F: concrete.Node(4), T: concrete.Node(10), W: 18},
			{F: concrete.Node(5), T: concrete.Node(6), W: 2},
			{F: concrete.Node(5), T: concrete.Node(7), W: 4},
			{F: concrete.Node(5), T: concrete.Node(9), W: 6},
			{F: concrete.Node(6), T: concrete.Node(7), W: 2},
			{F: concrete.Node(6), T: concrete.Node(8), W: 8},
			{F: concrete.Node(7), T: concrete.Node(8), W: 9},
			{F: concrete.Node(7), T: concrete.Node(9), W: 3},
			{F: concrete.Node(7), T: concrete.Node(10), W: 9},
			{F: concrete.Node(8), T: concrete.Node(10), W: 8},
			{F: concrete.Node(9), T: concrete.Node(10), W: 9},
		},

		want: 38,
		treeEdges: []concrete.Edge{
			{F: concrete.Node(1), T: concrete.Node(2), W: 4},
			{F: concrete.Node(1), T: concrete.Node(3), W: 1},
			{F: concrete.Node(2), T: concrete.Node(8), W: 7},
			{F: concrete.Node(3), T: concrete.Node(4), W: 3},
			{F: concrete.Node(5), T: concrete.Node(6), W: 2},
			{F: concrete.Node(6), T: concrete.Node(7), W: 2},
			{F: concrete.Node(6), T: concrete.Node(8), W: 8},
			{F: concrete.Node(7), T: concrete.Node(9), W: 3},
			{F: concrete.Node(8), T: concrete.Node(10), W: 8},
		},
	},

	{
		// https://upload.wikimedia.org/wikipedia/commons/2/2e/Boruvka%27s_algorithm_%28Sollin%27s_algorithm%29_Anim.gif
		// but with C--H and E--J cut.
		name:  "Borůvka WP example cut",
		graph: func() spanningGraph { return concrete.NewGraph(0, math.Inf(1)) },
		edges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('B'), W: 13},
			{F: concrete.Node('A'), T: concrete.Node('C'), W: 6},
			{F: concrete.Node('B'), T: concrete.Node('C'), W: 7},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 1},
			{F: concrete.Node('C'), T: concrete.Node('D'), W: 14},
			{F: concrete.Node('C'), T: concrete.Node('E'), W: 8},
			{F: concrete.Node('D'), T: concrete.Node('E'), W: 9},
			{F: concrete.Node('D'), T: concrete.Node('F'), W: 3},
			{F: concrete.Node('E'), T: concrete.Node('F'), W: 2},
			{F: concrete.Node('G'), T: concrete.Node('H'), W: 15},
			{F: concrete.Node('G'), T: concrete.Node('I'), W: 5},
			{F: concrete.Node('G'), T: concrete.Node('J'), W: 19},
			{F: concrete.Node('G'), T: concrete.Node('K'), W: 10},
			{F: concrete.Node('H'), T: concrete.Node('J'), W: 17},
			{F: concrete.Node('I'), T: concrete.Node('K'), W: 11},
			{F: concrete.Node('J'), T: concrete.Node('K'), W: 16},
			{F: concrete.Node('J'), T: concrete.Node('L'), W: 4},
			{F: concrete.Node('K'), T: concrete.Node('L'), W: 12},
		},

		want: 65,
		treeEdges: []concrete.Edge{
			{F: concrete.Node('A'), T: concrete.Node('C'), W: 6},
			{F: concrete.Node('B'), T: concrete.Node('C'), W: 7},
			{F: concrete.Node('B'), T: concrete.Node('D'), W: 1},
			{F: concrete.Node('D'), T: concrete.Node('F'), W: 3},
			{F: concrete.Node('E'), T: concrete.Node('F'), W: 2},
			{F: concrete.Node('G'), T: concrete.Node('H'), W: 15},
			{F: concrete.Node('G'), T: concrete.Node('I'), W: 5},
			{F: concrete.Node('G'), T: concrete.Node('K'), W: 10},
			{F: concrete.Node('J'), T: concrete.Node('L'), W: 4},
			{F: concrete.Node('K'), T: concrete.Node('L'), W: 12},
		},
	},
}

func testMinumumSpanning(mst func(dst graph.MutableUndirected, g spanningGraph) float64, t *testing.T) {
	for _, test := range spanningTreeTests {
		g := test.graph()
		for _, e := range test.edges {
			g.SetEdge(e)
		}

		dst := concrete.NewGraph(0, math.Inf(1))
		w := mst(dst, g)
		if w != test.want {
			t.Errorf("unexpected minimum spanning tree weight for %q: got: %f want: %f",
				test.name, w, test.want)
		}
		var got float64
		for _, e := range dst.Edges() {
			got += e.Weight()
		}
		if got != test.want {
			t.Errorf("unexpected minimum spanning tree edge weight sum for %q: got: %f want: %f",
				test.name, got, test.want)
		}

		gotEdges := dst.Edges()
		if len(gotEdges) != len(test.treeEdges) {
			t.Errorf("unexpected number of spanning tree edges for %q: got: %d want: %d",
				test.name, len(gotEdges), len(test.treeEdges))
		}
		for _, e := range test.treeEdges {
			w, ok := dst.Weight(e.From(), e.To())
			if !ok {
				t.Errorf("spanning tree edge not found in graph for %q: %+v",
					test.name, e)
			}
			if w != e.Weight() {
				t.Errorf("unexpected spanning tree edge weight for %q: got: %f want: %f",
					test.name, w, e.Weight())
			}
		}
	}
}

func TestKruskal(t *testing.T) {
	testMinumumSpanning(func(dst graph.MutableUndirected, g spanningGraph) float64 {
		return Kruskal(dst, g)
	}, t)
}

func TestPrim(t *testing.T) {
	testMinumumSpanning(func(dst graph.MutableUndirected, g spanningGraph) float64 {
		return Prim(dst, g)
	}, t)
}
