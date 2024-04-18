// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// Dptrfs improves the computed solution to a system of linear equations when
// the coefficient matrix is symmetric positive definite and tridiagonal, and
// provides error bounds and backward error estimates for the solution.
func (impl Implementation) Dptrfs(n, nrhs int, d, e, df, ef []float64, b []float64, ldb int, x []float64, ldx int, ferr, berr []float64, work []float64) {
	switch {
	case n < 0:
		panic(nLT0)
	case nrhs < 0:
		panic(nrhsLT0)
	case ldb < max(1, nrhs):
		panic(badLdB)
	case ldx < max(1, nrhs):
		panic(badLdX)
	}

	// Quick return if possible.
	if n == 0 || nrhs == 0 {
		for j := 0; j < nrhs; j++ {
			ferr[j] = 0
			berr[j] = 0
		}
		return
	}

	switch {
	case len(d) < n:
		panic(shortD)
	case len(df) < n:
		panic(shortDF)
	case len(e) < n-1:
		panic(shortE)
	case len(ef) < n-1:
		panic(shortEF)
	case len(b) < (n-1)*ldb+nrhs:
		panic(shortB)
	case len(x) < (n-1)*ldb+nrhs:
		panic(shortX)
	case len(ferr) < nrhs:
		panic(shortFERR)
	case len(berr) < nrhs:
		panic(shortBERR)
	case len(work) < 2*n:
		panic(shortWork)
	}

	const (
		iterMax = 5
		nz      = 4 // Maximum number of nonzero elements in each row of A, plus 1
		eps     = dlamchE
		safmin  = dlamchS
		safe1   = nz * safmin
		safe2   = safe1 / eps
	)

	bi := blas64.Implementation()
	res := work[n : 2*n]
	for j := 0; j < nrhs; j++ {
		iter := 1
		lastres := 3.0
		// Loop until stopping criterion is satisfied.
		for {
			// Compute residual R = B - A * X. Also compute abs(A)*abs(x) +
			// abs(b) for use in the backward error bound.
			if n == 1 {
				bi := b[j]
				dx := d[0] * x[j]
				res[0] = bi - dx
				work[0] = math.Abs(bi) + math.Abs(dx)
			} else {
				bi := b[j]
				dx := d[0] * x[j]
				ex := e[0] * x[ldx+j]
				res[0] = bi - dx - ex
				work[0] = math.Abs(bi) + math.Abs(dx) + math.Abs(ex)
				for i := 1; i < n-1; i++ {
					bi = b[i*ldb+j]
					cx := e[i-1] * x[(i-1)*ldx+j]
					dx = d[i] * x[i*ldx+j]
					ex = e[i] * x[(i+1)*ldx+j]
					res[i] = bi - cx - dx - ex
					work[i] = math.Abs(bi) + math.Abs(cx) + math.Abs(dx) + math.Abs(ex)
				}
				bi = b[(n-1)*ldb+j]
				cx := e[n-2] * x[(n-2)*ldx+j]
				dx = d[n-1] * x[(n-1)*ldx+j]
				res[n-1] = bi - cx - dx
				work[n-1] = math.Abs(bi) + math.Abs(cx) + math.Abs(dx)
			}

			// Compute componentwise relative backward error from formula
			//
			//  max over i (abs(R[i])/(abs(A)*abs(X) + abs(B))[i])
			//
			// where abs(Z) is the componentwise absolute value of the matrix or
			// vector Z. If the i-th component of the denominator is less than
			// safe2, then safe1 is added to the i-th components of the
			// numerator and denominator before dividing.
			var s float64
			for i, worki := range work[:n] {
				if worki > safe2 {
					s = math.Max(s, math.Abs(res[i])/worki)
				} else {
					s = math.Max(s, (math.Abs(res[i])+safe1)/(worki+safe1))
				}
			}
			berr[j] = s

			// Test stopping criterion. Continue iterating if
			//  1. The residual berr[j] is larger than machine epsilon, and
			//  2. berr[j] decreased by at least a factor of 2 during the last iteration, and
			//  3. At most itmax iterations tried.
			if berr[j] <= eps || 2*berr[j] > lastres || iter > iterMax {
				break
			}

			// Update solution and try again.
			impl.Dpttrs(n, 1, df, ef, res, 1)
			bi.Daxpy(n, 1, res, 1, x[j:], ldx)
			lastres = berr[j]
			iter++
		}

		// Bound error from formula
		//
		// 	norm(X - XTRUE)/norm(X) <= ferr = norm(abs(inv(A)) * (abs(R) + nz*eps*(abs(A)*abs(X)+abs(B))))/norm(X)
		//
		// where
		//   norm(Z) is the magnitude of the largest component of Z
		//   inv(A) is the inverse of A
		//   abs(Z) is the componentwise absolute value of the matrix or vector Z
		//   nz is the maximum number of nonzeros in any row of A, plus 1
		//   eps is machine epsilon
		//
		// The i-th component of abs(R)+nz*eps*(abs(A)*abs(X)+abs(B)) is
		// incremented by safe1 if the i-th component of abs(A)*abs(X) + abs(B)
		// is less than safe2.
		for i, worki := range work[:n] {
			if worki > safe2 {
				work[i] = math.Abs(work[n+i]) + nz*eps*worki
			} else {
				work[i] = math.Abs(work[n+i]) + nz*eps*worki + safe1
			}
		}
		ix := bi.Idamax(n, work, 1)
		ferr[j] = work[ix]

		// Estimate the norm of inv(A).
		//
		// Solve M(A) * x = e, where M(A) = (m[i,j]) is given by
		//
		// 	m[i,j] =  abs(A[i,j]), i == j,
		// 	m[i,j] = -abs(A[i,j]), i != j,
		//
		// and e = [1,1,...,1]ᵀ. Note M(A) = M(L)*D*M(L)ᵀ.
		//
		// Solve M(L) * b = e.
		work[0] = 1
		for i := 1; i < n; i++ {
			work[i] = 1 + work[i-1]*math.Abs(ef[i-1])
		}
		// Solve D * M(L)ᵀ * x = b.
		work[n-1] /= d[n-1]
		for i := n - 2; i >= 0; i-- {
			work[i] = work[i]/df[i] + work[i+1]*math.Abs(ef[i])
		}
		// Compute norm(inv(A)) = max(x(i)), 0<=i<n.
		ix = bi.Idamax(n, work, 1)
		ferr[j] *= math.Abs(work[ix])

		// Normalize error.
		lastres = 0
		for i := 0; i < n; i++ {
			lastres = math.Max(lastres, math.Abs(x[i*ldx+j]))
		}
		if lastres != 0 {
			ferr[j] /= lastres
		}
	}
}
