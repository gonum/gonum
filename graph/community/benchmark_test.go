// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package community

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"path"
	"strings"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/graph6"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/internal/order"
)

// The original data sets are from https://www.kaggle.com/datasets/dataup1/ogbn-arxiv
//
//go:embed testdata/*.graph6.gz
var testdataGraphs embed.FS

// TestLeidenVsLouvainArxiv compares Leiden and Louvain on ogbn-arxiv
// subgraphs across multiple resolutions. Leiden is expected to identify
// more communities by splitting weakly connected ones that Louvain merges,
// while maintaining comparable modularity.
func TestLeidenVsLouvainArxiv(t *testing.T) {
	for _, test := range []arXivTestGraph{
		arXivTestGraphs[0],
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, γ := range []float64{0.5, 1.0, 1.5} {
				t.Run(fmt.Sprintf("γ=%.1f", γ), func(t *testing.T) {
					const iterations = 3

					bestQLouvain, bestQLeiden := math.Inf(-1), math.Inf(-1)
					var bestCommLouvain, bestCommLeiden int

					for i := 0; i < iterations; i++ {
						srcL := rand.New(rand.NewPCG(uint64(i), 0))
						srcD := rand.New(rand.NewPCG(uint64(i), 0))

						rLouvain, err := Modularize(test.g, γ, srcL)
						if err != nil {
							t.Fatalf("unexpected error from Modularize: %v", err)
						}
						rLeiden, err := Leiden(test.g, γ, srcD)
						if err != nil {
							t.Fatalf("unexpected error from Leiden: %v", err)
						}

						qL := Q(rLouvain, nil, γ)
						qD := Q(rLeiden, nil, γ)

						if qL > bestQLouvain {
							bestQLouvain = qL
							bestCommLouvain = len(rLouvain.Communities())
						}
						if qD > bestQLeiden {
							bestQLeiden = qD
							bestCommLeiden = len(rLeiden.Communities())
						}
					}

					t.Logf("Louvain: Q=%.6f (%d communities)", bestQLouvain, bestCommLouvain)
					t.Logf("Leiden:  Q=%.6f (%d communities)", bestQLeiden, bestCommLeiden)
					t.Logf("ΔQ=%+.6f  ΔC=%+d",
						bestQLeiden-bestQLouvain, bestCommLeiden-bestCommLouvain)

					// Leiden should not produce drastically worse modularity.
					if bestQLeiden < bestQLouvain-0.05 {
						t.Errorf("Leiden Q (%.6f) significantly worse than Louvain Q (%.6f)",
							bestQLeiden, bestQLouvain)
					}
				})
			}
		})
	}
}

// TestLeidenVsLouvainDisconnected demonstrates Louvain's known weakness:
// it can produce communities that are not internally connected.
//
// This uses a 100-node duplication-divergence graph (seed=19) at γ=1.5
// where Louvain merges nodes {8,51,64,70,80,92} into one community,
// but nodes {8,70,80} are not reachable from {51,64,92} via edges
// internal to that community. Leiden, by design, guarantees connected
// communities and also achieves better or comparable modularity.
func TestLeidenVsLouvainDisconnected(t *testing.T) {
	g := simple.NewUndirectedGraph()
	err := gen.Duplication(g, 100, 0.8, 0.1, 0.5, rand.New(rand.NewPCG(19, 0)))
	if err != nil {
		t.Fatalf("failed to generate graph: %v", err)
	}

	const γ = 1.5

	srcLouvain := rand.New(rand.NewPCG(4, 19))
	rLouvain, err := Modularize(g, γ, srcLouvain)
	if err != nil {
		t.Fatalf("unexpected error from Modularize: %v", err)
	}
	louvainComms := rLouvain.Communities()
	qLouvain := Q(rLouvain, nil, γ)

	srcLeiden := rand.New(rand.NewPCG(4, 19))
	rLeiden, err := Leiden(g, γ, srcLeiden)
	if err != nil {
		t.Fatalf("unexpected error from Leiden: %v", err)
	}
	leidenComms := rLeiden.Communities()
	qLeiden := Q(rLeiden, nil, γ)

	t.Logf("Louvain: %d communities, Q=%.6f", len(louvainComms), qLouvain)
	t.Logf("Leiden:  %d communities, Q=%.6f", len(leidenComms), qLeiden)

	// Count and report disconnected communities produced by Louvain.
	louvainDisconnected := 0
	for _, comm := range louvainComms {
		if len(comm) <= 1 {
			continue
		}
		if !isCommunityConnected(g, comm) {
			louvainDisconnected++
			order.ByID(comm)
			ids := make([]int64, len(comm))
			for j, n := range comm {
				ids[j] = n.ID()
			}
			t.Logf("Louvain disconnected community (size %d): %v", len(comm), ids)
		}
	}

	if louvainDisconnected == 0 {
		t.Error("expected at least one disconnected Louvain community on this graph")
	}

	// All Leiden communities must be connected.
	for i, comm := range leidenComms {
		if len(comm) <= 1 {
			continue
		}
		if !isCommunityConnected(g, comm) {
			order.ByID(comm)
			t.Errorf("Leiden community %d (size %d) is not connected", i, len(comm))
		}
	}

	// Leiden should achieve at least comparable modularity.
	if qLeiden < qLouvain-0.05 {
		t.Errorf("Leiden Q (%.6f) significantly worse than Louvain Q (%.6f)", qLeiden, qLouvain)
	}
}

