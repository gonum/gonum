// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

func TestNewTriBand(t *testing.T) {
	t.Parallel()
	for cas, test := range []struct {
		data  []float64
		n, k  int
		kind  TriKind
		mat   *TriBandDense
		dense *Dense
	}{
		{
			data: []float64{1, 2, 3},
			n:    3, k: 0,
			kind: Upper,
			mat: &TriBandDense{
				mat: blas64.TriangularBand{
					Diag: blas.NonUnit,
					Uplo: blas.Upper,
					N:    3, K: 0,
					Data:   []float64{1, 2, 3},
					Stride: 1,
				},
			},
			dense: NewDense(3, 3, []float64{
				1, 0, 0,
				0, 2, 0,
				0, 0, 3,
			}),
		},
		{
			data: []float64{
				1, 2,
				3, 4,
				5, 6,
				7, 8,
				9, 10,
				11, -1,
			},
			n: 6, k: 1,
			kind: Upper,
			mat: &TriBandDense{
				mat: blas64.TriangularBand{
					Diag: blas.NonUnit,
					Uplo: blas.Upper,
					N:    6, K: 1,
					Data: []float64{
						1, 2,
						3, 4,
						5, 6,
						7, 8,
						9, 10,
						11, -1,
					},
					Stride: 2,
				},
			},
			dense: NewDense(6, 6, []float64{
				1, 2, 0, 0, 0, 0,
				0, 3, 4, 0, 0, 0,
				0, 0, 5, 6, 0, 0,
				0, 0, 0, 7, 8, 0,
				0, 0, 0, 0, 9, 10,
				0, 0, 0, 0, 0, 11,
			}),
		},
		{
			data: []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10, 11, 12,
				13, 14, -1,
				15, -1, -1,
			},
			n: 6, k: 2,
			kind: Upper,
			mat: &TriBandDense{
				mat: blas64.TriangularBand{
					Diag: blas.NonUnit,
					Uplo: blas.Upper,
					N:    6, K: 2,
					Data: []float64{
						1, 2, 3,
						4, 5, 6,
						7, 8, 9,
						10, 11, 12,
						13, 14, -1,
						15, -1, -1,
					},
					Stride: 3,
				},
			},
			dense: NewDense(6, 6, []float64{
				1, 2, 3, 0, 0, 0,
				0, 4, 5, 6, 0, 0,
				0, 0, 7, 8, 9, 0,
				0, 0, 0, 10, 11, 12,
				0, 0, 0, 0, 13, 14,
				0, 0, 0, 0, 0, 15,
			}),
		},
		{
			data: []float64{
				-1, 1,
				2, 3,
				4, 5,
				6, 7,
				8, 9,
				10, 11,
			},
			n: 6, k: 1,
			kind: Lower,
			mat: &TriBandDense{
				mat: blas64.TriangularBand{
					Diag: blas.NonUnit,
					Uplo: blas.Lower,
					N:    6, K: 1,
					Data: []float64{
						-1, 1,
						2, 3,
						4, 5,
						6, 7,
						8, 9,
						10, 11,
					},
					Stride: 2,
				},
			},
			dense: NewDense(6, 6, []float64{
				1, 0, 0, 0, 0, 0,
				2, 3, 0, 0, 0, 0,
				0, 4, 5, 0, 0, 0,
				0, 0, 6, 7, 0, 0,
				0, 0, 0, 8, 9, 0,
				0, 0, 0, 0, 10, 11,
			}),
		},
		{
			data: []float64{
				-1, -1, 1,
				-1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10, 11, 12,
				13, 14, 15,
			},
			n: 6, k: 2,
			kind: Lower,
			mat: &TriBandDense{
				mat: blas64.TriangularBand{
					Diag: blas.NonUnit,
					Uplo: blas.Lower,
					N:    6, K: 2,
					Data: []float64{
						-1, -1, 1,
						-1, 2, 3,
						4, 5, 6,
						7, 8, 9,
						10, 11, 12,
						13, 14, 15,
					},
					Stride: 3,
				},
			},
			dense: NewDense(6, 6, []float64{
				1, 0, 0, 0, 0, 0,
				2, 3, 0, 0, 0, 0,
				4, 5, 6, 0, 0, 0,
				0, 7, 8, 9, 0, 0,
				0, 0, 10, 11, 12, 0,
				0, 0, 0, 13, 14, 15,
			}),
		},
	} {
		triBand := NewTriBandDense(test.n, test.k, test.kind, test.data)
		r, c := triBand.Dims()
		n, k, kind := triBand.TriBand()
		if n != test.n {
			t.Errorf("unexpected triband size for test %d: got: %d want: %d", cas, n, test.n)
		}
		if k != test.k {
			t.Errorf("unexpected triband bandwidth for test %d: got: %d want: %d", cas, k, test.k)
		}
		if kind != test.kind {
			t.Errorf("unexpected triband bandwidth for test %v: got: %v want: %v", cas, kind, test.kind)
		}
		if r != n {
			t.Errorf("unexpected number of rows for test %d: got: %d want: %d", cas, r, n)
		}
		if c != n {
			t.Errorf("unexpected number of cols for test %d: got: %d want: %d", cas, c, n)
		}
		if !reflect.DeepEqual(triBand, test.mat) {
			t.Errorf("unexpected value via reflect for test %d: got: %v want: %v", cas, triBand, test.mat)
		}
		if !Equal(triBand, test.mat) {
			t.Errorf("unexpected value via mat.Equal for test %d: got: %v want: %v", cas, triBand, test.mat)
		}
		if !Equal(triBand, test.dense) {
			t.Errorf("unexpected value via mat.Equal(band, dense) for test %d:\ngot:\n% v\nwant:\n% v", cas, Formatted(triBand), Formatted(test.dense))
		}
	}
}

