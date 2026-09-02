// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dtrsener interface {
	Dtrexcer
	Dtrsen(ijob int, wantq bool, selected []bool, n int, t []float64, ldt int, q []float64, ldq int, wr, wi []float64, work []float64, lwork int, iwork []int, liwork int) (m int, s, sep float64, ok bool)
}

func DtrsenTest(t *testing.T, impl Dtrsener) {
	rnd := rand.New(rand.NewPCG(1, 1))

	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20} {
		for _, extra := range []int{0, 3} {
			for _, selectPattern := range []string{"none", "first", "last", "alternating", "all"} {
				for _, wantq := range []bool{false, true} {
					for _, ijob := range []int{0, 1, 2, 3} {
						testDtrsen(t, impl, n, extra, selectPattern, wantq, ijob, rnd)
					}
				}
			}
		}
	}

	testDtrsenWorkspace(t, impl)
	testDtrsenWorkspaceOutputs(t, impl)
	testDtrsenSelectSecondOfPair(t, impl)
}

func testDtrsenWorkspaceOutputs(t *testing.T, impl Dtrsener) {
	work := []float64{math.NaN(), math.NaN()}
	iwork := []int{-1}
	_, _, _, ok := impl.Dtrsen(0, false, []bool{false, false}, 2,
		[]float64{1, 0, 0, 2}, 2, nil, 1,
		make([]float64, 2), make([]float64, 2), work, len(work), iwork, len(iwork))
	if !ok {
		t.Fatal("trivial reorder failed")
	}
	if work[0] != 2 || iwork[0] != 1 {
		t.Errorf("workspace outputs=(%v,%d), want (2,1)", work[0], iwork[0])
	}
}

func testDtrsenSelectSecondOfPair(t *testing.T, impl Dtrsener) {
	work := make([]float64, 2)
	iwork := make([]int, 1)
	m, _, _, ok := impl.Dtrsen(0, false, []bool{false, true}, 2,
		[]float64{1, 2, -3, 1}, 2, nil, 1,
		make([]float64, 2), make([]float64, 2), work, len(work), iwork, len(iwork))
	if !ok {
		t.Fatal("complex pair selection failed")
	}
	if m != 2 {
		t.Fatalf("selected subspace dimension=%d, want 2", m)
	}
}

type eigPair struct{ wr, wi float64 }

