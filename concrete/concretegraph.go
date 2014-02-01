package concrete

import (
	"sort"

	gr "github.com/gonum/graph"
)

// A simple int alias.
type GonumNode int

func (node GonumNode) ID() int {
	return int(node)
}

// Just a collection of two nodes
type GonumEdge struct {
	H, T gr.Node
}

func (edge GonumEdge) Head() gr.Node {
	return edge.H
}

func (edge GonumEdge) Tail() gr.Node {
	return edge.T
}

// A GonumGraph is a very generalized graph that can handle an arbitrary number of vertices and edges -- as well as act as either directed or undirected.
//
// Internally, it uses a map of successors AND predecessors, to speed up some operations (such as getting all successors/predecessors). It also speeds up thing like adding edges (assuming both edges exist).
//
// However, its generality is also its weakness (and partially a flaw in needing to satisfy MutableGraph). For most purposes, creating your own graph is probably better. For instance, see discrete.TileGraph for an example
// of an immutable 2D grid of tiles that also implements the Graph interface, but would be more suitable if all you needed was a simple undirected 2D grid.
type GonumGraph struct {
	successors   map[int]map[int]float64
	predecessors map[int]map[int]float64
	nodeMap      map[int]gr.Node
	directed     bool
}

func NewGonumGraph(directed bool) *GonumGraph {
	return &GonumGraph{
		successors:   make(map[int]map[int]float64),
		predecessors: make(map[int]map[int]float64),
		nodeMap:      make(map[int]gr.Node),
		directed:     directed,
	}
}

func NewPreAllocatedGonumGraph(directed bool, numVertices int) *GonumGraph {
	return &GonumGraph{
		successors:   make(map[int]map[int]float64, numVertices),
		predecessors: make(map[int]map[int]float64, numVertices),
		nodeMap:      make(map[int]gr.Node, numVertices),
		directed:     directed,
	}
}

/* Mutable Graph implementation */

func (graph *GonumGraph) NewNode(successors []gr.Node) (node gr.Node) {
	nodeList := graph.NodeList()
	ids := make([]int, len(nodeList))
	for i, node := range nodeList {
		ids[i] = node.ID()
	}

	nodes := sort.IntSlice(ids)
	sort.Sort(&nodes)
	for i, node := range nodes {
		if i != node {
			graph.AddNode(GonumNode(i), successors)
			return GonumNode(i)
		}
	}

	newID := len(nodes)
	graph.AddNode(GonumNode(newID), successors)
	return GonumNode(newID)
}

func (graph *GonumGraph) AddNode(node gr.Node, successors []gr.Node) {
	id := node.ID()
	if _, ok := graph.successors[id]; ok {
		return
	}

	graph.nodeMap[id] = node

	graph.successors[id] = make(map[int]float64, len(successors))
	if !graph.directed {
		graph.predecessors[id] = make(map[int]float64, len(successors))
	} else {
		graph.predecessors[id] = make(map[int]float64)
	}
	for _, successor := range successors {
		succ := successor.ID()
		graph.successors[id][succ] = 1.0

		// Always add the reciprocal node to the graph
		if _, ok := graph.successors[succ]; !ok {
			graph.nodeMap[succ] = successor
			graph.predecessors[succ] = make(map[int]float64)
			graph.successors[succ] = make(map[int]float64)
		}

		graph.predecessors[succ][id] = 1.0

		// But only add the reciprocal edge if we're undirected
		if !graph.directed {
			graph.successors[succ][id] = 1.0
			graph.predecessors[id][succ] = 1.0
		}
	}
}

func (graph *GonumGraph) AddEdge(e gr.Edge) {
	id := e.Head().ID()
	successor := e.Tail().ID()
	if _, ok := graph.successors[id]; !ok {
		return
	}

	if _, ok := graph.successors[successor]; !ok {
		graph.nodeMap[successor] = e.Tail()
		graph.successors[successor] = make(map[int]float64)
		graph.predecessors[successor] = make(map[int]float64)
	}

	graph.successors[id][successor] = 1.0
	graph.predecessors[successor][id] = 1.0

	if !graph.directed {
		graph.successors[successor][id] = 1.0
		graph.predecessors[id][successor] = 1.0
	}
}

func (graph *GonumGraph) SetEdgeCost(e gr.Edge, cost float64) {
	id := e.Head().ID()
	successor := e.Tail().ID()
	// Normally I'd use graph.vertices.Contains(id) as above, but this is equivalent and a bit easier to read here
	if _, ok := graph.successors[id]; !ok {
		return
	} else if _, ok := graph.successors[id][successor]; !ok {
		return
	}
	graph.successors[id][successor] = cost
	graph.predecessors[successor][id] = cost

	// By the spec, only the empty graph will be toggled between directed and undirected. Therefore we can be sure the reciprocal edge exists
	if !graph.directed {
		graph.successors[successor][id] = cost
		graph.predecessors[id][successor] = cost
	}
}

