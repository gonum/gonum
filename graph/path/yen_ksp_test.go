// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

var yenShortestPathTests = []struct {
	name  string
	graph func() graph.WeightedEdgeAdder
	edges []simple.WeightedEdge

	query     simple.Edge
	k         int
	wantPaths [][]int64
}{
	// Positive weighted graphs.
	// https://en.wikipedia.org/wiki/Yen%27s_algorithm
	{
		name:  "wikipedia example",
		graph: func() graph.WeightedEdgeAdder { return simple.NewWeightedDirectedGraph(0, math.Inf(1)) },
		edges: []simple.WeightedEdge{
			{F: simple.Node(0), T: simple.Node(1), W: 3},
			{F: simple.Node(0), T: simple.Node(2), W: 2},
			{F: simple.Node(2), T: simple.Node(1), W: 1},
			{F: simple.Node(1), T: simple.Node(3), W: 4},
			{F: simple.Node(2), T: simple.Node(3), W: 2},
			{F: simple.Node(2), T: simple.Node(4), W: 3},
			{F: simple.Node(3), T: simple.Node(4), W: 2},
			{F: simple.Node(3), T: simple.Node(5), W: 1},
			{F: simple.Node(4), T: simple.Node(5), W: 2},
		},
		query: simple.Edge{F: simple.Node(0), T: simple.Node(5)},
		k:     3,
		wantPaths: [][]int64{
			{0, 2, 3, 5},
			{0, 2, 4, 5},
			{0, 1, 3, 5},
		},
	},
	{
		name:  "1 Edge graph",
		graph: func() graph.WeightedEdgeAdder { return simple.NewWeightedDirectedGraph(0, math.Inf(1)) },
		edges: []simple.WeightedEdge{
			{F: simple.Node(0), T: simple.Node(1), W: 3},
		},
		query: simple.Edge{F: simple.Node(0), T: simple.Node(1)},
		k:     1,
		wantPaths: [][]int64{
			{0, 1},
		},
	},
	{
		name:      "Empty Graph",
		graph:     func() graph.WeightedEdgeAdder { return simple.NewWeightedDirectedGraph(0, math.Inf(1)) },
		edges:     []simple.WeightedEdge{},
		query:     simple.Edge{F: simple.Node(0), T: simple.Node(1)},
		k:         1,
		wantPaths: nil,
	},
	{
		name:  "N Star Graph",
		graph: func() graph.WeightedEdgeAdder { return simple.NewWeightedDirectedGraph(0, math.Inf(1)) },
		edges: []simple.WeightedEdge{
			{F: simple.Node(0), T: simple.Node(1), W: 3},
			{F: simple.Node(0), T: simple.Node(2), W: 3},
			{F: simple.Node(0), T: simple.Node(3), W: 3},
		},
		query: simple.Edge{F: simple.Node(0), T: simple.Node(1)},
		k:     1,
		wantPaths: [][]int64{
			{0, 1},
		},
	},
}

func toIntPath(nodePaths [][]graph.Node) [][]int64 {
	var paths [][]int64

	for _, nodes := range nodePaths {
		var path []int64
		for _, node := range nodes {
			path = append(path, node.ID())
		}
		paths = append(paths, path)
	}

	return paths
}

func TestYenKSP(t *testing.T) {
	for _, test := range yenShortestPathTests {
		g := test.graph()
		for _, e := range test.edges {
			g.SetWeightedEdge(e)
		}

		got := YenKShortestPath(g.(graph.Graph), test.k, test.query.From(), test.query.To())
		gotIds := toIntPath(got)

		if !reflect.DeepEqual(test.wantPaths, gotIds) {
			t.Errorf("unexpected result for %q:\ngot: %v\nwant:%v", test.name, gotIds, test.wantPaths)
		}
	}
}
