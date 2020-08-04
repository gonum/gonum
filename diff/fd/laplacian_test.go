// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fd

import (
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
)

func TestLaplacian(t *testing.T) {
	t.Parallel()
	for cas, test := range hessianTestCases() {
		// Modify the test cases where the formula is set.
		settings := test.settings
		if settings != nil && !settings.Formula.isZero() {
			settings.Formula = Forward2nd
		}

		n := len(test.x)
		got := Laplacian(test.h.Func, test.x, test.settings)
		hess := mat.NewSymDense(n, nil)
		test.h.Hess(hess, test.x)
		var want float64
		for i := 0; i < n; i++ {
			want += hess.At(i, i)
		}
		if !scalar.EqualWithinAbsOrRel(got, want, test.tol, test.tol) {
			t.Errorf("Cas %d: Laplacian mismatch. got %v, want %v", cas, got, want)
		}

		// Test that concurrency works.
		if settings == nil {
			settings = &Settings{}
		}
		settings.Concurrent = true
		got2 := Laplacian(test.h.Func, test.x, settings)
		if !scalar.EqualWithinAbsOrRel(got, got2, 1e-5, 1e-5) {
			t.Errorf("Cas %d: Laplacian mismatch. got %v, want %v", cas, got2, got)
		}
	}
}
