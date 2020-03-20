// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linsolve

import (
	"fmt"
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/lapack/testlapack"
	"gonum.org/v1/gonum/linsolve/internal/triplet"
	"gonum.org/v1/gonum/mat"
)

const defaultTol = 1e-13

type testCase struct {
	name string

	mulVecTo func(*mat.VecDense, bool, mat.Vector) // Matrix-vector multiplication

	b    []float64 // Right-hand side vector
	diag []float64 // Diagonal for the Jacobi preconditioner
	tol  float64   // Tolerance for the convergence criterion

	want []float64 // Expected solution
}

func (tc *testCase) MulVecTo(dst *mat.VecDense, trans bool, x mat.Vector) {
	tc.mulVecTo(dst, trans, x)
}

func (tc *testCase) PreconSolve(dst *mat.VecDense, trans bool, rhs mat.Vector) error {
	if tc.diag == nil {
		dst.CopyVec(rhs)
	} else {
		n := len(tc.diag)
		diag := mat.NewVecDense(n, tc.diag)
		dst.DivElemVec(rhs, diag)
	}
	return nil
}

func spdTestCases(rnd *rand.Rand) []testCase {
	return []testCase{
		newRandomSPD(1, rnd),
		newRandomSPD(2, rnd),
		newRandomSPD(3, rnd),
		newRandomSPD(4, rnd),
		newRandomSPD(5, rnd),
		newRandomSPD(10, rnd),
		newRandomSPD(20, rnd),
		newRandomSPD(50, rnd),
		newRandomDiagonal(2, rnd),
		newRandomDiagonal(3, rnd),
		newRandomDiagonal(4, rnd),
		newRandomDiagonal(5, rnd),
		newRandomDiagonal(10, rnd),
		newRandomDiagonal(20, rnd),
		newRandomDiagonal(50, rnd),
		newGreenbaum41(24, 0.001, 1, 0.4, rnd),
		newGreenbaum41(24, 0.001, 1, 0.6, rnd),
		newGreenbaum41(24, 0.001, 1, 0.8, rnd),
		newGreenbaum41(24, 0.001, 1, 1, rnd),
		newPoisson1D(32, random(rnd)),
		newPoisson2D(32, 32, one),
	}
}

// newRandomSPD returns a test case with a random symmetric positive-definite
// matrix of order n, and a random right-hand side.
func newRandomSPD(n int, rnd *rand.Rand) testCase {
	// Generate a random matrix C.
	c := make([]float64, n*n)
	for i := range c {
		c[i] = rnd.NormFloat64()
	}
	C := mat.NewDense(n, n, c)
	// Use C to generate a SPD matrix A.
	var A mat.SymDense
	A.SymOuterK(1, C)
	for i := 0; i < n; i++ {
		A.SetSym(i, i, A.At(i, i)+float64(n))
	}
	// Generate the right-hand side.
	b := make([]float64, n)
	for i := range b {
		b[i] = 1 / math.Sqrt(float64(n))
	}
	// Compute the solution using the Cholesky factorization.
	var chol mat.Cholesky
	ok := chol.Factorize(&A)
	if !ok {
		panic("bad test matrix")
	}
	want := make([]float64, n)
	chol.SolveVecTo(mat.NewVecDense(n, want), mat.NewVecDense(n, b))
	// Matrix-vector multiplication.
	mulVecTo := func(dst *mat.VecDense, _ bool, x mat.Vector) {
		if dst.Len() != n || x.Len() != n {
			panic("mismatched vector length")
		}
		dst.MulVec(&A, x)
	}
	// Store the diagonal for preconditioning.
	diag := make([]float64, n)
	for i := range diag {
		diag[i] = A.At(i, i)
	}
	return testCase{
		name:     fmt.Sprintf("Random SPD n=%v", n),
		mulVecTo: mulVecTo,
		b:        b,
		tol:      defaultTol,
		diag:     diag,
		want:     want,
	}
}

