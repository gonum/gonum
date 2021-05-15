// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math/bits"
	"sync"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/blas/cblas128"
)

// poolFor returns the ceiling of base 2 log of size. It provides an index
// into a pool array to a sync.Pool that will return values able to hold
// size elements.
func poolFor(size uint) int {
	if size == 0 {
		return 0
	}
	return bits.Len(size - 1)
}

var (
	// poolDense contains size stratified workspace Dense pools.
	// Each poolDense element i returns sized matrices with a data
	// slice capped at 1<<i.
	poolDense [63]sync.Pool

	// poolSymDense is the SymDense equivalent of poolDense.
	poolSymDense [63]sync.Pool

	// poolTriDense is the TriDense equivalent of poolDense.
	poolTriDense [63]sync.Pool

	// poolVecDense is the VecDense equivalent of poolDense.
	poolVecDense [63]sync.Pool

	// poolCDense is the CDense equivalent of poolDense.
	poolCDense [63]sync.Pool

	// poolFloat64s is the []float64 equivalent of poolDense.
	poolFloat64s [63]sync.Pool

	// poolInts is the []int equivalent of poolDense.
	poolInts [63]sync.Pool
)

func init() {
	for i := range poolDense {
		l := 1 << uint(i)
		// Real matrix pools.
		poolDense[i].New = func() interface{} {
			return &Dense{mat: blas64.General{
				Data: make([]float64, l),
			}}
		}
		poolSymDense[i].New = func() interface{} {
			return &SymDense{mat: blas64.Symmetric{
				Uplo: blas.Upper,
				Data: make([]float64, l),
			}}
		}
		poolTriDense[i].New = func() interface{} {
			return &TriDense{mat: blas64.Triangular{
				Data: make([]float64, l),
			}}
		}
		poolVecDense[i].New = func() interface{} {
			return &VecDense{mat: blas64.Vector{
				Inc:  1,
				Data: make([]float64, l),
			}}
		}

		// Complex matrix pools.
		poolCDense[i].New = func() interface{} {
			return &CDense{mat: cblas128.General{
				Data: make([]complex128, l),
			}}
		}

		// Helper pools.
		poolFloat64s[i].New = func() interface{} {
			s := make([]float64, l)
			return &s
		}
		poolInts[i].New = func() interface{} {
			s := make([]int, l)
			return &s
		}
	}
}

// getWorkspace returns a *Dense of size r×c and a data slice
// with a cap that is less than 2*r*c. If clear is true, the
// data slice visible through the Matrix interface is zeroed.
func getWorkspace(r, c int, clear bool) *Dense {
	l := uint(r * c)
	w := poolDense[poolFor(l)].Get().(*Dense)
	w.mat.Data = w.mat.Data[:l]
	if clear {
		zero(w.mat.Data)
	}
	w.mat.Rows = r
	w.mat.Cols = c
	w.mat.Stride = c
	w.capRows = r
	w.capCols = c
	return w
}

// putWorkspace replaces a used *Dense into the appropriate size
// workspace pool. putWorkspace must not be called with a matrix
// where references to the underlying data slice have been kept.
func putWorkspace(w *Dense) {
	poolDense[poolFor(uint(cap(w.mat.Data)))].Put(w)
}

// getWorkspaceSym returns a *SymDense of size n and a cap that
// is less than 2*n. If clear is true, the data slice visible
// through the Matrix interface is zeroed.
func getWorkspaceSym(n int, clear bool) *SymDense {
	l := uint(n)
	l *= l
	s := poolSymDense[poolFor(l)].Get().(*SymDense)
	s.mat.Data = s.mat.Data[:l]
	if clear {
		zero(s.mat.Data)
	}
	s.mat.N = n
	s.mat.Stride = n
	s.cap = n
	return s
}

// putWorkspaceSym replaces a used *SymDense into the appropriate size
// workspace pool. putWorkspaceSym must not be called with a matrix
// where references to the underlying data slice have been kept.
func putWorkspaceSym(s *SymDense) {
	poolSymDense[poolFor(uint(cap(s.mat.Data)))].Put(s)
}

