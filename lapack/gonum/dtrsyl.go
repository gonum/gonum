// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dtrsyl solves the real Sylvester matrix equation:
//
//	op(A)*X + isgn*X*op(B) = scale*C
//
// where op(A) = A or Aᵀ, op(B) = B or Bᵀ, and isgn = +1 or -1.
//
// A is an m×m upper quasi-triangular matrix in Schur canonical form from
// Dgees or Dhseqr. B is an n×n upper quasi-triangular matrix in Schur
// canonical form.
//
// The solution X overwrites C, and scale is an output scaling factor set
// less than or equal to 1 to avoid overflow in X.
//
// Dtrsyl returns ok=false if A and B have common or very close eigenvalues,
// making the problem ill-conditioned. In this case, A or B is perturbed to
// compute an approximate solution that is within machine precision of an
// exact solution to a perturbed problem.
//
// Dtrsyl is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dtrsyl(trana, tranb blas.Transpose, isgn, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int) (scale float64, ok bool) {
	switch trana {
	case blas.NoTrans, blas.Trans:
	default:
		panic(badTrans)
	}
	switch tranb {
	case blas.NoTrans, blas.Trans:
	default:
		panic(badTrans)
	}
	switch {
	case isgn != 1 && isgn != -1:
		panic("lapack: bad isgn")
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case lda < max(1, m):
		panic(badLdA)
	case ldb < max(1, n):
		panic(badLdB)
	case ldc < max(1, n):
		panic(badLdC)
	}

	scale = 1
	ok = true

	if m == 0 || n == 0 {
		return scale, ok
	}

	switch {
	case len(a) < (m-1)*lda+m:
		panic(shortA)
	case len(b) < (n-1)*ldb+n:
		panic(shortB)
	case len(c) < (m-1)*ldc+n:
		panic(shortC)
	}

	notrna := trana == blas.NoTrans
	notrnb := tranb == blas.NoTrans

	eps := dlamchP
	smlnum := dlamchS / eps
	sgn := float64(isgn)

	bi := blas64.Implementation()

	anrm := impl.Dlange(lapack.MaxAbs, m, m, a, lda, nil)
	bnrm := impl.Dlange(lapack.MaxAbs, n, n, b, ldb, nil)
	smin := math.Max(smlnum, eps*math.Max(anrm, bnrm))

	if notrna && notrnb {
		// op(A) = A, op(B) = B
		// Solve A*X + isgn*X*B = scale*C
		// Process columns left-to-right, rows bottom-to-top.
		// suml = A[k, k+1:m] · X[k+1:m, l] (contribution from rows below)
		// sumr = X[k, 0:l] · B[0:l, l] (contribution from columns left)
		lnext := 0
		for l := 0; l < n; {
			if l < n-1 && b[(l+1)*ldb+l] != 0 {
				lnext = l + 2
			} else {
				lnext = l + 1
			}
			l1, l2 := l, lnext-1

			knext := m
			for k := m - 1; k >= 0; {
				if k > 0 && a[k*lda+k-1] != 0 {
					knext = k - 1
				} else {
					knext = k
				}
				k1, k2 := knext, k

				if k1 == k2 && l1 == l2 {
					var suml, sumr float64
					if m > k1+1 {
						suml = bi.Ddot(m-k1-1, a[k1*lda+k1+1:], 1, c[(k1+1)*ldc+l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec := c[k1*ldc+l1] - (suml + sgn*sumr)

					scaloc := 1.0
					a11 := a[k1*lda+k1] + sgn*b[l1*ldb+l1]
					da11 := math.Abs(a11)
					if da11 <= smin {
						a11 = smin
						da11 = smin
						ok = false
					}
					db := math.Abs(vec)
					if da11 < 1 && db > 1 && db > bignum*da11 {
						scaloc = 1 / db
					}
					x11 := (vec * scaloc) / a11
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x11
				} else if k1 == k2 && l1 != l2 {
					var suml, sumr float64
					if m > k1+1 {
						suml = bi.Ddot(m-k1-1, a[k1*lda+k1+1:], 1, c[(k1+1)*ldc+l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k1+1 {
						suml = bi.Ddot(m-k1-1, a[k1*lda+k1+1:], 1, c[(k1+1)*ldc+l2:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l2:], ldb)
					}
					vec1 := c[k1*ldc+l2] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(false, false, isgn, 1, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
				} else if k1 != k2 && l1 == l2 {
					var suml, sumr float64
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k1*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k2*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k2*ldc:], 1, b[l1:], ldb)
					}
					vec1 := c[k2*ldc+l1] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(false, false, isgn, 2, 1,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 1, x[:], 1)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k2*ldc+l1] = x[1]
				} else {
					var suml, sumr float64
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k1*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec00 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k1*lda+k2+1:], 1, c[(k2+1)*ldc+l2:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l2:], ldb)
					}
					vec01 := c[k1*ldc+l2] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k2*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k2*ldc:], 1, b[l1:], ldb)
					}
					vec10 := c[k2*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k2*lda+k2+1:], 1, c[(k2+1)*ldc+l2:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k2*ldc:], 1, b[l2:], ldb)
					}
					vec11 := c[k2*ldc+l2] - (suml + sgn*sumr)

					var vec, x [4]float64
					vec[0], vec[1] = vec00, vec01
					vec[2], vec[3] = vec10, vec11
					scaloc, _, ierr := impl.Dlasy2(false, false, isgn, 2, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
					c[k2*ldc+l1] = x[2]
					c[k2*ldc+l2] = x[3]
				}

				k = knext - 1
			}
			l = lnext
		}
	} else if !notrna && notrnb {
		// op(A) = Aᵀ, op(B) = B
		// Solve Aᵀ*X + isgn*X*B = scale*C
		// Process columns left-to-right, rows top-to-bottom.
		// suml = A[0:k, k] · X[0:k, l] (contribution from rows above)
		// sumr = X[k, 0:l] · B[0:l, l] (contribution from columns left)
		lnext := 0
		for l := 0; l < n; {
			if l < n-1 && b[(l+1)*ldb+l] != 0 {
				lnext = l + 2
			} else {
				lnext = l + 1
			}
			l1, l2 := l, lnext-1

			knext := 0
			for k := 0; k < m; {
				if k < m-1 && a[(k+1)*lda+k] != 0 {
					knext = k + 2
				} else {
					knext = k + 1
				}
				k1, k2 := k, knext-1

				if k1 == k2 && l1 == l2 {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec := c[k1*ldc+l1] - (suml + sgn*sumr)

					scaloc := 1.0
					a11 := a[k1*lda+k1] + sgn*b[l1*ldb+l1]
					da11 := math.Abs(a11)
					if da11 <= smin {
						a11 = smin
						da11 = smin
						ok = false
					}
					db := math.Abs(vec)
					if da11 < 1 && db > 1 && db > bignum*da11 {
						scaloc = 1 / db
					}
					x11 := (vec * scaloc) / a11
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x11
				} else if k1 == k2 && l1 != l2 {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l2:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l2:], ldb)
					}
					vec1 := c[k1*ldc+l2] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(true, false, isgn, 1, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
				} else if k1 != k2 && l1 == l2 {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k2:], lda, c[l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k2*ldc:], 1, b[l1:], ldb)
					}
					vec1 := c[k2*ldc+l1] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(true, false, isgn, 2, 1,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 1, x[:], 1)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k2*ldc+l1] = x[1]
				} else {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l1:], ldb)
					}
					vec00 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l2:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k1*ldc:], 1, b[l2:], ldb)
					}
					vec01 := c[k1*ldc+l2] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k2:], lda, c[l1:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k2*ldc:], 1, b[l1:], ldb)
					}
					vec10 := c[k2*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k2:], lda, c[l2:], ldc)
					}
					if l1 > 0 {
						sumr = bi.Ddot(l1, c[k2*ldc:], 1, b[l2:], ldb)
					}
					vec11 := c[k2*ldc+l2] - (suml + sgn*sumr)

					var vec, x [4]float64
					vec[0], vec[1] = vec00, vec01
					vec[2], vec[3] = vec10, vec11
					scaloc, _, ierr := impl.Dlasy2(true, false, isgn, 2, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
					c[k2*ldc+l1] = x[2]
					c[k2*ldc+l2] = x[3]
				}

				k = knext
			}
			l = lnext
		}
	} else if notrna && !notrnb {
		// op(A) = A, op(B) = Bᵀ
		// Solve A*X + isgn*X*Bᵀ = scale*C
		// Process columns right-to-left, rows bottom-to-top.
		// suml = A[k, k+1:m] · X[k+1:m, l] (contribution from rows below)
		// sumr = X[k, l+1:n] · B[l, l+1:n] (contribution from columns right)
		lnext := n
		for l := n - 1; l >= 0; {
			if l > 0 && b[l*ldb+l-1] != 0 {
				lnext = l - 1
			} else {
				lnext = l
			}
			l1, l2 := lnext, l

			knext := m
			for k := m - 1; k >= 0; {
				if k > 0 && a[k*lda+k-1] != 0 {
					knext = k - 1
				} else {
					knext = k
				}
				k1, k2 := knext, k

				if k1 == k2 && l1 == l2 {
					var suml, sumr float64
					if m > k1+1 {
						suml = bi.Ddot(m-k1-1, a[k1*lda+k1+1:], 1, c[(k1+1)*ldc+l1:], ldc)
					}
					if n > l1+1 {
						sumr = bi.Ddot(n-l1-1, c[k1*ldc+l1+1:], 1, b[l1*ldb+l1+1:], 1)
					}
					vec := c[k1*ldc+l1] - (suml + sgn*sumr)

					scaloc := 1.0
					a11 := a[k1*lda+k1] + sgn*b[l1*ldb+l1]
					da11 := math.Abs(a11)
					if da11 <= smin {
						a11 = smin
						da11 = smin
						ok = false
					}
					db := math.Abs(vec)
					if da11 < 1 && db > 1 && db > bignum*da11 {
						scaloc = 1 / db
					}
					x11 := (vec * scaloc) / a11
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x11
				} else if k1 == k2 && l1 != l2 {
					var suml, sumr float64
					if m > k1+1 {
						suml = bi.Ddot(m-k1-1, a[k1*lda+k1+1:], 1, c[(k1+1)*ldc+l1:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l1*ldb+l2+1:], 1)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k1+1 {
						suml = bi.Ddot(m-k1-1, a[k1*lda+k1+1:], 1, c[(k1+1)*ldc+l2:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l2*ldb+l2+1:], 1)
					}
					vec1 := c[k1*ldc+l2] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(false, true, isgn, 1, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
				} else if k1 != k2 && l1 == l2 {
					var suml, sumr float64
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k1*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if n > l1+1 {
						sumr = bi.Ddot(n-l1-1, c[k1*ldc+l1+1:], 1, b[l1*ldb+l1+1:], 1)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k2*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if n > l1+1 {
						sumr = bi.Ddot(n-l1-1, c[k2*ldc+l1+1:], 1, b[l1*ldb+l1+1:], 1)
					}
					vec1 := c[k2*ldc+l1] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(false, true, isgn, 2, 1,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 1, x[:], 1)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k2*ldc+l1] = x[1]
				} else {
					var suml, sumr float64
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k1*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l1*ldb+l2+1:], 1)
					}
					vec00 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k1*lda+k2+1:], 1, c[(k2+1)*ldc+l2:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l2*ldb+l2+1:], 1)
					}
					vec01 := c[k1*ldc+l2] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k2*lda+k2+1:], 1, c[(k2+1)*ldc+l1:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k2*ldc+l2+1:], 1, b[l1*ldb+l2+1:], 1)
					}
					vec10 := c[k2*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if m > k2+1 {
						suml = bi.Ddot(m-k2-1, a[k2*lda+k2+1:], 1, c[(k2+1)*ldc+l2:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k2*ldc+l2+1:], 1, b[l2*ldb+l2+1:], 1)
					}
					vec11 := c[k2*ldc+l2] - (suml + sgn*sumr)

					var vec, x [4]float64
					vec[0], vec[1] = vec00, vec01
					vec[2], vec[3] = vec10, vec11
					scaloc, _, ierr := impl.Dlasy2(false, true, isgn, 2, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
					c[k2*ldc+l1] = x[2]
					c[k2*ldc+l2] = x[3]
				}

				k = knext - 1
			}
			l = lnext - 1
		}
	} else {
		// op(A) = Aᵀ, op(B) = Bᵀ
		// Solve Aᵀ*X + isgn*X*Bᵀ = scale*C
		// Process columns right-to-left, rows top-to-bottom.
		// suml = A[0:k, k] · X[0:k, l] (contribution from rows above)
		// sumr = X[k, l+1:n] · B[l, l+1:n] (contribution from columns right)
		lnext := n
		for l := n - 1; l >= 0; {
			if l > 0 && b[l*ldb+l-1] != 0 {
				lnext = l - 1
			} else {
				lnext = l
			}
			l1, l2 := lnext, l

			knext := 0
			for k := 0; k < m; {
				if k < m-1 && a[(k+1)*lda+k] != 0 {
					knext = k + 2
				} else {
					knext = k + 1
				}
				k1, k2 := k, knext-1

				if k1 == k2 && l1 == l2 {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if n > l1+1 {
						sumr = bi.Ddot(n-l1-1, c[k1*ldc+l1+1:], 1, b[l1*ldb+l1+1:], 1)
					}
					vec := c[k1*ldc+l1] - (suml + sgn*sumr)

					scaloc := 1.0
					a11 := a[k1*lda+k1] + sgn*b[l1*ldb+l1]
					da11 := math.Abs(a11)
					if da11 <= smin {
						a11 = smin
						da11 = smin
						ok = false
					}
					db := math.Abs(vec)
					if da11 < 1 && db > 1 && db > bignum*da11 {
						scaloc = 1 / db
					}
					x11 := (vec * scaloc) / a11
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x11
				} else if k1 == k2 && l1 != l2 {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l1*ldb+l2+1:], 1)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l2:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l2*ldb+l2+1:], 1)
					}
					vec1 := c[k1*ldc+l2] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(true, true, isgn, 1, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
				} else if k1 != k2 && l1 == l2 {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if n > l1+1 {
						sumr = bi.Ddot(n-l1-1, c[k1*ldc+l1+1:], 1, b[l1*ldb+l1+1:], 1)
					}
					vec0 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k2:], lda, c[l1:], ldc)
					}
					if n > l1+1 {
						sumr = bi.Ddot(n-l1-1, c[k2*ldc+l1+1:], 1, b[l1*ldb+l1+1:], 1)
					}
					vec1 := c[k2*ldc+l1] - (suml + sgn*sumr)

					var vec, x [2]float64
					vec[0], vec[1] = vec0, vec1
					scaloc, _, ierr := impl.Dlasy2(true, true, isgn, 2, 1,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 1, x[:], 1)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k2*ldc+l1] = x[1]
				} else {
					var suml, sumr float64
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l1:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l1*ldb+l2+1:], 1)
					}
					vec00 := c[k1*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k1:], lda, c[l2:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k1*ldc+l2+1:], 1, b[l2*ldb+l2+1:], 1)
					}
					vec01 := c[k1*ldc+l2] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k2:], lda, c[l1:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k2*ldc+l2+1:], 1, b[l1*ldb+l2+1:], 1)
					}
					vec10 := c[k2*ldc+l1] - (suml + sgn*sumr)

					suml, sumr = 0, 0
					if k1 > 0 {
						suml = bi.Ddot(k1, a[k2:], lda, c[l2:], ldc)
					}
					if n > l2+1 {
						sumr = bi.Ddot(n-l2-1, c[k2*ldc+l2+1:], 1, b[l2*ldb+l2+1:], 1)
					}
					vec11 := c[k2*ldc+l2] - (suml + sgn*sumr)

					var vec, x [4]float64
					vec[0], vec[1] = vec00, vec01
					vec[2], vec[3] = vec10, vec11
					scaloc, _, ierr := impl.Dlasy2(true, true, isgn, 2, 2,
						a[k1*lda+k1:], lda, b[l1*ldb+l1:], ldb, vec[:], 2, x[:], 2)
					if !ierr {
						ok = false
					}
					if scaloc != 1 {
						for j := 0; j < n; j++ {
							bi.Dscal(m, scaloc, c[j:], ldc)
						}
						scale *= scaloc
					}
					c[k1*ldc+l1] = x[0]
					c[k1*ldc+l2] = x[1]
					c[k2*ldc+l1] = x[2]
					c[k2*ldc+l2] = x[3]
				}

				k = knext
			}
			l = lnext - 1
		}
	}

	return scale, ok
}

const bignum = 1 / dlamchS
