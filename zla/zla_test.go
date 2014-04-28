package zla

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/dane-unltd/lapack/clapack"
	"github.com/dane-unltd/lapack/dla"
	"github.com/gonum/blas"
	"github.com/gonum/blas/cblas"
	"github.com/gonum/blas/dbw"
	"github.com/gonum/blas/zbw"
)

func init() {
	Register(clapack.La{})
	dla.Register(clapack.La{})
	zbw.Register(cblas.Blas{})
	dbw.Register(cblas.Blas{})
}

func fillRandn(a []complex128, mu complex128, sigmaSq float64) {
	fact := math.Sqrt(0.5 * sigmaSq)
	for i := range a {
		a[i] = complex(fact*rand.NormFloat64(), fact*rand.NormFloat64()) + mu
	}
}

func TestQR(t *testing.T) {
	A := zbw.NewGeneral(blas.ColMajor, 3, 2,
		[]complex128{complex(1, 0), complex(2, 0), complex(3, 0),
			complex(4, 0), complex(5, 0), complex(6, 0)})
	B := zbw.NewGeneral(blas.ColMajor, 3, 2,
		[]complex128{complex(1, 0), complex(1, 0), complex(1, 0), complex(2, 0), complex(2, 0), complex(2, 0)})

	tau := zbw.Allocate(2)

	f := QR(A, tau)

	//fmt.Println(B)
	f.Solve(B)
	//fmt.Println(B)
}

func TestLanczos(t *testing.T) {
	A := zbw.NewGeneral(blas.ColMajor, 3, 4, nil)
	fillRandn(A.Data, 0, 1)

	Acpy := zbw.NewGeneral(blas.ColMajor, 3, 4, nil)
	copy(Acpy.Data, A.Data)

	u0 := make([]complex128, 3)
	fillRandn(u0, 0, 1)

	Ul, Vl, a, b := LanczosBi(Acpy, u0, 3)

	fmt.Println(a, b)

	tmpc := zbw.NewGeneral(blas.ColMajor, 3, 3, nil)
	bidic := zbw.NewGeneral(blas.ColMajor, 3, 3, nil)

	zbw.Gemm(blas.NoTrans, blas.NoTrans, 1, A, Vl, 0, tmpc)
	zbw.Gemm(blas.ConjTrans, blas.NoTrans, 1, Ul, tmpc, 0, bidic)

	fmt.Println(bidic)

	Ur, s, Vr := dla.SVDbd(blas.Lower, a, b)

	tmp := dbw.NewGeneral(blas.ColMajor, 3, 3, nil)
	bidi := dbw.NewGeneral(blas.ColMajor, 3, 3, nil)

	copy(tmp.Data, Ur.Data)
	for i := 0; i < 3; i++ {
		dbw.Scal(s[i], tmp.Col(i))
	}

	dbw.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, Vr, 0, bidi)

	fmt.Println(bidi)
	/*

		_ = Ul
		_ = Vl
			Uc := zbw.NewGeneral(blas.ColMajor, 3, 3, nil)
			zbw.Real2Cmplx(Ur.Data[:3*3], Uc.Data)

			fmt.Println(Uc.Data)

			U := zbw.NewGeneral(blas.ColMajor, M, K, nil)
			zbw.Gemm(blas.NoTrans, blas.NoTrans, 1, U1, Uc, 0, U)
	*/
}
