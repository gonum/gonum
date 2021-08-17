// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dggbal balances a pair of general real matrices (A,B).  This
// involves, first, permuting A and B by similarity transformations to
// isolate eigenvalues in the first 1 to ILO$-$1 and last IHI+1 to N
// elements on the diagonal; and second, applying a diagonal similarity
// transformation to rows and columns ILO to IHI to make the rows
// and columns as close in norm as possible. Both steps are optional.
// Balancing may reduce the 1-norm of the matrices, and improve the
// accuracy of the computed eigenvalues and/or eigenvectors in the
// generalized eigenvalue problem A*x = lambda*B*x.
//  n is the order of matrices A and B. n >= 0
//
// lscale, rscale must be of size n.
func (impl Implementation) Dggbal(job lapack.BalanceJob, n int, a []float64, lda int, b []float64, ldb int, lscale, rscale, work []float64) (ilo, ihi int) {
	var (
		alpha, pgamma, t, tc, sum, cmax, cor, basl float64
		i, j, kount, lm1, jp1, iflow, m, ip1       int // loop var
	)
	sclfac := 1.
	switch {
	case job != lapack.BalanceNone && job != lapack.Permute && job != lapack.PermuteScale && job != lapack.Scale:
		panic(badBalanceJob)
	case n < 0:
		panic(nLT0)
	}

	if lda < max(1, n) {
		panic(shortA)
	}
	if ldb < max(1, n) {
		panic(shortB)
	}

	// quick return if possible
	if n == 0 || n == 1 {
		ilo = 1
		ihi = n
		if n == 1 {
			rscale[0] = 1
			lscale[0] = 1
		}
		return
	}

	bi := blas64.Implementation()
	if job == lapack.BalanceNone {
		ilo = 1
		ihi = n
		for i = 0; i < n; i++ {
			lscale[i] = 1
			rscale[i] = 1
		}
		return
	}

	k := 1
	l := n
	if job == lapack.Scale {
		goto oneNinety
	}
	goto thirty
twenty:
	l = lm1
	if l == 1 {
		rscale[0] = 1
		lscale[0] = 1
		goto oneNinety
	}

thirty:
	lm1 = l - 1
eighty:
	for i = l; i >= 0; i-- {
		for j = 1; j < lm1; j++ {
			jp1 = j + 1
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				goto fifty
			}
		}
		j = l
		goto seventy
	fifty:
		for j = jp1; j < l; j++ {
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				continue eighty
			}
		}
		j = jp1 - 1
	seventy:
		m = l
		iflow = 1
		goto oneSixty
	}
	goto hundred

	// Find column with one nonzero in rows K through N.
ninety:
	k++
hundred:
	for j = k; j < l; j++ {
		for i = k; i < lm1; i++ {
			ip1++
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				goto oneTwenty
			}
		}
		i = l
		goto oneForty
	oneTwenty:
		for i = ip1; i < l; i++ {
			if a[i*lda+j] != 0 || b[i*lda+j] != 0 {
				continue hundred
			}
		}
		i = ip1 - 1
	oneForty:
		m = k
		iflow = 2
		goto oneSixty
	}
	goto oneNinety

	// Permute rows M and I
oneSixty:
	lscale[m] = float64(i)
	if i != m {
		bi.Dswap(n-k+1, a[i*lda+k:], 1, a[m*lda+k:], 1)
		bi.Dswap(n-k+1, b[i*ldb+k:], 1, b[m*ldb+k:], 1)
	}

	// Permute columns M and J
	rscale[m] = float64(j)
	if j != m {
		bi.Dswap(l, a[j:], lda, a[m:], lda)
		bi.Dswap(l, b[j:], ldb, b[m:], ldb)
	}
	switch iflow {
	case 1:
		goto twenty
	case 2:
		goto ninety
	}

oneNinety:
	ilo = k
	ihi = l
	if job == lapack.Permute {
		for i := ilo; i < ihi; i++ {
			lscale[i] = 1
			rscale[i] = 1
		}
		return
	}
	if ilo == ihi {
		return
	}

	// Balance the submatrix in rows ILO to IHI.
	nr := ihi - ilo + 1
	for i = ilo; i < ihi; i++ {
		rscale[i] = 0
		lscale[i] = 0

		work[i] = 0
		work[i+n] = 0
		work[i+2*n] = 0
		work[i+3*n] = 0
		work[i+4*n] = 0
		work[i+5*n] = 0
	}

	// Compute right side vector in resulting linear equations.
	basl = math.Log10(sclfac)
	for i = ilo; i < ihi; i++ {
		for j = ilo; j < ihi; j++ {
			tb := b[i*ldb+j]
			ta := a[i*lda+j]
			if ta != 0 {
				ta = math.Log10(math.Abs(ta)) / basl
			}
			if tb != 0 {
				tb = math.Log10(math.Abs(tb)) / basl
			}
			work[i+4*n] -= ta + tb
			work[j+5*n] -= ta + tb
		}
	}
	coef := 1 / float64(2*nr)
	coef2 := coef * coef
	coef5 := .5 * coef2
	nrp2 := nr + 2
	beta := 0.
	it := 1

	// Start generalized conjugate gradient iteration
