// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package community provides graph community detection functions.
package community

import "github.com/gonum/graph"

// community is a reduced graph node describing its membership.
type community struct {
	id int

	nodes  []graph.Node
	weight float64
}

func (n community) ID() int { return n.id }

// edge is a reduced graph edge.
type edge struct {
	from, to community
	weight   float64
}

func (e edge) From() graph.Node { return e.from }
func (e edge) To() graph.Node   { return e.to }
func (e edge) Weight() float64  { return e.weight }

// commIdx is an index of a node in a community held by a localMover.
type commIdx struct {
	community int
	node      int
}

// node is defined to avoid an import of .../graph/simple.
type node int

func (n node) ID() int { return int(n) }

const negativeWeight = "community: negative edge weight"

// weightFuncFor returns a constructed weight function for g.
func weightFuncFor(g graph.Graph) func(x, y graph.Node) float64 {
	if wg, ok := g.(graph.Weighter); ok {
		return func(x, y graph.Node) float64 {
			w, ok := wg.Weight(x, y)
			if !ok {
				return 0
			}
			if w < 0 {
				panic(negativeWeight)
			}
			return w
		}
	}
	return func(x, y graph.Node) float64 {
		e := g.Edge(x, y)
		if e == nil {
			return 0
		}
		w := e.Weight()
		if w < 0 {
			panic(negativeWeight)
		}
		return w
	}
}
