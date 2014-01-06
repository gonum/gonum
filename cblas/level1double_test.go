package cblas

import (
	"github.com/gonum/blas/testblas"
	"testing"
)

var blasser = Blas{}

func TestDasum(t *testing.T) {
	testblas.DasumTest(t, blasser)
}

func TestDaxpy(t *testing.T) {
	testblas.DaxpyTest(t, blasser)
}

func TestDdot(t *testing.T) {
	testblas.DdotTest(t, blasser)
}

func TestDnrm2(t *testing.T) {
	testblas.Dnrm2Test(t, blasser)
}

func TestIdamax(t *testing.T) {
	testblas.IdamaxTest(t, blasser)
}

func TestDswap(t *testing.T) {
	testblas.DswapTest(t, blasser)
}

func TestDcopy(t *testing.T) {
	testblas.DcopyTest(t, blasser)
}

func TestDrotg(t *testing.T) {
	testblas.DrotgTest(t, blasser)
}

func TestDrotmg(t *testing.T) {
	testblas.DrotmgTest(t, blasser)
}

func TestDrot(t *testing.T) {
	testblas.DrotTest(t, blasser)
}

func TestDrotm(t *testing.T) {
	testblas.DrotmTest(t, blasser)
}

func TestDscal(t *testing.T) {
	testblas.DscalTest(t, blasser)
}