// newRandomDiagonal returns a test case with a diagonal matrix with random positive elements,
// a random right-hand side and a known solution.
func newRandomDiagonal(n int, rnd *rand.Rand) testCase {
	// Generate a diagonal matrix with random positive elements.
	a := make([]float64, n)
	diag := make([]float64, n)
	for i := range a {
		a[i] = 1 + 10*rnd.Float64()
		diag[i] = a[i]
	}
	A := mat.NewDiagDense(n, a)
	// Generate the right-hand side.
	b := make([]float64, n)
	for i := range b {
		b[i] = 1 / math.Sqrt(float64(n))
	}
	// Compute the reference solution.
	want := make([]float64, n)
	for i := range want {
		want[i] = b[i] / a[i]
	}
	// Matrix-vector multiplication.
	mulVecTo := func(dst *mat.VecDense, _ bool, x mat.Vector) {
		if dst.Len() != n || x.Len() != n {
			panic("mismatched vector length")
		}
		dst.MulVec(A, x)
	}
	return testCase{
		name:     fmt.Sprintf("Random diagonal n=%v", n),
		mulVecTo: mulVecTo,
		b:        b,
		tol:      defaultTol,
		diag:     diag,
		want:     want,
	}
}

// newGreenbaum41 returns a test case with a symmetric positive definite matrix
// A defined as
//  A = U * D * Uᵀ,
// where U is a random orthogonal matrix and D is a diagonal matrix with entries
// given by
//  d_i = d_1 + (i-1)/(n-1)*(d_n-d_1)*rho^{n-i},   i = 2,...,n-1.
//
// This matrix is described in Section 4.1 of
//  Greenbaum, A. (1997). Iterative Methods for Solving Linear Systems. SIAM.
func newGreenbaum41(n int, d1, dn, rho float64, rnd *rand.Rand) testCase {
	if n < 2 || dn < d1 {
		panic("bad test")
	}
	// Generate the diagonal.
	d := make([]float64, n)
	d[0] = d1
	d[n-1] = dn
	for i := 1; i < n-1; i++ {
		d[i] = d1 + float64(i)/float64(n-1)*(dn-d1)*math.Pow(rho, float64(n-i-1))
	}
	// Generate the SPD matrix A.
	a := make([]float64, n*n)
	testlapack.Dlagsy(n, 0, d, a, n, rnd, make([]float64, 2*n))
	A := mat.NewSymDense(n, a)
	// Generate the right-hand side.
	b := make([]float64, n)
	for i := range b {
		b[i] = rnd.NormFloat64()
	}
	// Compute the solution using the Cholesky factorization.
	var chol mat.Cholesky
	ok := chol.Factorize(A)
	if !ok {
		panic("bad test matrix")
	}
	want := make([]float64, n)
	chol.SolveVecTo(mat.NewVecDense(n, want), mat.NewVecDense(n, b))
	// Matrix-vector multiplication.
	mulVecTo := func(dst *mat.VecDense, _ bool, x mat.Vector) {
		if dst.Len() != n || x.Len() != n {
			panic("mismatched vector length")
		}
		dst.MulVec(A, x)
	}
	// Store the diagonal for preconditioning.
	diag := make([]float64, n)
	for i := range diag {
		diag[i] = A.At(i, i)
	}
	return testCase{
		name:     fmt.Sprintf("Greenbaum 4.1 n=%v,d_1=%v,d_n=%v,rho=%v", n, d1, dn, rho),
		mulVecTo: mulVecTo,
		b:        b,
		tol:      defaultTol,
		diag:     diag,
		want:     want,
	}
}

func nonsym3x3() testCase {
	return testCase{
		name: "nonsym 3x3",
		mulVecTo: func(dst *mat.VecDense, trans bool, x mat.Vector) {
			A := mat.NewDense(3, 3, []float64{
				5, -1, 3,
				-1, 2, -2,
				3, -2, 3,
			})
			if trans {
				dst.MulVec(A.T(), x)
			} else {
				dst.MulVec(A, x)
			}
		},
		b:    []float64{7, -1, 4},
		diag: []float64{5, 2, -3},
		tol:  defaultTol,
		want: []float64{1, 1, 1},
	}
}

func nonsymTridiag(n int) testCase {
	A := triplet.NewMatrix(n, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			A.Append(i, i-1, -2)
		}
		A.Append(i, i, 4)
		if i < n-1 {
			A.Append(i, i+1, -1)
		}
	}
	b := make([]float64, n)
	for i := range b {
		switch i {
		case 0:
			b[i] = 3
		default:
			b[i] = 1
		case n - 1:
			b[i] = 2
		}
	}
	want := make([]float64, n)
	for i := range want {
		want[i] = 1
	}
	return testCase{
		name:     fmt.Sprintf("Nonsym tridiag n=%v", n),
		mulVecTo: A.MulVecTo,
		b:        b,
		tol:      defaultTol,
		want:     want,
	}
}

