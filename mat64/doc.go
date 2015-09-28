// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mat64 provides implementations of float64 matrix structures and
// linear algebra operations on them.
//
// Overiew
//
// This section provides a quick overview of the matrix package. The following
// sections provide more in depth commentary.
//
// mat64 provides:
//  - Interfaces for Matrix classes (Matrix, Symmetric, Triangular)
//  - Concrete implementations (Dense, SymDense, TriDense)
//  - Methods and functions for using matrix data (Add, Trace, SymRankOne)
//  - Types for constructing and using matrix factorizations (QR, LU)
//
// A matrix may be constructed through the corresponding New function. If no
// backing array is provided the matrix will be initialized to all zeros.
//  // Allocate a zeroed array of size 3×5
//  zero := mat64.NewDense(3, 5, nil)
// If a backing data slice is provided, the matrix will have those elements.
// Matrices are all stored in row-major format.
//  // Generate a 6×6 matrix of random values.
//  data := make([]float64, 36)
//  for i := range data {
//		data[i] = rand.NormFloat64()
//  }
//  a := mat64.NewDense(6, 6, data)
//
// Operations involving matrix data are implemented as functions when the values
// of the matrix remain unchanged
//  tr := mat64.Trace(a)
// and are implemented as methods when the operation modifies the receiver.
//  zero.Copy(a)
//
// Receivers must be the correct size for the matrix operations, otherwise the
// operation will panic. As a special case for convenience, a zero-sized matrix
// will be modified to have the correct size, allocating data if necessary.
//  var c mat64.Dense // construct a new zero-sized matrix
//  c.Mul(a, a)       // c is automatically adjusted to be 6×6
//
// The Matrix Interface(s)
//
// The Matrix interface is the common link between the concrete types. The Matrix
// interface is defined by three functions: Dims, which returns the dimensions
// of the Matrix, At, which returns the element in the specified location, and
// T for returning a Transpose (discussed later). All of the concrete types can
// perform these behaviors and so implement the interface. Methods and functions
// are designed to use this interface, so in particular the method
//  func (m *Dense) Mul(a, b Matrix)
// constructs a *Dense from the result of a multiplication with any Matrix types,
// not just *Dense. Where more restrictive requirements must be met, there are also the
// Symmetric and Triangular interfaces. For example, in
//  func (s *SymDense) AddSym(a, b Symmetric)
// the Symmetric interface guarantees a symmetric result.
//
// Transposes
//
// The T method is used for transposition. For example, c.Mul(a.T(), b) computes
// c = a^T * b. The mat64 types implement this method using an implicit transpose —
// see the Transpose type for more details. Note that some operations have a
// transpose as part of their definition, as in *SymDense.SymOuterK.
//
// Matrix Factorization
//
// Matrix factorizations, such as the LU decomposition, typically have their own
// specific data storage, and so are each implemented as a specific type. The
// factorization can be computed through a call to Factorize
//  var lu mat64.LU
//  lu.Factorize(a)
// The elements of the factorization can be extracted through methods on the
// appropriate type, i.e. *TriDense.LFromLU and *TriDense.UFromLU. Alternatively,
// they can be used directly, as in *Dense.SolveLU. Some factorizations can be
// updated directly, without needing to update the original matrix and refactorize,
// as in *LU.RankOne.
//
// BLAS and LAPACK
//
// BLAS and LAPACK are the standard APIs for linear algebra routines. Many
// operations in mat64 are implemented using calls to the wrapper functions
// in gonum/blas/blas64 and gonum/lapack/lapack64. By default, blas64 and
// lapack64 call the native Go implementations of the routines. Alternatively,
// it is possible to use C-based implementations of the APIs through the respective
// cgo packages and "Use" functions. The Go implementation of LAPACK makes calls
// through blas64, so if a cgo BLAS implementation is registered, the lapack64
// calls will be partially executed in Go and partially executed in C.
//
// Type Switching
//
// The Matrix abstraction enables efficiency as well as interoperability. Go's
// type reflection capabilities are used to choose the most efficient routine
// given the specific concrete types. For example, in
//  c.Mul(a, b)
// if a and b both implement RawMatrixer, that is, they can be represented as a
// blas64.General, blas64.Gemm (general matrix multiplication) is called, while
// instead if b is a RawSymmetricer blas64.Symm is used (general-symmetric
// multiplication), and if b is a *Vector blas64.Gemv is used.
//
// There are many possible type combinations and special cases. No specific guarantees
// are made about the performance of any method, and in particular, note that an
// abstract matrix type may be copied into a concrete type of the corresponding
// value. If there are specific special cases that are needed, please submit a
// pull-request or file an issue.
//
// Invariants
//
// Matrix input arguments to functions are never modified. If an operation would
// change Matrix data, that matrix will be the receiver of a function.
//
// For convenience, a matrix may be used as both a receiver and as an input, i.e.
//  a.Pow(a, 6)
//  v.SolveVec(a.T(), v)
// though in many cases this will cause an allocation. However, mat64 cannot
// detect arbitrary overlap, and so care should be taken when using matrix views.
// An exception to this rule is Copy, which does not allow a.Copy(a.T()).
package mat64
