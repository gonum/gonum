// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dtrexcer interface {
	Dtrexc(compq lapack.UpdateSchurComp, n int, t []float64, ldt int, q []float64, ldq int, ifst, ilst int, work []float64) (ifstOut, ilstOut int, ok bool)
}

func DtrexcTest(t *testing.T, impl Dtrexcer) {
	rnd := rand.New(rand.NewSource(1))

	for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 10, 18, 31, 53} {
		for _, extra := range []int{0, 3} {
			for cas := 0; cas < 100; cas++ {
				var ifst, ilst int
				if n > 0 {
					ifst = rnd.Intn(n)
					ilst = rnd.Intn(n)
				}
				dtrexcTest(t, impl, rnd, n, ifst, ilst, extra)
			}
		}
	}
}

func dtrexcTest(t *testing.T, impl Dtrexcer, rnd *rand.Rand, n, ifst, ilst, extra int) {
	const tol = 1e-13

	tmatCopy, _, _ := randomSchurCanonical(n, n+extra, false, rnd)
	tmat1 := cloneGeneral(tmatCopy)
	tmat2 := cloneGeneral(tmatCopy)

	fstSize, fstFirst := schurBlockSize(tmat1, ifst)
	lstSize, lstFirst := schurBlockSize(tmat1, ilst)

	name := fmt.Sprintf("Case n=%v,ifst=%v,nbfst=%v,ilst=%v,nblst=%v,extra=%v",
		n, ifst, fstSize, ilst, lstSize, extra)

	// 1. Test without accumulating Q.

	compq := lapack.UpdateSchurNone
	work := nanSlice(n)
	ifstGot1, ilstGot1, ok1 := impl.Dtrexc(compq, n, tmat1.Data, tmat1.Stride, nil, 1, ifst, ilst, work)

	if !generalOutsideAllNaN(tmat1) {
		t.Errorf("%v: out-of-range write to T1", name)
	}

	// 2. Test with accumulating Q.

	q := eye(n, n+extra)
	qCopy := cloneGeneral(q)

	compq = lapack.UpdateSchur
	work = nanSlice(n)
	ifstGot2, ilstGot2, ok2 := impl.Dtrexc(compq, n, tmat2.Data, tmat2.Stride, q.Data, q.Stride, ifst, ilst, work)

	if !generalOutsideAllNaN(q) {
		t.Errorf("%v: out-of-range write to Q", name)
	}

	// Check that outputs from cases 1. and 2. are exactly equal, then check one of them.
	if ifstGot1 != ifstGot2 {
		t.Errorf("%v: ifstGot1 != ifstGot2", name)
	}
	if ilstGot1 != ilstGot2 {
		t.Errorf("%v: ilstGot1 != ilstGot2", name)
	}
	if ok1 != ok2 {
		t.Errorf("%v: ok1 != ok2", name)
	}
	if !equalGeneral(tmat1, tmat2) {
		t.Errorf("%v: T1 != T2", name)
	}

	// Check that the index of the first block was correctly updated (if
	// necessary).
	ifstWant := ifst
	if !fstFirst {
		ifstWant = ifst - 1
	}
	if ifstWant != ifstGot1 {
		t.Errorf("%v: unexpected ifst=%v, want %v", name, ifstGot1, ifstWant)
	}

	// Check that the index of the last block is as expected when ok=true.
	// When ok=false, we don't know at which block the algorithm failed, so
	// we don't check.
	ilstWant := ilst
	if !lstFirst {
		ilstWant--
	}
	if ok1 {
		if ifstWant < ilstWant {
			// If the blocks are swapped backwards, these
			// adjustments are not necessary, the first row of the
			// last block will end up at ifst.
			switch {
			case fstSize == 2 && lstSize == 1:
				ilstWant--
			case fstSize == 1 && lstSize == 2:
				ilstWant++
			}
		}
		if ilstWant != ilstGot1 {
			t.Errorf("%v: unexpected ilst=%v, want %v", name, ilstGot1, ilstWant)
		}
	}

	if n <= 1 || ifstGot1 == ilstGot1 {
		// Too small matrix or no swapping.
		// Check that T was not modified.
		if !equalGeneral(tmat1, tmatCopy) {
			t.Errorf("%v: unexpected modification of T when no swapping", name)
		}
		// Check that Q was not modified.
		if !equalGeneral(q, qCopy) {
			t.Errorf("%v: unexpected modification of Q when no swapping", name)
		}
		// Nothing more to check
		return
	}

	if !isSchurCanonicalGeneral(tmat1) {
		t.Errorf("%v: T is not in Schur canonical form", name)
	}

	// Check that T was not modified except above the second subdiagonal in
	// rows and columns [modMin,modMax].
	modMin := min(ifstGot1, ilstGot1)
	modMax := max(ifstGot1, ilstGot1) + fstSize
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if modMin <= i && i < modMax && j+1 >= i {
				continue
			}
			if modMin <= j && j < modMax && j+1 >= i {
				continue
			}
			diff := tmat1.Data[i*tmat1.Stride+j] - tmatCopy.Data[i*tmatCopy.Stride+j]
			if diff != 0 {
				t.Errorf("%v: unexpected modification at T[%v,%v]", name, i, j)
			}
		}
	}

	// Check that the block at ifstGot was delivered to ilstGot correctly.
	if fstSize == 1 {
		// 1×1 blocks are swapped exactly.
		got := tmat1.Data[ilstGot1*tmat1.Stride+ilstGot1]
		want := tmatCopy.Data[ifstGot1*tmatCopy.Stride+ifstGot1]
		if want != got {
			t.Errorf("%v: unexpected 1×1 block at T[%v,%v]. Want %v, got %v",
				name, want, got, ilstGot1, ilstGot1)
		}
	} else {
		// Check that the swapped 2×2 block has the same eigenvalues.
		a, b, c, d := extract2x2Block(tmat1.Data[ilstGot1*tmat1.Stride+ilstGot1:], tmat1.Stride)
		ev1Got, ev2Got := schurBlockEigenvalues(a, b, c, d)
		// The block was originally located at T[ifstGot,ifstGot].
		a, b, c, d = extract2x2Block(tmatCopy.Data[ifstGot1*tmatCopy.Stride+ifstGot1:], tmatCopy.Stride)
		ev1Want, ev2Want := schurBlockEigenvalues(a, b, c, d)
		diff := cmplx.Abs(ev1Got - ev1Want)
		if diff > tol {
			t.Errorf("%v: unexpected first eigenvalue of 2×2 block [%v %v; %v %v] at T[%v,%v]; diff=%v, want<=%v",
				name, a, b, c, d, ilstGot1, ilstGot1, diff, tol)
		}
		diff = cmplx.Abs(ev2Got - ev2Want)
		if diff > tol {
			t.Errorf("%v: unexpected second eigenvalue of 2×2 block [%v %v; %v %v] at T[%v,%v]. Want %v, got %v",
				name, a, b, c, d, ilstGot1, ilstGot1, diff, tol)
		}
	}

	// Check that Q is orthogonal.
	resid := residualOrthogonal(q, false)
	if resid > float64(n)*tol {
		t.Errorf("%v: Q is not orthogonal; resid=%v, want<=%v", name, resid, float64(n)*tol)
	}
	// Check that Q is unchanged outside of columns [modMin,modMax].
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if modMin <= j && j < modMax {
				continue
			}
			if q.Data[i*q.Stride+j]-qCopy.Data[i*qCopy.Stride+j] != 0 {
				t.Errorf("%v: unexpected modification of Q[%v,%v]", name, i, j)
			}
		}
	}
	// Check that Qᵀ * TOrig * Q == T
	qt := zeros(n, n, n)
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, tmatCopy, 0, qt)
	qtq := cloneGeneral(tmat1)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, -1, qt, q, 1, qtq)
	resid = dlange(lapack.MaxColumnSum, n, n, qtq.Data, qtq.Stride)
	if resid > float64(n)*tol {
		t.Errorf("%v: mismatch between Qᵀ*(initial T)*Q and (final T); resid=%v, want<=%v",
			name, resid, float64(n)*tol)
	}
}

