package examples

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

// Solve2x2 demonstrates solving a small linear system with VecDense.SolveVec,
// computing the residual, and inspecting the condition number.
func Solve2x2() {
	// System: A x = b
	//   2x + 3y = 8
	//   4x + 5y = 13
	A := mat.NewDense(2, 2, []float64{
		2, 3,
		4, 5,
	})
	b := mat.NewVecDense(2, []float64{8, 13})

	// Allocate output vector.
	x := mat.NewVecDense(2, nil)

	// Solve.
	err := x.SolveVec(A, b)
	if err != nil {
		fmt.Printf("solve failed: %v\n", err)
		return
	}

	fmt.Printf("solution x = %v\n", x.RawVector())

	// Compute residual r = b - A*x.
	r := mat.NewVecDense(2, nil)
	r.SymmetricMult(A, x)
	r.Sub(b, r)
	residual := r.Norm(nil, 2)
	fmt.Printf("||r||_2 = %.2e\n", residual)

	// Condition number.
	var sv mat.SVD
	if err := sv.Factorize(A); err == nil {
		cond := sv.Cond()
		fmt.Printf("cond(A) = %.2e\n", cond)
	}

	// Expected solution: x = 1, y = 2.
	expected := []float64{1, 2}
	for i, v := range x.RawVector().Data {
		diff := math.Abs(v - expected[i])
		fmt.Printf("  x[%d] diff = %.2e\n", i, diff)
	}
}
