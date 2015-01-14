package zla

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/blas/cblas128"
	"github.com/gonum/lapack/clapack"
	"github.com/gonum/lapack/dla"
)

func init() {
	Register(clapack.Lapack{})
	dla.Register(clapack.Lapack{})
}

func fillRandn(a []complex128, mu complex128, sigmaSq float64) {
	fact := math.Sqrt(0.5 * sigmaSq)
	for i := range a {
		a[i] = complex(fact*rand.NormFloat64(), fact*rand.NormFloat64()) + mu
	}
}

func TestQR(t *testing.T) {
	A := cblas128.General{
		Rows:   3,
		Cols:   2,
		Stride: 2,
		Data:   []complex128{complex(1, 0), complex(2, 0), complex(3, 0), complex(4, 0), complex(5, 0), complex(6, 0)},
	}
	B := cblas128.General{
		Rows:   3,
		Cols:   2,
		Stride: 2,
		Data:   []complex128{complex(1, 0), complex(1, 0), complex(1, 0), complex(2, 0), complex(2, 0), complex(2, 0)},
	}

	tau := make([]complex128, 2)

	f := QR(A, tau)

	//fmt.Println(B)
	f.Solve(B)
	//fmt.Println(B)
}

func f64col(i int, a blas64.General) blas64.Vector {
	return blas64.Vector{
		Inc:  a.Stride,
		Data: a.Data[i:],
	}
}

func TestLanczos(t *testing.T) {
	A := cblas128.General{Rows: 3, Cols: 4, Stride: 4, Data: make([]complex128, 3*4)}
	fillRandn(A.Data, 0, 1)

	Acpy := cblas128.General{Rows: 3, Cols: 4, Stride: 4, Data: make([]complex128, 3*4)}
	copy(Acpy.Data, A.Data)

	u0 := make([]complex128, 3)
	fillRandn(u0, 0, 1)

	Ul, Vl, a, b := LanczosBi(Acpy, u0, 3)

	fmt.Println(a, b)

	tmpc := cblas128.General{Rows: 3, Cols: 3, Stride: 3, Data: make([]complex128, 3*3)}
	bidic := cblas128.General{Rows: 3, Cols: 3, Stride: 3, Data: make([]complex128, 3*3)}

	cblas128.Gemm(blas.NoTrans, blas.NoTrans, 1, A, Vl, 0, tmpc)
	cblas128.Gemm(blas.ConjTrans, blas.NoTrans, 1, Ul, tmpc, 0, bidic)

	fmt.Println(bidic)

	Ur, s, Vr := dla.SVDbd(blas.Lower, a, b)

	tmp := blas64.General{Rows: 3, Cols: 3, Stride: 3, Data: make([]float64, 3*3)}
	bidi := blas64.General{Rows: 3, Cols: 3, Stride: 3, Data: make([]float64, 3*3)}

	copy(tmp.Data, Ur.Data)
	for i := 0; i < 3; i++ {
		blas64.Scal(3, s[i], f64col(i, tmp))
	}

	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, Vr, 0, bidi)

	fmt.Println(bidi)
	/*

		_ = Ul
		_ = Vl
			Uc := zbw.NewGeneral( 3, 3, nil)
			zbw.Real2Cmplx(Ur.Data[:3*3], Uc.Data)

			fmt.Println(Uc.Data)

			U := zbw.NewGeneral( M, K, nil)
			zbw.Gemm(blas.NoTrans, blas.NoTrans, 1, U1, Uc, 0, U)
	*/
}
