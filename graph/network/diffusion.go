// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/mat"
)

// HeatDiffusion performs a heat diffusion network propagation across nodes
// of the undirected graph g using the initial heat query h provided, and
// diffusing for time t. The resulting heat diffusion is returned, written
// into the map h,
//  d = exp(-Lt)×h
// where L is the graph Laplacian or symmetric normalized Laplacian.
//
// The Laplacian is defined by D-A where D is a diagonal matrix holding the
// degree of each node and A is the graph adjacency matrix of the input graph.
// The symmetric normalized Laplacian is defined by I-D^(-1/2)AD^(-1/2).
//
// Nodes without corresponding entries in h are given an initial heat of zero,
// and entries in h without a corresponding node in g are not altered.
func HeatDiffusion(g graph.Undirected, h map[int64]float64, t float64, normalize bool) map[int64]float64 {
	// See https://doi.org/10.1371/journal.pcbi.1005598 for details.

	nodes := g.Nodes()
	sort.Sort(ordered.ByID(nodes))
	indexOf := make(map[int64]int, len(nodes))
	heat := make([]float64, len(nodes))
	for i, n := range nodes {
		id := n.ID()
		indexOf[id] = i
		heat[i] = h[id]
	}
	v := mat.NewVecDense(len(heat), heat)

	l := mat.NewSymDense(len(nodes), nil)
	if normalize {
		for j, u := range nodes {
			to := g.From(u)
			if len(to) != 0 {
				l.SetSym(j, j, -t)
			}
			uid := u.ID()
			udeg := math.Sqrt(float64(len(to)))
			for _, v := range to {
				vid := v.ID()
				if uid < vid {
					l.SetSym(indexOf[vid], j, t/(udeg*math.Sqrt(float64(len(g.From(v))))))
				}
			}
		}
	} else {
		for j, u := range nodes {
			to := g.From(u)
			l.SetSym(j, j, -t*float64(len(to)))
			uid := u.ID()
			for _, v := range to {
				vid := v.ID()
				if uid < vid {
					l.SetSym(indexOf[vid], j, t)
				}
			}
		}
	}

	var m mat.Dense
	m.Exp(l)
	v.MulVec(&m, v)

	for i, n := range v.RawVector().Data {
		h[nodes[i].ID()] = n
	}
	return h
}
