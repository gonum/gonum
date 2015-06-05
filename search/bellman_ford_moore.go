// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import "github.com/gonum/graph"

// BellmanFordFrom returns a shortest-path tree for a shortest path from u to all nodes in
// the graph g, or false indicating that a negative cycle exists in the graph. If weight
// is nil and the graph does not implement graph.Coster, UniformCost is used.
func BellmanFordFrom(u graph.Node, g graph.Graph, weight graph.CostFunc) (path Shortest, ok bool) {
	if !g.NodeExists(u) {
		return Shortest{from: u}, true
	}
	var (
		from   = g.Neighbors
		edgeTo func(graph.Node, graph.Node) graph.Edge
	)
	switch g := g.(type) {
	case graph.DirectedGraph:
		from = g.Successors
		edgeTo = g.EdgeTo
	default:
		edgeTo = g.EdgeBetween
	}
	if weight == nil {
		if g, ok := g.(graph.Coster); ok {
			weight = g.Cost
		} else {
			weight = UniformCost
		}
	}

	nodes := g.NodeList()

	path = newShortestFrom(u, nodes)
	path.dist[path.indexOf[u.ID()]] = 0

	// TODO(kortschak): Consider adding further optimisations
	// from http://arxiv.org/abs/1111.5414.
	for i := 1; i < len(nodes); i++ {
		changed := false
		for j, u := range nodes {
			for _, v := range from(u) {
				k := path.indexOf[v.ID()]
				joint := path.dist[j] + weight(edgeTo(u, v))
				if joint < path.dist[k] {
					path.set(k, joint, j)
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	for j, u := range nodes {
		for _, v := range from(u) {
			k := path.indexOf[v.ID()]
			if path.dist[j]+weight(edgeTo(u, v)) < path.dist[k] {
				return path, false
			}
		}
	}

	return path, true
}
