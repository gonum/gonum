// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import "github.com/gonum/blas/blas64"

const (
	// regionOverlap is the panic string used for the general case
	// of a matrix region overlap between a source and destination.
	regionOverlap = "mat64: bad region: overlap"

	// regionIdentity is the panic string used for the specific
	// case of complete agreement between a source and a destination.
	regionIdentity = "mat64: bad region: identical"

	// mismatchedStrides is the panic string used for overlapping
	// data slices with differing strides.
	mismatchedStrides = "mat64: bad region: different strides"
)

// checkOverlap returns false if the receiver does not overlap data elements
// referenced by the parameter and panics otherwise.
//
// checkOverlap methods return a boolean to allow the check call to be added to a
// boolean expression, making use of short-circuit operators.

func (m *Dense) checkOverlap(a blas64.General) bool {
	mat := m.RawMatrix()
	if cap(mat.Data) == 0 || cap(a.Data) == 0 {
		return false
	}

	off := offset(mat.Data[:1], a.Data[:1])

	if off == 0 {
		// At least one element overlaps.
		if mat.Cols == a.Cols && mat.Rows == a.Rows && mat.Stride == a.Stride {
			panic(regionIdentity)
		}
		panic(regionOverlap)
	}

	if off > 0 && len(mat.Data) <= off {
		// We know m is completely before a.
		return false
	}
	if off < 0 && len(a.Data) <= -off {
		// We know m is completely after a.
		return false
	}

	if mat.Stride != a.Stride {
		// Too hard, so assume the worst.
		panic(mismatchedStrides)
	}

	if off < 0 {
		off = -off
		mat.Cols, a.Cols = a.Cols, mat.Cols
	}
	if rectanglesOverlap(off, mat.Cols, a.Cols, mat.Stride) {
		panic(regionOverlap)
	}
	return false
}

// BUG(kortschak): Overlap detection for symmetric and triangular matrices is not
// precise; currently overlap is detected if the bounding rectangles overlap rather
// than exact overlap between visible elements.

func (s *SymDense) checkOverlap(a blas64.Symmetric) bool {
	mat := s.RawSymmetric()
	if cap(mat.Data) == 0 || cap(a.Data) == 0 {
		return false
	}

	off := offset(mat.Data[:1], a.Data[:1])

	if off == 0 {
		// At least one element overlaps.
		if mat.N == a.N && mat.Stride == a.Stride {
			panic(regionIdentity)
		}
		panic(regionOverlap)
	}

	if off > 0 && len(mat.Data) <= off {
		// We know s is completely before a.
		return false
	}
	if off < 0 && len(a.Data) <= -off {
		// We know s is completely after a.
		return false
	}

	if mat.Stride != a.Stride {
		// Too hard, so assume the worst.
		panic(mismatchedStrides)
	}

	// TODO(kortschak) Make this analysis more precise.
	if off > 0 {
		off = -off
		mat.N, a.N = a.N, mat.N
	}
	if rectanglesOverlap(off, mat.N, a.N, mat.Stride) {
		panic(regionOverlap)
	}
	return false
}

func (t *TriDense) checkOverlap(a blas64.Triangular) bool {
	mat := t.RawTriangular()
	if cap(mat.Data) == 0 || cap(a.Data) == 0 {
		return false
	}

	off := offset(mat.Data[:1], a.Data[:1])

	if off == 0 {
		// At least one element overlaps.
		if mat.N == a.N && mat.Stride == a.Stride {
			panic(regionIdentity)
		}
		panic(regionOverlap)
	}

	if off > 0 && len(mat.Data) <= off {
		// We know t is completely before a.
		return false
	}
	if off < 0 && len(a.Data) <= -off {
		// We know t is completely after a.
		return false
	}

	if mat.Stride != a.Stride {
		// Too hard, so assume the worst.
		panic(mismatchedStrides)
	}

	// TODO(kortschak) Make this analysis more precise.
	if off > 0 {
		off = -off
		mat.N, a.N = a.N, mat.N
	}
	if rectanglesOverlap(off, mat.N, a.N, mat.Stride) {
		panic(regionOverlap)
	}
	return false
}

func (v *Vector) checkOverlap(a blas64.Vector) bool {
	mat := v.mat
	if cap(mat.Data) == 0 || cap(a.Data) == 0 {
		return false
	}

	off := offset(mat.Data[:1], a.Data[:1])

	if off == 0 {
		// At least one element overlaps.
		if mat.Inc == a.Inc && len(mat.Data) == len(a.Data) {
			panic(regionIdentity)
		}
		panic(regionOverlap)
	}

	if off > 0 && len(mat.Data) <= off {
		// We know v is completely before a.
		return false
	}
	if off < 0 && len(a.Data) <= -off {
		// We know v is completely after a.
		return false
	}

	if mat.Inc != a.Inc {
		// Too hard, so assume the worst.
		panic(mismatchedStrides)
	}

	if mat.Inc == 1 || off&mat.Inc == 0 {
		panic(regionOverlap)
	}
	return false
}

// rectanglesOverlap returns whether the strided rectangles a and b overlap
// when b is offset by off elements after a but has at least one element before
// the end of a. a and b have aCols and bCols respectively.
func rectanglesOverlap(off, aCols, bCols, stride int) bool {
	if stride == 1 {
		return true
	}
	aTo := aCols
	bFrom := off % stride
	bTo := (bFrom + bCols) % stride
	if bFrom < bTo {
		return aTo > bFrom
	}
	return bTo > 0
}
