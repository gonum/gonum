// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestUniformPermutation(t *testing.T) {
	up := NewUniformPermutation(rand.NewSource(1))
	for _, n := range []int{10, 32, 64, 100} {
		m := mat.NewDense(n, n, nil)
		if m == nil {
			t.Error("Matrix failed")
		}
		up.PermTo(m)
		r, c := m.Dims()
		if r != n || c != n {
			t.Error("got back matrix of wrong size")
		}
		confirmPermMatrix(t, m)
	}

}

func confirmPermMatrix(t *testing.T, m mat.Matrix) {
	t.Helper()
	r, c := m.Dims()
	if r != c {
		t.Error("matrix not square")
	}
	rowSums := make([]float64, r)
	colSums := make([]float64, c)
	for i := 0; i < r; i++ {
		for j, v := range mat.Row(nil, i, m) {
			switch v {
			case 0, 1:
				rowSums[i] += v
				colSums[j] += v
			default:
				t.Errorf("unexpected value %f at position %d %d", v, i, j)
			}
		}
	}
	if floats.Max(rowSums) != 1 || floats.Min(rowSums) != 1 {
		t.Error("found non-1 row sum")
	}
	if floats.Max(colSums) != 1 || floats.Min(colSums) != 1 {
		t.Error("found non-1 row sum")
	}
}
