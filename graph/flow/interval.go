// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Find the unique set of intervals of a control flow graph, as described in
// Allen, Frances E., and John Cocke. "A program data flow analysis procedure."
// Communications of the ACM 19.3 (1976): 137 [1].
//
// [1]: https://amturing.acm.org/p137-allen.pdf

package flow

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

// --- [ Algorithm for Finding Intervals ] -------------------------------------

// Intervals returns the unique set of intervals of the given control flow
// graph. The returned set of intervals share the underlying nodes and edges of
// the provided control flow graph.
func Intervals(g Graph) []*Interval {
	var is []*Interval
	// From Algorithm for Finding Intervals [1].
	//
	// 1. Establish a set H for header nodes and initialize it with n_0, the
	//    unique entry node for the graph.
	H := newQueue()
	H.push(g.Entry())
	// 2. For h \in H find I(h) as follows:
	for !H.empty() {
		head := H.pop()
		// 2.1. Put h in I(h) as the first element of I(h).
		i := newInterval(g, head)
		for {
			// 2.2. Add to I(h) any node all of whose immediate predecessors are
			//      already in I(h).
			n, ok := findNodeWithImmPredsInInterval(g, i)
			if !ok {
				// 2.3. Repeat 2.2 until no more nodes can be added to I(h).
				break
			}
			i.addNode(n)
		}
		// 3. Add to H all nodes in G which are not already in H and which are not
		//    in I(h) but which have immediate predecessors in I(h). Therefore a
		//    node is added to H the first time any (but not all) of its immediate
		//    predecessors become members of an interval.
		for {
			n, ok := findUnusedNodeWithImmPredInInterval(g, I, H)
			if !ok {
				break
			}
			H.push(n)
		}
		// 4. Add I(h) to a set 𝓘 of intervals being developed.
		is = append(is, i)
		// 5. Select the next unprocessed node in H and repeat steps 2, 3, 4, 5.
		//    When there are no more unprocessed nodes in H, the procedure
		//    terminates.
	}
	return is
}

// --- [/ Algorithm for Finding Intervals ] ------------------------------------

// An Interval I(h) with header node h is a maximal single-entry subgraph in
// which h is the only entry node and all cycles contain h.
type Interval struct {
	g     Graph
	head  graph.Node
	nodes map[int64]graph.Node
}

// newInterval returns the scaffolding of the interval I(h) based on the given
// control flow graph and interval header node. The remaining nodes of the
// interval will be added by the Algorithm for Finding Intervals.
func newInterval(g Graph, head graph.Node) *Interval {
	return &Interval{
		g:    g,
		head: head,
		nodes: map[int64]graph.Node{
			head.ID(): head,
		},
	}
}

// addNode adds a node to the interval. addNode panics if the added node ID
// matches an existing node ID.
func (i *Interval) addNode(n graph.Node) {
	if prev, ok := i.nodes[n.ID()]; ok {
		panic(fmt.Errorf("node with ID %d already present in interval; prev DOTID %q, new DOTID %q", n.ID(), dotid(prev), dotid(n)))
	}
	i.nodes[n.ID()] = n
}

// Entry returns the header node of the interval.
func (i *Interval) Entry() graph.Node {
	return i.head
}

// Node returns the node with the given ID if it exists within the interval, and
// nil otherwise.
func (i *Interval) Node(id int64) graph.Node {
	panic("not yet implemented")
}

// Nodes returns all the nodes of the interval.
func (i *Interval) Nodes() graph.Nodes {
	panic("not yet implemented")
}

// From returns all nodes that can be reached directly from the node with the
// given ID within the interval.
func (i *Interval) From(id int64) graph.Nodes {
	panic("not yet implemented")
}

// HasEdgeBetween returns whether an edge exists between nodes with IDs xid and
// yid within the interval without considering direction.
func (i *Interval) HasEdgeBetween(xid, yid int64) bool {
	panic("not yet implemented")
}

// Edge returns the edge from u to v, with IDs uid and vid, if such an edge
// exists within the interval and nil otherwise. The node v must be directly
// reachable from u as defined by the From method.
func (i *Interval) Edge(uid, vid int64) graph.Edge {
	panic("not yet implemented")
}

// HasEdgeFromTo returns whether an edge exists within the interval from u to v
// with IDs uid and vid.
func (i *Interval) HasEdgeFromTo(uid, vid int64) bool {
	panic("not yet implemented")
}

// To returns all nodes that can reach directly to the node with the given ID
// within the interval.
func (i *Interval) To(id int64) graph.Nodes {
	panic("not yet implemented")
}

// dotid returns the DOT ID of the given if present.
func dotid(n graph.Node) string {
	if n, ok := n.(dot.Node); ok {
		return n.DOTID()
	}
	// DOT ID unknown.
	return ""
}
