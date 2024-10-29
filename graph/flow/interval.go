package flow

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/linear"
)

// Returns the set of intervals given by the control flow graph.
// IDs must be in reverse postorder
func Intervals(g graph.Directed, eid int64) []*Interval {
	var worklist linear.NodeQueue
	var intervals []*Interval
	var interval Interval
	var visited map[graph.Node]bool
	head := g.Node(eid)

	worklist.Enqueue(head)
	for {
		// exits when the worklist is empty
		if worklist.Len() <= 0 {
			break
		}

		n := worklist.Dequeue()
		visited = interval.findInterval(n, g, visited)
		intervals = append(intervals, &interval)
		// add all interval nodes to interval
		// can pass the map to the function to make it efficient

	}

	return intervals
}

// An Interval I(h) is the maximal, single entry subgraph for which h (head)
// is the entry node and in which all closed paths contain h.
type Interval struct {
	head  graph.Node
	nodes []graph.Node
}

func (i *Interval) Head() graph.Node {
	return i.head
}

func (i *Interval) Nodes() graph.Nodes {

}

// Returns the edge given 2 node id's if the edge exists.
// Else it returns null.
func (i *Interval) Edge(uid, vid int64) graph.Edge {

}

// Finds all interval nodes.
// Nodes are added to the interval if all their predecessors are in
// the interval or they are the header node.
func (i *Interval) findInterval(n graph.Node, g graph.Directed, visited map[graph.Node]bool) map[graph.Node]bool {
	i.nodes = append(i.nodes, n)

	for {
		// get the next node in reverse postorder
	}

	return visited
}

// Returns the nodes as a stack so they can be extracted in reverse postorder.
func dfsPostorder(g graph.Directed, eid int64, ns *linear.NodeStack) {
	succs := g.From(eid)
	for {
		if !succs.Next() {
			break
		}

		succ := succs.Node()
		dfsPostorder(g, succ.ID(), ns)
	}

	n := g.Node(eid)
	ns.Push(n)
}

// Extracts all nodes in stack into an array
func reversePostorder(ns linear.NodeStack) []graph.Node {

}
