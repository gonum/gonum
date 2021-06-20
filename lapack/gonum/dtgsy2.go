// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Dtgsy2 solves the generalized Sylvester equation:
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
//  Z = [ kron(In, A)  -kron(Bᵀ, Im) ]             (2)
//      [ kron(In, D)  -kron(Eᵀ, Im) ],
// Ik is the identity matrix of size k and Xᵀ is the transpose of X.
// kron(X, Y) is the Kronecker product between the matrices X and Y.
// In the process of solving (1), we solve a number of such systems
// where Dim(In), Dim(In) = 1 or 2.
// If trans = blas.Trans, solve the transposed system Zᵀ*y = scale*b for y,
// which is equivalent to solve for R and L in
//  Aᵀ * R  + Dᵀ * L   = scale * C           (3)
//  R  * Bᵀ + L  * Eᵀ  = scale * -F
// This case is used to compute an estimate of Dif[(A, D), (B, E)] =
// sigma_min(Z) using reverse communication with Dlacon.
// Dtgsy2 also (ijob >= 1) contributes to the computation in Dtgsyl
// of an upper bound on the separation between to matrix pairs. Then
// the input (A, D), (B, E) are sub-pencils of the matrix pair in
// Dtgsyl. See Dtgsyl for details.
//
// Dtgsy2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dtgsy2(trans blas.Transpose, ijob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, rdsum, rdscal float64, iwork []int) (scale, sumout, scalout float64, pq, info int) {
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
	for i := 0; i < ldz; i++ {
		ipiv[i] = -1
		jpiv[i] = -1
	}
	rhs := make([]float64, ldz)

	var p, q, k int // Index variables.
	// Determine block structure of A.
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
	iwork[p+1] = m

	// Determine block structure of B.
	q = p + 1
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
	pq = (p + 1) * (q - p - 1)

	// Solve (i, j) - subsystem
	//  A(i, i) * R(i, j) - L(i, j) * B(j, j) = C(i, j)
	//  D(i, i) * R(i, j) - L(i, j) * E(j, j) = F(i, j)
	// for i = p-1, p - 2, ..., 0; j = 0, 1, ..., q - 1
	bi := blas64.Implementation()
	scale = 1
	scaloc := 1.0
	var alpha float64
	var nb, mb int // Length variables.
	if trans == blas.NoTrans {
		for j := p + 2; j < q; j++ {
			js := iwork[j]
			jsp1 := js + 1
			je := iwork[j+1] - 1
			nb = je - js + 1
			for i := p; i >= 0; i-- {
				is := iwork[i]
				isp1 := is + 1
				ie := iwork[i+1] - 1
				mb = ie - is + 1
				zdim := mb * nb * 2
				if mb == 1 && nb == 1 {
					// Build a 2×2 system Z * x = RHS.
					z[0] = a[is*lda+is]
					z[ldz] = d[is*ldd+is]
					z[1] = -b[js*ldb+js]
					z[ldz+1] = -e[js*lde+js]

					// Set up right hand side(s).
					rhs[0] = c[is*ldc+js]
					rhs[1] = f[is*ldf+js]

					// Solve Z * x = RHS.
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					if ijob == 0 {
						scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
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
						bi.Daxpy(n-je-1, rhs[1], b[js*ldb+je+1:], 1, c[is*ldc+je+1:], 1)
						bi.Daxpy(n-je-1, rhs[1], e[js*lde+je+1:], 1, f[is*ldf+je+1:], 1)
					}
				} else if mb == 1 && nb == 2 {
					// Build a 4×4 system Z * x = RHS
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
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					if ijob == 0 {
						scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
						if scaloc != 1 {
							for k = 0; k < n; k++ {
								bi.Dscal(m, scaloc, c[k:], ldc)
								bi.Dscal(m, scaloc, f[k:], ldf)
							}
							scale *= scaloc
						}
					} else {
						rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv[:zdim], jpiv[:zdim])
					}

					// Unpack solution vector(s).
					c[is*ldc+js] = rhs[0]
					c[is*ldc+jsp1] = rhs[1]
					f[is*ldf+js] = rhs[2]
					f[is*ldf+jsp1] = rhs[3]

					// Substitute R(i,j) and L(i,j) into remaining equation.
					if i > 0 {
						bi.Dger(is, nb, -1, a[is:], lda, rhs, 1, c[js:], 1)
						bi.Dger(is, nb, -1, d[is:], ldd, rhs, 1, f[js:], 1)
					}
					if j < q {
						bi.Daxpy(n-je-1, rhs[2], b[js*ldb+je+1:], 1, c[is*ldc+je+1:], 1)
						bi.Daxpy(n-je-1, rhs[2], e[js*lde+je+1:], 1, f[is*ldf+je+1:], 1)
						bi.Daxpy(n-je-1, rhs[3], b[jsp1*ldb+je+1:], 1, c[is*ldc+je+1:], 1)
						bi.Daxpy(n-je-1, rhs[3], e[jsp1*lde+je+1:], 1, f[is*ldf+je+1:], 1)
					}
				} else if mb == 2 && nb == 1 {
					// Build a 4×4 system Z * x = RHS.
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
					k = impl.Dgetc2(zdim, z, ldz, ipiv, jpiv)
					if k >= 0 {
						info = k
					}
					if ijob == 0 {
						scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
						if scaloc != 1 {
							for k = 0; k < n; k++ {
								bi.Dscal(m, scaloc, c[k:], ldc)
								bi.Dscal(m, scaloc, f[k:], ldf)
							}
							scale *= scaloc
						}
					} else {
						rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv[:zdim], jpiv[:zdim])
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
					if j < q {
						bi.Dger(mb, n-je-1, 1, rhs[2:], 1, b[js*ldb+je+1:], 1,
							c[is*ldc+je+1:], ldc)
						bi.Dger(mb, n-je-1, 1, rhs[2:], 1, e[js*lde+je+1:], 1,
							f[is*ldf+je+1:], ldf)
					}
				} else if mb == 2 && nb == 2 {
					// Build 8×8 system Z * x = RHS
					impl.Dlaset(blas.All, ldz, ldz, 0, 0, z, ldz)

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

					// Set up right hand side(s).
					k = 0
					ii := mb * nb
					for jj := 0; jj < nb-1; jj++ {
						bi.Dcopy(mb, c[is*ldc+js+jj:], ldc, rhs[k:], 1)
						bi.Dcopy(mb, f[is*ldf+js+jj:], ldf, rhs[ii:], 1)
						k += mb
						ii += mb
					}

					// Solve Z * x = RHS
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					if ijob == 0 {
						scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
						if scaloc != 1 {
							for k = 0; k < n; k++ {
								bi.Dscal(m, scaloc, c[k:], ldc)
								bi.Dscal(m, scaloc, f[k:], ldf)
							}
							scale *= scaloc
						}
					} else {
						rdsum, rdscal = impl.Dlatdf(ijob, zdim, z, ldz, rhs, rdsum, rdscal, ipiv[:zdim], jpiv[:zdim])
					}

					// Unpack solution vectors(s).
					k = 0
					ii = mb * nb
					for jj := 0; jj < nb-1; jj++ {
						bi.Dcopy(mb, rhs[k:], 1, c[is*ldc+js+jj:], ldc)
						bi.Dcopy(mb, rhs[ii:], 1, f[is*ldf+js+jj:], ldf)
						k += mb
						ii += mb
					}

					// Substitute R(i, j) and L(i, j) into remaining equation.
					if i > 0 {
						bi.Dgemm(blas.NoTrans, blas.NoTrans, is, nb, mb, -1,
							a[is:], lda, rhs, mb, 1, c[js:], ldc)
						bi.Dgemm(blas.NoTrans, blas.NoTrans, is, nb, mb, -1,
							d[is:], ldd, rhs, mb, 1, f[js:], ldf)
					}
					if j < q {
						k = mb * nb
						bi.Dgemm(blas.NoTrans, blas.NoTrans, mb, n-je-1, nb, 1,
							rhs[k:], mb, b[js*ldb+je+1:], ldb, 1, c[is*ldc+je+1:], ldc)
						bi.Dgemm(blas.NoTrans, blas.NoTrans, mb, n-je-1, nb, 1,
							rhs[k:], mb, e[js*lde+je+1:], lde, 1, f[is*ldf+je+1:], ldf)
					}
				}
			}
		}
	} else {
		// trans == blas.Trans
		// Solve (i, j) - subsystem
		// 		A(i, i)ᵀ * R(i, j) + D(i, i)ᵀ * L(j, j)  =  C(i, j)
		// 		R(i, i)  * B(j, j) + L(i, j)  * E(j, j)  = -F(i, j)
		//    for i = 0, 1, ..., P-1, j = Q-1, Q - 2, ..., 0
		var alpha float64
		for i := 0; i < p; i++ {
			is := iwork[i]
			isp1 := is + 1
			ie := iwork[i+1] - 1
			mb = ie - is + 1 // mb is a length variable
			for j := q; j < p+2; j-- {
				js := iwork[j]
				jsp1 := js + 1
				je := iwork[j+1] - 1
				nb = je - js + 1 // nb is a length variable
				zdim := nb * mb * 2
				if nb == 1 && mb == 1 {
					// Build a 2×2 system Zᵀ * x = RHS.
					z[0] = a[is*lda+is]
					z[ldz] = -b[js*ldb+js]
					z[1] = d[is*ldd+is]
					z[ldz+1] = -e[js*lde+js]
					// Set up right hand side(s).
					rhs[0] = c[is*ldc+js]
					rhs[1] = f[is*ldf+js]

					// Solve Zᵀ * x = RHS.
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}

					// Unpack solution vector(s).
					c[is*ldc+js] = rhs[0]
					f[is*ldf+js] = rhs[1]

					// Substitute R(i, j) and L(i, j) into remaining equation.
					if j > p+2 {
						alpha = rhs[0]
						bi.Daxpy(js, alpha, b[js:], ldb, f[is*ldf:], 1)
						alpha = rhs[1]
						bi.Daxpy(js, alpha, e[js:], lde, f[is*ldf:], 1)
					}
					if i < p {
						alpha = -rhs[0]
						bi.Daxpy(m-ie-1, alpha, a[is*lda+ie+1:], 1, c[(ie+1)*ldc+js:], ldc)
						alpha = -rhs[1]
						bi.Daxpy(m-ie-1, alpha, d[is*ldd+ie+1:], 1, c[(ie+1)*ldf+js:], ldc)
					}

				} else if mb == 1 && nb == 2 {
					// Build a 4×4 system Zᵀ * x = RHS
					z[0] = a[is*lda+is]
					z[ldz] = 0
					z[2*ldz] = -b[js*ldb+js]
					z[3*ldz] = -b[jsp1*ldb+js]

					z[1] = 0
					z[ldz+1] = a[is*lda+is]
					z[2*ldz+1] = -b[js*ldb+jsp1]
					z[3*ldz+1] = -b[jsp1*ldb+jsp1]

					z[2] = d[is*ldd+is]
					z[ldz+2] = 0
					z[2*ldz+2] = -e[js*lde+js]
					z[3*ldz+2] = 0

					z[3] = 0
					z[ldz+3] = d[is*ldd+is]
					z[2*ldz+3] = -e[js*lde+jsp1]
					z[3*ldz+3] = -e[jsp1*lde+jsp1]

					// Set up right hand side(s)
					rhs[0] = c[is*ldc+js]
					rhs[1] = c[is*ldc+jsp1]
					rhs[2] = f[is*ldf+js]
					rhs[3] = f[is*ldf+jsp1]

					// Solve Zᵀ * x = RHS
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}

					// Unpack solution vector(s).
					c[is*ldc+js] = rhs[0]
					c[is*ldc+jsp1] = rhs[1]
					f[is*ldf+js] = rhs[2]
					f[is*ldf+jsp1] = rhs[3]

					// Substitute R(i,j) and L(i,j) into remaining equation.
					if j > p+2 {
						bi.Daxpy(js, rhs[0], b[js:], ldb, f[is*ldf:], 1)
						bi.Daxpy(js, rhs[1], b[jsp1:], ldb, f[is*ldf:], 1)
						bi.Daxpy(js, rhs[2], e[js:], lde, f[is*ldf:], 1)
						bi.Daxpy(js, rhs[3], e[jsp1:], lde, f[is*ldf:], 1)
					}
					if i < p {
						bi.Dger(m-ie-1, nb, -1, a[is*lda+ie+1:], 1, rhs, 1, c[(ie+1)*ldc+js:], 1)
						bi.Dger(m-ie-1, nb, -1, d[is*ldd+ie+1:], 1, rhs[2:], 1, c[(ie+1)*ldc+js:], 1)
					}
				} else if mb == 2 && nb == 1 {
					// Build a 4×4 system Zᵀ * x = RHS.
					z[0] = a[is*lda+is]
					z[ldz] = a[is*lda+isp1]
					z[2*ldz] = -b[js*ldb+js]
					z[3*ldz] = 0

					z[1] = a[isp1*lda+is]
					z[ldz+1] = a[isp1*lda+isp1]
					z[2*ldz+1] = 0
					z[3*ldz+1] = -b[js*ldb+js]

					z[2] = d[is*ldd+is]
					z[ldz+2] = d[is*ldd+isp1]
					z[2*ldz+2] = -e[js*lde+js]
					z[3*ldz+2] = 0

					z[3] = 0
					z[ldz+3] = d[isp1*ldb+isp1]
					z[2*ldz+3] = 0
					z[3*ldz+3] = -e[js*lde+js]

					// Set up right hand side(s).
					rhs[0] = c[is*ldc+js]
					rhs[1] = c[isp1*ldc+js]
					rhs[2] = f[is*ldf+js]
					rhs[3] = f[isp1*ldf+js]

					// Solve Zᵀ * x = RHS
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}

					// Unpack solution vectors
					c[i*ldc+js] = rhs[0]
					c[isp1*ldc+js] = rhs[1]
					f[is*ldf+js] = rhs[2]
					f[isp1*ldf+js] = rhs[3]

					// Substitute R(i, j) and L(i, j) into remaining equation.
					if j > p+2 {
						bi.Dger(mb, js, 1, rhs, 1, b[js:], ldb, f[is*ldf:], 1)
						bi.Dger(mb, js, 1, rhs[2:], 1, e[js:], lde, f[is*ldf:], 1)
					}
					if i < p {
						bi.Dgemv(blas.Trans, mb, m-ie-1, -1, a[is*lda+ie+1:], 1, rhs, 1,
							1, c[(ie+1)*lda+js:], ldc)
						bi.Dgemv(blas.Trans, mb, m-ie-1, -1, d[is*ldd+ie+1:], 1, rhs[2:], 1,
							1, c[(ie+1)*lda+js:], ldc)
					}

				} else if mb == 2 && nb == 2 {
					// Build 8×8 system Zᵀ * x = RHS
					impl.Dlaset(blas.All, ldz, ldz, 0, 0, z, ldz)

					z[0] = a[is*lda+is]
					z[ldz] = a[is*lda+isp1]
					z[4*ldz] = -b[js*ldb+js]
					z[6*ldz] = -b[jsp1*ldb+js]

					z[1] = a[isp1*lda+is]
					z[ldz+1] = a[isp1*lda+isp1]
					z[5*ldz+1] = -b[js*ldb+js]
					z[7*ldz+1] = -b[jsp1*ldb+js]

					z[2*ldz+2] = a[is*lda+is]
					z[3*ldz+2] = a[is*lda+isp1]
					z[4*ldz+2] = -b[js*ldb+jsp1]
					z[6*ldz+2] = -b[jsp1*ldb+jsp1]

					z[2*ldz+3] = a[isp1*lda+is]
					z[3*ldz+3] = a[isp1*lda+isp1]
					z[5*ldz+3] = -b[js*ldb+jsp1]
					z[7*ldz+3] = -b[jsp1*ldb+jsp1]

					z[4] = d[is*ldd+is]
					z[1*ldz+4] = d[is*ldd+isp1]
					z[4*ldz+4] = -e[js*lde+js]

					z[ldz+5] = d[isp1*ldd+isp1]
					z[5*ldz+5] = -e[js*lde+js]

					z[2*ldz+6] = d[is*ldd+is]
					z[3*ldz+6] = d[is*ldd+isp1]
					z[4*ldz+6] = -e[js*lde+jsp1]
					z[6*ldz+6] = -e[jsp1*lde+jsp1]

					z[3*ldz+7] = d[isp1*ldd+isp1]
					z[5*ldz+7] = -e[js*lde+jsp1]
					z[7*ldz+7] = -e[jsp1*lde+jsp1]

					// Set up right hand side(s).
					k = 0
					ii := nb * mb
					for jj := 0; jj < nb-1; jj++ {
						bi.Dcopy(mb, c[is*ldc+js+jj:], ldc, rhs[k:], 1)
						bi.Dcopy(mb, f[is*ldf+js+jj:], ldf, rhs[ii:], 1)
						k += nb
						ii += nb
					}

					// Solve Zᵀ * x = RHS
					k = impl.Dgetc2(zdim, z, ldz, ipiv[:zdim], jpiv[:zdim])
					if k >= 0 {
						info = k
					}
					scaloc = impl.Dgesc2(zdim, z, ldz, rhs, ipiv[:zdim], jpiv[:zdim])
					if scaloc != 1 {
						for k = 0; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scale *= scaloc
					}

					// Unpack solution vectors(s).
					k = 0
					ii = mb * nb
					for jj := 0; jj < nb-1; jj++ {
						bi.Dcopy(mb, rhs[k:], 1, c[is*ldc+js+jj:], ldc)
						bi.Dcopy(mb, rhs[ii:], 1, f[is*ldf+js+jj:], ldf)
						k += mb
						ii += mb
					}

					// Substitute R(i, j) and L(i, j) into remaining equation.
					if j > p+2 {
						bi.Dgemm(blas.NoTrans, blas.Trans, mb, js, nb, 1,
							c[is*ldc+js:], ldc, b[js:], ldb, 1, f[is*ldf:], ldf)
						bi.Dgemm(blas.NoTrans, blas.Trans, mb, js, nb, 1,
							f[is*ldf+js:], ldf, e[js:], lde, 1, f[is*ldf:], ldf)
					}
					if i < p {
						bi.Dgemm(blas.Trans, blas.NoTrans, m-ie-1, nb, mb, -1,
							a[is*lda+ie+1:], lda, c[is*ldc+js:], ldc, 1, c[(ie+1)*ldc+js:], ldc)
						bi.Dgemm(blas.Trans, blas.NoTrans, m-ie-1, nb, mb, -1,
							d[is*ldd+ie+1:], ldd, f[is*ldf+js:], ldf, 1, c[(ie+1)*ldc+js:], ldc)
					}
				}
			}
		}

	}

	return scale, rdsum, rdscal, pq, info
}
