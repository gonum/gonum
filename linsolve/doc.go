// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package linsolve provides iterative methods for solving linear systems.

Background

A system of linear equations can be written as

 A * x = b,

where A is a given n×n non-singular matrix, b is a given n-vector (the
right-hand side), and x is an unknown n-vector.

Direct methods such as the LU or QR decomposition compute (in the absence of
roundoff errors) the exact solution after a finite number of steps. For a
general matrix A they require O(n^3) arithmetic operations which becomes
infeasible for large n due to excessive memory and time cost.

Iterative methods, in contrast, generally do not compute the exact solution x.
Starting from an initial estimate x_0, they instead compute a sequence x_i of
increasingly accurate approximations to x. This iterative process is stopped
when the estimated difference between x_i and the true x becomes smaller than a
prescribed threshold. The number of iterations thus depends on the value of the
threshold. If the desired threshold is very small, then the iterative methods
may take as long or longer than a direct method. However, for many problems a
decent approximation can be found in a small number of steps.

The iterative methods implemented in this package do not access the elements of
A directly, they instead ask for the result of matrix-vector products with A.
For a general n×n matrix this requires O(n^2) operations, but can be much
cheaper depending on the structure of the matrix (sparse, banded,
block-stuctured, etc.). Such structure often arises in practical applications.
An iterative method can thus be significantly cheaper than a direct method, by
using a small number of iterations and taking advantage of matrix structure.

Iterative methods are most often useful in the following situations:

 - The system matrix A is sparse, blocked or has other special structure,
 - The problem size is sufficiently large that a dense factorization of A is
   costly in terms of compute time and/or memory storage,
 - Computing the product of A (or A^T, if necessary) with a vector can be done
   efficiently,
 - An approximate solution is all that is required.

Using linsolve

The two most important elements of the API are the MulVecToer interface and the
Iterative function.

MulVecToer interface

The MulVecToer interface represents the system matrix A. This abstracts the
details of any particular matrix storage, and allows the user to exploit the
properties of their particular matrix. Matrix types provided by gonum/mat and
github.com/james-bowman/sparse packages implement this interface.

Note that methods in this package have only limited means for checking whether
the provided MulVecToer represents a matrix that satisfies all assumptions made
by the chosen Method, for example if the matrix is actually symmetric positive
definite.

Iterative function

The Iterative function is the entry point to the functionality provided by this
package. It takes as parameters the matrix A (via the MulVecToer interface as
discussed above), the right-hand side vector b, the iterative method and
settings that control the iterative process and provide a way for reusing
memory.

Choosing an iterative method

The choice of an iterative method is typically guided by the properties of the
matrix A including symmetry, definiteness, sparsity, conditioning, and block
structure. In general, performance on symmetric matrices is well understood (see
the references below), with the conjugate gradient method being a good starting
point. Non-symmetric matrices are much more difficult to assess, where any
suggestion of a 'best' method is usually accompanied by a recommendation to use
trial-and-error.

Preconditioning

Preconditioning is a family of techniques that attempt to transform the linear
system into one that has the same solution but more favorable eigenspectrum. The
transformation matrix is called a preconditioner. A good preconditioner will
reduce the number of iterations needed to find a good approximate solution
(hopefully enough to overcome the cost of applying the preconditioning!), and in
some cases preconditioning is necessary to get any kind of convergence. In
linsolve a preconditioner is specified by Settings.PreconSolve.

Implementing Method interface

This package allows external implementations of iterative solvers by means of
the Method interface. It uses a reverse-communication style of API to
"outsource" operations such as matrix-vector multiplication, preconditioner
solve or convergence checks to the caller. The caller performs the commanded
operation and passes the result again to Method. The matrix A and the right-hand
side b are not directly available to Methods which encourages their cleaner
implementation. See the documentation for Method, Operation, and Context for
more information.

References

Further details about computational practice and mathematical theory of
iterative methods can be found in the following references:

 - Barrett, Richard et al. (1994). Templates for the Solution of Linear Systems:
   Building Blocks for Iterative Methods (2nd ed.). Philadelphia, PA: SIAM.
   Retrieved from http://www.netlib.org/templates/templates.pdf
 - Saad, Yousef (2003). Iterative methods for sparse linear systems (2nd ed.).
   Philadelphia, PA: SIAM. Retrieved from
   http://www-users.cs.umn.edu/~saad/IterMethBook_2ndEd.pdf
 - Greenbaum, A. (1997). Iterative methods for solving linear systems.
   Philadelphia, PA: SIAM.
*/
package linsolve
