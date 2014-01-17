package graph

import (
	"github.com/gonum/graph/set"
	"github.com/gonum/graph/xifo"
	"math"
	"sort"
)

type Node interface {
	ID() int
}

type Edge interface {
	Head() Node
	Tail() Node
}

// A Graph implements all methods necessary to run graph-specific algorithms on it. 90% of the time you want to actually implement DirectedGraph or UndirectedGraph, since the
// default adjacency functions are (somewhat deliberately) slow.
type Graph interface {
	NodeExists(node Node) bool // Returns whether a node with the given Node is currently in the graph
	Degree(node Node) int      // Degree is equivalent to len(Successors(node)) + len(Predecessors(node)); this means that reflexive edges are counted twice
	EdgeList() []Edge          // Returns a list of all edges in the graph. Edges in EdgeList() are always directed, even when only implementing UndirectedGraph.
	NodeList() []Node          // Returns a list of all node IDs in no particular order, useful for determining things like if a graph is fully connected. The caller is free to modify this list (so don't pass a reference to your own list)
}

type UndirectedGraph interface {
	Graph
	Neighbors(node Node) []Node          // Returns all nodes connected by any edge to this node
	IsNeighbor(node, neighbor Node) bool // Returns whether neighbor is connected by an edge to node
}

type DirectedGraph interface {
	UndirectedGraph
	Successors(node Node) []Node               // Gives the nodes connected by OUTBOUND edges, if the graph is an undirected graph, this set is equal to Predecessors
	IsSuccessor(node, successor Node) bool     // If successor shows up in the list returned by Successors(node), then it's a successor. If node doesn't exist, this should always return false
	Predecessors(node Node) []Node             // Gives the nodes connected by INBOUND edges, if the graph is an undirected graph, this set is equal to Successors
	IsPredecessor(node, predecessor Node) bool // If predecessor shows up in the list returned by Predecessors(node), then it's a predecessor. If node doesn't exist, this should always return false
}

// A crunch graph forces a sparse graph to become a dense graph. That is, if the node IDs are [1,4,9,7] it would "crunch" the ids into the contiguous block [0,1,2,3]
//
// All dense graphs should have the first ID at 0
type CrunchGraph interface {
	Graph
	Crunch()
}

// A Graph that implements Coster has an actual cost between adjacent nodes, also known as a weighted graph. If a graph implements coster and a function needs to read cost (e.g. A*), this function will
// take precedence over the Uniform Cost function (all weights are 1) if "nil" is passed in for the function argument
//
// Coster only need worry about the case when an edge from node 1 to node 2 exists (i.e. node2 is a successor to node1) -- asking for the weight in any other case is considered undefined behavior.
// The only possible exception to this is in D*-Lite, if an edge previously existed and then is removed when the graph changes between steps, a suitably discouraging cost such as Inf would likely produce the best behavior.
type Coster interface {
	Cost(node1, node2 Node) float64
}

type CostGraph interface {
	Coster
	Graph
}

// A graph that implements HeuristicCoster implements a heuristic between any two given nodes. Like Coster, if a graph implements this and a function needs a heuristic cost (e.g. A*), this function will
// take precedence over the Null Heuristic (always returns 0) if "nil" is passed in for the function argument
type HeuristicCoster interface {
	Coster
	HeuristicCost(node1, node2 Node) float64 // If HeuristicCost is not intended to be used, it can be implemented as the null heuristic (always returns 0)
}

// A Mutable Graph is a graph that can be changed in an arbitrary way. It is useful for several algorithms; for instance, Johnson's Algorithm requires adding a temporary node and changing edge weights.
// Another case where this is used is computing minimum spanning trees. Since trees are graphs, a minimum spanning tree can be created using this interface.
//
// Note that just because a graph does not implement MutableGraph does not mean that this package expects it to be invariant (though even a MutableGraph should be treated as invariant while an algorithm
// is operating on it), it simply means that without this interface this package can not properly handle the graph in order to, say, fill it with a minimum spanning tree.
//
// In functions that take a MutableGraph as an argument, it should not be the same as the Graph argument as concurrent modification will likely cause problems in most cases.
//
// Mutable graphs should always record the IDs as they are represented -- which means they are sparse by nature.
type MutableGraph interface {
	CostGraph
	NewNode(successors []Node) Node       // Adds a node with an arbitrary ID, and returns the new, unique ID used
	AddNode(node Node, successors []Node) // The graph itself is responsible for adding reciprocal edges if it's undirected. Likewise, the graph itself must add any non-existant nodes listed in successors.
	AddEdge(e Edge)                       // For a digraph, adds node1->node2; the graph is free to initialize this to any value it wishes. Node1 must exist, or it will result in undefined behavior, node2 must be created by the function if absent
	SetEdgeCost(e Edge, cost float64)     // The behavior is undefined if the edge has not been created with AddEdge (or the edge was removed before this function was called). For a directed graph only sets node1->node2
	RemoveNode(node Node)                 // The graph is reponsible for removing edges to a node that is removed
	RemoveEdge(e Edge)                    // The graph is responsible for removing reciprocal edges if it's undirected
	EmptyGraph()                          // Clears the graph of all nodes and edges
	SetDirected(bool)                     // This package will only call SetDirected on an empty graph, so there's no need to worry about the case where a graph suddenly becomes (un)directed
}

