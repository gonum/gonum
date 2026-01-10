package topo

import (
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

// --- Helpers ---

func cloneDirected(g graph.Directed) *simple.DirectedGraph {
	ng := simple.NewDirectedGraph()

	// Copy nodes
	it := g.Nodes()
	for it.Next() {
		ng.AddNode(simple.Node(it.Node().ID()))
	}

	// Copy edges
	it = g.Nodes()
	for it.Next() {
		u := it.Node()
		from := g.From(u.ID())
		for from.Next() {
			v := from.Node()
			ng.SetEdge(simple.Edge{F: simple.Node(u.ID()), T: simple.Node(v.ID())})
		}
	}
	return ng
}

func nodeIDs(g graph.Graph) []int64 {
	var ids []int64
	it := g.Nodes()
	for it.Next() {
		ids = append(ids, it.Node().ID())
	}
	return ids
}

// hasPath reports whether there exists a directed path from src to dst.
// This is a basic DFS; good enough for unit tests.
func hasPath(g graph.Directed, src, dst int64) bool {
	if src == dst {
		return true
	}
	seen := make(map[int64]bool, 16)
	stack := []int64{src}
	seen[src] = true

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		it := g.From(u)
		for it.Next() {
			v := it.Node().ID()
			if v == dst {
				return true
			}
			if !seen[v] {
				seen[v] = true
				stack = append(stack, v)
			}
		}
	}
	return false
}

func assertReachabilityPreserved(t *testing.T, before, after graph.Directed) {
	t.Helper()
	ids := nodeIDs(before)

	for _, u := range ids {
		for _, v := range ids {
			got := hasPath(after, u, v)
			want := hasPath(before, u, v)
			if got != want {
				t.Fatalf("reachability changed for (%d -> %d): before=%v after=%v", u, v, want, got)
			}
		}
	}
}

func assertEdgesSubset(t *testing.T, original, reduced graph.Directed) {
	t.Helper()
	it := reduced.Nodes()
	for it.Next() {
		u := it.Node().ID()
		from := reduced.From(u)
		for from.Next() {
			v := from.Node().ID()
			if !original.HasEdgeFromTo(u, v) {
				t.Fatalf("reduced graph contains edge %d->%d not present in original", u, v)
			}
		}
	}
}

func assertMinimal(t *testing.T, reduced graph.Directed) {
	t.Helper()

	// For each edge u->v in reduced, removing it must break reachability u=>v.
	// (i.e., the edge is necessary for preserving reachability.)
	it := reduced.Nodes()
	for it.Next() {
		u := it.Node().ID()
		from := reduced.From(u)
		for from.Next() {
			v := from.Node().ID()

			// Work on a clone to avoid needing EdgeAdder to restore.
			g2 := cloneDirected(reduced)
			g2.RemoveEdge(u, v)

			if hasPath(g2, u, v) {
				t.Fatalf("edge %d->%d is redundant: path still exists after removal", u, v)
			}
		}
	}
}

// --- Test graphs ---

func makeChainWithShortcuts(n int) *simple.DirectedGraph {
	g := simple.NewDirectedGraph()
	for i := 0; i < n; i++ {
		g.AddNode(simple.Node(int64(i)))
	}
	// Chain edges i -> i+1
	for i := 0; i < n-1; i++ {
		g.SetEdge(simple.Edge{F: simple.Node(int64(i)), T: simple.Node(int64(i + 1))})
	}
	// Add transitive shortcuts from 0 -> k for k >= 2
	for k := 2; k < n; k++ {
		g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(int64(k))})
	}
	return g
}

func makeDiamondWithShortcut() *simple.DirectedGraph {
	// 1->2, 1->3, 2->4, 3->4, plus shortcut 1->4
	g := simple.NewDirectedGraph()
	for _, id := range []int64{1, 2, 3, 4} {
		g.AddNode(simple.Node(id))
	}
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(3)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(4)})
	g.SetEdge(simple.Edge{F: simple.Node(3), T: simple.Node(4)})
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(4)}) // redundant
	return g
}

func makeAlreadyReduced() *simple.DirectedGraph {
	// Simple chain: no redundant edges.
	g := simple.NewDirectedGraph()
	for _, id := range []int64{10, 11, 12, 13} {
		g.AddNode(simple.Node(id))
	}
	g.SetEdge(simple.Edge{F: simple.Node(10), T: simple.Node(11)})
	g.SetEdge(simple.Edge{F: simple.Node(11), T: simple.Node(12)})
	g.SetEdge(simple.Edge{F: simple.Node(12), T: simple.Node(13)})
	return g
}

func makeDisconnected() *simple.DirectedGraph {
	// Two components:
	// A: 0->1->2 plus 0->2 (redundant)
	// B: 10->11
	g := simple.NewDirectedGraph()
	for _, id := range []int64{0, 1, 2, 10, 11} {
		g.AddNode(simple.Node(id))
	}
	g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(1)})
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(2)}) // redundant
	g.SetEdge(simple.Edge{F: simple.Node(10), T: simple.Node(11)})
	return g
}

// --- Tests ---

func TestTransitiveReduce_ChainWithShortcuts(t *testing.T) {
	orig := makeChainWithShortcuts(8)
	before := cloneDirected(orig)
	after := cloneDirected(orig)

	TransitiveReduce(after)

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	// Spot-check: in a chain-with-shortcuts, all 0->k (k>=2) should be removed.
	for k := int64(2); k < 8; k++ {
		if after.HasEdgeFromTo(0, k) {
			t.Fatalf("expected redundant edge 0->%d to be removed", k)
		}
	}
}

