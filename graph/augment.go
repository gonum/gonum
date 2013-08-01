package graph

import (
	"github.com/gonum/discrete/set"
	"sort"
)

type AugmentedGraph struct {
	graph                 Graph
	augmentedNodes        *set.Set
	augmentedSuccessors   map[int]map[int]float64
	augmentedPredecessors map[int]map[int]float64
	directed              bool
}

func NewAugmentedGraph(graph Graph, directed bool) *AugmentedGraph {
	return &AugmentedGraph{graph, set.NewSet(), make(map[int]map[int]float64), make(map[int]map[int]float64), directed}
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
	return graph.directed
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

}
