// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dtrevc3er interface {
	Dtrevc3(side lapack.EVSide, howmny lapack.EVHowMany, selected []bool, n int, t []float64, ldt int, vl []float64, ldvl int, vr []float64, ldvr int, mm int, work []float64, lwork int) int
}

func Dtrevc3Test(t *testing.T, impl Dtrevc3er) {
	rnd := rand.New(rand.NewSource(1))

	for _, side := range []lapack.EVSide{lapack.EVRight, lapack.EVLeft, lapack.EVBoth} {
		var name string
		switch side {
		case lapack.EVRight:
			name = "EVRigth"
		case lapack.EVLeft:
			name = "EVLeft"
		case lapack.EVBoth:
			name = "EVBoth"
		}
		t.Run(name, func(t *testing.T) {
			runDtrevc3Test(t, impl, rnd, side)
		})
	}
}

func runDtrevc3Test(t *testing.T, impl Dtrevc3er, rnd *rand.Rand, side lapack.EVSide) {
	for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 7, 10, 34} {
		for _, extra := range []int{0, 11} {
			for _, optwork := range []bool{true, false} {
				for cas := 0; cas < 10; cas++ {
					dtrevc3Test(t, impl, side, n, extra, optwork, rnd)
				}
			}
		}
	}
}

// dtrevc3Test tests Dtrevc3 by generating a random matrix T in Schur canonical
// form and performing the following checks:
//  1. Compute all eigenvectors of T and check that they are indeed correctly
//     normalized eigenvectors
//  2. Compute selected eigenvectors and check that they are exactly equal to
//     eigenvectors from check 1.
//  3. Compute all eigenvectors multiplied into a matrix Q and check that the
//     result is equal to eigenvectors from step 1 multiplied by Q and scaled
//     appropriately.
func dtrevc3Test(t *testing.T, impl Dtrevc3er, side lapack.EVSide, n, extra int, optwork bool, rnd *rand.Rand) {
	const tol = 1e-15

	name := fmt.Sprintf("n=%d,extra=%d,optwk=%v", n, extra, optwork)

	right := side != lapack.EVLeft
	left := side != lapack.EVRight

	// Generate a random matrix in Schur canonical form possibly with tiny or zero eigenvalues.
	// Zero elements of wi signify a real eigenvalue.
	tmat, wr, wi := randomSchurCanonical(n, n+extra, true, rnd)
	tmatCopy := cloneGeneral(tmat)

	//  1. Compute all eigenvectors of T and check that they are indeed correctly
	//     normalized eigenvectors

	howmny := lapack.EVAll

	var vr, vl blas64.General
	if right {
		// Fill VR and VL with NaN because they should be completely overwritten in Dtrevc3.
		vr = nanGeneral(n, n, n+extra)
	}
	if left {
		vl = nanGeneral(n, n, n+extra)
	}

	var work []float64
	if optwork {
		work = []float64{0}
		impl.Dtrevc3(side, howmny, nil, n, tmat.Data, tmat.Stride,
			vl.Data, max(1, vl.Stride), vr.Data, max(1, vr.Stride), n, work, -1)
		work = make([]float64, int(work[0]))
	} else {
		work = make([]float64, max(1, 3*n))
	}

	mGot := impl.Dtrevc3(side, howmny, nil, n, tmat.Data, tmat.Stride,
		vl.Data, max(1, vl.Stride), vr.Data, max(1, vr.Stride), n, work, len(work))

	if !generalOutsideAllNaN(tmat) {
		t.Errorf("%v: out-of-range write to T", name)
	}
	if !equalGeneral(tmat, tmatCopy) {
		t.Errorf("%v: unexpected modification of T", name)
	}
	if !generalOutsideAllNaN(vr) {
		t.Errorf("%v: out-of-range write to VR", name)
	}
	if !generalOutsideAllNaN(vl) {
		t.Errorf("%v: out-of-range write to VL", name)
	}

	mWant := n
	if mGot != mWant {
		t.Errorf("%v: unexpected value of m=%d, want %d", name, mGot, mWant)
	}

	if right {
		resid := residualRightEV(tmat, vr, wr, wi)
		if resid > tol {
			t.Errorf("%v: unexpected right eigenvectors; residual=%v, want<=%v", name, resid, tol)
		}
		resid = residualEVNormalization(vr, wi)
		if resid > tol {
			t.Errorf("%v: unexpected normalization of right eigenvectors; residual=%v, want<=%v", name, resid, tol)
		}
	}
	if left {
		resid := residualLeftEV(tmat, vl, wr, wi)
		if resid > tol {
			t.Errorf("%v: unexpected left eigenvectors; residual=%v, want<=%v", name, resid, tol)
		}
		resid = residualEVNormalization(vl, wi)
		if resid > tol {
			t.Errorf("%v: unexpected normalization of left eigenvectors; residual=%v, want<=%v", name, resid, tol)
		}
	}

	//  2. Compute selected eigenvectors and check that they are exactly equal to
	//     eigenvectors from check 1.

	howmny = lapack.EVSelected

	// Follow DCHKHS and select last max(1,n/4) real, max(1,n/4) complex
	// eigenvectors instead of selecting them randomly.
	selected := make([]bool, n)
	selectedWant := make([]bool, n)
	var nselr, nselc int
	for j := n - 1; j > 0; {
		if wi[j] == 0 {
			if nselr < max(1, n/4) {
				nselr++
				selected[j] = true
				selectedWant[j] = true
			}
			j--
		} else {
			if nselc < max(1, n/4) {
				nselc++
				// Select all columns to check that Dtrevc3 normalizes 'selected' correctly.
				selected[j] = true
				selected[j-1] = true
				selectedWant[j] = false
				selectedWant[j-1] = true
			}
			j -= 2
		}
	}
	mWant = nselr + 2*nselc

	var vrSel, vlSel blas64.General
	if right {
		vrSel = nanGeneral(n, mWant, n+extra)
	}
	if left {
		vlSel = nanGeneral(n, mWant, n+extra)
	}

	if optwork {
		// Reallocate optimal work in case it depends on howmny and selected.
		work = []float64{0}
		impl.Dtrevc3(side, howmny, selected, n, tmat.Data, tmat.Stride,
			vlSel.Data, max(1, vlSel.Stride), vrSel.Data, max(1, vrSel.Stride), mWant, work, -1)
		work = make([]float64, int(work[0]))
	}

	mGot = impl.Dtrevc3(side, howmny, selected, n, tmat.Data, tmat.Stride,
		vlSel.Data, max(1, vlSel.Stride), vrSel.Data, max(1, vrSel.Stride), mWant, work, len(work))

	if !generalOutsideAllNaN(tmat) {
		t.Errorf("%v: out-of-range write to T", name)
	}
	if !equalGeneral(tmat, tmatCopy) {
		t.Errorf("%v: unexpected modification of T", name)
	}
	if !generalOutsideAllNaN(vrSel) {
		t.Errorf("%v: out-of-range write to selected VR", name)
	}
	if !generalOutsideAllNaN(vlSel) {
		t.Errorf("%v: out-of-range write to selected VL", name)
	}

	if mGot != mWant {
		t.Errorf("%v: unexpected value of selected m=%d, want %d", name, mGot, mWant)
	}

	for i := range selected {
		if selected[i] != selectedWant[i] {
			t.Errorf("%v: unexpected selected[%v]", name, i)
		}
	}

	// Check that selected columns of vrSel are equal to the corresponding
	// columns of vr.
	var k int
	match := true
	if right {
	loopVR:
		for j := 0; j < n; j++ {
			if selected[j] && wi[j] == 0 {
				for i := 0; i < n; i++ {
					if vrSel.Data[i*vrSel.Stride+k] != vr.Data[i*vr.Stride+j] {
						match = false
						break loopVR
					}
				}
				k++
			} else if selected[j] && wi[j] != 0 {
				for i := 0; i < n; i++ {
					if vrSel.Data[i*vrSel.Stride+k] != vr.Data[i*vr.Stride+j] ||
						vrSel.Data[i*vrSel.Stride+k+1] != vr.Data[i*vr.Stride+j+1] {
						match = false
						break loopVR
					}
				}
				k += 2
			}
		}
	}
	if !match {
		t.Errorf("%v: unexpected selected VR", name)
	}

	// Check that selected columns of vlSel are equal to the corresponding
	// columns of vl.
	match = true
	k = 0
	if left {
	loopVL:
		for j := 0; j < n; j++ {
			if selected[j] && wi[j] == 0 {
				for i := 0; i < n; i++ {
					if vlSel.Data[i*vlSel.Stride+k] != vl.Data[i*vl.Stride+j] {
						match = false
						break loopVL
					}
				}
				k++
			} else if selected[j] && wi[j] != 0 {
				for i := 0; i < n; i++ {
					if vlSel.Data[i*vlSel.Stride+k] != vl.Data[i*vl.Stride+j] ||
						vlSel.Data[i*vlSel.Stride+k+1] != vl.Data[i*vl.Stride+j+1] {
						match = false
						break loopVL
					}
				}
				k += 2
			}
		}
	}
	if !match {
		t.Errorf("%v: unexpected selected VL", name)
	}

	//  3. Compute all eigenvectors multiplied into a matrix Q and check that the
	//     result is equal to eigenvectors from step 1 multiplied by Q and scaled
	//     appropriately.

	howmny = lapack.EVAllMulQ

	var vrMul, qr blas64.General
	var vlMul, ql blas64.General
	if right {
		vrMul = randomGeneral(n, n, n+extra, rnd)
		qr = cloneGeneral(vrMul)
	}
	if left {
		vlMul = randomGeneral(n, n, n+extra, rnd)
		ql = cloneGeneral(vlMul)
	}

	if optwork {
		// Reallocate optimal work in case it depends on howmny and selected.
		work = []float64{0}
		impl.Dtrevc3(side, howmny, nil, n, tmat.Data, tmat.Stride,
			vlMul.Data, max(1, vlMul.Stride), vrMul.Data, max(1, vrMul.Stride), n, work, -1)
		work = make([]float64, int(work[0]))
	}

	mGot = impl.Dtrevc3(side, howmny, selected, n, tmat.Data, tmat.Stride,
		vlMul.Data, max(1, vlMul.Stride), vrMul.Data, max(1, vrMul.Stride), n, work, len(work))

	if !generalOutsideAllNaN(tmat) {
		t.Errorf("%v: out-of-range write to T", name)
	}
	if !equalGeneral(tmat, tmatCopy) {
		t.Errorf("%v: unexpected modification of T", name)
	}
	if !generalOutsideAllNaN(vrMul) {
		t.Errorf("%v: out-of-range write to VRMul", name)
	}
	if !generalOutsideAllNaN(vlMul) {
		t.Errorf("%v: out-of-range write to VLMul", name)
	}

	mWant = n
	if mGot != mWant {
		t.Errorf("%v: unexpected value of m=%d, want %d", name, mGot, mWant)
	}

	if right {
		// Compute Q * VR explicitly and normalize to match Dtrevc3 output.
		qvWant := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qr, vr, 0, qvWant)
		normalizeEV(qvWant, wi)

		// Compute the difference between Dtrevc3 output and Q * VR.
		r := zeros(n, n, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				r.Data[i*r.Stride+j] = vrMul.Data[i*vrMul.Stride+j] - qvWant.Data[i*qvWant.Stride+j]
			}
		}
		qvNorm := dlange(lapack.MaxColumnSum, n, n, qvWant.Data, qvWant.Stride)
		resid := dlange(lapack.MaxColumnSum, n, n, r.Data, r.Stride) / qvNorm / float64(n)
		if resid > tol {
			t.Errorf("%v: unexpected VRMul; resid=%v, want <=%v", name, resid, tol)
		}
	}
	if left {
		// Compute Q * VL explicitly and normalize to match Dtrevc3 output.
		qvWant := zeros(n, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, ql, vl, 0, qvWant)
		normalizeEV(qvWant, wi)

		// Compute the difference between Dtrevc3 output and Q * VL.
		r := zeros(n, n, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				r.Data[i*r.Stride+j] = vlMul.Data[i*vlMul.Stride+j] - qvWant.Data[i*qvWant.Stride+j]
			}
		}
		qvNorm := dlange(lapack.MaxColumnSum, n, n, qvWant.Data, qvWant.Stride)
		resid := dlange(lapack.MaxColumnSum, n, n, r.Data, r.Stride) / qvNorm / float64(n)
		if resid > tol {
			t.Errorf("%v: unexpected VLMul; resid=%v, want <=%v", name, resid, tol)
		}
	}
}

