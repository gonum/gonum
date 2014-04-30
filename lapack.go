package lapack

import (
	"github.com/gonum/blas"
)

const None = 'N'

type Job byte

const (
	All       (Job) = 'A'
	Slim      (Job) = 'S'
	Overwrite (Job) = 'O'
)

type CompSV byte

const (
	Compact  (CompSV) = 'P'
	Explicit (CompSV) = 'I'
)

type Float64 interface {
	Dgeqrf(o blas.Order, m, n int, a []float64, lda int, tau []float64) int
	Dormqr(o blas.Order, s blas.Side, t blas.Transpose, m, n, k int, a []float64, lda int, tau []float64, c []float64, ldc int) int
	Dgesdd(o blas.Order, job Job, m, n int, a []float64, lda int, s []float64, u []float64, ldu int, vt []float64, ldvt int) int
	Dgebrd(o blas.Order, m, n int, a []float64, lda int, d, e, tauq, taup []float64) int
	Dbdsdc(o blas.Order, uplo blas.Uplo, compq CompSV, n int, d, e []float64, u []float64, ldu int, vt []float64, ldvt int, q []float64, iq []int32) int
}

type Complex128 interface {
	Float64

	Zgeqrf(o blas.Order, m, n int, a []complex128, lda int, tau []complex128) int
	Zunmqr(o blas.Order, s blas.Side, t blas.Transpose, m, n, k int, a []complex128, lda int, tau []complex128, c []complex128, ldc int) int
	Zgesdd(o blas.Order, job Job, m, n int, a []complex128, lda int, s []float64, u []complex128, ldu int, vt []complex128, ldvt int) int
	Zgebrd(o blas.Order, m, n int, a []complex128, lda int, d, e []float64, tauq, taup []complex128) int
}
