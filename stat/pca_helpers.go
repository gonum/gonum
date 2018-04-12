// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

// TransformDisjunctive transforms a Complete Disjunctive Table to a Transformed
// Complete Disjunctive Table in order to be able to perform Multiple
// Correspondence Analysis using PC.PrincipalComponents. TransformDisjunctive
// places the TCDT in dst and returns it.
//
// If dst is nil, a new mat.Dense is allocated. If dst is not a zero matrix,
// the dimensions of dst and cdt must match otherwise TransformDisjunctive will
// panic. If cdt contains values other 0 or 1 TransformDisjunctive will panic.
//
// It is safe to reuse cdt as dst.
func TransformDisjunctive(dst *mat.Dense, cdt mat.Matrix) *mat.Dense {
	r, c := cdt.Dims()

	if dst == nil {
		dst = mat.NewDense(r, c, nil)
	} else if dr, dc := dst.Dims(); !dst.IsZero() && (dr != r || dc != c) {
		panic(mat.ErrShape)
	}

	for j := 0; j < c; j++ {
		var p float64
		for i := 0; i < r; i++ {
			v := cdt.At(i, j)
			if v != 0 && v != 1 {
				panic("stat: input is not a complete disjunctive table")
			}
			p += v
		}
		for i := 0; i < r; i++ {
			dst.Set(i, j, cdt.At(i, j)/p-1)
		}
	}

	return dst
}

// TorgersonScaling converts a dissimilarity matrix to a matrix containing
// Euclidean coordinates in order to be able to perform Torgerson's Classical
// Multidimensional Scaling using PC.PrincipleComponenets. TorgersonScaling places
// the coordinates in dst and returns it.
//
// If dst is nil, a new mat.Dense is allocated. If dst is not a zero matrix,
// the dimensions of dst and dis must match otherwise TorgersonScaling will panic.
// The dis matrix must be square or TorgersonScaling will panic.
func TorgersonScaling(dst *mat.Dense, dis mat.Symmetric) *mat.Dense {
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
	for i := 0; i < n; i++ {
		s := -1 / float64(n)
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
	ed.Factorize(b, true)
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

	return dst
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
