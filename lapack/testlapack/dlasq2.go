// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlasq2er interface {
	Dlasq2(n int, z []float64) (info int)

	Dsyev(jobz lapack.EVJob, uplo blas.Uplo, n int, a []float64, lda int, w, work []float64, lwork int) (ok bool)
}

func Dlasq2Test(t *testing.T, impl Dlasq2er) {
	const tol = 1e-14

	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 25, 50} {
		for k := 0; k < 10; k++ {
			for typ := 0; typ <= 2; typ++ {
				name := fmt.Sprintf("n=%v,typ=%v", n, typ)

				want := make([]float64, n)
				z := make([]float64, 4*n)
				switch typ {
				case 0:
					// L is the identity, U has zero diagonal.
				case 1:
					// L is the identity, U has random diagonal, and so T is upper triangular.
					for i := 0; i < n; i++ {
						z[2*i] = rnd.Float64()
						want[i] = z[2*i]
					}
					sort.Float64s(want)
				case 2:
					// Random tridiagonal matrix
					for i := range z {
						z[i] = rnd.Float64()
					}
					// The slice 'want' is computed below.
				}
				zCopy := make([]float64, len(z))
				copy(zCopy, z)

				// Compute the eigenvalues of the symmetric positive definite
				// tridiagonal matrix associated with the slice z.
				info := impl.Dlasq2(n, z)
				if info != 0 {
					t.Fatalf("%v: Dlasq2 failed", name)
				}

				if n == 0 {
					continue
				}

				got := z[:n]

				if typ == 2 {
					// Compute the expected result.

					// Compute the non-symmetric tridiagonal matrix T = L*U where L and
					// U are represented by the slice z.
					ldt := n
					T := make([]float64, n*ldt)
					for i := 0; i < n; i++ {
						if i == 0 {
							T[0] = zCopy[0]
						} else {
							T[i*ldt+i] = zCopy[2*i-1] + zCopy[2*i]
						}
						if i < n-1 {
							T[i*ldt+i+1] = 1
							T[(i+1)*ldt+i] = zCopy[2*i+1] * zCopy[2*i]
						}
					}
					// Compute the symmetric tridiagonal matrix by applying a similarity
					// transformation on T: D^{-1}*T*D. See discussion and references in
					//  http://icl.cs.utk.edu/lapack-forum/viewtopic.php?f=5&t=4839
					d := make([]float64, n)
					d[0] = 1
					for i := 1; i < n; i++ {
						d[i] = d[i-1] * T[i*ldt+i-1] / T[(i-1)*ldt+i]
					}
					for i, di := range d {
						d[i] = math.Sqrt(di)
					}
					for i := 0; i < n; i++ {
						// Update only the upper triangle.
						for j := i; j <= min(i+1, n-1); j++ {
							T[i*ldt+j] *= d[j] / d[i]
						}
					}

					// Compute the eigenvalues of D^{-1}*T*D by using Dsyev. It's call
					// tree doesn't include Dlasq2.
					work := make([]float64, 3*n)
					ok := impl.Dsyev(lapack.EVNone, blas.Upper, n, T, ldt, want, work, len(work))
					if !ok {
						t.Fatalf("%v: Dsyev failed", name)
					}
				}

				sort.Float64s(got)
				diff := floats.Distance(got, want, math.Inf(1))
				if diff > tol {
					t.Errorf("%v: unexpected eigenvalues; diff=%v\n%v\n%v\n\n", name, diff, got, want)
				}
			}
		}
	}
}
