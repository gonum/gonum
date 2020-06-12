// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package product

import (
	"bytes"
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

func (n Node) DOTID() string { return fmt.Sprintf("(%d,%d)", n.A.ID(), n.B.ID()) }

func left() *simple.UndirectedGraph {
	edges := []simple.Edge{
		{F: simple.Node(-1), T: simple.Node(-2)},
		{F: simple.Node(-2), T: simple.Node(-3)},
		{F: simple.Node(-2), T: simple.Node(-4)},
		{F: simple.Node(-3), T: simple.Node(-5)},
		{F: simple.Node(-4), T: simple.Node(-5)},
	}
	g := simple.NewUndirectedGraph()
	for _, e := range edges {
		g.SetEdge(e)
	}
	return g
}

func right() *simple.UndirectedGraph {
	edges := []simple.Edge{
		{F: simple.Node(1), T: simple.Node(2)},
		{F: simple.Node(2), T: simple.Node(3)},
		{F: simple.Node(2), T: simple.Node(4)},
	}
	g := simple.NewUndirectedGraph()
	for _, e := range edges {
		g.SetEdge(e)
	}
	return g
}

func path(m int) *simple.UndirectedGraph {
	sign := 1
	if m < 0 {
		sign = -1
		m = -m
	}
	g := simple.NewUndirectedGraph()
	if m == 0 {
		g.AddNode(simple.Node(0))
	}
	for i := 1; i <= m; i++ {
		g.SetEdge(simple.Edge{F: simple.Node(sign * i), T: simple.Node(sign * (i + 1))})
	}
	return g
}

var productTests = []struct {
	name string
	a, b *simple.UndirectedGraph
}{
	{name: "paths", a: path(-1), b: path(1)},
	{name: "wp_mp", a: path(-2), b: path(2)},
	{name: "wp_gp", a: left(), b: right()},
	{name: "gnp_2×2", a: gnp(2, 0.5, rand.NewSource(1)), b: gnp(2, 0.5, rand.NewSource(2))},
	{name: "gnp_2×3", a: gnp(2, 0.5, rand.NewSource(1)), b: gnp(3, 0.5, rand.NewSource(2))},
	{name: "gnp_3×3", a: gnp(3, 0.5, rand.NewSource(1)), b: gnp(3, 0.5, rand.NewSource(2))},
	{name: "gnp_4×4", a: gnp(4, 0.5, rand.NewSource(1)), b: gnp(4, 0.5, rand.NewSource(2))},
}

func TestCartesian(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		Cartesian(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveCartesian(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		gotEdgesLen := len(graph.EdgesOf(got.Edges()))
		nA := len(graph.NodesOf(test.a.Nodes()))
		mA := len(graph.EdgesOf(test.a.Edges()))
		nB := len(graph.NodesOf(test.b.Nodes()))
		mB := len(graph.EdgesOf(test.b.Edges()))
		wantEdgesLen := mB*nA + mA*nB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Cartesian product of %s: got:%d want:%d",
				test.name, gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Cartesian product result for %s:\ngot:\n%s\nwant:\n%s",
				test.name, gotBytes, wantBytes)
		}
	}
}

