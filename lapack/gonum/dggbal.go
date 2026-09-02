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
//
//   - n is the order of matrices A and B. n >= 0
//   - lscale, rscale must be of size n.
//   - a and b of size lda*n and ldb*n respectively. On exit they are overwritten with balanced matrices.
//   - work is at least of size max(1, 6*n) when job is Scale/PermuteScale. Otherwise at least size 1.
//   - ilo, ihi are indices such that a(i,j) and b(i,j) are zero for j=1..ilo-1, i=ihi+1..n.
//
// Dggbal is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dggbal(job lapack.BalanceJob, n int, a []float64, lda int, b []float64, ldb int, lscale, rscale, work []float64) (ilo, ihi int) {
	_ = column(a, lda, 0, 0)
	var (
		alpha, pgamma, t, tc, sum, cmax, cor, basl float64
		i, j, kount, lm1, jp1, iflow, m, ip1       int // loop var
	)
	sclfac := 10.
	switch {
	case job != lapack.BalanceNone && job != lapack.Permute && job != lapack.PermuteScale && job != lapack.Scale:
		panic(badBalanceJob)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case ldb < max(1, n):
		panic(badLdB)
	case len(b) < (n-1)*ldb+n:
		panic(shortB)
	case len(lscale) < n || len(rscale) < n:
		panic(shortScale)
	case len(work) < 1 || len(work) < 6*n && (job == lapack.Scale || job == lapack.PermuteScale):
		panic(shortWork)
	}

	// quick return if possible
	if n == 0 || n == 1 || job == lapack.BalanceNone {
		ilo = 0
		ihi = n - 1
		for i = 0; i < n; i++ {
			lscale[i] = 1
			rscale[i] = 1
		}
		return ilo, ihi
	}

	bi := blas64.Implementation()

	k := 0
	l := n - 1
	if job == lapack.Scale {
		goto OneNinety
	}
	goto Thirty

	// Permute the matrices A and B to isolate the eigenvalues.
	// Find row with one nonzero in columns 1..L.

Twenty:
	l = lm1
	if l != 0 {
		goto Thirty
	}
	rscale[0] = 1
	lscale[0] = 1
	goto OneNinety

Thirty:
	lm1 = l - 1
OUTER:
	for i = l; i >= 0; i-- {
		for j = 0; j <= lm1; j++ {
			jp1 = j + 1
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				goto Fifty
			}
			// Forty:
		}
		j = l
		goto Seventy
	Fifty:
		for j = jp1; j <= l; j++ {
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				continue OUTER
			}
			// Sixty:
		}
		j = jp1 - 1
	Seventy:
		m = l
		iflow = 1
		goto OneSixty
	}

	goto Hundred

	// Find column with one nonzero in rows K through N.
Ninety:
	k++
Hundred:
	for j = k; j <= l; j++ {
		for i = k; i <= lm1; i++ {
			ip1 = i + 1
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				goto OneTwenty
			}
		} // 110
		i = l
		goto OneForty
	OneTwenty:
		for i = ip1; i <= l; i++ {
			if a[i*lda+j] != 0 || b[i*ldb+j] != 0 {
				continue Hundred // goto 150.
			}
		} // 130
		i = ip1 - 1
	OneForty:
		m = k
		iflow = 2
		goto OneSixty
	} // 150
	goto OneNinety

	// Permute rows M and I

OneSixty:
	lscale[m] = float64(i)
	if i != m {
		bi.Dswap(n-k, a[i*lda+k:], 1, a[m*lda+k:], 1)
		bi.Dswap(n-k, b[i*ldb+k:], 1, b[m*ldb+k:], 1)
	}

	// Permute columns M and J

	rscale[m] = float64(j)
	if j != m {
		bi.Dswap(l+1, a[j:], lda, a[m:], lda)
		bi.Dswap(l+1, b[j:], ldb, b[m:], ldb)
	}
	switch iflow {
	case 1:
		goto Twenty
	case 2:
		goto Ninety
	}

