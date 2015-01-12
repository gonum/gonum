package dla

import (
	"fmt"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack/clapack"
	"github.com/gonum/matrix/mat64"
)

type fm struct {
	*mat64.Dense
	margin int
}

func (m fm) Format(fs fmt.State, c rune) {
	if c == 'v' && fs.Flag('#') {
		fmt.Fprintf(fs, "%#v", m.Dense)
		return
	}
	mat64.Format(m.Dense, m.margin, '.', fs, c)
}

func init() {
	Register(clapack.Lapack{})
}

func TestQR(t *testing.T) {
	A := blas64.General{
		Rows:   3,
		Cols:   2,
		Stride: 2,
		Data:   []float64{1, 2, 3, 4, 5, 6},
	}
	B := blas64.General{
		Rows:   3,
		Cols:   2,
		Stride: 2,
		Data:   []float64{1, 1, 1, 2, 2, 2},
	}

	tau := make([]float64, 2)

	C := blas64.General{Rows: 2, Cols: 2, Stride: 2, Data: make([]float64, 2*2)}

	blas64.Gemm(blas.Trans, blas.NoTrans, 1, A, B, 0, C)

	fmt.Println(C)

	f := QR(A, tau)

	fmt.Println(B)
	fmt.Println(f)

	f.Solve(B)
	var pm mat64.Dense
	pm.SetRawMatrix(B)
	fmt.Println(fm{&pm, 0})
}
