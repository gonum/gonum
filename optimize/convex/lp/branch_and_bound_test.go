package lp

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestBNB(t *testing.T) {
	t.Parallel()

	// check known examples

	// Example based of example on wikipedia
	// https://en.wikipedia.org/wiki/Integer_programming
	// modified to have one optimal solution.
	_, x, err := BNB(
		[]float64{-0.5, -2},
		nil,
		[]float64{},
		mat.NewDense(3, 2, []float64{
			-1, +1,
			+3, +2,
			+2, +3,
		}),
		[]float64{1, 12, 12},
		[]bool{true, true},
		1e-10,
	)

	if err != nil {
		t.Errorf("unexpected error obtained: %s", err)
	}
	correctSolution := (2-1e-10 < x[0]) && (x[0] < 2+1e-10) && (2-1e-10 < x[1]) && (x[1] < 2+1e10)
	if !correctSolution {
		t.Errorf("Solution found isn't correct (found %f, %f)", x[0], x[1])
	}

	// Check for infeasibility
	// This happens if no whole number solution exists
	// test is a diamond with no whole numbers contained within.
	// Ideal solution would be (0, 0.5), but it's not a whole number.
	_, _, err = BNB(
		[]float64{1, 0},
		nil,
		[]float64{},
		mat.NewDense(4, 2, []float64{
			+1, +1,
			-1, +1,
			+1, -1,
			-1, -1,
		}),
		[]float64{
			+1.5,
			+0.5,
			+0.5,
			-0.5,
		},
		[]bool{true, true},
		1e-10,
	)

	if err != ErrInfeasible {
		t.Errorf("Infeasible problem not found as such")
	}
}
