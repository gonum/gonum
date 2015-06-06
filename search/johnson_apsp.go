// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"math"
	"math/rand"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

// JohnsonAllPaths returns a shortest-path tree for shortest paths in the graph g.
// If weight is nil and the graph does not implement graph.Coster, UniformCost is used.
//
// The time complexity of JohnsonAllPaths is O(|V|.|E|+|V|^2.log|V|).
func JohnsonAllPaths(g graph.Graph, weight graph.CostFunc) (paths AllShortest, ok bool) {
	jg := johnsonWeightAdjuster{
		g:      g,
		from:   g.Neighbors,
		to:     g.Neighbors,
		weight: weight,
	}
	switch g := g.(type) {
	case graph.DirectedGraph:
		jg.from = g.Successors
		jg.to = g.Predecessors
		jg.edgeTo = g.EdgeTo
	default:
		jg.edgeTo = g.EdgeBetween
	}
	if jg.weight == nil {
		if g, ok := g.(graph.Coster); ok {
			jg.weight = g.Cost
		} else {
			jg.weight = UniformCost
		}
	}

	nodes := g.NodeList()
	indexOf := make(map[int]int, len(nodes))
	for i, n := range nodes {
		indexOf[n.ID()] = i
	}
	sign := -1
	for {
		// Choose a random node ID until we find
		// one that is not in g.
		jg.q = sign * rand.Int()
		if _, exists := indexOf[jg.q]; !exists {
			break
		}
		sign *= -1
	}

	jg.bellmanFord = true
	jg.adjustBy, ok = BellmanFordFrom(johnsonGraphNode(jg.q), jg, nil)
	if !ok {
		return paths, false
	}

	jg.bellmanFord = false
	paths = DijkstraAllPaths(jg, nil)

	for i, u := range paths.nodes {
		hu := jg.adjustBy.WeightTo(u)
		for j, v := range paths.nodes {
			if i == j {
				continue
			}
			hv := jg.adjustBy.WeightTo(v)
			paths.dist.Set(i, j, paths.dist.At(i, j)-hu+hv)
		}
	}

	return paths, ok
}

type johnsonWeightAdjuster struct {
	q int
	g graph.Graph

	from, to func(graph.Node) []graph.Node
	edgeTo   func(graph.Node, graph.Node) graph.Edge
	weight   graph.CostFunc

	bellmanFord bool
	adjustBy    Shortest
}

var (
	_ graph.DirectedGraph = johnsonWeightAdjuster{}
	_ graph.Coster        = johnsonWeightAdjuster{}
)

func (g johnsonWeightAdjuster) NodeExists(n graph.Node) bool {
	if g.bellmanFord && n.ID() == g.q {
		return true
	}
	return g.g.NodeExists(n)

}

func (g johnsonWeightAdjuster) NodeList() []graph.Node {
	if g.bellmanFord {
		return append(g.g.NodeList(), johnsonGraphNode(g.q))
	}
	return g.g.NodeList()
}

func (g johnsonWeightAdjuster) Successors(n graph.Node) []graph.Node {
	if g.bellmanFord && n.ID() == g.q {
		return g.g.NodeList()
	}
	return g.from(n)
}

func (g johnsonWeightAdjuster) EdgeTo(u, v graph.Node) graph.Edge {
	if g.bellmanFord && u.ID() == g.q && g.g.NodeExists(v) {
		return concrete.Edge{johnsonGraphNode(g.q), v}
	}
	return g.edgeTo(u, v)
}

func (g johnsonWeightAdjuster) Cost(e graph.Edge) float64 {
	if g.bellmanFord {
		switch g.q {
		case e.From().ID():
			return 0
		case e.To().ID():
			return math.Inf(1)
		default:
			return g.weight(e)
		}
	}
	return g.weight(e) + g.adjustBy.WeightTo(e.From()) - g.adjustBy.WeightTo(e.To())
}

func (johnsonWeightAdjuster) Neighbors(graph.Node) []graph.Node {
	panic("search: unintended use of johnsonWeightAdjuster")
}
func (johnsonWeightAdjuster) EdgeBetween(_, _ graph.Node) graph.Edge {
	panic("search: unintended use of johnsonWeightAdjuster")
}
func (johnsonWeightAdjuster) Predecessors(graph.Node) []graph.Node {
	panic("search: unintended use of johnsonWeightAdjuster")
}
