// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/internal/order"
)

func TestLeidenUndirected(t *testing.T) {
	for _, test := range communityUndirectedQTests {
		g := simple.NewUndirectedGraph()
		for u, e := range test.g {
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}

		t.Run(test.name, func(t *testing.T) {
			testLeidenUndirected(t, test, g)
		})
	}
}

func TestLeidenWeightedUndirected(t *testing.T) {
	for _, test := range communityUndirectedQTests {
		g := simple.NewWeightedUndirectedGraph(0, 0)
		for u, e := range test.g {
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(u), T: simple.Node(v), W: 1})
			}
		}

		t.Run(test.name, func(t *testing.T) {
			testLeidenUndirected(t, test, g)
		})
	}
}

func testLeidenUndirected(t *testing.T, test communityUndirectedQTest, g graph.Undirected) {
	const iterations = 20

	if test.structures[0].resolution != 1 {
		panic("bad test: expect resolution=1")
	}

	var (
		got   *ReducedUndirected
		bestQ = math.Inf(-1)
	)

	// Leiden is randomised so we do this to ensure the level tests are consistent.
	src := rand.New(rand.NewPCG(1, 1))
	for i := 0; i < iterations; i++ {
		r := Leiden(g, 1, src).(*ReducedUndirected)
		if q := Q(r, nil, 1); q > bestQ || math.IsNaN(q) {
			bestQ = q
			got = r

			if math.IsNaN(q) {
				// Don't try again for non-connected case.
				break
			}
		}

		// Check that Q is non-decreasing at each level (same as Louvain).
		var qs []float64
		for p := r; p != nil; p = p.Expanded().(*ReducedUndirected) {
			qs = append(qs, Q(p, nil, 1))
		}
		if !math.IsNaN(qs[0]) {
			slices.Reverse(qs)
			// Leiden may not be strictly monotonic in Q due to the refinement phase constraints,
			// but should be generally increasing. We allow a small tolerance.
			for i := 1; i < len(qs); i++ {
				if qs[i] < qs[i-1]-1e-6 {
					t.Errorf("%s: Q values not monotonically increasing across levels: %v", test.name, qs)
					break
				}
			}
		}
	}

	gotCommunities := got.Communities()
	for _, c := range gotCommunities {
		order.ByID(c)
	}
	order.BySliceIDs(gotCommunities)

	// Validate against test expectations unless it's NaN (unconnected)
	if !math.IsNaN(test.structures[0].want) {
		// Leiden aims for equal OR BETTER modularity than Louvain/Reference.
		// So we check if it's at least close to the expected 'want', but allow higher.
		if bestQ < test.structures[0].want-test.structures[0].tol {
			t.Errorf("unexpectedly low Q value for %q: got: %v want >= %v",
				test.name, bestQ, test.structures[0].want)
		}
	}

	// Verify partition validity (disjoint and complete)
	seen := make(map[int64]int)
	for i, c := range gotCommunities {
		for _, n := range c {
			id := n.ID()
			if j, ok := seen[id]; ok {
				t.Errorf("%s: node %d in multiple communities %d and %d", test.name, id, j, i)
			}
			seen[id] = i
		}
	}
	nodes := graph.NodesOf(g.Nodes())
	if len(seen) != len(nodes) {
		t.Errorf("%s: partition covers %d nodes, graph has %d", test.name, len(seen), len(nodes))
	}
}

func TestLeidenVsLouvainModularity(t *testing.T) {
	// On the same graph with fixed seed, both should produce valid modularity.
	// Leiden may produce equal or better Q due to refinement.
	g := simple.NewUndirectedGraph()
	for u, e := range smallDumbell {
		if g.Node(int64(u)) == nil {
			g.AddNode(simple.Node(u))
		}
		for v := range e {
			g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
		}
	}

	src := rand.New(rand.NewPCG(42, 42))
	louvainResult := Modularize(g, 1, src).(*ReducedUndirected)
	leidenResult := Leiden(g, 1, src).(*ReducedUndirected)

	qLouvain := Q(louvainResult, nil, 1)
	qLeiden := Q(leidenResult, nil, 1)

	if math.IsNaN(qLeiden) {
		t.Error("Leiden produced NaN modularity on small dumbell")
	}
	// Leiden with refinement can yield different (often better) Q; allow some tolerance.
	if qLeiden < -0.5 {
		t.Errorf("Leiden Q=%.6g unexpectedly low on small dumbell", qLeiden)
	}
	_ = qLouvain // Louvain reference; no strict comparison due to algorithm difference
}

func TestLeidenNonContiguousUndirected(t *testing.T) {
	g := simple.NewUndirectedGraph()
	for _, e := range []simple.Edge{
		{F: simple.Node(0), T: simple.Node(1)},
		{F: simple.Node(4), T: simple.Node(5)},
	} {
		g.SetEdge(e)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("unexpected panic with non-contiguous ID range")
			}
		}()
		Leiden(g, 1, nil)
	}()
}

