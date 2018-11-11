// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simple_test

import (
	"math"
	"sort"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/testgraph"
)

func isZeroContiguousSet(nodes []graph.Node) bool {
	t := make([]graph.Node, len(nodes))
	copy(t, nodes)
	nodes = t
	sort.Sort(ordered.ByID(nodes))
	for i, n := range nodes {
		if int64(i) != n.ID() {
			return false
		}
	}
	return true
}

func directedMatrixBuilder(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool) {
	if len(nodes) == 0 {
		return
	}
	if !isZeroContiguousSet(nodes) {
		return
	}
	seen := make(set.Nodes)
	dg := simple.NewDirectedMatrix(len(nodes), absent, self, absent)
	for i := range nodes {
		seen.Add(simple.Node(i))
	}
	for _, edge := range edges {
		if edge.From().ID() == edge.To().ID() {
			continue
		}
		if !seen.Has(edge.From()) || !seen.Has(edge.To()) {
			continue
		}
		ce := simple.WeightedEdge{F: dg.Node(edge.From().ID()), T: dg.Node(edge.To().ID()), W: edge.Weight()}
		e = append(e, ce)
		dg.SetWeightedEdge(ce)
	}
	if len(e) == 0 && len(edges) != 0 {
		return nil, nil, nil, math.NaN(), math.NaN(), false
	}
	n = make([]graph.Node, 0, len(seen))
	for _, sn := range seen {
		n = append(n, sn)
	}
	return dg, n, e, self, absent, true
}

func TestDirectedMatrix(t *testing.T) {
	t.Run("AdjacencyMatrix", func(t *testing.T) {
		testgraph.AdjacencyMatrix(t, directedMatrixBuilder)
	})
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, directedMatrixBuilder)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, directedMatrixBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, directedMatrixBuilder)
	})
	t.Run("ReturnAllEdges", func(t *testing.T) {
		testgraph.ReturnAllEdges(t, directedMatrixBuilder)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, directedMatrixBuilder)
	})
	t.Run("ReturnAllWeightedEdges", func(t *testing.T) {
		testgraph.ReturnAllWeightedEdges(t, directedMatrixBuilder)
	})
	t.Run("ReturnEdgeSlice", func(t *testing.T) {
		testgraph.ReturnEdgeSlice(t, directedMatrixBuilder)
	})
	t.Run("ReturnWeightedEdgeSlice", func(t *testing.T) {
		testgraph.ReturnWeightedEdgeSlice(t, directedMatrixBuilder)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, directedMatrixBuilder)
	})
	t.Run("Weight", func(t *testing.T) {
		testgraph.Weight(t, directedMatrixBuilder)
	})
}

func directedMatrixFromBuilder(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool) {
	if len(nodes) == 0 {
		return
	}
	if !isZeroContiguousSet(nodes) {
		return
	}
	seen := make(set.Nodes)
	dg := simple.NewDirectedMatrixFrom(nodes, absent, self, absent)
	for _, n := range nodes {
		seen.Add(n)
	}
	for _, edge := range edges {
		if edge.From().ID() == edge.To().ID() {
			continue
		}
		if !seen.Has(edge.From()) || !seen.Has(edge.To()) {
			continue
		}
		ce := simple.WeightedEdge{F: dg.Node(edge.From().ID()), T: dg.Node(edge.To().ID()), W: edge.Weight()}
		e = append(e, ce)
		dg.SetWeightedEdge(ce)
	}
	if len(e) == 0 && len(edges) != 0 {
		return nil, nil, nil, math.NaN(), math.NaN(), false
	}
	n = make([]graph.Node, 0, len(seen))
	for _, sn := range seen {
		n = append(n, sn)
	}
	return dg, n, e, self, absent, true
}

func TestDirectedMatrixFrom(t *testing.T) {
	t.Run("AdjacencyMatrix", func(t *testing.T) {
		testgraph.AdjacencyMatrix(t, directedMatrixFromBuilder)
	})
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, directedMatrixFromBuilder)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnAllEdges", func(t *testing.T) {
		testgraph.ReturnAllEdges(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnAllWeightedEdges", func(t *testing.T) {
		testgraph.ReturnAllWeightedEdges(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnEdgeSlice", func(t *testing.T) {
		testgraph.ReturnEdgeSlice(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnWeightedEdgeSlice", func(t *testing.T) {
		testgraph.ReturnWeightedEdgeSlice(t, directedMatrixFromBuilder)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, directedMatrixFromBuilder)
	})
	t.Run("Weight", func(t *testing.T) {
		testgraph.Weight(t, directedMatrixFromBuilder)
	})
}

