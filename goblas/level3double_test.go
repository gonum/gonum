package goblas

import (
	"testing"

	"github.com/gonum/blas/testblas"
)

func TestDgemm(t *testing.T) {
	testblas.TestDgemm(t, blasser)
}
