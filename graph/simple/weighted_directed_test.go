// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple_test

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/testgraph"
)

func weightedDirectedBuilder(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool) {
	seen := make(set.Nodes)
	dg := simple.NewWeightedDirectedGraph(self, absent)
	for _, n := range nodes {
		seen.Add(n)
		dg.AddNode(n)
	}
	for _, edge := range edges {
		if edge.From().ID() == edge.To().ID() {
			continue
		}
		f := dg.Node(edge.From().ID())
		if f == nil {
			f = edge.From()
		}
		t := dg.Node(edge.To().ID())
		if t == nil {
			t = edge.To()
		}
		ce := simple.WeightedEdge{F: f, T: t, W: edge.Weight()}
		seen.Add(ce.F)
		seen.Add(ce.T)
		e = append(e, ce)
		dg.SetWeightedEdge(ce)
	}
	if len(e) == 0 && len(edges) != 0 {
		return nil, nil, nil, math.NaN(), math.NaN(), false
	}
	if len(seen) != 0 {
		n = make([]graph.Node, 0, len(seen))
	}
	for _, sn := range seen {
		n = append(n, sn)
	}
	return dg, n, e, self, absent, true
}

func TestWeightedDirected(t *testing.T) {
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, weightedDirectedBuilder)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, weightedDirectedBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, weightedDirectedBuilder)
	})
	t.Run("ReturnAllEdges", func(t *testing.T) {
		testgraph.ReturnAllEdges(t, weightedDirectedBuilder)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, weightedDirectedBuilder)
	})
	t.Run("ReturnAllWeightedEdges", func(t *testing.T) {
		testgraph.ReturnAllWeightedEdges(t, weightedDirectedBuilder)
	})
	t.Run("ReturnEdgeSlice", func(t *testing.T) {
		testgraph.ReturnEdgeSlice(t, weightedDirectedBuilder)
	})
	t.Run("ReturnWeightedEdgeSlice", func(t *testing.T) {
		testgraph.ReturnWeightedEdgeSlice(t, weightedDirectedBuilder)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, weightedDirectedBuilder)
	})
	t.Run("Weight", func(t *testing.T) {
		testgraph.Weight(t, weightedDirectedBuilder)
	})
}

// Tests Issue #27
func TestWeightedEdgeOvercounting(t *testing.T) {
	g := generateDummyGraph()

	if neigh := graph.NodesOf(g.From(int64(2))); len(neigh) != 2 {
		t.Errorf("Node 2 has incorrect number of neighbors got neighbors %v (count %d), expected 2 neighbors {0,1}", neigh, len(neigh))
	}
}

func generateDummyWeightedGraph() *simple.WeightedDirectedGraph {
	nodes := [4]struct{ srcID, targetID int }{
		{2, 1},
		{1, 0},
		{2, 0},
		{0, 2},
	}

	g := simple.NewWeightedDirectedGraph(0, math.Inf(1))

	for _, n := range nodes {
		g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(n.srcID), T: simple.Node(n.targetID), W: 1})
	}

	return g
}

// Test for issue #123 https://github.com/gonum/graph/issues/123
func TestIssue123WeightedDirectedGraph(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	g := simple.NewWeightedDirectedGraph(0, math.Inf(1))

	n0 := g.NewNode()
	g.AddNode(n0)

	n1 := g.NewNode()
	g.AddNode(n1)

	g.RemoveNode(n0.ID())

	n2 := g.NewNode()
	g.AddNode(n2)
}
