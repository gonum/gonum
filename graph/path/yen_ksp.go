// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"reflect"
	"sort"
	
	"gonum.org/v1/gonum/graph"
)

type YenShortest struct {
	path []graph.node,
	weight float64
}

func YenKSP(source graph.Node, sink graph.Node, g graph.Graph, k int) []Shortest {
	paths := make([]Shortest, k)

	rootPath, _ = DijkstraFrom(source, g).To(sink.ID())

	paths[0] = newShortestfrom(source, rootPath)
	
	var pot []YenShortest
	
	for i := 1; i < k; i++ {
		for n := 0; n < len(paths[i - 1].nodes) - 2; n++ {
			spur := paths[i - 1].nodes[i]
			root := paths[i - 1].nodes[:i]

			var rootWeight float64
			for x := 1; x < len(root); x++ {
				rootWeight += g.Weight(root[x - 1], root[x])
			}
			
			var edges []graph.Edge
			var nodes []graph.Node

			for _, p := range paths {
				if reflect.DeepEqual(p.nodes[:i], root) {
					append(edges, g.Edge(i, i + 1))
					g.RemoveEdge(i, i + 1)
				}
			}

			for _, node := range root {
				if (node != spur) {
					append(nodes, node)
					g.RemoveNode(node)
				}
			}

			spath, weight := DjikstraFrom(spur, g).To(sink.ID())
			var total []graph.Node
			append(total, root, spath)
			potential := YenShortest {total, weight + rootWeight}
			append(pot, potential)
			
			for _, edge := range edges {
				g.SetEdge(edge)
			}

		}

		if (len(pot) == 0) {
			break
		}

		sort.Slice(pot, func(a, b int) bool {
			return pot[a].weight < pot[b].weight
		})

		
		paths[i] = newShortestfrom(source, pot[0].path)
		pot = pot[1:]
	}
	
	return paths
}
