package concrete

import (
	"math"
	"sort"

	gr "github.com/gonum/graph"
)

// A simple int alias.
type Node int

func (node Node) ID() int {
	return int(node)
}

// Just a collection of two nodes
type Edge struct {
	H, T gr.Node
}

func (edge Edge) Head() gr.Node {
	return edge.H
}

func (edge Edge) Tail() gr.Node {
	return edge.T
}

type WeightedEdge struct {
	gr.Edge
	Cost float64
}

// A GonumGraph is a very generalized graph that can handle an arbitrary number of vertices and
// edges -- as well as act as either directed or undirected.
//
// Internally, it uses a map of successors AND predecessors, to speed up some operations (such as
// getting all successors/predecessors). It also speeds up thing like adding edges (assuming both
// edges exist).
//
// However, its generality is also its weakness (and partially a flaw in needing to satisfy
// MutableGraph). For most purposes, creating your own graph is probably better. For instance,
// see TileGraph for an example of an immutable 2D grid of tiles that also implements the Graph
// interface, but would be more suitable if all you needed was a simple undirected 2D grid.
type Graph struct {
	successors   map[int]map[int]WeightedEdge
	predecessors map[int]map[int]WeightedEdge
	nodeMap      map[int]gr.Node
}

func NewGraph() *Graph {
	return &Graph{
		successors:   make(map[int]map[int]WeightedEdge),
		predecessors: make(map[int]map[int]WeightedEdge),
		nodeMap:      make(map[int]gr.Node),
	}
}

func NewPreAllocatedGraph(numVertices int) *Graph {
	return &Graph{
		successors:   make(map[int]map[int]WeightedEdge, numVertices),
		predecessors: make(map[int]map[int]WeightedEdge, numVertices),
		nodeMap:      make(map[int]gr.Node, numVertices),
	}
}

/* Mutable Graph implementation */

func (graph *Graph) NewNode() (node gr.Node) {
	nodeList := graph.NodeList()
	ids := make([]int, len(nodeList))
	for i, node := range nodeList {
		ids[i] = node.ID()
	}

	nodes := sort.IntSlice(ids)
	sort.Sort(&nodes)
	for i, node := range nodes {
		if i != node {
			graph.AddNode(Node(i))
			return Node(i)
		}
	}

	newID := len(nodes)
	graph.AddNode(Node(newID))
	return Node(newID)
}

func (graph *Graph) AddNode(node gr.Node) {
	if _, ok := graph.nodeMap[node.ID()]; ok {
		return
	}

	graph.nodeMap[node.ID()] = node
	graph.successors[node.ID()] = make(map[int]WeightedEdge)
	graph.predecessors[node.ID()] = make(map[int]WeightedEdge)
}

func (graph *Graph) AddEdge(e gr.Edge, cost float64, directed bool) {
	head, tail := e.Head(), e.Tail()
	graph.AddNode(head)
	graph.AddNode(tail)

	graph.successors[head.ID()][tail.ID()] = WeightedEdge{Edge: e, Cost: cost}
	graph.predecessors[tail.ID()][head.ID()] = WeightedEdge{Edge: e, Cost: cost}
	if !directed {
		graph.successors[tail.ID()][head.ID()] = WeightedEdge{Edge: e, Cost: cost}
		graph.predecessors[head.ID()][tail.ID()] = WeightedEdge{Edge: e, Cost: cost}
	}
}

func (graph *Graph) RemoveNode(node gr.Node) {
	if _, ok := graph.nodeMap[node.ID()]; !ok {
		return
	}
	delete(graph.nodeMap, node.ID())

	for succ, _ := range graph.successors[node.ID()] {
		delete(graph.predecessors[succ], node.ID())
	}
	delete(graph.successors, node.ID())

	for pred, _ := range graph.predecessors[node.ID()] {
		delete(graph.successors[pred], node.ID())
	}
	delete(graph.predecessors, node.ID())

}

func (graph *Graph) RemoveEdge(e gr.Edge, directed bool) {
	head, tail := e.Head(), e.Tail()
	if _, ok := graph.nodeMap[head.ID()]; !ok {
		return
	} else if _, ok := graph.nodeMap[tail.ID()]; !ok {
		return
	}

	delete(graph.successors[head.ID()], tail.ID())
	delete(graph.predecessors[tail.ID()], head.ID())
	if !directed {
		delete(graph.successors[tail.ID()], head.ID())
		delete(graph.predecessors[head.ID()], tail.ID())
	}
}

