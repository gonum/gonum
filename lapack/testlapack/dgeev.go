// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dgeever interface {
	Dgeev(jobvl lapack.LeftEVJob, jobvr lapack.RightEVJob, n int, a []float64, lda int,
		wr, wi []float64, vl []float64, ldvl int, vr []float64, ldvr int, work []float64, lwork int) int
}

type dgeevTest struct {
	a      blas64.General
	evWant []complex128 // If nil, the eigenvalues are not known.
	valTol float64      // Tolerance for eigenvalue checks.
	vecTol float64      // Tolerance for eigenvector checks.
}

func DgeevTest(t *testing.T, impl Dgeever) {
	rnd := rand.New(rand.NewSource(1))

	for i, test := range []dgeevTest{
		{
			a:      A123{}.Matrix(),
			evWant: A123{}.Eigenvalues(),
		},

		dgeevTestForAntisymRandom(10, rnd),
		dgeevTestForAntisymRandom(11, rnd),
		dgeevTestForAntisymRandom(50, rnd),
		dgeevTestForAntisymRandom(51, rnd),
		dgeevTestForAntisymRandom(100, rnd),
		dgeevTestForAntisymRandom(101, rnd),

		{
			a:      Circulant(2).Matrix(),
			evWant: Circulant(2).Eigenvalues(),
		},
		{
			a:      Circulant(3).Matrix(),
			evWant: Circulant(3).Eigenvalues(),
		},
		{
			a:      Circulant(4).Matrix(),
			evWant: Circulant(4).Eigenvalues(),
		},
		{
			a:      Circulant(5).Matrix(),
			evWant: Circulant(5).Eigenvalues(),
		},
		{
			a:      Circulant(10).Matrix(),
			evWant: Circulant(10).Eigenvalues(),
		},
		{
			a:      Circulant(15).Matrix(),
			evWant: Circulant(15).Eigenvalues(),
			valTol: 1e-12,
		},
		{
			a:      Circulant(30).Matrix(),
			evWant: Circulant(30).Eigenvalues(),
			valTol: 1e-11,
		},
		{
			a:      Circulant(50).Matrix(),
			evWant: Circulant(50).Eigenvalues(),
			valTol: 1e-11,
		},
		{
			a:      Circulant(101).Matrix(),
			evWant: Circulant(101).Eigenvalues(),
			valTol: 1e-10,
		},
		{
			a:      Circulant(150).Matrix(),
			evWant: Circulant(150).Eigenvalues(),
			valTol: 1e-9,
		},

		{
			a:      Clement(2).Matrix(),
			evWant: Clement(2).Eigenvalues(),
		},
		{
			a:      Clement(3).Matrix(),
			evWant: Clement(3).Eigenvalues(),
		},
		{
			a:      Clement(4).Matrix(),
			evWant: Clement(4).Eigenvalues(),
		},
		{
			a:      Clement(5).Matrix(),
			evWant: Clement(5).Eigenvalues(),
		},
		{
			a:      Clement(10).Matrix(),
			evWant: Clement(10).Eigenvalues(),
		},
		{
			a:      Clement(15).Matrix(),
			evWant: Clement(15).Eigenvalues(),
		},
		{
			a:      Clement(30).Matrix(),
			evWant: Clement(30).Eigenvalues(),
			valTol: 1e-11,
		},
		{
			a:      Clement(50).Matrix(),
			evWant: Clement(50).Eigenvalues(),
			valTol: 1e-8,
		},

		{
			a:      Creation(2).Matrix(),
			evWant: Creation(2).Eigenvalues(),
		},
		{
			a:      Creation(3).Matrix(),
			evWant: Creation(3).Eigenvalues(),
		},
		{
			a:      Creation(4).Matrix(),
			evWant: Creation(4).Eigenvalues(),
		},
		{
			a:      Creation(5).Matrix(),
			evWant: Creation(5).Eigenvalues(),
		},
		{
			a:      Creation(10).Matrix(),
			evWant: Creation(10).Eigenvalues(),
		},
		{
			a:      Creation(15).Matrix(),
			evWant: Creation(15).Eigenvalues(),
		},
		{
			a:      Creation(30).Matrix(),
			evWant: Creation(30).Eigenvalues(),
		},
		{
			a:      Creation(50).Matrix(),
			evWant: Creation(50).Eigenvalues(),
		},
		{
			a:      Creation(101).Matrix(),
			evWant: Creation(101).Eigenvalues(),
		},
		{
			a:      Creation(150).Matrix(),
			evWant: Creation(150).Eigenvalues(),
		},

		{
			a:      Diagonal(0).Matrix(),
			evWant: Diagonal(0).Eigenvalues(),
		},
		{
			a:      Diagonal(10).Matrix(),
			evWant: Diagonal(10).Eigenvalues(),
		},
		{
			a:      Diagonal(50).Matrix(),
			evWant: Diagonal(50).Eigenvalues(),
		},
		{
			a:      Diagonal(151).Matrix(),
			evWant: Diagonal(151).Eigenvalues(),
		},

		{
			a:      Downshift(2).Matrix(),
			evWant: Downshift(2).Eigenvalues(),
		},
		{
			a:      Downshift(3).Matrix(),
			evWant: Downshift(3).Eigenvalues(),
		},
		{
			a:      Downshift(4).Matrix(),
			evWant: Downshift(4).Eigenvalues(),
		},
		{
			a:      Downshift(5).Matrix(),
			evWant: Downshift(5).Eigenvalues(),
		},
		{
			a:      Downshift(10).Matrix(),
			evWant: Downshift(10).Eigenvalues(),
		},
		{
			a:      Downshift(15).Matrix(),
			evWant: Downshift(15).Eigenvalues(),
		},
		{
			a:      Downshift(30).Matrix(),
			evWant: Downshift(30).Eigenvalues(),
		},
		{
			a:      Downshift(50).Matrix(),
			evWant: Downshift(50).Eigenvalues(),
		},
		{
			a:      Downshift(101).Matrix(),
			evWant: Downshift(101).Eigenvalues(),
		},
		{
			a:      Downshift(150).Matrix(),
			evWant: Downshift(150).Eigenvalues(),
		},

		{
			a:      Fibonacci(2).Matrix(),
			evWant: Fibonacci(2).Eigenvalues(),
		},
		{
			a:      Fibonacci(3).Matrix(),
			evWant: Fibonacci(3).Eigenvalues(),
		},
		{
			a:      Fibonacci(4).Matrix(),
			evWant: Fibonacci(4).Eigenvalues(),
		},
		{
			a:      Fibonacci(5).Matrix(),
			evWant: Fibonacci(5).Eigenvalues(),
		},
		{
			a:      Fibonacci(10).Matrix(),
			evWant: Fibonacci(10).Eigenvalues(),
		},
		{
			a:      Fibonacci(15).Matrix(),
			evWant: Fibonacci(15).Eigenvalues(),
		},
		{
			a:      Fibonacci(30).Matrix(),
			evWant: Fibonacci(30).Eigenvalues(),
		},
		{
			a:      Fibonacci(50).Matrix(),
			evWant: Fibonacci(50).Eigenvalues(),
		},
		{
			a:      Fibonacci(101).Matrix(),
			evWant: Fibonacci(101).Eigenvalues(),
		},
		{
			a:      Fibonacci(150).Matrix(),
			evWant: Fibonacci(150).Eigenvalues(),
		},

		{
			a:      Gear(2).Matrix(),
			evWant: Gear(2).Eigenvalues(),
		},
		{
			a:      Gear(3).Matrix(),
			evWant: Gear(3).Eigenvalues(),
		},
		{
			a:      Gear(4).Matrix(),
			evWant: Gear(4).Eigenvalues(),
			valTol: 1e-7,
			vecTol: 1e-8,
		},
		{
			a:      Gear(5).Matrix(),
			evWant: Gear(5).Eigenvalues(),
		},
		{
			a:      Gear(10).Matrix(),
			evWant: Gear(10).Eigenvalues(),
			valTol: 1e-8,
		},
		{
			a:      Gear(15).Matrix(),
			evWant: Gear(15).Eigenvalues(),
		},
		{
			a:      Gear(30).Matrix(),
			evWant: Gear(30).Eigenvalues(),
			valTol: 1e-8,
		},
		{
			a:      Gear(50).Matrix(),
			evWant: Gear(50).Eigenvalues(),
			valTol: 1e-8,
		},
		{
			a:      Gear(101).Matrix(),
			evWant: Gear(101).Eigenvalues(),
		},
		{
			a:      Gear(150).Matrix(),
			evWant: Gear(150).Eigenvalues(),
			valTol: 1e-8,
		},

		{
			a:      Grcar{N: 10, K: 3}.Matrix(),
			evWant: Grcar{N: 10, K: 3}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 10, K: 7}.Matrix(),
			evWant: Grcar{N: 10, K: 7}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 11, K: 7}.Matrix(),
			evWant: Grcar{N: 11, K: 7}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 50, K: 3}.Matrix(),
			evWant: Grcar{N: 50, K: 3}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 51, K: 3}.Matrix(),
			evWant: Grcar{N: 51, K: 3}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 50, K: 10}.Matrix(),
			evWant: Grcar{N: 50, K: 10}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 51, K: 10}.Matrix(),
			evWant: Grcar{N: 51, K: 10}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 50, K: 30}.Matrix(),
			evWant: Grcar{N: 50, K: 30}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 150, K: 2}.Matrix(),
			evWant: Grcar{N: 150, K: 2}.Eigenvalues(),
		},
		{
			a:      Grcar{N: 150, K: 148}.Matrix(),
			evWant: Grcar{N: 150, K: 148}.Eigenvalues(),
		},

		{
			a:      Hanowa{N: 6, Alpha: 17}.Matrix(),
			evWant: Hanowa{N: 6, Alpha: 17}.Eigenvalues(),
		},
		{
			a:      Hanowa{N: 50, Alpha: -1}.Matrix(),
			evWant: Hanowa{N: 50, Alpha: -1}.Eigenvalues(),
		},
		{
			a:      Hanowa{N: 100, Alpha: -1}.Matrix(),
			evWant: Hanowa{N: 100, Alpha: -1}.Eigenvalues(),
		},

		{
			a:      Lesp(2).Matrix(),
			evWant: Lesp(2).Eigenvalues(),
		},
		{
			a:      Lesp(3).Matrix(),
			evWant: Lesp(3).Eigenvalues(),
		},
		{
			a:      Lesp(4).Matrix(),
			evWant: Lesp(4).Eigenvalues(),
		},
		{
			a:      Lesp(5).Matrix(),
			evWant: Lesp(5).Eigenvalues(),
		},
		{
			a:      Lesp(10).Matrix(),
			evWant: Lesp(10).Eigenvalues(),
		},
		{
			a:      Lesp(15).Matrix(),
			evWant: Lesp(15).Eigenvalues(),
		},
		{
			a:      Lesp(30).Matrix(),
			evWant: Lesp(30).Eigenvalues(),
		},
		{
			a:      Lesp(50).Matrix(),
			evWant: Lesp(50).Eigenvalues(),
			valTol: 1e-12,
		},
		{
			a:      Lesp(101).Matrix(),
			evWant: Lesp(101).Eigenvalues(),
			valTol: 1e-12,
		},
		{
			a:      Lesp(150).Matrix(),
			evWant: Lesp(150).Eigenvalues(),
			valTol: 1e-12,
		},

		{
			a:      Rutis{}.Matrix(),
			evWant: Rutis{}.Eigenvalues(),
		},

		{
			a:      Tris{N: 74, X: 1, Y: -2, Z: 1}.Matrix(),
			evWant: Tris{N: 74, X: 1, Y: -2, Z: 1}.Eigenvalues(),
		},
		{
			a:      Tris{N: 74, X: 1, Y: 2, Z: -3}.Matrix(),
			evWant: Tris{N: 74, X: 1, Y: 2, Z: -3}.Eigenvalues(),
		},
		{
			a:      Tris{N: 75, X: 1, Y: 2, Z: -3}.Matrix(),
			evWant: Tris{N: 75, X: 1, Y: 2, Z: -3}.Eigenvalues(),
		},

		{
			a:      Wilk4{}.Matrix(),
			evWant: Wilk4{}.Eigenvalues(),
		},
		{
			a:      Wilk12{}.Matrix(),
			evWant: Wilk12{}.Eigenvalues(),
			valTol: 1e-7,
		},
		{
			a:      Wilk20(0).Matrix(),
			evWant: Wilk20(0).Eigenvalues(),
		},
		{
			a:      Wilk20(1e-10).Matrix(),
			evWant: Wilk20(1e-10).Eigenvalues(),
			valTol: 1e-12,
		},

		{
			a:      Zero(1).Matrix(),
			evWant: Zero(1).Eigenvalues(),
		},
		{
			a:      Zero(10).Matrix(),
			evWant: Zero(10).Eigenvalues(),
		},
		{
			a:      Zero(50).Matrix(),
			evWant: Zero(50).Eigenvalues(),
		},
		{
			a:      Zero(100).Matrix(),
			evWant: Zero(100).Eigenvalues(),
		},
	} {
		for _, jobvl := range []lapack.LeftEVJob{lapack.LeftEVCompute, lapack.LeftEVNone} {
			for _, jobvr := range []lapack.RightEVJob{lapack.RightEVCompute, lapack.RightEVNone} {
				for _, extra := range []int{0, 11} {
					for _, wl := range []worklen{minimumWork, mediumWork, optimumWork} {
						testDgeev(t, impl, strconv.Itoa(i), test, jobvl, jobvr, extra, wl)
					}
				}
			}
		}
	}

	for _, n := range []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 20, 50, 51, 100, 101} {
		for _, jobvl := range []lapack.LeftEVJob{lapack.LeftEVCompute, lapack.LeftEVNone} {
			for _, jobvr := range []lapack.RightEVJob{lapack.RightEVCompute, lapack.RightEVNone} {
				for cas := 0; cas < 10; cas++ {
					// Create a block diagonal matrix with
					// random eigenvalues of random multiplicity.
					ev := make([]complex128, n)
					tmat := zeros(n, n, n)
					for i := 0; i < n; {
						re := rnd.NormFloat64()
						if i == n-1 || rnd.Float64() < 0.5 {
							// Real eigenvalue.
							nb := rnd.Intn(min(4, n-i)) + 1
							for k := 0; k < nb; k++ {
								tmat.Data[i*tmat.Stride+i] = re
								ev[i] = complex(re, 0)
								i++
							}
							continue
						}
						// Complex eigenvalue.
						im := rnd.NormFloat64()
						nb := rnd.Intn(min(4, (n-i)/2)) + 1
						for k := 0; k < nb; k++ {
							// 2×2 block for the complex eigenvalue.
							tmat.Data[i*tmat.Stride+i] = re
							tmat.Data[(i+1)*tmat.Stride+i+1] = re
							tmat.Data[(i+1)*tmat.Stride+i] = -im
							tmat.Data[i*tmat.Stride+i+1] = im
							ev[i] = complex(re, im)
							ev[i+1] = complex(re, -im)
							i += 2
						}
					}

					// Compute A = Q T Qᵀ where Q is an
					// orthogonal matrix.
					q := randomOrthogonal(n, rnd)
					tq := zeros(n, n, n)
					blas64.Gemm(blas.NoTrans, blas.Trans, 1, tmat, q, 0, tq)
					a := zeros(n, n, n)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, tq, 0, a)

					test := dgeevTest{
						a:      a,
						evWant: ev,
						vecTol: 1e-7,
					}
					testDgeev(t, impl, "random", test, jobvl, jobvr, 0, optimumWork)
				}
			}
		}
	}
}

