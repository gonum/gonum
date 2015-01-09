package cgo

import (
	"testing"

	"github.com/gonum/blas/testblas"
)

func BenchmarkDgerSmSmInc1(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.SmallMat, testblas.SmallMat, 1, 1)
}

func BenchmarkDgerSmSmIncN(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.SmallMat, testblas.SmallMat, 2, 3)
}

func BenchmarkDgerMedMedInc1(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.MediumMat, testblas.MediumMat, 1, 1)
}

func BenchmarkDgerMedMedIncN(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.MediumMat, testblas.MediumMat, 2, 3)
}

func BenchmarkDgerLgLgInc1(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.LargeMat, testblas.LargeMat, 1, 1)
}

func BenchmarkDgerLgLgIncN(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.LargeMat, testblas.LargeMat, 2, 3)
}

func BenchmarkDgerLgSmInc1(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.LargeMat, testblas.SmallMat, 1, 1)
}

func BenchmarkDgerLgSmIncN(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.LargeMat, testblas.SmallMat, 2, 3)
}

func BenchmarkDgerSmLgInc1(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.SmallMat, testblas.LargeMat, 1, 1)
}

func BenchmarkDgerSmLgIncN(b *testing.B) {
	testblas.DgerBenchmark(b, impl, testblas.SmallMat, testblas.LargeMat, 2, 3)
}
