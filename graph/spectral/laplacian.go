// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spectral

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/mat"
)

// Laplacian is a graph Laplacian matrix.
type Laplacian struct {
	// Matrix holds the Laplacian matrix.
	mat.Matrix

	// Nodes holds the input graph nodes.
	Nodes []graph.Node

	// Index is a mapping from the graph
	// node IDs to row and column indices.
	Index map[int64]int
}

// NewLaplacian returns a Laplacian matrix for the simple undirected graph g.
// The Laplacian is defined as D-A where D is a diagonal matrix holding the
// degree of each node and A is the graph adjacency matrix of the input graph.
// If g contains self edges, NewLaplacian will panic.
func NewLaplacian(g graph.Undirected) Laplacian {
	nodes := graph.NodesOf(g.Nodes())
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
	}

	l := mat.NewSymDense(len(nodes), nil)
	for j, u := range nodes {
		uid := u.ID()
		to := graph.NodesOf(g.From(uid))
		l.SetSym(j, j, float64(len(to)))
		for _, v := range to {
			vid := v.ID()
			if uid == vid {
				panic("network: self edge in graph")
			}
			if uid < vid {
				l.SetSym(indexOf[vid], j, -1)
			}
		}
	}

	return Laplacian{Matrix: l, Nodes: nodes, Index: indexOf}
}

// NewSymNormLaplacian returns a symmetric normalized Laplacian matrix for the
// simple undirected graph g.
// The normalized Laplacian is defined as I-D^(-1/2)AD^(-1/2) where D is a
// diagonal matrix holding the degree of each node and A is the graph adjacency
// matrix of the input graph.
// If g contains self edges, NewSymNormLaplacian will panic.
func NewSymNormLaplacian(g graph.Undirected) Laplacian {
	nodes := graph.NodesOf(g.Nodes())
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
	}

	l := mat.NewSymDense(len(nodes), nil)
	for j, u := range nodes {
		uid := u.ID()
		to := graph.NodesOf(g.From(uid))
		if len(to) == 0 {
			continue
		}
		l.SetSym(j, j, 1)
		squdeg := math.Sqrt(float64(len(to)))
		for _, v := range to {
			vid := v.ID()
			if uid == vid {
				panic("network: self edge in graph")
			}
			if uid < vid {
				to := g.From(vid)
				k := to.Len()
				if k < 0 {
					k = len(graph.NodesOf(to))
				}
				l.SetSym(indexOf[vid], j, -1/(squdeg*math.Sqrt(float64(k))))
			}
		}
	}

	return Laplacian{Matrix: l, Nodes: nodes, Index: indexOf}
}

// NewRandomWalkLaplacian returns a damp-scaled random walk Laplacian matrix for
// the simple graph g.
// The random walk Laplacian is defined as I-D^(-1)A where D is a diagonal matrix
// holding the degree of each node and A is the graph adjacency matrix of the input
// graph.
// If g contains self edges, NewRandomWalkLaplacian will panic.
func NewRandomWalkLaplacian(g graph.Graph, damp float64) Laplacian {
	nodes := graph.NodesOf(g.Nodes())
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
	}

	l := mat.NewDense(len(nodes), len(nodes), nil)
	for j, u := range nodes {
		uid := u.ID()
		to := graph.NodesOf(g.From(uid))
		if len(to) == 0 {
			continue
		}
		l.Set(j, j, 1-damp)
		rudeg := (damp - 1) / float64(len(to))
		for _, v := range to {
			vid := v.ID()
			if uid == vid {
				panic("network: self edge in graph")
			}
			l.Set(indexOf[vid], j, rudeg)
		}
	}

	return Laplacian{Matrix: l, Nodes: nodes, Index: indexOf}
}
