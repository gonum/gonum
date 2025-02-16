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
		undirectedGraph.AddNode(nodes[i])    // Add it to the weightedUndirectedGraph
	}
	return nodes
}

func addNodesWeightedGraph(weightedUndirectedGraph *simple.WeightedUndirectedGraph, count int) []graph.Node {
	nodes := make([]graph.Node, count)
	for i := 0; i < count; i++ {
		nodes[i] = weightedUndirectedGraph.NewNode() // Create a new node
		weightedUndirectedGraph.AddNode(nodes[i])    // Add it to the weightedUndirectedGraph
	}
	return nodes
}

// Test ClosenessCentrality on an undirected weightedUndirectedGraph
func TestClosenessCentrality(test *testing.T) {
	tests := []struct {
		name            string
		undirectedGraph func() *simple.UndirectedGraph
		expectedResult  map[int64]float64
	}{
		{
			name: "Empty Graph",
			undirectedGraph: func() *simple.UndirectedGraph {
				return simple.NewUndirectedGraph() // No nodes, no edges
			},
			expectedResult: nil, // We expect a nil or empty result
		},
		{
			name: "Graph with one node only",
			undirectedGraph: func() *simple.UndirectedGraph {
				oneNode := simple.NewUndirectedGraph()
				node := oneNode.NewNode()
				oneNode.AddNode(node)
				return oneNode
			},
			expectedResult: nil, // We expect a nil or empty result
		},
		{
			name: "Cycle Graph (4 Nodes)",
			undirectedGraph: func() *simple.UndirectedGraph {
				cycleGraph := simple.NewUndirectedGraph()
				nodes := addNodes(cycleGraph, 4)
				cycleGraph.SetEdge(simple.Edge{F: nodes[0], T: nodes[1]})
				cycleGraph.SetEdge(simple.Edge{F: nodes[1], T: nodes[2]})
				cycleGraph.SetEdge(simple.Edge{F: nodes[2], T: nodes[3]})
				cycleGraph.SetEdge(simple.Edge{F: nodes[3], T: nodes[0]})
				return cycleGraph
			},
			expectedResult: map[int64]float64{0: 0.75, 1: 0.75, 2: 0.75, 3: 0.75},
		},
		{
			name: "Star Graph (5 Nodes)",
			undirectedGraph: func() *simple.UndirectedGraph {
				starGraph := simple.NewUndirectedGraph()
				nodes := addNodes(starGraph, 5)
				// Star topology (0 is the central node)
				for i := 1; i < 5; i++ {
					starGraph.SetEdge(simple.Edge{F: nodes[0], T: nodes[i]})
				}
				return starGraph
			},
			expectedResult: map[int64]float64{0: 1.0, 1: 4.0 / 7.0, 2: 4.0 / 7.0, 3: 4.0 / 7.0, 4: 4.0 / 7.0},
		},
		{
			name: "Line Graph (5 Nodes)",
			undirectedGraph: func() *simple.UndirectedGraph {
				lineGraph := simple.NewUndirectedGraph()
				nodes := addNodes(lineGraph, 5)
				// Line topology (A-B-C-D-E)
				for i := 0; i < 4; i++ {
					lineGraph.SetEdge(simple.Edge{F: nodes[i], T: nodes[i+1]})
				}
				return lineGraph
			},
			expectedResult: map[int64]float64{0: 4.0 / 10.0, 1: 4.0 / 7.0, 2: 4.0 / 6.0, 3: 4.0 / 7.0, 4: 4.0 / 10.0},
		},
	}

	const epsilon float64 = 1.e-8

	for _, testCase := range tests {
		test.Run(testCase.name, func(t *testing.T) {
			result := centrality.ClosenessCentrality(testCase.undirectedGraph())
			if result == nil && testCase.expectedResult != nil {
				t.Errorf("%s: Expected non-nil result, but got nil", testCase.name)
			} else if result != nil && testCase.expectedResult == nil {
				t.Errorf("%s: Expected nil result, but got non-nil", testCase.name)
			}
			for id, expectedValue := range testCase.expectedResult {
				if !numericalEqual(result[id], expectedValue, epsilon) {
					t.Errorf("%s: ClosenessCentrality(%d) = %f, expectedResult %f",
						testCase.name, id, result[id], expectedValue)
				}
			}
		})
	}
}

