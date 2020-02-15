// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dlaexcer interface {
	Dlaexc(wantq bool, n int, t []float64, ldt int, q []float64, ldq int, j1, n1, n2 int, work []float64) bool
}

func DlaexcTest(t *testing.T, impl Dlaexcer) {
	rnd := rand.New(rand.NewSource(1))

	for _, n := range []int{1, 2, 3, 4, 5, 6, 10, 18, 31, 53} {
		for _, extra := range []int{0, 3} {
			for cas := 0; cas < 100; cas++ {
				testDlaexc(t, impl, rnd, n, extra)
			}
		}
	}
}

func testDlaexc(t *testing.T, impl Dlaexcer, rnd *rand.Rand, n, extra int) {
	const tol = 1e-14

	// Generate random T in Schur canonical form.
	tmat, _, _ := randomSchurCanonical(n, n+extra, true, rnd)
	tmatCopy := cloneGeneral(tmat)

	// Randomly pick the index of the first block.
	j1 := rnd.Intn(n)
	if j1 > 0 && tmat.Data[j1*tmat.Stride+j1-1] != 0 {
		// Adjust j1 if it points to the second row of a 2x2 block.
		j1--
	}
	// Read sizes of the two blocks based on properties of T.
	var n1, n2 int
	switch j1 {
	case n - 1:
		n1, n2 = 1, 0
	case n - 2:
		if tmat.Data[(j1+1)*tmat.Stride+j1] == 0 {
			n1, n2 = 1, 1
		} else {
			n1, n2 = 2, 0
		}
	case n - 3:
		if tmat.Data[(j1+1)*tmat.Stride+j1] == 0 {
			n1, n2 = 1, 2
		} else {
			n1, n2 = 2, 1
		}
	default:
		if tmat.Data[(j1+1)*tmat.Stride+j1] == 0 {
			n1 = 1
			if tmat.Data[(j1+2)*tmat.Stride+j1+1] == 0 {
				n2 = 1
			} else {
				n2 = 2
			}
		} else {
			n1 = 2
			if tmat.Data[(j1+3)*tmat.Stride+j1+2] == 0 {
				n2 = 1
			} else {
				n2 = 2
			}
		}
	}

	name := fmt.Sprintf("Case n=%v,j1=%v,n1=%v,n2=%v,extra=%v", n, j1, n1, n2, extra)

	// 1. Test without accumulating Q.

	wantq := false

	work := nanSlice(n)

	ok := impl.Dlaexc(wantq, n, tmat.Data, tmat.Stride, nil, 1, j1, n1, n2, work)

	// 2. Test with accumulating Q.

	wantq = true

	tmat2 := cloneGeneral(tmatCopy)
	q := eye(n, n+extra)
	qCopy := cloneGeneral(q)
	work = nanSlice(n)

	ok2 := impl.Dlaexc(wantq, n, tmat2.Data, tmat2.Stride, q.Data, q.Stride, j1, n1, n2, work)

	if !generalOutsideAllNaN(tmat) {
		t.Errorf("%v: out-of-range write to T", name)
	}
	if !generalOutsideAllNaN(tmat2) {
		t.Errorf("%v: out-of-range write to T2", name)
	}
	if !generalOutsideAllNaN(q) {
		t.Errorf("%v: out-of-range write to Q", name)
	}

	// Check that outputs from cases 1. and 2. are exactly equal, then check one of them.
	if ok != ok2 {
		t.Errorf("%v: ok != ok2", name)
	}
	if !equalGeneral(tmat, tmat2) {
		t.Errorf("%v: T != T2", name)
	}

	if !ok {
		if n1 == 1 && n2 == 1 {
			t.Errorf("%v: unexpected failure", name)
		} else {
			t.Logf("%v: Dlaexc returned false", name)
		}
	}

	if !ok || n1 == 0 || n2 == 0 || j1+n1 >= n {
		// Check that T is not modified.
		if !equalGeneral(tmat, tmatCopy) {
			t.Errorf("%v: unexpected modification of T", name)
		}
		// Check that Q is not modified.
		if !equalGeneral(q, qCopy) {
			t.Errorf("%v: unexpected modification of Q", name)
		}
		return
	}

	// Check that T is not modified outside of rows and columns [j1:j1+n1+n2].
	for i := 0; i < n; i++ {
		if j1 <= i && i < j1+n1+n2 {
			continue
		}
		for j := 0; j < n; j++ {
			if j1 <= j && j < j1+n1+n2 {
				continue
			}
			diff := tmat.Data[i*tmat.Stride+j] - tmatCopy.Data[i*tmatCopy.Stride+j]
			if diff != 0 {
				t.Errorf("%v: unexpected modification of T[%v,%v]", name, i, j)
			}
		}
	}

	if !isSchurCanonicalGeneral(tmat) {
		t.Errorf("%v: T is not in Schur canonical form", name)
	}

	// Check that Q is orthogonal.
	resid := residualOrthogonal(q, false)
	if resid > tol {
		t.Errorf("%v: Q is not orthogonal; resid=%v, want<=%v", name, resid, tol)
	}

	// Check that Q is unchanged outside of columns [j1:j1+n1+n2].
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j1 <= j && j < j1+n1+n2 {
				continue
			}
			diff := q.Data[i*q.Stride+j] - qCopy.Data[i*qCopy.Stride+j]
			if diff != 0 {
				t.Errorf("%v: unexpected modification of Q[%v,%v]", name, i, j)
			}
		}
	}

	// Check that Qᵀ * TOrig * Q == T
	qt := zeros(n, n, n)
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, tmatCopy, 0, qt)
	qtq := cloneGeneral(tmat)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, qt, q, 1, qtq)
	resid = dlange(lapack.MaxColumnSum, n, n, qtq.Data, qtq.Stride)
	if resid > float64(n)*tol {
		t.Errorf("%v: mismatch between Qᵀ*(initial T)*Q and (final T); resid=%v, want<=%v",
			name, resid, float64(n)*tol)
	}
}