func TestTriBandAtSetUpper(t *testing.T) {
	t.Parallel()
	for _, kind := range []TriKind{Upper, Lower} {
		var band *TriBandDense
		var data []float64
		if kind {
			// 1  2  3  0  0  0
			// 0  4  5  6  0  0
			// 0  0  7  8  9  0
			// 0  0  0 10 11 12
			// 0  0  0  0 13 14
			// 0  0  0  0  0 15
			data = []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10, 11, 12,
				13, 14, -1,
				15, -1, -1,
			}
			band = NewTriBandDense(6, 2, kind, data)
		} else {
			// 1  0  0  0  0  0
			// 2  3  0  0  0  0
			// 4  5  6  0  0  0
			// 0  7  8  9  0  0
			// 0  0 10 11 12  0
			// 0  0  0 13 14 15
			data = []float64{
				-1, -1, 1,
				-1, 2, 3,
				4, 5, 6,
				7, 8, 9,
				10, 11, 12,
				13, 14, 15,
			}
			band = NewTriBandDense(6, 2, kind, data)
		}

		rows, cols := band.Dims()

		// Check At out of bounds.
		for _, row := range []int{-1, rows, rows + 1} {
			panicked, message := panics(func() { band.At(row, 0) })
			if !panicked || message != ErrRowAccess.Error() {
				t.Errorf("expected panic for invalid row access N=%d r=%d", rows, row)
			}
		}
		for _, col := range []int{-1, cols, cols + 1} {
			panicked, message := panics(func() { band.At(0, col) })
			if !panicked || message != ErrColAccess.Error() {
				t.Errorf("expected panic for invalid column access N=%d c=%d", cols, col)
			}
		}

		// Check Set out of bounds
		// First, check outside the matrix bounds.
		for _, row := range []int{-1, rows, rows + 1} {
			panicked, message := panics(func() { band.SetTriBand(row, 0, 1.2) })
			if !panicked || message != ErrRowAccess.Error() {
				t.Errorf("expected panic for invalid row access N=%d r=%d", rows, row)
			}
		}
		for _, col := range []int{-1, cols, cols + 1} {
			panicked, message := panics(func() { band.SetTriBand(0, col, 1.2) })
			if !panicked || message != ErrColAccess.Error() {
				t.Errorf("expected panic for invalid column access N=%d c=%d", cols, col)
			}
		}
		// Next, check outside the Triangular bounds.
		for _, s := range []struct{ r, c int }{
			{3, 2},
		} {
			if kind == Lower {
				s.r, s.c = s.c, s.r
			}
			panicked, message := panics(func() { band.SetTriBand(s.r, s.c, 1.2) })
			if !panicked || message != ErrTriangleSet.Error() {
				t.Errorf("expected panic for invalid triangular access N=%d, r=%d c=%d", cols, s.r, s.c)
			}
		}
		// Finally, check inside the triangle, but outside the band.
		for _, s := range []struct{ r, c int }{
			{1, 5},
		} {
			if kind == Lower {
				s.r, s.c = s.c, s.r
			}
			panicked, message := panics(func() { band.SetTriBand(s.r, s.c, 1.2) })
			if !panicked || message != ErrBandSet.Error() {
				t.Errorf("expected panic for invalid triangular access N=%d, r=%d c=%d", cols, s.r, s.c)
			}
		}

		// Test that At and Set work correctly.
		offset := 100.0
		dataCopy := make([]float64, len(data))
		copy(dataCopy, data)
		for i := 0; i < rows; i++ {
			for j := 0; j < rows; j++ {
				v := band.At(i, j)
				if v != 0 {
					band.SetTriBand(i, j, v+offset)
				}
			}
		}
		for i, v := range dataCopy {
			if v == -1 {
				if data[i] != -1 {
					t.Errorf("Set changed unexpected entry. Want %v, got %v", -1, data[i])
				}
			} else {
				if v != data[i]-offset {
					t.Errorf("Set incorrectly changed for %v. got %v, want %v", v, data[i], v+offset)
				}
			}
		}
	}
}

