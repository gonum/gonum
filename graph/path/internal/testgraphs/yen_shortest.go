// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testgraphs

import (
	"math"
	
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func init() {
	//for _, test := range YenShortestPathTests {
		//if len(test.WantPaths) !=  {
		//	panic(fmt.Sprintf("%q: bad shortest path test: non-unique paths marked unique", test.Name))
		//}
	//}
}

// ShortestPathTests are graphs used to test the static shortest path routines in path: BellmanFord,
// DijkstraAllPaths, DijkstraFrom, FloydWarshall and Johnson, and the static degenerate case for the
// dynamic shortest path routine in path/dynamic: DStarLite.
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
