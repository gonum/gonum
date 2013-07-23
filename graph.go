package discrete

import (
	"container/heap"
	"sort"
)

type Graph interface {
	Successors(node int) []int            // Gives the Successors connected by OUTBOUND edges, if the graph is a digraph
	IsSuccessor(node, successor int) bool // Is Successor
	Predecessors(node int) []int
	IsPredecessor(node, predecessor int) bool
	IsAdjacent(node, neighbor int) bool // IsSuccessor || IsPredecessor
	NodeExists(node int) bool
	Degree(node int) int // Degree is equivalent to len(Successors(node)) + len(Predecessors(node))
	EdgeList() [][2]int
	NodeList() []int // Returns a list of all node IDs in no particular order, useful for determining things like if a graph is fully connected
	IsDirected() bool
}

type Coster interface {
	Cost(node1, node2 int) float64 // A non-weighted graph should just return 1 (or some uniform positive value). Returns an error if the two nodes are not adjacent. A digraph returns the cost for node1->node2
}

type HeuristicCoster interface {
	HeuristicCost(node1, node2 int) float64 // If HeuristicCost is not intended to be used, it can be implemented as the null heuristic (always returns 0)
}

type MutableGraph interface {
	Graph
	AddNode(id int, successors []int)                 // The graph itself is responsible for adding reciprocal edges if it's undirected
	AddEdge(node1, node2 int)                         // For a digraph, adds node1->node2; the graph is free to initialize this to any value it wishes. An error should be returned if the edge has not been initialized with AddEdge
	SetEdgeCost(node1, node2 int, cost float64) error // Returns an error if the edge has not been created with AddEdge (or the edge was removed before this function was called)
	RemoveNode(node int)                              // The graph is reponsible for removing edges to a node that is removed
	RemoveEdge(node1, node2 int)                      // The graph is responsible for removing reciprocal edges if it's undirected
	EmptyGraph()                                      // Clears the graph of all nodes and edges
	SetDirected(bool)                                 // This package will only call SetDirected on the empty graph, so there's no need to worry about the case where a graph suddenly becomes (un)directed
}

// Returns true if, starting at path[0] and ending at path[len(path)-1], all nodes are traversable.
//
// Special case: a nil or zero length path is considered valid (true), a path of length 1 (only one node) is the trivial case
func IsPath(path []int, graph Graph) bool {
	if path == nil || len(path) <= 1 {
		return true
	}

	for i := 0; i < len(path)-1; i++ {
		if !graph.IsSuccessor(path[i], path[i+1]) {
			return false
		}
	}

	return true
}

func DepthFirstSearch(start, goal int, graph Graph) []int {
	closedSet := NewSet()
	openSet := Stack([]interface{}{start})
	predecessor := make(map[int]int)

	for !openSet.IsEmpty() {
		c, err := openSet.Pop()
		if err != nil {
			return nil
		}

		curr := c.(int)

		if closedSet.Contains(curr) {
			continue
		}

		if curr == goal {
			return rebuildPath(predecessor, goal)
		}

		closedSet.Add(curr)

		for _, neighbor := range graph.Successors(curr) {
			if closedSet.Contains(neighbor) {
				continue
			}

			predecessor[neighbor] = curr
			openSet.Push(neighbor)
		}
	}

	return nil
}

func NullHeuristic(a, b int) float64 {
	return 0.0
}

func UniformCost(a, b int) float64 {
	return 1.0
}

// Returns an ordered list consisting of the nodes between start and goal. The path will be the shortest path assuming the function heuristicCost is admissible
// The second return value is the cost
//
// Cost and HeuristicCost take precedence for evaluating cost/heuristic distance. If one is not present (i.e. nil) the function will check the graph's interface for the respective interface:
// Coster for Cost and HeuristicCoster for HeuristicCost. If the correct one is present, it will use the graph's function for evaluation.
//
// Finally, if neither the argument nor the interface is present, the function will assume discrete.UniformCost for Cost and discrete.NullHeuristic for HeuristicCost
//
// To run Uniform Cost Search, run A* with the NullHeuristic
//
// To run Breadth First Search, run A* with both the NullHeuristic and UniformCost (or any cost function that returns a uniform positive value)
func AStar(start, goal int, graph Graph, Cost, HeuristicCost func(int, int) float64) (path []int, cost float64) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	if HeuristicCost == nil {
		if hgraph, ok := graph.(HeuristicCoster); ok {
			HeuristicCost = hgraph.HeuristicCost
		} else {
			HeuristicCost = HeuristicCost
		}
	}

	closedSet := make(map[int]internalNode)
	openSet := make(aStarPriorityQueue, 0)
	heap.Init(openSet)
	node := internalNode{start, 0, HeuristicCost(start, goal)}
	heap.Push(openSet, node)
	predecessor := make(map[int]int)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(internalNode)

		// This isn't in most implementations of A*, it's a restructuring of the step "if node not in openSet, add it"
		// Instead of searching to check, we see if we already evaluated it. If we have we can ignore it
		if _, ok := closedSet[curr.int]; ok {
			continue
		}

		if curr.int == goal {
			return rebuildPath(predecessor, goal), curr.gscore
		}

		closedSet[curr.int] = curr

		for _, neighbor := range graph.Successors(curr.int) {
			g := curr.gscore + Cost(curr.int, neighbor)
			if _, ok := closedSet[neighbor]; ok && g >= closedSet[neighbor].gscore {
				continue
			}

			if _, ok := closedSet[neighbor]; !ok || g < closedSet[neighbor].gscore {
				node = internalNode{neighbor, g, g + HeuristicCost(neighbor, goal)}
				predecessor[node.int] = curr.int
				heap.Push(openSet, node)
			}
		}
	}

	return nil, 0.0
}

