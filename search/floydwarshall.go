// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"math"

	"github.com/gonum/graph"
	"github.com/gonum/matrix/mat64"
)

// FloydWarshall returns a shortest-path tree for the graph g or false indicating
// that a negative cycle exists in the graph. If weight is nil and the graph does not
// implement graph.Coster, UniformCost is used.
//
// The time complexity of FloydWarshall is O(|V|^3).
func FloydWarshall(g graph.Graph, weight graph.CostFunc) (paths ShortestPaths, ok bool) {
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

	indexOf := make(map[int]int, len(nodes))
	for i, n := range nodes {
		indexOf[n.ID()] = i
	}

	dist := make([]float64, len(nodes)*len(nodes))
	for i := range dist {
		dist[i] = math.Inf(1)
	}
	paths = ShortestPaths{
		nodes:   nodes,
		indexOf: indexOf,

		dist:    mat64.NewDense(len(nodes), len(nodes), dist),
		next:    make([][]int, len(nodes)*len(nodes)),
		forward: true,
	}
	for i, u := range nodes {
		paths.dist.Set(i, i, 0)
		for _, v := range from(u) {
			j := indexOf[v.ID()]
			paths.set(i, j, weight(edgeTo(u, v)), j)
		}
	}

	for k := range nodes {
		for i := range nodes {
			for j := range nodes {
				ij := paths.dist.At(i, j)
				joint := paths.dist.At(i, k) + paths.dist.At(k, j)
				if ij > joint {
					paths.set(i, j, joint, paths.at(i, k)...)
				} else if ij-joint == 0 {
					paths.add(i, j, paths.at(i, k)...)
				}
			}
		}
	}

	ok = true
	for i := range nodes {
		if paths.dist.At(i, i) < 0 {
			ok = false
			break
		}
	}

	return paths, ok
}
