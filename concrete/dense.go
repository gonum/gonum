package concrete

import (
	"github.com/gonum/graph"
	"math"
)

// A dense graph is a graph such that all IDs are in a contiguous block from 0 to
// TheNumberOfNodes-1. It uses an adjacency matrix and should be relatively fast for both access
// and writing.
//
// This graph implements the CrunchGraph, but since it's naturally dense this is superfluous.
type DenseGraph struct {
	adjacencyMatrix []float64
	numNodes        int
}

// Creates a dense graph with the proper number of nodes. If passable is true all nodes will have
// an edge with cost 1.0, otherwise every node will start unconnected (cost of +Inf.)
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
		nodes[i] = Node(i)
	}

	return nodes
}

func (dg *DenseGraph) DirectedEdgeList() []graph.Edge {
	edges := make([]graph.Edge, 0, len(dg.adjacencyMatrix))
	for i := 0; i < dg.numNodes; i++ {
		for j := 0; j < dg.numNodes; j++ {
			if dg.adjacencyMatrix[i*dg.numNodes+j] != math.Inf(1) {
				edges = append(edges, Edge{Node(i), Node(j)})
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
			neighbors = append(neighbors, Node(i))
		}
	}

	return neighbors
}

func (dg *DenseGraph) EdgeBetween(node, neighbor graph.Node) graph.Edge {
	if dg.adjacencyMatrix[neighbor.ID()*dg.numNodes+node.ID()] != math.Inf(1) ||
		dg.adjacencyMatrix[node.ID()*dg.numNodes+neighbor.ID()] != math.Inf(1) {
		return Edge{node, neighbor}
	}

	return nil
}

func (dg *DenseGraph) Successors(node graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < dg.numNodes; i++ {
		if dg.adjacencyMatrix[node.ID()*dg.numNodes+i] != math.Inf(1) {
			neighbors = append(neighbors, Node(i))
		}
	}

	return neighbors
}

func (dg *DenseGraph) EdgeTo(node, succ graph.Node) graph.Edge {
	if dg.adjacencyMatrix[node.ID()*dg.numNodes+succ.ID()] != math.Inf(1) {
		return Edge{node, succ}
	}

	return nil
}

func (dg *DenseGraph) Predecessors(node graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < dg.numNodes; i++ {
		if dg.adjacencyMatrix[i*dg.numNodes+node.ID()] != math.Inf(1) {
			neighbors = append(neighbors, Node(i))
		}
	}

	return neighbors
}

// DenseGraph is naturally dense, we don't need to do anything
func (dg *DenseGraph) Crunch() {
}

func (dg *DenseGraph) Cost(e graph.Edge) float64 {
	return dg.adjacencyMatrix[e.Head().ID()*dg.numNodes+e.Tail().ID()]
}

// Sets the cost of an edge. If the cost is +Inf, it will remove the edge,
// if directed is true, it will only remove the edge one way. If it's false it will change the cost
// of the edge from succ to node as well.
func (dg *DenseGraph) SetEdgeCost(e graph.Edge, cost float64, directed bool) {
	dg.adjacencyMatrix[e.Head().ID()*dg.numNodes+e.Tail().ID()] = cost
	if !directed {
		dg.adjacencyMatrix[e.Tail().ID()*dg.numNodes+e.Head().ID()] = cost
	}
}

// Equivalent to SetEdgeCost(edge, math.Inf(1), directed)
func (dg *DenseGraph) RemoveEdge(e graph.Edge, directed bool) {
	dg.SetEdgeCost(e, math.Inf(1), directed)
}
