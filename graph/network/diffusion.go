// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/mat"
)

// Diffuse performs a heat diffusion across nodes of the undirected
// graph described by the given Laplacian using the initial heat distribution,
// h, according to the Laplacian and diffusing for time t.
// The resulting heat distribution is returned, written into the map h,
//  d = exp(-Lt)×h
// where L is the graph Laplacian. Indexing into h is defined by the
// Laplacian Index field.
//
// Nodes without corresponding entries in h are given an initial heat of zero,
// and entries in h without a corresponding node in the original graph are
// not altered.
func Diffuse(h map[int64]float64, by Laplacian, t float64) map[int64]float64 {
	heat := make([]float64, len(by.Index))
	for id, i := range by.Index {
		heat[i] = h[id]
	}
	v := mat.NewVecDense(len(heat), heat)

	var m, tl mat.Dense
	tl.Scale(-t, by)
	m.Exp(&tl)
	v.MulVec(&m, v)

	for i, n := range v.RawVector().Data {
		h[by.Nodes[i].ID()] = n
	}
	return h
}

// DiffuseToEquilibrium performs a heat diffusion across nodes of the
// graph described by the given Laplacian using the initial heat
// distribution, h, according to the Laplacian until the update function
//  h_{n+1} = h_n - L×h_n
// results in a 2-norm update difference within tol, or iter updates have
// been made.
// The resulting heat distribution is returned as eq, written into the map h,
// and a boolean indicating whether the equilibrium converged.
// Indexing into h is defined by the Laplacian Index field.
//
// Nodes without corresponding entries in h are given an initial heat of zero,
// and entries in h without a corresponding node in the original graph are
// not altered.
func DiffuseToEquilibrium(h map[int64]float64, by Laplacian, tol float64, iters int) (eq map[int64]float64, ok bool) {
	heat := make([]float64, len(by.Index))
	for id, i := range by.Index {
		heat[i] = h[id]
	}
	v := mat.NewVecDense(len(heat), heat)

	last := make([]float64, len(by.Index))
	for id, i := range by.Index {
		heat[i] = h[id]
	}
	lastV := mat.NewVecDense(len(last), last)

	var tmp mat.VecDense
	for {
		lastV, v = v, lastV
		tmp.MulVec(by.Matrix, lastV)
		v.SubVec(lastV, &tmp)
		if normDiff(heat, last) < tol {
			ok = true
			break
		}
		iters--
		if iters < 0 {
			break
		}
	}

	for i, n := range v.RawVector().Data {
		h[by.Nodes[i].ID()] = n
	}
	return h, ok
}

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
// The Laplacian is defined by D-A where D is a diagonal matrix holding the
// degree of each node and A is the graph adjacency matrix of the input graph.
func NewLaplacian(g graph.Undirected) Laplacian {
	nodes := g.Nodes()
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
	}

	l := mat.NewSymDense(len(nodes), nil)
	for j, u := range nodes {
		to := g.From(u)
		l.SetSym(j, j, float64(len(to)))
		uid := u.ID()
		for _, v := range to {
			vid := v.ID()
			if uid < vid {
				l.SetSym(indexOf[vid], j, -1)
			}
		}
	}

	return Laplacian{Matrix: l, Nodes: nodes, Index: indexOf}
}

// NewSymNormLaplacian returns a symmetric normalized Laplacian matrix for the
// simple undirected graph g.
// The Laplacian is defined by I-D^(-1/2)AD^(-1/2) where D is a diagonal matrix holding the
// degree of each node and A is the graph adjacency matrix of the input graph.
func NewSymNormLaplacian(g graph.Undirected) Laplacian {
	nodes := g.Nodes()
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
	}

	l := mat.NewSymDense(len(nodes), nil)
	for j, u := range nodes {
		to := g.From(u)
		if len(to) == 0 {
			continue
		}
		l.SetSym(j, j, 1)
		uid := u.ID()
		udeg := math.Sqrt(float64(len(to)))
		for _, v := range to {
			vid := v.ID()
			if uid < vid {
				l.SetSym(indexOf[vid], j, -1/(udeg*math.Sqrt(float64(len(g.From(v))))))
			}
		}
	}

	return Laplacian{Matrix: l, Nodes: nodes, Index: indexOf}
}

// NewRandomWalkLaplacian returns a random walk Laplacian matrix for the
// simple graph g.
// The Laplacian is defined by I-D^(-1)A where D is a diagonal matrix holding the
// degree of each node and A is the graph adjacency matrix of the input graph.
func NewRandomWalkLaplacian(g graph.Graph, damp float64) Laplacian {
	nodes := g.Nodes()
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
	}

	l := mat.NewDense(len(nodes), len(nodes), nil)
	for j, u := range nodes {
		to := g.From(u)
		if len(to) == 0 {
			continue
		}
		l.Set(j, j, 1-damp)
		rudeg := (damp - 1) / float64(len(to))
		for _, v := range to {
			l.Set(indexOf[v.ID()], j, rudeg)
		}
	}

	return Laplacian{Matrix: l, Nodes: nodes, Index: indexOf}
}
