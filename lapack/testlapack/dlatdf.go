// Copyright ©2021 The Gonum Authors. All rights reserved.
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
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlatdfer interface {
	Dlatdf(ijob, n int, z []float64, ldz int, rhs []float64, rdsum, rdscal float64, ipiv, jpiv []int) (sum, scale float64)

	Dgetc2er
	Dgesc2er
}

func DlatdfTest(t *testing.T, impl Dlatdfer) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{3} {
		for _, ldz := range []int{n, n + 5} {
			testDlatdf(t, impl, rnd, 0, n, ldz, tol)
			testDlatdf(t, impl, rnd, 2, n, ldz, tol)
			break
		}
	}
}

func testDlatdf(t *testing.T, impl Dlatdfer, rnd *rand.Rand, ijob, n int, ldz int, tol float64) {
	name := fmt.Sprintf("n=%v,ldz=%v", n, ldz)
	// Random general matrix for Z
	z := randomGeneral(n, n, max(1, ldz), rnd)
	// Generate a random right hand side vector.
	b := randomGeneral(n, 1, 1, rnd)

	// Compute the LU part of the factorization of the n×n
	// matrix Z with Dgetc2:  Z = P * L * U * Q
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	lu := cloneGeneral(z)
	k := impl.Dgetc2(n, lu.Data, lu.Stride, ipiv, jpiv)
	if k > 0 {
		t.Errorf("%v: %d index of U was perturbed in Dgetc2", name, k)
	}
	rhs := cloneGeneral(b)
	ipivCopy := make([]int, len(ipiv))
	copy(ipivCopy, ipiv)
	jpivCopy := make([]int, len(jpiv))
	copy(jpivCopy, jpiv)
	scal := 1.
	if ijob == 0 {
		scal = impl.Dgesc2(n, z.Data, z.Stride, rhs.Data, ipiv, jpiv)
	}
	// From reference: rdscal (and rdsum) only makes
	// sense when Dtgsy2 is called by Dtgsyl.
	rdsum := 0.
	rdscal := scal
	// Call Dlatdf to solve Z*x = scale*b.
	x := cloneGeneral(rhs)
	luCopy := cloneGeneral(lu)
	sum, _ := impl.Dlatdf(ijob, n, lu.Data, lu.Stride, x.Data, rdsum, rdscal, ipiv, jpiv)
	if n == 0 {
		return
	}
	_, _ = sum, scal // are these used?
	// Check that const input to Dlatdf hasn't been modified.
	if !floats.Same(lu.Data, luCopy.Data) {
		t.Errorf("%v: unexpected modification in LU decompositon of Z", name)
	}
	if !intsEqual(ipiv, ipivCopy) {
		t.Errorf("%v: unexpected modification in ipiv", name)
	}
	if !intsEqual(jpiv, jpivCopy) {
		t.Errorf("%v: unexpected modification in jpiv", name)
	}

	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, z, x, -scal, b)
	diff := b

	// Compute the residual |Z*x - b| / |x|.
	xnorm := dlange(lapack.MaxColumnSum, n, 1, x.Data, 1)
	resid := dlange(lapack.MaxColumnSum, n, 1, diff.Data, 1) / xnorm
	if resid > tol || math.IsNaN(resid) {
		t.Errorf("%v: unexpected result;scal=%v, resid=%v, want<=%v", name, scal, resid, tol)
	}

}
