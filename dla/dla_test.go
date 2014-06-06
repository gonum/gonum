package dla

import (
	"fmt"
	"github.com/dane-unltd/lapack/clapack"
	"github.com/gonum/blas"
	"github.com/gonum/blas/cblas"
	"github.com/gonum/blas/dbw"
	"github.com/gonum/matrix/mat64"
	"testing"
)

type fm struct {
	mat64.Matrix
	margin int
}

func (m fm) Format(fs fmt.State, c rune) {
	if c == 'v' && fs.Flag('#') {
		fmt.Fprintf(fs, "%#v", m.Matrix)
		return
	}
	mat64.Format(m.Matrix, m.margin, '.', fs, c)
}

func init() {
	Register(clapack.La{})
	dbw.Register(cblas.Blas{})
}

func TestQR(t *testing.T) {
	A := dbw.NewGeneral(3, 2, []float64{1, 2, 3, 4, 5, 6})
	B := dbw.NewGeneral(3, 2, []float64{1, 1, 1, 2, 2, 2})

	tau := dbw.Allocate(2)

	f := QR(A, tau)

	fmt.Println(fm{B, 0})
	f.Solve(B)
	fmt.Println(fm{B, 0})
}
