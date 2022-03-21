// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dlapmrer interface {
	Dlapmr(forwrd bool, m, n int, x []float64, ldx int, k []int)
}

func DlapmrTest(t *testing.T, impl Dlapmrer) {
	rnd := rand.New(rand.NewSource(1))
	for _, fwd := range []bool{true, false} {
		for _, m := range []int{0, 1, 2, 3, 4, 5, 10} {
			for _, n := range []int{0, 1, 4} {
				for _, ldx := range []int{max(1, n), n + 3} {
					dlapmrTest(t, impl, rnd, fwd, m, n, ldx)
				}
			}
		}
	}
}

func dlapmrTest(t *testing.T, impl Dlapmrer, rnd *rand.Rand, fwd bool, m, n, ldx int) {
	name := fmt.Sprintf("forwrd=%v,m=%d,n=%d,ldx=%d", fwd, m, n, ldx)

	bi := blas64.Implementation()

	// Generate a random permutation and simultaneously apply it to the rows of the identity matrix.
	k := make([]int, m)
	for i := range k {
		k[i] = i
	}
	p := eye(m, m)
	for i := 0; i < m-1; i++ {
		j := i + rnd.Intn(m-i)
		k[i], k[j] = k[j], k[i]
		bi.Dswap(m, p.Data[i*p.Stride:], 1, p.Data[j*p.Stride:], 1)
	}
	kCopy := make([]int, len(k))
	copy(kCopy, k)

	// Generate a random matrix X.
	x := randomGeneral(m, n, ldx, rnd)

	// Applying the permutation k with Dlapmr is the same as multiplying X with P or Pᵀ from the left:
	//  - forward permutation:  P * X
	//  - backward permutation: Pᵀ* X
	trans := blas.NoTrans
	if !fwd {
		trans = blas.Trans
	}
	want := zeros(m, n, n)
	bi.Dgemm(trans, blas.NoTrans, m, n, m, 1, p.Data, p.Stride, x.Data, x.Stride, 0, want.Data, want.Stride)

	// Apply the permutation in k to X.
	impl.Dlapmr(fwd, m, n, x.Data, x.Stride, k)
	got := x

	// Check that k hasn't been modified in Dlapmr.
	var kmod bool
	for i, ki := range k {
		if ki != kCopy[i] {
			kmod = true
			break
		}
	}
	if kmod {
		t.Errorf("%s: unexpected modification of k", name)
	}

	// Check that Dlapmr yields the same result as multiplication with P.
	if !equalGeneral(got, want) {
		t.Errorf("%s: unexpected result", name)
	}
}