func residualOrthogonal(q blas64.General, rowwise bool) float64 {
	m, n := q.Rows, q.Cols
	if m == 0 || n == 0 {
		return 0
	}
	var transq blas.Transpose
	if m < n || (m == n && rowwise) {
		transq = blas.NoTrans
	} else {
		transq = blas.Trans
	}
	minmn := min(m, n)

	// Set work = I.
	work := blas64.Symmetric{
		Uplo:   blas.Upper,
		N:      minmn,
		Data:   make([]float64, minmn*minmn),
		Stride: minmn,
	}
	for i := 0; i < minmn; i++ {
		work.Data[i*work.Stride+i] = 1
	}

	// Compute
	//  work = work - Q * Qᵀ = I - Q * Qᵀ
	// or
	//  work = work - Qᵀ * Q = I - Qᵀ * Q
	blas64.Syrk(transq, -1, q, 1, work)
	return dlansy(lapack.MaxColumnSum, blas.Upper, work.N, work.Data, work.Stride)
}

func dlansy(norm lapack.MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int) float64 {
	if n == 0 {
		return 0
	}
	work := make([]float64, n)
	switch norm {
	case lapack.MaxAbs:
		if uplo == blas.Upper {
			var max float64
			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					v := math.Abs(a[i*lda+j])
					if math.IsNaN(v) {
						return math.NaN()
					}
					if v > max {
						max = v
					}
				}
			}
			return max
		}
		var max float64
		for i := 0; i < n; i++ {
			for j := 0; j <= i; j++ {
				v := math.Abs(a[i*lda+j])
				if math.IsNaN(v) {
					return math.NaN()
				}
				if v > max {
					max = v
				}
			}
		}
		return max
	case lapack.MaxRowSum, lapack.MaxColumnSum:
		// A symmetric matrix has the same 1-norm and ∞-norm.
		for i := 0; i < n; i++ {
			work[i] = 0
		}
		if uplo == blas.Upper {
			for i := 0; i < n; i++ {
				work[i] += math.Abs(a[i*lda+i])
				for j := i + 1; j < n; j++ {
					v := math.Abs(a[i*lda+j])
					work[i] += v
					work[j] += v
				}
			}
		} else {
			for i := 0; i < n; i++ {
				for j := 0; j < i; j++ {
					v := math.Abs(a[i*lda+j])
					work[i] += v
					work[j] += v
				}
				work[i] += math.Abs(a[i*lda+i])
			}
		}
		var max float64
		for i := 0; i < n; i++ {
			v := work[i]
			if math.IsNaN(v) {
				return math.NaN()
			}
			if v > max {
				max = v
			}
		}
		return max
	default:
		// lapack.Frobenius:
		panic("not implemented")
	}
}
