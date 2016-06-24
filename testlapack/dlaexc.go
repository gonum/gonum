// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dlaexcer interface {
	Dlaexc(wantq bool, n int, t []float64, ldt int, q []float64, ldq int, j1, n1, n2 int, work []float64) bool
}

func DlaexcTest(t *testing.T, impl Dlaexcer) {
	rnd := rand.New(rand.NewSource(1))

	for _, wantq := range []bool{true, false} {
		for _, n := range []int{1, 2, 3, 4, 5, 6, 10, 18, 31, 53} {
			for _, extra := range []int{0, 1, 11} {
				for cas := 0; cas < 100; cas++ {
					j1 := rnd.Intn(n)
					n1 := min(rnd.Intn(3), n-j1)
					n2 := min(rnd.Intn(3), n-j1-n1)
					testDlaexc(t, impl, wantq, n, j1, n1, n2, extra, rnd)
				}
			}
		}
	}
}

func testDlaexc(t *testing.T, impl Dlaexcer, wantq bool, n, j1, n1, n2, extra int, rnd *rand.Rand) {
	const tol = 1e-14

	tmat := randomGeneral(n, n, n+extra, rnd)
	// Zero out the lower triangle.
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			tmat.Data[i*tmat.Stride+j] = 0
		}
	}
	// Make any 2x2 diagonal block to be in Schur canonical form.
	if n1 == 2 {
		tmat.Data[(j1+1)*tmat.Stride+j1+1] = tmat.Data[j1*tmat.Stride+j1]
		tmat.Data[(j1+1)*tmat.Stride+j1] = tmat.Data[j1*tmat.Stride+j1+1]
	}
	if n2 == 2 {
		tmat.Data[(j1+n1+1)*tmat.Stride+j1+n1+1] = tmat.Data[(j1+n1)*tmat.Stride+j1+n1]
		tmat.Data[(j1+n1+1)*tmat.Stride+j1+n1] = tmat.Data[(j1+n1)*tmat.Stride+j1+n1+1]
	}
	tmatCopy := cloneGeneral(tmat)
	var q, qCopy blas64.General
	if wantq {
		q = eye(n, n+extra)
		qCopy = cloneGeneral(q)
	}
	work := nanSlice(n)

	ok := impl.Dlaexc(wantq, n, tmat.Data, tmat.Stride, q.Data, q.Stride, j1, n1, n2, work)

	prefix := fmt.Sprintf("Case n=%v, j1=%v, n1=%v, n2=%v, wantq=%v, extra=%v", n, j1, n1, n2, wantq, extra)
	if !generalOutsideAllNaN(tmat) {
		t.Errorf("%v: out-of-range write to T\n", prefix)
	}
	if wantq && !generalOutsideAllNaN(q) {
		t.Errorf("%v: out-of-range write to Q\n", prefix)
	}

	if !ok || n1 == 0 || n2 == 0 || j1+n1 >= n {
		// Check that T is not modified.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if tmat.Data[i*tmat.Stride+j] != tmatCopy.Data[i*tmatCopy.Stride+j] {
					t.Errorf("%v: ok == false but T[%v,%v] modified\n", prefix, i, j)
				}
			}
		}
		if !wantq {
			return
		}
		// Check that Q is not modified.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if q.Data[i*q.Stride+j] != qCopy.Data[i*qCopy.Stride+j] {
					t.Errorf("%v: ok == false but Q[%v,%v] modified\n", prefix, i, j)
				}
			}
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
				t.Errorf("%v: unexpected modification of T[%v,%v]\n", prefix, i, j)
			}
		}
	}
	// Check that T is modified at the diagonal as expected.
	for i := 0; i < n1; i++ {
		for j := 0; j < n1; j++ {
			got := tmat.Data[(j1+i)*tmat.Stride+j1+j]
			want := tmatCopy.Data[(j1+n1+i)*tmatCopy.Stride+j1+n1+j]
			if want != got {
				t.Errorf("%v: unexpected value of T[%v,%v]. want %v, got %v\n", prefix, j1+i, j1+j, want, got)
			}
		}
	}
	for i := 0; i < n2; i++ {
		for j := 0; j < n2; j++ {
			got := tmat.Data[(j1+n1+i)*tmat.Stride+j1+n1+j]
			want := tmatCopy.Data[(j1+i)*tmatCopy.Stride+j1+j]
			if want != got {
				t.Errorf("%v: unexpected value of T[%v,%v]. want %v, got %v\n", prefix, j1+n1+i, j1+n1+j, want, got)
			}
		}
	}

	if !wantq {
		return
	}

	if !isOrthonormal(q) {
		t.Errorf("%v: Q is not orthogonal\n", prefix)
	}
	// Check that Q is unchanged outside of columns [j1:j1+n1+n2].
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j1 <= j && j < j1+n1+n2 {
				continue
			}
			diff := q.Data[i*q.Stride+j] - qCopy.Data[i*qCopy.Stride+j]
			if diff != 0 {
				t.Errorf("%v: unexpected modification of Q[%v,%v]\n", prefix, i, j)
			}
		}
	}
	// Check that Q^T TOrig Q == T.
	tq := eye(n, n)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmatCopy, q, 0, tq)
	qtq := eye(n, n)
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, tq, 0, qtq)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			diff := qtq.Data[i*qtq.Stride+j] - tmat.Data[i*tmat.Stride+j]
			if math.Abs(diff) > tol {
				t.Errorf("%v: unexpected value of T[%v,%v]\n", prefix, i, j)
			}
		}
	}
}
