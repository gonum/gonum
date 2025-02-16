package centrality_test

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/centrality"
	"gonum.org/v1/gonum/graph/simple"
	"math"
	"testing"
)

func numericalEqual(a, b, epsilon float64) bool {
	if math.Abs(b) < epsilon {
		return math.Abs(a-b) < epsilon
	}
	return math.Abs(1.0-a/b) < epsilon
}

// addNodes creates `count` nodes, adds them to `g`, and returns the slice of nodes.
func addNodes(undirectedGraph *simple.UndirectedGraph, count int) []graph.Node {
	nodes := make([]graph.Node, count)
	for i := 0; i < count; i++ {
		nodes[i] = undirectedGraph.NewNode() // Create a new node
		undirectedGraph.AddNode(nodes[i])    // Add it to the undirectedGraph
	}
	return nodes
}

// Test ClosenessCentrality on an undirected undirectedGraph
func TestClosenessCentrality(test *testing.T) {
	tests := []struct {
		name            string
		undirectedGraph func() *simple.UndirectedGraph
		expectedResult  map[int64]float64
	}{
		{
			name: "Cycle Graph (4 Nodes)",
			undirectedGraph: func() *simple.UndirectedGraph {
				g := simple.NewUndirectedGraph()
				nodes := addNodes(g, 4)
				g.SetEdge(simple.Edge{F: nodes[0], T: nodes[1]})
				g.SetEdge(simple.Edge{F: nodes[1], T: nodes[2]})
				g.SetEdge(simple.Edge{F: nodes[2], T: nodes[3]})
				g.SetEdge(simple.Edge{F: nodes[3], T: nodes[0]})
				return g
			},
			expectedResult: map[int64]float64{0: 0.75, 1: 0.75, 2: 0.75, 3: 0.75},
		},
		{
			name: "Star Graph (5 Nodes)",
			undirectedGraph: func() *simple.UndirectedGraph {
				g := simple.NewUndirectedGraph()
				nodes := addNodes(g, 5)
				// Star topology (0 is the central node)
				for i := 1; i < 5; i++ {
					g.SetEdge(simple.Edge{F: nodes[0], T: nodes[i]})
				}
				return g
			},
			expectedResult: map[int64]float64{0: 1.0, 1: 4.0 / 7.0, 2: 4.0 / 7.0, 3: 4.0 / 7.0, 4: 4.0 / 7.0},
		},
		{
			name: "Line Graph (5 Nodes)",
			undirectedGraph: func() *simple.UndirectedGraph {
				g := simple.NewUndirectedGraph()
				nodes := addNodes(g, 5)
				// Line topology (A-B-C-D-E)
				for i := 0; i < 4; i++ {
					g.SetEdge(simple.Edge{F: nodes[i], T: nodes[i+1]})
				}
				return g
			},
			expectedResult: map[int64]float64{0: 4.0 / 10.0, 1: 4.0 / 7.0, 2: 4.0 / 6.0, 3: 4.0 / 7.0, 4: 4.0 / 10.0},
		},
	}

	const epsilon float64 = 1.e-8

	for _, testCase := range tests {
		test.Run(testCase.name, func(t *testing.T) {
			result := centrality.ClosenessCentrality(testCase.undirectedGraph())
			for id, expectedValue := range testCase.expectedResult {
				if !numericalEqual(result[id], expectedValue, epsilon) {
					t.Errorf("%s: ClosenessCentrality(%d) = %f, expectedResult %f",
						testCase.name, id, result[id], expectedValue)
				}
			}
		})
	}
}

// Test ClosenessCentralityWeighted on a weighted undirectedGraph
/*
func TestClosenessCentralityWeighted(t *testing.T) {
	g := simple.NewWeightedDirectedGraph(0, 0)

	const numberOfNodes = 4

	nodes := make([]undirectedGraph.Node, numberOfNodes)
	// Add nodes
	for i := 0; i < numberOfNodes; i++ {
		nodes[i] = g.NewNode()
		g.AddNode(nodes[i])

	}
	// Add edges
	g.SetWeightedEdge(simple.WeightedEdge{F: nodes[0], T: nodes[1], W: 1.0})
	g.SetWeightedEdge(simple.WeightedEdge{F: nodes[1], T: nodes[2], W: 2.0})
	g.SetWeightedEdge(simple.WeightedEdge{F: nodes[2], T: nodes[3], W: 3.0})
	g.SetWeightedEdge(simple.WeightedEdge{F: nodes[3], T: nodes[0], W: 4.0})

	// Compute centrality
	result := centrality.ClosenessCentrality(g)
	const expectedValueNode0 float64 = 0.75
	const expectedValueNode1 float64 = 0.75
	const expectedValueNode2 float64 = 0.75
	const expectedValueNode3 float64 = 0.75

	// Expected values (hand-calculated)
	expectedResult := map[int64]float64{
		nodes[0].ID(): expectedValue,
		nodes[1].ID(): expectedValue,
		nodes[2].ID(): expectedValue,
		nodes[3].ID(): expectedValue,
	}

	// Compare results
	for id, expectedValue := range expectedResult {
		if !numericalEqual(result[id], expectedValue, 0.001) {
			t.Errorf("ClosenessCentrality(%d) = %f, expectedResult %f", id, result[id], expectedValue)
		}
	}

	// Compare results
	for id, expectedValue := range expectedResult {
		if !numericalEqual(result[id], expectedValue, 0.001) {
			t.Errorf("ClosenessCentralityWeighted(%d) = %f, expectedResult %f", id, result[id], expectedValue)
		}
	}
}
*/