// Test ClosenessCentrality on an undirected weightedUndirectedGraph
func TestClosenessCentralityWeightedGraph(test *testing.T) {
	tests := []struct {
		name                    string
		weightedUndirectedGraph func() *simple.WeightedUndirectedGraph
		expectedResult          map[int64]float64
	}{
		{
			name: "Empty Graph",
			weightedUndirectedGraph: func() *simple.WeightedUndirectedGraph {
				return simple.NewWeightedUndirectedGraph(math.Inf(1), math.Inf(1)) // No nodes, no edges
			},
			expectedResult: nil, // We expect a nil or empty result
		},
		{
			name: "Graph with one node only",
			weightedUndirectedGraph: func() *simple.WeightedUndirectedGraph {
				oneNode := simple.NewWeightedUndirectedGraph(math.Inf(1), math.Inf(1))
				node := oneNode.NewNode()
				oneNode.AddNode(node)
				return oneNode
			},
			expectedResult: nil, // We expect a nil or empty result
		},
		{
			name: "Cycle Graph (4 Nodes)",
			weightedUndirectedGraph: func() *simple.WeightedUndirectedGraph {
				cycleGraph := simple.NewWeightedUndirectedGraph(math.Inf(1), math.Inf(1))
				nodes := addNodesWeightedGraph(cycleGraph, 4)
				weights := [4]float64{1.0, 2.0, 3.0, 4.0}
				cycleGraph.SetWeightedEdge(simple.WeightedEdge{F: nodes[0], T: nodes[1], W: weights[0]})
				cycleGraph.SetWeightedEdge(simple.WeightedEdge{F: nodes[1], T: nodes[2], W: weights[1]})
				cycleGraph.SetWeightedEdge(simple.WeightedEdge{F: nodes[2], T: nodes[3], W: weights[2]})
				cycleGraph.SetWeightedEdge(simple.WeightedEdge{F: nodes[3], T: nodes[0], W: weights[3]})
				return cycleGraph
			},
			expectedResult: map[int64]float64{0: 10.0 / 3.0, 1: 0.75, 2: 0.75, 3: 0.75},
		},
		{
			name: "Star Graph (5 Nodes)",
			weightedUndirectedGraph: func() *simple.UndirectedGraph {
				starGraph := simple.NewWeightedUndirectedGraph(math.Inf(1), math.Inf(1))
				nodes := addNodes(starGraph, 5)
				// Star topology (0 is the central node)
				for i := 1; i < 5; i++ {
					starGraph.SetEdge(simple.Edge{F: nodes[0], T: nodes[i]})
				}
				return starGraph
			},
			expectedResult: map[int64]float64{0: 1.0, 1: 4.0 / 7.0, 2: 4.0 / 7.0, 3: 4.0 / 7.0, 4: 4.0 / 7.0},
		},
		{
			name: "Line Graph (5 Nodes)",
			weightedUndirectedGraph: func() *simple.UndirectedGraph {
				lineGraph := simple.NewWeightedUndirectedGraph(math.Inf(1), math.Inf(1))
				nodes := addNodes(lineGraph, 5)
				// Line topology (A-B-C-D-E)
				for i := 0; i < 4; i++ {
					lineGraph.SetEdge(simple.Edge{F: nodes[i], T: nodes[i+1]})
				}
				return lineGraph
			},
			expectedResult: map[int64]float64{0: 4.0 / 10.0, 1: 4.0 / 7.0, 2: 4.0 / 6.0, 3: 4.0 / 7.0, 4: 4.0 / 10.0},
		},
	}

	const epsilon float64 = 1.e-8

	for _, testCase := range tests {
		test.Run(testCase.name, func(t *testing.T) {
			result := centrality.ClosenessCentrality(testCase.weightedUndirectedGraph())
			if result == nil && testCase.expectedResult != nil {
				t.Errorf("%s: Expected non-nil result, but got nil", testCase.name)
			} else if result != nil && testCase.expectedResult == nil {
				t.Errorf("%s: Expected nil result, but got non-nil", testCase.name)
			}
			for id, expectedValue := range testCase.expectedResult {
				if !numericalEqual(result[id], expectedValue, epsilon) {
					t.Errorf("%s: ClosenessCentrality(%d) = %f, expectedResult %f",
						testCase.name, id, result[id], expectedValue)
				}
			}
		})
	}
}