func (graph *GonumGraph) RemoveNode(node gr.Node) {
	id := node.ID()
	if _, ok := graph.successors[id]; !ok {
		return
	}
	delete(graph.nodeMap, id)

	for succ, _ := range graph.successors[id] {
		delete(graph.predecessors[succ], id)
	}
	delete(graph.successors, id)

	for pred, _ := range graph.predecessors[id] {
		delete(graph.successors[pred], id)
	}
	delete(graph.predecessors, id)

}

func (graph *GonumGraph) RemoveEdge(e gr.Edge) {
	id := e.Head().ID()
	succ := e.Tail().ID()
	if _, ok := graph.successors[id]; !ok {
		return
	} else if _, ok := graph.successors[succ]; !ok {
		return
	}

	delete(graph.successors[id], succ)
	delete(graph.predecessors[succ], id)
	if !graph.directed {
		delete(graph.predecessors[id], succ)
		delete(graph.successors[succ], id)
	}
}

func (graph *GonumGraph) EmptyGraph() {
	if len(graph.successors) == 0 {
		return
	}
	graph.successors = make(map[int]map[int]float64)
	graph.predecessors = make(map[int]map[int]float64)
	graph.nodeMap = make(map[int]gr.Node)
}

func (graph *GonumGraph) SetDirected(directed bool) {
	if len(graph.successors) > 0 {
		return
	}
	graph.directed = directed
}

/* Graph implementation */

func (graph *GonumGraph) Successors(node gr.Node) []gr.Node {
	id := node.ID()
	if _, ok := graph.successors[id]; !ok {
		return nil
	}

	successors := make([]gr.Node, 0, len(graph.successors[id]))
	for succ, _ := range graph.successors[id] {
		successors = append(successors, graph.nodeMap[succ])
	}

	return successors
}

func (graph *GonumGraph) IsSuccessor(node, successor gr.Node) bool {
	succ := successor.ID()
	id := node.ID()
	if _, ok := graph.successors[id]; !ok {
		return false
	}

	_, ok := graph.successors[id][succ]

	return ok
}

func (graph *GonumGraph) Predecessors(node gr.Node) []gr.Node {
	id := node.ID()
	if _, ok := graph.successors[id]; !ok {
		return nil
	}

	predecessors := make([]gr.Node, 0, len(graph.predecessors[id]))
	for pred, _ := range graph.predecessors[id] {
		predecessors = append(predecessors, graph.nodeMap[pred])
	}

	return predecessors
}

func (graph *GonumGraph) IsPredecessor(node, predecessor gr.Node) bool {
	id := node.ID()
	pred := predecessor.ID()
	if _, ok := graph.successors[id]; !ok {
		return false
	}

	_, ok := graph.predecessors[id][pred]

	return ok
}

func (graph *GonumGraph) Neighbors(node gr.Node) []gr.Node {
	id := node.ID()
	if _, ok := graph.successors[id]; !ok {
		return nil
	}

	neighbors := make([]gr.Node, 0, len(graph.predecessors[id])+len(graph.successors[id]))
	for succ, _ := range graph.successors[id] {
		neighbors = append(neighbors, graph.nodeMap[succ])
	}

	for pred, _ := range graph.predecessors[id] {
		// We should only add the predecessor if it wasn't already added from successors
		if _, ok := graph.successors[id][pred]; !ok {
			neighbors = append(neighbors, graph.nodeMap[pred])
		}
	}

	return neighbors
}

func (graph *GonumGraph) IsNeighbor(node, neigh gr.Node) bool {
	id := node.ID()
	neighbor := neigh.ID()
	if _, ok := graph.successors[id]; !ok {
		return false
	}

	_, succ := graph.predecessors[id][neighbor]
	_, pred := graph.predecessors[id][neighbor]

	return succ || pred
}

func (graph *GonumGraph) NodeExists(node gr.Node) bool {
	_, ok := graph.successors[node.ID()]

	return ok
}

func (graph *GonumGraph) Degree(node gr.Node) int {
	id := node.ID()
	if _, ok := graph.successors[id]; !ok {
		return 0
	}

	d := len(graph.successors[id])
	if graph.directed {
		return d + len(graph.predecessors[id])
	}
	if _, ok := graph.successors[id][id]; ok {
		d++
	}
	return d
}

func (graph *GonumGraph) EdgeList() []gr.Edge {
	eList := make([]gr.Edge, 0, len(graph.successors))
	for id, succMap := range graph.successors {
		for succ, _ := range succMap {
			if graph.directed || id <= succ {
				eList = append(eList, GonumEdge{graph.nodeMap[id], graph.nodeMap[succ]})
			}
		}
	}

	return eList
}

func (graph *GonumGraph) NodeList() []gr.Node {
	nodes := make([]gr.Node, 0, len(graph.successors))
	for _, node := range graph.nodeMap {
		nodes = append(nodes, node)
	}

	return nodes
}

func (graph *GonumGraph) IsDirected() bool {
	return graph.directed
}

func (graph *GonumGraph) Cost(node, succ gr.Node) float64 {
	return graph.successors[node.ID()][succ.ID()]
}
