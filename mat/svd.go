// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
	"gonum.org/v1/gonum/lapack/lapack64"
)

// SVDKind specifies the treatment of singular vectors during an SVD
// factorization.
type SVDKind int

const (
	// SVDThin computes the thin singular vectors. For the factorization
	//  A = U~ * Σ * V~^T
	// U~ is of size m×min(m,n) and V~ is of size n×min(m,n).
	SVDThin SVDKind = iota
	// SVDFull computes the full singular vectors.  For the factorization
	//  A = U * Σ * V^T
	// U is of size m×m, and V is of size n×n.
	SVDFull
	// SVDNone specifies that the singular vectors should not be computed during
	// the decomposition.
	SVDNone
)

// SVD is a type for creating and using the Singular Value Decomposition (SVD)
// of a matrix.
type SVD struct {
	U SVDKind
	V SVDKind

	computed bool

	s  []float64
	u  blas64.General
	vt blas64.General
}

// Factorize computes the singular value decomposition (SVD) of the input matrix A.
// The singular values of A are computed in all cases, and by default the thin
// singular vectors are computed. The singular vector computation behavior can
// be controlled with the U and V fields.
//
// The full singular value decomposition is a factorization of an m×n matrix A
// of the form
//  A = U * Σ * V^T
// where Σ is an m×n diagonal matrix, U is an m×m orthogonal matrix, and V is an
// n×n orthogonal matrix. The diagonal elements of Σ are the singular values of A.
// The first min(m,n) columns of U and V are, respectively, the left and right
// singular vectors of A. This factorization can be obtained by setting U and V
// to SVDFull.
//
// It is typically preferred to compute the thin SVD factorization, and is the
// default behavior. The thin representation saves a significant amount of memory
// if m >> n or m << n. If the singular vectors are not needed, time and memory
// can be saved by setting the respective field to SVDNone.
//
// Factorize returns whether the decomposition succeeded. If the decomposition
// failed, routines that require a successful factorization will panic.
func (svd *SVD) Factorize(a Matrix) (ok bool) {
	svd.computed = false
	m, n := a.Dims()
	var jobU, jobVT lapack.SVDJob
	switch svd.U {
	default:
		panic("svd: bad U kind")
	case SVDFull:
		// TODO(btracey): This code should be modified to have the smaller
		// matrix written in-place into aCopy when the lapack/native/dgesvd
		// implementation is complete.
		// This modification should be reflected in the V switch.
		svd.u = blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   use(svd.u.Data, m*m),
		}
		jobU = lapack.SVDAll
	case SVDThin:
		// TODO(btracey): This code should be modified to have the larger
		// matrix written in-place into aCopy when the lapack/native/dgesvd
		// implementation is complete.
		// This modification should be reflected in the V switch.
		svd.u = blas64.General{
			Rows:   m,
			Cols:   min(m, n),
			Stride: min(m, n),
			Data:   use(svd.u.Data, m*min(m, n)),
		}
		jobU = lapack.SVDStore
	case SVDNone:
		svd.u.Stride = 1
		jobU = lapack.SVDNone
	}
	switch svd.V {
	default:
		panic("svd: bad V kind")
	case SVDFull:
		svd.vt = blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   use(svd.vt.Data, n*n),
		}
		jobVT = lapack.SVDAll
	case SVDThin:
		svd.vt = blas64.General{
			Rows:   min(m, n),
			Cols:   n,
			Stride: n,
			Data:   use(svd.vt.Data, min(m, n)*n),
		}
		jobVT = lapack.SVDStore
	case SVDNone:
		svd.vt.Stride = 1
		jobVT = lapack.SVDNone
	}

	// A is destroyed on call, so copy the matrix.
	aCopy := DenseCopyOf(a)
	svd.s = use(svd.s, min(m, n))

	work := []float64{0}
	lapack64.Gesvd(jobU, jobVT, aCopy.mat, svd.u, svd.vt, svd.s, work, -1)
	work = getFloats(int(work[0]), false)
	ok = lapack64.Gesvd(jobU, jobVT, aCopy.mat, svd.u, svd.vt, svd.s, work, len(work))
	putFloats(work)
	svd.computed = true
	if !ok {
		svd.computed = false
	}
	return ok
}

// Cond returns the 2-norm condition number for the factorized matrix. Cond will
// panic if the receiver does not contain a successful factorization.
func (svd *SVD) Cond() float64 {
	if !svd.computed {
		panic("svd: no decomposition computed")
	}
	return svd.s[0] / svd.s[len(svd.s)-1]
}

// Values returns the singular values of the factorized matrix in descending order.
//
// If the input slice is non-nil, the values will be stored in-place into
// the slice. In this case, the slice must have length min(m,n), and Values will
// panic with ErrSliceLengthMismatch otherwise. If the input slice is nil, a new
// slice of the appropriate length will be allocated and returned.
//
// Values will panic if the receiver does not contain a successful factorization.
func (svd *SVD) Values(s []float64) []float64 {
	if !svd.computed {
		panic("svd: no decomposition computed")
	}
	if s == nil {
		s = make([]float64, len(svd.s))
	}
	if len(s) != len(svd.s) {
		panic(ErrSliceLengthMismatch)
	}
	copy(s, svd.s)
	return s
}

// UTo extracts the matrix U from the singular value decomposition. The first
// min(m,n) columns are the left singular vectors and correspond to the singular
// values as returned from SVD.Values.
//
// If dst is not nil, U is stored in-place into dst, and dst must have size
// m×m if svd.Kind() == SVDFull, size m×min(m,n) if svd.Kind() == SVDThin, and
// UTo panics otherwise. If dst is nil, a new matrix of the appropriate size is
// allocated and returned.
func (svd *SVD) UTo(dst *Dense) *Dense {
	kind := svd.U
	if kind != SVDFull && kind != SVDThin {
		panic("mat: improper SVD kind")
	}
	r := svd.u.Rows
	c := svd.u.Cols
	if dst == nil {
		dst = NewDense(r, c, nil)
	} else {
		dst.reuseAs(r, c)
	}

	tmp := &Dense{
		mat:     svd.u,
		capRows: r,
		capCols: c,
	}
	dst.Copy(tmp)

	return dst
}

// VTo extracts the matrix V from the singular value decomposition. The first
// min(m,n) columns are the right singular vectors and correspond to the singular
// values as returned from SVD.Values.
//
// If dst is not nil, V is stored in-place into dst, and dst must have size
// n×n if svd.Kind() == SVDFull, size n×min(m,n) if svd.Kind() == SVDThin, and
// VTo panics otherwise. If dst is nil, a new matrix of the appropriate size is
// allocated and returned.
func (svd *SVD) VTo(dst *Dense) *Dense {
	kind := svd.V
	if kind != SVDFull && kind != SVDThin {
		panic("mat: improper SVD kind")
	}
	r := svd.vt.Rows
	c := svd.vt.Cols
	if dst == nil {
		dst = NewDense(c, r, nil)
	} else {
		dst.reuseAs(c, r)
	}

	tmp := &Dense{
		mat:     svd.vt,
		capRows: r,
		capCols: c,
	}
	dst.Copy(tmp.T())

	return dst
}
