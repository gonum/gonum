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

// refineDirectedMultiplex refines the partition by running local moving on
// the sub-multiplex induced by each community (Leiden refinement phase).
// It uses the same weights and resolutions as the main multiplex mover.
func refineDirectedMultiplex(l *directedMultiplexLocalMover, weights, resolutions []float64, all bool, rnd func(int) int) [][]graph.Node {
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
		sub := &inducedDirectedMultiplex{g: l.g, ids: ids}
		subReduced := reduceDirectedMultiplex(sub, nil, weights)
		subMover := newDirectedMultiplexLocalMover(subReduced, subReduced.communities, weights, resolutions, all)
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

func leidenDirectedMultiplex(g DirectedMultiplex, weights, resolutions []float64, all bool, src rand.Source) *ReducedDirectedMultiplex {
	if weights != nil && len(weights) != g.Depth() {
		panic("community: weights vector length mismatch")
	}
	if resolutions != nil && len(resolutions) != 1 && len(resolutions) != g.Depth() {
		panic("community: resolutions vector length mismatch")
	}

	c := reduceDirectedMultiplex(g, nil, weights)
	rnd := rand.IntN
	if src != nil {
		rnd = rand.New(src).IntN
	}
	for iter := 0; iter < maxLeidenIterations; iter++ {
		l := newDirectedMultiplexLocalMover(c, c.communities, weights, resolutions, all)
		if l == nil {
			return c
		}
		done := l.localMovingHeuristic(rnd)
		if done {
			return c
		}
		refined := refineDirectedMultiplex(l, weights, resolutions, all, rnd)
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
		c = reduceDirectedMultiplex(c, refined, weights)
	}
	return c
}

// inducedDirectedMultiplex is an directed multiplex view over a subset of
// nodes of a ReducedDirectedMultiplex.
type inducedDirectedMultiplex struct {
	g   *ReducedDirectedMultiplex
	ids set.Ints[int64]
}

func (s *inducedDirectedMultiplex) Node(id int64) graph.Node {
	if !s.ids.Has(id) || id < 0 || id >= int64(len(s.g.nodes)) {
		return nil
	}
	return s.g.nodes[id]
}

func (s *inducedDirectedMultiplex) Nodes() graph.Nodes {
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

func (s *inducedDirectedMultiplex) Depth() int {
	return s.g.Depth()
}

func (s *inducedDirectedMultiplex) Layer(l int) graph.Directed {
	return &inducedMultiplexLayerDirected{multiplex: s.g, layer: l, ids: s.ids}
}

// inducedMultiplexLayerDirected is a single layer of an induced multiplex (graph.Directed).
type inducedMultiplexLayerDirected struct {
	multiplex *ReducedDirectedMultiplex
	layer     int
	ids       set.Ints[int64]
}

func (s *inducedMultiplexLayerDirected) Node(id int64) graph.Node {
	if !s.ids.Has(id) || id < 0 || id >= int64(len(s.multiplex.nodes)) {
		return nil
	}
	return s.multiplex.nodes[id]
}

func (s *inducedMultiplexLayerDirected) Nodes() graph.Nodes {
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

func (s *inducedMultiplexLayerDirected) From(uid int64) graph.Nodes {
	if !s.ids.Has(uid) {
		return graph.Empty
	}
	var nodes []graph.Node
	for _, vid := range s.multiplex.layers[s.layer].edgesFrom[uid] {
		if s.ids.Has(int64(vid)) {
			nodes = append(nodes, s.multiplex.nodes[vid])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

func (s *inducedMultiplexLayerDirected) To(vid int64) graph.Nodes {
	if !s.ids.Has(vid) {
		return graph.Empty
	}
	var nodes []graph.Node
	for _, uid := range s.multiplex.layers[s.layer].edgesTo[vid] {
		if s.ids.Has(int64(uid)) {
			nodes = append(nodes, s.multiplex.nodes[uid])
		}
	}
	slices.SortFunc(nodes, func(a, b graph.Node) int {
		return int(a.ID() - b.ID())
	})
	return iterator.NewOrderedNodes(nodes)
}

func (s *inducedMultiplexLayerDirected) HasEdgeBetween(xid, yid int64) bool {
	return s.ids.Has(xid) && s.ids.Has(yid) && s.multiplex.Layer(s.layer).(directedLayerHandle).HasEdgeBetween(xid, yid)
}

func (s *inducedMultiplexLayerDirected) HasEdgeFromTo(uid, vid int64) bool {
	return s.ids.Has(uid) && s.ids.Has(vid) && s.multiplex.Layer(s.layer).(directedLayerHandle).HasEdgeFromTo(uid, vid)
}

func (s *inducedMultiplexLayerDirected) Edge(uid, vid int64) graph.Edge {
	return s.WeightedEdge(uid, vid)
}

func (s *inducedMultiplexLayerDirected) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	if !s.ids.Has(uid) || !s.ids.Has(vid) {
		return nil
	}
	return s.multiplex.Layer(s.layer).(directedLayerHandle).WeightedEdge(uid, vid)
}

func (s *inducedMultiplexLayerDirected) Weight(xid, yid int64) (w float64, ok bool) {
	if !s.ids.Has(xid) || !s.ids.Has(yid) {
		return 0, false
	}
	return s.multiplex.Layer(s.layer).(directedLayerHandle).Weight(xid, yid)
}
