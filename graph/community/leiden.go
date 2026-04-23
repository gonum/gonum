// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/internal/order"
)

// Leiden returns the hierarchical modularization of g at the given resolution
// using the Leiden algorithm.
//
// The Leiden algorithm improves upon Louvain by guaranteeing well-connected
// communities through a refinement phase after each move phase. See
// Traag, Waltman & Van Eck, Sci Rep 9, 5233 (2019),https://doi.org/10.1038/s41598-019-41695-z
//
// If src is nil, rand.IntN is used as the random generator. Leiden will panic
// if g has any edge with negative edge weight.
//
// graph.Undirect may be used as a shim to allow modularization of directed graphs.
func Leiden(g graph.Graph, resolution float64, src rand.Source) ReducedGraph {
	switch g := g.(type) {
	case graph.Undirected:
		return leidenUndirected(g, resolution, src)
	case graph.Directed:
		return leidenDirected(g, resolution, src)
	default:
		panic(fmt.Sprintf("community: invalid graph type: %T", g))
	}
}

// maxLeidenIterations limits the main Leiden loop to avoid infinite iteration
// on pathological graphs (e.g. with very poor modularity structure).
const maxLeidenIterations = 1000

// leidenUndirected returns the hierarchical modularization of g at the given
// resolution using the Leiden algorithm.
func leidenUndirected(g graph.Undirected, resolution float64, src rand.Source) *ReducedUndirected {
	c := reduceUndirected(g, nil)
	rnd := rand.IntN
	if src != nil {
		rnd = rand.New(src).IntN
	}
	for iter := 0; iter < maxLeidenIterations; iter++ {
		l := newUndirectedLocalMover(c, c.communities, resolution)
		if l == nil {
			return c
		}
		done := l.localMovingHeuristic(rnd)
		if done {
			return c
		}
		refined := refineUndirected(l, resolution, rnd)
		// If the refinement phase resulted in no reduction in the number of
		// communities (i.e. all communities are singletons), we are done.
		nonEmpty := 0
		for _, comm := range refined {
			if len(comm) > 0 {
				nonEmpty++
			}
		}
		if nonEmpty == len(c.nodes) {
			return c
		}
		c = reduceUndirected(c, refined)
	}
	return c
}

// inducedUndirected is an undirected graph view over a subset of nodes of a
// ReducedUndirected (the induced subgraph).
type inducedUndirected struct {
	g   *ReducedUndirected
	ids set.Ints[int64]
}

// Node returns the node with the given ID if it exists in the subgraph.
func (s *inducedUndirected) Node(id int64) graph.Node {
	if !s.ids.Has(id) {
		return nil
	}
	return s.g.Node(id)
}

// Nodes returns all nodes in the subgraph.
func (s *inducedUndirected) Nodes() graph.Nodes {
	var nodes []graph.Node
	for id := range s.ids {
		nodes = append(nodes, s.g.Node(id))
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

// From returns all nodes in the subgraph that are adjacent to uid.
func (s *inducedUndirected) From(uid int64) graph.Nodes {
	if !s.ids.Has(uid) {
		return graph.Empty
	}
	var nodes []graph.Node
	for _, vid := range s.g.edges[uid] {
		if s.ids.Has(int64(vid)) {
			nodes = append(nodes, s.g.nodes[vid])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

// HasEdgeBetween returns whether an edge exists between xid and yid in the subgraph.
func (s *inducedUndirected) HasEdgeBetween(xid, yid int64) bool {
	return s.ids.Has(xid) && s.ids.Has(yid) && s.g.HasEdgeBetween(xid, yid)
}

// Edge returns the edge between xid and yid if it exists in the subgraph.
func (s *inducedUndirected) Edge(uid, vid int64) graph.Edge {
	return s.WeightedEdgeBetween(uid, vid)
}

// EdgeBetween returns the edge between xid and yid if it exists in the subgraph.
func (s *inducedUndirected) EdgeBetween(xid, yid int64) graph.Edge {
	return s.WeightedEdgeBetween(xid, yid)
}

// WeightedEdge returns the weighted edge from uid to vid if it exists.
func (s *inducedUndirected) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	return s.WeightedEdgeBetween(uid, vid)
}

// WeightedEdgeBetween returns the weighted edge between xid and yid if it exists in the subgraph.
func (s *inducedUndirected) WeightedEdgeBetween(xid, yid int64) graph.WeightedEdge {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return nil
	}
	return s.g.WeightedEdgeBetween(xid, yid)
}

// Weight returns the weight of the edge between xid and yid.
func (s *inducedUndirected) Weight(xid, yid int64) (w float64, ok bool) {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return 0, false
	}
	return s.g.Weight(xid, yid)
}

// refineUndirected refines the partition in l by running local moving on the
// subgraph induced by each community. This yields a partition with well-connected
// communities and is the Leiden refinement phase.
func refineUndirected(l *undirectedLocalMover, resolution float64, rnd func(int) int) [][]graph.Node {
	var refined [][]graph.Node
	for _, comm := range l.communities {
		if len(comm) <= 1 {
			refined = append(refined, comm)
			continue
		}
		ids := make(set.Ints[int64])
		for _, n := range comm {
			ids.Add(n.ID())
		}
		sub := &inducedUndirected{g: l.g, ids: ids}
		subReduced := reduceUndirected(sub, nil)
		subMover := newUndirectedLocalMover(subReduced, subReduced.communities, resolution)
		if subMover == nil {
			for _, n := range comm {
				refined = append(refined, []graph.Node{n})
			}
			continue
		}
		// reduceUndirected sorts nodes by ID, so subReduced node i = i-th node in sorted comm.
		sortedComm := make([]graph.Node, len(comm))
		copy(sortedComm, comm)
		order.ByID(sortedComm)
		_ = subMover.localMovingHeuristic(rnd)
		for _, subComm := range subMover.communities {
			refinedComm := make([]graph.Node, len(subComm))
			for j, n := range subComm {
				refinedComm[j] = sortedComm[n.ID()]
			}
			refined = append(refined, refinedComm)
		}
	}
	return refined
}

// LeidenMultiplex returns the hierarchical modularization of g at the given
// resolution using the Leiden algorithm. If all is true and g has negatively
// weighted layers, all communities will be searched. If src is nil, rand.IntN
// is used. LeidenMultiplex will panic if g has any edge with weight that does
// not sign-match the layer weight.
//
// Directed multiplex is not yet implemented; use UndirectedMultiplex or
// graph.Undirect for directed layers.
func LeidenMultiplex(g Multiplex, weights, resolutions []float64, all bool, src rand.Source) ReducedMultiplex {
	if weights != nil && len(weights) != g.Depth() {
		panic("community: weights vector length mismatch")
	}
	if resolutions != nil && len(resolutions) != 1 && len(resolutions) != g.Depth() {
		panic("community: resolutions vector length mismatch")
	}
	switch g := g.(type) {
	case UndirectedMultiplex:
		return leidenUndirectedMultiplex(g, weights, resolutions, all, src)
	case DirectedMultiplex:
		return leidenDirectedMultiplex(g, weights, resolutions, all, src)
	default:
		panic(fmt.Sprintf("community: invalid graph type: %T", g))
	}
}