// getWorkspaceTri returns a *TriDense of size n and a cap that
// is less than 2*n. If clear is true, the data slice visible
// through the Matrix interface is zeroed.
func getWorkspaceTri(n int, kind TriKind, clear bool) *TriDense {
	l := uint(n)
	l *= l
	t := poolTriDense[poolFor(l)].Get().(*TriDense)
	t.mat.Data = t.mat.Data[:l]
	if clear {
		zero(t.mat.Data)
	}
	t.mat.N = n
	t.mat.Stride = n
	if kind == Upper {
		t.mat.Uplo = blas.Upper
	} else if kind == Lower {
		t.mat.Uplo = blas.Lower
	} else {
		panic(ErrTriangle)
	}
	t.mat.Diag = blas.NonUnit
	t.cap = n
	return t
}

// putWorkspaceTri replaces a used *TriDense into the appropriate size
// workspace pool. putWorkspaceTri must not be called with a matrix
// where references to the underlying data slice have been kept.
func putWorkspaceTri(t *TriDense) {
	poolTriDense[poolFor(uint(cap(t.mat.Data)))].Put(t)
}

// getWorkspaceVec returns a *VecDense of length n and a cap that
// is less than 2*n. If clear is true, the data slice visible
// through the Matrix interface is zeroed.
func getWorkspaceVec(n int, clear bool) *VecDense {
	l := uint(n)
	v := poolVecDense[poolFor(l)].Get().(*VecDense)
	v.mat.Data = v.mat.Data[:l]
	if clear {
		zero(v.mat.Data)
	}
	v.mat.N = n
	return v
}

// putWorkspaceVec replaces a used *VecDense into the appropriate size
// workspace pool. putWorkspaceVec must not be called with a matrix
// where references to the underlying data slice have been kept.
func putWorkspaceVec(v *VecDense) {
	poolVecDense[poolFor(uint(cap(v.mat.Data)))].Put(v)
}

// getFloats returns a []float64 of length l and a cap that is
// less than 2*l. If clear is true, the slice visible is zeroed.
func getFloats(l int, clear bool) []float64 {
	w := *poolFloat64s[poolFor(uint(l))].Get().(*[]float64)
	w = w[:l]
	if clear {
		zero(w)
	}
	return w
}

// putFloats replaces a used []float64 into the appropriate size
// workspace pool. putFloats must not be called with a slice
// where references to the underlying data have been kept.
func putFloats(w []float64) {
	poolFloat64s[poolFor(uint(cap(w)))].Put(&w)
}

// getInts returns a []ints of length l and a cap that is
// less than 2*l. If clear is true, the slice visible is zeroed.
func getInts(l int, clear bool) []int {
	w := *poolInts[poolFor(uint(l))].Get().(*[]int)
	w = w[:l]
	if clear {
		for i := range w {
			w[i] = 0
		}
	}
	return w
}

// putInts replaces a used []int into the appropriate size
// workspace pool. putInts must not be called with a slice
// where references to the underlying data have been kept.
func putInts(w []int) {
	poolInts[poolFor(uint(cap(w)))].Put(&w)
}

// getWorkspaceCmplx returns a *CDense of size r×c and a data slice
// with a cap that is less than 2*r*c. If clear is true, the
// data slice visible through the CMatrix interface is zeroed.
func getWorkspaceCmplx(r, c int, clear bool) *CDense {
	l := uint(r * c)
	w := poolCDense[poolFor(l)].Get().(*CDense)
	w.mat.Data = w.mat.Data[:l]
	if clear {
		zeroC(w.mat.Data)
	}
	w.mat.Rows = r
	w.mat.Cols = c
	w.mat.Stride = c
	w.capRows = r
	w.capCols = c
	return w
}

// putWorkspaceCmplx replaces a used *CDense into the appropriate size
// workspace pool. putWorkspace must not be called with a matrix
// where references to the underlying data slice have been kept.
func putWorkspaceCmplx(w *CDense) {
	poolCDense[poolFor(uint(cap(w.mat.Data)))].Put(w)
}
