// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"testing"
)

type Dlaln2er interface {
	Dlaln2(trans bool, na, nw int, smin, ca float64, a []float64, lda int, d1, d2 float64, b []float64, ldb int, wr, wi float64, x []float64, ldx int) (scale, xnorm float64, ok bool)
}

func Dlaln2Test(t *testing.T, impl Dlaln2er) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, trans := range []bool{true, false} {
		for _, na := range []int{1, 2} {
			for _, nw := range []int{1, 2} {
				for _, extra := range []int{0, 1, 2, 13} {
					for cas := 0; cas < 1000; cas++ {
						testDlaln2(t, impl, trans, na, nw, extra, rnd)
					}
				}
			}
		}
	}
}

func testDlaln2(t *testing.T, impl Dlaln2er, trans bool, na, nw, extra int, rnd *rand.Rand) {
	const tol = 1e-11

	// Generate random input scalars.
	ca := rnd.NormFloat64()
	d1 := rnd.NormFloat64()
	d2 := rnd.NormFloat64()

	var w complex128
	if nw == 1 {
		w = complex(rnd.NormFloat64(), 0)
	} else {
		w = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}
	smin := dlamchP * (math.Abs(real(w)) + math.Abs(imag(w)))

	// Generate random input matrices.
	a := randomGeneral(na, na, na+extra, rnd)
	b := randomGeneral(na, nw, nw+extra, rnd)
	x := randomGeneral(na, nw, nw+extra, rnd)

	scale, xnormGot, ok := impl.Dlaln2(trans, na, nw, smin, ca, a.Data, a.Stride, d1, d2, b.Data, b.Stride, real(w), imag(w), x.Data, x.Stride)

	prefix := fmt.Sprintf("Case trans=%t, na=%v, nw=%v, extra=%v", trans, na, nw, extra)

	if !generalOutsideAllNaN(a) {
		t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
	}
	if !generalOutsideAllNaN(b) {
		t.Errorf("%v: out-of-range write to B\n%v", prefix, b.Data)
	}
	if !generalOutsideAllNaN(x) {
		t.Errorf("%v: out-of-range write to X\n%v", prefix, x.Data)
	}

	// Scale is documented to be <= 1.
	if scale <= 0 || 1 < scale {
		t.Errorf("%v: invalid value of scale=%v", prefix, scale)
	}

	// Calculate the infinity norm of X explicitly.
	var xnormWant float64
	for i := 0; i < na; i++ {
		var rowsum float64
		for j := 0; j < nw; j++ {
			rowsum += math.Abs(x.Data[i*x.Stride+j])
		}
		if rowsum > xnormWant {
			xnormWant = rowsum
		}
	}
	if xnormWant != xnormGot {
		t.Errorf("Case %v: unexpected xnorm with scale=%v. Want %v, got %v", prefix, scale, xnormWant, xnormGot)
	}

	if !ok {
		// If ok is false, the matrix has been perturbed but we don't
		// know how. Return without comparing both sides of the
		// equation.
		return
	}

	// Compute a complex matrix
	//  M := ca * A - w * D
	// or
	//  M := ca * Aᵀ - w * D.
	m := make([]complex128, na*na)
	if trans {
		// M = ca * Aᵀ
		for i := 0; i < na; i++ {
			for j := 0; j < na; j++ {
				m[i*na+j] = complex(ca*a.Data[j*a.Stride+i], 0)
			}
		}
	} else {
		// M = ca * Aᵀ
		for i := 0; i < na; i++ {
			for j := 0; j < na; j++ {
				m[i*na+j] = complex(ca*a.Data[i*a.Stride+j], 0)
			}
		}
	}
	// Subtract the diagonal matrix w * D.
	m[0] -= w * complex(d1, 0)
	if na == 2 {
		m[3] -= w * complex(d2, 0)
	}

	// Convert real na×2 matrices X and scale*B into complex na-vectors.
	cx := make([]complex128, na)
	cb := make([]complex128, na)
	switch nw {
	case 1:
		for i := 0; i < na; i++ {
			cx[i] = complex(x.Data[i*x.Stride], 0)
			cb[i] = complex(scale*b.Data[i*x.Stride], 0)
		}
	case 2:
		for i := 0; i < na; i++ {
			cx[i] = complex(x.Data[i*x.Stride], x.Data[i*x.Stride+1])
			cb[i] = complex(scale*b.Data[i*b.Stride], scale*b.Data[i*b.Stride+1])
		}
	}

	// Compute M * X.
	mx := make([]complex128, na)
	for i := 0; i < na; i++ {
		for j := 0; j < na; j++ {
			mx[i] += m[i*na+j] * cx[j]
		}
	}
	// Check whether |M * X - scale * B|_max <= tol.
	for i := 0; i < na; i++ {
		if cmplx.Abs(mx[i]-cb[i]) > tol {
			t.Errorf("Case %v: unexpected value of left-hand side at row %v with scale=%v. Want %v, got %v", prefix, i, scale, cb[i], mx[i])
		}
	}
}
