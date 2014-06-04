package testblas

import "github.com/gonum/blas"

type Dgbmver interface {
	Dgbmv(tA blas.Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

// TODO: Redo this test with row major implementation
/*
func DgbmvTest(t *testing.T, blasser Dgbmver) {
	A := []float64{0, 0, 3, 6, 0, 1, 4, 7, 9, 2, 5, 8}
	kL := 1
	kU := 2
	m := 4
	n := 3

	x1 := []float64{2, 3, 4, 5}
	incX1 := 1
	x2 := []float64{2, 1, 3, 1, 4, 1, 5}
	incX2 := 2

	solNoTrans := []float64{45, 32, 41, 32}

	in := make([]float64, len(x1))
	out := make([]float64, m*incX1)
	copy(in, x1)
	blasser.Dgbmv(blas.NoTrans, m, n, kL, kU, 1, A, m, in, incX1, 0, out, incX1)

	if !dStridedSliceTolEqual(m, out, incX1, solNoTrans, 1) {
		t.Error("Wrong Dgbmv result for: ColMajor, NoTrans, IncX==1")
	}

	in = make([]float64, len(x2))
	out = make([]float64, m*incX2)
	copy(in, x2)
	blasser.Dgbmv(blas.NoTrans, m, n, kL, kU, 1, A, m, in, incX2, 0, out, incX2)

	if !dStridedSliceTolEqual(m, out, incX2, solNoTrans, 1) {
		t.Error("Wrong Dgbmv result for: ColMajor, NoTrans, IncX==2")
	}

	solTrans := []float64{24, 42, 84}

	in = make([]float64, len(x1))
	out = make([]float64, n*incX1)
	copy(in, x1)
	blasser.Dgbmv(blas.Trans, m, n, kL, kU, 1, A, m, in, incX1, 0, out, incX1)

	if !dStridedSliceTolEqual(n, out, incX1, solTrans, 1) {
		t.Error("Wrong Dgbmv result for: ColMajor, Trans, IncX==1")
	}

	in = make([]float64, len(x2))
	out = make([]float64, n*incX2)
	copy(in, x2)
	blasser.Dgbmv(blas.Trans, m, n, kL, kU, 1, A, m, in, incX2, 0, out, incX2)

	if !dStridedSliceTolEqual(n, out, incX2, solTrans, 1) {
		t.Error("Wrong Dgbmv result for: ColMajor, Trans, IncX==2")
	}
}
*/