func testDgeev(t *testing.T, impl Dgeever, tc string, test dgeevTest, jobvl lapack.LeftEVJob, jobvr lapack.RightEVJob, extra int, wl worklen) {
	const defaultTol = 1e-13
	valTol := test.valTol
	if valTol == 0 {
		valTol = defaultTol
	}
	vecTol := test.vecTol
	if vecTol == 0 {
		vecTol = defaultTol
	}

	a := cloneGeneral(test.a)
	n := a.Rows

	var vl blas64.General
	if jobvl == lapack.LeftEVCompute {
		vl = nanGeneral(n, n, n)
	} else {
		vl.Stride = 1
	}

	var vr blas64.General
	if jobvr == lapack.RightEVCompute {
		vr = nanGeneral(n, n, n)
	} else {
		vr.Stride = 1
	}

	wr := make([]float64, n)
	wi := make([]float64, n)

	var lwork int
	switch wl {
	case minimumWork:
		if jobvl == lapack.LeftEVCompute || jobvr == lapack.RightEVCompute {
			lwork = max(1, 4*n)
		} else {
			lwork = max(1, 3*n)
		}
	case mediumWork:
		work := make([]float64, 1)
		impl.Dgeev(jobvl, jobvr, n, a.Data, a.Stride, wr, wi, vl.Data, vl.Stride, vr.Data, vr.Stride, work, -1)
		if jobvl == lapack.LeftEVCompute || jobvr == lapack.RightEVCompute {
			lwork = (int(work[0]) + 4*n) / 2
		} else {
			lwork = (int(work[0]) + 3*n) / 2
		}
		lwork = max(1, lwork)
	case optimumWork:
		work := make([]float64, 1)
		impl.Dgeev(jobvl, jobvr, n, a.Data, a.Stride, wr, wi, vl.Data, vl.Stride, vr.Data, vr.Stride, work, -1)
		lwork = int(work[0])
	}
	work := make([]float64, lwork)

	first := impl.Dgeev(jobvl, jobvr, n, a.Data, a.Stride, wr, wi,
		vl.Data, vl.Stride, vr.Data, vr.Stride, work, len(work))

	prefix := fmt.Sprintf("Case #%v: n=%v, jobvl=%c, jobvr=%c, extra=%v, work=%v",
		tc, n, jobvl, jobvr, extra, wl)

	if !generalOutsideAllNaN(vl) {
		t.Errorf("%v: out-of-range write to VL", prefix)
	}
	if !generalOutsideAllNaN(vr) {
		t.Errorf("%v: out-of-range write to VR", prefix)
	}

	if first > 0 {
		t.Logf("%v: all eigenvalues haven't been computed, first=%v", prefix, first)
	}

	// Check that conjugate pair eigenvalues are ordered correctly.
	for i := first; i < n; {
		if wi[i] == 0 {
			i++
			continue
		}
		if wr[i] != wr[i+1] {
			t.Errorf("%v: real parts of %vth conjugate pair not equal", prefix, i)
		}
		if wi[i] < 0 || wi[i+1] >= 0 {
			t.Errorf("%v: unexpected ordering of %vth conjugate pair", prefix, i)
		}
		i += 2
	}

	// Check the computed eigenvalues against provided known eigenvalues.
	if test.evWant != nil {
		used := make([]bool, n)
		for i := first; i < n; i++ {
			evGot := complex(wr[i], wi[i])
			idx := -1
			for k, evWant := range test.evWant {
				if !used[k] && cmplx.Abs(evWant-evGot) < valTol {
					idx = k
					used[k] = true
					break
				}
			}
			if idx == -1 {
				t.Errorf("%v: unexpected eigenvalue %v", prefix, evGot)
			}
		}
	}

	if first > 0 || (jobvl == lapack.LeftEVNone && jobvr == lapack.RightEVNone) {
		// No eigenvectors have been computed.
		return
	}

	// Check that the columns of VL and VR are eigenvectors that:
	//  - correspond to the computed eigenvalues
	//  - have Euclidean norm equal to 1
	//  - have the largest component real
	bi := blas64.Implementation()
	if jobvr == lapack.RightEVCompute {
		resid := residualRightEV(test.a, vr, wr, wi)
		if resid > vecTol {
			t.Errorf("%v: unexpected right eigenvectors; residual=%v, want<=%v", prefix, resid, vecTol)
		}

		for j := 0; j < n; j++ {
			nrm := 1.0
			if wi[j] == 0 {
				nrm = bi.Dnrm2(n, vr.Data[j:], vr.Stride)
			} else if wi[j] > 0 {
				nrm = math.Hypot(bi.Dnrm2(n, vr.Data[j:], vr.Stride), bi.Dnrm2(n, vr.Data[j+1:], vr.Stride))
			}
			diff := math.Abs(nrm - 1)
			if diff > defaultTol {
				t.Errorf("%v: unexpected Euclidean norm of right eigenvector; |VR[%v]-1|=%v, want<=%v",
					prefix, j, diff, defaultTol)
			}

			if wi[j] > 0 {
				var vmax float64  // Largest component in the column
				var vrmax float64 // Largest real component in the column
				for i := 0; i < n; i++ {
					vtest := math.Hypot(vr.Data[i*vr.Stride+j], vr.Data[i*vr.Stride+j+1])
					vmax = math.Max(vmax, vtest)
					if vr.Data[i*vr.Stride+j+1] == 0 {
						vrmax = math.Max(vrmax, math.Abs(vr.Data[i*vr.Stride+j]))
					}
				}
				if vrmax/vmax < 1-defaultTol {
					t.Errorf("%v: largest component of %vth right eigenvector is not real", prefix, j)
				}
			}
		}
	}
	if jobvl == lapack.LeftEVCompute {
		resid := residualLeftEV(test.a, vl, wr, wi)
		if resid > vecTol {
			t.Errorf("%v: unexpected left eigenvectors; residual=%v, want<=%v", prefix, resid, vecTol)
		}

		for j := 0; j < n; j++ {
			nrm := 1.0
			if wi[j] == 0 {
				nrm = bi.Dnrm2(n, vl.Data[j:], vl.Stride)
			} else if wi[j] > 0 {
				nrm = math.Hypot(bi.Dnrm2(n, vl.Data[j:], vl.Stride), bi.Dnrm2(n, vl.Data[j+1:], vl.Stride))
			}
			diff := math.Abs(nrm - 1)
			if diff > defaultTol {
				t.Errorf("%v: unexpected Euclidean norm of left eigenvector; |VL[%v]-1|=%v, want<=%v",
					prefix, j, diff, defaultTol)
			}

			if wi[j] > 0 {
				var vmax float64  // Largest component in the column
				var vrmax float64 // Largest real component in the column
				for i := 0; i < n; i++ {
					vtest := math.Hypot(vl.Data[i*vl.Stride+j], vl.Data[i*vl.Stride+j+1])
					vmax = math.Max(vmax, vtest)
					if vl.Data[i*vl.Stride+j+1] == 0 {
						vrmax = math.Max(vrmax, math.Abs(vl.Data[i*vl.Stride+j]))
					}
				}
				if vrmax/vmax < 1-defaultTol {
					t.Errorf("%v: largest component of %vth left eigenvector is not real", prefix, j)
				}
			}
		}
	}
}

