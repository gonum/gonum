package network

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

// almostEqual tests float64 equality within a small epsilon.
func almostEqual(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) <= eps
}

func TestMaxFlowDinicSameSourceAndTarget(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	// Add a single node
	graph.AddNode(simple.Node(0))
	_, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(0))
	if err == nil {
		t.Fatal("Expected error when source and target are the same, got nil")
	}
	want := "source and target must be different"
	if err.Error() != want {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

func TestThreeDisjointPaths(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	// Add nodes 0..4
	for i := int64(0); i < 5; i++ {
		graph.AddNode(simple.Node(i))
	}
	// Build three disjoint paths of capacity 1
	edges := []struct{ u, v int64 }{{0, 1}, {1, 4}, {0, 2}, {2, 4}, {0, 3}, {3, 4}}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), 1.0))
	}
	flow, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(4))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(flow, 3.0) {
		t.Errorf("flow = %v, want %v", flow, 3.0)
	}
}

func TestCycleWithTailGraph(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 4; i++ {
		graph.AddNode(simple.Node(i))
	}
	// Cycle: 0->1->2->0 and tail 2->3
	edges := []struct {
		u, v int64
		w    float64
	}{
		{0, 1, 0.3}, {1, 2, 0.6}, {2, 0, 0.9}, {2, 3, 0.7},
	}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), edge.w))
	}

	flow03, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(flow03, 0.3) {
		t.Errorf("flow 0->3 = %v, want %v", flow03, 0.3)
	}

	flow13, err := MaxFlowDinic(graph, simple.Node(1), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(flow13, 0.6) {
		t.Errorf("flow 1->3 = %v, want %v", flow13, 0.6)
	}
}
