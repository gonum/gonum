package discrete

import (
	"sort"
)

type GonumGraph struct {
	successors   map[int]map[int]float64
	predecessors map[int]map[int]float64
	vertices     *Set
	directed     bool
}

func NewGonumGraph(directed bool) *GonumGraph {
	return &GonumGraph{
		successors:   make(map[int]map[int]float64),
		predecessors: make(map[int]map[int]float64),
		vertices:     NewSet(),
		directed:     directed,
	}
}

func NewPreAllocatedGonumGraph(directed bool, numVertices int) *GonumGraph {
	return &GonumGraph{
		successors:   make(map[int]map[int]float64, numVertices),
		predecessors: make(map[int]map[int]float64, numVertices),
		vertices:     NewSet(),
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
	if graph.vertices.Contains(id) {
		return
	}
	graph.vertices.Add(id)
	graph.successors[id] = make(map[int]float64, len(successors))
	if !graph.directed {
		graph.predecessors[id] = make(map[int]float64, len(successors))
	} else {
		graph.predecessors[id] = make(map[int]float64)
	}
	for _, succ := range successors {
		graph.successors[id][succ] = 1.0

		// Always add the reciprocal node to the graph
		if !graph.vertices.Contains(succ) {
			graph.vertices.Add(succ)
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
	if !graph.vertices.Contains(id) {
		return
	}

	if !graph.vertices.Contains(successor) {
		graph.vertices.Add(successor)
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
	if !graph.vertices.Contains(id) {
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
	if !graph.vertices.Contains(id) {
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

	graph.vertices.Remove(id)

}

func (graph *GonumGraph) RemoveEdge(id, succ int) {
	if !graph.vertices.Contains(id) {
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
	if graph.vertices.Cardinality() == 0 {
		return
	}
	graph.vertices = NewSet()
	graph.successors = make(map[int]map[int]float64)
	graph.predecessors = make(map[int]map[int]float64)
}

func (graph *GonumGraph) SetDirected(directed bool) {
	if graph.vertices.Cardinality() > 0 {
		return
	}
	graph.directed = directed
}

/* Graph implementation */

func (graph *GonumGraph) Successors(id int) []int {
	if !graph.vertices.Contains(id) {
		return nil
	}

	successors := make([]int, len(graph.successors[id]))
	for succ, _ := range graph.successors[id] {
		successors = append(successors, succ)
	}

	return successors
}

func (graph *GonumGraph) IsSuccessor(id, succ int) bool {
	if !graph.vertices.Contains(id) {
		return false
	}

	_, ok := graph.successors[id][succ]

	return ok
}

func (graph *GonumGraph) Predecessors(id int) []int {
	if !graph.vertices.Contains(id) {
		return nil
	}

	predecessors := make([]int, len(graph.predecessors[id]))
	for pred, _ := range graph.predecessors[id] {
		predecessors = append(predecessors, pred)
	}

	return predecessors
}

func (graph *GonumGraph) IsPredecessor(id, pred int) bool {
	if !graph.vertices.Contains(id) {
		return false
	}

	_, ok := graph.predecessors[id][pred]

	return ok
}

func (graph *GonumGraph) IsAdjacent(id, neighbor int) bool {
	if !graph.vertices.Contains(id) {
		return false
	}

	_, succ := graph.predecessors[id][neighbor]
	_, pred := graph.predecessors[id][neighbor]

	return succ || pred
}

func (graph *GonumGraph) NodeExists(id int) bool {
	return graph.vertices.Contains(id)
}

func (graph *GonumGraph) Degree(id int) int {
	if !graph.vertices.Contains(id) {
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
	rawNodes := graph.vertices.Elements()
	nodes := make([]int, 0, len(rawNodes))
	for _, rawNode := range rawNodes {
		nodes = append(nodes, rawNode.(int))
	}

	return nodes
}

func (graph *GonumGraph) IsDirected() bool {
	return graph.directed
}

func (graph *GonumGraph) Cost(id, succ int) float64 {
	return graph.successors[id][succ]
}
