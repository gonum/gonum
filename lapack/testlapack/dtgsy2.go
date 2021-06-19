// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dtgsy2er interface {
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
	Dtgsy2(trans blas.Transpose, ijob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, rdsum, rdscal float64, iwork []int) (scale, sumout, scalout float64, pq, info int)
}

func Dtgsy2Test(t *testing.T, impl Dtgsy2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{2} {
		for _, m := range []int{2} {
			for _, lda := range []int{m} {
				for _, ldb := range []int{n} {
					for _, ldc := range []int{n} {
						for _, ldd := range []int{m} {
							for _, lde := range []int{n} {
								for _, ldf := range []int{n} {
									testSolveDtgsy2(t, impl, rnd, blas.NoTrans, m, n, lda, ldb, ldc, ldd, lde, ldf, false)
									testSolveDtgsy2(t, impl, rnd, blas.NoTrans, m, n, lda, ldb, ldc, ldd, lde, ldf, true)
								}
							}
						}
					}
				}
			}
		}
	}
}

func testSolveDtgsy2(t *testing.T, impl Dtgsy2er, rnd *rand.Rand, tp blas.Transpose, m, n, lda, ldb, ldc, ldd, lde, ldf int, bad bool) {
	const tol = 1e-12
	name := fmt.Sprintf("n=%v,m=%v,lda=%v,ldb=%v,ldc=%v,ldd=%v,lde=%v,ldf=%v,bad=%t", n, m, lda, ldb, ldc, ldd, lde, ldf, bad)
	lda = min(lda, m)
	ldb = min(ldb, n)
	ldc = min(ldc, n)
	ldd = min(ldd, m)
	lde = min(lde, n)
	ldf = min(ldf, n)
	// Generate random matrices (A, D) and (B, E) which must be
	// in generalized Schur canonical form, i.e. A, B are upper
	// quasi triangular and D, E are upper triangular.
	a := randomUpperQuasiTriangular(m, m, lda, max(1, m/2), rnd)
	b := randomUpperQuasiTriangular(n, n, lda, max(1, n/2), rnd)
	d := randomUpperTriangular(m, ldd, rnd)
	e := randomUpperTriangular(n, lde, rnd)
	// Generate random general matrix.
	c := randomGeneral(m, n, ldc, rnd)
	f := randomGeneral(m, n, ldf, rnd)
	cCopy := cloneGeneral(c)
	fCopy := cloneGeneral(f)
	// rdsum and rdscal only makes sense when Dtgsy2 is called by	Dtgsyl.
	rdsum, rdscal := 0., 0.
	iwork := make([]int, m+n+2)
	scale, sum, scalout, pq, info := impl.Dtgsy2(tp, 0, m, n, a.Data, lda, b.Data, ldb, c.Data, ldc, d.Data, ldd,
		e.Data, lde, f.Data, ldf, rdsum, rdscal, iwork)
	if info != 0 {
		t.Errorf("%v: got non-zero exit number. info=%d", name, info)
	}
	if scale == 0 {
		t.Errorf("%v: unexpected homogenous system solution", name)
	}
	_, _, _, _ = scale, sum, scalout, pq
	// solutions are overwritten (R,L)->(C,F).
	r := c
	l := f

	// Calculate residuals
	// | A * R - L * B - scale * C |  from (1)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, a, r, -scale, cCopy)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, l, b, 1, cCopy)
	res := dlange(lapack.MaxColumnSum, m, n, cCopy.Data, cCopy.Stride)
	if res > tol {
		t.Errorf("%v: | A * R - L * B - scale * C | residual large %v", name, res)
	}

	// | D * R - L * E - scale * F |  from (1)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, d, r, -scale, fCopy)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, l, e, 1, fCopy)
	res = dlange(lapack.MaxColumnSum, m, n, fCopy.Data, fCopy.Stride)
	if res > tol {
		t.Errorf("%v: | D * R - L * E - scale * F | residual large %v", name, res)
	}

}

// randomUpperQuasiTriangular returns a random, upper quasi triangular matrix, which is
// to say this is a random matrix with zeros in the subarray A[k:m, 0:k].
func randomUpperQuasiTriangular(r, c, stride, k int, rnd *rand.Rand) blas64.General {
	ans := randomGeneral(r, c, stride, rnd)
	for i := k; i < r; i++ {
		for j := k - 1; j >= 0; j-- {
			ans.Data[i*ans.Stride+j] = 0
		}
	}
	return ans
}

func randomUpperTriangular(n, stride int, rnd *rand.Rand) blas64.General {
	ans := nanGeneral(n, n, stride)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i <= j {
				ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
			} else {
				ans.Data[i*ans.Stride+j] = 0
			}
		}
	}
	return ans
}
