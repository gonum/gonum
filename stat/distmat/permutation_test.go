// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"math/cmplx"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/floats"

	"gonum.org/v1/gonum/mat"
)

func TestUniformPermutation(t *testing.T) {
	up := NewUniformPermutation(rand.NewSource(1))
	for _, n := range []int{10, 32, 64, 100} {
		m := up.Matrix(n)
		if m == nil {
			t.Error("Matrix failed")
		}
		if !confirmEigenvaluesAreRootsOfUnity(m) {
			t.Error("eigenvalue not a root of unity")
		}
		if !confirmEachColumnOneNZ(m) {
			t.Error("didnt get back a permutation matrix")
		}
	}

}

func confirmEigenvaluesAreRootsOfUnity(m mat.Matrix) bool {
	n, _ := m.Dims()
	var e mat.Eigen
	e.Factorize(m, mat.EigenLeft)
	values := make([]complex128, n)
	e.Values(values)
	for _, v := range values {
		a := cmplx.Abs(v)
		if !floats.EqualWithinAbs(a, 1.0, 1e-12) {
			return false
		}
	}
	return true
}

func confirmEachColumnOneNZ(m mat.Matrix) bool {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		count := 0
		for j := 0; j < c; j++ {
			if !floats.EqualWithinAbs(m.At(i, j), 0.0, 1e-12) {
				count++
			}
		}
		if count != 1 {
			return false
		}
	}
	return true
}
