package cblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/testblas"
)

func BenchmarkDgemmSmSmSm(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmSmall,
		testblas.DgemmSmall,
		testblas.DgemmSmall,
		blas.NoTrans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmMedMedMed(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		blas.NoTrans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmLgLgLg(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmLarge,
		testblas.DgemmLarge,
		testblas.DgemmLarge,
		blas.NoTrans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmLgSmLg(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmLarge,
		testblas.DgemmSmall,
		testblas.DgemmLarge,
		blas.NoTrans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmLgLgSm(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmLarge,
		testblas.DgemmSmall,
		testblas.DgemmLarge,
		blas.NoTrans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmHgHgSm(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmHuge,
		testblas.DgemmSmall,
		testblas.DgemmHuge,
		blas.NoTrans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmMedMedMedTNT(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		blas.Trans,
		blas.NoTrans,
	)
}

func BenchmarkDgemmMedMedMedNTT(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		blas.NoTrans,
		blas.Trans,
	)
}

func BenchmarkDgemmMedMedMedNTNT(b *testing.B) {
	testblas.DgemmBenchmark(b,
		Blas{},
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		testblas.DgemmMedium,
		blas.NoTrans,
		blas.NoTrans,
	)
}
