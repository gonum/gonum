// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mat64 provides basic linear algebra operations for float64 matrices.
//
// mat64 provides a set of concrete types that implement different classes of
// matrices (Dense, Symmetric, etc.) and operations on them. In most cases,
// an operation which results in a matrix value is a method on the value being
// produced. As an example, the method *Dense.Add(a,b) performs element-wise
// addition of the matrices a and b and stores the result into the method receiver.
// In all operations that assign to the receiver, the receiver may either have
// the correct dimensions for the result of the method or may be a zero-sized matrix.
// In the latter case, the receiver will be modified to the correct size, allocating
// data if necessary. The T() method is used for transposition; c.Mul(a.T(), b) computes
// c = a^T * b. Again, if c is a zero-sized matrix it will be modified to be the
// correct size, otherwise c must be size ac×bc, that is the number of columns in
// a (the number of rows in a^T) by the number of columns in b.
//
// Many operations are performed by calling into the gonum/blas/blas64 and
// gonum/lapack/lapack64 packages which implement BLAS and LAPACK routines
// respectively. By default, blas64 and lapack64 use native Go implementations
// of the routines. Alternatively, it is possible to use c-based libraries for
// the BLAS routines and/or the LAPACK routines with the respective cgo packages
// and Use routines. The Go implementation of LAPACK itself makes calls to blas64,
// so if a cgo BLAS implementation is used, the lapack64 calls will be partially
// executed in Go and partially executed in c.
//
// Along with the concrete matrix types are corresponding interface types. The
// Matrix interface, for example, represents an arbitrary matrix of float64 values,
// and the Symmetricer interface represents an arbitrary symmetric matrix. This
// abstract representation enables efficiency. In c.Mul(a,b), for example, if
// a is a RawSymmetricer and b is a RawMatrixer the specialized function for
// symmetric-dense multiplication (blas64.Symm) will be called. No specific
// guarantees are made about the performance of any method, as many special cases
// are possible. In particular, note that an abstract matrix type may be copied
// into a concrete type of the corresponding value. If there are specific special
// cases that are needed, please submit a pull-request or file an issue.
package mat64