// newPoisson1D returns a test case that arises from a finite-difference discretization
// of the Poisson equation
//  - ∂_x ∂_x u = f
// on the interval [0,1].
func newPoisson1D(nx int, f func(float64, float64) float64) testCase {
	tc := newPDE(nx, 1, negOne, nil, zero, nil, zero, f)
	tc.name = fmt.Sprintf("Poisson 1D nx=%v", nx)
	return tc
}

// newPoisson2D returns a test case that arises from a finite-difference discretization
// of the Poisson equation
//  - Δu = f
// on the unit square [0,1]×[0,1].
func newPoisson2D(nx, ny int, f func(float64, float64) float64) testCase {
	tc := newPDE(nx, ny, negOne, negOne, zero, zero, zero, f)
	tc.name = fmt.Sprintf("Poisson 2D nx=%v,ny=%v", nx, ny)
	tc.tol = 1e-12
	return tc
}

// newGreenbaum54 returns a test case with a general unsymmetric matrix
// A defined as
//  A = V*D*V^{-1},
// where V is a random matrix and D is a block-diagonal matrix with n1 complex
// and n2 real eigenvalues.
//
// This matrix is described in Section 5.4 of
//  Greenbaum, A. (1997). Iterative Methods for Solving Linear Systems. SIAM.
func newGreenbaum54(n1, n2 int, rnd *rand.Rand) testCase {
	n := 2*n1 + n2
	// Generate the block-diagonal matrix D.
	d := make([]float64, n*3)
	// Generate n1 blocks with complex eigenvalues.
	for i := 0; i < 2*n1; i += 2 {
		// The 2x2 block will have eigenvalues a±b*i.
		// Real part of the eigenvalue is in [1,2).
		a := rnd.Float64() + 1
		// Imaginary part is in [-1,1).
		b := 2*rnd.Float64() - 1
		d[i*3+1] = a
		d[i*3+2] = b
		d[(i+1)*3] = -b
		d[(i+1)*3+1] = a
	}
	// Generate n2 real eigenvalues.
	for i := 2 * n1; i < n; i++ {
		r := 9*rnd.Float64() + 1
		if rnd.Intn(2) == 0 {
			r *= -1
		}
		d[i*3+1] = r
	}
	D := mat.NewBandDense(n, n, 1, 1, d)
	// Generate a random matrix V.
	v := make([]float64, n*n)
	for i := range v {
		v[i] = rnd.NormFloat64()
	}
	V := mat.NewDense(n, n, v)
	var luV mat.LU
	luV.Factorize(V)
	// Generate the right-hand side.
	b := make([]float64, n)
	for i := range b {
		b[i] = rnd.NormFloat64()
	}
	// Compute V*D and (V*D)^{-1} for computing the reference solution and for the matrix-vector operation.
	var VD mat.Dense
	VD.Mul(V, D)
	var luVD mat.LU
	luVD.Factorize(&VD)
	// Compute the solution of V*D*V^{-1}*x = b.
	// First, compute the solution of V*D*y = b.
	want := make([]float64, n)
	wantVec := mat.NewVecDense(n, want)
	err := luVD.SolveVecTo(wantVec, false, mat.NewVecDense(n, b))
	if err != nil {
		panic("luVD.SolveVecTo(notrans) failed")
	}
	// Second, compute the solution of V^{-1}*x = y, which amounts to just
	// computing x = V*y.
	wantVec.MulVec(V, wantVec)
	// Matrix-vector multiplication.
	mulVecTo := func(dst *mat.VecDense, trans bool, x mat.Vector) {
		if dst.Len() != n || x.Len() != n {
			panic("mismatched vector length")
		}
		if trans {
			// Multiply (V*D*V^{-1})ᵀ * x which can be
			// rewritten as V^{-1}ᵀ * (V*D)ᵀ * x.
			dst.MulVec(VD.T(), x)
			err := luV.SolveVecTo(dst, true, dst)
			if err != nil {
				panic("luV.SolveVecTo(trans) failed")
			}
		} else {
			// Multiply V*D*V^{-1} * x.
			err := luV.SolveVecTo(dst, false, x)
			if err != nil {
				panic("luV.SolveVecTo(notrans) failed")
			}
			dst.MulVec(&VD, dst)
		}
	}
	return testCase{
		name:     fmt.Sprintf("Greenbaum 5.4 n=%v,n1=%v,n2=%v", n, n1, n2),
		mulVecTo: mulVecTo,
		b:        b,
		tol:      defaultTol,
		want:     want,
	}
}

