// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import (
	"maps"
	"slices"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/linear"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

/*
Intervals are defined as the maximal, single entry subgraph for which h (head)
is the entry node and in which all closed paths contain h.

This is the algorithm given for computing them:

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

// Intervals returns the IntervalGraph containing the individual intervals
// of the directed graph in g when starting from the entry node identified
// by eid.
func Intervals(g graph.Directed, eid int64) IntervalGraph {
	var intervals []*Interval

	var ns linear.NodeStack
	visited := make(map[int64]bool)

	dfsPostorder(g, eid, &ns, visited)
	slices.Reverse(ns)
	var worklist linear.NodeQueue
	worklist.Enqueue(g.Node(eid))
	inInterval := make(map[int64]graph.Node)
	node2interval := make(map[int64]*Interval)
	id := int64(0)

	for worklist.Len() != 0 {
		var interval Interval
		map1, map2 := interval.findInterval(worklist.Dequeue(), g)
		maps.Copy(inInterval, map1)
		maps.Copy(node2interval, map2)

		interval.id = id
		id++
		intervals = append(intervals, &interval)

		for _, node := range ns {
			if inInterval[node.ID()] != nil {
				continue
			}

			preds := g.To(node.ID())
			predsLength := preds.Len()
			x := 0
			for preds.Next() {
				if interval.nodes[preds.Node().ID()] != nil {
					x++
				}
			}

			if 0 < x && x < predsLength {
				worklist.Enqueue(node)
				break
			}
		}
	}

	return linkIntervals(intervals, g, node2interval)
}

// IntervalGraph contains the intervals and the edges between the intervals.
type IntervalGraph struct {
	Intervals map[int64]*Interval
	nodes     map[int64]graph.Node
	head      graph.Node
	from      map[int64]map[int64]graph.Edge
	to        map[int64]map[int64]graph.Edge
}

// Head returns header node for an IntervalGraph.
func (ig *IntervalGraph) Head() graph.Node {
	return ig.head
}

// Nodes returns a node iterator for an IntervalGraph.
func (ig *IntervalGraph) Nodes() graph.Nodes {
	if len(ig.nodes) == 0 {
		return graph.Empty
	}
	return iterator.NewNodes(ig.nodes)
}

// Edge returns the edge given 2 node id's if the edge exists.
// Else it returns null.
func (ig *IntervalGraph) Edge(uid, vid int64) graph.Edge {
	edge, ok := ig.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}

// From returns all nodes in g that can be reached directly from n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (ig *IntervalGraph) From(id int64) graph.Nodes {
	if len(ig.from[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByEdge(ig.nodes, ig.from[id])
}

// To returns all nodes in g that can reach directly to n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (ig *IntervalGraph) To(id int64) graph.Nodes {
	if len(ig.to[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByEdge(ig.nodes, ig.to[id])
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (ig *IntervalGraph) HasEdgeBetween(xid int64, yid int64) bool {
	if _, ok := ig.from[xid][yid]; ok {
		return true
	}
	_, ok := ig.from[yid][xid]
	return ok
}

// Node returns the node with the given ID if it exists in the IntervalGraph,
// and nil otherwise.
func (ig *IntervalGraph) Node(id int64) graph.Node {
	return ig.nodes[id]
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (ig *IntervalGraph) HasEdgeFromTo(uid, vid int64) bool {
	if _, ok := ig.from[uid][vid]; !ok {
		return false
	}
	return true
}

// setEdge adds e, an edge from one node to another.
func (ig *IntervalGraph) setEdge(e graph.Edge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	ig.nodes[fid] = from
	ig.nodes[tid] = to
	ig.from[fid] = map[int64]graph.Edge{tid: e}
	ig.to[tid] = map[int64]graph.Edge{fid: e}
}

// Interval I(h) is the maximal, single entry subgraph for which h (head)
// is the entry node and in which all closed paths contain h.
type Interval struct {
	head  graph.Node
	id    int64
	nodes map[int64]graph.Node
	from  map[int64]map[int64]graph.Edge
	to    map[int64]map[int64]graph.Edge
}

// ID returns the ID number of the node.
func (i *Interval) ID() int64 {
	return i.id
}

// Head returns header node for an interval.
func (i *Interval) Head() graph.Node {
	return i.head
}

// Nodes returns a node iterator for an interval.
func (i *Interval) Nodes() graph.Nodes {
	if len(i.nodes) == 0 {
		return graph.Empty
	}
	return iterator.NewNodes(i.nodes)
}

// Edge returns the edge given 2 node id's if the edge exists.
// Else it returns null.
func (i *Interval) Edge(uid, vid int64) graph.Edge {
	edge, ok := i.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}

// From returns all nodes in g that can be reached directly from n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (i *Interval) From(id int64) graph.Nodes {
	if len(i.from[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByEdge(i.nodes, i.from[id])
}

// To returns all nodes in g that can reach directly to n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (i *Interval) To(id int64) graph.Nodes {
	if len(i.to[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByEdge(i.nodes, i.to[id])
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (i *Interval) HasEdgeBetween(xid int64, yid int64) bool {
	if _, ok := i.from[xid][yid]; ok {
		return true
	}
	_, ok := i.from[yid][xid]
	return ok
}

// Node returns the node with the given ID if it exists in the interval,
// and nil otherwise.
func (i *Interval) Node(id int64) graph.Node {
	return i.nodes[id]
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (i *Interval) HasEdgeFromTo(uid, vid int64) bool {
	if _, ok := i.from[uid][vid]; !ok {
		return false
	}
	return true
}

// findInterval finds all interval nodes.
// Nodes are added to the interval if all their predecessors are in
// the interval or they are the header node.
func (i *Interval) findInterval(h graph.Node, g graph.Directed) (map[int64]graph.Node, map[int64]*Interval) {
	i.head = h
	var nq linear.NodeQueue
	nq.Enqueue(h)
	node2interval := make(map[int64]*Interval)
	node2interval[h.ID()] = i
	i.nodes = make(map[int64]graph.Node)
	i.nodes[h.ID()] = h
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
				node2interval[succs.Node().ID()] = i
			}
		}
	}

	return i.nodes, node2interval
}

// dfsPosterorder puts nodes into a stack in postorder.
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

	ns.Push(g.Node(eid))
}

// linkIntervals computes the internal and external edges for the intervals.
func linkIntervals(intervals []*Interval, g graph.Directed, node2interval map[int64]*Interval) IntervalGraph {
	ig := IntervalGraph{
		Intervals: make(map[int64]*Interval),
		nodes:     make(map[int64]graph.Node),
		head:      intervals[0].head,
		from:      make(map[int64]map[int64]graph.Edge),
		to:        make(map[int64]map[int64]graph.Edge),
	}

	for i, interval := range intervals {
		ig.nodes[int64(i)] = interval
		ig.Intervals[int64(i)] = interval

		interval.from = make(map[int64]map[int64]graph.Edge)
		interval.to = make(map[int64]map[int64]graph.Edge)
		for _, node := range interval.nodes {
			succs := g.From(node.ID())
			for succs.Next() {
				succNode := succs.Node()
				if interval.nodes[succNode.ID()] != nil {
					// edge contained within interval
					if interval.from[node.ID()] == nil {
						interval.from[node.ID()] = make(map[int64]graph.Edge)
					}

					if interval.to[succNode.ID()] == nil {
						interval.to[succNode.ID()] = make(map[int64]graph.Edge)
					}

					interval.from[node.ID()][succNode.ID()] = g.Edge(node.ID(), succNode.ID())
					interval.to[succNode.ID()][node.ID()] = g.Edge(node.ID(), succNode.ID())
				} else {
					// edge from one interval to another interval
					succInterval := node2interval[succNode.ID()]
					if ig.to[succInterval.ID()] == nil {
						ig.to[succInterval.ID()] = make(map[int64]graph.Edge)
					}

					if ig.from[interval.ID()] == nil {
						ig.from[interval.ID()] = make(map[int64]graph.Edge)
					}

					ig.setEdge(simple.Edge{F: interval, T: succInterval})
				}
			}
		}
	}

	return ig
}
