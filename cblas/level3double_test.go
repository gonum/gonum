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

func TestDsyrk(t *testing.T) {
	testblas.DsyrkTest(t, blasser)
}
