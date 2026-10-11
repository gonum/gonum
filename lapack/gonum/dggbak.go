// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dggbak forms the right or left eigenvectors of a real generalized
// eigenvalue problem
//
//	A*x = lambda*B*x,
//
// by backward transformation on
// the computed eigenvectors of the balanced pair of matrices output by
// Dggbal.
//
// Parameters:
//   - job specifies the type of backward transformation required.
//     lapack.NoBalancing: do nothing, return immediately.
//     lapack.Permute: do backwards transformation for permutation only.
//     lapack.Scale: do backwards transformation for scaling only.
//     lapack.PermuteScale: do backwards transformations for both scaling and permutation.
//   - side specifies whether V contains right (lapack.EVRight) or left (lapack.EVLeft) eigenvectors.
//   - n is number of rows in V. m is number of columns in V.
//   - ilo and ihi determined by Dggbal. 0<=ilo<=ihi for n>0; ilo=0, ihi=-1 for n=0.
//   - v contains V matrix with right or left eigenvectors to be transformed as returned by Dtgevc
//     On exit v is overwritten by the transformed eigenvectors.
//   - lscale and rscale are of length n and contain details of permutations and/or scaling factors
//     applied to left and right sides of A and B respectively.
//
// Dggbal is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dggbak(job lapack.BalanceJob, side lapack.EVSide, n, ilo, ihi int, lscale, rscale []float64, m int, v []float64, ldv int) {
	switch {
	case job != lapack.BalanceNone && job != lapack.Permute && job != lapack.Scale && job != lapack.PermuteScale:
		panic(badBalanceJob)
	case side != lapack.EVRight && side != lapack.EVLeft:
		panic(badEVSide)
	case n < 0:
		panic(nLT0)
	case m < 0:
		panic(mLT0)
	case ilo < 0 || ilo > n-1:
		panic(badIlo)
	case ihi > n-1 || ihi < -1 || ((ihi == -1 || ihi < ilo) && n != 0):
		panic(badIhi)
	case ldv < max(1, m):
		panic(badLdV)
	case len(lscale) < n:
		panic(badLenScale)
	case len(rscale) < n:
		panic(badLenScale)
	case len(v) < (n-1)*ldv+m:
		panic(shortV)
	}

	if n == 0 || m == 0 || job == lapack.BalanceNone {
		return // Quick return.
	}

	bi := blas64.Implementation()
	if job == lapack.Scale || job == lapack.PermuteScale {
		// Backward balance.

		if side == lapack.EVRight {
			// Backward transformation on right eigenvectors.
			for i := ilo; i <= ihi; i++ {
				bi.Dscal(m, rscale[i], v[i*ldv:], 1)
			}
		}

		if side == lapack.EVLeft {
			// Backward transformation on left eigenvectors.
			for i := ilo; i <= ihi; i++ {
				bi.Dscal(m, lscale[i], v[i*ldv:], 1)
			}
		}
	}

	if job == lapack.Permute || job == lapack.PermuteScale {
		// Backward permutation.
		if side == lapack.EVRight {
			// Backward transformation on right eigenvectors.
			if ilo != 0 {
				for i := ilo - 1; i >= 0; i-- {
					k := int(rscale[i])
					if k != i {
						bi.Dswap(m, v[i*ldv:], 1, v[k*ldv:], 1)
					}
				}
			}
			if ihi != n-1 {
				for i := ihi + 1; i < n; i++ {
					k := int(rscale[i])
					if k != i {
						bi.Dswap(m, v[i*ldv:], 1, v[k*ldv:], 1)
					}
				}
			}
		}
		if side == lapack.EVLeft {
			// Backward transformation on left eigenvectors.
			if ilo != 0 {
				for i := ilo - 1; i >= 0; i-- {
					k := int(lscale[i])
					if k != i {
						bi.Dswap(m, v[i*ldv:], 1, v[k*ldv:], 1)
					}
				}
			}
			if ihi != n-1 {
				for i := ihi + 1; i < n; i++ {
					k := int(lscale[i])
					if k != i {
						bi.Dswap(m, v[i*ldv:], 1, v[k*ldv:], 1)
					}
				}
			}
		}
	}
}
