// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgtsver interface {
	Dgtsv(n, nrhs int, dl, d, du []float64, b []float64, ldb int) (ok bool)
}

func DgtsvTest(t *testing.T, impl Dgtsver) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 25, 50} {
		for _, nrhs := range []int{0, 1, 2, 3, 4, 10} {
			for _, ldb := range []int{max(1, nrhs), nrhs + 3} {
				dgtsvTest(t, impl, rnd, n, nrhs, ldb)
			}
		}
	}
}

func dgtsvTest(t *testing.T, impl Dgtsver, rnd *rand.Rand, n, nrhs, ldb int) {
	const (
		tol   = 1e-14
		extra = 10
	)

	name := fmt.Sprintf("Case n=%d,nrhs=%d,ldb=%d", n, nrhs, ldb)

	if n == 0 {
		ok := impl.Dgtsv(n, nrhs, nil, nil, nil, nil, ldb)
		if !ok {
			t.Errorf("%v: unexpected failure for zero size matrix", name)
		}
		return
	}

	// Generate three random diagonals.
	var (
		d, dCopy   []float64
		dl, dlCopy []float64
		du, duCopy []float64
	)
	d = randomSlice(n+1+extra, rnd)
	dCopy = make([]float64, len(d))
	copy(dCopy, d)
	if n > 1 {
		dl = randomSlice(n+extra, rnd)
		dlCopy = make([]float64, len(dl))
		copy(dlCopy, dl)

		du = randomSlice(n+extra, rnd)
		duCopy = make([]float64, len(du))
		copy(duCopy, du)
	}

	b := randomGeneral(n, nrhs, ldb, rnd)
	got := cloneGeneral(b)

	ok := impl.Dgtsv(n, nrhs, dl, d, du, got.Data, got.Stride)
	if !ok {
		t.Fatalf("%v: unexpected failure in Dgtsv", name)
		return
	}

	// Compute A*X - B.
	dlagtm(blas.NoTrans, n, nrhs, 1, dlCopy, dCopy, duCopy, got.Data, got.Stride, -1, b.Data, b.Stride)

	anorm := dlangt(lapack.MaxColumnSum, n, dlCopy, dCopy, duCopy)
	bi := blas64.Implementation()
	var resid float64
	for j := 0; j < nrhs; j++ {
		bnorm := bi.Dasum(n, b.Data[j:], b.Stride)
		xnorm := bi.Dasum(n, got.Data[j:], got.Stride)
		resid = math.Max(resid, bnorm/anorm/xnorm)
	}
	if resid > tol {
		t.Errorf("%v: unexpected result; resid=%v,want<=%v", name, resid, tol)
	}
}
