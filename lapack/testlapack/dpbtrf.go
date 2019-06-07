// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dpbtrfer interface {
	Dpbtrf(uplo blas.Uplo, n, kd int, ab []float64, ldab int) (ok bool)
}

// DpbtrfTest tests a band Cholesky factorization on random symmetric positive definite
// band matrices by checking that the Cholesky factors multiply back to the original matrix.
func DpbtrfTest(t *testing.T, impl Dpbtrfer) {
	const tol = 1e-12

	rnd := rand.New(rand.NewSource(1))

	// The values of n and kd are chosen to assure that the blocked code path is taken.
	// With the current implementation of Ilaenv this happens if kd > 64.
	// Unfortunately, with the block size nb=32 this also means that in Dpbtrf
	// it never happens that i2<=0.
	for _, n := range []int{0, 1, 2, 3, 4, 5, 64, 65, 66, 91, 96, 97, 101, 128, 130} {
		for _, kd := range []int{0, (5*n + 1) / 4, (3*n - 1) / 4, (n + 1) / 4} {
			if kd+1 > n && n != 0 && kd != 0 {
				continue
			}
			for _, uplo := range []blas.Uplo{blas.Upper} {
				for _, ldextra := range []int{0, 7} {
					ldab := kd + 1 + ldextra
					name := fmt.Sprintf("uplo=%v,n=%v,kd=%v,ldab=%v", uplo, n, kd, ldab)

					// Allocate a band symmetric matrix A and fill it with random
					// numbers.
					ab := make([]float64, n*ldab)
					for i := range ab {
						ab[i] = rnd.Float64()
					}
					// Make sure that the matrix is diagonally dominant, this guarantees
					// positive definiteness.
					switch uplo {
					case blas.Upper:
						for i := 0; i < n; i++ {
							ab[i*ldab] = float64(2*kd) + rnd.Float64()
						}
					case blas.Lower:
						for i := 0; i < n; i++ {
							ab[i*ldab+kd] = float64(2*kd) + rnd.Float64()
						}
					}

					abFac := make([]float64, len(ab))
					copy(abFac, ab)

					// Compute the Cholesky decomposition of the symmetric band matrix A.
					ok := impl.Dpbtrf(uplo, n, kd, abFac, ldab)
					if !ok {
						t.Fatalf("%v: Dpbtrf failed", name)
					}

					if n == 0 {
						continue
					}

					bi := blas64.Implementation()
					switch uplo {
					case blas.Upper:
						// Compute the product U^T * U.
						for k := n - 1; k >= 0; k-- {
							kc := min(k, kd)
							// Compute the diagonal [k,k] element.
							abFac[k*ldab] = bi.Ddot(kc+1, abFac[(k-kc)*ldab+kc:], ldab-1, abFac[(k-kc)*ldab+kc:], ldab-1)
							// Compute the rest of column k.
							if kc > 0 {
								bi.Dtrmv(blas.Upper, blas.Trans, blas.NonUnit, kc,
									abFac[(k-kc)*ldab:], ldab-1, abFac[(k-kc)*ldab+kc:], ldab-1)
							}
							//              0 1 2 3 4   n=5 kd=2
							// a - - - - )( a|a|a|0|0 0  1
							// a a - - - )( - a|a|a|0 1  2
							// a a t - - )( - - a|a|a 2  3   kc=1
							// 0 a t t - )( - - - a|a 3  4   klen=2
							// 0 0 a a a )( - - - - a 4  5
							//              1 2 3 4 5
						}
					case blas.Lower:
						// Compute the product L * L^T.
					}

					// Compute and check the max-norm distance between got and A.
					var diff float64
					switch uplo {
					case blas.Upper:
						for i := 0; i < n; i++ {
							for j := 0; j < min(kd+1, n-i); j++ {
								diff = math.Max(diff, math.Abs(abFac[i*ldab+j]-ab[i*ldab+j]))
							}
						}
					case blas.Lower:
						for i := 0; i < n; i++ {
							for j := max(0, i-kd); j <= i; j++ {
								// diff = math.Max(diff, math.Abs(got[i*n+j]-abCopy[i*ldab+kd-i+j]))
							}
						}
					}
					if diff > tol {
						t.Errorf("%v: unexpected result, diff=%v", name, diff)
					}
				}
			}
		}
	}
}
