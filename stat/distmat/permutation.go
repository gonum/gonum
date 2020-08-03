package distmat

import (
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

// UniformPermutation is a distribution over permutation matrices
type UniformPermutation struct {
	src rand.Source
}

// NewUniformPermutation constructs a new permutation matrix
// generator using the given random source.
func NewUniformPermutation(src rand.Source) *UniformPermutation {
	up := UniformPermutation{src: src}
	return &up
}

// Matrix draws a random permutuation matrix of dimension n.
func (up *UniformPermutation) Matrix(n int) *mat.Dense {
	m := mat.NewDense(n, n, nil)
	up.MatrixTo(m)
	return m
}

// MatrixTo sets the given matrix to be a random permutation matrix
// panics if m is not square
// Note this does not zero the matrix -- before calling this ensure it
// is full of 0s.
func (up *UniformPermutation) MatrixTo(m *mat.Dense) {
	r, c := m.Dims()
	if r != c {
		panic(mat.ErrShape)
	}
	if r == 0 {
		return
	}
	iList := make([]int, r)
	for i := range iList {
		iList[i] = i
	}
	rnd := rand.New(up.src)
	rnd.Shuffle(r, func(i, j int) { iList[i], iList[j] = iList[j], iList[i] })
	for i, j := range iList {
		m.Set(i, j, 1.0)
	}
}