func undirectedMatrixBuilder(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool) {
	if len(nodes) == 0 {
		return
	}
	if !isZeroContiguousSet(nodes) {
		return
	}
	seen := make(set.Nodes)
	dg := simple.NewUndirectedMatrix(len(nodes), absent, self, absent)
	for i := range nodes {
		seen.Add(simple.Node(i))
	}
	for _, edge := range edges {
		if edge.From().ID() == edge.To().ID() {
			continue
		}
		if !seen.Has(edge.From()) || !seen.Has(edge.To()) {
			continue
		}
		ce := simple.WeightedEdge{F: dg.Node(edge.From().ID()), T: dg.Node(edge.To().ID()), W: edge.Weight()}
		e = append(e, ce)
		dg.SetWeightedEdge(ce)
	}
	if len(e) == 0 && len(edges) != 0 {
		return nil, nil, nil, math.NaN(), math.NaN(), false
	}
	n = make([]graph.Node, 0, len(seen))
	for _, sn := range seen {
		n = append(n, sn)
	}
	return dg, n, e, self, absent, true
}

func TestUnirectedMatrix(t *testing.T) {
	t.Run("AdjacencyMatrix", func(t *testing.T) {
		testgraph.AdjacencyMatrix(t, undirectedMatrixBuilder)
	})
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, undirectedMatrixBuilder)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnAllEdges", func(t *testing.T) {
		testgraph.ReturnAllEdges(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnAllWeightedEdges", func(t *testing.T) {
		testgraph.ReturnAllWeightedEdges(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnEdgeSlice", func(t *testing.T) {
		testgraph.ReturnEdgeSlice(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnWeightedEdgeSlice", func(t *testing.T) {
		testgraph.ReturnWeightedEdgeSlice(t, undirectedMatrixBuilder)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, undirectedMatrixBuilder)
	})
	t.Run("Weight", func(t *testing.T) {
		testgraph.Weight(t, undirectedMatrixBuilder)
	})
}

func undirectedMatrixFromBuilder(nodes []graph.Node, edges []graph.WeightedLine, self, absent float64) (g graph.Graph, n []graph.Node, e []graph.Edge, s, a float64, ok bool) {
	if len(nodes) == 0 {
		return
	}
	if !isZeroContiguousSet(nodes) {
		return
	}
	seen := make(set.Nodes)
	dg := simple.NewUndirectedMatrixFrom(nodes, absent, self, absent)
	for _, n := range nodes {
		seen.Add(n)
	}
	for _, edge := range edges {
		if edge.From().ID() == edge.To().ID() {
			continue
		}
		if !seen.Has(edge.From()) || !seen.Has(edge.To()) {
			continue
		}
		ce := simple.WeightedEdge{F: dg.Node(edge.From().ID()), T: dg.Node(edge.To().ID()), W: edge.Weight()}
		e = append(e, ce)
		dg.SetWeightedEdge(ce)
	}
	if len(e) == 0 && len(edges) != 0 {
		return nil, nil, nil, math.NaN(), math.NaN(), false
	}
	n = make([]graph.Node, 0, len(seen))
	for _, sn := range seen {
		n = append(n, sn)
	}
	return dg, n, e, self, absent, true
}

