// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlarfer interface {
	Dlarf(side blas.Side, m, n int, v []float64, incv int, tau float64, c []float64, ldc int, work []float64)
}

func DlarfTest(t *testing.T, impl Dlarfer) {
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		name := "Right"
		if side == blas.Left {
			name = "Left"
		}
		t.Run(name, func(t *testing.T) {
			runDlarfTest(t, impl, side)
		})
	}
}

func runDlarfTest(t *testing.T, impl Dlarfer, side blas.Side) {
	rnd := rand.New(rand.NewSource(1))
	for _, m := range []int{0, 1, 2, 3, 4, 5, 10} {
		for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
			for _, incv := range []int{1, 4} {
				for _, ldc := range []int{max(1, n), n + 3} {
					for _, nnzv := range []int{0, 1, 2} {
						for _, nnzc := range []int{0, 1, 2} {
							for _, tau := range []float64{0, rnd.NormFloat64()} {
								dlarfTest(t, impl, rnd, side, m, n, incv, ldc, nnzv, nnzc, tau)
							}
						}
					}
				}
			}
		}
	}
}

func dlarfTest(t *testing.T, impl Dlarfer, rnd *rand.Rand, side blas.Side, m, n, incv, ldc, nnzv, nnzc int, tau float64) {
	const tol = 1e-14

	c := make([]float64, m*ldc)
	for i := range c {
		c[i] = rnd.NormFloat64()
	}
	switch nnzc {
	case 0:
		// Zero out all of C.
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				c[i*ldc+j] = 0
			}
		}
	case 1:
		// Zero out right or bottom half of C.
		if side == blas.Left {
			for i := 0; i < m; i++ {
				for j := n / 2; j < n; j++ {
					c[i*ldc+j] = 0
				}
			}
		} else {
			for i := m / 2; i < m; i++ {
				for j := 0; j < n; j++ {
					c[i*ldc+j] = 0
				}
			}
		}
	default:
		// Leave C with random content.
	}
	cCopy := make([]float64, len(c))
	copy(cCopy, c)

	var work []float64
	if side == blas.Left {
		work = make([]float64, n)
	} else {
		work = make([]float64, m)
	}

	vlen := n
	if side == blas.Left {
		vlen = m
	}
	vlen = max(1, vlen)
	v := make([]float64, 1+(vlen-1)*incv)
	for i := range v {
		v[i] = rnd.NormFloat64()
	}
	switch nnzv {
	case 0:
		// Zero out all of v.
		for i := 0; i < vlen; i++ {
			v[i*incv] = 0
		}
	case 1:
		// Zero out half of v.
		for i := vlen / 2; i < vlen; i++ {
			v[i*incv] = 0
		}
	default:
		// Leave v with random content.
	}
	vCopy := make([]float64, len(v))
	copy(vCopy, v)

	impl.Dlarf(side, m, n, v, incv, tau, c, ldc, work)
	got := c

	name := fmt.Sprintf("m=%d,n=%d,incv=%d,tau=%f,ldc=%d", m, n, incv, tau, ldc)

	if !floats.Equal(v, vCopy) {
		t.Errorf("%v: unexpected modification of v", name)
	}
	if tau == 0 && !floats.Equal(got, cCopy) {
		t.Errorf("%v: unexpected modification of C", name)
	}

	if m == 0 || n == 0 || tau == 0 {
		return
	}

	bi := blas64.Implementation()

	want := make([]float64, len(cCopy))
	if side == blas.Left {
		// Compute want = (I - tau * v * vᵀ) * C

		// vtc = -tau * vᵀ * C = -tau * Cᵀ * v
		vtc := make([]float64, n)
		bi.Dgemv(blas.Trans, m, n, -tau, cCopy, ldc, v, incv, 0, vtc, 1)

		// want = C + v * vtcᵀ
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				want[i*ldc+j] = cCopy[i*ldc+j] + v[i*incv]*vtc[j]
			}
		}
	} else {
		// Compute want = C * (I - tau * v * vᵀ)

		// cv = -tau * C * v
		cv := make([]float64, m)
		bi.Dgemv(blas.NoTrans, m, n, -tau, cCopy, ldc, v, incv, 0, cv, 1)

		// want = C + cv * vᵀ
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				want[i*ldc+j] = cCopy[i*ldc+j] + cv[i]*v[j*incv]
			}
		}
	}
	diff := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			diff[i*n+j] = got[i*ldc+j] - want[i*ldc+j]
		}
	}
	resid := dlange(lapack.MaxColumnSum, m, n, diff, n)
	if resid > tol*float64(max(m, n)) {
		t.Errorf("%v: unexpected result; resid=%v, want<=%v", name, resid, tol*float64(max(m, n)))
	}
}
