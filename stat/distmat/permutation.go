// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
)

// UniformPermutation is a uniform distribution over the n!
// permutation matrices of size n×n for a given n.
type UniformPermutation struct {
	rnd     *rand.Rand
	indices []int
}

// NewUniformPermutation constructs a new permutation matrix
// generator using the given random source.
func NewUniformPermutation(src rand.Source) *UniformPermutation {
	return &UniformPermutation{rnd: rand.New(src)}
}

// PermTo sets the given matrix to be a random permutation matrix.
// It does not zero the destination's elements, so it is the responsibility
// of the caller to ensure it is correctly conditioned prior to the call.
//
// PermTo panics if dst is not square.
func (p *UniformPermutation) PermTo(dst *mat.Dense) {
	r, c := dst.Dims()
	if r != c {
		panic(mat.ErrShape)
	}
	if r == 0 {
		return
	}
	if len(p.indices) != r {
		p.indices = make([]int, r)
		for k := range p.indices {
			p.indices[k] = k
		}
	}
	p.rnd.Shuffle(r, func(i, j int) { p.indices[i], p.indices[j] = p.indices[j], p.indices[i] })
	for i, j := range p.indices {
		dst.Set(i, j, 1)
	}
}