func TestUndirectedMatrixFrom(t *testing.T) {
	t.Run("AdjacencyMatrix", func(t *testing.T) {
		testgraph.AdjacencyMatrix(t, undirectedMatrixFromBuilder)
	})
	t.Run("EdgeExistence", func(t *testing.T) {
		testgraph.EdgeExistence(t, undirectedMatrixFromBuilder)
	})
	t.Run("NodeExistence", func(t *testing.T) {
		testgraph.NodeExistence(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnAdjacentNodes", func(t *testing.T) {
		testgraph.ReturnAdjacentNodes(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnAllEdges", func(t *testing.T) {
		testgraph.ReturnAllEdges(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnAllNodes", func(t *testing.T) {
		testgraph.ReturnAllNodes(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnAllWeightedEdges", func(t *testing.T) {
		testgraph.ReturnAllWeightedEdges(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnEdgeSlice", func(t *testing.T) {
		testgraph.ReturnEdgeSlice(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnWeightedEdgeSlice", func(t *testing.T) {
		testgraph.ReturnWeightedEdgeSlice(t, undirectedMatrixFromBuilder)
	})
	t.Run("ReturnNodeSlice", func(t *testing.T) {
		testgraph.ReturnNodeSlice(t, undirectedMatrixFromBuilder)
	})
	t.Run("Weight", func(t *testing.T) {
		testgraph.Weight(t, undirectedMatrixFromBuilder)
	})
}

func TestBasicDenseImpassable(t *testing.T) {
	dg := simple.NewUndirectedMatrix(5, math.Inf(1), 0, math.Inf(1))
	if dg == nil {
		t.Fatal("Directed graph could not be made")
	}

	for i := 0; i < 5; i++ {
		if dg.Node(int64(i)) == nil {
			t.Errorf("Node that should exist doesn't: %d", i)
		}

		if degree := dg.From(int64(i)).Len(); degree != 0 {
			t.Errorf("Node in impassable graph has a neighbor. Node: %d Degree: %d", i, degree)
		}
	}

	for i := 5; i < 10; i++ {
		if dg.Node(int64(i)) != nil {
			t.Errorf("Node exists that shouldn't: %d", i)
		}
	}
}

func TestBasicDensePassable(t *testing.T) {
	dg := simple.NewUndirectedMatrix(5, 1, 0, math.Inf(1))
	if dg == nil {
		t.Fatal("Directed graph could not be made")
	}

	for i := 0; i < 5; i++ {
		if dg.Node(int64(i)) == nil {
			t.Errorf("Node that should exist doesn't: %d", i)
		}

		if degree := dg.From(int64(i)).Len(); degree != 4 {
			t.Errorf("Node in passable graph missing neighbors. Node: %d Degree: %d", i, degree)
		}
	}

	for i := 5; i < 10; i++ {
		if dg.Node(int64(i)) != nil {
			t.Errorf("Node exists that shouldn't: %d", i)
		}
	}
}

func TestDirectedDenseAddRemove(t *testing.T) {
	dg := simple.NewDirectedMatrix(10, math.Inf(1), 0, math.Inf(1))
	dg.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(0), T: simple.Node(2), W: 1})

	if neighbors := graph.NodesOf(dg.From(int64(0))); len(neighbors) != 1 || neighbors[0].ID() != 2 ||
		dg.Edge(int64(0), int64(2)) == nil {
		t.Errorf("Adding edge didn't create successor")
	}

	dg.RemoveEdge(int64(0), int64(2))

	if neighbors := graph.NodesOf(dg.From(int64(0))); len(neighbors) != 0 || dg.Edge(int64(0), int64(2)) != nil {
		t.Errorf("Removing edge didn't properly remove successor")
	}

	if neighbors := graph.NodesOf(dg.To(int64(2))); len(neighbors) != 0 || dg.Edge(int64(0), int64(2)) != nil {
		t.Errorf("Removing directed edge wrongly kept predecessor")
	}

	dg.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(0), T: simple.Node(2), W: 2})
	// I figure we've torture tested From/To at this point
	// so we'll just use the bool functions now
	if dg.Edge(int64(0), int64(2)) == nil {
		t.Fatal("Adding directed edge didn't change successor back")
	}
	c1, _ := dg.Weight(int64(2), int64(0))
	c2, _ := dg.Weight(int64(0), int64(2))
	if c1 == c2 {
		t.Error("Adding directed edge affected cost in undirected manner")
	}
}

func TestUndirectedDenseAddRemove(t *testing.T) {
	dg := simple.NewUndirectedMatrix(10, math.Inf(1), 0, math.Inf(1))
	dg.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(2)})

	if neighbors := graph.NodesOf(dg.From(int64(0))); len(neighbors) != 1 || neighbors[0].ID() != 2 ||
		dg.EdgeBetween(int64(0), int64(2)) == nil {
		t.Errorf("Couldn't add neighbor")
	}

	if neighbors := graph.NodesOf(dg.From(int64(2))); len(neighbors) != 1 || neighbors[0].ID() != 0 ||
		dg.EdgeBetween(int64(2), int64(0)) == nil {
		t.Errorf("Adding an undirected neighbor didn't add it reciprocally")
	}
}

func TestDenseLists(t *testing.T) {
	dg := simple.NewDirectedMatrix(15, 1, 0, math.Inf(1))
	nodes := graph.NodesOf(dg.Nodes())

	if len(nodes) != 15 {
		t.Fatalf("Wrong number of nodes: got:%v want:%v", len(nodes), 15)
	}

	sort.Sort(ordered.ByID(nodes))

	for i, node := range graph.NodesOf(dg.Nodes()) {
		if int64(i) != node.ID() {
			t.Errorf("Node list doesn't return properly id'd nodes")
		}
	}

	edges := graph.EdgesOf(dg.Edges())
	if len(edges) != 15*14 {
		t.Errorf("Improper number of edges for passable dense graph")
	}

	dg.RemoveEdge(int64(12), int64(11))
	edges = graph.EdgesOf(dg.Edges())
	if len(edges) != (15*14)-1 {
		t.Errorf("Removing edge didn't affect edge listing properly")
	}
}
