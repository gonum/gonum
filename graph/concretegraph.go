package graph

import (
	"sort"
)

// A GonumGraph is a very generalized graph that can handle an arbitrary number of vertices and edges -- as well as act as either directed or undirected.
//
// Internally, it uses a map of successors AND predecessors, to speed up some operations (such as getting all successors/predecessors). It also speeds up thing like adding edges (assuming both edges exist).
//
// However, its generality is also its weakness (and partially a flaw in needing to satisfy MutableGraph). For most purposes, creating your own graph is probably better. For instance, see discrete.TileGraph for an example
// of an immutable 2D grid of tiles that also implements the Graph interface, but would be more suitable if all you needed was a simple undirected 2D grid.
type GonumGraph struct {
	successors   map[int]map[int]float64
	predecessors map[int]map[int]float64
	directed     bool
}

func NewGonumGraph(directed bool) *GonumGraph {
	return &GonumGraph{
		successors:   make(map[int]map[int]float64),
		predecessors: make(map[int]map[int]float64),
		directed:     directed,
	}
}

func NewPreAllocatedGonumGraph(directed bool, numVertices int) *GonumGraph {
	return &GonumGraph{
		successors:   make(map[int]map[int]float64, numVertices),
		predecessors: make(map[int]map[int]float64, numVertices),
		directed:     directed,
	}
}

/* Mutable Graph implementation */

func (graph *GonumGraph) NewNode(successors []int) (id int) {
	nodes := sort.IntSlice(graph.NodeList())
	sort.Sort(&nodes)
	for i, node := range nodes {
		if i != node {
			graph.AddNode(i, successors)
			return i
		}
	}

	newID := len(nodes)
	graph.AddNode(newID, successors)
	return newID
}

func (graph *GonumGraph) AddNode(id int, successors []int) {
	if _, ok := graph.successors[id]; ok {
		return
	}

	graph.successors[id] = make(map[int]float64, len(successors))
	if !graph.directed {
		graph.predecessors[id] = make(map[int]float64, len(successors))
	} else {
		graph.predecessors[id] = make(map[int]float64)
	}
	for _, succ := range successors {
		graph.successors[id][succ] = 1.0

		// Always add the reciprocal node to the graph
		if _, ok := graph.successors[succ]; !ok {
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

func (graph *GonumGraph) AddEdge(id, successor int) {
	if _, ok := graph.successors[id]; !ok {
		return
	}

	if _, ok := graph.successors[successor]; !ok {
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

func (graph *GonumGraph) SetEdgeCost(id, successor int, cost float64) {
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

func (graph *GonumGraph) RemoveNode(id int) {
	if _, ok := graph.successors[id]; ok {
		return
	}

	for succ, _ := range graph.successors[id] {
		delete(graph.predecessors[succ], id)
	}
	delete(graph.successors, id)

	for pred, _ := range graph.predecessors[id] {
		delete(graph.successors[pred], id)
	}
	delete(graph.predecessors, id)

}

func (graph *GonumGraph) RemoveEdge(id, succ int) {
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
}

func (graph *GonumGraph) SetDirected(directed bool) {
	if len(graph.successors) > 0 {
		return
	}
	graph.directed = directed
}

/* Graph implementation */

func (graph *GonumGraph) Successors(id int) []int {
	if _, ok := graph.successors[id]; !ok {
		return nil
	}

	successors := make([]int, len(graph.successors[id]))
	for succ, _ := range graph.successors[id] {
		successors = append(successors, succ)
	}

	return successors
}

func (graph *GonumGraph) IsSuccessor(id, succ int) bool {
	if _, ok := graph.successors[id]; !ok {
		return false
	}

	_, ok := graph.successors[id][succ]

	return ok
}

func (graph *GonumGraph) Predecessors(id int) []int {
	if _, ok := graph.successors[id]; !ok {
		return nil
	}

	predecessors := make([]int, len(graph.predecessors[id]))
	for pred, _ := range graph.predecessors[id] {
		predecessors = append(predecessors, pred)
	}

	return predecessors
}

func (graph *GonumGraph) IsPredecessor(id, pred int) bool {
	if _, ok := graph.successors[id]; !ok {
		return false
	}

	_, ok := graph.predecessors[id][pred]

	return ok
}

func (graph *GonumGraph) IsAdjacent(id, neighbor int) bool {
	if _, ok := graph.successors[id]; !ok {
		return false
	}

	_, succ := graph.predecessors[id][neighbor]
	_, pred := graph.predecessors[id][neighbor]

	return succ || pred
}

func (graph *GonumGraph) NodeExists(id int) bool {
	_, ok := graph.successors[id]

	return ok
}

func (graph *GonumGraph) Degree(id int) int {
	if _, ok := graph.successors[id]; !ok {
		return 0
	}

	return len(graph.successors[id]) + len(graph.predecessors[id])
}

func (graph *GonumGraph) EdgeList() [][2]int {
	eList := make([][2]int, 0, len(graph.successors))
	for id, succMap := range graph.successors {
		for succ, _ := range succMap {
			eList = append(eList, [2]int{id, succ})
		}
	}

	return eList
}

func (graph *GonumGraph) NodeList() []int {
	nodes := make([]int, 0, len(graph.successors))
	for node, _ := range graph.successors {
		nodes = append(nodes, node)
	}

	return nodes
}

func (graph *GonumGraph) IsDirected() bool {
	return graph.directed
}

func (graph *GonumGraph) Cost(id, succ int) float64 {
	return graph.successors[id][succ]
}
