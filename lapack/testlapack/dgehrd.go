// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dgehrder interface {
	Dgehrd(n, ilo, ihi int, a []float64, lda int, tau, work []float64, lwork int)

	Dorgqr(m, n, k int, a []float64, lda int, tau, work []float64, lwork int)
}

func DgehrdTest(t *testing.T, impl Dgehrder) {
	rnd := rand.New(rand.NewSource(1))

	// Randomized tests for small matrix sizes that will most likely
	// use the unblocked algorithm.
	for _, n := range []int{1, 2, 3, 4, 5, 10, 34} {
		for _, extra := range []int{0, 13} {
			for _, optwork := range []bool{true, false} {
				for cas := 0; cas < 10; cas++ {
					ilo := rnd.Intn(n)
					ihi := rnd.Intn(n)
					if ilo > ihi {
						ilo, ihi = ihi, ilo
					}
					testDgehrd(t, impl, n, ilo, ihi, extra, optwork, rnd)
				}
			}
		}
	}

	// These are selected tests for larger matrix sizes to test the blocked
	// algorithm. Use sizes around several powers of two because that is
	// where the blocked path will most likely start to be taken. For
	// example, at present the blocked algorithm is used for sizes larger
	// than 129.
	for _, test := range []struct {
		n, ilo, ihi int
	}{
		{0, 0, -1},

		{68, 0, 63},
		{68, 0, 64},
		{68, 0, 65},
		{68, 0, 66},
		{68, 0, 67},

		{132, 2, 129},
		{132, 1, 129}, // Size = 129, unblocked.
		{132, 0, 129}, // Size = 130, blocked.
		{132, 1, 130},
		{132, 0, 130},
		{132, 1, 131},
		{132, 0, 131},

		{260, 2, 257},
		{260, 1, 257},
		{260, 0, 257},
		{260, 0, 258},
		{260, 0, 259},
	} {
		for _, extra := range []int{0, 13} {
			for _, optwork := range []bool{true, false} {
				testDgehrd(t, impl, test.n, test.ilo, test.ihi, extra, optwork, rnd)
			}
		}
	}
}

func testDgehrd(t *testing.T, impl Dgehrder, n, ilo, ihi, extra int, optwork bool, rnd *rand.Rand) {
	a := randomGeneral(n, n, n+extra, rnd)
	aCopy := a
	aCopy.Data = make([]float64, len(a.Data))
	copy(aCopy.Data, a.Data)

	var tau []float64
	if n > 1 {
		tau = nanSlice(n - 1)
	}

	var work []float64
	if optwork {
		work = nanSlice(1)
		impl.Dgehrd(n, ilo, ihi, a.Data, a.Stride, tau, work, -1)
		work = nanSlice(int(work[0]))
	} else {
		work = nanSlice(max(1, n))
	}

	impl.Dgehrd(n, ilo, ihi, a.Data, a.Stride, tau, work, len(work))

	if n == 0 {
		// Just make sure there is no panic.
		return
	}

	prefix := fmt.Sprintf("Case n=%v, ilo=%v, ihi=%v, extra=%v", n, ilo, ihi, extra)

	// Check any invalid modifications of a.
	if !generalOutsideAllNaN(a) {
		t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
	}
	for i := ilo; i <= ihi; i++ {
		for j := 0; j < min(ilo, i); j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification of A[%v,%v]", prefix, i, j)
			}
		}
	}
	for i := ihi + 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification of A[%v,%v]", prefix, i, j)
			}
		}
	}
	for i := 0; i <= ilo; i++ {
		for j := i; j < ilo+1; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification at A[%v,%v]", prefix, i, j)
			}
		}
		for j := ihi + 1; j < n; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification at A[%v,%v]", prefix, i, j)
			}
		}
	}
	for i := ihi + 1; i < n; i++ {
		for j := i; j < n; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification at A[%v,%v]", prefix, i, j)
			}
		}
	}

	// Check that tau has been assigned properly.
	for i, v := range tau {
		if math.IsNaN(v) {
			t.Errorf("%v: unexpected NaN at tau[%v]", prefix, i)
		}
	}

	// Extract Q and check that it is orthogonal.
	q := eye(n, n)
	if ilo != ihi {
		for i := ilo + 2; i <= ihi; i++ {
			for j := ilo + 1; j < ihi; j++ {
				q.Data[i*q.Stride+j] = a.Data[i*a.Stride+j-1]
			}
		}
		nh := ihi - ilo
		impl.Dorgqr(nh, nh, nh, q.Data[(ilo+1)*q.Stride+ilo+1:], q.Stride, tau[ilo:ihi], work, len(work))
	}
	if !isOrthogonal(q) {
		t.Errorf("%v: Q is not orthogonal\nQ=%v", prefix, q)
	}

	// Construct Q^T * AOrig * Q and check that it is upper Hessenberg.
	aq := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aCopy, q, 0, aq)
	qaq := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, aq, 0, qaq)
	for i := 0; i <= ilo; i++ {
		for j := ilo + 1; j <= ihi; j++ {
			qaqij := qaq.Data[i*qaq.Stride+j]
			diff := qaqij - a.Data[i*a.Stride+j]
			if math.Abs(diff) > 1e-13 {
				t.Errorf("%v: Q^T*AOrig*Q and A are not equal, diff at [%v,%v]=%v", prefix, i, j, diff)
			}
		}
	}
	for i := ilo + 1; i <= ihi; i++ {
		for j := ilo; j < n; j++ {
			qaqij := qaq.Data[i*qaq.Stride+j]
			if j < i-1 {
				if math.Abs(qaqij) > 1e-13 {
					t.Errorf("%v: Q^T*AOrig*Q is not upper Hessenberg, [%v,%v]=%v", prefix, i, j, qaqij)
				}
				continue
			}
			diff := qaqij - a.Data[i*a.Stride+j]
			if math.Abs(diff) > 1e-13 {
				t.Errorf("%v: Q^T*AOrig*Q and A are not equal, diff at [%v,%v]=%v", prefix, i, j, diff)
			}
		}
	}
}