func dgeevTestForAntisymRandom(n int, rnd *rand.Rand) dgeevTest {
	a := NewAntisymRandom(n, rnd)
	return dgeevTest{
		a:      a.Matrix(),
		evWant: a.Eigenvalues(),
	}
}

// residualRightEV returns the residual
//
//	| A E - E W|_1 / ( |A|_1 |E|_1 )
//
// where the columns of E contain the right eigenvectors of A and W is a block diagonal matrix with
// a 1×1 block for each real eigenvalue and a 2×2 block for each complex conjugate pair.
func residualRightEV(a, e blas64.General, wr, wi []float64) float64 {
	// The implementation follows DGET22 routine from the Reference LAPACK's
	// testing suite.

	n := a.Rows
	if n == 0 {
		return 0
	}

	bi := blas64.Implementation()
	ldr := n
	r := make([]float64, n*ldr)
	var (
		wmat  [4]float64
		ipair int
	)
	for j := 0; j < n; j++ {
		if ipair == 0 && wi[j] != 0 {
			ipair = 1
		}
		switch ipair {
		case 0:
			// Real eigenvalue, multiply j-th column of E with it.
			bi.Daxpy(n, wr[j], e.Data[j:], e.Stride, r[j:], ldr)
		case 1:
			// First of complex conjugate pair of eigenvalues
			wmat[0], wmat[1] = wr[j], wi[j]
			wmat[2], wmat[3] = -wi[j], wr[j]
			bi.Dgemm(blas.NoTrans, blas.NoTrans, n, 2, 2, 1, e.Data[j:], e.Stride, wmat[:], 2, 0, r[j:], ldr)
			ipair = 2
		case 2:
			// Second of complex conjugate pair of eigenvalues
			ipair = 0
		}
	}
	bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, a.Data, a.Stride, e.Data, e.Stride, -1, r, ldr)

	const eps = dlamchE
	anorm := math.Max(dlange(lapack.MaxColumnSum, n, n, a.Data, a.Stride), safmin)
	enorm := math.Max(dlange(lapack.MaxColumnSum, n, n, e.Data, e.Stride), eps)
	errnorm := dlange(lapack.MaxColumnSum, n, n, r, ldr) / enorm
	if anorm > errnorm {
		return errnorm / anorm
	}
	if anorm < 1 {
		return math.Min(errnorm, anorm) / anorm
	}
	return math.Min(errnorm/anorm, 1)
}

