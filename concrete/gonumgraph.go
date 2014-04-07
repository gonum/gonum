package concrete

import (
	"math"
	"sort"

	"github.com/gonum/graph"
)

// A simple int alias.
type Node int

func (node Node) ID() int {
	return int(node)
}

// Just a collection of two nodes
type Edge struct {
	H, T graph.Node
}

func (edge Edge) Head() graph.Node {
	return edge.H
}

func (edge Edge) Tail() graph.Node {
	return edge.T
}

type WeightedEdge struct {
	graph.Edge
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
	nodeMap      map[int]graph.Node
}

func NewGraph() *Graph {
	return &Graph{
		successors:   make(map[int]map[int]WeightedEdge),
		predecessors: make(map[int]map[int]WeightedEdge),
		nodeMap:      make(map[int]graph.Node),
	}
}

func NewPreAllocatedGraph(numVertices int) *Graph {
	return &Graph{
		successors:   make(map[int]map[int]WeightedEdge, numVertices),
		predecessors: make(map[int]map[int]WeightedEdge, numVertices),
		nodeMap:      make(map[int]graph.Node, numVertices),
	}
}

/* Mutable Graph implementation */

func (gr *Graph) NewNode() (node graph.Node) {
	nodeList := gr.NodeList()
	ids := make([]int, len(nodeList))
	for i, node := range nodeList {
		ids[i] = node.ID()
	}

	nodes := sort.IntSlice(ids)
	sort.Sort(&nodes)
	for i, node := range nodes {
		if i != node {
			gr.AddNode(Node(i))
			return Node(i)
		}
	}

	newID := len(nodes)
	gr.AddNode(Node(newID))
	return Node(newID)
}

func (gr *Graph) AddNode(node graph.Node) {
	if _, ok := gr.nodeMap[node.ID()]; ok {
		return
	}

	gr.nodeMap[node.ID()] = node
	gr.successors[node.ID()] = make(map[int]WeightedEdge)
	gr.predecessors[node.ID()] = make(map[int]WeightedEdge)
}

func (gr *Graph) AddEdge(e graph.Edge, cost float64, directed bool) {
	head, tail := e.Head(), e.Tail()
	gr.AddNode(head)
	gr.AddNode(tail)

	gr.successors[head.ID()][tail.ID()] = WeightedEdge{Edge: e, Cost: cost}
	gr.predecessors[tail.ID()][head.ID()] = WeightedEdge{Edge: e, Cost: cost}
	if !directed {
		gr.successors[tail.ID()][head.ID()] = WeightedEdge{Edge: e, Cost: cost}
		gr.predecessors[head.ID()][tail.ID()] = WeightedEdge{Edge: e, Cost: cost}
	}
}

func (gr *Graph) RemoveNode(node graph.Node) {
	if _, ok := gr.nodeMap[node.ID()]; !ok {
		return
	}
	delete(gr.nodeMap, node.ID())

	for succ, _ := range gr.successors[node.ID()] {
		delete(gr.predecessors[succ], node.ID())
	}
	delete(gr.successors, node.ID())

	for pred, _ := range gr.predecessors[node.ID()] {
		delete(gr.successors[pred], node.ID())
	}
	delete(gr.predecessors, node.ID())

}

func (gr *Graph) RemoveEdge(e graph.Edge, directed bool) {
	head, tail := e.Head(), e.Tail()
	if _, ok := gr.nodeMap[head.ID()]; !ok {
		return
	} else if _, ok := gr.nodeMap[tail.ID()]; !ok {
		return
	}

	delete(gr.successors[head.ID()], tail.ID())
	delete(gr.predecessors[tail.ID()], head.ID())
	if !directed {
		delete(gr.successors[tail.ID()], head.ID())
		delete(gr.predecessors[head.ID()], tail.ID())
	}
}

func (gr *Graph) EmptyGraph() {
	gr.successors = make(map[int]map[int]WeightedEdge)
	gr.predecessors = make(map[int]map[int]WeightedEdge)
	gr.nodeMap = make(map[int]graph.Node)
}

/* Graph implementation */

func (gr *Graph) Successors(node graph.Node) []graph.Node {
	if _, ok := gr.successors[node.ID()]; !ok {
		return nil
	}

	successors := make([]graph.Node, len(gr.successors[node.ID()]))
	i := 0
	for succ, _ := range gr.successors[node.ID()] {
		successors[i] = gr.nodeMap[succ]
		i++
	}

	return successors
}

func (gr *Graph) EdgeTo(node, succ graph.Node) graph.Edge {
	if _, ok := gr.nodeMap[node.ID()]; !ok {
		return nil
	} else if _, ok := gr.nodeMap[succ.ID()]; !ok {
		return nil
	}

	edge, ok := gr.successors[node.ID()][succ.ID()]
	if !ok {
		return nil
	}
	return edge
}

func (gr *Graph) Predecessors(node graph.Node) []graph.Node {
	if _, ok := gr.successors[node.ID()]; !ok {
		return nil
	}

	predecessors := make([]graph.Node, len(gr.predecessors[node.ID()]))
	i := 0
	for succ, _ := range gr.predecessors[node.ID()] {
		predecessors[i] = gr.nodeMap[succ]
		i++
	}

	return predecessors
}

func (gr *Graph) Neighbors(node graph.Node) []graph.Node {
	if _, ok := gr.successors[node.ID()]; !ok {
		return nil
	}

	neighbors := make([]graph.Node, len(gr.predecessors[node.ID()])+len(gr.successors[node.ID()]))
	i := 0
	for succ, _ := range gr.successors[node.ID()] {
		neighbors[i] = gr.nodeMap[succ]
		i++
	}

	for pred, _ := range gr.predecessors[node.ID()] {
		// We should only add the predecessor if it wasn't already added from successors
		if _, ok := gr.successors[node.ID()][pred]; !ok {
			neighbors[i] = gr.nodeMap[pred]
			i++
		}
	}

	return neighbors
}

func (gr *Graph) EdgeBetween(node, neigh graph.Node) graph.Edge {
	e := gr.EdgeTo(node, neigh)
	if e != nil {
		return e
	}

	e = gr.EdgeTo(neigh, node)
	if e != nil {
		return e
	}

	return nil
}

func (gr *Graph) NodeExists(node graph.Node) bool {
	_, ok := gr.nodeMap[node.ID()]

	return ok
}

func (gr *Graph) Degree(node graph.Node) int {
	if _, ok := gr.nodeMap[node.ID()]; !ok {
		return 0
	}

	return len(gr.successors[node.ID()]) + len(gr.predecessors[node.ID()])
}

func (gr *Graph) NodeList() []graph.Node {
	nodes := make([]graph.Node, len(gr.successors))
	i := 0
	for _, node := range gr.nodeMap {
		nodes[i] = node
		i++
	}

	return nodes
}

func (gr *Graph) Cost(e graph.Edge) float64 {
	if s, ok := gr.successors[e.Head().ID()]; ok {
		if we, ok := s[e.Tail().ID()]; ok {
			return we.Cost
		}
	}
	return math.Inf(1)
}

func (gr *Graph) EdgeList() []graph.Edge {
	edgeList := make([]graph.Edge, 0, len(gr.successors))
	edgeMap := make(map[int]map[int]struct{}, len(gr.successors))
	for node, succMap := range gr.successors {
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
