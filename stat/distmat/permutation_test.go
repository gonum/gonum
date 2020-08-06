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
		if !confirmEachRowAndColumnOneNZ(m) {
			t.Error("didnt get back a permutation matrix")
		}
	}

}

func confirmEachRowAndColumnOneNZ(m mat.Matrix) bool {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		rowNZ := 0
		colNZ := 0
		rowN1 := 0
		colN1 := 0
		for j := 0; j < c; j++ {
			if floats.EqualWithinAbs(m.At(i, j), 0, 1e-12) {
				rowNZ++
			} else if floats.EqualWithinAbs(m.At(i, j), 1, 1e-12) {
				rowN1++
			} else {
				return false
			}
			if floats.EqualWithinAbs(m.At(j, i), 0, 1e-12) {
				colNZ++
			} else if floats.EqualWithinAbs(m.At(j, i), 1, 1e-12) {
				colN1++
			} else {
				return false
			}
		}
		if rowNZ != (r-1) || colNZ != (r-1) {
			return false
		}
		if rowN1 != 1 || colN1 != 1 {
			return false
		}
	}
	return true
}