func rebuildPath(predecessors map[int]int, goal int) []int {
	path := []int{goal}
	curr := goal
	for prev, ok := predecessors[curr]; ok; prev, ok = predecessors[curr] {
		path = append([]int{prev}, path...)
		curr = prev
	}

	return path
}

// Finds the shortest path to every (connected) node in the graph from a single source -- no edges may have negative weights
func Dijkstra(source int, graph Graph) (paths map[int][]int, costs map[int]float64) {
	return nil, nil
}

// Same as Dijkstra, but handles negative edge weights
func BellmanFord(source int, graph Graph) (paths map[int][]int, costs map[int]float64) {
	return nil, nil
}

/* Basic Graph tests */

// Checks if every node in the graph has a degree of at least one, unless it's an empty graph or a graph with a single node in which case it's considered trivially connected
func FullyConnected(graph Graph) bool {
	nlist := graph.NodeList()
	if nlist == nil || len(nlist) <= 1 {
		return true
	}

	for _, node := range graph.NodeList() {
		if graph.Degree(node) == 0 {
			return false
		}
	}

	return true
}

/* Implements minimum-spanning tree algorithms; puts the resulting minimum spanning tree in the dst graph */

func Prim(dst MutableGraph, graph Graph, Cost func(int, int) float64) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	dst.EmptyGraph()
	dst.SetDirected(false)

	nlist := graph.NodeList()

	if nlist == nil || len(nlist) == 0 {
		return
	}

	dst.AddNode(nlist[0], nil)
	remainingNodes := NewSet()
	for _, node := range nlist[1:] {
		remainingNodes.Add(node)
	}

	edgeList := graph.EdgeList()
	for remainingNodes.Cardinality() != 0 {
		edgeWeights := make(edgeSorter, 0)
		for _, edge := range edgeList {
			if dst.NodeExists(edge[0]) && remainingNodes.Contains(edge[1]) {
				edgeWeights = append(edgeWeights, edgeWeight{edge, Cost(edge[0], edge[1])})
			}
		}

		sort.Sort(edgeWeights)
		myEdge := edgeWeights[0]

		dst.AddNode(myEdge.edge[0], []int{myEdge.edge[1]})
		dst.SetEdgeCost(myEdge.edge[0], myEdge.edge[1], myEdge.weight)

		remainingNodes.Remove(myEdge.edge[1])
	}

}

func Kruskal(dst MutableGraph, graph Graph, Cost func(int, int) float64) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	dst.EmptyGraph()
	dst.SetDirected(false)

	edgeList := graph.EdgeList()
	edgeWeights := make(edgeSorter, 0, len(edgeList))
	for _, edge := range edgeList {
		edgeWeights = append(edgeWeights, edgeWeight{edge, Cost(edge[0], edge[1])})
	}

	sort.Sort(edgeWeights)

	ds := NewDisjointSet()
	for _, node := range graph.NodeList() {
		ds.MakeSet(node)
	}

	for _, edge := range edgeWeights {
		if s1, s2 := ds.Find(edge.edge[0]), ds.Find(edge.edge[1]); s1 != s2 {
			ds.Union(s1, s2)
			dst.AddNode(edge.edge[0], []int{edge.edge[1]})
			dst.SetEdgeCost(edge.edge[0], edge.edge[1], edge.weight)
		}
	}
}

/* Control flow graph stuff */

func Dominators(start int, graph Graph) map[int]*Set {
	allNodes := NewSet()
	nlist := graph.NodeList()
	dominators := make(map[int]*Set, len(nlist))
	for _, node := range nlist {
		allNodes.Add(node)
	}

	for _, node := range nlist {
		dominators[node] = NewSet()
		if node == start {
			dominators[node].Add(start)
		} else {
			dominators[node].Copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node == start {
				continue
			}
			preds := graph.Predecessors(node)
			if len(preds) == 0 {
				continue
			}
			tmp := NewSet().Copy(dominators[preds[0]])
			for _, pred := range preds[1:] {
				tmp.Intersection(tmp, dominators[pred])
			}

			dom := NewSet()
			dom.Add(node)

			dom.Union(dom, tmp)
			if !dom.Equal(dominators[node]) {
				dominators[node] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

type edgeWeight struct {
	edge   [2]int
	weight float64
}

type edgeSorter []edgeWeight

func (el edgeSorter) Len() int {
	return len(el)
}

func (el edgeSorter) Less(i, j int) bool {
	return el[i].weight < el[j].weight
}

func (el edgeSorter) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}

type internalNode struct {
	int
	gscore, fscore float64
}

type aStarPriorityQueue []internalNode

func (pq aStarPriorityQueue) Less(i, j int) bool {
	return -pq[i].fscore < -pq[j].fscore
}

func (pq aStarPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq aStarPriorityQueue) Len() int {
	return len(pq)
}

func (pq aStarPriorityQueue) Push(x interface{}) {
	el, ok := x.(internalNode)
	if !ok {
		return
	}
	pq = append(pq, el)
}

func (pq aStarPriorityQueue) Pop() interface{} {
	if len(pq) == 0 {
		return nil
	}

	x := pq[len(pq)-1]
	pq = pq[:len(pq)-1]

	return x
}
