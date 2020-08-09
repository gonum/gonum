// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
)

func TestUniformPermutation(t *testing.T) {
	up := NewUniformPermutation(rand.NewSource(1))
	for _, n := range []int{10, 32, 64, 100} {
		m := mat.NewDense(n, n, nil)
		up.PermTo(m)
		// Test that each row and column satisfies the permutation matrix
		// invariant that all rows and columns have a single unit element
		// and the remaining elements are zero.
		for i := 0; i < n; i++ {
			checkHasSingleUnitElement(t, "row", i, mat.Row(nil, i, m))
			checkHasSingleUnitElement(t, "col", i, mat.Col(nil, i, m))
		}
	}
}

func checkHasSingleUnitElement(t *testing.T, dir string, n int, v []float64) {
	t.Helper()
	var sum float64
	for i, x := range v {
		switch x {
		case 0, 1:
			sum += x
		default:
			t.Errorf("unexpected value in %s %d position %d: %v", dir, n, i, v)
		}
	}
	if sum != 1 {
		t.Errorf("%s %d is not a valid vector: %v", dir, n, v)
	}
}
