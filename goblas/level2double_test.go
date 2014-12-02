package goblas

import (
	"testing"

	"github.com/gonum/blas/testblas"
)

func TestDgemv(t *testing.T) {
	testblas.DgemvTest(t, blasser)
}

func TestDger(t *testing.T) {
	testblas.DgerTest(t, blasser)
}

func TestDtxmv(t *testing.T) {
	testblas.DtxmvTest(t, blasser)
}

func TestDtbmv(t *testing.T) {
	testblas.DtbmvTest(t, blasser)
}
