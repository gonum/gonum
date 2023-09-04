// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dggbaler interface {
	Dggbal(job lapack.BalanceJob, n int, a []float64, lda int, b []float64, ldb int, lscale, rscale, work []float64) (ilo, ihi int)
}

func DggbalTest(t *testing.T, impl Dggbaler) {
	rnd := rand.New(rand.NewSource(2))
	const extra = 0
	for _, job := range []lapack.BalanceJob{lapack.PermuteScale, lapack.BalanceNone, lapack.Scale, lapack.Scale} {
		for _, n := range []int{2, 3, 4, 5, 6, 10, 18, 31, 53, 100} {
			for _, lda := range []int{n + extra, n} {
				for _, ldb := range []int{n + extra, n} {
					for cas := 0; cas < 10; cas++ {
						testDggbal(t, rnd, impl, job, n, lda, ldb)
					}
				}
			}
		}
	}
}

func testDggbal(t *testing.T, rnd *rand.Rand, impl Dggbaler, job lapack.BalanceJob, n, lda, ldb int) {
	a := unbalancedSparseGeneral(n, n, lda, 2*n, rnd)
	b := unbalancedSparseGeneral(n, n, ldb, 2*n, rnd)
	extra := a.Stride - n

	var rscale, lscale []float64
	if n > 0 {
		rscale = nanSlice(n)
		lscale = nanSlice(n)
	}
	lwork := 1
	if n != 0 && (job == lapack.Scale || job == lapack.PermuteScale) {
		lwork = 6 * n
	}
	work := nanSlice(lwork)
	aCopy := cloneGeneral(a)
	bCopy := cloneGeneral(b)

	ilo, ihi := impl.Dggbal(job, n, a.Data, a.Stride, b.Data, b.Stride, lscale, rscale, work)

	prefix := fmt.Sprintf("Case job=%c, n=%v, extra=%v", job, n, extra)

	if !generalOutsideAllNaN(a) {
		t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
	}

	if n == 0 {
		if ilo != 0 {
			t.Errorf("%v: unexpected ilo when n=0. Want 0, got %v", prefix, ilo)
		}
		if ihi != -1 {
			t.Errorf("%v: unexpected ihi when n=0. Want -1, got %v", prefix, ihi)
		}
		return
	}
	testMatrixBalancing(t, a, aCopy, job, ilo, ihi, lscale)
	testMatrixBalancing(t, b, bCopy, job, ilo, ihi, rscale)
}

func testMatrixBalancing(t *testing.T, a, aBeforeBalance blas64.General, job lapack.BalanceJob, ilo, ihi int, scale []float64) {
	n := a.Rows
	prefix := fmt.Sprintf("Case job=%c, n=%d, lda=%d ilo=%v, ihi=%v", job, n, a.Stride, ilo, ihi)

	if job == lapack.BalanceNone {
		if ilo != 0 {
			t.Errorf("%v: unexpected ilo when job=BalanceNone. Want 0, got %v", prefix, ilo)
		}
		if ihi != n-1 {
			t.Errorf("%v: unexpected ihi when job=BalanceNone. Want %v, got %v", prefix, n-1, ihi)
		}
		k := -1
		for i := range scale {
			if scale[i] != 1 {
				k = i
				break
			}
		}
		if k != -1 {
			t.Errorf("%v: unexpected lscale[%v] when job=BalanceNone. Want 1, got %v", prefix, k, scale[k])
		}
		if !equalApproxGeneral(a, aBeforeBalance, 0) {
			t.Errorf("%v: unexpected modification of A when job=BalanceNone", prefix)
		}

		return
	}

	if ilo < 0 || ihi < ilo || n <= ihi {
		t.Errorf("%v: invalid ordering of ilo=%v and ihi=%v", prefix, ilo, ihi)
	}

	if ilo >= 2 && !isUpperTriangular(blas64.General{Rows: ilo - 1, Cols: ilo - 1, Data: a.Data, Stride: a.Stride}) {
		t.Errorf("%v: T1 is not upper triangular", prefix)
	}
	m := n - ihi - 1 // Order of T2.
	k := ihi + 1
	if m >= 2 && !isUpperTriangular(blas64.General{Rows: m, Cols: m, Data: a.Data[k*a.Stride+k:], Stride: a.Stride}) {
		t.Errorf("%v: T2 is not upper triangular", prefix)
	}
}
