// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

// Dgees computes for an n×n real nonsymmetric matrix A, the eigenvalues, the
// real Schur form T, and, optionally, the matrix of Schur vectors Z. This
// gives the Schur factorization
//
//	A = Z * T * Zᵀ
//
// where T is upper quasi-triangular (the Schur form). A matrix is in real Schur
// form if it is upper quasi-triangular with 1×1 and 2×2 blocks. 2×2 blocks will
// be standardized in the form
//
//	[ a  b ]
//	[ c  a ]
//
// where b*c < 0. The eigenvalues of such a block are a ± sqrt(bc).
//
// Optionally, the routine can order the eigenvalues on the diagonal of the real
// Schur form so that selected eigenvalues are at the top left. The leading
// columns of Z then form an orthonormal basis for the invariant subspace
// corresponding to the selected eigenvalues.
//
// A matrix is in Schur canonical form if it is quasi-triangular and every 2×2
// diagonal block has the form shown above.
//
// If jobvs is lapack.SchurNone, no Schur vectors are computed.
// If jobvs is lapack.SchurHess, the Schur vectors are computed and returned in
// vs.
// For other values of jobvs, Dgees will panic.
//
// If sort is lapack.SortNone, eigenvalues are not ordered.
// If sort is lapack.SortSelected, eigenvalues are reordered so that selected
// eigenvalues appear in the leading diagonal blocks of the Schur form.
//
// selctg is a function that selects eigenvalues. It is only used when sort is
// lapack.SortSelected. An eigenvalue wr[j]+i*wi[j] is selected if
// selctg(wr[j], wi[j]) is true. Note that a selected complex eigenvalue may no
// longer satisfy selctg after ordering, since ordering may change the value of
// complex eigenvalues (especially if the eigenvalue is ill-conditioned).
//
// On entry, a contains the n×n matrix A. On return, a has been overwritten by
// its real Schur form T.
//
// wr and wi will contain the real and imaginary parts, respectively, of the
// computed eigenvalues in the same order that they appear on the diagonal of
// the output Schur form T. Complex conjugate pairs of eigenvalues will appear
// consecutively with the eigenvalue having positive imaginary part first.
// wr and wi must have length n, otherwise Dgees will panic.
//
// If jobvs is lapack.SchurHess, vs must have leading dimension ldvs >= n and
// have length at least (n-1)*ldvs+n, and on return will contain the orthogonal
// matrix Z of Schur vectors.
// If jobvs is lapack.SchurNone, vs is not referenced.
//
// work must have length at least lwork and lwork must be at least max(1,3*n).
// For good performance, lwork should generally be larger. On return, work[0]
// will contain the optimal value of lwork.
//
// If lwork is -1, instead of performing Dgees, only the optimal value of lwork
// is computed and stored into work[0].
//
// If sort is lapack.SortSelected, bwork must have length n; otherwise bwork is
// not referenced and can be nil.
//
// On return, sdim contains the number of eigenvalues (after sorting) for which
// selctg is true. If sort is lapack.SortNone, sdim will be 0.
//
// On return, ok reports whether all eigenvalues have been computed. If ok is
// false, a has been partially reduced and some eigenvalues may not have
// converged. In this case the eigenvalues in wr[0:n] and wi[0:n] are correct
// for indices where they were computed, but T may not be in Schur form.
//
// Dgees is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dgees(jobvs lapack.SchurComp, sort lapack.SchurSort, selctg func(wr, wi float64) bool, n int, a []float64, lda int, wr, wi []float64, vs []float64, ldvs int, work []float64, lwork int, bwork []bool) (sdim int, ok bool) {
	wantvs := jobvs == lapack.SchurHess
	wantst := sort == lapack.SortSelected

	minwrk := max(1, 3*n)

	switch {
	case jobvs != lapack.SchurHess && jobvs != lapack.SchurNone:
		panic(badSchurComp)
	case sort != lapack.SortNone && sort != lapack.SortSelected:
		panic(badSchurSort)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	case ldvs < 1 || (wantvs && ldvs < n):
		panic(badLdZ)
	case lwork < minwrk && lwork != -1:
		panic(badLWork)
	case len(work) < max(1, lwork):
		panic(shortWork)
	}

	// Compute workspace.
	var maxwrk int
	if n == 0 {
		maxwrk = 1
	} else {
		maxwrk = 2*n + n*impl.Ilaenv(1, "DGEHRD", " ", n, 1, n, 0)
		if wantvs {
			maxwrk = max(maxwrk, 2*n+(n-1)*impl.Ilaenv(1, "DORGHR", " ", n, 1, n, -1))
			impl.Dhseqr(lapack.EigenvaluesAndSchur, lapack.SchurHess, n, 0, n-1,
				nil, n, nil, nil, nil, max(1, ldvs), work, -1)
		} else {
			impl.Dhseqr(lapack.EigenvaluesAndSchur, lapack.SchurNone, n, 0, n-1,
				nil, n, nil, nil, nil, 1, work, -1)
		}
		hswork := int(work[0])

		trsenwork := 0
		if wantst {
			var iwork [1]int
			impl.Dtrsen(0, wantvs, nil, n, nil, n, nil, n, nil, nil, work, -1, iwork[:], -1)
			trsenwork = int(work[0])
		}
		maxwrk = max(maxwrk, max(1, 2*n+hswork), 2*n+trsenwork)
	}

	if lwork == -1 {
		work[0] = float64(maxwrk)
		return 0, true
	}

	// Quick return if possible.
	if n == 0 {
		work[0] = 1
		return 0, true
	}

	switch {
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case len(wr) != n:
		panic(badLenWr)
	case len(wi) != n:
		panic(badLenWi)
	case wantvs && len(vs) < (n-1)*ldvs+n:
		panic(shortZ)
	case wantst && len(bwork) < n:
		panic("lapack: insufficient length of bwork")
	case wantst && selctg == nil:
		panic("lapack: selctg is nil but sort is SortSelected")
	}

	// Get machine constants.
	eps := dlamchP
	smlnum := math.Sqrt(dlamchS) / eps
	bignum := 1 / smlnum

	// Scale A if max element outside range [smlnum,bignum].
	anrm := impl.Dlange(lapack.MaxAbs, n, n, a, lda, nil)
	var scalea bool
	var cscale float64
	if anrm > 0 && anrm < smlnum {
		scalea = true
		cscale = smlnum
	} else if anrm > bignum {
		scalea = true
		cscale = bignum
	}
	if scalea {
		impl.Dlascl(lapack.General, 0, 0, anrm, cscale, n, n, a, lda)
	}

	// Balance the matrix.
	// Use permutation-only balancing (not scaling) because scaling would
	// affect the Schur vectors. This matches LAPACK's DGEES.
	workbal := work[:n]
	ilo, ihi := impl.Dgebal(lapack.Permute, n, a, lda, workbal)

	// Reduce to upper Hessenberg form.
	iwrk := 2 * n
	tau := work[n : iwrk-1]
	impl.Dgehrd(n, ilo, ihi, a, lda, tau, work[iwrk:], lwork-iwrk)

	if wantvs {
		// Copy Householder vectors to vs.
		impl.Dlacpy(blas.Lower, n, n, a, lda, vs, ldvs)
		// Generate orthogonal matrix in vs.
		impl.Dorghr(n, ilo, ihi, vs, ldvs, tau, work[iwrk:], lwork-iwrk)
	}

	// Perform QR iteration, accumulating Schur vectors in vs if desired.
	iwrk = n
	var compz lapack.SchurComp
	if wantvs {
		compz = lapack.SchurOrig
	} else {
		compz = lapack.SchurNone
	}
	ieval := impl.Dhseqr(lapack.EigenvaluesAndSchur, compz, n, ilo, ihi,
		a, lda, wr, wi, vs, ldvs, work[iwrk:], lwork-iwrk)
	if ieval > 0 {
		ok = false
	} else {
		ok = true
	}

	// Sort eigenvalues if desired.
	sdim = 0
	if ok && wantst {
		// Populate bwork[] using selctg callback.
		for i := 0; i < n; {
			if i < n-1 && wi[i] != 0 {
				// Complex conjugate pair at positions i and i+1.
				// Both must be selected if either is selected.
				sel1 := selctg(wr[i], wi[i])
				sel2 := selctg(wr[i+1], wi[i+1])
				bwork[i] = sel1 || sel2
				bwork[i+1] = sel1 || sel2
				i += 2
			} else {
				// Real eigenvalue.
				bwork[i] = selctg(wr[i], wi[i])
				i++
			}
		}

		// Call Dtrsen to reorder eigenvalues.
		liwork := 1
		iwork := make([]int, liwork)

		var trsenM int
		var trsenOk bool
		trsenM, _, _, trsenOk = impl.Dtrsen(0, wantvs, bwork, n,
			a, lda, vs, ldvs, wr, wi, work[iwrk:], lwork-iwrk, iwork, liwork)

		if !trsenOk {
			ok = false
		}
		sdim = trsenM
	}

	if wantvs {
		// Undo balancing (permutation only).
		impl.Dgebak(lapack.Permute, lapack.EVRight, n, ilo, ihi, workbal, n, vs, ldvs)
	}

	if scalea {
		// Undo scaling for the Schur form of A.
		impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n, n, a, lda)
		impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n, 1, wr, n)
		impl.Dlascl(lapack.General, 0, 0, cscale, anrm, n, 1, wi, n)
	}

	work[0] = float64(maxwrk)
	return sdim, ok
}
