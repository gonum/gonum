// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
)

func TestUnitVector(t *testing.T) {
	u := NewUnitVector(rand.NewPCG(1, 1))
	for _, d := range []int{10, 32, 64, 100} {
		v := mat.NewVecDense(d, nil)
		u.UnitVecTo(v)
		l := mat.Norm(v, 2)
		if !scalar.EqualWithinAbs(l, 1.0, 1e-12) {
			t.Errorf("expected length 1 but got %f", l)
		}
	}
}

func TestUnitVectorStats(t *testing.T) {
	n := 1e7
	u := NewUnitVector(rand.NewPCG(1, 1))
	for _, d := range []int{1, 2, 3} {
		v := mat.NewVecDense(d, nil)
		tot := mat.NewVecDense(d, nil)
		for i := 0; i < int(n); i++ {
			u.UnitVecTo(v)
			tot.AddVec(tot, v)
		}
		tot.ScaleVec(1/n, tot)
		// Each dimension should average out to 0.
		for i := 0; i < d; i++ {
			if !scalar.EqualWithinAbs(tot.AtVec(i), 0.0, 1e-3) {
				t.Errorf("expected average entry 0 but got %f", tot.AtVec(i))
			}
		}
		l := mat.Norm(tot, 2)
		// And the length should be 0.
		if !scalar.EqualWithinAbs(l, 0.0, 1e-3) {
			t.Errorf("expected length 0 but got %f", l)
		}
	}
}
