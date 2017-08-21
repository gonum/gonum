// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/graph"
)

var (
	directedGraph = (*DirectedGraph)(nil)

	_ graph.Graph            = directedGraph
	_ graph.Directed         = directedGraph
	_ graph.WeightedDirected = directedGraph
)

// Tests Issue #27
func TestEdgeOvercounting(t *testing.T) {
	g := generateDummyGraph()

	if neigh := g.From(Node(Node(2))); len(neigh) != 2 {
		t.Errorf("Node 2 has incorrect number of neighbors got neighbors %v (count %d), expected 2 neighbors {0,1}", neigh, len(neigh))
	}
}

func generateDummyGraph() *DirectedGraph {
	nodes := [4]struct{ srcID, targetID int }{
		{2, 1},
		{1, 0},
		{2, 0},
		{0, 2},
	}

	g := NewDirectedGraph(0, math.Inf(1))

	for _, n := range nodes {
		g.SetEdge(Edge{F: Node(n.srcID), T: Node(n.targetID), W: 1})
	}

	return g
}

// Test for issue #123 https://github.com/gonum/graph/issues/123
func TestIssue123DirectedGraph(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	g := NewDirectedGraph(0, math.Inf(1))

	n0 := g.NewNode()
	g.AddNode(n0)

	n1 := g.NewNode()
	g.AddNode(n1)

	g.RemoveNode(n0)

	n2 := g.NewNode()
	g.AddNode(n2)
}
