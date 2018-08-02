// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"sort"

	"gonum.org/v1/gonum/graph"
)

// YenKShortestPath returns the k-shortest loopless paths from s to t in g. YenKShortestPath will
// panic if g contains a negative edge weight.
func YenKShortestPath(g graph.Graph, k int, s, t graph.Node) [][]graph.Node {
	yk := yenKSPAdjuster{
		g:       g,
		visited: make(map[[2]int64]bool),
	}

	if wg, ok := g.(Weighted); ok {
		yk.weight = wg.Weight
	} else {
		yk.weight = UniformCost(g)
	}

	shortest, _ := DijkstraFrom(s, yk).To(t.ID())
	if len(shortest) == 0 {
		return nil

	}
	paths := [][]graph.Node{shortest}

	var pot []yenShortest
	for i := int64(1); i < int64(k); i++ {
		for n := 0; n < (len(paths[i-1]) - 1); n++ {
			spur := paths[i-1][n]
			root := make([]graph.Node, len(paths[i-1][:n+1]))
			copy(root, paths[i-1][:n+1])

			var rootWeight float64
			for x := 1; x < len(root); x++ {
				w, _ := yk.weight(root[x-1].ID(), root[x].ID())
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
			if len(root) > 1 {
				root = append(root[:len(root)-1], spath...)
				pot = append(pot, yenShortest{root, weight + rootWeight})
			} else {
				pot = append(pot, yenShortest{spath, weight})
			}

			yk.visited = make(map[[2]int64]bool)
		}

		if len(pot) == 0 {
			break
		}

		sort.Sort(byPathWeight(pot))
		paths = append(paths, pot[0].path)
		pot = pot[1:]
	}

	return paths
}

type yenShortest struct {
	path   []graph.Node
	weight float64
}

type byPathWeight []yenShortest

func (s byPathWeight) Len() int           { return len(s) }
func (s byPathWeight) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byPathWeight) Less(i, j int) bool { return s[i].weight < s[j].weight }

// A wrapper for the graph imported to the YenKShotestPath functions.
// This wrapper modifies the From function by not returning edges in the visited map.
// The visited map is updated through each iteration of the YenKShortestPath function.
type yenKSPAdjuster struct {
	g       graph.Graph
	visited map[[2]int64]bool
	weight  Weighting
}

func (g yenKSPAdjuster) From(id int64) []graph.Node {
	nodes := g.g.From(id)

	for i := 0; i < len(nodes); i++ {
		if g.visited[[2]int64{id, nodes[i].ID()}] {
			nodes[int64(i)] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
			i--
		}
	}

	return nodes
}

func (g yenKSPAdjuster) addVisited(u, v int64) {
	g.visited[[2]int64{u, v}] = true

	/*
		if g.isDirected {
			g.visited[[2]int64{v, u}] = true
		}*/
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
