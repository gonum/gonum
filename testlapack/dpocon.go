package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
)

type Dpoconer interface {
	Dpotrfer
	Dgeconer
	Dlansy(norm lapack.MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int, work []float64) float64
	Dpocon(uplo blas.Uplo, n int, a []float64, lda int, anorm float64, work []float64, iwork []int) float64
}

func DpoconTest(t *testing.T, impl Dpoconer) {
	for _, test := range []struct {
		a    []float64
		n    int
		cond float64
		uplo blas.Uplo
	}{
		{
			a: []float64{
				89, 59, 77,
				0, 107, 59,
				0, 0, 89,
			},
			uplo: blas.Upper,
			n:    3,
			cond: 0.050052137643379,
		},
		{
			a: []float64{
				89, 0, 0,
				59, 107, 0,
				77, 59, 89,
			},
			uplo: blas.Lower,
			n:    3,
			cond: 0.050052137643379,
		},
	} {
		n := test.n
		a := make([]float64, len(test.a))
		copy(a, test.a)
		lda := n
		uplo := test.uplo
		work := make([]float64, 3*n)
		anorm := impl.Dlansy(lapack.MaxColumnSum, uplo, n, a, lda, work)
		// Compute cholesky decomposition
		ok := impl.Dpotrf(uplo, n, a, lda)
		if !ok {
			t.Errorf("Bad test, matrix not positive definite")
			continue
		}
		iwork := make([]int, n)
		cond := impl.Dpocon(uplo, n, a, lda, anorm, work, iwork)
		if math.Abs(cond-test.cond) > 1e-14 {
			t.Errorf("Cond mismatch. Want %v, got %v.", test.cond, cond)
		}
	}
	bi := blas64.Implementation()
	// Randomized tests compared against Dgecon.
	for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
		for _, test := range []struct {
			n, lda int
		}{
			{3, 0},
			{3, 5},
		} {
			for trial := 0; trial < 100; trial++ {
				n := test.n
				lda := test.lda
				if lda == 0 {
					lda = n
				}
				a := make([]float64, n*lda)
				for i := range a {
					a[i] = rand.NormFloat64()
				}

				// Multiply a by itself to make it symmetric positive definite.
				aCopy := make([]float64, len(a))
				copy(aCopy, a)
				bi.Dgemm(blas.Trans, blas.NoTrans, n, n, n, 1, aCopy, lda, aCopy, lda, 0, a, lda)

				aDense := make([]float64, len(a))
				if uplo == blas.Upper {
					for i := 0; i < n; i++ {
						for j := i; j < n; j++ {
							v := a[i*lda+j]
							aDense[i*lda+j] = v
							aDense[j*lda+i] = v
						}
					}
				} else {
					for i := 0; i < n; i++ {
						for j := 0; j <= i; j++ {
							v := a[i*lda+j]
							aDense[i*lda+j] = v
							aDense[j*lda+i] = v
						}
					}
				}
				work := make([]float64, 4*n)
				iwork := make([]int, n)

				anorm := impl.Dlansy(lapack.MaxColumnSum, uplo, n, a, lda, work)
				ok := impl.Dpotrf(uplo, n, a, lda)
				if !ok {
					t.Errorf("Bad test, matrix not positive definite")
					continue
				}
				got := impl.Dpocon(uplo, n, a, lda, anorm, work, iwork)

				denseNorm := impl.Dlange(lapack.MaxColumnSum, n, n, aDense, lda, work)
				ipiv := make([]int, n)
				impl.Dgetrf(n, n, aDense, lda, ipiv)
				want := impl.Dgecon(lapack.MaxColumnSum, n, aDense, lda, denseNorm, work, iwork)
				if math.Abs(got-want) > 1e-14 {
					t.Errorf("Cond mismatch. Want %v, got %v.", want, got)
				}
			}
		}
	}
}
