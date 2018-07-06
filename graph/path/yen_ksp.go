// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
)

type yenShortest struct {
	p   []graph.Node
	weight float64
}

func YenKSP(source graph.Node, sink graph.Node, g graph.Weighted, k int) [][]graph.Node {
	yk := yenKSPAdjuster{
		g:       g,
		visited: make(map[int64]set.Int64s),
		from:    g.From,
	}

	paths := make([][]graph.Node, k)

	paths[0], _ = DijkstraFrom(source, yk).To(sink.ID())
	
	var pot []yenShortest

	for i := int64(1); i < int64(k); i++ {
		for n := 0; n < (len(paths[i-1])-1); n++ {
			spur := paths[i-1][n]
			root := make([]graph.Node, len(paths[i-1][:n + 1]))
			copy(root, paths[i-1][:n + 1])

			var rootWeight float64
			for x := 1; x < len(root); x++ {
				w, _ := g.Weight(root[x-1].ID(), root[x].ID())
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
						yk.AddVisited(path[n].ID(), path[n + 1].ID())
					}
				}
			}

			spath, weight := DijkstraFrom(spur, yk).To(sink.ID())
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

		sort.Slice(pot, func(a, b int) bool {
			return pot[a].weight < pot[b].weight
		})

		paths[i] = pot[0].p
		
		pot = pot[1:]
	}

	return paths
}

type yenKSPAdjuster struct {
	g graph.Weighted
	visited map[int64]set.Int64s

	from   func(id int64) []graph.Node
}

func (g yenKSPAdjuster) From(id int64) []graph.Node {
	nodes := g.from(id)

	for i := 0; i < len(nodes); i++{
		if g.visited[id].Has(int64(nodes[i].ID())) {
			nodes[int64(i)] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
			i--;
		}
	}

	return nodes
}

func (g yenKSPAdjuster) AddVisited(parent, id int64) {
	if g.visited[parent] == nil {
		g.visited[parent] = make(set.Int64s)
	}
	
	g.visited[parent].Add(id)
}

func (g yenKSPAdjuster) Edge(uid, vid int64) graph.Edge {
	return g.g.Edge(uid, vid)
}

func (g yenKSPAdjuster) Weight(xid, yid int64) (w float64, ok bool) {
	return g.g.Weight(xid, yid)
}

func (g yenKSPAdjuster) HasEdgeBetween(xid, yid int64) bool {
	return g.g.HasEdgeBetween(xid, yid)
}
	

func contains(visited []int64, id int64) bool {
	for _, n := range visited {
		if n == id {
			return true
		}
	}
	return false
}
