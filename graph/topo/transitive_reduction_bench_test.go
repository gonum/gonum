package topo

import (
	"fmt"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/graph/simple"
)

func BenchmarkTransitiveReduce(b *testing.B) {
	benchReducers(b, "TransitiveReduction1", TransitiveReduce)
}

func BenchmarkTransitiveReduce2(b *testing.B) {
	benchReducers(b, "TransitiveReduction2", TransitiveReduce2)
}

type reducerFn func(g GraphReducer) error

func benchReducers(b *testing.B, label string, fn reducerFn) {
	sizes := []int{50, 250, 500}
	densities := []float64{0.02, 0.10, 0.30} // sparse -> dense-ish

	seed := int64(1)

	// 1) Random DAGs at multiple densities
	for _, n := range sizes {
		for _, p := range densities {
			name := fmt.Sprintf("%s/RandomDAG/n=%d/p=%d%%", label, n, int(p*100+0.5))
			edges := makeRandomDAGEdges(n, p, seed+int64(n*1000)+int64(p*100))

			b.Run(name, func(b *testing.B) {
				runReduceBenchmark(b, n, edges, fn)
			})
		}
	}

	// 2) Worst-case-ish: complete DAG (edges i->j for all i<j).
	// Reduction should shrink it to just the chain edges (i->i+1).
	for _, n := range []int{50, 100, 200} {
		name := fmt.Sprintf("%s/CompleteDAG/n=%d", label, n)
		edges := makeCompleteDAGEdges(n)

		b.Run(name, func(b *testing.B) {
			runReduceBenchmark(b, n, edges, fn)
		})
	}
}

// ---- Core runner ----

type edge struct{ from, to int64 }

// runReduceBenchmark builds a fresh graph per iteration from a fixed edge list, then runs fn(g).
// Graph construction is done outside the timed section (best-effort), so you primarily measure reduction cost.
func runReduceBenchmark(b *testing.B, n int, edges []edge, fn reducerFn) {
	b.ReportAllocs()

	// Pre-create nodes once; we’ll reuse IDs when constructing graphs.
	nodeIDs := make([]int64, n)
	for i := 0; i < n; i++ {
		nodeIDs[i] = int64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Build graph outside timer.
		b.StopTimer()
		g := simple.NewDirectedGraph()

		for _, id := range nodeIDs {
			g.AddNode(simple.Node(id))
		}
		for _, e := range edges {
			g.SetEdge(g.NewEdge(simple.Node(e.from), simple.Node(e.to)))
		}
		b.StartTimer()

		if err := fn(g); err != nil {
			b.Fatalf("reduction error: %v", err)
		}
	}
}

// ---- Graph families ----

// makeRandomDAGEdges creates a DAG by only adding edges from i->j where i<j with probability p.
// Deterministic due to seed.
func makeRandomDAGEdges(n int, p float64, seed int64) []edge {
	rng := rand.New(rand.NewSource(seed))
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
