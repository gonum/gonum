package network

import (
	"math"
	"strings"
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
	expected := "source and target must be different"
	if err.Error() != expected {
		t.Fatalf("unexpected error message: got %q, expected %q", err.Error(), expected)
	}
}

func TestNegativeCapacityRaisesError(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 3; i++ {
		graph.AddNode(simple.Node(i))
	}
	edges := []struct {
		u, v int64
		w    float64
	}{
		{0, 1, 0.3}, {1, 2, -0.6},
	}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), edge.w))
	}

	_, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(1))
	if err == nil {
		t.Fatal("expected an error when graph contains a negative capacity, got nil")
	}
	if !strings.Contains(err.Error(), "edge weights (capacities) can not be negative") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestMaxFlowDinicThreeDisjointPaths(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 5; i++ {
		graph.AddNode(simple.Node(i))
	}
	edges := []struct{ u, v int64 }{{0, 1}, {1, 4}, {0, 2}, {2, 4}, {0, 3}, {3, 4}}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), 1.0))
	}
	maxFlow, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(4))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow, 3.0) {
		t.Errorf("maxFlow = %v, expected %v", maxFlow, 3.0)
	}
}

func TestMaxFlowDinicThreeDisjointPathsWithParallelEdges(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 5; i++ {
		graph.AddNode(simple.Node(i))
	}
	edges := []struct{ u, v int64 }{{0, 1}, {1, 0}, {1, 4}, {0, 2}, {2, 0}, {2, 4}, {0, 3}, {3, 0}, {3, 4}}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), 1.0))
	}
	maxFlow, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(4))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow, 3.0) {
		t.Errorf("maxFlow = %v, expected %v", maxFlow, 3.0)
	}
}

func TestMaxFlowDinicCycleWithTailGraph(t *testing.T) {
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

	maxFlow03, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow03, 0.3) {
		t.Errorf("maxFlow 0->3 = %v, expected %v", maxFlow03, 0.3)
	}

	maxFlow13, err := MaxFlowDinic(graph, simple.Node(1), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow13, 0.6) {
		t.Errorf("maxFlow 1->3 = %v, expected %v", maxFlow13, 0.6)
	}
}

func TestMaxFlowDinicCycleWithTailGraphWithParallelEdges(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 4; i++ {
		graph.AddNode(simple.Node(i))
	}
	edges := []struct {
		u, v int64
		w    float64
	}{
		{0, 1, 0.3}, {1, 2, 0.6}, {2, 0, 0.9}, {2, 3, 0.7},
		{1, 0, 1.3}, {2, 1, 1.6}, {0, 2, 1.9},
	}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), edge.w))
	}

	maxFlow03, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow03, 0.7) {
		t.Errorf("maxFlow 0->3 = %v, expected %v", maxFlow03, 0.7)
	}

	maxFlow13, err := MaxFlowDinic(graph, simple.Node(1), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow13, 0.7) {
		t.Errorf("maxFlow 1->3 = %v, expected %v", maxFlow13, 0.7)
	}
}

func TestMaxFlowDinicFourLayeredDAG(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 8; i++ {
		graph.AddNode(simple.Node(i))
	}
	// Layers: 0->{1,2,3}, {1,2}->{4,5}, {3}->{5,6}, {4,5,6}->{7}
	edges := []struct{ u, v int64 }{
		{0, 1}, {0, 2}, {0, 3},
		{1, 4}, {2, 4}, {2, 5}, {3, 5}, {3, 6},
		{4, 7}, {5, 7}, {6, 7},
	}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), 1.0))
	}

	testCases := []struct {
		s, t, expected float64
	}{
		{0, 7, 3.0},
		{3, 7, 2.0},
		{0, 5, 2.0},
		{2, 4, 1.0},
	}
	for _, tc := range testCases {
		maxFlow, err := MaxFlowDinic(graph, simple.Node(int64(tc.s)), simple.Node(int64(tc.t)))
		if err != nil {
			t.Fatalf("Unexpected error for %v->%v: %v", tc.s, tc.t, err)
		}
		if !almostEqual(maxFlow, tc.expected) {
			t.Errorf("maxFlow %v->%v = %v, expected %v", tc.s, tc.t, maxFlow, tc.expected)
		}
	}
}

func TestMaxFlowDinicDiamondWithCrossGraph(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 4; i++ {
		graph.AddNode(simple.Node(i))
	}
	edges := []struct {
		u, v int64
		w    float64
	}{
		{0, 1, 10.0}, {0, 2, 10.0}, {1, 2, 5.0}, {1, 3, 10.0}, {2, 3, 10.0},
	}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), edge.w))
	}

	maxFlow03, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(3))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow03, 20.0) {
		t.Errorf("maxFlow 0->3 = %v, expected %v", maxFlow03, 20.0)
	}

	maxFlow02, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(2))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow02, 15.0) {
		t.Errorf("maxFlow 0->2 = %v, expected %v", maxFlow02, 15.0)
	}
}

func TestMaxFlowDinicDisconnectedGraphs(t *testing.T) {
	graph := simple.NewWeightedDirectedGraph(0, 0)
	for i := int64(0); i < 7; i++ {
		graph.AddNode(simple.Node(i))
	}
	edges := []struct {
		u, v int64
		w    float64
	}{
		{0, 1, 10.0}, {1, 2, 5.0}, {2, 3, 7.0},
		{4, 5, 11.0}, {5, 6, 10.0},
	}
	for _, edge := range edges {
		graph.SetWeightedEdge(graph.NewWeightedEdge(simple.Node(edge.u), simple.Node(edge.v), edge.w))
	}

	maxFlow, err := MaxFlowDinic(graph, simple.Node(0), simple.Node(5))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !almostEqual(maxFlow, 0.0) {
		t.Errorf("maxFlow = %v, expected %v", maxFlow, 0.0)
	}
}
