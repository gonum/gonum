// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"math"
	"math/rand"

	"github.com/gonum/graph"
	"github.com/gonum/matrix/mat64"
)

// FloydWarshall returns a shortest-path tree for the graph g or false indicating
// that a negative cycle exists in the graph. If weight is nil and the graph does not
// implement graph.Coster, UniformCost is used.
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

		dist: mat64.NewDense(len(nodes), len(nodes), dist),
		next: make([][]int, len(nodes)*len(nodes)),
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
				} else if i != k && k != j && ij-joint == 0 {
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

// ShortestPaths is a shortest-path tree created by the FloydWarshall function.
type ShortestPaths struct {
	// nodes hold the nodes of the analysed
	// graph.
	nodes []graph.Node
	// indexOf contains a mapping between
	// the id-dense representation of the
	// graph and the potentially id-sparse
	// nodes held in nodes.
	indexOf map[int]int

	// dist and next represent the shortest
	// paths between nodes.
	//
	// Indices into dist and next are
	// mapped through indexOf.
	//
	// dist contains the pairwise
	// distances between nodes.
	dist *mat64.Dense
	// next contains the shortest-path
	// tree of the graph. The first index
	// is a linear mapping of from-dense-id
	// and to-dense-id; the second is the
	// intermediate nodes by dense id.
	next [][]int
}

func (p ShortestPaths) at(from, to int) (mid []int) {
	return p.next[from+to*len(p.nodes)]
}

func (p ShortestPaths) set(from, to int, weight float64, mid ...int) {
	p.dist.Set(from, to, weight)
	p.next[from+to*len(p.nodes)] = append(p.next[from+to*len(p.nodes)][:0], mid...)
}

func (p ShortestPaths) add(from, to int, mid ...int) {
loop: // These are likely to be rare, so just loop over collisions.
	for _, k := range mid {
		for _, v := range p.next[from+to*len(p.nodes)] {
			if k == v {
				continue loop
			}
		}
		p.next[from+to*len(p.nodes)] = append(p.next[from+to*len(p.nodes)], k)
	}
}

// Weight returns the weight of the minimum path between u and v.
func (p ShortestPaths) Weight(u, v graph.Node) float64 {
	from, to := p.indexOf[u.ID()], p.indexOf[v.ID()]
	if from < 0 || from >= len(p.nodes) || to < 0 || to >= len(p.nodes) {
		return math.Inf(1)
	}
	return p.dist.At(p.indexOf[u.ID()], p.indexOf[v.ID()])
}

// Between returns a shortest path from u to v and the weight of the path. If more than
// one shortest path exists between u and v, a randomly chosen path will be returned and
// unique is returned false.
func (p ShortestPaths) Between(u, v graph.Node) (path []graph.Node, weight float64, unique bool) {
	from, fromOK := p.indexOf[u.ID()]
	to, toOK := p.indexOf[v.ID()]
	if !fromOK || !toOK || len(p.at(from, to)) == 0 {
		return nil, math.Inf(1), false
	}
	path = []graph.Node{p.nodes[from]}
	unique = true
	for from != to {
		c := p.at(from, to)
		if len(c) != 1 {
			unique = false
		}
		from = c[rand.Intn(len(c))]
		path = append(path, p.nodes[from])
	}
	return path, p.dist.At(p.indexOf[u.ID()], p.indexOf[v.ID()]), unique
}

// AllBetween returns all shortest paths from u to v and the weight of the paths.
func (p ShortestPaths) AllBetween(u, v graph.Node) (paths [][]graph.Node, weight float64) {
	from, fromOK := p.indexOf[u.ID()]
	to, toOK := p.indexOf[v.ID()]
	if !fromOK || !toOK || len(p.at(from, to)) == 0 {
		return nil, math.Inf(1)
	}
	paths = p.allBetween(from, to, []graph.Node{p.nodes[from]}, nil)
	return paths, p.dist.At(p.indexOf[u.ID()], p.indexOf[v.ID()])
}

func (p ShortestPaths) allBetween(from, to int, path []graph.Node, paths [][]graph.Node) [][]graph.Node {
	if from == to {
		return append(paths, path)
	}
	for i, from := range p.at(from, to) {
		if i != 0 {
			path = append([]graph.Node(nil), path...)
		}
		paths = p.allBetween(from, to, append(path, p.nodes[from]), paths)
	}
	return paths
}
