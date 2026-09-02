// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

type Zlasrer interface {
	Zlasr(side blas.Side, pivot lapack.Pivot, direct lapack.Direct,
		m, n int, c, s []float64, a []complex128, lda int)
}

// ZlasrTest verifies that Zlasr applies a sequence of plane rotations matching
// a naive reference implementation built from individual rotations.
func ZlasrTest(t *testing.T, impl Zlasrer) {
	const tol = 1e-13
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		for _, pivot := range []lapack.Pivot{lapack.Variable, lapack.Top, lapack.Bottom} {
			for _, direct := range []lapack.Direct{lapack.Forward, lapack.Backward} {
				for _, sz := range []struct{ m, n int }{
					{3, 4}, {5, 5}, {2, 6}, {6, 2}, {1, 5}, {5, 1}, {7, 9},
				} {
					name := fmt.Sprintf("side=%c pivot=%d dir=%c m=%d n=%d",
						sideChar(side), pivot, directChar(direct), sz.m, sz.n)
					zlasrCheck(t, impl, rnd, side, pivot, direct, sz.m, sz.n, name, tol)
				}
			}
		}
	}
}

func zlasrCheck(t *testing.T, impl Zlasrer, rnd *rand.Rand, side blas.Side,
	pivot lapack.Pivot, direct lapack.Direct, m, n int, name string, tol float64) {
	lda := n + 2
	a := make([]complex128, m*lda)
	for i := range a {
		a[i] = complex(rnd.NormFloat64(), rnd.NormFloat64())
	}
	aCopy := make([]complex128, len(a))
	copy(aCopy, a)

	var rotCount int
	if side == blas.Left {
		rotCount = m - 1
	} else {
		rotCount = n - 1
	}
	if rotCount < 0 {
		rotCount = 0
	}
	cArr := make([]float64, rotCount)
	sArr := make([]float64, rotCount)
	for i := 0; i < rotCount; i++ {
		theta := rnd.Float64() * 2 * math.Pi
		cArr[i] = math.Cos(theta)
		sArr[i] = math.Sin(theta)
	}

	impl.Zlasr(side, pivot, direct, m, n, cArr, sArr, a, lda)

	// Reference: apply each rotation individually in the prescribed order.
	want := make([]complex128, len(aCopy))
	copy(want, aCopy)

	applyRot := func(idx int) {
		ct := cArr[idx]
		st := sArr[idx]
		if ct == 1 && st == 0 {
			return
		}
		var p, q int
		// p, q are the two indices of the planes rotated.
		if side == blas.Left {
			switch pivot {
			case lapack.Variable:
				p = idx
				q = idx + 1
			case lapack.Top:
				p = 0
				q = idx + 1
			case lapack.Bottom:
				p = idx
				q = m - 1
			}
		} else {
			switch pivot {
			case lapack.Variable:
				p = idx
				q = idx + 1
			case lapack.Top:
				p = 0
				q = idx + 1
			case lapack.Bottom:
				p = idx
				q = n - 1
			}
		}
		if side == blas.Left {
			// For Variable and Top pivots the rotation is
			//   (new_q, new_p) = (c*old_q - s*old_p, s*old_q + c*old_p)
			// For Bottom pivot the rotation is
			//   (new_p, new_q) = (s*old_q + c*old_p, c*old_q - s*old_p)
			for i := 0; i < n; i++ {
				tmp := want[q*lda+i]
				tmp2 := want[p*lda+i]
				if pivot == lapack.Bottom {
					want[p*lda+i] = complex(st, 0)*tmp + complex(ct, 0)*tmp2
					want[q*lda+i] = complex(ct, 0)*tmp - complex(st, 0)*tmp2
				} else {
					want[q*lda+i] = complex(ct, 0)*tmp - complex(st, 0)*tmp2
					want[p*lda+i] = complex(st, 0)*tmp + complex(ct, 0)*tmp2
				}
			}
		} else {
			for i := 0; i < m; i++ {
				tmp := want[i*lda+q]
				tmp2 := want[i*lda+p]
				if pivot == lapack.Bottom {
					want[i*lda+p] = complex(st, 0)*tmp + complex(ct, 0)*tmp2
					want[i*lda+q] = complex(ct, 0)*tmp - complex(st, 0)*tmp2
				} else {
					want[i*lda+q] = complex(ct, 0)*tmp - complex(st, 0)*tmp2
					want[i*lda+p] = complex(st, 0)*tmp + complex(ct, 0)*tmp2
				}
			}
		}
	}

	if direct == lapack.Forward {
		for i := 0; i < rotCount; i++ {
			applyRot(i)
		}
	} else {
		for i := rotCount - 1; i >= 0; i-- {
			applyRot(i)
		}
	}

	maxDiff := 0.0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			d := cmplx.Abs(a[i*lda+j] - want[i*lda+j])
			if d > maxDiff {
				maxDiff = d
			}
		}
	}
	bound := tol * float64(max(1, max(m, n))) * float64(max(1, rotCount))
	if maxDiff > bound {
		t.Errorf("%s: mismatch vs reference; maxDiff=%v, tol=%v", name, maxDiff, bound)
	}
}
