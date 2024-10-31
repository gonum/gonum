// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import (
	"log"

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
	var ns linear.NodeStack
	visited := make(map[int64]bool)
	dfsPostorder(g, eid, &ns, visited)
	reversePostorderNodes, reversePostorderMap := reversePostorder(ns)
	n := g.Node(eid)
	worklist.Enqueue(n)
	for worklist.Len() > 0 {
		n = worklist.Dequeue()
		interval.findInterval(&n, g, reversePostorderMap, reversePostorderNodes)
		log.Println(interval.nodes)
		intervals = append(intervals, &interval)
		if n == nil {
			break
		}

		worklist.Enqueue(n)
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

// func (i *Interval) Nodes() graph.Nodes {

// }

// Returns the edge given 2 node id's if the edge exists.
// Else it returns null.
// func (i *Interval) Edge(uid, vid int64) graph.Edge {

// }

// Finds all interval nodes.
// Nodes are added to the interval if all their predecessors are in
// the interval or they are the header node.
func (i *Interval) findInterval(head *graph.Node, g graph.Directed, reversePostorderMap map[*graph.Node]int, reversePostorderArray []*graph.Node) {
	i.head = *head
	var nq linear.NodeQueue
	nq.Enqueue(*head)
	i.nodes = append(i.nodes, *head)
	intervalMap := make(map[graph.Node]bool)
	intervalMap[*head] = true
	var node graph.Node
	for nq.Len() > 0 {
		node = nq.Dequeue()
		succs := g.From(node.ID())

		for succs.Next() {
			if intervalMap[succs.Node()] {
				continue
			}

			preds := g.To(succs.Node().ID())
			predsLength := preds.Len()
			x := 0
			for preds.Next() {
				if intervalMap[preds.Node()] {
					x++
				}
			}

			if x == predsLength {
				nq.Enqueue(succs.Node())
				i.nodes = append(i.nodes, succs.Node())
				intervalMap[succs.Node()] = true
			}
		}
	}
}

// Put nodes into the stack in postorder.
func dfsPostorder(g graph.Directed, eid int64, ns *linear.NodeStack, visited map[int64]bool) {
	succs := g.From(eid)
	visited[eid] = true
	for {
		if !succs.Next() {
			break
		}

		succ := succs.Node()
		if !visited[succ.ID()] {
			dfsPostorder(g, succ.ID(), ns, visited)
		}

	}

	n := g.Node(eid)
	ns.Push(n)
}

// Extracts all nodes in the stack into an array in reverse postorder.
func reversePostorder(ns linear.NodeStack) ([]*graph.Node, map[*graph.Node]int) {
	var nodes []*graph.Node
	nodePosition := make(map[*graph.Node]int)
	stackLength := ns.Len()
	for i := 0; i < stackLength; i++ {
		n := ns.Pop()
		nodePosition[&n] = i
		nodes = append(nodes, &n)
	}

	return nodes, nodePosition
}
