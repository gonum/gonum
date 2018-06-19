// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/graph"
)

type YenShortest struct {
	path   []graph.Node
	weight float64
}

func YenKSP(source graph.Node, sink graph.Node, g graph.Graph, k int) [][]graph.Node {
	yk := yenKSPAdjuster{
		g:       g,
		visited: make([]int64, 0),
		from:    g.From,
		weight:  g.(graph.Weighted).Weight,
	}

	paths := make([][]graph.Node, k)

	paths[0], _ = DijkstraFrom(source, yk).To(sink.ID())

	var pot []YenShortest

	for i := int64(1); i < int64(k); i++ {
		fmt.Printf("%d \n", i)
		for n := 0; n < len(paths[i-1])-2; n++ {
			fmt.Printf("%d \n", n)
			spur := paths[i-1][n]
			root := paths[i-1][:n]

			var rootWeight float64
			for x := 1; x < len(root); x++ {
				w, _ := g.(graph.Weighted).Weight(root[x-1].ID(), root[x].ID())
				rootWeight += w
			}

			//var edges []graph.Edge
			//var nodes []graph.Node

			for _, path := range paths {
				for x := 0; x < len(root); x++ {
					if x < len(path) {
						if path[x].ID() != root[x].ID() {
							break
						}
					}
				}

				//if contains {
				//	edges = append(edges, g.Edge(i, i + 1))
				//	g.(graph.EdgeRemover).RemoveEdge(i, i + 1)
				//}
			}

			for _, node := range root {
				if node != spur {
					//nodes = append(nodes, node)
					yk.AddVisited(node.ID())
					//g.(graph.NodeRemover).RemoveNode(node.ID())
				}
			}

			spath, weight := DijkstraFrom(spur, yk).To(sink.ID())
			var total []graph.Node
			total = append(root, spath...)
			potential := YenShortest{total, weight + rootWeight}
			pot = append(pot, potential)

			//for _, edge := range edges {
			//	g.(graph.WeightedEdgeAdder).SetWeightedEdge(edge.(graph.WeightedEdge))
			//}

		}

		if len(pot) == 0 {
			break
		}

		sort.Slice(pot, func(a, b int) bool {
			return pot[a].weight < pot[b].weight
		})

		paths[i] = pot[0].path
		pot = pot[1:]
	}

	return paths
}

// TODO Wrap graph in YenKSPAdjuster interface and override the From method

type yenKSPAdjuster struct {
	g       graph.Graph
	visited []int64

	from   func(id int64) []graph.Node
	edgeTo func(uid, vid int64) graph.Edge
	weight func(xid, yid int64) (w float64, ok bool)
}

func (g yenKSPAdjuster) From(id int64) []graph.Node {
	if contains(g.visited, id) {
		return nil
	} else {
		nodes := g.from(id)
		for i, _ := range nodes {
			if contains(g.visited, int64(i)) {
				nodes[int64(i)] = nodes[len(nodes)-1]
				nodes = nodes[:len(nodes)-1]
			}
		}
		return nodes
	}
}

func (g yenKSPAdjuster) AddVisited(id int64) {
	g.visited = append(g.visited, id)
}

func (g yenKSPAdjuster) Edge(uid, vid int64) graph.Edge {
	return g.edgeTo(uid, vid)
}

func (g yenKSPAdjuster) Weight(xid, yid int64) (w float64, ok bool) {
	return g.weight(xid, yid)
}

func contains(visited []int64, id int64) bool {
	for _, n := range visited {
		if n == id {
			return true
		}
	}
	return false
}
