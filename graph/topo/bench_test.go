// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package topo

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	gnpDirected_10_tenth   = gnpDirected(10, 0.1)
	gnpDirected_100_tenth  = gnpDirected(100, 0.1)
	gnpDirected_1000_tenth = gnpDirected(1000, 0.1)
	gnpDirected_10_half    = gnpDirected(10, 0.5)
	gnpDirected_100_half   = gnpDirected(100, 0.5)
	gnpDirected_1000_half  = gnpDirected(1000, 0.5)
	pathDirected_10        = pathDirected(10)
	pathDirected_1000      = pathDirected(1000)
	pathDirected_100000    = pathDirected(100_000)
)

func gnpDirected(n int, p float64) graph.Directed {
	g := simple.NewDirectedGraph()
	err := gen.Gnp(g, n, p, nil)
	if err != nil {
		panic(fmt.Sprintf("topo: bad test: %v", err))
	}
	return g
}

func pathDirected(n int) graph.Directed {
	g := simple.NewDirectedGraph()
	var idSet gen.IDSet
	for i := 0; i < n; i++ {
		idSet = append(idSet, int64(i))
	}
	gen.Path(g, idSet)
	return g
}

func benchmarkTarjanSCC(b *testing.B, g graph.Directed) {
	var sccs [][]graph.Node
	for i := 0; i < b.N; i++ {
		sccs = TarjanSCC(g)
	}
	if len(sccs) == 0 {
		b.Fatal("unexpected number zero-sized SCC set")
	}
}

func BenchmarkTarjanSCCGnp_10_tenth(b *testing.B) {
	benchmarkTarjanSCC(b, gnpDirected_10_tenth)
}
func BenchmarkTarjanSCCGnp_100_tenth(b *testing.B) {
	benchmarkTarjanSCC(b, gnpDirected_100_tenth)
}
func BenchmarkTarjanSCCGnp_1000_tenth(b *testing.B) {
	benchmarkTarjanSCC(b, gnpDirected_1000_tenth)
}
func BenchmarkTarjanSCCGnp_10_half(b *testing.B) {
	benchmarkTarjanSCC(b, gnpDirected_10_half)
}
func BenchmarkTarjanSCCGnp_100_half(b *testing.B) {
	benchmarkTarjanSCC(b, gnpDirected_100_half)
}
func BenchmarkTarjanSCCGnp_1000_half(b *testing.B) {
	benchmarkTarjanSCC(b, gnpDirected_1000_half)
}

func benchmarkSort(b *testing.B, g graph.Directed) {
	for i := 0; i < b.N; i++ {
		_, _ = Sort(g)
	}
}

func BenchmarkSortGnp_10_tenth(b *testing.B) {
	benchmarkSort(b, gnpDirected_10_tenth)
}
func BenchmarkSortGnp_100_tenth(b *testing.B) {
	benchmarkSort(b, gnpDirected_100_tenth)
}
func BenchmarkSortGnp_1000_tenth(b *testing.B) {
	benchmarkSort(b, gnpDirected_1000_tenth)
}
func BenchmarkSortGnp_10_half(b *testing.B) {
	benchmarkSort(b, gnpDirected_10_half)
}
func BenchmarkSortGnp_100_half(b *testing.B) {
	benchmarkSort(b, gnpDirected_100_half)
}
func BenchmarkSortGnp_1000_half(b *testing.B) {
	benchmarkSort(b, gnpDirected_1000_half)
}
func BenchmarkSortPath_10(b *testing.B) {
	benchmarkSort(b, pathDirected_10)
}
func BenchmarkSortPath_1000(b *testing.B) {
	benchmarkSort(b, pathDirected_1000)
}
func BenchmarkSortPath_100000(b *testing.B) {
	benchmarkSort(b, pathDirected_100000)
}

func benchmarkSortStabilized(b *testing.B, g graph.Directed) {
	for i := 0; i < b.N; i++ {
		_, _ = SortStabilized(g, nil)
	}
}

func BenchmarkSortStabilizedGnp_10_tenth(b *testing.B) {
	benchmarkSortStabilized(b, gnpDirected_10_tenth)
}
func BenchmarkSortStabilizedGnp_100_tenth(b *testing.B) {
	benchmarkSortStabilized(b, gnpDirected_100_tenth)
}
func BenchmarkSortStabilizedGnp_1000_tenth(b *testing.B) {
	benchmarkSortStabilized(b, gnpDirected_1000_tenth)
}
func BenchmarkSortStabilizedGnp_10_half(b *testing.B) {
	benchmarkSortStabilized(b, gnpDirected_10_half)
}
func BenchmarkSortStabilizedGnp_100_half(b *testing.B) {
	benchmarkSortStabilized(b, gnpDirected_100_half)
}
func BenchmarkSortStabilizedGnp_1000_half(b *testing.B) {
	benchmarkSortStabilized(b, gnpDirected_1000_half)
}
func BenchmarkSortStabilizedPath_10(b *testing.B) {
	benchmarkSortStabilized(b, pathDirected_10)
}
func BenchmarkSortStabilizedPath_1000(b *testing.B) {
	benchmarkSortStabilized(b, pathDirected_1000)
}
func BenchmarkSortStabilizedPath_100000(b *testing.B) {
	benchmarkSortStabilized(b, pathDirected_100000)
}
