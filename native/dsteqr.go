// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
)

// Dsteqr computes the eigenvalues and optionally the eigenvectors of a symmetric
// tridiagonal matrix using the implicit QL or QR method. The eigenvectors of a
// full or band symmetric matrix can also be found if Dsytrd, Dsptrd, or Dsbtrd
// have been used to reduce this matrix to tridiagonal form.
//
// d, on entry, contains the diagonal elements of the tridiagonal matrix. On exit,
// d contains the eigenvalues in ascending order. d must have length n and
// Dsteqr will panic otherwise.
//
// e, on entry, contains the off-diagonal elements of the tridiagonal matrix on
// entry, and is overwritten during the call to Dsteqr. e must have length n-1 and
// Dsteqr will panic otherwise.
//
// z, on entry, contains the n×n orthogonal matrix used in the reduction to
// tridiagonal form if compz == lapack.EigDecomp. On exit, if
// compz == lapack.EigBoth, z contains the orthonormal eigenvectors of the
// original symmetric matrix, and if compz == lapack.EigDecomp, z contains the
// orthonormal eigenvectors of the symmetric tridiagonal matrix. z is not used
// if compz == lapack.EigValueOnly.
//
// work must have length at least max(1, 2*n-2) if the eigenvectors are computed,
// and Dsteqr will panic otherwise.
func (impl Implementation) Dsteqr(compz lapack.EigComp, n int, d, e, z []float64, ldz int, work []float64) (ok bool) {
	if len(d) < n {
		panic(badD)
	}
	if len(e) < n-1 {
		panic(badE)
	}
	if compz != lapack.EigValueOnly && compz != lapack.EigBoth && compz != lapack.EigDecomp {
		panic(badEigComp)
	}
	if compz != lapack.EigValueOnly {
		if len(work) < max(1, 2*n-2) {
			panic(badWork)
		}
		checkMatrix(n, n, z, ldz)
	}

	var icompz int
	if compz == lapack.EigDecomp {
		icompz = 1
	} else if compz == lapack.EigBoth {
		icompz = 2
	}

	if n == 0 {
		return true
	}
	if n == 1 {
		if icompz == 2 {
			z[0] = 1
		}
		return true
	}

	bi := blas64.Implementation()

	eps := dlamchE
	eps2 := eps * eps
	safmin := dlamchS
	safmax := 1 / safmin
	ssfmax := math.Sqrt(safmax) / 3
	ssfmin := math.Sqrt(safmin) / eps2

	// Compute the eigenvalues and eigenvectors of the tridiagonal matrix.
	if icompz == 2 {
		impl.Dlaset(blas.All, n, n, 0, 1, z, ldz)
	}
	const maxit = 30
	nmaxit := n * maxit

	jtot := 0

	// Determine where the matrix splits and choose QL or QR iteration for each
	// block, according to whether top or bottom diagonal element is smaller.
	// TODO(btracey): Reformat the gotos to use structural flow techniques.
	// TODO(btracey): Move these variables closer to actual usage when
	// gotos are removed.
	l1 := 0
	nm1 := n - 1
	var l, lend, lendsv, lsv, m int
	var anorm, c, g, p, r, s float64

	type scaletype int
	const (
		none scaletype = iota
		down
		up
	)
	var iscale scaletype

Ten:
	if l1 > n-1 {
		goto OneSixty
	}
	if l1 > 0 {
		e[l1-1] = 0
	}
	if l1 <= nm1 {
		for m = l1; m < nm1; m++ {
			test := math.Abs(e[m])
			if test == 0 {
				goto Thirty
			}
			if test <= (math.Sqrt(math.Abs(d[m]))*math.Sqrt(math.Abs(d[m+1])))*eps {
				e[m] = 0
				goto Thirty
			}
		}
	}
	m = n - 1
Thirty:
	l = l1
	lsv = l
	lend = m
	lendsv = lend
	l1 = m + 1
	if lend == l {
		goto Ten
	}

	// Scale submatrix in rows and columns L to Lend
	anorm = impl.Dlanst(lapack.MaxAbs, lend-l+1, d[l:], e[l:])
	switch {
	case anorm == 0:
		goto Ten
	case anorm > ssfmax:
		iscale = down
		// TODO(btracey): Why is lda n?
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmax, lend-l+1, 1, d[l:], n)
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmax, lend-l, 1, e[l:], n)
	case anorm < ssfmin:
		iscale = up
		// TODO(btracey): Why is lda n?
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmin, lend-l+1, 1, d[l:], n)
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmin, lend-l, 1, e[l:], n)
	}

	// Choose between QL and QR.
	if math.Abs(d[lend]) < math.Abs(d[l]) {
		lend = lsv
		l = lendsv
	}
	if lend > l {
		// QL Iteration. Look for small subdiagonal element.
	Fourty:
		if l != lend {
			lendm1 := lend - 1
			for m = l; m <= lendm1; m++ {
				v := math.Abs(e[m])
				if v*v <= (eps2*math.Abs(d[m]))*math.Abs(d[m+1])+safmin {
					goto Sixty
				}
			}
		}
		m = lend
	Sixty:
		if m < lend {
			e[m] = 0
		}
		p = d[l]
		if m == l {
			goto Eighty
		}

		// If remaining matrix is 2×2, use Dlae2 to compute its eigensystem.
		if m == l+1 {
			if icompz > 0 {
				d[l], d[l+1], work[l], work[n-1+l] = impl.Dlaev2(d[l], e[l], d[l+1])
				impl.Dlasr(blas.Right, lapack.Variable, lapack.Backward,
					n, 2, work[l:], work[n-1+l:], z[l:], ldz)
			} else {
				d[l], d[l+1] = impl.Dlae2(d[l], e[l], d[l+1])
			}
			e[l] = 0
			l += 2
			if l <= lend {
				goto Fourty
			}
			goto OneFourty
		}

		if jtot == nmaxit {
			goto OneFourty
		}
		jtot++

		// Form shift
		g = (d[l+1] - p) / (2 * e[l])
		r = impl.Dlapy2(g, 1)
		g = d[m] - p + e[l]/(g+math.Copysign(r, g))
		s = 1
		c = 1
		p = 0

		// Inner loop
		for i := m - 1; i >= l; i-- {
			f := s * e[i]
			b := c * e[i]
			c, s, r = impl.Dlartg(g, f)
			if i != m-1 {
				e[i+1] = r
			}
			g = d[i+1] - p
			r = (d[i]-g)*s + 2*c*b
			p = s * r
			d[i+1] = g + p
			g = c*r - b

			// If eigenvectors are desired, then save rotations.
			if icompz > 0 {
				work[i] = c
				work[n-1+i] = -s
			}
		}
		// If eigenvectors are desired, then apply saved rotations.
		if icompz > 0 {
			mm := m - l + 1
			impl.Dlasr(blas.Right, lapack.Variable, lapack.Backward,
				n, mm, work[l:], work[n-1+l:], z[l:], ldz)
		}
		d[l] -= p
		e[l] = g
		goto Fourty

		// Eigenvalue found.
	Eighty:
		d[l] = p
		l++
		if l <= lend {
			goto Fourty
		}
		goto OneFourty
	} else {
		// QR Iteration.
		// Look for small superdiagonal element.
	Ninety:
		if l != lend {
			lendp1 := lend + 1
			for m = l; m >= lendp1; m-- {
				v := math.Abs(e[m-1])
				if v*v <= (eps2*math.Abs(d[m])*math.Abs(d[m-1]) + safmin) {
					goto OneTen
				}
			}
		}
		m = lend
	OneTen:
		if m > lend {
			e[m-1] = 0
		}
		p = d[l]
		if m == l {
			goto OneThirty
		}

		// If remaining matrix is 2×2, use Dlae2 to compute its eigenvalues.
		if m == l-1 {
			if icompz > 0 {
				d[l-1], d[l], work[m], work[n-1+m] = impl.Dlaev2(d[l-1], e[l-1], d[l])
				impl.Dlasr(blas.Right, lapack.Variable, lapack.Forward,
					n, 2, work[m:], work[n-1+m:], z[l-1:], ldz)
			} else {
				d[l-1], d[l] = impl.Dlae2(d[l-1], e[l-1], d[l])
			}
			e[l-1] = 0
			l -= 2
			if l >= lend {
				goto Ninety
			}
			goto OneFourty
		}
		if jtot == nmaxit {
			goto OneFourty
		}
		jtot++

		// Form shift.
		g = (d[l-1] - p) / (2 * e[l-1])
		r = impl.Dlapy2(g, 1)
		g = d[m] - p + (e[l-1])/(g+math.Copysign(r, g))

		s = 1
		c = 1
		p = 0

		// Inner loop.
		for i := m; i < l; i++ {
			f := s * e[i]
			b := c * e[i]
			c, s, r = impl.Dlartg(g, f)
			if i != m {
				e[i-1] = r
			}
			g = d[i] - p
			r = (d[i+1]-g)*s + 2*c*b
			p = s * r
			d[i] = g + p
			g = c*r - b

			// If eigenvectors are desired, then save rotations.
			if icompz > 0 {
				work[i] = c
				work[n-1+i] = s
			}
		}

		// If eigenvectors are desired, then apply saved rotations.
		if icompz > 0 {
			mm := l - m + 1
			impl.Dlasr(blas.Right, lapack.Variable, lapack.Forward,
				n, mm, work[m:], work[n-1+m:], z[m:], ldz)
		}
		d[l] -= p
		e[l-1] = g
		goto Ninety

		// Eigenvalue found
	OneThirty:
		d[l] = p
		l--
		if l >= lend {
			goto Ninety
		}
		goto OneFourty
	}

	// Undo scaling if necessary.
