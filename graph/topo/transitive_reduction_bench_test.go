// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

type edge struct{ from, to int64 }

func BenchmarkTransitiveReduce(b *testing.B) {

	run := func(b *testing.B, n int, edges []edge) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			g := simple.NewDirectedGraph()
			for id := 0; id < n; id++ {
				g.AddNode(simple.Node(int64(id)))
			}
			for _, e := range edges {
				g.SetEdge(g.NewEdge(simple.Node(e.from), simple.Node(e.to)))
			}
			b.StartTimer()

			if err := TransitiveReduce(g); err != nil {
				b.Fatalf("reduction error: %v", err)
			}
		}
	}

	sizes := []int{50, 250, 500}
	densities := []float64{0.02, 0.10, 0.30}
	seed := uint64(1)

	for _, n := range sizes {
		for _, p := range densities {
			name := fmt.Sprintf("TransitiveReduction/RandomDAG/n=%d/p=%d%%", n, int(p*100+0.5))
			edges := makeRandomDAGEdges(n, p, seed+uint64(n*1000)+uint64(p*100))
			b.Run(name, func(b *testing.B) { run(b, n, edges) })
		}
	}

	for _, n := range []int{50, 100, 200} {
		name := fmt.Sprintf("TransitiveReduction/CompleteDAG/n=%d", n)
		edges := makeCompleteDAGEdges(n)
		b.Run(name, func(b *testing.B) { run(b, n, edges) })
	}
}

// ---- Graph families ----

// makeRandomDAGEdges creates a DAG by only adding edges from i->j where i<j with probability p.
// Deterministic due to seed.
func makeRandomDAGEdges(n int, p float64, seed uint64) []edge {
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	edges := make([]edge, 0, int(float64(n*n)*p/2))

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float64() < p {
				edges = append(edges, edge{from: int64(i), to: int64(j)})
			}
		}
	}
	return edges
}

// makeCompleteDAGEdges returns all edges i->j for i<j.
// This is dense and produces many redundant edges.
func makeCompleteDAGEdges(n int) []edge {
	edges := make([]edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, edge{from: int64(i), to: int64(j)})
		}
	}
	return edges
}