func TestTriBandDenseZero(t *testing.T) {
	t.Parallel()
	// Elements that equal 1 should be set to zero, elements that equal -1
	// should remain unchanged.
	for _, test := range []*TriBandDense{
		{
			mat: blas64.TriangularBand{
				Uplo:   blas.Upper,
				N:      6,
				K:      2,
				Stride: 5,
				Data: []float64{
					1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
					1, 1, -1, -1, -1,
					1, -1, -1, -1, -1,
				},
			},
		},
		{
			mat: blas64.TriangularBand{
				Uplo:   blas.Lower,
				N:      6,
				K:      2,
				Stride: 5,
				Data: []float64{
					-1, -1, 1, -1, -1,
					-1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
					1, 1, 1, -1, -1,
				},
			},
		},
	} {
		dataCopy := make([]float64, len(test.mat.Data))
		copy(dataCopy, test.mat.Data)
		test.Zero()
		for i, v := range test.mat.Data {
			if dataCopy[i] != -1 && v != 0 {
				t.Errorf("Matrix not zeroed in bounds")
			}
			if dataCopy[i] == -1 && v != -1 {
				t.Errorf("Matrix zeroed out of bounds")
			}
		}
	}
}

func TestTriBandDiagView(t *testing.T) {
	t.Parallel()
	for cas, test := range []*TriBandDense{
		NewTriBandDense(1, 0, Upper, []float64{1}),
		NewTriBandDense(4, 0, Upper, []float64{1, 2, 3, 4}),
		NewTriBandDense(6, 2, Upper, []float64{
			1, 2, 3,
			4, 5, 6,
			7, 8, 9,
			10, 11, 12,
			13, 14, -1,
			15, -1, -1,
		}),
		NewTriBandDense(1, 0, Lower, []float64{1}),
		NewTriBandDense(4, 0, Lower, []float64{1, 2, 3, 4}),
		NewTriBandDense(6, 2, Lower, []float64{
			-1, -1, 1,
			-1, 2, 3,
			4, 5, 6,
			7, 8, 9,
			10, 11, 12,
			13, 14, 15,
		}),
	} {
		testDiagView(t, cas, test)
	}
}

