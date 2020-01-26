// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package product

import (
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/stat/combin"
)

// Node is a product of two graph nodes.
// All graph products return this type directly via relevant graph.Graph method
// call, or indirectly via calls to graph.Edge methods from returned edges.
type Node struct {
	UID int64

	// A and B hold the nodes from graph a and b in a
	// graph product constructed by a product function.
	A, B graph.Node
}

// ID implements the graph.Node interface.
func (n Node) ID() int64 { return n.UID }

// Cartesian constructs the Cartesian product of a and b in dst.
//
// The Cartesian product of G₁ and G₂, G₁□G₂ has edges (u₁, u₂)~(v₁, v₂) when
// (u₁=v₁ and u₂~v₂) or (u₁~v₁ and u₂=v₂). The Cartesian product has size m₂n₁+m₁n₂
// where m is the size of the input graphs and n is their order.
func Cartesian(dst graph.Builder, a, b graph.Graph) {
	aNodes, bNodes, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	indexOfA := indexOf(aNodes)
	indexOfB := indexOf(bNodes)

	for _, p := range product {
		dst.AddNode(p)
	}

	dims := []int{len(aNodes), len(bNodes)}
	for i, uA := range aNodes {
		for j, uB := range bNodes {
			toB := b.From(uB.ID())
			for toB.Next() {
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, j}, dims)],
					product[combin.IdxFor([]int{i, indexOfB[toB.Node().ID()]}, dims)],
				))
			}

			toA := a.From(uA.ID())
			for toA.Next() {
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, j}, dims)],
					product[combin.IdxFor([]int{indexOfA[toA.Node().ID()], j}, dims)],
				))
			}
		}
	}
}

// Tensor constructs the Tensor product of a and b in dst.
//
// The Tensor product of G₁ and G₂, G₁⨯G₂ has edges (u₁, u₂)~(v₁, v₂) when
// u₁~v₁ and u₂~v₂. The Tensor product has size 2m₁m₂ where m is the size
// of the input graphs.
func Tensor(dst graph.Builder, a, b graph.Graph) {
	aNodes, bNodes, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	indexOfA := indexOf(aNodes)
	indexOfB := indexOf(bNodes)

	for _, p := range product {
		dst.AddNode(p)
	}

	dims := []int{len(aNodes), len(bNodes)}
	for i, uA := range aNodes {
		toA := a.From(uA.ID())
		for toA.Next() {
			j := indexOfA[toA.Node().ID()]
			for k, uB := range bNodes {
				toB := b.From(uB.ID())
				for toB.Next() {
					dst.SetEdge(dst.NewEdge(
						product[combin.IdxFor([]int{i, k}, dims)],
						product[combin.IdxFor([]int{j, indexOfB[toB.Node().ID()]}, dims)],
					))
				}
			}
		}
	}
}

// Lexicographical constructs the Lexicographical product of a and b in dst.
//
// The Lexicographical product of G₁ and G₂, G₁·G₂ has edges (u₁, u₂)~(v₁, v₂) when
// u₁~v₁ or (u₁=v₁ and u₂~v₂). The Lexicographical product has size m₂n₁+m₁n₂²
// where m is the size of the input graphs and n is their order.
func Lexicographical(dst graph.Builder, a, b graph.Graph) {
	aNodes, bNodes, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	indexOfA := indexOf(aNodes)
	indexOfB := indexOf(bNodes)

	for _, p := range product {
		dst.AddNode(p)
	}

	dims := []int{len(aNodes), len(bNodes)}
	p := make([]int, 2)
	for i, uA := range aNodes {
		toA := a.From(uA.ID())
		for toA.Next() {
			j := indexOfA[toA.Node().ID()]
			gen := combin.NewCartesianGenerator([]int{len(bNodes), len(bNodes)})
			for gen.Next() {
				p = gen.Product(p)
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, p[0]}, dims)],
					product[combin.IdxFor([]int{j, p[1]}, dims)],
				))
			}
		}

		for j, uB := range bNodes {
			toB := b.From(uB.ID())
			for toB.Next() {
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, j}, dims)],
					product[combin.IdxFor([]int{i, indexOfB[toB.Node().ID()]}, dims)],
				))
			}
		}
	}
}

// Strong constructs the Strong product of a and b in dst.
//
// The Strong product of G₁ and G₂, G₁⊠G₂ has edges (u₁, u₂)~(v₁, v₂) when
// (u₁=v₁ and u₂~v₂) or (u₁~v₁ and u₂=v₂) or (u₁~v₁ and u₂~v₂). The Strong
// product has size n₁m₂+n₂m₁+2m₁m₂ where m is the size of the input graphs
// and n is their order.
func Strong(dst graph.Builder, a, b graph.Graph) {
	aNodes, bNodes, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	indexOfA := indexOf(aNodes)
	indexOfB := indexOf(bNodes)

	for _, p := range product {
		dst.AddNode(p)
	}

	dims := []int{len(aNodes), len(bNodes)}
	for i, uA := range aNodes {
		for j, uB := range bNodes {
			toB := b.From(uB.ID())
			for toB.Next() {
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, j}, dims)],
					product[combin.IdxFor([]int{i, indexOfB[toB.Node().ID()]}, dims)],
				))
			}

			toA := a.From(uA.ID())
			for toA.Next() {
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, j}, dims)],
					product[combin.IdxFor([]int{indexOfA[toA.Node().ID()], j}, dims)],
				))
			}
		}

		toA := a.From(uA.ID())
		for toA.Next() {
			for j, uB := range bNodes {
				toB := b.From(uB.ID())
				for toB.Next() {
					dst.SetEdge(dst.NewEdge(
						product[combin.IdxFor([]int{i, j}, dims)],
						product[combin.IdxFor([]int{indexOfA[toA.Node().ID()], indexOfB[toB.Node().ID()]}, dims)],
					))
				}
			}
		}
	}
}

