// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"gonum.org/v1/gonum/graph/spectral"
	"gonum.org/v1/gonum/mat"
)

// Diffuse performs a heat diffusion across nodes of the undirected
// graph described by the given Laplacian using the initial heat distribution,
// h, according to the Laplacian with a diffusion time of t.
// The resulting heat distribution is returned, written into the map dst and
// returned,
//
//	d = exp(-Lt)×h
//
// where L is the graph Laplacian. Indexing into h and dst is defined by the
// Laplacian Index field. If dst is nil, a new map is created.
//
// Nodes without corresponding entries in h are given an initial heat of zero,
// and entries in h without a corresponding node in the original graph are
// not altered when written to dst.
func Diffuse(dst, h map[int64]float64, by spectral.Laplacian, t float64) map[int64]float64 {
	heat := make([]float64, len(by.Index))
	for id, i := range by.Index {
		heat[i] = h[id]
	}
	v := mat.NewVecDense(len(heat), heat)

	var m, tl mat.Dense
	tl.Scale(-t, by)
	m.Exp(&tl)
	v.MulVec(&m, v)

	if dst == nil {
		dst = make(map[int64]float64)
	}
	for i, n := range heat {
		dst[by.Nodes[i].ID()] = n
	}
	return dst
}

// DiffuseToEquilibrium performs a heat diffusion across nodes of the
// graph described by the given Laplacian using the initial heat
// distribution, h, according to the Laplacian until the update function
//
//	h_{n+1} = h_n - L×h_n
//
// results in a 2-norm update difference within tol, or iters updates have
// been made.
// The resulting heat distribution is returned as eq, written into the map dst,
// and a boolean indicating whether the equilibrium converged to within tol.
// Indexing into h and dst is defined by the Laplacian Index field. If dst
// is nil, a new map is created.
//
// Nodes without corresponding entries in h are given an initial heat of zero,
// and entries in h without a corresponding node in the original graph are
// not altered when written to dst.
func DiffuseToEquilibrium(dst, h map[int64]float64, by spectral.Laplacian, tol float64, iters int) (eq map[int64]float64, ok bool) {
	heat := make([]float64, len(by.Index))
	for id, i := range by.Index {
		heat[i] = h[id]
	}
	v := mat.NewVecDense(len(heat), heat)

	last := make([]float64, len(by.Index))
	for id, i := range by.Index {
		last[i] = h[id]
	}
	lastV := mat.NewVecDense(len(last), last)

	var tmp mat.VecDense
	for {
		iters--
		if iters < 0 {
			break
		}
		lastV, v = v, lastV
		tmp.MulVec(by.Matrix, lastV)
		v.SubVec(lastV, &tmp)
		if normDiff(heat, last) < tol {
			ok = true
			break
		}
	}

	if dst == nil {
		dst = make(map[int64]float64)
	}
	for i, n := range v.RawVector().Data {
		dst[by.Nodes[i].ID()] = n
	}
	return dst, ok
}
