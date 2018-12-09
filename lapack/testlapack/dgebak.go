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

type Dgebaker interface {
	Dgebak(job lapack.BalanceJob, side lapack.EVSide, n, ilo, ihi int, scale []float64, m int, v []float64, ldv int)
}

func DgebakTest(t *testing.T, impl Dgebaker) {
	rnd := rand.New(rand.NewSource(1))

	for _, job := range []lapack.BalanceJob{lapack.BalanceNone, lapack.Permute, lapack.Scale, lapack.PermuteScale} {
		for _, side := range []lapack.EVSide{lapack.EVLeft, lapack.EVRight} {
			for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 10, 18, 31, 53} {
				for _, extra := range []int{0, 11} {
					for cas := 0; cas < 100; cas++ {
						m := rnd.Intn(n + 1)
						v := randomGeneral(n, m, m+extra, rnd)
						var ilo, ihi int
						if v.Rows > 0 {
							ihi = rnd.Intn(n)
							ilo = rnd.Intn(ihi + 1)
						} else {
							ihi = -1
						}
						testDgebak(t, impl, job, side, ilo, ihi, v, rnd)
					}
				}
			}
		}
	}
}

func testDgebak(t *testing.T, impl Dgebaker, job lapack.BalanceJob, side lapack.EVSide, ilo, ihi int, v blas64.General, rnd *rand.Rand) {
	const tol = 1e-15
	n := v.Rows
	m := v.Cols
	extra := v.Stride - v.Cols

	// Create D and D^{-1} by generating random scales between ilo and ihi.
	d := eye(n, n)
	dinv := eye(n, n)
	scale := nanSlice(n)
	if job == lapack.Scale || job == lapack.PermuteScale {
		if ilo == ihi {
			scale[ilo] = 1
		} else {
			for i := ilo; i <= ihi; i++ {
				scale[i] = 2 * rnd.Float64()
				d.Data[i*d.Stride+i] = scale[i]
				dinv.Data[i*dinv.Stride+i] = 1 / scale[i]
			}
		}
	}

	// Create P by generating random column swaps.
	p := eye(n, n)
	if job == lapack.Permute || job == lapack.PermuteScale {
		// Make up some random permutations.
		for i := n - 1; i > ihi; i-- {
			scale[i] = float64(rnd.Intn(i + 1))
			blas64.Swap(blas64.Vector{N: n, Data: p.Data[i:], Inc: p.Stride},
				blas64.Vector{N: n, Data: p.Data[int(scale[i]):], Inc: p.Stride})
		}
		for i := 0; i < ilo; i++ {
			scale[i] = float64(i + rnd.Intn(ihi-i+1))
			blas64.Swap(blas64.Vector{N: n, Data: p.Data[i:], Inc: p.Stride},
				blas64.Vector{N: n, Data: p.Data[int(scale[i]):], Inc: p.Stride})
		}
	}

	got := cloneGeneral(v)
	impl.Dgebak(job, side, n, ilo, ihi, scale, m, got.Data, got.Stride)

	prefix := fmt.Sprintf("Case job=%v, side=%v, n=%v, ilo=%v, ihi=%v, m=%v, extra=%v",
		job, side, n, ilo, ihi, m, extra)

	if !generalOutsideAllNaN(got) {
		t.Errorf("%v: out-of-range write to V\n%v", prefix, got.Data)
	}

	// Compute D*V or D^{-1}*V and store into dv.
	dv := zeros(n, m, m)
	if side == lapack.EVRight {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, d, v, 0, dv)
	} else {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, dinv, v, 0, dv)
	}
	// Compute P*D*V or P*D^{-1}*V and store into want.
	want := zeros(n, m, m)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, p, dv, 0, want)

	if !equalApproxGeneral(want, got, tol) {
		t.Errorf("%v: unexpected value of V", prefix)
	}
}
