package graph

import (
	"github.com/gonum/discrete/set"
	"sort"
)

// TODO: Make ability to be undirected

type AugmentedGraph struct {
	graph                 Graph
	augmentedNodes        *set.Set
	augmentedSuccessors   map[int]map[int]float64
	augmentedPredecessors map[int]map[int]float64
}

func NewAugmentedGraph(graph Graph) *AugmentedGraph {
	return &AugmentedGraph{graph, set.NewSet(), make(map[int]map[int]float64), make(map[int]map[int]float64)}
}

func (graph *AugmentedGraph) Successors(node int) (succs []int) {
	if graph.augmentedNodes.Contains(node) {
		succs = make([]int, 0, len(graph.augmentedSuccessors[node]))
		for node, _ := range graph.augmentedSuccessors[node] {
			succs = append(succs, node)
		}
	} else {
		succs = graph.graph.Successors(node)
		for node, _ := range graph.augmentedSuccessors[node] {
			succs = append(succs, node)
		}
	}
	return succs
}

func (graph *AugmentedGraph) IsSuccessor(node, successor int) bool {
	var ok bool
	if graph.augmentedNodes.Contains(node) {
		_, ok = graph.augmentedSuccessors[node][successor]
	}
	return ok || graph.graph.IsSuccessor(node, successor)
}

func (graph *AugmentedGraph) Predecessors(node int) (preds []int) {
	if graph.augmentedNodes.Contains(node) {
		preds = make([]int, 0, len(graph.augmentedPredecessors[node]))
		for node, _ := range graph.augmentedPredecessors[node] {
			preds = append(preds, node)
		}
	} else {
		preds = graph.graph.Predecessors(node)
		for node, _ := range graph.augmentedPredecessors[node] {
			preds = append(preds, node)
		}
	}
	return preds
}

func (graph *AugmentedGraph) IsPredecessor(node, predecessor int) bool {
	var ok bool
	if graph.augmentedNodes.Contains(node) {
		_, ok = graph.augmentedPredecessors[node][predecessor]
	}
	return ok || graph.graph.IsPredecessor(node, predecessor)
}

func (graph *AugmentedGraph) IsAdjacent(node, neighbor int) bool {
	return graph.IsSuccessor(node, neighbor) || graph.IsPredecessor(node, neighbor)
}

func (graph *AugmentedGraph) NodeExists(node int) bool {
	return graph.augmentedNodes.Contains(node) || graph.graph.NodeExists(node)
}

func (graph *AugmentedGraph) Degree(node int) int {
	return len(graph.Successors(node)) + len(graph.Predecessors(node))
}

func (graph *AugmentedGraph) EdgeList() [][2]int {
	eList := graph.graph.EdgeList()
	for node, succList := range graph.augmentedSuccessors {
		for succ, _ := range succList {
			eList = append(eList, [2]int{node, succ})
		}
	}

	return eList
}

func (graph *AugmentedGraph) NodeList() []int {
	nodeList := graph.graph.NodeList()
	for _, rawNode := range graph.augmentedNodes.Elements() {
		nodeList = append(nodeList, rawNode.(int))
	}

	return nodeList
}

func (graph *AugmentedGraph) IsDirected() bool {
	return true
}

func (graph *AugmentedGraph) NewNode(succs []int) int {
	nodes := sort.IntSlice(graph.NodeList())
	sort.Sort(nodes)

	for i, node := range nodes {
		if i != node {
			graph.AddNode(i, succs)
			return i
		}
	}

	// This shouldn't ever occur, just to keep the compiler from complaining
	return -1
}

func (graph *AugmentedGraph) AddNode(node int, succs []int) {
	if graph.augmentedNodes.Contains(node) || graph.graph.NodeExists(node) {
		return
	}

	if !graph.augmentedNodes.Contains(node) {
		graph.augmentedNodes.Add(node)
	}

	for succ := range succs {
		if graph.IsSuccessor(node, succ) {
			continue
		} else if !graph.graph.NodeExists(node) {
			graph.augmentedNodes.Add(node)
		}

		if _, ok := graph.augmentedSuccessors[node]; !ok {
			graph.augmentedSuccessors[node] = make(map[int]float64, 1)
		}

		if _, ok := graph.augmentedPredecessors[succ]; !ok {
			graph.augmentedPredecessors[succ] = make(map[int]float64, 1)
		}

		graph.augmentedSuccessors[node][succ] = 1.0
		graph.augmentedPredecessors[succ][node] = 1.0
	}
}

func (graph *AugmentedGraph) AddEdge(node, succ int) {
	if graph.IsSuccessor(node, succ) {
		return
	}

	if !graph.NodeExists(node) {
		return
	}

	if !graph.NodeExists(succ) {
		graph.AddNode(succ, nil)
	}

	if _, ok := graph.augmentedSuccessors[node]; !ok {
		graph.augmentedSuccessors[node] = make(map[int]float64, 1)
	}

	if _, ok := graph.augmentedPredecessors[succ]; !ok {
		graph.augmentedPredecessors[succ] = make(map[int]float64, 1)
	}

	graph.augmentedSuccessors[node][succ] = 1.0
	graph.augmentedPredecessors[succ][node] = 1.0
}

func (graph *AugmentedGraph) SetEdgeCost(node, succ int, cost float64) {
	if !graph.IsSuccessor(node, succ) {
		return
	}

	if _, ok := graph.augmentedSuccessors[node]; !ok {
		graph.augmentedSuccessors[node] = make(map[int]float64, 1)
	}

	graph.augmentedSuccessors[node][succ] = cost

	if _, ok := graph.augmentedPredecessors[succ]; !ok {
		graph.augmentedPredecessors[succ] = make(map[int]float64, 1)
	}

	graph.augmentedPredecessors[succ][node] = cost
}

func (graph *AugmentedGraph) IsAugmentedNode(node int) bool {
	return graph.augmentedNodes.Contains(node)
}

func (graph *AugmentedGraph) IsAugmentedEdge(node, succ int) bool {
	if augset, ok := graph.augmentedSuccessors[node]; ok {
		_, ok = augset[succ]

		return ok && !graph.graph.IsSuccessor(node, succ)
	}

	return false
}

func (graph *AugmentedGraph) IsOverriddenEdge(node, succ int) bool {
	if augset, ok := graph.augmentedSuccessors[node]; ok {
		_, ok = augset[succ]

		return ok && graph.graph.IsSuccessor(node, succ)
	}

	return false
}

func (graph *AugmentedGraph) KillAugmentedNode(node int) {
	graph.augmentedNodes.Remove(node)

	for succ, _ := range graph.augmentedSuccessors[node] {
		if succ == node {
			continue
		}
		delete(graph.augmentedPredecessors[succ], node)
	}

	delete(graph.augmentedSuccessors, node)

	for pred, _ := range graph.augmentedPredecessors[node] {
		if pred == node {
			continue
		}
		delete(graph.augmentedSuccessors[pred], node)
	}

	delete(graph.augmentedPredecessors, node)
}

func (graph *AugmentedGraph) KillAugmentedEdge(node, succ int) {
	if eList, ok := graph.augmentedSuccessors[node]; ok {
		delete(eList, succ)
	}

	if pList, ok := graph.augmentedPredecessors[succ]; ok {
		delete(pList, node)
	}
}

func (graph *AugmentedGraph) Clear() {
	graph.augmentedNodes = set.NewSet()
	graph.augmentedSuccessors = make(map[int]map[int]float64)
	graph.augmentedPredecessors = make(map[int]map[int]float64)
}
