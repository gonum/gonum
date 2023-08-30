package testlapack

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
)

type Dtgsyler interface {
	Dtgsyl(trans blas.Transpose, ijob, m, n int, a []float64, lda int, b []float64, ldb int, c []float64, ldc int, d []float64, ldd int, e []float64, lde int, f []float64, ldf int, work []float64, iwork []int, workspaceQuery bool) (difOut, scaleOut float64, infoOut int)
}

func DtgsylTest(t *testing.T, impl Dtgsyler) {
	const ldAdd = 5
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{4, 9, 20} {
		for _, m := range []int{4, 9, 20} {
			for _, lda := range []int{m, m + ldAdd} {
				for _, ldb := range []int{n, n + ldAdd} {
					for _, ldc := range []int{n, n + ldAdd} {
						for _, ldd := range []int{m, m + ldAdd} {
							for _, lde := range []int{n, n + ldAdd} {
								for _, ldf := range []int{n, n + ldAdd} {
									for _, ijob := range []int{2, 1, 0} {
										testSolveDtgsyl(t, impl, rnd, blas.NoTrans, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf)
										return
										testSolveDtgsyl(t, impl, rnd, blas.Trans, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func testSolveDtgsyl(t *testing.T, impl Dtgsyler, rnd *rand.Rand, trans blas.Transpose, ijob, m, n, lda, ldb, ldc, ldd, lde, ldf int) {
	const tol = 1e-12
	name := fmt.Sprintf("trans=%v,ijob=%v,n=%v,m=%v,lda=%v,ldb=%v,ldc=%v,ldd=%v,lde=%v,ldf=%v", string(trans), ijob, n, m, lda, ldb, ldc, ldd, lde, ldf)
	lda = max(lda, max(1, m))
	ldb = max(ldb, max(1, n))
	ldc = max(ldc, max(1, n))
	ldd = max(ldd, max(1, m))
	lde = max(lde, max(1, n))
	ldf = max(ldf, max(1, n))
	notrans := trans == blas.NoTrans
	// Generate random matrices (A, D) and (B, E) which must be
	// in generalized Schur canonical form, i.e. A, B are upper
	// quasi triangular and D, E are upper triangular.
	var a, b, c, d, e, f blas64.General
	a, _, _ = randomSchurCanonical(m, lda, false, rnd)
	b, _, _ = randomSchurCanonical(n, ldb, false, rnd)

	d = randomUpperTriangular(m, ldd, rnd)
	e = randomUpperTriangular(n, lde, rnd)

	// Generate random general matrix.
	c = randomGeneral(m, n, ldc, rnd)
	f = randomGeneral(m, n, ldf, rnd)
	// Query for optimum workspace size.
	var query [1]float64
	impl.Dtgsyl(trans, ijob, m, n, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride, d.Data, d.Stride, e.Data, e.Stride, f.Data, f.Stride, query[:], nil, true)
	lwork := int(query[0] + dlamchE)
	if lwork < 1 {
		t.Fatalf("%v: bad workspace query lwork=%d", name, lwork)
	}
	if notrans && (ijob == 1 || ijob == 2) {
		if lwork < 2*m*n {
			t.Fatalf("%v: bad workspace query lwork=%d, expected >=%d", name, lwork, 2*m*n)
		}
	}
	iwork := make([]int, m+n+6)
	work := make([]float64, lwork)
	dif, scale, info := impl.Dtgsyl(trans, ijob, m, n, a.Data, a.Stride, b.Data, b.Stride, c.Data, c.Stride, d.Data, d.Stride, e.Data, e.Stride, f.Data, f.Stride, work, iwork, false)

	t.Error(name, dif, scale, info)
}