OneNinety:
	ilo = k
	ihi = l
	if job == lapack.Permute {
		for i := ilo; i <= ihi; i++ {
			lscale[i] = 1
			rscale[i] = 1
		}
		return ilo, ihi
	}
	if ilo == ihi {
		return ilo, ihi
	}

	// Balance the submatrix in rows ILO to IHI.
	nr := ihi - ilo + 1
	for i = ilo; i <= ihi; i++ {
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
	for i = ilo; i <= ihi; i++ {
		for j = ilo; j <= ihi; j++ {
			tb := b[i*ldb+j]
			ta := a[i*lda+j]
			if ta != 0 {
				ta = math.Log10(math.Abs(ta)) / basl
			}
			if tb != 0 {
				tb = math.Log10(math.Abs(tb)) / basl
			}
			work[i+4*n] -= (ta + tb)
			work[j+5*n] -= (ta + tb)
		}
	}
	coef := 1 / float64(2*nr)
	coef2 := coef * coef
	coef5 := .5 * coef2
	nrp2 := nr + 2
	beta := 0.
	it := 1

	// Start generalized conjugate gradient iteration
TwoFifty:
	gamma := bi.Ddot(nr, work[ilo+4*n:], 1, work[ilo+4*n:], 1) +
		bi.Ddot(nr, work[ilo+5*n:], 1, work[ilo+5*n:], 1)
	var ew, ewc float64
	for i = ilo; i <= ihi; i++ {
		ew += work[i+4*n]
		ewc += work[i+5*n]
	}
	ewmewc := ew - ewc
	gamma = coef*gamma - coef2*(ew*ew+ewc*ewc) - coef5*(ewmewc*ewmewc)
	if gamma == 0 {
		goto End
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

	for i = ilo; i <= ihi; i++ {
		work[i] += tc
		work[i+n] += t
	}

	// Apply matrix to vector.

	for i = ilo; i <= ihi; i++ {
		kount = 0
		sum = 0
		for j = ilo; j <= ihi; j++ {
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
	for j = ilo; j <= ihi; j++ {
		kount = 0
		sum = 0
		for i = ilo; i <= ihi; i++ {
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
	for i = ilo; i <= ihi; i++ {
		cor = alpha * work[i+n]
		cmax = math.Max(math.Abs(cor), cmax)
		lscale[i] += cor
		cor = alpha * work[i]
		cmax = math.Max(math.Abs(cor), cmax)
		rscale[i] += cor
	}
	if cmax < 0.5 {
		goto End
	}
	bi.Daxpy(nr, -alpha, work[ilo+2*n:], 1, work[ilo+4*n:], 1)
	bi.Daxpy(nr, -alpha, work[ilo+3*n:], 1, work[ilo+5*n:], 1)
	pgamma = gamma
	it++
	if it <= nrp2 {
		goto TwoFifty
	}

	// End generalized conjugate gradient iteration.

End: // LABEL 350
	sfmin := dlamchS
	sfmax := 1 / sfmin
	lsfmin := int(math.Log10(sfmin)/basl) + 1
	lsfmax := int(math.Log10(sfmax) / basl)
	var irab, lrab, ir, icab, lcab, jc int
	var rab, cab float64
	for i = ilo; i <= ihi; i++ {
		irab = bi.Idamax(n-ilo, a[i*lda+ilo:], 1)
		rab = math.Abs(a[i*lda+irab+ilo])
		irab = bi.Idamax(n-ilo, b[i*ldb+ilo:], 1)
		rab = math.Max(rab, math.Abs(b[i*ldb+irab+ilo]))
		lrab = int(math.Log10(rab+sfmin)/basl) + 1
		ir = int(lscale[i] + math.Copysign(.5, lscale[i]))
		ir = min(min(max(ir, lsfmin), lsfmax), lsfmax-lrab)
		lscale[i] = math.Pow(sclfac, float64(ir))

		icab = bi.Idamax(ihi+1, a[i:], lda)
		cab = math.Abs(a[icab*lda+i])
		icab = bi.Idamax(ihi+1, b[i:], ldb)
		cab = math.Max(cab, math.Abs(b[icab*ldb+i]))
		lcab = int(math.Log10(cab+sfmin)/basl) + 1
		jc = int(rscale[i] + math.Copysign(.5, rscale[i]))
		jc = min(min(max(jc, lsfmin), lsfmax), lsfmax-lcab)
		rscale[i] = math.Pow(sclfac, float64(jc))
	}

	// Row scaling of matrices A and B.
	for i = ilo; i <= ihi; i++ {
		bi.Dscal(n-ilo, lscale[i], a[i*lda+ilo:], 1)
		bi.Dscal(n-ilo, lscale[i], b[i*ldb+ilo:], 1)
	}

	// Column scaling of matrices A and B.
	for j = ilo; j <= ihi; j++ {
		bi.Dscal(ihi+1, rscale[j], a[j:], lda)
		bi.Dscal(ihi+1, rscale[j], b[j:], ldb)
	}
	return ilo, ihi
}

func column(z []float64, ldz, j, m int) []float64 {
	v := make([]float64, m)
	for i := range v {
		v[i] = z[i*ldz+j]
	}
	return v
}