// CoNormal constructs the Co-normal product of a and b in dst.
//
// The Co-normal product of G₁ and G₂, G₁*G₂ (or G₁[G₂]) has edges (u₁, u₂)~(v₁, v₂)
// when u₁~v₁ or u₂~v₂. The Co-normal product is non-commutative.
func CoNormal(dst graph.Builder, a, b graph.Graph) {
	aNodes, bNodes, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	indexOfA := indexOf(aNodes)
	indexOfB := indexOf(bNodes)

	for _, p := range product {
		dst.AddNode(p)
	}

	dims := []int{len(aNodes), len(bNodes)}
	p := make([]int, 2)
	for i, u := range aNodes {
		to := a.From(u.ID())
		for to.Next() {
			j := indexOfA[to.Node().ID()]
			gen := combin.NewCartesianGenerator([]int{len(bNodes), len(bNodes)})
			for gen.Next() {
				p = gen.Product(p)
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{i, p[0]}, dims)],
					product[combin.IdxFor([]int{j, p[1]}, dims)],
				))
			}
		}
	}
	for i, u := range bNodes {
		to := b.From(u.ID())
		for to.Next() {
			j := indexOfB[to.Node().ID()]
			gen := combin.NewCartesianGenerator([]int{len(aNodes), len(aNodes)})
			for gen.Next() {
				p = gen.Product(p)
				dst.SetEdge(dst.NewEdge(
					product[combin.IdxFor([]int{p[0], i}, dims)],
					product[combin.IdxFor([]int{p[1], j}, dims)],
				))
			}
		}
	}
}

// Modular constructs the Modular product of a and b in dst.
//
// The Modular product of G₁ and G₂, G₁◊G₂ has edges (u₁, u₂)~(v₁, v₂)
// when (u₁~v₁ and u₂~v₂) or (u₁≁v₁ and u₂≁v₂), and (u₁≠v₁ and u₂≠v₂).
//
// Modular is O(n^2) time where n is the order of the Cartesian product
// of a and b.
func Modular(dst graph.Builder, a, b graph.Graph) {
	_, _, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	for _, p := range product {
		dst.AddNode(p)
	}

	_, aUndirected := a.(graph.Undirected)
	_, bUndirected := b.(graph.Undirected)
	_, dstUndirected := dst.(graph.Undirected)
	undirected := aUndirected && bUndirected && dstUndirected

	n := len(product)
	if undirected {
		n--
	}
	for i, u := range product[:n] {
		var m int
		if undirected {
			m = i + 1
		}
		for _, v := range product[m:] {
			if u.A.ID() == v.A.ID() || u.B.ID() == v.B.ID() {
				// No self-loops.
				continue
			}
			inA := a.Edge(u.A.ID(), v.A.ID()) != nil
			inB := b.Edge(u.B.ID(), v.B.ID()) != nil
			if inA == inB {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

// ModularExt constructs the Modular product of a and b in dst with
// additional control over assessing edge agreement.
//
// In addition to the modular product conditions, agree(u₁v₁, u₂v₂) must
// return true when (u₁~v₁ and u₂~v₂) for an edge to be added between
// (u₁, u₂) and (v₁, v₂) in dst. If agree is nil, Modular is called.
//
// ModularExt is O(n^2) time where n is the order of the Cartesian product
// of a and b.
func ModularExt(dst graph.Builder, a, b graph.Graph, agree func(eA, eB graph.Edge) bool) {
	if agree == nil {
		Modular(dst, a, b)
		return
	}

	_, _, product := cartesianNodes(a, b)
	if len(product) == 0 {
		return
	}

	for _, p := range product {
		dst.AddNode(p)
	}

	_, aUndirected := a.(graph.Undirected)
	_, bUndirected := b.(graph.Undirected)
	_, dstUndirected := dst.(graph.Undirected)
	undirected := aUndirected && bUndirected && dstUndirected

	n := len(product)
	if undirected {
		n--
	}
	for i, u := range product[:n] {
		var m int
		if undirected {
			m = i + 1
		}
		for _, v := range product[m:] {
			if u.A.ID() == v.A.ID() || u.B.ID() == v.B.ID() {
				// No self-loops.
				continue
			}
			eA := a.Edge(u.A.ID(), v.A.ID())
			eB := b.Edge(u.B.ID(), v.B.ID())
			inA := eA != nil
			inB := eB != nil
			if (inA && inB && agree(eA, eB)) || (!inA && !inB) {
				dst.SetEdge(dst.NewEdge(u, v))
			}
		}
	}
}

// cartesianNodes returns the Cartesian product of the nodes in a and b.
func cartesianNodes(a, b graph.Graph) (aNodes, bNodes []graph.Node, product []Node) {
	aNodes = lexicalNodes(a)
	bNodes = lexicalNodes(b)

	lens := []int{len(aNodes), len(bNodes)}
	product = make([]Node, combin.Card(lens))
	gen := combin.NewCartesianGenerator(lens)
	p := make([]int, 2)
	for id := int64(0); gen.Next(); id++ {
		p = gen.Product(p)
		product[id] = Node{UID: id, A: aNodes[p[0]], B: bNodes[p[1]]}
	}
	return aNodes, bNodes, product
}

// lexicalNodes returns the nodes in g sorted lexically by node ID.
func lexicalNodes(g graph.Graph) []graph.Node {
	nodes := graph.NodesOf(g.Nodes())
	sort.Sort(ordered.ByID(nodes))
	return nodes
}

func indexOf(nodes []graph.Node) map[int64]int {
	idx := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		idx[n.ID()] = i
	}
	return idx
}
