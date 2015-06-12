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

// Prim generates a minimum spanning tree of g by greedy tree extension, placing
// the result in the destination. The destination is not cleared first.
func Prim(dst graph.MutableGraph, g graph.EdgeListGraph) {
	var weight graph.CostFunc
	if g, ok := g.(graph.Coster); ok {
		weight = g.Cost
	} else {
		weight = graph.UniformCost
	}

	nlist := g.Nodes()

	if nlist == nil || len(nlist) == 0 {
		return
	}

	dst.AddNode(nlist[0])
	remainingNodes := make(internal.IntSet)
	for _, node := range nlist[1:] {
		remainingNodes.Add(node.ID())
	}

	edgeList := g.Edges()
	for remainingNodes.Count() != 0 {
		var edges []concrete.WeightedEdge
		for _, edge := range edgeList {
			if (dst.Has(edge.From()) && remainingNodes.Has(edge.To().ID())) ||
				(dst.Has(edge.To()) && remainingNodes.Has(edge.From().ID())) {

				edges = append(edges, concrete.WeightedEdge{Edge: edge, Cost: weight(edge)})
			}
		}

		sort.Sort(byWeight(edges))
		myEdge := edges[0]

		dst.AddUndirectedEdge(myEdge.Edge, myEdge.Cost)
		remainingNodes.Remove(myEdge.Edge.From().ID())
	}

}

// Kruskal generates a minimum spanning tree of g by greedy tree coalesence, placing
// the result in the destination. The destination is not cleared first.
func Kruskal(dst graph.MutableGraph, g graph.EdgeListGraph) {
	var weight graph.CostFunc
	if g, ok := g.(graph.Coster); ok {
		weight = g.Cost
	} else {
		weight = graph.UniformCost
	}

	edgeList := g.Edges()
	edges := make([]concrete.WeightedEdge, 0, len(edgeList))
	for _, edge := range edgeList {
		edges = append(edges, concrete.WeightedEdge{Edge: edge, Cost: weight(edge)})
	}

	sort.Sort(byWeight(edges))

	ds := newDisjointSet()
	for _, node := range g.Nodes() {
		ds.makeSet(node.ID())
	}

	for _, edge := range edges {
		// The disjoint set doesn't really care for which is head and which is tail so this
		// should work fine without checking both ways
		if s1, s2 := ds.find(edge.Edge.From().ID()), ds.find(edge.Edge.To().ID()); s1 != s2 {
			ds.union(s1, s2)
			dst.AddUndirectedEdge(edge.Edge, edge.Cost)
		}
	}
}

type byWeight []concrete.WeightedEdge

func (e byWeight) Len() int {
	return len(e)
}

func (e byWeight) Less(i, j int) bool {
	return e[i].Cost < e[j].Cost
}

func (e byWeight) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