func (graph *Graph) EmptyGraph() {
	graph.successors = make(map[int]map[int]WeightedEdge)
	graph.predecessors = make(map[int]map[int]WeightedEdge)
	graph.nodeMap = make(map[int]gr.Node)
}

/* Graph implementation */

func (graph *Graph) Successors(node gr.Node) []gr.Node {
	if _, ok := graph.successors[node.ID()]; !ok {
		return nil
	}

	successors := make([]gr.Node, len(graph.successors[node.ID()]))
	i := 0
	for succ, _ := range graph.successors[node.ID()] {
		successors[i] = graph.nodeMap[succ]
		i++
	}

	return successors
}

func (graph *Graph) EdgeTo(node, succ gr.Node) gr.Edge {
	if _, ok := graph.nodeMap[node.ID()]; !ok {
		return nil
	} else if _, ok := graph.nodeMap[succ.ID()]; !ok {
		return nil
	}

	edge, ok := graph.successors[node.ID()][succ.ID()]
	if !ok {
		return nil
	}
	return edge
}

func (graph *Graph) Predecessors(node gr.Node) []gr.Node {
	if _, ok := graph.successors[node.ID()]; !ok {
		return nil
	}

	predecessors := make([]gr.Node, len(graph.predecessors[node.ID()]))
	i := 0
	for succ, _ := range graph.predecessors[node.ID()] {
		predecessors[i] = graph.nodeMap[succ]
		i++
	}

	return predecessors
}

func (graph *Graph) Neighbors(node gr.Node) []gr.Node {
	if _, ok := graph.successors[node.ID()]; !ok {
		return nil
	}

	neighbors := make([]gr.Node, len(graph.predecessors[node.ID()])+len(graph.successors[node.ID()]))
	i := 0
	for succ, _ := range graph.successors[node.ID()] {
		neighbors[i] = graph.nodeMap[succ]
		i++
	}

	for pred, _ := range graph.predecessors[node.ID()] {
		// We should only add the predecessor if it wasn't already added from successors
		if _, ok := graph.successors[node.ID()][pred]; !ok {
			neighbors[i] = graph.nodeMap[pred]
			i++
		}
	}

	return neighbors
}

func (graph *Graph) EdgeBetween(node, neigh gr.Node) gr.Edge {
	e := graph.EdgeTo(node, neigh)
	if e != nil {
		return e
	}

	e = graph.EdgeTo(neigh, node)
	if e != nil {
		return e
	}

	return nil
}

func (graph *Graph) NodeExists(node gr.Node) bool {
	_, ok := graph.nodeMap[node.ID()]

	return ok
}

func (graph *Graph) Degree(node gr.Node) int {
	if _, ok := graph.nodeMap[node.ID()]; !ok {
		return 0
	}

	return len(graph.successors[node.ID()]) + len(graph.predecessors[node.ID()])
}

func (graph *Graph) NodeList() []gr.Node {
	nodes := make([]gr.Node, len(graph.successors))
	i := 0
	for _, node := range graph.nodeMap {
		nodes[i] = node
		i++
	}

	return nodes
}

func (graph *Graph) Cost(e gr.Edge) float64 {
	if s, ok := graph.successors[e.Head().ID()]; ok {
		if we, ok := s[e.Tail().ID()]; ok {
			return we.Cost
		}
	}
	return math.Inf(1)
}

func (graph *Graph) EdgeList() []gr.Edge {
	edgeList := make([]gr.Edge, 0, len(graph.successors))
	edgeMap := make(map[int]map[int]struct{}, len(graph.successors))
	for node, succMap := range graph.successors {
		edgeMap[node] = make(map[int]struct{}, len(succMap))
		for succ, edge := range succMap {
			if doneMap, ok := edgeMap[succ]; ok {
				if _, ok := doneMap[node]; ok {
					continue
				}
			}
			edgeList = append(edgeList, edge)
			edgeMap[node][succ] = struct{}{}
		}
	}

	return edgeList
}
