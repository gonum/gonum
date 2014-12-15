package cblas

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

func TestDtbmv(t *testing.T) {
	testblas.DtbmvTest(t, blasser)
}

func TestDtxmv(t *testing.T) {
	testblas.DtxmvTest(t, blasser)
}

func TestDgbmv(t *testing.T) {
	testblas.DgbmvTest(t, blasser)
}

func TestDtbsv(t *testing.T) {
	testblas.DtbsvTest(t, blasser)
}

func TestDsbmv(t *testing.T) {
	testblas.DsbmvTest(t, blasser)
}

func TestDsyr(t *testing.T) {
	testblas.DsyrTest(t, blasser)
}

func TestDsymv(t *testing.T) {
	testblas.DsymvTest(t, blasser)
}

func TestDtrmv(t *testing.T) {
	testblas.DtrmvTest(t, blasser)
}

func TestDsyr2(t *testing.T) {
	testblas.Dsyr2Test(t, blasser)
}
