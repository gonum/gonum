// Copyright ©2023 The Gonum Authors. All rights reserved.
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
	if job == lapack.BalanceNone {
		return
	}
	want := cloneGeneral(aCopy)
	// LSCALE is DOUBLE PRECISION array, dimension (N)
	// Details of the permutations and scaling factors applied
	// to the left side of A and B.  If P(j) is the index of the
	// row interchanged with row j, and D(j)
	// is the scaling factor applied to row j, then
	//    LSCALE(j) = P(j)    for J = 1,...,ILO-1
	//              = D(j)    for J = ILO,...,IHI
	//              = P(j)    for J = IHI+1,...,N.
	// The order in which the interchanges are made is N to IHI+1,
	// then 1 to ILO-1.

	if job == lapack.Permute || job == lapack.PermuteScale {
		// Create the left permutation matrix Pl.
		pl := eye(n, n)
		for j := n - 1; j > ihi; j-- {
			blas64.Swap(blas64.Vector{N: n, Data: pl.Data[j:], Inc: pl.Stride},
				blas64.Vector{N: n, Data: pl.Data[int(lscale[j]):], Inc: pl.Stride})
		}
		for j := 0; j < ilo; j++ {
			blas64.Swap(blas64.Vector{N: n, Data: pl.Data[j:], Inc: pl.Stride},
				blas64.Vector{N: n, Data: pl.Data[int(lscale[j]):], Inc: pl.Stride})
		}

		// Create the right permutation matrix Pr.
		pr := eye(n, n)
		for j := n - 1; j > ihi; j-- {
			blas64.Swap(blas64.Vector{N: n, Data: pr.Data[j:], Inc: pr.Stride},
				blas64.Vector{N: n, Data: pr.Data[int(rscale[j]):], Inc: pr.Stride})
		}
		for j := 0; j < ilo; j++ {
			blas64.Swap(blas64.Vector{N: n, Data: pr.Data[j:], Inc: pr.Stride},
				blas64.Vector{N: n, Data: pr.Data[int(rscale[j]):], Inc: pr.Stride})
		}

		// Compute Plᵀ*A*Pl and store into want.
		ap := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, want, pl, 0, ap)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, pl, ap, 0, want)
	}

	if job == lapack.Scale || job == lapack.PermuteScale {
		return // TODO(soypat): Test this case!
		// Modify want by Dl and Dl^{-1}.
		dl := eye(n, n)
		dlinv := eye(n, n)
		for i := ilo; i <= ihi; i++ {
			dl.Data[i*dl.Stride+i] = lscale[i]
			dlinv.Data[i*dlinv.Stride+i] = 1 / lscale[i]
		}
		dr := eye(n, n)
		drinv := eye(n, n)
		for i := ilo; i <= ihi; i++ {
			dr.Data[i*dr.Stride+i] = rscale[i]
			drinv.Data[i*drinv.Stride+i] = 1 / rscale[i]
		}
		ad := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, want, dl, 0, ad)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, dlinv, ad, 0, want)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, want, dr, 0, ad)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, drinv, ad, 0, want)
	}
	if !equalApproxGeneral(want, a, 1e-5) {
		t.Errorf("%v: unexpected value of A, ilo=%v, ihi=%v", prefix, ilo, ihi)
	}
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
