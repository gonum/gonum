// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dlantber interface {
	Dlantb(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, n, k int, a []float64, lda int, work []float64) float64
}

func DlantbTest(t *testing.T, impl Dlantber) {
	rnd := rand.New(rand.NewSource(1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
		for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
			for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
				name := normToString(norm) + uploToString(uplo) + diagToString(diag)
				t.Run(name, func(t *testing.T) {
					for _, n := range []int{0, 1, 2, 3, 4, 5, 10} {
						for _, k := range []int{0, 1, 2, 3, n, n + 2} {
							for _, lda := range []int{k + 1, k + 3} {
								for iter := 0; iter < 10; iter++ {
									dlantbTest(t, impl, rnd, norm, uplo, diag, n, k, lda)
								}
							}
						}
					}
				})
			}
		}
	}
}

func dlantbTest(t *testing.T, impl Dlantber, rnd *rand.Rand, norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, n, k, lda int) {
	const tol = 1e-14

	name := fmt.Sprintf("n=%v,k=%v,lda=%v", n, k, lda)

	// Deal with zero-sized matrices early.
	if n == 0 {
		got := impl.Dlantb(norm, uplo, diag, n, k, nil, lda, nil)
		if got != 0 {
			t.Errorf("%v: unexpected result for zero-sized matrix", name)
		}
		return
	}

	a := make([]float64, max(0, (n-1)*lda+k+1))
	if rnd.Float64() < 0.5 {
		// Sometimes fill A with elements between -0.5 and 0.5 so that for
		// blas.Unit matrices the largest element is the 1 on the main diagonal.
		for i := range a {
			// Between -0.5 and 0.5.
			a[i] = rnd.Float64() - 0.5
		}
	} else {
		for i := range a {
			// Between -2 and 2.
			a[i] = 4*rnd.Float64() - 2
		}
	}
	// Sometimes put a NaN into A.
	if rnd.Float64() < 0.5 {
		a[rnd.Intn(len(a))] = math.NaN()
	}
	// Make a copy of A for later comparison.
	aCopy := make([]float64, len(a))
	copy(aCopy, a)

	var work []float64
	if norm == lapack.MaxColumnSum {
		work = make([]float64, n)
	}
	// Fill work with random garbage.
	for i := range work {
		work[i] = rnd.NormFloat64()
	}

	got := impl.Dlantb(norm, uplo, diag, n, k, a, lda, work)

	if !floats.Same(a, aCopy) {
		t.Fatalf("%v: unexpected modification of a", name)
	}

	// Generate a dense representation of A and compute the wanted result.
	ldaGen := n
	aGen := make([]float64, n*ldaGen)
	if uplo == blas.Upper {
		for i := 0; i < n; i++ {
			for j := 0; j < min(n-i, k+1); j++ {
				aGen[i*ldaGen+i+j] = a[i*lda+j]
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := max(0, k-i); j < k+1; j++ {
				aGen[i*ldaGen+i-(k-j)] = a[i*lda+j]
			}
		}
	}
	if diag == blas.Unit {
		for i := 0; i < n; i++ {
			aGen[i*ldaGen+i] = 1
		}
	}
	want := dlange(norm, n, n, aGen, ldaGen)

	if math.IsNaN(want) {
		if !math.IsNaN(got) {
			t.Errorf("%v: unexpected result with NaN element; got %v, want %v\n%v\n%v", name, got, want, a, aGen)
		}
		return
	}

	if norm == lapack.MaxAbs {
		if got != want {
			t.Errorf("%v: unexpected result; got %v, want %v", name, got, want)
		}
		return
	}
	diff := math.Abs(got - want)
	if diff > tol {
		t.Errorf("%v: unexpected result; got %v, want %v, diff=%v", name, got, want, diff)
	}
}
