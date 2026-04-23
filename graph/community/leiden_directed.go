// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"math/rand/v2"
	"slices"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/internal/order"
)

// leidenDirected returns the hierarchical modularization of g at the given
// resolution using the Leiden algorithm.
func leidenDirected(g graph.Directed, resolution float64, src rand.Source) *ReducedDirected {
	c := reduceDirected(g, nil)
	rnd := rand.IntN
	if src != nil {
		rnd = rand.New(src).IntN
	}
	for iter := 0; iter < maxLeidenIterations; iter++ {
		l := newDirectedLocalMover(c, c.communities, resolution)
		if l == nil {
			return c
		}
		done := l.localMovingHeuristic(rnd)
		if done {
			return c
		}
		refined := refineDirected(l, resolution, rnd)
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
		c = reduceDirected(c, refined)
	}
	return c
}

// inducedDirected is a directed graph view over a subset of nodes of a
// ReducedDirected (the induced subgraph).
type inducedDirected struct {
	g   *ReducedDirected
	ids set.Ints[int64]
}

// Node returns the node with the given ID if it exists
func (s *inducedDirected) Node(id int64) graph.Node {
	if !s.ids.Has(id) {
		return nil
	}
	return s.g.Node(id)
}

// Nodes returns all nodes in the subgraph.
func (s *inducedDirected) Nodes() graph.Nodes {
	var nodes []graph.Node
	for id := range s.ids {
		nodes = append(nodes, s.g.Node(id))
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

// From returns all nodes in the subgraph that are adjacent from uid.
func (s *inducedDirected) From(uid int64) graph.Nodes {
	if !s.ids.Has(uid) {
		return graph.Empty
	}
	var nodes []graph.Node
	for _, vid := range s.g.edgesFrom[uid] {
		if s.ids.Has(int64(vid)) {
			nodes = append(nodes, s.g.nodes[vid])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

// To returns all nodes in the subgraph that are adjacent to vid.
func (s *inducedDirected) To(vid int64) graph.Nodes {
	if !s.ids.Has(vid) {
		return graph.Empty
	}
	var nodes []graph.Node
	for _, uid := range s.g.edgesTo[vid] {
		if s.ids.Has(int64(uid)) {
			nodes = append(nodes, s.g.nodes[uid])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

// HasEdgeBetween returns whether an edge exists between xid and
func (s *inducedDirected) HasEdgeBetween(xid, yid int64) bool {
	return s.ids.Has(xid) && s.ids.Has(yid) && s.g.HasEdgeBetween(xid, yid)
}

// HasEdgeFromTo returns whether an edge exists from uid to vid.
func (s *inducedDirected) HasEdgeFromTo(uid, vid int64) bool {
	return s.ids.Has(uid) && s.ids.Has(vid) && s.g.HasEdgeFromTo(uid, vid)
}

// Edge returns the edge from uid to vid if it exists.
func (s *inducedDirected) Edge(uid, vid int64) graph.Edge {
	return s.WeightedEdge(uid, vid)
}

func (s *inducedDirected) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	return s.WeightedEdgeBetween(uid, vid)
}

// WeightedEdgeBetween returns the weighted edge from xid to yid if it exists.
// For directed graphs this is the edge in the from→to direction only.
func (s *inducedDirected) WeightedEdgeBetween(xid, yid int64) graph.WeightedEdge {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return nil
	}
	return s.g.WeightedEdge(xid, yid)
}

func (s *inducedDirected) Weight(xid, yid int64) (w float64, ok bool) {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return 0, false
	}
	return s.g.Weight(xid, yid)
}

// refineDirected refines the partition in l by running local moving on the
// subgraph induced by each community (Leiden refinement phase).
func refineDirected(l *directedLocalMover, resolution float64, rnd func(int) int) [][]graph.Node {
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
		sub := &inducedDirected{g: l.g, ids: ids}
		subReduced := reduceDirected(sub, nil)
		subMover := newDirectedLocalMover(subReduced, subReduced.communities, resolution)
		if subMover == nil {
			for _, n := range comm {
				refined = append(refined, []graph.Node{n})
			}
			continue
		}
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
