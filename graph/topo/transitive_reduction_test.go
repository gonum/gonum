// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

var transitiveReduceTests = []struct {
	name      string
	build     func() *simple.DirectedGraph
	wantPanic bool
	check     func(t *testing.T, before, after graph.Directed) // optional spot-checks
}{
	{
		name: "ChainWithShortcuts",
		build: func() *simple.DirectedGraph {
			const n = 8
			g := simple.NewDirectedGraph()
			// Chain edges i -> i+1.
			for i := 0; i < n-1; i++ {
				g.SetEdge(simple.Edge{F: simple.Node(int64(i)), T: simple.Node(int64(i + 1))})
			}
			// Add transitive shortcuts from 0 -> k for k >= 2.
			for k := 2; k < n; k++ {
				g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(int64(k))})
			}
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			// Spot-check: all 0->k (k>=2) should be removed.
			for k := int64(2); k < 8; k++ {
				if after.HasEdgeFromTo(0, k) {
					t.Errorf("expected redundant edge 0->%d to be removed", k)
				}
			}
		},
	},
	{
		name: "Diamond",
		build: func() *simple.DirectedGraph {
			// 1->2, 1->3, 2->4, 3->4, plus shortcut 1->4.
			g := simple.NewDirectedGraph()
			g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
			g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(3)})
			g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(4)})
			g.SetEdge(simple.Edge{F: simple.Node(3), T: simple.Node(4)})
			g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(4)}) // redundant
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if after.HasEdgeFromTo(1, 4) {
				t.Errorf("expected redundant edge 1->4 to be removed")
			}
		},
	},
	{
		name: "AlreadyReduced",
		build: func() *simple.DirectedGraph {
			// Simple chain: no redundant edges.
			g := simple.NewDirectedGraph()
			g.SetEdge(simple.Edge{F: simple.Node(10), T: simple.Node(11)})
			g.SetEdge(simple.Edge{F: simple.Node(11), T: simple.Node(12)})
			g.SetEdge(simple.Edge{F: simple.Node(12), T: simple.Node(13)})
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if edgeCount(after) != edgeCount(before) {
				t.Errorf("expected edge count unchanged: before=%d after=%d", edgeCount(before), edgeCount(after))
			}
		},
	},
	{
		name: "Disconnected",
		build: func() *simple.DirectedGraph {
			// Two components:
			// A: 0->1->2 plus 0->2 (redundant)
			// B: 10->11
			g := simple.NewDirectedGraph()
			g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(1)})
			g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
			g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(2)}) // redundant
			g.SetEdge(simple.Edge{F: simple.Node(10), T: simple.Node(11)})
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if after.HasEdgeFromTo(0, 2) {
				t.Errorf("expected redundant edge 0->2 to be removed")
			}
			if !after.HasEdgeFromTo(10, 11) {
				t.Errorf("expected edge 10->11 to remain")
			}
		},
	},
	{
		name: "EmptyGraph",
		build: func() *simple.DirectedGraph {
			return simple.NewDirectedGraph()
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if after.Nodes().Len() != 0 {
				t.Errorf("expected 0 nodes, got %d", after.Nodes().Len())
			}
			if edgeCount(after) != 0 {
				t.Errorf("expected 0 edges, got %d", edgeCount(after))
			}
		},
	},
	{
		name: "SingleNodeNoEdges",
		build: func() *simple.DirectedGraph {
			g := simple.NewDirectedGraph()
			g.AddNode(simple.Node(1))
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if edgeCount(after) != 0 {
				t.Errorf("expected 0 edges, got %d", edgeCount(after))
			}
		},
	},
	{
		name: "TwoNodesSingleEdge",
		build: func() *simple.DirectedGraph {
			g := simple.NewDirectedGraph()
			g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if !after.HasEdgeFromTo(1, 2) {
				t.Errorf("expected edge 1->2 to remain")
			}
			if edgeCount(after) != 1 {
				t.Errorf("expected 1 edge, got %d", edgeCount(after))
			}
		},
	},
	{
		name: "CyclePanics",
		build: func() *simple.DirectedGraph {
			g := simple.NewDirectedGraph()
			g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
			g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(3)})
			g.SetEdge(simple.Edge{F: simple.Node(3), T: simple.Node(1)}) // cycle
			return g
		},
		wantPanic: true,
	},
	{
		name: "StarNoTransitiveEdges",
		build: func() *simple.DirectedGraph {
			// 0 -> {1,2,3,4} and no other edges => nothing is redundant.
			g := simple.NewDirectedGraph()
			for i := int64(1); i <= 4; i++ {
				g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(i)})
			}
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if edgeCount(after) != edgeCount(before) {
				t.Errorf("expected no edge removals: before=%d after=%d", edgeCount(before), edgeCount(after))
			}
		},
	},
	{
		name: "LayeredDAGManyRedundantEdges",
		build: func() *simple.DirectedGraph {
			// Layers:
			// 0 -> (1,2,3)
			// (1,2,3) -> (4,5,6)
			// (4,5,6) -> 7
			// plus redundant shortcuts: 0->(4,5,6), 0->7, (1,2,3)->7
			g := simple.NewDirectedGraph()
			// 0 -> 1,2,3
			for _, v := range []int64{1, 2, 3} {
				g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(v)})
			}
			// 1,2,3 -> 4,5,6
			for _, u := range []int64{1, 2, 3} {
				for _, v := range []int64{4, 5, 6} {
					g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
				}
			}
			// 4,5,6 -> 7
			for _, u := range []int64{4, 5, 6} {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(7)})
			}
			// Redundant shortcuts
			for _, v := range []int64{4, 5, 6, 7} {
				g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(v)})
			}
			for _, u := range []int64{1, 2, 3} {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(7)})
			}
			return g
		},
		check: func(t *testing.T, before, after graph.Directed) {
			if after.HasEdgeFromTo(0, 7) {
				t.Errorf("expected redundant edge 0->7 to be removed")
			}
			for _, v := range []int64{4, 5, 6} {
				if after.HasEdgeFromTo(0, v) {
					t.Errorf("expected redundant edge 0->%d to be removed", v)
				}
			}
		},
	},
	{
		name: "RandomSmallDAG",
		build: func() *simple.DirectedGraph {
			// Deterministic "random" DAG:
			// Add edge i->j for i<j if (i*17 + j*31) % 5 == 0, and then add a backbone.
			const n = 12
			g := simple.NewDirectedGraph()
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					if (i*17+j*31)%5 == 0 {
						g.SetEdge(simple.Edge{F: simple.Node(int64(i)), T: simple.Node(int64(j))})
					}
				}
			}
			// Ensure connected-ish backbone.
			for i := 0; i < n-1; i++ {
				g.SetEdge(simple.Edge{F: simple.Node(int64(i)), T: simple.Node(int64(i + 1))})
			}
			return g
		},
	},
}

