// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"fmt"
	"math/rand/v2"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

// ExampleLeiden_connectedCommunities demonstrates the key behavioural
// difference between Leiden and Louvain: Leiden guarantees that every
// community it produces is internally connected, while Louvain can merge
// nodes into a single community even when those nodes are only reachable
// from each other through a different community.
//
// The graph is a 100-node duplication-divergence network generated with
// seed 19. At resolution γ=1.5, Louvain places nodes {8,51,64,70,80,92}
// into one community, but nodes {8,70,80} and {51,64,92} share no internal
// edges — the community is disconnected. Leiden's refinement phase detects
// and corrects this, producing only connected communities.
//
// See [Traag, Waltman & Van Eck, Sci Rep 9, 5233 (2019)] for details.
//
// [Traag, Waltman & Van Eck, Sci Rep 9, 5233 (2019)]: https://doi.org/10.1038/s41598-019-41695-z
func ExampleLeiden_connectedCommunities() {
	// Build a 100-node duplication-divergence graph.
	// This seed is chosen because Louvain demonstrably produces a
	// disconnected community at γ=1.5 on this particular graph.
	g := simple.NewUndirectedGraph()
	err := gen.Duplication(g, 100, 0.8, 0.1, 0.5, rand.New(rand.NewPCG(19, 0)))
	if err != nil {
		panic(err)
	}

	const γ = 1.5
	src1 := rand.New(rand.NewPCG(4, 19))
	src2 :=  rand.New(rand.NewPCG(4, 19))
	rLouvain, err := Modularize(g, γ, src1)
	if err != nil {
		panic(err)
	}
	rLeiden, err := Leiden(g, γ, src2)
	if err != nil {
		panic(err)
	}

	// Check whether each algorithm produced any disconnected community.
	louvainOK := allConnected(g, rLouvain.Communities())
	leidenOK := allConnected(g, rLeiden.Communities())

	fmt.Printf("Louvain: %d communities, all connected: %v\n",
		len(rLouvain.Communities()), louvainOK)
	fmt.Printf("Leiden:  %d communities, all connected: %v\n",
		len(rLeiden.Communities()), leidenOK)

	// Output:
	// Louvain: 11 communities, all connected: false
	// Leiden:  13 communities, all connected: true
}

// allConnected reports whether every community in comms forms a connected
// subgraph of g.
func allConnected(g interface {
	From(int64) graph.Nodes
}, comms [][]graph.Node) bool {
	for _, comm := range comms {
		if len(comm) <= 1 {
			continue
		}
		members := make(map[int64]bool, len(comm))
		for _, n := range comm {
			members[n.ID()] = true
		}
		visited := make(map[int64]bool, len(comm))
		queue := []int64{comm[0].ID()}
		visited[comm[0].ID()] = true
		for len(queue) > 0 {
			uid := queue[0]
			queue = queue[1:]
			to := g.From(uid)
			for to.Next() {
				vid := to.Node().ID()
				if members[vid] && !visited[vid] {
					visited[vid] = true
					queue = append(queue, vid)
				}
			}
		}
		if len(visited) != len(comm) {
			return false
		}
	}
	return true
}
