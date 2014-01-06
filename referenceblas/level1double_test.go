package naivegoblas

import (
	"blas/blastest"
	"testing"
)

var blasser = Blas{}

func TestDasum(t *testing.T) {
	blastest.DasumTest(t, blasser)
}

func TestDaxpy(t *testing.T) {
	blastest.DaxpyTest(t, blasser)
}

func TestDdot(t *testing.T) {
	blastest.DdotTest(t, blasser)
}

func TestDnrm2(t *testing.T) {
	blastest.Dnrm2Test(t, blasser)
}

func TestIdamax(t *testing.T) {
	blastest.IdamaxTest(t, blasser)
}

func TestDswap(t *testing.T) {
	blastest.DswapTest(t, blasser)
}

func TestDcopy(t *testing.T) {
	blastest.DcopyTest(t, blasser)
}

func TestDrotg(t *testing.T) {
	blastest.DrotgTest(t, blasser)
}

func TestDrotmg(t *testing.T) {
	blastest.DrotmgTest(t, blasser)
}

func TestDrot(t *testing.T) {
	blastest.DrotTest(t, blasser)
}

func TestDrotm(t *testing.T) {
	blastest.DrotmTest(t, blasser)
}

func TestDscal(t *testing.T) {
	blastest.DscalTest(t, blasser)
}
