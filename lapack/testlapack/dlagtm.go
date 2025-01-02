// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlagtmer interface {
	Dlagtm(trans blas.Transpose, m, n int, alpha float64, dl, d, du []float64, b []float64, ldb int, beta float64, c []float64, ldc int)
}

func DlagtmTest(t *testing.T, impl Dlagtmer) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans, blas.ConjTrans} {
		t.Run(transToString(trans), func(t *testing.T) {
			for _, m := range []int{0, 1, 2, 3, 4, 5, 10} {
				for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
					for _, ldb := range []int{max(1, n), n + 3} {
						for _, ldc := range []int{max(1, n), n + 4} {
							for _, alpha := range []float64{0, 1, rnd.NormFloat64()} {
								for _, beta := range []float64{0, 1, rnd.NormFloat64()} {
									dlagtmTest(t, impl, rnd, trans, m, n, ldb, ldc, alpha, beta)
								}
							}
						}
					}
				}
			}
		})
	}
}

func dlagtmTest(t *testing.T, impl Dlagtmer, rnd *rand.Rand, trans blas.Transpose, m, n int, ldb, ldc int, alpha, beta float64) {
	const (
		tol   = 1e-14
		extra = 10
	)

	name := fmt.Sprintf("Case m=%v,n=%v,ldb=%v,ldc=%v,alpha=%v,beta=%v", m, n, ldb, ldc, alpha, beta)

	// Generate three random diagonals.
	dl := randomSlice(n+extra, rnd)
	dlCopy := make([]float64, len(dl))
	copy(dlCopy, dl)

	d := randomSlice(n+1+extra, rnd)
	dCopy := make([]float64, len(d))
	copy(dCopy, d)

	du := randomSlice(n+extra, rnd)
	duCopy := make([]float64, len(du))
	copy(duCopy, du)

	b := randomGeneral(m, n, ldb, rnd)
	bCopy := cloneGeneral(b)

	got := randomGeneral(m, n, ldc, rnd)
	want := cloneGeneral(got)

	// Deal with zero-sized matrices early.
	if m == 0 || n == 0 {
		impl.Dlagtm(trans, m, n, alpha, dl, d, du, b.Data, b.Stride, beta, got.Data, got.Stride)
		if !floats.Same(dl, dlCopy) {
			t.Errorf("%v: unexpected modification in dl", name)
		}
		if !floats.Same(d, dCopy) {
			t.Errorf("%v: unexpected modification in d", name)
		}
		if !floats.Same(du, duCopy) {
			t.Errorf("%v: unexpected modification in du", name)
		}
		if !floats.Same(b.Data, bCopy.Data) {
			t.Errorf("%v: unexpected modification in B", name)
		}
		if !floats.Same(got.Data, want.Data) {
			t.Errorf("%v: unexpected modification in C", name)
		}
		return
	}

	impl.Dlagtm(trans, m, n, alpha, dl, d, du, b.Data, b.Stride, beta, got.Data, got.Stride)

	if !floats.Same(dl, dlCopy) {
		t.Errorf("%v: unexpected modification in dl", name)
	}
	if !floats.Same(d, dCopy) {
		t.Errorf("%v: unexpected modification in d", name)
	}
	if !floats.Same(du, duCopy) {
		t.Errorf("%v: unexpected modification in du", name)
	}
	if !floats.Same(b.Data, bCopy.Data) {
		t.Errorf("%v: unexpected modification in B", name)
	}

	// Generate a dense representation of the matrix and compute the wanted result.
	a := zeros(m, m, m)
	for i := 0; i < m-1; i++ {
		a.Data[i*a.Stride+i] = d[i]
		a.Data[i*a.Stride+i+1] = du[i]
		a.Data[(i+1)*a.Stride+i] = dl[i]
	}
	a.Data[(m-1)*a.Stride+m-1] = d[m-1]

	blas64.Gemm(trans, blas.NoTrans, alpha, a, b, beta, want)

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			got.Data[i*got.Stride+j] -= want.Data[i*want.Stride+j]
		}
	}
	diff := dlange(lapack.MaxColumnSum, got.Rows, got.Cols, got.Data, got.Stride)
	if diff > tol {
		t.Errorf("%v: unexpected result; diff=%v, want<=%v", name, diff, tol)
	}
}
