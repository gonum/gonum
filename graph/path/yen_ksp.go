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
	path []graph.Node
	weight float64
}

func YenKSP(source graph.Node, sink graph.Node, g graph.Graph, k int) [][]graph.Node {
	paths := make([][]graph.Node, k)

	paths[0], _ = DijkstraFrom(source, g).To(sink.ID())

	var pot []YenShortest
	
	for i := int64(1); i < int64(k); i++ {
		fmt.Printf("%d \n", i)
		for n := 0; n < len(paths[i - 1]) - 2; n++ {
			fmt.Printf("%d \n", n)
			spur := paths[i - 1][n]
			root := paths[i - 1][:n]

			var rootWeight float64
			for x := 1; x < len(root); x++ {
				w, _ := g.(graph.Weighted).Weight(root[x - 1].ID(), root[x].ID())
				rootWeight += w
			}
			
			var edges []graph.Edge
			var nodes []graph.Node

			for _, path := range paths {
				contains := true
				for x := 0; x < len(root); x++ {
					if path[x].ID() != root[x].ID() {
						contains = false
						break
					}
				}
				
				if contains {
					edges = append(edges, g.Edge(i, i + 1))
					g.(graph.EdgeRemover).RemoveEdge(i, i + 1)
				}
			}

			for _, node := range root {
				if (node != spur) {
					nodes = append(nodes, node)
					g.(graph.NodeRemover).RemoveNode(node.ID())
				}
			}

			spath, weight := DijkstraFrom(spur, g).To(sink.ID())
			var total []graph.Node
			total = append(root, spath...)
			potential := YenShortest {total, weight + rootWeight}
			pot = append(pot, potential)
			
			for _, edge := range edges {
				g.(graph.WeightedEdgeAdder).SetWeightedEdge(edge.(graph.WeightedEdge))
			}

		}

		if (len(pot) == 0) {
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
