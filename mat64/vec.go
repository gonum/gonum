// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

var (
	vector *Vector

	_ Matrix  = vector
	_ Mutable = vector

	// _ Cloner      = vector
	// _ Viewer      = vector
	// _ Subvectorer = vector

	// _ Adder     = vector
	// _ Suber     = vector
	// _ Muler = vector
	// _ Dotter    = vector
	// _ ElemMuler = vector

	// _ Scaler  = vector
	// _ Applyer = vector

	// _ Normer = vector
	// _ Sumer  = vector

	// _ Stacker   = vector
	// _ Augmenter = vector

	// _ Equaler       = vector
	// _ ApproxEqualer = vector

	// _ RawMatrixLoader = vector
	// _ RawMatrixer     = vector
)

// Vector represents a column vector.
type Vector struct {
	mat blas64.Vector
	n   int
	// A BLAS vector can have a negative increment, but allowing this
	// in the mat64 type complicates a lot of code, and doesn't gain anything.
	// Vector must have positive increment in this package.
}

func NewVector(n int, data []float64) *Vector {
	if len(data) != n && data != nil {
		panic(ErrShape)
	}
	if data == nil {
		data = make([]float64, n)
	}
	return &Vector{
		mat: blas64.Vector{
			Inc:  1,
			Data: data,
		},
		n: n,
	}
}

func (m *Vector) ViewVec(i, n int) *Vector {
	if i+n > m.n {
		panic(ErrIndexOutOfRange)
	}
	return &Vector{
		n: n,
		mat: blas64.Vector{
			Inc:  m.mat.Inc,
			Data: m.mat.Data[i*m.mat.Inc:],
		},
	}
}

func (m *Vector) At(r, c int) float64 {
	if c != 0 || r < 0 || r >= m.n {
		panic(ErrIndexOutOfRange)
	}
	return m.mat.Data[r*m.mat.Inc]
}

func (m *Vector) Set(r, c int, v float64) {
	if c != 0 || r < 0 || r >= m.n {
		panic(ErrIndexOutOfRange)
	}
	m.mat.Data[r*m.mat.Inc] = v
}

func (m *Vector) Dims() (r, c int) { return m.n, 1 }

func (m *Vector) Reset() {
	m.mat.Data = m.mat.Data[:0]
	m.mat.Inc = 0
	m.n = 0
}

func (m *Vector) RawVector() blas64.Vector {
	return m.mat
}

func (m *Vector) MulVec(a Matrix, trans bool, b *Vector) {
	// TODO (btracey): should there be some kind of interface for Vector?
	ar, ac := a.Dims()
	br, _ := b.Dims()
	if ac != br {
		panic(ErrShape)
	}

	var w Vector
	if m != a && m != b {
		w = *m
	}
	if w.n == 0 {
		w.mat.Data = use(w.mat.Data, ar)
	} else if ar != w.n {
		panic(ErrShape)
	}

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
		t := blas.NoTrans
		if trans {
			t = blas.Trans
		}
		blas64.Gemv(t,
			1, amat, b.mat,
			0, w.mat,
		)
		*m = w
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		for r := 0; r < ar; r++ {
			w.mat.Data[r*m.mat.Inc] = blas64.Dot(ac,
				blas64.Vector{Inc: 1, Data: a.Row(row, r)},
				b.mat,
			)
		}
		*m = w
		return
	}

	row := make([]float64, ac)
	for r := 0; r < ar; r++ {
		for i := range row {
			row[i] = a.At(r, i)
		}
		var v float64
		for i, e := range row {
			v += e * b.mat.Data[i*b.mat.Inc]
		}
		w.mat.Data[r*m.mat.Inc] = v
	}
	*m = w
}
