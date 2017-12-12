// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

// Undirect converts a directed graph to an undirected graph.
type Undirect struct {
	G Directed
}

var _ Undirected = Undirect{}

// Has returns whether the node exists within the graph.
func (g Undirect) Has(n Node) bool { return g.G.Has(n) }

// Nodes returns all the nodes in the graph.
func (g Undirect) Nodes() []Node { return g.G.Nodes() }

// From returns all nodes in g that can be reached directly from u.
func (g Undirect) From(u Node) []Node {
	var nodes []Node
	seen := make(map[int64]struct{})
	for _, n := range g.G.From(u) {
		seen[n.ID()] = struct{}{}
		nodes = append(nodes, n)
	}
	for _, n := range g.G.To(u) {
		id := n.ID()
		if _, ok := seen[id]; ok {
			continue
		}
		seen[n.ID()] = struct{}{}
		nodes = append(nodes, n)
	}
	return nodes
}

// HasEdgeBetween returns whether an edge exists between nodes x and y.
func (g Undirect) HasEdgeBetween(x, y Node) bool { return g.G.HasEdgeBetween(x, y) }

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
// If an edge exists, the Edge returned is an EdgePair. The weight of
// the edge is determined by applying the Merge func to the weights of the
// edges between u and v.
func (g Undirect) Edge(u, v Node) Edge { return g.EdgeBetween(u, v) }

// EdgeBetween returns the edge between nodes x and y. If an edge exists, the
// Edge returned is an EdgePair. The weight of the edge is determined by
// applying the Merge func to the weights of edges between x and y.
func (g Undirect) EdgeBetween(x, y Node) Edge {
	fe := g.G.Edge(x, y)
	re := g.G.Edge(y, x)
	if fe == nil && re == nil {
		return nil
	}

	return EdgePair{fe, re}
}

// UndirectWeighted converts a directed weighted graph to an undirected weighted graph,
// resolving edge weight conflicts.
type UndirectWeighted struct {
	G WeightedDirected

	// Absent is the value used to
	// represent absent edge weights
	// passed to Merge if the reverse
	// edge is present.
	Absent float64

	// Merge defines how discordant edge
	// weights in G are resolved. A merge
	// is performed if at least one edge
	// exists between the nodes being
	// considered. The edges corresponding
	// to the two weights are also passed,
	// in the same order.
	// The order of weight parameters
	// passed to Merge is not defined, so
	// the function should be commutative.
	// If Merge is nil, the arithmetic
	// mean is used to merge weights.
	Merge func(x, y float64, xe, ye Edge) float64
}

var (
	_ Undirected         = UndirectWeighted{}
	_ WeightedUndirected = UndirectWeighted{}
)

// Has returns whether the node exists within the graph.
func (g UndirectWeighted) Has(n Node) bool { return g.G.Has(n) }

// Nodes returns all the nodes in the graph.
func (g UndirectWeighted) Nodes() []Node { return g.G.Nodes() }

// From returns all nodes in g that can be reached directly from u.
func (g UndirectWeighted) From(u Node) []Node {
	var nodes []Node
	seen := make(map[int64]struct{})
	for _, n := range g.G.From(u) {
		seen[n.ID()] = struct{}{}
		nodes = append(nodes, n)
	}
	for _, n := range g.G.To(u) {
		id := n.ID()
		if _, ok := seen[id]; ok {
			continue
		}
		seen[n.ID()] = struct{}{}
		nodes = append(nodes, n)
	}
	return nodes
}

// HasEdgeBetween returns whether an edge exists between nodes x and y.
func (g UndirectWeighted) HasEdgeBetween(x, y Node) bool { return g.G.HasEdgeBetween(x, y) }

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
// If an edge exists, the Edge returned is an EdgePair. The weight of
// the edge is determined by applying the Merge func to the weights of the
// edges between u and v.
func (g UndirectWeighted) Edge(u, v Node) Edge { return g.WeightedEdgeBetween(u, v) }

// WeightedEdge returns the weighted edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
// If an edge exists, the Edge returned is an EdgePair. The weight of
// the edge is determined by applying the Merge func to the weights of the
// edges between u and v.
func (g UndirectWeighted) WeightedEdge(u, v Node) WeightedEdge { return g.WeightedEdgeBetween(u, v) }

// EdgeBetween returns the edge between nodes x and y. If an edge exists, the
// Edge returned is an EdgePair. The weight of the edge is determined by
// applying the Merge func to the weights of edges between x and y.
func (g UndirectWeighted) EdgeBetween(x, y Node) Edge {
	return g.WeightedEdgeBetween(x, y)
}

// WeightedEdgeBetween returns the weighted edge between nodes x and y. If an edge exists, the
// Edge returned is an EdgePair. The weight of the edge is determined by
// applying the Merge func to the weights of edges between x and y.
func (g UndirectWeighted) WeightedEdgeBetween(x, y Node) WeightedEdge {
	fe := g.G.Edge(x, y)
	re := g.G.Edge(y, x)
	if fe == nil && re == nil {
		return nil
	}

	f, ok := g.G.Weight(x, y)
	if !ok {
		f = g.Absent
	}
	r, ok := g.G.Weight(y, x)
	if !ok {
		r = g.Absent
	}

	var w float64
	if g.Merge == nil {
		w = (f + r) / 2
	} else {
		w = g.Merge(f, r, fe, re)
	}
	return WeightedEdgePair{EdgePair: [2]Edge{fe, re}, W: w}
}

// Weight returns the weight for the edge between x and y if Edge(x, y) returns a non-nil Edge.
// If x and y are the same node the internal node weight is returned. If there is no joining
// edge between the two nodes the weight value returned is zero. Weight returns true if an edge
// exists between x and y or if x and y have the same ID, false otherwise.
func (g UndirectWeighted) Weight(x, y Node) (w float64, ok bool) {
	fe := g.G.Edge(x, y)
	re := g.G.Edge(y, x)

	f, fOk := g.G.Weight(x, y)
	if !fOk {
		f = g.Absent
	}
	r, rOK := g.G.Weight(y, x)
	if !rOK {
		r = g.Absent
	}
	ok = fOk || rOK

	if g.Merge == nil {
		return (f + r) / 2, ok
	}
	return g.Merge(f, r, fe, re), ok
}

// EdgePair is an opposed pair of directed edges.
type EdgePair [2]Edge

// From returns the from node of the first non-nil edge, or nil.
func (e EdgePair) From() Node {
	if e[0] != nil {
		return e[0].From()
	} else if e[1] != nil {
		return e[1].From()
	}
	return nil
}

// To returns the to node of the first non-nil edge, or nil.
func (e EdgePair) To() Node {
	if e[0] != nil {
		return e[0].To()
	} else if e[1] != nil {
		return e[1].To()
	}
	return nil
}

// WeightedEdgePair is an opposed pair of directed edges.
type WeightedEdgePair struct {
	EdgePair
	W float64
}

// Weight returns the merged edge weights of the two edges.
func (e WeightedEdgePair) Weight() float64 { return e.W }
