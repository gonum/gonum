// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"math"

	"github.com/gonum/lapack"
)

// Dsterf computes all eigenvalues of a symmetric tridiagonal matrix using the
// Pal-Walker-Kahan variant of the QL or QR algorithm.
//
// d contains the diagonal elements of the tridiagonal matrix on entry, and
// contains the eigenvalues in ascending order on exit. d must have length at
// least n, or Dsterf will panic.
//
// e contains the off-diagonal elements of the tridiagonal matrix on entry, and is
// overwritten during the call to Dsterf. e must have length of at least n-1 or
// Dsterf will panic.
func (impl Implementation) Dsterf(n int, d, e []float64) (ok bool) {
	if n < 0 {
		panic(nLT0)
	}
	if n == 0 {
		return true
	}
	if len(d) < n {
		panic(badD)
	}
	if len(e) < n-1 {
		panic(badE)
	}

	const (
		none = 0 // The values are not scaled.
		down = 1 // The values are scaled below ssfmax threshold.
		up   = 2 // The values are scaled below ssfmin threshold.
	)

	// Determine the unit roundoff for this environment.
	eps := dlamchE
	eps2 := eps * eps
	safmin := dlamchS
	safmax := 1 / safmin
	ssfmax := math.Sqrt(safmax) / 3
	ssfmin := math.Sqrt(safmin) / eps2

	// Compute the eigenvalues of the tridiagonal matrix.
	maxit := 30
	nmaxit := n * maxit
	jtot := 0

	l1 := 0

	// TODO(btracey): Define these closer to use when gotos are removed.
	var anorm, c, gamma, r, rte, s, sigma float64
	var iscale, l, lend, lendsv, lsv, m int
	var el []float64

	// TOOD(btracey): Replace these goto statements with imperative flow control
	// structures.
Ten:
	if l1 > n-1 {
		goto OneSeventy
	}
	if l1 > 0 {
		e[l1-1] = 0
	}
	for m = l1; m < n-1; m++ {
		if math.Abs(e[m]) <= math.Sqrt(math.Abs(d[m]))*math.Sqrt(math.Abs(d[m+1]))*eps {
			e[m] = 0
			goto Thirty
		}
	}
	m = n - 1

Thirty:
	l = l1
	lsv = l
	lend = m
	lendsv = lend
	l1 = m + 1
	if lend == 0 {
		goto Ten
	}

	// Scale submatrix in rows and columns l to lend.
	anorm = impl.Dlanst(lapack.MaxAbs, lend-l+1, d[l:], e[l:])
	iscale = none
	if anorm == 0 {
		goto Ten
	}
	if anorm > ssfmax {
		iscale = down
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmax, lend-l+1, 1, d[l:], n)
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmax, lend-l, 1, e[l:], n)
	} else if anorm < ssfmin {
		iscale = up
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmin, lend-l+1, 1, d[l:], n)
		impl.Dlascl(lapack.General, 0, 0, anorm, ssfmin, lend-l, 1, e[l:], n)
	}

	el = e[l:lend]
	for i, v := range el {
		el[i] *= v
	}

	// Choose between QL and QR iteration.
	if math.Abs(d[lend]) < math.Abs(d[l]) {
		lend = lsv
		l = lendsv
	}
	if lend >= l {
		// QL Iteration.
		// Look for small sub-diagonal element.
	Fifty:
		if l != lend {
			for m = l; m < lend; m++ {
				if math.Abs(e[m]) <= eps2*(math.Abs(d[m]*d[m+1])) {
					goto Seventy
				}
			}
		}
		m = lend
	Seventy:
		if m < lend {
			e[m] = 0
		}
		p := d[l]
		if m == l {
			goto Ninety
		}
		// If remaining matrix is 2 by 2, use Dlae2 to compute its eigenvalues.
		if m == l+1 {
			rte = math.Sqrt(e[l])
			d[l], d[l+1] = impl.Dlae2(d[l], rte, d[l+1])
			e[l] = 0
			l += 2
			if l <= lend {
				goto Fifty
			}
			goto OneFifty
		}
		if jtot == nmaxit {
			goto OneFifty
		}
		jtot++

		// Form shift.
		rte = math.Sqrt(e[l])
		sigma = (d[l+1] - p) / (2 * rte)
		r = impl.Dlapy2(sigma, 1)
		sigma = p - (rte / (sigma + math.Copysign(r, sigma)))

		c = 1
		s = 0
		gamma = d[m] - sigma
		p = gamma * gamma

		// Inner loop.
		for i := m - 1; i >= l; i-- {
			bb := e[i]
			r := p + bb
			if i != m-1 {
				e[i+1] = s * r
			}
			oldc := c
			c = p / r
			s = bb / r
			oldgam := gamma
			alpha := d[i]
			gamma = c*(alpha-sigma) - s*oldgam
			d[i+1] = oldgam + (alpha - gamma)
			if c != 0 {
				p = (gamma * gamma) / c
			} else {
				p = oldc * bb
			}
		}
		e[l] = s * p
		d[l] = sigma + gamma
		goto Fifty

	Ninety:
		// Eigenvalue found.
		d[l] = p
		l++
		if l <= lend {
			goto Fifty
		}
		goto OneFifty
	} else {
	OneHundred:
		// QR Iteration.
		// Look for small super-diagonal element.
		for m = l; m >= lend+1; m-- {
			if math.Abs(e[m-1]) <= eps2*math.Abs(d[m]*d[m-1]) {
				goto OneTwenty
			}
		}
		m = lend
	OneTwenty:
		if m > lend {
			e[m-1] = 0
		}
		p := d[l]
		if m == l {
			goto OneFourty
		}

		// If remaining matrix is 2 by 2, use Dlae2 to compute its eigenvalues.
		if m == l-1 {
			rte = math.Sqrt(e[l-1])
			d[l], d[l-1] = impl.Dlae2(d[l], rte, d[l-1])
			e[l-1] = 0
			l -= 2
			if l >= lend {
				goto OneHundred
			}
			goto OneFifty
		}
		if jtot == nmaxit {
			goto OneFifty
		}
		jtot++

		// Form shift.
		rte = math.Sqrt(e[l-1])
		sigma = (d[l-1] - p) / (2 * rte)
		r = impl.Dlapy2(sigma, 1)
		sigma = p - (rte / (sigma + math.Copysign(r, sigma)))

		c = 1
		s = 0
		gamma = d[m] - sigma
		p = gamma * gamma

		// Inner loop.
		for i := m; i < l; i++ {
			bb := e[i]
			r := p + bb
			if i != m {
				e[i-1] = s * r
			}
			oldc := c
			c = p / r
			s = bb / r
			oldgam := gamma
			alpha := d[i+1]
			gamma = c*(alpha-sigma) - s*oldgam
			d[i] = oldgam + alpha - gamma
			if c != 0 {
				p = (gamma * gamma) / c
			} else {
				p = oldc * bb
			}
		}
		e[l-1] = s * p
		d[l] = sigma + gamma
		goto OneHundred

	OneFourty:
		// Eigenvalue found.
		d[l] = p
		l--
		if l >= lend {
			goto OneHundred
		}
		goto OneFifty
	}
OneFifty:
	// Undo scaling if necessary
	switch iscale {
	case down:
		impl.Dlascl(lapack.General, 0, 0, ssfmax, anorm, lendsv-lsv+1, 1, d[lsv:], n)
	case up:
		impl.Dlascl(lapack.General, 0, 0, ssfmin, anorm, lendsv-lsv+1, 1, d[lsv:], n)
	}

	// Check for no convergence to an eigenvalue after a total of n*maxit iterations.
	if jtot < nmaxit {
		goto Ten
	}
	for _, v := range e[:n-1] {
		if v != 0 {
			return false
		}
	}
OneSeventy:
	impl.Dlasrt(lapack.SortIncreasing, n, d)
	return true
}