twoFiddy:
	gamma := bi.Ddot(nr, work[ilo+4*n:], 1, work[ilo+4*n:], 1) +
		bi.Ddot(nr, work[ilo+5*n:], 1, work[ilo+5*n:], 1)
	var ew, ewc float64
	for i = ilo; i < ihi; i++ {
		ew += work[i+4*n]
		ewc += work[i+5*n]
	}
	gamma = coef*gamma - coef2*(math.Pow(ew, 2)+math.Pow(ewc, 2)) - coef5*math.Pow(ew-ewc, 2)
	if gamma == 0 {
		goto end
	}
	if it != 1 {
		beta = gamma / pgamma
	}
	t = coef5 * (ewc - 3*ew)
	tc = coef5 * (ew - 3*ewc)
	bi.Dscal(nr, beta, work[ilo:], 1)
	bi.Dscal(nr, beta, work[ilo+n:], 1)

	bi.Daxpy(nr, coef, work[ilo+4*n:], 1, work[ilo+n:], 1)
	bi.Daxpy(nr, coef, work[ilo+5*n:], 1, work[ilo:], 1)

	for i = ilo; i < ihi; i++ {
		work[i] += tc
		work[i+n] += t
	}

	for i = ilo; i < ihi; i++ {
		kount = 0
		sum = 0
		for j = ilo; j < ihi; j++ {
			if a[i*lda+j] != 0 {
				kount++
				sum += work[j]
			}
			if b[i*ldb+j] != 0 {
				kount++
				sum += work[j]
			}
		}
		work[i+2*n] = float64(kount)*work[i+n] + sum
	}
	for j = ilo; j < ihi; j++ {
		kount = 0
		sum = 0
		for i = ilo; i < ihi; i++ {
			if a[i*lda+j] != 0 {
				kount++
				sum += work[i+n]
			}
			if b[i*ldb+j] != 0 {
				kount++
				sum += work[i+n]
			}
		}
		work[j+3*n] = float64(kount)*work[j] + sum
	}
	sum = bi.Ddot(nr, work[ilo+n:], 1, work[ilo+2*n:], 1) +
		bi.Ddot(nr, work[ilo:], 1, work[ilo+3*n:], 1)
	alpha = gamma / sum

	// Determine correction to current iteration.
	cmax = 0
	for i = ilo; i < ihi; i++ {
		cor = alpha * work[i+n]
		cmax = math.Max(math.Abs(cor), cmax)
		lscale[i] += cor
		cor = alpha * work[i]
		cmax = math.Max(math.Abs(cor), cmax)
		rscale[i] += cor
	}
	if cmax < 0.5 {
		goto end
	}
	bi.Daxpy(nr, -alpha, work[ilo+2*n:], 1, work[ilo+4*n:], 1)
	bi.Daxpy(nr, -alpha, work[ilo+3*n:], 1, work[ilo+5*n:], 1)
	pgamma = gamma
	it++
	if it <= nrp2 {
		goto twoFiddy
	}
	// End generalized conjugate gradient iteration.

end: // LABEL 350
	sfmin := dlamchS
	sfmax := 1 / sfmin
	lsfmin := int(math.Log10(sfmin)/basl) + 1
	lsfmax := int(math.Log10(sfmax) / basl)
	var irab, lrab, ir, icab, lcab, jc int
	var rab, cab float64
	for i = ilo; i < ihi; i++ {
		irab = bi.Idamax(n-ilo+1, a[i*lda+ilo:], 1)
		rab = math.Abs(a[i*lda+irab+ilo-1])
		irab = bi.Idamax(n-ilo+1, b[i*ldb+ilo:], 1)
		rab = math.Max(rab, math.Abs(b[i*ldb+irab+ilo-1]))
		lrab = int(math.Log10(rab+sfmin)/basl) + 1
		ir = int(lscale[i] + math.Copysign(.5, lscale[i]))
		ir = min(min(max(ir, lsfmin), lsfmax), lsfmax-lrab)
		lscale[i] = math.Pow(sclfac, float64(ir))
		icab = bi.Idamax(ihi, a[i:], lda)
		cab = math.Abs(a[icab*lda+i])
		icab = bi.Idamax(ihi, b[i:], ldb)
		cab = math.Max(cab, math.Abs(b[icab*ldb+i]))
		lcab = int(math.Log10(cab+sfmin)/basl) + 1
		jc = int(rscale[i] + math.Copysign(.5, rscale[i]))
		jc = min(min(max(jc, lsfmin), lsfmax), lsfmax-lcab)
		rscale[i] = math.Pow(sclfac, float64(jc))
	}

	// Row scaling of matrices A and B.
	for i = ilo; i < ihi; i++ {
		bi.Dscal(n-ilo+1, lscale[i], a[i*lda+ilo:], 1)
		bi.Dscal(n-ilo+1, lscale[i], b[i*ldb+ilo:], 1)
	}

	// Column scaling of matrices A and B.
	for j = ilo; j < ihi; j++ {
		bi.Dscal(ihi, rscale[j], a[j:], lda)
		bi.Dscal(ihi, rscale[j], b[j:], ldb)
	}
	return
}
