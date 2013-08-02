package graph

import (
	"github.com/gonum/discrete/set"
)

type CutterGraph struct {
	graph      Graph
	nodeCutset *set.Set
	edgeCutset map[int]*set.Set
}

func NewCutterGraph(graph Graph) *CutterGraph {
	return &CutterGraph{graph, set.NewSet(), make(map[int]*set.Set)}
}

func (graph *CutterGraph) Successors(node int) []int {
	if graph.nodeCutset.Contains(node) {
		return nil
	}

	successors := graph.graph.Successors(node)
	if _, ok := graph.edgeCutset[node]; !ok {
		return successors
	}

	realSuccs := make([]int, 0, len(successors)-len(*graph.edgeCutset[node]))
	for _, succ := range successors {
		if !graph.edgeCutset[node].Contains(succ) {
			realSuccs = append(realSuccs, succ)
		}
	}

	return realSuccs
}

func (graph *CutterGraph) IsSuccessor(node, successor int) bool {
	if graph.nodeCutset.Contains(node) || graph.nodeCutset.Contains(successor) {
		return false
	}

	cutset, ok := graph.edgeCutset[node]

	return !ok && !cutset.Contains(successor) && graph.graph.IsSuccessor(node, successor)
}

func (graph *CutterGraph) Predecessors(node int) []int {
	if graph.nodeCutset.Contains(node) {
		return nil
	}

	predecessors := graph.graph.Predecessors(node)
	if _, ok := graph.edgeCutset[node]; !ok {
		return predecessors
	}

	realPreds := make([]int, 0, len(predecessors)-len(*graph.edgeCutset[node]))
	for _, pred := range predecessors {
		if !graph.edgeCutset[node].Contains(pred) {
			realPreds = append(realPreds, pred)
		}
	}

	return realPreds
}

func (graph *CutterGraph) IsPredecessor(node, predecessor int) bool {
	if graph.nodeCutset.Contains(node) || graph.nodeCutset.Contains(predecessor) {
		return false
	}

	cutset, ok := graph.edgeCutset[node]

	return !ok && !cutset.Contains(predecessor) && graph.graph.IsPredecessor(node, predecessor)
}

func (graph *CutterGraph) IsAdjacent(node, neighbor int) bool {
	return graph.IsSuccessor(node, neighbor) || graph.IsPredecessor(node, neighbor)
}

func (graph *CutterGraph) NodeExists(node int) bool {
	return !graph.nodeCutset.Contains(node) && graph.graph.NodeExists(node)
}

func (graph *CutterGraph) Degree(node int) int {
	if graph.nodeCutset.Contains(node) {
		return 0
	}

	return len(graph.Successors(node)) + len(graph.Predecessors(node))
}

func (graph *CutterGraph) EdgeList() [][2]int {
	eList := graph.graph.EdgeList()
	realEList := make([][2]int, 0)
	for _, edge := range eList {
		if !graph.edgeCutset[edge[0]].Contains(edge[1]) && !graph.nodeCutset.Contains(edge[0]) && !graph.nodeCutset.Contains(edge[1]) {
			realEList = append(realEList, edge)
		}
	}

	return realEList
}

func (graph *CutterGraph) NodeList() []int {
	nodeList := graph.graph.NodeList()
	realNodeList := make([]int, 0, len(nodeList)-len(*graph.nodeCutset))
	for _, node := range nodeList {
		if !graph.nodeCutset.Contains(node) {
			realNodeList = append(realNodeList, node)
		}
	}

	return realNodeList
}

func (graph *CutterGraph) IsDirected() bool {
	return graph.graph.IsDirected()
}

func (graph *CutterGraph) CutNode(node int) {
	graph.nodeCutset.Add(node)
}

func (graph *CutterGraph) UncutNode(node int) {
	graph.nodeCutset.Remove(node)
}

func (graph *CutterGraph) CutEdge(node, succ int) {
	if _, ok := graph.edgeCutset[node]; !ok {
		graph.edgeCutset[node] = set.NewSet()
	}
	graph.edgeCutset[node].Add(succ)
	if graph.graph.IsDirected() {
		if _, ok := graph.edgeCutset[succ]; !ok {
			graph.edgeCutset[succ] = set.NewSet()
		}
		graph.edgeCutset[succ].Add(node)
	}
}

func (graph *CutterGraph) UncutEdge(node, succ int) {
	if _, ok := graph.edgeCutset[node]; ok {
		graph.edgeCutset[node].Remove(succ)
	}
	if graph.graph.IsDirected() {
		if _, ok := graph.edgeCutset[succ]; ok {
			graph.edgeCutset[succ].Remove(node)
		}
	}
}

func (graph *CutterGraph) IsCutNode(node int) bool {
	return graph.nodeCutset.Contains(node)
}

func (graph *CutterGraph) IsCutEdge(node, succ int) bool {
	if cutset, ok := graph.edgeCutset[node]; ok {
		return cutset.Contains(succ)
	}

	return false
}

func (graph *CutterGraph) CutAllEdges(node int) {
	for _, succ := range graph.Successors(node) {
		graph.CutEdge(node, succ)
	}

	for _, pred := range graph.Predecessors(node) {
		graph.CutEdge(pred, node)
	}
}

func (graph *CutterGraph) UncutAllEdges(node int) {
	delete(graph.edgeCutset, node)

	for _, cutset := range graph.edgeCutset {
		if cutset.Contains(node) {
			cutset.Remove(node)
		}
	}

}