func TestTransitiveReduce_Diamond(t *testing.T) {
	orig := makeDiamondWithShortcut()
	before := cloneDirected(orig)
	after := cloneDirected(orig)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	// Spot-check: 1->4 should be removed.
	if after.HasEdgeFromTo(1, 4) {
		t.Fatalf("expected redundant edge 1->4 to be removed")
	}
}

func TestTransitiveReduce_AlreadyReduced(t *testing.T) {
	orig := makeAlreadyReduced()
	before := cloneDirected(orig)
	after := cloneDirected(orig)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	// Should remain unchanged in edge count.
	if after.Edges().Len() != before.Edges().Len() {
		t.Fatalf("expected edge count unchanged: before=%d after=%d", before.Edges().Len(), after.Edges().Len())
	}
}

func TestTransitiveReduce_Disconnected(t *testing.T) {
	orig := makeDisconnected()
	before := cloneDirected(orig)
	after := cloneDirected(orig)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	// Spot-check: redundant 0->2 should be removed, and 10->11 should remain.
	if after.HasEdgeFromTo(0, 2) {
		t.Fatalf("expected redundant edge 0->2 to be removed")
	}
	if !after.HasEdgeFromTo(10, 11) {
		t.Fatalf("expected edge 10->11 to remain")
	}
}

func TestTransitiveReduce_EmptyGraph(t *testing.T) {
	g := simple.NewDirectedGraph()
	before := cloneDirected(g)
	after := cloneDirected(g)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	// Minimality is vacuously true.
	assertMinimal(t, after)

	if after.Nodes().Len() != 0 {
		t.Fatalf("expected 0 nodes, got %d", after.Nodes().Len())
	}
}

func TestTransitiveReduce_SingleNodeNoEdges(t *testing.T) {
	g := simple.NewDirectedGraph()
	g.AddNode(simple.Node(1))

	before := cloneDirected(g)
	after := cloneDirected(g)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	if after.Edges().Len() != 0 {
		t.Fatalf("expected 0 edges, got %d", after.Edges().Len())
	}
}

func TestTransitiveReduce_TwoNodesSingleEdge(t *testing.T) {
	g := simple.NewDirectedGraph()
	g.AddNode(simple.Node(1))
	g.AddNode(simple.Node(2))
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})

	before := cloneDirected(g)
	after := cloneDirected(g)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	if !after.HasEdgeFromTo(1, 2) {
		t.Fatalf("expected edge 1->2 to remain")
	}
}

func TestTransitiveReduce_CycleReturnsError(t *testing.T) {
	g := simple.NewDirectedGraph()
	for _, id := range []int64{1, 2, 3} {
		g.AddNode(simple.Node(id))
	}
	g.SetEdge(simple.Edge{F: simple.Node(1), T: simple.Node(2)})
	g.SetEdge(simple.Edge{F: simple.Node(2), T: simple.Node(3)})
	g.SetEdge(simple.Edge{F: simple.Node(3), T: simple.Node(1)}) // cycle

	if err := TransitiveReduce(g); err == nil {
		t.Fatalf("expected error for cyclic graph, got nil")
	}
}

func TestTransitiveReduce_StarNoTransitiveEdges(t *testing.T) {
	// 0 -> {1,2,3,4} and no other edges => nothing is redundant.
	g := simple.NewDirectedGraph()
	for i := int64(0); i <= 4; i++ {
		g.AddNode(simple.Node(i))
	}
	for i := int64(1); i <= 4; i++ {
		g.SetEdge(simple.Edge{F: simple.Node(0), T: simple.Node(i)})
	}

	before := cloneDirected(g)
	after := cloneDirected(g)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	if after.Edges().Len() != before.Edges().Len() {
		t.Fatalf("expected no edge removals: before=%d after=%d", before.Edges().Len(), after.Edges().Len())
	}
}

func TestTransitiveReduce_LayeredDAGManyRedundantEdges(t *testing.T) {
	// Layers:
	// 0 -> (1,2,3)
	// (1,2,3) -> (4,5,6)
	// (4,5,6) -> 7
	// plus lots of redundant shortcuts: 0->(4,5,6), 0->7, (1,2,3)->7
	g := simple.NewDirectedGraph()
	for i := int64(0); i <= 7; i++ {
		g.AddNode(simple.Node(i))
	}
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

	before := cloneDirected(g)
	after := cloneDirected(g)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)

	// Spot-check: the direct 0->7 must be removed.
	if after.HasEdgeFromTo(0, 7) {
		t.Fatalf("expected redundant edge 0->7 to be removed")
	}
	// Spot-check: all 0->(4,5,6) must be removed (reachable via 0->(1,2,3)->(4,5,6))
	for _, v := range []int64{4, 5, 6} {
		if after.HasEdgeFromTo(0, v) {
			t.Fatalf("expected redundant edge 0->%d to be removed", v)
		}
	}
}

func TestTransitiveReduce_RandomSmallDAG(t *testing.T) {
	// Deterministic "random" DAG:
	// Add edge i->j for i<j if (i*17 + j*31) % 5 == 0, and then sprinkle some extra.
	const n = 12
	g := simple.NewDirectedGraph()
	for i := 0; i < n; i++ {
		g.AddNode(simple.Node(int64(i)))
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if (i*17+j*31)%5 == 0 {
				g.SetEdge(simple.Edge{F: simple.Node(int64(i)), T: simple.Node(int64(j))})
			}
		}
	}
	// Ensure connected-ish backbone
	for i := 0; i < n-1; i++ {
		g.SetEdge(simple.Edge{F: simple.Node(int64(i)), T: simple.Node(int64(i + 1))})
	}

	before := cloneDirected(g)
	after := cloneDirected(g)

	if err := TransitiveReduce(after); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEdgesSubset(t, before, after)
	assertReachabilityPreserved(t, before, after)
	assertMinimal(t, after)
}
