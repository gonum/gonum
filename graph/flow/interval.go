// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import (
	"maps"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/linear"
	"gonum.org/v1/gonum/graph/iterator"
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

(https://dl.acm.org/doi/pdf/10.1145/360018.360025)
*/

// Returns the set of intervals given by the directed graph.
func Intervals(g graph.Directed, eid int64) IntervalGraph {
	var worklist linear.NodeQueue
	var intervals []*Interval
	var ns linear.NodeStack
	visited := make(map[int64]bool)
	inInterval := make(map[int64]graph.Node)

	dfsPostorder(g, eid, &ns, visited)
	reversePostorderNodes := reversePostorder(ns)
	n := g.Node(eid)
	worklist.Enqueue(n)

	for worklist.Len() > 0 {
		var interval Interval
		n = worklist.Dequeue()
		maps.Copy(inInterval, interval.findInterval(&n, g))
		intervals = append(intervals, &interval)
		if n == nil {
			break
		}

		for _, node := range reversePostorderNodes {
			if inInterval[(*node).ID()] != nil {
				continue
			}

			preds := g.To((*node).ID())
			predsLength := preds.Len()
			x := 0
			for preds.Next() {
				if interval.nodes[preds.Node().ID()] != nil {
					x++
				}
			}

			if x > 0 && x < predsLength {
				worklist.Enqueue(*node)
				break
			}
		}
	}

	ig := linkIntervals(intervals, g)
	return ig
}

// An Interval I(h) is the maximal, single entry subgraph for which h (head)
// is the entry node and in which all closed paths contain h.
type Interval struct {
	head  graph.Node
	nodes map[int64]graph.Node
	from  map[int64]map[int64]graph.Edge
}

// Returns header node for an interval.
func (i *Interval) Head() graph.Node {
	return i.head
}

// Returns a node iterator for an interval.
func (i *Interval) Nodes() graph.Nodes {
	if len(i.nodes) == 0 {
		return graph.Empty
	}
	return iterator.NewNodes(i.nodes)
}

// Returns the edge given 2 node id's if the edge exists.
// Else it returns null.
func (i *Interval) Edge(uid, vid int64) graph.Edge {
	edge, ok := i.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}

func (i *Interval) From(id int64) graph.Nodes {
	if len(i.from[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByEdge(i.nodes, i.from[id])
}

func (i *Interval) HasEdgeBetween(xid int64, yid int64) bool {
	if _, ok := i.from[xid][yid]; ok {
		return true
	}
	_, ok := i.from[yid][xid]
	return ok
}

func (i *Interval) Node(id int64) graph.Node {
	return i.nodes[id]
}

// Finds all interval nodes.
// Nodes are added to the interval if all their predecessors are in
// the interval or they are the header node.
func (i *Interval) findInterval(head *graph.Node, g graph.Directed) map[int64]graph.Node {
	i.head = *head
	var nq linear.NodeQueue
	nq.Enqueue(*head)
	i.nodes = make(map[int64]graph.Node)
	i.nodes[(*head).ID()] = *head
	var node graph.Node
	for nq.Len() > 0 {
		node = nq.Dequeue()
		succs := g.From(node.ID())

		for succs.Next() {
			if i.nodes[succs.Node().ID()] != nil {
				continue
			}

			preds := g.To(succs.Node().ID())
			predsLength := preds.Len()
			x := 0
			for preds.Next() {
				if i.nodes[preds.Node().ID()] != nil {
					x++
				}
			}

			if x == predsLength {
				nq.Enqueue(succs.Node())
				i.nodes[succs.Node().ID()] = succs.Node()
			}
		}
	}

	return i.nodes
}

// Put nodes into a stack in postorder.
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
func reversePostorder(ns linear.NodeStack) []*graph.Node {
	var nodes []*graph.Node
	stackLength := ns.Len()
	for i := 0; i < stackLength; i++ {
		n := ns.Pop()
		nodes = append(nodes, &n)
	}

	return nodes
}

// Contains the intervals and the edges between the intervals
type IntervalGraph struct {
	Intervals []*Interval
	from      map[int64]map[int64]graph.Edge
}

// Computes the internal and external edges for the intervals.
func linkIntervals(intervals []*Interval, g graph.Directed) IntervalGraph {
	ig := IntervalGraph{
		Intervals: intervals,
		from:      make(map[int64]map[int64]graph.Edge),
	}

	for _, interval := range intervals {
		interval.from = make(map[int64]map[int64]graph.Edge)
		for _, node := range interval.nodes {
			succs := g.From(node.ID())
			for succs.Next() {
				succNode := succs.Node()
				if interval.nodes[succNode.ID()] != nil {
					if interval.from[node.ID()] == nil {
						interval.from[node.ID()] = make(map[int64]graph.Edge)
					}

					interval.from[node.ID()][succNode.ID()] = g.Edge(node.ID(), succNode.ID())
				} else {
					if ig.from[node.ID()] == nil {
						ig.from[node.ID()] = make(map[int64]graph.Edge)
					}

					ig.from[node.ID()][succNode.ID()] = g.Edge(node.ID(), succNode.ID())
				}
			}
		}
	}

	return ig
}
