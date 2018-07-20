// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
)

func YenKShortestPath(g graph.Graph, k int, s, t graph.Node) [][]graph.Node {
	yk := yenKSPAdjuster{
		g:       g,
		visited: make(map[int64]set.Int64s),
	}

	var weight Weighting
	if wg, ok := g.(Weighted); ok {
		weight = wg.Weight
	} else {
		weight = UniformCost(g)
	}

	yk.weight = weight

	var paths [][]graph.Node

	shortest, _ := DijkstraFrom(s, yk).To(t.ID())
	if len(shortest) != 0 {
		paths = append(paths, shortest)
	} else {
		return paths
	}

	var pot []yenShortest

	for i := int64(1); i < int64(k); i++ {
		for n := 0; n < (len(paths[i-1]) - 1); n++ {
			spur := paths[i-1][n]
			root := make([]graph.Node, len(paths[i-1][:n+1]))
			copy(root, paths[i-1][:n+1])

			var rootWeight float64
			for x := 1; x < len(root); x++ {
				w, _ := weight(root[x-1].ID(), root[x].ID())
				rootWeight += w
			}

			for _, path := range paths {
				if len(path) > n {
					ok := true
					for x := 0; x < len(root); x++ {
						if path[x].ID() != root[x].ID() {
							ok = false
							break
						}
					}
					if ok {
						yk.addVisited(path[n].ID(), path[n+1].ID())
					}
				}
			}

			spath, weight := DijkstraFrom(spur, yk).To(t.ID())
			size := len(root) - 1

			if len(root) > 1 {
				nroot := root[:size]
				nroot = append(nroot, spath...)
				potential := yenShortest{nroot, weight + rootWeight}
				pot = append(pot, potential)
			} else {
				potential := yenShortest{spath, weight}
				pot = append(pot, potential)
			}

			yk.visited = make(map[int64]set.Int64s)
		}

		if len(pot) == 0 {
			break
		}

		sort.Sort(byPathWeight(pot))

		paths = append(paths, pot[0].p)

		pot = pot[1:]
	}

	return paths
}

type yenShortest struct {
	p      []graph.Node
	weight float64
}

type byPathWeight []yenShortest

func (s byPathWeight) Len() int {
	return len(s)
}

func (s byPathWeight) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPathWeight) Less(i, j int) bool {
	return s[i].weight < s[j].weight
}

type yenKSPAdjuster struct {
	g       graph.Graph
	visited map[int64]set.Int64s
	weight  Weighting
}

func (g yenKSPAdjuster) From(id int64) []graph.Node {
	nodes := g.g.From(id)

	for i := 0; i < len(nodes); i++ {
		if g.visited[id].Has(int64(nodes[i].ID())) {
			nodes[int64(i)] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
			i--
		}
	}

	return nodes
}

func (g yenKSPAdjuster) addVisited(parent, id int64) {
	visited, ok := g.visited[parent]
	if !ok {
		visited = make(set.Int64s)
		g.visited[parent] = visited
	}

	visited.Add(id)
}

func (g yenKSPAdjuster) Edge(uid, vid int64) graph.Edge {
	return g.g.Edge(uid, vid)
}

func (g yenKSPAdjuster) Weight(xid, yid int64) (w float64, ok bool) {
	return g.weight(xid, yid)
}

func (g yenKSPAdjuster) HasEdgeBetween(xid, yid int64) bool {
	return g.g.HasEdgeBetween(xid, yid)
}