// A package that contains an edge (as from EdgeList), and a Weight (as if Cost(Edge.Head(), Edge.Tail()) had been called)
type WeightedEdge struct {
	Edge
	Weight float64
}

type GonumNode int

func (node GonumNode) ID() int {
	return int(node)
}

type GonumEdge struct {
	H, T Node
}

func (edge GonumEdge) Head() Node {
	return edge.H
}

func (edge GonumEdge) Tail() Node {
	return edge.T
}

/* Slow functions to replace the guarantee of a graph being directed or undirected */

func NeighborsFunc(graph Graph) func(node Node) []Node {
	return func(node Node) []Node {
		neighbors := []Node{}

		for _, edge := range graph.EdgeList() {
			if edge.Head().ID() == node.ID() {
				neighbors = append(neighbors, edge.Tail())
			} else if edge.Tail().ID() == node.ID() {
				neighbors = append(neighbors, edge.Head())
			}
		}

		return neighbors
	}
}

func SuccessorsFunc(graph Graph) func(node Node) []Node {
	return func(node Node) []Node {
		neighbors := []Node{}

		for _, edge := range graph.EdgeList() {
			if edge.Head().ID() == node.ID() {
				neighbors = append(neighbors, edge.Tail())
			}
		}

		return neighbors
	}
}

func PredecessorsFunc(graph Graph) func(node Node) []Node {
	return func(node Node) []Node {
		neighbors := []Node{}

		for _, edge := range graph.EdgeList() {
			if edge.Tail().ID() == node.ID() {
				neighbors = append(neighbors, edge.Head())
			}
		}

		return neighbors
	}
}

func IsSuccessorFunc(graph Graph) func(Node, Node) bool {
	return func(node, succ Node) bool {
		for _, edge := range graph.EdgeList() {
			if edge.Head().ID() == node.ID() && edge.Tail().ID() == succ.ID() {
				return true
			}
		}

		return false
	}
}

func IsPredecessorFunc(graph Graph) func(Node, Node) bool {
	return func(node, pred Node) bool {
		for _, edge := range graph.EdgeList() {
			if edge.Tail().ID() == node.ID() && edge.Head().ID() == pred.ID() {
				return true
			}
		}

		return false
	}
}

func IsNeighborFunc(graph Graph) func(Node, Node) bool {
	return func(node, succ Node) bool {
		for _, edge := range graph.EdgeList() {
			if (edge.Tail().ID() == node.ID() || edge.Head().ID() == node.ID()) && (edge.Tail().ID() == succ.ID() || edge.Head().ID() == succ.ID()) {
				return true
			}
		}

		return false
	}
}

func setupFuncs(graph Graph, cost, heuristicCost func(Node, Node) float64) (successorsFunc, predecessorsFunc, neighborsFunc func(Node) []Node, isSuccessorFunc, isPredecessorFunc, isNeighborFunc func(Node, Node) bool, costFunc, heuristicCostFunc func(Node, Node) float64) {
	switch g := graph.(type) {
	case DirectedGraph:
		successorsFunc = g.Successors
		predecessorsFunc = g.Predecessors
		neighborsFunc = g.Neighbors
		isSuccessorFunc = g.IsSuccessor
		isPredecessorFunc = g.IsPredecessor
		isNeighborFunc = g.IsNeighbor
	case UndirectedGraph:
		successorsFunc = g.Neighbors
		predecessorsFunc = g.Neighbors
		neighborsFunc = g.Neighbors
		isSuccessorFunc = g.IsNeighbor
		isPredecessorFunc = g.IsNeighbor
		isNeighborFunc = g.IsNeighbor
	default:
		successorsFunc = SuccessorsFunc(graph)
		predecessorsFunc = PredecessorsFunc(graph)
		neighborsFunc = NeighborsFunc(graph)
		isSuccessorFunc = IsNeighborFunc(graph)
		isPredecessorFunc = IsNeighborFunc(graph)
		isNeighborFunc = IsNeighborFunc(graph)
	}

	if heuristicCost != nil {
		heuristicCostFunc = heuristicCost
	} else {
		if g, ok := graph.(HeuristicCoster); ok {
			heuristicCostFunc = g.HeuristicCost
		} else {
			heuristicCostFunc = NullHeuristic
		}
	}

	if cost != nil {
		costFunc = cost
	} else {
		if g, ok := graph.(Coster); ok {
			costFunc = g.Cost
		} else {
			costFunc = UniformCost
		}
	}

	return
}

