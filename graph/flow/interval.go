// Find the unique set of intervals of a control flow graph, as described in
// Allen, Frances E., and John Cocke. "A program data flow analysis procedure."
// Communications of the ACM 19.3 (1976): 137 [1].
//
// [1]: https://amturing.acm.org/p137-allen.pdf

package flow

import "gonum.org/v1/gonum/graph"

// Intervals returns the unique set of intervals of the given control flow
// graph. The returned set of intervals share the underlying nodes and edges of
// the provided control flow graph.
func Intervals(g Graph) []*Interval {
	panic("not yet implemented")
}

// An Interval I(h) with header node h is a maximal single-entry subgraph in
// which h is the only entry node and all cycles contain h.
type Interval struct {
	g     Graph
	head  graph.Node
	nodes map[int64]graph.Node
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
//
// Nodes must not return nil.
func (i *Interval) Nodes() graph.Nodes {
	panic("not yet implemented")
}

// From returns all nodes that can be reached directly from the node with the
// given ID within the interval.
//
// From must not return nil.
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
//
// To must not return nil.
func (i *Interval) To(id int64) graph.Nodes {
	panic("not yet implemented")
}
