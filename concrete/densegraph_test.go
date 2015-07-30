// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete

import (
	"math"
	"sort"
	"testing"

	"github.com/gonum/graph"
)

var (
	_ graph.Graph    = (*UndirectedDenseGraph)(nil)
	_ graph.Directed = (*DirectedDenseGraph)(nil)
)

func TestBasicDenseImpassable(t *testing.T) {
	dg := NewUndirectedDenseGraph(5, false, math.Inf(1))
	if dg == nil {
		t.Fatal("Directed graph could not be made")
	}

	for i := 0; i < 5; i++ {
		if !dg.Has(Node(i)) {
			t.Errorf("Node that should exist doesn't: %d", i)
		}

		if degree := dg.Degree(Node(i)); degree != 0 {
			t.Errorf("Node in impassable graph has a neighbor. Node: %d Degree: %d", i, degree)
		}
	}

	for i := 5; i < 10; i++ {
		if dg.Has(Node(i)) {
			t.Errorf("Node exists that shouldn't: %d", i)
		}
	}
}

func TestBasicDensePassable(t *testing.T) {
	dg := NewUndirectedDenseGraph(5, true, math.Inf(1))
	if dg == nil {
		t.Fatal("Directed graph could not be made")
	}

	for i := 0; i < 5; i++ {
		if !dg.Has(Node(i)) {
			t.Errorf("Node that should exist doesn't: %d", i)
		}

		if degree := dg.Degree(Node(i)); degree != 4 {
			t.Errorf("Node in passable graph missing neighbors. Node: %d Degree: %d", i, degree)
		}
	}

	for i := 5; i < 10; i++ {
		if dg.Has(Node(i)) {
			t.Errorf("Node exists that shouldn't: %d", i)
		}
	}
}

func TestDirectedDenseAddRemove(t *testing.T) {
	dg := NewDirectedDenseGraph(10, false, math.Inf(1))
	dg.SetEdgeWeight(Edge{Node(0), Node(2)}, 1)

	if neighbors := dg.From(Node(0)); len(neighbors) != 1 || neighbors[0].ID() != 2 ||
		dg.Edge(Node(0), Node(2)) == nil {
		t.Errorf("Adding edge didn't create successor")
	}

	dg.RemoveEdge(Edge{Node(0), Node(2)})

	if neighbors := dg.From(Node(0)); len(neighbors) != 0 || dg.Edge(Node(0), Node(2)) != nil {
		t.Errorf("Removing edge didn't properly remove successor")
	}

	if neighbors := dg.To(Node(2)); len(neighbors) != 0 || dg.Edge(Node(0), Node(2)) != nil {
		t.Errorf("Removing directed edge wrongly kept predecessor")
	}

	dg.SetEdgeWeight(Edge{Node(0), Node(2)}, 2)
	// I figure we've torture tested From/To at this point
	// so we'll just use the bool functions now
	if dg.Edge(Node(0), Node(2)) == nil {
		t.Error("Adding directed edge didn't change successor back")
	} else if c1, c2 := dg.Weight(Edge{Node(2), Node(0)}), dg.Weight(Edge{Node(0), Node(2)}); math.Abs(c1-c2) < .000001 {
		t.Error("Adding directed edge affected cost in undirected manner")
	}
}

func TestUndirectedDenseAddRemove(t *testing.T) {
	dg := NewUndirectedDenseGraph(10, false, math.Inf(1))
	dg.SetEdgeWeight(Edge{Node(0), Node(2)}, 1)

	if neighbors := dg.From(Node(0)); len(neighbors) != 1 || neighbors[0].ID() != 2 ||
		dg.EdgeBetween(Node(0), Node(2)) == nil {
		t.Errorf("Couldn't add neighbor")
	}

	if neighbors := dg.From(Node(2)); len(neighbors) != 1 || neighbors[0].ID() != 0 ||
		dg.EdgeBetween(Node(2), Node(0)) == nil {
		t.Errorf("Adding an undirected neighbor didn't add it reciprocally")
	}
}

type byID []graph.Node

func (n byID) Len() int           { return len(n) }
func (n byID) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byID) Less(i, j int) bool { return n[i].ID() < n[j].ID() }

func TestDenseLists(t *testing.T) {
	dg := NewDirectedDenseGraph(15, true, math.Inf(1))
	nodes := dg.Nodes()

	if len(nodes) != 15 {
		t.Fatalf("Wrong number of nodes")
	}

	sort.Sort(byID(nodes))

	for i, node := range dg.Nodes() {
		if i != node.ID() {
			t.Errorf("Node list doesn't return properly id'd nodes")
		}
	}

	edges := dg.Edges()
	if len(edges) != 15*14 {
		t.Errorf("Improper number of edges for passable dense graph")
	}

	dg.RemoveEdge(Edge{Node(12), Node(11)})
	edges = dg.Edges()
	if len(edges) != (15*14)-1 {
		t.Errorf("Removing edge didn't affect edge listing properly")
	}
}
