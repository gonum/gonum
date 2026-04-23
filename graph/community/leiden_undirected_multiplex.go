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

// leidenUndirectedMultiplex returns the hierarchical modularization of g at the
// given resolution using the Leiden algorithm.
func leidenUndirectedMultiplex(g UndirectedMultiplex, weights, resolutions []float64, all bool, src rand.Source) *ReducedUndirectedMultiplex {
	if weights != nil && len(weights) != g.Depth() {
		panic("community: weights vector length mismatch")
	}
	if resolutions != nil && len(resolutions) != 1 && len(resolutions) != g.Depth() {
		panic("community: resolutions vector length mismatch")
	}

	c := reduceUndirectedMultiplex(g, nil, weights)
	rnd := rand.IntN
	if src != nil {
		rnd = rand.New(src).IntN
	}
	for iter := 0; iter < maxLeidenIterations; iter++ {
		l := newUndirectedMultiplexLocalMover(c, c.communities, weights, resolutions, all)
		if l == nil {
			return c
		}
		done := l.localMovingHeuristic(rnd)
		if done {
			return c
		}
		refined := refineUndirectedMultiplex(l, weights, resolutions, all, rnd)
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
		c = reduceUndirectedMultiplex(c, refined, weights)
	}
	return c
}

// inducedUndirectedMultiplex is an undirected multiplex view over a subset of
// nodes of a ReducedUndirectedMultiplex.
type inducedUndirectedMultiplex struct {
	g   *ReducedUndirectedMultiplex
	ids set.Ints[int64]
}

func (s *inducedUndirectedMultiplex) Node(id int64) graph.Node {
	if !s.ids.Has(id) || id < 0 || id >= int64(len(s.g.nodes)) {
		return nil
	}
	return s.g.nodes[id]
}

func (s *inducedUndirectedMultiplex) Nodes() graph.Nodes {
	var nodes []graph.Node
	for id := range s.ids {
		if id >= 0 && id < int64(len(s.g.nodes)) {
			nodes = append(nodes, s.g.nodes[id])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

func (s *inducedUndirectedMultiplex) Depth() int {
	return s.g.Depth()
}

func (s *inducedUndirectedMultiplex) Layer(l int) graph.Undirected {
	return &inducedMultiplexLayer{multiplex: s.g, layer: l, ids: s.ids}
}

// inducedMultiplexLayer is a single layer of an induced multiplex (graph.Undirected).
type inducedMultiplexLayer struct {
	multiplex *ReducedUndirectedMultiplex
	layer     int
	ids       set.Ints[int64]
}

func (s *inducedMultiplexLayer) Node(id int64) graph.Node {
	if !s.ids.Has(id) || id < 0 || id >= int64(len(s.multiplex.nodes)) {
		return nil
	}
	return s.multiplex.nodes[id]
}

func (s *inducedMultiplexLayer) Nodes() graph.Nodes {
	var nodes []graph.Node
	for id := range s.ids {
		if id >= 0 && id < int64(len(s.multiplex.nodes)) {
			nodes = append(nodes, s.multiplex.nodes[id])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

func (s *inducedMultiplexLayer) From(uid int64) graph.Nodes {
	if !s.ids.Has(uid) {
		return graph.Empty
	}
	var nodes []graph.Node
	for _, vid := range s.multiplex.layers[s.layer].edges[uid] {
		if s.ids.Has(int64(vid)) {
			nodes = append(nodes, s.multiplex.nodes[vid])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

func (s *inducedMultiplexLayer) HasEdgeBetween(xid, yid int64) bool {
	return s.ids.Has(xid) && s.ids.Has(yid) && s.multiplex.Layer(s.layer).(undirectedLayerHandle).HasEdgeBetween(xid, yid)
}

func (s *inducedMultiplexLayer) Edge(uid, vid int64) graph.Edge {
	return s.WeightedEdgeBetween(uid, vid)
}

func (s *inducedMultiplexLayer) EdgeBetween(xid, yid int64) graph.Edge {
	return s.WeightedEdgeBetween(xid, yid)
}

func (s *inducedMultiplexLayer) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	return s.WeightedEdgeBetween(uid, vid)
}

func (s *inducedMultiplexLayer) WeightedEdgeBetween(xid, yid int64) graph.WeightedEdge {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return nil
	}
	return s.multiplex.Layer(s.layer).(undirectedLayerHandle).WeightedEdgeBetween(xid, yid)
}

func (s *inducedMultiplexLayer) Weight(xid, yid int64) (w float64, ok bool) {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return 0, false
	}
	return s.multiplex.Layer(s.layer).(undirectedLayerHandle).Weight(xid, yid)
}

// refineUndirectedMultiplex refines the partition by running local moving on
// the sub-multiplex induced by each community (Leiden refinement phase).
// It uses the same weights and resolutions as the main multiplex mover.
func refineUndirectedMultiplex(l *undirectedMultiplexLocalMover, weights, resolutions []float64, all bool, rnd func(int) int) [][]graph.Node {
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
		sub := &inducedUndirectedMultiplex{g: l.g, ids: ids}
		subReduced := reduceUndirectedMultiplex(sub, nil, weights)
		subMover := newUndirectedMultiplexLocalMover(subReduced, subReduced.communities, weights, resolutions, all)
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