func testDtrsen(t *testing.T, impl Dtrsener, n, extra int, selectPattern string, wantq bool, ijob int, rnd *rand.Rand) {
	const tol = 100

	prefix := fmt.Sprintf("n=%d, extra=%d, select=%s, wantq=%v, ijob=%d", n, extra, selectPattern, wantq, ijob)

	if n == 0 {
		var work [1]float64
		var iwork [1]int
		m, s, sep, ok := impl.Dtrsen(ijob, wantq, nil, 0, nil, 1, nil, 1, nil, nil, work[:], 1, iwork[:], 1)
		if m != 0 || s != 1 || sep != 0 || !ok {
			t.Errorf("%s: unexpected result for n=0", prefix)
		}
		return
	}

	stride := n + extra
	tMat, wrOrig, wiOrig := randomSchurCanonical(n, stride, false, rnd)
	tOrig := cloneGeneral(tMat)

	// Generate selection array and record selected eigenvalues.
	selected := make([]bool, n)
	var selectedEigs []eigPair
	expectedM := 0

	switch selectPattern {
	case "none":
	case "first":
		selected[0] = true
		selectedEigs = append(selectedEigs, eigPair{wrOrig[0], wiOrig[0]})
		expectedM = 1
		if n > 1 && tMat.Data[1*tMat.Stride+0] != 0 {
			selected[1] = true
			selectedEigs = append(selectedEigs, eigPair{wrOrig[1], wiOrig[1]})
			expectedM = 2
		}
	case "last":
		k := n - 1
		if k > 0 && tMat.Data[k*tMat.Stride+k-1] != 0 {
			k--
			selected[k] = true
			selected[k+1] = true
			selectedEigs = append(selectedEigs, eigPair{wrOrig[k], wiOrig[k]})
			selectedEigs = append(selectedEigs, eigPair{wrOrig[k+1], wiOrig[k+1]})
			expectedM = 2
		} else {
			selected[k] = true
			selectedEigs = append(selectedEigs, eigPair{wrOrig[k], wiOrig[k]})
			expectedM = 1
		}
	case "alternating":
		sel := true
		for k := 0; k < n; {
			if k < n-1 && tMat.Data[(k+1)*tMat.Stride+k] != 0 {
				if sel {
					selected[k] = true
					selected[k+1] = true
					selectedEigs = append(selectedEigs, eigPair{wrOrig[k], wiOrig[k]})
					selectedEigs = append(selectedEigs, eigPair{wrOrig[k+1], wiOrig[k+1]})
					expectedM += 2
				}
				k += 2
			} else {
				if sel {
					selected[k] = true
					selectedEigs = append(selectedEigs, eigPair{wrOrig[k], wiOrig[k]})
					expectedM++
				}
				k++
			}
			sel = !sel
		}
	case "all":
		for i := 0; i < n; i++ {
			selected[i] = true
			selectedEigs = append(selectedEigs, eigPair{wrOrig[i], wiOrig[i]})
		}
		expectedM = n
	}

	// Set up Q.
	ldq := 1
	var qData []float64
	if wantq {
		ldq = stride
		q := eye(n, ldq)
		qData = q.Data
	}

	wr := make([]float64, n)
	wi := make([]float64, n)

	// Workspace query.
	var worksize [1]float64
	var iworksize [1]int
	impl.Dtrsen(ijob, wantq, selected, n, tMat.Data, max(1, tMat.Stride), nil, ldq, nil, nil, worksize[:], -1, iworksize[:], -1)
	lwork := int(worksize[0])
	liwork := iworksize[0]

	work := make([]float64, max(1, lwork))
	iwork := make([]int, max(1, liwork))

	m, s, sep, ok := impl.Dtrsen(ijob, wantq, selected, n, tMat.Data, tMat.Stride, qData, ldq, wr, wi, work, lwork, iwork, liwork)

	if !ok {
		t.Logf("%s: reordering failed (ok=false)", prefix)
		return
	}

	if m != expectedM {
		t.Errorf("%s: m=%d, want %d", prefix, m, expectedM)
	}

	// Verify T is still in Schur canonical form.
	if !isSchurCanonicalGeneral(tMat) {
		t.Errorf("%s: result is not in Schur canonical form", prefix)
	}

	// Verify eigenvalues in wr/wi match diagonal blocks.
	evTol := 1e-10
	checkSchurEigenvalues(t, prefix, tMat, wr, wi, evTol)

	// Verify selected eigenvalues moved to top.
	if m > 0 && m < n {
		topEigs := make([]eigPair, m)
		for i := 0; i < m; i++ {
			topEigs[i] = eigPair{wr[i], wi[i]}
		}
		if !eigMultisetsMatch(selectedEigs, topEigs, evTol) {
			t.Errorf("%s: selected eigenvalues not at top\n  selected: %v\n  top:      %v", prefix, selectedEigs, topEigs)
		}
	}

	// Verify condition number estimates.
	wants := ijob == 1 || ijob == 3
	wantsp := ijob == 2 || ijob == 3
	if m == 0 || m == n {
		if wants && s != 1 {
			t.Errorf("%s: s=%v, want 1 for m=0 or m=n", prefix, s)
		}
		wantSep := dlange(lapack.MaxColumnSum, n, n, tMat.Data, tMat.Stride)
		if wantsp && sep != wantSep {
			t.Errorf("%s: sep=%v, want %v for m=0 or m=n", prefix, sep, wantSep)
		}
	} else {
		if wants {
			if s <= 0 || s > 1 {
				t.Errorf("%s: s=%v, want in (0, 1]", prefix, s)
			}
		}
		if wantsp {
			if sep < 0 {
				t.Errorf("%s: sep=%v, want >= 0", prefix, sep)
			}
		}
	}

	if !wantq {
		return
	}

	q := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: ldq,
		Data:   qData,
	}

	// Verify Q is orthogonal.
	orthoTol := float64(n) * dlamchE
	if orthoTol < 1e-13 {
		orthoTol = 1e-13
	}
	resid := residualOrthogonal(q, false)
	if resid > orthoTol {
		t.Errorf("%s: Q not orthogonal, residual=%e, want<=%e", prefix, resid, orthoTol)
	}

	// Verify Q*T*Qᵀ = Torig.
	residSchur := residualSchurFactorization(tOrig, tMat, q)
	anorm := dlange(lapack.MaxColumnSum, n, n, tOrig.Data, tOrig.Stride)
	if anorm == 0 {
		anorm = 1
	}
	normRes := residSchur / (anorm * float64(n) * dlamchE)
	if normRes > float64(tol) {
		t.Errorf("%s: ||T_orig - Q*T*Qᵀ||/(||T||*n*eps)=%v, want<=%v", prefix, normRes, tol)
	}
}

func testDtrsenWorkspace(t *testing.T, impl Dtrsener) {
	for _, n := range []int{1, 5, 10} {
		var worksize [1]float64
		var iworksize [1]int
		impl.Dtrsen(0, true, nil, n, nil, n, nil, n, nil, nil, worksize[:], -1, iworksize[:], -1)
		lwork := int(worksize[0])
		liwork := iworksize[0]
		if lwork < 1 {
			t.Errorf("n=%d: workspace query returned lwork=%d, want >= 1", n, lwork)
		}
		if liwork < 1 {
			t.Errorf("n=%d: workspace query returned liwork=%d, want >= 1", n, liwork)
		}
	}

	var worksize [1]float64
	var iworksize [1]int
	impl.Dtrsen(2, false, []bool{false, true, false, false}, 4,
		[]float64{
			1, 2, 0, 0,
			-3, 1, 0, 0,
			0, 0, 4, 0,
			0, 0, 0, 5,
		}, 4, nil, 1, nil, nil, worksize[:], -1, iworksize[:], -1)
	if worksize[0] != 8 || iworksize[0] != 4 {
		t.Fatalf("complex-pair workspace=(%v,%d), want (8,4)", worksize[0], iworksize[0])
	}
}

func eigMultisetsMatch(a, b []eigPair, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	used := make([]bool, len(b))
	for _, ea := range a {
		found := false
		for j, eb := range b {
			if used[j] {
				continue
			}
			if math.Abs(ea.wr-eb.wr) < tol && math.Abs(ea.wi-eb.wi) < tol {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
