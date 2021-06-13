// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// DTGSY2 solves the generalized Sylvester equation:
//  A * R - L * B = scale * C                (1)
//  D * R - L * E = scale * F,
// using Level 1 and 2 BLAS. where R and L are unknown m×n matrices,
// (A, D), (B, E) and (C, F) are given matrix pairs of size m×m,
// n×n and m×n, respectively, with real entries. (A, D) and (B, E)
// must be in generalized Schur canonical form, i.e. A, B are upper
// quasi triangular and D, E are upper triangular. The solution (R, L)
// overwrites (C, F). 0 <= scale <= 1 is an output scaling factor
// chosen to avoid overflow.
//
// In matrix notation solving equation (1) corresponds to solve
// Z*x = scale*b, where Z is defined as
//  Z = [ kron(In, A)  -kron(B**T, Im) ]             (2)
//      [ kron(In, D)  -kron(E**T, Im) ],
// Ik is the identity matrix of size k and X**T is the transpose of X.
// kron(X, Y) is the Kronecker product between the matrices X and Y.
// In the process of solving (1), we solve a number of such systems
// where Dim(In), Dim(In) = 1 or 2.
// If TRANS = 'T', solve the transposed system Z**T*y = scale*b for y,
// which is equivalent to solve for R and L in
//  A**T * R  + D**T * L   = scale * C           (3)
//  R  * B**T + L  * E**T  = scale * -F
// This case is used to compute an estimate of Dif[(A, D), (B, E)] =
// sigma_min(Z) using reverse communication with DLACON.
// Dtgsy2 also (IJOB >= 1) contributes to the computation in Dtgsyl
// of an upper bound on the separation between to matrix pairs. Then
// the input (A, D), (B, E) are sub-pencils of the matrix pair in
// Dtgsyl. See Dtgsyl for details.
func (impl Implementation) Dtgsy2(trans blas.Transpose, ijob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, rdsum, rdscal float64, iwork []int) (scale, sumout, scalout float64, pq int) {
	switch {
	case trans != blas.NoTrans && trans != blas.Trans:
		panic(badTrans)
	case trans == blas.NoTrans && (ijob < 0 || ijob > 2):
		panic(badIjob)
	case m <= 0:
		panic(mLT0)
	case n <= 0:
		panic(nLT0)
	case lda < max(1, m):
		panic(badLdA)
	case ldb < max(1, n):
		panic(badLdB)
	case ldc < max(1, m):
		panic(badLdC)
	case ldd < max(1, m):
		panic(badLdD)
	case lde < max(1, n):
		panic(badLdE)
	case ldf < max(1, m):
		panic(badLdF)
	case len(iwork) < m+n+2:
		panic(badLWork)
	}
	ldz := 8
	z := make([]float64, ldz*ldz)
	ipiv, jpiv := make([]int, ldz), make([]int, ldz)
	rhs := make([]float64, ldz)
	// Determine block structure of A.
	p := 0

	for i := 0; i < m; p++ {
		iwork[p] = i
		if i == m-1 {
			break
		}
		if a[(i+1)*lda+i] != 0 {
			i += 2
		} else {
			i++
		}
	}
	iwork[p+1] = m + 1

	// determine block structure of B.
	q := p + 1
	for j := 0; j < n; {
		q++
		iwork[q] = j
		if j == n-1 {
			break
		}
		if b[(j+1)*ldb+j] != 0 {
			j += 2
		} else {
			j++
		}
	}
	iwork[q+1] = n
	pq = p * (q - p - 1)
	bi := blas64.Implementation()
	// Solve (I, J) - subsystem
	//  A(I, I) * R(I, J) - L(I, J) * B(J, J) = C(I, J)
	//  D(I, I) * R(I, J) - L(I, J) * E(J, J) = F(I, J)
	// for I = P-1, P - 2, ..., 0; J = 0, 1, ..., Q - 1
	// NO TRANS PART
	// if trans == blas.NoTrans
	scale = 1
	scaloc := 1.
	var info int
	var alpha float64
	for j := p + 2; j < q; {
		js := iwork[j]
		jsp1 := js + 1
		je := iwork[j+1] - 1
		nb := je - js + 1
		for i := p; i >= 0; i-- {
			is := iwork[i]
			isp1 := is + 1
			ie := iwork[1] - 1
			mb := ie - is + 1
			zdim := mb * nb * 2
			if mb == 1 && nb == 1 {
				// Build a 2-by-2 system Z * x = RHS.
				z[0] = a[is*lda+is]
				z[ldz] = d[is*ldd+is]
				z[1] = -b[js*ldb+js]
				z[ldz+1] = -e[js*lde+js]
				// Set up right hand side(s).
				rhs[0] = c[is*ldc+js]
				rhs[1] = f[is*ldc+js]

				// Solve Z * x = RHS.
				k := impl.Dgetc2(zdim, z, ldz, ipiv, jpiv)
				if k > -1 {
					info = k
				}
				if ijob == 0 {
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv, jpiv)
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}
				} else {
					rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv, jpiv)
				}
				// Unpack solution vector(s).
				c[is*ldc+js] = rhs[0]
				f[is*ldf+js] = rhs[1]

				// Substitute R(I, J) and L(I, J) into remaining equation.
				if i > 0 {
					alpha = -rhs[0]
					bi.Daxpy(is, alpha, a[is:], lda, c[js:], ldc)
					bi.Daxpy(is, alpha, d[is:], ldd, f[js:], ldf)
				}
				if j < q {
					bi.Daxpy(n-je, rhs[1], b[js*ldb+je+1:], 1, c[is*ldc+je+1:], 1)
					bi.Daxpy(n-je, rhs[1], e[js*lde+je+1:], 1, f[is*ldf+je+1:], 1)
				}
			} else if mb == 1 && nb == 2 {
				// Build a 4-by-4 system Z * x = RHS
				z[0] = a[is*lda+is]
				z[ldz] = 0
				z[2*ldz] = d[is*ldd+is]
				z[3*ldz] = 0

				z[1] = 0
				z[ldz+1] = a[is*lda+is]
				z[2*ldz+1] = 0
				z[3*ldz+1] = d[is*ldd+is]

				z[2] = -b[js*ldb+js]
				z[ldz+2] = -b[js*ldb+jsp1]
				z[2*ldz+2] = -e[js*lde+js]
				z[3*ldz+2] = -e[js*lde+jsp1]

				z[3] = -b[jsp1*ldb+js]
				z[ldz+3] = -b[jsp1*ldb+jsp1]
				z[2*ldz+3] = 0
				z[3*ldz+3] = -e[jsp1*lde+jsp1]

				// Set up right hand side(s)
				rhs[0] = c[is*ldc+js]
				rhs[1] = c[is*ldc+jsp1]
				rhs[2] = f[is*ldf+js]
				rhs[3] = f[is*ldf+jsp1]

				// Solve Z * x = RHS
				k := impl.Dgetc2(zdim, z, ldz, ipiv, jpiv)
				if k > -1 {
					info = k
				}
				if ijob == 0 {
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv, jpiv)
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}
				} else {
					rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv, jpiv)
				}

				// Unpack solution vector(s).
				c[is*ldc+js] = rhs[0]
				c[is*ldc+jsp1] = rhs[1]
				f[is*ldf+js] = rhs[2]
				f[is*ldf+jsp1] = rhs[3]

				// Substitute R(i,j) and L(i,j) into remaining equation.
				if i > 1 {
					bi.Dger(is, nb, -1, a[is:], lda, rhs, 1, c[js:], 1)
					bi.Dger(is, nb, -1, d[is:], ldd, rhs, 1, f[js:], 1)
				}

				if j < q {
					bi.Daxpy(n-je, rhs[3], b[js*ldb+je+1:], 1, c[is*ldc+je+1:], 1)
					bi.Daxpy(n-je, rhs[3], e[js*lde+je+1:], 1, f[is*ldf+je+1:], 1)

					bi.Daxpy(n-je, rhs[4], b[jsp1*ldb+je+1:], 1, c[is*ldc+je+1:], 1)
					bi.Daxpy(n-je, rhs[4], e[jsp1*lde+je+1:], 1, f[is*ldf+je+1:], 1)
				}
			} else if mb == 2 && nb == 1 {
				// Build a 4x4 system Z * x = RHS.
				z[0] = a[is*lda+is]
				z[ldz] = a[isp1*lda+is]
				z[2*ldz] = d[is*ldd+is]
				z[3*ldz] = 0

				z[1] = a[is*lda+isp1]
				z[ldz+1] = a[isp1*lda+isp1]
				z[2*ldz+1] = d[is*ldd+isp1]
				z[3*ldz+1] = d[isp1*ldd+isp1]

				z[2] = -b[js*ldb+js]
				z[ldz+2] = 0
				z[2*ldz+2] = -e[js*lde+js]
				z[3*ldz+2] = 0

				z[3] = 0
				z[ldz+3] = -b[js*ldb+js]
				z[2*ldz+3] = 0
				z[3*ldz+3] = -e[js*lde+js]

				// Set up right hand side(s).
				rhs[0] = c[is*ldc+js]
				rhs[1] = c[isp1*ldc+js]
				rhs[2] = f[is*ldf+js]
				rhs[3] = f[isp1*ldf+js]

				// Solve Z * x = RHS
				k := impl.Dgetc2(zdim, z, ldz, ipiv, jpiv)
				if k > -1 {
					info = k
				}

				if ijob == 0 {
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv, jpiv)
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}
				} else {
					rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv, jpiv)
				}

				// Unpack solution vectors
				c[is*ldc+js] = rhs[0]
				c[isp1*ldc+js] = rhs[1]
				f[is*ldf+js] = rhs[2]
				f[isp1*ldf+js] = rhs[3]

				// Substitute R(I, J) and L(I, J) into remaining equation.
				if i > 0 {
					bi.Dgemv(blas.NoTrans, is, mb, -1, a[is:], lda, rhs, 1, 1, c[js:], ldc)
					bi.Dgemv(blas.NoTrans, is, mb, -1, d[is:], ldd, rhs, 1, 1, f[js:], ldf)
				}

			} else if mb == 2 && nb == 2 {
				// Build 8x8 system Z * x = RHS
				impl.Dlaset('F', ldz, ldz, 0, 0, z, ldz)

				z[0] = a[is*lda+is]
				z[ldz] = a[isp1*lda+is]
				z[4*ldz] = d[is*ldd+is]

				z[1] = a[is*lda+isp1]
				z[ldz+1] = a[isp1*lda+isp1]
				z[4*ldz+1] = d[is*ldd+isp1]
				z[5*ldz+1] = d[isp1*ldd+isp1]

				z[2*ldz+2] = a[is*lda+is]
				z[3*ldz+2] = a[isp1*lda+is]
				z[6*ldz+2] = d[is*ldd+is]

				z[2*ldz+3] = a[is*lda+isp1]
				z[3*ldz+3] = a[isp1*lda+isp1]
				z[6*ldz+3] = d[is*ldd+isp1]
				z[7*ldz+3] = d[isp1*ldd+isp1]

				z[4] = -b[js*ldb+js]
				z[2*ldz+4] = -b[js*ldb+jsp1]
				z[4*ldz+4] = -e[js*lde+js]
				z[6*ldz+4] = -e[js*lde+jsp1]

				z[ldz+5] = -b[js*ldb+js]
				z[3*ldz+5] = -b[js*ldb+jsp1]
				z[5*ldz+5] = -e[js*lde+js]
				z[7*ldz+5] = -e[js*lde+jsp1]

				z[6] = -b[jsp1*ldb+js]
				z[2*ldz+6] = -b[jsp1*ldb+jsp1]
				z[6*ldz+6] = -e[jsp1*lde+jsp1]

				z[ldz+7] = -b[jsp1*ldb+js]
				z[3*ldz+7] = -b[jsp1*ldb+jsp1]
				z[7*ldz+7] = -e[jsp1*lde+jsp1]

				k := 0
				ii := mb*nb + 1
				for jj := 0; jj < nb-1; jj++ {
					bi.Dcopy(mb, c[is*ldc+js+jj:], ldc, rhs[k:], 1)
					bi.Dcopy(mb, f[is*ldf+js+jj:], ldf, rhs[ii:], 1)
					k += mb
					ii += mb
				}

				// Solve Z * x = RHS
				k = impl.Dgetc2(zdim, z, ldz, ipiv, jpiv)
				if k > -1 {
					info = k
				}
				if ijob == 0 {
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv, jpiv)
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}
				} else {
					rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv, jpiv)
				}

				// Unpack solution vectors(s).
				k = 0
				ii = mb*nb + 1
				for jj := 0; jj < nb-1; jj++ {
					bi.Dcopy(mb, rhs[k:], 1, c[is*ldc+js+jj:], ldc)
					bi.Dcopy(mb, rhs[ii:], 1, f[is*ldf+js+jj:], ldf)
					k += mb
					ii += mb
				}

				// Substitute R(I, J) and L(I, J) into remaining equation.
				if i > 1 {
					bi.Dgemm(blas.NoTrans, blas.NoTrans, is, nb, mb, -1,
						a[is:], lda, rhs, mb, 1, c[js:], ldc)
					bi.Dgemm(blas.NoTrans, blas.NoTrans, is, nb, mb, -1,
						d[is:], ldd, rhs, mb, 1, f[js:], ldf)
				}

				if j < q {
					k = mb*nb + 1
					bi.Dgemm(blas.NoTrans, blas.NoTrans, mb, n-je, nb, 1,
						rhs[k:], mb, b[js*ldb+je+1:], ldb, 1, c[is*ldc+je+1:], ldc)
					bi.Dgemm(blas.NoTrans, blas.NoTrans, mb, n-je, nb, 1,
						rhs[k:], mb, e[js*lde+je+1:], lde, 1, f[is*ldf+je+1:], ldf)
				}
			}
		}
	}

	// Solve (I, J) - subsystem
	//  A(I, I) * R(I, J) - L(I, J) * B(J, J) = C(I, J)
	//  D(I, I) * R(I, J) - L(I, J) * E(J, J) = F(I, J)
	// for I = P - 1, P - 2, ..., 0; J = 0, 1, ..., Q - 1
	return scale, rdsum, rdscal, pq
}
