// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
)

type Dpbtf2er interface {
	Dpbtf2(uplo blas.Uplo, n, kd int, ab []float64, ldab int) (ok bool)
}

// Dpbtf2Test tests Dpbtf2 on random symmetric positive definite band matrices
// by checking that the Cholesky factors multiply back to the original matrix.
func Dpbtf2Test(t *testing.T, impl Dpbtf2er) {
	// TODO(vladimir-ch): include expected-failure test case.
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, kd := range []int{0, (n + 1) / 4, (3*n - 1) / 4, (5*n + 1) / 4} {
			for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
				for _, ldab := range []int{kd + 1, kd + 1 + 7} {
					dpbtf2Test(t, impl, rnd, uplo, n, kd, ldab)
				}
			}
		}
	}
}

func dpbtf2Test(t *testing.T, impl Dpbtf2er, rnd *rand.Rand, uplo blas.Uplo, n, kd int, ldab int) {
	const tol = 1e-12

	name := fmt.Sprintf("uplo=%v,n=%v,kd=%v,ldab=%v", string(uplo), n, kd, ldab)

	// Allocate a band matrix and fill it with random numbers.
	ab := make([]float64, n*ldab)
	for i := range ab {
		ab[i] = rnd.NormFloat64()
	}
	// Make sure that the matrix U or L has a sufficiently positive diagonal.
	switch uplo {
	case blas.Upper:
		for i := 0; i < n; i++ {
			ab[i*ldab] = 2 + rnd.Float64()
		}
	case blas.Lower:
		for i := 0; i < n; i++ {
			ab[i*ldab+kd] = 2 + rnd.Float64()
		}
	}
	// Compute U^T*U or L*L^T. The resulting (symmetric) matrix A will be positive definite.
	dsbmm(uplo, n, kd, ab, ldab)

	// Compute the Cholesky decomposition of A.
	abFac := make([]float64, len(ab))
	copy(abFac, ab)
	ok := impl.Dpbtf2(uplo, n, kd, abFac, ldab)
	if !ok {
		t.Fatalf("%v: bad test matrix, Dpbtf2 failed", name)
	}

	if n == 0 {
		return
	}

	// Reconstruct an symmetric band matrix from the U^T*U or L*L^T factorization, overwriting abFac.
	dsbmm(uplo, n, kd, abFac, ldab)

	// Compute and check the max-norm distance between the reconstructed and original matrix A.
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
			for j := max(0, kd-i); j < kd+1; j++ {
				diff = math.Max(diff, math.Abs(abFac[i*ldab+j]-ab[i*ldab+j]))
			}
		}
	}
	if diff > tol {
		t.Errorf("%v: unexpected result, diff=%v", name, diff)
	}
}
