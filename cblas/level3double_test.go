package cblas

import (
	"testing"

	"github.com/gonum/blas/testblas"
)

func TestDsymm(t *testing.T) {
	testblas.DsymmTest(t, blasser)
}
