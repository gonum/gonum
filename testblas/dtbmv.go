package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type Dtbmver interface {
	Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int)
}

func DtbmvTest(t *testing.T, blasser Dtbmver) {
	for i, test := range []struct {
		ul  blas.Uplo
		tA  blas.Transpose
		d   blas.Diag
		n   int
		k   int
		a   [][]float64
		lda int
		x   []float64
	}{
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{},
			},
			x: []float64{},
		},
	} {
		aFlat := flattenBanded(test.a, test.kU, test.kL)
		yCopy := sliceCopy(test.y)
	}
}
