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

func TestUnitVector(t *testing.T) {
	u := NewUnitVector(rand.NewSource(1))
	for _, n := range []int{10, 32, 64, 100} {
		v := mat.NewVecDense(n, nil)
		u.UnitVecTo(v)
		l := mat.Norm(v, 2)
		if !floats.EqualWithinAbs(l, 1.0, 1e-12) {
			t.Errorf("expected length 1 but got %f", l)
		}
	}
}
