package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dsymmer interface {
	Dsymm(s blas.Side, ul blas.Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
}

func DsymmTest(t *testing.T, blasser Dsymmer) {
	for i, test := range []struct {
		m     int
		n     int
		side  blas.Side
		ul    blas.Uplo
		a     [][]float64
		b     [][]float64
		c     [][]float64
		alpha float64
		beta  float64
		ans   [][]float64
	}{
		{
			side: blas.Left,
			ul:   blas.Upper,
			m:    3,
			n:    4,
			a: [][]float64{
				{2, 3, 4},
				{0, 6, 7},
				{0, 0, 10},
			},
			b: [][]float64{
				{2, 3, 4, 8},
				{5, 6, 7, 15},
				{8, 9, 10, 20},
			},
			c: [][]float64{
				{8, 12, 2, 1},
				{9, 12, 9, 9},
				{12, 1, -1, 5},
			},
			alpha: 2,
			beta:  3,
			ans: [][]float64{
				{126, 156, 144, 285},
				{211, 252, 275, 535},
				{282, 291, 327, 689},
			},
		},
	} {
		aFlat := flatten(test.a)
		bFlat := flatten(test.b)
		cFlat := flatten(test.c)
		ansFlat := flatten(test.ans)
		blasser.Dsymm(test.side, test.ul, test.m, test.n, test.alpha, aFlat, len(test.a[0]), bFlat, test.n, test.beta, cFlat, test.n)
		if !floats.EqualApprox(cFlat, ansFlat, 1e-14) {
			t.Errorf("Case %v: Want %v, got %v.", i, ansFlat, cFlat)
		}
	}
}