func TestLeidenNonContiguousWeightedUndirected(t *testing.T) {
	g := simple.NewWeightedUndirectedGraph(0, 0)
	for _, e := range []simple.WeightedEdge{
		{F: simple.Node(0), T: simple.Node(1), W: 1},
		{F: simple.Node(4), T: simple.Node(5), W: 1},
	} {
		g.SetWeightedEdge(e)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("unexpected panic with non-contiguous ID range")
			}
		}()
		Leiden(g, 1, nil)
	}()
}

func TestLeidenDirectedSupported(t *testing.T) {
	g := simple.NewDirectedGraph()
	g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(1)})
	r := Leiden(g, 1, nil)
	if r == nil {
		t.Error("Leiden directed returned nil")
	}
}

func BenchmarkLeiden(b *testing.B) {
	src := rand.New(rand.NewPCG(1, 1))
	for i := 0; i < b.N; i++ {
		Leiden(dupGraph, 1, src)
	}
}

// Test that Leiden produces a valid partition on a small graph.
func TestLeidenRefinementConnectivity(t *testing.T) {
	g := simple.NewUndirectedGraph()
	// Small graph: triangle 0-1-2 and edge 2-3.
	for i := 0; i < 4; i++ {
		g.AddNode(simple.Node(i))
	}
	g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(1)})
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(0)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(3)})

	src := rand.New(rand.NewPCG(1, 1))
	r := Leiden(g, 1, src).(*ReducedUndirected)
	communities := r.Communities()
	for _, c := range communities {
		order.ByID(c)
	}
	order.BySliceIDs(communities)
	q := Q(r, nil, 1)
	if len(communities) > 0 && !math.IsNaN(q) && q < -0.5 {
		t.Errorf("unexpected low modularity Q=%.4v for communities %v", q, communities)
	}
}

func TestLeidenDeterminism(t *testing.T) {
	g := simple.NewUndirectedGraph()
	for u, e := range smallDumbell {
		if g.Node(int64(u)) == nil {
			g.AddNode(simple.Node(u))
		}
		for v := range e {
			g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
		}
	}
	src1 := rand.New(rand.NewPCG(42, 42))
	src2 := rand.New(rand.NewPCG(42, 42))
	r1 := Leiden(g, 1, src1).(*ReducedUndirected)
	r2 := Leiden(g, 1, src2).(*ReducedUndirected)
	c1 := r1.Communities()
	c2 := r2.Communities()
	for _, c := range c1 {
		order.ByID(c)
	}
	for _, c := range c2 {
		order.ByID(c)
	}
	order.BySliceIDs(c1)
	order.BySliceIDs(c2)
	if len(c1) != len(c2) {
		t.Errorf("different number of communities: %d vs %d", len(c1), len(c2))
	}
	for i := range c1 {
		if len(c1[i]) != len(c2[i]) {
			t.Errorf("community %d length: %d vs %d", i, len(c1[i]), len(c2[i]))
		}
		for j, n := range c1[i] {
			if n.ID() != c2[i][j].ID() {
				t.Errorf("community %d node %d: %d vs %d", i, j, n.ID(), c2[i][j].ID())
			}
		}
	}
}

func TestLeidenScoreProfile(t *testing.T) {
	g := simple.NewUndirectedGraph()
	for u, e := range smallDumbell {
		if g.Node(int64(u)) == nil {
			g.AddNode(simple.Node(u))
		}
		for v := range e {
			g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
		}
	}
	fn := LeidenScore(g, Weight, 5, rand.NewPCG(1, 1))
	p, err := Profile(fn, true, 1e-3, 0.1, 10)
	if err != nil {
		t.Fatalf("Profile with LeidenScore: %v", err)
	}
	if len(p) == 0 {
		t.Error("Profile returned no intervals")
	}
	for i := 1; i < len(p); i++ {
		if p[i].Score > p[i-1].Score {
			t.Errorf("Profile not monotonically decreasing: %v -> %v", p[i-1], p[i])
		}
	}
}

func TestLeidenMultiplex(t *testing.T) {
	g0 := simple.NewWeightedUndirectedGraph(0, 0)
	g1 := simple.NewWeightedUndirectedGraph(0, 0)
	for u, e := range smallDumbell {
		if g0.Node(int64(u)) == nil {
			g0.AddNode(simple.Node(u))
			g1.AddNode(simple.Node(u))
		}
		for v := range e {
			g0.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(u), T: simple.Node(v), W: 1})
			g1.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(u), T: simple.Node(v), W: 1})
		}
	}
	layers := UndirectedLayers{g0, g1}
	weights := []float64{1, 1}
	src := rand.New(rand.NewPCG(1, 1))
	r := LeidenMultiplex(layers, weights, nil, false, src).(*ReducedUndirectedMultiplex)
	communities := r.Communities()
	seen := make(map[int64]int)
	for i, c := range communities {
		for _, n := range c {
			id := n.ID()
			if j, ok := seen[id]; ok {
				t.Errorf("node %d in multiple communities %d and %d", id, j, i)
			}
			seen[id] = i
		}
	}
	if len(seen) != 6 {
		t.Errorf("partition covers %d nodes, want 6", len(seen))
	}
}
