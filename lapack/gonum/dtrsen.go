// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

// Dtrsen reorders the real Schur factorization of a real matrix A = Q*T*Qᵀ, so
// that a selected cluster of eigenvalues appears in the leading diagonal
// blocks of the upper quasi-triangular matrix T, and the matrix Q of Schur
// vectors is updated.
//
// The ijob parameter specifies what is computed:
//   - ijob=0: reorder only, no condition estimates
//   - ijob=1: compute reciprocal condition number for average of selected eigenvalues
//   - ijob=2: compute reciprocal condition number for right invariant subspace
//   - ijob=3: compute both 1 and 2
//
// If wantq is true, the matrix Q of Schur vectors is updated.
//
// selected specifies which eigenvalues are to be moved to the top-left of T.
// For a complex conjugate pair, either or both corresponding elements may be
// true; the pair is always moved together.
//
// On return, m is the dimension of the specified invariant subspace.
// s is the reciprocal condition number for the average of the selected
// eigenvalues (valid only if ijob=1 or 3).
// sep is the reciprocal condition number for the right invariant subspace
// (valid only if ijob=2 or 3).
//
// work must have length at least lwork. if lwork is -1, a workspace query is
// performed.
//
// iwork must have length at least liwork. if liwork is -1, a workspace query is
// performed.
//
// Dtrsen returns ok=false if the reordering failed because some eigenvalues are
// too close to swap.
func (impl Implementation) Dtrsen(ijob int, wantq bool, selected []bool, n int, t []float64, ldt int, q []float64, ldq int, wr, wi []float64, work []float64, lwork int, iwork []int, liwork int) (m int, s, sep float64, ok bool) {
	switch {
	case ijob < 0 || ijob > 3:
		panic(badIJob)
	case n < 0:
		panic(nLT0)
	case ldt < max(1, n):
		panic(badLdT)
	case wantq && ldq < max(1, n):
		panic(badLdQ)
	}

	// Workspace requirements.
	var lwmin, liwmin int
	msel := 0
	if n == 0 {
		lwmin = 1
		liwmin = 1
	} else {
		if ijob > 0 || (lwork != -1 && liwork != -1) {
			switch {
			case len(selected) < n:
				panic(badLenSelected)
			case len(t) < (n-1)*ldt+n:
				panic(shortT)
			}
			for k := 0; k < n; {
				if k < n-1 && t[(k+1)*ldt+k] != 0 {
					if selected[k] || selected[k+1] {
						msel += 2
					}
					k += 2
					continue
				}
				if selected[k] {
					msel++
				}
				k++
			}
		}
		switch ijob {
		case 0:
			lwmin = max(1, n)
			liwmin = 1
		case 1:
			lwmin = max(1, n, msel*(n-msel))
			liwmin = 1
		case 2, 3:
			lwmin = max(1, n, 2*msel*(n-msel))
			liwmin = max(1, msel*(n-msel))
		}
	}

	if lwork == -1 || liwork == -1 {
		if lwork == -1 {
			work[0] = float64(lwmin)
		}
		if liwork == -1 {
			iwork[0] = liwmin
		}
		return msel, 0, 0, true
	}

	if n == 0 {
		work[0] = float64(lwmin)
		iwork[0] = liwmin
		return 0, 1, 0, true
	}

	switch {
	case len(selected) < n:
		panic(badLenSelected)
	case len(t) < (n-1)*ldt+n:
		panic(shortT)
	case wantq && len(q) < (n-1)*ldq+n:
		panic(shortQ)
	case len(wr) < n:
		panic(badLenWr)
	case len(wi) < n:
		panic(badLenWi)
	case lwork < lwmin:
		panic(badLWork)
	case len(work) < lwork:
		panic(shortWork)
	case liwork < liwmin:
		panic(badLWork)
	case len(iwork) < liwork:
		panic(shortIWork)
	}
	defer func() {
		work[0] = float64(lwmin)
		iwork[0] = liwmin
	}()

	ok = true
	m = msel
	ks := 0
	pair := false
	for k := 0; k < n; k++ {
		if pair {
			pair = false
			continue
		}
		n2 := 0
		if k < n-1 && t[(k+1)*ldt+k] != 0 {
			n2 = 1
		}

		swap := selected[k]
		if n2 == 1 {
			swap = swap || selected[k+1]
		}
		if swap {
			if k > ks {
				compq := lapack.UpdateSchurNone
				if wantq {
					compq = lapack.UpdateSchur
				}
				_, _, swapOk := impl.Dtrexc(compq, n, t, ldt, q, ldq, k, ks, work)
				if !swapOk {
					ok = false
					break
				}
			}
			if n2 == 1 {
				ks += 2
			} else {
				ks++
			}
		}
		if n2 == 1 {
			pair = true
		}
	}

	// Update wr and wi.
	for i := 0; i < n; i++ {
		if i < n-1 && t[(i+1)*ldt+i] != 0 {
			wr[i] = t[i*ldt+i]
			wr[i+1] = t[i*ldt+i]
			wi[i] = math.Sqrt(math.Abs(t[i*ldt+i+1])) * math.Sqrt(math.Abs(t[(i+1)*ldt+i]))
			wi[i+1] = -wi[i]
			i++
		} else {
			wr[i] = t[i*ldt+i]
			wi[i] = 0
		}
	}

	s = 1
	sep = 0
	wants := ijob == 1 || ijob == 3
	wantsp := ijob == 2 || ijob == 3
	if (m == 0 || m == n) && wantsp {
		sep = impl.Dlange(lapack.MaxColumnSum, n, n, t, ldt, work)
	}
	if ijob == 0 || m == 0 || m == n || !ok {
		return m, s, sep, ok
	}

	n1 := m
	n2 := n - m
	nn := n1 * n2

	if wants {
		impl.Dlacpy(blas.All, n1, n2, t[n1:], ldt, work, n2)
		scale, _ := impl.Dtrsyl(blas.NoTrans, blas.NoTrans, -1, n1, n2,
			t, ldt, t[n1*ldt+n1:], ldt, work, n2)
		rnorm := impl.Dlange(lapack.Frobenius, n1, n2, work, n2, nil)
		if rnorm == 0 {
			s = 1
		} else {
			s = scale / (math.Sqrt(scale*scale/rnorm+rnorm) * math.Sqrt(rnorm))
		}
	}

	if wantsp {
		est := 0.0
		kase := 0
		isave := [3]int{}
		scale := 0.0
		for {
			est, kase = impl.Dlacn2(nn, work[nn:], work, iwork, est, kase, &isave)
			if kase == 0 {
				break
			}
			if kase == 1 {
				scale, _ = impl.Dtrsyl(blas.NoTrans, blas.NoTrans, -1, n1, n2,
					t, ldt, t[n1*ldt+n1:], ldt, work, n2)
			} else {
				scale, _ = impl.Dtrsyl(blas.Trans, blas.Trans, -1, n1, n2,
					t, ldt, t[n1*ldt+n1:], ldt, work, n2)
			}
		}
		sep = scale / est
	}

	return m, s, sep, ok
}
