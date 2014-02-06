package concrete

import (
	"github.com/gonum/graph"
	"math"
)

type DenseGraph struct {
	adjacencyMatrix []float64
	numNodes        int
}

func NewDenseGraph(numNodes int, passable bool) graph.Graph {
	dg := &DenseGraph{adjacencyMatrix: make([]float64, numNodes*numNodes), numNodes: numNodes}
	if passable {
		for i := range dg.adjacencyMatrix {
			dg.adjacencyMatrix[i] = 1.0
		}
	} else {
		for i := range dg.adjacencyMatrix {
			dg.adjacencyMatrix[i] = math.Inf(1)
		}
	}

	return dg
}

func (dg *DenseGraph) NodeExists(node graph.Node) bool {
	return node.ID() < dg.numNodes
}

func (dg *DenseGraph) Degree(node graph.Node) int {
	deg := 0
	for i := 0; i < dg.numNodes; i++ {
		if dg.adjacencyMatrix[i*dg.numNodes+node.ID()] != math.Inf(1) {
			deg++
		}

		if dg.adjacencyMatrix[node.ID()*dg.numNodes+i] != math.Inf(1) {
			deg++
		}
	}

	return deg
}

func (dg *DenseGraph) NodeList() []graph.Node {
	nodes := make([]graph.Node, dg.numNodes)
	for i := 0; i < dg.numNodes; i++ {
		nodes[i] = GonumNode(i)
	}

	return nodes
}

func (dg *DenseGraph) DirectedEdgeList() []graph.Edge {
	edges := make([]graph.Edge, 0, len(dg.adjacencyMatrix))
	for i := 0; i < dg.numNodes; i++ {
		for j := 0; j < dg.numNodes; j++ {
			if dg.adjacencyMatrix[i*dg.numNodes+j] != math.Inf(1) {
				edges = append(edges, GonumEdge{GonumNode(i), GonumNode(j)})
			}
		}
	}

	return edges
}

func (dg *DenseGraph) Neighbors(node graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < dg.numNodes; i++ {
		if dg.adjacencyMatrix[i*dg.numNodes+node.ID()] != math.Inf(1) ||
			dg.adjacencyMatrix[node.ID()*dg.numNodes+i] != math.Inf(1) {
			neighbors = append(neighbors, GonumNode(i))
		}
	}

	return neighbors
}

func (dg *DenseGraph) IsNeighbor(node, neighbor graph.Node) bool {
	return dg.adjacencyMatrix[neighbor.ID()*dg.numNodes+node.ID()] != math.Inf(1) ||
		dg.adjacencyMatrix[node.ID()*dg.numNodes+neighbor.ID()] != math.Inf(1)
}

func (dg *DenseGraph) Successors(node graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < dg.numNodes; i++ {
		if dg.adjacencyMatrix[node.ID()*dg.numNodes+i] != math.Inf(1) {
			neighbors = append(neighbors, GonumNode(i))
		}
	}

	return neighbors
}

func (dg *DenseGraph) IsSuccessor(node, succ graph.Node) bool {
	return dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()] != math.Inf(1)
}

func (dg *DenseGraph) Predecessors(node graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < dg.numNodes; i++ {
		if dg.adjacencyMatrix[i*dg.numNodes+node.ID()] != math.Inf(1) {
			neighbors = append(neighbors, GonumNode(i))
		}
	}

	return neighbors
}

func (dg *DenseGraph) IsPredecessor(node, pred graph.Node) bool {
	return dg.adjacencyMatrix[pred.ID()*dg.numNodes+node.ID()] != math.Inf(1)
}

// Naturally dense, we don't need to do anything
func (dg *DenseGraph) Crunch() {
}

func (dg *DenseGraph) SetEdgeCost(node, succ graph.Node, cost float64, directed bool) {
	dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()] = cost
	if !directed {
		dg.adjacencyMatrix[succ.ID()*dg.numNodes+node.ID()] = cost
	}
}

// More or less equivalent to SetEdgeCost(node, succ, math.Inf(1), directed)
func (dg *DenseGraph) RemoveEdge(node, succ graph.Node, directed bool) {
	dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()] = math.Inf(1)
	if !directed {
		dg.adjacencyMatrix[succ.ID()*dg.numNodes+node.ID()] = math.Inf(1)
	}
}
