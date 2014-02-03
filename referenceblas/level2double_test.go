package referenceblas

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

func TestDtbmv(t *testing.T) {
	testblas.DtbmvTest(t, blasser)
}

func TestDgbmv(t *testing.T) {
	testblas.DgbmvTest(t, blasser)
}