/* Simple operations */

func CopyGraph(dst MutableGraph, src Graph) {
	dst.EmptyGraph()
	dst.SetDirected(false)

	var Cost func(Node, Node) float64
	if cgraph, ok := src.(Coster); ok {
		Cost = cgraph.Cost
	}

	for _, edge := range src.EdgeList() {
		if !dst.NodeExists(edge.Head()) {
			dst.AddNode(edge.Head(), []Node{edge.Tail()})
		} else {
			dst.AddEdge(edge)
		}

		if Cost != nil {
			dst.SetEdgeCost(edge, Cost(edge.Head(), edge.Tail()))
		}
	}
}

/* Basic Graph tests */

// Also known as Tarjan's Strongly Connected Components Algorithm. This returns all the strongly connected components in the graph.
//
// A strongly connected component of a graph is a set of vertices where it's possible to reach any vertex in the set from any other (meaning there's a cycle between them)
//
// Generally speaking, a directed graph where the number of strongly connected components is equal to the number of nodes is acyclic, unless you count reflexive edges as a cycle (which requires only a little extra testing)
//
// An undirected graph should end up with as many SCCs as there are "islands" (or subgraphs) of connections, meaning having more than one strongly connected component implies that your graph is not fully connected.
func Tarjan(graph Graph) (sccs [][]Node) {
	index := 0
	vStack := &xifo.GonumStack{}
	stackSet := set.NewSet()
	sccs = make([][]Node, 0)

	nodes := graph.NodeList()
	lowlinks := make(map[int]int, len(nodes))
	indices := make(map[int]int, len(nodes))

	successors, _, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)

	var strongconnect func(Node) []Node

	strongconnect = func(node Node) []Node {
		indices[node.ID()] = index
		lowlinks[node.ID()] = index
		index += 1

		vStack.Push(node)
		stackSet.Add(node.ID())

		for _, succ := range successors(node) {
			if _, ok := indices[succ.ID()]; !ok {
				strongconnect(succ)
				lowlinks[node.ID()] = int(math.Min(float64(lowlinks[node.ID()]), float64(lowlinks[succ.ID()])))
			} else if stackSet.Contains(succ) {
				lowlinks[node.ID()] = int(math.Min(float64(lowlinks[node.ID()]), float64(lowlinks[succ.ID()])))
			}
		}

		if lowlinks[node.ID()] == indices[node.ID()] {
			scc := make([]Node, 0)
			for {
				v := vStack.Pop()
				stackSet.Remove(v.(Node).ID())
				scc = append(scc, v.(Node))
				if v.(Node).ID() == node.ID() {
					return scc
				}
			}
		}

		return nil
	}

	for _, n := range nodes {
		if _, ok := indices[n.ID()]; !ok {
			sccs = append(sccs, strongconnect(n))
		}
	}

	return sccs
}

// Returns true if, starting at path[0] and ending at path[len(path)-1], all nodes between are valid neighbors. That is, for each element path[i], path[i+1] is a valid successor
//
// Special case: a nil or zero length path is considered valid (true), a path of length 1 (only one node) is the trivial case, but only if the node listed in path exists.
func IsPath(path []Node, graph Graph) bool {
	_, _, _, isSuccessor, _, _, _, _ := setupFuncs(graph, nil, nil)
	if path == nil || len(path) == 0 {
		return true
	} else if len(path) == 1 {
		return graph.NodeExists(path[0])
	}

	for i := 0; i < len(path)-1; i++ {
		if !isSuccessor(path[i], path[i+1]) {
			return false
		}
	}

	return true
}

/* Implements minimum-spanning tree algorithms; puts the resulting minimum spanning tree in the dst graph */

