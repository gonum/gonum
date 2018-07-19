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


var YenShortestPathTests = []struct {
	Name              string
	Graph             func() graph.WeightedEdgeAdder
	Edges             []simple.WeightedEdge

	Query         simple.Edge
	K int
	WantPaths     [][]int64
}{
	// Positive weighted graphs.
	{
		Name:  "simple graph",
		Graph: func() graph.WeightedEdgeAdder { return simple.NewWeightedDirectedGraph(0, math.Inf(1)) },
		Edges: []simple.WeightedEdge{
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
		Query:  simple.Edge{F: simple.Node(0), T: simple.Node(5)},
		K: 3,
		WantPaths: [][]int64{
			{0, 2, 3, 5},
			{0, 2, 4, 5},
			{0, 1, 3, 5},
		},	
	},
}

func toIntPath(nodePaths [][]graph.Node) [][]int64 {
	paths := make([][]int64, 0)
	for _, nodes := range nodePaths {
		path := make([]int64, 0)
		for _, node := range nodes {
			path = append(path, node.ID())
		}
		paths = append(paths, path)
	}

	return paths
}

func TestYenKSP(t *testing.T) {
	for _, test := range YenShortestPathTests {
		g := test.Graph()
		for _, e := range test.Edges {
			g.SetWeightedEdge(e)
		}

		paths := toIntPath(YenKShortestPath(g.(graph.Graph), test.K, test.Query.From(), test.Query.To()))

		expected := test.WantPaths
		
		if len(expected) != len(paths) {
			t.Errorf("ERROR: expected %d paths, got %d paths", len(expected), len(paths))
		} else {
			for i := 0; i < len(paths); i++ {
				if !reflect.DeepEqual(expected[i], paths[i]) {
					t.Errorf("ERROR: path #%d expected: %d, got: %d", i+1, expected[i], paths[i])
				}
			}
		}
	}
}
