// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

type Dlasq1er interface {
	Dlasq1(n int, d, e, work []float64) int

	Dgebrd(m, n int, a []float64, lda int, d, e, tauQ, tauP, work []float64, lwork int)
}

func Dlasq1Test(t *testing.T, impl Dlasq1er) {
	const tol = 1e-14
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 8, 10, 30, 50} {
		for typ := 0; typ <= 7; typ++ {
			name := fmt.Sprintf("n=%v,typ=%v", n, typ)

			// Generate a diagonal matrix D with positive entries.
			d := make([]float64, n)
			switch typ {
			case 0:
				// The zero matrix.
			case 1:
				// The identity matrix.
				for i := range d {
					d[i] = 1
				}
			case 2:
				// A diagonal matrix with evenly spaced entries 1, ..., eps.
				for i := 0; i < n; i++ {
					if i == 0 {
						d[0] = 1
					} else {
						d[i] = 1 - (1-dlamchE)*float64(i)/float64(n-1)
					}
				}
			case 3, 4, 5:
				// A diagonal matrix with geometrically spaced entries 1, ..., eps.
				for i := 0; i < n; i++ {
					if i == 0 {
						d[0] = 1
					} else {
						d[i] = math.Pow(dlamchE, float64(i)/float64(n-1))
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
				// A diagonal matrix with "clustered" entries 1, eps, ..., eps.
				for i := range d {
					if i == 0 {
						d[i] = 1
					} else {
						d[i] = dlamchE
					}
				}
			case 7:
				// Diagonal matrix with random entries.
				for i := range d {
					d[i] = math.Abs(rnd.NormFloat64())
				}
			}

			dWant := make([]float64, n)
			copy(dWant, d)
			sort.Sort(sort.Reverse(sort.Float64Slice(dWant)))

			// Allocate work slice to the maximum length needed below.
			work := make([]float64, max(1, 4*n))

			// Generate an n×n matrix A by pre- and post-multiplying D with
			// random orthogonal matrices:
			//  A = U*D*V.
			lda := max(1, n)
			a := make([]float64, n*lda)
			Dlagge(n, n, 0, 0, d, a, lda, rnd, work)

			// Reduce A to bidiagonal form B represented by the diagonal d and
			// off-diagonal e.
			tauQ := make([]float64, n)
			tauP := make([]float64, n)
			e := make([]float64, max(0, n-1))
			impl.Dgebrd(n, n, a, lda, d, e, tauQ, tauP, work, len(work))

			// Compute the singular values of B.
			for i := range work {
				work[i] = math.NaN()
			}
			info := impl.Dlasq1(n, d, e, work)
			if info != 0 {
				t.Fatalf("%v: Dlasq1 returned non-zero info=%v", name, info)
			}

			if n == 0 {
				continue
			}

			if !sort.IsSorted(sort.Reverse(sort.Float64Slice(d))) {
				t.Errorf("%v: singular values not sorted", name)
			}

			diff := floats.Distance(d, dWant, math.Inf(1))
			if diff > tol*floats.Max(dWant) {
				t.Errorf("%v: unexpected result; diff=%v", name, diff)
			}
		}
	}
}
