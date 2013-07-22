package discrete

import ()

type Graph interface {
	Neighbors(node int) []int // Gives the neighbors FROM this node, if the graph is a digraph
	NodeExists(node int) bool
	Degree(node int) int                    // Degree is equivalent to len(Neighbors(node)) for an undirected graph, but may be higher in a digraph
	NodeList() []int                        // Returns a list of all node IDs, useful for determining things like if a graph is fully connected
	Cost(node1, node2 int) (float64, error) // A non-weighted graph should just return 1. Returns an error if the two nodes are not adjacent. A digraph returns the cost for node1->node2
	IsDirected() bool
}

type MutableGraph interface {
	Graph
	AddNode(neighbors []int) int                      // Returns the ID of the new node, the graph itself is responsible for adding
	AddEdge(node1, node2 int)                         // For a digraph, adds node1->node2
	SetEdgeCost(node1, node2 int, cost float64) error // Returns an error if the edge has not been created with AddEdge (or the edge was removed before this function was called)
	RemoveNode(node int)
	RemoveEdge(node int)
	EmptyGraph() // Clears the graph of all nodes and edges
}

/* Searches */

// Only Difference between BFS and DFS is stack or queue
// So we just implement it once and pass in a stack or queue depending
//
// If maxDepth is -1, we don't depth limit searches
func xFirstSearch(start, goal int, graph Graph, dataStruct XInFirstOut, maxDepth int) []int {
	return nil
}

func BreadthFirstSearch(start, goal int, graph Graph) []int {
	return xFirstSearch(start, goal, graph, make(Queue, 0), -1)
}

func DepthFirstSearch(start, goal int, graph Graph) []int {
	return xFirstSearch(start, goal, graph, make(Stack, 0), -1)
}

// DFID is BFS in terms of node expansion order, but it's slower (redoes the search a lot) with the tradeoff of only using as much memory as DFS
// It works by running DFS at successively lower depths
func DepthFirstIterativeDeepening(start, goal int, graph Graph, maxDepth int) []int {
	if maxDepth < 0 {
		// Just keep going until we find something, with DFID you're kind of SOL if you don't know the max depth possible
		for {
			path := xFirstSearch(start, goal, graph, make(Stack, 0), i)
			if path != nil {
				return path
			}
		}
	} else {
		for i := 0; i <= maxDepth; i++ {
			path := xFirstSearch(start, goal, graph, make(Stack, 0), i)
			if path != nil {
				return path
			}
		}
	}

	return nil
}

/* Technically UCS and A* are just BFS with a priority queue, however we can't implement a priority queue with the XInFirstOut interface since Push necessarily requires an extra argument (the priority) */

// Returns an ordered list consisting of the nodes between start and goal. The path will be the shortest path assuming the function heuristicCost is admissible
// The second return value is the cost
func AStar(start, goal int, graph Graph, heuristicCost func(int, int) float64) (path []int, cost float64) {
	return nil, 0
}

func UniformCostSearch(start, goal int, graph Graph) (path []int, cost float64) {
	// UCS is just AStar with the null heuristic
	return AStar(start, goal, graph, func(int, int) float64 { return 0.0 })
}

func Dijkstra(source int, graph Graph) (paths map[int][]int, costs map[int]float64) {
	return nil, nil
}

func BellmanFord(source int, graph Graph) (paths map[int][]int, costs map[int]float64) {
	return nil, nil
}

/* Basic Graph tests */

func Acyclic(graph Graph) bool {
	return true
}

func FullyConnected(graph Graph) bool {
	return true
}

/* Implements minimum-spanning tree algorithms */

func Prim(dst MutableGraph, graph Graph) {

}

func Kruskal(dst MutableGraph, graph Graph) {

}

/* Control flow graph stuff */

func Dominates(dominator, dominated int, graph Graph) bool {
	return true
}
