package flow

import "gonum.org/v1/gonum/graph"

// returns the set of intervals given by the control flow graph
func Intervals(g graph.Directed, eid int64) []*Interval {

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

// returns the edge given 2 node id's if the edge exists.
// else it returns null
func (i *Interval) Edge(uid, vid int64) graph.Edge {

}
