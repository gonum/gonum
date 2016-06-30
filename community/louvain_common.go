// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package community provides graph community detection functions.
package community

import (
	"fmt"

	"github.com/gonum/graph"
)

// Q returns the modularity Q score of the graph g subdivided into the
// given communities at the given resolution. If communities is nil, the
// unclustered modularity score is returned. The resolution parameter
// is γ as defined in Reichardt and Bornholdt doi:10.1103/PhysRevE.74.016110.
// Q will panic if g has any edge with negative edge weight.
//
// If g is undirected, Q is calculated according to
//  Q = 1/2m \sum_{ij} [ A_{ij} - (\gamma k_i k_j)/2m ] \delta(c_i,c_j),
// If g is directed, it is calculated according to
//  Q = 1/m \sum_{ij} [ A_{ij} - (\gamma k_i^in k_j^out)/m ] \delta(c_i,c_j).
//
// graph.Undirect may be used as a shim to allow calculation of Q for
// directed graphs with the undirected modularity function.
func Q(g graph.Graph, communities [][]graph.Node, resolution float64) float64 {
	switch g := g.(type) {
	case graph.Undirected:
		return qUndirected(g, communities, resolution)
	case graph.Directed:
		return qDirected(g, communities, resolution)
	default:
		panic(fmt.Sprintf("community: invalid graph type: %T", g))
	}
}

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