func TestTransitiveReduce(t *testing.T) {
	for _, test := range transitiveReduceTests {
		t.Run(test.name, func(t *testing.T) {
			orig := test.build()
			before := cloneDirected(orig)
			after := cloneDirected(orig)

			if test.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("expected panic, got none")
					}
				}()
			}
			TransitiveReduce(after)
			if test.wantPanic {
				return // we expected to panic; if we got here, recover didn't trigger
			}

			checkEdgesSubset(t, before, after)
			checkReachabilityPreserved(t, before, after)
			checkMinimal(t, after)

			if test.check != nil {
				test.check(t, before, after)
			}
		})
	}
}

func checkEdgesSubset(t *testing.T, original, reduced graph.Directed) {
	t.Helper()

	it := reduced.Nodes()
	for it.Next() {
		u := it.Node().ID()
		from := reduced.From(u)
		for from.Next() {
			v := from.Node().ID()
			if !original.HasEdgeFromTo(u, v) {
				t.Errorf("reduced graph contains edge %d->%d not present in original", u, v)
			}
		}
	}
}

func checkReachabilityPreserved(t *testing.T, before, after graph.Directed) {
	t.Helper()

	nodes := graph.NodesOf(before.Nodes())
	for _, u := range nodes {
		for _, v := range nodes {
			got := PathExistsIn(after, after.Node(u.ID()), after.Node(v.ID()))
			want := PathExistsIn(before, before.Node(u.ID()), before.Node(v.ID()))
			if got != want {
				t.Errorf("reachability changed for (%d -> %d): before=%v after=%v", u.ID(), v.ID(), want, got)
			}
		}
	}
}

func checkMinimal(t *testing.T, reduced graph.Directed) {
	t.Helper()

	// For each edge u->v, removing it must break reachability u=>v.
	it := reduced.Nodes()
	for it.Next() {
		u := it.Node().ID()
		from := reduced.From(u)
		for from.Next() {
			v := from.Node().ID()

			g2 := cloneDirected(reduced)
			g2.RemoveEdge(u, v)

			u2 := g2.Node(u)
			v2 := g2.Node(v)
			if u2 == nil || v2 == nil {
				t.Fatalf("missing node after clone: u=%d v=%d", u, v)
			}
			if PathExistsIn(g2, u2, v2) {
				t.Errorf("edge %d->%d is redundant: path still exists after removal", u, v)
			}
		}
	}
}

func cloneDirected(g graph.Directed) *simple.DirectedGraph {
	ng := simple.NewDirectedGraph()
	graph.Copy(ng, g)
	return ng
}

func edgeCount(g graph.Directed) int {
	n := 0
	it := g.Nodes()
	for it.Next() {
		u := it.Node().ID()
		from := g.From(u)
		for from.Next() {
			n++
		}
	}
	return n
}
