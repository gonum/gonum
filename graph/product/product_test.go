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

func vertical() *simple.UndirectedGraph {
	edges := []simple.Edge{
		{F: simple.Node(-1), T: simple.Node(-2)},
	}
	g := simple.NewUndirectedGraph()
	for _, e := range edges {
		g.SetEdge(e)
	}
	return g
}

func horizontal() *simple.UndirectedGraph {
	edges := []simple.Edge{
		{F: simple.Node(1), T: simple.Node(2)},
	}
	g := simple.NewUndirectedGraph()
	for _, e := range edges {
		g.SetEdge(e)
	}
	return g
}

var productTests = []struct {
	a, b *simple.UndirectedGraph
}{
	{a: vertical(), b: horizontal()},
	{a: left(), b: right()},
}

func TestCartesian(t *testing.T) {
	for _, test := range productTests {
		got := simple.NewUndirectedGraph()
		Cartesian(got, test.a, test.b)
		gotBytes, _ := dot.Marshal(got, "", "", "  ")

		want := simple.NewUndirectedGraph()
		naiveCartesian(want, test.a, test.b)
		wantBytes, _ := dot.Marshal(want, "", "", "  ")

		gotEdgesLen := got.Edges().Len()
		nA := test.a.Nodes().Len()
		mA := test.a.Edges().Len()
		nB := test.b.Nodes().Len()
		mB := test.b.Edges().Len()
		wantEdgesLen := mB*nA + mA*nB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Cartesian product: got:%d want:%d", gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Cartesian product result: got:\n%s\nwant:\n%s", gotBytes, wantBytes)
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

		gotEdgesLen := got.Edges().Len()
		mA := test.a.Edges().Len()
		mB := test.b.Edges().Len()
		wantEdgesLen := 2 * mA * mB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Tensor product: got:%d want:%d", gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Tensor product result: got:\n%s\nwant:\n%s", gotBytes, wantBytes)
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

		gotEdgesLen := got.Edges().Len()
		nA := test.a.Nodes().Len()
		mA := test.a.Edges().Len()
		nB := test.b.Nodes().Len()
		mB := test.b.Edges().Len()
		wantEdgesLen := mB*nA + mA*nB*nB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Lexicographical product: got:%d want:%d", gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Lexicographical product result: got:\n%s\nwant:\n%s", gotBytes, wantBytes)
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

		gotEdgesLen := got.Edges().Len()
		nA := test.a.Nodes().Len()
		mA := test.a.Edges().Len()
		nB := test.b.Nodes().Len()
		mB := test.b.Edges().Len()
		wantEdgesLen := nA*mB + nB*mA + 2*mA*mB
		if gotEdgesLen != wantEdgesLen {
			t.Errorf("unexpected number of edges for Strong product: got:%d want:%d", gotEdgesLen, wantEdgesLen)
		}

		if !bytes.Equal(gotBytes, wantBytes) {
			t.Errorf("unexpected Strong product result: got:\n%s\nwant:\n%s", gotBytes, wantBytes)
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
			t.Errorf("unexpected Co-normal product result: got:\n%s\nwant:\n%s", gotBytes, wantBytes)
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
			t.Errorf("unexpected Modular product result: got:\n%s\nwant:\n%s", gotBytes, wantBytes)
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

func BenchmarkModular(b *testing.B) {
	g1 := gnp(50, 0.5, nil)
	g2 := gnp(50, 0.5, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dst := simple.NewDirectedGraph()
		Modular(dst, g1, g2)
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