// isCommunityConnected reports whether all nodes in comm form a connected
// subgraph within g. It runs a BFS from the first node, restricted to
// edges whose both endpoints are in the community.
func isCommunityConnected(g graph.Undirected, comm []graph.Node) bool {
	if len(comm) <= 1 {
		return true
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

	return len(visited) == len(comm)
}

// BenchmarkLeidenResolution benchmarks the Leiden algorithm across
// multiple resolutions on the duplication-divergence graph.
func BenchmarkLeidenResolution(b *testing.B) {
	type benchCase struct {
		name string
		g    graph.Undirected
	}
	graphs := []benchCase{
		{name: "dupGraph", g: dupGraph},
	}

	for _, bg := range graphs {
		for _, γ := range []float64{0.5, 1, 2, 5, 10} {
			b.Run(fmt.Sprintf("%s_γ=%.1f", bg.name, γ), func(b *testing.B) {
				src := rand.NewPCG(1, 1)
				for i := 0; i < b.N; i++ {
					Leiden(bg.g, γ, src)
				}
			})
		}
	}
}

// BenchmarkLouvainResolution benchmarks the Louvain algorithm across
// multiple resolutions on the duplication-divergence graph.
func BenchmarkLouvainResolution(b *testing.B) {
	for _, γ := range []float64{0.5, 1, 2, 5, 10} {
		b.Run(fmt.Sprintf("dupGraph_γ=%.1f", γ), func(b *testing.B) {
			src := rand.NewPCG(1, 1)
			for i := 0; i < b.N; i++ {
				Modularize(dupGraph, γ, src)
			}
		})
	}
}

func BenchmarkArXivGraphs(b *testing.B) {
	// BenchmarkLeidenArxiv benchmarks the Leiden algorithm on ogbn-arxiv subgraphs.
	b.Run("leiden", func(b *testing.B) {
		for _, test := range arXivTestGraphs {
			for _, γ := range []float64{1, 2, 5} {
				b.Run(fmt.Sprintf("E=%d_γ=%.1f", test.size, γ), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						src := rand.NewPCG(42, uint64(i))
						Leiden(test.g, γ, src)
					}
				})
			}
		}
	})

	// BenchmarkLouvainArxiv benchmarks the Louvain algorithm on ogbn-arxiv subgraphs.
	b.Run("louvain", func(b *testing.B) {
		for _, test := range arXivTestGraphs {
			for _, γ := range []float64{1, 2, 5} {
				b.Run(fmt.Sprintf("E=%d_γ=%.1f", test.size, γ), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						src := rand.NewPCG(42, uint64(i))
						Modularize(test.g, γ, src)
					}
				})
			}
		}
	})
}

type arXivTestGraph struct {
	name string
	g    graph.Undirected
	size int
}

var arXivTestGraphs []arXivTestGraph

func init() {
	for _, file := range []struct {
		name string
		size int
	}{
		0: {name: "arxiv_2500.graph6.gz", size: 2500},
		1: {name: "arxiv_5000.graph6.gz", size: 5000},
	} {
		g, err := loadGraph6(file.name)
		if err != nil {
			panic(err)
		}
		name, _, ok := strings.Cut(file.name, ".")
		if !ok {
			panic(fmt.Sprintf("invalid arxiv graph test file name: %s", file.name))
		}
		arXivTestGraphs = append(arXivTestGraphs, arXivTestGraph{
			name: name,
			g:    g,
			size: file.size,
		})
	}
}

// loadGraph6 loads a graph6-encoded graph from embedded testdata.
// The filename should be relative to testdata (e.g., "arxiv_5000.graph6.gz").
func loadGraph6(file string) (graph.Undirected, error) {
	f, err := testdataGraphs.Open(path.Join("testdata", file))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	g6, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	u := simple.NewUndirectedGraph()
	copyUndirected(u, graph6.Graph(bytes.TrimSpace(g6)))
	return u, nil
}

// copyUndirected is a specialised version of graph.Copy that does not
// make unnecessary edge set calls undirected graphs.
func copyUndirected(dst graph.Builder, src graph.Graph) {
	nodes := src.Nodes()
	for nodes.Next() {
		dst.AddNode(nodes.Node())
	}
	nodes.Reset()
	for nodes.Next() {
		u := nodes.Node()
		uid := u.ID()
		to := src.From(uid)
		for to.Next() {
			vid := to.Node().ID()
			if uid < vid {
				dst.SetEdge(src.Edge(uid, vid))
			}
		}
	}
}