// newGreenbaum73 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  - Δu + 40*(x*∂_x u + y*∂_y u) - 100*u = f
//
// This test problem is described in Section 7.3 of
//  Greenbaum, A. (1997). Iterative Methods for Solving Linear Systems. SIAM.
func newGreenbaum73(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny,
		negOne, negOne, // - ∂_x ∂_x u - ∂_y ∂_y u
		func(x, _ float64) float64 { return 40 * x }, // 40 * x * ∂_x u
		func(_, y float64) float64 { return 40 * y }, // 40 * y * ∂_y u
		constant(-100), // -100 * u
		random(rnd))
	tc.name = fmt.Sprintf("Greenbaum 7.3 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDENonsymmetric returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu + henk*∂_x u + (∂_x henk/2)*u = f
// where henk(x,y) := 20*exp(3.5*(x^2 + y^2))
//
// This test case is adapted from
//  http://www.netlib.org/templates/dftemplates.tgz
func newPDENonsymmetric(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny, one, one, henk, zero, dhenkdx, random(rnd))
	tc.name = fmt.Sprintf("PDE Nonsymmetric nx=%v,ny=%v", nx, ny)
	return tc
}

func henk(x, y float64) float64 {
	return 20 * math.Exp(3.5*(x*x+y*y))
}

func dhenkdx(x, y float64) float64 {
	return 70 * x * math.Exp(3.5*(x*x+y*y))
}

// newPDEYang47 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu + 1000*∂_x u = f
// The large coefficient of ∂_xu causes a loss of diagonal dominance in the system matrix.
//
// This test case corresponds to Eq. 4.7 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang47(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny, one, one, constant(1000), zero, zero, random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.7 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang48 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu + 1000*∂_y u = f
// The large coefficient of ∂_yu causes a loss of diagonal dominance in the system matrix.
//
// This test case corresponds to Eq. 4.8 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang48(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny, one, one, zero, constant(1000), zero, random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.8 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang49 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu + 1000*exp(x*y)*(∂_x u - ∂_y u) = f
//
// This test case corresponds to Eq. 4.9 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang49(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(16, 16,
		one, one,
		func(x, y float64) float64 { return 1000 * math.Exp(x*y) },
		func(x, y float64) float64 { return -1000 * math.Exp(x*y) },
		zero, random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.9 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang410 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  ∂_x (exp(x*y) * ∂_x u) + ∂_y (exp(-x*y) * ∂_y u) - 500*(x + y)*∂_x u - (250*(1 + y) + 1/(1 + x + y))*u = f
//
// This test case corresponds to Eq. 4.10 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang410(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny,
		func(x, y float64) float64 { return math.Exp(x * y) },  // ∂_x (exp(x*y) * ∂_x u)
		func(x, y float64) float64 { return math.Exp(-x * y) }, // ∂_y (exp(-x*y) * ∂_y u)
		func(x, y float64) float64 { return -500 * (x + y) },   // -500*(x + y)*∂_x u
		zero,
		func(x, y float64) float64 { return -250*(1+y) - 1/(1+x+y) }, // - (250*(1 + y) + 1/(1+x+y)) * u
		random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.10 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang412 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu - 100000*x^2*(∂_x u + ∂_y u) = f
//
// This test case corresponds to Eq. 4.12 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang412(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny,
		one, one,
		func(x, _ float64) float64 { return -100000 * x * x },
		func(x, _ float64) float64 { return -100000 * x * x },
		zero, random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.12 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang413 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu - 1000*(1 + x^2)*∂_x u + 100*∂_y u = f
