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

func TestDgbmv(t *testing.T) {
	testblas.DgbmvTest(t, blasser)
}

func TestDtbsv(t *testing.T) {
	testblas.DtbsvTest(t, blasser)
}

func TestDsbmv(t *testing.T) {
	testblas.DsbmvTest(t, blasser)
}

func TestDtbmv(t *testing.T) {
	testblas.DtbmvTest(t, blasser)
}

func TestDtrsv(t *testing.T) {
	testblas.DtrsvTest(t, blasser)
}

func TestDtrmv(t *testing.T) {
	testblas.DtrmvTest(t, blasser)
}

func TestDsymv(t *testing.T) {
	testblas.DsymvTest(t, blasser)
}

func TestDsyr(t *testing.T) {
	testblas.DsyrTest(t, blasser)
}

func TestDsyr2(t *testing.T) {
	testblas.Dsyr2Test(t, blasser)
}

func TestDspr(t *testing.T) {
	testblas.DsprTest(t, blasser)
}

func TestDspmv(t *testing.T) {
	testblas.DspmvTest(t, blasser)
}

func TestDtpsv(t *testing.T) {
	testblas.DtpsvTest(t, blasser)
}

func TestDtpmv(t *testing.T) {
	testblas.DtpmvTest(t, blasser)
}
