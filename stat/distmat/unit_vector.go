// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

// UnitVector is a uniform distribtion over the surface of a sphere.
type UnitVector struct {
	norm distuv.Normal
}

// NewUnitVector constructs a unit vector generator.
func NewUnitVector(src rand.Source) *UnitVector {
	return &UnitVector{norm: distuv.Normal{Mu: 0, Sigma: 1, Src: src}}
}

// UnitVecTo sets the given n-dimension vector to be a random
// unit-length n-dimension vector.
//
// This uses the algorithm of Mueller from:
// https://dl.acm.org/doi/10.1145/377939.377946
// and summarized on an non-paywalled page at:
// https://mathworld.wolfram.com/HyperspherePointPicking.html
//
// UnitVecTo panics if dst has 0 length.
func (u *UnitVector) UnitVecTo(dst *mat.VecDense) {
	r := dst.Len()
	if r == 0 {
		panic(mat.ErrZeroLength)
	}
	for i := 0; i < r; i++ {
		dst.SetVec(i, u.norm.Rand())
	}
	l := mat.Norm(dst, 2)
	dst.ScaleVec(1/l, dst)
}
