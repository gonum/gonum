// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

type Dptsver interface {
	Dptsv(n, nrhs int, d, e []float64, b []float64, ldb int) (ok bool)
}

func DptsvTest(t *testing.T, impl Dptsver) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50, 51, 52, 53, 54, 100} {
		for _, nrhs := range []int{0, 1, 2, 3, 4, 5, 10, 20, 50} {
			for _, ldb := range []int{max(1, nrhs), nrhs + 3} {
				dptsvTest(t, impl, rnd, n, nrhs, ldb)
			}
		}
	}
}

func dptsvTest(t *testing.T, impl Dptsver, rnd *rand.Rand, n, nrhs, ldb int) {
	const tol = 1e-15

	name := fmt.Sprintf("n=%v", n)

	// Generate a random diagonally dominant symmetric tridiagonal matrix A.
	d, e := newRandomSymTridiag(n, rnd)

	// Generate a random solution matrix X.
	xWant := randomGeneral(n, nrhs, ldb, rnd)

	// Compute the right-hand side.
	b := zeros(n, nrhs, ldb)
	dstmm(n, nrhs, d, e, xWant.Data, xWant.Stride, b.Data, b.Stride)

	// Solve A*X=B.
	ok := impl.Dptsv(n, nrhs, d, e, b.Data, b.Stride)
	if !ok {
		t.Errorf("%v: Dptsv failed", name)
		return
	}

	resid := dpttrsResidual(b, xWant)
	if resid > tol {
		t.Errorf("%v: unexpected solution: |diff| = %v, want <= %v", name, resid, tol)
	}
}