// (u₁=v₁ and u₂~v₂) or (u₁~v₁ and u₂=v₂).
func naiveCartesian(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)

	for _, p := range product {
		dst.AddNode(p)
	}

	for _, u := range product {
		for _, v := range product {
			edgeInA := a.Edge(u.A.ID(), v.A.ID()) != nil
			edgeInB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if (u.A.ID() == v.A.ID() && edgeInB) || (edgeInA && u.B.ID() == v.B.ID()) {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

func TestTensor(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		Tensor(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveTensor(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		gotEdgesLen := len(graph.EdgesOf(got.Edges()))
		mA := len(graph.EdgesOf(test.a.Edges()))
		mB := len(graph.EdgesOf(test.b.Edges()))
		wantEdgesLen := 2 * mA * mB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Tensor product of %s: got:%d want:%d",
				test.name, gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Tensor product result for %s:\ngot:\n%s\nwant:\n%s",
				test.name, gotBytes, wantBytes)
		}
	}
}

// u₁~v₁ and u₂~v₂.
func naiveTensor(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)

	for _, p := range product {
		dst.AddNode(p)
	}

	for _, u := range product {
		for _, v := range product {
			edgeInA := a.Edge(u.A.ID(), v.A.ID()) != nil
			edgeInB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if edgeInA && edgeInB {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

func TestLexicographical(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		Lexicographical(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveLexicographical(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		gotEdgesLen := len(graph.EdgesOf(got.Edges()))
		nA := len(graph.NodesOf(test.a.Nodes()))
		mA := len(graph.EdgesOf(test.a.Edges()))
		nB := len(graph.NodesOf(test.b.Nodes()))
		mB := len(graph.EdgesOf(test.b.Edges()))
		wantEdgesLen := mB*nA + mA*nB*nB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Lexicographical product of %s: got:%d want:%d",
				test.name, gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Lexicographical product result for %s:\ngot:\n%s\nwant:\n%s",
				test.name, gotBytes, wantBytes)
		}
	}
}

// u₁~v₁ or (u₁=v₁ and u₂~v₂).
func naiveLexicographical(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)

	for _, p := range product {
		dst.AddNode(p)
	}

	for _, u := range product {
		for _, v := range product {
			edgeInA := a.Edge(u.A.ID(), v.A.ID()) != nil
			edgeInB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if edgeInA || (u.A.ID() == v.A.ID() && edgeInB) {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

func TestStrong(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		Strong(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveStrong(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		gotEdgesLen := len(graph.EdgesOf(got.Edges()))
		nA := len(graph.NodesOf(test.a.Nodes()))
		mA := len(graph.EdgesOf(test.a.Edges()))
		nB := len(graph.NodesOf(test.b.Nodes()))
		mB := len(graph.EdgesOf(test.b.Edges()))
		wantEdgesLen := nA*mB + nB*mA + 2*mA*mB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Strong product of %s: got:%d want:%d",
				test.name, gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Strong product result for %s:\ngot:\n%s\nwant:\n%s",
				test.name, gotBytes, wantBytes)
		}
	}
}

// (u₁=v₁ and u₂~v₂) or (u₁~v₁ and u₂=v₂) or (u₁~v₁ and u₂~v₂).
func naiveStrong(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)

	for _, p := range product {
		dst.AddNode(p)
	}

	for _, u := range product {
		for _, v := range product {
			edgeInA := a.Edge(u.A.ID(), v.A.ID()) != nil
			edgeInB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if (u.A.ID() == v.A.ID() && edgeInB) || (edgeInA && u.B.ID() == v.B.ID()) || (edgeInA && edgeInB) {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

func TestCoNormal(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		CoNormal(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveCoNormal(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Co-normal product result for %s:\ngot:\n%s\nwant:\n%s",
				test.name, gotBytes, wantBytes)
		}
	}
}

// u₁~v₁ or u₂~v₂.
func naiveCoNormal(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)

	for _, p := range product {
		dst.AddNode(p)
	}

	for _, u := range product {
		for _, v := range product {
			edgeInA := a.Edge(u.A.ID(), v.A.ID()) != nil
			edgeInB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if edgeInA || edgeInB {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

func TestModular(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		Modular(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveModular(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Modular product result for %s:\ngot:\n%s\nwant:\n%s", test.name, gotBytes, wantBytes)
		}
	}
}

func TestModularExt(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		ModularExt(got, test.a, test.b, func(a, b graph.Edge) bool { return true })
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveModular(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected ModularExt product result for %s:\ngot:\n%s\nwant:\n%s", test.name, gotBytes, wantBytes)
		}
	}
}

// (u₁~v₁ and u₂~v₂) or (u₁≁v₁ and u₂≁v₂).
func naiveModular(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)

	for _, p := range product {
		dst.AddNode(p)
	}

	for i, u := range product {
		for j, v := range product {
			if i == j || u.A.ID() == v.A.ID() || u.B.ID() == v.B.ID() {
				// No self-edges.
				continue
			}
			edgeInA := a.Edge(u.A.ID(), v.A.ID()) != nil
			edgeInB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if (edgeInA && edgeInB) || (!edgeInA && !edgeInB) {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

func BenchmarkProduct(b *testing.B) {
	for seed, bench := range []struct {
		name    string
		product func(dst graph.Builder, a, b graph.Graph)
		len     []int
	}{
		{"Cartesian", Cartesian, []int{50, 100}},
		{"Cartesian naive", naiveCartesian, []int{50, 100}},
		{"CoNormal", CoNormal, []int{50}},
		{"CoNormal naive", naiveCoNormal, []int{50}},
		{"Lexicographical", Lexicographical, []int{50}},
		{"Lexicographical naive", naiveLexicographical, []int{50}},
		{"Modular", Modular, []int{50}},
		{"Modular naive", naiveModular, []int{50}},
		{"Strong", Strong, []int{50}},
		{"Strong naive", naiveStrong, []int{50}},
		{"Tensor", Tensor, []int{50}},
		{"Tensor naive", naiveTensor, []int{50}},
	} {
		for _, p := range []float64{0.05, 0.25, 0.5, 0.75, 0.95} {
			for _, n := range bench.len {
				src := rand.NewSource(uint64(seed))
				b.Run(fmt.Sprintf("%s %d-%.2f", bench.name, n, p), func(b *testing.B) {
					g1 := gnp(n, p, src)
					g2 := gnp(n, p, src)
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						dst := simple.NewDirectedGraph()
						bench.product(dst, g1, g2)
					}
				})
			}
		}
	}
}

func gnp(n int, p float64, src rand.Source) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()
	err := gen.Gnp(g, n, p, src)
	if err != nil {
		panic(fmt.Sprintf("gnp: bad test: %v", err))
	}
	return g
}
