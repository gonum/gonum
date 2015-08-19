// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"sort"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/internal"
)

// EdgeListerGraph is an undirected graph than returns its complete set of edges.
type EdgeListerGraph interface {
	graph.Undirected
	Edges() []graph.Edge
}

// Prim generates a minimum spanning tree of g by greedy tree extension, placing
// the result in the destination. The destination is not cleared first.
func Prim(dst graph.MutableUndirected, g EdgeListerGraph) {
	var weight Weighting
	if wg, ok := g.(graph.Weighter); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	nlist := g.Nodes()

	if nlist == nil || len(nlist) == 0 {
		return
	}

	dst.AddNode(nlist[0])
	remain := make(internal.IntSet)
	for _, node := range nlist[1:] {
		remain.Add(node.ID())
	}

	edgeList := g.Edges()
	for remain.Count() != 0 {
		var edges []concrete.Edge
		for _, e := range edgeList {
			u := e.From()
			v := e.To()
			if (dst.Has(u) && remain.Has(v.ID())) || (dst.Has(v) && remain.Has(u.ID())) {
				w, ok := weight(u, v)
				if !ok {
					panic("prim: unexpected invalid weight")
				}
				edges = append(edges, concrete.Edge{F: u, T: v, W: w})
			}
		}

		sort.Sort(byWeight(edges))
		min := edges[0]
		dst.SetEdge(min)
		remain.Remove(min.From().ID())
	}

}

// Kruskal generates a minimum spanning tree of g by greedy tree coalesence, placing
// the result in the destination. The destination is not cleared first.
func Kruskal(dst graph.MutableUndirected, g EdgeListerGraph) {
	var weight Weighting
	if wg, ok := g.(graph.Weighter); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	edgeList := g.Edges()
	edges := make([]concrete.Edge, 0, len(edgeList))
	for _, e := range edgeList {
		u := e.From()
		v := e.To()
		w, ok := weight(u, v)
		if !ok {
			panic("kruskal: unexpected invalid weight")
		}
		edges = append(edges, concrete.Edge{F: u, T: v, W: w})
	}

	sort.Sort(byWeight(edges))

	ds := newDisjointSet()
	for _, node := range g.Nodes() {
		ds.makeSet(node.ID())
	}

	for _, e := range edges {
		// The disjoint set doesn't really care for which is head and which is tail so this
		// should work fine without checking both ways
		if s1, s2 := ds.find(e.From().ID()), ds.find(e.To().ID()); s1 != s2 {
			ds.union(s1, s2)
			dst.SetEdge(e)
		}
	}
}

type byWeight []concrete.Edge

func (e byWeight) Len() int           { return len(e) }
func (e byWeight) Less(i, j int) bool { return e[i].Weight() < e[j].Weight() }
func (e byWeight) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