// residualLeftEV returns the residual
//
//	| Aᵀ E - E Wᵀ|_1 / ( |Aᵀ|_1 |E|_1 )
//
// where the columns of E contain the left eigenvectors of A and W is a block diagonal matrix with
// a 1×1 block for each real eigenvalue and a 2×2 block for each complex conjugate pair.
func residualLeftEV(a, e blas64.General, wr, wi []float64) float64 {
	// The implementation follows DGET22 routine from the Reference LAPACK's
	// testing suite.

	n := a.Rows
	if n == 0 {
		return 0
	}

	bi := blas64.Implementation()
	ldr := n
	r := make([]float64, n*ldr)
	var (
		wmat  [4]float64
		ipair int
	)
	for j := 0; j < n; j++ {
		if ipair == 0 && wi[j] != 0 {
			ipair = 1
		}
		switch ipair {
		case 0:
			// Real eigenvalue, multiply j-th column of E with it.
			bi.Daxpy(n, wr[j], e.Data[j:], e.Stride, r[j:], ldr)
		case 1:
			// First of complex conjugate pair of eigenvalues
			wmat[0], wmat[1] = wr[j], wi[j]
			wmat[2], wmat[3] = -wi[j], wr[j]
			bi.Dgemm(blas.NoTrans, blas.Trans, n, 2, 2, 1, e.Data[j:], e.Stride, wmat[:], 2, 0, r[j:], ldr)
			ipair = 2
		case 2:
			// Second of complex conjugate pair of eigenvalues
			ipair = 0
		}
	}
	bi.Dgemm(blas.Trans, blas.NoTrans, n, n, n, 1, a.Data, a.Stride, e.Data, e.Stride, -1, r, ldr)

	const eps = dlamchE
	anorm := math.Max(dlange(lapack.MaxRowSum, n, n, a.Data, a.Stride), safmin)
	enorm := math.Max(dlange(lapack.MaxColumnSum, n, n, e.Data, e.Stride), eps)
	errnorm := dlange(lapack.MaxColumnSum, n, n, r, ldr) / enorm
	if anorm > errnorm {
		return errnorm / anorm
	}
	if anorm < 1 {
		return math.Min(errnorm, anorm) / anorm
	}
	return math.Min(errnorm/anorm, 1)
}
