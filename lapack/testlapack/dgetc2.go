// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"
)

type Dgetc2er interface {
	Dgetc2(n int, a []float64, lda int, ipiv, jpiv []int) (k int)
}

func Dgetc2Test(t *testing.T, impl Dgetc2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, lda := range []int{n} {
			dgetc2Test(t, impl, rnd, n, lda)
		}
	}
}

func dgetc2Test(t *testing.T, impl Dgetc2er, rnd *rand.Rand, n, lda int) {
	name := fmt.Sprintf("n=%v,lda=%v", n, lda)
	if lda == 0 {
		lda = 1
	}
	// Generate a random symmetric positive definite band matrix.
	// apd := randSymBand(uplo, n, kd, ldab, rnd)
	ap := randomGeneral(n, n, lda, rnd)
	// Copy to store output
	aout := make([]float64, len(ap.Data))
	// ipib and jpiv are outputs.
	ipiv := make([]int, n)
	jpiv := make([]int, n)
	copy(aout, ap.Data)
	k := impl.Dgetc2(n, aout, lda, ipiv, jpiv)
	if k >= 0 {
		t.Fatalf("%v: matrix was perturbed at %d", name, k)
	}

	// Check uniqueness of pivot indices.
	// This is done using an int as storage for
	// an int tuple
	pivmap := make(map[int]struct{})

	errcount := 0
	for i := range ipiv {
		tuple := ipiv[i]<<16 | jpiv[i]
		if _, present := pivmap[tuple]; present {
			t.Errorf("%v: ipiv repeated value found (%d,%d)", name, ipiv[i], jpiv[i])
			errcount++
		}
		pivmap[tuple] = struct{}{}
	}
	if errcount > 0 {
		t.Errorf("ipiv:%d", ipiv)
		t.Errorf("jpiv:%d", jpiv)
	}
}
