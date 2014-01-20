package referenceblas

import (
	"github.com/gonum/blas/testblas"
	"testing"
)

func TestGemv(t *testing.T) {
	testblas.DgemvTest(t, blasser)
}
