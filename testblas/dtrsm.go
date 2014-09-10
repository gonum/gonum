package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type Dtrsmer interface {
	Dtrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag,
		m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
}

func TestDtrsm(t *testing.T, blasser Dtrsmer) {

	for i, test := range []struct {
		s        blas.Side
		ul       blas.Uplo
		tA       blas.Transpose
		d        blas.Diag
		m, n     int
		alpha    float64
		a, b     [][]float64
		lda, ldb int

		ans    [][]float64
		panics bool
	}{
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     3,
			n:     1,
			alpha: 2,
			lda:   3,
			ldb:   1,

			a: [][]float64{
				{27.760766912759156, -0.021508733517253534, -0.5705006213444302},
				{0.0004901483858841176, 0.018718921572188584, 0.13804515090321257},
				{-0.021508733517253534, -0.5705006213444302, 15.0000324839208237},
			},
			b: [][]float64{
				{7},
				{8},
				{9},
			},

			ans: [][]float64{
				{0.504308834262264},
				{854.7368902709834},
				{33.70917860157599},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     3,
			n:     1,
			alpha: 2,
			lda:   3,
			ldb:   1,

			a: [][]float64{
				{27.760766912759156, -0.021508733517253534, -0.5705006213444302},
				{0.0004901483858841176, 0.018718921572188584, 0.13804515090321257},
				{-0.021508733517253534, -0.5705006213444302, 15.0000324839208237},
			},
			b: [][]float64{{7}},

			panics: true,
		},
	} {

		const name = "RowMajor"

		aFlat := flatten(test.a)
		aCopy := flatten(test.a)
		bFlat := flatten(test.b)
		ansFlat := flatten(test.ans)

		fn := func() {
			blasser.Dtrsm(
				test.s, test.ul, test.tA, test.d,
				test.m, test.n,
				test.alpha,
				aFlat, test.lda,
				bFlat, test.ldb,
			)
		}

		if panics(fn) != test.panics {
			if test.panics {
				t.Errorf("Test %d expected panic.", i)
			} else {
				t.Errorf("Test %d unexpected panic.", i)
			}
		}
		if test.panics {
			continue
		}

		// Check that a is unchanged
		if !dSliceEqual(aFlat, aCopy) {
			t.Errorf("Test %d case %v: a changed during Dtrsm", i, name)
		}
		if !dSliceTolEqual(ansFlat, bFlat) {
			t.Errorf("Test %d case %v: answer mismatch. Expected %v, Found %v. ", i, name, ansFlat, bFlat)
		}
	}
}