// residualEVNormalization returns the maximum normalization error in E:
//
//	max |max-norm(E[:,j]) - 1|
func residualEVNormalization(emat blas64.General, wi []float64) float64 {
	n := emat.Rows
	if n == 0 {
		return 0
	}
	var (
		e      = emat.Data
		lde    = emat.Stride
		enrmin = math.Inf(1)
		enrmax float64
		ipair  int
	)
	for j := 0; j < n; j++ {
		if ipair == 0 && j < n-1 && wi[j] != 0 {
			ipair = 1
		}
		var nrm float64
		switch ipair {
		case 0:
			// Real eigenvector
			for i := 0; i < n; i++ {
				nrm = math.Max(nrm, math.Abs(e[i*lde+j]))
			}
			enrmin = math.Min(enrmin, nrm)
			enrmax = math.Max(enrmax, nrm)
		case 1:
			// Complex eigenvector
			for i := 0; i < n; i++ {
				nrm = math.Max(nrm, math.Abs(e[i*lde+j])+math.Abs(e[i*lde+j+1]))
			}
			enrmin = math.Min(enrmin, nrm)
			enrmax = math.Max(enrmax, nrm)
			ipair = 2
		case 2:
			ipair = 0
		}
	}
	return math.Max(math.Abs(enrmin-1), math.Abs(enrmin-1))
}

// normalizeEV normalizes eigenvectors in the columns of E so that the element
// of largest magnitude has magnitude 1.
func normalizeEV(emat blas64.General, wi []float64) {
	n := emat.Rows
	if n == 0 {
		return
	}
	var (
		bi    = blas64.Implementation()
		e     = emat.Data
		lde   = emat.Stride
		ipair int
	)
	for j := 0; j < n; j++ {
		if ipair == 0 && j < n-1 && wi[j] != 0 {
			ipair = 1
		}
		switch ipair {
		case 0:
			// Real eigenvector
			ii := bi.Idamax(n, e[j:], lde)
			remax := 1 / math.Abs(e[ii*lde+j])
			bi.Dscal(n, remax, e[j:], lde)
		case 1:
			// Complex eigenvector
			var emax float64
			for i := 0; i < n; i++ {
				emax = math.Max(emax, math.Abs(e[i*lde+j])+math.Abs(e[i*lde+j+1]))
			}
			bi.Dscal(n, 1/emax, e[j:], lde)
			bi.Dscal(n, 1/emax, e[j+1:], lde)
			ipair = 2
		case 2:
			ipair = 0
		}
	}
}