//
// This test case corresponds to Eq. 4.13 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang413(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny,
		one, one,
		func(x, _ float64) float64 { return -1000 * (1 + x*x) },
		constant(100),
		zero, random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.13 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang414 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu - 1000*x^2*∂_x u + 1000*u = f
//
// This test case corresponds to Eq. 4.14 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang414(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny,
		one, one,
		func(x, _ float64) float64 { return -1000 * x * x },
		zero,
		constant(1000), random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.14 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDEYang415 returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  Δu - 1000*(1 - 2*x)*∂_x u - 1000*(1 - 2*y)*∂_y u = f
//
// This test case corresponds to Eq. 4.15 in
//  Ulrike Meier Yang (1994), Preconditioned Conjugate Gradient-Like Methods for Nonsymmetric Linear Systems.
func newPDEYang415(nx, ny int, rnd *rand.Rand) testCase {
	tc := newPDE(nx, ny,
		one, one,
		func(x, _ float64) float64 { return -1000 * (1 - 2*x) },
		func(_, y float64) float64 { return -1000 * (1 - 2*y) },
		zero, random(rnd))
	tc.name = fmt.Sprintf("PDE Yang, Eq. 4.15 nx=%v,ny=%v", nx, ny)
	return tc
}

// newPDE returns a test case that arises from a finite-difference discretization
// of the partial differential equation
//  ∂_x (a ∂_x u) + ∂_y (b ∂_y u) + c ∂_x u + d ∂_y u + e u = f
// on the unit square [0,1]×[0,1] with zero Dirichlet boundary conditions.
//
// nx and ny must be positive. If ny is 1, a 1D variant of the equation on [0,1]×{0} will be used,
// and b and d will not be referenced.
func newPDE(nx, ny int, a, b, c, d, e, f func(float64, float64) float64) testCase {
	if nx <= 0 || ny <= 0 {
		panic("invalid mesh size")
	}

	var (
		A    *triplet.Matrix
		rhs  []float64
		diag []float64
	)
	if ny == 1 {
		A, rhs, diag = newPDESystem1D(nx, a, c, e, f)
	} else {
		A, rhs, diag = newPDESystem2D(nx, ny, a, b, c, d, e, f)
	}
	// Use a dense solver to obtain a reference solution.
	Ad := A.DenseCopy()
	var lu mat.LU
	lu.Factorize(Ad)
	n := len(rhs)
	want := make([]float64, n)
	err := lu.SolveVecTo(mat.NewVecDense(n, want), false, mat.NewVecDense(n, rhs))
	if err != nil {
		panic("lu.SolveVec failed")
	}
	return testCase{
		mulVecTo: A.MulVecTo,
		b:        rhs,
		tol:      defaultTol,
		diag:     diag,
		want:     want,
	}
}

// newPDESystem1D assembles and returns the matrix A, the right-hand side vector,
// and the diagonal of A for a 1-dimensional PDE problem.
func newPDESystem1D(nx int, a, b, c, f func(float64, float64) float64) (A *triplet.Matrix, rhs, diag []float64) {
	h := 1 / float64(nx+1)
	A = triplet.NewMatrix(nx, nx)
	rhs = make([]float64, nx)
	diag = make([]float64, nx)
	var i int
	for ix := 0; ix < nx; ix++ {
		s := newStencil1D(ix, h, a, b, c)
		// Add stencil elements to the system matrix, skipping boundary nodes.
		if ix > 0 {
			A.Append(i, i-1, s.left)
		}
		A.Append(i, i, s.center)
		diag[i] = s.center
		if ix < nx-1 {
			A.Append(i, i+1, s.right)
		}
		// Assemble the right-hand side.
		x := float64(ix+1) * h
		rhs[i] = f(x, 0) * h * h
		i++
	}
	return A, rhs, diag
}

type stencil1D struct {
	left, right float64
	center      float64
}

// newStencil1D returns a finite difference stencil that approximates the differential operator
//  ∂_x (a ∂_x u) + b ∂_x u + c u
// at point [(i+1)*h] using the formula
//    a(i+1/2,0)*(u(i+1) - u(i)) + a(i-1/2,0)*(u(i-1) - u(i))
//  + (h/2)*b(i,0)*(u(i+1) - u(i-1))
//  + h^2*c(i,0)*u(i)
func newStencil1D(i int, h float64, a, b, c func(float64, float64) float64) (s stencil1D) {
	x := float64(i+1) * h

	coeff := a(x+0.5*h, 0)
	s.center -= coeff
	s.right = coeff
	coeff = a(x-0.5*h, 0)
	s.center -= coeff
	s.left = coeff

	coeff = b(x, 0)
	s.right += 0.5 * h * coeff
	s.left -= 0.5 * h * coeff

	s.center += h * h * c(x, 0)

	return s
}

