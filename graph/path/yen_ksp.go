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
	p   []graph.Node
	weight float64
}

func YenKSP(source graph.Node, sink graph.Node, g graph.Weighted, k int) [][]graph.Node {
	yk := yenKSPAdjuster{
		g:       g,
		visited: make(map[int64][]int64),
		from:    g.From,
	}


	paths := make([][]graph.Node, k)

	paths[0], _ = DijkstraFrom(source, yk).To(sink.ID())
	
	var pot []YenShortest

	for i := int64(1); i < int64(k); i++ {
		for n := 0; n < (len(paths[i-1])-1); n++ {
			spur := paths[i-1][n]
			//fmt.Printf("SPUR: %d \n", spur)
			root := make([]graph.Node, len(paths[i-1][:n + 1]))
			copy(root, paths[i-1][:n + 1])

			//fmt.Println(root)
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

			if spur.ID() == 2 {
				fmt.Println(yk.visited[2])
			}

			spath, weight := DijkstraFrom(spur, yk).To(sink.ID())
			size := len(root) - 1
			//fmt.Println(spath)
			//fmt.Println(weight)
			if len(root) > 1 {
				nroot := root[:size]
				nroot = append(nroot, spath...)
				potential := YenShortest{nroot, weight + rootWeight}
				if i == 2 {
					fmt.Println(potential)
				}
				pot = append(pot, potential)
			} else {
				potential := YenShortest{spath, weight}
				if i == 2 {
					fmt.Println(potential)
				}
				pot = append(pot, potential)
			}

			//fmt.Println(pot)

			yk.visited = make(map[int64][]int64)
		}

		//fmt.Printf("POTENTIAL: ")
		//fmt.Println(pot)
		if len(pot) == 0 {
			break
		}

		sort.Slice(pot, func(a, b int) bool {
			return pot[a].weight < pot[b].weight
		})

		//fmt.Printf("SORTED POTENTIAL: ")
		//fmt.Println(pot)

		paths[i] = pot[0].p
		
		pot = pot[1:]
	}

	return paths
}

type yenKSPAdjuster struct {
	g graph.Weighted
	visited map[int64][]int64

	from   func(id int64) []graph.Node
}

func (g yenKSPAdjuster) From(id int64) []graph.Node {
	nodes := g.from(id)

	for i := 0; i < len(nodes); i++{
		if contains(g.visited[id], int64(nodes[i].ID())) {
			nodes[int64(i)] = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
			i--;
		}
	}

	if id == 2 {
		fmt.Println("Final 2 List: ")
		fmt.Println(nodes)
	}
	return nodes
}

func (g yenKSPAdjuster) AddVisited(parent, id int64) {
	if !contains(g.visited[parent], id) {
		fmt.Printf("Visited %d , %d \n", parent, id)
		g.visited[parent] = append(g.visited[parent], id)
	}
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
