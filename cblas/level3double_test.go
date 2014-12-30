package cblas

import (
	"testing"

	"github.com/gonum/blas/testblas"
)

func TestDgemm(t *testing.T) {
	testblas.TestDgemm(t, blasser)
}

func TestDsymm(t *testing.T) {
	testblas.DsymmTest(t, blasser)
}

func TestDtrsm(t *testing.T) {
	testblas.DtrsmTest(t, blasser)
}

func TestDsyrk(t *testing.T) {
	testblas.DsyrkTest(t, blasser)
}

func TestDsyr2k(t *testing.T) {
	testblas.Dsyr2kTest(t, blasser)
}

func TestDtrmm(t *testing.T) {
	testblas.DtrmmTest(t, blasser)
}