OneFourty:
	switch iscale {
	case down:
		impl.Dlascl(lapack.General, 0, 0, ssfmax, anorm, lendsv-lsv+1, 1, d[lsv:], n)
		impl.Dlascl(lapack.General, 0, 0, ssfmax, anorm, lendsv-lsv, 1, e[lsv:], n)
	case up:
		impl.Dlascl(lapack.General, 0, 0, ssfmin, anorm, lendsv-lsv+1, 1, e[lsv:], n)
		impl.Dlascl(lapack.General, 0, 0, ssfmin, anorm, lendsv-lsv, 1, e[lsv:], n)
	}

	// Check for no convergence to an eigenvalue after a total of n*maxit iterations.
	if jtot < nmaxit {
		goto Ten
	}
	for i := 0; i < n-1; i++ {
		if e[i] != 0 {
			return false
		}
	}
	return true

	// Order eigenvalues and eigenvectors.
OneSixty:
	if icompz == 0 {
		impl.Dlasrt(lapack.SortIncreasing, n, d)
	} else {
		// TODO(btracey): Consider replacing this sort with a call to sort.Sort.
		for ii := 1; ii < n; ii++ {
			i := ii - 1
			k := i
			p := d[i]
			for j := ii; j < n; j++ {
				if d[j] < p {
					k = j
					p = d[j]
				}
			}
			if k != i {
				d[k] = d[i]
				d[i] = p
				bi.Dswap(n, z[i:], ldz, z[k:], ldz)
			}
		}
	}
	return true
}
