// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"math"

	"github.com/gonum/lapack"
)

// Dgesvd computes the singular value decomposition of the input matrix A.
// The singular value decomposition is
//  A = U * Sigma * V^T
// where Sigma is an m×n diagonal matrix containing the singular values of A,
// U is an m×m orthogonal matrix and V is an n×n orthogonal matrix. The first
// min(m,n) columns of U and V are the left and right singular vectors of A
// respectively.
//
// jobU and jobVT are options for computing the singular vectors. The behavior
// is as follows
//  jobU == lapack.SVDAll		All M columns of U are returned in u
//  jobU == lapack.SVDInPlace	The first min(m,n) columns are returned in u
//  jobU == lapack.SVDOverwrite	The first min(m,n) columns of U are written into a
//	jobU == lapack.SVDNone		The columns of U are not computed.
// The behavior is the same for jobVT and the rows of V^T. At most one of jobU
// and jobVT can equal lapack.SVDOverwrite.
//
// On entry, a contains the data for the m×n matrix A. During the call to Dgesvd
// the data is overwritten. On exit, A contains the appropriate singular vectors
// if either job is lapack.SVDOverwrite.
//
// s is a slice of length at least min(m,n) and on exit contains the singular
// values in decreasing order.
//
// u contains the left singular vectors on exit, stored columnwise. If
// jobU == lapack.SVDAll, u is of size m×m. If jobU == lapack.SVDInPlace u is
// of size m×min(m,n). If jobU == lapack.SVDOverwrite or lapack.SVDNone, u is
// not used.
//
// vt contains the left singular vectors on exit, stored rowwise. If
// jobV == lapack.SVDAll, vt is of size n×m. If jobV == lapack.SVDInPlace vt is
// of size min(m,n)×n. If jobU == lapack.SVDOverwrite or lapack.SVDNone, vt is
// not used.
//
// work is a slice for storing temporary memory, and lwork is the usable size of
// the slice. lwork must be at least max(5*min(m,n), 3*min(m,n)+max(m,n)).
// If lwork == -1, instead of performing Dgesvd, the optimal work length will be
// stored into work[0]. Dgesvd will panic if the working memory has insufficient
// storage.
//
// Dgesvd returns whether the decomposition successfully completed.
func (impl Implementation) Dgesvd(jobU, jobVT lapack.SVDJob, m, n int, a []float64, lda int, s, u []float64, ldu int, vt []float64, ldvt int, work []float64, lwork int) (ok bool) {
	checkMatrix(m, n, a, lda)
	if jobU == lapack.SVDAll {
		checkMatrix(m, m, u, ldu)
	} else if jobU == lapack.SVDInPlace {
		checkMatrix(m, min(m, n), u, ldu)
	}
	if jobVT == lapack.SVDAll {
		checkMatrix(n, n, vt, ldvt)
	} else if jobVT == lapack.SVDInPlace {
		checkMatrix(min(m, n), n, vt, ldvt)
	}
	if jobU == lapack.SVDOverwrite && jobVT == lapack.SVDOverwrite {
		panic("lapack: both jobU and jobV are lapack.SVDOverwrite")
	}
	if len(s) < min(m, n) {
		panic(badS)
	}
	minWork := max(5*min(m, n), 3*min(m, n)+max(m, n))
	if lwork != -1 {
		if len(work) < lwork {
			panic(badWork)
		}
		if lwork < minWork {
			panic(badWork)
		}
	}
	if m == 0 || n == 0 {
		return true
	}

	minmn := min(m, n)

	wantua := jobU == lapack.SVDAll
	wantus := jobU == lapack.SVDInPlace
	wantuas := wantua || wantus
	wantuo := jobU == lapack.SVDOverwrite
	wantun := jobU == lapack.None

	wantva := jobVT == lapack.SVDAll
	wantvs := jobVT == lapack.SVDInPlace
	wantvas := wantva || wantvs
	wantvo := jobVT == lapack.SVDOverwrite
	wantvn := jobVT == lapack.None

	var mnthr int
	dum := []float64{0}
	// Compute optimal space for subroutines.
	maxwrk := 1
	opts := string(jobU) + string(jobVT)
	if m >= n {
		mnthr = impl.Ilaenv(6, "DGESVD", opts, m, n, 0, 0)
		bdspac := 5 * n
		impl.Dgeqrf(m, n, a, lda, dum, dum, -1)
		lwork_dgeqrf := int(dum[0])
		impl.Dorgqr(m, n, n, a, lda, dum, dum, -1)
		lwork_dorgqr_n := int(dum[0])
		impl.Dorgqr(m, m, n, a, lda, dum, dum, -1)
		lwork_dorgqr_m := int(dum[0])
		impl.Dgebrd(n, n, a, lda, s, dum, dum, dum, dum, -1)
		lwork_dgebrd := int(dum[0])
		impl.Dorgbr(lapack.ApplyP, n, n, n, a, lda, dum, dum, -1)
		lwork_dorgbr_p := int(dum[0])
		impl.Dorgbr(lapack.ApplyQ, n, n, n, a, lda, dum, dum, -1)
		lwork_dorgbr_q := int(dum[0])

		if m >= mnthr {
			// m >> n
			if wantun {
				// Path 1
				maxwrk = n + lwork_dgeqrf
				maxwrk = max(maxwrk, 3*n+lwork_dgebrd)
				if wantvo || wantvas {
					maxWork = max(maxwrk, 3*n+lwork_dorgbr_p)
				}
				maxwrk = max(maxwrk, bdspac)
			} else if wantuo && wantvn {
				// Path 2
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_n)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = max(n*n+wrkbl, n*n+m*n+n)
			} else if wantuo && wantvs {
				// Path 3
				// or lapack.All
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_n)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = max(n*n+wrkbl, n*n+m*n+n)
			} else if wantus && wantvn {
				// Path 4
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_n)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = n*n + wrkbl
			} else if wantus && wantvo {
				// Path 5
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_n)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = 2*n*n + wrkbl
			} else if wantus && wantvas {
				// Path 6
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_n)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = n*n + wrkbl
			} else if wantua && wantvn {
				// Path 7
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_m)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = n*n + wrkbl
			} else if wantua && wantvo {
				// Path 8
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_m)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = 2*n*n + wrkbl
			} else if wantua && wantvas {
				// Path 9
				wrkbl := n + lwork_dgeqrf
				wrkbl = max(wrkbl, n+lwork_dorgqr_m)
				wrkbl = max(wrkbl, 3*n+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_q)
				wrkbl = max(wrkbl, 3*n+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = n*n + wrkbl
			}
		} else {
			// Path 10: m > n
			impl.Dgebrd(m, n, a, lda, s, dum, dum, dum, dum, -1)
			lwork_dgebrd := int(dum[0])
			maxwrk = 3*n + lwork_dgebrd
			if wantus || wantuo {
				impl.Dorgbr(lapack.ApplyQ, m, n, n, a, lda, dum, dum, -1)
				lwork_dorgbr_q = int(dum[0])
				maxwrk = max(maxwrk, 3*n+lwork_dorgbr_q)
			}
			if wantua {
				impl.Dorgbr(lapack.ApplyQ, m, m, n, a, lda, dum, dum, -1)
				lwork_dorgbr_q := int(dum[0])
				maxwrk = max(maxwrk, 3*n+lwork_dorgbr_p)
			}
			if !wantvn {
				maxwrk = max(maxwrk, 3*n+lwork_dorgbr_p)
			}
			maxwrk = max(maxwrk, bdspac)
		}
	} else {
		mnthr = impl.Ilaenv(6, "DGESVD", opts, m, n, 0, 0)
		bdspac := 5 * m
		impl.Dgelqf(m, n, a, lda, dum, dum, -1)
		lwork_dgelqf := int(dum[0])
		impl.Dorglq(n, n, m, dum, n, dum, dum, -1)
		lwork_dorglq_n := int(dum[0])
		impl.Dorglq(m, n, m, a, lda, dum, dum, -1)
		lwork_dorglq_m := int(dum[0])
		impl.Dgebrd(m, m, a, lda, s, dum, dum, dum, dum, -1)
		lwork_dgebrd := int(dum[0])
		impl.Dorgbr(lapack.ApplyP, m, m, m, a, n, dum, dum, -1)
		lwork_dorgbr_p := int(dum[0])
		impl.Dorgbr(lapack.ApplyQ, m, m, m, a, n, dum, dum, -1)
		lwork_dorgbr_q := int(dum[0])
		if n >= mnthr {
			// n >> m
			if wantvn {
				// Path 1t
				maxwrk = m + lwork_dgelqf
				maxwrk = max(maxwrk, 3*m+lwork_dgebrd)
				if wntuo.OR.wntuas {
					maxwrk = max(maxwrk, 3*m+lwork_dorgbr_q)
				}
				maxwrk = max(maxwrk, bdspac)
			} else if wantvo && wantun {
				// Path 2t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_m)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = max(m*m+wrkbl, m*m+m*n+m)
			} else if wantvo && wantuas {
				// Path 3t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_m)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = max(m*m+wrkbl, m*m+m*n+m)
			} else if wantvs && wantun {
				// Path 4t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_m)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = m*m + wrkbl
			} else if wantvs && wantuo {
				// Path 5t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_m)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = 2*m*m + wrkbl
			} else if wantvs && wantuas {
				// Path 6t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_m)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = m*m + wrkbl
			} else if wantva && wantun {
				// Path 7t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_n)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = m*m + wrkbl
			} else if wantva && wantuo {
				// Path 8t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_n)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = 2*m*m + wrkbl
			} else if wantva && wantuas {
				// Path 9t
				wrkbl := m + lwork_dgelqf
				wrkbl = max(wrkbl, m+lwork_dorglq_n)
				wrkbl = max(wrkbl, 3*m+lwork_dgebrd)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_p)
				wrkbl = max(wrkbl, 3*m+lwork_dorgbr_q)
				wrkbl = max(wrkbl, bdspac)
				maxwrk = m*m + wrkbl
			}
		} else {
			// Path 10t, n > m
			impl.Dgebrd(m, n, a, lda, s, dum, dum, dum, dum, -1)
			lwork_dgebrd = int(dum[0])
			maxwrk := 3*m + lwork_dgebrd
			if wantvs || wantvo {
				impl.Dorgbr(lapack.ApplyP, m, n, m, a, n, dum, dum, -1)
				lwork_dorgbr_p = int(dum[0])
				maxwrk = max(maxwrk, 3*m+lwork_dorgbr_p)
			}
			if wantva {
				impl.Dorgbr(lapack.ApplyP, n, n, m, a, n, dum, dum, -1)
				lwork_dorgbr_p = int(dum[0])
				maxwrk = max(maxwrk, 3*m+lwork_dorgbr_p)
			}
			if !wantun {
				maxwrk = max(maxwrk, 3*m+lwork_dorgbr_q)
			}
			maxwrk = max(maxwrk, bdspac)
		}
	}
	maxwrk = max(maxwrk, minWork)
	work[0] = maxwrk
	if lwork == -1 {
		return true
	}

	// Perform decomposition.
	eps := dlamchE
	smlnum := math.Sqrt(dlamchS) / eps
	bignum := 1 / smlnum

	// Scale A if max element outside range [smlnum, bignum]
	anrm := impl.Dlange(lapack.MaxAbs, m, n, a, lda, dum)
	iscl := 0
	if anrm > 0 && anrm < smlnum {
		iscl = 1
		impl.Dlascl(lapack.General, 0, 0, anrm, smlnum, m, n, a, lda)
	} else if anrm > bignum {
		iscl = 1
		impl.Dlascl(lapack.General, 0, 0, anrm, bignum, m, n, a, lda)
	}

	// Line 671
	if m >= n {
		// If A has sufficiently more rows than columns, use the QR decomposition.
		if m >= mnthr {
			if wantun {
				// Path 1
				itau = 1
				iwo
			}
		}
	}
}
