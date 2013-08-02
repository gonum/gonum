package graph

type PseudomutableGraph struct {
	*CutterGraph
}

// TODO: Add ability to make undirected (difficult)

func NewPseudomutableGraph(graph Graph) PseudomutableGraph {
	ag := NewAugmentedGraph(graph)
	return PseudomutableGraph{NewCutterGraph(ag)}
}

func (graph PseudomutableGraph) NewNode(successors []int) int {
	id := graph.graph.(AugmentedGraph).NewNode(successors)

	for _, succ := range successors {
		if graph.IsCutNode(succ) {
			graph.UncutNode(node)
		}
	}

	return id
}

func (graph PseudomutableGraph) AddNode(node int, succs []int) {
	if graph.IsCutNode(node) {
		graph.UncutNode(node)
	} else {
		graph.graph.(AugmentedGraph).AddNode(node, succs)
	}

	for _, succ := range succs {
		if graph.IsCutNode(node) {
			graph.UncutNode(succ)
		}

		if graph.IsCutEdge(node, succ) {
			graph.UncutEdge(node, succ)
		}
	}

}

func (graph PseudomutableGraph) AddEdge(node, succ int) {
	if !graph.NodeExists(node) {
		return
	}

	if graph.IsCutNode(succ) {
		graph.UncutNode(succ)
	}

	if graph.IsCutEdge(node, succ) {
		graph.UncutEdge(node, succ)
	} else {
		graph.graph.(AugmentedGraph).AddEdge(node, succ)
	}
}

func (graph PseudomutableGraph) RemoveNode(node int) {
	if graph.graph.(AugmentedGraph).IsAugmentedNode(node) {
		graph.graph.(AugmentedGraph).KillAugmentedNode(node)
	} else {
		graph.CutAllEdges(node) // Need to cut all edges in case it gets re-added
		graph.CutNode(node)
	}
}

func (graph PseudomutableGraph) RemoveEdge(node, succ int) {
	ag := graph.graph.(AugmentedGraph)

	if ag.IsAugmentedEdge(node, succ) || ag.IsOverriddenEdge(node, succ) {
		ag.KillAugmentedEdge(node, succ)
	}

	// If it still exists, then we need to cut it
	if graph.IsSuccessor(node, succ) {
		graph.CutEdge(node, succ)
	}
}

func (graph PseudomutableGraph) SetEdgeCost(node, succ int, cost float64) {
	if !graph.IsSuccessor(node, succ) {
		return
	}

	graph.graph.(AugmentedGraph).SetEdgeCost(node, succ, cost)
}

func (graph PseudomutableGraph) EmptyGraph() {
	for node := range graph.NodeList() {
		graph.CutNode(node)
		graph.CutAllEdges(node)
	}
}

func (graph PseudomutableGraph) SetDirected(directed bool) {
	//TODO
}
