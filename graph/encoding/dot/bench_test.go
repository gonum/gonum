// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dot

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	gnpDirected_10_tenth   = gnpDirected(10, 0.1)
	gnpDirected_100_tenth  = gnpDirected(100, 0.1)
	gnpDirected_1000_tenth = gnpDirected(1000, 0.1)
	gnpDirected_10_half    = gnpDirected(10, 0.5)
	gnpDirected_100_half   = gnpDirected(100, 0.5)
	gnpDirected_1000_half  = gnpDirected(1000, 0.5)

	powerLawMultiDirected_10_tenth   = powerLawMultiDirected(10, 1)
	powerLawMultiDirected_100_tenth  = powerLawMultiDirected(100, 10)
	powerLawMultiDirected_1000_tenth = powerLawMultiDirected(1000, 100)
	powerLawMultiDirected_10_half    = powerLawMultiDirected(10, 5)
	powerLawMultiDirected_100_half   = powerLawMultiDirected(100, 50)
	powerLawMultiDirected_1000_half  = powerLawMultiDirected(1000, 500)
)

func gnpDirected(n int, p float64) graph.Directed {
	g := simple.NewDirectedGraph()
	err := gen.Gnp(g, n, p, rand.NewPCG(1, 1))
	if err != nil {
		panic(fmt.Sprintf("dot: bad test: %v", err))
	}
	return g
}

func powerLawMultiDirected(n, d int) graph.DirectedMultigraph {
	g := multi.NewDirectedGraph()
	err := gen.PowerLaw(g, n, d, rand.NewPCG(1, 1))
	if err != nil {
		panic(fmt.Sprintf("dot: bad test: %v", err))
	}
	return g
}

func benchmarkUnmarshal(b *testing.B, g graph.Directed) {
	marshalled, err := Marshal(g, "g", "", "")
	if err != nil {
		b.Fatalf("dot: bad Marshal input: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(marshalled, simple.NewDirectedGraph()); err != nil {
			b.Fatalf("dot: bad Unmarshal input: %v", err)
		}
	}
}

func BenchmarkUnmarshalGnp_10_tenth(b *testing.B) {
	benchmarkUnmarshal(b, gnpDirected_10_tenth)
}
func BenchmarkUnmarshalGnp_100_tenth(b *testing.B) {
	benchmarkUnmarshal(b, gnpDirected_100_tenth)
}
func BenchmarkUnmarshalGnp_1000_tenth(b *testing.B) {
	benchmarkUnmarshal(b, gnpDirected_1000_tenth)
}
func BenchmarkUnmarshalGnp_10_half(b *testing.B) {
	benchmarkUnmarshal(b, gnpDirected_10_half)
}
func BenchmarkUnmarshalGnp_100_half(b *testing.B) {
	benchmarkUnmarshal(b, gnpDirected_100_half)
}
func BenchmarkUnmarshalGnp_1000_half(b *testing.B) {
	benchmarkUnmarshal(b, gnpDirected_1000_half)
}

func benchmarkUnmarshalMulti(b *testing.B, g graph.DirectedMultigraph) {
	marshalled, err := MarshalMulti(g, "g", "", "")
	if err != nil {
		b.Fatalf("dot: bad Marshal input: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := UnmarshalMulti(marshalled, multi.NewDirectedGraph()); err != nil {
			b.Fatalf("dot: bad Unmarshal input: %v", err)
		}
	}
}

func BenchmarkUnmarshalMultiPowerLaw_10_tenth(b *testing.B) {
	benchmarkUnmarshalMulti(b, powerLawMultiDirected_10_tenth)
}
func BenchmarkUnmarshalMultiPowerLaw_100_tenth(b *testing.B) {
	benchmarkUnmarshalMulti(b, powerLawMultiDirected_100_tenth)
}
func BenchmarkUnmarshalMultiPowerLaw_1000_tenth(b *testing.B) {
	benchmarkUnmarshalMulti(b, powerLawMultiDirected_1000_tenth)
}
func BenchmarkUnmarshalMultiPowerLaw_10_half(b *testing.B) {
	benchmarkUnmarshalMulti(b, powerLawMultiDirected_10_half)
}
func BenchmarkUnmarshalMultiPowerLaw_100_half(b *testing.B) {
	benchmarkUnmarshalMulti(b, powerLawMultiDirected_100_half)
}
func BenchmarkUnmarshalMultiPowerLaw_1000_half(b *testing.B) {
	benchmarkUnmarshalMulti(b, powerLawMultiDirected_1000_half)
}
