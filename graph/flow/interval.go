package flow

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/linear"
)

/*
1. Establish a set H for header nodes and initialize it with no, the
unique entry node for the graph.

2. For h in H find l(h) as follows:
	2.1. Put h in l(h) as the first element of l(h).

	2.2. Add to l(h) any node all of whose immediate predecessors
	are already in l(h).

	2.3. Repeat 2.2 until no more nodes can be added to l(h).

3. Add to H all nodes in G which are not already in H and which
are not in l(h) but which have immediate predecessors in l(h).
Therefore a node is added to H the first time any (but not all)
of its immediate predecessors become members of an interval.

4. Add l(h) to a set S of intervals being developed.

5. Select the next unprocessed node in H and repeat steps 2, 3, 4, 5.
When there are no more unprocessed nodes in H, the procedure
terminates
*/

// Returns the set of intervals given by the control flow graph.
// IDs must be in reverse postorder
func Intervals(g graph.Directed, eid int64) []*Interval {
	var worklist linear.NodeQueue
	var intervals []*Interval
	var interval Interval
	var inInterval map[graph.Node]bool
	var ns linear.NodeStack
	var intNodeMap map[graph.Node]int
	dfsPostorder(g, eid, &ns)
	reversePostorderNodes, reversePostorderMap := reversePostorder(ns)
	head := g.Node(eid)

	worklist.Enqueue(head)
	for {
		// exits when the worklist is empty
		if worklist.Len() <= 0 {
			break
		}

		n := worklist.Dequeue()
		worklist.Enqueue(interval.findInterval(n, g, reversePostorderMap, reversePostorderNodes))
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
func (i *Interval) findInterval(n graph.Node, g graph.Directed, reversePostorderMap map[*graph.Node]int, reversePostorderArray []*graph.Node) graph.Node {
	i.head = n
	i.nodes = append(i.nodes, n)
	var intervalMap map[graph.Node]bool
	intervalMap[n] = true
	for {
		nPos := reversePostorderMap[&n]
		if nPos > len(reversePostorderArray) {
			break
		}

		n = *reversePostorderArray[nPos+1]
		preds := g.To(n.ID())
		x := 0
		for preds.Next() {
			if intervalMap[preds.Node()] {
				x++
			}
		}

		if 0 < x && x < preds.Len() {
			// n is a new interval header
			return n
		}

		// n is a part of the interval
		i.nodes = append(i.nodes, n)
		intervalMap[n] = true
	}

	// run out of nodes
	return nil
}

// Put nodes into the stack in postorder.
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

// Extracts all nodes in the stack into an array in reverse postorder.
func reversePostorder(ns linear.NodeStack) ([]*graph.Node, map[*graph.Node]int) {
	var nodes []*graph.Node
	var nodePosition map[*graph.Node]int
	for i := 0; i < ns.Len(); i++ {
		n := ns.Pop()
		nodePosition[&n] = i
		nodes = append(nodes, &n)
	}

	return nodes, nodePosition
}
