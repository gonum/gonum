// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mds

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

// TorgersonScaling converts a dissimilarity matrix to a matrix containing
// Euclidean coordinates. TorgersonScaling places the coordinates in dst and
// returns it and true if successful. If the scaling is not successful, dst
// is returned, but will not be a valid scaling.
//
// If dst is nil, a new mat.Dense is allocated. If dst is not a zero matrix,
// the dimensions of dst and dis must match otherwise TorgersonScaling will panic.
// The dis matrix must be square or TorgersonScaling will panic.
func TorgersonScaling(dst *mat.Dense, dis mat.Symmetric) (mds *mat.Dense, ok bool) {
	// https://doi.org/10.1007/0-387-28981-X_12

	n := dis.Symmetric()
	if dst == nil {
		dst = mat.NewDense(n, n, nil)
	} else if r, c := dst.Dims(); !dst.IsZero() && (r != n || c != n) {
		panic(mat.ErrShape)
	}

	b := mat.NewSymDense(n, nil)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			v := dis.At(i, j)
			v *= v
			b.SetSym(i, j, v)
		}
	}
	c := mat.NewSymDense(n, nil)
	s := -1 / float64(n)
	for i := 0; i < n; i++ {
		c.SetSym(i, i, 1+s)
		for j := i + 1; j < n; j++ {
			c.SetSym(i, j, s)
		}
	}
	dst.Product(c, b, c)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			b.SetSym(i, j, -0.5*dst.At(i, j))
		}
	}

	var ed mat.EigenSym
	ok = ed.Factorize(b, true)
	if !ok {
		return dst, false
	}
	dst.EigenvectorsSym(&ed)
	vals := ed.Values(nil)
	sort.Sort(byValues{
		values:  vals,
		vectors: dst.RawMatrix(),
	})
	for i, v := range vals {
		if v < 0 {
			vals[i] = 0
			continue
		}
		vals[i] = math.Sqrt(v)
	}
	dst.Mul(dst, mat.NewDiagonal(len(vals), vals))

	return dst, true
}

type byValues struct {
	values  []float64
	vectors blas64.General
}

func (e byValues) Len() int           { return len(e.values) }
func (e byValues) Less(i, j int) bool { return e.values[i] > e.values[j] }
func (e byValues) Swap(i, j int) {
	e.values[i], e.values[j] = e.values[j], e.values[i]
	blas64.Swap(e.vectors.Rows,
		blas64.Vector{Inc: e.vectors.Stride, Data: e.vectors.Data[i:]},
		blas64.Vector{Inc: e.vectors.Stride, Data: e.vectors.Data[j:]},
	)
}
