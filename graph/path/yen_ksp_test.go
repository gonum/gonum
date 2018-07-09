// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph"
)


func TestYenKSP(t *testing.T) {
	g := simple.NewWeightedDirectedGraph(0, math.Inf(1))

	edges := []simple.WeightedEdge{
		{F: simple.Node(0), T: simple.Node(1), W: 3},
		{F: simple.Node(0), T: simple.Node(2), W: 2},
		{F: simple.Node(2), T: simple.Node(1), W: 1},
		{F: simple.Node(1), T: simple.Node(3), W: 4},
		{F: simple.Node(2), T: simple.Node(3), W: 2},
		{F: simple.Node(2), T: simple.Node(4), W: 3},
		{F: simple.Node(3), T: simple.Node(4), W: 2},
		{F: simple.Node(3), T: simple.Node(5), W: 1},
		{F: simple.Node(4), T: simple.Node(5), W: 2},
	}

	for _, edge := range edges {
		g.SetWeightedEdge(edge)
	}

	shortests := YenKShortestPath(g, 3, simple.Node(0), simple.Node(5))
	expected := make([][]graph.Node, 3)

	expected[0] = make([]graph.Node, 4)
	expected[0][0] = simple.Node(0)
	expected[0][1] = simple.Node(2)
	expected[0][2] = simple.Node(3)
	expected[0][3] = simple.Node(5)

	expected[1] = make([]graph.Node, 4)
	expected[1][0] = simple.Node(0)
	expected[1][1] = simple.Node(2)
	expected[1][2] = simple.Node(4)
	expected[1][3] = simple.Node(5)

	expected[2] = make([]graph.Node, 4)
	expected[2][0] = simple.Node(0)
	expected[2][1] = simple.Node(1)
	expected[2][2] = simple.Node(3)
	expected[2][3] = simple.Node(5)

	for i, sp := range shortests {
		e := expected[i]
		for n, p := range sp {
			if (e[n].ID() != p.ID()) {
				t.Errorf("ERROR: path #%d expected: %d, got: %d", i+1, e[n].ID(), p.ID())
			}
		}
	}
}
