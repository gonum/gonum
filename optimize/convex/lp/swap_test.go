package lp

import (
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestSwap(t *testing.T) {
	for _, test := range []struct {
		swap  *Swap
		b     []float64
		ans   []float64
		trans []float64
		err   error
	}{
		{
			swap:  &Swap{3, []int{1, 0, 2, 0}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, 0},
			b:     []float64{1, 1, 1},
			ans:   []float64{0.10972222, -0.22083333, -1.45555556},
			trans: []float64{-3.2, 0.7, 0.93333333},
			err:   nil,
		},
		{
			swap:  &Swap{3, []int{1, 0, 2, 0}, []float64{1, 0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, 0},
			b:     []float64{1, 1, 1},
			ans:   []float64{0.10972222, -0.22083333, -1.45555556},
			trans: []float64{-3.2, 0.7, 0.93333333},
			err:   ErrSwapSingular,
		},
	} {
		v := make([]float64, test.swap.Dim)
		vVec := mat.NewVecDense(len(v), v)
		bVec := mat.NewVecDense(test.swap.Dim, test.b)
		vVec.CopyVec(bVec)

		err := test.swap.SolveVec(nil, false, vVec)
		if err != nil {
			if test.err == nil {
				t.Errorf("Unexpected error: %s", err)
			} else if err != test.err {
				t.Errorf("Error mismatch. Want %v, got %v", test.err, err)
			}
			continue
		}
		if err == nil && test.err != nil {
			t.Errorf("Did not error during solve.")
			continue
		}
		if !floats.EqualApprox(v, test.ans, 1e-07) {
			t.Errorf("Solution mismatch. Want %v, got %v", test.ans, v)
		}

		err = test.swap.SolveVec(vVec, true, bVec)
		if err != nil {
			if test.err == nil {
				t.Errorf("Unexpected error: %s", err)
			} else if err != test.err {
				t.Errorf("Error mismatch. Want %v, got %v", test.err, err)
			}
			continue
		}
		if err == nil && test.err != nil {
			t.Errorf("Did not error during transpose solve.")
			continue
		}
		if !floats.EqualApprox(v, test.trans, 1e-07) {
			t.Errorf("Transpose solution mismatch. Want %v, got %v", test.ans, v)
		}

	}
}