func TestTriBandDenseSolveTo(t *testing.T) {
	t.Parallel()

	const tol = 1e-15

	for tc, test := range []struct {
		a *TriBandDense
		b *Dense
	}{
		{
			a: NewTriBandDense(5, 2, Upper, []float64{
				-0.34, -0.49, -0.51,
				-0.25, -0.5, 1.03,
				-1.1, 0.3, -0.82,
				1.69, 0.69, -2.22,
				-0.62, 1.22, -0.85,
			}),
			b: NewDense(5, 2, []float64{
				0.44, 1.34,
				0.07, -1.45,
				-0.32, -0.88,
				-0.09, -0.15,
				-1.17, -0.19,
			}),
		},
		{
			a: NewTriBandDense(5, 2, Lower, []float64{
				0, 0, -0.34,
				0, -0.49, -0.25,
				-0.51, -0.5, -1.1,
				1.03, 0.3, 1.69,
				-0.82, 0.69, -0.62,
			}),
			b: NewDense(5, 2, []float64{
				0.44, 1.34,
				0.07, -1.45,
				-0.32, -0.88,
				-0.09, -0.15,
				-1.17, -0.19,
			}),
		},
	} {
		a := test.a
		for _, trans := range []bool{false, true} {
			for _, dstSameAsB := range []bool{false, true} {
				name := fmt.Sprintf("Case %d,trans=%v,dstSameAsB=%v", tc, trans, dstSameAsB)

				n, nrhs := test.b.Dims()
				var dst Dense
				var err error
				if dstSameAsB {
					dst = *NewDense(n, nrhs, nil)
					dst.Copy(test.b)
					err = a.SolveTo(&dst, trans, &dst)
				} else {
					tmp := NewDense(n, nrhs, nil)
					tmp.Copy(test.b)
					err = a.SolveTo(&dst, trans, asBasicMatrix(tmp))
				}

				if err != nil {
					t.Fatalf("%v: unexpected error from SolveTo", name)
				}

				var resid Dense
				if trans {
					resid.Mul(a.T(), &dst)
				} else {
					resid.Mul(a, &dst)
				}
				resid.Sub(&resid, test.b)
				diff := Norm(&resid, 1)
				if diff > tol {
					t.Errorf("%v: unexpected result; diff=%v,want<=%v", name, diff, tol)
				}
			}
		}
	}
}

func TestTriBandDenseSolveVecTo(t *testing.T) {
	t.Parallel()

	const tol = 1e-15

	for tc, test := range []struct {
		a *TriBandDense
		b *VecDense
	}{
		{
			a: NewTriBandDense(5, 2, Upper, []float64{
				-0.34, -0.49, -0.51,
				-0.25, -0.5, 1.03,
				-1.1, 0.3, -0.82,
				1.69, 0.69, -2.22,
				-0.62, 1.22, -0.85,
			}),
			b: NewVecDense(5, []float64{
				0.44,
				0.07,
				-0.32,
				-0.09,
				-1.17,
			}),
		},
		{
			a: NewTriBandDense(5, 2, Lower, []float64{
				0, 0, -0.34,
				0, -0.49, -0.25,
				-0.51, -0.5, -1.1,
				1.03, 0.3, 1.69,
				-0.82, 0.69, -0.62,
			}),
			b: NewVecDense(5, []float64{
				0.44,
				0.07,
				-0.32,
				-0.09,
				-1.17,
			}),
		},
	} {
		a := test.a
		for _, trans := range []bool{false, true} {
			for _, dstSameAsB := range []bool{false, true} {
				name := fmt.Sprintf("Case %d,trans=%v,dstSameAsB=%v", tc, trans, dstSameAsB)

				n, _ := test.b.Dims()
				var dst VecDense
				var err error
				if dstSameAsB {
					dst = *NewVecDense(n, nil)
					dst.CopyVec(test.b)
					err = a.SolveVecTo(&dst, trans, &dst)
				} else {
					tmp := NewVecDense(n, nil)
					tmp.CopyVec(test.b)
					err = a.SolveVecTo(&dst, trans, asBasicVector(tmp))
				}

				if err != nil {
					t.Fatalf("%v: unexpected error from SolveVecTo", name)
				}

				var resid VecDense
				if trans {
					resid.MulVec(a.T(), &dst)
				} else {
					resid.MulVec(a, &dst)
				}
				resid.SubVec(&resid, test.b)
				diff := Norm(&resid, 1)
				if diff > tol {
					t.Errorf("%v: unexpected result; diff=%v,want<=%v", name, diff, tol)
				}
			}
		}
	}
}
