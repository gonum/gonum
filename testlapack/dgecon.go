package testlapack

import (
	"math"
	"testing"

	"github.com/gonum/lapack"
)

type Dgeconer interface {
	Dlanger
	Dgetrfer
	Dgecon(norm lapack.MatrixNorm, n int, a []float64, lda int, anorm float64, work []float64, iwork []int) float64
}

func DgeconTest(t *testing.T, impl Dgeconer) {
	for _, test := range []struct {
		m       int
		n       int
		a       []float64
		condOne float64
		condInf float64
	}{
		{
			a: []float64{
				8, 1, 6,
				3, 5, 7,
				4, 9, 2,
			},
			m:       3,
			n:       3,
			condOne: 3.0 / 16,
			condInf: 3.0 / 16,
		},
		{
			a: []float64{
				2, 9, 3, 2,
				10, 9, 9, 3,
				1, 1, 5, 2,
				8, 4, 10, 2,
			},
			m:       4,
			n:       4,
			condOne: 0.024740155174938,
			condInf: 0.012034465570035,
		},
	} {
		m := test.m
		n := test.n
		lda := n
		a := make([]float64, len(test.a))
		copy(a, test.a)
		ipiv := make([]int, min(m, n))

		// Find the norms of the original matrix.
		work := make([]float64, 4*n)
		oneNorm := impl.Dlange(lapack.MaxColumnSum, m, n, a, lda, work)
		infNorm := impl.Dlange(lapack.MaxRowSum, m, n, a, lda, work)

		// Compute LU factorization of a.
		impl.Dgetrf(m, n, a, lda, ipiv)

		// Compute the condition number
		iwork := make([]int, n)
		condOne := impl.Dgecon(lapack.MaxColumnSum, n, a, lda, oneNorm, work, iwork)
		condInf := impl.Dgecon(lapack.MaxRowSum, n, a, lda, infNorm, work, iwork)
		if math.Abs(condOne-test.condOne) > 1e-13 {
			t.Errorf("One norm mismatch. Want %v, got %v.", test.condOne, condOne)
		}
		if math.Abs(condInf-test.condInf) > 1e-13 {
			t.Errorf("Inf norm mismatch. Want %v, got %v.", test.condInf, condInf)
		}
	}
}
