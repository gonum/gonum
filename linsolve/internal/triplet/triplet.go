// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package triplet provides triplet representation for sparse matrices.
package triplet

import "gonum.org/v1/gonum/mat"

type triplet struct {
	i, j int
	v    float64
}

type Matrix struct {
	r, c int
	data []triplet
}

func NewMatrix(r, c int) *Matrix {
	if r <= 0 || c <= 0 {
		panic("triplet: invalid shape")
	}
	return &Matrix{
		r: r,
		c: c,
	}
}

func (m *Matrix) Dims() (r, c int) {
	return m.r, m.c
}

// Append appends a non-zero element to the list of matrix elements without
// checking whether it already exists.
func (m *Matrix) Append(i, j int, v float64) {
	if i < 0 || m.r <= i {
		panic("triplet: row index out of range")
	}
	if j < 0 || m.c <= j {
		panic("triplet: column index out of range")
	}
	if v == 0 {
		return
	}
	m.data = append(m.data, triplet{i, j, v})
}

func (m *Matrix) MulVecTo(dst *mat.VecDense, trans bool, x mat.Vector) {
	dst.Zero()
	if trans {
		if m.c != dst.Len() || m.r != x.Len() {
			panic("triplet: dimension mismatch")
		}
		for _, aij := range m.data {
			dst.SetVec(aij.j, dst.AtVec(aij.j)+aij.v*x.AtVec(aij.i))
		}
		return
	}
	if m.c != x.Len() || m.r != dst.Len() {
		panic("triplet: dimension mismatch")
	}
	for _, aij := range m.data {
		dst.SetVec(aij.i, dst.AtVec(aij.i)+aij.v*x.AtVec(aij.j))
	}
}

func (m *Matrix) DenseCopy() *mat.Dense {
	d := mat.NewDense(m.r, m.c, nil)
	for _, aij := range m.data {
		d.Set(aij.i, aij.j, aij.v)
	}
	return d
}
