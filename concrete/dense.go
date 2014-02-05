package concrete

import (
	"github.com/gonum/graph"
	"math"
)

// A dense graph is a graph such that all IDs are in a contiguous block from 0 to TheNumberOfNodes-1
// it uses an adjacency matrix and should be relatively fast for both access and writing.
//
// This graph implements the CrunchGraph, but since it's naturally dense this is superfluous
type DenseGraph struct {
	adjacencyMatrix []float64
	numNodes        int
}

// Creates a dense graph with the proper number of nodes. If passable is true all nodes will have
// an edge with cost 1.0, otherwise every node will start unconnected (cost of +Inf)
func NewDenseGraph(numNodes int, passable bool) *DenseGraph {
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

// DenseGraph is naturally dense, we don't need to do anything
func (dg *DenseGraph) Crunch() {
}

func (dg *DenseGraph) Cost(node, succ graph.Node) float64 {
	return dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()]
}

// Sets the cost of the edge between node and succ. If the cost is +Inf, it will remove the edge,
// if directed is true, it will only remove the edge one way. If it's false it will change the cost
// of the edge from succ to node as well.
func (dg *DenseGraph) SetEdgeCost(node, succ graph.Node, cost float64, directed bool) {
	dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()] = cost
	if !directed {
		dg.adjacencyMatrix[succ.ID()*dg.numNodes+node.ID()] = cost
	}
}

<<<<<<< HEAD
// More or less equivalent to SetEdgeCost(node, succ, math.Inf(1), directed)
=======
>>>>>>> Made a basic dense graph
func (dg *DenseGraph) RemoveEdge(node, succ graph.Node, directed bool) {
	dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()] = math.Inf(1)
	if !directed {
		dg.adjacencyMatrix[succ.ID()*dg.numNodes+node.ID()] = math.Inf(1)
	}
}
