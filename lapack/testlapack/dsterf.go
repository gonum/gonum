// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/lapack"
)

type Dsterfer interface {
	Dsteqrer
	Dlansyer
	Dsterf(n int, d, e []float64) (ok bool)
}

func DsterfTest(t *testing.T, impl Dsterfer) {
	const tol = 1e-14

	// Tests with precomputed eigenvalues.
	for cas, test := range []struct {
		d []float64
		e []float64
		n int

		want []float64
	}{
		{
			d: []float64{1, 3, 4, 6},
			e: []float64{2, 4, 5},
			n: 4,
			// Computed from original Fortran code.
			want: []float64{11.046227528488854, 4.795922173417400, -2.546379458290125, 0.704229756383872},
		},
	} {
		n := test.n
		got := make([]float64, len(test.d))
		copy(got, test.d)
		e := make([]float64, len(test.e))
		copy(e, test.e)
		ok := impl.Dsterf(n, got, e)
		if !ok {
			t.Errorf("Case %d, n=%v: Dsterf failed", cas, n)
			continue
		}
		want := make([]float64, len(test.want))
		copy(want, test.want)
		sort.Float64s(want)
		diff := floats.Distance(got, want, math.Inf(1))
		if diff > tol {
			t.Errorf("Case %d, n=%v: unexpected result, |dGot-dWant|=%v", cas, n, diff)
		}
	}

	rnd := rand.New(rand.NewSource(1))
	// Probabilistic tests.
	for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 10, 50} {
		for typ := 0; typ <= 8; typ++ {
			d := make([]float64, n)
			var e []float64
			if n > 1 {
				e = make([]float64, n-1)
			}
			// Generate a tridiagonal matrix A.
			switch typ {
			case 0:
				// The zero matrix.
			case 1:
				// The identity matrix.
				for i := range d {
					d[i] = 1
				}
			case 2:
				// A diagonal matrix with evenly spaced entries
				// 1, ..., eps  and random signs.
				for i := 0; i < n; i++ {
					if i == 0 {
						d[i] = 1
					} else {
						d[i] = 1 - (1-dlamchE)*float64(i)/float64(n-1)
					}
					if rnd.Float64() < 0.5 {
						d[i] *= -1
					}
				}
			case 3, 4, 5:
				// A diagonal matrix with geometrically spaced entries
				// 1, ..., eps  and random signs.
				for i := 0; i < n; i++ {
					if i == 0 {
						d[i] = 1
					} else {
						d[i] = math.Pow(dlamchE, float64(i)/float64(n-1))
					}
					if rnd.Float64() < 0.5 {
						d[i] *= -1
					}
				}
				switch typ {
				case 4:
					// Multiply by SQRT(overflow threshold).
					floats.Scale(math.Sqrt(1/dlamchS), d)
				case 5:
					// Multiply by SQRT(underflow threshold).
					floats.Scale(math.Sqrt(dlamchS), d)
				}
			case 6:
				// A diagonal matrix with "clustered" entries 1, eps, ..., eps
				// and random signs.
				for i := range d {
					if i == 0 {
						d[i] = 1
					} else {
						d[i] = dlamchE
					}
				}
				for i := range d {
					if rnd.Float64() < 0.5 {
						d[i] *= -1
					}
				}
			case 7:
				// Diagonal matrix with random entries.
				for i := range d {
					d[i] = rnd.NormFloat64()
				}
			case 8:
				// Random symmetric tridiagonal matrix.
				for i := range d {
					d[i] = rnd.NormFloat64()
				}
				for i := range e {
					e[i] = rnd.NormFloat64()
				}
			}
			eCopy := make([]float64, len(e))
			copy(eCopy, e)

			name := fmt.Sprintf("n=%d,type=%d", n, typ)

			// Compute the eigenvalues of A using Dsterf.
			dGot := make([]float64, len(d))
			copy(dGot, d)
			ok := impl.Dsterf(n, dGot, e)
			if !ok {
				t.Errorf("%v: Dsterf failed", name)
				continue
			}

			if n == 0 {
				continue
			}

			// Test that the eigenvalues are sorted.
			if !sort.Float64sAreSorted(dGot) {
				t.Errorf("%v: eigenvalues are not sorted", name)
				continue
			}

			// Compute the expected eigenvalues of A using Dsteqr.
			dWant := make([]float64, len(d))
			copy(dWant, d)
			copy(e, eCopy)
			z := nanGeneral(n, n, n)
			ok = impl.Dsteqr(lapack.EVTridiag, n, dWant, e, z.Data, z.Stride, make([]float64, 2*n))
			if !ok {
				t.Errorf("%v: computing reference solution using Dsteqr failed", name)
				continue
			}

			if resid := residualOrthogonal(z, false); resid > tol*float64(n) {
				t.Errorf("%v: Z is not orthogonal; resid=%v, want<=%v", name, resid, tol*float64(n))
			}

			// Check whether eigenvalues from Dsteqr and Dsterf (which use
			// different algorithms) are equal.
			var diff, dMax float64
			for i, di := range dGot {
				diffAbs := math.Abs(di - dWant[i])
				diff = math.Max(diff, diffAbs)
				dAbs := math.Max(math.Abs(di), math.Abs(dWant[i]))
				dMax = math.Max(dMax, dAbs)
			}
			dMax = math.Max(dlamchS, dMax)
			if diff > tol*dMax {
				t.Errorf("%v: unexpected result; |dGot-dWant|=%v", name, diff)
			}

			// Construct A as a symmetric dense matrix and compute its 1-norm.
			copy(e, eCopy)
			lda := n
			a := make([]float64, n*lda)
			var anorm, tmp float64
			for i := 0; i < n-1; i++ {
				a[i*lda+i] = d[i]
				a[i*lda+i+1] = e[i]
				tmp2 := math.Abs(e[i])
				anorm = math.Max(anorm, math.Abs(d[i])+tmp+tmp2)
				tmp = tmp2
			}
			a[(n-1)*lda+n-1] = d[n-1]
			anorm = math.Max(anorm, math.Abs(d[n-1])+tmp)

			// Compute A - Z D Zᵀ. The result should be the zero matrix.
			bi := blas64.Implementation()
			for i := 0; i < n; i++ {
				bi.Dsyr(blas.Upper, n, -dGot[i], z.Data[i:], z.Stride, a, lda)
			}

			// Compute |A - Z D Zᵀ|.
			wnorm := impl.Dlansy(lapack.MaxColumnSum, blas.Upper, n, a, lda, make([]float64, n))

			// Compute diff := |A - Z D Zᵀ| / (|A| N).
			if anorm > wnorm {
				diff = wnorm / anorm / float64(n)
			} else {
				if anorm < 1 {
					diff = math.Min(wnorm, float64(n)*anorm) / anorm / float64(n)
				} else {
					diff = math.Min(wnorm/anorm, float64(n)) / float64(n)
				}
			}

			// Check whether diff is small.
			if diff > tol {
				t.Errorf("%v: unexpected result; |A - Z D Zᵀ|/(|A| n)=%v", name, diff)
			}
		}
	}
}
