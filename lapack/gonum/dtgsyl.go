package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

// Dtgsyl solves the generalized Sylvester equation:
//
//	A * R - L * B = scaleOut * C                 (1)
//	D * R - L * E = scaleOut * F
//
// where R and L are unknown m-by-n matrices, (A, D), (B, E) and
// (C, F) are given matrix pairs of size m×m, n×n and m×n,
// respectively, with real entries. (A, D) and (B, E) must be in
// generalized (real) Schur canonical form, i.e. A, B are upper quasi
// triangular and D, E are upper triangular.
//
//	The solution (R, L) overwrites (C, F). 0 <= scaleOut <= 1 is an output
//	scaling factor chosen to avoid overflow.
//
// In matrix notation (1) is equivalent to solve  Zx = scale b, where
// Z is defined as
//
//	Z = [ kron(In, A)  -kron(Bᵀ, Im) ]  (2)
//	    [ kron(In, D)  -kron(Eᵀ, Im) ].
//
// Here Ik is the identity matrix of size k and Xᵀ is the transpose of
// X. kron(X, Y) is the Kronecker product between the matrices X and Y.
//
// If trans==blas.Trans, Dtgsyl solves the transposed system Zᵀ*y = scale*b,
// which is equivalent to solve for R and L in
//
//	Aᵀ * R  + Dᵀ * L  = scaleOut * C           (3)
//	R  * Bᵀ + L  * Eᵀ = scaleOut * -F
//
// This case trans==blas.Trans is used to compute an one-norm-based estimate
// of Dif[(A,D), (B,E)], the separation between the matrix pairs (A,D)
// and (B,E), using [Implementation.Dlacon].
//
// Notes on arguments:
// - If ijob >= 1, Dtgsyl computes a Frobenius norm-based estimate
// of Dif[(A,D),(B,E)]. That is, the reciprocal of a lower bound on the
// reciprocal of the smallest singular value of Z. See [1-2] for more
// information.
//   - iwork is int array, dimension (M+N+6)
//   - work is float array, dimension max(1, lwork). On exit work[0] contains the optimal lwork.
//   - lwork is the dimension of the array work. lwork >= 1.
//     If ijob is 1 or 2 and trans==blas.NoTrans then lwork >= max(1, 2*m*n)
//   - If workspaceQuery is true then the routine
//     only calculates the optimal size of the WORK array, returns
//     this value as the first entry of the work array.
//
// Dtgsyl is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dtgsyl(trans blas.Transpose, ijob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, work []float64, iwork []int, workspaceQuery bool) (difOut, scaleOut float64, infoOut int) {
	infoOut = -1
	lwork := len(work)
	notran := trans == blas.NoTrans

	lwmin := 1
	if notran && (ijob == 1 || ijob == 2) {
		lwmin = max(1, 2*m*n)
	}

	switch {
	case !notran && trans != blas.Trans:
		panic(badTrans)
	case notran && (ijob < 0 || ijob > 4):
		panic(badIJob)
	case n < 1:
		panic(nLT1)
	case lda < max(1, m):
		panic(badLdA)
	case ldb < max(1, n):
		panic(badLdB)
	case ldc < max(1, n): // ldc and ldf are inverted w.r.t reference due to row-major storage.
		panic(badLdC)
	case ldd < max(1, m):
		panic(badLdD)
	case lde < max(1, n):
		panic(badLdE)
	case ldf < max(1, n):
		panic(badLdF)
	case lwork < lwmin && !workspaceQuery:
		panic(badLWork)
	case workspaceQuery:
		work[0] = float64(lwmin)
		return 0, 0, infoOut
	}

	work[0] = float64(lwmin)
	if m == 0 || n == 0 {
		// Early return.
		scaleOut = 1
		return 0, scaleOut, 0
	}

	// Determine optimal block sizes mb and nb.
	trs := "T"
	if notran {
		trs = "N"
	}
	mb := impl.Ilaenv(2, "DTGSYL", trs, m, n, -1, -1)
	nb := impl.Ilaenv(5, "DTGSYL", trs, m, n, -1, -1)

	isolve := 1
	ifunc := 0

	if notran && ijob >= 3 {
		ifunc = ijob - 2
		impl.Dlaset(blas.All, m, n, 0, 0, c, ldc)
		impl.Dlaset(blas.All, m, n, 0, 0, f, ldf)
	} else if notran && ijob >= 1 {
		isolve = 2
	}
	ldw := n
	var pq int
	if (mb <= 1 && nb <= 1) || (mb >= m && nb >= n) {
		for iround := 1; iround <= isolve; iround++ {
			// Use unblocked level 2 solver.
			dscale := 0.0
			dsum := 1.0
			scaleOut, dscale, dsum, pq, infoOut = impl.Dtgsy2(trans, ifunc, m, n, a, lda, b, ldb, c, ldc,
				d, ldd, e, lde, f, ldf, dsum, dscale, iwork)
			if dscale != 0 {
				divisor := dscale * math.Sqrt(dsum)
				if ijob == 1 || ijob == 3 {
					difOut = math.Sqrt(float64(2*m*n)) / divisor
				} else {
					difOut = math.Sqrt(float64(pq)) / divisor
				}
			}

			if isolve == 2 && iround == 1 {
				if notran {
					ifunc = ijob
				}
				impl.Dlacpy(blas.All, m, n, c, ldc, work, ldw)
				impl.Dlacpy(blas.All, m, n, f, ldf, work[m*n:], ldw)
				impl.Dlaset(blas.All, m, n, 0, 0, c, ldc)
				impl.Dlaset(blas.All, m, n, 0, 0, f, ldf)
			} else if isolve == 2 && iround == 2 {
				// scaleOut = scale2 // scale2 is undefined?
				impl.Dlacpy(blas.All, m, n, work, ldw, c, ldc)
				impl.Dlacpy(blas.All, m, n, work[m*n:], ldw, f, ldf)
			}
		}
		// End of unblocked level 2 solver.
		return difOut, scaleOut, infoOut
	}

	// Determine block structure of A and B.
	p := blockStructureL(iwork, mb, m, a, lda)
	if iwork[p-2] == iwork[p-1] {
		p--
	}

	q := blockStructureL(iwork[p:], nb, n, b, ldb)
	if iwork[p:][q-2] == iwork[p:][q-1] {
		q--
	}

	p -= 2
	q += p

	bi := blas64.Implementation()
	var ppqq, linfo int
	var dscale, dsum, scaloc, scale2 float64
	if notran {
		// Solve (I, J)-subsystem
		//     A(I, I) * R(I, J) - L(I, J) * B(J, J) = C(I, J)
		//     D(I, I) * R(I, J) - L(I, J) * E(J, J) = F(I, J)
		// for I = P, P - 1,..., 1; J = 1, 2,..., Q
		for iround := 1; iround <= isolve; iround++ {
			dsum = 1
			pq = 0
			scaleOut = 1
			for j := p + 2; j <= q; j++ {
				js := iwork[j]
				je := iwork[j+1] - 1
				nb := je - js + 1
				for i := p; i >= 0; i-- {
					is := iwork[i]
					ie := iwork[i+1] - 1
					mb := ie - is + 1
					scaloc, dscale, dsum, ppqq, linfo = impl.Dtgsy2(trans, ifunc, mb, nb, a[is*lda+is:], lda,
						b[js*ldb+js:], ldb, c[is*ldc+js:], ldc, d[is*ldd+is:], ldd,
						e[js*lde+js:], lde, f[is*ldf+js:], ldf, dsum, dscale, iwork[q+2:])
					if linfo > 0 {
						infoOut = linfo
					}

					pq = pq + ppqq
					if scaloc != 1 {
						for k := 0; k <= js-1; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						for k := js; k <= je; k++ {
							bi.Dscal(is, scaloc, c[k:], ldc)
							bi.Dscal(is, scaloc, f[k:], ldf)
						}
						for k := js; k <= je; k++ {
							bi.Dscal(m-ie-1, scaloc, c[(ie+1)*ldc+k:], ldc)
							bi.Dscal(m-ie-1, scaloc, f[(ie+1)*ldf+k:], ldf)
						}
						for k := je + 1; k < n; k++ {
							bi.Dscal(m, scaloc, c[k:], ldc)
							bi.Dscal(m, scaloc, f[k:], ldf)
						}
						scaleOut *= scaloc
					}

					// Substitute R(I, J) and L(I, J) into remaining equation.

					if i > 0 {
						bi.Dgemm(blas.NoTrans, blas.NoTrans, is, nb, mb, -1,
							a[is:], lda, c[is*ldc+js:], ldc, 1, c[js:], ldc)
						bi.Dgemm(blas.NoTrans, blas.NoTrans, is, nb, mb, -1,
							d[is:], ldd, c[is*ldc+js:], ldc, 1, f[js:], ldf)
					}
					if j < q {
						bi.Dgemm(blas.NoTrans, blas.NoTrans, mb, n-je-1, nb, 1,
							f[is*ldf+js:], ldf, b[js*ldb+je+1:], ldb, 1, c[is*ldc+je+1:], ldc)
						bi.Dgemm(blas.NoTrans, blas.NoTrans, mb, n-je-1, nb, 1,
							f[is*ldf+js:], ldf, e[js*lde+je+1:], lde, 1, f[is*ldf+je+1:], ldf)
					}
				}
			}
			// Continue iround loop.
			if dscale != 0 {
				divisor := dscale * math.Sqrt(dsum)
				if ijob == 1 || ijob == 3 {
					difOut = math.Sqrt(float64(2*m*n)) / divisor
				} else {
					difOut = math.Sqrt(float64(pq)) / divisor
				}
			}
			if isolve == 2 && iround == 1 {
				if notran {
					ifunc = ijob
				}
				scale2 = scaleOut
				impl.Dlacpy(blas.All, m, n, c, ldc, work, ldw)
				impl.Dlacpy(blas.All, m, n, f, ldf, work[m*n:], ldw)
				impl.Dlaset(blas.All, m, n, 0, 0, c, ldc)
				impl.Dlaset(blas.All, m, n, 0, 0, f, ldf)
			} else if isolve == 2 && iround == 2 {
				scaleOut = scale2
				impl.Dlacpy(blas.All, m, n, work, ldw, c, ldc)
				impl.Dlacpy(blas.All, m, n, work[m*n:], ldw, f, ldf)
			}
		}
		work[0] = float64(lwmin)
		// End of notran code.
		return difOut, scaleOut, infoOut
	}

	// Solve transposed (I, J)-subsystem
	//      A(I, I)ᵀ * R(I, J)  + D(I, I)ᵀ * L(I, J)  =  C(I, J)
	//      R(I, J)  * B(J, J)ᵀ + L(I, J)  * E(J, J)ᵀ = -F(I, J)
	// for I = 1,2,..., P; J = Q, Q-1,..., 1
	scaleOut = 1
	for i := 0; i <= p; i++ {
		is := iwork[i]
		ie := iwork[i+1] - 1
		mb := ie - is + 1
		for j := q; j >= p+2; j-- {
			js := iwork[j]
			je := iwork[j+1] - 1
			nb := je - js + 1
			scaloc, dscale, dsum, _, linfo = impl.Dtgsy2(trans, ifunc, mb, nb, a[is*lda+is:], lda,
				b[js*ldb+js:], ldb, c[is*ldc+js:], ldc, d[is*ldd+is:], ldd,
				e[js*lde+js:], lde, f[is*ldf+js:], ldf, dsum, dscale, iwork[q+2:])
			if linfo > 0 {
				infoOut = linfo
			}
			if scaloc != 1 {
				for k := 0; k <= js-1; k++ {
					bi.Dscal(m, scaloc, c[k:], ldc)
					bi.Dscal(m, scaloc, f[k:], ldf)
				}
				for k := js; k <= je; k++ {
					bi.Dscal(is, scaloc, c[k:], ldc)
					bi.Dscal(is, scaloc, f[k:], ldf)
				}
				for k := js; k <= je; k++ {
					bi.Dscal(m-ie-1, scaloc, c[(ie+1)*ldc+k:], ldc)
					bi.Dscal(m-ie-1, scaloc, f[(ie+1)*ldf+k:], ldf)
				}
				for k := je + 1; k < n; k++ {
					bi.Dscal(m, scaloc, c[k:], ldc)
					bi.Dscal(m, scaloc, f[k:], ldf)
				}
				scaleOut *= scaloc
			}

			// Substitute R(I, J) and L(I, J) into remaining equation.

			if j > p+2 {
				bi.Dgemm(blas.NoTrans, blas.Trans, mb, js, nb, 1, c[is*ldc+js:], ldc,
					b[js:], ldb, 1, f[is*ldf:], ldf)
				bi.Dgemm(blas.NoTrans, blas.Trans, mb, js, nb, 1, f[is*ldf+js:], ldf,
					e[js:], lde, 1, f[is*ldf:], ldf)
			}
			if i < p {
				bi.Dgemm(blas.Trans, blas.NoTrans, m-ie-1, nb, mb, -1, a[is*lda+ie+1:], lda,
					c[is*ldc+js:], ldc, 1, c[(ie+1)*ldc+js:], ldc)
				bi.Dgemm(blas.Trans, blas.NoTrans, m-ie-1, nb, mb, -1, d[is*ldd+ie+1:], ldd,
					f[is*ldf+js:], ldf, 1, c[(ie+1)*ldc+js:], ldc)
			}
		}
	}
	work[0] = float64(lwmin)
	return difOut, scaleOut, infoOut
}

// blockStructure computes the block structure of a matrix A for Dtgsy2.
// On exit, the block ranges are stored in dst[0:p] and the number of
// blocks is returned, which is equal to p, the return value of blockStructure.
func blockStructure(dst []int, blocksize, cols int, a []float64, lda int) int {
	p := -1
	for i := 0; i < cols; {
		p++
		dst[p] = i
		if i == cols-1 {
			break
		}
		if a[(i+1)*lda+i] != 0 {
			i += blocksize
		} else {
			i++
		}
	}
	dst[p+1] = cols
	return p + 2
}

// blockStructureL is like blockStructure but for Dtgsyl.
func blockStructureL(dst []int, blocksize, cols int, a []float64, lda int) int {
	p := -1
	for i := 0; i < cols; {
		p++
		dst[p] = i
		i += blocksize
		if i >= cols-1 {
			break
		}
		if a[i*lda+i-1] != 0 {
			i++
		}
	}
	dst[p+1] = cols
	return p + 2
}
