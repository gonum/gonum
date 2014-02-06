package graph

// All a node needs to do is identify itself. This allows the user to pass in nodes more interesting than an int,
// but also allow us to reap the benefits of having a map-storable, ==able type.
type Node interface {
	ID() int
}

// Allows edges to do something more interesting that just be a group of nodes. While the methods are called Head and Tail,
// they are not considered directed unless the given interface specifies otherwise
type Edge interface {
	Head() Node
	Tail() Node
}

// A Graph ensures the behavior of an undirected graph, necessary to run certain algorithms on it.
//
// The Graph interface is directed. This means that EdgeList() should return an edge where Head always goes towards Tail. If your graph is undirected and you only maintain edges for one direction,
// simply return two edges for each one of your edges, with the Head and Tail swapped in each one.
type Graph interface {
	NodeExists(node Node) bool           // Returns whether a node with the given Node is currently in the graph
	Degree(node Node) int                // Degree is equivalent to len(Successors(node)) + len(Predecessors(node)); this means that reflexive edges are counted twice
	NodeList() []Node                    // Returns a list of all node IDs in no particular order, useful for determining things like if a graph is fully connected. The caller is free to modify this list (so don't pass a reference to your own list)
	Neighbors(node Node) []Node          // Returns all nodes connected by any edge to this node
	IsNeighbor(node, neighbor Node) bool // Returns whether neighbor is connected by an edge to node
}

// Directed graphs are characterized by having seperable Heads and Tails in their edges. That is, if node1 goes to node2, that does not necessarily imply that node2 goes to node1.
//
// While it's possible for a directed graph to have fully reciprocal edges (i.e. the graph is symmetric) -- it is not required to be. The graph is also required to implement UndirectedGraph
// because it can be useful to know all neighbors regardless of direction; not because this graph treats directed graphs as special cases of undirected ones (the truth is, in fact, the opposite)
type DirectedGraph interface {
	Graph
	Successors(node Node) []Node               // Gives the nodes connected by OUTBOUND edges, if the graph is an undirected graph, this set is equal to Predecessors
	IsSuccessor(node, successor Node) bool     // If successor shows up in the list returned by Successors(node), then it's a successor. If node doesn't exist, this should always return false
	Predecessors(node Node) []Node             // Gives the nodes connected by INBOUND edges, if the graph is an undirected graph, this set is equal to Successors
	IsPredecessor(node, predecessor Node) bool // If predecessor shows up in the list returned by Predecessors(node), then it's a predecessor. If node doesn't exist, this should always return false
}

// Returns all undirected edges in the graph
type EdgeLister interface {
	EdgeList() []Edge
}

type EdgeListGraph interface {
	Graph
	EdgeLister
}

type DirectedEdgeLister interface {
	DirectedEdgeList() []Edge
}

type DirectedEdgeListGraph interface {
	Graph
	DirectedEdgeLister
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
// If no edge exists between node1 and node2, the cost should be taken to be +inf (can be gotten by math.Inf(1))
type Coster interface {
	Cost(node1, node2 Node) float64
}

// Guarantees that something implementing Coster is also a Graph
type CostGraph interface {
	Coster
	Graph
}

// A graph that implements HeuristicCoster implements a heuristic between any two given nodes. Like Coster, if a graph implements this and a function needs a heuristic cost (e.g. A*), this function will
// take precedence over the Null Heuristic (always returns 0) if "nil" is passed in for the function argument
type HeuristicCoster interface {
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

// A DStarGraph is a special interface that allows the DStarLite function to be used on a graph
//
// D*-lite is an algorithm that allows for the graph representation to change when actions are taken, whether this be from actions taken by the agent or simply new information gathered.
// As such, there's a Move function, that allows the graph to take into account an agent moving to the next node. This is always followed by a call to ChangedEdges.
//
// Traditionally in D*-lite, the algorithm would scan every edge to see if the cost changed, and then update its information if it detected any changes. This slightly remixed step
// allows the graph to provide notification of any changes, and even provide an alternate cost function if it needs to. This can be used to speed up the algorithm significantly
// since the graph no longer has to scan for changes, and only updates when told to. If changedEdges is nil or of len 0, no updates will be performed. If changedEdges is not nil, it
// will update the internal representation. If newCostFunc is non-nil it will be swapped with dStar's current cost function if and only if changedEdges is non-nil/len>0, however,
// newCostFunc is not required to be non-nil if updates are present. DStar will continue using the current cost function if that is the case.
type DStarGraph interface {
	Graph
	Move(target Node)
	ChangedEdges() (newCostFunc func(Node, Node) float64, changedEdges []Edge)
}

// A function that returns the cost from one node to another
type CostFunc func(Node, Node) float64
