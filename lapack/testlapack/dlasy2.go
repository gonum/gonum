// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dlasy2er interface {
	Dlasy2(tranl, tranr bool, isgn, n1, n2 int, tl []float64, ldtl int, tr []float64, ldtr int, b []float64, ldb int, x []float64, ldx int) (scale, xnorm float64, ok bool)
}

func Dlasy2Test(t *testing.T, impl Dlasy2er) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, tranl := range []bool{true, false} {
		for _, tranr := range []bool{true, false} {
			for _, isgn := range []int{1, -1} {
				for _, n1 := range []int{0, 1, 2} {
					for _, n2 := range []int{0, 1, 2} {
						for _, extra := range []int{0, 3} {
							for cas := 0; cas < 100; cas++ {
								var big bool
								if cas%2 == 0 {
									big = true
								}
								testDlasy2(t, impl, tranl, tranr, isgn, n1, n2, extra, big, rnd)
							}
						}
					}
				}
			}
		}
	}
}

func testDlasy2(t *testing.T, impl Dlasy2er, tranl, tranr bool, isgn, n1, n2, extra int, big bool, rnd *rand.Rand) {
	const tol = 1e-14

	name := fmt.Sprintf("Case n1=%v, n2=%v, isgn=%v, big=%v", n1, n2, isgn, big)

	tl := randomGeneral(n1, n1, n1+extra, rnd)
	tr := randomGeneral(n2, n2, n2+extra, rnd)
	x := randomGeneral(n1, n2, n2+extra, rnd)
	b := randomGeneral(n1, n2, n2+extra, rnd)
	if big {
		for i := 0; i < n1; i++ {
			for j := 0; j < n2; j++ {
				b.Data[i*b.Stride+j] *= bignum
			}
		}
	}

	tlCopy := cloneGeneral(tl)
	trCopy := cloneGeneral(tr)
	bCopy := cloneGeneral(b)

	scale, xnorm, ok := impl.Dlasy2(tranl, tranr, isgn, n1, n2, tl.Data, tl.Stride, tr.Data, tr.Stride, b.Data, b.Stride, x.Data, x.Stride)

	// Check any invalid modifications in read-only input.
	if !equalGeneral(tl, tlCopy) {
		t.Errorf("%v: unexpected modification in TL", name)
	}
	if !equalGeneral(tr, trCopy) {
		t.Errorf("%v: unexpected modification in TR", name)
	}
	if !equalGeneral(b, bCopy) {
		t.Errorf("%v: unexpected modification in B", name)
	}

	// Check any invalid modifications of x.
	if !generalOutsideAllNaN(x) {
		t.Errorf("%v: out-of-range write to x\n%v", name, x.Data)
	}

	if n1 == 0 || n2 == 0 {
		return
	}

	if scale <= 0 || 1 < scale {
		t.Errorf("%v: invalid value of scale, want in (0,1], got %v", name, scale)
	}

	xnormWant := dlange(lapack.MaxRowSum, x.Rows, x.Cols, x.Data, x.Stride)
	if xnormWant != xnorm {
		t.Errorf("%v: unexpected xnorm: want %v, got %v", name, xnormWant, xnorm)
	}

	if !ok {
		t.Logf("%v: Dlasy2 returned ok=false", name)
		return
	}

	// Compute diff := op(TL)*X + sgn*X*op(TR) - scale*B.
	diff := zeros(n1, n2, n2)
	if tranl {
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, tl, x, 0, diff)
	} else {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tl, x, 0, diff)
	}
	if tranr {
		blas64.Gemm(blas.NoTrans, blas.Trans, float64(isgn), x, tr, 1, diff)
	} else {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, float64(isgn), x, tr, 1, diff)
	}
	for i := 0; i < n1; i++ {
		for j := 0; j < n2; j++ {
			diff.Data[i*diff.Stride+j] -= scale * b.Data[i*b.Stride+j]
		}
	}
	// Check that residual |op(TL)*X + sgn*X*op(TR) - scale*B| / |X| is small.
	resid := dlange(lapack.MaxColumnSum, n1, n2, diff.Data, diff.Stride) / xnorm
	if resid > tol {
		t.Errorf("%v: unexpected result, resid=%v, want<=%v", name, resid, tol)
	}
}
