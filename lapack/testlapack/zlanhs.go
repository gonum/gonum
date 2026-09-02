// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/lapack"
)

type Zlanhser interface {
	Zlanhs(norm lapack.MatrixNorm, n int, a []complex128, lda int, work []float64) float64
}

func ZlanhsTest(t *testing.T, impl Zlanhser) {
	const tol = 1e-14
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{0, 1, 2, 4, 9} {
		for _, lda := range []int{max(1, n), n + 5} {
			// Random n×n matrix (we zero the sub-subdiagonal in the reference,
			// matching how Dlanhs_test.go operates).
			a := make([]complex128, n*lda)
			for i := range a {
				a[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
			}
			for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxRowSum, lapack.MaxColumnSum, lapack.Frobenius} {
				var work []float64
				if norm == lapack.MaxColumnSum {
					work = make([]float64, n)
				}
				got := impl.Zlanhs(norm, n, a, lda, work)

				// Reference: zero the elements below the first subdiagonal and
				// compute the full-matrix norm on the result.
				aRef := make([]complex128, len(a))
				copy(aRef, a)
				for i := 2; i < n; i++ {
					for j := 0; j < i-1; j++ {
						aRef[i*lda+j] = 0
					}
				}
				want := zlangeRef(norm, n, n, aRef, lda)
				if want == 0 {
					if math.Abs(got-want) > tol {
						t.Errorf("n=%v lda=%v norm=%v: got %v want %v", n, lda, normToString(norm), got, want)
					}
					continue
				}
				if math.Abs(got-want) > tol*want {
					t.Errorf("n=%v lda=%v norm=%v: got %v want %v", n, lda, normToString(norm), got, want)
				}
			}
		}
	}
}

// zlangeRef is a naive reference implementation of Zlange used for cross
// checking. It is intentionally unoptimised and mirrors cmplx.Abs-based
// accumulation.
func zlangeRef(norm lapack.MatrixNorm, m, n int, a []complex128, lda int) float64 {
	if m == 0 || n == 0 {
		return 0
	}
	var value float64
	switch norm {
	case lapack.MaxAbs:
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				value = math.Max(value, cmplx.Abs(a[i*lda+j]))
			}
		}
	case lapack.MaxColumnSum:
		for j := 0; j < n; j++ {
			var sum float64
			for i := 0; i < m; i++ {
				sum += cmplx.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
	case lapack.MaxRowSum:
		for i := 0; i < m; i++ {
			var sum float64
			for j := 0; j < n; j++ {
				sum += cmplx.Abs(a[i*lda+j])
			}
			value = math.Max(value, sum)
		}
	case lapack.Frobenius:
		var sum float64
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				abs := cmplx.Abs(a[i*lda+j])
				sum += abs * abs
			}
		}
		value = math.Sqrt(sum)
	default:
		panic("invalid norm")
	}
	return value
}
