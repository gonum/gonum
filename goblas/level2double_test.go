package goblas

import (
	"github.com/gonum/blas/testblas"
	"testing"
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

func TestDgbmv(t *testing.T) {
	testblas.DgbmvTest(t, blasser)
}
