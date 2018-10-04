// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"reflect"
	"testing"
)

func TestNewDiagonal(t *testing.T) {
	for i, test := range []struct {
		data  []float64
		n     int
		mat   *DiagDense
		dense *Dense
	}{
		{
			data: []float64{1, 2, 3, 4, 5, 6},
			n:    6,
			mat: &DiagDense{
				data: []float64{1, 2, 3, 4, 5, 6},
			},
			dense: NewDense(6, 6, []float64{
				1, 0, 0, 0, 0, 0,
				0, 2, 0, 0, 0, 0,
				0, 0, 3, 0, 0, 0,
				0, 0, 0, 4, 0, 0,
				0, 0, 0, 0, 5, 0,
				0, 0, 0, 0, 0, 6,
			}),
		},
	} {
		band := NewDiagonal(test.n, test.data)
		rows, cols := band.Dims()

		if rows != test.n {
			t.Errorf("unexpected number of rows for test %d: got: %d want: %d", i, rows, test.n)
		}
		if cols != test.n {
			t.Errorf("unexpected number of cols for test %d: got: %d want: %d", i, cols, test.n)
		}
		if !reflect.DeepEqual(band, test.mat) {
			t.Errorf("unexpected value via reflect for test %d: got: %v want: %v", i, band, test.mat)
		}
		if !Equal(band, test.mat) {
			t.Errorf("unexpected value via mat.Equal for test %d: got: %v want: %v", i, band, test.mat)
		}
		if !Equal(band, test.dense) {
			t.Errorf("unexpected value via mat.Equal(band, dense) for test %d:\ngot:\n% v\nwant:\n% v", i, Formatted(band), Formatted(test.dense))
		}
	}
}

func TestDiagonalAtSet(t *testing.T) {
	for _, n := range []int{1, 3, 8} {
		for _, nilstart := range []bool{true, false} {
			var diag *DiagDense
			if nilstart {
				diag = NewDiagonal(n, nil)
			} else {
				data := make([]float64, n)
				diag = NewDiagonal(n, data)
				// Test the data is used.
				for i := range data {
					data[i] = -float64(i) - 1
					v := diag.At(i, i)
					if v != data[i] {
						t.Errorf("Diag shadow mismatch. Got %v, want %v", v, data[i])
					}
				}
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if i != j {
						if diag.At(i, j) != 0 {
							t.Errorf("Diag returned non-zero off diagonal element at %d, %d", i, j)
						}
					}
					v := float64(i) + 1
					diag.SetDiag(i, v)
					v2 := diag.At(i, i)
					if v2 != v {
						t.Errorf("Diag at/set mismatch. Got %v, want %v", v, v2)
					}
				}
			}
		}
	}
}