// newPDESystem2D assembles and returns the matrix A, the right-hand side vector,
// and the diagonal of A for a 2-dimensional PDE problem.
func newPDESystem2D(nx, ny int, a, b, c, d, e, f func(float64, float64) float64) (A *triplet.Matrix, rhs, diag []float64) {
	// Finite difference stencil:
	//             * (ix,iy+1)
	//             |
	//             | (ix,iy)
	//      * ---- * ---- *
	// (ix-1,iy)   |   (ix+1,iy)
	//             |
	//             * (ix,iy-1)
	// Node (ix,iy) is mapped to the index ix+iy*k.

	h := 1 / float64(nx+1)
	n := nx * ny
	A = triplet.NewMatrix(n, n)
	rhs = make([]float64, n)
	diag = make([]float64, n)
	var i int
	for iy := 0; iy < ny; iy++ {
		y := float64(iy+1) * h
		for ix := 0; ix < nx; ix++ {
			s := newStencil2D(ix, iy, h, a, b, c, d, e)
			// Add the coefficients from the stencil to the system matrix, skipping boundary nodes.
			if iy > 0 {
				A.Append(i, i-nx, s.down)
			}
			if ix > 0 {
				A.Append(i, i-1, s.left)
			}
			A.Append(i, i, s.center)
			diag[i] = s.center
			if ix < nx-1 {
				A.Append(i, i+1, s.right)
			}
			if iy < ny-1 {
				A.Append(i, i+nx, s.up)
			}
			// Assemble the right-hand side.
			x := float64(ix+1) * h
			rhs[i] = f(x, y) * h * h
			i++
		}
	}
	return A, rhs, diag
}

type stencil2D struct {
	left, right float64
	up, down    float64
	center      float64
}

// newStencil2D returns a finite difference stencil that approximates the differential operator
//  ∂_x (a ∂_x u) + ∂_y (b ∂_y u) + c ∂_x u + d ∂_y u + e u
// at point [(i+1)*h,(j+1)*h] using the formula
//    a(i+1/2,j)*(u(i+1,j) - u(i,j)) + a(i-1/2,j)*(u(i-1,j) - u(i,j))
//  + b(i,j+1/2)*(u(i,j+1) - u(i,j)) + b(i,j-1/2)*(u(i,j-1) - u(i,j))
//  + (h/2)*c(i,j)*(u(i+1,j) - u(i-1,j)) + (h/2)*d(i,j)*(u(i,j+1) - u(i,j-1))
//  + h^2*e(i,j)*u(i,j)
func newStencil2D(i, j int, h float64, a, b, c, d, e func(float64, float64) float64) (s stencil2D) {
	x := float64(i+1) * h
	y := float64(j+1) * h

	coeff := a(x+0.5*h, y)
	s.center -= coeff
	s.right = coeff
	coeff = a(x-0.5*h, y)
	s.center -= coeff
	s.left = coeff
	coeff = b(x+0.5*h, y)
	s.center -= coeff
	s.up = coeff
	coeff = b(x-0.5*h, y)
	s.center -= coeff
	s.down = coeff

	coeff = c(x, y)
	s.right += 0.5 * h * coeff
	s.left -= 0.5 * h * coeff
	coeff = d(x, y)
	s.up += 0.5 * h * coeff
	s.down -= 0.5 * h * coeff

	s.center += h * h * c(x, y)

	return s
}

func zero(_, _ float64) float64 {
	return 0
}

func one(_, _ float64) float64 {
	return 1
}

func negOne(_, _ float64) float64 {
	return -1
}

func constant(c float64) func(_, _ float64) float64 {
	return func(_, _ float64) float64 {
		return c
	}
}

func random(rnd *rand.Rand) func(_, _ float64) float64 {
	return func(_, _ float64) float64 {
		return rnd.NormFloat64()
	}
}
