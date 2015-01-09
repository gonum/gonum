package native

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/testblas"
)

func BenchmarkDgemvSmSmNoTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.SmallMat, testblas.SmallMat, 1, 1)
}

func BenchmarkDgemvSmSmNoTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.SmallMat, testblas.SmallMat, 2, 3)
}

func BenchmarkDgemvSmSmTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.SmallMat, testblas.SmallMat, 1, 1)
}

func BenchmarkDgemvSmSmTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.SmallMat, testblas.SmallMat, 2, 3)
}

func BenchmarkDgemvMedMedNoTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.MediumMat, testblas.MediumMat, 1, 1)
}

func BenchmarkDgemvMedMedNoTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.MediumMat, testblas.MediumMat, 2, 3)
}

func BenchmarkDgemvMedMedTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.MediumMat, testblas.MediumMat, 1, 1)
}

func BenchmarkDgemvMedMedTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.MediumMat, testblas.MediumMat, 2, 3)
}

func BenchmarkDgemvLgLgNoTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.LargeMat, testblas.LargeMat, 1, 1)
}

func BenchmarkDgemvLgLgNoTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.LargeMat, testblas.LargeMat, 2, 3)
}

func BenchmarkDgemvLgLgTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.LargeMat, testblas.LargeMat, 1, 1)
}

func BenchmarkDgemvLgLgTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.LargeMat, testblas.LargeMat, 2, 3)
}

func BenchmarkDgemvLgSmNoTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.LargeMat, testblas.SmallMat, 1, 1)
}

func BenchmarkDgemvLgSmNoTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.LargeMat, testblas.SmallMat, 2, 3)
}

func BenchmarkDgemvLgSmTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.LargeMat, testblas.SmallMat, 1, 1)
}

func BenchmarkDgemvLgSmTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.LargeMat, testblas.SmallMat, 2, 3)
}

func BenchmarkDgemvSmLgNoTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.SmallMat, testblas.LargeMat, 1, 1)
}

func BenchmarkDgemvSmLgNoTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.NoTrans, testblas.SmallMat, testblas.LargeMat, 2, 3)
}

func BenchmarkDgemvSmLgTransInc1(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.SmallMat, testblas.LargeMat, 1, 1)
}

func BenchmarkDgemvSmLgTransIncN(b *testing.B) {
	testblas.DgemvBenchmark(b, impl, blas.Trans, testblas.SmallMat, testblas.LargeMat, 2, 3)
}