// Generates a minimum spanning tree with sets.
//
// As with other algorithms that use Cost, the order of precedence is Argument > Interface > UniformCost
func Prim(dst MutableGraph, graph Graph, Cost func(Node, Node) float64) {
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
	remainingNodes := set.NewSet()
	for _, node := range nlist[1:] {
		remainingNodes.Add(node.ID())
	}

	edgeList := graph.EdgeList()
	for remainingNodes.Cardinality() != 0 {
		edgeWeights := make(edgeSorter, 0)
		for _, edge := range edgeList {
			if dst.NodeExists(edge.Head()) && remainingNodes.Contains(edge.Tail().ID()) {
				edgeWeights = append(edgeWeights, WeightedEdge{Edge: edge, Weight: Cost(edge.Head(), edge.Tail())})
			}
		}

		sort.Sort(edgeWeights)
		myEdge := edgeWeights[0]

		if !dst.NodeExists(myEdge.Head()) {
			dst.AddNode(myEdge.Head(), []Node{myEdge.Tail()})
		} else {
			dst.AddEdge(myEdge.Edge)
		}
		dst.SetEdgeCost(myEdge.Edge, myEdge.Weight)

		remainingNodes.Remove(myEdge.Edge.Head())
	}

}

// Generates a minimum spanning tree for a graph using discrete.DisjointSet
//
// As with other algorithms with Cost, the precedence goes Argument > Interface > UniformCost
func Kruskal(dst MutableGraph, graph Graph, Cost func(Node, Node) float64) {
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
		edgeWeights = append(edgeWeights, WeightedEdge{Edge: edge, Weight: Cost(edge.Head(), edge.Tail())})
	}

	sort.Sort(edgeWeights)

	ds := set.NewDisjointSet()
	for _, node := range graph.NodeList() {
		ds.MakeSet(node.ID())
	}

	for _, edge := range edgeWeights {
		if s1, s2 := ds.Find(edge.Edge.Head().ID()), ds.Find(edge.Edge.Tail().ID); s1 != s2 {
			ds.Union(s1, s2)
			if !dst.NodeExists(edge.Edge.Head()) {
				dst.AddNode(edge.Edge.Head(), []Node{edge.Edge.Tail()})
			} else {
				dst.AddEdge(edge.Edge)
			}
			dst.SetEdgeCost(edge.Edge, edge.Weight)
		}
	}
}

/* Control flow graph stuff */

// A dominates B if and only if the only path through B travels through A
//
// This returns all possible dominators for all nodes, it does not prune for strict dominators, immediate dominators etc
//
// The int map[int]*set.Set is the node's ID
func Dominators(start Node, graph Graph) map[int]*set.Set {
	allNodes := set.NewSet()
	nlist := graph.NodeList()
	dominators := make(map[int]*set.Set, len(nlist))
	for _, node := range nlist {
		allNodes.Add(node.ID())
	}

	_, predecessors, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)

	for _, node := range nlist {
		dominators[node.ID()] = set.NewSet()
		if node.ID() == start.ID() {
			dominators[node.ID()].Add(start.ID())
		} else {
			dominators[node.ID()].Copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == start.ID() {
				continue
			}
			preds := predecessors(node)
			if len(preds) == 0 {
				continue
			}
			tmp := set.NewSet().Copy(dominators[preds[0].ID()])
			for _, pred := range preds[1:] {
				tmp.Intersection(tmp, dominators[pred.ID()])
			}

			dom := set.NewSet()
			dom.Add(node.ID())

			dom.Union(dom, tmp)
			if !set.Equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

// A Postdominates B if and only if all paths from B travel through A
//
// This returns all possible post-dominators for all nodes, it does not prune for strict postdominators, immediate postdominators etc
func PostDominators(end Node, graph Graph) map[int]*set.Set {
	successors, _, _, _, _, _, _, _ := setupFuncs(graph, nil, nil)
	allNodes := set.NewSet()
	nlist := graph.NodeList()
	dominators := make(map[int]*set.Set, len(nlist))
	for _, node := range nlist {
		allNodes.Add(node.ID())
	}

	for _, node := range nlist {
		dominators[node.ID()] = set.NewSet()
		if node.ID() == end.ID() {
			dominators[node.ID()].Add(end.ID())
		} else {
			dominators[node.ID()].Copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == end.ID() {
				continue
			}
			succs := successors(node)
			if len(succs) == 0 {
				continue
			}
			tmp := set.NewSet().Copy(dominators[succs[0].ID()])
			for _, succ := range succs[1:] {
				tmp.Intersection(tmp, dominators[succ.ID()])
			}

			dom := set.NewSet()
			dom.Add(node.ID())

			dom.Union(dom, tmp)
			if !set.Equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

/* Purely internal data structures and functions (mostly for sorting) */

/** Sorts a list of edges by weight, agnostic to repeated edges as well as direction **/

type edgeSorter []WeightedEdge

func (el edgeSorter) Len() int {
	return len(el)
}

func (el edgeSorter) Less(i, j int) bool {
	return el[i].Weight < el[j].Weight
}

func (el edgeSorter) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}
